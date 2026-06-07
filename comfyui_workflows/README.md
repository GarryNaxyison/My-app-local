# ComfyUI Ad Shorts Workflow

This folder contains the disk-conscious ComfyUI chain for Poliglot AI ads in Shorts, Reels, and TikTok format.

## Files

- `ad_shorts_reels_tiktok_native.json` - main ComfyUI workflow copied to `C:\Users\Admin\Documents\ComfyUI\user\default\workflows`.
- `ad_shorts_reels_tiktok_wan_stage_api.json` - API prompt for only the Wan/RIFE/video-combine stage, useful if a keyframe already exists.
- `poliglot_elite_20s_ad_prompt.md` - English 20-second voiceover, XTTS direction, and 4-shot visual plan.
- `poliglot_2x5s_test_ad_variants.md` - first two 5-second test concepts: desktop web app and mobile voice/pronunciation.
- `test_variants/` - ComfyUI workflow variants for the first two tests.
- `run_two_test_ads.mjs` - queues the first two visual test clips through a running ComfyUI API server.
- `build_ad_shorts_workflow.py` - regenerates the workflow and placeholder ComfyUI input files.
- `install_ad_short_models.ps1` - downloads only missing compact video/interpolation models.
- `model_manifest.json` - model inventory and disk budget.

## Quick Start

1. Install missing video models:

```powershell
PowerShell -ExecutionPolicy Bypass -File .\comfyui_workflows\install_ad_short_models.ps1
```

2. Open ComfyUI Desktop and load:

```text
C:\Users\Admin\Documents\ComfyUI\user\default\workflows\ad_shorts_reels_tiktok_native.json
```

3. Replace these placeholder inputs in `C:\Users\Admin\Documents\ComfyUI\input` when producing a real ad:

- `ad_product_reference.png` - product/logo/app visual lock.
- `ad_structure_reference.png` - composition/pose/layout reference.
- `ad_voiceover.wav` - XTTS-generated voiceover.

4. Run the main 480x832 branch first. Run the 1080x1920 upscale branch only when the creative is approved.

The workflow includes a `20s English ad script for XTTS` text node. For best quality, make the final 20-second ad as four approved 5-second Wan clips and stitch them, rather than one very long Wan pass.

The first two tests intentionally cover both interface styles:

- desktop web app on a laptop with presenter;
- mobile web voice/pronunciation screen on a smartphone with presenter.

## Chain

```text
Product reference
-> FLUX.2 Dev keyframe
-> CLIP Vision image lock
-> Wan 2.1 I2V 480p FP8 motion
-> VAE Decode
-> RIFE 4.26 frame interpolation
-> Video Combine with XTTS voice file
-> Save MP4
-> optional RealESRGAN 1080x1920 upscale
-> Whisper subtitles after export
```

## Why This Version

The official Wan 2.1 docs prefer fp16 for best quality, but this machine already stores heavy Flux.2 files. The default workflow uses `wan2.1_i2v_480p_14B_fp8_scaled.safetensors` to keep the new download set around 25 GB instead of adding another 32 GB fp16 model. The upscale branch lets you generate motion at a manageable size and only upscale final picks.
