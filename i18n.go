package main

import (
	"fmt"
	"strings"
)

type uiCopy struct {
	MenuButton       string
	StopButton       string
	Back             string
	BackMenu         string
	MainMenuTitle    string
	MainMenuBody     string
	Learning         string
	Words            string
	Stats            string
	Settings         string
	Tools            string
	NewLesson        string
	Practice         string
	AITutor          string
	Shadowing        string
	LevelTest        string
	LearnWords       string
	WordGame         string
	Spelling         string
	Vocabulary       string
	Phrasebook       string
	Mistakes         string
	Progress         string
	Leaders          string
	Limits           string
	BotLanguage      string
	LearningLanguage string
	Notifications    string
	Premium          string
	Referral         string
	ChooseBotLang    string
	ChooseLearnLang  string
	ChooseTimezone   string
	TimezoneHint     string
	BotLangSet       string
	LearnLangSet     string
	UnknownButton    string
	Stopped          string
	WordQuestion     string
	ChooseAnswer     string
	Correct          string
	TryAgain         string
	NextWord         string
	NextPage         string
	AlreadyLearned   string
	AddedToVocab     string
	TotalLearned     string
	WriteWord        string
	Hint             string
	GoodSpelling     string
	Tool             toolUICopy
}

type toolUICopy struct {
	VoiceToText            string
	ImageTranslate         string
	Translator             string
	GPTAgent               string
	VoicePrompt            string
	ImagePrompt            string
	TranslatorPrompt       string
	VoiceModePrompt        string
	ImageModePrompt        string
	TranslatorModePrompt   string
	ImageOpenToolsPrompt   string
	ImageSinglePhotoPrompt string
	TranscriptLabel        string
	TranslationLabel       string
	SourceLanguage         string
	TargetLanguage         string
	AutoDetect             string
	WebApp                 string
	VoiceDiscussPrompt     string
	ImageDiscussPrompt     string
	VoicePremiumRequired   string
	ImagePremiumRequired   string
	VoiceLimitReached      string
	VoiceTooLong           string
	VoiceDownloadFailed    string
	VoiceTranscribeFailed  string
	VoiceTranslationFailed string
	ImageFileMissing       string
	ImageDownloadFailed    string
	ImageReadFailed        string
}

