$appName = "glipboard"

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " WARNING: Uninstalling Glipboard..." -ForegroundColor Yellow
Write-Host "==========================================================" -ForegroundColor Cyan

# 1. Çalışan glipboard süreçlerini güvenle sonlandır
$process = Get-Process -Name $appName -ErrorAction SilentlyContinue
if ($process) {
    Write-Host "-> Running glipboard process found, terminating..." -ForegroundColor Yellow
    Stop-Process -Name $appName -Force
    Start-Sleep -Seconds 1
    Write-Host "-> Process terminated." -ForegroundColor Green
}

# 2. Başlangıç (Startup) VBScript dosyasını kaldır
$startupFolder = [Environment]::GetFolderPath('Startup')
$vbsPath = Join-Path -Path $startupFolder -ChildPath "$appName-daemon.vbs"

if (Test-Path $vbsPath) {
    Remove-Item $vbsPath -Force
    Write-Host "-> Removed Startup daemon script ($vbsPath)" -ForegroundColor Green
}

# 3. Masaüstü Kısayolunu kaldır (TUI)
$desktopFolder = [Environment]::GetFolderPath('Desktop')
$desktopShortcutPath = Join-Path -Path $desktopFolder -ChildPath "$appName TUI.lnk"

if (Test-Path $desktopShortcutPath) {
    Remove-Item $desktopShortcutPath -Force
    Write-Host "-> Removed Desktop shortcut ($desktopShortcutPath)" -ForegroundColor Green
}

# 4. .config/glipboard klasörünü ve içindekileri tamamen kaldır
$userConfigPath = Join-Path $HOME ".config\$appName"
if (Test-Path $userConfigPath) {
    Remove-Item $userConfigPath -Recurse -Force
    Write-Host "-> User config folder removed ($userConfigPath)" -ForegroundColor Green
}

# 5. Kurulum ve Yapılandırma Klasörlerini (AppData / LocalAppData) kaldır
$localAppDataPath = Join-Path $env:LOCALAPPDATA $appName
if (Test-Path $localAppDataPath) {
    Remove-Item $localAppDataPath -Recurse -Force
    Write-Host "-> Application installation folder removed ($localAppDataPath)" -ForegroundColor Green
}

$appDataPath = Join-Path $env:APPDATA $appName
if (Test-Path $appDataPath) {
    Remove-Item $appDataPath -Recurse -Force
    Write-Host "-> Application data folder removed ($appDataPath)" -ForegroundColor Green
}

# 6. Yerel dizindeki olası kalıntıları temizle
$targetFiles = @(
    ".\glipboard.exe",
    ".\glipboard-windows-amd64.exe",
    ".\glipboard-windows-arm64.exe",
    ".\icon.ico",
    ".\glipboard.db",
    ".\glipboard.db-wal",
    ".\glipboard.db-shm"
)

foreach ($file in $targetFiles) {
    if (Test-Path $file) {
        Remove-Item $file -Force
        Write-Host "-> Removed local artifact: $file" -ForegroundColor Gray
    }
}

Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host " Cleanup complete! All glipboard components have been removed." -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Cyan
