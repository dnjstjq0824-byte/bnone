package notifications

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bookingnetwork/bn-one/services/api/internal/database"
	"github.com/bookingnetwork/bn-one/services/api/internal/middleware"
	"github.com/bookingnetwork/bn-one/services/api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

// ---------------- Notifications ----------------

func (h *Handler) ListNotifications(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	query := `SELECT id, user_id, type, title, body, channel, status, scheduled_at, sent_at, read_at, created_at
	          FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT 50`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	list := make([]models.Notification, 0)
	for rows.Next() {
		var n models.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.Channel, &n.Status, &n.ScheduledAt, &n.SentAt, &n.ReadAt, &n.CreatedAt); err == nil {
			list = append(list, n)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, list)
}

func (h *Handler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	notifID := chi.URLParam(r, "id")

	now := time.Now().UTC()
	_, err := h.db.Exec("UPDATE notifications SET status = 'READ', read_at = $1 WHERE id = $2 AND user_id = $3", now, notifID, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	middleware.RespondJSON(w, r, http.StatusOK, map[string]string{"message": "Marked as read"})
}

func (h *Handler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var body struct {
		Platform           string `json:"platform"`
		DeviceIdentifier   string `json:"device_identifier"`
		PushTokenReference string `json:"push_token_reference"`
		AppVersion         string `json:"app_version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	devID := uuid.New().String()
	now := time.Now().UTC()
	query := `INSERT INTO devices (id, user_id, platform, device_identifier_hash, push_token_reference, app_version, last_seen_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	h.db.Exec(query, devID, userID, body.Platform, body.DeviceIdentifier, body.PushTokenReference, body.AppVersion, now)

	middleware.RespondJSON(w, r, http.StatusCreated, map[string]string{"device_id": devID, "status": "REGISTERED"})
}

// ---------------- Support Tickets ----------------

func (h *Handler) ListSupportTickets(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	query := `SELECT id, user_id, subject, category, description, priority, status, assigned_admin_id, created_at, updated_at, resolved_at
	          FROM support_tickets WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}
	defer rows.Close()

	tickets := make([]models.SupportTicket, 0)
	for rows.Next() {
		var t models.SupportTicket
		if err := rows.Scan(&t.ID, &t.UserID, &t.Subject, &t.Category, &t.Description, &t.Priority, &t.Status, &t.AssignedAdminID, &t.CreatedAt, &t.UpdatedAt, &t.ResolvedAt); err == nil {
			tickets = append(tickets, t)
		}
	}

	middleware.RespondJSON(w, r, http.StatusOK, tickets)
}

func (h *Handler) CreateSupportTicket(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	var body struct {
		Subject     string `json:"subject"`
		Category    string `json:"category"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Malformed JSON body", nil)
		return
	}

	if body.Subject == "" || body.Description == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "VALIDATION_ERROR", "Subject and description required", nil)
		return
	}
	if body.Category == "" {
		body.Category = "GENERAL"
	}
	if body.Priority == "" {
		body.Priority = "NORMAL"
	}

	ticketID := uuid.New().String()
	now := time.Now().UTC()

	query := `INSERT INTO support_tickets (id, user_id, subject, category, description, priority, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, 'OPEN', $7, $8)`
	_, err := h.db.Exec(query, ticketID, userID, body.Subject, body.Category, body.Description, body.Priority, now, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	// Create initial message
	msgID := uuid.New().String()
	h.db.Exec("INSERT INTO support_messages (id, ticket_id, sender_type, sender_id, message, visibility, created_at) VALUES ($1, $2, 'USER', $3, $4, 'PUBLIC', $5)",
		msgID, ticketID, userID, body.Description, now)

	middleware.RespondJSON(w, r, http.StatusCreated, map[string]interface{}{
		"ticket_id": ticketID,
		"status":    "OPEN",
		"subject":   body.Subject,
		"message":   "Your inquiry has been submitted. Our support team will respond shortly.",
	})
}

func (h *Handler) GetSupportTicket(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	ticketID := chi.URLParam(r, "id")

	var t models.SupportTicket
	query := `SELECT id, user_id, subject, category, description, priority, status, assigned_admin_id, created_at, updated_at, resolved_at
	          FROM support_tickets WHERE id = $1 AND user_id = $2`
	err := h.db.QueryRow(query, ticketID, userID).Scan(&t.ID, &t.UserID, &t.Subject, &t.Category, &t.Description, &t.Priority, &t.Status, &t.AssignedAdminID, &t.CreatedAt, &t.UpdatedAt, &t.ResolvedAt)
	if err != nil {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Ticket not found", nil)
		return
	}

	// Messages
	msgRows, err := h.db.Query("SELECT id, ticket_id, sender_type, sender_id, message, visibility, created_at FROM support_messages WHERE ticket_id = $1 AND visibility = 'PUBLIC' ORDER BY created_at ASC", ticketID)
	if err == nil {
		msgs := make([]models.SupportMessage, 0)
		for msgRows.Next() {
			var m models.SupportMessage
			if err := msgRows.Scan(&m.ID, &m.TicketID, &m.SenderType, &m.SenderID, &m.Message, &m.Visibility, &m.CreatedAt); err == nil {
				msgs = append(msgs, m)
			}
		}
		msgRows.Close()
		t.Messages = msgs
	}

	middleware.RespondJSON(w, r, http.StatusOK, t)
}

func (h *Handler) PostTicketMessage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	ticketID := chi.URLParam(r, "id")

	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Message == "" {
		middleware.RespondError(w, r, http.StatusBadRequest, "INVALID_REQUEST", "Message required", nil)
		return
	}

	var exists bool
	err := h.db.QueryRow("SELECT 1 FROM support_tickets WHERE id = $1 AND user_id = $2", ticketID, userID).Scan(&exists)
	if err != nil || !exists {
		middleware.RespondError(w, r, http.StatusNotFound, "NOT_FOUND", "Ticket not found", nil)
		return
	}

	msgID := uuid.New().String()
	now := time.Now().UTC()
	_, err = h.db.Exec("INSERT INTO support_messages (id, ticket_id, sender_type, sender_id, message, visibility, created_at) VALUES ($1, $2, 'USER', $3, $4, 'PUBLIC', $5)",
		msgID, ticketID, userID, body.Message, now)
	if err != nil {
		middleware.RespondError(w, r, http.StatusInternalServerError, "DATABASE_ERROR", err.Error(), nil)
		return
	}

	// Update ticket updated_at
	h.db.Exec("UPDATE support_tickets SET updated_at = $1, status = 'OPEN' WHERE id = $2", now, ticketID)

	middleware.RespondJSON(w, r, http.StatusCreated, map[string]string{"message_id": msgID, "status": "SENT"})
}
