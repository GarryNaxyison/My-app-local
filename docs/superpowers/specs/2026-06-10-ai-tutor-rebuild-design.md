# AI Tutor Rebuild Design

## Status

Approved for design on 2026-06-10.

## Goal

Rebuild AI Tutor from scratch around AI-generated reusable lessons. The old local tutor course bank and old stored tutor lesson base are no longer the product source.

The new tutor must provide the same complete interactive lesson in both Web and Telegram:

1. AI prepares a short 5-7 sentence level-appropriate story in the target learning language, including target-language audio text for playback.
2. The learner studies the story.
3. The learner retells what they understood in their own words.
4. The learner answers three comprehension questions in the target learning language.
5. The learner studies six AI-selected lesson words with interface-language translations, examples, and audio playback for the word and example.
6. The learner recalls the six words through multiple-choice test options.
7. The learner writes two or three target-language sentences using the lesson words, then receives final recommendations from AI.
8. The learner rates the lesson and chooses whether to review it tomorrow, in three days, in one week, or not at all.
9. Only high-quality lessons are promoted into the shared lesson bank.

## Product Decisions

- Telegram must run the full interactive lesson, not a compact card.
- Web and Telegram must share one backend lesson/session engine.
- Old `tutor_lessons` and `tutor_user_lessons` data will be cleared or replaced by migration.
- The deterministic local course-bank generator in `course_tutor.go` is not the new source of lessons.
- New lessons are AI-generated in strict JSON and stored with lifecycle status.
- Generated lesson JSON must include audio-ready target text for the story, every word, and every word example.
- Word recall is a test with answer variants, not a free-form old repeat card.
- The production checker must surface final recommendations in the interface language.
- Lesson generation must use the user's configured topic/theme preference as the prompt topic seed when it is present.
- Lessons promoted after quality review are stored in a separate reusable SQLite lesson bank, not only in the main user-progress database.
- Lessons are grouped for reuse by learning language, interface language, and level band: `A1-A2`, `B1-B2`, `C1-C2`.
- The exact user level still matters during generation and validation.
- A lesson is reusable only after preflight validation and post-lesson quality review.

## Recommended Approach

Use a shared AI Tutor Engine with thin Web and Telegram adapters.

Rejected alternatives:

- Separate Web and Telegram tutors: faster initially, but duplicates stage logic, checking rules, and lesson state.
- Telegram-first tutor engine: makes the bot easier, but forces Web into Telegram-shaped UX.
- Single prompt with no quality gate: simpler, but unsafe for a shared lesson bank because low-quality AI output can be reused.

The shared engine is recommended because it makes lesson state, answer checking, quality scoring, review scheduling, and future analytics consistent across both surfaces.

## Architecture

The new system has five backend units:

1. Lesson Generator
   - Calls AI with a strict JSON prompt.
   - Produces one complete lesson package.
   - Uses learning language, interface language, exact CEFR level, level band, topic seed, and recent lesson fingerprints.

2. Lesson Validator
   - Runs deterministic local checks first.
   - Runs AI quality preflight second.
   - Rejects malformed, off-level, wrong-language, unsafe, duplicate, or incomplete lessons.

3. Session Engine
   - Owns the active lesson stage.
   - Accepts learner input.
   - Calls the correct AI checker for free-text stages.
   - Advances, retries, or finishes the session.

4. Surface Adapters
   - Web adapter exposes JSON API endpoints.
   - Telegram adapter renders the same `next_step` and sends text/callback input back to the engine.

5. Quality Promoter
   - Runs after lesson completion.
   - Combines AI quality review with learner completion data.
   - Promotes only top lessons into the approved shared bank.

## Lesson Lifecycle

```text
missing approved lesson
-> generate draft lesson
-> local validation
-> AI preflight validation
-> create user session
-> run full interactive lesson
-> collect stage answers and metrics
-> AI post-lesson quality check
-> approved for reuse or rejected
```

Lesson statuses:

- `draft`: generated but not shared.
- `preflight_rejected`: failed local or AI validation.
- `in_trial`: being used by a first learner but not yet reusable.
- `approved`: may be reused by future users in the same language/interface/level band.
- `rejected`: completed but not good enough for reuse.
- `archived`: removed from rotation after later quality or product concerns.

