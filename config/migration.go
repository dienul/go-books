package config

import (
	"database/sql"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

// RunMigrations menerapkan migration yang belum dijalankan.
func RunMigrations(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{Dir: "migrations"}

	applied, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("menjalankan migration: %w", err)
	}

	if applied > 0 {
		fmt.Printf("migration diterapkan: %d\n", applied)
	}

	return nil
}
