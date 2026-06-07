package main

import "strings"

type learningLanguage struct {
	Code              string
	Name              string
	NativeName        string
	Adjective         string
	TeacherNoun       string
	PracticeDirection string
	InterfaceName     string
	ExplanationName   string
}

var russianLanguage = learningLanguage{Code: "ru", Name: "русский", NativeName: "Russian", Adjective: "русский", TeacherNoun: "Russian teacher", PracticeDirection: "по-русски", InterfaceName: "Русский", ExplanationName: "Russian"}

var learningLanguages = []learningLanguage{
	{Code: "en", Name: "английский", NativeName: "English", Adjective: "английский", TeacherNoun: "English teacher", PracticeDirection: "по-английски", InterfaceName: "English", ExplanationName: "English"},
	russianLanguage,
	{Code: "es", Name: "испанский", NativeName: "Spanish", Adjective: "испанский", TeacherNoun: "Spanish teacher", PracticeDirection: "по-испански", InterfaceName: "Español", ExplanationName: "Spanish"},
	{Code: "de", Name: "немецкий", NativeName: "German", Adjective: "немецкий", TeacherNoun: "German teacher", PracticeDirection: "по-немецки", InterfaceName: "Deutsch", ExplanationName: "German"},
	{Code: "fr", Name: "французский", NativeName: "French", Adjective: "французский", TeacherNoun: "French teacher", PracticeDirection: "по-французски", InterfaceName: "Français", ExplanationName: "French"},
	{Code: "it", Name: "итальянский", NativeName: "Italian", Adjective: "итальянский", TeacherNoun: "Italian teacher", PracticeDirection: "по-итальянски", InterfaceName: "Italiano", ExplanationName: "Italian"},
	{Code: "zh", Name: "китайский", NativeName: "Chinese", Adjective: "китайский", TeacherNoun: "Chinese teacher", PracticeDirection: "по-китайски", InterfaceName: "中文", ExplanationName: "Chinese"},
	{Code: "ja", Name: "японский", NativeName: "Japanese", Adjective: "японский", TeacherNoun: "Japanese teacher", PracticeDirection: "по-японски", InterfaceName: "日本語", ExplanationName: "Japanese"},
	{Code: "ko", Name: "корейский", NativeName: "Korean", Adjective: "корейский", TeacherNoun: "Korean teacher", PracticeDirection: "по-корейски", InterfaceName: "한국어", ExplanationName: "Korean"},
	{Code: "tg", Name: "таджикский", NativeName: "Tajik", Adjective: "таджикский", TeacherNoun: "Tajik teacher", PracticeDirection: "по-таджикски", InterfaceName: "Тоҷикӣ", ExplanationName: "Tajik"},
	{Code: "uz", Name: "узбекский", NativeName: "Uzbek", Adjective: "узбекский", TeacherNoun: "Uzbek teacher", PracticeDirection: "по-узбекски", InterfaceName: "O‘zbekcha", ExplanationName: "Uzbek"},
	{Code: "tt", Name: "татарский", NativeName: "Tatar", Adjective: "татарский", TeacherNoun: "Tatar teacher", PracticeDirection: "по-татарски", InterfaceName: "Татарча", ExplanationName: "Tatar"},
	{Code: "hy", Name: "армянский", NativeName: "Armenian", Adjective: "армянский", TeacherNoun: "Armenian teacher", PracticeDirection: "по-армянски", InterfaceName: "Հայերեն", ExplanationName: "Armenian"},
	{Code: "kk", Name: "казахский", NativeName: "Kazakh", Adjective: "казахский", TeacherNoun: "Kazakh teacher", PracticeDirection: "по-казахски", InterfaceName: "Қазақша", ExplanationName: "Kazakh"},
	{Code: "ky", Name: "кыргызский", NativeName: "Kyrgyz", Adjective: "кыргызский", TeacherNoun: "Kyrgyz teacher", PracticeDirection: "по-кыргызски", InterfaceName: "Кыргызча", ExplanationName: "Kyrgyz"},
	{Code: "ka", Name: "грузинский", NativeName: "Georgian", Adjective: "грузинский", TeacherNoun: "Georgian teacher", PracticeDirection: "по-грузински", InterfaceName: "ქართული", ExplanationName: "Georgian"},
	{Code: "uk", Name: "украинский", NativeName: "Ukrainian", Adjective: "украинский", TeacherNoun: "Ukrainian teacher", PracticeDirection: "по-украински", InterfaceName: "Українська", ExplanationName: "Ukrainian"},
	{Code: "pl", Name: "польский", NativeName: "Polish", Adjective: "польский", TeacherNoun: "Polish teacher", PracticeDirection: "по-польски", InterfaceName: "Polski", ExplanationName: "Polish"},
	{Code: "ro", Name: "румынский", NativeName: "Romanian", Adjective: "румынский", TeacherNoun: "Romanian teacher", PracticeDirection: "по-румынски", InterfaceName: "Română", ExplanationName: "Romanian"},
	{Code: "pt", Name: "португальский", NativeName: "Portuguese", Adjective: "португальский", TeacherNoun: "Portuguese teacher", PracticeDirection: "по-португальски", InterfaceName: "Português", ExplanationName: "Portuguese"},
	{Code: "ar", Name: "арабский", NativeName: "Arabic", Adjective: "арабский", TeacherNoun: "Arabic teacher", PracticeDirection: "по-арабски", InterfaceName: "العربية", ExplanationName: "Arabic"},
	{Code: "bn", Name: "бенгальский", NativeName: "Bengali", Adjective: "бенгальский", TeacherNoun: "Bengali teacher", PracticeDirection: "по-бенгальски", InterfaceName: "বাংলা", ExplanationName: "Bengali"},
	{Code: "cs", Name: "чешский", NativeName: "Czech", Adjective: "чешский", TeacherNoun: "Czech teacher", PracticeDirection: "по-чешски", InterfaceName: "Čeština", ExplanationName: "Czech"},
	{Code: "el", Name: "греческий", NativeName: "Greek", Adjective: "греческий", TeacherNoun: "Greek teacher", PracticeDirection: "по-гречески", InterfaceName: "Ελληνικά", ExplanationName: "Greek"},
	{Code: "hi", Name: "хинди", NativeName: "Hindi", Adjective: "хинди", TeacherNoun: "Hindi teacher", PracticeDirection: "на хинди", InterfaceName: "हिंदी", ExplanationName: "Hindi"},
	{Code: "hu", Name: "венгерский", NativeName: "Hungarian", Adjective: "венгерский", TeacherNoun: "Hungarian teacher", PracticeDirection: "по-венгерски", InterfaceName: "Magyar", ExplanationName: "Hungarian"},
	{Code: "id", Name: "индонезийский", NativeName: "Indonesian", Adjective: "индонезийский", TeacherNoun: "Indonesian teacher", PracticeDirection: "по-индонезийски", InterfaceName: "Bahasa Indonesia", ExplanationName: "Indonesian"},
	{Code: "nl", Name: "нидерландский", NativeName: "Dutch", Adjective: "нидерландский", TeacherNoun: "Dutch teacher", PracticeDirection: "по-нидерландски", InterfaceName: "Nederlands", ExplanationName: "Dutch"},
	{Code: "sv", Name: "шведский", NativeName: "Swedish", Adjective: "шведский", TeacherNoun: "Swedish teacher", PracticeDirection: "по-шведски", InterfaceName: "svenska", ExplanationName: "Swedish"},
	{Code: "ta", Name: "тамильский", NativeName: "Tamil", Adjective: "тамильский", TeacherNoun: "Tamil teacher", PracticeDirection: "по-тамильски", InterfaceName: "தமிழ்", ExplanationName: "Tamil"},
	{Code: "te", Name: "телугу", NativeName: "Telugu", Adjective: "телугу", TeacherNoun: "Telugu teacher", PracticeDirection: "на телугу", InterfaceName: "తెలుగు", ExplanationName: "Telugu"},
	{Code: "th", Name: "тайский", NativeName: "Thai", Adjective: "тайский", TeacherNoun: "Thai teacher", PracticeDirection: "по-тайски", InterfaceName: "ภาษาไทย", ExplanationName: "Thai"},
	{Code: "tl", Name: "тагалог", NativeName: "Tagalog", Adjective: "тагалог", TeacherNoun: "Tagalog teacher", PracticeDirection: "на тагалоге", InterfaceName: "Tagalog", ExplanationName: "Tagalog"},
	{Code: "tr", Name: "турецкий", NativeName: "Turkish", Adjective: "турецкий", TeacherNoun: "Turkish teacher", PracticeDirection: "по-турецки", InterfaceName: "Türkçe", ExplanationName: "Turkish"},
	{Code: "vi", Name: "вьетнамский", NativeName: "Vietnamese", Adjective: "вьетнамский", TeacherNoun: "Vietnamese teacher", PracticeDirection: "по-вьетнамски", InterfaceName: "Tiếng Việt", ExplanationName: "Vietnamese"},
}

