#!/usr/bin/env bash

set -e

APP_NAME="glipboard"
BIN_DIR="$HOME/.local/bin"
EXECUTABLE_PATH="$BIN_DIR/$APP_NAME"

# --- GITHUB REPOSITORY DETAILS ---
GITHUB_USER="bedirmirac"
GITHUB_REPO="glipboard"
# ---------------------------------

echo "=========================================================="
echo " WARNING: Glipboard Installation"
echo "=========================================================="
echo "This application will be configured to start automatically"
echo "in the background (as a Daemon/Service) upon system boot"
echo "to manage clipboard operations."
echo ""
echo "If you do not want the application to start automatically,"
echo "you can abort the installation now."
echo "=========================================================="
read -p "Do you approve the automatic startup and wish to continue? [y/N]: " response

if [[ ! "$response" =~ ^[yY](es)?$ ]]; then
  echo "Installation aborted."
  exit 0
fi

# 1. Detect OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
Linux*) OS_NAME="linux" ;;
Darwin*) OS_NAME="darwin" ;;
*)
  echo "ERROR: Unsupported OS: $OS"
  exit 1
  ;;
esac

case "$ARCH" in
x86_64) OS_ARCH="amd64" ;;
aarch64) OS_ARCH="arm64" ;;
arm64) OS_ARCH="arm64" ;;
*)
  echo "ERROR: Unsupported architecture: $ARCH"
  exit 1
  ;;
esac

# 2. Download Binary and Icons
DOWNLOAD_URL="https://github.com/${GITHUB_USER}/${GITHUB_REPO}/releases/latest/download/${APP_NAME}-${OS_NAME}-${OS_ARCH}"
ICON_URL_BASE="https://raw.githubusercontent.com/${GITHUB_USER}/${GITHUB_REPO}/main/assets"

echo "-> Downloading Glipboard binary from GitHub..."
mkdir -p "$BIN_DIR"

if curl -sL --fail "$DOWNLOAD_URL" -o "$EXECUTABLE_PATH"; then
  chmod +x "$EXECUTABLE_PATH"
  echo "-> Binary downloaded successfully to: $EXECUTABLE_PATH"
else
  echo "ERROR: Failed to download the binary."
  echo "Please ensure the repository is public and the release assets are named correctly."
  exit 1
fi

# 3. Setup System Daemons and Desktop Entries
if [ "$OS" = "Linux" ]; then
  echo "-> Linux detected. Setting up Systemd service and .desktop entry..."

  # Download and set up the icon
  ICON_DIR="$HOME/.local/share/icons"
  mkdir -p "$ICON_DIR"
  curl -sL --fail "$ICON_URL_BASE/icon.png" -o "$ICON_DIR/$APP_NAME.png" || echo "Warning: Failed to download icon.png"

  # Systemd Service (Background daemon)
  SYSTEMD_DIR="$HOME/.config/systemd/user"
  mkdir -p "$SYSTEMD_DIR"

  cat >"$SYSTEMD_DIR/$APP_NAME.service" <<EOF
[Unit]
Description=Glipboard Daemon
After=network.target

[Service]
ExecStart=$EXECUTABLE_PATH
Restart=on-failure

[Install]
WantedBy=default.target
EOF

  systemctl --user daemon-reload
  systemctl --user enable --now "$APP_NAME.service"
  echo "-> Systemd service enabled and started."

  # Desktop Entry (For TUI)
  DESKTOP_DIR="$HOME/.local/share/applications"
  mkdir -p "$DESKTOP_DIR"

  cat >"$DESKTOP_DIR/$APP_NAME.desktop" <<EOF
[Desktop Entry]
Name=Glipboard TUI
Comment=Terminal UI for Glipboard
Exec=$EXECUTABLE_PATH --tui=true
Icon=$APP_NAME
Terminal=true
Type=Application
Categories=Utility;
EOF
  echo "-> Desktop application (.desktop) created."

elif [ "$OS" = "Darwin" ]; then
  echo "-> macOS detected. Setting up LaunchAgent and .app..."

  # LaunchAgent (Background daemon)
  PLIST_DIR="$HOME/Library/LaunchAgents"
  mkdir -p "$PLIST_DIR"
  PLIST_FILE="$PLIST_DIR/com.user.$APP_NAME.plist"

  cat >"$PLIST_FILE" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.user.$APP_NAME</string>
    <key>ProgramArguments</key>
    <array>
        <string>$EXECUTABLE_PATH</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF

  launchctl load "$PLIST_FILE"
  echo "-> LaunchAgent enabled and started."

  # macOS Desktop App (TUI)
  APP_DIR="/Applications/$APP_NAME TUI.app"
  osacompile -e "tell application \"Terminal\" to do script \"$EXECUTABLE_PATH --tui=true\"" -o "$APP_DIR"

  # Replace default applet icon with custom .icns
  APP_ICON_DIR="$APP_DIR/Contents/Resources"
  if curl -sL --fail "$ICON_URL_BASE/icon.icns" -o "$APP_ICON_DIR/applet.icns"; then
    touch "$APP_DIR" # Refresh macOS icon cache
  else
    echo "Warning: Failed to download icon.icns"
  fi

  echo "-> App created at '$APP_DIR'."
fi

echo "=========================================================="
echo " Installation Complete! "
echo " Glipboard is now running in the background."
echo " You can open 'Glipboard TUI' from your application menu."
echo "=========================================================="
