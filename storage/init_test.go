package storage

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewStorage(t *testing.T) {
	tempDir := t.TempDir()

	if runtime.GOOS == "windows" {
		t.Setenv("USERPROFILE", tempDir)
	} else {
		t.Setenv("HOME", tempDir)
	}

	s, err := NewStorage()
	if err != nil {
		t.Fatalf("newStorage returned error: %v", err)
	}
	t.Cleanup(func() {
		s.db.Close()
	})

	expectedPath := filepath.Join(tempDir, ".config", "glipboard", "clipboard.db")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("database file was not created at the expected path: %v", expectedPath)
	}

	var name string
	err = s.db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='clipboard';").Scan(&name)
	if err != nil {
		t.Errorf("clipboard couldn't be found in database: %v", err)
	}
}
