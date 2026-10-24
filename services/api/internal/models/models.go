package models

import (
	"time"
)

// Standard API Response envelope
type ApiResponse struct {
	Data interface{} `json:"data,omitempty"`
	Meta *MetaInfo   `json:"meta,omitempty"`
}

type ApiErrorResponse struct {
	Error ApiError `json:"error"`
}

type ApiError struct {
	Code      string        `json:"code"`
	Message   string        `json:"message"`
	Details   []interface{} `json:"details,omitempty"`
	RequestID string        `json:"request_id,omitempty"`
}

type MetaInfo struct {
	RequestID  string    `json:"request_id"`
	ServerTime time.Time `json:"server_time"`
	TotalCount *int      `json:"total_count,omitempty"`
}

// User entity
type User struct {
	ID                   string     `json:"id"`
	Email                string     `json:"email"`
	PasswordHash         string     `json:"-"`
	DisplayName          string     `json:"display_name"`
	CountryCode          string     `json:"country_code"`
	Locale               string     `json:"locale"`
	BaseCurrency         string     `json:"base_currency"`
	Timezone             string     `json:"timezone"`
	Status               string     `json:"status"`
	EmailVerifiedAt      *time.Time `json:"email_verified_at,omitempty"`
	OnboardingCompletedAt *time.Time `json:"onboarding_completed_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at,omitempty"`
}

// User Session
type UserSession struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	RefreshTokenHash string     `json:"-"`
	DeviceID         string     `json:"device_id"`
	DeviceName       string     `json:"device_name"`
	Platform         string     `json:"platform"`
	ExpiresAt        time.Time  `json:"expires_at"`
	RevokedAt        *time.Time `json:"revoked_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// User Preferences
type UserPreferences struct {
	ID                   string `json:"id"`
	UserID               string `json:"user_id"`
	MaskAmounts          bool   `json:"mask_amounts"`
	DarkMode             bool   `json:"dark_mode"`
	StartOfWeek          string `json:"start_of_week"`
	NotificationSettings string `json:"notification_settings"`
	AnalyticsConsent     bool   `json:"analytics_consent"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Finance Account
