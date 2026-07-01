package cmd

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/bedirmirac/glipboard/storage"
)

func setupLogger() *os.File {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	configDir := filepath.Join(homeDir, ".config", "glipboard")

	err = os.MkdirAll(configDir, 0o755)
	if err != nil {
		return nil
	}

	logPath := filepath.Join(configDir, ".glipboard.log")
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err == nil {
		log.SetOutput(logFile)
		return logFile
	}

	return nil
}

func StartDaemon() {
	listener, err := net.Listen("tcp", "127.0.0.1:49321")
	if err != nil {
		return
	}

	defer listener.Close()

	logFile := setupLogger()
	if logFile != nil {
		defer logFile.Close()
	}

	s, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("error during starting the database: %v", err)
	}

	var text storage.Clipboard
	lastCopied, _ := clipboard.ReadAll()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		newContent, err := clipboard.ReadAll()
		if err == nil && newContent != "" && newContent != lastCopied {
			text.Context = newContent
			err := s.Save(text.Context)
			if err != nil {
				log.Printf("error during saving the content: %v", err)
			} else {
				lastCopied = newContent
			}
			isExceeded, err := s.IsLimitExceeded()
			if err != nil {
				log.Printf("error after saving: %v", err)
			}
			if isExceeded {
				err := s.DeleteOldestRecord()
				if err != nil {
					log.Printf("error while deleting the oldest record: %v", err)
				}
			}

		}
	}
}
