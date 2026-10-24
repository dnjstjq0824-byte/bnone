package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

type Worker struct {
	db     *sql.DB
	driver string
}

func main() {
	dbURL := getEnv("DATABASE_URL", "bn_one.db")
	driver := getEnv("DB_DRIVER", "sqlite")
	if strings.HasPrefix(dbURL, "postgres") {
		driver = "pgx"
	}

	db, err := sql.Open(driver, dbURL)
	if err != nil {
		log.Fatalf("[Worker] Failed to connect to DB: %v", err)
	}
	defer db.Close()

	w := &Worker{db: db, driver: driver}
	log.Println("[BN ONE Worker] Background worker initialized. Starting scheduled jobs...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run periodic tasks
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Initial run
	w.runRenewalAlerts()
	w.runBudgetAlerts()
	w.processDataExports()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				w.runRenewalAlerts()
				w.runBudgetAlerts()
				w.processDataExports()
			}
		}
	}()

	<-stop
	log.Println("[BN ONE Worker] Shutting down cleanly...")
}

func (w *Worker) runRenewalAlerts() {
	now := time.Now().UTC()
	thresholdDate := now.AddDate(0, 0, 3).Format("2006-01-02")

	query := `SELECT us.id, us.user_id, COALESCE(ss.name, us.custom_service_name, 'Subscription'), us.amount_minor, us.currency_code, us.next_renewal_at
	          FROM user_subscriptions us
	          LEFT JOIN subscription_services ss ON us.service_id = ss.id
	          WHERE us.status = 'ACTIVE' AND us.reminder_enabled = TRUE 
	            AND us.next_renewal_at <= $1`
	rows, err := w.db.Query(query, thresholdDate)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var subID, userID, name, curr, nextRenew string
		var amount int64
		if err := rows.Scan(&subID, &userID, &name, &amount, &curr, &nextRenew); err == nil {
			dedupKey := fmt.Sprintf("renew_%s_%s", subID, nextRenew)

			// Check if notification already delivered for this renewal date
			var exists bool
			w.db.QueryRow("SELECT 1 FROM subscription_reminders WHERE deduplication_key = $1", dedupKey).Scan(&exists)
			if !exists {
				// Create notification
				notifID := uuid.New().String()
				title := fmt.Sprintf("[%s] 구독 갱신 예정 안내", name)
				body := fmt.Sprintf("%s 정기결제 갱신일(%s)이 다가옵니다. 예정금액: %s %d", name, nextRenew, curr, amount)
				w.db.Exec(`INSERT INTO notifications (id, user_id, type, title, body, channel, status, scheduled_at, sent_at, created_at)
				           VALUES ($1, $2, 'RENEWAL_REMINDER', $3, $4, 'IN_APP', 'DELIVERED', $5, $5, $5)`,
					notifID, userID, title, body, now)

				// Record reminder
				remindID := uuid.New().String()
				w.db.Exec(`INSERT INTO subscription_reminders (id, subscription_id, scheduled_at, status, sent_at, deduplication_key, created_at)
				           VALUES ($1, $2, $3, 'SENT', $3, $4, $3)`,
					remindID, subID, now, dedupKey)

				log.Printf("[Worker] Sent subscription renewal alert to user %s for %s", userID, name)
			}
		}
	}
}

