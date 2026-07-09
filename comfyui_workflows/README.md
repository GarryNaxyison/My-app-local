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
- `build_poliglot_outro_workflow.py` - builds the 2-second Poliglot AI outro workflow, exact URL reference card, and silent audio placeholder.
- `poliglot_outro_2s_prompt.md` - final FLUX.2/Wan prompts for the Shorts/Reels/TikTok end insert.
- `poliglot_outro_2s_native.json` - native ComfyUI Desktop workflow for the 2-second end insert.
- `poliglot_outro_2s_wan_stage_api.json` - API prompt for queueing the same 2-second outro through a running ComfyUI server.
- `poliglot_outro_reference_1080x1920.png` - exact brand/logo/URL start card copied into ComfyUI input.
- `run_poliglot_outro_2s.mjs` - queues the outro API prompt and writes the generated MP4 path to `tmp/comfy-outro`.
- `render_poliglot_ae_outro.py` - deterministic 1080x1920 post-render that uses the more dynamic generated MP4 as motion background, masks model text, keeps the generated logo motion, and overlays only the exact `poliglotAI.online` URL with an After Effects-style float/reveal.
- `render_poliglot_tech_outro_v11.py` - kinetic 2-second 1080x1920 tech-warp outro renderer; writes the final MP4, preview frames, a ComfyUI input reference, and a matching Wan/RIFE API prompt.
- `render_poliglot_tech_outro_v11_composite.py` - final v11 compositor that uses the ComfyUI Wan output as moving source/background and overlays the exact animated logo/URL lockup.
- `render_poliglot_tech_outro_v12_comet_text.py` - refined final compositor that masks model text artifacts and assembles the URL from a comet-like letter/particle pass.
- `render_poliglot_tech_outro_v13_motion_logo_text.py` - previous compositor with post-materialization logo rotation/glint and a stronger per-character comet URL assembly.
- `render_poliglot_tech_outro_v14_coin_logo_upright_text.py` - previous compositor with upright comet URL assembly, no separate text panel, and a simple one-turn coin spin for the materialized logo.
- `render_poliglot_tech_outro_v15_clean_vector_text.py` - previous Poliglot compositor that keeps the v14 motion but resolves the URL into one clean deterministic TrueType text layer instead of a per-character final lockup.
- `render_neriva_tech_outro_v16_clean_text_slow_spin.py` - Neriva rebrand compositor that renders both `neriva.ru` and `NERIVA.RU` clean-text variants with a slower 360-degree logo coin spin.
- `render_neriva_tech_outro_v17_upper_comet_text_lockup.py` - current `NERIVA.RU` variant where the comet-flown glyphs remain as the final lockup instead of being replaced by a second full-line overlay.
- `render_neriva_luxury_minimal_outro_v18.py` - five premium minimalist `NERIVA.RU` direction candidates: noir, ivory, platinum, gallery, and emerald.
- `neriva_dzen_article_visuals.md` - nine ComfyUI-ready still-image prompts for the Yandex Zen article about why learners study English for years but do not speak, with 3 variants for each article visual block.
- `neriva_dzen_article_visuals_batch.json` - machine-readable batch manifest for the same 9 Neriva article prompts.
- `generated/neriva_dzen_article_visuals/` - generated 1536x864 PNG set for the Neriva Yandex Zen article, including all 9 variants and a contact sheet.
- `poliglot_outro_2s_tech_warp_v11_api.json` - v11 ComfyUI API prompt for the new dynamic technology direction.
- `poliglot_outro_2s_tech_warp_v11_prompt.md` - v11 prompt notes and output paths.
- `poliglot_outro_tech_warp_v11_reference_1080x1920.png` - v11 high-energy reference frame copied into ComfyUI input.
- `assets/fonts/Manrope-wght.ttf` and `assets/fonts/OFL-Manrope.txt` - bundled open-source font used for the clean readable outro URL typography.
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

## 2s Outro

Regenerate the outro workflow and input files:

```powershell
python .\comfyui_workflows\build_poliglot_outro_workflow.py
```

Open this workflow in ComfyUI Desktop:

```text
C:\Users\Admin\Documents\ComfyUI\user\default\workflows\poliglot_outro_2s_native.json
```

Or queue it through a running ComfyUI API server:

```powershell
node .\comfyui_workflows\run_poliglot_outro_2s.mjs
```

The outro uses `poliglot_outro_reference_1080x1920.png` as the exact logo and `poliglotAI.online` text lock. Keep motion subtle; if the generated video bends the URL, use the generated MP4 as the moving background and overlay the exact reference logo/text in the editor.

For the more kinetic approved direction, render the AE-style reveal pass after the ComfyUI API output exists:

```powershell
& 'C:\Users\Admin\Documents\ComfyUI\.venv\Scripts\python.exe' .\comfyui_workflows\render_poliglot_ae_outro.py
```

Current polished output with the cleaner readable typography pass:

```text
C:\Users\Admin\Documents\ComfyUI\output\outro\poliglot_outro_ae_reveal_v9_clean.mp4
```

Render the kinetic replacement for the rejected static/flat-space direction:

```powershell
& 'C:\Users\Admin\Documents\ComfyUI\.venv\Scripts\python.exe' .\comfyui_workflows\render_poliglot_tech_outro_v11.py
```

Queue `poliglot_outro_2s_tech_warp_v11_api.json` through ComfyUI, then render the refined final composite:

```powershell
& 'C:\Users\Admin\Documents\ComfyUI\.venv\Scripts\python.exe' .\comfyui_workflows\render_neriva_luxury_minimal_outro_v18.py
```

Current final output:

```text
C:\Users\Admin\Documents\ComfyUI\output\outro\neriva_luxury_minimal_v18_01_noir.mp4
C:\Users\Admin\Documents\ComfyUI\output\outro\neriva_luxury_minimal_v18_02_ivory.mp4
C:\Users\Admin\Documents\ComfyUI\output\outro\neriva_luxury_minimal_v18_03_platinum.mp4
C:\Users\Admin\Documents\ComfyUI\output\outro\neriva_luxury_minimal_v18_04_gallery.mp4
C:\Users\Admin\Documents\ComfyUI\output\outro\neriva_luxury_minimal_v18_05_emerald.mp4
```

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
