#!/usr/bin/env bash

APP_NAME="glipboard"
BIN_DIR="$HOME/.local/bin"
EXECUTABLE_PATH="$BIN_DIR/$APP_NAME"

echo "=========================================================="
echo " Removing Glipboard from the system..."
echo "=========================================================="

OS="$(uname -s)"

if [ "$OS" = "Linux" ]; then
    echo "-> Stopping and removing Systemd service..."
    systemctl --user stop "$APP_NAME.service" 2>/dev/null || true
    systemctl --user disable "$APP_NAME.service" 2>/dev/null || true
    rm -f "$HOME/.config/systemd/user/$APP_NAME.service"
    systemctl --user daemon-reload
elif [ "$OS" = "Darwin" ]; then
    echo "-> Stopping and removing macOS LaunchAgent..."
    PLIST_FILE="$HOME/Library/LaunchAgents/com.user.$APP_NAME.plist"
    launchctl unload "$PLIST_FILE" 2>/dev/null || true
    rm -f "$PLIST_FILE"
fi

echo "-> Terminating running $APP_NAME processes..."
pkill -x "$APP_NAME" 2>/dev/null || true

echo "-> Removing desktop entries and icons..."
if [ "$OS" = "Linux" ]; then
    rm -f "$HOME/.local/share/applications/$APP_NAME.desktop"
    rm -f "$HOME/.local/share/icons/$APP_NAME.png"
elif [ "$OS" = "Darwin" ]; then
    rm -rf "/Applications/$APP_NAME TUI.app"
fi

echo "-> Removing executable files..."
rm -f "$EXECUTABLE_PATH"
rm -f "/usr/local/bin/$APP_NAME"

echo "-> Removing configuration and database directories..."
rm -rf "$HOME/.config/$APP_NAME"
rm -rf "$HOME/.local/share/$APP_NAME"

echo "-> Cleaning up local repository remnants..."
rm -rf "./$APP_NAME" "./$APP_NAME-linux-amd64" "./$APP_NAME-darwin-arm64" "./$APP_NAME-darwin-amd64"
rm -f ./glipboard.db ./glipboard.db-wal ./glipboard.db-shm

echo "=========================================================="
echo " Cleanup complete! All Glipboard components are removed."
echo "=========================================================="
