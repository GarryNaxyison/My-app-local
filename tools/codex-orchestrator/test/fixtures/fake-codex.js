let prompt = '';
const threadId = 'thread-test-0001';
const turnId = 'turn-test-0001';

function event(type, fields = {}) { process.stdout.write(`${JSON.stringify({ type, ...fields })}\n`); }
function item(type, fields = {}) { return { id: `item-${Math.random().toString(16).slice(2)}`, type, ...fields }; }

process.stdin.setEncoding('utf8');
process.stdin.on('data', (chunk) => { prompt += chunk; });
process.stdin.on('end', () => {
  event('thread.started', { thread_id: threadId });
  event('turn.started', { turn_id: turnId });
  const subagent = item('collab_agent_tool_call', {
    tool: 'spawn_agent',
    sender_thread_id: threadId,
    receiver_thread_ids: ['subagent-qwen'],
    agents_states: { 'subagent-qwen': { status: 'running', message: 'Анализ задачи' } },
    model: 'qwen/qwen3.8-27b:free',
    prompt: 'Независимо проверь поставленную задачу',
    status: 'in_progress',
  });
  event('item.started', { item: subagent });
  if (prompt.toLowerCase().includes('hold')) { setInterval(() => {}, 1000); return; }
  setTimeout(() => {
    event('item.completed', { item: { ...subagent, agents_states: { 'subagent-qwen': { status: 'completed', message: 'Проверка завершена' } }, status: 'completed' } });
    event('item.completed', { item: item('agent_message', { text: 'Готово: задача выполнена, Qwen Coder проверил изменения.' }) });
    event('turn.completed', { turn_id: turnId, status: 'completed', usage: { input_tokens: 20, output_tokens: 12 } });
  }, 2200);
});
