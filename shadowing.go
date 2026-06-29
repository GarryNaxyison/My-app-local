package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"
	"unicode"
)

const (
	modeShadowingPrefix     = "shadowing:"
	modePronunciationPrefix = "pronunciation:"
	shadowingMaxPhraseRunes = 160
	shadowingDeckSize       = 10000
)

func shadowingMode(phrase string) string {
	phrase = sanitizeShadowingPhrase(phrase)
	if phrase == "" {
		phrase = "Could you say that a little slower, please?"
	}
	return modeShadowingPrefix + base64.RawURLEncoding.EncodeToString([]byte(phrase))
}

func parseShadowingMode(mode string) (string, bool) {
	return parseEncodedPhraseMode(mode, modeShadowingPrefix)
}

func pronunciationPracticeMode(phrase string) string {
	phrase = sanitizeShadowingPhrase(phrase)
	if phrase == "" {
		phrase = "Could you say that a little slower, please?"
	}
	return modePronunciationPrefix + base64.RawURLEncoding.EncodeToString([]byte(phrase))
}

func parsePronunciationMode(mode string) (string, bool) {
	return parseEncodedPhraseMode(mode, modePronunciationPrefix)
}

func parseEncodedPhraseMode(mode string, prefix string) (string, bool) {
	if !strings.HasPrefix(mode, prefix) {
		return "", false
	}
	raw := strings.TrimPrefix(mode, prefix)
	if raw == "" {
		return "", false
	}
	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return "", false
	}
	phrase := sanitizeShadowingPhrase(string(decoded))
	return phrase, phrase != ""
}

func sanitizeShadowingPhrase(raw string) string {
	raw = strings.TrimSpace(cleanModelReply(raw))
	if raw == "" {
		return ""
	}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = strings.TrimPrefix(line, "-")
		line = strings.TrimSpace(line)
		if dot := strings.Index(line, ". "); dot > 0 && dot <= 2 {
			allDigits := true
			for _, r := range line[:dot] {
				if !unicode.IsDigit(r) {
					allDigits = false
					break
				}
			}
			if allDigits {
				line = strings.TrimSpace(line[dot+1:])
			}
		}
		lower := strings.ToLower(line)
		for _, prefix := range []string{"phrase:", "target:", "shadowing:", "фраза:", "цель:"} {
			if strings.HasPrefix(lower, prefix) {
				line = strings.TrimSpace(line[len(prefix):])
				break
			}
		}
		line = strings.Trim(line, " \t\r\n\"'`*“”«»")
		if line == "" {
			continue
		}
		runes := []rune(line)
		if len(runes) > shadowingMaxPhraseRunes {
			line = string(runes[:shadowingMaxPhraseRunes])
		}
		return strings.TrimSpace(line)
	}
	return ""
}

var shadowingContexts = []string{
	"ordering something simple", "asking for directions", "checking a reservation", "meeting a new colleague", "buying a ticket",
	"confirming a time", "describing a daily routine", "making a small request", "apologizing politely", "asking for help",
	"talking about the weather", "planning a short trip", "shopping for clothes", "explaining a small problem", "calling a service desk",
	"speaking with a teacher", "chatting before work", "choosing food", "talking about health", "arranging a delivery",
}

var shadowingMoves = []string{
	"ask one clear question", "make one polite request", "confirm one detail", "decline gently", "offer one option",
	"say what you need", "ask for repetition", "give a short reason", "compare two choices", "invite a response",
	"state a preference", "ask about price", "ask about location", "say what happened", "suggest a next step",
	"thank and close", "check understanding", "ask for permission", "describe a plan", "give a simple warning",
	"ask for an opinion", "offer help", "say you are ready", "ask about timing", "summarize one decision",
}

var shadowingDetails = []string{
	"keep it friendly", "use present tense if possible", "use a natural spoken rhythm", "make it useful for travel",
	"make it useful at work", "include a time word", "include a place word", "include a polite softener",
	"include a common verb", "avoid rare vocabulary", "sound like a real conversation", "make it easy to repeat aloud",
	"fit a beginner learner", "use one compact phrase", "include a small contrast", "make it emotionally neutral",
	"make it slightly urgent", "make it calm and confident", "avoid idioms", "avoid names and brands",
}

func shadowingDeckSpec(index int) string {
	if index < 0 {
		index = -index
	}
	index %= shadowingDeckSize
	context := shadowingContexts[index%len(shadowingContexts)]
	move := shadowingMoves[(index/len(shadowingContexts))%len(shadowingMoves)]
	detail := shadowingDetails[(index/(len(shadowingContexts)*len(shadowingMoves)))%len(shadowingDetails)]
	return fmt.Sprintf("Card %04d/%d. Situation: %s. Communicative move: %s. Constraint: %s.", index+1, shadowingDeckSize, context, move, detail)
}

func shadowingDeckIndex(user userState) int {
	seed := uint64(user.TelegramID)
	for _, r := range normalizeLearningLanguage(user.LearningLanguage) + "|" + strings.TrimSpace(user.Level) {
		seed = seed*131 + uint64(r)
	}
	now := time.Now().UTC()
	seed += uint64(user.PracticeCount*37 + user.VoiceCount*53 + user.LessonCount*17)
	seed += uint64(now.YearDay()*97 + now.Hour()*11 + now.Minute()/6)
	return int(seed % shadowingDeckSize)
}

func shadowingPhrasePrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, card string, previous string) []chatMessage {
	avoid := ""
	if strings.TrimSpace(previous) != "" {
		avoid = "\nDo not repeat or lightly paraphrase this previous phrase: " + strings.TrimSpace(previous)
	}
	userPrompt := "Create exactly one natural " + language.NativeName + " phrase for a listening-and-repeat drill at CEFR level " + level + ". " +
		"Use this semantic card so the drill comes from a 10,000-card phrase deck shared across all languages:\n" + card + avoid + "\n\n" +
		"The learner will listen, repeat aloud, and compare speech recognition to the phrase. " +
		"Requirements: 5-11 words when the language uses spaces; one everyday spoken chunk; useful rhythm; no rare vocabulary; no translation; no explanation; no numbering; no quotation marks; no Markdown. " +
		"Return only the phrase in " + language.NativeName + ". If the target language is not English, do not use English words unless they are unavoidable international words. " +
		"Return only the target-language phrase."
	return []chatMessage{
		{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)},
		{
			Role: "user",
			Content: renderAppPrompt("listening.phrase.user", userPrompt, mergePromptVars(commonPromptVars(language, interfaceLanguage), map[string]string{
				"level":           level,
				"deck_card":       card,
				"previous_phrase": previous,
			})),
		},
	}
}

func pronunciationPhrasePrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, card string, previous string) []chatMessage {
	avoid := ""
	if strings.TrimSpace(previous) != "" {
		avoid = "\nDo not repeat or lightly paraphrase this previous phrase: " + previous
	}
	userPrompt := "Create exactly one natural " + language.NativeName + " phrase for a pronunciation drill at CEFR level " + level + ". " +
		"The learner will listen to an audio model, record the same phrase, and receive pronunciation scoring.\n" +
		"Semantic card:\n" + card + avoid + "\n\n" +
		"Requirements: 5-10 words when the language uses spaces; everyday spoken phrase; useful sounds and rhythm; no rare vocabulary; no translation; no explanation; no numbering; no quotation marks; no Markdown. " +
		"Return only the phrase in " + language.NativeName + "."
	return []chatMessage{
		{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)},
		{
			Role: "user",
			Content: renderAppPrompt("pronunciation.phrase.user", userPrompt, mergePromptVars(commonPromptVars(language, interfaceLanguage), map[string]string{
				"level":           level,
				"deck_card":       card,
				"previous_phrase": previous,
			})),
		},
	}
}

func shadowingFeedbackPrompt(language learningLanguage, interfaceLanguage learningLanguage, level string, target string, transcript string, score int, missing string, fromVoice bool) []chatMessage {
	source := "Voice transcript"
	if !fromVoice {
		source = "Typed fallback, not a real pronunciation sample"
	}
	userPrompt := "Target language: " + language.NativeName + "\nInterface language: " + interfaceLanguage.NativeName + "\nLevel: " + level + "\n" +
		"Input source: " + source + "\nListening target phrase:\n" + target + "\n\nSpeech recognizer heard or learner typed:\n" + transcript + "\n\nLocal strict similarity score: " + itoa(score) + "/100\nMissing or weak words: " + missing + "\n\n" +
		"Give compact listening feedback in " + interfaceLanguage.NativeName + ". " +
		"Do not repeat the numeric score; the app shows it separately. " +
		"If the source is typed fallback, say that pronunciation cannot be judged from text and ask for a voice repeat next time. " +
		"If the score is below 70, do not praise it as excellent or perfect; clearly say it needs another repeat. " +
		"Return exactly three short lines: 1) what matched, 2) one concrete pronunciation or rhythm tip, 3) repeat line with the exact target phrase. " +
		"No Markdown, no asterisks, no headings, no JSON, no follow-up question."
	return []chatMessage{
		{Role: "system", Content: coachSystemPrompt(language, interfaceLanguage)},
		{
			Role: "user",
			Content: renderAppPrompt("listening.feedback.user", userPrompt, mergePromptVars(commonPromptVars(language, interfaceLanguage), map[string]string{
				"level":        level,
				"input_source": source,
				"target":       target,
				"transcript":   transcript,
				"score":        itoa(score),
				"missing":      missing,
			})),
		},
	}
}

