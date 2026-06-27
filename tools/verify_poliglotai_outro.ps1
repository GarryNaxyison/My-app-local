[CmdletBinding()]
param(
  [string]$VideoPath = "assets/video-outro/poliglotai-outro-logo-burst.mp4",
  [string]$PreviewPath = "assets/video-outro/poliglotai-outro-logo-burst-preview.png"
)

$ErrorActionPreference = "Stop"

function Resolve-RepoPath {
  param([string]$Path)
  if ([System.IO.Path]::IsPathRooted($Path)) {
    return $Path
  }
  return Join-Path (Get-Location) $Path
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

Assert-True (Test-Path $videoFullPath) "Missing outro MP4: $videoFullPath"
Assert-True (Test-Path $previewFullPath) "Missing outro preview PNG: $previewFullPath"

$ffprobe = Get-Command ffprobe -ErrorAction SilentlyContinue
Assert-True ($null -ne $ffprobe) "ffprobe is not available in PATH."

$probeRaw = & $ffprobe.Source -v error -print_format json -show_streams -show_format $videoFullPath
Assert-True ($LASTEXITCODE -eq 0) "ffprobe failed for $videoFullPath"

$probe = $probeRaw | ConvertFrom-Json
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

$frameRateParts = ([string]$video.avg_frame_rate).Split("/")
Assert-Equal $frameRateParts.Count 2 "Unexpected avg_frame_rate format."
$frameRate = [double]::Parse($frameRateParts[0], [System.Globalization.CultureInfo]::InvariantCulture) / [double]::Parse($frameRateParts[1], [System.Globalization.CultureInfo]::InvariantCulture)
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
