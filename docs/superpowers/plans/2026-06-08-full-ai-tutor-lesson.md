# Full AI Tutor Lesson Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the approved full AI Tutor lesson flow: teaching point, slot-aware practice, final word check, summary, web rendering, Telegram rendering, and regression coverage.

**Architecture:** Keep `course_tutor.go` as the deterministic lesson source and add backward-compatible JSON fields to `tutorLesson`. Web `TutorView` consumes the new fields when present and falls back to existing lesson fields. Telegram remains a compact card that previews the same shared backend structure.

**Tech Stack:** Go backend, React/TypeScript web app, Playwright e2e, existing local tutor course bank, existing Telegram MarkdownV2 rendering.

---

## File Structure

- Modify `course_tutor.go`: add lesson-session structs, attach them to `tutorLesson`, build teaching points, final word checks, slot feedback, and summary from scenario slots.
- Modify `course_tutor_test.go`: add backend red tests for doctor teaching point, final word check, slot feedback, and lesson-step alignment.
- Modify `bot.go`: include teaching point and final-check preview in Telegram tutor cards.
- Modify `telegram_test.go`: add Telegram rendering assertions for teaching point and final-check preview.
- Modify `web-react/src/App.tsx`: add TypeScript fields, final-check state, slot-aware feedback, teaching point rendering, final-check rendering, summary rendering, and updated stage list.
- Modify `web-react/src/styles/app.css`: add compact final-check and summary styling only if existing choice/review styles are not enough.
- Modify `web-react/e2e/web-smoke.spec.ts`: update AI Tutor mocked lesson and e2e flow to cover teaching point, final check before SRS, and summary.

## Task 1: Backend Lesson Model And Teaching Point

**Files:**
- Modify: `course_tutor.go`
- Test: `course_tutor_test.go`

- [ ] **Step 1: Write the failing backend test**

Add this test near the other tutor scenario tests in `course_tutor_test.go`:

```go
func TestTutorDoctorLessonHasTeachingPointAndFinalWordCheck(t *testing.T) {
	user := userState{TelegramID: 191, FirstName: "demo", InterfaceLanguage: "en", LearningLanguage: "en", Level: "A1"}
	lesson, err := buildTutorLessonForSequence(user, 4)
	if err != nil {
		t.Fatalf("buildTutorLessonForSequence(doctor) error = %v", err)
	}
	if lesson.TeachingPoint.Pattern != "Hello. I need to see a doctor today, please." {
		t.Fatalf("teaching pattern = %q", lesson.TeachingPoint.Pattern)
	}
	wantOrder := []string{"action", "doctor", "today"}
	if strings.Join(lesson.TeachingPoint.SceneOrder, "|") != strings.Join(wantOrder, "|") {
		t.Fatalf("scene order = %#v, want %#v", lesson.TeachingPoint.SceneOrder, wantOrder)
	}
	if lesson.TeachingPoint.ModelAnswer != "Hello. I need to see a doctor today, please." {
		t.Fatalf("teaching model answer = %q", lesson.TeachingPoint.ModelAnswer)
	}
	if len(lesson.FinalWordCheck.Items) < 3 {
		t.Fatalf("final word check items = %d, want at least 3", len(lesson.FinalWordCheck.Items))
	}
	if lesson.FinalWordCheck.RequiredCorrect < 2 {
		t.Fatalf("required correct = %d, want at least 2", lesson.FinalWordCheck.RequiredCorrect)
	}
	combined := strings.ToLower(strings.Join(tutorChoiceTexts(lesson.FinalWordCheck.Items), "\n"))
	for _, want := range []string{"doctor", "today"} {
		if !strings.Contains(combined, want) {
			t.Fatalf("final word check misses %q:\n%s", want, combined)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```powershell
go test ./... -run TestTutorDoctorLessonHasTeachingPointAndFinalWordCheck
```

Expected: compile failure because `TeachingPoint` and `FinalWordCheck` are not defined on `tutorLesson`.

- [ ] **Step 3: Add backend JSON structs**

In `course_tutor.go`, after `type tutorScenarioSlots struct`, add:

```go
type tutorTeachingPoint struct {
	Title       string   `json:"title"`
	Pattern     string   `json:"pattern"`
	SceneOrder  []string `json:"scene_order"`
	Explanation string   `json:"explanation"`
	ModelAnswer string   `json:"model_answer"`
}

