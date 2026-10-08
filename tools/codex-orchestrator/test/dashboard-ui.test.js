import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';

const root = new URL('../', import.meta.url);
const html = readFileSync(new URL('public/index.html', root), 'utf8');
const client = readFileSync(new URL('public/app.js', root), 'utf8');

test('dashboard declares task, workspace, live status, and agent areas', () => {
  for (const id of ['task-form', 'task-input', 'cwd-input', 'agent-grid', 'activity-list', 'answer-panel', 'stop-button']) {
    assert.match(html, new RegExp(`id="${id}"`), `missing ${id}`);
  }
});

test('client escapes model text before inserting it into dashboard markup', () => {
  assert.match(client, /function escapeHtml\(value\)/);
  assert.match(client, /escapeHtml\(agent\.name\)/);
  assert.match(client, /escapeHtml\(detail\)/);
});
