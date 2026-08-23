package storage

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/bedirmirac/glipboard/helper"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	appDir, err := helper.GetConfigFolder()
	dbPath := filepath.Join(appDir, "clipboard.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS clipboard (
		hash TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		context TEXT,
		file_path TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("error creating table: %v", err)
	}

	return &Storage{
		db: db,
	}, nil
}
