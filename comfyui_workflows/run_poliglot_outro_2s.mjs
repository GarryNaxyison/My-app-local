import { randomUUID } from 'node:crypto';
import { mkdir, readFile, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

const COMFY_BASE = process.env.COMFY_BASE || 'http://127.0.0.1:8000';
const PROMPT_PATH = 'comfyui_workflows/poliglot_outro_2s_wan_stage_api.json';
const OUT_DIR = 'tmp/comfy-outro';

async function checkReady() {
  const response = await fetch(`${COMFY_BASE}/system_stats`);
  if (!response.ok) {
    throw new Error(`ComfyUI is not ready at ${COMFY_BASE}: HTTP ${response.status}`);
  }
}

async function queue(prompt) {
  const response = await fetch(`${COMFY_BASE}/prompt`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ client_id: randomUUID(), prompt }),
  });
  const data = await response.json();
  if (!response.ok || data.error) {
    throw new Error(`ComfyUI rejected the outro prompt: ${JSON.stringify(data, null, 2)}`);
  }
  return data.prompt_id;
}

async function history(promptId) {
  const response = await fetch(`${COMFY_BASE}/history/${promptId}`);
  if (!response.ok) {
    return null;
  }
  const data = await response.json();
  return data[promptId] || null;
}

async function waitForOutputs(promptId) {
  const started = Date.now();
  while (Date.now() - started < 2 * 60 * 60 * 1000) {
    const item = await history(promptId);
    const video = item?.outputs?.['32']?.images?.[0] || null;
    const keyframe = item?.outputs?.['33']?.images?.[0] || null;
    if (video) {
      return { video, keyframe };
    }
    await new Promise((resolve) => setTimeout(resolve, 5000));
  }
  throw new Error('Timed out waiting for the Poliglot AI outro render.');
}

function localOutputPath(result) {
  const subfolder = result.subfolder ? `${result.subfolder}\\` : '';
  return `C:\\Users\\Admin\\Documents\\ComfyUI\\output\\${subfolder}${result.filename}`;
}

async function main() {
  await mkdir(OUT_DIR, { recursive: true });
  await checkReady();

  const prompt = JSON.parse(await readFile(PROMPT_PATH, 'utf8'));
  const promptId = await queue(prompt);
  console.log(`Queued Poliglot AI outro: ${promptId}`);

  const outputs = await waitForOutputs(promptId);
  const result = {
    prompt_id: promptId,
    video: outputs.video ? localOutputPath(outputs.video) : null,
    keyframe: outputs.keyframe ? localOutputPath(outputs.keyframe) : null,
  };

  await writeFile(join(OUT_DIR, 'poliglot_outro_2s_result.json'), `${JSON.stringify(result, null, 2)}\n`, 'utf8');
  console.log(JSON.stringify(result, null, 2));
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
