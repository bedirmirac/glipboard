package storage

import (
	"database/sql"
	"fmt"
)

func (s *Storage) Save(context string) error {
	instertSQL := `INSERT INTO clipboard (context) VALUES (?)`

	_, err := s.db.Exec(instertSQL, context)
	if err != nil {
		return fmt.Errorf("error during inserting context to database: %v", err)
	}
	return nil
}

func (s *Storage) Delete(id int) error {
	query := `DELETE FROM clipboard WHERE id = ?`

	res, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error during deleting: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error during deleting: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("there is probably no context with id %v. Nothing deleted", id)
	}
	return nil
}

func (s *Storage) DeleteAll() error {
	query := `DELETE FROM clipboard;`

	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error during deleting all: %v", err)
	}

	return nil
}

func (s *Storage) IsLimitExceeded() (bool, error) {
	var dummy int

	query := `SELECT 1 FROM clipboard LIMIT 1 OFFSET 50`

	err := s.db.QueryRow(query).Scan(&dummy)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error during checking if limit is exceeded: %v", err)
	}
	return true, nil
}

func (s *Storage) DeleteOldestRecord() error {
	query := `
		DELETE FROM clipboard 
		WHERE id = (
			SELECT id FROM clipboard ORDER BY id ASC LIMIT 1
		)
	`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error during deleting the oldest record: %v", err)
	}

	return nil
}
