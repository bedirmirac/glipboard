package storage

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteOldestRecord(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	tests := []struct {
		hash    string
		context string
	}{
		{"a_hash", "text1"},
		{"b_hash", "text2"},
		{"c_hash", "text3"},
	}

	q := `INSERT INTO clipboard (hash, type, context, file_path) VALUES (?, 'text', ?, '')`
	for _, test := range tests {
		_, err := db.Exec(q, test.hash, test.context)
		if err != nil {
			t.Fatalf("error during inserting test data: %v", err)
		}
	}

	err := s.DeleteOldestRecord()
	if err != nil {
		t.Fatalf("no error expected, but there is an error: %v", err)
	}

	rows, err := db.Query("SELECT hash, context FROM clipboard ORDER BY hash ASC")
	if err != nil {
		t.Fatalf("error during fetching data: %v", err)
	}
	t.Cleanup(func() {
		rows.Close()
	})

	var items []Clipboard
	for rows.Next() {
		var item Clipboard
		err := rows.Scan(&item.Hash, &item.Context)
		if err != nil {
			t.Fatalf("error during fetching data from a row: %v", err)
		}
		items = append(items, item)
	}

	err = rows.Err()
	if err != nil {
		t.Fatalf("error during scanning: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("expected number of data is 2, but %v number of record was founded", len(items))
	}

	for _, item := range items {
		if item.Hash == "a_hash" {
			t.Fatalf("function didn't work correctly, oldest record still exists")
		}
	}
}

func TestIsLimitExceeded(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}

	q := `INSERT INTO clipboard (hash, type, context, file_path) VALUES (?, 'text', ?, '')`
	str := "text"

	for i := 0; i < 50; i++ {
		_, err := db.Exec(q, fmt.Sprintf("hash_%d", i), str)
		if err != nil {
			t.Fatalf("error during inserting test data: %v", err)
		}
	}

	isExceeded, err := s.IsLimitExceeded()
	if err != nil {
		t.Fatalf("no error expected, but there is an error: %v", err)
	}
	if isExceeded {
		t.Errorf("it shouldn't be returned true")
	}

	_, err = db.Exec(q, "hash_50", str)
	if err != nil {
		t.Fatalf("error during inserting test data: %v", err)
	}

	isExceeded, err = s.IsLimitExceeded()
	if err != nil {
		t.Fatalf("no error expected, but there is an error: %v", err)
	}
	if !isExceeded {
		t.Fatalf("it should be returned true, but returned false")
	}
}

func TestDeleteAll(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	q := `INSERT INTO clipboard (hash, type, context, file_path) VALUES (?, 'text', ?, '')`
	tests := []string{"text", "test1", "test2", "test3"}

	for i, test := range tests {
		_, err := db.Exec(q, fmt.Sprintf("hash_%d", i), test)
		if err != nil {
			t.Fatalf("error during inserting test data to mock database: %v", err)
		}
	}

	err := s.DeleteAll()
	if err != nil {
		t.Fatalf("function (DeleteAll) didn't run successfully: %v", err)
	}

	var exists bool
	qCheck := `SELECT EXISTS (SELECT 1 FROM clipboard);`

	err = db.QueryRow(qCheck).Scan(&exists)
	if err != nil {
		t.Fatalf("error during scanning if any data exists: %v", err)
	}

	if exists {
		t.Fatalf("expected no data exists, but there are some data exist")
	}

	err = s.DeleteAll()
	if err != nil {
		t.Fatalf("there shouldn't be an error: %v", err)
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	tests := []struct {
		name      string
		text      string
		hash      string
		expectErr bool
	}{
		{"Non-exists hash", "Text", "iAmNotARealHash", true},
		{"Correct value", "Text1", "realHash1", false},
	}

	q := `INSERT INTO clipboard (hash, type, context, file_path) VALUES (?, 'text', ?, '')`
	for _, test := range tests {
		if !test.expectErr {
			_, err := db.Exec(q, test.hash, test.text)
			if err != nil {
				t.Fatalf("error during saving the test data to test db: %v", err)
			}
		}
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := s.Delete(test.hash)

			if test.expectErr {
				if err == nil {
					t.Errorf("[%s] expected error, but ran successfully", test.name)
				}
				return
			} else {
				if err != nil {
					t.Errorf("[%s] unexpected error: %v", test.name, err)
				}
			}

			var ctx string
			err = db.QueryRow(`SELECT context FROM clipboard WHERE hash = ?`, test.hash).Scan(&ctx)

			if err == nil {
				t.Errorf("[%s] data with hash %v should be deleted but it still exists", test.name, test.hash)
			} else if !errors.Is(err, sql.ErrNoRows) {
				t.Errorf("[%s] expected sql.ErrNoRows, but got: %v", test.name, err)
			}
		})
	}
}

func TestSave(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	tempDir := t.TempDir()
	validImgPath := filepath.Join(tempDir, "test_image.png")

	err := os.WriteFile(validImgPath, []byte("fake image data"), 0o644)
	if err != nil {
		t.Errorf("there's an error during writing the fake img to temp dir: %v", err)
	}

	tests := []struct {
		name        string
		dataType    string
		rawData     []byte
		text        string
		filepath    string
		expectError bool
	}{
		{"Empty string", "text", nil, "", "", false},
		{"Correct value", "text", []byte("testValue"), "testValue", "", false},
		{"Emojis", "text", []byte("Hello 🌍"), "Hello 🌍", "", false},
		{"SQL Injection", "text", []byte(`It's a "test" string; DROP TABLE;`), `It's a "test" string; DROP TABLE;`, "", false},
		{"Image Test", "image", []byte("fake image data"), "", validImgPath, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := s.Save(test.dataType, test.rawData, test.text, test.filepath)

			if test.expectError {
				if err == nil {
					t.Errorf("expected error, but there is no error")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			hashBytes := sha256.Sum256(test.rawData)
			expectedHash := fmt.Sprintf("%x", hashBytes)

			switch test.dataType {
			case "text":
				var actualString string
				query := `SELECT context FROM clipboard WHERE hash = ?`
				err = db.QueryRow(query, expectedHash).Scan(&actualString)
				if err != nil {
					t.Fatalf("error during scanning data from database: %v", err)
				}
				if actualString != test.text {
					t.Errorf("expected value %v, but value saved is %v", test.text, actualString)
				}

			case "image":
				var actualPath string
				query := `SELECT file_path FROM clipboard WHERE hash = ?`
				err = db.QueryRow(query, expectedHash).Scan(&actualPath)
				if err != nil {
					t.Fatalf("error during scanning data from database: %v", err)
				}
				if actualPath != test.filepath {
					t.Errorf("expected value %v, but value saved is %v", test.filepath, actualPath)
				}

			default:
				t.Errorf("unsupported data type")
			}
		})
	}
}
