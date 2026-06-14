# Web Lesson Ribbon Clipboard Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Restore mobile Lesson action visibility, desktop ribbon arrow clickability, and clipboard image paste in Tools image mode.

**Architecture:** Keep the current React/Vite app structure. Add focused Playwright regressions first, then make scoped CSS layout changes and a small clipboard event plumbing change in `App.tsx`.

**Tech Stack:** React 19, TypeScript, Vite, Playwright, Nx workspace.

---

## File Map

- Modify `web-react/e2e/web-smoke.spec.ts`: strengthen regressions around the reported user scenarios.
- Modify `web-react/src/App.tsx`: route paste events from the tools work area into `ImageUploadControl`.
- Modify `web-react/src/styles/app.css`: reserve stable ribbon arrow lanes and keep mobile Lesson actions above the bottom nav.

### Task 1: Failing Regression Coverage

- [x] Add a mobile assertion that `.lesson-new-button-v2` has a bounding box inside the viewport and above `.mobile-bottom-nav-v2`.
- [x] Extend the desktop function-ribbon assertion to verify both arrow centers resolve to an arrow button, not a `.function-chip`, with visible gap from the ribbon.
- [x] Add an image-tool paste assertion that dispatches a `ClipboardEvent("paste")` on `.context-display--tools` after image mode is selected and expects `.file-chip-v2` to show the pasted file.
- [x] Run `npm --prefix web-react run e2e -- --grep "desktop function ribbon arrows stay inside the menu frame|mobile lesson keeps output readable|image tool accepts a pasted clipboard image" --reporter=list` and confirm the new clipboard assertion fails before production changes.

### Task 2: Minimal Implementation

- [x] Update `ImageUploadControl` so the same image extraction logic can be used by parent paste handlers.
- [x] Add a paste handler on the Tools workspace while image mode is active, passing pasted image files into the upload control state.
- [x] Adjust `.function-ribbon-shell`, `.function-ribbon`, and `.ribbon-morph-arrow` so arrow lanes have stable width and menu chips cannot overlap the arrows during hover expansion.
- [x] Adjust mobile Lesson layout so `.lesson-new-button-v2` remains in document flow, visible, and clear of the fixed bottom navigation.

### Task 3: Verification and Sync

- [x] Run `npm --prefix web-react run build`.
- [x] Run the targeted Playwright grep from Task 1.
- [x] Inspect `git diff -- web-react/e2e/web-smoke.spec.ts web-react/src/App.tsx web-react/src/styles/app.css docs/superpowers`.
- [ ] Commit and push the verified changes.
