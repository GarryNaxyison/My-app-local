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
Add-Type -AssemblyName System.Drawing

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

function Assert-Tool {
  param([string]$Name)
  $tool = Get-Command $Name -ErrorAction SilentlyContinue
  Assert-True ($null -ne $tool) "$Name is not available in PATH."
  return $tool.Source
}

function Clamp01 {
  param([double]$Value)
  if ($Value -lt 0) { return 0.0 }
  if ($Value -gt 1) { return 1.0 }
  return $Value
}

function Lerp {
  param(
    [double]$A,
    [double]$B,
    [double]$T
  )
  return $A + (($B - $A) * $T)
}

function SmoothStep {
  param([double]$T)
  $t = Clamp01 $T
  return $t * $t * (3 - (2 * $t))
}

function EaseOutCubic {
  param([double]$T)
  $t = Clamp01 $T
  return 1 - [Math]::Pow(1 - $t, 3)
}

function EaseOutBack {
  param([double]$T)
  $t = Clamp01 $T
  $c1 = 1.70158
  $c3 = $c1 + 1
  return 1 + ($c3 * [Math]::Pow($t - 1, 3)) + ($c1 * [Math]::Pow($t - 1, 2))
}

function EaseOutExpo {
  param([double]$T)
  $t = Clamp01 $T
  if ($t -ge 1) { return 1.0 }
  return 1 - [Math]::Pow(2, -10 * $t)
}

function ColorFromHex {
  param(
    [int]$Alpha,
    [string]$Hex
  )
  $c = [System.Drawing.ColorTranslator]::FromHtml($Hex)
  return [System.Drawing.Color]::FromArgb($Alpha, $c.R, $c.G, $c.B)
}

function Setup-Graphics {
  param([System.Drawing.Graphics]$Graphics)
  $Graphics.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
  $Graphics.InterpolationMode = [System.Drawing.Drawing2D.InterpolationMode]::HighQualityBicubic
  $Graphics.PixelOffsetMode = [System.Drawing.Drawing2D.PixelOffsetMode]::HighQuality
  $Graphics.TextRenderingHint = [System.Drawing.Text.TextRenderingHint]::AntiAliasGridFit
}

function Draw-GlowText {
  param(
    [System.Drawing.Graphics]$Graphics,
    [string]$Text,
    [System.Drawing.Font]$Font,
    [System.Drawing.Brush]$GlowBrush,
    [System.Drawing.Brush]$FillBrush,
    [float]$X,
    [float]$Y
  )
  $offsets = @(
    @(0.0, 0.0, 24),
    @(2.0, 0.0, 20),
    @(-2.0, 0.0, 20),
    @(0.0, 2.0, 20),
    @(0.0, -2.0, 20),
    @(4.0, 3.0, 12),
    @(-4.0, 3.0, 12)
  )
  foreach ($entry in $offsets) {
    $dx = [float]$entry[0]
    $dy = [float]$entry[1]
    $alpha = [int]$entry[2]
    $c = $GlowBrush.Color
    $gBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb($alpha, $c.R, $c.G, $c.B))
    try {
      $Graphics.DrawString($Text, $Font, $gBrush, $X + $dx, $Y + $dy)
    } finally {
      $gBrush.Dispose()
    }
  }
  $Graphics.DrawString($Text, $Font, $FillBrush, $X, $Y)
}

function Save-Bitmap {
  param(
    [System.Drawing.Bitmap]$Bitmap,
    [string]$Path
  )
  $Bitmap.Save($Path, [System.Drawing.Imaging.ImageFormat]::Png)
}

$ffmpeg = Assert-Tool "ffmpeg"
$ffprobe = Assert-Tool "ffprobe"

$outputFull = Resolve-RepoPath $OutputDir
$framesDir = Join-Path $outputFull "_frames"
$videoPath = Join-Path $outputFull "poliglotai-outro-logo-burst.mp4"
$previewPath = Join-Path $outputFull "poliglotai-outro-logo-burst-preview.png"

$primaryLogoFull = Resolve-RepoPath $PrimaryLogo
$fallbackLogoFull = Resolve-RepoPath $FallbackLogo
$logoPath = if (Test-Path -Path $primaryLogoFull -PathType Leaf) { $primaryLogoFull } elseif (Test-Path -Path $fallbackLogoFull -PathType Leaf) { $fallbackLogoFull } else { throw "Neither logo asset exists: $primaryLogoFull or $fallbackLogoFull" }

