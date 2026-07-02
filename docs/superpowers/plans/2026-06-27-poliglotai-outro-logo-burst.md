# NERIVA Outro Logo Burst Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a reusable 1080 x 1920 MP4 outro with the existing NERIVA logo and the exact text `NERIVA.online`.

**Architecture:** Add one verifier script and one deterministic generator script. The verifier is the TDD acceptance test: it fails while the MP4 is missing, then passes after the generator renders the video and preview frame. The generator renders PNG frames with Windows `System.Drawing`, then encodes them with `ffmpeg` to H.264 MP4.

**Tech Stack:** PowerShell, Windows `System.Drawing`, `ffmpeg`, `ffprobe`, existing PNG logo assets.

---

## File Structure

- Create `tools/verify_poliglotai_outro.ps1`: media acceptance test. It checks file existence, dimensions, codec, frame rate, duration, lack of audio, and preview image dimensions.
- Create `tools/create_poliglotai_outro.ps1`: deterministic renderer. It draws the Logo Burst frames from the existing logo and encodes them to MP4.
- Create generated `assets/video-outro/poliglotai-outro-logo-burst.mp4`: final reusable outro.
- Create generated `assets/video-outro/poliglotai-outro-logo-burst-preview.png`: still frame from the final readable hold.

Use `web/assets/brand-logo.png` as the primary logo. It is a 512 x 512 PNG with alpha. Use `web/assets/brand-logo-mini.png` only as fallback if the primary logo is missing.

## Task 1: Add The Failing Media Acceptance Test

**Files:**
- Create: `tools/verify_poliglotai_outro.ps1`

- [ ] **Step 1: Write the verifier script**

Create `tools/verify_poliglotai_outro.ps1` with this content:

```powershell
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

Write-Host "NERIVA outro verification passed."
Write-Host "Video: $videoFullPath"
Write-Host "Duration: $duration"
Write-Host "Frame rate: $frameRate"
```

- [ ] **Step 2: Run the verifier and watch it fail**

Run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools/verify_poliglotai_outro.ps1
```

Expected result:

```text
Missing outro MP4: ...\assets\video-outro\poliglotai-outro-logo-burst.mp4
```

The command must exit non-zero. That failure proves the acceptance test catches the missing generated asset.

- [ ] **Step 3: Commit the failing verifier**

Run:

```bash
git add tools/verify_poliglotai_outro.ps1
git commit -m "test: add NERIVA outro verifier"
```

Expected result: one commit containing only `tools/verify_poliglotai_outro.ps1`.

## Task 2: Add The Deterministic Outro Generator

**Files:**
- Create: `tools/create_poliglotai_outro.ps1`
- Generate: `assets/video-outro/poliglotai-outro-logo-burst.mp4`
- Generate: `assets/video-outro/poliglotai-outro-logo-burst-preview.png`

- [ ] **Step 1: Write the generator script**

Create `tools/create_poliglotai_outro.ps1` with this content:

```powershell
[CmdletBinding()]
param(
  [string]$OutputDir = "assets/video-outro",
  [string]$PrimaryLogo = "web/assets/brand-logo.png",
  [string]$FallbackLogo = "web/assets/brand-logo-mini.png",
  [int]$Width = 1080,
  [int]$Height = 1920,
  [int]$Fps = 30,
  [double]$DurationSeconds = 2.2
)

$ErrorActionPreference = "Stop"

function Resolve-RepoPath {
  param([string]$Path)
  if ([System.IO.Path]::IsPathRooted($Path)) {
    return $Path
  }
  return Join-Path (Get-Location) $Path
}

function Assert-Tool {
  param([string]$Name)
  $command = Get-Command $Name -ErrorAction SilentlyContinue
  if ($null -eq $command) {
    throw "$Name is not available in PATH."
  }
  return $command.Source
}

function Ease-Out-Cubic {
  param([double]$Value)
  $v = [Math]::Min(1.0, [Math]::Max(0.0, $Value))
  return 1.0 - [Math]::Pow(1.0 - $v, 3.0)
}

