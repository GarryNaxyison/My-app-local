import { spawn } from 'node:child_process';
import { createServer } from 'node:http';
import { existsSync, readFileSync, statSync } from 'node:fs';
import { resolve, extname, join, normalize } from 'node:path';
import { fileURLToPath } from 'node:url';
import { randomUUID } from 'node:crypto';

const appRoot = fileURLToPath(new URL('.', import.meta.url));
const projectRoot = resolve(process.env.CODEX_ORCHESTRATOR_WORKSPACE || 'C:\\Users\\Admin\\Documents\\New project');
const codexHome = process.env.CODEX_HOME || join(process.env.USERPROFILE || '', '.codex');
const port = Number(process.env.CODEX_ORCHESTRATOR_PORT || 4178);
const codexCommand = process.env.CODEX_CLI_PATH || 'codex';
const profile = 'codex-orchestrator';
const agentCatalog = [
  { id: 'coordinator', name: 'Оркестратор', role: 'Главный агент', provider: 'OpenRouter', model: '@preset/codex-free-coding', icon: '✳' },
  { id: 'free_qwen_coder', name: 'Qwen Coder', role: 'Реализация и анализ', provider: 'OpenRouter', model: 'qwen/qwen3.8-27b:free', icon: 'Q' },
  { id: 'free_laguna', name: 'Laguna', role: 'Независимая проверка', provider: 'OpenRouter', model: 'poolside/laguna-s-2.1:free', icon: 'L' },
  { id: 'free_auto', name: 'Free Auto', role: 'Свободный вспомогательный агент', provider: 'OpenRouter', model: 'openrouter/free', icon: 'A' },
  { id: 'local_qwen', name: 'Qwen Local', role: 'Локальный анализ', provider: 'Ollama · локально', model: 'qwen3.6:35b-a3b-coding-32k', icon: '⌂' },
];

const state = {
  profile,
  profileReady: existsSync(join(codexHome, `${profile}.config.toml`)),
  connected: false,
  status: 'idle',
  threadId: null,
  turnId: null,
  cwd: projectRoot,
  prompt: '',
  answer: '',
  startedAt: null,
  error: null,
  agents: Object.fromEntries(agentCatalog.map((agent) => [agent.id, { ...agent, status: 'ready', activeTask: '', lastActivity: null }])),
  events: [],
};

let child = null;
let stdoutBuffer = '';
const clients = new Set();

function snapshot() {
  return {
    profile: state.profile,
    profileReady: state.profileReady,
    connected: state.connected,
    status: state.status,
    threadId: state.threadId,
    turnId: state.turnId,
    cwd: state.cwd,
    prompt: state.prompt,
    answer: state.answer,
    startedAt: state.startedAt,
    error: state.error,
    agents: Object.values(state.agents),
    events: state.events,
  };
}

function publish(type, data = {}) {
  const event = { id: randomUUID(), type, at: new Date().toISOString(), ...data };
  state.events.unshift(event);
  state.events = state.events.slice(0, 120);
  for (const client of clients) client.write(`event: update\ndata: ${JSON.stringify({ event, snapshot: snapshot() })}\n\n`);
  return event;
}

function labelAgent(model) {
  if (!model) return null;
  const normalizedModel = String(model).toLowerCase();
  if (normalizedModel.includes('qwen3.6:35b')) return state.agents.local_qwen;
  if (normalizedModel.includes('qwen3.8-27b')) return state.agents.free_qwen_coder;
  if (normalizedModel.includes('laguna')) return state.agents.free_laguna;
  if (normalizedModel === 'openrouter/free') return state.agents.free_auto;
  return null;
}

function setAgentStatus(agent, status, task = '') {
  if (!agent) return;
  agent.status = status;
  agent.activeTask = task;
  agent.lastActivity = Date.now();
}

function field(item, snakeName, camelName) { return item?.[snakeName] ?? item?.[camelName]; }
function itemKind(item) { return String(item?.type || '').replaceAll('_', '').toLowerCase(); }

