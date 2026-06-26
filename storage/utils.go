package storage

import "fmt"

func (s *Storage) Save(context Clipboard) error {
	instertSQL := `INSERT INTO clipboard (context) VALUES (?)`

	_, err := s.db.Exec(instertSQL, context.Context)
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
