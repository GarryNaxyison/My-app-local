param(
    [string]$ComfyBase = "C:\Users\Admin\Documents\ComfyUI",
    [switch]$IncludeOptional,
    [switch]$SkipFreeSpaceCheck
)

$ErrorActionPreference = "Stop"

$ManifestPath = Join-Path $PSScriptRoot "model_manifest.json"
if (!(Test-Path -LiteralPath $ManifestPath)) {
    throw "Missing manifest: $ManifestPath"
}

$Manifest = Get-Content -Raw -LiteralPath $ManifestPath | ConvertFrom-Json
$Models = @($Manifest.models | Where-Object { $_.download -eq $true -or ($IncludeOptional -and $_.url) })

function Get-FreeBytesForPath {
    param([string]$Path)
    $root = [System.IO.Path]::GetPathRoot((Resolve-Path -LiteralPath (Split-Path -Parent $Path)).Path)
    $driveName = $root.Substring(0, 1)
    return (Get-PSDrive -Name $driveName).Free
}

function Resolve-Target {
    param([string]$RelativePath)
    return Join-Path $ComfyBase $RelativePath
}

function Test-ModelPresent {
    param([string]$Path, [double]$ExpectedGB)
    if (!(Test-Path -LiteralPath $Path)) {
        return $false
    }
    $item = Get-Item -LiteralPath $Path
    $minBytes = [int64]([Math]::Max(1MB, $ExpectedGB * 1GB * 0.80))
    return $item.Length -ge $minBytes
}

$missingBytes = [int64]0
foreach ($model in $Models) {
    $target = Resolve-Target $model.relativePath
    if (!(Test-ModelPresent -Path $target -ExpectedGB ([double]$model.sizeGB))) {
        $missingBytes += [int64](([double]$model.sizeGB) * 1GB)
    }
}

if (!$SkipFreeSpaceCheck -and $missingBytes -gt 0) {
    $probeTarget = Resolve-Target $Models[0].relativePath
    $freeBytes = Get-FreeBytesForPath $probeTarget
    $targetFreeBytes = [int64](([double]$Manifest.diskPlan.targetFreeAfterInstallGB) * 1GB)
    if (($freeBytes - $missingBytes) -lt $targetFreeBytes) {
        $freeGB = [Math]::Round($freeBytes / 1GB, 2)
        $missingGB = [Math]::Round($missingBytes / 1GB, 2)
        $targetFreeGB = $Manifest.diskPlan.targetFreeAfterInstallGB
        throw "Not enough free disk space. Free: ${freeGB} GB, needed downloads: ${missingGB} GB, target reserve: ${targetFreeGB} GB."
    }
}

function Download-File {
    param(
        [string]$Url,
        [string]$Target
    )

    $parent = Split-Path -Parent $Target
    New-Item -ItemType Directory -Force -Path $parent | Out-Null

    $tmp = "$Target.part"
    if (Test-Path -LiteralPath $tmp) {
        Remove-Item -LiteralPath $tmp -Force
    }

    Write-Host "Downloading $Url"
    Write-Host " -> $Target"
    try {
        Start-BitsTransfer -Source $Url -Destination $tmp -TransferType Download -Description "ComfyUI ad short model"
    } catch {
        Write-Warning "BITS failed, falling back to Invoke-WebRequest: $($_.Exception.Message)"
        Invoke-WebRequest -Uri $Url -OutFile $tmp -Headers @{ "User-Agent" = "ComfyUI-AdShorts-Installer" }
    }

    if (!(Test-Path -LiteralPath $tmp)) {
        throw "Download did not create temp file: $tmp"
    }

    Move-Item -LiteralPath $tmp -Destination $Target -Force
}

foreach ($model in $Models) {
    $target = Resolve-Target $model.relativePath
    if (Test-ModelPresent -Path $target -ExpectedGB ([double]$model.sizeGB)) {
        Write-Host "OK: $($model.name)"
        continue
    }
    Download-File -Url $model.url -Target $target
}

Write-Host "Done. Restart ComfyUI if it was open so the new model dropdowns refresh."