function eventText(item) {
  switch (itemKind(item)) {
    case 'commandexecution': return 'Агент выполняет команду в рабочей папке';
    case 'filechange': return 'Агент изменяет файлы проекта';
    case 'websearch': return 'Агент ищет информацию в интернете';
    case 'mcptoolcall': return `Вызван инструмент MCP: ${item.server || 'инструмент'}`;
    case 'collabagenttoolcall': return ['spawnagent', 'spawn_agent'].includes(item.tool) ? 'Запуск субагента' : ['wait', 'waitagents', 'wait_agents'].includes(item.tool) ? 'Ожидание ответа субагента' : 'Обновление команды агентов';
    case 'agentmessage': return 'Получено сообщение агента';
    case 'reasoning': return 'Агент анализирует задачу';
    case 'plan': return 'Агент составил план';
    default: return 'Обновление статуса';
  }
}

function applyCollaborators(item, completed) {
  const receivers = field(item, 'receiver_thread_ids', 'receiverThreadIds') || [];
  const agentStates = field(item, 'agents_states', 'agentsStates') || {};
  const modelAgent = labelAgent(item.model);
  const promptText = item.prompt ? String(item.prompt).replace(/\s+/g, ' ').slice(0, 110) : '';
  for (const receiverId of receivers) {
    const agent = modelAgent || state.agents[`thread:${receiverId}`] || {
      id: `thread:${receiverId}`,
      name: item.model ? 'Субагент' : 'Агент Codex',
      role: 'Задача оркестратора',
      provider: item.model ? 'Codex model' : 'Определяется Codex',
      model: item.model || 'модель по конфигурации',
      icon: '✳',
    };
    if (!modelAgent) state.agents[`thread:${receiverId}`] = agent;
    const agentState = agentStates[receiverId] || {};
    const reported = agentState.status;
    const nextStatus = completed
      ? (reported === 'errored' ? 'error' : reported === 'completed' ? 'done' : 'running')
      : 'running';
    setAgentStatus(agent, nextStatus, promptText || agentState.message || 'Получил подзадачу');
  }
  for (const [receiverId, agentState] of Object.entries(agentStates)) {
    const agent = state.agents[`thread:${receiverId}`] || Object.values(state.agents).find((entry) => entry.model === item.model);
    if (agent) setAgentStatus(agent, agentState.status === 'running' ? 'running' : agentState.status === 'errored' ? 'error' : agentState.status === 'completed' ? 'done' : 'waiting', agentState.message || agent.activeTask);
  }
}

function handleExecEvent(event) {
  if (!event || typeof event !== 'object') return;
  const type = String(event.type || '').toLowerCase();
  if (type === 'thread.started') {
    state.threadId = event.thread_id || event.threadId || state.threadId;
    publish('thread-started', { threadId: state.threadId });
    return;
  }
  if (type === 'turn.started') {
    state.turnId = event.turn_id || event.turnId || state.turnId;
    state.status = 'running';
    state.agents.coordinator.status = 'running';
    publish('turn-started', { turnId: state.turnId });
    return;
  }
  if (type === 'item.started' || type === 'item.completed' || type === 'item.updated') {
    const item = event.item || {};
    const kind = itemKind(item);
    const completed = type === 'item.completed';
    if (kind === 'collabagenttoolcall') applyCollaborators(item, completed);
    if (kind === 'agentmessage' && typeof item.text === 'string') {
      state.answer = item.text;
      publish('answer-delta', { text: item.text });
    }
    const detail = kind === 'agentmessage' ? String(item.text || '').slice(0, 240) : '';
    publish(completed ? 'item-completed' : 'item-started', { label: eventText(item), detail, itemType: item.type || 'unknown' });
    return;
  }
  if (type === 'item.delta' || type === 'item.agent_message.delta' || type === 'item/agentMessage/delta') {
    const delta = event.delta || event.text || '';
    state.answer += delta;
    publish('answer-delta', { text: delta });
    return;
  }
  if (type === 'turn.completed') {
    state.status = event.status === 'interrupted' ? 'interrupted' : 'completed';
    state.agents.coordinator.status = state.status === 'completed' ? 'done' : 'interrupted';
    for (const agent of Object.values(state.agents)) if (agent.id !== 'coordinator' && agent.status === 'running') setAgentStatus(agent, 'done', 'Задача завершена');
    publish('turn-completed', { turnId: event.turn_id || event.turnId || state.turnId, status: state.status });
    return;
  }
  if (type === 'turn.failed' || type === 'error') {
    state.status = 'error';
    state.error = event.message || event.error?.message || event.error || 'Codex завершил задачу с ошибкой';
    state.agents.coordinator.status = 'error';
    publish('error', { message: state.error });
  }
}