type FinanceAccount struct {
	ID                     string     `json:"id"`
	UserID                 string     `json:"user_id"`
	Name                   string     `json:"name"`
	AccountType            string     `json:"account_type"` // CASH, MANUAL_BANK, SAVINGS, E_WALLET_MANUAL, OTHER_ASSET, LIABILITY
	CurrencyCode           string     `json:"currency_code"`
	OpeningBalanceMinor    int64      `json:"opening_balance_minor"`
	OpeningBalanceDate     string     `json:"opening_balance_date"`
	CalculatedBalanceMinor int64      `json:"calculated_balance_minor"`
	Status                 string     `json:"status"` // ACTIVE, ARCHIVED, HIDDEN
	Color                  string     `json:"color"`
	Icon                   string     `json:"icon"`
	SortOrder              int        `json:"sort_order"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	ArchivedAt             *time.Time `json:"archived_at,omitempty"`
	SourceType             string     `json:"source_type,omitempty"` // "USER_ENTERED"
}

// Account Entry
type AccountEntry struct {
	ID                string    `json:"id"`
	TransactionID     string    `json:"transaction_id"`
	AccountID         string    `json:"account_id"`
	AmountMinorSigned int64     `json:"amount_minor_signed"`
	CurrencyCode      string    `json:"currency_code"`
	CreatedAt         time.Time `json:"created_at"`
}

// Manual Asset
type Asset struct {
	ID                  string    `json:"id"`
	UserID              string    `json:"user_id"`
	Name                string    `json:"name"`
	AssetType           string    `json:"asset_type"` // REAL_ESTATE, VEHICLE, INVESTMENT_MANUAL, OTHER
	CurrencyCode        string    `json:"currency_code"`
	EstimatedValueMinor int64     `json:"estimated_value_minor"`
	ValuationDate       string    `json:"valuation_date"`
	ValuationSource     string    `json:"valuation_source"`
	Notes               string    `json:"notes"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Manual Liability
type Liability struct {
	ID                     string    `json:"id"`
	UserID                 string    `json:"user_id"`
	Name                   string    `json:"name"`
	LiabilityType          string    `json:"liability_type"` // MORTGAGE, PERSONAL_LOAN, CREDIT_CARD_DEBT, OTHER
	CurrencyCode           string    `json:"currency_code"`
	OutstandingAmountMinor int64     `json:"outstanding_amount_minor"`
	ValuationDate          string    `json:"valuation_date"`
	InterestRate           *float64  `json:"interest_rate_optional,omitempty"`
	Notes                  string    `json:"notes"`
	Status                 string    `json:"status"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// Category
type Category struct {
	ID           string    `json:"id"`
	OwnerUserID  *string   `json:"owner_user_id,omitempty"`
	ParentID     *string   `json:"parent_id,omitempty"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	CategoryType string    `json:"category_type"` // INCOME, EXPENSE, TRANSFER
	Icon         string    `json:"icon"`
	Color        string    `json:"color"`
	IsSystem     bool      `json:"is_system"`
	IsActive     bool      `json:"is_active"`
	SortOrder    int       `json:"sort_order"`
	CreatedAt    time.Time `json:"created_at"`
}

// Transaction
type Transaction struct {
	ID              string             `json:"id"`
	UserID          string             `json:"user_id"`
	TransactionType string             `json:"transaction_type"` // INCOME, EXPENSE, TRANSFER, ADJUSTMENT
	OccurredAt      time.Time          `json:"occurred_at"`
	Description     string             `json:"description"`
	MerchantName    string             `json:"merchant_name"`
	Notes           string             `json:"notes"`
	TripID          *string            `json:"trip_id,omitempty"`
	Status          string             `json:"status"` // POSTED, VOIDED
	Version         int                `json:"version"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	AccountID       string             `json:"account_id,omitempty"`
	AmountMinor     int64              `json:"amount_minor"`
	CurrencyCode    string             `json:"currency_code"`
	CategoryName    string             `json:"category_name,omitempty"`
	CategoryID      string             `json:"category_id,omitempty"`
	Splits          []TransactionSplit `json:"splits,omitempty"`
	Entries         []AccountEntry     `json:"entries,omitempty"`
}

type TransactionSplit struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	CategoryID    string    `json:"category_id"`
	CategoryName  string    `json:"category_name,omitempty"`
	AmountMinor   int64     `json:"amount_minor"`
	CurrencyCode  string    `json:"currency_code"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// Budget
type Budget struct {
	ID             string           `json:"id"`
	UserID         string           `json:"user_id"`
	Name           string           `json:"name"`
	PeriodType     string           `json:"period_type"` // MONTHLY, YEARLY, CUSTOM
	StartDate      string           `json:"start_date"`
	EndDate        string           `json:"end_date"`
	CurrencyCode   string           `json:"currency_code"`
	LimitMinor     int64            `json:"limit_minor"`
	SpentMinor     int64            `json:"spent_minor"`
	RemainingMinor int64            `json:"remaining_minor"`
	UsagePercent   float64          `json:"usage_percent"`
	Status         string           `json:"status"` // ACTIVE, ARCHIVED
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	Categories     []BudgetCategory `json:"categories,omitempty"`
}

type BudgetCategory struct {
	ID                   string `json:"id"`
	BudgetID             string `json:"budget_id"`
	CategoryID           string `json:"category_id"`
	CategoryName         string `json:"category_name,omitempty"`
	AllocatedAmountMinor int64  `json:"allocated_amount_minor"`
	SpentAmountMinor     int64  `json:"spent_amount_minor,omitempty"`
}

// Subscriptions
type SubscriptionService struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Slug          string    `json:"slug"`
	Category      string    `json:"category"`
	Website       string    `json:"website"`
	LogoAssetID   string    `json:"logo_asset_id"`
	CatalogStatus string    `json:"catalog_status"`
	CreatedAt     time.Time `json:"created_at"`
}

type UserSubscription struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	ServiceID          *string   `json:"service_id,omitempty"`
	ServiceName        string    `json:"service_name"`
	CustomServiceName  *string   `json:"custom_service_name,omitempty"`
	PaymentAccountID   *string   `json:"payment_account_id,omitempty"`
	PaymentAccountName string    `json:"payment_account_name,omitempty"`
	AmountMinor        int64     `json:"amount_minor"`
	CurrencyCode       string    `json:"currency_code"`
	BillingCycle       string    `json:"billing_cycle"` // MONTHLY, YEARLY, WEEKLY
	BillingInterval    int       `json:"billing_interval"`
	StartDate          string    `json:"start_date"`
	NextRenewalAt      string    `json:"next_renewal_at"`
	Status             string    `json:"status"` // ACTIVE, PAUSED, CANCELLED_RECORD
	ReminderEnabled    bool      `json:"reminder_enabled"`
	Notes              string    `json:"notes"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// FX
type Currency struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	DecimalDigits int    `json:"decimal_digits"`
	IsActive      bool   `json:"is_active"`
}

type ExchangeRate struct {
	ID            string    `json:"id"`
	ProviderID    string    `json:"provider_id"`
	ProviderName  string    `json:"provider_name,omitempty"`
	BaseCurrency  string    `json:"base_currency"`
	QuoteCurrency string    `json:"quote_currency"`
	Rate          float64   `json:"rate"`
	EffectiveAt   time.Time `json:"effective_at"`
	RetrievedAt   time.Time `json:"retrieved_at"`
	Status        string    `json:"status"` // VERIFIED, ESTIMATED, STALE
}

// Travel
type Trip struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	CountryCode  string    `json:"country_code"`
	CityName     string    `json:"city_name"`
	StartDate    string    `json:"start_date"`
	EndDate      string    `json:"end_date"`
	BaseCurrency string    `json:"base_currency"`
	Status       string    `json:"status"` // PLANNED, ACTIVE, COMPLETED, CANCELLED
	Notes        string    `json:"notes"`
	TotalBudget  int64     `json:"total_budget_minor"`
	TotalSpent   int64     `json:"total_spent_minor"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Notification
type Notification struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Body        string     `json:"body"`
	Channel     string     `json:"channel"`
	Status      string     `json:"status"`
	ScheduledAt time.Time  `json:"scheduled_at"`
	SentAt      time.Time  `json:"sent_at"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Support Ticket
type SupportTicket struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	UserEmail       string           `json:"user_email,omitempty"`
	Subject         string           `json:"subject"`
	Category        string           `json:"category"`
	Description     string           `json:"description"`
	Priority        string           `json:"priority"`
	Status          string           `json:"status"`
	AssignedAdminID *string          `json:"assigned_admin_id,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	ResolvedAt      *time.Time       `json:"resolved_at,omitempty"`
	Messages        []SupportMessage `json:"messages,omitempty"`
}

type SupportMessage struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticket_id"`
	SenderType string    `json:"sender_type"` // USER, ADMIN, SYSTEM
	SenderID   string    `json:"sender_id"`
	Message    string    `json:"message"`
	Visibility string    `json:"visibility"`
	CreatedAt  time.Time `json:"created_at"`
}

// Admin
type AdminUser struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	FullName   string    `json:"full_name"`
	Status     string    `json:"status"`
	MFAEnabled bool      `json:"mfa_enabled"`
	Roles      []string  `json:"roles,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID               string    `json:"id"`
	ActorType        string    `json:"actor_type"`
	ActorID          string    `json:"actor_id"`
	Action           string    `json:"action"`
	ResourceType     string    `json:"resource_type"`
	ResourceID       string    `json:"resource_id"`
	Result           string    `json:"result"`
	MetadataRedacted string    `json:"metadata_redacted"`
	RequestID        string    `json:"request_id"`
	CreatedAt        time.Time `json:"created_at"`
}

type FeatureFlag struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	Environment string    `json:"environment"`
	Enabled     bool      `json:"enabled"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

