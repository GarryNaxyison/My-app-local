const byId = (id) => document.getElementById(id);
const state = { snapshot: null, localEventCount: 0, startedAt: 0, timer: null };
const providerIcons = { OpenRouter: '◈', 'Ollama · локально': '⌂', 'Codex model': '✳', 'Определяется Codex': '✳' };
const statusNames = { ready: 'Ожидает', starting: 'Запускается', running: 'Работает', done: 'Готово', error: 'Ошибка', interrupted: 'Остановлен', waiting: 'Ждёт ответа' };

function escapeHtml(value) { return String(value ?? '').replace(/[&<>"']/g, (char) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[char]); }
function shortPath(value) { const normalized = String(value || '').replaceAll('\\', '/').replace(/\/$/, ''); const parts = normalized.split('/'); return parts.length > 3 ? `…/${parts.slice(-2).join('/')}` : normalized; }
function renderProviders(agents) {
  byId('provider-list').innerHTML = agents.map((agent) => `<div class="provider-row ${agent.provider.includes('локально') ? 'local' : ''}" data-agent-id="${escapeHtml(agent.id)}"><span class="provider-icon">${providerIcons[agent.provider] || '✳'}</span><div><strong>${escapeHtml(agent.name)}</strong><small>${escapeHtml(agent.provider)}</small></div><span class="provider-state ${agent.status === 'running' ? 'running' : ''}">${escapeHtml(statusNames[agent.status] || 'Ожидает')}</span></div>`).join('');
}
function renderAgents(agents) {
  const children = agents.filter((agent) => agent.id !== 'coordinator');
  byId('agent-grid').innerHTML = children.map((agent) => `<article class="agent-card ${escapeHtml(agent.status)}" data-agent-id="${escapeHtml(agent.id)}" data-provider="${escapeHtml(agent.provider)}"><span class="agent-status" aria-label="${escapeHtml(statusNames[agent.status] || agent.status)}"></span><small class="agent-status-text">${escapeHtml(statusNames[agent.status] || agent.status)}</small><div class="agent-head"><span class="agent-avatar">${escapeHtml(agent.icon)}</span><div><strong>${escapeHtml(agent.name)}</strong><small class="agent-provider">${escapeHtml(agent.provider)}</small></div></div><div class="agent-task">${escapeHtml(agent.activeTask || agent.role)}</div></article>`).join('');
  renderProviders(agents);
  const active = children.some((agent) => agent.status === 'running');
  byId('flow-panel')?.classList.toggle('busy', active);
}
function timeLabel(value) { return new Intl.DateTimeFormat('ru-RU', { hour: '2-digit', minute: '2-digit' }).format(new Date(value)); }
function describeEvent(event) {
  if (event.type === 'task-submitted') return ['task', 'Задача принята оркестратором', event.detail];
  if (event.type === 'server-connected') return ['system', 'Codex подключён', event.detail];
  if (event.type === 'server-log') return ['system', event.label, event.detail];
  if (event.type === 'turn-started') return ['system', 'Оркестратор начал работу', 'Профиль codex-orchestrator'];
  if (event.type === 'turn-completed') return [event.status === 'completed' ? 'success' : 'error', event.status === 'completed' ? 'Задача завершена' : `Задача: ${statusNames[event.status] || event.status}`, event.error || ''];
  if (event.type === 'turn-interrupted') return ['error', 'Запрошена остановка задачи', ''];
  if (event.type === 'error') return ['error', 'Ошибка запуска', event.message];
  if (event.type === 'session-reset') return ['system', event.label, ''];
  if (event.type === 'answer-delta') return null;
  if (event.type === 'item-started' || event.type === 'item-completed') return [event.itemType === 'collabAgentToolCall' ? 'agent' : 'system', `${event.type === 'item-started' ? 'В работе' : 'Готово'} · ${event.label}`, event.detail || ''];
  return null;
}
function renderEvents(events) {
  const list = byId('activity-list');
  byId('event-count').textContent = `${events.length} ${events.length === 1 ? 'событие' : events.length < 5 ? 'события' : 'событий'}`;
  list.innerHTML = events.map((event) => {
    const described = describeEvent(event);
    if (!described) return '';
    const [kind, title, detail] = described;
    const icon = kind === 'agent' ? '↗' : kind === 'error' ? '!' : kind === 'success' ? '✓' : kind === 'task' ? '✳' : '·';
    return `<div class="activity-row ${kind}"><span class="activity-icon">${icon}</span><div><div class="activity-label">${escapeHtml(title)}</div>${detail ? `<div class="activity-detail">${escapeHtml(detail)}</div>` : ''}</div><time class="activity-clock">${timeLabel(event.at)}</time></div>`;
  }).join('');
  byId('event-time').textContent = events[0] ? `ОБНОВЛЕНО ${timeLabel(events[0].at)}` : '';
}
function render(snapshot) {
  state.snapshot = snapshot;
  const agents = snapshot.agents || [];
  renderAgents(agents);
  byId('workspace-short').textContent = shortPath(snapshot.cwd);
  byId('cwd-input').value = snapshot.cwd;
  byId('connection-label').textContent = snapshot.connected ? 'CODEX ПОДКЛЮЧЁН' : snapshot.profileReady ? 'ПРОФИЛЬ ГОТОВ' : 'ПРОФИЛЬ НЕ НАЙДЕН';
  byId('connection-label').parentElement.classList.toggle('disconnected', !snapshot.connected && !snapshot.profileReady);
  byId('session-status').textContent = ({ idle: 'Готов к работе', starting: 'Запускаю агентов', running: 'Идёт работа', completed: 'Задача выполнена', interrupted: 'Задача остановлена', error: 'Требует внимания', approval: 'Ожидает подтверждения' })[snapshot.status] || 'Готов к работе';
  const statusText = ({ idle: 'В ожидании задачи', starting: 'Запускаю оркестратор', running: 'Агенты выполняют задачу', completed: 'Задача завершена', interrupted: 'Задача остановлена', error: 'Задача завершилась с ошибкой', approval: 'Нужно подтверждение' })[snapshot.status] || 'В ожидании задачи';
  byId('run-state-label').textContent = statusText;
  byId('run-state-label').parentElement.classList.toggle('active', snapshot.status === 'running' || snapshot.status === 'starting');
  byId('flow-panel').classList.toggle('busy', snapshot.status === 'running' || snapshot.status === 'starting');
  byId('run-button').disabled = ['starting', 'running'].includes(snapshot.status);
  byId('stop-button').hidden = !['running', 'starting'].includes(snapshot.status);
  byId('task-input').value = snapshot.status === 'idle' ? byId('task-input').value : snapshot.prompt || '';
  const taskSummary = byId('task-summary');
  if (snapshot.prompt) { taskSummary.classList.remove('empty-summary'); taskSummary.innerHTML = `<span class="summary-spark">✳</span><span>${escapeHtml(snapshot.prompt)}</span>`; }
  else { taskSummary.classList.add('empty-summary'); taskSummary.innerHTML = '<span class="summary-spark">✳</span><span>Когда запустишь задачу, события появятся здесь в реальном времени.</span>'; }
  const answerPanel = byId('answer-panel');
  answerPanel.hidden = !snapshot.answer;
  byId('answer-text').textContent = snapshot.answer || '';
  byId('answer-status').textContent = snapshot.status === 'completed' ? 'Завершено' : snapshot.status === 'error' ? 'Не удалось завершить' : 'Готовит ответ';
  if (snapshot.answer) byId('answer-text').scrollTop = byId('answer-text').scrollHeight;
  const error = byId('composer-error');
  error.hidden = !snapshot.error;
  error.textContent = snapshot.error || '';
  renderEvents(snapshot.events || []);
  if (snapshot.startedAt && ['running', 'starting'].includes(snapshot.status)) {
    state.startedAt = snapshot.startedAt;
    if (!state.timer) state.timer = setInterval(() => { const elapsed = Math.floor((Date.now() - state.startedAt) / 1000); byId('session-status').textContent = `Идёт работа · ${Math.floor(elapsed / 60)}:${String(elapsed % 60).padStart(2, '0')}`; }, 1000);
  } else { clearInterval(state.timer); state.timer = null; }
}
async function post(path, body = {}) {
  const response = await fetch(path, { method: 'POST', headers: { 'content-type': 'application/json' }, body: JSON.stringify(body) });
  const result = await response.json();
  if (!response.ok) throw new Error(result.error || 'Не удалось выполнить запрос');
  return result;
}
async function refresh() { const response = await fetch('/api/status'); render(await response.json()); }

byId('task-form').addEventListener('submit', async (event) => {
  event.preventDefault();
  const prompt = byId('task-input').value.trim();
  if (!prompt) return;
  byId('composer-error').hidden = true;
  byId('run-button').disabled = true;
  byId('run-button').querySelector('span').textContent = 'Запускаю…';
  try { const result = await post('/api/start', { prompt, cwd: byId('cwd-input').value.trim() }); render(result.snapshot); }
  catch (error) { byId('composer-error').textContent = error.message; byId('composer-error').hidden = false; }
  finally { byId('run-button').disabled = false; byId('run-button').querySelector('span').textContent = 'Запустить агентов'; }
});
byId('task-input').addEventListener('keydown', (event) => { if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') byId('task-form').requestSubmit(); });
byId('cwd-edit').addEventListener('click', () => {
  const input = byId('cwd-input');
  if (input.readOnly) { input.readOnly = false; input.classList.add('editable'); input.focus(); input.select(); byId('cwd-edit').textContent = 'Готово'; }
  else { input.readOnly = true; input.classList.remove('editable'); byId('cwd-edit').textContent = 'Изменить'; }
});
byId('cwd-input').readOnly = true;
byId('stop-button').addEventListener('click', async () => { try { const result = await post('/api/stop'); render(result.snapshot); } catch (error) { byId('composer-error').textContent = error.message; byId('composer-error').hidden = false; } });
byId('reset-button').addEventListener('click', async () => { if (state.snapshot?.status === 'running' && !confirm('Остановить текущую задачу и начать новую сессию?')) return; try { const result = await post('/api/reset'); render(result.snapshot); byId('task-input').value = ''; } catch (error) { byId('composer-error').textContent = error.message; byId('composer-error').hidden = false; } });
const events = new EventSource('/api/events');
events.addEventListener('update', (message) => { try { const payload = JSON.parse(message.data); if (payload.snapshot) render(payload.snapshot); } catch {} });
events.onerror = () => { byId('connection-label').textContent = 'ПЕРЕПОДКЛЮЧАЮСЬ'; };
refresh();