var uiCopies = map[string]uiCopy{
	"ru": {
		MenuButton: "📋 Меню", StopButton: "⏹ Стоп", Back: "◀ Назад", BackMenu: "◀ В меню",
		MainMenuTitle: "Главное меню", MainMenuBody: "Выбери действие:",
		Learning: "Учиться", Words: "Слова", Stats: "Прогресс", Settings: "Настройки", Tools: "Инструменты",
		NewLesson: "Новый урок", Practice: "Практика", AITutor: "AI Репетитор", LevelTest: "Определить уровень", LearnWords: "Учить слова", WordGame: "Повторяй-ка", Spelling: "Правописание", Vocabulary: "Словарик", Phrasebook: "Разговорник", Mistakes: "Словарь ошибок", Progress: "Мой прогресс", Leaders: "Лидеры", Limits: "Мои лимиты", BotLanguage: "Язык бота", LearningLanguage: "Язык изучения", Notifications: "Уведомления", Premium: "Premium",
		ChooseBotLang: "Выбери язык бота", ChooseLearnLang: "Какой язык хочешь изучать?", ChooseTimezone: "Выбери часовой пояс", TimezoneHint: "Я буду присылать ежедневное напоминание примерно в 19:00 по твоему времени.",
		BotLangSet: "Готово! Язык бота: %s.", LearnLangSet: "Готово! Теперь изучаем: %s.", UnknownButton: "Не понял кнопку. Нажми 📋 Меню.", Stopped: "Остановлено. Нажми 📋 Меню чтобы выбрать следующее действие.",
		WordQuestion: "Как будет %s:", ChooseAnswer: "Выбери правильный вариант:", Correct: "Правильно!", TryAgain: "Пока не то. Попробуй ещё раз:", NextWord: "Следующее слово", NextPage: "Вперёд ▶", AlreadyLearned: "Это слово уже было в твоем словаре.", AddedToVocab: "Слово добавлено в изученный словарь.", TotalLearned: "Всего изучено слов: %d", WriteWord: "Напиши %s:", Hint: "Подсказка", GoodSpelling: "Отлично, написано верно!",
	},
	"en": englishUICopy(),
	"es": {MenuButton: "📋 Menú", StopButton: "⏹ Parar", Back: "◀ Atrás", BackMenu: "◀ Al menú", MainMenuTitle: "Menú principal", MainMenuBody: "Elige una acción:", Learning: "Aprender", Words: "Palabras", Stats: "Progreso", Settings: "Ajustes", Tools: "Herramientas", NewLesson: "Nueva lección", Practice: "Práctica", LevelTest: "Nivel", LearnWords: "Aprender palabras", WordGame: "Repasar", Spelling: "Ortografía", Vocabulary: "Vocabulario", Mistakes: "Errores", Progress: "Mi progreso", Leaders: "Líderes", Limits: "Límites", BotLanguage: "Idioma del bot", LearningLanguage: "Idioma de estudio", Notifications: "Notificaciones", Premium: "Premium", ChooseBotLang: "Elige el idioma del bot", ChooseLearnLang: "¿Qué idioma quieres estudiar?", ChooseTimezone: "Elige zona horaria", TimezoneHint: "Enviaré recordatorios diarios sobre las 19:00 de tu hora.", BotLangSet: "Listo. Idioma del bot: %s.", LearnLangSet: "Listo. Ahora estudiamos: %s.", UnknownButton: "No entendí el botón. Pulsa 📋 Menú.", Stopped: "Detenido. Pulsa 📋 Menú para elegir otra acción.", WordQuestion: "¿Cómo se dice %s?", ChooseAnswer: "Elige la opción correcta:", Correct: "¡Correcto!", TryAgain: "Aún no. Inténtalo otra vez:", NextWord: "Siguiente palabra", AlreadyLearned: "Esta palabra ya está en tu vocabulario.", AddedToVocab: "Palabra añadida al vocabulario aprendido.", TotalLearned: "Palabras aprendidas: %d", WriteWord: "Escribe %s:", Hint: "Pista", GoodSpelling: "¡Muy bien escrito!"},
	"de": {MenuButton: "📋 Menü", StopButton: "⏹ Stopp", Back: "◀ Zurück", BackMenu: "◀ Zum Menü", MainMenuTitle: "Hauptmenü", MainMenuBody: "Wähle eine Aktion:", Learning: "Lernen", Words: "Wörter", Stats: "Fortschritt", Settings: "Einstellungen", Tools: "Werkzeuge", NewLesson: "Neue Lektion", Practice: "Übung", LevelTest: "Niveau", LearnWords: "Wörter lernen", WordGame: "Wiederholen", Spelling: "Rechtschreibung", Vocabulary: "Wortschatz", Mistakes: "Fehler", Progress: "Mein Fortschritt", Leaders: "Bestenliste", Limits: "Limits", BotLanguage: "Bot-Sprache", LearningLanguage: "Lernsprache", Notifications: "Benachrichtigungen", Premium: "Premium", ChooseBotLang: "Wähle die Bot-Sprache", ChooseLearnLang: "Welche Sprache möchtest du lernen?", ChooseTimezone: "Wähle die Zeitzone", TimezoneHint: "Ich sende tägliche Erinnerungen ungefähr um 19:00 deiner Zeit.", BotLangSet: "Fertig. Bot-Sprache: %s.", LearnLangSet: "Fertig. Wir lernen jetzt: %s.", UnknownButton: "Ich habe die Taste nicht verstanden. Drücke 📋 Menü.", Stopped: "Gestoppt. Drücke 📋 Menü für die nächste Aktion.", WordQuestion: "Wie sagt man %s?", ChooseAnswer: "Wähle die richtige Antwort:", Correct: "Richtig!", TryAgain: "Noch nicht. Versuch es noch einmal:", NextWord: "Nächstes Wort", AlreadyLearned: "Dieses Wort ist schon in deinem Wortschatz.", AddedToVocab: "Wort zum gelernten Wortschatz hinzugefügt.", TotalLearned: "Gelernte Wörter: %d", WriteWord: "Schreibe %s:", Hint: "Hinweis", GoodSpelling: "Sehr gut geschrieben!"},
	"fr": {MenuButton: "📋 Menu", StopButton: "⏹ Stop", Back: "◀ Retour", BackMenu: "◀ Au menu", MainMenuTitle: "Menu principal", MainMenuBody: "Choisis une action :", Learning: "Apprendre", Words: "Mots", Stats: "Progrès", Settings: "Réglages", Tools: "Outils", NewLesson: "Nouvelle leçon", Practice: "Pratique", LevelTest: "Niveau", LearnWords: "Apprendre des mots", WordGame: "Réviser", Spelling: "Orthographe", Vocabulary: "Vocabulaire", Mistakes: "Erreurs", Progress: "Mes progrès", Leaders: "Classement", Limits: "Limites", BotLanguage: "Langue du bot", LearningLanguage: "Langue étudiée", Notifications: "Notifications", Premium: "Premium", ChooseBotLang: "Choisis la langue du bot", ChooseLearnLang: "Quelle langue veux-tu apprendre ?", ChooseTimezone: "Choisis le fuseau horaire", TimezoneHint: "J'enverrai un rappel quotidien vers 19:00, heure locale.", BotLangSet: "C'est prêt. Langue du bot : %s.", LearnLangSet: "C'est prêt. Nous apprenons : %s.", UnknownButton: "Je n'ai pas compris ce bouton. Appuie sur 📋 Menu.", Stopped: "Arrêté. Appuie sur 📋 Menu pour choisir la suite.", WordQuestion: "Comment dit-on %s ?", ChooseAnswer: "Choisis la bonne réponse :", Correct: "Correct !", TryAgain: "Pas encore. Essaie encore :", NextWord: "Mot suivant", AlreadyLearned: "Ce mot est déjà dans ton vocabulaire.", AddedToVocab: "Mot ajouté au vocabulaire appris.", TotalLearned: "Mots appris : %d", WriteWord: "Écris %s :", Hint: "Indice", GoodSpelling: "Très bien écrit !"},
	"it": {MenuButton: "📋 Menu", StopButton: "⏹ Stop", Back: "◀ Indietro", BackMenu: "◀ Al menu", MainMenuTitle: "Menu principale", MainMenuBody: "Scegli un'azione:", Learning: "Studiare", Words: "Parole", Stats: "Progressi", Settings: "Impostazioni", Tools: "Strumenti", NewLesson: "Nuova lezione", Practice: "Pratica", LevelTest: "Livello", LearnWords: "Impara parole", WordGame: "Ripasso", Spelling: "Ortografia", Vocabulary: "Vocabolario", Mistakes: "Errori", Progress: "I miei progressi", Leaders: "Classifica", Limits: "Limiti", BotLanguage: "Lingua del bot", LearningLanguage: "Lingua da studiare", Notifications: "Notifiche", Premium: "Premium", ChooseBotLang: "Scegli la lingua del bot", ChooseLearnLang: "Che lingua vuoi studiare?", ChooseTimezone: "Scegli il fuso orario", TimezoneHint: "Manderò promemoria giornalieri circa alle 19:00 del tuo orario.", BotLangSet: "Fatto. Lingua del bot: %s.", LearnLangSet: "Fatto. Ora studiamo: %s.", UnknownButton: "Non ho capito il pulsante. Premi 📋 Menu.", Stopped: "Fermato. Premi 📋 Menu per scegliere la prossima azione.", WordQuestion: "Come si dice %s?", ChooseAnswer: "Scegli la risposta corretta:", Correct: "Corretto!", TryAgain: "Non ancora. Riprova:", NextWord: "Parola successiva", AlreadyLearned: "Questa parola è già nel tuo vocabolario.", AddedToVocab: "Parola aggiunta al vocabolario imparato.", TotalLearned: "Parole imparate: %d", WriteWord: "Scrivi %s:", Hint: "Suggerimento", GoodSpelling: "Scritto benissimo!"},
	"pl": {MenuButton: "📋 Menu", StopButton: "⏹ Stop", Back: "◀ Wstecz", BackMenu: "◀ Do menu", MainMenuTitle: "Menu główne", MainMenuBody: "Wybierz działanie:", Learning: "Nauka", Words: "Słowa", Stats: "Postęp", Settings: "Ustawienia", Tools: "Narzędzia", NewLesson: "Nowa lekcja", Practice: "Praktyka", LevelTest: "Poziom", LearnWords: "Ucz się słów", WordGame: "Powtórka", Spelling: "Pisownia", Vocabulary: "Słownik", Mistakes: "Błędy", Progress: "Mój postęp", Leaders: "Ranking", Limits: "Limity", BotLanguage: "Język bota", LearningLanguage: "Język nauki", Notifications: "Powiadomienia", Premium: "Premium", ChooseBotLang: "Wybierz język bota", ChooseLearnLang: "Jakiego języka chcesz się uczyć?", ChooseTimezone: "Wybierz strefę czasową", TimezoneHint: "Będę wysyłać codzienne przypomnienia około 19:00 twojego czasu.", BotLangSet: "Gotowe. Język bota: %s.", LearnLangSet: "Gotowe. Teraz uczymy się: %s.", UnknownButton: "Nie rozumiem przycisku. Naciśnij 📋 Menu.", Stopped: "Zatrzymano. Naciśnij 📋 Menu, aby wybrać kolejne działanie.", WordQuestion: "Jak powiedzieć %s?", ChooseAnswer: "Wybierz poprawną odpowiedź:", Correct: "Dobrze!", TryAgain: "Jeszcze nie. Spróbuj ponownie:", NextWord: "Następne słowo", AlreadyLearned: "To słowo jest już w twoim słowniku.", AddedToVocab: "Słowo dodane do nauczonego słownika.", TotalLearned: "Nauczone słowa: %d", WriteWord: "Napisz %s:", Hint: "Podpowiedź", GoodSpelling: "Świetnie napisane!"},
	"pt": {MenuButton: "📋 Menu", StopButton: "⏹ Parar", Back: "◀ Voltar", BackMenu: "◀ Ao menu", MainMenuTitle: "Menu principal", MainMenuBody: "Escolhe uma ação:", Learning: "Aprender", Words: "Palavras", Stats: "Progresso", Settings: "Definições", Tools: "Ferramentas", NewLesson: "Nova lição", Practice: "Prática", LevelTest: "Nível", LearnWords: "Aprender palavras", WordGame: "Revisão", Spelling: "Ortografia", Vocabulary: "Vocabulário", Mistakes: "Erros", Progress: "Meu progresso", Leaders: "Classificação", Limits: "Limites", BotLanguage: "Idioma do bot", LearningLanguage: "Idioma de estudo", Notifications: "Notificações", Premium: "Premium", ChooseBotLang: "Escolhe o idioma do bot", ChooseLearnLang: "Que idioma queres estudar?", ChooseTimezone: "Escolhe o fuso horário", TimezoneHint: "Vou enviar lembretes diários por volta das 19:00 no teu horário.", BotLangSet: "Pronto. Idioma do bot: %s.", LearnLangSet: "Pronto. Agora estudamos: %s.", UnknownButton: "Não entendi o botão. Toca em 📋 Menu.", Stopped: "Parado. Toca em 📋 Menu para escolher a próxima ação.", WordQuestion: "Como se diz %s?", ChooseAnswer: "Escolhe a resposta correta:", Correct: "Correto!", TryAgain: "Ainda não. Tenta outra vez:", NextWord: "Próxima palavra", AlreadyLearned: "Esta palavra já está no teu vocabulário.", AddedToVocab: "Palavra adicionada ao vocabulário aprendido.", TotalLearned: "Palavras aprendidas: %d", WriteWord: "Escreve %s:", Hint: "Dica", GoodSpelling: "Muito bem escrito!"},
	"ro": {MenuButton: "📋 Meniu", StopButton: "⏹ Stop", Back: "◀ Înapoi", BackMenu: "◀ La meniu", MainMenuTitle: "Meniu principal", MainMenuBody: "Alege o acțiune:", Learning: "Învățare", Words: "Cuvinte", Stats: "Progres", Settings: "Setări", Tools: "Instrumente", NewLesson: "Lecție nouă", Practice: "Practică", LevelTest: "Nivel", LearnWords: "Învață cuvinte", WordGame: "Repetare", Spelling: "Ortografie", Vocabulary: "Vocabular", Mistakes: "Greșeli", Progress: "Progresul meu", Leaders: "Clasament", Limits: "Limite", BotLanguage: "Limba botului", LearningLanguage: "Limba de studiu", Notifications: "Notificări", Premium: "Premium", ChooseBotLang: "Alege limba botului", ChooseLearnLang: "Ce limbă vrei să înveți?", ChooseTimezone: "Alege fusul orar", TimezoneHint: "Voi trimite mementouri zilnice în jur de 19:00, ora ta.", BotLangSet: "Gata. Limba botului: %s.", LearnLangSet: "Gata. Acum învățăm: %s.", UnknownButton: "Nu am înțeles butonul. Apasă 📋 Meniu.", Stopped: "Oprit. Apasă 📋 Meniu pentru următoarea acțiune.", WordQuestion: "Cum se spune %s?", ChooseAnswer: "Alege răspunsul corect:", Correct: "Corect!", TryAgain: "Încă nu. Încearcă din nou:", NextWord: "Următorul cuvânt", AlreadyLearned: "Acest cuvânt este deja în vocabularul tău.", AddedToVocab: "Cuvânt adăugat în vocabularul învățat.", TotalLearned: "Cuvinte învățate: %d", WriteWord: "Scrie %s:", Hint: "Indiciu", GoodSpelling: "Foarte bine scris!"},
	"uk": {MenuButton: "📋 Меню", StopButton: "⏹ Стоп", Back: "◀ Назад", BackMenu: "◀ До меню", MainMenuTitle: "Головне меню", MainMenuBody: "Вибери дію:", Learning: "Навчання", Words: "Слова", Stats: "Прогрес", Settings: "Налаштування", Tools: "Інструменти", NewLesson: "Новий урок", Practice: "Практика", LevelTest: "Рівень", LearnWords: "Вчити слова", WordGame: "Повторення", Spelling: "Правопис", Vocabulary: "Словник", Mistakes: "Помилки", Progress: "Мій прогрес", Leaders: "Лідери", Limits: "Ліміти", BotLanguage: "Мова бота", LearningLanguage: "Мова навчання", Notifications: "Сповіщення", Premium: "Premium", ChooseBotLang: "Вибери мову бота", ChooseLearnLang: "Яку мову хочеш вивчати?", ChooseTimezone: "Вибери часовий пояс", TimezoneHint: "Я надсилатиму щоденні нагадування приблизно о 19:00 за твоїм часом.", BotLangSet: "Готово. Мова бота: %s.", LearnLangSet: "Готово. Тепер вивчаємо: %s.", UnknownButton: "Не зрозумів кнопку. Натисни 📋 Меню.", Stopped: "Зупинено. Натисни 📋 Меню, щоб вибрати наступну дію.", WordQuestion: "Як сказати %s?", ChooseAnswer: "Вибери правильний варіант:", Correct: "Правильно!", TryAgain: "Поки не так. Спробуй ще раз:", NextWord: "Наступне слово", AlreadyLearned: "Це слово вже є у твоєму словнику.", AddedToVocab: "Слово додано до вивченого словника.", TotalLearned: "Вивчено слів: %d", WriteWord: "Напиши %s:", Hint: "Підказка", GoodSpelling: "Чудово написано!"},
	"kk": {MenuButton: "📋 Мәзір", StopButton: "⏹ Тоқтату", Back: "◀ Артқа", BackMenu: "◀ Мәзірге", MainMenuTitle: "Басты мәзір", MainMenuBody: "Әрекетті таңда:", Learning: "Оқу", Words: "Сөздер", Stats: "Прогресс", Settings: "Баптаулар", Tools: "Құралдар", NewLesson: "Жаңа сабақ", Practice: "Практика", LevelTest: "Деңгей", LearnWords: "Сөз үйрену", WordGame: "Қайталау", Spelling: "Емле", Vocabulary: "Сөздік", Mistakes: "Қателер", Progress: "Менің прогресім", Leaders: "Көшбасшылар", Limits: "Лимиттер", BotLanguage: "Бот тілі", LearningLanguage: "Оқу тілі", Notifications: "Хабарламалар", Premium: "Premium", ChooseBotLang: "Бот тілін таңда", ChooseLearnLang: "Қай тілді үйренгің келеді?", ChooseTimezone: "Уақыт белдеуін таңда", TimezoneHint: "Күнделікті еске салуды сенің уақытыңмен шамамен 19:00-де жіберемін.", BotLangSet: "Дайын. Бот тілі: %s.", LearnLangSet: "Дайын. Енді үйренеміз: %s.", UnknownButton: "Түймені түсінбедім. 📋 Мәзірді бас.", Stopped: "Тоқтатылды. Келесі әрекет үшін 📋 Мәзірді бас.", WordQuestion: "%s қалай айтылады?", ChooseAnswer: "Дұрыс жауапты таңда:", Correct: "Дұрыс!", TryAgain: "Әлі дұрыс емес. Қайта көр:", NextWord: "Келесі сөз", AlreadyLearned: "Бұл сөз сенің сөздігіңде бар.", AddedToVocab: "Сөз үйренілген сөздікке қосылды.", TotalLearned: "Үйренілген сөздер: %d", WriteWord: "%s жаз:", Hint: "Көмек", GoodSpelling: "Өте жақсы жаздың!"},
	"ky": {MenuButton: "📋 Меню", StopButton: "⏹ Токтотуу", Back: "◀ Артка", BackMenu: "◀ Менюга", MainMenuTitle: "Башкы меню", MainMenuBody: "Аракетти танда:", Learning: "Окуу", Words: "Сөздөр", Stats: "Прогресс", Settings: "Жөндөөлөр", Tools: "Куралдар", NewLesson: "Жаңы сабак", Practice: "Практика", LevelTest: "Деңгээлди аныктоо", LearnWords: "Сөз үйрөнүү", WordGame: "Кайталоо", Spelling: "Туура жазуу", Vocabulary: "Сөздүк", Mistakes: "Каталар", Progress: "Менин прогрессим", Leaders: "Лидерлер", Limits: "Лимиттер", BotLanguage: "Бот тили", LearningLanguage: "Үйрөнүү тили", Notifications: "Эскертмелер", Premium: "Premium", ChooseBotLang: "Боттун тилин танда", ChooseLearnLang: "Кайсы тилди үйрөнгүң келет?", ChooseTimezone: "Убакыт алкагын танда", TimezoneHint: "Күн сайын сенин убактың боюнча болжол менен 19:00дө эскертме жиберем.", BotLangSet: "Даяр. Бот тили: %s.", LearnLangSet: "Даяр. Эми үйрөнөбүз: %s.", UnknownButton: "Бул баскычты түшүнгөн жокмун. 📋 Менюну бас.", Stopped: "Токтотулду. Кийинки аракет үчүн 📋 Менюну бас.", WordQuestion: "%s кандай айтылат?", ChooseAnswer: "Туура жоопту танда:", Correct: "Туура!", TryAgain: "Азырынча туура эмес. Дагы аракет кыл:", NextWord: "Кийинки сөз", AlreadyLearned: "Бул сөз сенин сөздүгүңдө бар.", AddedToVocab: "Сөз үйрөнүлгөн сөздүккө кошулду.", TotalLearned: "Үйрөнүлгөн сөздөр: %d", WriteWord: "%s жаз:", Hint: "Кеңеш", GoodSpelling: "Абдан жакшы жазылды!"},
	"ka": {MenuButton: "📋 მენიუ", StopButton: "⏹ გაჩერება", Back: "◀ უკან", BackMenu: "◀ მენიუში", MainMenuTitle: "მთავარი მენიუ", MainMenuBody: "აირჩიე მოქმედება:", Learning: "სწავლა", Words: "სიტყვები", Stats: "პროგრესი", Settings: "პარამეტრები", Tools: "ინსტრუმენტები", NewLesson: "ახალი გაკვეთილი", Practice: "პრაქტიკა", LevelTest: "დონის განსაზღვრა", LearnWords: "სიტყვების სწავლა", WordGame: "გამეორება", Spelling: "მართლწერა", Vocabulary: "ლექსიკონი", Mistakes: "შეცდომები", Progress: "ჩემი პროგრესი", Leaders: "ლიდერები", Limits: "ლიმიტები", BotLanguage: "ბოტის ენა", LearningLanguage: "სასწავლო ენა", Notifications: "შეხსენებები", Premium: "Premium", ChooseBotLang: "აირჩიე ბოტის ენა", ChooseLearnLang: "რომელი ენის სწავლა გინდა?", ChooseTimezone: "აირჩიე დროის სარტყელი", TimezoneHint: "ყოველდღე დაახლოებით 19:00-ზე, შენი დროით, შეგახსენებ.", BotLangSet: "მზადაა. ბოტის ენა: %s.", LearnLangSet: "მზადაა. ახლა ვსწავლობთ: %s.", UnknownButton: "ეს ღილაკი ვერ გავიგე. დააჭირე 📋 მენიუს.", Stopped: "გაჩერებულია. შემდეგი მოქმედებისთვის დააჭირე 📋 მენიუს.", WordQuestion: "როგორ იქნება %s?", ChooseAnswer: "აირჩიე სწორი პასუხი:", Correct: "სწორია!", TryAgain: "ჯერ არა. სცადე კიდევ ერთხელ:", NextWord: "შემდეგი სიტყვა", AlreadyLearned: "ეს სიტყვა უკვე შენს ლექსიკონშია.", AddedToVocab: "სიტყვა დაემატა ნასწავლ ლექსიკონს.", TotalLearned: "ნასწავლი სიტყვები: %d", WriteWord: "დაწერე %s:", Hint: "მინიშნება", GoodSpelling: "ძალიან კარგად დაწერე!"},
	"uz": {MenuButton: "📋 Menyu", StopButton: "⏹ To‘xtatish", Back: "◀ Orqaga", BackMenu: "◀ Menyuga", MainMenuTitle: "Asosiy menyu", MainMenuBody: "Amalni tanlang:", Learning: "O‘rganish", Words: "So‘zlar", Stats: "Natija", Settings: "Sozlamalar", Tools: "Vositalar", NewLesson: "Yangi dars", Practice: "Mashq", LevelTest: "Daraja", LearnWords: "So‘z o‘rganish", WordGame: "Takrorlash", Spelling: "Imlo", Vocabulary: "Lug‘at", Mistakes: "Xatolar", Progress: "Mening natijam", Leaders: "Reyting", Limits: "Limitlar", BotLanguage: "Bot tili", LearningLanguage: "O‘rganish tili", Notifications: "Eslatmalar", Premium: "Premium", ChooseBotLang: "Bot tilini tanlang", ChooseLearnLang: "Qaysi tilni o‘rganmoqchisiz?", ChooseTimezone: "Vaqt mintaqasini tanlang", TimezoneHint: "Har kuni taxminan soat 19:00 da eslatma yuboraman.", BotLangSet: "Tayyor. Bot tili: %s.", LearnLangSet: "Tayyor. Endi o‘rganamiz: %s.", UnknownButton: "Tugmani tushunmadim. 📋 Menyu ni bosing.", Stopped: "To‘xtatildi. Keyingi amal uchun 📋 Menyu ni bosing.", WordQuestion: "%s qanday aytiladi?", ChooseAnswer: "To‘g‘ri javobni tanlang:", Correct: "To‘g‘ri!", TryAgain: "Hali emas. Yana urinib ko‘ring:", NextWord: "Keyingi so‘z", AlreadyLearned: "Bu so‘z lug‘atingizda bor.", AddedToVocab: "So‘z o‘rganilgan lug‘atga qo‘shildi.", TotalLearned: "O‘rganilgan so‘zlar: %d", WriteWord: "%s yozing:", Hint: "Yordam", GoodSpelling: "Juda yaxshi yozildi!"},
	"tt": {MenuButton: "📋 Меню", StopButton: "⏹ Туктату", Back: "◀ Артка", BackMenu: "◀ Менюга", MainMenuTitle: "Төп меню", MainMenuBody: "Гамәл сайла:", Learning: "Өйрәнү", Words: "Сүзләр", Stats: "Алга китеш", Settings: "Көйләүләр", Tools: "Кораллар", NewLesson: "Яңа дәрес", Practice: "Күнегү", LevelTest: "Дәрәҗә", LearnWords: "Сүз өйрәнү", WordGame: "Кабатлау", Spelling: "Дөрес язу", Vocabulary: "Сүзлек", Mistakes: "Хаталар", Progress: "Минем алга китеш", Leaders: "Лидерлар", Limits: "Лимитлар", BotLanguage: "Бот теле", LearningLanguage: "Өйрәнү теле", Notifications: "Искәртмәләр", Premium: "Premium", ChooseBotLang: "Бот телен сайла", ChooseLearnLang: "Кайсы телне өйрәнәсең килә?", ChooseTimezone: "Вакыт поясын сайла", TimezoneHint: "Көн саен синең вакыт белән якынча 19:00 искәртү җибәрәм.", BotLangSet: "Әзер. Бот теле: %s.", LearnLangSet: "Әзер. Хәзер өйрәнәбез: %s.", UnknownButton: "Төймәне аңламадым. 📋 Менюга бас.", Stopped: "Туктатылды. Киләсе гамәл өчен 📋 Менюга бас.", WordQuestion: "%s ничек әйтелә?", ChooseAnswer: "Дөрес җавапны сайла:", Correct: "Дөрес!", TryAgain: "Әлегә түгел. Кабатлап кара:", NextWord: "Киләсе сүз", AlreadyLearned: "Бу сүз синең сүзлегеңдә инде бар.", AddedToVocab: "Сүз өйрәнелгән сүзлеккә өстәлде.", TotalLearned: "Өйрәнелгән сүзләр: %d", WriteWord: "%s яз:", Hint: "Киңәш", GoodSpelling: "Бик яхшы язылган!"},
	"tg": {MenuButton: "📋 Меню", StopButton: "⏹ Ист", Back: "◀ Бозгашт", BackMenu: "◀ Ба меню", MainMenuTitle: "Менюи асосӣ", MainMenuBody: "Амалро интихоб кунед:", Learning: "Омӯзиш", Words: "Калимаҳо", Stats: "Пешрафт", Settings: "Танзимот", Tools: "Абзорҳо", NewLesson: "Дарси нав", Practice: "Машқ", LevelTest: "Муайян кардани сатҳ", LearnWords: "Омӯхтани калимаҳо", WordGame: "Такрор", Spelling: "Имло", Vocabulary: "Луғат", Mistakes: "Хатоҳо", Progress: "Пешрафти ман", Leaders: "Пешсафон", Limits: "Лимитҳо", BotLanguage: "Забони бот", LearningLanguage: "Забони омӯзиш", Notifications: "Ёдраскуниҳо", Premium: "Premium", ChooseBotLang: "Забони ботро интихоб кунед", ChooseLearnLang: "Кадом забонро омӯхтан мехоҳед?", ChooseTimezone: "Минтақаи вақтро интихоб кунед", TimezoneHint: "Ман ҳар рӯз тақрибан соати 19:00 бо вақти шумо ёдрас мефиристам.", BotLangSet: "Тайёр. Забони бот: %s.", LearnLangSet: "Тайёр. Акнун меомӯзем: %s.", UnknownButton: "Тугмаро нафаҳмидам. 📋 Менюро пахш кунед.", Stopped: "Қатъ шуд. Барои амали нав 📋 Менюро пахш кунед.", WordQuestion: "%s чӣ тавр гуфта мешавад?", ChooseAnswer: "Ҷавоби дурустро интихоб кунед:", Correct: "Дуруст!", TryAgain: "Ҳоло не. Боз кӯшиш кунед:", NextWord: "Калимаи навбатӣ", AlreadyLearned: "Ин калима аллакай дар луғати шумост.", AddedToVocab: "Калима ба луғати омӯхташуда илова шуд.", TotalLearned: "Калимаҳои омӯхташуда: %d", WriteWord: "%s нависед:", Hint: "Маслиҳат", GoodSpelling: "Хеле хуб навишта шуд!"},
	"hy": {MenuButton: "📋 Մենյու", StopButton: "⏹ Կանգ", Back: "◀ Հետ", BackMenu: "◀ Մենյու", MainMenuTitle: "Գլխավոր մենյու", MainMenuBody: "Ընտրիր գործողությունը:", Learning: "Ուսուցում", Words: "Բառեր", Stats: "Առաջընթաց", Settings: "Կարգավորումներ", Tools: "Գործիքներ", NewLesson: "Նոր դաս", Practice: "Պրակտիկա", LevelTest: "Որոշել մակարդակը", LearnWords: "Սովորել բառեր", WordGame: "Կրկնություն", Spelling: "Ուղղագրություն", Vocabulary: "Բառարան", Mistakes: "Սխալներ", Progress: "Իմ առաջընթացը", Leaders: "Առաջատարներ", Limits: "Սահմանաչափեր", BotLanguage: "Բոտի լեզու", LearningLanguage: "Ուսուցման լեզու", Notifications: "Հիշեցումներ", Premium: "Premium", ChooseBotLang: "Ընտրիր բոտի լեզուն", ChooseLearnLang: "Ո՞ր լեզուն ես ուզում սովորել:", ChooseTimezone: "Ընտրիր ժամային գոտին", TimezoneHint: "Ամեն օր մոտավորապես 19:00-ին քո ժամով հիշեցում կուղարկեմ:", BotLangSet: "Պատրաստ է։ Բոտի լեզուն՝ %s:", LearnLangSet: "Պատրաստ է։ Հիմա սովորում ենք՝ %s:", UnknownButton: "Չհասկացա կոճակը։ Սեղմիր 📋 Մենյու:", Stopped: "Կանգնեցված է։ Հաջորդ գործողության համար սեղմիր 📋 Մենյու:", WordQuestion: "%s ինչպես է ասվում?", ChooseAnswer: "Ընտրիր ճիշտ պատասխանը:", Correct: "Ճիշտ է:", TryAgain: "Դեռ ոչ։ Փորձիր նորից:", NextWord: "Հաջորդ բառը", AlreadyLearned: "Այս բառը արդեն քո բառարանում է:", AddedToVocab: "Բառը ավելացվեց սովորած բառարանին:", TotalLearned: "Սովորած բառեր՝ %d", WriteWord: "Գրիր %s:", Hint: "Հուշում", GoodSpelling: "Շատ լավ է գրված:"},
}

