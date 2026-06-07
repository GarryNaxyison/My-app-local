# AI Tutor Scenario Slots Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace AI Tutor word-index-driven scenario generation with curated scenario slots so cafe, hotel, doctor, travel, and service desk lessons produce coherent explanations, choices, and mini-dialogues.

**Architecture:** Add `scenario_slots` to the tutor lesson JSON and introduce internal curated templates in `course_tutor.go`. Keep vocabulary selection unchanged, but build visible teaching content from the scenario template instead of `words[0..3]`.

**Tech Stack:** Go 1.22 backend, existing local tutor lesson generator, React/Vite web app, Playwright e2e tests.

---

## File Structure

- Modify `course_tutor.go`: add slot/template structs, resolve templates by topic, wire template output into success criteria, mini explanation, checks, writing, answer variants, listening, and dialogue.
- Modify `course_tutor_test.go`: add Go regressions for cafe slot behavior, RU service-string leakage, and conversational choice answers.
- Modify `web-react/e2e/web-smoke.spec.ts`: update the tutor mock to a cafe lesson with `football` as an extra word and assert the visible flow stays scenario-based.
- Create/modify no new runtime UI components. `scenario_slots` can be added to TypeScript type metadata later only if TypeScript needs it.

## Task 1: Add Failing Go Regression Tests

**Files:**
- Modify: `course_tutor_test.go`
- Test: `course_tutor_test.go`

- [ ] **Step 1: Add helper functions and failing tests**

Append these tests near the other tutor lesson tests in `course_tutor_test.go`:

```go
func TestTutorCafeScenarioUsesSlotsNotRandomWords(t *testing.T) {
	topic := tutorTestTopicByCode(t, "food")
	function := tutorCourseFunctions[0]
	words := []tutorLessonWord{
		{ID: "en:menu", Word: "menu", Translation: "меню", Level: "A1"},
		{ID: "en:football", Word: "football", Translation: "футбол", Level: "A1"},
		{ID: "en:coffee", Word: "coffee", Translation: "кофе", Level: "A1"},
		{ID: "en:breakfast", Word: "breakfast", Translation: "завтрак", Level: "A1"},
	}

	checks := tutorScenarioChecks(words, topic, function, "ru")
	variants := tutorAnswerVariants(words, topic, function, "en", "ru", false)
	dialogueVariants := tutorAnswerVariants(words, topic, function, "en", "ru", true)

	combined := strings.ToLower(strings.Join(append(tutorChoiceTexts(checks), tutorVariantTexts(append(variants, dialogueVariants...))...), "\n"))
	for _, want := range []string{"coffee", "for breakfast", "could i see the menu"} {
		if !strings.Contains(combined, want) {
			t.Fatalf("cafe scenario content misses %q:\n%s", want, combined)
		}
	}
	for _, bad := range []string{"have football", "like football", "football for breakfast"} {
		if strings.Contains(combined, bad) {
			t.Fatalf("cafe scenario used football as an order/detail via %q:\n%s", bad, combined)
		}
	}
}

func TestTutorRussianMiniExplanationDoesNotLeakServiceCriteria(t *testing.T) {
	user := userState{TelegramID: 91, FirstName: "demo", InterfaceLanguage: "ru", LearningLanguage: "en", Level: "A1"}
	lesson, err := buildTutorLessonForSequence(user, 1)
	if err != nil {
		t.Fatalf("buildTutorLessonForSequence(food) error = %v", err)
	}
	visible := strings.Join(append([]string{lesson.MiniExplanation}, lesson.SuccessCriteria...), "\n")
	for _, bad := range []string{"Use the lesson goal", "Include ", "Variant focus"} {
		if strings.Contains(visible, bad) {
			t.Fatalf("visible RU tutor text leaked service criteria %q:\n%s", bad, visible)
		}
	}
	if !strings.Contains(visible, "попросить меню") || !strings.Contains(visible, "coffee") || !strings.Contains(visible, "for breakfast") {
		t.Fatalf("RU cafe explanation does not expose the scenario pattern:\n%s", visible)
	}
}

func TestTutorScenarioChoicesAreDialogueReplies(t *testing.T) {
	topic := tutorTestTopicByCode(t, "food")
	function := tutorCourseFunctions[0]
	words := []tutorLessonWord{
		{ID: "en:menu", Word: "menu", Translation: "меню", Level: "A1"},
		{ID: "en:football", Word: "football", Translation: "футбол", Level: "A1"},
		{ID: "en:coffee", Word: "coffee", Translation: "кофе", Level: "A1"},
		{ID: "en:breakfast", Word: "breakfast", Translation: "завтрак", Level: "A1"},
	}

	checks := tutorScenarioChecks(words, topic, function, "ru")
	if len(checks) < 2 {
		t.Fatalf("checks = %d, want at least 2 scenario checks", len(checks))
	}
	best := tutorCorrectChoiceText(t, checks[0])
	detail := tutorCorrectChoiceText(t, checks[1])
	if best != "Could I see the menu and have coffee for breakfast, please?" {
		t.Fatalf("best cafe reply = %q", best)
	}
	if detail != "for breakfast" {
		t.Fatalf("detail check = %q, want for breakfast", detail)
	}
	if strings.Contains(best, " / ") || strings.Contains(detail, ":") {
		t.Fatalf("correct choices should be dialogue replies/details, got best=%q detail=%q", best, detail)
	}
}

func tutorTestTopicByCode(t *testing.T, code string) tutorTopic {
	t.Helper()
	for _, topic := range tutorScenarioCourseTopics {
		if topic.Code == code {
			return topic
		}
	}
	t.Fatalf("topic %q not found", code)
	return tutorTopic{}
}

func tutorChoiceTexts(checks []tutorLessonChoice) []string {
	var texts []string
	for _, check := range checks {
		texts = append(texts, check.Prompt)
		for _, option := range check.Options {
			texts = append(texts, option.Text, option.Label, option.Why)
		}
	}
	return texts
}

func tutorVariantTexts(variants []tutorAnswerVariant) []string {
	texts := make([]string, 0, len(variants)*3)
	for _, variant := range variants {
		texts = append(texts, variant.Text, variant.Label, variant.Why)
	}
	return texts
}

func tutorCorrectChoiceText(t *testing.T, choice tutorLessonChoice) string {
	t.Helper()
	for _, option := range choice.Options {
		if option.ID == choice.CorrectAnswerID {
			return option.Text
		}
	}
	t.Fatalf("correct option %q not found in %#v", choice.CorrectAnswerID, choice.Options)
	return ""
}
```

