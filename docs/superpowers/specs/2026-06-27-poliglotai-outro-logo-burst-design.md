# PoliglotAI Logo Burst Outro Design

Date: 2026-06-27

## Goal

Create a reusable short vertical MP4 outro for Reels, TikTok, and YouTube Shorts. The outro appears at the end of user videos and makes the brand and destination clear in roughly two seconds.

Primary visible text:

```text
PoliglotAI.online
```

## Approved Direction

Style: TikTok energy.

Chosen concept: Logo Burst.

The animation should feel fast, bright, and social-video native, while keeping the final brand frame readable after mobile platform compression.

## Output

- Format: MP4
- Resolution: 1080 x 1920
- Aspect ratio: 9:16
- Duration: about 2.2 seconds
- Frame rate: 30 fps
- Codec: H.264
- Audio: none
- Target file: `assets/video-outro/poliglotai-outro-logo-burst.mp4`

## Brand Assets

Use the existing project logo assets, preferring:

- `web/assets/brand-logo.png`
- fallback: `web/assets/brand-logo-mini.png`

The mini logo is 512 x 512 PNG with alpha and is suitable for compositing over a dark animated background.

Use the existing brand palette:

- deep navy / near black: `#050914`, `#07111f`
- electric blue: `#2d5bff`
- cyan: `#12b8d7`
- mint highlight: `#8ff1d0`
- warm accent: `#f3b84b`

## Animation Timeline

The final video should use this structure:

- 0.00s to 0.25s: dark background with a quick cyan / blue radial flash.
- 0.25s to 0.75s: logo enters the center with fast zoom and slight bounce.
- 0.55s to 1.15s: bright ring or pulse expands behind the logo.
- 0.75s to 1.35s: `PoliglotAI.online` slides in below the logo with glow.
- 1.35s to 2.20s: stable final frame for readability and editing safety.

## Layout

Keep the composition centered and safe for mobile overlays:

- logo center around 44% of frame height;
- text center below the logo around 59% to 64% of frame height;
- keep important content away from the bottom 250 px and top 180 px;
- use strong contrast against the dark background;
- do not add explanatory copy, CTA text, handles, or extra labels.

## Implementation Shape

Build a small local generator rather than hand-editing a video file. The generator can use `ffmpeg` filters and the existing PNG logo to render the outro deterministically.

Recommended files:

- `assets/video-outro/poliglotai-outro-logo-burst.mp4`
- `assets/video-outro/poliglotai-outro-logo-burst-preview.png`
- `tools/create_poliglotai_outro.ps1`

The script should:

- create the output directory when missing;
- composite the logo over a procedural dark animated background;
- render text with a system font available on Windows;
- export a preview frame from the final stable section;
- fail clearly if `ffmpeg` or the logo file is missing.

## Verification

Use local media tooling for verification:

- `ffprobe` confirms width 1080, height 1920, duration near 2.2 seconds, and video codec H.264.
- `ffmpeg` exports a still frame from the stable final section for visual inspection.
- The final frame shows the logo and `PoliglotAI.online` clearly.
- The MP4 plays without audio and is suitable for importing into mobile video editors.

The Playwright MCP browser is not required for this artifact. In this environment, Chromium remote debugging was blocked by system policy during startup, so media verification should rely on `ffmpeg` / `ffprobe` and exported preview frames.

## Acceptance Criteria

- A reusable vertical MP4 outro exists at `assets/video-outro/poliglotai-outro-logo-burst.mp4`.
- The video uses the existing PoliglotAI logo, not a placeholder.
- The visible text is exactly `PoliglotAI.online`.
- The animation follows the approved Logo Burst direction.
- The final readable hold lasts at least 0.7 seconds.
- Verification commands and results are reported before claiming the asset is complete.

## Out Of Scope

- Sound design or music.
- Multiple style variants.
- A full advertising video.
- Server deployment.
- Changes to the public site, web app, bot, or payment flows.