var shadowingFallbackPhrases = map[string][]string{
	"en": {"Please say that a little slower.", "I need a few minutes to check the details.", "Let us meet near the main entrance."},
	"ru": {"Повторите, пожалуйста, немного медленнее.", "Мне нужно несколько минут, чтобы проверить детали.", "Встретимся у главного входа."},
	"es": {"Por favor, repítelo un poco más despacio.", "Necesito unos minutos para revisar los detalles.", "Quedemos cerca de la entrada principal."},
	"de": {"Sagen Sie das bitte etwas langsamer.", "Ich brauche ein paar Minuten, um die Details zu prüfen.", "Treffen wir uns am Haupteingang."},
	"fr": {"Répétez un peu plus lentement, s’il vous plaît.", "J’ai besoin de quelques minutes pour vérifier les détails.", "Retrouvons-nous près de l’entrée principale."},
	"it": {"Ripetilo un po’ più lentamente, per favore.", "Mi servono alcuni minuti per controllare i dettagli.", "Incontriamoci vicino all’ingresso principale."},
	"zh": {"请说得稍微慢一点。", "我需要几分钟来核对细节。", "我们在主入口附近见面吧。"},
	"ja": {"もう少しゆっくり言ってください。", "詳細を確認するのに数分ください。", "正面入口の近くで会いましょう。"},
	"ko": {"조금 더 천천히 말해 주세요.", "세부 사항을 확인할 시간이 몇 분 필요합니다.", "정문 근처에서 만나요."},
	"tg": {"Лутфан, инро каме сусттар такрор кунед.", "Барои санҷидани ҷузъиёт ба ман чанд дақиқа лозим аст.", "Биёед назди даромадгоҳи асосӣ вохӯрем."},
	"uz": {"Iltimos, biroz sekinroq takrorlang.", "Tafsilotlarni tekshirish uchun bir necha daqiqa kerak.", "Asosiy kirish yonida uchrashamiz."},
	"tt": {"Зинһар, моны бераз акрынрак кабатла.", "Детальләрне тикшерү өчен миңа берничә минут кирәк.", "Төп керү урыны янында очрашыйк."},
	"hy": {"Խնդրում եմ, մի փոքր ավելի դանդաղ կրկնիր։", "Մանրամասները ստուգելու համար մի քանի րոպե է պետք։", "Հանդիպենք գլխավոր մուտքի մոտ։"},
	"kk": {"Өтінемін, мұны сәл баяуырақ қайтала.", "Мәліметтерді тексеру үшін маған бірнеше минут керек.", "Басты кіреберістің жанында кездесейік."},
	"ky": {"Суранам, муну бир аз жайыраак кайтала.", "Чоо-жайын текшерүүгө мага бир нече мүнөт керек.", "Негизги кире бериштин жанынан жолугалы."},
	"ka": {"გთხოვ, ეს ცოტა ნელა გაიმეორე.", "დეტალების შესამოწმებლად რამდენიმე წუთი მჭირდება.", "მთავარ შესასვლელთან შევხვდეთ."},
	"uk": {"Повтори, будь ласка, трохи повільніше.", "Мені потрібно кілька хвилин, щоб перевірити деталі.", "Зустріньмося біля головного входу."},
	"pl": {"Powtórz trochę wolniej, proszę.", "Potrzebuję kilku minut, żeby sprawdzić szczegóły.", "Spotkajmy się przy głównym wejściu."},
	"ro": {"Repetă puțin mai încet, te rog.", "Am nevoie de câteva minute ca să verific detaliile.", "Să ne întâlnim lângă intrarea principală."},
	"pt": {"Repita um pouco mais devagar, por favor.", "Preciso de alguns minutos para verificar os detalhes.", "Vamos nos encontrar perto da entrada principal."},
	"ar": {"من فضلك قل ذلك ببطء أكثر قليلا.", "أحتاج إلى بضع دقائق لأراجع التفاصيل.", "لنلتق قرب المدخل الرئيسي."},
	"bn": {"দয়া করে একটু ধীরে বলুন।", "বিস্তারিত যাচাই করতে আমার কয়েক মিনিট দরকার।", "চলুন প্রধান প্রবেশদ্বারের কাছে দেখা করি।"},
	"cs": {"Řekni to prosím trochu pomaleji.", "Potřebuji pár minut na kontrolu detailů.", "Sejdeme se u hlavního vchodu."},
	"el": {"Πες το λίγο πιο αργά, σε παρακαλώ.", "Χρειάζομαι λίγα λεπτά για να ελέγξω τις λεπτομέρειες.", "Ας συναντηθούμε κοντά στην κύρια είσοδο."},
	"hi": {"कृपया इसे थोड़ा धीरे बोलें।", "विवरण जाँचने के लिए मुझे कुछ मिनट चाहिए।", "मुख्य प्रवेश द्वार के पास मिलते हैं।"},
	"hu": {"Mondd egy kicsit lassabban, kérlek.", "Kell néhány perc, hogy ellenőrizzem a részleteket.", "Találkozzunk a főbejáratnál."},
	"id": {"Tolong ucapkan sedikit lebih pelan.", "Saya perlu beberapa menit untuk memeriksa detailnya.", "Kita bertemu di dekat pintu masuk utama."},
	"nl": {"Zeg dat alsjeblieft iets langzamer.", "Ik heb een paar minuten nodig om de details te controleren.", "Laten we afspreken bij de hoofdingang."},
	"sv": {"Säg det lite långsammare, tack.", "Jag behöver några minuter för att kontrollera detaljerna.", "Vi ses vid huvudentrén."},
	"ta": {"தயவுசெய்து அதை கொஞ்சம் மெதுவாக சொல்லுங்கள்.", "விவரங்களைச் சரிபார்க்க எனக்கு சில நிமிடங்கள் தேவை.", "முதன்மை நுழைவாயிலருகே சந்திப்போம்."},
	"te": {"దయచేసి దాన్ని కొంచెం నెమ్మదిగా చెప్పండి.", "వివరాలను చూసేందుకు నాకు కొన్ని నిమిషాలు కావాలి.", "ప్రధాన ప్రవేశద్వారం దగ్గర కలుద్దాం."},
	"th": {"กรุณาพูดให้ช้าลงอีกเล็กน้อย.", "ฉันต้องใช้เวลาสักครู่เพื่อตรวจสอบรายละเอียด.", "เจอกันใกล้ทางเข้าหลัก."},
	"tl": {"Pakisabi iyon nang kaunti pang mabagal.", "Kailangan ko ng ilang minuto para suriin ang mga detalye.", "Magkita tayo malapit sa pangunahing pasukan."},
	"tr": {"Lütfen bunu biraz daha yavaş söyle.", "Ayrıntıları kontrol etmek için birkaç dakikaya ihtiyacım var.", "Ana girişin yakınında buluşalım."},
	"vi": {"Vui lòng nói câu đó chậm hơn một chút.", "Tôi cần vài phút để kiểm tra chi tiết.", "Mình gặp nhau gần lối vào chính."},
}

func shadowingFallbackPhrase(language learningLanguage, index int) string {
	phrases := shadowingFallbackPhrases[normalizeLearningLanguage(language.Code)]
	if len(phrases) == 0 {
		phrases = shadowingFallbackPhrases["en"]
	}
	if len(phrases) == 0 {
		return "Could you say that a little slower, please?"
	}
	if index < 0 {
		index = -index
	}
	return phrases[index%len(phrases)]
}