func englishUICopy() uiCopy {
	return uiCopy{
		MenuButton: "📋 Menu", StopButton: "⏹ Stop", Back: "◀ Back", BackMenu: "◀ Menu",
		MainMenuTitle: "Main Menu", MainMenuBody: "Choose an action:",
		Learning: "Learn", Words: "Words", Stats: "Progress", Settings: "Settings", Tools: "Tools",
		NewLesson: "New lesson", Practice: "Practice", AITutor: "AI Tutor", Shadowing: "Listening", LevelTest: "Find level", LearnWords: "Learn words", WordGame: "Review", Spelling: "Spelling", Vocabulary: "Vocabulary", Phrasebook: "Phrasebook", Mistakes: "Mistakes", Progress: "My progress", Leaders: "Leaders", Limits: "Limits", BotLanguage: "Bot language", LearningLanguage: "Learning language", Notifications: "Notifications", Premium: "Premium", Referral: "Referrals",
		ChooseBotLang: "Choose bot language", ChooseLearnLang: "Which language do you want to learn?", ChooseTimezone: "Choose time zone", TimezoneHint: "I will send daily reminders around 19:00 in your local time.",
		BotLangSet: "Done. Bot language: %s.", LearnLangSet: "Done. Now learning: %s.", UnknownButton: "I did not understand that button. Press 📋 Menu.", Stopped: "Stopped. Press 📋 Menu to choose the next action.",
		WordQuestion: "How do you say %s?", ChooseAnswer: "Choose the correct answer:", Correct: "Correct!", TryAgain: "Not quite. Try again:", NextWord: "Next word", NextPage: "Next ▶", AlreadyLearned: "This word is already in your vocabulary.", AddedToVocab: "Word added to your learned vocabulary.", TotalLearned: "Learned words: %d", WriteWord: "Write %s:", Hint: "Hint", GoodSpelling: "Great spelling!",
		Tool: englishToolUICopy(),
	}
}

var uiAliases = map[string]string{}

type compactTelegramTerms struct {
	Menu, Stop, Back, MainMenu, ChooseAction                string
	Learning, Words, Progress, Settings, Tools              string
	NewLesson, Practice, Listening, Level, LearnWords       string
	Review, Spelling, Vocabulary, Mistakes, Leaders         string
	Limits, BotLanguage, LearningLanguage, Notifications    string
	Premium, Referrals, ChooseBotLang, ChooseLearnLang      string
	ChooseTimezone, TimezoneHint, Ready, UnknownButton      string
	Stopped, WordQuestion, ChooseAnswer, Correct            string
	TryAgain, NextWord, NextPage, AlreadyLearned            string
	AddedToVocab, TotalLearned, WriteWord, Hint             string
	GoodSpelling, VoiceToText, ImageTranslate, Translator   string
	GPTAgent, VoicePrompt, ImagePrompt, TranslatorPrompt    string
	VoiceModePrompt, ImageModePrompt, TranslatorModePrompt  string
	ImageOpenToolsPrompt, ImageSinglePhotoPrompt            string
	Transcript, Translation, SourceLanguage, TargetLanguage string
	Auto, WebApp, VoiceDiscuss, ImageDiscuss                string
	VoicePremium, ImagePremium, VoiceLimit, VoiceTooLong    string
	VoiceDownloadFailed, VoiceTranscribeFailed              string
	VoiceTranslationFailed, ImageFileMissing                string
	ImageDownloadFailed, ImageReadFailed                    string
}

