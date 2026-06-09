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
	if !strings.HasPrefix(mode, modeShadowingPrefix) {
		return "", false
	}
	raw := strings.TrimPrefix(mode, modeShadowingPrefix)
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
	"en": {"Could you say that a little slower, please?", "I need a few minutes to check the details.", "Let us meet near the main entrance."},
	"ru": {"?????????, ??????????, ??????? ?????????.", "??? ????? ????????? ?????, ????? ????????? ??????.", "??????? ?????????? ? ???????? ?????."},
	"es": {"?Puedes repetirlo un poco m?s despacio, por favor?", "Necesito unos minutos para revisar los detalles.", "Vamos a vernos cerca de la entrada principal."},
	"de": {"K?nnen Sie das bitte etwas langsamer sagen?", "Ich brauche ein paar Minuten, um die Details zu pr?fen.", "Treffen wir uns bitte am Haupteingang."},
	"fr": {"Pouvez-vous r?p?ter un peu plus lentement ?", "J'ai besoin de quelques minutes pour v?rifier les d?tails.", "Retrouvons-nous pr?s de l'entr?e principale."},
	"it": {"Puoi ripeterlo un po' pi? lentamente?", "Mi servono alcuni minuti per controllare i dettagli.", "Incontriamoci vicino all'ingresso principale."},
	"zh": {"???????????", "???????????", "???????????"},
	"ja": {"????????????????", "???????????????????", "???????????????"},
	"ko": {"?? ? ??? ?? ???.", "?? ??? ??? ??? ?? ????.", "?? ???? ???."},
	"tg": {"??????, ???? ??????? ?????? ?????.", "????? ????????? ??????? ???? ?????? ????? ???.", "????? ????? ??????????? ????? ???????."},
	"uz": {"Iltimos, biroz sekinroq takrorlang.", "Tafsilotlarni tekshirish uchun bir necha daqiqa kerak.", "Asosiy kirish yonida uchrashamiz."},
	"tt": {"??????, ????? ???????? ??????????.", "??????????? ??????? ???? ??????? ????? ?????.", "??? ???? ????? ?????? ????????."},
	"hy": {"??????? ?? ?????? ?? ???? ????? ???????", "???????????? ????????? ????? ??? ?? ???? ???? ? ?????", "???? ????????? ??????? ?????? ????"},
	"kk": {"????????, ??? ???????? ??????????.", "???????????? ????????? ??????? ????? ?????.", "??????? ???????????? ??????? ??????????."},
	"ky": {"???????, ??? ?? ???????? ??????????.", "?????????? ?????????? ??? ???? ????? ?????.", "????? ???? ???????? ??????? ????????."},
	"ka": {"??????, ???? ???? ?????????.", "????????? ????????????? ????????? ???? ????????.", "?????? ???????????? ????????."},
	"uk": {"?????????, ???? ?????, ????? ??????????.", "???? ???????? ?????? ??????, ??? ?????????? ??????.", "???????????? ???? ????????? ?????."},
	"pl": {"Mo?esz powt?rzy? troch? wolniej?", "Potrzebuj? kilku minut, ?eby sprawdzi? szczeg??y.", "Spotkajmy si? przy g??wnym wej?ciu."},
	"ro": {"Po?i repeta pu?in mai ?ncet, te rog?", "Am nevoie de c?teva minute ca s? verific detaliile.", "S? ne ?nt?lnim l?ng? intrarea principal?."},
	"pt": {"Voc? pode repetir um pouco mais devagar?", "Preciso de alguns minutos para verificar os detalhes.", "Vamos nos encontrar perto da entrada principal."},
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

