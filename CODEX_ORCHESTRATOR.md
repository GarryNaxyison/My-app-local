# Codex Orchestrator for Windows

This local dashboard starts each submitted task with the configured `codex-orchestrator` profile and opens a browser app window. It is designed for entering a project task and seeing the main agent, active subagents, provider labels, progress events, and the final response.

## Start the app

Double-click `tools\codex-orchestrator\start-codex-orchestrator.cmd`, or run its PowerShell script directly:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File "C:\Users\Admin\Documents\New project\tools\codex-orchestrator\start-codex-orchestrator.ps1"
```

The script starts the dashboard in the background and opens it in Microsoft Edge app mode when Edge is available. Otherwise, it opens the default browser. The dashboard listens only on `127.0.0.1:4178`.

The user-requested behavior is recorded in [docs/tracking/CODEX_ORCHESTRATOR_WISHES.md](docs/tracking/CODEX_ORCHESTRATOR_WISHES.md).

Requirements:

- Node.js and the Codex CLI must be installed and available to the current Windows user.
- The Codex profile file must exist at `%USERPROFILE%\.codex\codex-orchestrator.config.toml`.
- Profile provider credentials and local model services must already be configured. The dashboard never reads or displays API keys.

For another working directory, change the **Папка проекта** field in the dashboard before starting a task. Each submission runs as a new `codex exec --profile codex-orchestrator --json` process, so each task starts a fresh Codex run. Tasks use the `workspace-write` sandbox for that directory and the CLI's non-interactive `never` approval policy; edits and commands inside the selected project can run without approval dialogs.

## What the provider labels mean

- The main orchestrator uses the provider configured by the profile. This profile selects OpenRouter and the `@preset/codex-free-coding` model alias.
- Agent cards use the agent model reported in Codex collaboration events when it is present, matched against the local agent configuration: Qwen Coder, Laguna, Free Auto, or local Qwen through Ollama.
- Agent status changes come from Codex app-server notifications. Agents that have not been started remain marked as waiting. When Codex does not report a subagent model, the card says that the provider is not identified instead of guessing.

The dashboard reads the local Codex CLI's JSONL event stream; it does not proxy model traffic through the dashboard server. Codex CLI 0.137 does not accept `--profile` on `app-server`, so this app uses `exec`, where the profile is supported. The profile's obsolete scalar agent settings were removed to let this installed CLI load its named agent configurations.

## Tests

From the repository root:

```powershell
npm --prefix tools/codex-orchestrator install
npm --prefix tools/codex-orchestrator run test:install-browser
npm --prefix tools/codex-orchestrator test
npm --prefix tools/codex-orchestrator run test:e2e
npm --prefix tools/codex-orchestrator run test:profile
```

Playwright uses a local fake Codex CLI, so it covers dashboard behavior without sending a model request. `test:profile` checks that the installed CLI loads the configured profile with empty input, also without starting a model request.