var compactTelegramI18n = map[string]compactTelegramTerms{
	"ar": {Menu: "القائمة", Stop: "إيقاف", Back: "رجوع", MainMenu: "القائمة الرئيسية", ChooseAction: "اختر إجراء:", Learning: "تعلّم", Words: "كلمات", Progress: "التقدم", Settings: "الإعدادات", Tools: "الأدوات", NewLesson: "درس جديد", Practice: "تدريب", Listening: "استماع", Level: "تحديد المستوى", LearnWords: "تعلّم الكلمات", Review: "مراجعة", Spelling: "إملاء", Vocabulary: "المفردات", Mistakes: "الأخطاء", Leaders: "المتصدرون", Limits: "الحدود", BotLanguage: "لغة البوت", LearningLanguage: "لغة التعلم", Notifications: "الإشعارات", Premium: "Premium", Referrals: "الإحالات", ChooseBotLang: "اختر لغة البوت", ChooseLearnLang: "أي لغة تريد أن تتعلم؟", ChooseTimezone: "اختر المنطقة الزمنية", TimezoneHint: "سأرسل تذكيرا يوميا حوالي 19:00 حسب وقتك.", Ready: "تم.", UnknownButton: "لم أفهم الزر. اضغط القائمة.", Stopped: "تم الإيقاف. اضغط القائمة لاختيار الإجراء التالي.", WordQuestion: "كيف نقول %s؟", ChooseAnswer: "اختر الإجابة الصحيحة:", Correct: "صحيح!", TryAgain: "ليس تماما. حاول مرة أخرى:", NextWord: "الكلمة التالية", NextPage: "التالي ▶", AlreadyLearned: "هذه الكلمة موجودة بالفعل في مفرداتك.", AddedToVocab: "تمت إضافة الكلمة إلى مفرداتك المتعلمة.", TotalLearned: "الكلمات المتعلمة: %d", WriteWord: "اكتب %s:", Hint: "تلميح", GoodSpelling: "إملاء ممتاز!", VoiceToText: "الصوت إلى نص", ImageTranslate: "ترجمة نص الصورة", Translator: "مترجم", GPTAgent: "وكيل GPT", VoicePrompt: "أرسل رسالة صوتية وسأحولها إلى نص.", ImagePrompt: "أرسل صورة فيها نص وسأقرأه وأترجمه.", TranslatorPrompt: "اختر اللغات ثم أرسل نصا أو صوتا أو صورة.", VoiceModePrompt: "أرسل رسالة صوتية. سأحولها إلى نص دون بدء التدريب.", ImageModePrompt: "أرسل صورة فيها نص. سأقرأ النص وأترجمه.", TranslatorModePrompt: "المترجم جاهز. اختر لغة المصدر والهدف ثم أرسل المحتوى.", ImageOpenToolsPrompt: "لترجمة نص من صورة، افتح الأدوات واختر ترجمة نص الصورة.", ImageSinglePhotoPrompt: "أرسل صورة واحدة فقط عبر أداة ترجمة نص الصورة.", Transcript: "النص", Translation: "الترجمة", SourceLanguage: "لغة المصدر", TargetLanguage: "لغة الترجمة", Auto: "تلقائي", WebApp: "تطبيق الويب", VoiceDiscuss: "يمكنك مناقشة هذا بالنص وسأتابع التدريب.", ImageDiscuss: "يمكنك إرسال صورة أخرى أو كتابة نص للمتابعة.", VoicePremium: "التدريب الصوتي متاح في Premium.", ImagePremium: "ترجمة نص الصورة متاحة في Premium.", VoiceLimit: "انتهى حد الرسائل الصوتية اليوم.\n\nالصوت: %d/%d", VoiceTooLong: "الرسالة الصوتية طويلة جدا.\n\nالحد الأقصى: %d ثانية\nرسالتك: %d ثانية", VoiceDownloadFailed: "تعذر تنزيل الرسالة الصوتية: %s", VoiceTranscribeFailed: "تعذر التعرف على الصوت: %s", VoiceTranslationFailed: "تعذرت ترجمة النص: %s", ImageFileMissing: "تعذر العثور على ملف الصورة. حاول إرسالها مرة أخرى.", ImageDownloadFailed: "تعذر تنزيل الصورة: %s", ImageReadFailed: "تعذرت قراءة الصورة: %s"},
	"bn": {Menu: "মেনু", Stop: "থামুন", Back: "ফিরে যান", MainMenu: "প্রধান মেনু", ChooseAction: "একটি কাজ বেছে নিন:", Learning: "শেখা", Words: "শব্দ", Progress: "অগ্রগতি", Settings: "সেটিংস", Tools: "টুল", NewLesson: "নতুন পাঠ", Practice: "অনুশীলন", Listening: "শোনা", Level: "স্তর নির্ধারণ", LearnWords: "শব্দ শিখুন", Review: "পুনরাবৃত্তি", Spelling: "বানান", Vocabulary: "শব্দভান্ডার", Mistakes: "ভুল", Leaders: "লিডারবোর্ড", Limits: "সীমা", BotLanguage: "বটের ভাষা", LearningLanguage: "শেখার ভাষা", Notifications: "নোটিফিকেশন", Premium: "Premium", Referrals: "রেফারেল", ChooseBotLang: "বটের ভাষা বেছে নিন", ChooseLearnLang: "আপনি কোন ভাষা শিখতে চান?", ChooseTimezone: "টাইম জোন বেছে নিন", TimezoneHint: "আমি আপনার সময়ে প্রায় 19:00-এ দৈনিক মনে করিয়ে দেব।", Ready: "হয়ে গেছে.", UnknownButton: "বোতামটি বুঝিনি। মেনু চাপুন।", Stopped: "থামানো হয়েছে। পরের কাজের জন্য মেনু চাপুন।", WordQuestion: "%s কীভাবে বলা হয়?", ChooseAnswer: "সঠিক উত্তর বেছে নিন:", Correct: "সঠিক!", TryAgain: "এখনও নয়। আবার চেষ্টা করুন:", NextWord: "পরের শব্দ", NextPage: "পরের পৃষ্ঠা ▶", AlreadyLearned: "এই শব্দটি আপনার শব্দভান্ডারে আছে।", AddedToVocab: "শব্দটি শেখা শব্দভান্ডারে যোগ হয়েছে।", TotalLearned: "শেখা শব্দ: %d", WriteWord: "%s লিখুন:", Hint: "ইঙ্গিত", GoodSpelling: "দারুণ বানান!", VoiceToText: "ভয়েস থেকে টেক্সট", ImageTranslate: "ছবির টেক্সট অনুবাদ", Translator: "অনুবাদক", GPTAgent: "GPT এজেন্ট", VoicePrompt: "ভয়েস বার্তা পাঠান, আমি সেটি টেক্সটে বদলাব।", ImagePrompt: "টেক্সটসহ ছবি পাঠান, আমি পড়ে অনুবাদ করব।", TranslatorPrompt: "ভাষা বেছে নিয়ে টেক্সট, ভয়েস বা ছবি পাঠান।", VoiceModePrompt: "ভয়েস বার্তা পাঠান। আমি টেক্সটে বদলাব, অনুশীলন শুরু করব না।", ImageModePrompt: "টেক্সটসহ ছবি পাঠান। আমি পড়ে অনুবাদ করব।", TranslatorModePrompt: "অনুবাদক প্রস্তুত। উৎস ও লক্ষ্য ভাষা বেছে নিন।", ImageOpenToolsPrompt: "ছবির টেক্সট অনুবাদ করতে টুল খুলুন।", ImageSinglePhotoPrompt: "শুধু একটি ছবি পাঠান।", Transcript: "লিপি", Translation: "অনুবাদ", SourceLanguage: "উৎস ভাষা", TargetLanguage: "লক্ষ্য ভাষা", Auto: "স্বয়ংক্রিয়", WebApp: "ওয়েব অ্যাপ", VoiceDiscuss: "আপনি টেক্সটে আলোচনা করতে পারেন, আমি অনুশীলন চালিয়ে যাব।", ImageDiscuss: "আরেকটি ছবি পাঠাতে পারেন বা টেক্সট লিখতে পারেন।", VoicePremium: "ভয়েস অনুশীলন Premium-এ उपलब्ध.", ImagePremium: "ছবির টেক্সট অনুবাদ Premium-এ उपलब्ध.", VoiceLimit: "আজকের ভয়েস সীমা শেষ।\n\nভয়েস: %d/%d", VoiceTooLong: "ভয়েস বার্তাটি খুব দীর্ঘ।\n\nসর্বোচ্চ: %d সেকেন্ড\nআপনার বার্তা: %d সেকেন্ড", VoiceDownloadFailed: "ভয়েস বার্তা ডাউনলোড করা যায়নি: %s", VoiceTranscribeFailed: "ভয়েস শনাক্ত করা যায়নি: %s", VoiceTranslationFailed: "লিপি অনুবাদ করা যায়নি: %s", ImageFileMissing: "ছবি ফাইল পাওয়া যায়নি। আবার পাঠান।", ImageDownloadFailed: "ছবি ডাউনলোড করা যায়নি: %s", ImageReadFailed: "ছবি পড়া যায়নি: %s"},
	"cs": {Menu: "Menu", Stop: "Stop", Back: "Zpět", MainMenu: "Hlavní menu", ChooseAction: "Vyber akci:", Learning: "Učení", Words: "Slova", Progress: "Pokrok", Settings: "Nastavení", Tools: "Nástroje", NewLesson: "Nová lekce", Practice: "Procvičování", Listening: "Poslech", Level: "Zjistit úroveň", LearnWords: "Učit slova", Review: "Opakování", Spelling: "Pravopis", Vocabulary: "Slovník", Mistakes: "Chyby", Leaders: "Žebříček", Limits: "Limity", BotLanguage: "Jazyk bota", LearningLanguage: "Jazyk učení", Notifications: "Oznámení", Premium: "Premium", Referrals: "Doporučení", ChooseBotLang: "Vyber jazyk bota", ChooseLearnLang: "Jaký jazyk se chceš učit?", ChooseTimezone: "Vyber časové pásmo", TimezoneHint: "Denní připomínku pošlu kolem 19:00 tvého času.", Ready: "Hotovo.", UnknownButton: "Tomu tlačítku nerozumím. Stiskni Menu.", Stopped: "Zastaveno. Stiskni Menu pro další akci.", WordQuestion: "Jak se řekne %s?", ChooseAnswer: "Vyber správnou odpověď:", Correct: "Správně!", TryAgain: "Ještě ne. Zkus to znovu:", NextWord: "Další slovo", NextPage: "Další ▶", AlreadyLearned: "Tohle slovo už máš ve slovníku.", AddedToVocab: "Slovo bylo přidáno do naučeného slovníku.", TotalLearned: "Naučená slova: %d", WriteWord: "Napiš %s:", Hint: "Nápověda", GoodSpelling: "Výborný pravopis!", VoiceToText: "Hlas na text", ImageTranslate: "Překlad textu z obrázku", Translator: "Překladač", GPTAgent: "GPT Agent", VoicePrompt: "Pošli hlasovou zprávu a převedu ji na text.", ImagePrompt: "Pošli obrázek s textem, přečtu ho a přeložím.", TranslatorPrompt: "Vyber jazyky a pošli text, hlas nebo fotku.", VoiceModePrompt: "Pošli hlasovou zprávu. Přepíšu ji na text bez spuštění procvičování.", ImageModePrompt: "Pošli obrázek s textem. Přečtu ho a přeložím.", TranslatorModePrompt: "Překladač je připraven. Vyber zdrojový a cílový jazyk.", ImageOpenToolsPrompt: "Pro překlad textu z obrázku otevři Nástroje.", ImageSinglePhotoPrompt: "Pošli prosím jen jeden obrázek.", Transcript: "Přepis", Translation: "Překlad", SourceLanguage: "Zdrojový jazyk", TargetLanguage: "Cílový jazyk", Auto: "Auto", WebApp: "Webová aplikace", VoiceDiscuss: "Můžeš o tom psát textem a budu pokračovat v procvičování.", ImageDiscuss: "Můžeš poslat další fotku nebo napsat text.", VoicePremium: "Hlasové procvičování je dostupné v Premium.", ImagePremium: "Překlad textu z obrázku je dostupný v Premium.", VoiceLimit: "Dnešní limit hlasových zpráv skončil.\n\nHlas: %d/%d", VoiceTooLong: "Hlasová zpráva je moc dlouhá.\n\nMaximum: %d sekund\nTvoje zpráva: %d sekund", VoiceDownloadFailed: "Nepodařilo se stáhnout hlasovou zprávu: %s", VoiceTranscribeFailed: "Nepodařilo se rozpoznat hlas: %s", VoiceTranslationFailed: "Nepodařilo se přeložit přepis: %s", ImageFileMissing: "Soubor obrázku nebyl nalezen. Zkus fotku poslat znovu.", ImageDownloadFailed: "Nepodařilo se stáhnout obrázek: %s", ImageReadFailed: "Nepodařilo se přečíst obrázek: %s"},
	"el": {Menu: "Μενού", Stop: "Στοπ", Back: "Πίσω", MainMenu: "Κύριο μενού", ChooseAction: "Διάλεξε ενέργεια:", Learning: "Μάθηση", Words: "Λέξεις", Progress: "Πρόοδος", Settings: "Ρυθμίσεις", Tools: "Εργαλεία", NewLesson: "Νέο μάθημα", Practice: "Εξάσκηση", Listening: "Ακρόαση", Level: "Εύρεση επιπέδου", LearnWords: "Μάθε λέξεις", Review: "Επανάληψη", Spelling: "Ορθογραφία", Vocabulary: "Λεξιλόγιο", Mistakes: "Λάθη", Leaders: "Κατάταξη", Limits: "Όρια", BotLanguage: "Γλώσσα bot", LearningLanguage: "Γλώσσα μάθησης", Notifications: "Ειδοποιήσεις", Premium: "Premium", Referrals: "Προσκλήσεις", ChooseBotLang: "Διάλεξε γλώσσα bot", ChooseLearnLang: "Ποια γλώσσα θέλεις να μάθεις;", ChooseTimezone: "Διάλεξε ζώνη ώρας", TimezoneHint: "Θα στέλνω καθημερινή υπενθύμιση περίπου στις 19:00.", Ready: "Έτοιμο.", UnknownButton: "Δεν κατάλαβα το κουμπί. Πάτησε Μενού.", Stopped: "Σταμάτησε. Πάτησε Μενού για την επόμενη ενέργεια.", WordQuestion: "Πώς λέγεται %s;", ChooseAnswer: "Διάλεξε τη σωστή απάντηση:", Correct: "Σωστά!", TryAgain: "Όχι ακόμα. Δοκίμασε ξανά:", NextWord: "Επόμενη λέξη", NextPage: "Επόμενο ▶", AlreadyLearned: "Αυτή η λέξη υπάρχει ήδη στο λεξιλόγιό σου.", AddedToVocab: "Η λέξη προστέθηκε στο μαθημένο λεξιλόγιο.", TotalLearned: "Μαθημένες λέξεις: %d", WriteWord: "Γράψε %s:", Hint: "Υπόδειξη", GoodSpelling: "Πολύ καλή ορθογραφία!", VoiceToText: "Φωνή σε κείμενο", ImageTranslate: "Μετάφραση εικόνας", Translator: "Μεταφραστής", GPTAgent: "GPT Agent", VoicePrompt: "Στείλε φωνητικό μήνυμα και θα το κάνω κείμενο.", ImagePrompt: "Στείλε εικόνα με κείμενο και θα τη μεταφράσω.", TranslatorPrompt: "Διάλεξε γλώσσες και στείλε κείμενο, φωνή ή φωτογραφία.", VoiceModePrompt: "Στείλε φωνητικό μήνυμα. Θα το μεταγράψω χωρίς εξάσκηση.", ImageModePrompt: "Στείλε εικόνα με κείμενο. Θα τη διαβάσω και θα τη μεταφράσω.", TranslatorModePrompt: "Ο μεταφραστής είναι έτοιμος. Διάλεξε γλώσσες.", ImageOpenToolsPrompt: "Για μετάφραση εικόνας άνοιξε τα Εργαλεία.", ImageSinglePhotoPrompt: "Στείλε μόνο μία εικόνα.", Transcript: "Απομαγνητοφώνηση", Translation: "Μετάφραση", SourceLanguage: "Γλώσσα πηγής", TargetLanguage: "Γλώσσα μετάφρασης", Auto: "Αυτόματα", WebApp: "Web εφαρμογή", VoiceDiscuss: "Μπορείς να το συζητήσεις με κείμενο.", ImageDiscuss: "Μπορείς να στείλεις άλλη φωτογραφία ή κείμενο.", VoicePremium: "Η φωνητική εξάσκηση είναι διαθέσιμη στο Premium.", ImagePremium: "Η μετάφραση εικόνας είναι διαθέσιμη στο Premium.", VoiceLimit: "Το σημερινό όριο φωνητικών τελείωσε.\n\nΦωνή: %d/%d", VoiceTooLong: "Το φωνητικό μήνυμα είναι πολύ μεγάλο.\n\nΜέγιστο: %d δευτερόλεπτα\nΜήνυμα: %d δευτερόλεπτα", VoiceDownloadFailed: "Δεν έγινε λήψη φωνητικού: %s", VoiceTranscribeFailed: "Δεν αναγνωρίστηκε η φωνή: %s", VoiceTranslationFailed: "Δεν μεταφράστηκε το κείμενο: %s", ImageFileMissing: "Δεν βρέθηκε το αρχείο εικόνας.", ImageDownloadFailed: "Δεν έγινε λήψη εικόνας: %s", ImageReadFailed: "Δεν διαβάστηκε η εικόνα: %s"},
	"hi": {Menu: "मेनू", Stop: "रोकें", Back: "वापस", MainMenu: "मुख्य मेनू", ChooseAction: "कार्रवाई चुनें:", Learning: "सीखना", Words: "शब्द", Progress: "प्रगति", Settings: "सेटिंग्स", Tools: "टूल", NewLesson: "नया पाठ", Practice: "अभ्यास", Listening: "सुनना", Level: "स्तर खोजें", LearnWords: "शब्द सीखें", Review: "दोहराव", Spelling: "वर्तनी", Vocabulary: "शब्दावली", Mistakes: "गलतियां", Leaders: "लीडरबोर्ड", Limits: "सीमाएं", BotLanguage: "बॉट भाषा", LearningLanguage: "सीखने की भाषा", Notifications: "सूचनाएं", Premium: "Premium", Referrals: "रेफरल", ChooseBotLang: "बॉट की भाषा चुनें", ChooseLearnLang: "आप कौन सी भाषा सीखना चाहते हैं?", ChooseTimezone: "समय क्षेत्र चुनें", TimezoneHint: "मैं आपके समय के अनुसार लगभग 19:00 पर दैनिक याद दिलाऊंगा.", Ready: "हो गया.", UnknownButton: "मैंने बटन नहीं समझा. मेनू दबाएं.", Stopped: "रोक दिया गया. अगली कार्रवाई के लिए मेनू दबाएं.", WordQuestion: "%s कैसे कहते हैं?", ChooseAnswer: "सही उत्तर चुनें:", Correct: "सही!", TryAgain: "अभी नहीं. फिर कोशिश करें:", NextWord: "अगला शब्द", NextPage: "अगला ▶", AlreadyLearned: "यह शब्द आपकी शब्दावली में पहले से है.", AddedToVocab: "शब्द सीखी हुई शब्दावली में जोड़ दिया गया.", TotalLearned: "सीखे शब्द: %d", WriteWord: "%s लिखें:", Hint: "संकेत", GoodSpelling: "बहुत अच्छी वर्तनी!", VoiceToText: "आवाज से टेक्स्ट", ImageTranslate: "चित्र टेक्स्ट अनुवाद", Translator: "अनुवादक", GPTAgent: "GPT एजेंट", VoicePrompt: "आवाज संदेश भेजें, मैं उसे टेक्स्ट में बदल दूंगा.", ImagePrompt: "टेक्स्ट वाली छवि भेजें, मैं पढ़कर अनुवाद करूंगा.", TranslatorPrompt: "भाषाएं चुनें, फिर टेक्स्ट, आवाज या फोटो भेजें.", VoiceModePrompt: "आवाज संदेश भेजें. मैं टेक्स्ट बनाऊंगा और अभ्यास शुरू नहीं करूंगा.", ImageModePrompt: "टेक्स्ट वाली छवि भेजें. मैं पढ़कर अनुवाद करूंगा.", TranslatorModePrompt: "अनुवादक तैयार है. स्रोत और लक्ष्य भाषा चुनें.", ImageOpenToolsPrompt: "चित्र से टेक्स्ट अनुवाद के लिए टूल खोलें.", ImageSinglePhotoPrompt: "कृपया केवल एक छवि भेजें.", Transcript: "लिप्यंतरण", Translation: "अनुवाद", SourceLanguage: "स्रोत भाषा", TargetLanguage: "लक्ष्य भाषा", Auto: "स्वचालित", WebApp: "वेब ऐप", VoiceDiscuss: "आप टेक्स्ट में चर्चा कर सकते हैं, मैं अभ्यास जारी रखूंगा.", ImageDiscuss: "आप दूसरी फोटो भेज सकते हैं या टेक्स्ट लिख सकते हैं.", VoicePremium: "आवाज अभ्यास Premium में उपलब्ध है.", ImagePremium: "चित्र टेक्स्ट अनुवाद Premium में उपलब्ध है.", VoiceLimit: "आज की आवाज सीमा खत्म हो गई.\n\nआवाज: %d/%d", VoiceTooLong: "आवाज संदेश बहुत लंबा है.\n\nअधिकतम: %d सेकंड\nआपका संदेश: %d सेकंड", VoiceDownloadFailed: "आवाज संदेश डाउनलोड नहीं हुआ: %s", VoiceTranscribeFailed: "आवाज पहचानी नहीं गई: %s", VoiceTranslationFailed: "लिपि का अनुवाद नहीं हुआ: %s", ImageFileMissing: "छवि फ़ाइल नहीं मिली. फिर भेजें.", ImageDownloadFailed: "छवि डाउनलोड नहीं हुई: %s", ImageReadFailed: "छवि पढ़ी नहीं गई: %s"},
	"hu": {Menu: "Menü", Stop: "Stop", Back: "Vissza", MainMenu: "Főmenü", ChooseAction: "Válassz műveletet:", Learning: "Tanulás", Words: "Szavak", Progress: "Haladás", Settings: "Beállítások", Tools: "Eszközök", NewLesson: "Új lecke", Practice: "Gyakorlás", Listening: "Hallásértés", Level: "Szintfelmérés", LearnWords: "Szavak tanulása", Review: "Ismétlés", Spelling: "Helyesírás", Vocabulary: "Szókincs", Mistakes: "Hibák", Leaders: "Ranglista", Limits: "Limitek", BotLanguage: "Bot nyelve", LearningLanguage: "Tanult nyelv", Notifications: "Értesítések", Premium: "Premium", Referrals: "Ajánlások", ChooseBotLang: "Válaszd ki a bot nyelvét", ChooseLearnLang: "Melyik nyelvet szeretnéd tanulni?", ChooseTimezone: "Válassz időzónát", TimezoneHint: "Napi emlékeztetőt küldök kb. 19:00-kor.", Ready: "Kész.", UnknownButton: "Nem értettem a gombot. Nyomd meg a Menü gombot.", Stopped: "Leállítva. Nyomd meg a Menü gombot a folytatáshoz.", WordQuestion: "Hogyan mondjuk: %s?", ChooseAnswer: "Válaszd ki a helyes választ:", Correct: "Helyes!", TryAgain: "Még nem. Próbáld újra:", NextWord: "Következő szó", NextPage: "Következő ▶", AlreadyLearned: "Ez a szó már benne van a szókincsedben.", AddedToVocab: "A szó bekerült a megtanult szókincsbe.", TotalLearned: "Megtanult szavak: %d", WriteWord: "Írd le: %s", Hint: "Tipp", GoodSpelling: "Nagyon jó helyesírás!", VoiceToText: "Hangból szöveg", ImageTranslate: "Képszöveg fordítása", Translator: "Fordító", GPTAgent: "GPT ügynök", VoicePrompt: "Küldj hangüzenetet, és szöveggé alakítom.", ImagePrompt: "Küldj szöveges képet, elolvasom és lefordítom.", TranslatorPrompt: "Válassz nyelveket, majd küldj szöveget, hangot vagy fotót.", VoiceModePrompt: "Küldj hangüzenetet. Leírom szövegként, gyakorlás nélkül.", ImageModePrompt: "Küldj szöveges képet. Elolvasom és lefordítom.", TranslatorModePrompt: "A fordító készen áll. Válassz forrás- és célnyelvet.", ImageOpenToolsPrompt: "Képszöveg fordításához nyisd meg az Eszközöket.", ImageSinglePhotoPrompt: "Csak egy képet küldj.", Transcript: "Átirat", Translation: "Fordítás", SourceLanguage: "Forrásnyelv", TargetLanguage: "Célnyelv", Auto: "Automatikus", WebApp: "Webalkalmazás", VoiceDiscuss: "Szövegben megbeszélheted, folytatom a gyakorlást.", ImageDiscuss: "Küldhetsz új fotót vagy írhatsz szöveget.", VoicePremium: "A hangos gyakorlás Premium csomagban érhető el.", ImagePremium: "A képszöveg fordítása Premium csomagban érhető el.", VoiceLimit: "A mai hangüzenet-limit elfogyott.\n\nHang: %d/%d", VoiceTooLong: "A hangüzenet túl hosszú.\n\nMaximum: %d mp\nÜzeneted: %d mp", VoiceDownloadFailed: "Nem sikerült letölteni a hangüzenetet: %s", VoiceTranscribeFailed: "Nem sikerült felismerni a hangot: %s", VoiceTranslationFailed: "Nem sikerült lefordítani az átiratot: %s", ImageFileMissing: "Nem találom a képfájlt. Küldd el újra.", ImageDownloadFailed: "Nem sikerült letölteni a képet: %s", ImageReadFailed: "Nem sikerült elolvasni a képet: %s"},
	"id": {Menu: "Menu", Stop: "Berhenti", Back: "Kembali", MainMenu: "Menu utama", ChooseAction: "Pilih tindakan:", Learning: "Belajar", Words: "Kata", Progress: "Progres", Settings: "Pengaturan", Tools: "Alat", NewLesson: "Pelajaran baru", Practice: "Latihan", Listening: "Mendengar", Level: "Cari level", LearnWords: "Belajar kata", Review: "Ulasan", Spelling: "Ejaan", Vocabulary: "Kosakata", Mistakes: "Kesalahan", Leaders: "Peringkat", Limits: "Batas", BotLanguage: "Bahasa bot", LearningLanguage: "Bahasa belajar", Notifications: "Notifikasi", Premium: "Premium", Referrals: "Referal", ChooseBotLang: "Pilih bahasa bot", ChooseLearnLang: "Bahasa apa yang ingin kamu pelajari?", ChooseTimezone: "Pilih zona waktu", TimezoneHint: "Saya akan mengirim pengingat harian sekitar 19:00 waktu lokal.", Ready: "Selesai.", UnknownButton: "Saya tidak memahami tombol itu. Tekan Menu.", Stopped: "Dihentikan. Tekan Menu untuk memilih tindakan berikutnya.", WordQuestion: "Bagaimana mengatakan %s?", ChooseAnswer: "Pilih jawaban yang benar:", Correct: "Benar!", TryAgain: "Belum tepat. Coba lagi:", NextWord: "Kata berikutnya", NextPage: "Berikutnya ▶", AlreadyLearned: "Kata ini sudah ada di kosakatamu.", AddedToVocab: "Kata ditambahkan ke kosakata yang dipelajari.", TotalLearned: "Kata dipelajari: %d", WriteWord: "Tulis %s:", Hint: "Petunjuk", GoodSpelling: "Ejaan bagus!", VoiceToText: "Suara ke teks", ImageTranslate: "Terjemah teks gambar", Translator: "Penerjemah", GPTAgent: "Agen GPT", VoicePrompt: "Kirim pesan suara dan saya ubah menjadi teks.", ImagePrompt: "Kirim gambar berisi teks, saya baca dan terjemahkan.", TranslatorPrompt: "Pilih bahasa, lalu kirim teks, suara, atau foto.", VoiceModePrompt: "Kirim pesan suara. Saya transkripsikan tanpa memulai latihan.", ImageModePrompt: "Kirim gambar berisi teks. Saya baca dan terjemahkan.", TranslatorModePrompt: "Penerjemah siap. Pilih bahasa sumber dan target.", ImageOpenToolsPrompt: "Untuk menerjemahkan teks gambar, buka Alat.", ImageSinglePhotoPrompt: "Kirim hanya satu gambar.", Transcript: "Transkrip", Translation: "Terjemahan", SourceLanguage: "Bahasa sumber", TargetLanguage: "Bahasa target", Auto: "Otomatis", WebApp: "Aplikasi web", VoiceDiscuss: "Kamu bisa membahas ini lewat teks, saya lanjutkan latihan.", ImageDiscuss: "Kamu bisa kirim foto lain atau menulis teks.", VoicePremium: "Latihan suara tersedia di Premium.", ImagePremium: "Terjemah teks gambar tersedia di Premium.", VoiceLimit: "Batas pesan suara hari ini habis.\n\nSuara: %d/%d", VoiceTooLong: "Pesan suara terlalu panjang.\n\nMaksimum: %d detik\nPesanmu: %d detik", VoiceDownloadFailed: "Gagal mengunduh pesan suara: %s", VoiceTranscribeFailed: "Gagal mengenali suara: %s", VoiceTranslationFailed: "Gagal menerjemahkan transkrip: %s", ImageFileMissing: "File gambar tidak ditemukan. Coba kirim ulang.", ImageDownloadFailed: "Gagal mengunduh gambar: %s", ImageReadFailed: "Gagal membaca gambar: %s"},
	"nl": {Menu: "Menu", Stop: "Stop", Back: "Terug", MainMenu: "Hoofdmenu", ChooseAction: "Kies een actie:", Learning: "Leren", Words: "Woorden", Progress: "Voortgang", Settings: "Instellingen", Tools: "Tools", NewLesson: "Nieuwe les", Practice: "Oefenen", Listening: "Luisteren", Level: "Niveau bepalen", LearnWords: "Woorden leren", Review: "Herhalen", Spelling: "Spelling oefenen", Vocabulary: "Woordenschat", Mistakes: "Fouten", Leaders: "Ranglijst", Limits: "Limieten", BotLanguage: "Bottaal", LearningLanguage: "Leertaal", Notifications: "Meldingen", Premium: "Premium", Referrals: "Verwijzingen", ChooseBotLang: "Kies de bottaal", ChooseLearnLang: "Welke taal wil je leren?", ChooseTimezone: "Kies tijdzone", TimezoneHint: "Ik stuur dagelijks een herinnering rond 19:00 lokale tijd.", Ready: "Klaar.", UnknownButton: "Ik begreep die knop niet. Druk op Menu.", Stopped: "Gestopt. Druk op Menu voor de volgende actie.", WordQuestion: "Hoe zeg je %s?", ChooseAnswer: "Kies het juiste antwoord:", Correct: "Goed!", TryAgain: "Nog niet. Probeer opnieuw:", NextWord: "Volgend woord", NextPage: "Volgende ▶", AlreadyLearned: "Dit woord staat al in je woordenschat.", AddedToVocab: "Woord toegevoegd aan je geleerde woordenschat.", TotalLearned: "Geleerde woorden: %d", WriteWord: "Schrijf %s:", Hint: "Hint", GoodSpelling: "Prima spelling!", VoiceToText: "Spraak naar tekst", ImageTranslate: "Beeldtekst vertalen", Translator: "Vertaler", GPTAgent: "GPT-agent", VoicePrompt: "Stuur een spraakbericht en ik zet het om naar tekst.", ImagePrompt: "Stuur een afbeelding met tekst, ik lees en vertaal die.", TranslatorPrompt: "Kies talen en stuur tekst, stem of foto.", VoiceModePrompt: "Stuur een spraakbericht. Ik transcribeer zonder oefening te starten.", ImageModePrompt: "Stuur een afbeelding met tekst. Ik lees en vertaal die.", TranslatorModePrompt: "Vertaler is klaar. Kies bron- en doeltaal.", ImageOpenToolsPrompt: "Open Tools om tekst uit een afbeelding te vertalen.", ImageSinglePhotoPrompt: "Stuur slechts één afbeelding.", Transcript: "Transcript", Translation: "Vertaling", SourceLanguage: "Brontaal", TargetLanguage: "Doeltaal", Auto: "Auto", WebApp: "Webapp", VoiceDiscuss: "Je kunt dit via tekst bespreken, ik ga door met oefenen.", ImageDiscuss: "Je kunt nog een foto sturen of tekst schrijven.", VoicePremium: "Spraakoefening is beschikbaar in Premium.", ImagePremium: "Beeldtekst vertalen is beschikbaar in Premium.", VoiceLimit: "Je spraaklimiet voor vandaag is op.\n\nSpraak: %d/%d", VoiceTooLong: "Het spraakbericht is te lang.\n\nMaximum: %d seconden\nJe bericht: %d seconden", VoiceDownloadFailed: "Kon spraakbericht niet downloaden: %s", VoiceTranscribeFailed: "Kon spraak niet herkennen: %s", VoiceTranslationFailed: "Kon transcript niet vertalen: %s", ImageFileMissing: "Afbeeldingsbestand niet gevonden. Stuur de foto opnieuw.", ImageDownloadFailed: "Kon afbeelding niet downloaden: %s", ImageReadFailed: "Kon afbeelding niet lezen: %s"},
	"sv": {Menu: "Meny", Stop: "Stopp", Back: "Tillbaka", MainMenu: "Huvudmeny", ChooseAction: "Välj åtgärd:", Learning: "Lärande", Words: "Ord", Progress: "Framsteg", Settings: "Inställningar", Tools: "Verktyg", NewLesson: "Ny lektion", Practice: "Övning", Listening: "Lyssning", Level: "Hitta nivå", LearnWords: "Lär ord", Review: "Repetition", Spelling: "Stavning", Vocabulary: "Ordförråd", Mistakes: "Misstag", Leaders: "Topplista", Limits: "Gränser", BotLanguage: "Botens språk", LearningLanguage: "Språk att lära", Notifications: "Aviseringar", Premium: "Premium", Referrals: "Hänvisningar", ChooseBotLang: "Välj botens språk", ChooseLearnLang: "Vilket språk vill du lära dig?", ChooseTimezone: "Välj tidszon", TimezoneHint: "Jag skickar dagliga påminnelser runt 19:00 lokal tid.", Ready: "Klart.", UnknownButton: "Jag förstod inte knappen. Tryck på Meny.", Stopped: "Stoppat. Tryck på Meny för nästa åtgärd.", WordQuestion: "Hur säger man %s?", ChooseAnswer: "Välj rätt svar:", Correct: "Rätt!", TryAgain: "Inte riktigt. Försök igen:", NextWord: "Nästa ord", NextPage: "Nästa ▶", AlreadyLearned: "Det här ordet finns redan i ditt ordförråd.", AddedToVocab: "Ordet lades till i inlärt ordförråd.", TotalLearned: "Inlärda ord: %d", WriteWord: "Skriv %s:", Hint: "Tips", GoodSpelling: "Mycket bra stavning!", VoiceToText: "Röst till text", ImageTranslate: "Översätt bildtext", Translator: "Översättare", GPTAgent: "GPT-agent", VoicePrompt: "Skicka ett röstmeddelande så gör jag det till text.", ImagePrompt: "Skicka en bild med text, jag läser och översätter.", TranslatorPrompt: "Välj språk och skicka text, röst eller foto.", VoiceModePrompt: "Skicka ett röstmeddelande. Jag transkriberar utan att starta övning.", ImageModePrompt: "Skicka en bild med text. Jag läser och översätter.", TranslatorModePrompt: "Översättaren är redo. Välj käll- och målspråk.", ImageOpenToolsPrompt: "Öppna Verktyg för att översätta text från bild.", ImageSinglePhotoPrompt: "Skicka bara en bild.", Transcript: "Transkript", Translation: "Översättning", SourceLanguage: "Källspråk", TargetLanguage: "Målspråk", Auto: "Auto", WebApp: "Webbapp", VoiceDiscuss: "Du kan diskutera detta med text, jag fortsätter övningen.", ImageDiscuss: "Du kan skicka ett nytt foto eller skriva text.", VoicePremium: "Röstövning finns i Premium.", ImagePremium: "Bildtextöversättning finns i Premium.", VoiceLimit: "Dagens röstgräns är slut.\n\nRöst: %d/%d", VoiceTooLong: "Röstmeddelandet är för långt.\n\nMax: %d sekunder\nDitt meddelande: %d sekunder", VoiceDownloadFailed: "Kunde inte ladda ner röstmeddelandet: %s", VoiceTranscribeFailed: "Kunde inte känna igen rösten: %s", VoiceTranslationFailed: "Kunde inte översätta transkriptet: %s", ImageFileMissing: "Bildfilen hittades inte. Skicka fotot igen.", ImageDownloadFailed: "Kunde inte ladda ner bilden: %s", ImageReadFailed: "Kunde inte läsa bilden: %s"},
	"ta": {Menu: "பட்டி", Stop: "நிறுத்து", Back: "பின்", MainMenu: "முதன்மை பட்டி", ChooseAction: "செயலைத் தேர்வு செய்க:", Learning: "கற்றல்", Words: "சொற்கள்", Progress: "முன்னேற்றம்", Settings: "அமைப்புகள்", Tools: "கருவிகள்", NewLesson: "புதிய பாடம்", Practice: "பயிற்சி", Listening: "கேட்பு", Level: "நிலை கண்டறி", LearnWords: "சொற்கள் கற்று", Review: "மீள்பார்வு", Spelling: "எழுத்துப்பிழை", Vocabulary: "சொல்வளம்", Mistakes: "பிழைகள்", Leaders: "முன்னணி", Limits: "வரம்புகள்", BotLanguage: "பாட் மொழி", LearningLanguage: "கற்கும் மொழி", Notifications: "அறிவிப்புகள்", Premium: "Premium", Referrals: "பரிந்துரைகள்", ChooseBotLang: "பாட் மொழியைத் தேர்வு செய்க", ChooseLearnLang: "நீங்கள் எந்த மொழியை கற்க விரும்புகிறீர்கள்?", ChooseTimezone: "நேர மண்டலத்தைத் தேர்வு செய்க", TimezoneHint: "உங்கள் நேரப்படி சுமார் 19:00 மணிக்கு தினசரி நினைவூட்டுவேன்.", Ready: "முடிந்தது.", UnknownButton: "அந்த பொத்தானை புரிந்துகொள்ளவில்லை. பட்டியை அழுத்தவும்.", Stopped: "நிறுத்தப்பட்டது. அடுத்த செயலுக்கு பட்டியை அழுத்தவும்.", WordQuestion: "%s எப்படி சொல்வது?", ChooseAnswer: "சரியான பதிலைத் தேர்வு செய்க:", Correct: "சரி!", TryAgain: "இன்னும் இல்லை. மீண்டும் முயற்சி செய்க:", NextWord: "அடுத்த சொல்", NextPage: "அடுத்து ▶", AlreadyLearned: "இந்த சொல் ஏற்கனவே உங்கள் சொல்வளத்தில் உள்ளது.", AddedToVocab: "சொல் கற்ற சொல்வளத்தில் சேர்க்கப்பட்டது.", TotalLearned: "கற்ற சொற்கள்: %d", WriteWord: "%s எழுதுக:", Hint: "குறிப்பு", GoodSpelling: "மிக நல்ல எழுத்து!", VoiceToText: "குரல் முதல் உரை", ImageTranslate: "பட உரை மொழிபெயர்ப்பு", Translator: "மொழிபெயர்ப்பாளர்", GPTAgent: "GPT முகவர்", VoicePrompt: "குரல் செய்தி அனுப்புங்கள், அதை உரையாக்குவேன்.", ImagePrompt: "உரை உள்ள படத்தை அனுப்புங்கள், வாசித்து மொழிபெயர்ப்பேன்.", TranslatorPrompt: "மொழிகளைத் தேர்ந்து உரை, குரல் அல்லது படம் அனுப்புங்கள்.", VoiceModePrompt: "குரல் செய்தி அனுப்புங்கள். பயிற்சி தொடங்காமல் உரையாக்குவேன்.", ImageModePrompt: "உரை உள்ள படத்தை அனுப்புங்கள். வாசித்து மொழிபெயர்ப்பேன்.", TranslatorModePrompt: "மொழிபெயர்ப்பாளர் தயார். மூல மற்றும் இலக்கு மொழியைத் தேர்வு செய்க.", ImageOpenToolsPrompt: "பட உரையை மொழிபெயர்க்க கருவிகளைத் திறக்கவும்.", ImageSinglePhotoPrompt: "ஒரே ஒரு படத்தை மட்டும் அனுப்பவும்.", Transcript: "உரை வடிவம்", Translation: "மொழிபெயர்ப்பு", SourceLanguage: "மூல மொழி", TargetLanguage: "இலக்கு மொழி", Auto: "தானியங்கு", WebApp: "வலை பயன்பாடு", VoiceDiscuss: "இதைக் குறித்து உரையில் பேசலாம், நான் பயிற்சியை தொடர்வேன்.", ImageDiscuss: "மற்றொரு படத்தை அனுப்பலாம் அல்லது உரை எழுதலாம்.", VoicePremium: "குரல் பயிற்சி Premium-ல் கிடைக்கும்.", ImagePremium: "பட உரை மொழிபெயர்ப்பு Premium-ல் கிடைக்கும்.", VoiceLimit: "இன்றைய குரல் வரம்பு முடிந்தது.\n\nகுரல்: %d/%d", VoiceTooLong: "குரல் செய்தி மிக நீளமாக உள்ளது.\n\nஅதிகபட்சம்: %d விநாடி\nஉங்கள் செய்தி: %d விநாடி", VoiceDownloadFailed: "குரல் செய்தியை பதிவிறக்க முடியவில்லை: %s", VoiceTranscribeFailed: "குரலை அறிய முடியவில்லை: %s", VoiceTranslationFailed: "உரை வடிவத்தை மொழிபெயர்க்க முடியவில்லை: %s", ImageFileMissing: "பட கோப்பு கிடைக்கவில்லை. மீண்டும் அனுப்பவும்.", ImageDownloadFailed: "படத்தை பதிவிறக்க முடியவில்லை: %s", ImageReadFailed: "படத்தை வாசிக்க முடியவில்லை: %s"},
	"te": {Menu: "మెను", Stop: "ఆపు", Back: "వెనక్కి", MainMenu: "ప్రధాన మెను", ChooseAction: "చర్యను ఎంచుకోండి:", Learning: "నేర్చుకోవడం", Words: "పదాలు", Progress: "పురోగతి", Settings: "సెట్టింగులు", Tools: "సాధనాలు", NewLesson: "కొత్త పాఠం", Practice: "అభ్యాసం", Listening: "వినికిడి", Level: "స్థాయి కనుగొను", LearnWords: "పదాలు నేర్చుకోండి", Review: "పునశ్చరణ", Spelling: "అక్షర దోషం", Vocabulary: "పదకోశం", Mistakes: "తప్పులు", Leaders: "నాయకులు", Limits: "పరిమితులు", BotLanguage: "బాట్ భాష", LearningLanguage: "నేర్చుకునే భాష", Notifications: "నోటిఫికేషన్లు", Premium: "Premium", Referrals: "రెఫరల్స్", ChooseBotLang: "బాట్ భాషను ఎంచుకోండి", ChooseLearnLang: "మీరు ఏ భాష నేర్చుకోవాలి?", ChooseTimezone: "సమయ మండలాన్ని ఎంచుకోండి", TimezoneHint: "మీ స్థానిక సమయం ప్రకారం సుమారు 19:00కి రోజూ గుర్తు చేస్తాను.", Ready: "పూర్తైంది.", UnknownButton: "ఆ బటన్ అర్థం కాలేదు. మెను నొక్కండి.", Stopped: "ఆపబడింది. తదుపరి చర్య కోసం మెను నొక్కండి.", WordQuestion: "%s ఎలా చెబుతారు?", ChooseAnswer: "సరైన సమాధానం ఎంచుకోండి:", Correct: "సరైంది!", TryAgain: "ఇంకా కాదు. మళ్లీ ప్రయత్నించండి:", NextWord: "తదుపరి పదం", NextPage: "తదుపరి ▶", AlreadyLearned: "ఈ పదం ఇప్పటికే మీ పదకోశంలో ఉంది.", AddedToVocab: "పదం నేర్చుకున్న పదకోశానికి జోడించబడింది.", TotalLearned: "నేర్చుకున్న పదాలు: %d", WriteWord: "%s రాయండి:", Hint: "సూచన", GoodSpelling: "చాలా మంచి వర్ణక్రమం!", VoiceToText: "గళం నుంచి పాఠ్యం", ImageTranslate: "చిత్ర పాఠ్య అనువాదం", Translator: "అనువాదకుడు", GPTAgent: "GPT ఏజెంట్", VoicePrompt: "వాయిస్ సందేశం పంపండి, దాన్ని పాఠ్యంగా మారుస్తాను.", ImagePrompt: "పాఠ్యం ఉన్న చిత్రం పంపండి, చదివి అనువదిస్తాను.", TranslatorPrompt: "భాషలు ఎంచుకుని పాఠ్యం, గళం లేదా ఫోటో పంపండి.", VoiceModePrompt: "వాయిస్ సందేశం పంపండి. అభ్యాసం ప్రారంభించకుండా పాఠ్యంగా మారుస్తాను.", ImageModePrompt: "పాఠ్యం ఉన్న చిత్రం పంపండి. చదివి అనువదిస్తాను.", TranslatorModePrompt: "అనువాదకుడు సిద్ధంగా ఉన్నాడు. మూలం మరియు లక్ష్య భాష ఎంచుకోండి.", ImageOpenToolsPrompt: "చిత్ర పాఠ్యాన్ని అనువదించడానికి సాధనాలను తెరవండి.", ImageSinglePhotoPrompt: "ఒక చిత్రమే పంపండి.", Transcript: "లిప్యంతరీకరణ", Translation: "అనువాదం", SourceLanguage: "మూల భాష", TargetLanguage: "లక్ష్య భాష", Auto: "ఆటో", WebApp: "వెబ్ యాప్", VoiceDiscuss: "దీనిపై పాఠ్యంగా మాట్లాడవచ్చు, నేను అభ్యాసం కొనసాగిస్తాను.", ImageDiscuss: "మరొక ఫోటో పంపవచ్చు లేదా పాఠ్యం రాయవచ్చు.", VoicePremium: "వాయిస్ అభ్యాసం Premium-లో అందుబాటులో ఉంది.", ImagePremium: "చిత్ర పాఠ్య అనువాదం Premium-లో అందుబాటులో ఉంది.", VoiceLimit: "ఈరోజు వాయిస్ పరిమితి ముగిసింది.\n\nవాయిస్: %d/%d", VoiceTooLong: "వాయిస్ సందేశం చాలా పొడవుగా ఉంది.\n\nగరిష్ఠం: %d సెకన్లు\nమీ సందేశం: %d సెకన్లు", VoiceDownloadFailed: "వాయిస్ సందేశం డౌన్‌లోడ్ కాలేదు: %s", VoiceTranscribeFailed: "వాయిస్ గుర్తించలేకపోయాను: %s", VoiceTranslationFailed: "లిప్యంతరీకరణను అనువదించలేకపోయాను: %s", ImageFileMissing: "చిత్ర ఫైల్ కనబడలేదు. మళ్లీ పంపండి.", ImageDownloadFailed: "చిత్రం డౌన్‌లోడ్ కాలేదు: %s", ImageReadFailed: "చిత్రం చదవలేకపోయాను: %s"},
	"th": {Menu: "เมนู", Stop: "หยุด", Back: "กลับ", MainMenu: "เมนูหลัก", ChooseAction: "เลือกการทำงาน:", Learning: "เรียน", Words: "คำศัพท์", Progress: "ความคืบหน้า", Settings: "ตั้งค่า", Tools: "เครื่องมือ", NewLesson: "บทเรียนใหม่", Practice: "ฝึก", Listening: "การฟัง", Level: "หาระดับ", LearnWords: "เรียนคำศัพท์", Review: "ทบทวน", Spelling: "การสะกด", Vocabulary: "คลังคำศัพท์", Mistakes: "ข้อผิดพลาด", Leaders: "อันดับ", Limits: "ขีดจำกัด", BotLanguage: "ภาษาบอต", LearningLanguage: "ภาษาที่เรียน", Notifications: "การแจ้งเตือน", Premium: "Premium", Referrals: "แนะนำเพื่อน", ChooseBotLang: "เลือกภาษาบอต", ChooseLearnLang: "คุณอยากเรียนภาษาอะไร?", ChooseTimezone: "เลือกเขตเวลา", TimezoneHint: "ฉันจะส่งเตือนทุกวันประมาณ 19:00 ตามเวลาของคุณ.", Ready: "เรียบร้อย.", UnknownButton: "ฉันไม่เข้าใจปุ่มนี้ กดเมนู.", Stopped: "หยุดแล้ว กดเมนูเพื่อเลือกสิ่งถัดไป.", WordQuestion: "%s พูดว่าอย่างไร?", ChooseAnswer: "เลือกคำตอบที่ถูกต้อง:", Correct: "ถูกต้อง!", TryAgain: "ยังไม่ใช่ ลองอีกครั้ง:", NextWord: "คำถัดไป", NextPage: "ถัดไป ▶", AlreadyLearned: "คำนี้อยู่ในคลังคำศัพท์ของคุณแล้ว.", AddedToVocab: "เพิ่มคำลงในคลังคำที่เรียนแล้ว.", TotalLearned: "คำที่เรียนแล้ว: %d", WriteWord: "เขียน %s:", Hint: "คำใบ้", GoodSpelling: "สะกดได้ดีมาก!", VoiceToText: "เสียงเป็นข้อความ", ImageTranslate: "แปลข้อความจากภาพ", Translator: "นักแปล", GPTAgent: "ตัวแทน GPT", VoicePrompt: "ส่งข้อความเสียง แล้วฉันจะแปลงเป็นข้อความ.", ImagePrompt: "ส่งภาพที่มีข้อความ ฉันจะอ่านและแปล.", TranslatorPrompt: "เลือกภาษา แล้วส่งข้อความ เสียง หรือภาพ.", VoiceModePrompt: "ส่งข้อความเสียง ฉันจะแปลงเป็นข้อความโดยไม่เริ่มฝึก.", ImageModePrompt: "ส่งภาพที่มีข้อความ ฉันจะอ่านและแปล.", TranslatorModePrompt: "นักแปลพร้อมแล้ว เลือกภาษาต้นทางและปลายทาง.", ImageOpenToolsPrompt: "เปิดเครื่องมือเพื่อแปลข้อความจากภาพ.", ImageSinglePhotoPrompt: "กรุณาส่งเพียงหนึ่งภาพ.", Transcript: "ถอดเสียง", Translation: "คำแปล", SourceLanguage: "ภาษาต้นทาง", TargetLanguage: "ภาษาเป้าหมาย", Auto: "อัตโนมัติ", WebApp: "เว็บแอป", VoiceDiscuss: "คุณคุยต่อด้วยข้อความได้ ฉันจะฝึกต่อ.", ImageDiscuss: "คุณส่งภาพใหม่หรือเขียนข้อความได้.", VoicePremium: "การฝึกเสียงมีใน Premium.", ImagePremium: "แปลข้อความจากภาพมีใน Premium.", VoiceLimit: "ขีดจำกัดเสียงวันนี้หมดแล้ว.\n\nเสียง: %d/%d", VoiceTooLong: "ข้อความเสียงยาวเกินไป.\n\nสูงสุด: %d วินาที\nข้อความของคุณ: %d วินาที", VoiceDownloadFailed: "ดาวน์โหลดเสียงไม่สำเร็จ: %s", VoiceTranscribeFailed: "จำเสียงไม่ได้: %s", VoiceTranslationFailed: "แปลข้อความถอดเสียงไม่สำเร็จ: %s", ImageFileMissing: "ไม่พบไฟล์ภาพ ลองส่งอีกครั้ง.", ImageDownloadFailed: "ดาวน์โหลดภาพไม่สำเร็จ: %s", ImageReadFailed: "อ่านภาพไม่สำเร็จ: %s"},
	"tl": {Menu: "Menu", Stop: "Itigil", Back: "Bumalik", MainMenu: "Pangunahing menu", ChooseAction: "Pumili ng aksyon:", Learning: "Pag-aaral", Words: "Mga salita", Progress: "Progreso", Settings: "Mga setting", Tools: "Mga tool", NewLesson: "Bagong aralin", Practice: "Pagsasanay", Listening: "Pakikinig", Level: "Hanapin ang antas", LearnWords: "Matuto ng salita", Review: "Balik-aral", Spelling: "Pagbaybay", Vocabulary: "Talasalitaan", Mistakes: "Mga mali", Leaders: "Ranggo", Limits: "Mga limitasyon", BotLanguage: "Wika ng bot", LearningLanguage: "Wikang pinag-aaralan", Notifications: "Mga abiso", Premium: "Premium", Referrals: "Mga referral", ChooseBotLang: "Piliin ang wika ng bot", ChooseLearnLang: "Anong wika ang gusto mong matutunan?", ChooseTimezone: "Piliin ang time zone", TimezoneHint: "Magpapadala ako ng paalala araw-araw bandang 19:00.", Ready: "Tapos na.", UnknownButton: "Hindi ko naintindihan ang button. Pindutin ang Menu.", Stopped: "Itinigil. Pindutin ang Menu para sa susunod.", WordQuestion: "Paano sabihin ang %s?", ChooseAnswer: "Piliin ang tamang sagot:", Correct: "Tama!", TryAgain: "Hindi pa. Subukan muli:", NextWord: "Susunod na salita", NextPage: "Susunod ▶", AlreadyLearned: "Nasa talasalitaan mo na ang salitang ito.", AddedToVocab: "Naidagdag ang salita sa natutunang talasalitaan.", TotalLearned: "Natutunang salita: %d", WriteWord: "Isulat ang %s:", Hint: "Pahiwatig", GoodSpelling: "Mahusay ang baybay!", VoiceToText: "Boses sa teksto", ImageTranslate: "Salin ng teksto sa larawan", Translator: "Tagasalin", GPTAgent: "GPT Agent", VoicePrompt: "Magpadala ng voice message at gagawin ko itong teksto.", ImagePrompt: "Magpadala ng larawang may teksto, babasahin at isasalin ko.", TranslatorPrompt: "Pumili ng mga wika at magpadala ng teksto, boses, o larawan.", VoiceModePrompt: "Magpadala ng voice message. Ita-transcribe ko ito nang hindi nagsisimula ng practice.", ImageModePrompt: "Magpadala ng larawang may teksto. Babasahin at isasalin ko.", TranslatorModePrompt: "Handa na ang tagasalin. Piliin ang source at target na wika.", ImageOpenToolsPrompt: "Buksan ang Mga tool para magsalin ng teksto sa larawan.", ImageSinglePhotoPrompt: "Isang larawan lang ang ipadala.", Transcript: "Transkripsyon", Translation: "Salin", SourceLanguage: "Source na wika", TargetLanguage: "Target na wika", Auto: "Auto", WebApp: "Web app", VoiceDiscuss: "Maaari mo itong pag-usapan sa teksto at magpapatuloy ako sa practice.", ImageDiscuss: "Maaari kang magpadala ng bagong larawan o magsulat ng teksto.", VoicePremium: "Available ang voice practice sa Premium.", ImagePremium: "Available ang image text translation sa Premium.", VoiceLimit: "Tapos na ang voice limit ngayong araw.\n\nBoses: %d/%d", VoiceTooLong: "Masyadong mahaba ang voice message.\n\nMaximum: %d segundo\nMensahe mo: %d segundo", VoiceDownloadFailed: "Hindi ma-download ang voice message: %s", VoiceTranscribeFailed: "Hindi makilala ang boses: %s", VoiceTranslationFailed: "Hindi maisalin ang transcript: %s", ImageFileMissing: "Hindi mahanap ang image file. Ipadala muli.", ImageDownloadFailed: "Hindi ma-download ang larawan: %s", ImageReadFailed: "Hindi mabasa ang larawan: %s"},
	"tr": {Menu: "Menü", Stop: "Dur", Back: "Geri", MainMenu: "Ana menü", ChooseAction: "Bir işlem seç:", Learning: "Öğrenme", Words: "Kelimeler", Progress: "İlerleme", Settings: "Ayarlar", Tools: "Araçlar", NewLesson: "Yeni ders", Practice: "Pratik", Listening: "Dinleme", Level: "Seviye bul", LearnWords: "Kelime öğren", Review: "Tekrar", Spelling: "Yazım", Vocabulary: "Kelime hazinesi", Mistakes: "Hatalar", Leaders: "Liderler", Limits: "Limitler", BotLanguage: "Bot dili", LearningLanguage: "Öğrenme dili", Notifications: "Bildirimler", Premium: "Premium", Referrals: "Referanslar", ChooseBotLang: "Bot dilini seç", ChooseLearnLang: "Hangi dili öğrenmek istiyorsun?", ChooseTimezone: "Saat dilimini seç", TimezoneHint: "Yerel saatine göre yaklaşık 19:00'da günlük hatırlatma göndereceğim.", Ready: "Hazır.", UnknownButton: "Bu düğmeyi anlamadım. Menüye bas.", Stopped: "Durduruldu. Sonraki işlem için Menüye bas.", WordQuestion: "%s nasıl söylenir?", ChooseAnswer: "Doğru cevabı seç:", Correct: "Doğru!", TryAgain: "Henüz değil. Tekrar dene:", NextWord: "Sonraki kelime", NextPage: "Sonraki ▶", AlreadyLearned: "Bu kelime zaten kelime hazinende.", AddedToVocab: "Kelime öğrenilen kelime hazinesine eklendi.", TotalLearned: "Öğrenilen kelimeler: %d", WriteWord: "%s yaz:", Hint: "İpucu", GoodSpelling: "Çok iyi yazım!", VoiceToText: "Sesten metne", ImageTranslate: "Görsel metin çevirisi", Translator: "Çevirmen", GPTAgent: "GPT Ajanı", VoicePrompt: "Sesli mesaj gönder, metne çevireyim.", ImagePrompt: "Metin içeren bir görsel gönder, okuyup çevireyim.", TranslatorPrompt: "Dilleri seç, sonra metin, ses veya fotoğraf gönder.", VoiceModePrompt: "Sesli mesaj gönder. Pratik başlatmadan metne çevireceğim.", ImageModePrompt: "Metin içeren görsel gönder. Okuyup çevireceğim.", TranslatorModePrompt: "Çevirmen hazır. Kaynak ve hedef dili seç.", ImageOpenToolsPrompt: "Görselden metin çevirmek için Araçlar'ı aç.", ImageSinglePhotoPrompt: "Lütfen yalnızca bir görsel gönder.", Transcript: "Döküm", Translation: "Çeviri", SourceLanguage: "Kaynak dil", TargetLanguage: "Hedef dil", Auto: "Otomatik", WebApp: "Web uygulaması", VoiceDiscuss: "Bunu metinle konuşabilirsin, pratiğe devam ederim.", ImageDiscuss: "Başka fotoğraf gönderebilir veya metin yazabilirsin.", VoicePremium: "Ses pratiği Premium'da kullanılabilir.", ImagePremium: "Görsel metin çevirisi Premium'da kullanılabilir.", VoiceLimit: "Bugünkü ses limiti doldu.\n\nSes: %d/%d", VoiceTooLong: "Sesli mesaj çok uzun.\n\nMaksimum: %d saniye\nMesajın: %d saniye", VoiceDownloadFailed: "Sesli mesaj indirilemedi: %s", VoiceTranscribeFailed: "Ses tanınamadı: %s", VoiceTranslationFailed: "Döküm çevrilemedi: %s", ImageFileMissing: "Görsel dosyası bulunamadı. Tekrar gönder.", ImageDownloadFailed: "Görsel indirilemedi: %s", ImageReadFailed: "Görsel okunamadı: %s"},
	"vi": {Menu: "Menu", Stop: "Dừng", Back: "Quay lại", MainMenu: "Menu chính", ChooseAction: "Chọn hành động:", Learning: "Học", Words: "Từ", Progress: "Tiến độ", Settings: "Cài đặt", Tools: "Công cụ", NewLesson: "Bài học mới", Practice: "Luyện tập", Listening: "Nghe", Level: "Tìm trình độ", LearnWords: "Học từ", Review: "Ôn tập", Spelling: "Chính tả", Vocabulary: "Từ vựng", Mistakes: "Lỗi sai", Leaders: "Bảng xếp hạng", Limits: "Giới hạn", BotLanguage: "Ngôn ngữ bot", LearningLanguage: "Ngôn ngữ học", Notifications: "Thông báo", Premium: "Premium", Referrals: "Giới thiệu", ChooseBotLang: "Chọn ngôn ngữ bot", ChooseLearnLang: "Bạn muốn học ngôn ngữ nào?", ChooseTimezone: "Chọn múi giờ", TimezoneHint: "Tôi sẽ gửi nhắc nhở hằng ngày khoảng 19:00 giờ địa phương.", Ready: "Xong.", UnknownButton: "Tôi không hiểu nút đó. Nhấn Menu.", Stopped: "Đã dừng. Nhấn Menu để chọn hành động tiếp theo.", WordQuestion: "Nói %s như thế nào?", ChooseAnswer: "Chọn đáp án đúng:", Correct: "Đúng!", TryAgain: "Chưa đúng. Thử lại:", NextWord: "Từ tiếp theo", NextPage: "Tiếp ▶", AlreadyLearned: "Từ này đã có trong từ vựng của bạn.", AddedToVocab: "Từ đã được thêm vào từ vựng đã học.", TotalLearned: "Từ đã học: %d", WriteWord: "Viết %s:", Hint: "Gợi ý", GoodSpelling: "Viết chính tả rất tốt!", VoiceToText: "Giọng nói thành văn bản", ImageTranslate: "Dịch chữ trong ảnh", Translator: "Trình dịch", GPTAgent: "GPT Agent", VoicePrompt: "Gửi tin nhắn thoại, tôi sẽ chuyển thành văn bản.", ImagePrompt: "Gửi ảnh có chữ, tôi sẽ đọc và dịch.", TranslatorPrompt: "Chọn ngôn ngữ rồi gửi văn bản, giọng nói hoặc ảnh.", VoiceModePrompt: "Gửi tin nhắn thoại. Tôi sẽ chép lại mà không bắt đầu luyện tập.", ImageModePrompt: "Gửi ảnh có chữ. Tôi sẽ đọc và dịch.", TranslatorModePrompt: "Trình dịch đã sẵn sàng. Chọn ngôn ngữ nguồn và đích.", ImageOpenToolsPrompt: "Mở Công cụ để dịch chữ trong ảnh.", ImageSinglePhotoPrompt: "Vui lòng chỉ gửi một ảnh.", Transcript: "Bản chép", Translation: "Bản dịch", SourceLanguage: "Ngôn ngữ nguồn", TargetLanguage: "Ngôn ngữ đích", Auto: "Tự động", WebApp: "Ứng dụng web", VoiceDiscuss: "Bạn có thể trao đổi bằng văn bản, tôi sẽ tiếp tục luyện tập.", ImageDiscuss: "Bạn có thể gửi ảnh khác hoặc viết văn bản.", VoicePremium: "Luyện giọng nói có trong Premium.", ImagePremium: "Dịch chữ trong ảnh có trong Premium.", VoiceLimit: "Hết giới hạn tin nhắn thoại hôm nay.\n\nGiọng nói: %d/%d", VoiceTooLong: "Tin nhắn thoại quá dài.\n\nTối đa: %d giây\nTin của bạn: %d giây", VoiceDownloadFailed: "Không tải được tin nhắn thoại: %s", VoiceTranscribeFailed: "Không nhận dạng được giọng nói: %s", VoiceTranslationFailed: "Không dịch được bản chép: %s", ImageFileMissing: "Không tìm thấy tệp ảnh. Hãy gửi lại.", ImageDownloadFailed: "Không tải được ảnh: %s", ImageReadFailed: "Không đọc được ảnh: %s"},
}

