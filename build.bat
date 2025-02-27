@echo off
echo Building gitignore CLI tool...
go build -o gitignore.exe main.go
if %ERRORLEVEL% NEQ 0 (
    echo Build failed!
    exit /b %ERRORLEVEL%
)
echo Build successful! You can now use gitignore.exe
