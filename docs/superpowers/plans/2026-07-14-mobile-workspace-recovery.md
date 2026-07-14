# Mobile Workspace Recovery Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make every affected mobile learning screen scroll and compose in normal document flow, keep all controls reachable above the bottom navigation, and always offer a visible way to continue a word/review round.

**Architecture:** Treat the mobile app as a two-row shell: one independently scrollable `.context-display` and one non-overlapping `.mobile-bottom-nav-v2`. Replace the accumulated, conflicting mobile geometry rules with one final owner for each workspace (trainer, roleplay, chat composer, tools composer, and note suggestions). Keep desktop layout unchanged by restricting behavioral CSS to `@media (max-width: 760px)` and preserving existing desktop selectors.

**Tech Stack:** React 19, TypeScript, CSS, Vite, Playwright.

---

## Evidence and root-cause map

| User report | Observable root cause in current source | Planned outcome |
| --- | --- | --- |
| 1–2. Word and review result is cut off; heading says “Start a new round” but no reachable continuation action | `ChoiceTrainer` clears `wordChallenge` after a result, so its heading falls back to the start copy. `TrainerResultBox` does render a next button for successful results, but it is placed after lengthy result content inside a constrained inner scroller. | A dedicated, labelled continuation action is present for both result and empty-round states, and it is reachable by scrolling the one page-level content region. |
| 3–4. Roleplay input and send button overlap or drift after another answer | The same mobile selectors are declared repeatedly. Earlier rules make a fixed-height two-row roleplay workspace; later rules change the composer to a two-column icon button, while final rules partly switch it back to flow. The cascade leaves incompatible `height`, `overflow`, grid, and button-width constraints. | Roleplay output, textarea, send action, file controls, and saved phrases form a single vertical flow; only the output history itself can scroll when necessary. |
| 5. Large gap in spelling | Generic `trainer-display` and mobile viewport-height rules are mixed with a full-height content frame, so short trainer states inherit unused height rather than content-sized layout. | Spelling content is content-sized and starts directly below its heading; its result remains reachable without creating a synthetic gap. |
| 6–7 and 10. Translator labels, send action, microphone/error controls, and lower area overlap | `.tool-input-shell-v2` is alternately absolute-positioned and grid-positioned across several mobile media blocks. Several rules reserve padding for an overlay button after the button has been returned to document flow. | All tool controls use one-column flow on mobile: input, send button, device controls/error, language selectors, then notes. No component reserves space for a non-existent overlay. |
| 8–9. “Save to notes” covers phrase chips, and phrases/input do not receive enough usable width | Mobile composer styles cap/hide quick-save strips and force horizontal chip rows; this competes with the parent’s fixed workspace height. | The label occupies its own row, chips are visible in a full-width wrapping/scrolling row, and the composer remains above it. |

## File structure and ownership

- `web-react/src/App.tsx`
  - Gives the word/review flow explicit state-specific actions and stable selectors.
  - Gives roleplay/tools/note regions stable data attributes only where Playwright needs an unambiguous semantic target.
- `web-react/src/styles/app.css`
  - Removes contradictory mobile rules for the affected workspaces and creates one final mobile layout contract.
  - Does not alter desktop declarations or selectors outside the `max-width: 760px` breakpoint.
- `web-react/e2e/web-smoke.spec.ts`
  - Adds API-mocked mobile regression cases for each reported view and geometric assertions for overlap, horizontal overflow, reachability, and control order.

### Task 1: Add a reproducible mobile geometry test harness

**Files:**

- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add small geometry helpers near the existing Playwright helpers.**

```ts
type Rect = { left: number; top: number; right: number; bottom: number; width: number; height: number };

const rectOf = async (locator: Locator): Promise<Rect> => {
  const box = await locator.boundingBox();
  expect(box).not.toBeNull();
  return {
    left: box!.x,
    top: box!.y,
    right: box!.x + box!.width,
    bottom: box!.y + box!.height,
    width: box!.width,
    height: box!.height,
  };
};

const expectNoHorizontalOverflow = async (page: Page) => {
  const overflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);
  expect(overflow).toBeLessThanOrEqual(1);
};
```

