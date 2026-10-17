# Image Preview Pan Design

## Goal

When a product screenshot preview is zoomed with the mouse wheel, the user can hold the mouse button and drag the image around. Native scrollbars must not appear in the preview.

## Design

The public landing preview in `site-react/src/EnglishSparkLanding.tsx` keeps the existing click-to-open and wheel-to-zoom behavior. The preview stores a pan offset in React state and applies it to the image with CSS transform alongside the existing scale.

The preview frame becomes a clipped viewport with `overflow: hidden`. When the scale is greater than `1`, pointer down starts a drag, pointer move updates the pan offset, and pointer up/cancel stops the drag. Opening, closing, or zooming back to scale `1` resets the pan offset to the centered image.

## Testing

Extend `site-react/e2e/landing-v3.spec.ts` so the screenshot preview is opened, zoomed, checked for no scroll overflow, dragged by mouse, and verified to have a non-zero pan transform.