var shadowingUICopies = map[string]shadowingUICopy{
	"en": {"Shadowing is not configured yet: speech recognition and text-to-speech models are required.", "Shadowing: repeat like the speaker", "Listen first. Then send a voice message and copy the rhythm, word order, and sound.", "Phrase", "Send a voice message here. Text can only be used as a limited fallback.", "no obvious gaps", "Match: the recognizer heard a close repeat.", "Match: this needs another repeat.", "Tip: repeat the whole rhythm without long pauses between words.", "Tip: text cannot check pronunciation; send voice for a real score.", "Repeat", "Target", "I heard", "Accuracy", "Bonus", "Next shadowing"},
	"ru": {"Shadowing ???? ?? ????????: ????? ?????? ????????????? ? ???????.", "Shadowing: ??????? ??? ??????", "??????? ????????? ?????. ????? ??????? ????????? ? ???????? ????, ??????? ???? ? ????????.", "?????", "????????? ?????????. ????? ????? ?????? ???????????? ???????? ????????.", "??? ????? ?????????", "??????????: ????????????? ???????? ??????? ??????.", "??????????: ????? ????????? ??? ???.", "?????: ??????? ??????? ???? ??? ??????? ???? ????? ???????.", "?????: ?? ?????? ???????????? ?? ?????????; ??? ???????? ?????? ????? ?????.", "???????", "????", "? ???????", "????????", "?????", "????????? shadowing"},
	"es": {"Shadowing a?n no est? configurado: faltan modelos de voz y lectura.", "Shadowing: repite como el hablante", "Escucha primero. Luego env?a voz y copia ritmo, orden y sonido.", "Frase", "Env?a voz. El texto solo sirve como revisi?n limitada.", "sin huecos claros", "Coincidencia: el reconocimiento oy? una repetici?n cercana.", "Coincidencia: hay que repetirlo otra vez.", "Consejo: repite el ritmo completo sin pausas largas.", "Consejo: el texto no eval?a pronunciaci?n; env?a voz.", "Repite", "Objetivo", "Escuch?", "Precisi?n", "Bono", "Siguiente shadowing"},
	"de": {"Shadowing ist noch nicht eingerichtet: Spracherkennung und TTS fehlen.", "Shadowing: wie der Sprecher wiederholen", "H?re zuerst zu. Sende dann eine Sprachnachricht und kopiere Rhythmus, Wortfolge und Klang.", "Satz", "Sende eine Sprachnachricht. Text ist nur eine begrenzte Notl?sung.", "keine klaren L?cken", "Treffer: die Erkennung h?rte eine nahe Wiederholung.", "Treffer: bitte noch einmal wiederholen.", "Tipp: sprich den ganzen Rhythmus ohne lange Pausen.", "Tipp: Text pr?ft keine Aussprache; sende Stimme.", "Wiederhole", "Ziel", "Ich h?rte", "Genauigkeit", "Bonus", "N?chstes Shadowing"},
	"fr": {"Shadowing n'est pas encore configur? : reconnaissance vocale et synth?se requises.", "Shadowing : r?p?te comme le locuteur", "?coute d'abord. Envoie ensuite un vocal et copie le rythme, l'ordre et le son.", "Phrase", "Envoie un vocal. Le texte n'est qu'un d?pannage limit?.", "pas de manque ?vident", "Correspondance : la reconnaissance a entendu une r?p?tition proche.", "Correspondance : il faut r?p?ter encore.", "Conseil : garde le rythme complet sans longues pauses.", "Conseil : le texte ne v?rifie pas la prononciation ; envoie un vocal.", "R?p?te", "Cible", "J'ai entendu", "Pr?cision", "Bonus", "Shadowing suivant"},
	"it": {"Shadowing non ? ancora configurato: servono riconoscimento vocale e sintesi.", "Shadowing: ripeti come lo speaker", "Prima ascolta. Poi invia un vocale e copia ritmo, ordine delle parole e suono.", "Frase", "Invia un vocale. Il testo ? solo un controllo limitato.", "nessuna lacuna chiara", "Corrispondenza: il riconoscimento ha sentito una ripetizione vicina.", "Corrispondenza: serve un altro tentativo.", "Consiglio: ripeti il ritmo intero senza pause lunghe.", "Consiglio: il testo non controlla la pronuncia; invia voce.", "Ripeti", "Obiettivo", "Ho sentito", "Precisione", "Bonus", "Shadowing successivo"},
	"zh": {"Shadowing ?????????????????", "Shadowing?????????", "???????????????????????", "??", "???????????????????", "??????", "??????????????", "????????????", "????????????????????", "??????????????????", "??", "??", "???", "???", "??", "??? Shadowing"},
	"ja": {"Shadowing ?????????????????????????", "Shadowing???????????", "????????????????????????????????", "????", "?????????????????????????????", "????????????", "????????????????", "???????????????", "????????????????????????????", "???????????????????????????????", "????", "??", "??????", "???", "????", "?? Shadowing"},
	"ko": {"Shadowing? ?? ???? ????: ?? ??? TTS ??? ????.", "Shadowing: ???? ?? ???", "?? ?? ???. ??? ???? ??, ??, ??? ?? ? ???.", "??", "??? ?? ???. ???? ???? ?? ??? ???.", "??? ?? ??", "??: ?? ??? ??? ? ????.", "??: ? ? ? ??? ????.", "?: ?? ??? ?? ?? ?? ?? ???? ?? ???.", "?: ????? ??? ??? ? ???. ??? ?? ???.", "??", "??", "?? ??", "???", "???", "?? Shadowing"},
	"tg": {"Shadowing ???? ?????? ?????????: ???????? ??????? ???? ?? ?????? ????????.", "Shadowing: ????? ?????? ?????? ???", "????? ??? ?????. ???? ???? ???????? ?? ????, ??????? ???????? ?? ?????? ????? ?????.", "?????", "???? ????????. ???? ????? ??????? ??????? ??????? ???.", "???????? ?????? ????", "?????????: ?????? ??????? ???????? ?????.", "?????????: ??? ?? ??? ?????? ????? ???.", "????????: ?????? ??????, ?? ??????????? ????? ?????? ?????.", "????????: ???? ?????????? ??????????; ???? ????????.", "??????", "?????", "??? ???????", "??????", "?????", "Shadowing-? ???????"},
	"uz": {"Shadowing hali sozlanmagan: nutqni tanish va ovozlantirish modellari kerak.", "Shadowing: spiker kabi takrorlang", "Avval tinglang. Keyin ovoz yuboring va ritm, so'z tartibi va tovushni ko'chiring.", "Ibora", "Ovoz yuboring. Matn faqat cheklangan zaxira tekshiruv.", "aniq bo'shliq yo'q", "Moslik: tanib olish yaqin takrorni eshitdi.", "Moslik: yana bir bor takrorlash kerak.", "Maslahat: so'zlar orasida uzun pauzasiz butun ritmni ayting.", "Maslahat: matn talaffuzni tekshirmaydi; ovoz yuboring.", "Takrorlang", "Maqsad", "Eshitildi", "Aniqlik", "Bonus", "Keyingi shadowing"},
	"tt": {"Shadowing ??? ???????????: ??????? ???? ??? ???????? ?????????? ?????.", "Shadowing: ?????? ????? ???????", "????? ?????. ?????? ????? ????? ??? ??????, ??? ????????, ????????? ???????.", "?????", "????? ?????????. ????? ???? ??? ????? ????? ???????.", "???? ????? ???? ??", "???? ????: ???? ???? ?????????? ??????.", "???? ????: ????? ??? ???????? ?????.", "?????: ?????? ???????? ???? ???????? ????? ?????? ???.", "?????: ????? ????????? ????????; ????? ?????.", "???????", "??????", "??? ???????", "????????", "?????", "?????? shadowing"},
	"hy": {"Shadowing-? ??? ??????????? ??? ???? ?? ????? ???????? ? ??????????? ?????????", "Shadowing? ?????? ?????? ??????", "??? ????? ???? ???????? ???? ? ?????? ?????, ?????? ????? ? ?????????", "???????????????", "???????? ????? ?????? ????? ?????????? ??????????? ???????? ??", "??????? ????????? ???", "??????????? ????????? ???? ??? ????????????", "??????????? ???? ? ?? ??? ????? ???????", "????????? ????? ?????? ????? ????? ????? ??????????", "????????? ??????? ??????????????? ?? ?????????, ???????? ?????", "??????", "??????", "?? ?????", "????????????", "??????", "?????? shadowing"},
	"kk": {"Shadowing ??? ???????????: ???????? ???? ???? ???????? ?????????? ?????.", "Shadowing: ?????? ?????? ???????", "??????? ????????. ????? ????? ???????, ?????, ??? ???? ???? ??????? ??????????.", "?????", "????? ?????????. ????? ??? ???????? ??????? ???????.", "????? ??????? ???", "?????????: ???? ????? ?????????? ??????.", "?????????: ???? ??? ??? ???????? ?????.", "?????: ??? ???????? ???? ????????, ????? ???????? ???????.", "?????: ????? ????????? ???????????; ????? ?????????.", "???????", "??????", "??? ???????", "??????", "?????", "?????? shadowing"},
	"ky": {"Shadowing ???????? ????????? ????: ?? ?????? ???? ?? ??????? ????????? ?????.", "Shadowing: ????????? ???????", "???????? ??. ????? ????? ?? ???????, ??????, ??? ???????? ???? ???? ???????.", "?????", "?? ?????????. ????? ????????? ??????? ???????? ????.", "???? ???????? ???", "??? ?????: ?????? ????? ?????????? ????.", "??? ?????: ???? ??? ???? ???????? ?????.", "?????: ????????? ????????? ???? ???????? ??? ?????? ???.", "?????: ????? ????????? ??????????; ?? ?????????.", "???????", "??????", "??? ?????", "??????", "?????", "??????? shadowing"},
	"ka": {"Shadowing ??? ?? ???? ?????????: ??????? ??????????? ???????? ?? ??????? ??????.", "Shadowing: ???????? ?????? ????????", "??? ?????????. ?????? ????????? ??? ?? ????????? ?????, ???????? ???? ?? ?????????.", "?????", "????????? ???. ?????? ?????? ????????? ?????????? ??????????.", "?????? ?????????? ?? ????", "?????????: ????????? ???? ????????? ????????.", "?????????: ????? ?????? ?????????? ??????.", "?????: ????? ????? ??????, ???????? ????? ?????? ???????? ??????.", "?????: ?????? ?????????? ??? ????????; ????????? ???.", "????????", "??????", "?? ??????", "???????", "??????", "??????? shadowing"},
	"uk": {"Shadowing ?? ?? ???????????: ???????? ?????? ????????????? ? ?????????.", "Shadowing: ??????? ?? ??????", "???????? ????????? ?????. ????? ??????? ???????? ? ??????? ????, ??????? ???? ? ????????.", "?????", "????????? ????????. ????? ??? ???? ???????? ??????? ?????????.", "????? ????? ?????????", "????: ????????????? ?????? ???????? ??????.", "????: ???????? ????????? ?? ???.", "??????: ??????? ???????? ???? ??? ?????? ???? ??? ???????.", "??????: ????? ?? ????????? ??????; ????????? ?????.", "???????", "????", "? ?????", "????????", "?????", "????????? shadowing"},
	"pl": {"Shadowing nie jest jeszcze skonfigurowany: potrzebne s? modele rozpoznawania i syntezy mowy.", "Shadowing: powt?rz jak lektor", "Najpierw pos?uchaj. Potem wy?lij g?os i skopiuj rytm, szyk s??w oraz brzmienie.", "Fraza", "Wy?lij g?os. Tekst to tylko ograniczona kontrola awaryjna.", "brak wyra?nych luk", "Zgodno??: rozpoznawanie us?ysza?o bliskie powt?rzenie.", "Zgodno??: trzeba powt?rzy? jeszcze raz.", "Wskaz?wka: powt?rz ca?y rytm bez d?ugich pauz.", "Wskaz?wka: tekst nie sprawdza wymowy; wy?lij g?os.", "Powt?rz", "Cel", "Us?ysza?em", "Dok?adno??", "Bonus", "Nast?pny shadowing"},
	"ro": {"Shadowing nu este ?nc? setat: sunt necesare modele de recunoa?tere ?i redare vocal?.", "Shadowing: repet? ca vorbitorul", "Ascult? mai ?nt?i. Apoi trimite voce ?i copiaz? ritmul, ordinea cuvintelor ?i sunetul.", "Fraz?", "Trimite voce. Textul este doar o verificare limitat? de rezerv?.", "f?r? lipsuri evidente", "Potrivire: recunoa?terea a auzit o repetare apropiat?.", "Potrivire: trebuie repetat ?nc? o dat?.", "Sfat: repet? ritmul complet f?r? pauze lungi.", "Sfat: textul nu verific? pronun?ia; trimite voce.", "Repet?", "?int?", "Am auzit", "Acurate?e", "Bonus", "Urm?torul shadowing"},
	"pt": {"Shadowing ainda n?o est? configurado: s?o necess?rios modelos de reconhecimento e voz.", "Shadowing: repita como o locutor", "Ou?a primeiro. Depois envie voz e copie ritmo, ordem das palavras e som.", "Frase", "Envie voz. Texto ? s? uma verifica??o limitada de reserva.", "sem falhas claras", "Correspond?ncia: o reconhecimento ouviu uma repeti??o pr?xima.", "Correspond?ncia: ? preciso repetir mais uma vez.", "Dica: repita o ritmo inteiro sem pausas longas.", "Dica: texto n?o verifica pron?ncia; envie voz.", "Repita", "Alvo", "Ouvi", "Precis?o", "B?nus", "Pr?ximo shadowing"},
}

