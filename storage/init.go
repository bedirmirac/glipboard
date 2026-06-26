package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error locating home directory: %v", err)
	}

	appDir := filepath.Join(homeDir, ".config", "glipboard")

	err = os.MkdirAll(appDir, 0o755)
	if err != nil {
		return nil, fmt.Errorf("error creating directory: %v", err)
	}

	dbPath := filepath.Join(appDir, "clipboard.db")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS clipboard (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		context TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return &Storage{
		db: db,
	}, nil
}