func (w *Worker) runBudgetAlerts() {
	now := time.Now().UTC()
	currentDate := now.Format("2006-01-02")

	query := `SELECT b.id, b.user_id, b.name, b.limit_minor, b.currency_code
	          FROM budgets b
	          WHERE b.status = 'ACTIVE' AND b.start_date <= $1 AND b.end_date >= $1`
	rows, err := w.db.Query(query, currentDate)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var bID, userID, bName, curr string
		var limitMinor int64
		if err := rows.Scan(&bID, &userID, &bName, &limitMinor, &curr); err == nil && limitMinor > 0 {
			var spent sql.NullInt64
			w.db.QueryRow(`SELECT -SUM(ae.amount_minor_signed)
			               FROM transactions t
			               JOIN account_entries ae ON t.id = ae.transaction_id
			               WHERE t.user_id = $1 AND t.transaction_type = 'EXPENSE' AND t.status = 'POSTED' AND t.deleted_at IS NULL`, userID).Scan(&spent)

			if spent.Valid {
				usagePercent := float64(spent.Int64) / float64(limitMinor) * 100.0
				if usagePercent >= 80.0 {
					dedupKey := fmt.Sprintf("budget_%s_%s", bID, now.Format("2006-01"))
					var exists bool
					w.db.QueryRow("SELECT 1 FROM notifications WHERE user_id = $1 AND type = 'BUDGET_ALERT' AND body LIKE '%' || $2 || '%'", userID, bName).Scan(&exists)
					if !exists {
						notifID := uuid.New().String()
						title := fmt.Sprintf("[%s] 예산 사용량 80%% 초과 경고", bName)
						body := fmt.Sprintf("현재 예산 사용률이 %.1f%%에 도달했습니다. (사용: %s %d / 한도: %s %d)", usagePercent, curr, spent.Int64, curr, limitMinor)
						w.db.Exec(`INSERT INTO notifications (id, user_id, type, title, body, channel, status, scheduled_at, sent_at, created_at)
						           VALUES ($1, $2, 'BUDGET_ALERT', $3, $4, 'IN_APP', 'DELIVERED', $5, $5, $5)`,
							notifID, userID, title, body, now)
						log.Printf("[Worker] Budget alert triggered for user %s on budget %s", userID, dedupKey)
					}
				}
			}
		}
	}
}

func (w *Worker) processDataExports() {
	query := `SELECT id, user_id, format FROM data_export_jobs WHERE status = 'PENDING' LIMIT 5`
	rows, err := w.db.Query(query)
	if err != nil {
		return
	}
	defer rows.Close()

	type ExpJob struct {
		ID     string
		UserID string
		Format string
	}
	jobs := make([]ExpJob, 0)
	for rows.Next() {
		var j ExpJob
		if err := rows.Scan(&j.ID, &j.UserID, &j.Format); err == nil {
			jobs = append(jobs, j)
		}
	}

	for _, job := range jobs {
		now := time.Now().UTC()
		exportDir := filepath.Join(os.TempDir(), "bn_one_exports")
		os.MkdirAll(exportDir, 0755)
		filePath := filepath.Join(exportDir, fmt.Sprintf("export_%s.csv", job.ID))

		f, err := os.Create(filePath)
		if err == nil {
			cw := csv.NewWriter(f)
			cw.Write([]string{"Transaction ID", "Date", "Type", "Description", "Merchant", "Amount Minor Signed", "Currency", "Status"})

			txRows, _ := w.db.Query(`SELECT t.id, t.occurred_at, t.transaction_type, t.description, t.merchant_name, ae.amount_minor_signed, ae.currency_code, t.status
			                         FROM transactions t
			                         LEFT JOIN account_entries ae ON t.id = ae.transaction_id
			                         WHERE t.user_id = $1 AND t.deleted_at IS NULL ORDER BY t.occurred_at DESC`, job.UserID)
			if txRows != nil {
				for txRows.Next() {
					var id, occurred, txType, desc, merch, curr, status string
					var amount int64
					if err := txRows.Scan(&id, &occurred, &txType, &desc, &merch, &amount, &curr, &status); err == nil {
						cw.Write([]string{id, occurred, txType, desc, merch, fmt.Sprintf("%d", amount), curr, status})
					}
				}
				txRows.Close()
			}
			cw.Flush()
			f.Close()

			w.db.Exec("UPDATE data_export_jobs SET status = 'COMPLETED', storage_reference = $1, completed_at = $2 WHERE id = $3", filePath, now, job.ID)
			log.Printf("[Worker] Data export completed for job %s", job.ID)
		}
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