func shadowingLooksWrongScript(phrase string, language learningLanguage) bool {
	code := normalizeLearningLanguage(language.Code)
	if code == "en" || code == "es" || code == "de" || code == "fr" || code == "it" || code == "pt" || code == "pl" || code == "ro" || code == "uz" {
		return false
	}
	hasExpected := false
	for _, r := range phrase {
		switch code {
		case "zh":
			hasExpected = hasExpected || (r >= '\u4e00' && r <= '\u9fff')
		case "ja":
			hasExpected = hasExpected || (r >= '\u3040' && r <= '\u30ff') || (r >= '\u4e00' && r <= '\u9fff')
		case "ko":
			hasExpected = hasExpected || (r >= '\uac00' && r <= '\ud7af')
		case "hy":
			hasExpected = hasExpected || (r >= '\u0530' && r <= '\u058f')
		case "ka":
			hasExpected = hasExpected || (r >= '\u10a0' && r <= '\u10ff')
		default:
			hasExpected = hasExpected || (r >= '\u0400' && r <= '\u052f')
		}
		if hasExpected {
			return false
		}
	}
	return true
}

func (b *bot) shadowingConfigured() bool {
	return b != nil &&
		b.openrouter != nil &&
		strings.TrimSpace(b.cfg.OpenRouterSTTModel) != "" &&
		strings.TrimSpace(b.cfg.OpenRouterTTSModel) != ""
}

type shadowingUICopy struct {
	Unavailable string
	Title       string
	Start       string
	Phrase      string
	TextHint    string
	NoGaps      string
	MatchGood   string
	MatchWeak   string
	TipVoice    string
	TipText     string
	Repeat      string
	Target      string
	Heard       string
	Accuracy    string
	Bonus       string
	Next        string
}

var shadowingUICopies = buildShadowingUICopies()

func shadowingCopy(user userState) shadowingUICopy {
	if copy, ok := shadowingUICopies[normalizeInterfaceLanguage(user.InterfaceLanguage)]; ok {
		return copy
	}
	return shadowingUICopies["en"]
}

func shadowingUnavailableText(user userState) string {
	return shadowingCopy(user).Unavailable
}

func shadowingTitle(user userState) string {
	return shadowingCopy(user).Title
}

func shadowingStartMessage(user userState, phrase string) string {
	copy := shadowingCopy(user)
	return copy.Title + "\n\n" + copy.Start
}

func shadowingTextAnswerHint(user userState) string {
	return shadowingCopy(user).TextHint
}

func pronunciationStartMessage(user userState, phrase string) string {
	copy := ui(user)
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		return "🗣 " + copy.Pronunciation + "\n\nФраза:\n" + strings.TrimSpace(phrase) + "\n\nПослушай образец, прочитай фразу и отправь голосовое сообщение. Я оценю именно произношение, ритм и совпадение с этой фразой."
	}
	return "🗣 " + copy.Pronunciation + "\n\nPhrase:\n" + strings.TrimSpace(phrase) + "\n\nListen to the model, read the phrase, and send a voice message. I will score pronunciation, rhythm, and match to this exact phrase."
}

func pronunciationVoiceHint(user userState) string {
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		return "Отправь голосовое сообщение с этой фразой. Текстом произношение не проверяется."
	}
	return "Send a voice message with this phrase. Text cannot check pronunciation."
}

func pronunciationVoicePrompt(user userState) string {
	return pronunciationVoiceHint(user)
}

func shadowingFallbackFeedback(user userState, target string, transcript string, score int, missing string, fromVoice bool) string {
	copy := shadowingCopy(user)
	if missing == "" {
		missing = copy.NoGaps
	}
	match := copy.MatchGood
	if score < 70 {
		match = copy.MatchWeak + " " + missing
	}
	tip := copy.TipVoice
	if !fromVoice {
		tip = copy.TipText
	}
	return match + "\n" + tip + "\n" + copy.Repeat + ": " + target
}

func shadowingResultMessage(user userState, target string, transcript string, score int, feedback string, bonusXP int, assessment *pronunciationAssessment) string {
	copy := shadowingCopy(user)
	result := fmt.Sprintf("%s\n\n%s:\n%s\n\n%s:\n%s\n\n%s: %d/100\n%s: +%d XP",
		copy.Title, copy.Target, target, copy.Heard, strings.TrimSpace(transcript), copy.Accuracy, score, copy.Bonus, bonusXP)
	if assessment != nil {
		result += "\n" + pronunciationShortDetails(user, assessment)
	}
	result += "\n\n" + strings.TrimSpace(feedback)
	return result
}

func pronunciationShortDetails(user userState, assessment *pronunciationAssessment) string {
	if assessment == nil {
		return ""
	}
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		return "\u0421\u0438\u043b\u0430 \u0430\u043a\u0446\u0435\u043d\u0442\u0430: " + itoa(assessment.AccentStrength) + "/100\n" +
			"\u041f\u043b\u0430\u0432\u043d\u043e\u0441\u0442\u044c: " + itoa(assessment.Fluency) + "/100" +
			pronunciationProblemWordsLine(user, assessment)
	}
	return "Accent strength: " + itoa(assessment.AccentStrength) + "/100\n" +
		"Fluency: " + itoa(assessment.Fluency) + "/100" +
		pronunciationProblemWordsLine(user, assessment)
}