## Lesson Stages

Canonical stage order:

1. `story_intro`
   - Shows title, level, theme, goal, grammar focus, and story.
   - Story is 5-7 sentences in the target learning language.

2. `retell`
   - User describes what they understood.
   - Answer language should be the target learning language.
   - AI checks comprehension and gives concise corrections in the interface language.

3. `question_1`
4. `question_2`
5. `question_3`
   - Each question is in the target learning language.
   - User answers in the target learning language.
   - AI checks each answer separately, corrects mistakes, and gives one useful tip.

6. `word_learn_1` through `word_learn_6`
   - Each target word or chunk is shown with interface-language translation, example, short note, word audio, and example audio.
   - This is the "learn words" stage.

7. `word_recall_1` through `word_recall_6`
   - The learner recalls the words through multiple-choice answer variants with one correct option.
   - This is the "repeat/check words" stage.

8. `production`
   - The learner writes two or three target-language sentences using at least three of the six words.
   - AI checks usage, grammar, naturalness, and returns final recommendations in the interface language.

9. `lesson_feedback`
   - Asks how the lesson felt.
   - Stores rating, difficulty, weak words, and learner note if present.

10. `review_schedule`
   - User chooses tomorrow, in three days, in one week, or no review.

11. `complete`
   - Stores completion metrics and triggers post-lesson quality review.

## Database Design

Add new tables instead of stretching the old tutor tables.

```sql
CREATE TABLE IF NOT EXISTS ai_tutor_lessons (
  id TEXT PRIMARY KEY,
  learning_language TEXT NOT NULL,
  interface_language TEXT NOT NULL,
  exact_level TEXT NOT NULL,
  level_band TEXT NOT NULL,
  theme TEXT NOT NULL DEFAULT '',
  status TEXT NOT NULL,
  payload_json TEXT NOT NULL,
  fingerprint TEXT NOT NULL DEFAULT '',
  preflight_score INTEGER NOT NULL DEFAULT 0,
  post_score INTEGER NOT NULL DEFAULT 0,
  completion_count INTEGER NOT NULL DEFAULT 0,
  average_rating REAL NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS ai_tutor_sessions (
  id TEXT PRIMARY KEY,
  telegram_id INTEGER NOT NULL,
  lesson_id TEXT NOT NULL,
  surface TEXT NOT NULL,
  current_stage TEXT NOT NULL,
  status TEXT NOT NULL,
  started_at TEXT NOT NULL,
  completed_at TEXT NOT NULL DEFAULT '',
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS ai_tutor_answers (
  session_id TEXT NOT NULL,
  stage TEXT NOT NULL,
  answer_text TEXT NOT NULL DEFAULT '',
  answer_json TEXT NOT NULL DEFAULT '',
  feedback_json TEXT NOT NULL DEFAULT '',
  correct INTEGER NOT NULL DEFAULT 0,
  created_at TEXT NOT NULL,
  PRIMARY KEY (session_id, stage)
);

CREATE TABLE IF NOT EXISTS ai_tutor_reviews (
  id TEXT PRIMARY KEY,
  telegram_id INTEGER NOT NULL,
  lesson_id TEXT NOT NULL,
  due_at TEXT NOT NULL,
  interval_code TEXT NOT NULL,
  status TEXT NOT NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS ai_tutor_quality_checks (
  id TEXT PRIMARY KEY,
  lesson_id TEXT NOT NULL,
  session_id TEXT NOT NULL DEFAULT '',
  kind TEXT NOT NULL,
  score INTEGER NOT NULL,
  approved INTEGER NOT NULL,
  issues_json TEXT NOT NULL DEFAULT '[]',
  result_json TEXT NOT NULL,
  created_at TEXT NOT NULL
);
```

Migration rules:

- Clear or drop old `tutor_lessons` and `tutor_user_lessons`.
- Keep user profile, level, language, XP, mistakes, learned words, and reminders.
- Remove old local course-bank source from the AI Tutor start path.
- Keep old `lesson` micro-task feature separate unless explicitly replaced later.

## Lesson JSON Shape

The generator must return strict JSON only.

