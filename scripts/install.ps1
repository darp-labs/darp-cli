#!/usr/bin/env pwsh
#Requires -Version 7.6

[CmdletBinding()]
param(
    [string]$Version = "",
    [string]$Repository = "darp-labs/darp-cli"
)

$ErrorActionPreference = "Stop"
$MinimumCosignVersion = [System.Management.Automation.SemanticVersion]"3.1.4"
$HelperName = "darp-trust-bootstrap_windows_amd64.exe"
$HelperSha256 = "fade33f32b8c03a8bdaa48bd601099f81fea7392dbe8adf98b4e352b93ec2ed7"

function Fail([string]$Message) {
    throw "darp installer: $Message"
}

function Test-SemVer([string]$Value, [switch]$AllowPrerelease) {
    $pattern = if ($AllowPrerelease) {
        '^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$'
    } else {
        '^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$'
    }
    return $Value -match $pattern
}

function Get-Sha256([string]$Path) {
    return (Get-FileHash -Algorithm SHA256 -LiteralPath $Path).Hash.ToLowerInvariant()
}

function Invoke-HttpsDownload([string]$Uri, [string]$Destination) {
    if (-not $Uri.StartsWith("https://", [System.StringComparison]::OrdinalIgnoreCase)) {
        Fail "refusing insecure download: $Uri"
    }
    try {
        Invoke-WebRequest -Uri $Uri -OutFile $Destination -MaximumRedirection 5
    } catch {
        Fail "failed to download $Uri`: $($_.Exception.Message)"
    }
    if (-not (Test-Path -LiteralPath $Destination) -or (Get-Item -LiteralPath $Destination).Length -eq 0) {
        Fail "download is missing or empty: $Uri"
    }
}

function Resolve-Release([string]$Requested) {
    $base = "https://api.github.com/repos/$Repository"
    $api = if ([string]::IsNullOrWhiteSpace($Requested)) {
        "$base/releases/latest"
    } else {
        $tag = if ($Requested.StartsWith("v")) { $Requested } else { "v$Requested" }
        if (-not (Test-SemVer $Requested -AllowPrerelease)) {
            Fail "invalid explicit version '$Requested'"
        }
        "$base/releases/tags/$tag"
    }
    try {
        $release = Invoke-RestMethod -Uri $api -Headers @{ Accept = "application/vnd.github+json"; "X-Github-Api-Version" = "2022-11-28" } -MaximumRedirection 5
    } catch {
        if ([string]::IsNullOrWhiteSpace($Requested)) {
            Fail "could not resolve the latest stable release: $($_.Exception.Message)"
        }
        Fail "explicit release '$Requested' was not found; no fallback is used: $($_.Exception.Message)"
    }
    if (-not $release.tag_name -or -not (Test-SemVer $release.tag_name -AllowPrerelease)) {
        Fail "release returned an invalid tag"
    }
    $tag = [string]$release.tag_name
    $number = $tag.TrimStart("v")
    if (-not [string]::IsNullOrWhiteSpace($Requested)) {
        $requestedNumber = $Requested.TrimStart("v")
        if ($number -ne $requestedNumber) {
            Fail "resolved release $tag does not match explicit version $Requested"
        }
    }
    if (-not $release.assets) {
        Fail "release $tag has no assets"
    }
    $sourceSha = $null
    try {
        # The commits endpoint dereferences both lightweight and annotated tags.
        $commit = Invoke-RestMethod -Uri "$base/commits/$tag" -Headers @{ Accept = "application/vnd.github+json"; "X-Github-Api-Version" = "2022-11-28" }
        $sourceSha = [string]$commit.sha
    } catch {
        Fail "could not resolve the source revision for release $tag"
    }
    if ($sourceSha -notmatch '^[0-9a-fA-F]{40}$') {
        Fail "release $tag returned an invalid source revision"
    }
    return @{ Tag = $tag; Number = $number; SourceSha = $sourceSha.ToLowerInvariant() }
}

if ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture -ne [System.Runtime.InteropServices.Architecture]::X64) {
    Fail "unsupported architecture; DARP supports Windows AMD64 only"
}
$psVersion = $PSVersionTable.PSVersion
if ($psVersion.Major -lt 7 -or ($psVersion.Major -eq 7 -and $psVersion.Minor -lt 6)) {
    Fail "PowerShell 7.6 or newer is required"
}