func shadowingCopy(user userState) shadowingUICopy {
	if copy, ok := shadowingUICopies[normalizeInterfaceLanguage(user.InterfaceLanguage)]; ok {
		if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
			copy.Title = "\u0410\u0443\u0434\u0438\u0440\u043e\u0432\u0430\u043d\u0438\u0435"
			copy.Start = "\u0421\u043d\u0430\u0447\u0430\u043b\u0430 \u043f\u0440\u043e\u0441\u043b\u0443\u0448\u0430\u0439 \u0444\u0440\u0430\u0437\u0443. \u041f\u043e\u0442\u043e\u043c \u043e\u0442\u043f\u0440\u0430\u0432\u044c \u0433\u043e\u043b\u043e\u0441 \u0438 \u043f\u043e\u0432\u0442\u043e\u0440\u0438 \u0440\u0438\u0442\u043c, \u0441\u043b\u043e\u0432\u0430 \u0438 \u0437\u0432\u0443\u0447\u0430\u043d\u0438\u0435."
			copy.Phrase = "\u0424\u0440\u0430\u0437\u0430"
			copy.TextHint = "\u0414\u043b\u044f \u0440\u0435\u0430\u043b\u044c\u043d\u043e\u0439 \u043e\u0446\u0435\u043d\u043a\u0438 \u043f\u0440\u043e\u0438\u0437\u043d\u043e\u0448\u0435\u043d\u0438\u044f \u043d\u0443\u0436\u0435\u043d \u0433\u043e\u043b\u043e\u0441."
			copy.Next = "\u0421\u043b\u0435\u0434\u0443\u044e\u0449\u0435\u0435 \u0430\u0443\u0434\u0438\u0440\u043e\u0432\u0430\u043d\u0438\u0435"
			return copy
		}
		copy.Unavailable = strings.ReplaceAll(copy.Unavailable, "Shadowing", "Listening")
		copy.Title = strings.ReplaceAll(copy.Title, "Shadowing", "Listening")
		copy.Next = strings.ReplaceAll(copy.Next, "Shadowing", "Listening")
		return copy
	}
	copy := shadowingUICopies["en"]
	copy.Unavailable = strings.ReplaceAll(copy.Unavailable, "Shadowing", "Listening")
	copy.Title = "Listening: repeat what you hear"
	copy.Next = "Next listening"
	return copy
}

func shadowingUnavailableText(user userState) string {
	return shadowingCopy(user).Unavailable
}

func shadowingTitle(user userState) string {
	return shadowingCopy(user).Title
}

func shadowingStartMessage(user userState, phrase string) string {
	copy := shadowingCopy(user)
	return copy.Title + "\n\n" + copy.Start + "\n\n" + copy.Phrase + ":\n" + phrase
}

func shadowingTextAnswerHint(user userState) string {
	return shadowingCopy(user).TextHint
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
	if !fromVoice {
		return b.telegram.sendMessageWithCopy(ctx, chatID, shadowingTextAnswerHint(user), ui(user))
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
