# AI Tutor Rebuild Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the old local AI Tutor with a shared AI-generated lesson engine that runs the full interactive lesson in both Web and Telegram.

**Architecture:** Build a backend-owned tutor engine with strict JSON lesson generation, deterministic validation, persistent sessions, AI answer checkers, and lesson quality promotion. Web and Telegram become adapters over the same `next_step` session API instead of owning separate progress logic.

**Tech Stack:** Go, modernc SQLite, OpenRouter chat completions, Telegram Bot API, React/Vite, Playwright, existing runtime prompt system in `app_prompts.json`.

---

## File Structure

- Create `ai_tutor_core.go`: shared constants, lesson payload structs, stage model, validation helpers, level-band helpers, fingerprint helper.
- Create `ai_tutor_core_test.go`: core validation and stage order tests.
- Create `ai_tutor_store_test.go`: SQLite round-trip and migration tests for new AI Tutor tables.
- Create `ai_tutor_prompts.go`: generation, check, and quality prompt builders.
- Create `ai_tutor_prompts_test.go`: prompt contract tests.
- Create `ai_tutor_engine.go`: `aiTutorEngine`, session start, answer submission, stage transitions, quality promotion.
- Create `ai_tutor_engine_test.go`: fake-AI engine tests.
- Modify `storage.go`: add AI Tutor store types and methods to the `store` interface and JSON store.
- Modify `sqlite_store.go`: add AI Tutor tables, migration cleanup for old tutor lesson rows, and SQLite store methods.
- Modify `web_api.go`: add `/api/ai-tutor/*` endpoints and route old `/api/tutor/start` to the new engine response during transition.
- Modify `web_api_feature_test.go`: add API coverage for start, answer, review, and resume.
- Modify `bot.go`: route `menu_tutor`, `ai_tutor:*` text, and `ait|*` callbacks through the shared engine.
- Modify `telegram_test.go`: add Telegram full-flow tests with fake engine/store hooks.
- Modify `app_prompts.json`: add AI Tutor generator, checker, and quality prompt keys.
- Modify `web-react/src/lib/types.ts`: add server-driven AI Tutor DTO types.
- Leave `web-react/src/lib/api.ts` unchanged; use the existing generic `api()` helper in `App.tsx`.
- Modify `web-react/src/App.tsx`: replace old local TutorView state machine with server-driven `next_step`.
- Modify `web-react/e2e/web-smoke.spec.ts`: add a mocked smoke path for server-driven AI Tutor.

## Shared Contracts

Use these names consistently across tasks:

```go
const (
	aiTutorStatusDraft             = "draft"
	aiTutorStatusPreflightRejected = "preflight_rejected"
	aiTutorStatusInTrial           = "in_trial"
	aiTutorStatusApproved          = "approved"
	aiTutorStatusRejected          = "rejected"
	aiTutorStatusArchived          = "archived"

	aiTutorSessionActive   = "active"
	aiTutorSessionComplete = "complete"

	aiTutorStageStoryIntro     = "story_intro"
	aiTutorStageRetell         = "retell"
	aiTutorStageQuestion1      = "question_1"
	aiTutorStageQuestion2      = "question_2"
	aiTutorStageQuestion3      = "question_3"
	aiTutorStageProduction     = "production"
	aiTutorStageLessonFeedback = "lesson_feedback"
	aiTutorStageReviewSchedule = "review_schedule"
	aiTutorStageComplete       = "complete"
)
```

Use `word_learn_1` through `word_learn_6` and `word_recall_1` through `word_recall_6` as generated stage IDs.

### Task 1: Core Payload Types And Validation

**Files:**
- Create: `ai_tutor_core.go`
- Create: `ai_tutor_core_test.go`

- [ ] **Step 1: Write the failing core validation tests**

Add `ai_tutor_core_test.go`:

```go
package main

import (
	"strings"
	"testing"
)

func validAITutorLessonPayloadForTest() aiTutorLessonPayload {
	return aiTutorLessonPayload{
		Title:             "A Morning Visit",
		Level:             "A1",
		LevelBand:         "A1-A2",
		TargetLanguage:    "en",
		InterfaceLanguage: "ru",
		Theme:             "daily life",
		LessonGoal:        "Understand and retell a short everyday story.",
		GrammarFocus: aiTutorGrammarFocus{
			Name:                      "Simple present",
			ShortExplanationInterface: "Use simple present for routines.",
			ModelSentenceTarget:       "I go to the shop.",
		},
		Story: aiTutorStory{
			TextTarget:    "Mia wakes up early. She drinks water. She walks to a small shop. She buys bread. She says thank you. She goes home.",
			SentenceCount: 6,
		},
		Words: []aiTutorWord{
			{ID: "w1", Target: "wake up", InterfaceTranslation: "prosypatsya", PartOfSpeech: "verb", ExampleSentenceTarget: "I wake up early.", ExampleTranslationInterface: "Ya prosypayus rano."},
			{ID: "w2", Target: "water", InterfaceTranslation: "voda", PartOfSpeech: "noun", ExampleSentenceTarget: "I drink water.", ExampleTranslationInterface: "Ya p'yu vodu."},
			{ID: "w3", Target: "shop", InterfaceTranslation: "magazin", PartOfSpeech: "noun", ExampleSentenceTarget: "The shop is small.", ExampleTranslationInterface: "Magazin malenkiy."},
			{ID: "w4", Target: "bread", InterfaceTranslation: "hleb", PartOfSpeech: "noun", ExampleSentenceTarget: "She buys bread.", ExampleTranslationInterface: "Ona pokupaet hleb."},
			{ID: "w5", Target: "thank you", InterfaceTranslation: "spasibo", PartOfSpeech: "chunk", ExampleSentenceTarget: "Thank you for the bread.", ExampleTranslationInterface: "Spasibo za hleb."},
			{ID: "w6", Target: "home", InterfaceTranslation: "dom", PartOfSpeech: "noun", ExampleSentenceTarget: "She goes home.", ExampleTranslationInterface: "Ona idet domoy."},
		},
		RetellTask: aiTutorRetellTask{
			InstructionInterface: "Retell the story in your own words.",
			MinSentences:         2,
			FeedbackRubric:       []string{"mentions Mia", "mentions shop", "uses target language"},
		},
		ComprehensionQuestions: []aiTutorQuestion{
			{ID: "q1", QuestionTarget: "When does Mia wake up?", ExpectedPoints: []string{"early"}, FeedbackRuleInterface: "Check time detail."},
			{ID: "q2", QuestionTarget: "What does Mia buy?", ExpectedPoints: []string{"bread"}, FeedbackRuleInterface: "Check object detail."},
			{ID: "q3", QuestionTarget: "Where does Mia go at the end?", ExpectedPoints: []string{"home"}, FeedbackRuleInterface: "Check final place."},
		},
		WordLearning: []aiTutorWordTask{{WordID: "w1"}, {WordID: "w2"}, {WordID: "w3"}, {WordID: "w4"}, {WordID: "w5"}, {WordID: "w6"}},
		WordRecall:   []aiTutorWordTask{{WordID: "w1"}, {WordID: "w2"}, {WordID: "w3"}, {WordID: "w4"}, {WordID: "w5"}, {WordID: "w6"}},
		ProductionTask: aiTutorProductionTask{
			InstructionInterface: "Write 2-3 sentences with at least three lesson words.",
			RequiredWordCount:    3,
			SentenceCount:        "2-3",
			EvaluationCriteria:   []string{"2-3 sentences", "uses at least three words", "target language"},
		},
		ReviewOptions: []string{"tomorrow", "3_days", "1_week", "no_review"},
		QualitySelfCheck: aiTutorQualitySelfCheck{
			CEFRReason:    "Short A1 sentences and high-frequency words.",
			WhyReusable:   "Everyday routine topic.",
			DuplicateRisk: "low",
		},
	}
}

func TestAITutorLevelBand(t *testing.T) {
	cases := map[string]string{
		"A1": "A1-A2",
		"A2": "A1-A2",
		"B1": "B1-B2",
		"B2": "B1-B2",
		"C1": "C1-C2",
		"C2": "C1-C2",
		"":   "A1-A2",
	}
	for level, want := range cases {
		if got := aiTutorLevelBand(level); got != want {
			t.Fatalf("aiTutorLevelBand(%q) = %q, want %q", level, got, want)
		}
	}
}

func TestValidateAITutorLessonAcceptsValidPayload(t *testing.T) {
	lesson := validAITutorLessonPayloadForTest()
	if issues := validateAITutorLessonPayload(lesson); len(issues) != 0 {
		t.Fatalf("valid lesson rejected: %#v", issues)
	}
}

func TestValidateAITutorLessonRejectsWrongCounts(t *testing.T) {
	lesson := validAITutorLessonPayloadForTest()
	lesson.Words = lesson.Words[:5]
	lesson.ComprehensionQuestions = lesson.ComprehensionQuestions[:2]
	lesson.Story.TextTarget = "One. Two. Three. Four."
	lesson.Story.SentenceCount = 4
	issues := strings.Join(validateAITutorLessonPayload(lesson), "\n")
	for _, want := range []string{"exactly 6 words", "exactly 3 questions", "5-7 story sentences"} {
		if !strings.Contains(issues, want) {
			t.Fatalf("issues miss %q:\n%s", want, issues)
		}
	}
}

func TestAITutorCanonicalStages(t *testing.T) {
	stages := aiTutorCanonicalStages()
	if stages[0] != aiTutorStageStoryIntro {
		t.Fatalf("first stage = %q", stages[0])
	}
	if stages[len(stages)-1] != aiTutorStageComplete {
		t.Fatalf("last stage = %q", stages[len(stages)-1])
	}
	if len(stages) != 21 {
		t.Fatalf("stage count = %d, want 21", len(stages))
	}
}
```

- [ ] **Step 2: Run the tests and verify they fail**

Run:

```powershell
go test ./... -run "TestAITutor(LevelBand|CanonicalStages)|TestValidateAITutorLesson" -count=1
```

Expected: compile failure with undefined `aiTutorLessonPayload`, `aiTutorLevelBand`, `validateAITutorLessonPayload`, and `aiTutorCanonicalStages`.

- [ ] **Step 3: Add the core implementation**

Create `ai_tutor_core.go`:

