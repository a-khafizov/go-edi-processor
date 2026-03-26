package deps

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(cfg *Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBUser, cfg.DBPassword, cfg.DBSSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
