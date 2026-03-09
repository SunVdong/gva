@echo off
set GOOS=linux
set GOARCH=amd64
cd /d "%~dp0server"
go build -o "..\app-linux" .
cd /d "%~dp0"
echo build done: app-linux
pause
