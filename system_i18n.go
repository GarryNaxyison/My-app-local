package main

import "fmt"

type systemUICopy struct {
	LessonKind              string
	PracticeKind            string
	SpamWait                string
	UnknownInterfaceLang    string
	UnknownLearningLang     string
	UnknownTimezone         string
	LevelManualButton       string
	LevelStartButton        string
	LevelStartText          string
	ManualLevelText         string
	LevelQuestionText       string
	LevelCurrentScore       string
	LevelDontKnow           string
	LevelManualSet          string
	LevelTestSet            string
	LevelUnavailable        string
	LevelOpenFailed         string
	LevelUnknownChoice      string
	LevelUnknownAnswer      string
	LevelInactive           string
	LevelStale              string
	LessonAnswerInstruction string
	LessonDone              string
	PracticeStarted         string
	LimitReached            string
	ReminderTitle           string
	ReminderStatus          string
	ReminderEnabled         string
	ReminderDisabled        string
	ReminderTimezone        string
	ReminderDefault         string
	ReminderHour            string
	ReminderHint            string
	ReminderOnButton        string
	ReminderOffButton       string
	ReminderFooter          string
	ReminderMessages        []string
	VocabLevelQuestion      string
}

var systemUICopies = map[string]systemUICopy{
	"en": {
		LessonKind: "lessons", PracticeKind: "practice",
		SpamWait:             "Too many messages in a row. Wait %d seconds and try again.",
		UnknownInterfaceLang: "I did not understand the bot language. Please choose it again with /uilanguage.",
		UnknownLearningLang:  "I did not understand the learning language. Please choose it again with /language.",
		UnknownTimezone:      "I did not understand the time zone. Please choose it again with /timezone.",
		LevelManualButton:    "Choose my level", LevelStartButton: "Start test",
		LevelStartText:          "Level check: %s\n\nIf you already know your level, choose it manually. If you are not sure, take a short test: I will ask %d questions from A1 to C2, set your study level, and lessons with words will adapt to it.",
		ManualLevelText:         "Choose your study level.\n\nIf you are unsure, it is better to go back and take the test.",
		LevelQuestionText:       "Question %d/%d\n\n%s",
		LevelCurrentScore:       "Current score: %d",
		LevelDontKnow:           "I don't know",
		LevelManualSet:          "Done. I set your study level to %s.\n\nI will use it in lessons, practice, and new words. You can retake the level test anytime with /level.",
		LevelTestSet:            "Done. Your study level is %s.\n\nI will use it in lessons, practice, and new words. The word-learning section will start choosing words from this level and above.",
		LevelUnavailable:        "A separate level test for %s is not connected yet.\n\nFor now, choose the level manually, and AI will use it in lessons and practice.",
		LevelOpenFailed:         "Could not open the test question. Try starting the test again.",
		LevelUnknownChoice:      "I did not understand the selected level. Try choosing it again with /level.",
		LevelUnknownAnswer:      "I did not understand the test answer. Try starting the level check again.",
		LevelInactive:           "This test is no longer active. Start the level check again with /level.",
		LevelStale:              "This question is no longer current. Answer the latest test message.",
		LessonAnswerInstruction: "Reply with one message. After checking, the lesson will finish automatically.",
		LessonDone:              "Lesson complete. Press %s for the next step.",
		PracticeStarted:         "Practice mode is on. Write any phrase in %s. I will answer simply and explain in %s. To exit, press %s.",
		LimitReached:            "The limit for %s is over for today.\n\n%s\n\nTomorrow the limit will reset automatically.",
		ReminderTitle:           "Notification settings", ReminderStatus: "Status", ReminderEnabled: "enabled", ReminderDisabled: "disabled",
		ReminderTimezone: "Time zone", ReminderDefault: "default", ReminderHour: "Send time", ReminderHint: "You can enable or disable daily notifications and choose a time zone.",
		ReminderOnButton: "Enable notifications", ReminderOffButton: "Disable notifications",
		ReminderFooter: "Open %s and choose the next step. You can turn reminders off with /reminderoff.",
		ReminderMessages: []string{
			"A small language step today is better than a perfect plan for tomorrow. Open a lesson.",
			"Your vocabulary grows when you come back to it. Do one short exercise today.",
			"Five minutes of practice is enough to keep the rhythm alive.",
		},
		VocabLevelQuestion: "Choose the correct %s word for: %s",
	},
	"ru": {
		LessonKind: "уроки", PracticeKind: "практику",
		SpamWait:             "Слишком много сообщений подряд. Подожди %d секунд и попробуй снова.",
		UnknownInterfaceLang: "Не понял выбранный язык интерфейса. Попробуй выбрать его заново через /uilanguage.",
		UnknownLearningLang:  "Не понял выбранный язык. Попробуй выбрать его заново через /language.",
		UnknownTimezone:      "Не понял часовой пояс. Попробуй выбрать его ещё раз через /timezone.",
		LevelManualButton:    "Выбрать свой уровень", LevelStartButton: "Начать тест",
		LevelStartText:          "Определение уровня: %s\n\nЕсли ты уже знаешь свой уровень, можешь выбрать его сам. Если не уверен, пройди короткий тест: я задам %d вопросов от A1 до C2, выставлю учебный уровень, и уроки со словами будут подстраиваться под него.",
		ManualLevelText:         "Выбери свой учебный уровень.\n\nЕсли сомневаешься, лучше вернуться и пройти тест.",
		LevelQuestionText:       "Вопрос %d/%d\n\n%s",
		LevelCurrentScore:       "Текущий балл: %d",
		LevelDontKnow:           "Я не знаю",
		LevelManualSet:          "Готово! Я поставил учебный уровень: %s\n\nБуду учитывать его в уроках, практике и новых словах. Если захочешь перепроверить уровень тестом, он всегда доступен через /level.",
		LevelTestSet:            "Готово! Твой учебный уровень: %s\n\nЯ буду учитывать его в уроках, практике и новых словах. Раздел «Учить слова» начнет подбирать слова с этого уровня и выше.",
		LevelUnavailable:        "Для языка «%s» отдельный тест уровня ещё не подключён.\n\nПока можно выбрать уровень вручную, а ИИ будет учитывать его в уроках и практике.",
		LevelOpenFailed:         "Не получилось открыть вопрос теста. Попробуй начать тест заново.",
		LevelUnknownChoice:      "Не понял выбранный уровень. Попробуй выбрать его заново через /level.",
		LevelUnknownAnswer:      "Не понял ответ в тесте. Попробуй начать определение уровня заново.",
		LevelInactive:           "Этот тест уже не активен. Начни определение уровня заново через /level.",
		LevelStale:              "Этот вопрос уже не актуален. Отвечай на последнее сообщение теста.",
		LessonAnswerInstruction: "Ответь одним сообщением. После проверки урок закончится автоматически.",
		LessonDone:              "Урок завершен. Нажми %s для следующего шага.",
		PracticeStarted:         "Режим практики включен. Напиши любую фразу %s. Я отвечу просто и объясню на %s. Чтобы выйти, нажми %s.",
		LimitReached:            "Лимит на %s сегодня закончился.\n\n%s\n\nЗавтра лимит обновится автоматически.",
		ReminderTitle:           "Настройки уведомлений", ReminderStatus: "Статус", ReminderEnabled: "включены", ReminderDisabled: "выключены",
		ReminderTimezone: "Часовой пояс", ReminderDefault: "по умолчанию", ReminderHour: "Время рассылки", ReminderHint: "Можно включить или отключить ежедневные уведомления и выбрать часовой пояс.",
		ReminderOnButton: "Включить уведомления", ReminderOffButton: "Отключить уведомления",
		ReminderFooter: "Открой %s и выбери следующий шаг. Напоминания можно выключить командой /reminderoff.",
		ReminderMessages: []string{
			"Полиглотский ритуал начинается просто: 5 минут языка сегодня ценнее, чем большой план на завтра. Загляни в урок.",
			"Слова запоминают тех, кто возвращается. Открой тренировку и закрепи одно новое слово.",
			"Язык не требует героизма. Ему нужна маленькая ежедневная победа. Сделаем одну?",
			"Один короткий урок сегодня - ещё один кирпичик в твою башню полиглота.",
			"Пора размять языковую мышцу: урок, практика или повторяйка уже ждут тебя.",
			"Полиглоты растут не рывками, а серией маленьких касаний языка. Вернись на пару минут.",
			"Сегодня хороший день, чтобы добавить одно слово в личный словарь и забрать немного XP.",
			"Вернись к языку на 5 минут. Будущий ты скажет спасибо.",
			"Мини-квест дня: одно слово, один ответ, один маленький шаг к свободной речи.",
			"Язык любит регулярность. Давай сделаем короткий заход без давления и длинных обещаний.",
			"Твой словарь скучает. Пора дать ему новое слово.",
			"Даже 3 минуты практики держат мозг в языковом режиме.",
			"Полиглотский режим: открыть бота, сделать маленький шаг, забрать XP.",
			"Сегодня можно не учить много. Достаточно не потерять цепочку.",
			"Короткий языковой спринт: без давления, только прогресс.",
			"Одно слово в день - это уже 365 слов в год. Красиво же?",
			"Твой уровень ждет XP. Урок или повторяйка?",
			"Маленькое повторение сегодня спасает от большого забывания завтра.",
			"Вернись к языку: 5 минут, и серия продолжается.",
			"Пусть сегодня практика будет лёгкой: одно задание, один шаг, немного уверенности.",
			"Полиглот не тот, кто знает всё. Полиглот тот, кто возвращается к языку снова и снова.",
			"Время открыть словарь и поймать новое слово.",
			"Сегодняшняя миссия: сделать язык чуть ближе.",
			"Пара фраз в практике, и мозг снова вспоминает маршрут.",
			"Не жди настроения. Начни с одного слова, настроение подтянется.",
			"Язык зовет на короткую тренировку. Всего несколько минут.",
			"Уровни сами себя не прокачают. Заберём немного XP?",
			"Пять минут языка - это инвестиция в уверенность.",
			"Добавим в день маленький полиглотский штрих: слово, фраза или мини-игра.",
			"Закрепи язык сегодня, и завтра начнём новый круг сильнее.",
		},
		VocabLevelQuestion: "Выбери правильное слово на языке «%s» для: %s",
	},
}

