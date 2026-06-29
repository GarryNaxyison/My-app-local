package main

import (
	"encoding/json"
	"strings"
)

func aiTutorLessonGenerationPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, topicSeed string, recentFingerprints []string) []chatMessage {
	level = normalizeCEFRLevel(level)
	vars := mergePromptVars(commonPromptVars(language, interfaceLanguage), map[string]string{
		"level":                      level,
		"level_band":                 aiTutorLevelBand(level),
		"topic_seed":                 firstNonEmpty(strings.TrimSpace(topicSeed), "everyday life"),
		"recent_lesson_fingerprints": strings.Join(recentFingerprints, ", "),
	})
	userPrompt := renderAppPrompt("ai_tutor.lesson.generate.user", aiTutorLessonGenerationPromptFallback, vars)
	return []chatMessage{
		{Role: "system", Content: "You are an expert CEFR language lesson designer. Return valid JSON only."},
		{Role: "user", Content: userPrompt},
	}
}

func aiTutorRetellCheckPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, lesson aiTutorLessonPayload, learnerAnswer string) []chatMessage {
	payload := aiTutorJSON(map[string]any{
		"lesson":         lesson,
		"learner_answer": strings.TrimSpace(learnerAnswer),
	})
	return aiTutorCheckerPrompt("ai_tutor.retell.check.user", aiTutorRetellCheckPromptFallback, language, interfaceLanguage, level, payload)
}

func aiTutorQuestionCheckPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, lesson aiTutorLessonPayload, question aiTutorQuestion, learnerAnswer string) []chatMessage {
	payload := aiTutorJSON(map[string]any{
		"story":          lesson.Story,
		"question":       question,
		"learner_answer": strings.TrimSpace(learnerAnswer),
	})
	return aiTutorCheckerPrompt("ai_tutor.question.check.user", aiTutorQuestionCheckPromptFallback, language, interfaceLanguage, level, payload)
}

func aiTutorProductionCheckPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, lesson aiTutorLessonPayload, learnerAnswer string) []chatMessage {
	payload := aiTutorJSON(map[string]any{
		"words":           lesson.Words,
		"production_task": lesson.ProductionTask,
		"learner_answer":  strings.TrimSpace(learnerAnswer),
	})
	return aiTutorCheckerPrompt("ai_tutor.production.check.user", aiTutorProductionCheckPromptFallback, language, interfaceLanguage, level, payload)
}

func aiTutorQualityPrompt(language learningLanguage, interfaceLanguage learningLanguage, lesson aiTutorLessonPayload, sessionSummary any, kind string) []chatMessage {
	kind = strings.TrimSpace(kind)
	if kind == "" {
		kind = "preflight"
	}
	key := "ai_tutor.lesson.quality_preflight.user"
	fallback := aiTutorQualityPreflightPromptFallback
	if kind == "post" || kind == "post_lesson" {
		key = "ai_tutor.lesson.quality_post.user"
		fallback = aiTutorQualityPostPromptFallback
	}
	payload := aiTutorJSON(map[string]any{
		"kind":            kind,
		"lesson":          lesson,
		"session_summary": sessionSummary,
	})
	return aiTutorCheckerPrompt(key, fallback, language, interfaceLanguage, lesson.Level, payload)
}

func aiTutorCheckerPrompt(key string, fallback string, language learningLanguage, interfaceLanguage learningLanguage, level string, payloadJSON string) []chatMessage {
	vars := mergePromptVars(commonPromptVars(language, interfaceLanguage), map[string]string{
		"level":        normalizeCEFRLevel(level),
		"payload_json": payloadJSON,
	})
	userPrompt := renderAppPrompt(key, fallback, vars)
	return []chatMessage{
		{Role: "system", Content: "You are a strict language-learning evaluator. Return valid JSON only."},
		{Role: "user", Content: userPrompt},
	}
}

func aiTutorJSON(value any) string {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return "{}"
	}
	return string(data)
}

