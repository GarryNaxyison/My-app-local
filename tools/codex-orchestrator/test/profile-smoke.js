import { spawn } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import { resolve } from 'node:path';

const workspace = resolve(fileURLToPath(new URL('../../..', import.meta.url)));
const codex = spawn('codex', ['--ask-for-approval', 'never', 'exec', '--profile', 'codex-orchestrator', '--json', '--sandbox', 'workspace-write', '--cd', workspace, '-'], {
  cwd: workspace,
  stdio: ['pipe', 'pipe', 'pipe'],
  windowsHide: true,
});
let output = '';
let finished = false;
const timeout = setTimeout(() => finish(new Error('Codex profile validation timed out')), 15_000);

function finish(error) {
  if (finished) return;
  finished = true;
  clearTimeout(timeout);
  if (!codex.killed) codex.kill();
  if (error) {
    process.stderr.write(`${error.message}${output ? `\n${output.slice(-1200)}` : ''}\n`);
    process.exitCode = 1;
  } else {
    process.stdout.write('Codex profile loaded for exec --json (empty input stopped before any model request).\n');
  }
}

for (const stream of [codex.stdout, codex.stderr]) {
  stream.setEncoding('utf8');
  stream.on('data', (chunk) => { output += chunk; if (output.length > 4000) output = output.slice(-4000); });
}
codex.on('error', (error) => finish(new Error(`Could not start Codex CLI: ${error.message}`)));
codex.on('exit', (code) => {
  if (finished) return;
  if (code !== 0 && /No prompt provided via stdin\./i.test(output)) finish();
  else finish(new Error(`Codex did not stop on the expected empty-input check (exit ${code ?? 'unknown'}).`));
});
codex.stdin.end();
