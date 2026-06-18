# Light Hero Shader Design

## Goal

The landing hero must keep the new very light theme while making the hero animation visible again. Dark theme must keep its current shader look.

## Design

The landing receives the current site theme from `PublicSiteApp` and passes a shader preset into `GenerativeArtScene`. The default preset remains `dark` so existing uses keep their behavior. The light preset uses a separate shader palette and stronger mesh/particle contrast, intended to read like a blueprint/neural map over the pale hero.

## Boundaries

- Change only the public landing hero path.
- Preserve dark theme colors, animation gating, and one-canvas structure.
- Keep the light hero readable with dark text and light glass panels.
- Add a visible DOM marker for regression tests: `data-hero-preset`.

## Verification

- Playwright checks that dark loads `data-hero-preset="dark"` and light loads `data-hero-preset="light"`.
- The light theme test checks that the visual matter layer is no longer heavily faded or multiply-blended.
- Existing hero animation frame-budget and readability tests must still pass.
