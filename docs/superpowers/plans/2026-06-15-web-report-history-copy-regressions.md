# Web Report, AI Tutor History, and Admin Fix Copy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make desktop bug reporting visible, remove AI Tutor lesson-goal text from history bodies, and clarify Telegram admin correction order.

**Architecture:** Keep existing UI and backend flow. Add small regression tests first, then make scoped changes in the existing React app shell, AI Tutor history message construction, and Telegram word-report copy.

**Tech Stack:** React 19 + Vite + Playwright for web, Go backend tests for Telegram/admin report behavior.

---

### Task 1: Desktop Bug Report Visibility

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`
- Modify: `web-react/src/styles/app.css`

- [x] **Step 1: Write the failing Playwright assertion**

In `web-react/e2e/web-smoke.spec.ts`, extend the existing `bug report paste keeps one clipboard image and no generic section text` test so desktop must prove the report button is visible, named correctly, and clickable before using it:

```ts
const reportButton = page.locator(isMobile ? ".mobile-report-button-v2" : ".v2-report-button");
await expect(reportButton).toBeVisible();
await expect(reportButton).toHaveAccessibleName("Сообщить об ошибке");
await reportButton.click();
```

- [x] **Step 2: Run the focused e2e test to verify the desktop assertion**

Run:

```bash
npm --prefix web-react run e2e -- --project=desktop-chromium -g "bug report paste keeps one clipboard image"
```

Expected before the CSS fix: the test can fail if the desktop report button is clipped or not considered visible.

- [x] **Step 3: Implement minimal CSS resilience**

In `web-react/src/styles/app.css`, keep `.v2-report-button` in the desktop topbar but make the action row shrink safely:

```css
.v2-topbar__actions {
  min-width: 0;
  flex-wrap: wrap;
}

.v2-report-button {
  flex: 0 0 auto;
}