function isLoopbackRequest(request) {
  const address = request.socket.remoteAddress || '';
  return address === '127.0.0.1' || address === '::1' || address === '::ffff:127.0.0.1';
}

function respond(response, status, body) {
  response.writeHead(status, { 'content-type': 'application/json; charset=utf-8', 'cache-control': 'no-store' });
  response.end(JSON.stringify(body));
}

async function readJson(request) {
  let body = '';
  for await (const chunk of request) {
    body += chunk;
    if (body.length > 80_000) throw new Error('Запрос слишком большой');
  }
  return JSON.parse(body || '{}');
}

function serveAsset(pathname, response) {
  const decoded = decodeURIComponent(pathname === '/' ? '/index.html' : pathname);
  const file = resolve(join(appRoot, 'public'), `.${normalize(decoded)}`);
  const publicRoot = resolve(appRoot, 'public');
  if (!file.startsWith(`${publicRoot}\\`) && file !== publicRoot) return response.writeHead(403).end();
  if (!existsSync(file) || !statSync(file).isFile()) return response.writeHead(404).end('Not found');
  const mime = { '.html': 'text/html; charset=utf-8', '.css': 'text/css; charset=utf-8', '.js': 'text/javascript; charset=utf-8', '.svg': 'image/svg+xml' }[extname(file)] || 'application/octet-stream';
  response.writeHead(200, { 'content-type': mime, 'cache-control': 'no-cache', 'x-content-type-options': 'nosniff' });
  response.end(readFileSync(file));
}

function startCodex(prompt, cwd) {
  stdoutBuffer = '';
  const args = ['--ask-for-approval', 'never', 'exec', '--profile', profile, '--json', '--sandbox', 'workspace-write', '--cd', cwd, '-'];
  const fakeServer = process.env.CODEX_ORCHESTRATOR_FAKE === '1';
  const processChild = fakeServer
    ? spawn(process.execPath, [join(appRoot, 'test', 'fixtures', 'fake-codex.js')], { cwd, stdio: ['pipe', 'pipe', 'pipe'], windowsHide: true })
    : spawn(codexCommand, args, { cwd, stdio: ['pipe', 'pipe', 'pipe'], windowsHide: true });
  child = processChild;
  processChild.stdout.setEncoding('utf8');
  processChild.stdout.on('data', (chunk) => {
    stdoutBuffer += chunk;
    let newline;
    while ((newline = stdoutBuffer.indexOf('\n')) >= 0) {
      const line = stdoutBuffer.slice(0, newline).trim();
      stdoutBuffer = stdoutBuffer.slice(newline + 1);
      if (!line) continue;
      try { handleExecEvent(JSON.parse(line)); } catch (error) { publish('server-log', { label: 'Не удалось прочитать событие Codex', detail: error.message }); }
    }
  });
  processChild.stderr.setEncoding('utf8');
  processChild.stderr.on('data', (chunk) => {
    const detail = String(chunk).trim();
    if (detail) publish('server-log', { label: 'Codex CLI', detail: detail.slice(0, 280) });
  });
  processChild.on('error', (error) => {
    if (child !== processChild) return;
    child = null;
    state.connected = false;
    state.status = 'error';
    state.error = `Не удалось запустить Codex CLI: ${error.message}`;
    state.agents.coordinator.status = 'error';
    publish('error', { message: state.error });
  });
  processChild.on('exit', (code) => {
    if (child !== processChild) return;
    child = null;
    state.connected = false;
    if (stdoutBuffer.trim()) {
      try { handleExecEvent(JSON.parse(stdoutBuffer)); } catch {}
      stdoutBuffer = '';
    }
    if (state.status === 'running' || state.status === 'starting') {
      state.status = code === 0 ? 'completed' : 'error';
      state.agents.coordinator.status = code === 0 ? 'done' : 'error';
      if (code !== 0) state.error ||= `Codex CLI завершился с кодом ${code ?? 'неизвестен'}`;
      publish(code === 0 ? 'turn-completed' : 'error', code === 0 ? { status: state.status } : { message: state.error });
    }
    publish('runner-exit', { label: 'Запуск Codex завершён', detail: `Код выхода: ${code ?? 'неизвестен'}` });
  });
  processChild.stdin.end(prompt);
}

