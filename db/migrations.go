package db

import (
	"database/sql"
	"fmt"
)

func Migrate(db *sql.DB) error {
	// Create profiles table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY,
			user_id TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			username TEXT,
			bio TEXT,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating profiles table: %w", err)
	}

	// Create quotes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS quotes (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			quote TEXT NOT NULL,
			approved BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("error creating quotes table: %w", err)
	}

	fmt.Println("Migration successful.")
	return nil
}
