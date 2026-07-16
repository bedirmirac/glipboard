package storage

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
)

const (
	ErrSqliteConstraintPrimaryKey = 1555
	ErrSqliteConstraintUnique     = 2067
)

var ErrAlreadyExists = errors.New("It already exists")

func (s *Storage) Save(dataType string, rawData []byte, textContent string, filePath string) error {
	hashBytes := sha256.Sum256(rawData)
	hashString := fmt.Sprintf("%x", hashBytes)

	insertQuery := `INSERT INTO clipboard (hash, type, context, file_path) VALUES (?, ?, ?, ?)`

	_, err := s.db.Exec(insertQuery, hashString, dataType, textContent, filePath)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) {
			errCode := sqliteErr.Code()
			if errCode == ErrSqliteConstraintPrimaryKey || errCode == ErrSqliteConstraintUnique {
				return ErrAlreadyExists
			} else {
				return fmt.Errorf("error during saving to database: %v", err)
			}
		}
	}
	return nil
}

func (s *Storage) Delete(hash string) error {
	query := `DELETE FROM clipboard WHERE hash = ?`

	res, err := s.db.Exec(query, hash)
	if err != nil {
		return fmt.Errorf("error during deleting: %v", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error during deleting: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("there is probably no context with hash %v. Nothing deleted", hash)
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
		WHERE hash = (
			SELECT hash FROM clipboard ORDER BY created_at ASC, rowid ASC LIMIT 1
		)
	`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error during deleting the oldest record: %v", err)
	}

	return nil
}