```go
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode"
)

const (
	aiTutorStatusDraft             = "draft"
	aiTutorStatusPreflightRejected = "preflight_rejected"
	aiTutorStatusInTrial           = "in_trial"
	aiTutorStatusApproved          = "approved"
	aiTutorStatusRejected          = "rejected"
	aiTutorStatusArchived          = "archived"

	aiTutorSessionActive   = "active"
	aiTutorSessionComplete = "complete"

	aiTutorStageStoryIntro     = "story_intro"
	aiTutorStageRetell         = "retell"
	aiTutorStageQuestion1      = "question_1"
	aiTutorStageQuestion2      = "question_2"
	aiTutorStageQuestion3      = "question_3"
	aiTutorStageProduction     = "production"
	aiTutorStageLessonFeedback = "lesson_feedback"
	aiTutorStageReviewSchedule = "review_schedule"
	aiTutorStageComplete       = "complete"
)

type aiTutorLessonPayload struct {
	Title                  string                  `json:"title"`
	Level                  string                  `json:"level"`
	LevelBand              string                  `json:"level_band"`
	TargetLanguage         string                  `json:"target_language"`
	InterfaceLanguage      string                  `json:"interface_language"`
	Theme                  string                  `json:"theme"`
	LessonGoal             string                  `json:"lesson_goal"`
	GrammarFocus           aiTutorGrammarFocus     `json:"grammar_focus"`
	Story                  aiTutorStory            `json:"story"`
	Words                  []aiTutorWord           `json:"words"`
	RetellTask             aiTutorRetellTask       `json:"retell_task"`
	ComprehensionQuestions []aiTutorQuestion       `json:"comprehension_questions"`
	WordLearning           []aiTutorWordTask       `json:"word_learning"`
	WordRecall             []aiTutorWordTask       `json:"word_recall"`
	ProductionTask         aiTutorProductionTask   `json:"production_task"`
	ReviewOptions          []string                `json:"review_options"`
	QualitySelfCheck       aiTutorQualitySelfCheck `json:"quality_self_check"`
}

type aiTutorGrammarFocus struct {
	Name                      string `json:"name"`
	ShortExplanationInterface string `json:"short_explanation_interface"`
	ModelSentenceTarget       string `json:"model_sentence_target"`
}

type aiTutorStory struct {
	TextTarget    string `json:"text_target"`
	SentenceCount int    `json:"sentence_count"`
}

type aiTutorWord struct {
	ID                          string `json:"id"`
	Target                      string `json:"target"`
	InterfaceTranslation        string `json:"interface_translation"`
	PartOfSpeech                string `json:"part_of_speech"`
	ExampleSentenceTarget       string `json:"example_sentence_target"`
	ExampleTranslationInterface string `json:"example_translation_interface"`
	DifficultyNoteInterface     string `json:"difficulty_note_interface,omitempty"`
}

type aiTutorRetellTask struct {
	InstructionInterface string   `json:"instruction_interface"`
	MinSentences         int      `json:"min_sentences"`
	FeedbackRubric       []string `json:"feedback_rubric"`
}

type aiTutorQuestion struct {
	ID                    string   `json:"id"`
	QuestionTarget        string   `json:"question_target"`
	ExpectedPoints        []string `json:"expected_points"`
	FeedbackRuleInterface string   `json:"feedback_rule_interface"`
}

type aiTutorWordTask struct {
	WordID string `json:"word_id"`
}

type aiTutorProductionTask struct {
	InstructionInterface string   `json:"instruction_interface"`
	RequiredWordCount    int      `json:"required_word_count"`
	SentenceCount        string   `json:"sentence_count"`
	EvaluationCriteria   []string `json:"evaluation_criteria"`
}

type aiTutorQualitySelfCheck struct {
	CEFRReason    string `json:"cefr_reason"`
	WhyReusable   string `json:"why_reusable"`
	DuplicateRisk string `json:"duplicate_risk"`
}

type aiTutorLessonRecord struct {
	ID                string
	LearningLanguage  string
	InterfaceLanguage string
	ExactLevel        string
	LevelBand         string
	Theme             string
	Status            string
	Payload           aiTutorLessonPayload
	Fingerprint       string
	PreflightScore    int
	PostScore         int
	CompletionCount   int
	AverageRating     float64
	CreatedAt         string
	UpdatedAt         string
}

type aiTutorSessionRecord struct {
	ID          string
	TelegramID  int64
	LessonID    string
	Surface     string
	CurrentStage string
	Status      string
	StartedAt   string
	CompletedAt string
	UpdatedAt   string
}

type aiTutorAnswerRecord struct {
	SessionID    string
	Stage        string
	AnswerText   string
	AnswerJSON   string
	FeedbackJSON string
	Correct      bool
	CreatedAt    string
}

type aiTutorQualityCheckRecord struct {
	ID         string
	LessonID   string
	SessionID  string
	Kind       string
	Score      int
	Approved   bool
	IssuesJSON string
	ResultJSON string
	CreatedAt  string
}

func aiTutorLevelBand(level string) string {
	switch normalizeCEFRLevel(level) {
	case "B1", "B2":
		return "B1-B2"
	case "C1", "C2":
		return "C1-C2"
	default:
		return "A1-A2"
	}
}

func aiTutorCanonicalStages() []string {
	stages := []string{aiTutorStageStoryIntro, aiTutorStageRetell, aiTutorStageQuestion1, aiTutorStageQuestion2, aiTutorStageQuestion3}
	for index := 1; index <= 6; index++ {
		stages = append(stages, aiTutorWordLearnStage(index))
	}
	for index := 1; index <= 6; index++ {
		stages = append(stages, aiTutorWordRecallStage(index))
	}
	stages = append(stages, aiTutorStageProduction, aiTutorStageLessonFeedback, aiTutorStageReviewSchedule, aiTutorStageComplete)
	return stages
}

func aiTutorWordLearnStage(index int) string {
	return "word_learn_" + itoa(index)
}

func aiTutorWordRecallStage(index int) string {
	return "word_recall_" + itoa(index)
}

func validateAITutorLessonPayload(lesson aiTutorLessonPayload) []string {
	var issues []string
	if strings.TrimSpace(lesson.Title) == "" {
		issues = append(issues, "title is required")
	}
	if strings.TrimSpace(lesson.Level) == "" {
		issues = append(issues, "level is required")
	}
	if strings.TrimSpace(lesson.TargetLanguage) == "" || strings.TrimSpace(lesson.InterfaceLanguage) == "" {
		issues = append(issues, "target and interface languages are required")
	}
	if count := len(lesson.Words); count != 6 {
		issues = append(issues, "lesson must contain exactly 6 words")
	}
	if count := len(lesson.ComprehensionQuestions); count != 3 {
		issues = append(issues, "lesson must contain exactly 3 questions")
	}
	storyCount := lesson.Story.SentenceCount
	if counted := countAITutorStorySentences(lesson.Story.TextTarget); counted > 0 {
		storyCount = counted
	}
	if storyCount < 5 || storyCount > 7 {
		issues = append(issues, "story must contain 5-7 story sentences")
	}
	if lesson.ProductionTask.RequiredWordCount < 3 {
		issues = append(issues, "production task must require at least 3 words")
	}
	if !hasAITutorReviewOptions(lesson.ReviewOptions) {
		issues = append(issues, "review options must include tomorrow, 3_days, 1_week, no_review")
	}
	for index, word := range lesson.Words {
		if strings.TrimSpace(word.ID) == "" || strings.TrimSpace(word.Target) == "" || strings.TrimSpace(word.InterfaceTranslation) == "" {
			issues = append(issues, "word "+itoa(index+1)+" requires id, target, and interface translation")
		}
	}
	return issues
}

func countAITutorStorySentences(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}
	count := 0
	for _, r := range text {
		if r == '.' || r == '!' || r == '?' || r == '。' || r == '！' || r == '？' {
			count++
		}
	}
	if count == 0 && strings.FieldsFunc(text, unicode.IsSpace) != nil {
		return 1
	}
	return count
}

func hasAITutorReviewOptions(options []string) bool {
	seen := map[string]bool{}
	for _, option := range options {
		seen[strings.TrimSpace(option)] = true
	}
	return seen["tomorrow"] && seen["3_days"] && seen["1_week"] && seen["no_review"]
}

func aiTutorFingerprint(lesson aiTutorLessonPayload) string {
	seed := strings.ToLower(strings.Join([]string{
		lesson.LevelBand,
		lesson.TargetLanguage,
		lesson.InterfaceLanguage,
		lesson.Theme,
		lesson.Story.TextTarget,
	}, "|"))
	sum := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(sum[:])
}
```

- [ ] **Step 4: Run the core tests and verify they pass**

Run:

```powershell
go test ./... -run "TestAITutor(LevelBand|CanonicalStages)|TestValidateAITutorLesson" -count=1
```

Expected: PASS for the new core tests.

- [ ] **Step 5: Commit Task 1**

```powershell
git add ai_tutor_core.go ai_tutor_core_test.go
git commit -m "Add AI tutor core model"
```

### Task 2: SQLite Schema And Store Methods

**Files:**
- Modify: `storage.go`
- Modify: `sqlite_store.go`
- Create: `ai_tutor_store_test.go`

- [ ] **Step 1: Write failing SQLite store tests**

Add `ai_tutor_store_test.go`:

