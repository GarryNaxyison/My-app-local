param(
  [switch]$Persist
)

$ErrorActionPreference = 'Stop'

function Add-PathEntry {
  param([string]$PathEntry)
  if ([string]::IsNullOrWhiteSpace($PathEntry) -or -not (Test-Path -LiteralPath $PathEntry)) {
    return
  }
  $parts = [Environment]::GetEnvironmentVariable('PATH', 'Process') -split ';' | Where-Object { $_ }
  if ($parts -notcontains $PathEntry) {
    $env:PATH = $PathEntry + ';' + $env:PATH
  }
}

function Add-UserPathEntry {
  param([string]$PathEntry)
  if (-not $Persist -or [string]::IsNullOrWhiteSpace($PathEntry) -or -not (Test-Path -LiteralPath $PathEntry)) {
    return
  }
  $current = [Environment]::GetEnvironmentVariable('PATH', 'User')
  $parts = $current -split ';' | Where-Object { $_ }
  if ($parts -notcontains $PathEntry) {
    $next = if ([string]::IsNullOrWhiteSpace($current)) { $PathEntry } else { $current.TrimEnd(';') + ';' + $PathEntry }
    [Environment]::SetEnvironmentVariable('PATH', $next, 'User')
  }
}

function Find-ExecutableDir {
  param(
    [string]$Executable,
    [string[]]$CandidateDirs
  )
  foreach ($dir in $CandidateDirs) {
    if ([string]::IsNullOrWhiteSpace($dir)) {
      continue
    }
    $candidate = Join-Path $dir $Executable
    if (Test-Path -LiteralPath $candidate) {
      return $dir
    }
  }
  return $null
}

function Get-CommandSource {
  param([string]$CommandName)
  $cmd = Get-Command $CommandName -ErrorAction SilentlyContinue
  if ($cmd -and $cmd.Source -and ($cmd.Source -notlike '*\WindowsApps\*')) {
    return $cmd.Source
  }
  return $null
}

$goDirs = @(
  (Join-Path $env:ProgramFiles 'Go\bin'),
  (Join-Path ${env:ProgramFiles(x86)} 'Go\bin'),
  (Join-Path $env:LOCALAPPDATA 'Programs\Go\bin'),
  (Join-Path $env:USERPROFILE 'go\bin')
) | Where-Object { $_ }

$pythonDirs = @()
$pythonRoots = @(
  (Join-Path $env:LOCALAPPDATA 'Programs\Python'),
  $env:ProgramFiles,
  ${env:ProgramFiles(x86)}
) | Where-Object { $_ -and (Test-Path -LiteralPath $_) }

foreach ($root in $pythonRoots) {
  Get-ChildItem -LiteralPath $root -Directory -Filter 'Python*' -ErrorAction SilentlyContinue |
    ForEach-Object {
      $pythonDirs += $_.FullName
      $pythonDirs += (Join-Path $_.FullName 'Scripts')
    }
}

$goDir = Find-ExecutableDir -Executable 'go.exe' -CandidateDirs $goDirs
if ($goDir) {
  Add-PathEntry $goDir
  Add-UserPathEntry $goDir
}

$pythonDir = Find-ExecutableDir -Executable 'python.exe' -CandidateDirs $pythonDirs
if ($pythonDir) {
  Add-PathEntry $pythonDir
  Add-UserPathEntry $pythonDir
  $scriptsDir = Join-Path (Split-Path $pythonDir -Parent) 'Scripts'
  Add-PathEntry $scriptsDir
  Add-UserPathEntry $scriptsDir
}

$go = Get-CommandSource 'go'
$gofmt = Get-CommandSource 'gofmt'
$python = Get-CommandSource 'python'

if ($go) {
  & $go version
} else {
  Write-Warning 'go.exe was not found. Install Go, then rerun this script.'
}

if ($gofmt) {
  Write-Host "gofmt: $gofmt"
} else {
  Write-Warning 'gofmt.exe was not found. It is installed with Go and should be in the same bin directory as go.exe.'
}

if ($python) {
  & $python --version
} else {
  Write-Warning 'A real python.exe was not found. WindowsApps Store aliases are ignored; install Python from python.org or winget.'
}

if ($Persist) {
  Write-Host 'User PATH was updated for found tools. Open a new terminal to inherit it.'
} else {
  Write-Host 'PATH was updated for this PowerShell process only. Pass -Persist to update User PATH.'
}
