package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// ConnectDatabase membuka koneksi PostgreSQL menggunakan DATABASE_URL.
func ConnectDatabase() (*sql.DB, error) {
	// File .env digunakan saat development lokal. Di Railway, konfigurasi
	// dibaca langsung dari environment variables yang disediakan platform.
	_ = godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL belum diatur")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("membuka koneksi database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("menghubungi database: %w", err)
	}

	return db, nil
}
