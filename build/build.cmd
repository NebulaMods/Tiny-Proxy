@echo off
setlocal enabledelayedexpansion

:menu
echo Choose the OS you want to build for:
echo 1. Windows
echo 2. Linux
echo 3. Darwin (macOS)
echo 4. All
set /p choice="Enter your choice (1-4): "

if "%choice%"=="1" goto build_windows
if "%choice%"=="2" goto build_linux
if "%choice%"=="3" goto build_darwin
if "%choice%"=="4" goto build_all

echo Invalid choice.
goto menu

:build_windows
set GOOS=windows
set GOARCH=amd64
set FILETYPE=.exe
goto build

:build_linux
set GOOS=linux
set GOARCH=amd64
set FILETYPE=
goto build

:build_darwin
set GOOS=darwin
set GOARCH=amd64
set FILETYPE=
goto build

:build_all
call :build_windows
call :build_linux
call :build_darwin
goto end

:build
echo Building for %GOOS%...
set CGO_ENABLED=0
set GOGC=100
go build -o ..\bin\%GOOS%\tiny-proxy%FILETYPE% -ldflags "-s -w" ..\src\cmd\main.go
if errorlevel 1 (
    echo Build failed for %GOOS%!
    goto end
)
echo Build succeeded for %GOOS%!
goto end

:end
echo Build process completed.
pause