var extraInterfaceLanguages = []learningLanguage{}

func normalizeLearningLanguage(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "rus" {
		code = "ru"
	}
	if code == "zho" || code == "chi" || code == "cmn" {
		code = "zh"
	}
	if code == "jpn" {
		code = "ja"
	}
	if code == "kor" {
		code = "ko"
	}
	for _, language := range learningLanguages {
		if code == language.Code {
			return language.Code
		}
	}
	return "en"
}

func learningLanguageByCode(code string) learningLanguage {
	code = normalizeLearningLanguage(code)
	for _, language := range learningLanguages {
		if language.Code == code {
			return language
		}
	}
	return learningLanguages[0]
}

func userLearningLanguage(user userState) learningLanguage {
	return learningLanguageByCode(user.LearningLanguage)
}

func normalizeInterfaceLanguage(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "rus" {
		code = "ru"
	}
	if code == "rom" {
		code = "ro"
	}
	if code == "por" {
		code = "pt"
	}
	if code == "kir" {
		code = "ky"
	}
	if code == "geo" || code == "kat" {
		code = "ka"
	}
	if code == "zho" || code == "chi" || code == "cmn" {
		code = "zh"
	}
	if code == "jpn" {
		code = "ja"
	}
	if code == "kor" {
		code = "ko"
	}
	if code == "gr" || code == "ell" || code == "gre" {
		code = "el"
	}
	if code == "cz" || code == "cze" || code == "ces" {
		code = "cs"
	}
	if code == "in" || code == "ind" {
		code = "id"
	}
	if code == "se" || code == "swe" {
		code = "sv"
	}
	if code == "fil" || code == "tgl" {
		code = "tl"
	}
	if code == "vn" || code == "vie" {
		code = "vi"
	}
	for _, language := range interfaceLanguages() {
		if code == language.Code {
			return language.Code
		}
	}
	return "ru"
}

