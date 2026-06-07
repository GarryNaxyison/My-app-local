# Project Agent Instructions

## Primary Workspace

Work in this repository by default: `E:\PROJECTS\New project\My app local`. Treat the parent folder `E:\PROJECTS\New project` as a container unless the user explicitly asks to work there.

## Universal MCP Startup

At the start of this or any future project, verify or start the MCP connections described in `MCP_STARTUP.md`: Figma, 21st.dev Magic, Ref, shadcn, Nx MCP, Linux/Ubuntu access when server work may be needed, and any project-specific MCP from Codex config.

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
npx -y nx-mcp@latest "C:\Users\Admin\Documents\New project" --transport http --port 9921 --no-minimal --disableTelemetry
```

Recommended background launch with logs:

```powershell
$workspace = 'C:\Users\Admin\Documents\New project'
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
Get-Content "C:\Users\Admin\Documents\New project\tmp\mcp\nx-mcp-http.out.log" -Raw
```

Expected startup log:

```text
Starting Nx MCP server
Running in standalone mode (no IDE connection)
Nx MCP server (Streamable HTTP) listening on port 9921
```

The plain Codex config command below uses stdio and exits immediately when run manually without an MCP client attached:

```powershell
npx nx mcp --workspacePath "C:\Users\Admin\Documents\New project" --no-minimal --disableTelemetry
```