function Ease-Out-Back {
  param([double]$Value)
  $v = [Math]::Min(1.0, [Math]::Max(0.0, $Value))
  $c1 = 1.70158
  $c3 = $c1 + 1.0
  return 1.0 + $c3 * [Math]::Pow($v - 1.0, 3.0) + $c1 * [Math]::Pow($v - 1.0, 2.0)
}

function Smooth-Window {
  param(
    [double]$Time,
    [double]$Start,
    [double]$End
  )
  return Ease-Out-Cubic (($Time - $Start) / ($End - $Start))
}

Add-Type -AssemblyName System.Drawing
Add-Type -AssemblyName System.Windows.Forms

$ffmpeg = Assert-Tool "ffmpeg"
$ffprobe = Assert-Tool "ffprobe"

$outputFullDir = Resolve-RepoPath $OutputDir
New-Item -ItemType Directory -Force -Path $outputFullDir | Out-Null

$videoPath = Join-Path $outputFullDir "poliglotai-outro-logo-burst.mp4"
$previewPath = Join-Path $outputFullDir "poliglotai-outro-logo-burst-preview.png"
$framesDir = Join-Path $outputFullDir "_frames_logo_burst"

if (Test-Path $framesDir) {
  Remove-Item -LiteralPath $framesDir -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $framesDir | Out-Null

$logoPath = Resolve-RepoPath $PrimaryLogo
if (-not (Test-Path $logoPath)) {
  $logoPath = Resolve-RepoPath $FallbackLogo
}
if (-not (Test-Path $logoPath)) {
  throw "No logo file found. Checked '$PrimaryLogo' and '$FallbackLogo'."
}

$logo = [System.Drawing.Image]::FromFile($logoPath)
$frameCount = [int][Math]::Round($DurationSeconds * $Fps)
$fontFamily = New-Object System.Drawing.FontFamily("Segoe UI")
$domainText = "NERIVA.online"

try {
  for ($frame = 0; $frame -lt $frameCount; $frame++) {
    $time = $frame / $Fps
    $bitmap = New-Object System.Drawing.Bitmap($Width, $Height, [System.Drawing.Imaging.PixelFormat]::Format32bppArgb)
    $graphics = [System.Drawing.Graphics]::FromImage($bitmap)
    $graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
    $graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
    $graphics.TextRenderingHint = [System.Drawing.Text.TextRenderingHint]::AntiAliasGridFit

    $bgRect = New-Object System.Drawing.Rectangle(0, 0, $Width, $Height)
    $bgBrush = New-Object System.Drawing.Drawing2D.LinearGradientBrush($bgRect, [System.Drawing.Color]::FromArgb(255, 5, 9, 20), [System.Drawing.Color]::FromArgb(255, 7, 17, 31), 90)
    $graphics.FillRectangle($bgBrush, $bgRect)
    $bgBrush.Dispose()

    $flashProgress = [Math]::Max(0.0, 1.0 - ($time / 0.38))
    if ($flashProgress -gt 0) {
      for ($i = 0; $i -lt 8; $i++) {
        $radius = 90 + ($i * 115) + ((1.0 - $flashProgress) * 220)
        $alpha = [int](75 * $flashProgress * (1.0 - ($i / 10.0)))
        $brush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb($alpha, 18, 184, 215))
        $graphics.FillEllipse($brush, ($Width / 2) - $radius, ($Height * 0.43) - $radius, $radius * 2, $radius * 2)
        $brush.Dispose()
      }
    }

    for ($stripe = 0; $stripe -lt 7; $stripe++) {
      $x = (($time * 820) + ($stripe * 235)) % ($Width + 420) - 210
      $path = New-Object System.Drawing.Drawing2D.GraphicsPath
      $path.AddPolygon(@(
        [System.Drawing.PointF]::new([float]$x, 0),
        [System.Drawing.PointF]::new([float]($x + 82), 0),
        [System.Drawing.PointF]::new([float]($x - 300), [float]$Height),
        [System.Drawing.PointF]::new([float]($x - 382), [float]$Height)
      ))
      $stripeBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb(22, 45, 91, 255))
      $graphics.FillPath($stripeBrush, $path)
      $stripeBrush.Dispose()
      $path.Dispose()
    }

    $ringProgress = Smooth-Window $time 0.55 1.15
    if ($time -ge 0.55 -and $time -le 1.35) {
      $ringRadius = 150 + ($ringProgress * 430)
      $ringAlpha = [int](210 * (1.0 - [Math]::Min(1.0, $ringProgress)) + 24)
      $ringPen = New-Object System.Drawing.Pen([System.Drawing.Color]::FromArgb($ringAlpha, 143, 241, 208), 9)
      $graphics.DrawEllipse($ringPen, ($Width / 2) - $ringRadius, ($Height * 0.44) - ($ringRadius * 0.62), $ringRadius * 2, $ringRadius * 1.24)
      $ringPen.Dispose()
    }

    $logoProgress = Smooth-Window $time 0.25 0.75
    $logoScale = 0.18 + ((Ease-Out-Back $logoProgress) * 0.62)
    if ($time -lt 0.25) {
      $logoScale = 0.04
    }
    $logoSize = [int](512 * $logoScale)
    $logoX = [int](($Width - $logoSize) / 2)
    $logoY = [int](($Height * 0.44) - ($logoSize / 2))

    for ($glow = 5; $glow -ge 1; $glow--) {
      $glowSize = $logoSize + ($glow * 42)
      $alpha = [int](16 / $glow + 5)
      $glowBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb($alpha, 18, 184, 215))
      $graphics.FillEllipse($glowBrush, ($Width - $glowSize) / 2, ($Height * 0.44) - ($glowSize / 2), $glowSize, $glowSize)
      $glowBrush.Dispose()
    }

    $graphics.DrawImage($logo, $logoX, $logoY, $logoSize, $logoSize)

    $textProgress = Smooth-Window $time 0.75 1.35
    $textAlpha = [int](255 * $textProgress)
    $textYOffset = [int]((1.0 - $textProgress) * 90)
    $font = New-Object System.Drawing.Font($fontFamily, 86, [System.Drawing.FontStyle]::Bold, [System.Drawing.GraphicsUnit]::Pixel)
    $textSize = $graphics.MeasureString($domainText, $font)
    $textX = [float](($Width - $textSize.Width) / 2)
    $textY = [float](($Height * 0.61) + $textYOffset)

    for ($glow = 5; $glow -ge 1; $glow--) {
      $glowBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb([int](24 * $textProgress), 18, 184, 215))
      $graphics.DrawString($domainText, $font, $glowBrush, $textX - ($glow * 2), $textY)
      $graphics.DrawString($domainText, $font, $glowBrush, $textX + ($glow * 2), $textY)
      $graphics.DrawString($domainText, $font, $glowBrush, $textX, $textY - ($glow * 2))
      $graphics.DrawString($domainText, $font, $glowBrush, $textX, $textY + ($glow * 2))
      $glowBrush.Dispose()
    }

    $shadowBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb([int](155 * $textProgress), 0, 0, 0))
    $textBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb($textAlpha, 246, 252, 255))
    $accentPen = New-Object System.Drawing.Pen([System.Drawing.Color]::FromArgb([int](210 * $textProgress), 243, 184, 75), 4)
    $graphics.DrawString($domainText, $font, $shadowBrush, $textX + 4, $textY + 6)
    $graphics.DrawString($domainText, $font, $textBrush, $textX, $textY)
    $graphics.DrawLine($accentPen, [float]($textX + 40), [float]($textY + $textSize.Height + 16), [float]($textX + $textSize.Width - 40), [float]($textY + $textSize.Height + 16))

    $accentPen.Dispose()
    $textBrush.Dispose()
    $shadowBrush.Dispose()
    $font.Dispose()

    $framePath = Join-Path $framesDir ("frame_{0:D4}.png" -f $frame)
    $bitmap.Save($framePath, [System.Drawing.Imaging.ImageFormat]::Png)

    if ($frame -eq [int][Math]::Floor(1.65 * $Fps)) {
      $bitmap.Save($previewPath, [System.Drawing.Imaging.ImageFormat]::Png)
    }

    $graphics.Dispose()
    $bitmap.Dispose()
  }
} finally {
  $logo.Dispose()
  $fontFamily.Dispose()
}

