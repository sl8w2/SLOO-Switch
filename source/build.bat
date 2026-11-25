@echo off
setlocal
cd /d "%~dp0"
set "GOEXE=go"
where go >nul 2>nul && goto build
if exist "%ProgramFiles%\Go\bin\go.exe" set "GOEXE=%ProgramFiles%\Go\bin\go.exe" & goto build
echo Go is not installed. Download it from https://go.dev/dl/ and run this file again.
pause
exit /b 1
:build
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
"%GOEXE%" build -trimpath -ldflags "-s -w -H=windowsgui" -o "%~dp0..\SLOO-Switch.exe" .
if errorlevel 1 goto failed
echo.
echo Done. SLOO-Switch.exe was created next to the source folder. One file, no extras.
pause
exit /b 0
:failed
echo.
echo Build failed.
pause
exit /b 1