const aiTutorLessonGenerationPromptFallback = `Interface language: {{interface_language_native_name}}
Target learning language: {{learning_language_native_name}}
Exact learner level: {{level}}
Level band: {{level_band}}
Topic/theme seed: {{topic_seed}}
Recent lesson fingerprints to avoid: {{recent_lesson_fingerprints}}

Return strict JSON only. No Markdown. No prose outside JSON.

Create one complete reusable micro-lesson for a language learner.
All learner-facing fields that are not explicitly suffixed with _target must be written in the interface language, including title, theme, lesson_goal, grammar explanations, task instructions, evaluation criteria, recommendations, feedback rules, review text, and quality self-check text.

Requirements:
1. Create a short story in the target learning language.
   - Exactly 5-7 sentences.
   - Natural, everyday, culturally neutral.
   - Match the exact CEFR level strictly.
   - Do not include translations inside the story.
   - Create a concrete short story title in the interface language and put it in story.title_interface.
   - Also put an English-safe machine title in story.story_title.
   - Do not use generic titles such as Story, Continue, Next, Section, Раздел, or Дальше.
2. Select exactly 6 useful target words or short chunks from the story.
   - Each word must include target form and interface-language translation.
   - Each word must include audio_text_target for target-language word/chunk audio.
   - Each example sentence must include example_audio_text_target for target-language audio.
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
   - Include final recommendations in the interface language.
7. Create final reflection and review options.

Return exactly this JSON shape:
{
  "title": "",
  "level": "",
  "level_band": "",
  "target_language": "",
  "interface_language": "",
  "theme": "",
  "lesson_goal": "",
  "grammar_focus": {"name": "", "short_explanation_interface": "", "model_sentence_target": ""},
  "story": {"title_interface": "", "story_title": "", "text_target": "", "audio_text_target": "", "sentence_count": 0},
  "words": [{"id": "", "target": "", "interface_translation": "", "part_of_speech": "", "example_sentence_target": "", "example_translation_interface": "", "audio_text_target": "", "example_audio_text_target": "", "difficulty_note_interface": ""}],
  "retell_task": {"instruction_interface": "", "min_sentences": 2, "feedback_rubric": []},
  "comprehension_questions": [{"id": "", "question_target": "", "expected_points": [], "feedback_rule_interface": ""}],
  "word_learning": [{"word_id": ""}],
  "word_recall": [{"word_id": ""}],
  "production_task": {"instruction_interface": "", "required_word_count": 3, "sentence_count": "2-3", "evaluation_criteria": [], "recommendations_interface": []},
  "review_options": ["tomorrow", "3_days", "1_week", "no_review"],
  "quality_self_check": {"cefr_reason": "", "why_reusable": "", "duplicate_risk": ""}
}`

const aiTutorRetellCheckPromptFallback = `Target learning language: {{learning_language_native_name}}
Interface language: {{interface_language_native_name}}
Level: {{level}}

Evaluate the learner retell against the lesson story.
Return strict JSON only with keys: ok, comprehension_score, corrected_answer_target, advice_interface, mistakes.
Advice and mistake explanations must use the interface language. Corrected answer must use the target language.

Payload:
{{payload_json}}`

const aiTutorQuestionCheckPromptFallback = `Target learning language: {{learning_language_native_name}}
Interface language: {{interface_language_native_name}}
Level: {{level}}

Evaluate one comprehension answer.
Return strict JSON only with keys: ok, result, corrected_answer_target, tip_interface, mistakes.
result must be one of correct, partial, incorrect.

Payload:
{{payload_json}}`

const aiTutorProductionCheckPromptFallback = `Target learning language: {{learning_language_native_name}}
Interface language: {{interface_language_native_name}}
Level: {{level}}

Evaluate the learner's 2-3 sentence production task.
Return strict JSON only with keys: ok, used_words, missing_requirement, corrected_version_target, recommendations_interface, mistakes.

Payload:
{{payload_json}}`

const aiTutorQualityPreflightPromptFallback = `Target learning language: {{learning_language_native_name}}
Interface language: {{interface_language_native_name}}
Level: {{level}}

Review this generated lesson before any learner sees it.
Return strict JSON only with keys: approved, score, critical_issues, fix_suggestions, reasons.
Reject wrong language, wrong CEFR level, unsafe content, duplicates, wrong counts, bad translations, or questions not based on the story.

Payload:
{{payload_json}}`

const aiTutorQualityPostPromptFallback = `Target learning language: {{learning_language_native_name}}
Interface language: {{interface_language_native_name}}
Level: {{level}}

Review this completed trial lesson for shared-bank promotion.
Return strict JSON only with keys: approved, score, reasons, bank_promotion_decision, critical_issues.
Approve only strong reusable lessons with score at least 88 and no critical issues.

Payload:
{{payload_json}}`
