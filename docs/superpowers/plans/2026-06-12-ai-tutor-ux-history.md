# AI Tutor UX History Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add backend-backed AI Tutor completed lesson history with restart, plus the requested AI Tutor UI copy/navigation fixes.

**Architecture:** The backend exposes completed lesson history and same-lesson restart APIs on top of existing `ai_tutor_sessions` and `ai_tutor_lessons`. The React app loads backend history into the existing completed-lessons dialog, keeps local storage as fallback, and adds read-only preview navigation for completed steps.

**Tech Stack:** Go `net/http` API, SQLite/json store implementations, React 19 + TypeScript + Vite, Nx wrapper commands, Playwright/browser verification.

---

### Task 1: Backend Completed History API

**Files:**
- Modify: `storage.go`
- Modify: `sqlite_store.go`
- Modify: `web_api.go`
- Test: `web_api_feature_test.go`

- [ ] **Step 1: Write the failing API test**

Add a test after `TestWebAITutorReviewCompletesSession`:

```go
func TestWebAITutorCompletedLessonsReturnsFinishedSessions(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	payload := validAITutorLessonPayloadForTest()
	payload.Title = "Sports Sunday"
	payload.Theme = "sports"
	lesson := aiTutorLessonRecord{ID: "lesson-completed-history", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	session := aiTutorSessionRecord{ID: "session-completed-history", TelegramID: -42, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageComplete, Status: aiTutorSessionComplete, CompletedAt: "2026-06-12T12:00:00Z"}
	if err := store.createAITutorSession(session); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}

	body := requestJSON(t, api, cookie, http.MethodGet, "/api/ai-tutor/completed", nil)
	items, _ := body["items"].([]any)
	if len(items) != 1 {
		t.Fatalf("items = %#v, want one completed lesson", body["items"])
	}
	item, _ := items[0].(map[string]any)
	if item["lesson_id"] != lesson.ID || item["session_id"] != session.ID || item["title"] != "Sports Sunday" {
		t.Fatalf("completed lesson item mismatch: %#v", item)
	}
}
```

- [ ] **Step 2: Run the failing test**

Run: `npm run test -- --runInBand web_api_feature_test.go -run TestWebAITutorCompletedLessonsReturnsFinishedSessions`

Expected: fail because `/api/ai-tutor/completed` is not registered.

- [ ] **Step 3: Implement store/query/API**

Add:

```go
type aiTutorCompletedLessonRecord struct {
	SessionID   string               `json:"session_id"`
	LessonID    string               `json:"lesson_id"`
	Title       string               `json:"title"`
	Topic       string               `json:"topic"`
	Level       string               `json:"level"`
	CompletedAt string               `json:"completed_at"`
	Lesson      aiTutorLessonPayload `json:"lesson"`
}
```

Add `completedAITutorLessons(telegramID int64, limit int)` to `store`, implement it in SQLite by joining completed sessions to lessons, and register `GET /api/ai-tutor/completed` in `web_api.go`.

- [ ] **Step 4: Run the test again**

Run: `npm run test -- --runInBand web_api_feature_test.go -run TestWebAITutorCompletedLessonsReturnsFinishedSessions`

Expected: pass.

### Task 2: Backend Restart API

**Files:**
- Modify: `storage.go`
- Modify: `sqlite_store.go`
- Modify: `web_api.go`
- Test: `web_api_feature_test.go`

- [ ] **Step 1: Write the failing restart test**

Add:

```go
func TestWebAITutorRestartCreatesNewSessionForCompletedLesson(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-restart-history", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	if err := store.createAITutorSession(aiTutorSessionRecord{ID: "session-old-restart", TelegramID: -42, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageComplete, Status: aiTutorSessionComplete, CompletedAt: "2026-06-12T12:00:00Z"}); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}

	body := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/restart", map[string]any{"lesson_id": lesson.ID})
	if body["current_stage"] != aiTutorStageStoryIntro {
		t.Fatalf("current_stage = %#v response=%#v", body["current_stage"], body)
	}
	session, _ := body["session"].(map[string]any)
	if session["lesson_id"] != lesson.ID && session["LessonID"] != lesson.ID {
		t.Fatalf("restart returned wrong lesson session: %#v", session)
	}
}
```

- [ ] **Step 2: Run the failing test**

Run: `npm run test -- --runInBand web_api_feature_test.go -run TestWebAITutorRestartCreatesNewSessionForCompletedLesson`

Expected: fail because `/api/ai-tutor/restart` is not registered.