const server = createServer(async (request, response) => {
  if (!isLoopbackRequest(request)) return respond(response, 403, { error: 'Доступ разрешён только с этого компьютера' });
  const origin = request.headers.origin;
  if (origin && !/^https?:\/\/(127\.0\.0\.1|localhost)(:\d+)?$/i.test(origin)) return respond(response, 403, { error: 'Недопустимый источник запроса' });
  const url = new URL(request.url, `http://127.0.0.1:${port}`);

  if (request.method === 'GET' && url.pathname === '/api/status') return respond(response, 200, snapshot());
  if (request.method === 'GET' && url.pathname === '/api/events') {
    response.writeHead(200, { 'content-type': 'text/event-stream; charset=utf-8', 'cache-control': 'no-cache', connection: 'keep-alive', 'x-accel-buffering': 'no' });
    response.write(`event: update\ndata: ${JSON.stringify({ snapshot: snapshot() })}\n\n`);
    clients.add(response);
    request.on('close', () => clients.delete(response));
    return;
  }
  if (request.method === 'POST' && url.pathname === '/api/start') {
    let body;
    try { body = await readJson(request); } catch (error) { return respond(response, 400, { error: error.message }); }
    const prompt = String(body.prompt || '').trim();
    const cwd = resolve(String(body.cwd || projectRoot).trim());
    if (!prompt) return respond(response, 400, { error: 'Напиши задачу для оркестратора' });
    if (prompt.length > 20_000) return respond(response, 400, { error: 'Задача должна быть короче 20 000 символов' });
    if (!existsSync(cwd) || !statSync(cwd).isDirectory()) return respond(response, 400, { error: 'Папка проекта не найдена' });
    if (state.status === 'running' || state.status === 'starting') return respond(response, 409, { error: 'Сначала дождись текущей задачи или останови её' });
    state.cwd = cwd;
    state.prompt = prompt;
    state.answer = '';
    state.error = null;
    state.threadId = null;
    state.turnId = null;
    state.startedAt = Date.now();
    state.status = 'starting';
    state.agents.coordinator.status = 'starting';
    for (const agent of Object.values(state.agents)) if (agent.id !== 'coordinator') setAgentStatus(agent, 'ready', '');
    publish('task-submitted', { label: 'Задача передана оркестратору', detail: prompt.slice(0, 180), cwd });
    try {
      startCodex(prompt, cwd);
      state.connected = true;
      state.status = 'running';
      return respond(response, 202, { ok: true, snapshot: snapshot() });
    } catch (error) {
      state.status = 'error';
      state.error = error.message;
      state.agents.coordinator.status = 'error';
      publish('error', { message: state.error });
      return respond(response, 500, { error: error.message, snapshot: snapshot() });
    }
  }
  if (request.method === 'POST' && url.pathname === '/api/stop') {
    if (!child || !['running', 'starting'].includes(state.status)) return respond(response, 409, { error: 'Активную задачу остановить не получилось' });
    const runningProcess = child;
    child = null;
    runningProcess.kill();
    state.connected = false;
    state.status = 'interrupted';
    state.agents.coordinator.status = 'interrupted';
    for (const agent of Object.values(state.agents)) if (agent.status === 'running') setAgentStatus(agent, 'interrupted', 'Остановлен пользователем');
    publish('turn-interrupted', { label: 'Запуск Codex остановлен' });
    return respond(response, 200, { ok: true, snapshot: snapshot() });
  }
  if (request.method === 'POST' && url.pathname === '/api/reset') {
    const runningProcess = child;
    child = null;
    if (runningProcess) runningProcess.kill();
    state.connected = false;
    state.threadId = null;
    state.turnId = null;
    state.status = 'idle';
    state.answer = '';
    state.error = null;
    state.startedAt = null;
    for (const agent of Object.values(state.agents)) setAgentStatus(agent, 'ready', '');
    publish('session-reset', { label: 'Сессия Codex сброшена' });
    return respond(response, 200, { ok: true, snapshot: snapshot() });
  }
  if (request.method === 'GET') return serveAsset(url.pathname, response);
  return respond(response, 404, { error: 'Маршрут не найден' });
});

server.listen(port, '127.0.0.1', () => console.log(`Codex Orchestrator listening on http://127.0.0.1:${port}`));
