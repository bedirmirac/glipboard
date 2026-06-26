$appName = "glipboard"
$appExe = "glipboard.exe"

# --- GITHUB REPOSITORY DETAILS ---
$githubUser = "bedirmirac"
$githubRepo = "glipboard"
# ---------------------------------

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " WARNING: Glipboard Installation" -ForegroundColor Yellow
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "This application will be configured to start automatically"
Write-Host "in the background upon system boot to manage clipboard operations."
Write-Host ""
Write-Host "If you do not want the application to start automatically,"
Write-Host "you can abort the installation now."
Write-Host "==========================================================" -ForegroundColor Cyan

$response = Read-Host "Do you approve the automatic startup and wish to continue? [y/N]"

if ($response -notmatch "^[yY]") {
    Write-Host "Installation aborted." -ForegroundColor Red
    exit
}

# 1. Detect Architecture
$arch = $env:PROCESSOR_ARCHITECTURE.ToLower()
if ($arch -eq "amd64") {
    $osArch = "amd64"
} elseif ($arch -eq "arm64") {
    $osArch = "arm64"
} else {
    Write-Host "ERROR: Unsupported architecture: $arch" -ForegroundColor Red
    exit
}

# 2. Download Binary and Icon
$installDir = "$env:LOCALAPPDATA\$appName"
if (!(Test-Path -Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

$targetExe = Join-Path -Path $installDir -ChildPath $appExe
$downloadUrl = "https://github.com/$githubUser/$githubRepo/releases/latest/download/$appName-windows-$osArch.exe"
$iconUrl = "https://raw.githubusercontent.com/$githubUser/$githubRepo/main/assets/icon.ico"
$iconPath = Join-Path -Path $installDir -ChildPath "icon.ico"

Write-Host "-> Downloading Glipboard binary from GitHub..." -ForegroundColor Cyan

try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $targetExe -ErrorAction Stop
    Write-Host "-> Binary downloaded successfully to: $targetExe" -ForegroundColor Green
} catch {
    Write-Host "ERROR: Failed to download the binary." -ForegroundColor Red
    Write-Host "Please ensure the repository is public and the release assets are named correctly." -ForegroundColor Red
    exit
}

Write-Host "-> Downloading icon..." -ForegroundColor Cyan
try {
    Invoke-WebRequest -Uri $iconUrl -OutFile $iconPath -ErrorAction Stop
} catch {
    Write-Host "Warning: Failed to download icon.ico. Default icon will be used." -ForegroundColor Yellow
}

# 3. Setup system startup and desktop shortcuts
$WshShell = New-Object -comObject WScript.Shell

# Startup Shortcut (Background daemon)
$startupFolder = [Environment]::GetFolderPath('Startup')
$startupShortcutPath = Join-Path -Path $startupFolder -ChildPath "$appName.lnk"

$startupShortcut = $WshShell.CreateShortcut($startupShortcutPath)
$startupShortcut.TargetPath = $targetExe
$startupShortcut.WindowStyle = 7 # 7 = Minimized (Runs silently)
$startupShortcut.Save()
Write-Host "-> Added to Startup tasks." -ForegroundColor Green

# Desktop Shortcut (For TUI)
$desktopFolder = [Environment]::GetFolderPath('Desktop')
$desktopShortcutPath = Join-Path -Path $desktopFolder -ChildPath "$appName TUI.lnk"

$desktopShortcut = $WshShell.CreateShortcut($desktopShortcutPath)
$desktopShortcut.TargetPath = $targetExe
$desktopShortcut.Arguments = "--tui=true"

if (Test-Path -Path $iconPath) {
    $desktopShortcut.IconLocation = $iconPath
}

$desktopShortcut.Save()
Write-Host "-> Desktop shortcut created." -ForegroundColor Green

# 4. Start the daemon immediately
Start-Process -FilePath $targetExe -WindowStyle Hidden

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " Installation Complete! " -ForegroundColor Green
Write-Host " Glipboard is now running in the background."
Write-Host " You can use the 'Glipboard TUI' shortcut on your Desktop for the interface."
Write-Host "==========================================================" -ForegroundColor Cyan
