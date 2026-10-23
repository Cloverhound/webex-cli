# Webex CLI installer for Windows
# Usage: irm https://raw.githubusercontent.com/Cloverhound/webex-cli/main/install.ps1 | iex
#   Pin a version: $env:WEBEX_CLI_VERSION = "0.16.0" before running

$ErrorActionPreference = "Stop"

$Repo = "Cloverhound/webex-cli"
$Binary = "webex.exe"
$InstallDir = if ($env:INSTALL_DIR) { $env:INSTALL_DIR } else { "$env:LOCALAPPDATA\webex-cli" }
$Releases = "https://github.com/$Repo/releases"
$GoProxyLatest = "https://proxy.golang.org/github.com/!cloverhound/webex-cli/@latest"
$ApiLatest = "https://api.github.com/repos/$Repo/releases/latest"

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "Unsupported: 32-bit Windows is not supported"
    return
}

function Get-Text($Url, $Headers = @{}) {
    $resp = Invoke-WebRequest -Uri $Url -Headers $Headers -UseBasicParsing
    if ($resp.Content -is [byte[]]) {
        return [System.Text.Encoding]::UTF8.GetString($resp.Content)
    }
    return [string]$resp.Content
}

# The lookup avoids api.github.com where possible: sandboxes and shared CI
# hosts often block it or exhaust its unauthenticated rate limit.
function Get-LatestVersion {
    try {
        $sums = Get-Text "$Releases/latest/download/checksums.txt"
        $m = [regex]::Match($sums, "webex-cli_(\S+)_windows_${Arch}\.zip")
        if ($m.Success) { return $m.Groups[1].Value }
    } catch {}
    try {
        $info = (Get-Text $GoProxyLatest) | ConvertFrom-Json
        if ($info.Version) { return $info.Version -replace '^v', '' }
    } catch {}
    try {
        $headers = @{}
        if ($env:GITHUB_TOKEN) { $headers["Authorization"] = "Bearer $env:GITHUB_TOKEN" }
        $release = (Get-Text $ApiLatest $headers) | ConvertFrom-Json
        if ($release.tag_name) { return $release.tag_name -replace '^v', '' }
    } catch {}
    return $null
}

if ($env:WEBEX_CLI_VERSION) {
    $Version = $env:WEBEX_CLI_VERSION -replace '^v', ''
    Write-Host "Using pinned version: v$Version"
} else {
    Write-Host "Fetching latest release..."
    $Version = Get-LatestVersion
    if (-not $Version) {
        Write-Host "Could not determine the latest version. Tried:"
        Write-Host "  $Releases/latest/download/checksums.txt"
        Write-Host "  $GoProxyLatest"
        Write-Host "  $ApiLatest"
        Write-Host "Pin a version instead:"
        Write-Host "  `$env:WEBEX_CLI_VERSION = `"x.y.z`"; irm https://raw.githubusercontent.com/$Repo/main/install.ps1 | iex"
        Write-Error "Could not determine latest version"
        return
    }
    Write-Host "Latest version: v$Version"
}

# Download
$ZipName = "webex-cli_${Version}_windows_${Arch}.zip"
$Url = "$Releases/download/v$Version/$ZipName"

$TmpDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "webex-cli-install-$(Get-Random)")

try {
    Write-Host "Downloading $Url..."
    $ZipPath = Join-Path $TmpDir $ZipName
    Invoke-WebRequest -Uri $Url -OutFile $ZipPath -UseBasicParsing

    # Verify
    $Sums = Get-Text "$Releases/download/v$Version/checksums.txt"
    $Expected = $null
    foreach ($line in $Sums -split "`n") {
        $fields = $line.Trim() -split '\s+'
        if ($fields.Count -eq 2 -and $fields[1].TrimStart('*') -eq $ZipName) {
            $Expected = $fields[0].ToLower()
            break
        }
    }
    if (-not $Expected) {
        Write-Error "$ZipName is not listed in checksums.txt; refusing to install"
        return
    }
    $Actual = (Get-FileHash -Algorithm SHA256 -Path $ZipPath).Hash.ToLower()
    if ($Actual -ne $Expected) {
        Write-Error "Checksum mismatch for ${ZipName} (expected $Expected, got $Actual); refusing to install"
        return
    }
    Write-Host "Checksum verified."

    # Extract
    Expand-Archive -Path $ZipPath -DestinationPath $TmpDir -Force

    # Install
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Move-Item -Path (Join-Path $TmpDir $Binary) -Destination (Join-Path $InstallDir $Binary) -Force

    # Add to user PATH if not already there
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("Path", "$InstallDir;$UserPath", "User")
        $env:Path = "$InstallDir;$env:Path"
        Write-Host "Added $InstallDir to your PATH."
    }

    Write-Host ""
    Write-Host "Installed webex v$Version to $InstallDir\$Binary"
    Write-Host ""
    Write-Host "NOTE: Restart your terminal for PATH changes to take effect."
    Write-Host ""
    Write-Host "Get started:"
    Write-Host "  webex config set client-id <your-client-id>      # if not using built-in defaults"
    Write-Host "  webex config set client-secret <your-client-secret>"
    Write-Host "  webex login"
} finally {
    Remove-Item -Recurse -Force $TmpDir -ErrorAction SilentlyContinue
}
