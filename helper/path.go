package helper

import (
	"os"
	"path/filepath"
)

func GetConfigFolder() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	folder := filepath.Join(homeDir, ".config", "glipboard")

	err = os.MkdirAll(folder, 0o755)
	if err != nil {
		return "", err
	}

	return folder, nil
}
