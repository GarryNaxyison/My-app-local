# Neriva Dzen Article Visuals - Detailer V2

Generated on 2026-07-09 through local ComfyUI API at `127.0.0.1:8188`.

This is the stronger test workflow set. It uses a base SDXL photo generation pass, then a face detailer pass and a hand detailer pass. Text rendering is intentionally kept out of diffusion prompts; use deterministic text overlay nodes for Russian or English copy.

## Installed Nodes And Models

- Custom nodes:
  - `ComfyUI-Impact-Pack`
  - `ComfyUI-Impact-Subpack`
  - `ComfyUI_essentials`
  - `ComfyUI-TextOverlay`
- Detector/detailer models:
  - `models/ultralytics/bbox/face_yolov8m.pt`
  - `models/ultralytics/bbox/face_yolov8n.pt`
  - `models/ultralytics/bbox/hand_yolov8n.pt`
  - `models/sams/sam_vit_b_01ec64.pth`

## Workflow

API workflow sample: `comfyui_workflows/neriva_dzen_article_visuals_detailer_v2_api.json`

Prompt/seed pack: `comfyui_workflows/neriva_dzen_article_visuals_detailer_v2_prompts.json`

Chain:

```text
CheckpointLoaderSimple
-> CLIPTextEncode positive/negative
-> EmptyLatentImage 1536x864
-> KSampler
-> VAEDecode
-> FaceDetailer with face_yolov8m + SAM ViT-B
-> BboxDetectorSEGS with hand_yolov8n
-> DetailerForEach hand inpaint pass
-> SaveImage
```

## Generation Settings

- Checkpoint: `Juggernaut-XL_v9_RunDiffusionPhoto_v2.safetensors`
- Size: `1536x864`
- Base steps: `34`
- Base CFG: `7.2`
- Sampler: `dpmpp_2m_sde_gpu`
- Scheduler: `karras`
- Face detailer: `18` steps, denoise `0.25`
- Hand detailer: `18` steps, denoise `0.30`

## Text Policy

Do not ask the diffusion model to render exact Russian or English text. For headlines, domains, CTA, or labels, use:

- `Text Overlay`
- `DrawText+`

Use a real TTF font with Cyrillic support, for example the bundled Manrope font or a Windows system font. Diffusion prompts should keep screens/signs abstract: waves, circles, HUD panels, and glow only.

## Files

| ID | File | Seed | Article role |
| --- | --- | --- | --- |
| 1A | `neriva_dzen_01_cover_neon_desk_v01.png` | `2501000` | Cover variant: bright AI phone vs old textbooks |
| 1B | `neriva_dzen_01_cover_classroom_signal_v02.png` | `2501177` | Cover variant: classroom with tech interface |
| 1C | `neriva_dzen_01_cover_editorial_split_v03.png` | `2501354` | Cover variant: top-down editorial split |
| 2A | `neriva_dzen_02_old_methods_street_dictionary_v01.png` | `2501531` | Old-methods block: street fear barrier |
| 2B | `neriva_dzen_02_old_methods_platform_silence_v02.png` | `2501708` | Old-methods block: metro/platform silence |
| 2C | `neriva_dzen_02_old_methods_red_marks_v03.png` | `2501885` | Old-methods block: school correction pressure |
| 3A | `neriva_dzen_03_future_holographic_coach_v01.png` | `2502062` | Future/Neriva block: confident AI-assisted city walk |
| 3B | `neriva_dzen_03_future_phone_conversation_v02.png` | `2502239` | Future/Neriva block: accessible smartphone practice |
| 3C | `neriva_dzen_03_future_crossing_fluency_v03.png` | `2708331` | Future/Neriva block: confident crosswalk portrait |

`neriva_dzen_article_visuals_detailer_v2_contact_sheet.png` is a review sheet only; use the individual PNG files for publication.
