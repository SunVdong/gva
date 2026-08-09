@echo off
setlocal

set GOOS=linux
set GOARCH=amd64

echo [1/2] building backend (linux/amd64)...
cd /d "%~dp0server"
go build -o "..\app-linux" .
if errorlevel 1 (
  echo backend build failed
  pause
  exit /b 1
)
cd /d "%~dp0"
echo backend done: app-linux

echo [2/2] building frontend...
cd /d "%~dp0web"
where npm >nul 2>nul
if %errorlevel%==0 (
  call npm run build
) else (
  yarn build
)
if errorlevel 1 (
  echo frontend build failed
  pause
  exit /b 1
)
cd /d "%~dp0"
echo frontend done: web\dist

echo.
echo build done:
echo   - app-linux
echo   - web\dist
pause