```json
{
  "title": "",
  "level": "",
  "level_band": "",
  "target_language": "",
  "interface_language": "",
  "theme": "",
  "lesson_goal": "",
  "grammar_focus": {
    "name": "",
    "short_explanation_interface": "",
    "model_sentence_target": ""
  },
  "story": {
    "text_target": "",
    "sentence_count": 0
  },
  "words": [
    {
      "id": "",
      "target": "",
      "interface_translation": "",
      "part_of_speech": "",
      "example_sentence_target": "",
      "example_translation_interface": "",
      "difficulty_note_interface": ""
    }
  ],
  "retell_task": {
    "instruction_interface": "",
    "min_sentences": 2,
    "feedback_rubric": []
  },
  "comprehension_questions": [
    {
      "id": "",
      "question_target": "",
      "expected_points": [],
      "feedback_rule_interface": ""
    }
  ],
  "word_learning": [],
  "word_recall": [],
  "production_task": {
    "instruction_interface": "",
    "required_word_count": 3,
    "sentence_count": "2-3",
    "evaluation_criteria": []
  },
  "review_options": ["tomorrow", "3_days", "1_week", "no_review"],
  "quality_self_check": {
    "cefr_reason": "",
    "why_reusable": "",
    "duplicate_risk": ""
  }
}
```

Hard requirements:

- `words` must contain exactly six items.
- `comprehension_questions` must contain exactly three items.
- Story must contain five to seven sentences.
- User-answer stages must require target-language answers.
- Explanations, instructions, and feedback must use the interface language except where target-language output is explicitly required.

## Generator Prompt

The runtime prompt should be stored in `app_prompts.json` under a new key such as `ai_tutor.lesson.generate.user`.

Prompt intent:

```text
You are an expert CEFR language lesson designer.

Create one complete micro-lesson for a language learner.

Inputs:
- Interface language: {{interface_language_native_name}}
- Target learning language: {{learning_language_native_name}}
- Exact learner level: {{level}}
- Level bank: {{level_band}}
- Topic/theme seed: {{topic_seed}}
- Recent lesson fingerprints to avoid: {{recent_lesson_fingerprints}}

Return strict JSON only. No Markdown. No prose outside JSON.

The lesson must be level-appropriate, practical, short, and reusable for other learners with the same target language and level bank.

Requirements:
1. Create a short story in the target learning language.
   - Exactly 5-7 sentences.
   - Natural, everyday, culturally neutral.
   - Match the exact CEFR level strictly.
   - Do not include translations inside the story.
2. Select exactly 6 useful target words or short chunks from the story.
   - Each word must include target form and interface-language translation.
   - Prefer high-frequency useful vocabulary.
3. Create a learner retelling task.
4. Create exactly 3 comprehension questions.
   - Questions are in the target learning language.
   - Answers are expected in the target learning language.
   - Include expected answer points and correction guidance.
5. Create word learning and recall data for the 6 words.
6. Create a production task.
   - The learner writes 2-3 target-language sentences.
   - The learner must use at least 3 of the 6 words.
7. Create final reflection and review options.

Return the exact JSON shape requested by the application.
```

## Checker Prompts

Use separate prompts for separate checks.

1. `ai_tutor.retell.check.user`
   - Input: story, level, learner retell.
   - Output: JSON with comprehension score, corrected target-language answer, concise interface-language advice, mistakes.

2. `ai_tutor.question.check.user`
   - Input: one question, expected points, learner answer.
   - Output: JSON with correct/partial/incorrect, corrected answer, one tip, mistakes.

3. `ai_tutor.production.check.user`
   - Input: six words, production task, learner sentences.
   - Output: JSON with used words, missing requirement, corrected version, interface-language recommendations, mistakes.

4. `ai_tutor.lesson.quality_preflight.user`
   - Input: generated lesson JSON.
   - Output: approved, score 0-100, critical issues, fix suggestions.

5. `ai_tutor.lesson.quality_post.user`
   - Input: generated lesson JSON, stage results, learner rating, completion data.
   - Output: approved, score 0-100, reasons, bank promotion decision.

Promotion threshold:

- A lesson is promoted only when `approved = true`, score is at least 88, and there are no critical issues.

Critical preflight failures:

- Wrong target language.
- Wrong interface language for instructions.
- Wrong CEFR difficulty.
- Story is not five to seven sentences.
- Not exactly six words.
- Not exactly three questions.
- Questions are not based on the story.
- Translations are wrong or misleading.
- Unsafe, political, sexual, hateful, medical-advice-risk, or financial-advice-risk content.
- Duplicate or near-duplicate of a recent lesson fingerprint.

## Telegram Design

Telegram must be a full interactive client.

State:

- Store active mode as `ai_tutor:<session_id>`.
- The canonical current stage lives in `ai_tutor_sessions.current_stage`.
- Telegram callbacks send stage actions, word choices, and review choices.
- Plain messages are submitted as free-text answers for the active stage.

Rendering:

- Each stage sends one clear message.
- Long story messages can be split safely by paragraph.
- Audio controls are shown for story playback and for each word/example pair.
- Buttons are used for word learning, recall choices, continue, retry, and review schedule.
- Free-text stages show a clear prompt and then wait for the user's next message.

Telegram input handling:

- If mode starts with `ai_tutor:`, route text to `SubmitAITutorAnswer`.
- If callback starts with `ait|`, route callback payload to the same session engine.
- After each submission, render returned feedback and next step.

## Web Design

Web stays richer visually but does not own separate learning logic.

Endpoints:

- `POST /api/ai-tutor/start`
- `GET /api/ai-tutor/session`
- `POST /api/ai-tutor/answer`
- `POST /api/ai-tutor/review`
- `POST /api/ai-tutor/finish`

The response should always include:

```json
{
  "session": {},
  "lesson": {},
  "current_stage": "",
  "next_step": {},
  "feedback": {},
  "user": {}
}
```

The existing React Tutor view can be replaced or heavily simplified around server-driven `next_step`. LocalStorage should not be the source of truth for lesson progress.

## Error Handling

- If AI generation fails, retry with the same input up to a small limit, then show a localized error.
- If generated JSON is invalid, reject the lesson and retry generation.
- If preflight fails, store the quality check and retry generation with failure reasons.
- If a free-text checker fails, keep the learner on the same stage and show a retry-safe message.
- If Telegram receives text with no active session, offer to start a new AI Tutor lesson.
- If a session is interrupted, resume from `ai_tutor_sessions.current_stage`.
- If an approved lesson later shows weak metrics, archive it.

## Testing

Backend tests:

- Starts an AI Tutor session for Web and Telegram from the same engine.
- Rejects generated lessons with invalid JSON.
- Rejects generated lessons with fewer or more than six words.
- Rejects generated lessons with fewer or more than three questions.
- Rejects generated lessons whose story is outside five to seven sentences.
- Advances through all canonical stages in order.
- Stores free-text answers and checker feedback.
- Schedules review for tomorrow, three days, one week, or no review.
- Does not reuse rejected lessons.
- Promotes only approved lessons with score at least 88.
- Clears old tutor lesson tables during migration.

Telegram tests:

- `menu_tutor` starts a full interactive session.
- Telegram text routes to the active AI Tutor session.
- Callback choices advance word and recall stages.
- Review callback completes the session.
- Interrupted session resumes from stored stage.

Web tests:

- Web starts a session through `/api/ai-tutor/start`.
- Web submits each stage through `/api/ai-tutor/answer`.
- Web does not rely on LocalStorage as canonical progress.
- Final review scheduling updates backend state.

Verification after implementation:

```powershell
go test ./...
npm --prefix web-react run build
npm --prefix web-react run e2e
node tools/check_encoding_artifacts.mjs
```

## Completion Criteria

- AI Tutor no longer uses the old local course bank as the product lesson source.
- Old stored tutor lesson base is cleared or replaced by migration.
- Web and Telegram both run the same full interactive AI Tutor lesson.
- Lesson generation uses strict JSON.
- Lesson validation rejects malformed, unsafe, and off-spec AI output before user exposure.
- Post-lesson quality review promotes only strong lessons into the shared bank.
- User review scheduling supports tomorrow, three days, one week, and no review.
- Progress is resumable after interruption.
- Tests cover generation validation, stage progression, Telegram flow, Web flow, and promotion logic.
