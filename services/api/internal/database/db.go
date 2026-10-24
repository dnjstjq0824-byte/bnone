package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
	Driver string
}

func Connect(driver, dsn string) (*DB, error) {
	var sqlDriver string
	if driver == "postgres" || strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		sqlDriver = "pgx"
		driver = "postgres"
	} else {
		sqlDriver = "sqlite"
		driver = "sqlite"
	}

	db, err := sql.Open(sqlDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("[DB] Connected successfully using driver: %s", driver)
	return &DB{DB: db, Driver: driver}, nil
}

func (db *DB) Migrate(migrationsPath, fixturesPath string) error {
	log.Printf("[DB] Running migrations from %s", migrationsPath)

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("reading migrations dir: %w", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			fullPath := filepath.Join(migrationsPath, file.Name())
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("reading migration file %s: %w", file.Name(), err)
			}

			sqlStatements := string(content)
			// Execute statements
			if err := db.executeMultiSQL(sqlStatements); err != nil {
				return fmt.Errorf("executing migration %s: %w", file.Name(), err)
			}
			log.Printf("[DB] Applied migration: %s", file.Name())
		}
	}

	// Apply fixtures if available
	if fixturesPath != "" {
		if fixtureFiles, err := os.ReadDir(fixturesPath); err == nil {
			for _, file := range fixtureFiles {
				if strings.HasSuffix(file.Name(), ".sql") {
					fullPath := filepath.Join(fixturesPath, file.Name())
					content, err := os.ReadFile(fullPath)
					if err == nil {
						if err := db.executeMultiSQL(string(content)); err != nil {
							log.Printf("[DB] Note on fixture %s: %v", file.Name(), err)
						} else {
							log.Printf("[DB] Applied fixture: %s", file.Name())
						}
					}
				}
			}
		}
	}

	return nil
}

func (db *DB) executeMultiSQL(script string) error {
	// Split by semicolon, clean comments and empty statements
	stmts := strings.Split(script, ";")
	for _, stmt := range stmts {
		trimmed := strings.TrimSpace(stmt)
		if trimmed == "" {
			continue
		}
		// If SQLite, replace PostgreSQL specific syntax if any
		if db.Driver == "sqlite" {
			trimmed = strings.ReplaceAll(trimmed, "TIMESTAMP WITH TIME ZONE", "DATETIME")
			// SQLite supports ON CONFLICT (col) DO NOTHING
		}
		if _, err := db.Exec(trimmed); err != nil {
			// Ignore conflict errors during re-migration
			if !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "UNIQUE constraint failed") {
				return fmt.Errorf("failed executing stmt [%s]: %w", trimmed, err)
			}
		}
	}
	return nil
}

