$ErrorActionPreference = "Stop"

$GithubOwner = "neetozone"
$GithubRepo = "neeto-cal-cli"
$BinaryName = "neetocal"
$InstallDir = "$env:LOCALAPPDATA\Programs\neetocal"

# Detect architecture
$Arch = if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }

# Get latest version
Write-Host "Fetching latest version..."
$Release = Invoke-RestMethod "https://api.github.com/repos/$GithubOwner/$GithubRepo/releases/latest"
$Version = $Release.tag_name
$VersionNum = $Version.TrimStart("v")

# Download
$Filename = "${GithubRepo}_${VersionNum}_windows_${Arch}.zip"
$Url = "https://github.com/$GithubOwner/$GithubRepo/releases/download/$Version/$Filename"

Write-Host "Downloading $BinaryName $Version for Windows/$Arch..."
$TmpDir = Join-Path $env:TEMP "neetocal-install-$(Get-Random)"
New-Item -ItemType Directory -Path $TmpDir -Force | Out-Null
$ZipPath = Join-Path $TmpDir $Filename

Invoke-WebRequest -Uri $Url -OutFile $ZipPath

# Extract
Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

# Install
New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
Move-Item -Path (Join-Path $TmpDir "$BinaryName.exe") -Destination (Join-Path $InstallDir "$BinaryName.exe") -Force

# Add to user PATH if not already present
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    if ($UserPath) {
        [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    } else {
        [Environment]::SetEnvironmentVariable("Path", "$InstallDir", "User")
    }
    $env:Path = "$env:Path;$InstallDir"
    Write-Host "Added $InstallDir to your PATH."
}

# Cleanup
Remove-Item -Recurse -Force $TmpDir

Write-Host ""
Write-Host "$BinaryName $Version installed successfully."
Write-Host "Restart your terminal, then run '$BinaryName --help' to get started."