```go
package main

import (
	"path/filepath"
	"testing"
)

func TestSQLiteAITutorLessonAndSessionRoundTrip(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "app.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	lesson := aiTutorLessonRecord{
		ID:                "lesson-1",
		LearningLanguage:  "en",
		InterfaceLanguage: "ru",
		ExactLevel:        "A1",
		LevelBand:         "A1-A2",
		Theme:             "daily life",
		Status:            aiTutorStatusApproved,
		Payload:           validAITutorLessonPayloadForTest(),
		Fingerprint:       "fp-1",
		PreflightScore:    92,
	}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	found, ok, err := store.findApprovedAITutorLesson("en", "ru", "A1-A2", 42)
	if err != nil {
		t.Fatalf("findApprovedAITutorLesson() error = %v", err)
	}
	if !ok || found.ID != lesson.ID {
		t.Fatalf("found lesson = %#v ok=%v", found, ok)
	}
	session := aiTutorSessionRecord{
		ID:           "session-1",
		TelegramID:   42,
		LessonID:     lesson.ID,
		Surface:      "web",
		CurrentStage: aiTutorStageStoryIntro,
		Status:       aiTutorSessionActive,
	}
	if err := store.createAITutorSession(session); err != nil {
		t.Fatalf("createAITutorSession() error = %v", err)
	}
	loaded, ok, err := store.getAITutorSession("session-1")
	if err != nil || !ok {
		t.Fatalf("getAITutorSession() loaded=%#v ok=%v err=%v", loaded, ok, err)
	}
	if loaded.CurrentStage != aiTutorStageStoryIntro || loaded.TelegramID != 42 {
		t.Fatalf("loaded session mismatch: %#v", loaded)
	}
}

func TestSQLiteAITutorMigrationClearsLegacyTutorBase(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "app.sqlite"), "")
	if err != nil {
		t.Fatalf("newSQLiteStore() error = %v", err)
	}
	lesson, err := buildTutorLessonForSequence(tutorReusableLessonUser(userState{TelegramID: 1, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}), 0)
	if err != nil {
		t.Fatalf("legacy lesson build error = %v", err)
	}
	tx, err := store.db.Begin()
	if err != nil {
		t.Fatalf("begin error = %v", err)
	}
	if err := store.saveTutorLessonTx(tx, lesson); err != nil {
		t.Fatalf("save legacy lesson error = %v", err)
	}
	if _, err := tx.Exec(`INSERT INTO tutor_user_lessons (telegram_id, lesson_id, assigned_at, lesson_order) VALUES (1, ?, 'now', 1)`, lesson.ID); err != nil {
		t.Fatalf("save legacy assignment error = %v", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("commit error = %v", err)
	}
	if err := store.clearLegacyTutorLessonBase(); err != nil {
		t.Fatalf("clearLegacyTutorLessonBase() error = %v", err)
	}
	var lessons, assignments int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM tutor_lessons`).Scan(&lessons); err != nil {
		t.Fatalf("count tutor_lessons error = %v", err)
	}
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM tutor_user_lessons`).Scan(&assignments); err != nil {
		t.Fatalf("count tutor_user_lessons error = %v", err)
	}
	if lessons != 0 || assignments != 0 {
		t.Fatalf("legacy rows remain: lessons=%d assignments=%d", lessons, assignments)
	}
}
```

- [ ] **Step 2: Run the tests and verify they fail**

Run:

```powershell
go test ./... -run "TestSQLiteAITutor" -count=1
```

Expected: compile failure for missing AI Tutor store methods.

- [ ] **Step 3: Extend the store interface and JSON store**

Modify `storage.go` by adding methods to `type store interface`:

```go
	saveAITutorLesson(lesson aiTutorLessonRecord) error
	findApprovedAITutorLesson(language string, interfaceLanguage string, levelBand string, telegramID int64) (aiTutorLessonRecord, bool, error)
	createAITutorSession(session aiTutorSessionRecord) error
	getAITutorSession(sessionID string) (aiTutorSessionRecord, bool, error)
	updateAITutorSessionStage(sessionID string, stage string, status string, completedAt string) error
	saveAITutorAnswer(answer aiTutorAnswerRecord) error
	saveAITutorQualityCheck(check aiTutorQualityCheckRecord) error
	scheduleAITutorReview(telegramID int64, lessonID string, dueAt string, intervalCode string) error
	clearLegacyTutorLessonBase() error
```

Add in-memory fields to `jsonStore`:

```go
	aiTutorLessons       map[string]aiTutorLessonRecord
	aiTutorSessions      map[string]aiTutorSessionRecord
	aiTutorAnswers       map[string]aiTutorAnswerRecord
	aiTutorQualityChecks map[string]aiTutorQualityCheckRecord
```

Initialize them in `newJSONStore` next to existing maps.

Add JSON store methods with map storage. Use key `sessionID + "|" + stage` for answers. `findApprovedAITutorLesson` must return the first approved lesson matching language, interface language, and level band that does not already have a completed or active session for the same user.

- [ ] **Step 4: Add SQLite schema**

Modify `sqlite_store.go` `init()` statements by adding the five `CREATE TABLE IF NOT EXISTS ai_tutor_*` statements from the spec. Add these indexes after existing tutor indexes:

```go
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_ai_tutor_lessons_context ON ai_tutor_lessons(learning_language, interface_language, level_band, status, created_at)`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_ai_tutor_sessions_user_status ON ai_tutor_sessions(telegram_id, status, updated_at)`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_ai_tutor_reviews_due ON ai_tutor_reviews(status, due_at, telegram_id)`); err != nil {
		return err
	}
```

- [ ] **Step 5: Add SQLite store methods**

Add methods to `sqlite_store.go`. Use `formatDBTime(time.Now().UTC())` for timestamps. Marshal and unmarshal `aiTutorLessonPayload` through `payload_json`.

```go
func (s *sqliteStore) clearLegacyTutorLessonBase() error {
	_, err := s.db.Exec(`DELETE FROM tutor_user_lessons; DELETE FROM tutor_lessons`)
	return err
}
```

For `findApprovedAITutorLesson`, query:

```sql
SELECT id, learning_language, interface_language, exact_level, level_band, theme, status,
       payload_json, fingerprint, preflight_score, post_score, completion_count, average_rating,
       created_at, updated_at
FROM ai_tutor_lessons
WHERE learning_language = ? AND interface_language = ? AND level_band = ? AND status = ?
  AND NOT EXISTS (
    SELECT 1 FROM ai_tutor_sessions
    WHERE ai_tutor_sessions.telegram_id = ?
      AND ai_tutor_sessions.lesson_id = ai_tutor_lessons.id
  )
ORDER BY post_score DESC, preflight_score DESC, created_at ASC
LIMIT 1
```

Use `normalizeLearningLanguage`, `normalizeInterfaceLanguage`, and `aiTutorLevelBand`.

- [ ] **Step 6: Run SQLite tests**

Run:

```powershell
go test ./... -run "TestSQLiteAITutor" -count=1
```

Expected: PASS for both SQLite tests.

- [ ] **Step 7: Commit Task 2**

```powershell
git add storage.go sqlite_store.go ai_tutor_store_test.go
git commit -m "Add AI tutor SQLite storage"
```

### Task 3: Runtime Prompts

**Files:**
- Create: `ai_tutor_prompts.go`
- Create: `ai_tutor_prompts_test.go`
- Modify: `app_prompts.json`

- [ ] **Step 1: Write failing prompt tests**

Add `ai_tutor_prompts_test.go`:

```go
package main

import (
	"strings"
	"testing"
)

func TestAITutorGenerationPromptContract(t *testing.T) {
	language := learningLanguageByCode("en")
	interfaceLanguage := learningLanguageByCode("ru")
	messages := aiTutorLessonGenerationPrompt(language, interfaceLanguage, "A1", "food", []string{"old-fp"})
	if len(messages) != 2 {
		t.Fatalf("message count = %d", len(messages))
	}
	content := messages[1].Content
	for _, want := range []string{
		"Return strict JSON only",
		"Exactly 5-7 sentences",
		"Select exactly 6 useful target words",
		"Create exactly 3 comprehension questions",
		"old-fp",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("generation prompt misses %q:\n%s", want, content)
		}
	}
}

func TestAITutorCheckerPromptsReturnJSONOnly(t *testing.T) {
	language := learningLanguageByCode("en")
	interfaceLanguage := learningLanguageByCode("ru")
	prompts := [][]chatMessage{
		aiTutorRetellCheckPrompt(language, interfaceLanguage, "A1", validAITutorLessonPayloadForTest(), "Mia goes shop."),
		aiTutorQuestionCheckPrompt(language, interfaceLanguage, "A1", validAITutorLessonPayloadForTest(), validAITutorLessonPayloadForTest().ComprehensionQuestions[0], "Early."),
		aiTutorProductionCheckPrompt(language, interfaceLanguage, "A1", validAITutorLessonPayloadForTest(), "I drink water. I buy bread."),
		aiTutorQualityPrompt(language, interfaceLanguage, validAITutorLessonPayloadForTest(), nil, "preflight"),
	}
	for index, prompt := range prompts {
		if len(prompt) != 2 {
			t.Fatalf("prompt %d message count = %d", index, len(prompt))
		}
		if !strings.Contains(prompt[1].Content, "Return strict JSON only") {
			t.Fatalf("prompt %d does not force JSON:\n%s", index, prompt[1].Content)
		}
	}
}
```

- [ ] **Step 2: Run prompt tests and verify they fail**

Run:

```powershell
go test ./... -run "TestAITutor.*Prompt" -count=1
```

Expected: compile failure for missing prompt functions.

- [ ] **Step 3: Implement prompt builders**

Create `ai_tutor_prompts.go`:

```go
package main

import (
	"encoding/json"
	"strings"
)

func aiTutorLessonGenerationPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, topicSeed string, recentFingerprints []string) []chatMessage {
	level = normalizeCEFRLevel(level)
	if level == "" {
		level = "A1"
	}
	fallback := "You are an expert CEFR language lesson designer.\n\n" +
		"Create one complete micro-lesson for a language learner.\n\n" +
		"Inputs:\n" +
		"- Interface language: " + interfaceLanguage.NativeName + "\n" +
		"- Target learning language: " + language.NativeName + "\n" +
		"- Exact learner level: " + level + "\n" +
		"- Level bank: " + aiTutorLevelBand(level) + "\n" +
		"- Topic/theme seed: " + strings.TrimSpace(topicSeed) + "\n" +
		"- Recent lesson fingerprints to avoid: " + strings.Join(recentFingerprints, ", ") + "\n\n" +
		"Return strict JSON only. No Markdown. No prose outside JSON.\n\n" +
		"Requirements:\n" +
		"1. Create a short story in the target learning language. Exactly 5-7 sentences. Match the exact CEFR level strictly. Do not include translations inside the story.\n" +
		"2. Select exactly 6 useful target words or short chunks from the story. Include target form and interface-language translation.\n" +
		"3. Create a learner retelling task.\n" +
		"4. Create exactly 3 comprehension questions in the target learning language. Answers are expected in the target learning language.\n" +
		"5. Create word learning and recall data for the 6 words.\n" +
		"6. Create a production task where the learner writes 2-3 target-language sentences using at least 3 words.\n" +
		"7. Create review options: tomorrow, 3_days, 1_week, no_review.\n\n" +
		"Return the application JSON shape with keys: title, level, level_band, target_language, interface_language, theme, lesson_goal, grammar_focus, story, words, retell_task, comprehension_questions, word_learning, word_recall, production_task, review_options, quality_self_check."
	return []chatMessage{
		{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage) + "\nReturn machine-readable JSON when the user asks for JSON."},
		{Role: "user", Content: renderAppPrompt("ai_tutor.lesson.generate.user", fallback, commonPromptVars(language, interfaceLanguage))},
	}
}

func aiTutorRetellCheckPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, lesson aiTutorLessonPayload, answer string) []chatMessage {
	body, _ := json.Marshal(lesson)
	fallback := "Return strict JSON only. Check the learner retelling.\nLesson JSON:\n" + string(body) + "\nLearner level: " + level + "\nLearner answer:\n" + answer + "\nReturn keys: correct, comprehension_score, corrected_answer_target, feedback_interface, mistakes."
	return []chatMessage{{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)}, {Role: "user", Content: renderAppPrompt("ai_tutor.retell.check.user", fallback, commonPromptVars(language, interfaceLanguage))}}
}

func aiTutorQuestionCheckPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, lesson aiTutorLessonPayload, question aiTutorQuestion, answer string) []chatMessage {
	body, _ := json.Marshal(struct {
		Lesson   aiTutorLessonPayload `json:"lesson"`
		Question aiTutorQuestion      `json:"question"`
	}{Lesson: lesson, Question: question})
	fallback := "Return strict JSON only. Check one comprehension answer.\nPayload:\n" + string(body) + "\nLearner level: " + level + "\nLearner answer:\n" + answer + "\nReturn keys: correct, partial, corrected_answer_target, feedback_interface, mistakes."
	return []chatMessage{{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)}, {Role: "user", Content: renderAppPrompt("ai_tutor.question.check.user", fallback, commonPromptVars(language, interfaceLanguage))}}
}

func aiTutorProductionCheckPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, lesson aiTutorLessonPayload, answer string) []chatMessage {
	body, _ := json.Marshal(lesson)
	fallback := "Return strict JSON only. Check the learner production task.\nLesson JSON:\n" + string(body) + "\nLearner level: " + level + "\nLearner answer:\n" + answer + "\nReturn keys: correct, used_words, missing_requirement, corrected_answer_target, feedback_interface, mistakes."
	return []chatMessage{{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)}, {Role: "user", Content: renderAppPrompt("ai_tutor.production.check.user", fallback, commonPromptVars(language, interfaceLanguage))}}
}

func aiTutorQualityPrompt(language learningLanguage, interfaceLanguage learningLanguage, lesson aiTutorLessonPayload, sessionResults any, kind string) []chatMessage {
	payload, _ := json.Marshal(struct {
		Kind           string               `json:"kind"`
		Lesson         aiTutorLessonPayload `json:"lesson"`
		SessionResults any                  `json:"session_results,omitempty"`
	}{Kind: kind, Lesson: lesson, SessionResults: sessionResults})
	fallback := "Return strict JSON only. Judge this AI Tutor lesson for reuse in the shared bank.\nPayload:\n" + string(payload) + "\nReturn keys: approved, score, critical_issues, fix_suggestions, reasons."
	return []chatMessage{{Role: "system", Content: "You are a strict CEFR lesson quality reviewer. Return JSON only."}, {Role: "user", Content: renderAppPrompt("ai_tutor.lesson.quality_"+kind+".user", fallback, commonPromptVars(language, interfaceLanguage))}}
}
```

- [ ] **Step 4: Add runtime prompt keys**

Modify `app_prompts.json` by adding keys under `prompts`:

```json
"ai_tutor.lesson.generate.user": {
  "feature": "AI Tutor rebuild lesson generation",
  "tool": "ai_tutor",
  "function": "aiTutorLessonGenerationPrompt",
  "notes": "Creates one strict JSON reusable lesson package.",
  "template": "You are an expert CEFR language lesson designer.\n\nCreate one complete micro-lesson for a language learner.\n\nInputs:\n- Interface language: {{interface_language_native_name}}\n- Target learning language: {{learning_language_native_name}}\n- Exact learner level: {{level}}\n- Level bank: {{level_band}}\n- Topic/theme seed: {{topic_seed}}\n- Recent lesson fingerprints to avoid: {{recent_lesson_fingerprints}}\n\nReturn strict JSON only. No Markdown. No prose outside JSON.\n\nRequirements:\n1. Create a short story in the target learning language. Exactly 5-7 sentences. Match the exact CEFR level strictly. Do not include translations inside the story.\n2. Select exactly 6 useful target words or short chunks from the story. Include target form and interface-language translation.\n3. Create a learner retelling task.\n4. Create exactly 3 comprehension questions in the target learning language. Answers are expected in the target learning language.\n5. Create word learning and recall data for the 6 words.\n6. Create a production task where the learner writes 2-3 target-language sentences using at least 3 words.\n7. Create review options: tomorrow, 3_days, 1_week, no_review.\n\nReturn the application JSON shape with keys: title, level, level_band, target_language, interface_language, theme, lesson_goal, grammar_focus, story, words, retell_task, comprehension_questions, word_learning, word_recall, production_task, review_options, quality_self_check."
}
```

Add equivalent keys with explicit JSON-only templates for:

- `ai_tutor.retell.check.user`
- `ai_tutor.question.check.user`
- `ai_tutor.production.check.user`
- `ai_tutor.lesson.quality_preflight.user`
- `ai_tutor.lesson.quality_post.user`

Keep `function` values matching the Go prompt function names.

- [ ] **Step 5: Run prompt tests**

Run:

```powershell
go test ./... -run "TestAITutor.*Prompt" -count=1
node tools/check_encoding_artifacts.mjs
```

Expected: both commands pass.

- [ ] **Step 6: Commit Task 3**

```powershell
git add ai_tutor_prompts.go ai_tutor_prompts_test.go app_prompts.json
git commit -m "Add AI tutor runtime prompts"
```

### Task 4: Lesson Generation And Session Start Engine

**Files:**
- Create: `ai_tutor_engine.go`
- Create: `ai_tutor_engine_test.go`
- Modify: `storage.go`

- [ ] **Step 1: Write failing start-engine tests**

Add `ai_tutor_engine_test.go`:

```go
package main

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"
)

type fakeAITutorClient struct {
	responses []string
	calls     int
}

func (f *fakeAITutorClient) complete(ctx context.Context, messages []chatMessage, temperature float64, maxTokens int) (string, error) {
	f.calls++
	if len(f.responses) == 0 {
		return `{"approved":true,"score":91,"critical_issues":[],"fix_suggestions":[],"reasons":["ok"]}`, nil
	}
	out := f.responses[0]
	f.responses = f.responses[1:]
	return out, nil
}

func newTestJSONStore(t *testing.T) *jsonStore {
	t.Helper()
	store, err := newJSONStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatalf("newJSONStore() error = %v", err)
	}
	return store
}

func TestAITutorStartGeneratesValidatedSession(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	body, _ := json.Marshal(payload)
	ai := &fakeAITutorClient{responses: []string{
		string(body),
		`{"approved":true,"score":91,"critical_issues":[],"fix_suggestions":[],"reasons":["ok"]}`,
	}}
	engine := newAITutorEngine(store, ai)
	engine.now = func() time.Time { return time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC) }
	user := userState{TelegramID: 77, FirstName: "demo", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}
	result, err := engine.Start(context.Background(), user, "web")
	if err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if result.Session.CurrentStage != aiTutorStageStoryIntro {
		t.Fatalf("stage = %q", result.Session.CurrentStage)
	}
	if result.NextStep.Stage != aiTutorStageStoryIntro {
		t.Fatalf("next step = %#v", result.NextStep)
	}
	if ai.calls != 2 {
		t.Fatalf("AI calls = %d, want generation + preflight", ai.calls)
	}
}
```

- [ ] **Step 2: Run the test and verify it fails**

Run:

```powershell
go test ./... -run "TestAITutorStartGeneratesValidatedSession" -count=1
```

Expected: compile failure for missing `newAITutorEngine`, `Start`, and test JSON store helper.

- [ ] **Step 3: Add engine result types**

Add to `ai_tutor_engine.go`:

```go
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type aiTutorCompletionClient interface {
	complete(ctx context.Context, messages []chatMessage, temperature float64, maxTokens int) (string, error)
}

type aiTutorEngine struct {
	store store
	ai    aiTutorCompletionClient
	now   func() time.Time
}

type aiTutorStep struct {
	Stage       string               `json:"stage"`
	Kind        string               `json:"kind"`
	Title       string               `json:"title"`
	Instruction string               `json:"instruction"`
	Lesson      aiTutorLessonPayload `json:"lesson,omitempty"`
	Word        *aiTutorWord         `json:"word,omitempty"`
	Question    *aiTutorQuestion     `json:"question,omitempty"`
	Options     []string             `json:"options,omitempty"`
}

type aiTutorFeedback struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	JSON    string `json:"json,omitempty"`
}

type aiTutorResult struct {
	Session aiTutorSessionRecord `json:"session"`
	Lesson  aiTutorLessonRecord  `json:"lesson"`
	NextStep aiTutorStep         `json:"next_step"`
	Feedback aiTutorFeedback     `json:"feedback,omitempty"`
}

func newAITutorEngine(store store, ai aiTutorCompletionClient) *aiTutorEngine {
	return &aiTutorEngine{store: store, ai: ai, now: time.Now}
}
```

- [ ] **Step 4: Implement Start**

Add to `ai_tutor_engine.go`:

```go
func (e *aiTutorEngine) Start(ctx context.Context, user userState, surface string) (aiTutorResult, error) {
	language := normalizeLearningLanguage(user.LearningLanguage)
	interfaceLanguage := normalizeInterfaceLanguage(user.InterfaceLanguage)
	level := normalizeCEFRLevel(user.Level)
	levelBand := aiTutorLevelBand(level)
	if approved, ok, err := e.store.findApprovedAITutorLesson(language, interfaceLanguage, levelBand, user.TelegramID); err != nil {
		return aiTutorResult{}, err
	} else if ok {
		return e.createSessionForLesson(user, surface, approved)
	}
	lesson, err := e.generateLesson(ctx, user)
	if err != nil {
		return aiTutorResult{}, err
	}
	if err := e.store.saveAITutorLesson(lesson); err != nil {
		return aiTutorResult{}, err
	}
	return e.createSessionForLesson(user, surface, lesson)
}

func (e *aiTutorEngine) createSessionForLesson(user userState, surface string, lesson aiTutorLessonRecord) (aiTutorResult, error) {
	session := aiTutorSessionRecord{
		ID:           aiTutorNewID("aits"),
		TelegramID:   user.TelegramID,
		LessonID:     lesson.ID,
		Surface:      strings.TrimSpace(surface),
		CurrentStage: aiTutorStageStoryIntro,
		Status:       aiTutorSessionActive,
		StartedAt:    formatDBTime(e.now().UTC()),
		UpdatedAt:    formatDBTime(e.now().UTC()),
	}
	if session.Surface == "" {
		session.Surface = "web"
	}
	if err := e.store.createAITutorSession(session); err != nil {
		return aiTutorResult{}, err
	}
	return aiTutorResult{Session: session, Lesson: lesson, NextStep: aiTutorBuildStep(lesson.Payload, aiTutorStageStoryIntro)}, nil
}

