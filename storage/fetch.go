package storage

import (
	"fmt"
)

type Clipboard struct {
	ID      int
	Context string
}

func (s *Storage) Fetch() ([]Clipboard, error) {
	rows, err := s.db.Query("SELECT id, context FROM clipboard ORDER BY id DESC")
	if err != nil {
		return nil, fmt.Errorf("error during fetching data: %v", err)
	}
	defer rows.Close()
	var items []Clipboard
	for rows.Next() {
		var item Clipboard
		err := rows.Scan(&item.ID, &item.Context)
		if err != nil {
			return nil, fmt.Errorf("error during fetching data from a row: %v", err)
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error during scanning: %v", err)
	}
	return items, nil
}
