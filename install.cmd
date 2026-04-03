@echo off
setlocal enabledelayedexpansion

set GITHUB_OWNER=neetozone
set GITHUB_REPO=neeto-cal-cli
set BINARY_NAME=neetocal
set INSTALL_DIR=%LOCALAPPDATA%\Programs\neetocal

:: Detect architecture
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set ARCH=arm64
) else (
    set ARCH=amd64
)

:: Get latest version
echo Fetching latest version...
set TMP_DIR=%TEMP%\neetocal-install
if exist "%TMP_DIR%" rmdir /s /q "%TMP_DIR%"
mkdir "%TMP_DIR%"

curl -fsSL "https://api.github.com/repos/%GITHUB_OWNER%/%GITHUB_REPO%/releases/latest" -o "%TMP_DIR%\release.json" 2>nul
if errorlevel 1 (
    echo Error: could not fetch latest version. Check your internet connection.
    exit /b 1
)

for /f "tokens=2 delims=:," %%a in ('findstr "tag_name" "%TMP_DIR%\release.json"') do (
    set VERSION=%%~a
    set VERSION=!VERSION: =!
    set VERSION=!VERSION:"=!
)

if not defined VERSION (
    echo Error: could not determine latest version.
    exit /b 1
)

set VERSION_NUM=!VERSION:v=!

:: Download
set FILENAME=%GITHUB_REPO%_%VERSION_NUM%_windows_%ARCH%.zip
set URL=https://github.com/%GITHUB_OWNER%/%GITHUB_REPO%/releases/download/%VERSION%/%FILENAME%

echo Downloading %BINARY_NAME% %VERSION% for Windows/%ARCH%...
curl -fsSL "%URL%" -o "%TMP_DIR%\%FILENAME%" 2>nul
if errorlevel 1 (
    echo Error: download failed.
    exit /b 1
)

:: Extract
tar -xf "%TMP_DIR%\%FILENAME%" -C "%TMP_DIR%" 2>nul
if errorlevel 1 (
    echo Error: extraction failed.
    exit /b 1
)

:: Install
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
move /y "%TMP_DIR%\%BINARY_NAME%.exe" "%INSTALL_DIR%\%BINARY_NAME%.exe" >nul

:: Add to user PATH if not already present
set "USER_PATH="
for /f "tokens=2*" %%a in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "USER_PATH=%%b"

echo ;!USER_PATH!; | findstr /i /c:";%INSTALL_DIR%;" >nul 2>nul
if errorlevel 1 (
    if defined USER_PATH (
        setx PATH "!USER_PATH!;%INSTALL_DIR%" >nul 2>&1
    ) else (
        setx PATH "%INSTALL_DIR%" >nul 2>&1
    )
    echo Added %INSTALL_DIR% to your PATH.
)

:: Cleanup
rmdir /s /q "%TMP_DIR%"

echo.
echo %BINARY_NAME% %VERSION% installed successfully.
echo Restart your terminal, then run '%BINARY_NAME% --help' to get started.