- [ ] **Step 2: Add failing mobile tests for screenshots 1–2 (words and review).** Mock `/api/words/next`, `/api/words/answer`, `/api/word-game/next`, and `/api/word-game/answer` with a long successful result. Assert the result has exactly one `[data-testid="trainer-next-round"]`, the action is visible after `context-display.scrollTo(0, scrollHeight)`, and it lies above the mobile navigation.

```ts
const context = page.locator('.context-display--words');
await context.evaluate((node) => node.scrollTo(0, node.scrollHeight));
const nextRound = page.getByTestId('trainer-next-round');
const nav = page.locator('.mobile-bottom-nav-v2');
await expect(nextRound).toBeVisible();
expect((await rectOf(nextRound)).bottom).toBeLessThanOrEqual((await rectOf(nav)).top - 4);
```

- [ ] **Step 3: Add failing mobile roleplay tests for screenshots 3–4.** Start a mocked scenario, submit twice, then assert that output ends before the composer starts, textarea does not overlap the send button, and the send button is fully inside the viewport.

```ts
const output = page.locator('[data-testid="roleplay-output"]');
const textarea = page.locator('[data-testid="roleplay-composer"] textarea');
const send = page.locator('[data-testid="roleplay-send"]');
expect((await rectOf(output)).bottom).toBeLessThanOrEqual((await rectOf(textarea)).top + 1);
expect((await rectOf(textarea)).right).toBeLessThanOrEqual((await rectOf(send)).left + 1);
expect((await rectOf(send)).right).toBeLessThanOrEqual(390);
```

- [ ] **Step 4: Add failing tests for screenshots 5–11.** Cover spelling’s initial vertical gap, text-tool control order, Shadowing/Tutor note visibility, and microphone-error alignment. Use the existing session, tools, and audio mocks rather than a live service.

```ts
const spellingHeading = page.locator('.context-display--spelling h2');
const spellingForm = page.locator('.context-display--spelling .inline-form-v2');
expect((await rectOf(spellingForm)).top - (await rectOf(spellingHeading)).bottom).toBeLessThan(56);

const notes = page.locator('.context-display--tools .phrase-quick-save-v2');
const chips = notes.locator('.phrase-quick-save-v2__chips');
expect((await rectOf(chips)).top).toBeGreaterThanOrEqual((await rectOf(notes.locator('.eyebrow'))).bottom + 4);
```

- [ ] **Step 5: Run the focused tests and confirm they fail before the UI change.**

Run: `npm --prefix web-react exec playwright test e2e/web-smoke.spec.ts --project=mobile-chromium -g "mobile (word round|roleplay composer|spelling spacing|tools composer|quick-save)"`

Expected: failures showing the missing stable controls and one or more current overlap/gap assertions.

- [ ] **Step 6: Commit the tests separately.**

```bash
git add web-react/e2e/web-smoke.spec.ts
git commit -m "test: cover mobile workspace regressions"
```

### Task 2: Make word and review round continuation explicit

**Files:**

- Modify: `web-react/src/App.tsx:7549-7700` (`ChoiceTrainer` and `TrainerResultBox`)
- Test: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Preserve the distinction between an unavailable challenge and a completed round.** Introduce a `canStartRound` boolean in `ChoiceTrainer`; it is true only when there is no challenge, no result, and no request is in progress. Do not show it for the backend’s `challenge.empty` response.

```tsx
const canStartRound = !challenge && !result && !busy;

{canStartRound ? (
  <Button data-testid="trainer-start-round" type="button" onClick={() => void start()}>
    <Play size={17} />
    {copy("start_new_round", "Start a new round")}
  </Button>
) : null}
```

- [ ] **Step 2: Give the successful result its permanent, semantic next action.** Keep `TrainerResultBox` as the sole owner of the successful result action, but give it a `data-testid` and a label that is specific to its mode.