var systemUICopyOverrides = map[string]systemUICopy{
	"es": {
		LessonKind: "lecciones", PracticeKind: "práctica", LevelManualButton: "Elegir mi nivel", LevelStartButton: "Empezar test", LevelDontKnow: "No lo sé",
		LevelStartText:    "Prueba de nivel: %s\n\nSi ya conoces tu nivel, puedes elegirlo manualmente. Si no estás seguro, haz un test corto: haré %d preguntas de A1 a C2, definiré tu nivel de estudio y las lecciones se adaptarán a él.",
		ManualLevelText:   "Elige tu nivel de estudio.\n\nSi tienes dudas, vuelve y haz el test.",
		LevelQuestionText: "Pregunta %d/%d\n\n%s", LevelCurrentScore: "Puntuación actual: %d",
		LevelManualSet:          "Listo. He puesto tu nivel de estudio: %s.\n\nLo tendré en cuenta en lecciones, práctica y palabras nuevas.",
		LevelTestSet:            "Listo. Tu nivel de estudio es %s.\n\nLo tendré en cuenta en lecciones, práctica y palabras nuevas.",
		LessonAnswerInstruction: "Responde con un solo mensaje. Después de revisarlo, la lección terminará automáticamente.",
		LessonDone:              "Lección terminada. Pulsa %s para el siguiente paso.",
		PracticeStarted:         "Modo práctica activado. Escribe cualquier frase en %s. Responderé de forma sencilla y explicaré en %s. Para salir, pulsa %s.",
		ReminderTitle:           "Ajustes de notificaciones", ReminderStatus: "Estado", ReminderEnabled: "activadas", ReminderDisabled: "desactivadas", ReminderTimezone: "Zona horaria", ReminderDefault: "por defecto", ReminderHour: "Hora de envío", ReminderHint: "Puedes activar o desactivar las notificaciones diarias y elegir una zona horaria.", ReminderOnButton: "Activar notificaciones", ReminderOffButton: "Desactivar notificaciones", ReminderFooter: "Abre %s y elige el siguiente paso. Puedes desactivar los recordatorios con /reminderoff.",
		ReminderMessages:   []string{"Un pequeño paso de idioma hoy vale más que un plan perfecto para mañana.", "Vuelve a practicar unos minutos y mantén el ritmo.", "Una palabra, una frase, una pequeña victoria."},
		VocabLevelQuestion: "Elige la palabra correcta en %s para: %s",
	},
	"de": {
		LessonKind: "Lektionen", PracticeKind: "Übung", LevelManualButton: "Niveau wählen", LevelStartButton: "Test starten", LevelDontKnow: "Ich weiß es nicht",
		LevelStartText:    "Einstufungstest: %s\n\nWenn du dein Niveau kennst, wähle es manuell. Wenn nicht, mache einen kurzen Test mit %d Fragen von A1 bis C2.",
		ManualLevelText:   "Wähle dein Lernniveau.\n\nWenn du unsicher bist, gehe zurück und mache den Test.",
		LevelQuestionText: "Frage %d/%d\n\n%s", LevelCurrentScore: "Aktuelle Punktzahl: %d",
		LevelManualSet:          "Fertig. Dein Lernniveau ist %s.\n\nIch berücksichtige es in Lektionen, Übungen und neuen Wörtern.",
		LevelTestSet:            "Fertig. Dein Lernniveau ist %s.\n\nIch berücksichtige es in Lektionen, Übungen und neuen Wörtern.",
		LessonAnswerInstruction: "Antworte mit einer Nachricht. Nach der Prüfung endet die Lektion automatisch.",
		LessonDone:              "Lektion abgeschlossen. Drücke %s für den nächsten Schritt.",
		PracticeStarted:         "Übungsmodus ist aktiv. Schreibe einen Satz auf %s. Ich antworte einfach und erkläre auf %s. Zum Beenden drücke %s.",
		ReminderTitle:           "Benachrichtigungseinstellungen", ReminderStatus: "Status", ReminderEnabled: "aktiviert", ReminderDisabled: "deaktiviert", ReminderTimezone: "Zeitzone", ReminderDefault: "Standard", ReminderHour: "Sendezeit", ReminderHint: "Du kannst tägliche Erinnerungen aktivieren oder deaktivieren und die Zeitzone wählen.", ReminderOnButton: "Benachrichtigungen aktivieren", ReminderOffButton: "Benachrichtigungen deaktivieren", ReminderFooter: "Öffne %s und wähle den nächsten Schritt. Erinnerungen kannst du mit /reminderoff ausschalten.",
		ReminderMessages:   []string{"Ein kleiner Sprachschritt heute ist besser als ein perfekter Plan für morgen.", "Übe ein paar Minuten und halte den Rhythmus.", "Ein Wort, ein Satz, ein kleiner Fortschritt."},
		VocabLevelQuestion: "Wähle das richtige %s-Wort für: %s",
	},
	"fr": {
		LessonKind: "leçons", PracticeKind: "pratique", LevelManualButton: "Choisir mon niveau", LevelStartButton: "Commencer le test", LevelDontKnow: "Je ne sais pas",
		LevelStartText:    "Test de niveau : %s\n\nSi tu connais déjà ton niveau, choisis-le manuellement. Sinon, fais un court test de %d questions de A1 à C2.",
		ManualLevelText:   "Choisis ton niveau d'étude.\n\nSi tu hésites, reviens et passe le test.",
		LevelQuestionText: "Question %d/%d\n\n%s", LevelCurrentScore: "Score actuel : %d",
		LevelManualSet:          "C'est prêt. Ton niveau d'étude est %s.\n\nJe l'utiliserai dans les leçons, la pratique et les nouveaux mots.",
		LevelTestSet:            "C'est prêt. Ton niveau d'étude est %s.\n\nJe l'utiliserai dans les leçons, la pratique et les nouveaux mots.",
		LessonAnswerInstruction: "Réponds avec un seul message. Après la correction, la leçon se terminera automatiquement.",
		LessonDone:              "Leçon terminée. Appuie sur %s pour la suite.",
		PracticeStarted:         "Mode pratique activé. Écris une phrase en %s. Je répondrai simplement et j'expliquerai en %s. Pour sortir, appuie sur %s.",
		ReminderTitle:           "Réglages des notifications", ReminderStatus: "Statut", ReminderEnabled: "activées", ReminderDisabled: "désactivées", ReminderTimezone: "Fuseau horaire", ReminderDefault: "par défaut", ReminderHour: "Heure d'envoi", ReminderHint: "Tu peux activer ou désactiver les rappels quotidiens et choisir un fuseau horaire.", ReminderOnButton: "Activer les notifications", ReminderOffButton: "Désactiver les notifications", ReminderFooter: "Ouvre %s et choisis la suite. Tu peux désactiver les rappels avec /reminderoff.",
		ReminderMessages:   []string{"Un petit pas linguistique aujourd'hui vaut mieux qu'un plan parfait pour demain.", "Reviens pratiquer quelques minutes et garde le rythme.", "Un mot, une phrase, une petite victoire."},
		VocabLevelQuestion: "Choisis le bon mot en %s pour : %s",
	},
	"it": {
		LessonKind: "lezioni", PracticeKind: "pratica", LevelManualButton: "Scegli il mio livello", LevelStartButton: "Inizia test", LevelDontKnow: "Non lo so",
		LevelStartText:    "Test di livello: %s\n\nSe conosci già il tuo livello, sceglilo manualmente. Se non sei sicuro, fai un breve test con %d domande da A1 a C2.",
		ManualLevelText:   "Scegli il tuo livello di studio.\n\nSe hai dubbi, torna indietro e fai il test.",
		LevelQuestionText: "Domanda %d/%d\n\n%s", LevelCurrentScore: "Punteggio attuale: %d",
		LevelManualSet:          "Fatto. Ho impostato il tuo livello: %s.\n\nLo userò in lezioni, pratica e nuove parole.",
		LevelTestSet:            "Fatto. Il tuo livello è %s.\n\nLo userò in lezioni, pratica e nuove parole.",
		LessonAnswerInstruction: "Rispondi con un solo messaggio. Dopo la correzione, la lezione finirà automaticamente.",
		LessonDone:              "Lezione completata. Premi %s per il prossimo passo.",
		PracticeStarted:         "Modalità pratica attiva. Scrivi una frase in %s. Risponderò in modo semplice e spiegherò in %s. Per uscire, premi %s.",
		ReminderTitle:           "Impostazioni notifiche", ReminderStatus: "Stato", ReminderEnabled: "attive", ReminderDisabled: "disattivate", ReminderTimezone: "Fuso orario", ReminderDefault: "predefinito", ReminderHour: "Ora di invio", ReminderHint: "Puoi attivare o disattivare i promemoria giornalieri e scegliere il fuso orario.", ReminderOnButton: "Attiva notifiche", ReminderOffButton: "Disattiva notifiche", ReminderFooter: "Apri %s e scegli il prossimo passo. Puoi disattivare i promemoria con /reminderoff.",
		ReminderMessages:   []string{"Un piccolo passo linguistico oggi vale più di un piano perfetto per domani.", "Torna a praticare per qualche minuto e mantieni il ritmo.", "Una parola, una frase, una piccola vittoria."},
		VocabLevelQuestion: "Scegli la parola corretta in %s per: %s",
	},
}

