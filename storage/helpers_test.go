package storage

import (
	"database/sql"
	"testing"
)

func setupTestDb(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Test db couldn't be opened: %v", err)
	}

	query := `CREATE TABLE clipboard (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	context TEXT NOT NULL
	);
	`
	_, err = db.Exec(query)
	if err != nil {
		t.Fatalf("Table couldn't be created: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}
