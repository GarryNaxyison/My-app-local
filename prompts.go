package main

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
)

func coachSystemPrompt(language learningLanguage, interfaceLanguage learningLanguage) string {
	fallback := `You are an encouraging ` + language.NativeName + ` coach for learners who use ` + interfaceLanguage.NativeName + ` as their interface and explanation language.
Act like an expert ` + language.TeacherNoun + `, not a casual chatbot.
The learner's interface language is ` + interfaceLanguage.NativeName + `. Write explanations, instructions, feedback, and non-target-language UI text in ` + interfaceLanguage.NativeName + `.
Base explanations on standard language teaching practice: CEFR levels, core grammar patterns, common learner mistakes, pronunciation awareness, collocations, and practical classroom-style examples.
When you teach vocabulary, prefer high-frequency everyday ` + language.NativeName + ` and stable dictionary meanings. Do not invent rare meanings unless the learner asks.
When checking answers, explain the rule briefly and give one natural version a teacher would accept.
Keep replies practical, warm, and concise.
The learner is a beginner. Explain mostly in ` + interfaceLanguage.NativeName + `.
Use ` + language.NativeName + ` only for short examples, corrected sentences, and simple practice questions.
Do not start replies with long ` + language.NativeName + ` paragraphs.
Do not overwhelm the learner. Focus on one or two useful corrections.
Format the human-readable reply beautifully in plain Telegram text: short blocks, line breaks, and simple bullet symbols like "•".
You may use 0-2 friendly relevant emojis, but do not overuse them.
Never use Markdown syntax in the human-readable reply. Do not use asterisks, **bold**, headings with #, tables, or code fences.
Never claim to be a certified teacher or therapist.`
	return renderAppPrompt("coach.system", fallback, commonPromptVars(language, interfaceLanguage))
}

type studyPromptLabels struct {
	Situation   string
	Pattern     string
	Chunks      string
	Task        string
	Correction  string
	ModelPhrase string
	YourTurn    string
}