func (e *aiTutorEngine) generateLesson(ctx context.Context, user userState) (aiTutorLessonRecord, error) {
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	raw, err := e.ai.complete(ctx, aiTutorLessonGenerationPrompt(language, interfaceLanguage, user.Level, user.LearningFocus, nil), 0.4, 4000)
	if err != nil {
		return aiTutorLessonRecord{}, err
	}
	var payload aiTutorLessonPayload
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &payload); err != nil {
		return aiTutorLessonRecord{}, err
	}
	if issues := validateAITutorLessonPayload(payload); len(issues) > 0 {
		return aiTutorLessonRecord{}, errors.New(strings.Join(issues, "; "))
	}
	qualityRaw, err := e.ai.complete(ctx, aiTutorQualityPrompt(language, interfaceLanguage, payload, nil, "preflight"), 0.1, 1200)
	if err != nil {
		return aiTutorLessonRecord{}, err
	}
	quality := aiTutorParseQualityResult(qualityRaw)
	status := aiTutorStatusInTrial
	if !quality.Approved || quality.Score < 88 {
		status = aiTutorStatusPreflightRejected
	}
	record := aiTutorLessonRecord{
		ID:                aiTutorNewID("aitl"),
		LearningLanguage:  normalizeLearningLanguage(user.LearningLanguage),
		InterfaceLanguage: normalizeInterfaceLanguage(user.InterfaceLanguage),
		ExactLevel:        normalizeCEFRLevel(user.Level),
		LevelBand:         aiTutorLevelBand(user.Level),
		Theme:             payload.Theme,
		Status:            status,
		Payload:           payload,
		Fingerprint:       aiTutorFingerprint(payload),
		PreflightScore:    quality.Score,
		CreatedAt:         formatDBTime(e.now().UTC()),
		UpdatedAt:         formatDBTime(e.now().UTC()),
	}
	if err := e.store.saveAITutorQualityCheck(aiTutorQualityCheckRecord{ID: aiTutorNewID("aitq"), LessonID: record.ID, Kind: "preflight", Score: quality.Score, Approved: quality.Approved, ResultJSON: qualityRaw, CreatedAt: record.CreatedAt}); err != nil {
		return aiTutorLessonRecord{}, err
	}
	if status == aiTutorStatusPreflightRejected {
		return aiTutorLessonRecord{}, errors.New("generated lesson failed preflight")
	}
	return record, nil
}

func aiTutorNewID(prefix string) string {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return prefix + "-" + itoa(int(time.Now().UnixNano()))
	}
	return prefix + "-" + hex.EncodeToString(bytes[:])
}
```

Add `aiTutorParseQualityResult`:

```go
type aiTutorQualityResult struct {
	Approved bool     `json:"approved"`
	Score    int      `json:"score"`
	Issues   []string `json:"critical_issues"`
}

func aiTutorParseQualityResult(raw string) aiTutorQualityResult {
	var result aiTutorQualityResult
	_ = json.Unmarshal([]byte(strings.TrimSpace(raw)), &result)
	return result
}
```

- [ ] **Step 5: Add step builder**

Add:

```go
func aiTutorBuildStep(lesson aiTutorLessonPayload, stage string) aiTutorStep {
	step := aiTutorStep{Stage: stage, Lesson: lesson}
	switch stage {
	case aiTutorStageStoryIntro:
		step.Kind = "story"
		step.Title = lesson.Title
		step.Instruction = lesson.LessonGoal
	case aiTutorStageRetell:
		step.Kind = "free_text"
		step.Title = "Retell"
		step.Instruction = lesson.RetellTask.InstructionInterface
	case aiTutorStageQuestion1, aiTutorStageQuestion2, aiTutorStageQuestion3:
		index := map[string]int{aiTutorStageQuestion1: 0, aiTutorStageQuestion2: 1, aiTutorStageQuestion3: 2}[stage]
		step.Kind = "free_text"
		step.Title = "Question " + itoa(index+1)
		if index < len(lesson.ComprehensionQuestions) {
			question := lesson.ComprehensionQuestions[index]
			step.Question = &question
			step.Instruction = question.QuestionTarget
		}
	case aiTutorStageProduction:
		step.Kind = "free_text"
		step.Title = "Your sentences"
		step.Instruction = lesson.ProductionTask.InstructionInterface
	case aiTutorStageLessonFeedback:
		step.Kind = "rating"
		step.Title = "Lesson feedback"
		step.Instruction = "How did this lesson feel?"
		step.Options = []string{"easy", "good", "hard", "bad"}
	case aiTutorStageReviewSchedule:
		step.Kind = "review"
		step.Title = "Review"
		step.Instruction = "Choose review time."
		step.Options = []string{"tomorrow", "3_days", "1_week", "no_review"}
	case aiTutorStageComplete:
		step.Kind = "complete"
		step.Title = "Complete"
		step.Instruction = "Lesson complete."
	default:
		if index, ok := aiTutorStageIndex(stage, "word_learn_"); ok && index < len(lesson.Words) {
			word := lesson.Words[index]
			step.Kind = "word_learn"
			step.Title = word.Target
			step.Instruction = word.InterfaceTranslation
			step.Word = &word
		} else if index, ok := aiTutorStageIndex(stage, "word_recall_"); ok && index < len(lesson.Words) {
			word := lesson.Words[index]
			step.Kind = "word_recall"
			step.Title = word.InterfaceTranslation
			step.Instruction = "Recall the target word."
			step.Word = &word
		}
	}
	return step
}
```

Add the missing stage-index helper:

```go
func aiTutorStageIndex(stage string, prefix string) (int, bool) {
	if !strings.HasPrefix(stage, prefix) {
		return 0, false
	}
	value, err := strconv.Atoi(strings.TrimPrefix(stage, prefix))
	if err != nil || value <= 0 {
		return 0, false
	}
	return value - 1, true
}
```

Add `strconv` to the `ai_tutor_engine.go` imports.

- [ ] **Step 6: Run start-engine tests**

Run:

```powershell
go test ./... -run "TestAITutorStartGeneratesValidatedSession" -count=1
```

Expected: PASS.

- [ ] **Step 7: Commit Task 4**

```powershell
git add ai_tutor_engine.go ai_tutor_engine_test.go storage.go
git commit -m "Add AI tutor session start engine"
```

### Task 5: Stage Submission And Review Scheduling

**Files:**
- Modify: `ai_tutor_engine.go`
- Modify: `ai_tutor_engine_test.go`
- Modify: `storage.go`
- Modify: `sqlite_store.go`

- [ ] **Step 1: Write failing stage progression tests**

Add tests to `ai_tutor_engine_test.go`:

```go
func TestAITutorSubmitAdvancesThroughDeterministicStages(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatal(err)
	}
	session := aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageStoryIntro, Status: aiTutorSessionActive}
	if err := store.createAITutorSession(session); err != nil {
		t.Fatal(err)
	}
	engine := newAITutorEngine(store, &fakeAITutorClient{})
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Text: "continue"})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Session.CurrentStage != aiTutorStageRetell {
		t.Fatalf("stage = %q, want retell", result.Session.CurrentStage)
	}
}

func TestAITutorReviewScheduleCompletesSession(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	session := aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageReviewSchedule, Status: aiTutorSessionActive}
	_ = store.createAITutorSession(session)
	engine := newAITutorEngine(store, &fakeAITutorClient{})
	engine.now = func() time.Time { return time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC) }
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Choice: "3_days"})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Session.Status != aiTutorSessionComplete {
		t.Fatalf("status = %q", result.Session.Status)
	}
	if result.Session.CurrentStage != aiTutorStageComplete {
		t.Fatalf("stage = %q", result.Session.CurrentStage)
	}
}
```

- [ ] **Step 2: Run tests and verify failure**

Run:

```powershell
go test ./... -run "TestAITutor(Submit|Review)" -count=1
```

Expected: compile failure for missing `Submit` and `aiTutorSubmitInput`.

- [ ] **Step 3: Add submit input and stage helpers**

Add to `ai_tutor_engine.go`:

```go
type aiTutorSubmitInput struct {
	Text   string `json:"text"`
	Choice string `json:"choice"`
}

func aiTutorNextStage(stage string) string {
	stages := aiTutorCanonicalStages()
	for index, current := range stages {
		if current == stage && index+1 < len(stages) {
			return stages[index+1]
		}
	}
	return aiTutorStageComplete
}
```

- [ ] **Step 4: Implement Submit for deterministic stages**

Add:

```go
func (e *aiTutorEngine) Submit(ctx context.Context, user userState, sessionID string, input aiTutorSubmitInput) (aiTutorResult, error) {
	session, ok, err := e.store.getAITutorSession(sessionID)
	if err != nil {
		return aiTutorResult{}, err
	}
	if !ok || session.TelegramID != user.TelegramID {
		return aiTutorResult{}, errors.New("AI Tutor session not found")
	}
	lesson, ok, err := e.store.getAITutorLesson(session.LessonID)
	if err != nil {
		return aiTutorResult{}, err
	}
	if !ok {
		return aiTutorResult{}, errors.New("AI Tutor lesson not found")
	}
	if session.Status == aiTutorSessionComplete {
		return aiTutorResult{Session: session, Lesson: lesson, NextStep: aiTutorBuildStep(lesson.Payload, aiTutorStageComplete)}, nil
	}
	current := session.CurrentStage
	answerText := strings.TrimSpace(input.Text)
	if input.Choice != "" {
		answerText = strings.TrimSpace(input.Choice)
	}
	if err := e.store.saveAITutorAnswer(aiTutorAnswerRecord{SessionID: session.ID, Stage: current, AnswerText: answerText, Correct: true, CreatedAt: formatDBTime(e.now().UTC())}); err != nil {
		return aiTutorResult{}, err
	}
	next := aiTutorNextStage(current)
	status := aiTutorSessionActive
	completedAt := ""
	if current == aiTutorStageReviewSchedule || next == aiTutorStageComplete {
		status = aiTutorSessionComplete
		next = aiTutorStageComplete
		completedAt = formatDBTime(e.now().UTC())
		if input.Choice != "no_review" && input.Choice != "" {
			if err := e.store.scheduleAITutorReview(user.TelegramID, lesson.ID, aiTutorReviewDueAt(e.now().UTC(), input.Choice), input.Choice); err != nil {
				return aiTutorResult{}, err
			}
		}
	}
	if err := e.store.updateAITutorSessionStage(session.ID, next, status, completedAt); err != nil {
		return aiTutorResult{}, err
	}
	session.CurrentStage = next
	session.Status = status
	session.CompletedAt = completedAt
	return aiTutorResult{Session: session, Lesson: lesson, NextStep: aiTutorBuildStep(lesson.Payload, next), Feedback: aiTutorFeedback{OK: true, Message: "Saved."}}, nil
}
```

Add `getAITutorLesson` to the store interface and implement it in JSON and SQLite stores because `Submit` needs the payload by session lesson ID.

Add:

```go
func aiTutorReviewDueAt(now time.Time, interval string) string {
	switch interval {
	case "tomorrow":
		return formatDBTime(now.AddDate(0, 0, 1))
	case "3_days":
		return formatDBTime(now.AddDate(0, 0, 3))
	case "1_week":
		return formatDBTime(now.AddDate(0, 0, 7))
	default:
		return ""
	}
}
```

- [ ] **Step 5: Run progression tests**

Run:

```powershell
go test ./... -run "TestAITutor(Submit|Review)" -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit Task 5**