$release = Resolve-Release $Version
$tag = [string]$release.Tag
$number = [string]$release.Number
$sourceSha = [string]$release.SourceSha
$assetName = "darp_${number}_windows_amd64.zip"
$downloadBase = "https://github.com/$Repository/releases/download/$tag"
$tempDir = Join-Path ([System.IO.Path]::GetTempPath()) ("darp-install-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

$archive = Join-Path $tempDir $assetName
$checksums = Join-Path $tempDir "checksums.txt"
$bundle = Join-Path $tempDir "checksums.txt.sigstore.json"
try {
    Invoke-HttpsDownload "$downloadBase/$assetName" $archive
    Invoke-HttpsDownload "$downloadBase/checksums.txt" $checksums
    Invoke-HttpsDownload "$downloadBase/checksums.txt.sigstore.json" $bundle

    $cosign = Get-Command cosign -ErrorAction SilentlyContinue
    $cosignPath = $null
    if ($cosign) {
        $cosignVersionOutput = (& $cosign.Source version 2>&1 | Out-String).Trim()
        $cosignVersionMatch = [regex]::Match($cosignVersionOutput, '(?im)(?:GitVersion|Version):?\s*v?([0-9]+\.[0-9]+\.[0-9]+)')
        if ($cosignVersionMatch.Success) {
            $cosignVersion = [System.Management.Automation.SemanticVersion]$cosignVersionMatch.Groups[1].Value
            if ($cosignVersion -eq $MinimumCosignVersion) {
                $cosignPath = $cosign.Source
            }
        }
    }

    if (-not $cosignPath) {
        $helper = Join-Path $tempDir $HelperName
        Invoke-HttpsDownload "$downloadBase/$HelperName" $helper
        $helperHash = Get-Sha256 $helper
        if ($helperHash -ne $HelperSha256) {
            Fail "bootstrap helper SHA-256 mismatch"
        }
        $helperResult = & $helper "-version" "v3.1.4"
        if ($LASTEXITCODE -ne 0) {
            Fail "authenticated Cosign bootstrap failed"
        }
        $cosignPath = ($helperResult | Select-Object -Last 1).Trim()
        if (-not $cosignPath -or -not (Test-Path -LiteralPath $cosignPath)) {
            Fail "authenticated Cosign bootstrap did not return an executable"
        }
    }

    & $cosignPath verify-blob `
        --bundle $bundle `
        --certificate-identity "https://github.com/$Repository/.github/workflows/release.yml@refs/tags/$tag" `
        --certificate-oidc-issuer "https://token.actions.githubusercontent.com" `
        --certificate-github-workflow-repository $Repository `
        --certificate-github-workflow-ref "refs/tags/$tag" `
        --certificate-github-workflow-sha $sourceSha | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Fail "signature verification failed"
    }

    $manifestEntry = Get-Content -LiteralPath $checksums |
        Where-Object { $_ -match '\s+\*?' + [regex]::Escape($assetName) + '$' }
    if (@($manifestEntry).Count -ne 1) {
        Fail "checksum manifest is missing a unique entry for $assetName"
    }
    $expectedHash = (($manifestEntry | Select-Object -First 1).Trim() -split '\s+', 2)[0].ToLowerInvariant()
    if ($expectedHash -notmatch '^[0-9a-f]{64}$') {
        Fail "checksum manifest contains an invalid SHA-256 for $assetName"
    }
    if ((Get-Sha256 $archive) -ne $expectedHash) {
        Fail "artifact SHA-256 mismatch for $assetName"
    }

    $extractDir = Join-Path $tempDir "extract"
    Expand-Archive -LiteralPath $archive -DestinationPath $extractDir -Force
    $candidate = Join-Path $extractDir "darp.exe"
    if (-not (Test-Path -LiteralPath $candidate) -or (Get-Item -LiteralPath $candidate).Length -eq 0) {
        Fail "archive does not contain a valid root-level darp.exe"
    }

    $installDir = Join-Path $env:LOCALAPPDATA "Programs\darp"
    $target = Join-Path $installDir "darp.exe"
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null

    $backup = $null
    if (Test-Path -LiteralPath $target) {
        $backup = Join-Path $tempDir "darp.exe.backup"
        Move-Item -LiteralPath $target -Destination $backup -Force
    }
    try {
        Copy-Item -LiteralPath $candidate -Destination $target -Force
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        $pathParts = @($userPath -split ';' | Where-Object { $_ })
        if ($pathParts -notcontains $installDir) {
            $newUserPath = (@($pathParts) + $installDir) -join ';'
            [Environment]::SetEnvironmentVariable("Path", $newUserPath, "User")
        }
    } catch {
        if (Test-Path -LiteralPath $target) { Remove-Item -LiteralPath $target -Force -ErrorAction SilentlyContinue }
        if ($backup) { Move-Item -LiteralPath $backup -Destination $target -Force }
        Fail "failed to install darp.exe or update the current-user PATH: $($_.Exception.Message)"
    }

    Write-Host "Installed darp $tag to $target"
} catch {
    if ($_.Exception.Message -like "darp installer:*") {
        throw
    }
    Fail $_.Exception.Message
} finally {
    Remove-Item -LiteralPath $tempDir -Recurse -Force -ErrorAction SilentlyContinue
}
