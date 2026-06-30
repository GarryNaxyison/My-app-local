# MCP Startup

Use this checklist at the start of work when MCP context may matter. The primary workspace is:

```text
E:\PROJECTS\New project\My app local
```

## Required Checks

- Figma MCP: use when working from Figma links, inspecting designs, or writing screens back to Figma.
- 21st.dev Magic MCP: use when generating UI components or looking up component patterns.
- Ref MCP: use when framework/library documentation is needed.
- shadcn MCP: use when adding or checking shadcn/ui components.
- Nx MCP: use when project graph, affected targets, or Nx workspace context is needed.
- Open Design MCP: use for design projects, artifacts, project file access, plugin/design-system context, and handoff from Open Design.
- Linux/Ubuntu access: use when server deployment or production verification is needed.
- Project-specific Codex MCP: use any additional MCP configured for this project/session.

If an MCP is unavailable, continue with the best local fallback and report the missing connection only when it affects the task.

## Nx MCP

Start Nx MCP as a standalone Streamable HTTP server when a persistent endpoint is needed:

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

Validate:

```powershell
Get-NetTCPConnection -LocalPort 9921 -ErrorAction SilentlyContinue
Get-Content "E:\PROJECTS\New project\My app local\tmp\mcp\nx-mcp-http.out.log" -Raw
Get-Content "E:\PROJECTS\New project\My app local\tmp\mcp\nx-mcp-http.err.log" -Raw
```

Expected startup log:

```text
Starting Nx MCP server
Running in standalone mode (no IDE connection)
Nx MCP server (Streamable HTTP) listening on port 9921
```

Manual stdio command, only useful when an MCP client attaches to it:

```powershell
npx nx mcp --workspacePath "E:\PROJECTS\New project\My app local" --no-minimal --disableTelemetry
```

## Open Design MCP

Open Design is installed at:

```text
C:\Users\Admin\.codex\vendor_imports\open-design
```

Codex MCP server name:

```text
open-design
```

Start the Open Design daemon when design/artifact work needs it:

```powershell
$odRoot = 'C:\Users\Admin\.codex\vendor_imports\open-design'
$logDir = 'C:\Users\Admin\.codex\tmp'
New-Item -ItemType Directory -Force -Path $logDir | Out-Null
Start-Process -FilePath 'C:\Program Files\nodejs\node.exe' `
  -ArgumentList @((Join-Path $odRoot 'apps\daemon\dist\cli.js'),'daemon','start','--headless','--serve-web','--port','7456','--host','127.0.0.1') `
  -WorkingDirectory $odRoot `
  -WindowStyle Hidden `
  -RedirectStandardOutput (Join-Path $logDir 'open-design-daemon.out.log') `
  -RedirectStandardError (Join-Path $logDir 'open-design-daemon.err.log')
```

Validate:

```powershell
Get-NetTCPConnection -LocalPort 7456 -ErrorAction SilentlyContinue
(Invoke-WebRequest -UseBasicParsing -Uri 'http://127.0.0.1:7456/api/health').Content
& 'C:\Program Files\nodejs\node.exe' 'C:\Users\Admin\.codex\vendor_imports\open-design\apps\daemon\dist\cli.js' daemon status --json --daemon-url http://127.0.0.1:7456
```

Expected health response includes:

```json
{"ok":true,"version":"0.12.1"}
```

The Codex MCP config uses absolute paths:

```text
command = 'C:\Program Files\nodejs\node.exe'
args = ['C:\Users\Admin\.codex\vendor_imports\open-design\apps\daemon\dist\cli.js', "mcp", "--daemon-url", "http://127.0.0.1:7456"]
```

Do not rely on a bare `od` command in Git Bash: `/usr/bin/od` is the system octal-dump tool and can shadow Open Design's CLI.

## Browser MCP

Use Playwright/browser MCP for:

- checking UI behavior in a real browser;
- taking screenshots;
- inspecting console errors;
- validating responsive layouts.

## Production Snapshot MCP

Use the SQLite production snapshot MCP only for read-only production data inspection. Refresh the snapshot before relying on current data.

## Server Access

When deployment is requested, build the Linux binary locally or on the server as appropriate, upload changed deployment blocks, restart services, and verify the running service. Follow the project deployment instructions in `AGENTS.md`.