```powershell
git add ai_tutor_engine.go ai_tutor_engine_test.go storage.go sqlite_store.go
git commit -m "Add AI tutor stage progression"
```

### Task 6: AI Free-Text Checkers And Quality Promotion

**Files:**
- Modify: `ai_tutor_engine.go`
- Modify: `ai_tutor_engine_test.go`
- Modify: `sqlite_store.go`

- [ ] **Step 1: Write failing checker and promotion tests**

Add:

```go
func TestAITutorSubmitRetellUsesAIChecker(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageRetell, Status: aiTutorSessionActive})
	ai := &fakeAITutorClient{responses: []string{`{"correct":true,"comprehension_score":88,"corrected_answer_target":"Mia goes to the shop.","feedback_interface":"Good.","mistakes":[]}`}}
	engine := newAITutorEngine(store, ai)
	result, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Text: "Mia shop."})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	if result.Session.CurrentStage != aiTutorStageQuestion1 {
		t.Fatalf("stage = %q", result.Session.CurrentStage)
	}
	if result.Feedback.JSON == "" {
		t.Fatalf("missing checker JSON feedback")
	}
	if ai.calls != 1 {
		t.Fatalf("AI calls = %d", ai.calls)
	}
}

func TestAITutorCompletionPromotesApprovedTrialLesson(t *testing.T) {
	store := newTestJSONStore(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusInTrial, Payload: payload}
	_ = store.saveAITutorLesson(lesson)
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: 9, LessonID: lesson.ID, Surface: "web", CurrentStage: aiTutorStageReviewSchedule, Status: aiTutorSessionActive})
	ai := &fakeAITutorClient{responses: []string{`{"approved":true,"score":93,"critical_issues":[],"fix_suggestions":[],"reasons":["clear"]}`}}
	engine := newAITutorEngine(store, ai)
	_, err := engine.Submit(context.Background(), userState{TelegramID: 9, InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}, "session-1", aiTutorSubmitInput{Choice: "tomorrow"})
	if err != nil {
		t.Fatalf("Submit() error = %v", err)
	}
	updated, _, _ := store.getAITutorLesson("lesson-1")
	if updated.Status != aiTutorStatusApproved || updated.PostScore != 93 {
		t.Fatalf("lesson not promoted: %#v", updated)
	}
}
```

- [ ] **Step 2: Run tests and verify failure**

Run:

```powershell
go test ./... -run "TestAITutor(SubmitRetell|CompletionPromotes)" -count=1
```

Expected: failure because free-text stages are not using AI checkers and completion does not promote.

- [ ] **Step 3: Route free-text stages to checker prompts**

In `Submit`, before deterministic answer save, branch on `session.CurrentStage`:

```go
if aiTutorStageNeedsChecker(current) {
	checkRaw, err := e.checkFreeTextStage(ctx, user, lesson.Payload, current, answerText)
	if err != nil {
		return aiTutorResult{}, err
	}
	if err := e.store.saveAITutorAnswer(aiTutorAnswerRecord{SessionID: session.ID, Stage: current, AnswerText: answerText, FeedbackJSON: checkRaw, Correct: true, CreatedAt: formatDBTime(e.now().UTC())}); err != nil {
		return aiTutorResult{}, err
	}
	next := aiTutorNextStage(current)
	if err := e.store.updateAITutorSessionStage(session.ID, next, aiTutorSessionActive, ""); err != nil {
		return aiTutorResult{}, err
	}
	session.CurrentStage = next
	return aiTutorResult{Session: session, Lesson: lesson, NextStep: aiTutorBuildStep(lesson.Payload, next), Feedback: aiTutorFeedback{OK: true, Message: "Checked.", JSON: checkRaw}}, nil
}
```

Add helpers:

```go
func aiTutorStageNeedsChecker(stage string) bool {
	return stage == aiTutorStageRetell || stage == aiTutorStageQuestion1 || stage == aiTutorStageQuestion2 || stage == aiTutorStageQuestion3 || stage == aiTutorStageProduction
}

func (e *aiTutorEngine) checkFreeTextStage(ctx context.Context, user userState, lesson aiTutorLessonPayload, stage string, answer string) (string, error) {
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	switch stage {
	case aiTutorStageRetell:
		return e.ai.complete(ctx, aiTutorRetellCheckPrompt(language, interfaceLanguage, user.Level, lesson, answer), 0.2, 1200)
	case aiTutorStageQuestion1, aiTutorStageQuestion2, aiTutorStageQuestion3:
		index := map[string]int{aiTutorStageQuestion1: 0, aiTutorStageQuestion2: 1, aiTutorStageQuestion3: 2}[stage]
		if index >= len(lesson.ComprehensionQuestions) {
			return "", errors.New("question stage is out of range")
		}
		return e.ai.complete(ctx, aiTutorQuestionCheckPrompt(language, interfaceLanguage, user.Level, lesson, lesson.ComprehensionQuestions[index], answer), 0.2, 1200)
	case aiTutorStageProduction:
		return e.ai.complete(ctx, aiTutorProductionCheckPrompt(language, interfaceLanguage, user.Level, lesson, answer), 0.2, 1200)
	default:
		return "", errors.New("stage does not support AI checking")
	}
}
```

- [ ] **Step 4: Add post-lesson quality promotion**

When completing in `Submit`, call:

```go
if err := e.finishAITutorLesson(ctx, user, lesson, session.ID); err != nil {
	return aiTutorResult{}, err
}
```

Add:

```go
func (e *aiTutorEngine) finishAITutorLesson(ctx context.Context, user userState, lesson aiTutorLessonRecord, sessionID string) error {
	if lesson.Status != aiTutorStatusInTrial {
		return nil
	}
	raw, err := e.ai.complete(ctx, aiTutorQualityPrompt(userLearningLanguage(user), userInterfaceLanguage(user), lesson.Payload, map[string]string{"session_id": sessionID}, "post"), 0.1, 1200)
	if err != nil {
		return err
	}
	quality := aiTutorParseQualityResult(raw)
	status := aiTutorStatusRejected
	if quality.Approved && quality.Score >= 88 {
		status = aiTutorStatusApproved
	}
	if err := e.store.saveAITutorQualityCheck(aiTutorQualityCheckRecord{ID: aiTutorNewID("aitq"), LessonID: lesson.ID, SessionID: sessionID, Kind: "post", Score: quality.Score, Approved: quality.Approved, ResultJSON: raw, CreatedAt: formatDBTime(e.now().UTC())}); err != nil {
		return err
	}
	return e.store.updateAITutorLessonQuality(lesson.ID, status, quality.Score)
}
```

Add `updateAITutorLessonQuality` to `store`, JSON store, and SQLite store.

- [ ] **Step 5: Run checker and promotion tests**

Run:

```powershell
go test ./... -run "TestAITutor(SubmitRetell|CompletionPromotes)" -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit Task 6**

```powershell
git add ai_tutor_engine.go ai_tutor_engine_test.go storage.go sqlite_store.go
git commit -m "Add AI tutor checkers and promotion"
```

### Task 7: Web API Endpoints

**Files:**
- Modify: `web_api.go`
- Modify: `web_api_feature_test.go`

- [ ] **Step 1: Write failing Web API tests**

Add these tests to `web_api_feature_test.go`:

```go
func TestWebAITutorStartReturnsSessionStep(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-web-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	body := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/start", map[string]any{})
	if _, ok := body["session"].(map[string]any); !ok {
		t.Fatalf("missing session in response: %#v", body)
	}
	next, ok := body["next_step"].(map[string]any)
	if !ok || next["stage"] != aiTutorStageStoryIntro {
		t.Fatalf("next_step = %#v", body["next_step"])
	}
}

func TestWebAITutorAnswerAdvancesSession(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-web-2", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	start := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/start", map[string]any{})
	session := start["session"].(map[string]any)
	sessionID, _ := session["ID"].(string)
	if sessionID == "" {
		sessionID, _ = session["id"].(string)
	}
	answer := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/answer", map[string]any{"session_id": sessionID, "text": "continue"})
	if answer["current_stage"] != aiTutorStageRetell {
		t.Fatalf("current_stage = %#v response=%#v", answer["current_stage"], answer)
	}
}

