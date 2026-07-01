@echo off
chcp 65001 >nul
echo This script enables Windows features required by Docker Desktop.
echo Please allow the administrator permission prompt.
echo.
powershell -NoProfile -ExecutionPolicy Bypass -Command "Start-Process powershell.exe -Verb RunAs -ArgumentList @('-NoExit','-ExecutionPolicy','Bypass','-Command','dism.exe /online /enable-feature /featurename:Microsoft-Windows-Subsystem-Linux /all /norestart; dism.exe /online /enable-feature /featurename:VirtualMachinePlatform /all /norestart; wsl --update; wsl --set-default-version 2; Write-Host \"If DISM says restart required, restart Windows, then open Docker Desktop.\" -ForegroundColor Green; pause')"
echo.
echo If Windows asks you to restart, restart first.
pause