Assert-True (($outputFull -ne $null) -and ($outputFull.Length -gt 0)) "Output directory is invalid."

if (-not [System.IO.Path]::GetFullPath($framesDir).StartsWith([System.IO.Path]::GetFullPath($outputFull), [System.StringComparison]::OrdinalIgnoreCase)) {
  throw "Temporary frames directory must stay inside the output directory."
}

New-Item -ItemType Directory -Force -Path $outputFull | Out-Null
if (Test-Path -Path $framesDir) {
  Remove-Item -LiteralPath $framesDir -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $framesDir | Out-Null

$frameCount = [int][Math]::Round($DurationSeconds * $Fps)
Assert-True ($frameCount -ge 1) "Frame count must be positive."
$previewIndex = [Math]::Min($frameCount - 1, [int][Math]::Round(1.65 * $Fps))

$logoImage = [System.Drawing.Image]::FromFile($logoPath)

function Render-Frame {
  param(
    [double]$Time,
    [string]$Path,
    [switch]$Preview
  )

  $bmp = New-Object System.Drawing.Bitmap($Width, $Height, [System.Drawing.Imaging.PixelFormat]::Format32bppArgb)
  $gfx = [System.Drawing.Graphics]::FromImage($bmp)
  try {
    Setup-Graphics $gfx

    $bgRect = New-Object System.Drawing.Rectangle(0, 0, $Width, $Height)
    $bgBrush = New-Object System.Drawing.Drawing2D.LinearGradientBrush($bgRect, [System.Drawing.Color]::FromArgb(255, 5, 8, 18), [System.Drawing.Color]::FromArgb(255, 6, 24, 47), 90)
    try {
      $blend = New-Object System.Drawing.Drawing2D.ColorBlend
      $blend.Positions = [float[]](0.0, 0.52, 1.0)
      $blend.Colors = [System.Drawing.Color[]]@(
        [System.Drawing.Color]::FromArgb(255, 4, 7, 16),
        [System.Drawing.Color]::FromArgb(255, 8, 20, 46),
        [System.Drawing.Color]::FromArgb(255, 3, 5, 12)
      )
      $bgBrush.InterpolationColors = $blend
      $gfx.FillRectangle($bgBrush, $bgRect)
    } finally {
      $bgBrush.Dispose()
    }

    $vignettePath = New-Object System.Drawing.Drawing2D.GraphicsPath
    try {
      $vignettePath.AddEllipse(-120, 180, $Width + 240, $Height - 60)
      $vignette = New-Object System.Drawing.Drawing2D.PathGradientBrush($vignettePath)
      try {
        $vignette.CenterColor = [System.Drawing.Color]::FromArgb(0, 0, 0, 0)
        $vignette.SurroundColors = [System.Drawing.Color[]]@([System.Drawing.Color]::FromArgb(170, 1, 4, 11))
        $gfx.FillEllipse($vignette, -120, 180, $Width + 240, $Height - 60)
      } finally {
        $vignette.Dispose()
      }
    } finally {
      $vignettePath.Dispose()
    }

    $flash = [Math]::Exp(-[Math]::Pow(($Time - 0.14) / 0.085, 2))
    if ($flash -gt 0.002) {
      $flashPath = New-Object System.Drawing.Drawing2D.GraphicsPath
      try {
        $flashPath.AddEllipse(($Width * 0.18), ($Height * 0.19), ($Width * 0.68), ($Height * 0.46))
        $flashBrush = New-Object System.Drawing.Drawing2D.PathGradientBrush($flashPath)
        try {
          $flashBrush.CenterColor = [System.Drawing.Color]::FromArgb([int](200 * $flash), 110, 255, 255)
          $flashBrush.SurroundColors = [System.Drawing.Color[]]@([System.Drawing.Color]::FromArgb(0, 110, 255, 255))
          $gfx.FillEllipse($flashBrush, ($Width * 0.12), ($Height * 0.11), ($Width * 0.76), ($Height * 0.54))
        } finally {
          $flashBrush.Dispose()
        }
      } finally {
        $flashPath.Dispose()
      }
    }

    $savedState = $gfx.Save()
    $gfx.TranslateTransform(($Width / 2.0), ($Height / 2.0))
    $gfx.RotateTransform(-24)
    $stripeShift = ($Time * 250.0) % 240.0
    for ($i = -18; $i -le 18; $i++) {
      $stripeBase = ([double]$i * 170.0)
      $offset = $stripeBase - [double]$stripeShift
      $distance = [Math]::Abs($i) / 18.0
      $alpha = [int](14 + (22 * (1 - $distance)))
      $stripeColor = [System.Drawing.Color]::FromArgb($alpha, 41, 173, 255)
      $stripeBrush = New-Object System.Drawing.SolidBrush($stripeColor)
      try {
        $pts = New-Object 'System.Drawing.PointF[]' 4
        $pts[0] = New-Object System.Drawing.PointF([float]$offset, -2200.0)
        $pts[1] = New-Object System.Drawing.PointF([float]($offset + 58.0), -2200.0)
        $pts[2] = New-Object System.Drawing.PointF([float]($offset + 260.0), 2200.0)
        $pts[3] = New-Object System.Drawing.PointF([float]($offset + 202.0), 2200.0)
        $gfx.FillPolygon($stripeBrush, $pts)
      } finally {
        $stripeBrush.Dispose()
      }
    }
    $gfx.Restore($savedState)

    $logoIntro = EaseOutBack (Clamp01 (($Time - 0.08) / 0.72))
    $logoHold = EaseOutCubic (Clamp01 (($Time - 0.62) / 0.75))
    $logoMotion = (Lerp 0.0 1.0 ([Math]::Min($logoIntro, $logoHold)))
    $logoScale = (Lerp 0.76 1.02 $logoMotion)
    $logoPulse = 1 + (0.015 * [Math]::Sin($Time * 20))
    $logoSize = 420 * $logoScale * $logoPulse
    $logoX = ($Width - $logoSize) / 2
    $logoY = (0.44 * $Height) - ($logoSize / 2) + (18 * (1 - $logoMotion))

    $ringSeed = SmoothStep (Clamp01 (($Time - 0.56) / 0.82))
    if ($ringSeed -gt 0.001) {
      $ringOpacity = [int](150 * (1 - $ringSeed))
      $ringSize = $logoSize + (Lerp 120 340 $ringSeed)
      $ringX = ($Width - $ringSize) / 2
      $ringY = $logoY + ($logoSize * 0.12) - (($ringSize - ($logoSize * 0.76)) / 2)
      $ringPen = New-Object System.Drawing.Pen([System.Drawing.Color]::FromArgb($ringOpacity, 64, 215, 255), 18)
      try {
        $ringPen.StartCap = [System.Drawing.Drawing2D.LineCap]::Round
        $ringPen.EndCap = [System.Drawing.Drawing2D.LineCap]::Round
        $gfx.DrawEllipse($ringPen, $ringX, $ringY, $ringSize, $ringSize)
      } finally {
        $ringPen.Dispose()
      }
      $ringPen2 = New-Object System.Drawing.Pen([System.Drawing.Color]::FromArgb([int]($ringOpacity * 0.55), 14, 122, 255), 5)
      try {
        $ringPen2.StartCap = [System.Drawing.Drawing2D.LineCap]::Round
        $ringPen2.EndCap = [System.Drawing.Drawing2D.LineCap]::Round
        $gfx.DrawEllipse($ringPen2, $ringX - 22, $ringY - 22, $ringSize + 44, $ringSize + 44)
      } finally {
        $ringPen2.Dispose()
      }
    }

    $logoGlowSize = $logoSize + 94
    $logoGlowX = $logoX - 47
    $logoGlowY = $logoY - 47
    $glowBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb([int](100 * (0.55 + 0.45 * $logoMotion)), 44, 200, 255))
    try {
      $gfx.FillEllipse($glowBrush, $logoGlowX, $logoGlowY + 14, $logoGlowSize, $logoGlowSize)
    } finally {
      $glowBrush.Dispose()
    }

    for ($n = 4; $n -ge 1; $n--) {
      $alpha = [int](24 / $n)
      $softGlow = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb($alpha, 64, 217, 255))
      try {
        $inflate = 14 * $n
        $gfx.DrawImage($logoImage, [int]($logoX - ($inflate / 2)), [int]($logoY - ($inflate / 2)), [int]($logoSize + $inflate), [int]($logoSize + $inflate))
      } finally {
        $softGlow.Dispose()
      }
    }
    $gfx.DrawImage($logoImage, [int]$logoX, [int]$logoY, [int]$logoSize, [int]$logoSize)

    $textIntro = EaseOutCubic (Clamp01 (($Time - 0.28) / 0.72))
    $textAlpha = [int](255 * $textIntro)
    $textY = (0.615 * $Height) + ((1 - $textIntro) * 44)
    $textX = 0.0
    $textFont = New-Object System.Drawing.Font("Segoe UI Semibold", 72, [System.Drawing.FontStyle]::Bold, [System.Drawing.GraphicsUnit]::Pixel)
    $measure = $gfx.MeasureString("PoliglotAI.online", $textFont)
    $textX = ($Width - $measure.Width) / 2
    $shadowBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb([int](110 * $textIntro), 12, 175, 255))
    $fillBrush = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb($textAlpha, 248, 250, 252))
    try {
      Draw-GlowText -Graphics $gfx -Text "PoliglotAI.online" -Font $textFont -GlowBrush $shadowBrush -FillBrush $fillBrush -X $textX -Y $textY
    } finally {
      $shadowBrush.Dispose()
      $fillBrush.Dispose()
      $textFont.Dispose()
    }

    $underlineWidth = (Lerp 220 318 $textIntro)
    $underlineX = ($Width - $underlineWidth) / 2
    $underlineY = $textY + $measure.Height + 18
    $underlinePen = New-Object System.Drawing.Pen([System.Drawing.Color]::FromArgb([int](170 * $textIntro), 248, 179, 73), 10)
    try {
      $underlinePen.StartCap = [System.Drawing.Drawing2D.LineCap]::Round
      $underlinePen.EndCap = [System.Drawing.Drawing2D.LineCap]::Round
      $gfx.DrawLine($underlinePen, $underlineX, $underlineY, $underlineX + $underlineWidth, $underlineY)
    } finally {
      $underlinePen.Dispose()
    }

    $finalGlow = New-Object System.Drawing.Pen([System.Drawing.Color]::FromArgb([int](90 * $textIntro), 82, 208, 255), 2)
    try {
      $gfx.DrawEllipse($finalGlow, ($Width * 0.15), ($Height * 0.18), ($Width * 0.70), ($Height * 0.52))
    } finally {
      $finalGlow.Dispose()
    }

    Save-Bitmap -Bitmap $bmp -Path $Path
  } finally {
    $gfx.Dispose()
    $bmp.Dispose()
  }
}

