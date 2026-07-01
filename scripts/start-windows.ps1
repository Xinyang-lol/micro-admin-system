$ErrorActionPreference = "Stop"

$Root = Split-Path -Parent $PSScriptRoot
$Backend = Join-Path $Root "backend"
$Frontend = Join-Path $Root "frontend"

function Test-CommandExists {
    param(
        [string]$Name,
        [string]$InstallTip
    )

    if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
        Write-Host ""
        Write-Host "Missing command: $Name" -ForegroundColor Red
        Write-Host $InstallTip -ForegroundColor Yellow
        return $false
    }

    return $true
}

Write-Host "========================================"
Write-Host " Micro Admin System - Windows Launcher"
Write-Host "========================================"
Write-Host ""

$ok = $true
$ok = (Test-CommandExists "docker" "Install and start Docker Desktop first: https://www.docker.com/products/docker-desktop/") -and $ok
$ok = (Test-CommandExists "go" "Install Go 1.22 or later first: https://go.dev/dl/") -and $ok
$ok = (Test-CommandExists "node" "Install Node.js 20 or later first: https://nodejs.org/") -and $ok
$npmCommand = if (Get-Command "npm.cmd" -ErrorAction SilentlyContinue) { "npm.cmd" } else { "npm" }
$ok = (Test-CommandExists $npmCommand "Install Node.js first, then reopen PowerShell.") -and $ok

if (-not $ok) {
    Write-Host ""
    Write-Host "The environment is not ready. Install the missing tools, then run start-windows.bat again." -ForegroundColor Yellow
    exit 1
}

Write-Host "1/4 Starting MySQL, Redis and Consul..."
Push-Location $Backend
try {
    docker compose up -d
}
catch {
    Write-Host ""
    Write-Host "Docker failed to start services. Make sure Docker Desktop is running." -ForegroundColor Red
    throw
}
finally {
    Pop-Location
}

Write-Host "2/4 Opening backend service windows..."
$services = @(
    @{ Title = "user-service"; Command = "go run ./user-service" },
    @{ Title = "device-service"; Command = "go run ./device-service" },
    @{ Title = "file-service"; Command = "go run ./file-service" },
    @{ Title = "api-gateway"; Command = "go run ./api-gateway" }
)

foreach ($service in $services) {
    $cmd = "cd `"$Backend`"; `$Host.UI.RawUI.WindowTitle = `"$($service.Title)`"; $($service.Command)"
    Start-Process powershell.exe -ArgumentList @("-NoExit", "-NoProfile", "-Command", $cmd)
    Start-Sleep -Seconds 1
}

Write-Host "3/4 Installing frontend dependencies..."
Push-Location $Frontend
try {
    if (-not (Test-Path (Join-Path $Frontend "node_modules"))) {
        & $npmCommand install
    }
    else {
        Write-Host "frontend/node_modules already exists. Skipping npm install."
    }
}
finally {
    Pop-Location
}

Write-Host "4/4 Opening frontend service window..."
$frontendCmd = "cd `"$Frontend`"; `$Host.UI.RawUI.WindowTitle = `"frontend-vite`"; $npmCommand run dev"
Start-Process powershell.exe -ArgumentList @("-NoExit", "-NoProfile", "-Command", $frontendCmd)

Write-Host ""
Write-Host "Waiting for frontend..."
Start-Sleep -Seconds 8
Start-Process "http://127.0.0.1:5173"

Write-Host ""
Write-Host "Done. Open: http://127.0.0.1:5173" -ForegroundColor Green
Write-Host "Username: admin" -ForegroundColor Green
Write-Host "Password: admin123" -ForegroundColor Green
Write-Host ""