func userInterfaceLanguage(user userState) learningLanguage {
	return interfaceLanguageByCode(user.InterfaceLanguage)
}

func interfaceLanguages() []learningLanguage {
	languages := []learningLanguage{russianLanguage}
	for _, language := range learningLanguages {
		if language.Code != russianLanguage.Code {
			languages = append(languages, language)
		}
	}
	languages = append(languages, extraInterfaceLanguages...)
	return languages
}

func interfaceLanguageByCode(code string) learningLanguage {
	code = normalizeInterfaceLanguage(code)
	for _, language := range interfaceLanguages() {
		if language.Code == code {
			return language
		}
	}
	return interfaceLanguages()[0]
}

func interfaceLanguageSelectionText() string {
	var builder strings.Builder
	builder.WriteString("Choose bot language / Выбери язык бота")
	for _, language := range interfaceLanguages() {
		copy := ui(userState{InterfaceLanguage: language.Code, InterfaceSelected: true})
		builder.WriteString("\n")
		builder.WriteString(language.InterfaceName)
		builder.WriteString(": ")
		builder.WriteString(copy.ChooseBotLang)
	}
	return builder.String()
}

func parseInterfaceLanguageCallback(data string) (string, bool) {
	parts := strings.Split(data, "|")
	if len(parts) != 2 || parts[0] != "ui" {
		return "", false
	}
	code := normalizeInterfaceLanguage(parts[1])
	for _, language := range interfaceLanguages() {
		if code == language.Code {
			return code, true
		}
	}
	return "", false
}

func leaderboardLanguages(user userState) []string {
	seen := map[string]bool{}
	add := func(code string) {
		seen[normalizeLearningLanguage(code)] = true
	}

	add(user.LearningLanguage)
	for _, word := range user.LearnedWords {
		add(word.Language)
	}
	for _, mistake := range user.Mistakes {
		add(mistake.Language)
	}

	languages := make([]string, 0, len(seen))
	for _, language := range learningLanguages {
		if seen[language.Code] {
			languages = append(languages, language.NativeName)
		}
	}
	return languages
}

func languageSelectionText() string {
	return "Какой язык хочешь изучать?\n\nВыбери один вариант. Потом язык можно будет поменять в меню."
}

func parseLanguageCallback(data string) (string, bool) {
	parts := strings.Split(data, "|")
	if len(parts) != 2 || parts[0] != "lang" {
		return "", false
	}
	code := strings.ToLower(strings.TrimSpace(parts[1]))
	for _, language := range learningLanguages {
		if code == language.Code {
			return code, true
		}
	}
	return "", false
}