func englishToolUICopy() toolUICopy {
	return toolUICopy{
		VoiceToText:            "Voice to text",
		ImageTranslate:         "Image text translate",
		Translator:             "Translator",
		GPTAgent:               "GPT Agent",
		VoicePrompt:            "Send a voice message and I will turn it into text.",
		ImagePrompt:            "Send an image with text and I will read and translate it.",
		TranslatorPrompt:       "Choose languages, then send text, voice, or a photo with text.",
		VoiceModePrompt:        "Send a voice message. I will transcribe it to text and will not start practice.",
		ImageModePrompt:        "Send an image with text. I will read what is written there and translate it.",
		TranslatorModePrompt:   "Translator is ready. Choose source and target languages, then send text, voice, or a photo with text.",
		ImageOpenToolsPrompt:   "To translate text from an image, open 🧰 Tools and choose Image text translate.",
		ImageSinglePhotoPrompt: "Please send only 1 image. Open 🧰 Tools → Image text translate and send one photo without an album.",
		TranscriptLabel:        "Transcript",
		TranslationLabel:       "Translation",
		SourceLanguage:         "Source language",
		TargetLanguage:         "Target language",
		AutoDetect:             "Auto",
		WebApp:                 "Web app",
		VoiceDiscussPrompt:     "You can discuss this by text, and I will continue in practice mode.",
		ImageDiscussPrompt:     "You can send another photo, and I will process the new image. Or write by text, and I will continue in practice mode.",
		VoicePremiumRequired:   "Voice practice is available in Premium.\n\nFree currently works with text only.\nOpen 📋 Menu → Premium.",
		ImagePremiumRequired:   "Image text translation is available in Premium.\n\nOpen 📋 Menu → Premium.",
		VoiceLimitReached:      "Your voice-message limit for today is over.\n\nVoice messages: %d/%d\nTomorrow the limit will reset automatically.",
		VoiceTooLong:           "The voice message is too long.\n\nMaximum: %d seconds\nYour message: %d seconds",
		VoiceDownloadFailed:    "Could not download the voice message: %s",
		VoiceTranscribeFailed:  "Could not recognize the voice message: %s",
		VoiceTranslationFailed: "Could not translate the transcript: %s",
		ImageFileMissing:       "Could not find the image file. Try sending the photo again.",
		ImageDownloadFailed:    "Could not download the image: %s",
		ImageReadFailed:        "Could not read the image: %s",
	}
}