```tsx
<Button data-testid="trainer-next-round" type="button" onClick={onNext}>
  <ChevronRight size={17} />
  {compactWordResult ? copy("next_word", "Next") : copy("next", "Next")}
</Button>
```

- [ ] **Step 3: Keep result state intact until the next request starts.** Confirm `startWord`, `startWordGame`, and `startSpelling` only clear their respective result immediately before issuing their new request, never in response to an unrelated navigation render.

- [ ] **Step 4: Re-run the Task 1 word/review tests.**

Run: `npm --prefix web-react exec playwright test e2e/web-smoke.spec.ts --project=mobile-chromium -g "mobile word round"`

Expected: PASS; the action is rendered once and is reachable above the navigation rail.

- [ ] **Step 5: Commit the flow change.**

```bash
git add web-react/src/App.tsx web-react/e2e/web-smoke.spec.ts
git commit -m "fix: keep mobile word rounds actionable"
```

### Task 3: Replace conflicting mobile height and scrolling rules with one shell contract

**Files:**

- Modify: `web-react/src/styles/app.css:4684-10320`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Delete or fold the duplicated mobile declarations that give the same components incompatible geometry.** In particular, consolidate repeated rules for `.v2-app`, `.context-display`, `.chat-workspace`, `.tools-layout-v2`, `.composer-panel-v2`, `.tool-input-shell-v2`, `.roleplay-view-v2--mobile-session`, and `.roleplay-submit-v2`. Do not add a fourth corrective block on top of the existing three.

- [ ] **Step 2: Define the single owner of vertical scrolling.** On phones, the app grid must have a flexible content row and an auto-sized navigation row. The content panel owns `overflow-y: auto`; normal view descendants own their natural height.

```css
@media (max-width: 760px) {
  .v2-app {
    height: 100dvh;
    min-height: 100dvh;
    grid-template-rows: minmax(0, 1fr) auto;
    overflow: hidden;
  }

  .context-display {
    min-width: 0;
    min-height: 0;
    height: auto;
    overflow-x: hidden;
    overflow-y: auto;
    overscroll-behavior: contain;
    padding: 12px;
  }

  .mobile-bottom-nav-v2 {
    position: relative;
    grid-row: 2;
    width: 100%;
    max-width: 100%;
  }
}
```

- [ ] **Step 3: Make trainer and spelling content content-sized.** Remove mobile viewport/minimum heights from `.trainer-display` and spelling-specific descendants. Keep spacing in normal flow, and reserve no blank region for the navigation because the nav is now its own grid row.

```css
@media (max-width: 760px) {
  .context-display--words .trainer-display,
  .context-display--word-game .trainer-display,
  .context-display--spelling .trainer-display {
    min-height: 0;
    height: auto;
    overflow: visible;
    align-content: start;
  }
}
```

- [ ] **Step 4: Verify long-result scrolling and short spelling layout.**

Run: `npm --prefix web-react exec playwright test e2e/web-smoke.spec.ts --project=mobile-chromium -g "mobile (word round|spelling spacing)"`

Expected: PASS; page width has no overflow and no test needs a magic `dvh` offset to reach the bottom action.

- [ ] **Step 5: Commit the shell change.**

```bash
git add web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts
git commit -m "fix: unify mobile workspace scrolling"
```

### Task 4: Put roleplay output and composer into a stable vertical flow

**Files:**

- Modify: `web-react/src/App.tsx:6651-6774`
- Modify: `web-react/src/styles/app.css`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add stable, non-visual test hooks to the roleplay regions.**

```tsx
<div className="chat-workspace__output" data-testid="roleplay-output">
  ...
</div>
<div className="composer-textarea-shell-v2" data-testid="roleplay-composer">
  ...
  <Button data-testid="roleplay-send" className="composer-submit-v2 roleplay-submit-v2" ...>
```

- [ ] **Step 2: On mobile, use a single-column grid for the active roleplay session.** Do not give the roleplay panel a calculated viewport height. The output card grows with the conversation; only `.roleplay-dialog-scroll-v2` may use a bounded scroll when it needs to limit an unusually long history.

