param([int]$Port = 4178)
$ErrorActionPreference = 'Stop'
$appDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$workspace = Split-Path -Parent (Split-Path -Parent $appDir)
$env:CODEX_ORCHESTRATOR_WORKSPACE = $workspace
$env:CODEX_ORCHESTRATOR_PORT = [string]$Port
$url = "http://127.0.0.1:$Port"

try {
    $null = Invoke-RestMethod -Uri "$url/api/status" -TimeoutSec 2
} catch {
    $node = (Get-Command node.exe -ErrorAction Stop).Source
    Start-Process -FilePath $node -ArgumentList @('server.js') -WorkingDirectory $appDir -WindowStyle Hidden
    $ready = $false
    for ($attempt = 0; $attempt -lt 40; $attempt++) {
        Start-Sleep -Milliseconds 250
        try { $null = Invoke-RestMethod -Uri "$url/api/status" -TimeoutSec 2; $ready = $true; break } catch {}
    }
    if (-not $ready) { throw "Codex Orchestrator не запустился. Проверь, что порт $Port свободен." }
}

$edge = Get-Command msedge.exe -ErrorAction SilentlyContinue
$edgePath = if ($edge) { $edge.Source } else { $null }
if (-not $edgePath) {
    foreach ($candidate in @("$env:ProgramFiles(x86)\Microsoft\Edge\Application\msedge.exe", "$env:ProgramFiles\Microsoft\Edge\Application\msedge.exe")) {
        if (Test-Path -LiteralPath $candidate) { $edgePath = $candidate; break }
    }
}
if ($edgePath) { Start-Process -FilePath $edgePath -ArgumentList @('--app=' + $url) }
else { Start-Process $url }