$framePattern = Join-Path $framesDir "frame_%04d.png"
& $ffmpeg -y -framerate $Fps -i $framePattern -c:v libx264 -pix_fmt yuv420p -movflags +faststart -an $videoPath
if ($LASTEXITCODE -ne 0) {
  throw "ffmpeg failed to encode $videoPath"
}

& $ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 $videoPath | Out-Null
if ($LASTEXITCODE -ne 0) {
  throw "ffprobe failed to read $videoPath"
}

Remove-Item -LiteralPath $framesDir -Recurse -Force

Write-Host "Created NERIVA outro:"
Write-Host $videoPath
Write-Host $previewPath
```

- [ ] **Step 2: Run the generator**

Run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools/create_poliglotai_outro.ps1
```

Expected result:

```text
Created NERIVA outro:
...\assets\video-outro\poliglotai-outro-logo-burst.mp4
...\assets\video-outro\poliglotai-outro-logo-burst-preview.png
```

- [ ] **Step 3: Run the verifier and watch it pass**

Run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools/verify_poliglotai_outro.ps1
```

Expected result:

```text
NERIVA outro verification passed.
Video: ...\assets\video-outro\poliglotai-outro-logo-burst.mp4
Duration: 2.2...
Frame rate: 30
```

- [ ] **Step 4: Commit the generator and generated assets**

Run:

```bash
git add tools/create_poliglotai_outro.ps1 assets/video-outro/poliglotai-outro-logo-burst.mp4 assets/video-outro/poliglotai-outro-logo-burst-preview.png
git commit -m "feat: add NERIVA outro video"
```

Expected result: one commit containing the generator script, MP4, and preview PNG.

## Task 3: Final Verification And Push

**Files:**
- Verify: `tools/create_poliglotai_outro.ps1`
- Verify: `tools/verify_poliglotai_outro.ps1`
- Verify: `assets/video-outro/poliglotai-outro-logo-burst.mp4`
- Verify: `assets/video-outro/poliglotai-outro-logo-burst-preview.png`

- [ ] **Step 1: Re-run the full generation and verification cycle**

Run:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools/create_poliglotai_outro.ps1
powershell -NoProfile -ExecutionPolicy Bypass -File tools/verify_poliglotai_outro.ps1
```

