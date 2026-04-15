@echo off
setlocal

set "BASE_URL=https://neeto-downloads.s3.amazonaws.com/cli/NeetoCal/latest"
set "INSTALL_DIR=%LOCALAPPDATA%\Programs\neetocal"
set "ARCHIVE=neetocal_windows_amd64.zip"

echo Downloading NeetoCal CLI...
set "TMPDIR=%TEMP%\neetocal-install"
if exist "%TMPDIR%" rmdir /s /q "%TMPDIR%"
mkdir "%TMPDIR%"

curl -fsSL "%BASE_URL%/%ARCHIVE%" -o "%TMPDIR%\%ARCHIVE%"
if %errorlevel% neq 0 (
    echo Failed to download NeetoCal CLI.
    exit /b 1
)

echo Extracting...
powershell -Command "Expand-Archive -Path '%TMPDIR%\%ARCHIVE%' -DestinationPath '%TMPDIR%' -Force"

echo Installing to %INSTALL_DIR%...
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
copy /y "%TMPDIR%\neetocal.exe" "%INSTALL_DIR%\neetocal.exe" >nul

:: Add to user PATH if not already present
echo %PATH% | findstr /i /c:"%INSTALL_DIR%" >nul
if %errorlevel% neq 0 (
    for /f "tokens=2*" %%A in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "USER_PATH=%%B"
    setx PATH "%USER_PATH%;%INSTALL_DIR%" >nul
    echo Added %INSTALL_DIR% to user PATH.
)

rmdir /s /q "%TMPDIR%"

echo NeetoCal CLI installed successfully. Restart your terminal and run 'neetocal --help' to get started.
