# Poliglot AI 2s Outro Prompt

Use this insert at the end of TikTok, Reels, and Shorts videos.

## Exact On-Screen Text

`poliglotAI.online`

## FLUX.2 Keyframe Prompt

Optional creative-background variant. The default workflow keeps these FLUX.2 nodes disabled so the exact URL is not redrawn by the image model.

Vertical 9:16 premium social outro keyframe for Poliglot AI.
Use the supplied Poliglot AI brand card as the visual identity lock: central 3D sapphire logo, mint neon language-wave orbit, silver/glass highlights, dark luxury tech background.
Composition: clean TikTok/Reels/Shorts end-card, centered logo reveal, exact URL zone below the logo, high contrast, crisp premium product identity.
Keep the supplied URL layout stable. Do not invent new text, subtitles, people, app screenshots, badges, icons, flags, or extra brand marks.
Mood: fast modern TikTok end screen, glossy light sweep, polished SaaS ad finish, elite but readable.

## Wan 2.1 I2V Motion Prompt

Default text-safe route: animate the exact reference card directly with Wan 2.1 I2V.

2-second vertical TikTok/Reels/Shorts outro for Poliglot AI, starting from the supplied branded end-card.
Motion: quick premium logo pop-in feel, subtle push-in, sapphire gem glint, mint orbit ribbon gently waves, soft light sweep passes behind the URL, tiny particle sparkle around the logo, smooth energetic finish.
Keep the exact logo and the exact visible text "poliglotAI.online" stable and readable. No new captions, no extra words, no people, no UI screens, no watermark.

## Negative Prompt

misspelled URL, wrong text, extra text, subtitles, captions, watermark, duplicated logo, distorted logo, deformed letters, fake letters, unreadable text, random words, people, hands, faces, flags, app screenshots, low quality, jpeg artifacts, blur, flicker, noisy text, overexposed, bad composition, black frames

## Timing

- Wan length: `33` frames.
- RIFE multiplier: `2`.
- Output fps: `32`.
- Expected duration: `2.06s`.

## Files

- Native ComfyUI workflow: `poliglot_outro_2s_native.json`.
- API prompt: `poliglot_outro_2s_wan_stage_api.json`.
- Reference card: `poliglot_outro_reference_1080x1920.png`.
- Silent audio placeholder: `poliglot_outro_2s_silent.wav` in ComfyUI input.

## Production Note

For the cleanest final brand insert, keep the generated motion subtle. The default API/native route animates the exact reference card directly. If the URL becomes less readable, use the generated clip as animated background and overlay the exact logo/URL from the reference card in the editor.
