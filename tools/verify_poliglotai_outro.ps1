[CmdletBinding()]
param(
  [string]$VideoPath = "assets/video-outro/poliglotai-outro-logo-burst.mp4",
  [string]$PreviewPath = "assets/video-outro/poliglotai-outro-logo-burst-preview.png"
)

$ErrorActionPreference = "Stop"
$ScriptRoot = $PSScriptRoot
$RepoRoot = Split-Path -Parent $ScriptRoot

function Resolve-RepoPath {
  param([string]$Path)
  if ([System.IO.Path]::IsPathRooted($Path)) {
    return $Path
  }
  return Join-Path $RepoRoot $Path
}

function Assert-True {
  param(
    [bool]$Condition,
    [string]$Message
  )
  if (-not $Condition) {
    throw $Message
  }
}

function Assert-Equal {
  param(
    [object]$Actual,
    [object]$Expected,
    [string]$Message
  )
  if ($Actual -ne $Expected) {
    throw "$Message Expected '$Expected', got '$Actual'."
  }
}

$videoFullPath = Resolve-RepoPath $VideoPath
$previewFullPath = Resolve-RepoPath $PreviewPath

Assert-True (Test-Path -Path $videoFullPath -PathType Leaf) "Missing outro MP4: $videoFullPath"
Assert-True (Test-Path -Path $previewFullPath -PathType Leaf) "Missing outro preview PNG: $previewFullPath"

$ffprobe = Get-Command ffprobe -ErrorAction SilentlyContinue
Assert-True ($null -ne $ffprobe) "ffprobe is not available in PATH."

$probeRaw = & $ffprobe.Source -v error -print_format json -show_streams -show_format $videoFullPath
Assert-True ($LASTEXITCODE -eq 0) "ffprobe failed for $videoFullPath"

$probe = $null
try {
  $probe = $probeRaw | ConvertFrom-Json
} catch {
  throw "ffprobe returned invalid JSON for $videoFullPath"
}

Assert-True ($null -ne $probe.format) "ffprobe output is missing format data for $videoFullPath"
Assert-True ($null -ne $probe.format.duration -and [string]$probe.format.duration -ne "") "ffprobe output is missing format.duration for $videoFullPath"
$videoStreams = @($probe.streams | Where-Object { $_.codec_type -eq "video" })
$audioStreams = @($probe.streams | Where-Object { $_.codec_type -eq "audio" })
Assert-Equal $videoStreams.Count 1 "MP4 must contain exactly one video stream."
Assert-Equal $audioStreams.Count 0 "MP4 must not contain audio streams."

$video = $videoStreams[0]
Assert-Equal ([int]$video.width) 1080 "Video width mismatch."
Assert-Equal ([int]$video.height) 1920 "Video height mismatch."
Assert-Equal ([string]$video.codec_name) "h264" "Video codec mismatch."

$duration = [double]::Parse([string]$probe.format.duration, [System.Globalization.CultureInfo]::InvariantCulture)
Assert-True ($duration -ge 2.15 -and $duration -le 2.35) "Duration must be between 2.15 and 2.35 seconds. Actual: $duration"

$avgFrameRate = [string]$video.avg_frame_rate
Assert-True (-not [string]::IsNullOrWhiteSpace($avgFrameRate)) "ffprobe output is missing avg_frame_rate for $videoFullPath"
$frameRateParts = $avgFrameRate.Split("/")
Assert-Equal $frameRateParts.Count 2 "Unexpected avg_frame_rate format."
$numerator = [string]$frameRateParts[0]
$denominator = [string]$frameRateParts[1]
Assert-True (-not [string]::IsNullOrWhiteSpace($numerator) -and -not [string]::IsNullOrWhiteSpace($denominator)) "Unexpected avg_frame_rate value for $videoFullPath"
$denominatorValue = [double]::Parse($denominator, [System.Globalization.CultureInfo]::InvariantCulture)
Assert-True ($denominatorValue -ne 0) "avg_frame_rate denominator must not be zero for $videoFullPath"
$frameRate = [double]::Parse($numerator, [System.Globalization.CultureInfo]::InvariantCulture) / $denominatorValue
Assert-True ([Math]::Abs($frameRate - 30.0) -lt 0.01) "Frame rate must be 30 fps. Actual: $frameRate"

Add-Type -AssemblyName System.Drawing
$preview = [System.Drawing.Image]::FromFile($previewFullPath)
try {
  Assert-Equal $preview.Width 1080 "Preview width mismatch."
  Assert-Equal $preview.Height 1920 "Preview height mismatch."
} finally {
  $preview.Dispose()
}

$videoSize = (Get-Item $videoFullPath).Length
$previewSize = (Get-Item $previewFullPath).Length
Assert-True ($videoSize -gt 100000) "MP4 file is too small to be a rendered outro. Size: $videoSize bytes."
Assert-True ($previewSize -gt 10000) "Preview PNG is too small to be a rendered frame. Size: $previewSize bytes."

Write-Host "PoliglotAI outro verification passed."
Write-Host "Video: $videoFullPath"
Write-Host "Duration: $duration"
Write-Host "Frame rate: $frameRate"