var studyPromptLabelsByLanguage = map[string]studyPromptLabels{
	"en": {Situation: "Situation", Pattern: "Pattern", Chunks: "Useful chunks", Task: "Your task", Correction: "Correction", ModelPhrase: "Model phrase", YourTurn: "Your turn"},
	"ru": {Situation: "Ситуация", Pattern: "Шаблон", Chunks: "Полезные фразы", Task: "Твоя задача", Correction: "Исправление", ModelPhrase: "Пример фразы", YourTurn: "Твоя очередь"},
	"es": {Situation: "Situación", Pattern: "Patrón", Chunks: "Frases útiles", Task: "Tu tarea", Correction: "Corrección", ModelPhrase: "Frase modelo", YourTurn: "Tu turno"},
	"de": {Situation: "Situation", Pattern: "Muster", Chunks: "Nützliche Wendungen", Task: "Deine Aufgabe", Correction: "Korrektur", ModelPhrase: "Beispielsatz", YourTurn: "Du bist dran"},
	"fr": {Situation: "Situation", Pattern: "Structure", Chunks: "Expressions utiles", Task: "Ta tâche", Correction: "Correction", ModelPhrase: "Phrase modèle", YourTurn: "À toi"},
	"it": {Situation: "Situazione", Pattern: "Modello", Chunks: "Frasi utili", Task: "Il tuo compito", Correction: "Correzione", ModelPhrase: "Frase modello", YourTurn: "Tocca a te"},
	"zh": {Situation: "情境", Pattern: "句型", Chunks: "实用表达", Task: "你的任务", Correction: "纠正", ModelPhrase: "示例句", YourTurn: "轮到你"},
	"ja": {Situation: "場面", Pattern: "パターン", Chunks: "便利な表現", Task: "あなたの課題", Correction: "添削", ModelPhrase: "お手本の文", YourTurn: "あなたの番"},
	"ko": {Situation: "상황", Pattern: "패턴", Chunks: "유용한 표현", Task: "과제", Correction: "교정", ModelPhrase: "예시 문장", YourTurn: "당신 차례"},
	"tg": {Situation: "Вазъият", Pattern: "Намуна", Chunks: "Ибораҳои муфид", Task: "Вазифаи ту", Correction: "Ислоҳ", ModelPhrase: "Ҷумлаи намунавӣ", YourTurn: "Навбати ту"},
	"uz": {Situation: "Vaziyat", Pattern: "Andoza", Chunks: "Foydali iboralar", Task: "Vazifang", Correction: "Tuzatish", ModelPhrase: "Namuna jumla", YourTurn: "Navbat sizda"},
	"tt": {Situation: "Вазгыять", Pattern: "Үрнәк", Chunks: "Файдалы гыйбарәләр", Task: "Синең бирем", Correction: "Төзәтү", ModelPhrase: "Үрнәк җөмлә", YourTurn: "Синең чират"},
	"hy": {Situation: "Իրավիճակ", Pattern: "Կաղապար", Chunks: "Օգտակար արտահայտություններ", Task: "Քո առաջադրանքը", Correction: "Ուղղում", ModelPhrase: "Օրինակ նախադասություն", YourTurn: "Քո հերթն է"},
	"kk": {Situation: "Жағдай", Pattern: "Үлгі", Chunks: "Пайдалы тіркестер", Task: "Сенің тапсырмаң", Correction: "Түзету", ModelPhrase: "Үлгі сөйлем", YourTurn: "Сенің кезегің"},
	"ky": {Situation: "Кырдаал", Pattern: "Үлгү", Chunks: "Пайдалуу фразалар", Task: "Сенин тапшырмаң", Correction: "Түзөтүү", ModelPhrase: "Үлгү сүйлөм", YourTurn: "Сенин кезегиң"},
	"ka": {Situation: "სიტუაცია", Pattern: "ნიმუში", Chunks: "სასარგებლო ფრაზები", Task: "შენი დავალება", Correction: "შესწორება", ModelPhrase: "ნიმუში ფრაზა", YourTurn: "შენი ჯერია"},
	"uk": {Situation: "Ситуація", Pattern: "Шаблон", Chunks: "Корисні фрази", Task: "Твоє завдання", Correction: "Виправлення", ModelPhrase: "Зразкова фраза", YourTurn: "Твоя черга"},
	"pl": {Situation: "Sytuacja", Pattern: "Wzorzec", Chunks: "Przydatne zwroty", Task: "Twoje zadanie", Correction: "Korekta", ModelPhrase: "Przykładowe zdanie", YourTurn: "Twoja kolej"},
	"ro": {Situation: "Situație", Pattern: "Tipar", Chunks: "Expresii utile", Task: "Sarcina ta", Correction: "Corectare", ModelPhrase: "Frază model", YourTurn: "Rândul tău"},
	"pt": {Situation: "Situação", Pattern: "Padrão", Chunks: "Frases úteis", Task: "Sua tarefa", Correction: "Correção", ModelPhrase: "Frase modelo", YourTurn: "Sua vez"},
}

func localizedStudyPromptLabels(interfaceLanguage learningLanguage) studyPromptLabels {
	labels, ok := studyPromptLabelsByLanguage[normalizeInterfaceLanguage(interfaceLanguage.Code)]
	if !ok {
		return studyPromptLabelsByLanguage["en"]
	}
	return labels
}

func lessonPromptLabelList(labels studyPromptLabels) string {
	return labels.Situation + ":, " + labels.Pattern + ":, " + labels.Chunks + ":, " + labels.ModelPhrase + ":, " + labels.Task + ":"
}

func practicePromptLabelList(labels studyPromptLabels) string {
	return labels.Correction + ":, " + labels.ModelPhrase + ":, " + labels.YourTurn + ":"
}

func learningFocusInstruction(focus string) string {
	focus = strings.TrimSpace(focus)
	if focus == "" {
		return ""
	}
	return "User-selected global learning focus: " + focus + ". Use this as a topic preference when it fits naturally. Do not repeat the same airport/hotel scenario unless the focus asks for it. "
}

func cefrTopicPool(level string) string {
	switch normalizeCEFRLevel(level) {
	case "A1":
		return "greetings, names, countries, family, daily routine, time, food, cafe order, shopping basics, home, weather, simple directions, likes and dislikes"
	case "A2":
		return "plans, appointments, transport, hotel check-in, everyday problems, health, invitations, past weekend, work small talk, polite requests, simple opinions, comparing choices"
	case "B1":
		return "work updates, travel disruption, explaining reasons, advice, habits, experiences, plans, short storytelling, service complaints, preferences, problem solving, agreeing and disagreeing"
	case "B2":
		return "negotiation, interviews, presentations, detailed opinions, trade-offs, hypothetical situations, conflict resolution, culture, work decisions, priorities, persuasion, nuanced recommendations"
	default:
		return "abstract discussion, argument structure, nuance, professional scenarios, storytelling, debate, persuasion, idioms, register, precise corrections, complex problem solving"
	}
}