func pronunciationProblemWordsLine(user userState, assessment *pronunciationAssessment) string {
	if assessment == nil || len(assessment.ProblemWords) == 0 {
		return ""
	}
	label := "Problem words"
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		label = "\u0421\u043b\u0430\u0431\u044b\u0435 \u0441\u043b\u043e\u0432\u0430"
	}
	var parts []string
	for _, item := range assessment.ProblemWords {
		if strings.TrimSpace(item.Word) == "" {
			continue
		}
		parts = append(parts, strings.TrimSpace(item.Word))
		if len(parts) >= 4 {
			break
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return "\n" + label + ": " + strings.Join(parts, ", ")
}

func shadowingDoneKeyboard(user userState, copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": shadowingCopy(user).Next, "callback_data": "menu_shadowing"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func (b *bot) buildShadowingPhrase(ctx context.Context, user userState, previous string) string {
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	index := shadowingDeckIndex(user)
	if b.openrouter != nil {
		raw, err := b.openrouter.complete(ctx, shadowingPhrasePrompt(language, interfaceLanguage, user.Level, shadowingDeckSpec(index), previous), 0.7, 180)
		if err == nil {
			phrase := sanitizeShadowingPhrase(raw)
			if phrase != "" && !shadowingLooksWrongScript(phrase, language) && !strings.EqualFold(phrase, strings.TrimSpace(previous)) {
				return phrase
			}
		}
	}
	phrase := shadowingFallbackPhrase(language, index)
	if strings.EqualFold(phrase, strings.TrimSpace(previous)) {
		phrase = shadowingFallbackPhrase(language, index+17)
	}
	return phrase
}

func (b *bot) buildPronunciationPhrase(ctx context.Context, user userState, previous string) string {
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	index := int((time.Now().UTC().UnixNano() + int64(user.VoiceCount*97) + int64(user.LessonCount*31)) % int64(shadowingDeckSize))
	if index < 0 {
		index = -index
	}
	if b.openrouter != nil {
		raw, err := b.openrouter.completeWithModel(ctx, b.cfg.OpenRouterPronunciationModel, pronunciationPhrasePrompt(language, interfaceLanguage, user.Level, shadowingDeckSpec(index), previous), 0.8, 160)
		if err == nil {
			phrase := sanitizeShadowingPhrase(cleanModelReply(raw))
			if phrase != "" && !shadowingLooksWrongScript(phrase, language) && !strings.EqualFold(phrase, strings.TrimSpace(previous)) {
				return phrase
			}
		}
	}
	phrase := shadowingFallbackPhrase(language, index+1)
	if strings.EqualFold(phrase, strings.TrimSpace(previous)) {
		phrase = shadowingFallbackPhrase(language, index+17)
	}
	return phrase
}

func (b *bot) startShadowing(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	copy := ui(user)
	if !b.shadowingConfigured() {
		return b.telegram.sendMessageWithCopy(ctx, chatID, shadowingUnavailableText(user), copy)
	}
	if !user.isPremium(time.Now()) {
		return b.telegram.sendInlineMessage(ctx, chatID, copy.Tool.VoicePremiumRequired, premiumMenuShortcutKeyboard(copy))
	}
	voiceLimit := voiceLimitFor(user)
	if user.VoiceToday >= voiceLimit {
		return b.telegram.sendMessage(ctx, chatID, fmt.Sprintf(copy.Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
	}
	if !canUsePractice(user) {
		return b.telegram.sendMessageWithCopy(ctx, chatID, b.limitReachedText(systemUI(user).PracticeKind, user), copy)
	}
	_ = b.telegram.sendChatAction(ctx, chatID, "typing")
	previous, _ := parseShadowingMode(user.Mode)
	phrase := b.buildShadowingPhrase(ctx, user, previous)
	if err := b.store.setMode(user.TelegramID, shadowingMode(phrase)); err != nil {
		return err
	}
	if err := b.telegram.sendInlineMessage(ctx, chatID, shadowingStartMessage(user, phrase), backToMenuKeyboard(copy)); err != nil {
		return err
	}
	b.sendTextPronunciation(ctx, chatID, user, phrase, "shadowing.mp3", phrase, "shadowing target")
	return nil
}

func (b *bot) startPronunciation(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	copy := ui(user)
	if !b.shadowingConfigured() {
		return b.telegram.sendMessageWithCopy(ctx, chatID, shadowingUnavailableText(user), copy)
	}
	if !user.isPremium(time.Now()) {
		return b.telegram.sendInlineMessage(ctx, chatID, copy.Tool.VoicePremiumRequired, premiumMenuShortcutKeyboard(copy))
	}
	voiceLimit := voiceLimitFor(user)
	if user.VoiceToday >= voiceLimit {
		return b.telegram.sendMessage(ctx, chatID, fmt.Sprintf(copy.Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
	}
	if !canUsePractice(user) {
		return b.telegram.sendMessageWithCopy(ctx, chatID, b.limitReachedText(systemUI(user).PracticeKind, user), copy)
	}
	_ = b.telegram.sendChatAction(ctx, chatID, "typing")
	previous, _ := parsePronunciationMode(user.Mode)
	phrase := b.buildPronunciationPhrase(ctx, user, previous)
	if err := b.store.setMode(user.TelegramID, pronunciationPracticeMode(phrase)); err != nil {
		return err
	}
	if err := b.telegram.sendInlineMessage(ctx, chatID, pronunciationStartMessage(user, phrase), backToMenuKeyboard(copy)); err != nil {
		return err
	}
	b.sendTextPronunciation(ctx, chatID, user, phrase, "pronunciation.mp3", phrase, "pronunciation target")
	return nil
}

func premiumMenuShortcutKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": copy.Premium, "callback_data": "menu_premium"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func (b *bot) handleShadowingAnswer(ctx context.Context, chatID int64, user userState, transcript string, fromVoice bool, transcription *audioTranscription) error {
	transcript = strings.TrimSpace(transcript)
	if transcript == "" {
		return b.telegram.sendMessageWithCopy(ctx, chatID, shadowingTextAnswerHint(user), ui(user))
	}
	target, ok := parseShadowingMode(user.Mode)
	if !ok {
		return b.startShadowing(ctx, chatID, user)
	}
	if !canUsePractice(user) {
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.telegram.sendMessageWithCopy(ctx, chatID, b.limitReachedText(systemUI(user).PracticeKind, user), ui(user))
	}
	_ = b.telegram.sendChatAction(ctx, chatID, "typing")
	feedback, score, assessment := b.buildShadowingFeedback(ctx, user, target, transcript, fromVoice, transcription, nil)
	bonusXP, err := b.recordShadowingAttempt(user, target, transcript, score)
	if err != nil {
		return err
	}
	if err := b.telegram.sendInlineMessage(ctx, chatID, shadowingResultMessage(user, target, transcript, score, feedback, bonusXP, assessment), shadowingDoneKeyboard(user, ui(user))); err != nil {
		return err
	}
	if score < 78 {
		b.sendTextPronunciation(ctx, chatID, user, target, "shadowing-repeat.mp3", target, "shadowing repeat")
	}
	return b.maybePromoteLearningLevel(ctx, chatID, user.TelegramID, user.FirstName)
}

func pronunciationDoneKeyboard(user userState, copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": copy.NextWord, "callback_data": "menu_pronunciation"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func (b *bot) handlePronunciationAnswer(ctx context.Context, chatID int64, user userState, transcription audioTranscription) error {
	target, ok := parsePronunciationMode(user.Mode)
	if !ok {
		return b.startPronunciation(ctx, chatID, user)
	}
	transcript := strings.TrimSpace(transcription.Text)
	if transcript == "" {
		return b.telegram.sendMessageWithCopy(ctx, chatID, pronunciationVoiceHint(user), ui(user))
	}
	if !canUsePractice(user) {
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.telegram.sendMessageWithCopy(ctx, chatID, b.limitReachedText(systemUI(user).PracticeKind, user), ui(user))
	}
	assessment := b.buildPronunciationAssessment(ctx, user, target, transcription, pronunciationModeExact)
	historyEntry := fmt.Sprintf("pronunciation target: %s\nspoken: %s\nscore: %d/100", target, transcript, assessment.Score)
	history := trimPracticeHistory(append(user.PracticeHistory, historyEntry), practiceMemoryLimit)
	if err := b.store.incrementPractice(user.TelegramID); err != nil {
		return err
	}
	if err := b.store.savePracticeHistory(user.TelegramID, history); err != nil {
		return err
	}
	if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
		return err
	}
	text := pronunciationAssessmentText(user, &assessment)
	if strings.TrimSpace(text) == "" {
		text = pronunciationVoiceHint(user)
	}
	if err := b.telegram.sendInlineMessage(ctx, chatID, text, pronunciationDoneKeyboard(user, ui(user))); err != nil {
		return err
	}
	if assessment.Score < 78 {
		b.sendTextPronunciation(ctx, chatID, user, shadowingCorrectionAudioText(target, &assessment), "pronunciation-repeat.mp3", target, "pronunciation repeat")
	}
	return b.maybePromoteLearningLevel(ctx, chatID, user.TelegramID, user.FirstName)
}

func (b *bot) buildShadowingFeedback(ctx context.Context, user userState, target string, transcript string, fromVoice bool, transcription *audioTranscription, assessmentOverride *pronunciationAssessment) (string, int, *pronunciationAssessment) {
	if fromVoice {
		if assessmentOverride != nil {
			feedback := strings.TrimSpace(assessmentOverride.Feedback)
			if feedback == "" {
				feedback = shadowingFallbackFeedback(user, target, transcript, assessmentOverride.Score, shadowingMissingHint(target, transcript), fromVoice)
			}
			return feedback, assessmentOverride.Score, assessmentOverride
		}
		if transcription == nil {
			transcription = &audioTranscription{Text: transcript}
		}
		assessment := b.buildPronunciationAssessment(ctx, user, target, *transcription, pronunciationModeExact)
		feedback := strings.TrimSpace(assessment.Feedback)
		if feedback == "" {
			feedback = shadowingFallbackFeedback(user, target, transcript, assessment.Score, shadowingMissingHint(target, transcript), fromVoice)
		}
		return feedback, assessment.Score, &assessment
	}
	score := shadowingScoreForSource(target, transcript, fromVoice)
	missing := shadowingMissingHint(target, transcript)
	if b == nil || b.openrouter == nil {
		return shadowingFallbackFeedback(user, target, transcript, score, missing, fromVoice), score, nil
	}
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	raw, err := b.openrouter.complete(ctx, shadowingFeedbackPrompt(language, interfaceLanguage, user.Level, target, transcript, score, missing, fromVoice), 0.25, 500)
	if err != nil {
		return shadowingFallbackFeedback(user, target, transcript, score, missing, fromVoice), score, nil
	}
	feedback := strings.TrimSpace(cleanModelReply(raw))
	if feedback == "" {
		feedback = shadowingFallbackFeedback(user, target, transcript, score, missing, fromVoice)
	}
	return feedback, score, nil
}

func (b *bot) recordShadowingAttempt(user userState, target string, transcript string, score int) (int, error) {
	bonusXP := shadowingBonusXP(score)
	historyEntry := fmt.Sprintf("shadowing target: %s\nspoken: %s\nscore: %d/100", target, transcript, score)
	history := trimPracticeHistory(append(user.PracticeHistory, historyEntry), practiceMemoryLimit)
	if err := b.store.incrementPractice(user.TelegramID); err != nil {
		return 0, err
	}
	if err := b.store.savePracticeHistory(user.TelegramID, history); err != nil {
		return 0, err
	}
	if err := b.store.addXP(user.TelegramID, bonusXP); err != nil {
		return 0, err
	}
	if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
		return 0, err
	}
	return bonusXP, nil
}

func shadowingBonusXP(score int) int {
	switch {
	case score >= 90:
		return 10
	case score >= 75:
		return 7
	case score >= 55:
		return 4
	default:
		return 2
	}
}

func shadowingCorrectionAudioText(target string, assessment *pronunciationAssessment) string {
	if assessment != nil && strings.TrimSpace(assessment.CorrectedText) != "" {
		return strings.TrimSpace(assessment.CorrectedText)
	}
	return strings.TrimSpace(target)
}

func shadowingUnits(text string) []string {
	var units []string
	var token strings.Builder
	flush := func() {
		if token.Len() == 0 {
			return
		}
		units = append(units, token.String())
		token.Reset()
	}
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '\'' || r == '’' {
			token.WriteRune(r)
			continue
		}
		flush()
	}
	flush()
	if len(units) <= 1 {
		joined := strings.Join(units, "")
		if len([]rune(joined)) > 4 {
			units = units[:0]
			for _, r := range joined {
				units = append(units, string(r))
			}
		}
	}
	return units
}

func shadowingScore(target string, transcript string) int {
	targetUnits := shadowingUnits(target)
	spokenUnits := shadowingUnits(transcript)
	if len(targetUnits) == 0 || len(spokenUnits) == 0 {
		return 0
	}
	longer := math.Max(float64(len(targetUnits)), float64(len(spokenUnits)))
	shorter := math.Min(float64(len(targetUnits)), float64(len(spokenUnits)))
	tokenAcc := 1 - float64(editDistanceStrings(targetUnits, spokenUnits))/longer
	orderAcc := float64(lcsLength(targetUnits, spokenUnits)) / longer
	charAcc := 1 - float64(editDistanceRunes(shadowingCompactRunes(target), shadowingCompactRunes(transcript)))/
		math.Max(1, math.Max(float64(len(shadowingCompactRunes(target))), float64(len(shadowingCompactRunes(transcript)))))
	lengthFit := shorter / longer
	score := int(math.Round((tokenAcc*0.52 + charAcc*0.28 + orderAcc*0.14 + lengthFit*0.06) * 96))
	if strings.Join(targetUnits, " ") == strings.Join(spokenUnits, " ") && score < 96 {
		score = 96
	}
	if lengthFit < 0.55 && score > 52 {
		score = 52
	}
	if score < 0 {
		return 0
	}
	if score > 96 {
		return 96
	}
	return score
}

func shadowingScoreForSource(target string, transcript string, fromVoice bool) int {
	score := shadowingScore(target, transcript)
	if !fromVoice && score > 60 {
		return 60
	}
	return score
}

func shadowingCompactRunes(text string) []rune {
	var out []rune
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			out = append(out, r)
		}
	}
	return out
}

