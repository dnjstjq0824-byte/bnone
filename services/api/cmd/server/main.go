package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bookingnetwork/bn-one/services/api/internal/accounts"
	"github.com/bookingnetwork/bn-one/services/api/internal/admin"
	"github.com/bookingnetwork/bn-one/services/api/internal/analytics"
	"github.com/bookingnetwork/bn-one/services/api/internal/auth"
	"github.com/bookingnetwork/bn-one/services/api/internal/budgets"
	"github.com/bookingnetwork/bn-one/services/api/internal/config"
	"github.com/bookingnetwork/bn-one/services/api/internal/database"
	"github.com/bookingnetwork/bn-one/services/api/internal/forex"
	"github.com/bookingnetwork/bn-one/services/api/internal/middleware"
	"github.com/bookingnetwork/bn-one/services/api/internal/notifications"
	"github.com/bookingnetwork/bn-one/services/api/internal/subscriptions"
	"github.com/bookingnetwork/bn-one/services/api/internal/transactions"
	"github.com/bookingnetwork/bn-one/services/api/internal/trips"
	"github.com/bookingnetwork/bn-one/services/api/internal/users"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	log.Printf("[BN ONE] Starting API Server on port %s (env: %s)...", cfg.Port, cfg.Environment)

	// Connect to Database
	db, err := database.Connect(cfg.DBDriver, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[DB] Fatal connection error: %v", err)
	}
	defer db.Close()

	// Automatically run migrations and fixtures
	migrationsPath := "./database/migrations"
	fixturesPath := "./database/fixtures"
	if custom := os.Getenv("MIGRATIONS_PATH"); custom != "" {
		migrationsPath = custom
	} else if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		if _, err := os.Stat("./migrations"); err == nil {
			migrationsPath = "./migrations"
			fixturesPath = "./fixtures"
		} else {
			migrationsPath = "../../database/migrations"
			fixturesPath = "../../database/fixtures"
		}
	}
	if customFixtures := os.Getenv("FIXTURES_PATH"); customFixtures != "" {
		fixturesPath = customFixtures
	}
	if err := db.Migrate(migrationsPath, fixturesPath); err != nil {
		log.Printf("[DB] Migration warning: %v", err)
	}

	// Handlers
	authHandler := auth.NewHandler(db, cfg.JWTSecret)
	usersHandler := users.NewHandler(db)
	accountsHandler := accounts.NewHandler(db)
	txHandler := transactions.NewHandler(db)
	analyticsHandler := analytics.NewHandler(db)
	budgetsHandler := budgets.NewHandler(db)
	subsHandler := subscriptions.NewHandler(db)
	fxHandler := forex.NewHandler(db)
	tripsHandler := trips.NewHandler(db)
	notifHandler := notifications.NewHandler(db)
	adminHandler := admin.NewHandler(db)

	r := chi.NewRouter()

	// Base middlewares
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(60 * time.Second))
	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.SecurityHeadersMiddleware)

	// CORS configuration for Vue Web and Admin Console
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "http://localhost:5174", "http://127.0.0.1:3000", "http://127.0.0.1:5173", "http://127.0.0.1:5174", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID", "Idempotency-Key"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health and Readiness
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"HEALTHY","product":"BN ONE API","version":"1.0.0"}`))
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, `{"status":"NOT_READY"}`, http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"READY"}`))
	})

	// API v1 Routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public Auth
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Post("/auth/admin/login", authHandler.AdminLogin)

		// Public FX & Currencies Reference
		r.Get("/currencies", fxHandler.ListCurrencies)
		r.Get("/fx/rates", fxHandler.GetRates)
		r.Get("/fx/rates/{base}/{quote}", fxHandler.GetRatePair)
		r.Post("/fx/convert", fxHandler.ConvertCurrency)

		// Public Subscription Catalog List
		r.Get("/subscription-services", subsHandler.ListSubscriptionServices)

		// Authenticated User Routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

			r.Post("/auth/logout", authHandler.Logout)

			// User profile & preferences
			r.Get("/me", usersHandler.GetMe)
			r.Patch("/me", usersHandler.UpdateMe)
			r.Get("/me/preferences", usersHandler.GetPreferences)
			r.Patch("/me/preferences", usersHandler.UpdatePreferences)
			r.Get("/me/sessions", usersHandler.GetSessions)
			r.Delete("/me/sessions/{id}", usersHandler.RevokeSession)
			r.Post("/me/data-exports", usersHandler.CreateDataExport)
			r.Get("/me/data-exports/{id}", usersHandler.GetDataExport)
			r.Post("/me/deletion-requests", usersHandler.CreateDeletionRequest)

			// Accounts & Assets
			r.Get("/accounts", accountsHandler.ListAccounts)
			r.Post("/accounts", accountsHandler.CreateAccount)
			r.Get("/accounts/{id}", accountsHandler.GetAccount)
			r.Patch("/accounts/{id}", accountsHandler.UpdateAccount)
			r.Post("/accounts/{id}/archive", accountsHandler.ArchiveAccount)
			r.Get("/accounts/{id}/balance", accountsHandler.GetBalance)
			r.Get("/accounts/{id}/entries", accountsHandler.GetEntries)

			r.Get("/assets", accountsHandler.ListAssets)
			r.Post("/assets", accountsHandler.CreateAsset)
			r.Get("/liabilities", accountsHandler.ListLiabilities)
			r.Post("/liabilities", accountsHandler.CreateLiability)
			r.Get("/net-worth", accountsHandler.GetNetWorth)
			r.Get("/net-worth/history", accountsHandler.NetWorthHistory)

			// Transactions & Categories
			r.Get("/transactions", txHandler.ListTransactions)
			r.Post("/transactions", txHandler.CreateTransaction)
			r.Get("/transactions/{id}", txHandler.GetTransaction)
			r.Delete("/transactions/{id}", txHandler.DeleteTransaction)
			r.Post("/transfers/manual", txHandler.ManualTransfer)
			r.Get("/categories", txHandler.ListCategories)
			r.Post("/categories", txHandler.CreateCategory)

			// Spending Analytics
			r.Get("/analytics/overview", analyticsHandler.Overview)
			r.Get("/analytics/income-expense", analyticsHandler.IncomeExpense)
			r.Get("/analytics/category-spending", analyticsHandler.CategorySpending)
			r.Get("/analytics/monthly-trend", analyticsHandler.MonthlyTrend)
			r.Get("/analytics/net-worth", accountsHandler.GetNetWorth)
			r.Get("/analytics/account-balances", analyticsHandler.AccountBalances)

			// Budgets
			r.Get("/budgets", budgetsHandler.ListBudgets)
			r.Post("/budgets", budgetsHandler.CreateBudget)
			r.Get("/budgets/{id}", budgetsHandler.GetBudget)
			r.Patch("/budgets/{id}", budgetsHandler.UpdateBudget)
			r.Delete("/budgets/{id}", budgetsHandler.DeleteBudget)
			r.Get("/budgets/{id}/progress", budgetsHandler.GetBudgetProgress)

			// Subscriptions
			r.Get("/me/subscriptions", subsHandler.ListUserSubscriptions)
			r.Post("/me/subscriptions", subsHandler.CreateUserSubscription)
			r.Get("/me/subscriptions/{id}", subsHandler.GetSubscription)
			r.Patch("/me/subscriptions/{id}", subsHandler.UpdateSubscription)
			r.Delete("/me/subscriptions/{id}", subsHandler.DeleteSubscription)
			r.Post("/me/subscriptions/{id}/pause", subsHandler.PauseSubscription)
			r.Post("/me/subscriptions/{id}/resume", subsHandler.ResumeSubscription)
			r.Post("/me/subscriptions/{id}/cancel-record", subsHandler.CancelRecord)
			r.Get("/me/subscriptions/summary", subsHandler.Summary)
			r.Get("/me/subscriptions/renewals", subsHandler.Renewals)

			// Travel & Trips
			r.Get("/trips", tripsHandler.ListTrips)
			r.Post("/trips", tripsHandler.CreateTrip)
			r.Get("/trips/{id}", tripsHandler.GetTrip)
			r.Patch("/trips/{id}", tripsHandler.UpdateTrip)
			r.Delete("/trips/{id}", tripsHandler.DeleteTrip)
			r.Get("/trips/{id}/expenses", tripsHandler.GetTripExpenses)
			r.Post("/trips/{id}/expenses/link", tripsHandler.LinkExpense)
			r.Delete("/trips/{id}/expenses/{transaction_id}", tripsHandler.UnlinkExpense)
			r.Get("/trips/{id}/report", tripsHandler.GetTripReport)

			// Notifications
			r.Get("/me/notifications", notifHandler.ListNotifications)
			r.Patch("/me/notifications/{id}/read", notifHandler.MarkAsRead)
			r.Post("/me/devices", notifHandler.RegisterDevice)

			// Support
			r.Get("/me/support-tickets", notifHandler.ListSupportTickets)
			r.Post("/me/support-tickets", notifHandler.CreateSupportTicket)
			r.Get("/me/support-tickets/{id}", notifHandler.GetSupportTicket)
			r.Post("/me/support-tickets/{id}/messages", notifHandler.PostTicketMessage)
		})

		// Operations Admin Routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.AdminAuthMiddleware(cfg.JWTSecret))

			r.Get("/dashboard", adminHandler.Dashboard)
			r.Get("/users", adminHandler.ListUsers)
			r.Get("/users/{id}", adminHandler.GetUserDetail)
			r.Patch("/users/{id}/status", adminHandler.UpdateUserStatus)
			r.Post("/users/{id}/revoke-sessions", adminHandler.RevokeUserSessions)

			r.Get("/support-tickets", adminHandler.ListSupportTickets)
			r.Get("/support-tickets/{id}", adminHandler.GetSupportTicket)
			r.Post("/support-tickets/{id}/assign", adminHandler.AssignSupportTicket)
			r.Post("/support-tickets/{id}/reply", adminHandler.ReplySupportTicket)
			r.Post("/support-tickets/{id}/resolve", adminHandler.ResolveSupportTicket)

			r.Get("/fx/providers", adminHandler.ListFXProviders)
			r.Post("/fx/sync-jobs", adminHandler.TriggerFXSync)

			r.Get("/privacy/deletion-requests", adminHandler.ListDeletionRequests)
			r.Post("/privacy/deletion-requests/{id}/review", adminHandler.ReviewDeletionRequest)

			r.Get("/audit-logs", adminHandler.ListAuditLogs)
			r.Get("/feature-flags", adminHandler.ListFeatureFlags)
			r.Patch("/feature-flags/{key}", adminHandler.ToggleFeatureFlag)
		})
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[BN ONE API] HTTP Server running on http://127.0.0.1:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[BN ONE API] Could not listen: %v\n", err)
		}
	}()

	<-stop
	log.Println("[BN ONE API] Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[BN ONE API] Server forced to shutdown: %v", err)
	}

	log.Println("[BN ONE API] Server exited cleanly.")
}