func compactPromptSnippet(value string, limit int) string {
	value = strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	if value == "" || limit <= 0 {
		return ""
	}
	if len(value) > limit {
		return value[:limit] + "..."
	}
	return value
}

func promptVariationSeed(scope string, level string, count int, focus string, recent []string) string {
	scope = strings.ToLower(strings.TrimSpace(scope))
	if scope == "" {
		scope = "ai"
	}
	if count < 0 {
		count = 0
	}
	turn := count + 1
	level = normalizeCEFRLevel(level)
	recentParts := make([]string, 0, len(recent))
	for _, item := range recent {
		if snippet := compactPromptSnippet(item, 220); snippet != "" {
			recentParts = append(recentParts, snippet)
		}
	}
	source := strings.Join([]string{
		"scope=" + scope,
		"level=" + level,
		"turn=" + strconv.Itoa(turn),
		"focus=" + compactPromptSnippet(focus, 160),
		"recent=" + strings.Join(recentParts, " | "),
	}, "\n")
	sum := sha256.Sum256([]byte(source))
	return scope + "-" + level + "-" + strconv.Itoa(turn) + "-" + hex.EncodeToString(sum[:])[:10]
}

func promptVariationInstruction(seed string) string {
	seed = strings.TrimSpace(seed)
	if seed == "" {
		return ""
	}
	return "Variation seed: " + seed + ". Use it privately to vary the topic, scene, names, objects, grammar pattern, examples, and dialogue turns. Do not reveal the variation seed to the learner. "
}

func lessonHistoryForPrompt(user userState) []string {
	if len(user.LessonHistory) > 0 {
		return trimLessonHistory(user.LessonHistory)
	}
	if strings.TrimSpace(user.LastLessonPrompt) == "" {
		return nil
	}
	return []string{user.LastLessonPrompt}
}

func recentLessonContext(recentLessons []string) string {
	recentLessons = trimLessonHistory(recentLessons)
	if len(recentLessons) == 0 {
		return ""
	}
	lines := make([]string, 0, len(recentLessons))
	for index, lesson := range recentLessons {
		if snippet := compactPromptSnippet(lesson, 260); snippet != "" {
			lines = append(lines, strconv.Itoa(index+1)+") "+snippet)
		}
	}
	return strings.Join(lines, " | ")
}

func lessonTopicInstruction(level string, focus string, lessonCount int, recentLessons []string) string {
	if lessonCount < 0 {
		lessonCount = 0
	}
	parts := []string{
		"Lesson novelty rule: never repeat any of the last 10 lessons, their topics, situations, example sentences, or communicative tasks.",
		"Lesson number: " + strconv.Itoa(lessonCount+1) + "; use it as a rotation seed for topic choice.",
		"CEFR topic pool for this level: " + cefrTopicPool(level) + ".",
		"If no global learning focus is set, choose a fresh topic from the CEFR pool and rotate the topic every lesson.",
		"If a global learning focus is set, keep that preference but vary the sub-situation, grammar pattern, vocabulary set, and learner task every time.",
		"Do not default to airport, hotel, restaurant, or travel unless the rotation seed or learning focus clearly points there.",
	}
	if recent := recentLessonContext(recentLessons); recent != "" {
		parts = append(parts, "Recent lessons to avoid: "+recent)
	}
	return strings.Join(parts, " ") + " "
}

func practiceTopicInstruction(level string, focus string, practiceCount int, recentMessages []string) string {
	if practiceCount < 0 {
		practiceCount = 0
	}
	parts := []string{
		"Practice novelty rule: if you introduce a new mini-situation, do not reuse the same topic, example, or question pattern from recent context.",
		"Practice turn number: " + strconv.Itoa(practiceCount+1) + "; use it as a rotation seed.",
		"CEFR topic pool for optional examples: " + cefrTopicPool(level) + ".",
		"If the learner gives only a short or generic message, choose a fresh level-appropriate topic from that pool instead of repeating airport/hotel/travel.",
	}
	if strings.TrimSpace(focus) != "" {
		parts = append(parts, "The global focus is only a preference; vary the concrete scene and task inside it.")
	}
	if len(recentMessages) > 0 {
		parts = append(parts, "Recent context is listed below; avoid copying its scenario.")
	}
	return strings.Join(parts, " ") + " "
}

func lessonPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, focus string, lessonCount int, recentLessons []string) []chatMessage {
	labels := localizedStudyPromptLabels(interfaceLanguage)
	variationSeed := promptVariationSeed("lesson", level, lessonCount, focus, recentLessons)
	variationInstruction := promptVariationInstruction(variationSeed)
	userPrompt := "Create one short " + language.NativeName + " practice task for level " + level + ". " +
		learningFocusInstruction(focus) +
		lessonTopicInstruction(level, focus, lessonCount, recentLessons) +
		variationInstruction +
		"The task must fit Telegram and use active recall, not passive explanation. Include: a tiny situation, one tiny teacher notice about the pattern, 3-5 useful words or chunks, " +
		"one short model example sentence in " + language.NativeName + " only, " +
		"and one clear instruction for the learner to produce a fresh answer in " + language.NativeName + ". " +
		"Make the task feel like a useful real-life micro-scenario, not a school worksheet. " +
		"Write almost everything in " + interfaceLanguage.NativeName + ". Do not say that you are waiting for many answers. " +
		"Use clean plain-text formatting with exactly these localized short block labels: " + lessonPromptLabelList(labels) + ". Put each chunk on its own line with bullet symbols. Do not use English section labels unless the interface language is English. Do not use design/style words such as saturation, gradient, prompt, or UI. Do not use Markdown or asterisks. " +
		"The " + labels.ModelPhrase + " block must contain only one clean " + language.NativeName + " sentence for pronunciation. Do not add translations, transliteration, explanations, or parenthesized interface-language text inside that block. " +
		"Make it clear that this is one short task. Keep under 900 characters."
	return []chatMessage{
		{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)},
		{
			Role: "user",
			Content: renderAppPrompt("lesson.generate.user", userPrompt, mergePromptVars(commonPromptVars(language, interfaceLanguage), map[string]string{
				"level":                      level,
				"lesson_prompt_labels":       lessonPromptLabelList(labels),
				"learning_focus_instruction": learningFocusInstruction(focus),
				"lesson_topic_instruction":   lessonTopicInstruction(level, focus, lessonCount, recentLessons),
				"lesson_number":              strconv.Itoa(lessonCount + 1),
				"recent_lesson_context":      recentLessonContext(recentLessons),
				"variation_seed":             variationSeed,
				"variation_instruction":      variationInstruction,
			})),
		},
	}
}

// feedbackPrompt asks the AI to return feedback AND a JSON block with mistakes.
// The JSON block is appended after a sentinel line so we can split it off cleanly.
func feedbackPrompt(language learningLanguage, interfaceLanguage learningLanguage, level, task, answer string) []chatMessage {
	userPrompt := "Target language: " + language.NativeName + "\nLevel: " + level + "\nTask:\n" + task + "\n\nLearner answer:\n" + answer + "\n\n" +
		"Check the answer and finish the task. Reply mostly in " + interfaceLanguage.NativeName + ". Return exactly three compact numbered sections in " + interfaceLanguage.NativeName + ": " +
		"1) corrected version with only the natural target-language answer, 2) one highest-impact thing to improve with a short reason, 3) mini-summary with one repeatable phrase from the correction. " +
		"Prioritize a correction the learner can reuse tomorrow over a long grammar lecture. " +
		"Do not ask a follow-up question. Do not continue the lesson. Keep the feedback under 1000 characters. " +
		"The feedback must be plain text only: no Markdown, no asterisks, no # headings.\n\n" +
		"After the feedback, on a new line write exactly: ---MISTAKES---\n" +
		"Then write a JSON array of mistakes found in the learner's answer. " +
		"Each object must have these keys: \"word\" (the wrong word or phrase the learner wrote), " +
		"\"correction\" (the correct version in " + language.NativeName + "), \"explanation\" (one short sentence in " + interfaceLanguage.NativeName + " why it is wrong). " +
		"If there are no mistakes, write an empty array: []. Example:\n" +
		"[{\"word\":\"I goed\",\"correction\":\"I went\",\"explanation\":\"Use the irregular past form went.\"}]"
	return []chatMessage{
		{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)},
		{
			Role: "user",
			Content: renderAppPrompt("lesson.feedback.user", userPrompt, mergePromptVars(commonPromptVars(language, interfaceLanguage), map[string]string{
				"level":          level,
				"task":           task,
				"learner_answer": answer,
			})),
		},
	}
}