- [ ] **Step 3: Implement restart**

Add `userCanRestartAITutorLesson(telegramID int64, lessonID string)` to the store. Add `handleAITutorRestart` in `web_api.go` that validates access, loads the lesson, creates a new session with `CurrentStage: aiTutorStageStoryIntro`, and returns `api.aiTutorDTO`.

- [ ] **Step 4: Add rejection test**

Add a test that creates a lesson with no completed session for `-42`, posts restart, and expects `403`.

- [ ] **Step 5: Run restart tests**

Run: `npm run test -- --runInBand web_api_feature_test.go -run 'TestWebAITutor(Restart|Completed)'`

Expected: pass.

### Task 3: Frontend History And Restart

**Files:**
- Modify: `web-react/src/lib/types.ts`
- Modify: `web-react/src/App.tsx`
- Test/build: `web_app_shell_test.go`, `npm --prefix web-react run build`

- [ ] **Step 1: Add API contract test**

In `TestReactFrontendKeepsGoAPIContracts`, add:

```go
`/api/ai-tutor/completed`,
`/api/ai-tutor/restart`,
```

- [ ] **Step 2: Run failing contract test**

Run: `npm run test -- --runInBand web_app_shell_test.go -run TestReactFrontendKeepsGoAPIContracts`

Expected: fail until frontend references both endpoints.

- [ ] **Step 3: Add frontend types**

Add `AiTutorCompletedLesson` and `AiTutorCompletedLessonsResponse` with `session_id`, `lesson_id`, `title`, `topic`, `level`, `completed_at`, and `lesson`.

- [ ] **Step 4: Load backend history**

When the completed-lessons dialog opens, call `api<AiTutorCompletedLessonsResponse>("/api/ai-tutor/completed")`, map backend records into `TutorCompletedLessonRecord`, and merge them with local records.

- [ ] **Step 5: Restart from history**

On completed lesson click, call `api<AiTutorResponse>("/api/ai-tutor/restart", { method: "POST", body: { lesson_id } })`, then reuse the same state update path as the normal AI Tutor start/answer response.

- [ ] **Step 6: Run contract test and build**

Run:

```bash
npm run test -- --runInBand web_app_shell_test.go -run TestReactFrontendKeepsGoAPIContracts
npm --prefix web-react run build
```

Expected: both pass.

### Task 4: AI Tutor Step Preview And Copy Fixes

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/lib/i18n.ts`
- Modify: `web-react/src/styles/app.css`

- [ ] **Step 1: Story separator**

Change the story words render from:

```tsx
<span>{display(word.target)} ? {display(word.interface_translation)}</span>
```

to:

```tsx
<span>{display(word.target)} - {display(word.interface_translation)}</span>
```

- [ ] **Step 2: Placeholder key**

Add `tutor_answer_field_placeholder` to explicit English/Russian/German copy and derived copy. Use it in the AI Tutor textarea for non-production text stages.

- [ ] **Step 3: Completed-step preview**

Add local `previewStage` state. Enable buttons where `index < currentStageIndex`. Render selected previous stage material using the current lesson data, hide the composer in preview mode, and show a button to return to current step.

- [ ] **Step 4: Quick start blocks**

Update `app_guide_step_*` copy for English and Russian to describe product capabilities instead of procedural steps. Keep five cards and existing layout.

- [ ] **Step 5: Build**

Run: `npm --prefix web-react run build`

Expected: pass.

### Task 5: Verification, Commit, Push

**Files:**
- All modified files

- [ ] **Step 1: Run backend focused tests**

Run:

```bash
npm run test -- --runInBand web_api_feature_test.go -run 'TestWebAITutor(Completed|Restart|Review|Start|Answer)'
```

- [ ] **Step 2: Run frontend build**

Run:

```bash
npm --prefix web-react run build
```

- [ ] **Step 3: Run Nx check**

Run:

```bash
npm run check
```

- [ ] **Step 4: Browser verification**

Start dev server:

```bash
npm run web:dev
```

Open with Playwright/browser, inspect the quick-start dialog and AI Tutor placeholder.

- [ ] **Step 5: Commit and push**

Run:

```bash
git add storage.go sqlite_store.go web_api.go web_api_feature_test.go web_app_shell_test.go web-react/src/App.tsx web-react/src/lib/types.ts web-react/src/lib/i18n.ts web-react/src/styles/app.css docs/superpowers/plans/2026-06-12-ai-tutor-ux-history.md
git commit -m "feat: add ai tutor lesson history restart"
git push
```

