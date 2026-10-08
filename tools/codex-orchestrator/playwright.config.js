import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './e2e',
  timeout: 15_000,
  expect: { timeout: 5_000 },
  use: { baseURL: 'http://127.0.0.1:4278', browserName: 'chromium', headless: true },
  webServer: {
    command: 'node tools/codex-orchestrator/server.js',
    cwd: '../..',
    url: 'http://127.0.0.1:4278/api/status',
    reuseExistingServer: false,
    timeout: 10_000,
    env: { CODEX_ORCHESTRATOR_FAKE: '1', CODEX_ORCHESTRATOR_PORT: '4278', CODEX_ORCHESTRATOR_WORKSPACE: process.cwd() },
  },
});