const practiceMemoryLimit = 5
const lessonMemoryLimit = 10

// practicePrompt asks the AI to return the reply AND a JSON mistakes block.
func practicePrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, recentMessages []string, focus string, practiceCount int) []chatMessage {
	userMessage := ""
	if len(recentMessages) > 0 {
		userMessage = recentMessages[len(recentMessages)-1]
	}
	context := practiceContextText(recentMessages)
	labels := localizedStudyPromptLabels(interfaceLanguage)
	variationSeed := promptVariationSeed("practice", level, practiceCount, focus, recentMessages)
	variationInstruction := promptVariationInstruction(variationSeed)
	userPrompt := "The learner level is " + level + ". Continue a friendly " + language.NativeName + " practice chat. " +
		learningFocusInstruction(focus) +
		practiceTopicInstruction(level, focus, practiceCount, recentMessages) +
		variationInstruction +
		"Use the recent learner messages as short context, but answer the latest message. " +
		"Reply mostly in " + interfaceLanguage.NativeName + " with one short " + language.NativeName + " phrase or question. " +
		"If there are mistakes, briefly correct only the most useful one in " + interfaceLanguage.NativeName + "; if a previous mistake or weak phrase appears in the recent context, recycle it naturally once. " +
		"Keep the learner speaking: model one natural chunk, then ask for one small answer that reuses it. " +
		"The " + labels.ModelPhrase + " block must contain only one clean " + language.NativeName + " sentence for pronunciation. Do not add translations, transliteration, explanations, or parenthesized interface-language text inside that block. " +
		"Ask only one simple " + language.NativeName + " question at the end, and put that question on its own final line. " +
		"Use clean plain text with visibly separated blocks using exactly these localized labels: " + practicePromptLabelList(labels) + ". Use bullets when listing alternatives. Do not use English section labels unless the interface language is English. Do not use design/style words such as saturation, gradient, prompt, or UI. Do not send one dense paragraph. Use 0-2 relevant emojis. Do not use Markdown or asterisks in your reply. " +
		context +
		"Latest learner message: " + userMessage + "\n\n" +
		"After your reply, on a new line write exactly: ---MISTAKES---\n" +
		"Then write a JSON array of grammar or vocabulary mistakes found in the learner's message. " +
		"Each object: \"word\" (wrong phrase), \"correction\" (correct version in " + language.NativeName + "), " +
		"\"explanation\" (one sentence in " + interfaceLanguage.NativeName + "). If no mistakes, write []."
	return []chatMessage{
		{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)},
		{
			Role: "user",
			Content: renderAppPrompt("practice.chat.user", userPrompt, mergePromptVars(commonPromptVars(language, interfaceLanguage), map[string]string{
				"level":                      level,
				"recent_context":             context,
				"latest_learner_message":     userMessage,
				"practice_prompt_labels":     practicePromptLabelList(labels),
				"model_phrase_label":         labels.ModelPhrase,
				"learning_focus_instruction": learningFocusInstruction(focus),
				"practice_topic_instruction": practiceTopicInstruction(level, focus, practiceCount, recentMessages),
				"practice_turn_number":       strconv.Itoa(practiceCount + 1),
				"variation_seed":             variationSeed,
				"variation_instruction":      variationInstruction,
			})),
		},
	}
}

func mistakeExtractionPrompt(language learningLanguage, interfaceLanguage learningLanguage, source, task, learnerAnswer, feedback string) []chatMessage {
	source = strings.TrimSpace(source)
	if source == "" {
		source = "learning answer"
	}
	userPrompt := "Extract only the learner's real grammar, vocabulary, spelling, or phrase-choice mistakes from this " + source + ".\n" +
		"Target language: " + language.NativeName + "\n" +
		"Interface/explanation language: " + interfaceLanguage.NativeName + "\n" +
		"Task or context:\n" + strings.TrimSpace(task) + "\n\n" +
		"Learner answer:\n" + strings.TrimSpace(learnerAnswer) + "\n\n" +
		"Teacher feedback, if any:\n" + strings.TrimSpace(feedback) + "\n\n" +
		"Return JSON only. Return an array of objects with keys \"word\", \"correction\", \"explanation\". " +
		"\"word\" is the exact wrong word or phrase the learner wrote. \"correction\" is the corrected target-language phrase. " +
		"\"explanation\" is one short sentence in " + interfaceLanguage.NativeName + ". " +
		"If the learner answer is already correct or there is no clear mistake, return []."
	return []chatMessage{
		{Role: "system", Content: "You extract language-learning mistakes into strict JSON. Do not include Markdown, prose, or code fences."},
		{Role: "user", Content: userPrompt},
	}
}