func generatedToolUICopy(code string) (toolUICopy, bool) {
	t, ok := compactTelegramI18n[normalizeInterfaceLanguage(code)]
	if !ok {
		return toolUICopy{}, false
	}
	return toolUICopy{
		VoiceToText:            t.VoiceToText,
		ImageTranslate:         t.ImageTranslate,
		Translator:             t.Translator,
		GPTAgent:               t.GPTAgent,
		VoicePrompt:            t.VoicePrompt,
		ImagePrompt:            t.ImagePrompt,
		TranslatorPrompt:       t.TranslatorPrompt,
		VoiceModePrompt:        t.VoiceModePrompt,
		ImageModePrompt:        t.ImageModePrompt,
		TranslatorModePrompt:   t.TranslatorModePrompt,
		ImageOpenToolsPrompt:   t.ImageOpenToolsPrompt,
		ImageSinglePhotoPrompt: t.ImageSinglePhotoPrompt,
		TranscriptLabel:        t.Transcript,
		TranslationLabel:       t.Translation,
		SourceLanguage:         t.SourceLanguage,
		TargetLanguage:         t.TargetLanguage,
		AutoDetect:             t.Auto,
		WebApp:                 t.WebApp,
		VoiceDiscussPrompt:     t.VoiceDiscuss,
		ImageDiscussPrompt:     t.ImageDiscuss,
		VoicePremiumRequired:   t.VoicePremium,
		ImagePremiumRequired:   t.ImagePremium,
		VoiceLimitReached:      t.VoiceLimit,
		VoiceTooLong:           t.VoiceTooLong,
		VoiceDownloadFailed:    t.VoiceDownloadFailed,
		VoiceTranscribeFailed:  t.VoiceTranscribeFailed,
		VoiceTranslationFailed: t.VoiceTranslationFailed,
		ImageFileMissing:       t.ImageFileMissing,
		ImageDownloadFailed:    t.ImageDownloadFailed,
		ImageReadFailed:        t.ImageReadFailed,
	}, true
}

