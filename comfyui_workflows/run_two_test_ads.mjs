import { randomUUID } from 'node:crypto';
import { mkdir, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

const COMFY_BASE = process.env.COMFY_BASE || 'http://127.0.0.1:8000';
const OUT_DIR = 'tmp/comfy-test-ads';

const NEGATIVE_PROMPT = 'minor, underage, childlike, sensual pose, fake letters, unreadable text, random subtitles, watermark, duplicated logo, distorted logo, broken UI, cluttered background, overexposed, low quality, jpeg artifacts, static frozen scene, deformed hands, deformed face, extra limbs, noisy details, blurry product, bad composition';

const TESTS = [
  {
    slug: '01_web_hero_presenter',
    reference: 'poliglot_web_product_reference.png',
    seed: 5281200101,
    fluxPrompt: 'Vertical 9:16 luxury ad keyframe for NERIVA desktop web app. A warm confident young adult woman, age 18 or older, in an elegant business suit presents the NERIVA web version on a laptop in a premium modern office studio. The laptop screen clearly shows a beautiful desktop web interface: sidebar, learning cards, progress panel, practice area, premium dashboard, voice waveform. Use the project web assets as brand identity: logo, dark premium dashboard mood, app background. Elite SaaS advertisement, product-first, clean cinematic lighting, tasteful and expensive, no fake text, no sensual pose.',
    wanPrompt: '5-second vertical luxury SaaS ad shot. Smooth cinematic push-in toward a laptop showing the NERIVA desktop web interface while a friendly young adult woman, age 18 or older, in a business suit presents it with confident calm energy. Premium office studio, polished functional UI visible: sidebar, learning cards, progress, practice, voice waveform, stable logo/product identity, no captions inside the image.',
  },
  {
    slug: '02_mobile_voice_pronunciation',
    reference: 'poliglot_web_structure_reference.png',
    seed: 5281200102,
    fluxPrompt: 'Vertical 9:16 premium ad keyframe for NERIVA mobile web voice practice. A warm confident young adult woman, age 18 or older, in a tailored business suit holds a smartphone with the NERIVA mobile web interface clearly visible. The phone screen shows a polished mobile app shell: compact header, voice practice card, microphone control, waveform, pronunciation feedback, progress chip. Luxury tech commercial, high-end studio lighting, functional product UI visible, no fake words, no random subtitles, no sensual pose.',
    wanPrompt: '5-second vertical premium product shot. Slow parallax around the presenter and smartphone, showing the NERIVA mobile web interface with voice practice, waveform motion, pronunciation feedback rings, and progress energy. Young adult business presenter, age 18 or older, friendly and professional. Elite language-learning SaaS commercial, stable product identity, sharp polished lighting, no captions inside the image.',
  },
];

function promptFor(test) {
  return {
    '1': { class_type: 'LoadImage', inputs: { image: test.reference } },
    '2': { class_type: 'ImageScale', inputs: { image: ['1', 0], upscale_method: 'lanczos', width: 480, height: 832, crop: 'center' } },
    '3': { class_type: 'UNETLoader', inputs: { unet_name: 'flux2_dev_fp8mixed.safetensors', weight_dtype: 'default' } },
    '4': { class_type: 'CLIPLoader', inputs: { clip_name: 'mistral_3_small_flux2_bf16.safetensors', type: 'flux2', device: 'default' } },
    '5': { class_type: 'VAELoader', inputs: { vae_name: 'flux2-vae.safetensors' } },
    '6': { class_type: 'CLIPTextEncode', inputs: { clip: ['4', 0], text: test.fluxPrompt } },
    '7': { class_type: 'FluxGuidance', inputs: { conditioning: ['6', 0], guidance: 4.0 } },
    '8': { class_type: 'VAEEncode', inputs: { pixels: ['2', 0], vae: ['5', 0] } },
    '9': { class_type: 'ReferenceLatent', inputs: { conditioning: ['7', 0], latent: ['8', 0] } },
    '10': { class_type: 'BasicGuider', inputs: { model: ['3', 0], conditioning: ['9', 0] } },
    '11': { class_type: 'RandomNoise', inputs: { noise_seed: test.seed } },
    '12': { class_type: 'KSamplerSelect', inputs: { sampler_name: 'euler' } },
    '13': { class_type: 'Flux2Scheduler', inputs: { steps: 16, width: 480, height: 832 } },
    '14': { class_type: 'EmptyFlux2LatentImage', inputs: { width: 480, height: 832, batch_size: 1 } },
    '15': { class_type: 'SamplerCustomAdvanced', inputs: { noise: ['11', 0], guider: ['10', 0], sampler: ['12', 0], sigmas: ['13', 0], latent_image: ['14', 0] } },
    '16': { class_type: 'VAEDecode', inputs: { samples: ['15', 0], vae: ['5', 0] } },
    '17': { class_type: 'CLIPLoader', inputs: { clip_name: 'umt5_xxl_fp8_e4m3fn_scaled.safetensors', type: 'wan', device: 'default' } },
    '18': { class_type: 'VAELoader', inputs: { vae_name: 'wan_2.1_vae.safetensors' } },
    '19': { class_type: 'UNETLoader', inputs: { unet_name: 'wan2.1_i2v_480p_14B_fp8_scaled.safetensors', weight_dtype: 'default' } },
    '20': { class_type: 'ModelSamplingSD3', inputs: { model: ['19', 0], shift: 5.0 } },
    '21': { class_type: 'CLIPTextEncode', inputs: { clip: ['17', 0], text: test.wanPrompt } },
    '22': { class_type: 'CLIPTextEncode', inputs: { clip: ['17', 0], text: NEGATIVE_PROMPT } },
    '23': { class_type: 'CLIPVisionLoader', inputs: { clip_name: 'clip_vision_h.safetensors' } },
    '24': { class_type: 'CLIPVisionEncode', inputs: { clip_vision: ['23', 0], image: ['16', 0], crop: 'center' } },
    '25': { class_type: 'WanImageToVideo', inputs: { positive: ['21', 0], negative: ['22', 0], vae: ['18', 0], width: 480, height: 832, length: 81, batch_size: 1, clip_vision_output: ['24', 0], start_image: ['16', 0] } },
    '26': { class_type: 'KSampler', inputs: { model: ['20', 0], positive: ['25', 0], negative: ['25', 1], latent_image: ['25', 2], seed: test.seed + 99, steps: 12, cfg: 5.5, sampler_name: 'uni_pc', scheduler: 'simple', denoise: 1.0 } },
    '27': { class_type: 'VAEDecode', inputs: { samples: ['26', 0], vae: ['18', 0] } },
    '28': { class_type: 'FrameInterpolationModelLoader', inputs: { model_name: 'rife_v4.26.safetensors' } },
    '29': { class_type: 'FrameInterpolate', inputs: { interp_model: ['28', 0], images: ['27', 0], multiplier: 2 } },
    '30': { class_type: 'EmptyAudio', inputs: { duration: 5.1, sample_rate: 44100, channels: 2 } },
    '31': { class_type: 'CreateVideo', inputs: { images: ['29', 0], fps: 32, audio: ['30', 0] } },
    '32': { class_type: 'SaveVideo', inputs: { video: ['31', 0], filename_prefix: `ad_shorts_tests/${test.slug}`, format: 'mp4', codec: 'auto' } },
    '33': { class_type: 'SaveImage', inputs: { images: ['16', 0], filename_prefix: `ad_shorts_tests/${test.slug}_keyframe` } },
  };
}

async function checkReady() {
  const response = await fetch(`${COMFY_BASE}/system_stats`);
  if (!response.ok) throw new Error(`ComfyUI is not ready: HTTP ${response.status}`);
}

async function queue(test) {
  const response = await fetch(`${COMFY_BASE}/prompt`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ client_id: randomUUID(), prompt: promptFor(test) }),
  });
  const data = await response.json();
  if (!response.ok || data.error) {
    throw new Error(`ComfyUI rejected ${test.slug}: ${JSON.stringify(data, null, 2)}`);
  }
  return data.prompt_id;
}