type tutorFinalWordCheck struct {
	Prompt          string              `json:"prompt"`
	Items           []tutorLessonChoice `json:"items"`
	RequiredCorrect int                 `json:"required_correct"`
	SummaryPass     string              `json:"summary_pass"`
	SummaryRetry    string              `json:"summary_retry"`
}

type tutorSummaryBlock struct {
	CanSay      string   `json:"can_say"`
	StrongItems []string `json:"strong_items"`
	WeakItems   []string `json:"weak_items"`
	NextReview  string   `json:"next_review"`
}
```

In `type tutorLesson struct`, after `ScenarioSlots`, add:

```go
TeachingPoint    tutorTeachingPoint  `json:"teaching_point,omitempty"`
FinalWordCheck   tutorFinalWordCheck `json:"final_word_check,omitempty"`
TutorSummary     tutorSummaryBlock   `json:"tutor_summary,omitempty"`
```

- [ ] **Step 4: Add teaching point helper**

In `course_tutor.go`, after `tutorScenarioExpectedWords`, add:

```go
func tutorTeachingPointFor(template tutorScenarioTemplate, interfaceLanguage string) tutorTeachingPoint {
	item := strings.TrimSpace(template.Item)
	detail := strings.TrimSpace(template.Detail)
	model := strings.TrimSpace(template.ModelAnswer)
	order := []string{"action"}
	if item != "" {
		order = append(order, item)
	}
	if detail != "" {
		order = append(order, detail)
	}
	title := tutorLocalized(interfaceLanguage, "One useful pattern", "One useful pattern")
	explanation := tutorLocalized(interfaceLanguage, template.MiniExplanationRU, template.MiniExplanationEN)
	if strings.TrimSpace(explanation) == "" {
		explanation = fmt.Sprintf("Scene order: %s. Pattern: %s", strings.Join(order, " -> "), model)
	}
	return tutorTeachingPoint{
		Title:       title,
		Pattern:     model,
		SceneOrder:  order,
		Explanation: explanation,
		ModelAnswer: model,
	}
}
```

- [ ] **Step 5: Add final word check helper**

In `course_tutor.go`, after `tutorTeachingPointFor`, add:

```go
func tutorFinalWordCheck(words []tutorLessonWord, template tutorScenarioTemplate, interfaceLanguage string) tutorFinalWordCheck {
	items := tutorWordCheckItems(words, template, interfaceLanguage)
	required := tutorMinInt(len(items), tutorMaxInt(2, len(items)-1))
	return tutorFinalWordCheck{
		Prompt:          tutorLocalized(interfaceLanguage, "Final word check: prove you still remember the lesson words.", "Final word check: prove you still remember the lesson words."),
		Items:           items,
		RequiredCorrect: required,
		SummaryPass:     tutorLocalized(interfaceLanguage, "Word check passed. Now schedule the review.", "Word check passed. Now schedule the review."),
		SummaryRetry:    tutorLocalized(interfaceLanguage, "Repeat the weak words before finishing the lesson.", "Repeat the weak words before finishing the lesson."),
	}
}

