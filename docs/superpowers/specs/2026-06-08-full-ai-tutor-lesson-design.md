# Full AI Tutor Lesson Design

## Status

Approved for design on 2026-06-08.

## Goal

Rebuild AI Tutor from a sequence of exercise blocks into a complete guided tutor lesson. A lesson must teach a theme, introduce useful words, explain one practical pattern, guide the learner through controlled practice, ask for a learner-produced answer, include listening and pronunciation, run a short dialogue, and finish with an actual word check.

The lesson should feel like a tutor session, not like disconnected cards.

User-approved target example:

```text
Pattern: Hello. I need to see a doctor today, please.
Scene order: action, then doctor, then today.
```

## Current Problems

The current AI Tutor already has a local course bank, scenario slots, web rendering, Telegram rendering, pronunciation scoring, and tests. The issue is the learning shape:

- The visible flow is still mostly "step list plus task blocks".
- The final review is mainly SRS scheduling, not a final check of lesson vocabulary.
- Vocabulary is introduced and checked early, but the end of the lesson does not prove that the learner can still recall the words.
- Tutor feedback is generic in some stages and does not consistently name what the learner used or missed.
- The backend `tutorLessonSteps` and the frontend `TutorStageId` are not conceptually aligned enough for a single source of truth.
- Telegram gets a compact lesson card, but it does not yet represent the new full lesson arc.

## Research References

The redesign follows common language-learning patterns from real apps and teaching practice:

- Duolingo describes end-of-lesson mistake review, personalized practice, spaced repetition, and active recall as important for retaining words and grammar: https://blog.duolingo.com/spaced-repetition-for-learning/
- Babbel Review uses spaced repetition after lessons and offers vocabulary review through flashcards, writing, speaking, and listening: https://support.babbel.com/hc/en-us/articles/205600228-Review
- Busuu Conversations are level-appropriate, tied to lesson content, based on real scenarios, and finish with personalized feedback: https://www.busuu.com/en/languages/language-learning-with-busuu-conversations
- British Council TeachingEnglish describes a practical move from lesson introduction and language input into controlled practice, freer practice, review, and feedback: https://www.teachingenglish.org.uk/professional-development/teachers/planning-lessons-and-courses/articles/lesson-plans
- British Council also distinguishes controlled practice from free practice: controlled practice restricts the target language, while free practice lets learners produce language more freely: https://www.teachingenglish.org.uk/professional-development/teachers/teaching-knowledge-database/c/controlled-practice and https://www.teachingenglish.org.uk/en/article/free-practice

## Recommended Approach

Use a deterministic Tutor Session Engine inside the existing local course-bank path.

This is the recommended approach because it is:

- Testable without relying on live LLM output.
- Compatible with the existing 300-lesson course bank.
- Stable enough for Telegram and web to share.
- Flexible enough to add richer lesson fields without breaking existing clients.

Rejected alternatives:

- Lightly patching the current stages: faster, but it keeps the feeling of disconnected exercises.
- Fully LLM-generated lessons: more flexible, but unstable, harder to regression test, and riskier for a paid product flow.

## Lesson Arc

Each lesson should have this tutor-session structure:

1. Goal and scene
   - Tell the learner what they will be able to say by the end.
   - Keep it tied to one real scenario, such as clinic, cafe, hotel, work, travel, shopping, or service desk.

2. Lesson words
   - Show 4-6 high-value words or chunks for the scene.
   - These are not random dictionary cards. They are "words needed to solve this scene".
   - Start with active recall: target word, translation choices, and one example.

3. Mini teaching
   - One pattern only.
   - Explain the scene order and the reusable sentence shape.
   - Example for doctor: action -> doctor -> today.
   - Model answer: `Hello. I need to see a doctor today, please.`

4. Controlled practice
   - Multiple-choice checks that teach why an answer works.
   - Questions should check:
     - the best full line;
     - the required detail;
     - the next natural line.
   - Distractors should be useful mistakes, not random vocabulary.

5. Build your answer
   - The learner writes their own answer.
   - The tutor checks for required scene slots, not just "any lesson word".
   - Feedback must name what is present and what is missing.

6. Listening
   - The learner hears a model line without seeing the full text as the main surface.
   - The task checks 2-3 key details from the model.
   - This prepares the phrase before pronunciation.

7. Pronunciation
   - The learner repeats the target phrase.
   - Existing pronunciation scoring stays, but the target phrase must match the lesson's model answer or a focused variant of it.

8. Mini dialogue
   - The tutor gives one last realistic prompt.
   - The learner answers with the target pattern.
   - Variants should include weak, acceptable, and strong replies.

9. Final word check
   - New final test before SRS scheduling.
   - Checks the same lesson words again in mixed form:
     - translation recall;
     - word-in-scene meaning;
     - short phrase completion;
     - optional spoken/written recall where available.
   - The lesson is not finished until this check is complete.