func systemUI(user userState) systemUICopy {
	code := normalizeInterfaceLanguage(user.InterfaceLanguage)
	copy := systemUICopies["en"]
	if code == "ru" {
		copy = systemUICopies["ru"]
	}
	if code != "ru" && code != "en" {
		copy = mergeSystemUICopy(copy, compactLocalizedSystemFallback(code))
	}
	if override, ok := systemUICopyOverrides[code]; ok {
		copy = mergeSystemUICopy(copy, override)
	}
	if override, ok := cleanReminderSystemUICopyOverrides[code]; ok {
		copy = mergeSystemUICopy(copy, override)
	}
	return copy
}

func mergeSystemUICopy(base, override systemUICopy) systemUICopy {
	if override.LessonKind != "" {
		base.LessonKind = override.LessonKind
	}
	if override.PracticeKind != "" {
		base.PracticeKind = override.PracticeKind
	}
	if override.SpamWait != "" {
		base.SpamWait = override.SpamWait
	}
	if override.UnknownInterfaceLang != "" {
		base.UnknownInterfaceLang = override.UnknownInterfaceLang
	}
	if override.UnknownLearningLang != "" {
		base.UnknownLearningLang = override.UnknownLearningLang
	}
	if override.UnknownTimezone != "" {
		base.UnknownTimezone = override.UnknownTimezone
	}
	if override.LevelManualButton != "" {
		base.LevelManualButton = override.LevelManualButton
	}
	if override.LevelStartButton != "" {
		base.LevelStartButton = override.LevelStartButton
	}
	if override.LevelStartText != "" {
		base.LevelStartText = override.LevelStartText
	}
	if override.ManualLevelText != "" {
		base.ManualLevelText = override.ManualLevelText
	}
	if override.LevelQuestionText != "" {
		base.LevelQuestionText = override.LevelQuestionText
	}
	if override.LevelCurrentScore != "" {
		base.LevelCurrentScore = override.LevelCurrentScore
	}
	if override.LevelDontKnow != "" {
		base.LevelDontKnow = override.LevelDontKnow
	}
	if override.LevelManualSet != "" {
		base.LevelManualSet = override.LevelManualSet
	}
	if override.LevelTestSet != "" {
		base.LevelTestSet = override.LevelTestSet
	}
	if override.LevelUnavailable != "" {
		base.LevelUnavailable = override.LevelUnavailable
	}
	if override.LevelOpenFailed != "" {
		base.LevelOpenFailed = override.LevelOpenFailed
	}
	if override.LevelUnknownChoice != "" {
		base.LevelUnknownChoice = override.LevelUnknownChoice
	}
	if override.LevelUnknownAnswer != "" {
		base.LevelUnknownAnswer = override.LevelUnknownAnswer
	}
	if override.LevelInactive != "" {
		base.LevelInactive = override.LevelInactive
	}
	if override.LevelStale != "" {
		base.LevelStale = override.LevelStale
	}
	if override.LessonAnswerInstruction != "" {
		base.LessonAnswerInstruction = override.LessonAnswerInstruction
	}
	if override.LessonDone != "" {
		base.LessonDone = override.LessonDone
	}
	if override.PracticeStarted != "" {
		base.PracticeStarted = override.PracticeStarted
	}
	if override.LimitReached != "" {
		base.LimitReached = override.LimitReached
	}
	if override.ReminderTitle != "" {
		base.ReminderTitle = override.ReminderTitle
	}
	if override.ReminderStatus != "" {
		base.ReminderStatus = override.ReminderStatus
	}
	if override.ReminderEnabled != "" {
		base.ReminderEnabled = override.ReminderEnabled
	}
	if override.ReminderDisabled != "" {
		base.ReminderDisabled = override.ReminderDisabled
	}
	if override.ReminderTimezone != "" {
		base.ReminderTimezone = override.ReminderTimezone
	}
	if override.ReminderDefault != "" {
		base.ReminderDefault = override.ReminderDefault
	}
	if override.ReminderHour != "" {
		base.ReminderHour = override.ReminderHour
	}
	if override.ReminderHint != "" {
		base.ReminderHint = override.ReminderHint
	}
	if override.ReminderOnButton != "" {
		base.ReminderOnButton = override.ReminderOnButton
	}
	if override.ReminderOffButton != "" {
		base.ReminderOffButton = override.ReminderOffButton
	}
	if override.ReminderFooter != "" {
		base.ReminderFooter = override.ReminderFooter
	}
	if len(override.ReminderMessages) > 0 {
		base.ReminderMessages = override.ReminderMessages
	}
	if override.VocabLevelQuestion != "" {
		base.VocabLevelQuestion = override.VocabLevelQuestion
	}
	return base
}

