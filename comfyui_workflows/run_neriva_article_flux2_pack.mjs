import { randomUUID } from 'node:crypto';
import { copyFile, mkdir, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

const COMFY_BASE = process.env.COMFY_BASE || 'http://127.0.0.1:8188';
const COMFY_OUTPUT = process.env.COMFY_OUTPUT || 'C:\\Users\\Admin\\Documents\\ComfyUI\\output';
const OUT_DIR = 'comfyui_workflows/generated/neriva_article_flux2_dev_v1';
const WIDTH = 1536;
const HEIGHT = 864;
const STEPS = 24;

const GLOBAL_SUFFIX = [
  'Production key visual for a Russian article about Neriva AI, 16:9 landscape.',
  'Premium cinematic photorealism, high detail, editorial composition, strong article-cover readability.',
  'Leave clean readable UI surfaces where interface text appears; no watermark, no logo artifacts, no random extra text beyond the specified UI idea.',
  'Avoid deformed hands, extra fingers, distorted face, broken phone, unreadable gibberish in major focal areas, low quality, oversharpening.',
].join(' ');

const IMAGES = [
  {
    slug: '01_hero_crystal_start',
    file: 'neriva_01_hero_crystal_start_flux2.png',
    seed: 76100101,
    title: 'Hero Image / About / Start',
    prompt: [
      'Luminous multifaceted geometric crystal structure, central core pulsating with warm and cool light, detailed abstract symbols of diverse languages, glyphs, characters, speech bubbles floating around it, clean lines, intricate patterns, soft background of digital data, cinematic lighting, shallow depth of field, detailed textures, high resolution, photorealistic, gritty aesthetic, cool-toned color palette.',
      'The crystal is the clear central source of knowledge and technology, powerful but elegant, with a balanced amount of negative space for article layout.',
      GLOBAL_SUFFIX,
    ].join(' '),
  },
  {
    slug: '02_translation_context_phone',
    file: 'neriva_02_translation_context_phone_flux2.png',
    seed: 76100102,
    title: 'Translation & Context',
    prompt: [
      "Close-up of human hand holding glowing smartphone, detailed skin texture, phone display shows a clean dual-language translation UI with two readable text rows and small flag icon zones, surrounding transparent context overlays like 'Emotional interpretation' and 'Metaphor check' as intricate UI, wet urban pavement in background, cold blue and warm amber neon light, gritty film grain, highly detailed, photorealistic.",
      "The phone screen should face the camera and be large enough for deterministic final typography overlay: Russian sentence 'Таким Запомни' on top and English sentence 'Remember this moment' below.",
      GLOBAL_SUFFIX,
    ].join(' '),
  },
  {
    slug: '03_interactive_learning_barrier',
    file: 'neriva_03_interactive_learning_barrier_flux2.png',
    seed: 76100103,
    title: 'Interactive Learning / Barrier',
    prompt: [
      'Close-up of young adult person wearing modern sleek headphones, detailed face and focused expression, interacting with a floating articulated AI language partner projected as a light form, smooth flowing abstract energy connecting them, detailed interface with progress charts and speech waves, urban city environment at night, cold cinematic tones, bokeh lights, shallow depth of field, high resolution, intricate textures, photorealistic.',
      'The image should visualize overcoming the language barrier through smooth AI conversation practice, calm confident learning energy, no child, no teenager.',
      GLOBAL_SUFFIX,
    ].join(' '),
  },
  {
    slug: '04_creative_music_analysis',
    file: 'neriva_04_creative_music_analysis_flux2.png',
    seed: 76100104,
    title: 'Creative / Music Analysis',
    prompt: [
      "Close-up of audio producer at mixing console, multiple transparent monitors, detailed faders and screen textures, one screen displays a clean lyrics analysis interface, another screen shows spectral audio waveforms and code, floating glowing musical notes, intricate technical details, moody studio lighting, warm tungsten and cool LED tones, gritty aesthetic, highly detailed, photorealistic.",
      "The central monitor should leave a crisp rectangular area for deterministic text overlay: 'Lyrics Analysis - Thus Remember' and notes on specific lines.",
      GLOBAL_SUFFIX,
    ].join(' '),
  },
  {
    slug: '05_top_language_bots',
    file: 'neriva_05_top_language_bots_flux2.png',
    seed: 76100105,
    title: 'Marketing Trojan Horse / Top Lists',
    prompt: [
      "Detailed high-tech interface, glowing transparent screens showing a ranked list of top language bots, one entry highlighted with glowing green emphasis, surrounding entries include generic icons and abstract app glyphs, complex digital ecosystem background, intricate data visualizations, cool blue lighting, deep shadows, cinematic photography, high resolution, precise details.",
      "The screen should have a large clean top-list layout for deterministic final typography overlay: 'ТОП-5 ИИ ДЛЯ ЯЗЫКОВ 2026', 'BEST LANGUAGE BOTS 2026', and highlighted 'Neriva AI'.",
      GLOBAL_SUFFIX,
    ].join(' '),
  },
];

function promptFor(image) {
  return {
    '1': { class_type: 'UNETLoader', inputs: { unet_name: 'flux2_dev_fp8mixed.safetensors', weight_dtype: 'default' } },
    '2': { class_type: 'CLIPLoader', inputs: { clip_name: 'mistral_3_small_flux2_bf16.safetensors', type: 'flux2', device: 'default' } },
    '3': { class_type: 'VAELoader', inputs: { vae_name: 'flux2-vae.safetensors' } },
    '4': { class_type: 'CLIPTextEncode', inputs: { clip: ['2', 0], text: image.prompt } },
    '5': { class_type: 'FluxGuidance', inputs: { conditioning: ['4', 0], guidance: 4.0 } },
    '6': { class_type: 'BasicGuider', inputs: { model: ['1', 0], conditioning: ['5', 0] } },
    '7': { class_type: 'RandomNoise', inputs: { noise_seed: image.seed } },
    '8': { class_type: 'KSamplerSelect', inputs: { sampler_name: 'euler' } },
    '9': { class_type: 'Flux2Scheduler', inputs: { steps: STEPS, width: WIDTH, height: HEIGHT } },
    '10': { class_type: 'EmptyFlux2LatentImage', inputs: { width: WIDTH, height: HEIGHT, batch_size: 1 } },
    '11': { class_type: 'SamplerCustomAdvanced', inputs: { noise: ['7', 0], guider: ['6', 0], sampler: ['8', 0], sigmas: ['9', 0], latent_image: ['10', 0] } },
    '12': { class_type: 'VAEDecode', inputs: { samples: ['11', 0], vae: ['3', 0] } },
    '13': { class_type: 'SaveImage', inputs: { images: ['12', 0], filename_prefix: `neriva_article_flux2_dev_v1/raw/${image.slug}` } },
  };
}

async function checkReady() {
  const response = await fetch(`${COMFY_BASE}/system_stats`);
  if (!response.ok) throw new Error(`ComfyUI is not ready: HTTP ${response.status}`);
}

async function queue(image) {
  const response = await fetch(`${COMFY_BASE}/prompt`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ client_id: randomUUID(), prompt: promptFor(image) }),
  });
  const data = await response.json();
  if (!response.ok || data.error) {
    throw new Error(`ComfyUI rejected ${image.slug}: ${JSON.stringify(data, null, 2)}`);
  }
  return data.prompt_id;
}