func generatedUICopy(code string) (uiCopy, bool) {
	t, ok := compactTelegramI18n[normalizeInterfaceLanguage(code)]
	if !ok {
		return uiCopy{}, false
	}
	menuButton := "📋 " + t.Menu
	stopButton := "⏹ " + t.Stop
	backButton := "◀ " + t.Back
	copy := uiCopy{
		MenuButton:       menuButton,
		StopButton:       stopButton,
		Back:             backButton,
		BackMenu:         backButton + " " + t.Menu,
		MainMenuTitle:    t.MainMenu,
		MainMenuBody:     t.ChooseAction,
		Learning:         t.Learning,
		Words:            t.Words,
		Stats:            t.Progress,
		Settings:         t.Settings,
		Tools:            t.Tools,
		NewLesson:        t.NewLesson,
		Practice:         t.Practice,
		AITutor:          aiTutorButtonLabel(code),
		Shadowing:        t.Listening,
		LevelTest:        t.Level,
		LearnWords:       t.LearnWords,
		WordGame:         t.Review,
		Spelling:         t.Spelling,
		Vocabulary:       t.Vocabulary,
		Phrasebook:       phrasebookButtonLabel(code),
		Mistakes:         t.Mistakes,
		Progress:         t.Progress,
		Leaders:          t.Leaders,
		Limits:           t.Limits,
		BotLanguage:      t.BotLanguage,
		LearningLanguage: t.LearningLanguage,
		Notifications:    t.Notifications,
		Premium:          t.Premium,
		Referral:         t.Referrals,
		ChooseBotLang:    t.ChooseBotLang,
		ChooseLearnLang:  t.ChooseLearnLang,
		ChooseTimezone:   t.ChooseTimezone,
		TimezoneHint:     t.TimezoneHint,
		BotLangSet:       t.Ready + " " + t.BotLanguage + ": %s.",
		LearnLangSet:     t.Ready + " " + t.LearningLanguage + ": %s.",
		UnknownButton:    t.UnknownButton,
		Stopped:          t.Stopped,
		WordQuestion:     t.WordQuestion,
		ChooseAnswer:     t.ChooseAnswer,
		Correct:          t.Correct,
		TryAgain:         t.TryAgain,
		NextWord:         t.NextWord,
		NextPage:         t.NextPage,
		AlreadyLearned:   t.AlreadyLearned,
		AddedToVocab:     t.AddedToVocab,
		TotalLearned:     t.TotalLearned,
		WriteWord:        t.WriteWord,
		Hint:             t.Hint,
		GoodSpelling:     t.GoodSpelling,
	}
	if toolCopy, ok := generatedToolUICopy(code); ok {
		copy.Tool = toolCopy
	}
	return copy, true
}

var toolUICopies = map[string]toolUICopy{
	"ru": {
		VoiceToText: "Голос в текст", ImageTranslate: "Перевод с картинки", Translator: "Переводчик", GPTAgent: "GPT Агент",
		VoicePrompt: "Пришли голосовое сообщение, и я переведу его в текст.", ImagePrompt: "Пришли картинку с текстом, и я прочитаю ее и переведу.", TranslatorPrompt: "Выбери языки, затем отправь текст, голос или фото с текстом.",
		VoiceModePrompt:        "Пришли голосовое сообщение. Я расшифрую его в текст и не буду запускать практику.",
		ImageModePrompt:        "Пришли картинку с текстом. Я прочитаю, что на ней написано, и переведу.",
		TranslatorModePrompt:   "Переводчик готов. Выбери язык исходного текста и язык перевода, затем отправь текст, голос или фото с текстом.",
		ImageOpenToolsPrompt:   "Чтобы перевести текст с картинки, открой 🧰 Инструменты и выбери «Перевод с картинки».",
		ImageSinglePhotoPrompt: "Нужно отправить только 1 картинку. Открой 🧰 Инструменты → Перевод с картинки и пришли одно фото без альбома.",
		TranscriptLabel:        "Расшифровка", TranslationLabel: "Перевод",
		SourceLanguage: "Язык исходного текста", TargetLanguage: "Язык перевода", AutoDetect: "Авто", WebApp: "Веб-версия",
		VoiceDiscussPrompt:   "Можешь обсудить это текстом, и я продолжу в режиме практики.",
		ImageDiscussPrompt:   "Можешь прислать ещё фото — я обработаю новое изображение. Или напиши текстом, и я продолжу обсуждение в режиме практики.",
		VoicePremiumRequired: "Голосовая практика доступна в Premium.\n\nFree пока работает только с текстом.\nНажми 📋 Меню → Premium.",
		ImagePremiumRequired: "Перевод с картинки доступен в Premium.\n\nНажми 📋 Меню → Premium.",
		VoiceLimitReached:    "Лимит голосовых на сегодня закончился.\n\nГолосовые: %d/%d\nЗавтра лимит обновится автоматически.",
		VoiceTooLong:         "Голосовое слишком длинное.\n\nМаксимум: %d секунд\nТвое сообщение: %d секунд",
		VoiceDownloadFailed:  "Не получилось скачать голосовое: %s", VoiceTranscribeFailed: "Не получилось распознать голос: %s", VoiceTranslationFailed: "Не получилось перевести расшифровку: %s",
		ImageFileMissing: "Не получилось найти файл картинки. Попробуй отправить фото ещё раз.", ImageDownloadFailed: "Не получилось скачать картинку: %s", ImageReadFailed: "Не получилось прочитать картинку: %s",
	},
	"es": {VoiceToText: "Voz a texto", ImageTranslate: "Traducir imagen", GPTAgent: "Agente GPT", VoiceModePrompt: "Envía un mensaje de voz. Lo transcribiré a texto y no iniciaré la práctica.", ImageModePrompt: "Envía una imagen con texto. Leeré lo que dice y lo traduciré.", TranscriptLabel: "Transcripción", TranslationLabel: "Traducción"},
	"de": {VoiceToText: "Sprache zu Text", ImageTranslate: "Bildtext übersetzen", GPTAgent: "GPT-Agent", VoiceModePrompt: "Sende eine Sprachnachricht. Ich transkribiere sie als Text und starte keine Übung.", ImageModePrompt: "Sende ein Bild mit Text. Ich lese den Text und übersetze ihn.", TranscriptLabel: "Transkript", TranslationLabel: "Übersetzung"},
	"fr": {VoiceToText: "Voix en texte", ImageTranslate: "Traduire une image", GPTAgent: "Agent GPT", VoiceModePrompt: "Envoie un message vocal. Je le transcrirai en texte sans lancer la pratique.", ImageModePrompt: "Envoie une image avec du texte. Je lirai le texte et le traduirai.", TranscriptLabel: "Transcription", TranslationLabel: "Traduction"},
	"it": {VoiceToText: "Voce in testo", ImageTranslate: "Traduci immagine", GPTAgent: "Agente GPT", VoiceModePrompt: "Invia un messaggio vocale. Lo trascriverò in testo senza avviare la pratica.", ImageModePrompt: "Invia un'immagine con testo. La leggerò e la tradurrò.", TranscriptLabel: "Trascrizione", TranslationLabel: "Traduzione"},
	"pl": {VoiceToText: "Głos na tekst", ImageTranslate: "Tłumacz obraz", GPTAgent: "Agent GPT", VoiceModePrompt: "Wyślij wiadomość głosową. Przepiszę ją na tekst i nie uruchomię praktyki.", ImageModePrompt: "Wyślij obraz z tekstem. Odczytam go i przetłumaczę.", TranscriptLabel: "Transkrypcja", TranslationLabel: "Tłumaczenie"},
	"pt": {VoiceToText: "Voz para texto", ImageTranslate: "Traduzir imagem", GPTAgent: "Agente GPT", VoiceModePrompt: "Envia uma mensagem de voz. Vou transcrevê-la para texto sem iniciar a prática.", ImageModePrompt: "Envia uma imagem com texto. Vou ler e traduzir.", TranscriptLabel: "Transcrição", TranslationLabel: "Tradução"},
	"ro": {VoiceToText: "Voce în text", ImageTranslate: "Traducere din imagine", GPTAgent: "Agent GPT", VoiceModePrompt: "Trimite un mesaj vocal. Îl transcriu în text și nu pornesc practica.", ImageModePrompt: "Trimite o imagine cu text. O citesc și o traduc.", TranscriptLabel: "Transcriere", TranslationLabel: "Traducere"},
	"uk": {VoiceToText: "Голос у текст", ImageTranslate: "Переклад з картинки", GPTAgent: "GPT Агент", VoiceModePrompt: "Надішли голосове повідомлення. Я перетворю його на текст і не запускатиму практику.", ImageModePrompt: "Надішли картинку з текстом. Я прочитаю її і перекладу.", TranscriptLabel: "Розшифровка", TranslationLabel: "Переклад"},
	"kk": {VoiceToText: "Дауысты мәтінге", ImageTranslate: "Суреттен аудару", GPTAgent: "GPT Агент", VoiceModePrompt: "Дауыс хабарламасын жібер. Мен оны мәтінге айналдырамын және практиканы бастамаймын.", ImageModePrompt: "Мәтіні бар сурет жібер. Мен оқып, аударамын.", TranscriptLabel: "Транскрипция", TranslationLabel: "Аударма"},
	"ky": {VoiceToText: "Үндү текстке", ImageTranslate: "Сүрөттөн которуу", GPTAgent: "GPT Агент", VoiceModePrompt: "Үн билдирүүсүн жибер. Мен аны текстке айландырам жана практиканы баштабайм.", ImageModePrompt: "Тексти бар сүрөт жибер. Мен окуп, которуп берем.", TranscriptLabel: "Текстке түшүрүү", TranslationLabel: "Котормо"},
	"ka": {VoiceToText: "ხმა ტექსტად", ImageTranslate: "სურათის თარგმნა", GPTAgent: "GPT აგენტი", VoiceModePrompt: "გამომიგზავნე ხმოვანი შეტყობინება. ტექსტად გადავაქცევ და პრაქტიკას არ დავიწყებ.", ImageModePrompt: "გამომიგზავნე სურათი ტექსტით. წავიკითხავ და ვთარგმნი.", TranscriptLabel: "ტრანსკრიფცია", TranslationLabel: "თარგმანი"},
	"uz": {VoiceToText: "Ovozdan matn", ImageTranslate: "Rasmdan tarjima", GPTAgent: "GPT Agent", VoiceModePrompt: "Ovozli xabar yuboring. Uni matnga aylantiraman va mashqni boshlamayman.", ImageModePrompt: "Matnli rasm yuboring. Uni o‘qib tarjima qilaman.", TranscriptLabel: "Matn", TranslationLabel: "Tarjima"},
	"tt": {VoiceToText: "Тавыштан текст", ImageTranslate: "Рәсемнән тәрҗемә", GPTAgent: "GPT Агент", VoiceModePrompt: "Тавыш хәбәрен җибәр. Мин аны текстка әйләндерәм һәм күнегүне башламыйм.", ImageModePrompt: "Текстлы рәсем җибәр. Мин аны укып тәрҗемә итәм.", TranscriptLabel: "Текстка күчерү", TranslationLabel: "Тәрҗемә"},
	"tg": {VoiceToText: "Овоз ба матн", ImageTranslate: "Тарҷума аз тасвир", GPTAgent: "GPT Агент", VoiceModePrompt: "Паёми овозӣ фиристед. Ман онро ба матн табдил медиҳам ва машқро оғоз намекунам.", ImageModePrompt: "Тасвир бо матн фиристед. Ман онро мехонам ва тарҷума мекунам.", TranscriptLabel: "Матн", TranslationLabel: "Тарҷума"},
	"hy": {VoiceToText: "Ձայնը տեքստի", ImageTranslate: "Թարգմանել նկարից", GPTAgent: "GPT գործակալ", VoiceModePrompt: "Ուղարկիր ձայնային հաղորդագրություն։ Ես այն կվերածեմ տեքստի և պրակտիկա չեմ սկսի։", ImageModePrompt: "Ուղարկիր տեքստով նկար։ Ես կկարդամ և կթարգմանեմ։", TranscriptLabel: "Տառադարձում", TranslationLabel: "Թարգմանություն"},
}

