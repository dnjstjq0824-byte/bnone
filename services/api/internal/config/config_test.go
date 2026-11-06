package config

import (
	"testing"
)

func TestSplitOrigins_StripsWildcard(t *testing.T) {
	got := splitOrigins("http://localhost:3000, *, http://localhost:3001")
	if len(got) != 2 {
		t.Fatalf("expected wildcard stripped, got %v", got)
	}
	for _, o := range got {
		if o == "*" {
			t.Fatalf("wildcard origin leaked through: %v", got)
		}
	}
}

func TestSplitOrigins_EmptyEntries(t *testing.T) {
	if got := splitOrigins("  , "); len(got) != 0 {
		t.Fatalf("expected no origins, got %v", got)
	}
}

func TestLoad_ProductionRefusesMissingSecret(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("JWT_SECRET", "")
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when JWT_SECRET is missing in production")
		}
	}()
	Load()
}

func TestLoad_ProductionRefusesShortSecret(t *testing.T) {
	t.Setenv("ENVIRONMENT", "production")
	t.Setenv("JWT_SECRET", "too-short")
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when JWT_SECRET is shorter than 32 chars in production")
		}
	}()
	Load()
}

func TestLoad_DevelopmentFallsBack(t *testing.T) {
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("JWT_SECRET", "")
	cfg := Load()
	if cfg.JWTSecret == "" {
		t.Fatal("expected dev fallback secret")
	}
	if len(cfg.CORSAllowedOrigins) == 0 {
		t.Fatal("expected default CORS origins")
	}
}

func TestLoad_ReadsLogLevelAndBootstrapFlags(t *testing.T) {
	t.Setenv("ENVIRONMENT", "development")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("ADMIN_BOOTSTRAP_ENABLED", "true")
	t.Setenv("ADMIN_BOOTSTRAP_PASSWORD", "seed-only-secret")
	cfg := Load()
	if cfg.LogLevel != "debug" {
		t.Fatalf("expected LOG_LEVEL=debug, got %q", cfg.LogLevel)
	}
	if !cfg.AdminBootstrapEnabled || cfg.AdminBootstrapPassword != "seed-only-secret" {
		t.Fatal("expected bootstrap flags to be wired from env")
	}
}