- [ ] **Step 2: Run tests to verify RED**

Run:

```powershell
go test ./... -run 'TestTutorCafeScenarioUsesSlotsNotRandomWords|TestTutorRussianMiniExplanationDoesNotLeakServiceCriteria|TestTutorScenarioChoicesAreDialogueReplies'
```

Expected: FAIL. Current code should mention `football` in a cafe answer, leak at least one English service criterion, or return `word: translation` for the detail check.

## Task 2: Implement Scenario Slots and Template Wiring

**Files:**
- Modify: `course_tutor.go`
- Test: `course_tutor_test.go`

- [ ] **Step 1: Add structs and template resolver**

Add `tutorScenarioSlots`, `tutorScenarioTemplate`, `tutorScenarioChoiceLine`, `tutorScenarioVariantLine`, and `tutorScenarioTemplateFor()` near the tutor model and course topic definitions. The `food` template must include:

```go
TopicCode:       "food",
Role:            "guest",
SituationRU:     "Вы в кафе. Нужно попросить меню, заказать напиток и уточнить завтрак или ужин.",
SituationEN:     "You are in a cafe: ask for the menu, order a drink, and clarify breakfast or dinner.",
RequiredAction:  "ask for the menu",
Item:            "coffee",
Detail:          "for breakfast",
Politeness:      []string{"Could I...?", "I'd like..., please."},
ModelAnswer:     "Could I see the menu and have coffee for breakfast, please?",
LastTutorLine:   "Good morning. Would you like breakfast or dinner?",
```

Also add curated entries for `hotel`, `doctor`, `travel`, and `services`, plus a generic fallback that does not depend on `words[0..3]`.

- [ ] **Step 2: Add `ScenarioSlots` to `tutorLesson`**

Add this field to the struct:

```go
ScenarioSlots tutorScenarioSlots `json:"scenario_slots,omitempty"`
```

In `buildTutorLessonForSequence()`, resolve the template once:

```go
scenarioTemplate := tutorScenarioTemplateFor(topic, courseFunction, interfaceLanguage)
```

Set:

```go
ScenarioSlots: tutorScenarioSlotsFromTemplate(scenarioTemplate, interfaceLanguage),
```

- [ ] **Step 3: Rewire visible lesson builders**

Change the lesson construction to pass the template into builders that need it. The target shape is:

```go
checks := tutorScenarioChecks(words, topic, courseFunction, interfaceLanguage)
choice := checks[0]
...
SuccessCriteria: tutorSuccessCriteria(topic, courseFunction, variant, interfaceLanguage),
MiniExplanation: tutorMiniExplanation(level, topic, courseFunction, words, interfaceLanguage),
WritingTask: tutorWritingTask(first, second, topic, courseFunction, interfaceLanguage),
WritingExpected: tutorScenarioExpectedWords(scenarioTemplate),
AnswerVariants: tutorAnswerVariants(words, topic, courseFunction, language, interfaceLanguage, false),
ListeningText: listeningText,
ListeningExpected: listeningExpected,
DialogueVariants: tutorAnswerVariants(words, topic, courseFunction, language, interfaceLanguage, true),
```

The exported function signatures can stay close to the current ones by resolving the template inside each helper. This keeps existing call sites and tests small.

- [ ] **Step 4: Replace word-index content generation**

Update:

```go
func tutorSuccessCriteria(topic tutorTopic, function tutorCourseFunction, variant tutorCourseVariant, interfaceLanguage string) []string
func tutorScenarioChecks(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) []tutorLessonChoice
func tutorMiniExplanation(level string, topic tutorTopic, function tutorCourseFunction, words []tutorLessonWord, interfaceLanguage string) string
func tutorWritingTask(first tutorLessonWord, second tutorLessonWord, topic tutorTopic, function tutorCourseFunction, interfaceLanguage string) string
func tutorAnswerVariants(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string, dialogue bool) []tutorAnswerVariant
func tutorListeningBlock(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string) (string, string, []string)
func tutorScenarioDialogue(words []tutorLessonWord, topic tutorTopic, function tutorCourseFunction, language string, interfaceLanguage string) ([]string, string, string)
```

For the cafe path, the correct outputs must use `menu`, `coffee`, and `for breakfast` from the template even when `words[1]` is `football`.

- [ ] **Step 5: Run Go tests to verify GREEN**

Run:

```powershell
go test ./... -run 'TestTutorCafeScenarioUsesSlotsNotRandomWords|TestTutorRussianMiniExplanationDoesNotLeakServiceCriteria|TestTutorScenarioChoicesAreDialogueReplies'
```

Expected: PASS.

- [ ] **Step 6: Run the full Go suite**

Run:

```powershell
go test ./...
```

Expected: PASS.

## Task 3: Add Playwright Cafe AI Tutor Regression

**Files:**
- Modify: `web-react/e2e/web-smoke.spec.ts`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Update mocked tutor lesson to cafe slots**

In `mockApi()`, replace the mocked `/api/tutor/start` lesson with a cafe lesson that includes `football` in `words` but uses `coffee` and `for breakfast` in `mini_explanation`, `checks`, `answer_variants`, `dialogue`, and `dialogue_variants`.

The mock should include:

```ts
scenario_slots: {
  role: "guest",
  situation: "Вы в кафе. Нужно попросить меню, заказать напиток и уточнить завтрак или ужин.",
  required_action: "ask for the menu",
  item: "coffee",
  detail: "for breakfast",
  politeness: ["Could I...?", "I'd like..., please."],
  model_answer: "Could I see the menu and have coffee for breakfast, please?",
},
```

- [ ] **Step 2: Update visible assertions in the existing tutor test**

Rename the test to:

```ts
test("AI Tutor cafe scenario uses slots for explanation choices and dialogue", async ({ page }) => {
```

Assert these visible strings:

```ts
await expect(page.locator(".tutor-context-v2")).toContainText("В кафе нужен короткий порядок");
await expect(page.locator(".tutor-context-v2")).toContainText("Could I see the menu?");
await expect(page.locator(".tutor-context-v2")).not.toContainText("Use the lesson goal");
await expect(page.locator(".tutor-context-v2")).not.toContainText("Variant focus");
await expect(page.locator(".tutor-choice-grid-v2")).toContainText("Could I see the menu and have coffee for breakfast, please?");
await expect(page.locator(".tutor-choice-grid-v2")).toContainText("I need a ticket.");
await expect(page.locator(".tutor-choice-grid-v2")).toContainText("for breakfast");
await expect(page.locator(".tutor-dialogue")).toContainText("Would you like breakfast or dinner?");
await expect(page.locator(".tutor-answer-variants-v2")).toContainText("I'd like coffee for breakfast, please.");
await expect(page.locator(".tutor-answer-variants-v2")).not.toContainText("have football");
```

- [ ] **Step 3: Run Playwright test to verify RED if backend-only work was not wired into mock yet**

Run:

```powershell
npm --prefix web-react test -- --grep "AI Tutor cafe scenario uses slots"
```

Expected before mock update: FAIL because the old hotel mock lacks cafe text. Expected after mock update: PASS.

## Task 4: Final Verification and Sync

**Files:**
- Modify: all changed files
- Test: all changed tests

- [ ] **Step 1: Run focused Go tests**

Run:

```powershell
go test ./... -run 'Tutor'
```

Expected: PASS.

- [ ] **Step 2: Run full Go tests**

Run:

```powershell
go test ./...
```

Expected: PASS.

- [ ] **Step 3: Run focused Playwright test**

Run:

```powershell
npm --prefix web-react test -- --grep "AI Tutor cafe scenario uses slots"
```

Expected: PASS.

- [ ] **Step 4: Review diff**

Run:

```powershell
git diff -- course_tutor.go course_tutor_test.go web-react/e2e/web-smoke.spec.ts
```

Expected: changes are scoped to scenario slots, tests, and the tutor e2e mock.

- [ ] **Step 5: Commit and push**

Run:

```powershell
git add course_tutor.go course_tutor_test.go web-react/e2e/web-smoke.spec.ts docs/superpowers/plans/2026-06-07-ai-tutor-scenario-slots.md
git commit -m "fix: generate tutor lessons from scenario slots"
git push origin main
```

Expected: commit succeeds and pushes to `https://github.com/GarryNaxyison/My-app-local.git`.