@media (max-width: 1180px) {
  .v2-report-button span {
    display: none;
  }
}
```

Use the existing nearby topbar/report styles; do not create a new floating desktop button.

- [x] **Step 4: Re-run the focused e2e test**

Run:

```bash
npm --prefix web-react run e2e -- --project=desktop-chromium -g "bug report paste keeps one clipboard image"
```

Expected after the fix: the desktop test passes and the dialog opens from `.v2-report-button`.

### Task 2: AI Tutor History Body Cleanup

**Files:**
- Modify: `web_app_shell_test.go`
- Modify: `web-react/src/App.tsx`

- [x] **Step 1: Write the failing static regression test**

In `web_app_shell_test.go`, add a test that requires a helper marker and prevents the old history body join:

```go
func TestReactFrontendBuildsAITutorHistoryFromStoryOnly(t *testing.T) {
	appSource, err := os.ReadFile(filepath.Join("web-react", "src", "App.tsx"))
	if err != nil {
		t.Fatal(err)
	}
	source := string(appSource)
	for _, want := range []string{
		`function aiTutorHistoryBody`,
		`aiTutorHistoryBody(lessonPayload, step)`,
	} {
		if !strings.Contains(source, want) {
			t.Fatalf("missing AI Tutor history body cleanup marker %q", want)
		}
	}
	if strings.Contains(source, `[lessonPayload?.lesson_goal || step.instruction, lessonPayload?.story?.text_target]`) {
		t.Fatalf("AI Tutor history still prepends lesson_goal before story text")
	}
}
```

- [x] **Step 2: Run the focused Go test to verify it fails**

Run:

```bash
go test ./... -run TestReactFrontendBuildsAITutorHistoryFromStoryOnly
```

Expected before the implementation: FAIL because `aiTutorHistoryBody` does not exist and the old join still exists.

- [x] **Step 3: Implement the helper and use it in start/restart**

In `web-react/src/App.tsx`, add a small helper near the AI Tutor text helpers:

```ts
function aiTutorHistoryBody(lesson: AiTutorStep["lesson"] | undefined, step: AiTutorStep | null) {
  const story = cleanAppText(lesson?.story?.text_target).trim();
  if (story) return story;
  return cleanTutorPrompt(step?.instruction);
}
```

Use it in both `startTutor` and `restartTutorLesson` panel messages:

```ts
panelMessage(
  aiTutorHistoryBody(lessonPayload, step),
  "default",
  cleanAppText(lessonPayload?.title || step.title || copy("ai_tutor", "AI Tutor")),
  getRecord(payload),
  "tutor",
)
```

- [x] **Step 4: Re-run the focused Go test**

Run:

```bash
go test ./... -run TestReactFrontendBuildsAITutorHistoryFromStoryOnly
```

Expected after the implementation: PASS.

### Task 3: Telegram Admin Fix Copy

**Files:**
- Modify: `telegram_test.go`
- Modify: `ai_tutor_word_reports.go`

- [x] **Step 1: Write parser and prompt regression tests**

In `telegram_test.go`, add a parser unit test:

```go
func TestParseAITutorWordReportFixTextUsesLearningWordThenInterfaceTranslation(t *testing.T) {
	word, translation, ok := parseAITutorWordReportFixText("football match - футбольный матч")
	if !ok {
		t.Fatal("expected correction text to parse")
	}
	if word != "football match" || translation != "футбольный матч" {
		t.Fatalf("parsed word=%q translation=%q", word, translation)
	}
}
```

Update `TestTelegramAITutorWordReportFixReplyAppliesCorrection` server handler so it captures the prompt text and asserts it contains:

```go
"сначала слово/фраза на языке изучения"
"затем перевод на языке интерфейса"
"football match - футбольный матч"
```

- [x] **Step 2: Run the focused Telegram tests to verify prompt-copy failure**

Run:

```bash
go test ./... -run "TestParseAITutorWordReportFixTextUsesLearningWordThenInterfaceTranslation|TestTelegramAITutorWordReportFixReplyAppliesCorrection"
```

Expected before copy update: parser test passes, prompt-copy assertion fails.

- [x] **Step 3: Update Telegram prompt and fallback copy**

In `ai_tutor_word_reports.go`, replace ambiguous `word - translation` copy with explicit Russian admin text:

```go
const aiTutorWordReportFixFormat = "Формат: сначала слово/фраза на языке изучения, затем перевод на языке интерфейса: football match - футбольный матч"
```

Use it in `promptAITutorWordReportFix` for both AI Tutor and vocabulary reports, preserving the existing SQLite destination details. Use it in the invalid-format reply as well:

```go
return true, b.telegram.sendMessageToChat(ctx, message.Chat.ID, aiTutorWordReportFixFormat)
```

- [x] **Step 4: Re-run the focused Telegram tests**

Run:

```bash
go test ./... -run "TestParseAITutorWordReportFixTextUsesLearningWordThenInterfaceTranslation|TestTelegramAITutorWordReportFixReplyAppliesCorrection"
```

Expected after the update: PASS.

### Task 4: Full Verification, Browser Check, Commit, Push

**Files:**
- Verify all files changed by Tasks 1-3.

- [x] **Step 1: Run targeted Go regression tests**

Run:

```bash
go test ./... -run "TestReactFrontendBuildsAITutorHistoryFromStoryOnly|TestParseAITutorWordReportFixTextUsesLearningWordThenInterfaceTranslation|TestTelegramAITutorWordReportFixReplyAppliesCorrection|TestWebAppAITutorWordReportUIWiring"
```

Expected: PASS.

- [x] **Step 2: Run web build**

Run:

```bash
npm --prefix web-react run build
```

Expected: TypeScript and Vite build complete with exit code 0.

- [x] **Step 3: Run focused Playwright checks**

Run:

```bash
npm --prefix web-react run e2e -- --project=desktop-chromium -g "bug report paste keeps one clipboard image"
npm --prefix web-react run e2e -- --project=mobile-chromium -g "bug report paste keeps one clipboard image"
```

Expected: both focused Playwright tests pass.

- [x] **Step 4: Browser smoke verification attempted**

Note: direct browser smoke reached the local login screen, so authenticated shell behavior was verified through mocked Playwright e2e coverage.

Open the local app in the in-app browser and verify:

- desktop width shows a bug report button in the topbar and opens `.bug-report-dialog-v2`;
- mobile width still shows `.mobile-report-button-v2`;
- AI Tutor story history no longer starts with the lesson goal when mocked or available test data is used.

- [x] **Step 5: Commit and push scoped changes**

Stage only the files touched by this plan:

```bash
git add web-react/e2e/web-smoke.spec.ts web-react/src/styles/app.css web_app_shell_test.go web-react/src/App.tsx telegram_test.go ai_tutor_word_reports.go docs/superpowers/plans/2026-06-15-web-report-history-copy-regressions.md web/index.html web/assets/index-DBmr_BWk.css web/assets/index-WgKCCR-1.js
git commit -m "fix: stabilize report controls and tutor history copy"
git push origin codex/ai-tutor-rebuild-fix
```

### Task 5: Build Artifact Guard

**Files:**
- Verify: `web/index.html`
- Verify: `web/assets/index-WgKCCR-1.js`

- [x] **Step 1: Verify the generated web shell points at the built bundle**

Check that `web/index.html` loads the current built asset hash.

- [x] **Step 2: Verify the generated bundle keeps story-only history construction**

Check that the built JS bundle still contains the story-first helper and does not contain the old `lesson_goal` history join.