func compactLocalizedSystemFallback(code string) systemUICopy {
	copy := ui(userState{InterfaceLanguage: code})
	return systemUICopy{
		LessonKind:              copy.NewLesson,
		PracticeKind:            copy.Practice,
		SpamWait:                copy.UnknownButton + " %d.",
		UnknownInterfaceLang:    copy.ChooseBotLang,
		UnknownLearningLang:     copy.ChooseLearnLang,
		UnknownTimezone:         copy.ChooseTimezone,
		LevelManualButton:       copy.LevelTest,
		LevelStartButton:        copy.LevelTest,
		LevelStartText:          copy.LevelTest + ": %s\n\n" + copy.MainMenuBody + " %d A1-C2.",
		ManualLevelText:         copy.LevelTest + ". " + copy.ChooseAnswer,
		LevelQuestionText:       copy.WordQuestion + "\n\n%d/%d\n\n%s",
		LevelCurrentScore:       copy.Progress + ": %d",
		LevelDontKnow:           copy.Hint,
		LevelManualSet:          copy.Correct + " " + copy.LevelTest + ": %s",
		LevelTestSet:            copy.Correct + " " + copy.LevelTest + ": %s",
		LevelUnavailable:        copy.LevelTest + ": %s. " + copy.ChooseAnswer,
		LevelOpenFailed:         copy.LevelTest + ". " + copy.TryAgain,
		LevelUnknownChoice:      copy.LevelTest + ". " + copy.TryAgain,
		LevelUnknownAnswer:      copy.LevelTest + ". " + copy.TryAgain,
		LevelInactive:           copy.LevelTest + ". " + copy.TryAgain,
		LevelStale:              copy.LevelTest + ". " + copy.NextWord,
		LessonAnswerInstruction: copy.ChooseAnswer,
		LessonDone:              copy.NewLesson + ". " + copy.MenuButton + ": %s",
		PracticeStarted:         copy.Practice + ": %s / %s. " + copy.StopButton + ": %s",
		LimitReached:            copy.Limits + ": %s\n\n%s",
		ReminderTitle:           copy.Notifications,
		ReminderStatus:          copy.Progress,
		ReminderEnabled:         copy.Correct,
		ReminderDisabled:        copy.Stopped,
		ReminderTimezone:        copy.ChooseTimezone,
		ReminderDefault:         copy.TimezoneHint,
		ReminderHour:            copy.ChooseTimezone,
		ReminderHint:            copy.TimezoneHint,
		ReminderOnButton:        copy.Notifications,
		ReminderOffButton:       copy.StopButton,
		ReminderFooter:          copy.MenuButton + ": %s",
		ReminderMessages:        []string{copy.Notifications + ": " + copy.NewLesson + " / " + copy.Practice},
		VocabLevelQuestion:      copy.ChooseAnswer + " %s: %s",
	}
}

func systemSprintf(template string, args ...any) string {
	return fmt.Sprintf(template, args...)
}