func tutorWordCheckItems(words []tutorLessonWord, template tutorScenarioTemplate, interfaceLanguage string) []tutorLessonChoice {
	items := make([]tutorLessonChoice, 0, 4)
	seen := map[string]bool{}
	addWord := func(wordText string) {
		wordText = strings.TrimSpace(wordText)
		if wordText == "" || seen[strings.ToLower(wordText)] {
			return
		}
		for _, word := range words {
			if !strings.EqualFold(strings.TrimSpace(word.Word), wordText) {
				continue
			}
			items = append(items, tutorChoiceFromWords(append([]tutorLessonWord{word}, words...), interfaceLanguage))
			seen[strings.ToLower(wordText)] = true
			return
		}
	}
	addWord(template.Item)
	addWord(template.Detail)
	for _, word := range words {
		addWord(word.Word)
		if len(items) >= 4 {
			break
		}
	}
	return items
}
```

- [ ] **Step 6: Add summary helper**

In `course_tutor.go`, after `tutorFinalWordCheck`, add:

```go
func tutorSummaryForLesson(words []tutorLessonWord, template tutorScenarioTemplate, function tutorCourseFunction, interfaceLanguage string) tutorSummaryBlock {
	strong := make([]string, 0, 4)
	for _, item := range []string{template.Item, template.Detail} {
		item = strings.TrimSpace(item)
		if item != "" {
			strong = append(strong, item)
		}
	}
	for _, word := range words {
		if len(strong) >= 4 {
			break
		}
		text := strings.TrimSpace(word.Word)
		if text != "" {
			strong = append(strong, text)
		}
	}
	return tutorSummaryBlock{
		CanSay:      strings.TrimSpace(template.ModelAnswer),
		StrongItems: strong,
		WeakItems:   []string{},
		NextReview:  tutorLocalized(interfaceLanguage, "Choose a review button based on how hard the final check felt.", "Choose a review button based on how hard the final check felt."),
	}
}
```

- [ ] **Step 7: Attach new fields in lesson builder**

In `buildTutorLessonForSequence`, inside the returned `tutorLesson{...}`, add these fields after `ScenarioSlots`:

```go
TeachingPoint:    tutorTeachingPointFor(scenarioTemplate, interfaceLanguage),
FinalWordCheck:   tutorFinalWordCheck(words, scenarioTemplate, interfaceLanguage),
TutorSummary:     tutorSummaryForLesson(words, scenarioTemplate, courseFunction, interfaceLanguage),
```

- [ ] **Step 8: Run test to verify it passes**

Run:

```powershell
go test ./... -run TestTutorDoctorLessonHasTeachingPointAndFinalWordCheck
```

Expected: PASS.

- [ ] **Step 9: Commit backend model**

Run:

```powershell
git add course_tutor.go course_tutor_test.go
git commit -m "feat: add tutor lesson session fields"
```

## Task 2: Backend Slot Feedback And Step Alignment

**Files:**
- Modify: `course_tutor.go`
- Test: `course_tutor_test.go`

- [ ] **Step 1: Write failing feedback and step tests**

Add this test near the backend tutor tests:

```go
func TestTutorAnswerSlotFeedbackNamesMissingDetail(t *testing.T) {
	topic := tutorTestTopicByCode(t, "doctor")
	template := tutorScenarioTemplateFor(topic, tutorCourseFunctions[0], "en")
	feedback := tutorAnswerSlotFeedback("Hello. I need to see a doctor, please.", template, "en")
	if !strings.Contains(strings.ToLower(feedback), "today") {
		t.Fatalf("feedback should name missing detail today, got %q", feedback)
	}
	if !strings.Contains(strings.ToLower(feedback), "doctor") {
		t.Fatalf("feedback should mention present item doctor, got %q", feedback)
	}
}

