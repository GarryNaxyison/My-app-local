# Telegram AI Tutor And Pronunciation Stage Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make AI Tutor work from Telegram and make the web tutor pronunciation stage the primary desktop work area without changing other tutor stages.

**Architecture:** Reuse the existing `store.nextTutorLesson` and `buildTutorLessonForSequence` path for Telegram. Add a small Telegram renderer and localized `AITutor` UI label. Scope web layout changes to a pronunciation-only wrapper under `.tutor-context-v2`.

**Tech Stack:** Go Telegram bot, existing JSON/SQLite stores, React/Vite app, CSS, Playwright.

---

## File Structure

- Modify `i18n.go`: add `AITutor` to `uiCopy`, runtime fallback helper, generated copy assignment, explicit English/Russian copy values.
- Modify `telegram.go`: use `copy.AITutor` in main and learning menus; add a tutor keyboard.
- Modify `bot.go`: handle `menu_tutor`, add `startTutorLesson`, add Telegram tutor message formatting helpers.
- Modify `telegram_test.go`: add regression tests for `menu_tutor`, all-language button labels, and non-regression of `/start`.
- Modify `web-react/src/App.tsx`: wrap only the tutor pronunciation stage in a primary/side layout.
- Modify `web-react/src/styles/app.css`: add desktop two-column pronunciation stage rules and mobile stack override.
- Modify `web-react/e2e/web-smoke.spec.ts`: extend existing AI Tutor smoke test with desktop layout assertions and mobile stack assertions.

## Task 1: Telegram Failing Tests

**Files:**
- Modify: `telegram_test.go`

- [ ] **Step 1: Add failing test for menu_tutor callback**

Add a test near the other callback/menu tests:

```go
func TestMenuTutorCallbackStartsTutorLesson(t *testing.T) {
	var methods []string
	var payloads []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		methods = append(methods, r.URL.Path)
		payloads = append(payloads, payload)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	store := &jsonStore{
		path:  filepath.Join(t.TempDir(), "users.json"),
		users: map[int64]userState{},
	}
	user := userState{
		TelegramID:        123,
		FirstName:         "Test",
		InterfaceLanguage: "ru",
		InterfaceSelected: true,
		TimezoneSelected:  true,
		LanguageSelected:  true,
		LearningLanguage:  "en",
		Level:             "A1",
	}
	store.users[user.TelegramID] = user
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}

	err := b.handleCallbackQuery(context.Background(), callbackQuery{
		ID:   "cb-1",
		From: telegramUser{ID: user.TelegramID, FirstName: user.FirstName},
		Data: "menu_tutor",
	})
	if err != nil {
		t.Fatalf("handleCallbackQuery(menu_tutor) error = %v", err)
	}
	if len(methods) < 2 || methods[0] != "/answerCallbackQuery" || methods[1] != "/sendMessage" {
		t.Fatalf("expected callback answer and tutor message, got methods=%v payloads=%#v", methods, payloads)
	}
	text, _ := payloads[1]["text"].(string)
	if strings.Contains(text, ui(user).UnknownButton) {
		t.Fatalf("menu_tutor fell through to unknown button: %q", text)
	}
	if !strings.Contains(text, "AI Репетитор") || !strings.Contains(text, "В кафе") || !strings.Contains(text, "Could I see the menu") {
		t.Fatalf("expected rendered tutor lesson, got %q", text)
	}
}
```

- [ ] **Step 2: Add failing test for localized AI Tutor labels**

Add:

```go
func TestAITutorMenuLabelIsLocalizedForEveryInterfaceLanguage(t *testing.T) {
	for _, language := range interfaceLanguages() {
		copy := ui(userState{InterfaceLanguage: language.Code, InterfaceSelected: true})
		if strings.TrimSpace(copy.AITutor) == "" {
			t.Fatalf("missing AI Tutor label for %s", language.Code)
		}
		keyboard := mainMenuInlineKeyboard(copy)
		rows := keyboard["inline_keyboard"].([][]map[string]any)
		found := false
		for _, row := range rows {
			for _, button := range row {
				if button["callback_data"] == "menu_tutor" {
					found = true
					text, _ := button["text"].(string)
					if !strings.Contains(text, copy.AITutor) {
						t.Fatalf("menu_tutor label for %s = %q, want copy label %q", language.Code, text, copy.AITutor)
					}
					if language.Code != "ru" && strings.Contains(text, "AI Репетитор") {
						t.Fatalf("menu_tutor label for %s leaked Russian hardcode: %q", language.Code, text)
					}
				}
			}
		}
		if !found {
			t.Fatalf("missing menu_tutor for %s", language.Code)
		}
	}
}
```

- [ ] **Step 3: Run failing tests**

Run:

```bash
go test ./... -run 'TestMenuTutorCallbackStartsTutorLesson|TestAITutorMenuLabelIsLocalizedForEveryInterfaceLanguage|TestStartCommandShowsPrivacyPolicyBeforeOnboarding'
```