func editDistanceStrings(a []string, b []string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = minInt(cur[j-1]+1, minInt(prev[j]+1, prev[j-1]+cost))
		}
		prev = cur
	}
	return prev[len(b)]
}

func editDistanceRunes(a []rune, b []rune) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = minInt(cur[j-1]+1, minInt(prev[j]+1, prev[j-1]+cost))
		}
		prev = cur
	}
	return prev[len(b)]
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func lcsLength(a []string, b []string) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	prev := make([]int, len(b)+1)
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				cur[j] = prev[j-1] + 1
			} else if prev[j] > cur[j-1] {
				cur[j] = prev[j]
			} else {
				cur[j] = cur[j-1]
			}
		}
		prev = cur
	}
	return prev[len(b)]
}

func shadowingMissingHint(target string, transcript string) string {
	targetUnits := shadowingUnits(target)
	spokenUnits := shadowingUnits(transcript)
	if len(targetUnits) == 0 || len(spokenUnits) == 0 {
		return ""
	}
	seen := map[string]int{}
	for _, unit := range spokenUnits {
		seen[unit]++
	}
	var missing []string
	for _, unit := range targetUnits {
		if seen[unit] > 0 {
			seen[unit]--
			continue
		}
		missing = append(missing, unit)
		if len(missing) >= 5 {
			break
		}
	}
	return strings.Join(missing, ", ")
}

