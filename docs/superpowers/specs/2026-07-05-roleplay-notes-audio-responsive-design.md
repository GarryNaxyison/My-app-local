# Roleplay notes, audio, and desktop scaling design

## Scope

Fix four user-visible issues in the NERIVA web app:

- Roleplay scenario cards start too low on desktop.
- Quick save to notes offers too few phrases, especially only the first sentence from roleplay replies.
- Some vocabulary audio buttons fail, and generated audio for words such as `water` can reuse a bad cached clip.
- Large desktop monitors, including 2K 27-inch screens at 100% browser zoom, need a larger, denser but readable web UI.

## Design

Use the existing React/CSS structure. Do not introduce a new layout system.

Roleplay scenario selection should use a compact top-aligned panel: reduce empty vertical space, keep the panel header and card grid close to the top, and preserve mobile scroll behavior.

Quick save candidates should include structured fields plus useful sentences parsed from assistant message bodies. The extractor must dedupe candidates, ignore labels and interface-only lines, and keep the existing six-candidate cap.

Vocabulary audio should remain protected for learned words, but the frontend may fall back to text-to-speech for the displayed word when the learned-word endpoint fails. Backend audio cache keys should include enough requested audio text/model/voice information to avoid stale bad clips.

Desktop scaling should use bounded responsive tokens for large screens: increase app font size, panel padding, toolbar dimensions, and main content density above 1920px without changing mobile breakpoints.

## Verification

Run Go tests for audio-related helpers and web API behavior, run the React build, and inspect the roleplay/vocabulary screens in a desktop browser viewport.
