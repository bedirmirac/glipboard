Write-Host "Removing Glipboard from the system..." -ForegroundColor Cyan

# 1. Safely terminate any running glipboard processes
$process = Get-Process -Name "glipboard" -ErrorAction SilentlyContinue
if ($process) {
    Write-Host "Running glipboard process found, terminating..." -ForegroundColor Yellow
    Stop-Process -Name "glipboard" -Force
    Start-Sleep -Seconds 1
}

# 2. Remove binaries and database files from the local directory
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

# 3. Remove potential AppData configuration folders
$appDataPath = Join-Path $env:APPDATA "glipboard"
if (Test-Path $appDataPath) {
    Remove-Item $appDataPath -Recurse -Force
    Write-Host "Application data folder removed ($appDataPath)" -ForegroundColor Gray
}

$localAppDataPath = Join-Path $env:LOCALAPPDATA "glipboard"
if (Test-Path $localAppDataPath) {
    Remove-Item $localAppDataPath -Recurse -Force
    Write-Host "Local application data folder removed ($localAppDataPath)" -ForegroundColor Gray
}

Write-Host "Cleanup complete! All glipboard components have been removed." -ForegroundColor Green
