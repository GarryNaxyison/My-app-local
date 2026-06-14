# AI Tutor Lesson UI Cleanup Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Clean up the AI Tutor and normal Lesson frontend states so reports, completion, notes, errors, localization, and next-lesson actions match the approved user-facing behavior.

**Architecture:** Keep this as a frontend-focused patch using existing React state, `copy(...)` localization, and CSS classes. Add source-level regression tests in the existing Go web shell test file where this repository already validates frontend wiring, then run the React TypeScript/Vite build and Browser/Playwright smoke checks.

**Tech Stack:** Go test suite for source-level frontend regression checks, React/TypeScript in `web-react/src/App.tsx`, app copy in `web-react/src/lib/i18n.ts`, CSS in `web-react/src/styles/app.css`, Vite build.

---

## File Structure

Files to modify:

```text
web_app_shell_test.go
web-react/src/App.tsx
web-react/src/lib/i18n.ts
web-react/src/styles/app.css
```

Reference spec:

```text
docs/superpowers/specs/2026-06-14-ai-tutor-lesson-ui-cleanup-design.md
```

---

### Task 1: Add Frontend Regression Tests

**Files:**

- Modify: `web_app_shell_test.go`

- [x] **Step 1: Add source assertions for the approved UI behavior**

Add tests that read `web-react/src/App.tsx` and `web-react/src/lib/i18n.ts` and assert:

```go
func TestWebAppAITutorLessonUICleanupWiring(t *testing.T) {
    app := readWebAppFile(t, "web-react/src/App.tsx")
    i18n := readWebAppFile(t, "web-react/src/lib/i18n.ts")

    assertContains(t, app, "isTutorWordLearnStage")
    assertContains(t, app, "tutor-word-report-action-v2")
    assertContains(t, app, "tutor-note-strip-v2")
    assertContains(t, app, "tutorActualErrorLines")
    assertContains(t, app, "ai_tutor_next_lesson")
    assertContains(t, i18n, "ai_tutor_next_lesson")
    assertContains(t, i18n, "Следующий урок")
}
```

- [x] **Step 2: Run the targeted test and confirm it fails before implementation**

Run:

```powershell
go test ./... -run TestWebAppAITutorLessonUICleanupWiring
```

Expected before implementation:

```text
missing expected frontend wiring
```

---

### Task 2: Fix AI Tutor Stage, Completion, Notes, and Errors

**Files:**

- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`

- [x] **Step 1: Add robust stage detection**

Use a stage-name helper:

```ts
const isTutorWordLearnStage = (stage?: string | null) => Boolean(stage && /^word_learn_\d+$/.test(stage));
```

- [x] **Step 2: Show the word report button for active word stages**

Render the existing report button when the active backend session stage is a `word_learn_N` stage, even if `kind` is missing:

```tsx
const canReportActiveWord = !isPreviewing && isTutorWordLearnStage(aiTutorStep?.stage);
```

- [x] **Step 3: Remove duplicate completion copy**

Do not render `tutor_complete_body` in both the generic instruction paragraph and the completion summary. Keep one localized completion body and one localized XP/status line.

- [x] **Step 4: Split notes from actual errors**

Create separate derived lists:

```ts
const tutorNoteLines = collectTutorNoteLines(...);
const tutorActualErrorLines = collectTutorActualErrorLines(...);
```

Render notes once in:

```tsx
<div className="tutor-note-strip-v2">...</div>
```

Render `Ошибки` only when `tutorActualErrorLines.length > 0`.

- [x] **Step 5: Add compact note strip styles**

Use horizontal scrolling and stable spacing:

```css
.tutor-note-strip-v2 {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  white-space: nowrap;
}
```

---

### Task 3: Fix Lesson Completed CTA

**Files:**

- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/lib/i18n.ts`

- [x] **Step 1: Add localized next-lesson copy**

Add `ai_tutor_next_lesson` or an existing shared key with:

```ts
ai_tutor_next_lesson: "Next lesson"
ai_tutor_next_lesson: "Следующий урок"
```

- [x] **Step 2: Use it after a normal lesson answer is complete**

In the normal Lesson view, when there is no active answer input because the last lesson answer produced feedback, keep the chat readable and label the CTA:

```tsx
copy("ai_tutor_next_lesson", "Next lesson")
```

---

### Task 4: Strengthen Localization Fallbacks for AI Tutor UI

**Files:**

- Modify: `web-react/src/lib/i18n.ts`
- Modify: `web-react/src/App.tsx`

- [x] **Step 1: Add missing AI Tutor copy keys**

Ensure known UI keys for stage labels, option labels, note actions, errors, completion, and next lesson exist in the explicit English/Russian copy and the generated locale copy block.

- [x] **Step 2: Prefer localized known-stage UI copy over backend English**

For known AI Tutor stages, render labels and instructions from local `copy(...)` keys. Keep generated target-language lesson content unchanged.

- [x] **Step 3: Prevent `Раздел` from appearing as tutor completion UI**

Do not let generic `section` fallback fill tutor completion or tutor view labels.

---

### Task 5: Verification and Build

**Files:**

- Test/build only

- [x] **Step 1: Run targeted source tests**

Run:

```powershell
go test ./... -run 'TestWebApp.*Tutor|TestWebApp.*Lesson'
```

Expected:

```text
PASS
```

- [x] **Step 2: Build the frontend**

Run:

```powershell
npm --prefix web-react run build
```

Expected:

```text
✓ built
```

- [x] **Step 3: Verify in browser**

Serve the app locally with the existing project server or Vite preview and use Browser/Playwright to confirm:

- the active AI Tutor word card has the report button;
- the notes row is single and horizontally scrollable;
- completion text is not duplicated;
- the Lesson completed CTA reads `Следующий урок`.

- [ ] **Step 4: Commit and push**

Run:

```powershell
git status --short
git add docs/superpowers/specs/2026-06-14-ai-tutor-lesson-ui-cleanup-design.md docs/superpowers/plans/2026-06-14-ai-tutor-lesson-ui-cleanup.md web_app_shell_test.go web-react/src/App.tsx web-react/src/lib/i18n.ts web-react/src/styles/app.css web
git commit -m "fix: clean up ai tutor lesson ui"
git push origin HEAD
```
