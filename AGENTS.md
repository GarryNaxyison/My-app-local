# Project Agent Instructions

## Primary Workspace

Work in this repository by default: `E:\PROJECTS\New project\My app local`. Treat the parent folder `E:\PROJECTS\New project` as a container unless the user explicitly asks to work there.

## On-Demand MCP Usage

Do not verify, start, or call MCP servers by default. Use MCP tools only when the current task directly needs that capability or local shell/file inspection is insufficient. Prefer local repository context first, and use `MCP_STARTUP.md` as an on-demand trigger map rather than a startup checklist.

Headroom MCP is configured in Codex as `headroom`. Keep it available for context savings, but call `headroom_stats` only when MCP availability matters or before compression-heavy work. Use `headroom_compress` for large tool outputs, logs, search results, or files, and `headroom_retrieve` when the original content is needed. If the current Codex process does not expose the Headroom tools yet, restart Codex so it reloads `C:\Users\Admin\.codex\config.toml`.

Use Headroom MCP for large context payloads: tool outputs over ~200 lines, files over ~20 KB, large JSON/CSS/TSX/logs/search results, and repeated diagnostics. Do not compress short diffs, concise errors, or exact code fragments needed for editing unless they are unusually large.

MarkItDown MCP is configured in Codex as `markitdown`. Use it only when converting documents, URLs, PDFs, or other resources into Markdown is useful for the current task. It runs from the isolated venv at `C:\Users\Admin\.codex\vendor_imports\markitdown-mcp-venv` via `python.exe -m markitdown_mcp` and exposes `convert_to_markdown(uri)` for `file:`, `http:`, `https:`, and `data:` URIs. Treat inputs as trusted or sanitized because the server reads resources with the current user's privileges.

## Shell Preference

Use Git Bash for shell commands when practical. On this machine Git Bash is installed at:

```powershell
C:\Program Files\Git\bin\bash.exe
```

When `bash` is not available in `PATH`, invoke it explicitly from PowerShell:

```powershell
& 'C:\Program Files\Git\bin\bash.exe' -lc "git status --short"
```

## GitHub Sync

After verified file changes in this repository, commit and push them to the configured GitHub remote. The expected remote is `https://github.com/GarryNaxyison/My-app-local.git`. If push authentication is missing or fails, report that clearly.

## Nx MCP Startup

For this workspace, start Nx MCP as a standalone Streamable HTTP server when a persistent MCP endpoint is needed:

```powershell
npx -y nx-mcp@latest "E:\PROJECTS\New project\My app local" --transport http --port 9921 --no-minimal --disableTelemetry
```

Recommended background launch with logs:

```powershell
$workspace = 'E:\PROJECTS\New project\My app local'
$logDir = Join-Path $workspace 'tmp\mcp'
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
Start-Process -FilePath 'C:\Program Files\nodejs\npx.cmd' `
  -ArgumentList @('-y','nx-mcp@latest',$workspace,'--transport','http','--port','9921','--no-minimal','--disableTelemetry') `
  -WorkingDirectory $workspace `
  -WindowStyle Hidden `
  -RedirectStandardOutput (Join-Path $logDir 'nx-mcp-http.out.log') `
  -RedirectStandardError (Join-Path $logDir 'nx-mcp-http.err.log')
```

Validate the server:

```powershell
Get-NetTCPConnection -LocalPort 9921 -ErrorAction SilentlyContinue
Get-Content "E:\PROJECTS\New project\My app local\tmp\mcp\nx-mcp-http.out.log" -Raw
```

Expected startup log:

```text
Starting Nx MCP server
Running in standalone mode (no IDE connection)
Nx MCP server (Streamable HTTP) listening on port 9921
```

The plain Codex config command below uses stdio and exits immediately when run manually without an MCP client attached:

```powershell
npx nx mcp --workspacePath "E:\PROJECTS\New project\My app local" --no-minimal --disableTelemetry
```