func roleplayPrompt(language learningLanguage, interfaceLanguage learningLanguage, scenarioPayload string, focus string) []chatMessage {
	payload := strings.TrimSpace(scenarioPayload)
	if instruction := learningFocusInstruction(focus); instruction != "" {
		payload += "\n" + instruction
	}
	variationSeed := promptVariationSeed("roleplay", "A2", 0, focus, []string{payload})
	variationInstruction := promptVariationInstruction(variationSeed)
	if variationInstruction != "" {
		payload += "\n" + variationInstruction
	}
	fallback := payload + "\n\nAfter the learner-facing roleplay answer, on a new line write exactly: ---MISTAKES---\nThen write a JSON array of any learner mistakes if the payload includes a learner line; otherwise write []."
	return []chatMessage{
		{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage) + "\nYou are now inside the dedicated AI roleplay tool. Do not use the generic practice-chat structure unless the roleplay payload explicitly requests it."},
		{
			Role: "user",
			Content: renderAppPrompt("roleplay.scenario.user", fallback, mergePromptVars(commonPromptVars(language, interfaceLanguage), map[string]string{
				"roleplay_payload":      payload,
				"variation_seed":        variationSeed,
				"variation_instruction": variationInstruction,
			})),
		},
	}
}

func practiceContextText(recentMessages []string) string {
	if len(recentMessages) == 0 {
		return ""
	}
	text := "Recent context and learner messages:\n"
	for i, message := range recentMessages {
		text += itoa(i+1) + ") " + message + "\n"
	}
	return text + "\n"
}

func practiceImageMessage(learnerText, imageContext string) string {
	text := strings.TrimSpace(learnerText)
	context := strings.TrimSpace(imageContext)
	if context == "" {
		return text
	}
	if text == "" {
		return "Learner attached an image for practice.\nCoach-only image context, do not grade this as learner writing:\n" + context
	}
	return "Learner attached an image for practice.\nLearner note:\n" + text + "\n\nCoach-only image context, do not grade this as learner writing:\n" + context
}

func translationPrompt(text string, interfaceLanguage learningLanguage) []chatMessage {
	language := learningLanguageByCode("en")
	userPrompt := "Translate this text into natural " + interfaceLanguage.NativeName + ". Return only the translation, no Markdown, no explanations.\n\nText:\n" + text
	return []chatMessage{
		{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)},
		{
			Role: "user",
			Content: renderAppPrompt("tools.translation_plain.user", userPrompt, map[string]string{
				"interface_language_native_name": interfaceLanguage.NativeName,
				"text":                           text,
			}),
		},
	}
}

func vocabularyHintPrompt(word vocabWord, prompt string, learningLanguage learningLanguage, interfaceLanguage learningLanguage, mode string) []chatMessage {
	mode = strings.TrimSpace(mode)
	if mode == "" {
		mode = "word study"
	}
	systemPrompt := "You write ultra-short context clues for hidden target-language words in a language-learning bot. " +
		"Reply only in " + interfaceLanguage.NativeName + ". " +
		"Describe the target word's sense or usage, not the interface-language prompt text. " +
		"Do not reveal or write the target-language answer. No Markdown. No lists."
	userPrompt := "Interface language: " + interfaceLanguage.NativeName + "\n" +
		"Target learning language: " + learningLanguage.NativeName + "\n" +
		"Mode: " + mode + "\n" +
		"Interface prompt shown to learner: " + prompt + "\n" +
		"Hidden target word to learn: " + word.English + "\n\n" +
		"Give one tiny context clue for the hidden target word, maximum 8 words. " +
		"Help the learner understand which " + learningLanguage.NativeName + " word is meant by usage or context. " +
		"Do not explain the interface prompt itself. Do not translate into " + learningLanguage.NativeName + ". Do not include the hidden target word."
	return []chatMessage{
		{
			Role:    "system",
			Content: renderAppPrompt("vocabulary.hint.system", systemPrompt, commonPromptVars(learningLanguage, interfaceLanguage)),
		},
		{
			Role: "user",
			Content: renderAppPrompt("vocabulary.hint.user", userPrompt, mergePromptVars(commonPromptVars(learningLanguage, interfaceLanguage), map[string]string{
				"mode":             mode,
				"interface_prompt": prompt,
				"target_word":      word.English,
			})),
		},
	}
}