```css
@media (max-width: 760px) {
  .context-display--roleplay .roleplay-view-v2--mobile-session,
  .context-display--roleplay .chat-workspace--roleplay {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: auto auto;
    height: auto;
    min-height: 0;
    max-height: none;
    overflow: visible;
  }

  .context-display--roleplay .composer-textarea-shell-v2 {
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    gap: 10px;
  }
}
```

- [ ] **Step 3: Make the mobile roleplay send control a full-width flow sibling.** Remove the 56px two-column override and every roleplay-only padding reservation meant for an overlaid send button.

```css
@media (max-width: 760px) {
  .context-display--roleplay .roleplay-submit-v2 {
    position: static;
    width: 100%;
    min-width: 0;
    min-height: 50px;
    justify-self: stretch;
  }
}
```

- [ ] **Step 4: Re-run the two-message roleplay test from Task 1.**

Run: `npm --prefix web-react exec playwright test e2e/web-smoke.spec.ts --project=mobile-chromium -g "mobile roleplay composer"`

Expected: PASS; textarea, send action, output, and mobile navigation have pairwise non-overlapping rectangles.

- [ ] **Step 5: Commit the roleplay layout.**

```bash
git add web-react/src/App.tsx web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts
git commit -m "fix: keep mobile roleplay composer in flow"
```

### Task 5: Normalize translator and generic composer control order

**Files:**

- Modify: `web-react/src/App.tsx:8729-8894`
- Modify: `web-react/src/styles/app.css`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add stable test hooks around the text-tool form and the tool quick-save strip.**

```tsx
<div className="tool-input-shell-v2" data-testid="tools-composer">
  ...
  <Button data-testid="tools-send" className="tools-submit-v2" ...>
```

- [ ] **Step 2: Keep all phone composer buttons in document flow.** The textarea must use normal padding on phones, followed by the full-width submit button. Remove mobile declarations that set `right`, `bottom`, absolute positioning, or a 118px right padding for `.tools-submit-v2`, `.composer-submit-v2`, and `.roleplay-submit-v2`.

```css
@media (max-width: 760px) {
  .tool-input-shell-v2,
  .composer-textarea-shell-v2 {
    position: static;
    display: grid;
    grid-template-columns: minmax(0, 1fr);
    gap: 10px;
  }

  .tool-input-shell-v2 textarea,
  .composer-textarea-shell-v2 textarea {
    min-height: 128px;
    padding: 14px;
  }

  .tools-submit-v2,
  .composer-submit-v2 {
    position: static;
    width: 100%;
    min-width: 0;
  }
}
```

- [ ] **Step 3: Order translator controls semantically.** Keep the source selector, language-swap button, and target selector after the send button; stack them into one column on phones. Keep microphone/image controls directly after the send action in their modes, including an error message on its own full-width row.

- [ ] **Step 4: Re-run the text/voice/image tool test cases.**

Run: `npm --prefix web-react exec playwright test e2e/web-smoke.spec.ts --project=mobile-chromium -g "mobile tools composer"`

Expected: PASS; each subsequent control begins below the preceding control, all content remains within the viewport width, and the lower panel has no artificial blank-height requirement.

- [ ] **Step 5: Commit the tools layout.**

```bash
git add web-react/src/App.tsx web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts
git commit -m "fix: stack mobile tool controls safely"
```

### Task 6: Keep “Save to notes” readable and independent from the composer

**Files:**

- Modify: `web-react/src/App.tsx:7328-7378` (`PhraseQuickSave` and `TutorNoteStrip`)
- Modify: `web-react/src/styles/app.css:5376-5420` and mobile overrides
- Test: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Separate the quick-save label from its chips semantically.** Associate the chip row with its label and expose a stable selector without changing the visible copy.