Expected: the first two tests fail because `menu_tutor` is unhandled and `uiCopy.AITutor` does not exist yet.

## Task 2: Telegram Implementation

**Files:**
- Modify: `i18n.go`
- Modify: `telegram.go`
- Modify: `bot.go`
- Test: `telegram_test.go`

- [ ] **Step 1: Add UI copy field and runtime label helper**

Implement:

```go
type uiCopy struct {
	// existing fields...
	Practice string
	AITutor  string
	Shadowing string
	// existing fields...
}
```

Set explicit copies:

```go
AITutor: "AI Репетитор",
```

for Russian, and:

```go
AITutor: "AI Tutor",
```

for English. In `generatedUICopy`, assign:

```go
AITutor: aiTutorButtonLabel(code),
```

Add:

```go
func aiTutorButtonLabel(code string) string {
	switch normalizeInterfaceLanguage(code) {
	case "ru":
		return "AI Репетитор"
	case "es":
		return "Tutor IA"
	case "de":
		return "KI-Tutor"
	case "fr":
		return "Tuteur IA"
	case "it":
		return "Tutor IA"
	case "pt":
		return "Tutor IA"
	case "pl":
		return "Tutor AI"
	case "uk":
		return "AI Репетитор"
	default:
		return "AI Tutor"
	}
}
```

In `withRuntimeUICopy`, fill empty labels:

```go
if copy.AITutor == "" {
	copy.AITutor = aiTutorButtonLabel(code)
}
```

- [ ] **Step 2: Use localized button labels**

Replace both hardcoded menu buttons in `telegram.go`:

```go
{"text": "🤖 " + copy.AITutor, "callback_data": "menu_tutor"}
```

- [ ] **Step 3: Add tutor callback handler**

In `handleCallbackQuery`, add:

```go
case "menu_tutor":
	return b.startTutorLesson(ctx, chatID, user)
```

- [ ] **Step 4: Add Telegram tutor renderer**

Add helpers near the other bot learning flow methods:

```go
func (b *bot) startTutorLesson(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	_ = b.telegram.sendChatAction(ctx, chatID, "typing")
	lesson, err := b.store.nextTutorLesson(user, func(sequence int) (tutorLesson, error) {
		return buildTutorLessonForSequence(tutorReusableLessonUser(user), sequence)
	})
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeMessage(user, "lesson_task_failed", err.Error()), ui(user))
	}
	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, formatTelegramTutorLesson(lesson, user), tutorLessonKeyboard(ui(user)))
}
```

The formatter must escape MarkdownV2 and include title, scenario, mini-explanation, target phrase, one choice, and one dialogue reply:

```go
func formatTelegramTutorLesson(lesson tutorLesson, user userState) string {
	copy := ui(user)
	var b strings.Builder
	title := strings.TrimSpace(lesson.Title)
	if title == "" {
		title = copy.AITutor
	}
	b.WriteString("*" + escapeMarkdownV2(title) + "*")
	if lesson.LessonNumber > 0 && lesson.CourseSize > 0 {
		b.WriteString(" · " + escapeMarkdownV2(fmt.Sprintf("%d/%d", lesson.LessonNumber, lesson.CourseSize)))
	}
	appendTelegramTutorLine(&b, "🎯", lesson.Goal)
	appendTelegramTutorLine(&b, "📍", lesson.Scenario)
	appendTelegramTutorLine(&b, "🧩", lesson.MiniExplanation)
	target := strings.TrimSpace(lesson.PronunciationText)
	if target == "" {
		target = strings.TrimSpace(lesson.ScenarioSlots.ModelAnswer)
	}
	appendTelegramTutorLine(&b, "🗣", target)
	appendTelegramTutorChoice(&b, lesson)
	appendTelegramTutorDialogue(&b, lesson)
	return b.String()
}
```

- [ ] **Step 5: Run Telegram tests**

Run:

```bash
go test ./... -run 'TestMenuTutorCallbackStartsTutorLesson|TestAITutorMenuLabelIsLocalizedForEveryInterfaceLanguage|TestStartCommandShowsPrivacyPolicyBeforeOnboarding'
```

Expected: all selected tests pass.

- [ ] **Step 6: Commit Telegram implementation**

```bash
git add i18n.go telegram.go bot.go telegram_test.go
git commit -m "fix: start ai tutor from telegram"
```

## Task 3: Web Pronunciation Stage Tests

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add desktop layout assertions to existing AI Tutor smoke**

After the pronunciation stage becomes visible, assert:

```ts
const pronunciationStage = page.locator(".tutor-pronunciation-stage-v2");
await expect(pronunciationStage).toBeVisible();
const targetBox = await page.locator(".tutor-pronunciation-primary-v2").boundingBox();
const toolsBox = await page.locator(".tutor-pronunciation-tools-v2").boundingBox();
expect(targetBox?.width || 0).toBeGreaterThan(toolsBox?.width || 0);
expect(targetBox?.height || 0).toBeGreaterThan(180);
```

