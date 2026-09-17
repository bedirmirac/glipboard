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

	isExceeded, err := s.IsLimitExceeded(50)
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

	isExceeded, err = s.IsLimitExceeded(50)
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

func CountTest(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	tests := []struct {
		name string
		text string
		hash string
	}{
		{"test", "Text", "as23214as"},
		{"test1", "Text1", "aggf2312a"},
	}

	q := `INSERT INTO clipboard (hash, type, context, file_path) VALUES (?, 'text', ?, '')`
	for _, test := range tests {
		_, err := db.Exec(q, test.hash, test.text)
		if err != nil {
			t.Fatalf("error during saving the test data to test db: %v", err)
		}
	}
	functionCount, err := s.Count()
	var testCount int
	if err != nil {
		t.Fatalf("error Count(): %v", err)
	}
	query := `SELECT COUNT(*) FROM clipboard`
	err = db.QueryRow(query).Scan(&testCount)
	if err != nil {
		t.Fatalf("error during counting test data: %v", err)
	}

	if functionCount != testCount {
		t.Fatalf("function Count() and test Count() are not the same: %v", err)
	}
}
func TestTrimToLimit(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}

	insertQuery := `INSERT INTO clipboard (hash, type, context, file_path) VALUES (?, ?, ?, ?)`

	// 1'den 10'a kadar 10 adet test kaydı ekle (rowid'ler: 1, 2, ..., 10)
	for i := 1; i <= 10; i++ {
		hash := fmt.Sprintf("test_hash_%d", i)
		_, err := db.Exec(insertQuery, hash, "text", "test_context", "test_path")
		if err != nil {
			t.Fatalf("failed to insert test data: %v", err)
		}
	}

	// Senaryo 1: 10 kayıt varken limiti 5'e düşür (en eski 1..5 silinmeli, 6..10 kalmalı)
	newLimit := 5
	if err := s.TrimToLimit(newLimit); err != nil {
		t.Fatalf("TrimToLimit(%d) failed: %v", newLimit, err)
	}

	// Kalan toplam satır sayısını doğrula
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM clipboard`).Scan(&count); err != nil {
		t.Fatalf("failed to query count: %v", err)
	}
	if count != newLimit {
		t.Fatalf("expected %d rows after trim, got %d", newLimit, count)
	}

	// Kalan satırların gerçekten en güncel 5 satır (6, 7, 8, 9, 10) olduğunu doğrula
	rows, err := db.Query(`SELECT rowid FROM clipboard ORDER BY rowid ASC`)
	if err != nil {
		t.Fatalf("failed to query remaining rowids: %v", err)
	}
	defer rows.Close()

	var remainingIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("failed to scan rowid: %v", err)
		}
		remainingIDs = append(remainingIDs, id)
	}

	expectedIDs := []int{6, 7, 8, 9, 10}
	if len(remainingIDs) != len(expectedIDs) {
		t.Fatalf("expected remaining IDs len %d, got %d", len(expectedIDs), len(remainingIDs))
	}
	for i, expected := range expectedIDs {
		if remainingIDs[i] != expected {
			t.Errorf("at index %d: expected rowid %d, got %d", i, expected, remainingIDs[i])
		}
	}

	// Senaryo 2: Tabloda 5 kayıt varken daha büyük bir limit (örn. 10) ver
	// Hiçbir kayıt silinmemeli, mevcut 5 kayıt korunmalı
	if err := s.TrimToLimit(10); err != nil {
		t.Fatalf("TrimToLimit(10) failed: %v", err)
	}

	if err := db.QueryRow(`SELECT COUNT(*) FROM clipboard`).Scan(&count); err != nil {
		t.Fatalf("failed to query count after second trim: %v", err)
	}
	if count != 5 {
		t.Errorf("expected count to remain 5 when limit > count, got %d", count)
	}
}
