# App Platform Layout Stability Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make every completed pronunciation report appear below its recording workspace without desktop/mobile layout cross-effects.

**Architecture:** Keep `data-platform` as the platform boundary. Move the completed report outside the two-column recording tools aside so it becomes a normal full-width block after the workbench. Consolidate pronunciation CSS into platform-scoped desktop and mobile rules; the mobile rules must not be selected merely because browser zoom narrows the effective viewport.

**Tech Stack:** React 19, TypeScript, Vite, CSS, Playwright.

---

### Task 1: Capture Report Geometry Regressions

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts:2674-2806`

- [ ] **Step 1: Write failing desktop geometry coverage**

Add a test that uploads the existing pronunciation fixture and asserts the completed report is below the workbench:

```ts
expect(reportBox!.y).toBeGreaterThanOrEqual(workbenchBox!.y + workbenchBox!.height);
expect(reportBox!.x + reportBox!.width).toBeLessThanOrEqual(viewport.width);
```

- [ ] **Step 2: Run the focused test to verify it fails**

Run: `npm --prefix web-react run e2e -- --grep "pronunciation report follows the desktop workspace"`

Expected: FAIL because the report is currently rendered inside the right-side controls aside.

- [ ] **Step 3: Write failing mobile long-report coverage**

Extend the pronunciation API fixture with several weak words and long tips, then assert target, controls, and report have non-overlapping vertical rectangles and the report fits the viewport width.

- [ ] **Step 4: Run the focused mobile test to verify it fails**

Run: `npm --prefix web-react run e2e -- --grep "long pronunciation report keeps mobile flow"`

Expected: FAIL because the test asserts the new standalone report boundary.

### Task 2: Put Completed Reports In Normal Flow

**Files:**
- Modify: `web-react/src/App.tsx:6899-6932`

- [ ] **Step 1: Move only the completed report markup**

Keep recording controls in `.pronunciation-tools-v2`. Render the existing report below `.pronunciation-workbench-v2` only when `practiceResult` exists:

```tsx
{practiceResult ? (
  <section className="pronunciation-result-window-v2 v2-panel">
    <PronunciationReport report={practiceResult} copy={copy} />
  </section>
) : null}
```

- [ ] **Step 2: Run focused geometry tests**

Run: `npm --prefix web-react run e2e -- --grep "pronunciation report follows|long pronunciation report"`

Expected: PASS.

### Task 3: Isolate Desktop And Mobile Pronunciation CSS

**Files:**
- Modify: `web-react/src/styles/app.css:6764-6944,7655-7708,9652-9743,10308-10464`

- [ ] **Step 1: Remove conflicting pronunciation layout overrides**

Replace repeated viewport-only pronunciation workbench rules with one desktop and one mobile block scoped by `data-platform`. Preserve the desktop two-column workbench and set the report to a full-width normal-flow block.

```css
.v2-shell--desktop .pronunciation-workbench-v2 { grid-template-columns: minmax(420px, 1fr) minmax(280px, 360px); }
.v2-shell--desktop .pronunciation-result-window-v2 { width: 100%; }
.v2-shell--mobile .pronunciation-workbench-v2 { display: block; }
.v2-shell--mobile .pronunciation-result-window-v2 { margin-top: 12px; overflow-wrap: anywhere; }
```

- [ ] **Step 2: Run the focused pronunciation suite**

Run: `npm --prefix web-react run e2e -- --grep "pronunciation"`

Expected: PASS with no layout-overlap assertions failing.

### Task 4: Verify Platform Isolation And Build

**Files:**
- Verify: `web-react/e2e/web-smoke.spec.ts`
- Verify: `web-react/src/App.tsx`
- Verify: `web-react/src/styles/app.css`

- [ ] **Step 1: Run the complete web smoke suite**

Run: `npm --prefix web-react run e2e`

Expected: PASS.

- [ ] **Step 2: Run the production build**

Run: `npm --prefix web-react run build`

Expected: Vite completes with exit code 0 and emits `web/index.html`.

- [ ] **Step 3: Review the diff and commit**

Run: `git diff --check` and review only `App.tsx`, `app.css`, the focused e2e test, and this plan. Commit with `fix(web): isolate pronunciation result layout`.