try {
  for ($index = 0; $index -lt $frameCount; $index++) {
    $time = $index / [double]$Fps
    $framePath = Join-Path $framesDir ("frame-{0:D4}.png" -f $index)
    Render-Frame -Time $time -Path $framePath
    if ($index -eq $previewIndex) {
      Copy-Item -LiteralPath $framePath -Destination $previewPath -Force
    }
  }

  if (-not (Test-Path -Path $previewPath -PathType Leaf)) {
    $previewTime = 1.65
    $previewFramePath = Join-Path $framesDir "preview-frame.png"
    Render-Frame -Time $previewTime -Path $previewFramePath
    Copy-Item -LiteralPath $previewFramePath -Destination $previewPath -Force
  }

  if (Test-Path -Path $videoPath) {
    Remove-Item -LiteralPath $videoPath -Force
  }

  & $ffmpeg -y -framerate $Fps -i (Join-Path $framesDir 'frame-%04d.png') -c:v libx264 -pix_fmt yuv420p -movflags +faststart -an $videoPath
  Assert-True ($LASTEXITCODE -eq 0) "ffmpeg failed while encoding the outro video."

  & $ffprobe -v error -show_streams -show_format $videoPath | Out-Null
  Assert-True ($LASTEXITCODE -eq 0) "ffprobe failed to read the encoded outro video."
} catch {
  Write-Host $_.Exception.Message
  if ($_.InvocationInfo -and $_.InvocationInfo.PositionMessage) {
    Write-Host $_.InvocationInfo.PositionMessage
  }
  throw
} finally {
  if (Test-Path -Path $framesDir) {
    Remove-Item -LiteralPath $framesDir -Recurse -Force
  }
  if ($logoImage) {
    $logoImage.Dispose()
  }
}

Write-Host "Generated video: $videoPath"
Write-Host "Generated preview: $previewPath"
