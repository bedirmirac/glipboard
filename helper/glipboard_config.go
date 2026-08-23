package helper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type LimitConfig struct {
	Limit int `json:"limit"`
}

func WriteLimit(limit int) error {
	folderPath, err := GetConfigFolder()

	if err != nil {
		return err
	}

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return fmt.Errorf("config.json couldn't be created: %v", err)
	}

	fullPath := filepath.Join(folderPath, "config.json")

	config := LimitConfig{
		Limit: limit,
	}

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fullPath, jsonData, 0644)
}

func ReadLimit() (int, error) {
	folderPath, err := GetConfigFolder()

	if err != nil {
		return 50, err
	}
	fullPath := filepath.Join(folderPath, "config.json")

	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return 50, err
	}

	var config LimitConfig
	if err := json.Unmarshal(fileData, &config); err != nil {
		return 50, err
	}

	return config.Limit, nil
}