func TestWebAITutorReviewCompletesSession(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	payload := validAITutorLessonPayloadForTest()
	lesson := aiTutorLessonRecord{ID: "lesson-web-3", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}
	if err := store.saveAITutorLesson(lesson); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	start := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/start", map[string]any{})
	session := start["session"].(map[string]any)
	sessionID, _ := session["ID"].(string)
	if sessionID == "" {
		sessionID, _ = session["id"].(string)
	}
	if err := store.updateAITutorSessionStage(sessionID, aiTutorStageReviewSchedule, aiTutorSessionActive, ""); err != nil {
		t.Fatalf("updateAITutorSessionStage() error = %v", err)
	}
	done := requestJSON(t, api, cookie, http.MethodPost, "/api/ai-tutor/review", map[string]any{"session_id": sessionID, "choice": "no_review"})
	if done["current_stage"] != aiTutorStageComplete {
		t.Fatalf("current_stage = %#v response=%#v", done["current_stage"], done)
	}
}
```

- [ ] **Step 2: Run tests and verify failure**

Run:

```powershell
go test ./... -run "TestWebAITutor" -count=1
```

Expected: 404 for missing `/api/ai-tutor/*` routes or compile failure if fake engine injection is missing.

- [ ] **Step 3: Add routes**

Modify `web_api.go` route setup:

```go
	mux.HandleFunc("/api/ai-tutor/start", api.handleAITutorStart)
	mux.HandleFunc("/api/ai-tutor/session", api.handleAITutorSession)
	mux.HandleFunc("/api/ai-tutor/answer", api.handleAITutorAnswer)
	mux.HandleFunc("/api/ai-tutor/review", api.handleAITutorReview)
	mux.HandleFunc("/api/ai-tutor/finish", api.handleAITutorFinish)
```

Keep `/api/tutor/start` as a transition alias by changing `handleTutorStart` to call `handleAITutorStart`.

- [ ] **Step 4: Add API handlers**

Add handlers near the old tutor handlers:

```go
func (api *webAPI) handleAITutorStart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	result, err := api.bot.aiTutorEngine().Start(r.Context(), user, "web")
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, api.aiTutorDTO(result, user))
}
```

Add `handleAITutorAnswer`:

```go
func (api *webAPI) handleAITutorAnswer(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		SessionID string `json:"session_id"`
		Text      string `json:"text"`
		Choice    string `json:"choice"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	result, err := api.bot.aiTutorEngine().Submit(r.Context(), user, req.SessionID, aiTutorSubmitInput{Text: req.Text, Choice: req.Choice})
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, api.aiTutorDTO(result, user))
}
```

Add `handleAITutorReview` as a thin wrapper around `Submit` with `Choice`. Add `handleAITutorSession` to load current active session for user if `session_id` is provided; otherwise return 404 with a localized start message. Add `handleAITutorFinish` as an alias for a `no_review` submission when the current stage is `review_schedule`.

Add:

```go
func (api *webAPI) aiTutorDTO(result aiTutorResult, user userState) map[string]any {
	return map[string]any{
		"session":       result.Session,
		"lesson":        result.Lesson.Payload,
		"lesson_status": result.Lesson.Status,
		"current_stage": result.Session.CurrentStage,
		"next_step":     result.NextStep,
		"feedback":      result.Feedback,
		"user":          api.userDTO(user),
	}
}
```

- [ ] **Step 5: Add bot engine accessor**

Add to `bot.go`:

```go
func (b *bot) aiTutorEngine() *aiTutorEngine {
	return newAITutorEngine(b.store, b.openrouter)
}
```

If tests need a fake engine, add an optional field:

```go
	aiTutor *aiTutorEngine
```

and make the accessor return `b.aiTutor` when non-nil.

- [ ] **Step 6: Run Web API tests**

Run:

```powershell
go test ./... -run "TestWebAITutor" -count=1
```

Expected: PASS.

- [ ] **Step 7: Commit Task 7**

```powershell
git add web_api.go web_api_feature_test.go bot.go
git commit -m "Add AI tutor web API"
```

### Task 8: Telegram Full Interactive Flow

**Files:**
- Modify: `bot.go`
- Modify: `telegram_test.go`

- [ ] **Step 1: Write failing Telegram tests**

Replace `TestMenuTutorCallbackStartsTutorLesson` with:

```go
func TestTelegramAITutorMenuStartsInteractiveSession(t *testing.T) {
	var payloads []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		payloads = append(payloads, payload)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	store := &jsonStore{path: filepath.Join(t.TempDir(), "users.json"), users: map[int64]userState{}, aiTutorLessons: map[string]aiTutorLessonRecord{}, aiTutorSessions: map[string]aiTutorSessionRecord{}, aiTutorAnswers: map[string]aiTutorAnswerRecord{}}
	user := userState{TelegramID: 123, FirstName: "Test", InterfaceLanguage: "ru", InterfaceSelected: true, TimezoneSelected: true, LanguageSelected: true, LearningLanguage: "en", Level: "A1"}
	store.users[user.TelegramID] = user
	if err := store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-tg-1", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()}); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}
	err := b.handleCallbackQuery(context.Background(), callbackQuery{ID: "cb-1", From: telegramUser{ID: user.TelegramID, FirstName: user.FirstName}, Data: "menu_tutor"})
	if err != nil {
		t.Fatalf("handleCallbackQuery(menu_tutor) error = %v", err)
	}
	refreshed := store.users[user.TelegramID]
	if !strings.HasPrefix(refreshed.Mode, "ai_tutor:") {
		t.Fatalf("mode = %q, want ai_tutor session", refreshed.Mode)
	}
	text, _ := payloads[len(payloads)-1]["text"].(string)
	if !strings.Contains(text, "A Morning Visit") {
		t.Fatalf("expected story step title, got %q", text)
	}
}
```

Add text and callback tests:

```go
func TestTelegramAITutorTextRoutesToActiveSession(t *testing.T) {
	store := &jsonStore{path: filepath.Join(t.TempDir(), "users.json"), users: map[int64]userState{}, aiTutorLessons: map[string]aiTutorLessonRecord{}, aiTutorSessions: map[string]aiTutorSessionRecord{}, aiTutorAnswers: map[string]aiTutorAnswerRecord{}}
	user := userState{TelegramID: 123, FirstName: "Test", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1", Mode: "ai_tutor:session-1"}
	store.users[user.TelegramID] = user
	_ = store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-tg-2", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()})
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: user.TelegramID, LessonID: "lesson-tg-2", Surface: "telegram", CurrentStage: aiTutorStageStoryIntro, Status: aiTutorSessionActive})
	var sent string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		_ = json.NewDecoder(r.Body).Decode(&payload)
		sent, _ = payload["text"].(string)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}
	if err := b.handleAITutorText(context.Background(), user.TelegramID, user, "continue"); err != nil {
		t.Fatalf("handleAITutorText() error = %v", err)
	}
	session, _, _ := store.getAITutorSession("session-1")
	if session.CurrentStage != aiTutorStageRetell {
		t.Fatalf("stage = %q, want retell", session.CurrentStage)
	}
	if !strings.Contains(sent, "Retell") {
		t.Fatalf("sent text = %q", sent)
	}
}

func TestTelegramAITutorCallbackRoutesToActiveSession(t *testing.T) {
	store := &jsonStore{path: filepath.Join(t.TempDir(), "users.json"), users: map[int64]userState{}, aiTutorLessons: map[string]aiTutorLessonRecord{}, aiTutorSessions: map[string]aiTutorSessionRecord{}, aiTutorAnswers: map[string]aiTutorAnswerRecord{}}
	user := userState{TelegramID: 123, FirstName: "Test", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1", Mode: "ai_tutor:session-1"}
	store.users[user.TelegramID] = user
	_ = store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-tg-3", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()})
	_ = store.createAITutorSession(aiTutorSessionRecord{ID: "session-1", TelegramID: user.TelegramID, LessonID: "lesson-tg-3", Surface: "telegram", CurrentStage: aiTutorStageReviewSchedule, Status: aiTutorSessionActive})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte(`{"ok":true}`)) }))
	defer server.Close()
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}
	err := b.handleCallbackQuery(context.Background(), callbackQuery{ID: "cb-1", From: telegramUser{ID: user.TelegramID, FirstName: user.FirstName}, Data: "ait|session-1|choice|no_review"})
	if err != nil {
		t.Fatalf("handleCallbackQuery() error = %v", err)
	}
	session, _, _ := store.getAITutorSession("session-1")
	if session.Status != aiTutorSessionComplete {
		t.Fatalf("status = %q, want complete", session.Status)
	}
}
```

- [ ] **Step 2: Run tests and verify failure**

Run:

```powershell
go test ./... -run "TestTelegramAITutor" -count=1
```

Expected: failure because Telegram still sends a compact tutor card.

- [ ] **Step 3: Replace `startTutorLesson`**

Modify `startTutorLesson` in `bot.go`:

```go
func (b *bot) startTutorLesson(ctx context.Context, chatID int64, user userState) error {
	result, err := b.aiTutorEngine().Start(ctx, user, "telegram")
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, err.Error(), ui(user))
	}
	if err := b.store.setMode(user.TelegramID, "ai_tutor:"+result.Session.ID); err != nil {
		return err
	}
	return b.sendAITutorStep(ctx, chatID, user, result)
}
```

- [ ] **Step 4: Add Telegram step renderer**

Add:

```go
func (b *bot) sendAITutorStep(ctx context.Context, chatID int64, user userState, result aiTutorResult) error {
	text := formatTelegramAITutorStep(result.NextStep, result.Feedback, user)
	keyboard := aiTutorTelegramKeyboard(result.Session.ID, result.NextStep, ui(user))
	if keyboard != nil {
		return b.telegram.sendInlineMessage(ctx, chatID, text, keyboard)
	}
	return b.telegram.sendMessage(ctx, chatID, text)
}

func formatTelegramAITutorStep(step aiTutorStep, feedback aiTutorFeedback, user userState) string {
	var builder strings.Builder
	if feedback.Message != "" {
		builder.WriteString(feedback.Message)
		builder.WriteString("\n\n")
	}
	if step.Title != "" {
		builder.WriteString(step.Title)
		builder.WriteString("\n\n")
	}
	if step.Instruction != "" {
		builder.WriteString(step.Instruction)
	}
	if step.Kind == "story" && step.Lesson.Story.TextTarget != "" {
		builder.WriteString("\n\n")
		builder.WriteString(step.Lesson.Story.TextTarget)
	}
	if step.Word != nil {
		builder.WriteString("\n\n")
		builder.WriteString(step.Word.Target)
		builder.WriteString(" - ")
		builder.WriteString(step.Word.InterfaceTranslation)
	}
	return strings.TrimSpace(builder.String())
}
```

Add `aiTutorTelegramKeyboard` that returns callback rows for `continue`, `word_learn`, `word_recall`, `rating`, and `review` stages. Use callback data:

```text
ait|<session_id>|choice|<value>
```

- [ ] **Step 5: Route callbacks**

In `handleCallbackQuery` default branch before translator callbacks:

```go
if strings.HasPrefix(query.Data, "ait|") {
	return b.handleAITutorCallback(ctx, chatID, user, query.Data)
}
```

Add:

```go
func (b *bot) handleAITutorCallback(ctx context.Context, chatID int64, user userState, data string) error {
	parts := strings.Split(data, "|")
	if len(parts) < 4 {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
	sessionID := parts[1]
	value := parts[3]
	result, err := b.aiTutorEngine().Submit(ctx, user, sessionID, aiTutorSubmitInput{Choice: value})
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, err.Error(), ui(user))
	}
	if result.Session.Status == aiTutorSessionComplete {
		_ = b.store.setMode(user.TelegramID, "idle")
	}
	return b.sendAITutorStep(ctx, chatID, user, result)
}
```

- [ ] **Step 6: Route text input**

Add a helper and call it from the main text handling path before generic practice or unknown mode handling.

```go
func (b *bot) handleAITutorText(ctx context.Context, chatID int64, user userState, text string) error {
	sessionID := strings.TrimPrefix(user.Mode, "ai_tutor:")
	result, err := b.aiTutorEngine().Submit(ctx, user, sessionID, aiTutorSubmitInput{Text: text})
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, err.Error(), ui(user))
	}
	if result.Session.Status == aiTutorSessionComplete {
		_ = b.store.setMode(user.TelegramID, "idle")
	}
	return b.sendAITutorStep(ctx, chatID, user, result)
}
```

In the existing text path, add:

```go
if strings.HasPrefix(user.Mode, "ai_tutor:") {
	return b.handleAITutorText(ctx, chatID, user, text)
}
```

- [ ] **Step 7: Run Telegram tests**

Run:

```powershell
go test ./... -run "TestTelegramAITutor" -count=1
```

Expected: PASS.

- [ ] **Step 8: Commit Task 8**

```powershell
git add bot.go telegram_test.go
git commit -m "Add Telegram AI tutor session flow"
```

### Task 9: React Server-Driven Tutor View

**Files:**
- Modify: `web-react/src/lib/types.ts`
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Add TypeScript types**

In `web-react/src/lib/types.ts`, add:

```ts
export type AiTutorStep = {
  stage: string;
  kind: string;
  title?: string;
  instruction?: string;
  lesson?: unknown;
  word?: { id?: string; target?: string; interface_translation?: string };
  question?: { id?: string; question_target?: string };
  options?: string[];
};

export type AiTutorSession = {
  ID?: string;
  id?: string;
  CurrentStage?: string;
  current_stage?: string;
  Status?: string;
  status?: string;
};

export type AiTutorResponse = {
  session?: AiTutorSession;
  lesson?: unknown;
  current_stage?: string;
  next_step?: AiTutorStep;
  feedback?: { ok?: boolean; message?: string; json?: string };
  user?: unknown;
};
```

- [ ] **Step 2: Replace local TutorView progress source**

In `web-react/src/App.tsx`, keep the visual shell but replace local stage arrays with server-driven state:

```tsx
const [aiTutorSessionId, setAiTutorSessionId] = useState("");
const [aiTutorStep, setAiTutorStep] = useState<AiTutorStep | null>(null);
const [aiTutorFeedback, setAiTutorFeedback] = useState("");
```

Update `startTutor` to call `/api/ai-tutor/start` and set `next_step`:

```tsx
const payload = await api<AiTutorResponse>("/api/ai-tutor/start", { method: "POST", body: {} });
setAiTutorSessionId(String(payload.session?.id || payload.session?.ID || ""));
setAiTutorStep(payload.next_step || null);
setAiTutorFeedback(payload.feedback?.message || "");
```

Add submit handler:

```tsx
const submitAiTutorStep = async (text: string, choice = "") => {
  const payload = await api<AiTutorResponse>("/api/ai-tutor/answer", {
    method: "POST",
    body: { session_id: aiTutorSessionId, text, choice },
  });
  setAiTutorStep(payload.next_step || null);
  setAiTutorFeedback(payload.feedback?.message || "");
};
```

Render by `aiTutorStep.kind`:

- `story`: show story and continue button.
- `free_text`: show textbox.
- `word_learn`: show word card and continue button.
- `word_recall`: show textbox or choice buttons.
- `rating`: show rating buttons.
- `review`: show review option buttons.
- `complete`: show completion summary.

- [ ] **Step 3: Add Web smoke test**

In `web-react/e2e/web-smoke.spec.ts`, mock responses if the test suite already mocks API. The smoke path:

```ts
test("AI Tutor renders server-driven step", async ({ page }) => {
  await page.route("**/api/ai-tutor/start", async (route) => {
    await route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        session: { id: "s1", current_stage: "story_intro", status: "active" },
        next_step: {
          stage: "story_intro",
          kind: "story",
          title: "A Morning Visit",
          instruction: "Read the story.",
          lesson: { story: { text_target: "Mia wakes up early. She buys bread. She goes home." } },
        },
        feedback: {},
      }),
    });
  });
  await page.goto("/");
  await page.getByText(/AI Tutor|AI Репетитор/).click();
  await expect(page.getByText("A Morning Visit")).toBeVisible();
});
```

Adjust selectors to the existing test app entry path if `page.goto("/")` is not how current tests start.

- [ ] **Step 4: Run frontend checks**

Run:

```powershell
npm --prefix web-react run build
npm --prefix web-react run e2e
```

Expected: build passes and the AI Tutor smoke test passes with existing smoke tests.

- [ ] **Step 5: Commit Task 9**

```powershell
git add web-react/src/lib/types.ts web-react/src/App.tsx web-react/e2e/web-smoke.spec.ts
git commit -m "Add server-driven AI tutor web view"
```

### Task 10: Remove Old Tutor Source From Product Flow

**Files:**
- Modify: `course_tutor.go`
- Modify: `course_tutor_test.go`
- Modify: `web_api.go`
- Modify: `bot.go`
- Modify: `web-react/src/App.tsx`

- [ ] **Step 1: Write regression checks for old-source removal**

Add these regression tests:

```go
func TestWebTutorStartUsesAITutorEngine(t *testing.T) {
	api, store, cookie := newTestWebAPI(t)
	payload := validAITutorLessonPayloadForTest()
	if err := store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-compat", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: payload}); err != nil {
		t.Fatalf("saveAITutorLesson() error = %v", err)
	}
	body := requestJSON(t, api, cookie, http.MethodPost, "/api/tutor/start", map[string]any{})
	if body["next_step"] == nil || body["tutor_lesson"] != nil {
		t.Fatalf("/api/tutor/start did not return new AI Tutor shape: %#v", body)
	}
}

