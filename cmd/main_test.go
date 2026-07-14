package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetLocalImagePath(t *testing.T) {
	tempDir := t.TempDir()

	validImgPath := filepath.Join(tempDir, "test_image.png")
	invalidExtPath := filepath.Join(tempDir, "test_doc.txt")

	err := os.WriteFile(validImgPath, []byte("fake image data"), 0o644)
	if err != nil {
		t.Fatalf("Could not create temporary test file: %v", err)
	}
	err = os.WriteFile(invalidExtPath, []byte("fake text data"), 0o644)
	if err != nil {
		t.Fatalf("Could not create temporary test file: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		wantPath string
		wantOk   bool
	}{
		{
			name:     "Valid image with file:// format",
			input:    "file://" + validImgPath + "\nextra line",
			wantPath: validImgPath,
			wantOk:   true,
		},
		{
			name:     "Valid image with absolute path",
			input:    validImgPath + "\n",
			wantPath: validImgPath,
			wantOk:   true,
		},
		{
			name:     "File with invalid extension",
			input:    invalidExtPath,
			wantPath: "",
			wantOk:   false,
		},
		{
			name:     "Non-existent file",
			input:    filepath.Join(tempDir, "does_not_exist.png"),
			wantPath: "",
			wantOk:   false,
		},
		{
			name:     "Directory path (not an image)",
			input:    tempDir,
			wantPath: "",
			wantOk:   false,
		},
		{
			name:     "Random text",
			input:    "hello world this is not an image",
			wantPath: "",
			wantOk:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotOk := getLocalImagePath(tt.input)
			if gotOk != tt.wantOk {
				t.Errorf("getLocalImagePath() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotPath != tt.wantPath {
				t.Errorf("getLocalImagePath() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("getEnv test is only intended for Linux environments.")
	}

	tests := []struct {
		name        string
		setupEnvs   func(t *testing.T)
		wantWayland bool
	}{
		{
			name: "Wayland Session Type",
			setupEnvs: func(t *testing.T) {
				t.Setenv("XDG_SESSION_TYPE", "wayland")
			},
			wantWayland: true,
		},
		{
			name: "X11 Session Type",
			setupEnvs: func(t *testing.T) {
				t.Setenv("XDG_SESSION_TYPE", "x11")
			},
			wantWayland: false,
		},
		{
			name: "WAYLAND_DISPLAY set",
			setupEnvs: func(t *testing.T) {
				t.Setenv("XDG_SESSION_TYPE", "")
				t.Setenv("WAYLAND_DISPLAY", "wayland-0")
			},
			wantWayland: true,
		},
		{
			name: "DISPLAY set only",
			setupEnvs: func(t *testing.T) {
				t.Setenv("XDG_SESSION_TYPE", "")
				t.Setenv("WAYLAND_DISPLAY", "")
				t.Setenv("DISPLAY", ":0")
			},
			wantWayland: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnvs(t)
			got := getEnv()
			if got != tt.wantWayland {
				t.Errorf("getEnv() = %v, want %v", got, tt.wantWayland)
			}
		})
	}
}

func TestGetPath(t *testing.T) {
	path, err := getPath()
	if err != nil {
		t.Fatalf("getPath() returned an unexpected error: %v", err)
	}

	if !strings.Contains(path, filepath.Join(".config", "glipboard", "Pictures", "img_")) {
		t.Errorf("getPath() generated an unexpected file path: %v", path)
	}

	if !strings.HasSuffix(path, ".png") {
		t.Errorf("getPath() generated file should have a .png extension: %v", path)
	}

	dir := filepath.Dir(path)
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		t.Errorf("getPath() could not create the required directories: %v", dir)
	}
}

func TestSetupLogger(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)
	logFile := setupLogger()
	if logFile == nil {
		t.Fatal("setupLogger() returned nil, expected *os.File.")
	}

	defer logFile.Close()

	expectedPath := filepath.Join(tempHome, ".config", "glipboard", ".glipboard.log")
	info, err := os.Stat(expectedPath)
	if err != nil {
		t.Errorf("setupLogger() could not create the log file: %v", err)
	}

	if info.IsDir() {
		t.Errorf("The expected log file should not be a directory.")
	}
}