func (api *webAPI) handleShadowingStart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.bot.shadowingConfigured() {
		writeAPIError(w, http.StatusServiceUnavailable, shadowingUnavailableText(user))
		return
	}
	if !user.isPremium(time.Now()) {
		writeAPIError(w, http.StatusPaymentRequired, ui(user).Tool.VoicePremiumRequired)
		return
	}
	voiceLimit := voiceLimitFor(user)
	if user.VoiceToday >= voiceLimit {
		writeAPIError(w, http.StatusTooManyRequests, fmt.Sprintf(ui(user).Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
		return
	}
	if !canUsePractice(user) {
		writeAPIError(w, http.StatusTooManyRequests, api.bot.limitReachedText(systemUI(user).PracticeKind, user))
		return
	}
	var req struct {
		Previous string `json:"previous"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	previous := strings.TrimSpace(req.Previous)
	if previous == "" {
		previous, _ = parseShadowingMode(user.Mode)
	}
	phrase := api.bot.buildShadowingPhrase(r.Context(), user, previous)
	if err := api.bot.store.setMode(user.TelegramID, shadowingMode(phrase)); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"target":      phrase,
		"instruction": shadowingStartMessage(user, phrase),
		"user":        api.userDTO(refreshed),
	})
}

func (api *webAPI) handleShadowingAnswer(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !canUsePractice(user) {
		_ = api.bot.store.setMode(user.TelegramID, "idle")
		writeAPIError(w, http.StatusTooManyRequests, api.bot.limitReachedText(systemUI(user).PracticeKind, user))
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, webMaxVoiceUploadBytes+webMaxMultipartOverhead)
	if err := r.ParseMultipartForm(webMaxVoiceUploadBytes); err != nil {
		writeAPIError(w, http.StatusBadRequest, "Could not read shadowing form.")
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	target, ok := parseShadowingMode(user.Mode)
	if !ok {
		target = sanitizeShadowingPhrase(r.FormValue("target"))
	}
	if target == "" {
		writeAPIError(w, http.StatusConflict, "No active shadowing phrase. Start shadowing again.")
		return
	}
	transcript := strings.TrimSpace(r.FormValue("text"))
	var transcription *audioTranscription
	var assessmentOverride *pronunciationAssessment
	voiceBytes, voiceName, err := readOptionalUploadedFile(r, "voice", webMaxVoiceUploadBytes)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	fromVoice := false
	if len(voiceBytes) > 0 {
		fromVoice = true
		if !api.bot.shadowingConfigured() {
			writeAPIError(w, http.StatusServiceUnavailable, shadowingUnavailableText(user))
			return
		}
		if !user.isPremium(time.Now()) {
			writeAPIError(w, http.StatusPaymentRequired, ui(user).Tool.VoicePremiumRequired)
			return
		}
		voiceLimit := voiceLimitFor(user)
		if user.VoiceToday >= voiceLimit {
			writeAPIError(w, http.StatusTooManyRequests, fmt.Sprintf(ui(user).Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
			return
		}
		format := audioUploadFormat(voiceName, http.DetectContentType(voiceBytes))
		if assessment, detailed, ok := api.bot.assessPronunciationFromAudio(r.Context(), user, voiceBytes, format, target, pronunciationModeExact); ok {
			assessmentOverride = &assessment
			transcription = &detailed
			transcript = detailed.Text
		} else {
			detailed, err := api.bot.transcribeLearningVoice(r.Context(), user, voiceBytes, format, target)
			if err != nil {
				writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(ui(user).Tool.VoiceTranscribeFailed, err.Error()))
				return
			}
			transcription = &detailed
			transcript = detailed.Text
		}
		if err := api.bot.store.incrementVoice(user.TelegramID); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	transcript = strings.TrimSpace(transcript)
	if transcript == "" {
		writeAPIError(w, http.StatusBadRequest, "Record voice or type what you said.")
		return
	}
	feedback, score, assessment := api.bot.buildShadowingFeedback(r.Context(), user, target, transcript, fromVoice, transcription, assessmentOverride)
	bonusXP, err := api.bot.recordShadowingAttempt(user, target, transcript, score)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	promotedTo, err := api.bot.maybePromoteLearningLevelForWeb(r.Context(), user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"target":                target,
		"transcript":            transcript,
		"score":                 score,
		"feedback":              feedback,
		"pronunciation":         pronunciationAssessmentDTO(assessment),
		"correction_audio_text": shadowingCorrectionAudioText(target, assessment),
		"bonus_xp":              bonusXP,
		"promoted_to":           promotedTo,
		"user":                  api.userDTO(refreshed),
	})
}

func (api *webAPI) handlePronunciationStart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		Previous string `json:"previous"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	phrase := api.bot.buildPronunciationPhrase(r.Context(), user, req.Previous)
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"target":      phrase,
		"instruction": "Listen to the model, repeat the exact phrase, then send your recording for pronunciation scoring.",
		"user":        api.userDTO(refreshed),
	})
}