func TestMenuTutorDoesNotUseLocalCourseBank(t *testing.T) {
	var text string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]any
		_ = json.NewDecoder(r.Body).Decode(&payload)
		text, _ = payload["text"].(string)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()
	store := &jsonStore{path: filepath.Join(t.TempDir(), "users.json"), users: map[int64]userState{}, aiTutorLessons: map[string]aiTutorLessonRecord{}, aiTutorSessions: map[string]aiTutorSessionRecord{}, aiTutorAnswers: map[string]aiTutorAnswerRecord{}}
	user := userState{TelegramID: 123, FirstName: "Test", InterfaceLanguage: "ru", InterfaceSelected: true, TimezoneSelected: true, LanguageSelected: true, LearningLanguage: "en", Level: "A1"}
	store.users[user.TelegramID] = user
	_ = store.saveAITutorLesson(aiTutorLessonRecord{ID: "lesson-no-local", LearningLanguage: "en", InterfaceLanguage: "ru", ExactLevel: "A1", LevelBand: "A1-A2", Status: aiTutorStatusApproved, Payload: validAITutorLessonPayloadForTest()})
	b := &bot{store: store, telegram: &telegramClient{baseURL: server.URL, http: server.Client()}}
	if err := b.handleCallbackQuery(context.Background(), callbackQuery{ID: "cb-1", From: telegramUser{ID: user.TelegramID, FirstName: user.FirstName}, Data: "menu_tutor"}); err != nil {
		t.Fatalf("handleCallbackQuery(menu_tutor) error = %v", err)
	}
	for _, bad := range []string{"local-a1-a2-course-core", "Pattern:", "Final word check"} {
		if strings.Contains(text, bad) {
			t.Fatalf("Telegram tutor still uses old local flow marker %q in %q", bad, text)
		}
	}
}
```

- [ ] **Step 2: Run tests and verify failure**

Run:

```powershell
go test ./... -run "Test(MenuTutorDoesNotUseLocalCourseBank|WebTutorStartUsesAITutorEngine)" -count=1
```

Expected: failure until old references are removed or routed to the new engine.

- [ ] **Step 3: Remove old product path**

Keep `course_tutor.go` in the repo only until no compile references require it. Remove these product calls:

```go
buildTutorLessonForSequence(...)
store.nextTutorLesson(...)
formatTelegramTutorLesson(...)
```

from:

- `handleTutorStart`
- `startTutorLesson`
- React Tutor display copy

Make `/api/tutor/start` call `handleAITutorStart` so old frontend clients do not 404 during rollout.

- [ ] **Step 4: Run old-source removal tests**

Run:

```powershell
go test ./... -run "Test(MenuTutorDoesNotUseLocalCourseBank|WebTutorStartUsesAITutorEngine)" -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit Task 10**

```powershell
git add course_tutor.go course_tutor_test.go web_api.go bot.go web-react/src/App.tsx
git commit -m "Route tutor product flow to new AI engine"
```

### Task 11: Full Verification And Push

**Files:**
- All files changed by previous tasks.

- [ ] **Step 1: Run backend tests**

Run:

```powershell
go test ./...
```

Expected: PASS.

- [ ] **Step 2: Run frontend build**

Run:

```powershell
npm --prefix web-react run build
```

Expected: PASS.

- [ ] **Step 3: Run frontend e2e**

Run:

```powershell
npm --prefix web-react run e2e
```

Expected: PASS.

- [ ] **Step 4: Run encoding check**

Run:

```powershell
node tools/check_encoding_artifacts.mjs
```

Expected: `Encoding artifacts check passed.`

- [ ] **Step 5: Inspect final diff**

Run:

```powershell
git status --short
git log --oneline -5
```

Expected: only intended files are changed or the branch is clean after commits. Existing unrelated `web-react/playwright-report/index.html` can remain unstaged if it was present before this work.

- [ ] **Step 6: Push**

Run:

```powershell
git push
```

Expected: branch pushes to `origin/main` unless execution happens on a feature branch.

## Self-Review Checklist

- Spec coverage: Tasks cover new model, validation, prompts, SQLite, engine, Web API, Telegram full flow, React view, old source removal, quality promotion, and verification.
- No old lesson bank: Task 2 clears old rows; Task 10 removes old product routing.
- Shared engine: Tasks 4-8 make Web and Telegram use `aiTutorEngine`.
- Quality gate: Tasks 3, 4, and 6 add preflight and post-lesson checks with score threshold 88.
- Review scheduling: Task 5 stores tomorrow, three days, one week, and no review behavior.
- Resumability: Task 2 stores sessions; Task 8 routes Telegram by `ai_tutor:<session_id>`; Task 7 exposes session endpoint.