10. Tutor summary
   - Summarize what the learner can now say.
   - Name strong words, weak words, and the next review schedule.
   - Then show SRS buttons.

## Backend Design

Keep `buildTutorLessonForSequence` as the source of lesson generation. Add explicit lesson-session fields while preserving existing fields for compatibility.

Proposed public JSON additions:

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

Add to `tutorLesson`:

```go
TeachingPoint  tutorTeachingPoint   `json:"teaching_point,omitempty"`
FinalWordCheck tutorFinalWordCheck  `json:"final_word_check,omitempty"`
TutorSummary   tutorSummaryBlock    `json:"tutor_summary,omitempty"`
```

Internal helpers:

- `tutorTeachingPointFor(template, topic, function, interfaceLanguage)`
- `tutorFinalWordCheck(words, template, interfaceLanguage)`
- `tutorWordCheckItems(words, template, interfaceLanguage)`
- `tutorAnswerSlotFeedback(answer, template, interfaceLanguage)`
- `tutorSummaryForLesson(words, template, function, interfaceLanguage)`

The existing `ScenarioSlots` remains the source of truth for required action, item, detail, politeness, and model answer.

## Frontend Web Design

Update `TutorView` to make the visible flow feel like a session:

- Rename the final stage from memory review to final check/review.
- Add a real final word-check state before SRS buttons.
- Use `teaching_point` for the mini explanation when present.
- Keep existing `mini_explanation`, `checks`, `writing_expected`, `listening_expected`, and `dialogue_variants` as fallback.
- Show tutor feedback that names missing scene slots:
  - missing action;
  - missing item;
  - missing detail;
  - missing politeness if expected.
- Do not expose internal service labels.
- Keep pronunciation layout work already present in `.tutor-pronunciation-stage-v2`.

Visible lesson stages on web:

```text
Words -> Pattern -> Guided practice -> Your answer -> Listening -> Pronunciation -> Mini dialogue -> Final check -> Summary
```

The left plan can still look like step navigation, but the central transcript should read like the tutor is guiding the learner.

## Telegram Design

Telegram remains a compact lesson card, not a full interactive clone.

It should show:

- lesson title and course number;
- goal and scene;
- 4-6 lesson words;
- one teaching point with pattern and model answer;
- one guided practice question;
- one mini dialogue prompt and model answer;
- final check preview: which words will be reviewed;
- next lesson/menu buttons.

If Telegram later becomes interactive, it should reuse `final_word_check` instead of creating a separate Telegram-only flow.

## Error Handling

- If `FinalWordCheck` cannot be built, fallback to current `review` words and SRS buttons.
- If a lesson has fewer than 4 words, keep the existing "not enough vocabulary" backend error.
- If `teaching_point` is absent, web and Telegram fallback to `mini_explanation` and `pronunciation_text`.
- If pronunciation check fails, keep the learner on pronunciation and show the existing request error.
- If the learner fails the final word check, keep the stage active and show the weak words instead of finishing the lesson.

## Testing

Backend tests:

- Doctor lesson teaching point contains the exact order action -> doctor -> today and model answer `Hello. I need to see a doctor today, please.`
- Final word check includes the lesson's key words and does not use unrelated distractors as correct answers.
- Final word check requires multiple correct answers before summary/SRS.
- Tutor feedback names missing slots for writing/dialogue answers.
- Existing scenario-slot tests remain passing.
- `tutorLessonSteps` and the frontend stage model remain conceptually aligned.

Frontend Playwright tests:

- AI Tutor shows the new stage flow.
- Mini teaching uses `teaching_point` when returned by API.
- Writing stage rejects an answer missing the required detail.
- Mini dialogue accepts a correct model-pattern answer.
- Final word check appears after dialogue and before SRS.
- Lesson cannot finish before final word check passes.
- Summary shows can-say, strong words, weak words, and SRS actions.

Telegram tests:

- `formatTelegramTutorLesson` includes teaching point and model answer.
- Telegram card includes lesson words and final-check preview.
- Telegram rendering does not leak internal fields.

Full verification after implementation:

```powershell
go test ./...
npm --prefix web-react run build
npm --prefix web-react run e2e
node tools/check_encoding_artifacts.mjs
```

## Completion Criteria

- AI Tutor feels like a complete lesson: teach, practise, produce, listen, pronounce, dialogue, final vocabulary check, summary.
- The final block checks words instead of only scheduling SRS.
- Doctor lesson can produce the user-requested pattern exactly.
- Web and Telegram share the same backend lesson structure.
- Existing local course-bank behavior remains deterministic and testable.
- No visible UI text contains mojibake, service labels, or random internal criteria.
