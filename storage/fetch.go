package storage

import (
	"fmt"
)

type Clipboard struct {
	Hash     string
	Context  string
	Type     string
	FilePath string
}

func (s *Storage) Fetch() ([]Clipboard, error) {
	rows, err := s.db.Query("SELECT hash, context, type, file_path FROM clipboard ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("error during fetching data: %v", err)
	}
	defer rows.Close()
	var items []Clipboard
	for rows.Next() {
		var item Clipboard
		err := rows.Scan(&item.Hash, &item.Context, &item.Type, &item.FilePath)
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
