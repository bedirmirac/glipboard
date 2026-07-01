package storage

import (
	"testing"

	_ "modernc.org/sqlite"
)

func TestFetch(t *testing.T) {
	db := setupTestDb(t)
	s := &Storage{db: db}
	tests := []struct {
		name string
		text string
	}{
		{"Test1", "text1"},
		{"Test2", "text2"},
	}
	q := `INSERT INTO clipboard (context) VALUES (?)`
	for _, test := range tests {
		_, err := db.Exec(q, test.text)
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
		if tests[expectedIndex].text != results[i].Context {
			t.Errorf("expected %v, but fetched value is %v", tests[expectedIndex].text, results[i].Context)
		}
	}
}
