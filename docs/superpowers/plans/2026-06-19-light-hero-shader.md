# Light Hero Shader Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a separate light-theme shader preset for the public landing hero while preserving the dark-theme shader path.

**Architecture:** Pass the existing `SiteTheme` state from `PublicSiteApp` into `EnglishSparkLanding`, then map it to `GenerativeArtScene variant="dark" | "light"`. Keep default props dark-compatible and tune only light CSS overlays so the light shader remains visible behind readable text.

**Tech Stack:** React 19, TypeScript, Three.js, CSS, Playwright.

---

### Task 1: Regression Test

**Files:**
- Modify: `site-react/e2e/public-site.spec.ts`

- [ ] Add assertions to the existing English landing and light theme tests for `data-hero-preset`.
- [ ] Assert light `.landing-hero__matter` is visible enough: opacity at least `0.62`, no `mix-blend-mode: multiply`, and contrast not below `1`.
- [ ] Run `npm --prefix site-react run e2e -- public-site.spec.ts -g "landing light theme keeps"` and confirm the new assertion fails before implementation.

### Task 2: Theme Prop Wiring

**Files:**
- Modify: `site-react/src/PublicSiteApp.tsx`
- Modify: `site-react/src/EnglishSparkLanding.tsx`

- [ ] Add `siteTheme?: "light" | "dark"` to `EnglishSparkLanding`.
- [ ] Pass `theme` from `PublicSiteApp` into `EnglishSparkLanding`.
- [ ] Add `data-hero-preset={siteTheme}` to `.landing-hero`.
- [ ] Pass `variant={siteTheme}` to `GenerativeArtScene`.

### Task 3: Light Shader Variant

**Files:**
- Modify: `site-react/src/components/ui/anomalous-matter-hero.tsx`
- Modify: `site-react/src/englishSparkLanding.css`

- [ ] Add `variant?: "dark" | "light"` to `GenerativeArtSceneProps`.
- [ ] Keep the current shader constants as the dark variant.
- [ ] Add light shader uniforms for brighter background, darker blue line structure, stronger particles, and slightly higher mesh opacity.
- [ ] Tune light CSS so the shader is not washed out: higher matter opacity, normal/screen blending, and a less opaque right-side veil.

### Task 4: Verification And Sync

**Files:**
- Build output under `Сайт полиглота для бота/`

- [ ] Run targeted Playwright tests for landing light/dark hero.
- [ ] Run `npm --prefix site-react run build`.
- [ ] Run relevant built-output Playwright checks.
- [ ] Review git diff, commit, and push to the configured GitHub remote.