func toolUICopyFor(code string) toolUICopy {
	normalizedCode := normalizeInterfaceLanguage(code)
	copy := englishToolUICopy()
	if generated, ok := generatedToolUICopy(normalizedCode); ok {
		copy = mergeToolUICopy(copy, generated)
	}
	if localized, ok := toolUICopies[normalizedCode]; ok {
		copy = mergeToolUICopy(copy, localized)
	}
	if label := localizedTranslatorToolLabel(normalizedCode); label != "" {
		copy.Translator = label
	}
	return copy
}

func localizedTranslatorToolLabel(code string) string {
	switch normalizeInterfaceLanguage(code) {
	case "ru":
		return "Переводчик"
	case "es":
		return "Traductor"
	case "de":
		return "Übersetzer"
	case "fr":
		return "Traducteur"
	case "it":
		return "Traduttore"
	case "zh":
		return "翻译器"
	case "ja":
		return "翻訳"
	case "ko":
		return "번역기"
	case "tg":
		return "Тарҷумон"
	case "uz":
		return "Tarjimon"
	case "tt":
		return "Тәрҗемәче"
	case "hy":
		return "Թարգմանիչ"
	case "kk":
		return "Аудармашы"
	case "ky":
		return "Котормочу"
	case "ka":
		return "მთარგმნელი"
	case "uk":
		return "Перекладач"
	case "pl":
		return "Tłumacz"
	case "ro":
		return "Traducător"
	case "pt":
		return "Tradutor"
	default:
		return "Translator"
	}
}

func mergeToolUICopy(base toolUICopy, override toolUICopy) toolUICopy {
	if override.VoiceToText != "" {
		base.VoiceToText = override.VoiceToText
	}
	if override.ImageTranslate != "" {
		base.ImageTranslate = override.ImageTranslate
	}
	if override.Translator != "" {
		base.Translator = override.Translator
	}
	if override.GPTAgent != "" {
		base.GPTAgent = override.GPTAgent
	}
	if override.VoicePrompt != "" {
		base.VoicePrompt = override.VoicePrompt
	} else if override.VoiceModePrompt != "" {
		base.VoicePrompt = override.VoiceModePrompt
	}
	if override.ImagePrompt != "" {
		base.ImagePrompt = override.ImagePrompt
	} else if override.ImageModePrompt != "" {
		base.ImagePrompt = override.ImageModePrompt
	}
	if override.TranslatorPrompt != "" {
		base.TranslatorPrompt = override.TranslatorPrompt
	} else if override.TranslatorModePrompt != "" {
		base.TranslatorPrompt = override.TranslatorModePrompt
	}
	if override.VoiceModePrompt != "" {
		base.VoiceModePrompt = override.VoiceModePrompt
	}
	if override.ImageModePrompt != "" {
		base.ImageModePrompt = override.ImageModePrompt
	}
	if override.TranslatorModePrompt != "" {
		base.TranslatorModePrompt = override.TranslatorModePrompt
	}
	if override.ImageOpenToolsPrompt != "" {
		base.ImageOpenToolsPrompt = override.ImageOpenToolsPrompt
	}
	if override.ImageSinglePhotoPrompt != "" {
		base.ImageSinglePhotoPrompt = override.ImageSinglePhotoPrompt
	}
	if override.TranscriptLabel != "" {
		base.TranscriptLabel = override.TranscriptLabel
	}
	if override.TranslationLabel != "" {
		base.TranslationLabel = override.TranslationLabel
	}
	if override.SourceLanguage != "" {
		base.SourceLanguage = override.SourceLanguage
	}
	if override.TargetLanguage != "" {
		base.TargetLanguage = override.TargetLanguage
	}
	if override.AutoDetect != "" {
		base.AutoDetect = override.AutoDetect
	}
	if override.WebApp != "" {
		base.WebApp = override.WebApp
	}
	if override.VoiceDiscussPrompt != "" {
		base.VoiceDiscussPrompt = override.VoiceDiscussPrompt
	}
	if override.ImageDiscussPrompt != "" {
		base.ImageDiscussPrompt = override.ImageDiscussPrompt
	}
	if override.VoicePremiumRequired != "" {
		base.VoicePremiumRequired = override.VoicePremiumRequired
	}
	if override.ImagePremiumRequired != "" {
		base.ImagePremiumRequired = override.ImagePremiumRequired
	}
	if override.VoiceLimitReached != "" {
		base.VoiceLimitReached = override.VoiceLimitReached
	}
	if override.VoiceTooLong != "" {
		base.VoiceTooLong = override.VoiceTooLong
	}
	if override.VoiceDownloadFailed != "" {
		base.VoiceDownloadFailed = override.VoiceDownloadFailed
	}
	if override.VoiceTranscribeFailed != "" {
		base.VoiceTranscribeFailed = override.VoiceTranscribeFailed
	}
	if override.VoiceTranslationFailed != "" {
		base.VoiceTranslationFailed = override.VoiceTranslationFailed
	}
	if override.ImageFileMissing != "" {
		base.ImageFileMissing = override.ImageFileMissing
	}
	if override.ImageDownloadFailed != "" {
		base.ImageDownloadFailed = override.ImageDownloadFailed
	}
	if override.ImageReadFailed != "" {
		base.ImageReadFailed = override.ImageReadFailed
	}
	return base
}

func shadowingButtonLabel(code string) string {
	switch normalizeInterfaceLanguage(code) {
	case "ru":
		return "\u0410\u0443\u0434\u0438\u0440\u043e\u0432\u0430\u043d\u0438\u0435"
	case "es":
		return "Audici\u00f3n"
	case "de":
		return "H\u00f6rtraining"
	case "fr":
		return "\u00c9coute"
	case "it":
		return "Ascolto"
	case "zh":
		return "\u542c\u529b"
	case "ja":
		return "\u30ea\u30b9\u30cb\u30f3\u30b0"
	case "ko":
		return "\ub4e3\uae30"
	case "tg":
		return "\u0428\u0443\u043d\u0430\u0432\u043e\u04e3"
	case "uz":
		return "Tinglash"
	case "tt":
		return "\u0422\u044b\u04a3\u043b\u0430\u0443"
	case "hy":
		return "\u053c\u057d\u0578\u0582\u0574"
	case "kk":
		return "\u0422\u044b\u04a3\u0434\u0430\u043b\u044b\u043c"
	case "ky":
		return "\u0423\u0433\u0443\u0443"
	case "ka":
		return "\u10db\u10dd\u10e1\u10db\u10d4\u10dc\u10d0"
	case "uk":
		return "\u0410\u0443\u0434\u0456\u044e\u0432\u0430\u043d\u043d\u044f"
	case "pl":
		return "S\u0142uchanie"
	case "ro":
		return "Ascultare"
	case "pt":
		return "Audi\u00e7\u00e3o"
	default:
		return "Listening"
	}
}

func aiTutorButtonLabel(code string) string {
	switch normalizeInterfaceLanguage(code) {
	case "ru":
		return "AI Репетитор"
	case "es":
		return "Tutor IA"
	case "de":
		return "KI-Tutor"
	case "fr":
		return "Tuteur IA"
	case "it":
		return "Tutor IA"
	case "pt":
		return "Tutor IA"
	case "pl":
		return "Tutor AI"
	case "uk":
		return "AI-репетитор"
	case "zh":
		return "AI 导师"
	case "ja":
		return "AIチューター"
	case "ko":
		return "AI 튜터"
	default:
		return "AI Tutor"
	}
}

func withRuntimeUICopy(code string, copy uiCopy) uiCopy {
	if copy.NextPage == "" {
		copy.NextPage = englishUICopy().NextPage
	}
	if copy.AITutor == "" {
		copy.AITutor = aiTutorButtonLabel(code)
	}
	if copy.Shadowing == "" {
		copy.Shadowing = shadowingButtonLabel(code)
	}
	if copy.Phrasebook == "" {
		copy.Phrasebook = phrasebookButtonLabel(code)
	}
	if copy.Referral == "" {
		switch normalizeInterfaceLanguage(code) {
		case "ru":
			copy.Referral = "\u0420\u0435\u0444\u0435\u0440\u0430\u043b\u044b"
		default:
			copy.Referral = englishUICopy().Referral
		}
	}
	copy.Tool = toolUICopyFor(code)
	return copy
}

func phrasebookButtonLabel(code string) string {
	switch normalizeInterfaceLanguage(code) {
	case "ru":
		return "Разговорник"
	case "es":
		return "Frases"
	case "de":
		return "Sprachführer"
	case "fr":
		return "Phrases"
	case "it":
		return "Frasario"
	case "pl":
		return "Rozmówki"
	case "pt":
		return "Frases"
	case "ro":
		return "Fraze"
	case "uk":
		return "Розмовник"
	case "kk":
		return "Сөйлескіш"
	case "ky":
		return "Сүйлөшмө"
	case "ka":
		return "სასაუბრო"
	case "uz":
		return "So'zlashgich"
	case "tt":
		return "Сөйләшкеч"
	case "tg":
		return "Гуфтугӯнома"
	case "hy":
		return "Զրուցարան"
	case "zh":
		return "短语本"
	case "ja":
		return "フレーズ集"
	case "ko":
		return "표현집"
	case "ar":
		return "كتاب العبارات"
	case "bn":
		return "ফ্রেজবুক"
	case "cs":
		return "Konverzace"
	case "el":
		return "Φράσεις"
	case "hi":
		return "वाक्य-पुस्तिका"
	case "hu":
		return "Kifejezések"
	case "id":
		return "Buku frasa"
	case "nl":
		return "Zinnenboek"
	case "sv":
		return "Frasbok"
	case "ta":
		return "சொற்றொடர்கள்"
	case "te":
		return "పదబంధాలు"
	case "th":
		return "สมุดวลี"
	case "tl":
		return "Mga parirala"
	case "tr":
		return "Deyimler"
	case "vi":
		return "Sổ cụm từ"
	default:
		return "Phrasebook"
	}
}

func replaceButtonMentions(text string, oldMenuButton string, oldStopButton string, menuButton string, stopButton string) string {
	if oldMenuButton != "" && oldMenuButton != menuButton {
		text = strings.ReplaceAll(text, oldMenuButton, menuButton)
	}
	if oldStopButton != "" && oldStopButton != stopButton {
		text = strings.ReplaceAll(text, oldStopButton, stopButton)
	}
	return text
}

func replaceToolButtonMentions(copy toolUICopy, oldMenuButton string, oldStopButton string, menuButton string, stopButton string) toolUICopy {
	copy.VoiceToText = replaceButtonMentions(copy.VoiceToText, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImageTranslate = replaceButtonMentions(copy.ImageTranslate, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.Translator = replaceButtonMentions(copy.Translator, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.GPTAgent = replaceButtonMentions(copy.GPTAgent, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.VoicePrompt = replaceButtonMentions(copy.VoicePrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImagePrompt = replaceButtonMentions(copy.ImagePrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.TranslatorPrompt = replaceButtonMentions(copy.TranslatorPrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.VoiceModePrompt = replaceButtonMentions(copy.VoiceModePrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImageModePrompt = replaceButtonMentions(copy.ImageModePrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.TranslatorModePrompt = replaceButtonMentions(copy.TranslatorModePrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImageOpenToolsPrompt = replaceButtonMentions(copy.ImageOpenToolsPrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImageSinglePhotoPrompt = replaceButtonMentions(copy.ImageSinglePhotoPrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.TranscriptLabel = replaceButtonMentions(copy.TranscriptLabel, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.TranslationLabel = replaceButtonMentions(copy.TranslationLabel, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.SourceLanguage = replaceButtonMentions(copy.SourceLanguage, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.TargetLanguage = replaceButtonMentions(copy.TargetLanguage, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.AutoDetect = replaceButtonMentions(copy.AutoDetect, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.WebApp = replaceButtonMentions(copy.WebApp, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.VoiceDiscussPrompt = replaceButtonMentions(copy.VoiceDiscussPrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImageDiscussPrompt = replaceButtonMentions(copy.ImageDiscussPrompt, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.VoicePremiumRequired = replaceButtonMentions(copy.VoicePremiumRequired, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImagePremiumRequired = replaceButtonMentions(copy.ImagePremiumRequired, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.VoiceLimitReached = replaceButtonMentions(copy.VoiceLimitReached, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.VoiceTooLong = replaceButtonMentions(copy.VoiceTooLong, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.VoiceDownloadFailed = replaceButtonMentions(copy.VoiceDownloadFailed, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.VoiceTranscribeFailed = replaceButtonMentions(copy.VoiceTranscribeFailed, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.VoiceTranslationFailed = replaceButtonMentions(copy.VoiceTranslationFailed, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImageFileMissing = replaceButtonMentions(copy.ImageFileMissing, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImageDownloadFailed = replaceButtonMentions(copy.ImageDownloadFailed, oldMenuButton, oldStopButton, menuButton, stopButton)
	copy.ImageReadFailed = replaceButtonMentions(copy.ImageReadFailed, oldMenuButton, oldStopButton, menuButton, stopButton)
	return copy
}

func ui(user userState) uiCopy {
	code := normalizeInterfaceLanguage(user.InterfaceLanguage)
	if copy, ok := uiCopies[code]; ok {
		return withRuntimeUICopy(code, copy)
	}
	if copy, ok := generatedUICopy(code); ok {
		return withRuntimeUICopy(code, copy)
	}
	if alias := uiAliases[code]; alias != "" {
		if copy, ok := uiCopies[alias]; ok {
			return withRuntimeUICopy(alias, copy)
		}
	}
	return withRuntimeUICopy("en", englishUICopy())
}

func uiFromOptionalUser(users []userState) uiCopy {
	if len(users) == 0 {
		return englishUICopy()
	}
	return ui(users[0])
}

func (c uiCopy) format(template string, args ...any) string {
	return fmt.Sprintf(template, args...)
}