func (api *webAPI) handlePronunciationCheck(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, webMaxVoiceUploadBytes+webMaxMultipartOverhead)
	if err := r.ParseMultipartForm(webMaxVoiceUploadBytes); err != nil {
		writeAPIError(w, http.StatusBadRequest, "Could not read pronunciation form.")
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	target := sanitizeShadowingPhrase(r.FormValue("target"))
	if target == "" {
		writeAPIError(w, http.StatusBadRequest, "target is required")
		return
	}
	voiceBytes, voiceName, err := readOptionalUploadedFile(r, "voice", webMaxVoiceUploadBytes)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(voiceBytes) == 0 {
		writeAPIError(w, http.StatusBadRequest, "Record or upload audio for pronunciation check.")
		return
	}
	if !api.bot.shadowingConfigured() {
		writeAPIError(w, http.StatusServiceUnavailable, shadowingUnavailableText(user))
		return
	}
	if !user.isPremium(time.Now()) {
		writeAPIError(w, http.StatusPaymentRequired, ui(user).Tool.VoicePremiumRequired)
		return
	}
	voiceLimit := voiceLimitFor(user)
	if user.VoiceToday >= voiceLimit {
		writeAPIError(w, http.StatusTooManyRequests, fmt.Sprintf(ui(user).Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
		return
	}
	format := audioUploadFormat(voiceName, http.DetectContentType(voiceBytes))
	var assessment pronunciationAssessment
	var transcription audioTranscription
	if assessed, detailed, ok := api.bot.assessPronunciationFromAudio(r.Context(), user, voiceBytes, format, target, pronunciationModeExact); ok {
		assessment = assessed
		transcription = detailed
	} else {
		detailed, err := api.bot.transcribeLearningVoice(r.Context(), user, voiceBytes, format, target)
		if err != nil {
			writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(ui(user).Tool.VoiceTranscribeFailed, err.Error()))
			return
		}
		transcription = detailed
		assessment = api.bot.buildPronunciationAssessment(r.Context(), user, target, detailed, pronunciationModeExact)
	}
	if err := api.bot.store.incrementVoice(user.TelegramID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"target":                target,
		"transcript":            strings.TrimSpace(transcription.Text),
		"from_voice":            true,
		"pronunciation":         pronunciationAssessmentDTO(&assessment),
		"correction_audio_text": shadowingCorrectionAudioText(target, &assessment),
		"user":                  api.userDTO(refreshed),
	})
}