func vocabularyGenerationPrompt(user userState, promptLanguage learningLanguage, forbidden []string) []chatMessage {
	learningLanguage := userLearningLanguage(user)
	level := normalizeCEFRLevel(user.Level)
	cleanForbidden := make([]string, 0, len(forbidden))
	seen := map[string]bool{}
	for _, item := range forbidden {
		item = cleanDictionaryDisplay(item)
		if item == "" || seen[strings.ToLower(item)] {
			continue
		}
		seen[strings.ToLower(item)] = true
		cleanForbidden = append(cleanForbidden, item)
		if len(cleanForbidden) >= 30 {
			break
		}
	}
	forbiddenText := "none"
	if len(cleanForbidden) > 0 {
		forbiddenText = strings.Join(cleanForbidden, ", ")
	}
	systemPrompt := "You are a professional CEFR lexicographer for a language-learning app. " +
		"Return strict JSON only. No Markdown. No explanations. " +
		"Generate one useful high-frequency vocabulary card, not a rare dictionary item. " +
		"The card must be safe for a multiple-choice lesson: the translation must not reveal or repeat the target word. " +
		"Never choose a forbidden word, spelling variant, inflected form, proper name, abbreviation, phrase, or vulgar word."
	userPrompt := "Target learning language: " + learningLanguage.NativeName + "\n" +
		"Learner prompt/interface language: " + promptLanguage.NativeName + "\n" +
		"Requested CEFR level: " + level + "\n" +
		"Forbidden words: " + forbiddenText + "\n\n" +
		"Return exactly this JSON shape: " +
		"{\"word\":\"target lemma\",\"translation\":\"short meaning in " + promptLanguage.NativeName + "\",\"russian\":\"short Russian meaning\",\"context\":\"optional short sense note in " + promptLanguage.NativeName + "\",\"level\":\"" + level + "\",\"topic\":\"everyday topic\",\"part_of_speech\":\"noun|verb|adjective|adverb\"}.\n" +
		"Rules: word must be in " + learningLanguage.NativeName + "; translation must be in " + promptLanguage.NativeName + "; russian must be a Russian meaning; level must be one of A1, A2, B1, B2, C1, C2; do not include examples."
	return []chatMessage{
		{
			Role:    "system",
			Content: renderAppPrompt("vocabulary.generation.system", systemPrompt, commonPromptVars(learningLanguage, promptLanguage)),
		},
		{
			Role: "user",
			Content: renderAppPrompt("vocabulary.generation.user", userPrompt, mergePromptVars(commonPromptVars(learningLanguage, promptLanguage), map[string]string{
				"level":     level,
				"forbidden": forbiddenText,
			})),
		},
	}
}

