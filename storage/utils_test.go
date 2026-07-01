package storage

import (
	"database/sql"
	"errors"
	"testing"
)

func TestDeleteOldestRecord(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	tests := []string{"text", "text1", "text2"}
	for _, test := range tests {
		q := `INSERT INTO clipboard (context) VALUES (?)`
		_, err := db.Exec(q, test)
		if err != nil {
			t.Fatalf("error during inserting test data: %v", err)
		}
	}
	err := s.DeleteOldestRecord()
	if err != nil {
		t.Fatalf("no error expected, but there is an error: %v", err)
	}

	rows, err := db.Query("SELECT id, context FROM clipboard ORDER BY id DESC")
	if err != nil {
		t.Fatalf("error during fetching data: %v", err)
	}
	t.Cleanup(func() {
		rows.Close()
	})
	var items []Clipboard
	for rows.Next() {
		var item Clipboard
		err := rows.Scan(&item.ID, &item.Context)
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
		t.Fatalf("expected number of data is 2 , but %v number of record was founded", len(items))
	}

	for _, item := range items {
		if item.Context == "text" {
			t.Fatalf("function didn't work correctly")
		}
	}
}

func TestIsLimitExceeded(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}

	q := `INSERT INTO clipboard (context) VALUES (?)`
	str := "text"
	for range 50 {
		_, err := db.Exec(q, str)
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
	_, err = db.Exec(q, str)
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
	q := `INSERT INTO clipboard (context) VALUES (?)`
	tests := []string{"text", "test1", "test2", "test3"}
	for _, test := range tests {
		_, err := db.Exec(q, test)
		if err != nil {
			t.Fatalf("error during insrting test data to mock database: %v", err)
		}
	}
	err := s.DeleteAll()
	if err != nil {
		t.Fatalf("function (DeleteAll) didn't run successfully: %v", err)
	}
	var exists bool
	q = `SELECT EXISTS (SELECT 1 FROM clipboard);`

	err = db.QueryRow(q).Scan(&exists)
	if err != nil {
		t.Fatalf("error during scanning if any data exists: %v", err)
	}

	if exists {
		t.Fatalf("expected no data exists, but there are some data exist: %v", err)
	}
}

func TestDelete(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	tests := []struct {
		name      string
		text      string
		id        int
		expectErr bool
	}{
		{"Non-exists id", "Text", -1, true},
		{"Correct value", "Text1", 1, false},
	}

	q := `INSERT INTO clipboard (context) VALUES (?)`
	for i := 1; i < len(tests); i++ {
		_, err := db.Exec(q, tests[i].text)
		if err != nil {
			t.Fatalf("error during saving the test data to test db: %v", err)
		}
	}

	for _, test := range tests {
		err := s.Delete(test.id)

		if test.expectErr {
			if err == nil {
				t.Errorf("[%s] expected error, but ran successfully", test.name)
			}
			continue
		} else {
			if err != nil {
				t.Errorf("[%s] unexpected error: %v", test.name, err)
			}
		}

		var ctx string
		err = db.QueryRow(`SELECT context FROM clipboard WHERE id = ?`, test.id).Scan(&ctx)

		if err == nil {
			t.Errorf("[%s] data with id %d should be deleted but it still exists", test.name, test.id)
		} else if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("[%s] expected sql.ErrNoRows, but got: %v", test.name, err)
		}
	}
}

func TestSave(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	tests := []struct {
		name        string
		text        string
		expectError bool
	}{
		{"Empty string", "", false},
		{"Correct value", "testValue", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := s.Save(test.text)
			if test.expectError {
				if err == nil {
					t.Errorf("expected error, but there is no error")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			var actualString string
			query := `SELECT context FROM clipboard ORDER BY id DESC LIMIT 1`
			err = db.QueryRow(query).Scan(&actualString)
			if err != nil {
				t.Fatalf("error during scanning data from database: %v", err)
			}

			if actualString != test.text {
				t.Errorf("expected value %v, but value saved is %v", test.text, actualString)
			}
		})
	}
}