func TestTutorLessonStepsIncludeFinalCheckBeforeReview(t *testing.T) {
	steps := tutorLessonSteps("en")
	var codes []string
	for _, step := range steps {
		codes = append(codes, step.Code)
	}
	got := strings.Join(codes, "|")
	want := "words|explain|choice|writing|listening|pronunciation|dialogue|final-check|review"
	if got != want {
		t.Fatalf("tutor step codes = %s, want %s", got, want)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run:

```powershell
go test ./... -run "TestTutorAnswerSlotFeedbackNamesMissingDetail|TestTutorLessonStepsIncludeFinalCheckBeforeReview"
```

Expected: compile failure for `tutorAnswerSlotFeedback`, then step-order failure after the function exists.

- [ ] **Step 3: Add slot feedback helper**

In `course_tutor.go`, after `tutorSummaryForLesson`, add:

```go
func tutorAnswerSlotFeedback(answer string, template tutorScenarioTemplate, interfaceLanguage string) string {
	normalized := strings.ToLower(answer)
	present := make([]string, 0, 3)
	missing := make([]string, 0, 3)
	for _, slot := range []string{template.Item, template.Detail} {
		slot = strings.TrimSpace(slot)
		if slot == "" {
			continue
		}
		if strings.Contains(normalized, strings.ToLower(slot)) {
			present = append(present, slot)
		} else {
			missing = append(missing, slot)
		}
	}
	if len(missing) == 0 {
		return tutorLocalized(interfaceLanguage, "Good: your answer has the required scene words.", "Good: your answer has the required scene words.")
	}
	return tutorLocalized(interfaceLanguage, "Good parts: ", "Good parts: ") + strings.Join(present, ", ") + ". " +
		tutorLocalized(interfaceLanguage, "Add: ", "Add: ") + strings.Join(missing, ", ") + "."
}
```

- [ ] **Step 4: Update backend lesson steps**

In `tutorLessonSteps`, add the `final-check` step before `review`.

For the Russian branch, use:

```go
{Code: "final-check", Title: "Final word check", Summary: "lesson words"},
{Code: "review", Title: "SRS", Summary: "memory schedule"},
```

For the English branch, use:

```go
{Code: "final-check", Title: "Final word check", Summary: "lesson words"},
{Code: "review", Title: "SRS", Summary: "memory schedule"},
```

- [ ] **Step 5: Run tests to verify they pass**

Run:

```powershell
go test ./... -run "TestTutorAnswerSlotFeedbackNamesMissingDetail|TestTutorLessonStepsIncludeFinalCheckBeforeReview"
```

Expected: PASS.

- [ ] **Step 6: Commit backend feedback and steps**

Run:

```powershell
git add course_tutor.go course_tutor_test.go
git commit -m "feat: add tutor final check step"
```

## Task 3: Telegram Tutor Card Uses New Lesson Fields

**Files:**
- Modify: `bot.go`
- Test: `telegram_test.go`

- [ ] **Step 1: Write failing Telegram rendering test**

Add this test to `telegram_test.go`:

```go
func TestFormatTelegramTutorLessonIncludesTeachingPointAndFinalCheck(t *testing.T) {
	user := userState{InterfaceLanguage: "en", LearningLanguage: "en", Level: "A1"}
	lesson, err := buildTutorLessonForSequence(user, 4)
	if err != nil {
		t.Fatalf("buildTutorLessonForSequence(doctor): %v", err)
	}
	text := formatTelegramTutorLesson(lesson, user)
	for _, want := range []string{
		"Hello\\. I need to see a doctor today, please\\.",
		"Final word check",
		"doctor",
		"today",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("telegram tutor lesson missing %q:\n%s", want, text)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```powershell
go test ./... -run TestFormatTelegramTutorLessonIncludesTeachingPointAndFinalCheck
```

Expected: failure because final-check preview is not rendered yet.

- [ ] **Step 3: Render teaching point in Telegram**

In `formatTelegramTutorLesson`, replace:

```go
appendTelegramTutorLine(&builder, "Pattern:", lesson.MiniExplanation)
```

with:

```go
teaching := strings.TrimSpace(lesson.TeachingPoint.Explanation)
if teaching == "" {
	teaching = lesson.MiniExplanation
}
appendTelegramTutorLine(&builder, "Pattern:", teaching)
if strings.TrimSpace(lesson.TeachingPoint.ModelAnswer) != "" {
	appendTelegramTutorLine(&builder, "Model:", lesson.TeachingPoint.ModelAnswer)
}
```

- [ ] **Step 4: Add Telegram final-check preview helper**

After `appendTelegramTutorDialogue`, add:

```go
func appendTelegramTutorFinalCheck(builder *strings.Builder, lesson tutorLesson) {
	if len(lesson.FinalWordCheck.Items) == 0 {
		return
	}
	builder.WriteString("\n\n")
	builder.WriteString(escapeMarkdownV2("Final word check"))
	limit := tutorMinInt(3, len(lesson.FinalWordCheck.Items))
	for index := 0; index < limit; index++ {
		item := lesson.FinalWordCheck.Items[index]
		text := strings.TrimSpace(item.Prompt)
		if text == "" {
			continue
		}
		builder.WriteString("\n")
		builder.WriteString(escapeMarkdownV2("• " + text))
	}
}
```

Call it in `formatTelegramTutorLesson` after `appendTelegramTutorDialogue(&builder, lesson)`:

```go
appendTelegramTutorFinalCheck(&builder, lesson)
```

- [ ] **Step 5: Run test to verify it passes**

Run:

```powershell
go test ./... -run TestFormatTelegramTutorLessonIncludesTeachingPointAndFinalCheck
```

Expected: PASS.

- [ ] **Step 6: Commit Telegram rendering**

Run:

```powershell
git add bot.go telegram_test.go
git commit -m "feat: show tutor teaching point in telegram"
```

## Task 4: Web Types And Final-Check State

**Files:**
- Modify: `web-react/src/App.tsx`
- Test: `web-react/e2e/web-smoke.spec.ts`

- [ ] **Step 1: Update mocked tutor lesson for e2e red test**

In `web-react/e2e/web-smoke.spec.ts`, inside the mocked `tutor_lesson`, add:

```ts
teaching_point: {
  title: "One pattern",
  pattern: "Could I see the menu and have coffee for breakfast, please?",
  scene_order: ["action", "coffee", "for breakfast"],
  explanation: "Scene order: action, then coffee, then for breakfast. Pattern: Could I see the menu and have coffee for breakfast, please.",
  model_answer: "Could I see the menu and have coffee for breakfast, please?",
},
final_word_check: {
  prompt: "Final word check",
  required_correct: 2,
  summary_pass: "Words checked. Now choose review timing.",
  summary_retry: "Repeat the weak words before finishing.",
  items: [
    {
      prompt: "Choose the lesson word: coffee",
      correct_answer_id: "en:coffee",
      options: [
        { id: "en:coffee", text: "coffee", label: "Correct", why: "Key order word." },
        { id: "en:menu", text: "menu", label: "Close, but not this", why: "This is another step." },
      ],
    },
    {
      prompt: "Choose the lesson word: breakfast",
      correct_answer_id: "en:breakfast",
      options: [
        { id: "en:breakfast", text: "breakfast", label: "Correct", why: "Key detail." },
        { id: "en:football", text: "football", label: "Unrelated", why: "Not an order detail." },
      ],
    },
  ],
},
tutor_summary: {
  can_say: "Could I see the menu and have coffee for breakfast, please?",
  strong_items: ["coffee", "for breakfast"],
  weak_items: [],
  next_review: "Choose how hard it was to recall the words.",
},
```

Then update the e2e expected step count from:

```ts
await expect(page.locator(".tutor-step-v2")).toHaveCount(8);
```

to:

```ts
await expect(page.locator(".tutor-step-v2")).toHaveCount(9);
```

After dialogue, update the flow to expect final check before SRS:

```ts
await tutorSubmit.click();
await expect(page.locator(".tutor-context-v2")).toContainText("Final word check");
await page.locator(".tutor-choice-grid-v2 button").filter({ hasText: "coffee" }).click();
await tutorSubmit.click();
await page.locator(".tutor-choice-grid-v2 button").filter({ hasText: "breakfast" }).click();
await tutorSubmit.click();
await expect(page.locator(".tutor-context-v2")).toContainText("Повторение");
```

- [ ] **Step 2: Run e2e red test**

Run:

```powershell
npm --prefix web-react run e2e -- --grep "AI Tutor cafe scenario uses slots"
```

Expected: failure because the app has no final-check stage or type support yet.

- [ ] **Step 3: Add TypeScript lesson types**

In `web-react/src/App.tsx`, after `type TutorLessonChoice`, add:

```ts
type TutorTeachingPoint = {
  title?: string;
  pattern?: string;
  scene_order?: string[];
  explanation?: string;
  model_answer?: string;
};

type TutorFinalWordCheck = {
  prompt?: string;
  items?: TutorLessonChoice[];
  required_correct?: number;
  summary_pass?: string;
  summary_retry?: string;
};

type TutorSummaryBlock = {
  can_say?: string;
  strong_items?: string[];
  weak_items?: string[];
  next_review?: string;
};
```

In `type TutorLesson`, after `scenario?: string;`, add:

```ts
teaching_point?: TutorTeachingPoint;
final_word_check?: TutorFinalWordCheck;
tutor_summary?: TutorSummaryBlock;
```

Change:

```ts
type TutorStageId = "words" | "explain" | "choice" | "writing" | "listening" | "pronunciation" | "dialogue" | "review";
```

to:

```ts
type TutorStageId = "words" | "explain" | "choice" | "writing" | "listening" | "pronunciation" | "dialogue" | "final-check" | "review";
```

- [ ] **Step 4: Add final-check React state**

In `TutorView`, near choice state, add:

```ts
const [finalCheckIndex, setFinalCheckIndex] = useState(0);
const [finalCheckSelection, setFinalCheckSelection] = useState("");
const [finalCheckResults, setFinalCheckResults] = useState<Array<{ id: string; answer: string }>>([]);
```

In the lesson reset effect, add:

```ts
setFinalCheckIndex(0);
setFinalCheckSelection("");
setFinalCheckResults([]);
```

Add derived values after `choiceOptions`:

```ts
const finalCheckItems = useMemo(
  () => (tutorLesson?.final_word_check?.items || []).filter((item) => (item.options || []).length && item.correct_answer_id),
  [tutorLesson],
);
const activeFinalCheck = finalCheckItems[Math.min(finalCheckIndex, Math.max(finalCheckItems.length - 1, 0))];
const finalCheckOptions = useMemo(
  () => tutorStableShuffle(activeFinalCheck?.options || [], `${tutorLesson?.id || "tutor"}:${activeFinalCheck?.prompt || "final"}:${finalCheckIndex}`),
  [activeFinalCheck, tutorLesson?.id, finalCheckIndex],
);
```

- [ ] **Step 5: Add final-check stage entry**

In `tutorStages`, insert before review:

```ts
{
  id: "final-check",
  title: copy("tutor_step_final_check", "Final word check"),
  instruction: copy("tutor_final_check_instruction", "Check the lesson words once more before scheduling review."),
},
```

- [ ] **Step 6: Run TypeScript build red check**

Run:

```powershell
npm --prefix web-react run build
```

Expected: fail if any new state/type names are inconsistent; otherwise proceed to Task 5.

## Task 5: Web Final-Check Behavior And Summary Rendering

**Files:**
- Modify: `web-react/src/App.tsx`
- Modify: `web-react/src/styles/app.css`

- [ ] **Step 1: Add final-check action branch**

In `handleTutorAction`, before the `review` branch, add:

```ts
if (activeStage.id === "final-check") {
  if (!activeFinalCheck) {
    advanceTutorStage(copy("tutor_final_check_result", "Final word check completed"));
    return;
  }
  if (!finalCheckSelection) {
    setFeedback(copy("tutor_word_need_choice", "Choose the translation first."));
    setFeedbackTone("error");
    return;
  }
  if (finalCheckSelection !== activeFinalCheck.correct_answer_id) {
    setFeedback(display(tutorLesson.final_word_check?.summary_retry) || copy("tutor_feedback_wrong", "Not yet. Look at the theme words and choose again."));
    setFeedbackTone("error");
    return;
  }
  const option = finalCheckOptions.find((item) => item.id === finalCheckSelection);
  const nextResults = [...finalCheckResults, { id: activeFinalCheck.correct_answer_id || finalCheckSelection, answer: display(option?.text) }];
  setFinalCheckResults(nextResults);
  if (finalCheckIndex + 1 < finalCheckItems.length) {
    setFinalCheckIndex((index) => index + 1);
    setFinalCheckSelection("");
    setFeedback(copy("tutor_word_correct_next", "Correct. Next word."));
    setFeedbackTone("success");
    return;
  }
  advanceTutorStage(display(tutorLesson.final_word_check?.summary_pass) || `${copy("tutor_final_check_result", "Final word check completed")}: ${nextResults.length}/${Math.max(finalCheckItems.length, 1)}`);
  return;
}
```

Update `primaryActionLabel`:

```ts
if (activeStage.id === "final-check") return copy("tutor_check_answer", "Check answer");
```

Update button icon condition:

```tsx
{tutorVoiceChecking ? <Spinner size="small" className="button-spinner-v2" /> : activeStage?.id === "choice" || activeStage?.id === "final-check" ? <Check size={16} /> : activeStage?.id === "review" ? <ListChecks size={16} /> : <ChevronRight size={16} />}
```

- [ ] **Step 2: Render teaching point**

In `renderTutorMaterial`, replace explain material with:

```tsx
if (activeStage.id === "explain") {
  const point = tutorLesson.teaching_point;
  const order = (point?.scene_order || []).filter(Boolean).join(" -> ");
  return (
    <div className="tutor-lesson-card-v2 tutor-mini-explanation-v2">
      {point?.title ? <strong>{display(point.title)}</strong> : null}
      <p>{display(point?.explanation || tutorLesson.mini_explanation || tutorLesson.grammar || tutorLesson.scenario || "")}</p>
      {order ? <span className="tutor-pattern-line-v2">{order}</span> : null}
      {point?.pattern ? <b className="tutor-pattern-line-v2">{display(point.pattern)}</b> : null}
    </div>
  );
}
```

- [ ] **Step 3: Render final check**

In `renderTutorMaterial`, before review rendering, add:

```tsx
if (activeStage.id === "final-check") {
  return (
    <>
      <p className="tutor-task-copy-v2">
        <strong>{copy("tutor_final_check_progress", "Final word")} {Math.min(finalCheckIndex + 1, Math.max(finalCheckItems.length, 1))}/{Math.max(finalCheckItems.length, 1)}</strong>
        <span>{display(activeFinalCheck?.prompt) || display(tutorLesson.final_word_check?.prompt) || copy("tutor_final_check_instruction", "Check the lesson words once more.")}</span>
      </p>
      <div className="choice-grid-v2 tutor-choice-grid-v2">
        {finalCheckOptions.map((option) => (
          <button
            key={option.id}
            type="button"
            className={cn(
              finalCheckSelection === option.id && "is-selected",
              feedbackTone === "error" && finalCheckSelection === option.id && "is-wrong",
            )}
            onClick={() => {
              setFinalCheckSelection(option.id);
              setFeedback("");
              setFeedbackTone("idle");
            }}
          >
            {option.label ? <small>{display(option.label)}</small> : null}
            <span>{display(option.text)}</span>
            {option.why ? <em>{display(option.why)}</em> : null}
          </button>
        ))}
      </div>
      {finalCheckResults.length ? (
        <div className="tutor-result-list-v2">
          {finalCheckResults.map((item) => <span key={item.id}>{item.answer}</span>)}
        </div>
      ) : null}
    </>
  );
}
```

- [ ] **Step 4: Render summary fields**

In the `viewingCompleteSummary` block, before `review_summary`, add:

```tsx
{tutorLesson.tutor_summary?.can_say ? <p>{display(tutorLesson.tutor_summary.can_say)}</p> : null}
{(tutorLesson.tutor_summary?.strong_items || []).length ? (
  <div className="tutor-review-summary-v2">
    {(tutorLesson.tutor_summary?.strong_items || []).map((item) => <span key={item}>{display(item)}</span>)}
  </div>
) : null}
```

In the review stage rendering, before the SRS buttons, add:

```tsx
{tutorLesson.tutor_summary?.next_review ? <p className="tutor-task-copy-v2">{display(tutorLesson.tutor_summary.next_review)}</p> : null}
```

- [ ] **Step 5: Add CSS for pattern lines if needed**

In `web-react/src/styles/app.css`, near `.tutor-lesson-card-v2`, add:

```css
.tutor-pattern-line-v2 {
  display: block;
  margin-top: 10px;
  color: var(--text-strong);
  font-weight: 700;
  line-height: 1.45;
}
```

- [ ] **Step 6: Run e2e target to verify green**

Run:

```powershell
npm --prefix web-react run e2e -- --grep "AI Tutor cafe scenario uses slots"
```

Expected: PASS.

- [ ] **Step 7: Run build**

Run:

```powershell
npm --prefix web-react run build
```

Expected: PASS.

- [ ] **Step 8: Commit web final check**

Run:

```powershell
git add web-react/src/App.tsx web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts
git commit -m "feat: add tutor final word check"
```

## Task 6: Focused Full Verification And Sync

**Files:**
- No planned code edits unless verification reveals a regression.

- [ ] **Step 1: Run Go tutor tests**

Run:

```powershell
go test ./... -run "Tutor|menu_tutor|AITutor"
```

Expected: PASS.

- [ ] **Step 2: Run full Go test suite**

Run:

```powershell
go test ./...
```

Expected: PASS.

- [ ] **Step 3: Run frontend build**

Run:

```powershell
npm --prefix web-react run build
```

Expected: PASS.

- [ ] **Step 4: Run focused e2e**

Run:

```powershell
npm --prefix web-react run e2e -- --grep "AI Tutor"
```

Expected: PASS.

- [ ] **Step 5: Run encoding artifact check**

Run:

```powershell
node tools/check_encoding_artifacts.mjs
```

Expected: PASS with no reported artifacts.

- [ ] **Step 6: Inspect final diff**

Run:

```powershell
git status --short
git diff --stat
```

Expected: only planned implementation files are modified, plus unrelated pre-existing dirty files that were present before the implementation.

- [ ] **Step 7: Commit any verification-only fixes**

If a verification command exposed a defect and a fix was needed, commit the fix with:

```powershell
git add course_tutor.go course_tutor_test.go bot.go telegram_test.go web-react/src/App.tsx web-react/src/styles/app.css web-react/e2e/web-smoke.spec.ts
git commit -m "fix: stabilize tutor final lesson flow"
```

Do not stage unrelated pre-existing files.

- [ ] **Step 8: Push implementation commits**

Run:

```powershell
git push
```

Expected: commits push to `https://github.com/GarryNaxyison/My-app-local.git`.
