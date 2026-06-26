#!/bin/bash

echo "Removing Glipboard from the system..."

# 1. Terminate any running glipboard processes in the background
pkill -f glipboard 2>/dev/null

# 2. Remove binaries and database files from the local directory
rm -f ./glipboard
rm -f ./glipboard-linux-amd64
rm -f ./glipboard.db
rm -f ./glipboard.db-wal
rm -f ./glipboard.db-shm

# 3. Remove potential global installation and configuration directories
rm -rf "$HOME/.config/glipboard"
rm -rf "$HOME/.local/share/glipboard"
rm -f "$HOME/.local/bin/glipboard"
rm -f "/usr/local/bin/glipboard"

echo "Cleanup complete! All glipboard components have been removed."
