package cmd

import (
	"context"
	"fmt"
	_ "image/jpeg"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/bedirmirac/glipboard/storage"
	"golang.design/x/clipboard"
	_ "golang.org/x/image/webp"
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
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	s, err := storage.NewStorage()
	if err != nil {
		log.Fatalf("error during starting the database: %v", err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:49321")
	if err != nil {
		return
	}

	defer listener.Close()
	/* --- HTTP SERVER: Listens for requests from the TUI --- */
	mux := http.NewServeMux()
	mux.HandleFunc("/request", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is supported", http.StatusMethodNotAllowed)
			return
		}

		dataType := r.URL.Query().Get("type")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error: Could not read request body: %v", err)
			http.Error(w, "Could not read body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		switch dataType {
		case "text":
			clipboard.Write(clipboard.FmtText, body)
			log.Println("Text from client successfully written to clipboard.")
			w.WriteHeader(http.StatusOK)

		case "image":
			filePath := string(body)
			imgData, err := os.ReadFile(filePath)
			if err != nil {
				log.Printf("Error: Could not read from image file (%v): %v", filePath, err)
				http.Error(w, "Could not read image file", http.StatusInternalServerError)
				return
			}
			clipboard.Write(clipboard.FmtImage, imgData)
			log.Println("Image from client successfully written to clipboard.")
			w.WriteHeader(http.StatusOK)
		case "delete":
			hash := string(body)
			err := s.Delete(hash)
			if err != nil {
				log.Printf("Error: Could not delete the data  (%v): %v", hash, err)
				http.Error(w, "Could not delete the data", http.StatusInternalServerError)
				return
			}
		case "deleteImagePath":
			path := string(body)
			err = os.Remove(path)
			if err != nil {
				log.Printf("Error: Could not delete the image  (%v): %v", path, err)
				http.Error(w, "Could not delete the image", http.StatusInternalServerError)
				return
			}
		case "deleteAll":
			err := s.DeleteAll()
			if err != nil {
				log.Println("Error: Could not run the delete all function")
				http.Error(w, "Could not run the delete all function", http.StatusInternalServerError)
				return
			}
			path, err := getPath()
			path = filepath.Dir(path)
			if err != nil {
				log.Printf("Error: Could not get the image path  (%v): %v", path, err)
				http.Error(w, "Could not find the image path", http.StatusInternalServerError)
				return
			}
			err = os.RemoveAll(path)
			if err != nil {
				log.Printf("Error: Could not remove the Pictures folder: %v", err)
				http.Error(w, "Could not remove the Pictures folder", http.StatusInternalServerError)
				return
			}
		default:
			log.Printf("Error: Unknown data type received: %s", dataType)
			http.Error(w, "Unknown data type", http.StatusBadRequest)
		}
	})

	go func() {
		if err := http.Serve(listener, mux); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	logFile := setupLogger()
	if logFile != nil {
		defer logFile.Close()
	}

	os := runtime.GOOS
	isWayland := getEnv()
	switch os {
	case "linux":
		if isWayland {
			eventDriven(s)
		} else {
			polling(s)
		}
	default:
		polling(s)
	}
}

func eventDriven(s *storage.Storage) {
	ch := clipboard.Watch(context.TODO())

	for data := range ch {
		switch data.Format {
		case clipboard.FmtText:
			textContent := string(data.Bytes)
			if originalPath, isImage := getLocalImagePath(textContent); isImage {
				imgData, err := os.ReadFile(originalPath)
				if err == nil {
					imgPath, err := getPath()
					if err == nil {
						err = os.WriteFile(imgPath, imgData, 0o644)
						if err != nil {
							log.Printf("error during writing local image: %v", err)
						}
						err = s.Save("image", imgData, "", imgPath)
						if err != nil {
							log.Printf("error during saving data to database: %v", err)
						}
						continue
					}
				}
			}
			err := s.Save("text", data.Bytes, textContent, "")
			if err != nil {
				log.Printf("error in save function: %v", err)
			}
		case clipboard.FmtImage:

			imgPath, err := getPath()
			if err != nil {
				log.Printf("error during getting image path: %v", err)
			}
			err = os.WriteFile(imgPath, data.Bytes, 0o644)
			if err != nil {
				log.Printf("error during writing the image: %v", err)
			}
			err = s.Save("image", data.Bytes, "", imgPath)
			if err != nil {
				log.Printf("error in save function: %v", err)
			}
		}
	}
}

func polling(s *storage.Storage) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C
		newText := clipboard.Read(clipboard.FmtText)
		newImg := clipboard.Read(clipboard.FmtImage)
		if string(newText) != "" {
			if originalPath, isImage := getLocalImagePath(string(newText)); isImage {
				imgData, err := os.ReadFile(originalPath)
				if err == nil {
					imgPath, err := getPath()
					if err == nil {
						err := os.WriteFile(imgPath, imgData, 0o644)
						if err != nil {
							log.Printf("error during writing local image: %v", err)
						}
						err = s.Save("image", imgData, "", imgPath)
						if err != nil {
							log.Printf("error during saving data to database: %v", err)
						}
						isExceeded, _ := s.IsLimitExceeded()
						if isExceeded {
							err := s.DeleteOldestRecord()
							if err != nil {
								log.Printf("error during deletoing the oldest record: %v", err)
							}
						}
						continue
					}
				}
			}
			err := s.Save("text", []byte(newText), string(newText), "")
			if err != nil {
				log.Printf("error during saving the content: %v", err)
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

		} else if len(newImg) > 0 {
			imgPath, err := getPath()
			if err != nil {
				log.Printf("error during getting image path: %v", err)
			}
			err = s.Save("image", newImg, "", imgPath)
			if err != nil {
				log.Printf("error during saving the content: %v", err)
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

func getPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	picturesDir := filepath.Join(homeDir, ".config", "glipboard", "Pictures")

	err = os.MkdirAll(picturesDir, 0o755)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("img_%d.png", time.Now().UnixNano())

	finalPath := filepath.Join(picturesDir, fileName)
	return finalPath, nil
}

func getEnv() bool {
	if runtime.GOOS != "linux" {
		return false
	}

	sessionType := os.Getenv("XDG_SESSION_TYPE")

	if sessionType == "wayland" {
		return true
	} else if sessionType == "x11" {
		return false
	}

	waylandDisplay := os.Getenv("WAYLAND_DISPLAY")
	if waylandDisplay != "" {
		return true
	}

	display := os.Getenv("DISPLAY")
	if display != "" {
		return false
	}

	return false
}

/* getLocalImagePath checks if the copied text is actually a file path to an image. */
func getLocalImagePath(text string) (string, bool) {
	text = strings.TrimSpace(text)
	var path string

	if strings.HasPrefix(text, "file://") {
		lines := strings.Split(text, "\n")
		firstLine := strings.TrimSpace(lines[0])

		u, err := url.Parse(firstLine)
		if err != nil {
			return "", false
		}
		path = u.Path
	} else if strings.HasPrefix(text, "/") {
		lines := strings.Split(text, "\n")
		path = strings.TrimSpace(lines[0])
	} else {
		return "", false
	}

	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return "", false
	}

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp", ".gif", ".bmp":
		return path, true
	}

	return "", false
}