Expected result:

```text
Created NERIVA outro:
...\assets\video-outro\poliglotai-outro-logo-burst.mp4
...\assets\video-outro\poliglotai-outro-logo-burst-preview.png
NERIVA outro verification passed.
```

- [ ] **Step 2: Inspect the final media metadata**

Run:

```powershell
ffprobe -v error -select_streams v:0 -show_entries stream=codec_name,width,height,avg_frame_rate -show_entries format=duration -of default=noprint_wrappers=1 assets/video-outro/poliglotai-outro-logo-burst.mp4
```

Expected result includes:

```text
codec_name=h264
width=1080
height=1920
avg_frame_rate=30/1
duration=2.200000
```

- [ ] **Step 3: Confirm git status only contains intended files**

Run:

```bash
git status --short
```

Expected result may still include pre-existing unrelated workspace changes such as `site-react/playwright-report/index.html` or `english-coach-bot`, but the outro work should be limited to:

```text
tools/create_poliglotai_outro.ps1
tools/verify_poliglotai_outro.ps1
assets/video-outro/poliglotai-outro-logo-burst.mp4
assets/video-outro/poliglotai-outro-logo-burst-preview.png
```

- [ ] **Step 4: Push the current branch**

Run:

```bash
git push origin HEAD
```

Expected result: current branch is pushed to `https://github.com/GarryNaxyison/My-app-local.git`.

## Self-Review Notes

- Spec coverage: the plan covers output format, 9:16 resolution, duration, 30 fps, H.264, no audio, existing logo use, exact text, Logo Burst animation, final hold, preview frame, and `ffprobe` verification.
- Placeholder scan: the plan contains no unresolved placeholder terms or undefined implementation steps.
- Type consistency: script names, output paths, logo paths, dimensions, duration, and verification ranges match across all tasks.
