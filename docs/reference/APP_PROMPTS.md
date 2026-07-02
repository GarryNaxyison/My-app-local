# NERIVA Prompt Registry

Purpose: this is the root project index for prompts that are editable outside the Go binary. Runtime templates live in `app_prompts.json` and are uploaded together with the binary and environment files.

Editing rule: every prompt must keep a clear header with the feature, source function, model/tool path, variables, and a short note explaining where the prompt is used.

## Target Deployable File

- Runtime file: `app_prompts.json`
- Server location: next to the deployed binary and `.env` files
- Runtime behavior: `APP_PROMPTS_FILE` defaults to `app_prompts.json`; when a key is missing, the app falls back to the compiled Go prompt.
- Required navigation: prompt headers must make it possible to find a prompt by feature name, not only by code function.

## Prompt Index

### coach.system

- Feature: shared AI coach identity
- Source: `prompts.go` -> `coachSystemPrompt`
- Used by: lesson generation, lesson feedback, practice chat, plain translation, Listening, pronunciation-aware flows
- Variables: `learning_language.native_name`, `learning_language.teacher_noun`, `interface_language.native_name`
- Note: this controls tone, explanation language, Telegram/plain-text formatting, and anti-overload behavior.

### lesson.generate

- Feature: new lesson / active-recall micro-task
- Source: `prompts.go` -> `lessonPrompt`
- Used by: Telegram lesson start, Web `/api/lesson/start`
- Variables: `learning_language.native_name`, `interface_language.native_name`, `level`, localized study labels, `learning_focus_instruction`, `lesson_topic_instruction`, `lesson_number`, `recent_lesson_context`, `variation_seed`, `variation_instruction`
- Note: generates one short real-life task with situation, pattern notice, chunks, and learner instruction. It also receives the optional Settings `learning_focus` instruction plus a CEFR topic rotation and anti-repeat context for the last 10 lessons.

### lesson.feedback

- Feature: lesson answer correction
- Source: `prompts.go` -> `feedbackPrompt`
- Used by: Telegram lesson answer, Web `/api/lesson/answer`
- Variables: `learning_language.native_name`, `interface_language.native_name`, `level`, `task`, `answer`
- Note: returns compact feedback plus `---MISTAKES---` JSON for the mistake dictionary.

### practice.chat

- Feature: adaptive practice chat
- Source: `prompts.go` -> `practicePrompt`
- Used by: Telegram practice, Web `/api/practice`
- Variables: `learning_language.native_name`, `interface_language.native_name`, `level`, recent learner messages, localized study labels, `learning_focus_instruction`, `practice_topic_instruction`, `practice_turn_number`, `variation_seed`, `variation_instruction`
- Note: must keep one model phrase and one final target-language question, then return mistakes JSON. It receives the optional Settings `learning_focus`, the last 5 practice messages, and a practice topic rotation so vague/generic exchanges do not fall back to the same scene.

### roleplay.scenario

- Feature: AI roleplay scenario
- Source: `prompts.go` -> `roleplayPrompt`; structured payload is prepared by `web-react/src/App.tsx` -> `buildRoleplayScenarioPrompt`
- Used by: Web `/app/v2` Roleplay through `/api/practice` payloads beginning with `ROLEPLAY_TOOL_V2`
- Variables: `learning_language.native_name`, `interface_language.native_name`, `roleplay_payload`, `variation_seed`, `variation_instruction`
- Note: this is intentionally separate from `practice.chat`; it keeps roleplay scene/setup/result structure, respects the optional Settings `learning_focus`, and still appends the `---MISTAKES---` JSON block.

### practice.image_context

- Feature: practice with attached image
- Source: `prompts.go` -> `practiceImageMessage`
- Used by: Telegram/web image-in-practice paths before `practicePrompt`
- Variables: `learner_text`, `image_context`
- Note: tells the coach not to grade image OCR/vision context as learner writing.

### tools.translation_plain

- Feature: quick translation from voice transcript or plain text
- Source: `prompts.go` -> `translationPrompt`
- Used by: Telegram voice/text translation helper, Web tool fallback paths
- Variables: `interface_language.native_name`, `text`
- Note: returns only the natural translation, without explanations.

### tools.translation_pair

- Feature: translator tool with source and target languages
- Source: `translator.go` -> `translationToolPrompt`
- Used by: Telegram/web translator tool
- Variables: `source_language`, `target_language`, `interface_language`, `text`
- Note: this is the main prompt for the explicit source-to-target translator.

### vocabulary.hint

- Feature: hidden-word clue
- Source: `prompts.go` -> `vocabularyHintPrompt`
- Used by: word learning/review/spelling preparation
- Variables: `word`, `interface_prompt`, `learning_language.native_name`, `interface_language.native_name`, `mode`
- Note: must not reveal the answer.

### vocabulary.example

- Feature: CEFR-matched example sentence
- Source: `prompts.go` -> `vocabularyExamplePrompt`
- Used by: vocabulary cards, word learning, spelling/review examples
- Variables: `word`, `learning_language.native_name`, `interface_language.native_name`, `level`, `mode`
- Note: must return one short sentence matching the word level.

### listening.phrase

- Feature: Listening target phrase
- Source: `shadowing.go` -> `shadowingPhrasePrompt`
- Used by: Telegram/web Listening start
- Variables: `learning_language.native_name`, `interface_language.native_name`, `level`, `deck_card`
- Note: generates the phrase the learner hears and repeats.

### listening.feedback

- Feature: Listening repeat feedback
- Source: `shadowing.go` -> `shadowingFeedbackPrompt`
- Used by: Telegram/web Listening answer
- Variables: `target`, `transcript`, `score`, `missing`, `from_voice`, `level`
- Note: explains pronunciation/listening result briefly.

### pronunciation.coach_report

- Feature: pronunciation score explanation
- Source: `pronunciation.go` -> `pronunciationCoachPrompt`
- Used by: pronunciation assessment after STT/technical scoring
- Variables: technical report with transcript, expected text, scores, weak words, fluency
- Note: turns technical scoring into learner-friendly feedback.

### level.assessment_localization

- Feature: level-test prompt localization
- Source: `level_assessment.go` -> `localizedAssessmentPrompt`, `assessmentPromptPayload`
- Used by: CEFR assessment questions
- Variables: `interface_language`, `learning_language`, question payload
- Note: this is assessment copy/payload shaping rather than a free coaching prompt.