async function history(promptId) {
  const response = await fetch(`${COMFY_BASE}/history/${promptId}`);
  if (!response.ok) return null;
  const data = await response.json();
  return data[promptId] || null;
}

async function waitForImage(promptId, image) {
  const started = Date.now();
  while (Date.now() - started < 60 * 60 * 1000) {
    const item = await history(promptId);
    const output = item?.outputs?.['13']?.images?.[0] || null;
    if (output) return output;
    await new Promise((resolve) => setTimeout(resolve, 5000));
  }
  throw new Error(`Timed out waiting for ${image.slug}`);
}

function localOutputPath(result) {
  const subfolder = result.subfolder ? `${result.subfolder}\\` : '';
  return `${COMFY_OUTPUT}\\${subfolder}${result.filename}`;
}

async function main() {
  await mkdir(OUT_DIR, { recursive: true });
  await checkReady();

  const workflowSample = promptFor(IMAGES[0]);
  await writeFile(join(OUT_DIR, 'flux2_dev_api_workflow_sample.json'), `${JSON.stringify(workflowSample, null, 2)}\n`, 'utf8');
  await writeFile(join(OUT_DIR, 'prompts.json'), `${JSON.stringify({ width: WIDTH, height: HEIGHT, steps: STEPS, images: IMAGES }, null, 2)}\n`, 'utf8');

  const results = [];
  for (const image of IMAGES) {
    console.log(`[queue] ${image.slug}`);
    const promptId = await queue(image);
    console.log(`[wait] ${image.slug} ${promptId}`);
    const output = await waitForImage(promptId, image);
    const source = localOutputPath(output);
    const destination = join(OUT_DIR, image.file);
    await copyFile(source, destination);
    const result = {
      slug: image.slug,
      title: image.title,
      seed: image.seed,
      prompt_id: promptId,
      comfy_output: source,
      project_output: destination,
      prompt: image.prompt,
    };
    console.log(`[done] ${image.slug} -> ${destination}`);
    results.push(result);
    await writeFile(join(OUT_DIR, 'manifest.json'), `${JSON.stringify({ generated_at: new Date().toISOString(), width: WIDTH, height: HEIGHT, steps: STEPS, results }, null, 2)}\n`, 'utf8');
  }
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