```tsx
<section className="phrase-quick-save-v2" data-testid="phrase-quick-save">
  <span id="phrase-quick-save-label" className="eyebrow">...</span>
  <div className="phrase-quick-save-v2__chips" aria-labelledby="phrase-quick-save-label">
```

- [ ] **Step 2: Make the strip full width and vertically safe on phones.** Do not clip its height. The label stays on its own line; chips may wrap normally, and only an intentionally long chip row may horizontally scroll without hiding its first chip.

```css
@media (max-width: 760px) {
  .composer-panel-v2 .phrase-quick-save-v2,
  .roleplay-dialog-card-v2 .phrase-quick-save-v2 {
    width: 100%;
    max-height: none;
    overflow: visible;
    margin-top: 8px;
    padding-top: 10px;
  }

  .phrase-quick-save-v2__chips {
    flex-wrap: wrap;
    overflow: visible;
  }

  .phrase-quick-save-v2__chips button {
    max-width: 100%;
  }
}
```

- [ ] **Step 3: Verify quick-save in Shadowing, Roleplay, and Tools.** Assert the label bottom is before the chip row top, chips remain visible, and neither region overlaps the submit action.

Run: `npm --prefix web-react exec playwright test e2e/web-smoke.spec.ts --project=mobile-chromium -g "mobile quick-save"`

Expected: PASS; the phrase text is visible and selectable, with no clipping or button collision.

- [ ] **Step 4: Commit quick-save layout.**

```bash
git add web-react/src/App.tsx web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts
git commit -m "fix: preserve mobile quick-save phrases"
```

### Task 7: Full verification, visual review, and delivery

**Files:**

- Modify only if test-driven changes require it: `web-react/e2e/web-smoke.spec.ts`, `web-react/src/App.tsx`, `web-react/src/styles/app.css`

- [ ] **Step 1: Run the entire mobile Playwright project.**

Run: `npm --prefix web-react exec playwright test e2e/web-smoke.spec.ts --project=mobile-chromium`

Expected: PASS with no retries masking geometry failures.

- [ ] **Step 2: Run desktop checks to prove the desktop web app was not changed.**

Run: `npm --prefix web-react exec playwright test e2e/web-smoke.spec.ts --project=desktop-chromium -g "desktop|roleplay|tools|learn words"`

Expected: PASS; desktop roleplay continues to use its intended two-column workspace where applicable.

- [ ] **Step 3: Build the production bundle.**

Run: `npm --prefix web-react run build`

Expected: Vite completes successfully with a generated `web/` bundle.

- [ ] **Step 4: Inspect the final diff for scope and CSS duplication.**

Run: `git diff --check && rg -n "calc\\(100dvh - 196px|grid-template-columns: minmax\\(0, 1fr\\) 56px|padding-right: 118px" web-react/src/styles/app.css`

Expected: `git diff --check` is silent; the obsolete roleplay and overlay-button declarations are absent from mobile CSS.

- [ ] **Step 5: Manually inspect at 390×844 and 412×915.** Verify words, review, spelling, roleplay after two messages, text/voice/image tools, Shadowing, and notes while the browser chrome is visible. Confirm scrolling occurs within content, not the document; confirm the bottom nav remains visible but never covers a required control.

- [ ] **Step 6: Commit and push the verified implementation.**

```bash
git add web-react/src/App.tsx web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts web/
git commit -m "fix: recover mobile learning workspaces"
git push origin HEAD
```

## Acceptance criteria

- On words and review, a successful result always has one visible `Next` action after the user scrolls to the end; an empty, non-loading round has a `Start a new round` action.
- At 390px and 412px wide, `document.documentElement.scrollWidth <= window.innerWidth + 1` in every affected view.
- No required control overlaps another: output → textarea → send → device/language controls → notes → bottom navigation.
- The spelling form begins within 56px of its heading’s bottom and has no viewport-height spacer.
- “Save to notes” label and its phrase chips are both visible; at least the first chip has a positive-width, on-screen rectangle.
- All mobile-specific rules remain in `@media (max-width: 760px)`; desktop Playwright checks remain green.