async function history(promptId) {
  const response = await fetch(`${COMFY_BASE}/history/${promptId}`);
  if (!response.ok) return null;
  const data = await response.json();
  return data[promptId] || null;
}

async function waitForOutputs(promptId, test) {
  const started = Date.now();
  while (Date.now() - started < 60 * 60 * 1000) {
    const item = await history(promptId);
    const video = item?.outputs?.['32']?.images?.[0] || null;
    const keyframe = item?.outputs?.['33']?.images?.[0] || null;
    if (video) return { video, keyframe };
    await new Promise((resolve) => setTimeout(resolve, 5000));
  }
  throw new Error(`Timed out waiting for ${test.slug}`);
}

function localOutputPath(result) {
  const subfolder = result.subfolder ? `${result.subfolder}\\` : '';
  return `C:\\Users\\Admin\\Documents\\ComfyUI\\output\\${subfolder}${result.filename}`;
}

async function main() {
  await mkdir(OUT_DIR, { recursive: true });
  await checkReady();
  const results = [];
  for (const test of TESTS) {
    console.log(`queue ${test.slug}`);
    const promptId = await queue(test);
    const outputs = await waitForOutputs(promptId, test);
    const result = {
      slug: test.slug,
      prompt_id: promptId,
      video: outputs.video ? localOutputPath(outputs.video) : null,
      keyframe: outputs.keyframe ? localOutputPath(outputs.keyframe) : null,
    };
    console.log(JSON.stringify(result, null, 2));
    results.push(result);
  }
  await writeFile(join(OUT_DIR, 'results.json'), `${JSON.stringify(results, null, 2)}\n`, 'utf8');
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
