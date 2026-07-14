package storage

import (
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestFetch(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	tests := []struct {
		hash     string
		dataType string
		context  string
		filePath string
	}{
		{"hash1", "text", "text1", "/path/1"},
		{"hash2", "text", "text2", "/path/2"},
	}

	q := `INSERT INTO clipboard (hash, type, context, file_path, created_at) VALUES (?, ?, ?, ?, ?)`

	for i, test := range tests {
		createdAt := time.Now().Add(time.Duration(i) * time.Second)

		_, err := db.Exec(q, test.hash, test.dataType, test.context, test.filePath, createdAt)
		if err != nil {
			t.Fatalf("error during saving test values to database: %v", err)
		}
	}

	results, err := s.Fetch()
	if err != nil {
		t.Fatalf("fetch function didn't run without error: %v", err)
	}

	if len(results) != len(tests) {
		t.Fatalf("expected number of data %v, but %v many data returned", len(tests), len(results))
	}

	for i := 0; i < len(results); i++ {
		expectedIndex := len(tests) - 1 - i

		if tests[expectedIndex].context != results[i].Context {
			t.Errorf("expected context %v, but fetched value is %v", tests[expectedIndex].context, results[i].Context)
		}
		if tests[expectedIndex].hash != results[i].Hash {
			t.Errorf("expected hash %v, but fetched value is %v", tests[expectedIndex].hash, results[i].Hash)
		}
	}
}