func vocabularyTranslationPrompt(word vocabWord, learningLanguage learningLanguage, targetLanguage learningLanguage, mode string) []chatMessage {
	mode = strings.TrimSpace(mode)
	if mode == "" {
		mode = "word study"
	}
	known := []string{}
	if value := cleanDictionaryDisplay(word.Russian); value != "" {
		known = append(known, "ru: "+value)
	}
	if word.Translations != nil {
		codes := make([]string, 0, len(word.Translations))
		for code := range word.Translations {
			codes = append(codes, normalizeInterfaceLanguage(code))
		}
		sort.Strings(codes)
		for _, code := range codes {
			value := cleanDictionaryDisplay(word.Translations[code])
			if value != "" {
				known = append(known, code+": "+value)
			}
		}
	}
	knownText := "none"
	if len(known) > 0 {
		knownText = strings.Join(known, " | ")
	}
	systemPrompt := "You are a professional lexicographer for a CEFR language-learning app. " +
		"Return strict JSON only. No Markdown. No explanations. " +
		"Translate the hidden target-language lemma into the learner's native/interface language. " +
		"Choose the common dictionary sense that best matches the part of speech, CEFR level, topic and existing glosses. " +
		"Translate the meaning, not the spelling; false friends, hyphenated lemmas and borrowed-looking words still need a native semantic equivalent. " +
		"If the word is ambiguous, return up to 3 short native-language equivalents separated by '; '. " +
		"Never return the hidden target-language word itself. If the native language commonly borrows the same surface form, use a native paraphrase or short explanation instead. " +
		"Do not include the hidden target-language word, transliteration, examples, grammar notes, or source-language labels."
	userPrompt := "Target learning language: " + learningLanguage.NativeName + "\n" +
		"Learner native/interface language: " + targetLanguage.NativeName + "\n" +
		"Mode: " + mode + "\n" +
		"CEFR level: " + strings.TrimSpace(word.Level) + "\n" +
		"Part of speech: " + strings.TrimSpace(word.PartOfSpeech) + "\n" +
		"Topic: " + strings.TrimSpace(word.Topic) + "\n" +
		"Hidden target word: " + strings.TrimSpace(word.English) + "\n" +
		"Existing dictionary glosses: " + knownText + "\n" +
		"Context: " + strings.TrimSpace(word.Context) + "\n\n" +
		"Return exactly this JSON shape: {\"translation\":\"main meaning; optional second meaning; optional third meaning\"}. " +
		"Each meaning must be natural in " + targetLanguage.NativeName + ", short enough for a multiple-choice prompt, and must not equal or repeat the hidden target word."
	return []chatMessage{
		{
			Role:    "system",
			Content: renderAppPrompt("vocabulary.translation.system", systemPrompt, commonPromptVars(learningLanguage, targetLanguage)),
		},
		{
			Role: "user",
			Content: renderAppPrompt("vocabulary.translation.user", userPrompt, mergePromptVars(commonPromptVars(learningLanguage, targetLanguage), map[string]string{
				"mode":        mode,
				"target_word": word.English,
				"level":       strings.TrimSpace(word.Level),
				"topic":       strings.TrimSpace(word.Topic),
				"known":       knownText,
			})),
		},
	}
}

func vocabularyExamplePrompt(word vocabWord, learningLanguage learningLanguage, interfaceLanguage learningLanguage, level string, mode string) []chatMessage {
	mode = strings.TrimSpace(mode)
	if mode == "" {
		mode = "word practice"
	}
	level = strings.TrimSpace(level)
	if level == "" {
		level = "A1"
	}
	systemPrompt := "You write one short natural example sentence for a language learner. " +
		"Use the target word exactly as appropriate in " + learningLanguage.NativeName + ". " +
		"No Markdown. No lists. No explanation."
	userPrompt := "Target learning language: " + learningLanguage.NativeName + "\n" +
		"Interface language: " + interfaceLanguage.NativeName + "\n" +
		"CEFR level: " + level + "\n" +
		"Mode: " + mode + "\n" +
		"Target word: " + word.English + "\n\n" +
		"Return one short everyday " + learningLanguage.NativeName + " example sentence using the target word. " +
		"Match the CEFR level: A1-A2 use very simple grammar and common context; B1-B2 can use a richer but still clear sentence; C1-C2 may sound natural and nuanced. " +
		"Keep it under 110 characters when possible. Return only the sentence."
	return []chatMessage{
		{
			Role:    "system",
			Content: renderAppPrompt("vocabulary.example.system", systemPrompt, commonPromptVars(learningLanguage, interfaceLanguage)),
		},
		{
			Role: "user",
			Content: renderAppPrompt("vocabulary.example.user", userPrompt, mergePromptVars(commonPromptVars(learningLanguage, interfaceLanguage), map[string]string{
				"level":       level,
				"mode":        mode,
				"target_word": word.English,
			})),
		},
	}
}

func trimStringHistory(history []string, limit int) []string {
	if limit <= 0 {
		return nil
	}
	trimmed := make([]string, 0, len(history))
	for _, message := range history {
		message = strings.TrimSpace(message)
		if message != "" {
			trimmed = append(trimmed, message)
		}
	}
	if len(trimmed) > limit {
		trimmed = trimmed[len(trimmed)-limit:]
	}
	return trimmed
}

func trimPracticeHistory(history []string, limit int) []string {
	return trimStringHistory(history, limit)
}

func trimLessonHistory(history []string) []string {
	return trimStringHistory(history, lessonMemoryLimit)
}
