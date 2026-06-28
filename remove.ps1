$appName = "glipboard"

Write-Host "Removing Glipboard from the system..." -ForegroundColor Cyan

# 1. Safely terminate any running glipboard processes
$process = Get-Process -Name $appName -ErrorAction SilentlyContinue
if ($process) {
    Write-Host "Running glipboard process found, terminating..." -ForegroundColor Yellow
    Stop-Process -Name $appName -Force
    Start-Sleep -Seconds 1
}

# 2. Remove Startup daemon script (VBScript) and Desktop Shortcut
$startupFolder = [Environment]::GetFolderPath('Startup')
$vbsPath = Join-Path -Path $startupFolder -ChildPath "$appName-daemon.vbs"

if (Test-Path $vbsPath) {
    Remove-Item $vbsPath -Force
    Write-Host "Removed Startup daemon script ($vbsPath)" -ForegroundColor Gray
}

$desktopFolder = [Environment]::GetFolderPath('Desktop')
$desktopShortcutPath = Join-Path -Path $desktopFolder -ChildPath "$appName TUI.lnk"

if (Test-Path $desktopShortcutPath) {
    Remove-Item $desktopShortcutPath -Force
    Write-Host "Removed Desktop shortcut ($desktopShortcutPath)" -ForegroundColor Gray
}

# 3. Remove binaries and database files from the local directory (Eğer çalıştırılan klasörde kalıntı varsa)
$targetFiles = @(
    ".\glipboard.exe",
    ".\glipboard-windows-amd64.exe",
    ".\glipboard.db",
    ".\glipboard.db-wal",
    ".\glipboard.db-shm"
)

foreach ($file in $targetFiles) {
    if (Test-Path $file) {
        Remove-Item $file -Force
        Write-Host "Removed $file." -ForegroundColor Gray
    }
}

# 4. Remove AppData configuration folders (Ana kurulum yerleri)
$appDataPath = Join-Path $env:APPDATA $appName
if (Test-Path $appDataPath) {
    Remove-Item $appDataPath -Recurse -Force
    Write-Host "Application data folder removed ($appDataPath)" -ForegroundColor Gray
}

$localAppDataPath = Join-Path $env:LOCALAPPDATA $appName
if (Test-Path $localAppDataPath) {
    Remove-Item $localAppDataPath -Recurse -Force
    Write-Host "Local application data folder removed ($localAppDataPath)" -ForegroundColor Gray
}

Write-Host "Cleanup complete! All glipboard components have been removed." -ForegroundColor Green
