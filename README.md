# Glipboard

Glipboard is a seamless, lightweight, and cross-platform clipboard management tool written in Go. Designed for efficiency and minimal resource consumption, it operates silently as a system daemon to handle your clipboard operations in the background. 

When you need to interact with your clipboard history, Glipboard provides a fast, keyboard-centric Terminal User Interface (TUI) that can be triggered instantly. Perfect for standalone window managers and minimal desktop environments.

## Features
* **Cross-Platform:** Native background service integration for Linux (Systemd), macOS (LaunchAgent), and Windows (Startup).
* **Silent Daemon:** Runs invisibly in the background on system boot without unnecessary bloat.
* **Terminal UI (TUI):** A clean, interactive terminal interface built for power users.
* **Local Storage:** Stores your clipboard history locally using SQLite.
* **Standalone Binary:** Compiled as a single, static executable with zero external dependencies.

## Installation

### Linux & macOS
Run the following script to automatically download the latest binary, configure the background daemon, and set up the desktop entry:
\`\`\`bash
curl -sL https://raw.githubusercontent.com/bedirmirac/glipboard/main/install.sh | bash
\`\`\`

### Windows
Open PowerShell and run the following command to install the background service and create a desktop shortcut:
\`\`\`powershell
Invoke-WebRequest -Uri https://raw.githubusercontent.com/bedirmirac/glipboard/main/install.ps1 -OutFile install.ps1; .\install.ps1
\`\`\`

## Usage

Glipboard automatically starts its daemon in the background upon system login. 

To open the interface and browse your clipboard history, run:
\`\`\`bash
glipboard --tui=true
\`\`\`
*Pro Tip: Bind this command to a custom keyboard shortcut in your window manager (e.g., Super+V) for instant access.*
Or if you installed glipborad using given scripts you can open the tui with desktop app created

## Acknowledgements

This project would not be possible without the incredible open-source work from the following repositories:
* [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) - For the powerful and elegant TUI framework.
* [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - For the robust SQLite driver powering the local history storage.
* [atotto/clipboard](https://github.com/atotto/clipboard) - For the seamless cross-platform clipboard interactions.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
