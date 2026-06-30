# MCP On-Demand Usage

Use this file as a trigger map, not as a startup checklist. Do not verify, start, or call MCP servers by default. Prefer local shell/file inspection first, then use MCP only when the current task directly needs the capability or local tools are insufficient. The primary workspace is:

```text
E:\PROJECTS\New project\My app local
```

## Trigger Map

- Headroom MCP: use for large outputs/files/logs/search results over ~200 lines or ~20 KB, large JSON/CSS/TSX payloads, and repeated diagnostics.
- Browser/Playwright MCP: use for UI behavior, screenshots, console checks, responsive verification, or local web testing.
- GitHub MCP: use for PRs, issues, reviews, releases, remote repository metadata, or explicit GitHub actions.
- Nx MCP: use when project graph, affected targets, workspace/project config, or Nx-specific documentation is needed.
- MarkItDown MCP: use to convert trusted local files, URLs, PDFs, or data URIs into Markdown for LLM-readable context.
- Figma MCP: use when working from Figma links, inspecting designs, or writing screens back to Figma.
- 21st.dev Magic MCP: use when generating UI components or looking up component patterns.
- Ref MCP: use when framework/library documentation is needed and local docs are insufficient.
- shadcn MCP: use when adding or checking shadcn/ui components.
- Linux/Ubuntu access: use when server deployment or production verification is needed.
- Project-specific Codex MCP: use additional MCP configured for this project/session only when the task needs it.

Do not start persistent MCP daemons unless the current task needs them. If an MCP is unavailable, continue with the best local fallback and report the missing connection only when it affects the task.

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

## MarkItDown MCP

MarkItDown MCP is installed in an isolated Python venv at:

```text
C:\Users\Admin\.codex\vendor_imports\markitdown-mcp-venv
```

Codex MCP server name:

```text
markitdown
```

The Codex MCP config uses:

```text
command = 'C:\Users\Admin\.codex\vendor_imports\markitdown-mcp-venv\Scripts\python.exe'
args = ["-m", "markitdown_mcp"]
```

Validate with an MCP client or after Codex restart by listing the `markitdown` tools. The expected tool is:

```text
convert_to_markdown
```

The tool accepts a `uri` value for `file:`, `http:`, `https:`, or `data:` resources and converts it to Markdown. Treat inputs as trusted or sanitized: MarkItDown runs with the current user's privileges and can read any resource the process can access.

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