After checking pronunciation, assert:

```ts
await expect(page.locator(".tutor-pronunciation-primary-v2 .tutor-pronunciation-target-v2")).toContainText("I'd like coffee for breakfast");
await expect(page.locator(".tutor-pronunciation-tools-v2 .pronunciation-report-v2")).toContainText("67/100");
```

- [ ] **Step 2: Add mobile stack assertion**

In the same test or a focused mobile branch, assert at mobile width:

```ts
await page.setViewportSize({ width: 390, height: 844 });
const mobileTarget = await page.locator(".tutor-pronunciation-primary-v2").boundingBox();
const mobileTools = await page.locator(".tutor-pronunciation-tools-v2").boundingBox();
expect((mobileTools?.y || 0)).toBeGreaterThan((mobileTarget?.y || 0));
```

- [ ] **Step 3: Run the focused test and verify it fails before implementation**

Run:

```bash
npm --prefix web-react run e2e -- --grep "AI Tutor cafe scenario uses slots"
```

Expected: fails because `.tutor-pronunciation-stage-v2` does not exist yet.

## Task 4: Web Pronunciation Stage Implementation

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Wrap pronunciation material only**

Replace the pronunciation branch in `renderTutorMaterial` with:

```tsx
return (
  <div className="tutor-pronunciation-stage-v2">
    <div className="tutor-pronunciation-primary-v2">
      <div className="tutor-pronunciation-target-v2">
        <span>{copy("tutor_pronunciation_target", "Repeat exactly")}</span>
        <strong>{target}</strong>
      </div>
      {target ? <AudioActionRow clips={[{ label: copy("spoken_model", "Spoken model"), text: target, targetLanguage: user.learning_language }]} /> : null}
      {tutorPronunciationCorrection ? <AudioActionRow clips={[{ label: interfaceLocale === "ru" ? "Исправленный образец" : "Corrected model", text: tutorPronunciationCorrection, targetLanguage: user.learning_language }]} /> : null}
    </div>
    <div className="tutor-pronunciation-tools-v2">
      {!viewingPastStage ? (
        <div className="tutor-listening-voice-v2">
          <FileControls voiceFile={voiceFile} imageFile={imageFile} setVoiceFile={setVoiceFile} setImageFile={setImageFile} allowImage={false} copy={copy} />
          <p>{voiceFile ? copy("tutor_voice_ready", "Voice answer is ready for pronunciation check.") : copy("tutor_voice_hint", "Record your repeat for pronunciation scoring.")}</p>
        </div>
      ) : null}
      {tutorPronunciation ? <PronunciationReport pronunciation={tutorPronunciation} copy={copy} compact /> : null}
    </div>
  </div>
);
```

- [ ] **Step 2: Add scoped CSS**

Add near existing tutor pronunciation CSS:

```css
.tutor-pronunciation-stage-v2 {
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(280px, 0.72fr);
  gap: 16px;
  align-items: start;
  min-height: 360px;
}

.tutor-pronunciation-primary-v2 {
  display: grid;
  gap: 14px;
  min-height: 320px;
}

.tutor-pronunciation-primary-v2 .tutor-pronunciation-target-v2 {
  min-height: 230px;
  align-content: center;
  padding: 24px;
}

.tutor-pronunciation-primary-v2 .tutor-pronunciation-target-v2 strong {
  font-size: clamp(1.8rem, 2.2vw, 2.8rem);
  line-height: 1.22;
}

.tutor-pronunciation-tools-v2 {
  display: grid;
  gap: 12px;
  align-content: start;
}
```

Inside the existing mobile media query, add:

```css
.tutor-pronunciation-stage-v2 {
  grid-template-columns: 1fr;
  min-height: 0;
}

.tutor-pronunciation-primary-v2,
.tutor-pronunciation-primary-v2 .tutor-pronunciation-target-v2 {
  min-height: 0;
}
```

- [ ] **Step 3: Run focused Playwright test**

Run:

```bash
npm --prefix web-react run e2e -- --grep "AI Tutor cafe scenario uses slots"
```

Expected: desktop and mobile projects pass.

- [ ] **Step 4: Commit web implementation**

```bash
git add web-react/src/App.tsx web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts
git commit -m "fix: emphasize tutor pronunciation stage"
```

## Task 5: Full Verification And Push

**Files:**
- All changed files.

- [ ] **Step 1: Run Go tests**

```bash
go test ./...
```

Expected: pass.

- [ ] **Step 2: Run focused web e2e**

```bash
npm --prefix web-react run e2e -- --grep "AI Tutor cafe scenario uses slots"
```

Expected: pass for configured desktop and mobile projects.

- [ ] **Step 3: Inspect git status**

```bash
git status --short
```

Expected: clean.

- [ ] **Step 4: Push**

```bash
git push origin main
```

Expected: commits pushed to `https://github.com/GarryNaxyison/My-app-local.git`.
