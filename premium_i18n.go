package main

import "fmt"

type premiumUICopy struct {
	FreeVoiceUnavailable        string
	LessonsPerDay               string
	PracticeMessagesPerDay      string
	PremiumVoicesPerDay         string
	VoiceTextAndTranslation     string
	ImageTextTranslation        string
	VoicePhotoContextPractice   string
	BestValueLabel              string
	MaxAccessLabel              string
	PlatinumPriority            string
	MonthPrice                  string
	YearPrice                   string
	PlatinumMonthPrice          string
	PlatinumYearPrice           string
	InviteFree                  string
	InvoiceDescription          string
	PaymentUnrecognized         string
	PreCheckoutFailed           string
	ActivatedTitle              string
	ActivatedAvailable          string
	ActiveUntil                 string
	CheckLimits                 string
	YooKassaUnavailable         string
	YooKassaCreateFailed        string
	YooKassaPaymentText         string
	InviteCreateFailed          string
	InviteLinkText              string
	ReferralShareText           string
	ReferralInviter             string
	ReferralInviterLevel        string
	ReferralInvitee             string
	ReferralBalanceTitle        string
	ReferralInvitationsLabel    string
	ReferralBalanceHint         string
	ReferralWithdrawButton      string
	ReferralWithdrawUnavailable string
	LimitsTitle                 string
	PlanLabel                   string
	LessonsToday                string
	PracticeToday               string
	VoicesToday                 string
	MonthStarsButton            string
	MonthRubButton              string
	YearStarsButton             string
	YearRubButton               string
	PlanStarsButton             string
	PlanRubButton               string
	InviteFriendButton          string
	PaySBPButton                string
	PayStarsButton              string
	MonthTitle                  string
	YearTitle                   string
	PlatinumMonthTitle          string
	PlatinumYearTitle           string
	MonthDays                   string
	YearDays                    string
}

func englishPremiumUICopy() premiumUICopy {
	return premiumUICopy{
		FreeVoiceUnavailable:        "voice messages unavailable",
		LessonsPerDay:               "%d lessons per day",
		PracticeMessagesPerDay:      "%d practice messages per day",
		PremiumVoicesPerDay:         "%d AI audio actions per day up to %d seconds each",
		VoiceTextAndTranslation:     "voice to text and translation of what was said",
		ImageTextTranslation:        "text translation from images",
		VoicePhotoContextPractice:   "practice using voice or photo context",
		BestValueLabel:              "70% launch price",
		MaxAccessLabel:              "maximum access",
		PlatinumPriority:            "maximum daily limits for serious study sprints",
		MonthPrice:                  "30 days: %d ₽ or %d Stars",
		YearPrice:                   "Year: %d ₽ or %d Stars (about %d%% off)",
		PlatinumMonthPrice:          "Platinum 30 days: %d ₽ or %d Stars",
		PlatinumYearPrice:           "Platinum year: %d ₽ or %d Stars (about %d%% off)",
		InviteFree:                  "Invite friends: they get 7 days Premium, you get 7 days Premium when they reach level 3, plus referral balance from purchases",
		InvoiceDescription:          "Premium: %d lessons, %d practice messages, %d AI audio actions up to %d seconds, voice checks, TTS, and image text translation per day.",
		PaymentUnrecognized:         "Payment received, but I could not recognize the plan. Open %s -> Premium.",
		PreCheckoutFailed:           "Could not confirm the purchase.",
		ActivatedTitle:              "Premium activated for %s",
		ActivatedAvailable:          "Now available: %d lessons, %d practice messages and %d AI audio actions per day",
		ActiveUntil:                 "Premium is available until %s.",
		CheckLimits:                 "Open %s -> %s to check your remaining limits.",
		YooKassaUnavailable:         "YooKassa payment is not configured yet. Check YOOKASSA_SHOP_ID, YOOKASSA_SECRET_KEY and YOOKASSA_RETURN_URL in .env.",
		YooKassaCreateFailed:        "Could not create a YooKassa payment. Try again a little later.",
		YooKassaPaymentText:         "*%s* ⭐\n\nPrice: *%d ₽*\nPayment method: *SBP via YooKassa*\n\nAfter payment, Premium will turn on automatically.",
		InviteCreateFailed:          "Could not create the invite link: %s",
		InviteLinkText:              "Try Poliglot AI:\n%s\n\nPoliglot AI is a Telegram bot and web app for learning languages with AI: short lessons, conversation practice, vocabulary training, translator, voice-to-text, and photo text translation.\n\nOpen this link and you will receive 7 days of Premium in the bot.",
		ReferralShareText:           "Try Poliglot AI:\n%s\n\nPoliglot AI is a Telegram bot and web app for learning languages with AI: short lessons, conversation practice, vocabulary training, translator, voice-to-text, and photo text translation.\n\nOpen this link and you will receive 7 days of Premium in the bot.",
		ReferralInviter:             "A new learner joined through your link. You will receive 7 days of Premium when they reach level 3.",
		ReferralInviterLevel:        "Your invited learner reached level 3. Premium was extended by %d days, now until %s.",
		ReferralInvitee:             "You joined by invitation and received Premium for 7 days. Your friend will receive 7 days of Premium when you reach level 3.",
		ReferralBalanceTitle:        "Referral balance",
		ReferralInvitationsLabel:    "Invitations",
		ReferralBalanceHint:         "20% from direct referrals and 5% from their referrals. Withdrawal is available from 1000 ₽ or the USDT equivalent at the current USDT/RUB rate.",
		ReferralWithdrawButton:      "Withdraw",
		ReferralWithdrawUnavailable: "Withdrawal is available from 1000 ₽ or the USDT equivalent at the current USDT/RUB rate. Your balance: %s (%s).",
		LimitsTitle:                 "My limits",
		PlanLabel:                   "Plan",
		LessonsToday:                "Lessons today",
		PracticeToday:               "Practice today",
		VoicesToday:                 "Voice messages today",
		MonthStarsButton:            "30 days: %d Stars",
		MonthRubButton:              "30 days: %d ₽ via YooKassa",
		YearStarsButton:             "Year: %d Stars",
		YearRubButton:               "Year: %d ₽ via YooKassa",
		PlanStarsButton:             "%s: %d Stars",
		PlanRubButton:               "%s: %d ₽ via YooKassa",
		InviteFriendButton:          "Invite a friend",
		PaySBPButton:                "Pay via SBP",
		PayStarsButton:              "Pay %d Stars",
		MonthTitle:                  "Premium 30 days",
		YearTitle:                   "Premium for a year",
		PlatinumMonthTitle:          "Platinum 30 days",
		PlatinumYearTitle:           "Platinum for a year",
		MonthDays:                   "30 days",
		YearDays:                    "365 days",
	}
}

var premiumUICopyOverrides = map[string]premiumUICopy{
	"ru": {
		FreeVoiceUnavailable:        "голосовые недоступны",
		LessonsPerDay:               "%d уроков в день",
		PracticeMessagesPerDay:      "%d сообщений практики в день",
		PremiumVoicesPerDay:         "%d голосовых в день до %d секунд",
		VoiceTextAndTranslation:     "голос в текст и перевод услышанного",
		ImageTextTranslation:        "перевод текста с картинки",
		VoicePhotoContextPractice:   "практика по контексту голоса или фото",
		BestValueLabel:              "стартовая цена со скидкой 70%",
		MaxAccessLabel:              "максимальный доступ",
		PlatinumPriority:            "максимальные дневные лимиты для серьёзных учебных рывков",
		MonthPrice:                  "30 дней: %d ₽ или %d Stars",
		YearPrice:                   "Год: %d ₽ или %d Stars (скидка около %d%%)",
		PlatinumMonthPrice:          "Platinum 30 дней: %d ₽ или %d Stars",
		PlatinumYearPrice:           "Platinum год: %d ₽ или %d Stars (скидка около %d%%)",
		InviteFree:                  "Приглашай друзей: им 7 дней Premium, тебе 7 дней Premium, когда друг достигнет 3 уровня, и баланс с покупок",
		InvoiceDescription:          "Premium: %d уроков, %d сообщений практики, %d голосовых до %d секунд, голос в текст и перевод текста с картинки в день.",
		PaymentUnrecognized:         "Оплата получена, но тариф не распознан. Открой %s -> Premium.",
		PreCheckoutFailed:           "Не удалось подтвердить покупку.",
		ActivatedTitle:              "Premium активирован на %s",
		ActivatedAvailable:          "Теперь доступно: %d уроков, %d сообщений практики и %d голосовых в день",
		ActiveUntil:                 "Premium доступен до %s.",
		CheckLimits:                 "Открой %s -> %s, чтобы проверить остатки.",
		YooKassaUnavailable:         "Оплата через ЮKassa пока не настроена. Проверь YOOKASSA_SHOP_ID, YOOKASSA_SECRET_KEY и YOOKASSA_RETURN_URL в .env.",
		YooKassaCreateFailed:        "Не получилось создать платеж через ЮKassa. Попробуй чуть позже.",
		YooKassaPaymentText:         "*%s* ⭐\n\nЦена: *%d ₽*\nСпособ оплаты: *СБП через ЮKassa*\n\nПосле оплаты Premium включится автоматически.",
		InviteCreateFailed:          "Не получилось создать ссылку-приглашение: %s",
		InviteLinkText:              "Попробуй Poliglot AI:\n%s\n\nPoliglot AI - это Telegram-бот и веб-версия для изучения языков с AI: короткие уроки, разговорная практика, тренировка слов, переводчик, голос в текст и перевод текста с фото.\n\nОткрой эту ссылку и получи 7 дней Premium-подписки в боте.",
		ReferralShareText:           "Попробуй Poliglot AI:\n%s\n\nPoliglot AI - это Telegram-бот и веб-версия для изучения языков с AI: короткие уроки, разговорная практика, тренировка слов, переводчик, голос в текст и перевод текста с фото.\n\nОткрой эту ссылку и получи 7 дней Premium-подписки в боте.",
		ReferralInviter:             "По твоей ссылке пришел новый ученик. Ты получишь 7 дней Premium, когда он достигнет 3 уровня.",
		ReferralInviterLevel:        "Приглашенный ученик достиг 3 уровня. Premium продлен на %d дней, теперь до %s.",
		ReferralInvitee:             "Ты пришел по приглашению и получил Premium на 7 дней. Друг получит 7 дней Premium, когда ты достигнешь 3 уровня.",
		ReferralBalanceTitle:        "Реферальный баланс",
		ReferralInvitationsLabel:    "Приглашений",
		ReferralBalanceHint:         "20% с прямых приглашений и 5% со второго уровня. Вывод доступен от 1000 ₽ либо аналог в USDT по актуальному курсу USDT/RUB.",
		ReferralWithdrawButton:      "Вывод",
		ReferralWithdrawUnavailable: "Вывод доступен от 1000 ₽ либо аналог в USDT по актуальному курсу USDT/RUB. Твой баланс: %s (%s).",
		LimitsTitle:                 "Мои лимиты",
		PlanLabel:                   "Тариф",
		LessonsToday:                "Уроки сегодня",
		PracticeToday:               "Практика сегодня",
		VoicesToday:                 "Голосовые сегодня",
		MonthStarsButton:            "30 дней: %d Stars",
		MonthRubButton:              "30 дней: %d ₽ через ЮKassa",
		YearStarsButton:             "Год: %d Stars",
		YearRubButton:               "Год: %d ₽ через ЮKassa",
		PlanStarsButton:             "%s: %d Stars",
		PlanRubButton:               "%s: %d ₽ через ЮKassa",
		InviteFriendButton:          "Пригласить друга",
		PaySBPButton:                "Оплатить через СБП",
		PayStarsButton:              "Оплатить %d Stars",
		MonthTitle:                  "Premium 30 дней",
		YearTitle:                   "Premium на год",
		PlatinumMonthTitle:          "Platinum 30 дней",
		PlatinumYearTitle:           "Platinum на год",
		MonthDays:                   "30 дней",
		YearDays:                    "365 дней",
	},
	"es": {
		FreeVoiceUnavailable:      "mensajes de voz no disponibles",
		LessonsPerDay:             "%d lecciones al día",
		PracticeMessagesPerDay:    "%d mensajes de práctica al día",
		PremiumVoicesPerDay:       "%d mensajes de voz al día hasta %d segundos",
		VoiceTextAndTranslation:   "voz a texto y traducción de lo dicho",
		ImageTextTranslation:      "traducción de texto en imágenes",
		VoicePhotoContextPractice: "práctica con contexto de voz o foto",
		MonthPrice:                "30 días: %d ₽ o %d Stars",
		YearPrice:                 "Año: %d ₽ o %d Stars (aprox. %d%% de descuento)",
		InviteFree:                "Tu amigo recibe 7 días Premium; tú recibes 7 días cuando alcance el nivel 3",
		InvoiceDescription:        "Premium: %d lecciones, %d mensajes de práctica, %d mensajes de voz hasta %d segundos, voz a texto y traducción de imágenes al día.",
		PaymentUnrecognized:       "Pago recibido, pero no pude reconocer el plan. Abre %s -> Premium.",
		PreCheckoutFailed:         "No se pudo confirmar la compra.",
		ActivatedTitle:            "Premium activado por %s",
		ActivatedAvailable:        "Ahora tienes: %d lecciones, %d mensajes de práctica y %d mensajes de voz al día",
		ActiveUntil:               "Premium está disponible hasta %s.",
		CheckLimits:               "Abre %s -> %s para ver tus límites restantes.",
		YooKassaPaymentText:       "*%s* ⭐\n\nPrecio: *%d ₽*\nMétodo de pago: *SBP vía YooKassa*\n\nDespués del pago, Premium se activará automáticamente.",
		InviteLinkText:            "Invita a un amigo con este enlace:\n%s\n\nCuando un usuario nuevo abra el bot por primera vez con él, recibirás Premium por 7 días.",
		ReferralInviter:           "Un nuevo estudiante llegó por tu enlace. Recibirás 7 días de Premium cuando alcance el nivel 3.",
		ReferralInviterLevel:      "Tu estudiante invitado alcanzó el nivel 3. Premium se amplió %d días, ahora hasta %s.",
		ReferralInvitee:           "Entraste por invitación y recibiste Premium por 7 días. Tu amigo recibirá 7 días de Premium cuando alcances el nivel 3.",
		LimitsTitle:               "Mis límites",
		PlanLabel:                 "Plan",
		LessonsToday:              "Lecciones hoy",
		PracticeToday:             "Práctica hoy",
		VoicesToday:               "Voz hoy",
		MonthStarsButton:          "30 días: %d Stars",
		MonthRubButton:            "30 días: %d ₽ vía YooKassa",
		YearStarsButton:           "Año: %d Stars",
		YearRubButton:             "Año: %d ₽ vía YooKassa",
		InviteFriendButton:        "Invitar a un amigo",
		PaySBPButton:              "Pagar por SBP",
		PayStarsButton:            "Pagar %d Stars",
		MonthTitle:                "Premium 30 días",
		YearTitle:                 "Premium por un año",
		MonthDays:                 "30 días",
		YearDays:                  "365 días",
	},
	"de": {
		FreeVoiceUnavailable:      "Sprachnachrichten nicht verfügbar",
		LessonsPerDay:             "%d Lektionen pro Tag",
		PracticeMessagesPerDay:    "%d Übungsnachrichten pro Tag",
		PremiumVoicesPerDay:       "%d Sprachnachrichten pro Tag bis %d Sekunden",
		VoiceTextAndTranslation:   "Sprache zu Text und Übersetzung des Gesagten",
		ImageTextTranslation:      "Textübersetzung aus Bildern",
		VoicePhotoContextPractice: "Übung mit Sprach- oder Fotokontext",
		MonthPrice:                "30 Tage: %d ₽ oder %d Stars",
		YearPrice:                 "Jahr: %d ₽ oder %d Stars (ca. %d%% Rabatt)",
		InviteFree:                "Dein Freund erhält 7 Tage Premium; du erhältst 7 Tage, wenn er Level 3 erreicht",
		InvoiceDescription:        "Premium: %d Lektionen, %d Übungsnachrichten, %d Sprachnachrichten bis %d Sekunden, Sprache zu Text und Bildtextübersetzung pro Tag.",
		PaymentUnrecognized:       "Zahlung erhalten, aber der Tarif wurde nicht erkannt. Öffne %s -> Premium.",
		PreCheckoutFailed:         "Der Kauf konnte nicht bestätigt werden.",
		ActivatedTitle:            "Premium für %s aktiviert",
		ActivatedAvailable:        "Jetzt verfügbar: %d Lektionen, %d Übungsnachrichten und %d Sprachnachrichten pro Tag",
		ActiveUntil:               "Premium ist verfügbar bis %s.",
		CheckLimits:               "Öffne %s -> %s, um deine verbleibenden Limits zu prüfen.",
		YooKassaPaymentText:       "*%s* ⭐\n\nPreis: *%d ₽*\nZahlungsart: *SBP über YooKassa*\n\nNach der Zahlung wird Premium automatisch aktiviert.",
		InviteLinkText:            "Lade einen Freund mit diesem Link ein:\n%s\n\nWenn ein neuer Nutzer den Bot zum ersten Mal darüber öffnet, erhältst du Premium für 7 Tage.",
		ReferralInviter:           "Ein neuer Lernender kam über deinen Link. Du erhältst 7 Tage Premium, wenn er Level 3 erreicht.",
		ReferralInviterLevel:      "Dein eingeladener Lernender hat Level 3 erreicht. Premium wurde um %d Tage verlängert, jetzt bis %s.",
		ReferralInvitee:           "Du bist über eine Einladung gekommen und hast 7 Tage Premium erhalten. Dein Freund erhält 7 Tage Premium, wenn du Level 3 erreichst.",
		LimitsTitle:               "Meine Limits",
		PlanLabel:                 "Tarif",
		LessonsToday:              "Lektionen heute",
		PracticeToday:             "Übung heute",
		VoicesToday:               "Sprachnachrichten heute",
		MonthStarsButton:          "30 Tage: %d Stars",
		MonthRubButton:            "30 Tage: %d ₽ über YooKassa",
		YearStarsButton:           "Jahr: %d Stars",
		YearRubButton:             "Jahr: %d ₽ über YooKassa",
		InviteFriendButton:        "Freund einladen",
		PaySBPButton:              "Per SBP bezahlen",
		PayStarsButton:            "%d Stars bezahlen",
		MonthTitle:                "Premium 30 Tage",
		YearTitle:                 "Premium für ein Jahr",
		MonthDays:                 "30 Tage",
		YearDays:                  "365 Tage",
	},
	"fr": {
		FreeVoiceUnavailable:      "messages vocaux indisponibles",
		LessonsPerDay:             "%d leçons par jour",
		PracticeMessagesPerDay:    "%d messages de pratique par jour",
		PremiumVoicesPerDay:       "%d messages vocaux par jour jusqu'à %d secondes",
		VoiceTextAndTranslation:   "voix en texte et traduction de ce qui est dit",
		ImageTextTranslation:      "traduction du texte des images",
		VoicePhotoContextPractice: "pratique avec le contexte d'une voix ou d'une photo",
		MonthPrice:                "30 jours : %d ₽ ou %d Stars",
		YearPrice:                 "Année : %d ₽ ou %d Stars (environ %d%% de remise)",
		InviteFree:                "Ton ami reçoit 7 jours Premium ; tu reçois 7 jours quand il atteint le niveau 3",
		InvoiceDescription:        "Premium : %d leçons, %d messages de pratique, %d messages vocaux jusqu'à %d secondes, voix en texte et traduction d'images par jour.",
		PaymentUnrecognized:       "Paiement reçu, mais je n'ai pas reconnu le forfait. Ouvre %s -> Premium.",
		PreCheckoutFailed:         "Impossible de confirmer l'achat.",
		ActivatedTitle:            "Premium activé pour %s",
		ActivatedAvailable:        "Disponible maintenant : %d leçons, %d messages de pratique et %d messages vocaux par jour",
		ActiveUntil:               "Premium est disponible jusqu'au %s.",
		CheckLimits:               "Ouvre %s -> %s pour vérifier tes limites restantes.",
		YooKassaPaymentText:       "*%s* ⭐\n\nPrix : *%d ₽*\nMode de paiement : *SBP via YooKassa*\n\nAprès le paiement, Premium s'activera automatiquement.",
		InviteLinkText:            "Invite un ami avec ce lien :\n%s\n\nQuand un nouvel utilisateur ouvre le bot avec ce lien pour la première fois, tu reçois Premium pour 7 jours.",
		ReferralInviter:           "Un nouvel élève est arrivé par ton lien. Tu recevras 7 jours de Premium lorsqu'il atteindra le niveau 3.",
		ReferralInviterLevel:      "Ton élève invité a atteint le niveau 3. Premium a été prolongé de %d jours, maintenant jusqu'au %s.",
		ReferralInvitee:           "Tu es arrivé par invitation et tu as reçu 7 jours de Premium. Ton ami recevra 7 jours de Premium quand tu atteindras le niveau 3.",
		LimitsTitle:               "Mes limites",
		PlanLabel:                 "Forfait",
		LessonsToday:              "Leçons aujourd'hui",
		PracticeToday:             "Pratique aujourd'hui",
		VoicesToday:               "Voix aujourd'hui",
		MonthStarsButton:          "30 jours : %d Stars",
		MonthRubButton:            "30 jours : %d ₽ via YooKassa",
		YearStarsButton:           "Année : %d Stars",
		YearRubButton:             "Année : %d ₽ via YooKassa",
		InviteFriendButton:        "Inviter un ami",
		PaySBPButton:              "Payer par SBP",
		PayStarsButton:            "Payer %d Stars",
		MonthTitle:                "Premium 30 jours",
		YearTitle:                 "Premium pour un an",
		MonthDays:                 "30 jours",
		YearDays:                  "365 jours",
	},
	"it": {
		FreeVoiceUnavailable:      "messaggi vocali non disponibili",
		LessonsPerDay:             "%d lezioni al giorno",
		PracticeMessagesPerDay:    "%d messaggi di pratica al giorno",
		PremiumVoicesPerDay:       "%d messaggi vocali al giorno fino a %d secondi",
		VoiceTextAndTranslation:   "voce in testo e traduzione del parlato",
		ImageTextTranslation:      "traduzione del testo dalle immagini",
		VoicePhotoContextPractice: "pratica con contesto da voce o foto",
		MonthPrice:                "30 giorni: %d ₽ o %d Stars",
		YearPrice:                 "Anno: %d ₽ o %d Stars (circa %d%% di sconto)",
		InviteFree:                "Il tuo amico riceve 7 giorni Premium; tu ricevi 7 giorni quando raggiunge il livello 3",
		InvoiceDescription:        "Premium: %d lezioni, %d messaggi di pratica, %d messaggi vocali fino a %d secondi, voce in testo e traduzione da immagini al giorno.",
		PaymentUnrecognized:       "Pagamento ricevuto, ma non ho riconosciuto il piano. Apri %s -> Premium.",
		PreCheckoutFailed:         "Impossibile confermare l'acquisto.",
		ActivatedTitle:            "Premium attivato per %s",
		ActivatedAvailable:        "Ora disponibili: %d lezioni, %d messaggi di pratica e %d messaggi vocali al giorno",
		ActiveUntil:               "Premium è disponibile fino al %s.",
		CheckLimits:               "Apri %s -> %s per controllare i limiti rimasti.",
		YooKassaPaymentText:       "*%s* ⭐\n\nPrezzo: *%d ₽*\nMetodo di pagamento: *SBP tramite YooKassa*\n\nDopo il pagamento, Premium si attiverà automaticamente.",
		InviteLinkText:            "Invita un amico con questo link:\n%s\n\nQuando un nuovo utente apre il bot da lì per la prima volta, riceverai Premium per 7 giorni.",
		ReferralInviter:           "Un nuovo studente è arrivato dal tuo link. Riceverai 7 giorni di Premium quando raggiungerà il livello 3.",
		ReferralInviterLevel:      "Lo studente invitato ha raggiunto il livello 3. Premium è stato esteso di %d giorni, ora fino al %s.",
		ReferralInvitee:           "Sei arrivato tramite invito e hai ricevuto 7 giorni di Premium. Il tuo amico riceverà 7 giorni di Premium quando raggiungerai il livello 3.",
		LimitsTitle:               "I miei limiti",
		PlanLabel:                 "Piano",
		LessonsToday:              "Lezioni oggi",
		PracticeToday:             "Pratica oggi",
		VoicesToday:               "Vocali oggi",
		MonthStarsButton:          "30 giorni: %d Stars",
		MonthRubButton:            "30 giorni: %d ₽ tramite YooKassa",
		YearStarsButton:           "Anno: %d Stars",
		YearRubButton:             "Anno: %d ₽ tramite YooKassa",
		InviteFriendButton:        "Invita un amico",
		PaySBPButton:              "Paga con SBP",
		PayStarsButton:            "Paga %d Stars",
		MonthTitle:                "Premium 30 giorni",
		YearTitle:                 "Premium per un anno",
		MonthDays:                 "30 giorni",
		YearDays:                  "365 giorni",
	},
	"pt": {FreeVoiceUnavailable: "mensagens de voz indisponíveis", LessonsPerDay: "%d lições por dia", PracticeMessagesPerDay: "%d mensagens de prática por dia", PremiumVoicesPerDay: "%d mensagens de voz por dia até %d segundos", VoiceTextAndTranslation: "voz para texto e tradução do que foi dito", ImageTextTranslation: "tradução de texto em imagens", VoicePhotoContextPractice: "prática com contexto de voz ou foto", MonthPrice: "30 dias: %d ₽ ou %d Stars", YearPrice: "Ano: %d ₽ ou %d Stars (cerca de %d%% de desconto)", InviteFree: "Ganha 7 dias grátis convidando um amigo", InvoiceDescription: "Premium: %d lições, %d mensagens de prática, %d mensagens de voz até %d segundos, voz para texto e tradução de imagens por dia.", ActivatedTitle: "Premium ativado por %s", LimitsTitle: "Meus limites", PlanLabel: "Plano", LessonsToday: "Lições hoje", PracticeToday: "Prática hoje", VoicesToday: "Voz hoje", InviteFriendButton: "Convidar amigo", PaySBPButton: "Pagar por SBP", PayStarsButton: "Pagar %d Stars", MonthTitle: "Premium 30 dias", YearTitle: "Premium por um ano", MonthDays: "30 dias", YearDays: "365 dias"},
	"pl": {FreeVoiceUnavailable: "wiadomości głosowe niedostępne", LessonsPerDay: "%d lekcji dziennie", PracticeMessagesPerDay: "%d wiadomości praktyki dziennie", PremiumVoicesPerDay: "%d wiadomości głosowych dziennie do %d sekund", VoiceTextAndTranslation: "głos na tekst i tłumaczenie wypowiedzi", ImageTextTranslation: "tłumaczenie tekstu z obrazów", VoicePhotoContextPractice: "praktyka z kontekstem głosu lub zdjęcia", MonthPrice: "30 dni: %d ₽ albo %d Stars", YearPrice: "Rok: %d ₽ albo %d Stars (około %d%% zniżki)", InviteFree: "Zdobądź 7 dni gratis, zapraszając znajomego", InvoiceDescription: "Premium: %d lekcji, %d wiadomości praktyki, %d wiadomości głosowych do %d sekund, głos na tekst i tłumaczenie obrazów dziennie.", ActivatedTitle: "Premium aktywowane na %s", LimitsTitle: "Moje limity", PlanLabel: "Plan", LessonsToday: "Lekcje dziś", PracticeToday: "Praktyka dziś", VoicesToday: "Głos dziś", InviteFriendButton: "Zaproś znajomego", PaySBPButton: "Zapłać przez SBP", PayStarsButton: "Zapłać %d Stars", MonthTitle: "Premium 30 dni", YearTitle: "Premium na rok", MonthDays: "30 dni", YearDays: "365 dni"},
	"ro": {FreeVoiceUnavailable: "mesajele vocale nu sunt disponibile", LessonsPerDay: "%d lecții pe zi", PracticeMessagesPerDay: "%d mesaje de practică pe zi", PremiumVoicesPerDay: "%d mesaje vocale pe zi până la %d secunde", VoiceTextAndTranslation: "voce în text și traducerea celor spuse", ImageTextTranslation: "traducerea textului din imagini", VoicePhotoContextPractice: "practică folosind context vocal sau foto", MonthPrice: "30 de zile: %d ₽ sau %d Stars", YearPrice: "An: %d ₽ sau %d Stars (aprox. %d%% reducere)", InviteFree: "Primește 7 zile gratis invitând un prieten", InvoiceDescription: "Premium: %d lecții, %d mesaje de practică, %d mesaje vocale până la %d secunde, voce în text și traducere din imagini pe zi.", ActivatedTitle: "Premium activat pentru %s", LimitsTitle: "Limitele mele", PlanLabel: "Plan", LessonsToday: "Lecții azi", PracticeToday: "Practică azi", VoicesToday: "Voce azi", InviteFriendButton: "Invită un prieten", PaySBPButton: "Plătește prin SBP", PayStarsButton: "Plătește %d Stars", MonthTitle: "Premium 30 de zile", YearTitle: "Premium pentru un an", MonthDays: "30 de zile", YearDays: "365 de zile"},
	"uk": {FreeVoiceUnavailable: "голосові недоступні", LessonsPerDay: "%d уроків на день", PracticeMessagesPerDay: "%d повідомлень практики на день", PremiumVoicesPerDay: "%d голосових на день до %d секунд", VoiceTextAndTranslation: "голос у текст і переклад сказаного", ImageTextTranslation: "переклад тексту з картинки", VoicePhotoContextPractice: "практика за контекстом голосу або фото", MonthPrice: "30 днів: %d ₽ або %d Stars", YearPrice: "Рік: %d ₽ або %d Stars (знижка близько %d%%)", InviteFree: "Отримай 7 днів безкоштовно, запросивши друга", InvoiceDescription: "Premium: %d уроків, %d повідомлень практики, %d голосових до %d секунд, голос у текст і переклад з картинки на день.", ActivatedTitle: "Premium активовано на %s", LimitsTitle: "Мої ліміти", PlanLabel: "Тариф", LessonsToday: "Уроки сьогодні", PracticeToday: "Практика сьогодні", VoicesToday: "Голосові сьогодні", InviteFriendButton: "Запросити друга", PaySBPButton: "Оплатити через СБП", PayStarsButton: "Оплатити %d Stars", MonthTitle: "Premium 30 днів", YearTitle: "Premium на рік", MonthDays: "30 днів", YearDays: "365 днів"},
	"kk": {FreeVoiceUnavailable: "дауыс хабарламалары қолжетімсіз", LessonsPerDay: "күніне %d сабақ", PracticeMessagesPerDay: "күніне %d практика хабарламасы", PremiumVoicesPerDay: "күніне %d дауыс хабарламасы, %d секундқа дейін", VoiceTextAndTranslation: "дауысты мәтінге айналдыру және айтылғанды аудару", ImageTextTranslation: "суреттегі мәтінді аудару", VoicePhotoContextPractice: "дауыс немесе фото контекстімен практика", MonthPrice: "30 күн: %d ₽ немесе %d Stars", YearPrice: "Жыл: %d ₽ немесе %d Stars (шамамен %d%% жеңілдік)", InviteFree: "Досыңды шақырып, 7 күнді тегін ал", InvoiceDescription: "Premium: күніне %d сабақ, %d практика хабарламасы, %d дауыс хабарламасы %d секундқа дейін, дауыс мәтінге және суреттен аудару.", ActivatedTitle: "Premium %s мерзімге қосылды", LimitsTitle: "Менің лимиттерім", PlanLabel: "Тариф", LessonsToday: "Бүгінгі сабақтар", PracticeToday: "Бүгінгі практика", VoicesToday: "Бүгінгі дауыс", InviteFriendButton: "Дос шақыру", PaySBPButton: "СБП арқылы төлеу", PayStarsButton: "%d Stars төлеу", MonthTitle: "Premium 30 күн", YearTitle: "Premium бір жылға", MonthDays: "30 күн", YearDays: "365 күн"},
	"ky": {FreeVoiceUnavailable: "үн билдирүүлөр жеткиликсиз", LessonsPerDay: "күнүнө %d сабак", PracticeMessagesPerDay: "күнүнө %d практика билдирүүсү", PremiumVoicesPerDay: "күнүнө %d үн билдирүү, %d секундга чейин", VoiceTextAndTranslation: "үндү текстке айлантуу жана айтылганды которуу", ImageTextTranslation: "сүрөттөгү текстти которуу", VoicePhotoContextPractice: "үн же сүрөт контексти менен практика", MonthPrice: "30 күн: %d ₽ же %d Stars", YearPrice: "Жыл: %d ₽ же %d Stars (болжол менен %d%% арзандатуу)", InviteFree: "Досуңду чакырып 7 күн акысыз ал", InvoiceDescription: "Premium: күнүнө %d сабак, %d практика билдирүүсү, %d үн билдирүү %d секундга чейин, үн текстке жана сүрөттөн которуу.", ActivatedTitle: "Premium %s мөөнөткө кошулду", LimitsTitle: "Менин лимиттерим", PlanLabel: "Тариф", LessonsToday: "Бүгүнкү сабактар", PracticeToday: "Бүгүнкү практика", VoicesToday: "Бүгүнкү үн", InviteFriendButton: "Дос чакыруу", PaySBPButton: "СБП аркылуу төлөө", PayStarsButton: "%d Stars төлөө", MonthTitle: "Premium 30 күн", YearTitle: "Premium бир жылга", MonthDays: "30 күн", YearDays: "365 күн"},
	"ka": {FreeVoiceUnavailable: "ხმოვანი შეტყობინებები მიუწვდომელია", LessonsPerDay: "%d გაკვეთილი დღეში", PracticeMessagesPerDay: "%d პრაქტიკის შეტყობინება დღეში", PremiumVoicesPerDay: "%d ხმოვანი დღეში %d წამამდე", VoiceTextAndTranslation: "ხმის ტექსტად გადაყვანა და ნათქვამის თარგმნა", ImageTextTranslation: "სურათში ტექსტის თარგმნა", VoicePhotoContextPractice: "პრაქტიკა ხმის ან ფოტოს კონტექსტით", MonthPrice: "30 დღე: %d ₽ ან %d Stars", YearPrice: "წელი: %d ₽ ან %d Stars (დაახლოებით %d%% ფასდაკლება)", InviteFree: "მოიწვიე მეგობარი და მიიღე 7 დღე უფასოდ", InvoiceDescription: "Premium: დღეში %d გაკვეთილი, %d პრაქტიკის შეტყობინება, %d ხმოვანი %d წამამდე, ხმა ტექსტად და სურათის თარგმნა.", ActivatedTitle: "Premium გააქტიურდა %s-ით", LimitsTitle: "ჩემი ლიმიტები", PlanLabel: "გეგმა", LessonsToday: "დღეს გაკვეთილები", PracticeToday: "დღეს პრაქტიკა", VoicesToday: "დღეს ხმა", InviteFriendButton: "მეგობრის მოწვევა", PaySBPButton: "SBP-ით გადახდა", PayStarsButton: "%d Stars გადახდა", MonthTitle: "Premium 30 დღე", YearTitle: "Premium ერთი წლით", MonthDays: "30 დღე", YearDays: "365 დღე"},
	"uz": {FreeVoiceUnavailable: "ovozli xabarlar mavjud emas", LessonsPerDay: "kuniga %d dars", PracticeMessagesPerDay: "kuniga %d mashq xabari", PremiumVoicesPerDay: "kuniga %d ovozli xabar, %d soniyagacha", VoiceTextAndTranslation: "ovozni matnga aylantirish va aytilganni tarjima qilish", ImageTextTranslation: "rasmdagi matnni tarjima qilish", VoicePhotoContextPractice: "ovoz yoki rasm konteksti bilan mashq", MonthPrice: "30 kun: %d ₽ yoki %d Stars", YearPrice: "Yil: %d ₽ yoki %d Stars (taxminan %d%% chegirma)", InviteFree: "Do‘st taklif qilib 7 kun bepul oling", InvoiceDescription: "Premium: kuniga %d dars, %d mashq xabari, %d ovozli xabar %d soniyagacha, ovozni matnga va rasmdan tarjima.", ActivatedTitle: "Premium %s muddatga yoqildi", LimitsTitle: "Mening limitlarim", PlanLabel: "Tarif", LessonsToday: "Bugungi darslar", PracticeToday: "Bugungi mashq", VoicesToday: "Bugungi ovoz", InviteFriendButton: "Do‘st taklif qilish", PaySBPButton: "SBP orqali to‘lash", PayStarsButton: "%d Stars to‘lash", MonthTitle: "Premium 30 kun", YearTitle: "Premium bir yilga", MonthDays: "30 kun", YearDays: "365 kun"},
	"tt": {FreeVoiceUnavailable: "тавыш хәбәрләре юк", LessonsPerDay: "көненә %d дәрес", PracticeMessagesPerDay: "көненә %d күнегү хәбәре", PremiumVoicesPerDay: "көненә %d тавыш хәбәре, %d секундка кадәр", VoiceTextAndTranslation: "тавышны текстка әйләндерү һәм әйтелгәнне тәрҗемә итү", ImageTextTranslation: "рәсемдәге текстны тәрҗемә итү", VoicePhotoContextPractice: "тавыш яки фото контексты белән күнегү", MonthPrice: "30 көн: %d ₽ яки %d Stars", YearPrice: "Ел: %d ₽ яки %d Stars (якынча %d%% ташлама)", InviteFree: "Дустыңны чакырып 7 көн бушлай ал", InvoiceDescription: "Premium: көненә %d дәрес, %d күнегү хәбәре, %d тавыш хәбәре %d секундка кадәр, тавыштан текст һәм рәсемнән тәрҗемә.", ActivatedTitle: "Premium %s вакытка кушылды", LimitsTitle: "Минем лимитлар", PlanLabel: "Тариф", LessonsToday: "Бүген дәресләр", PracticeToday: "Бүген күнегү", VoicesToday: "Бүген тавыш", InviteFriendButton: "Дус чакыру", PaySBPButton: "СБП аша түләү", PayStarsButton: "%d Stars түләү", MonthTitle: "Premium 30 көн", YearTitle: "Premium бер елга", MonthDays: "30 көн", YearDays: "365 көн"},
	"tg": {FreeVoiceUnavailable: "паёмҳои овозӣ дастрас нестанд", LessonsPerDay: "дар як рӯз %d дарс", PracticeMessagesPerDay: "дар як рӯз %d паёми машқ", PremiumVoicesPerDay: "дар як рӯз %d паёми овозӣ то %d сония", VoiceTextAndTranslation: "овоз ба матн ва тарҷумаи гуфтаҳо", ImageTextTranslation: "тарҷумаи матн аз тасвир", VoicePhotoContextPractice: "машқ бо контексти овоз ё акс", MonthPrice: "30 рӯз: %d ₽ ё %d Stars", YearPrice: "Сол: %d ₽ ё %d Stars (тақрибан %d%% тахфиф)", InviteFree: "Дӯстро даъват карда 7 рӯз ройгон гиред", InvoiceDescription: "Premium: дар як рӯз %d дарс, %d паёми машқ, %d паёми овозӣ то %d сония, овоз ба матн ва тарҷума аз тасвир.", ActivatedTitle: "Premium барои %s фаъол шуд", LimitsTitle: "Лимитҳои ман", PlanLabel: "Тариф", LessonsToday: "Дарсҳои имрӯз", PracticeToday: "Машқи имрӯз", VoicesToday: "Овози имрӯз", InviteFriendButton: "Дӯстро даъват кунед", PaySBPButton: "Бо СБП пардохт кардан", PayStarsButton: "%d Stars пардохт кардан", MonthTitle: "Premium 30 рӯз", YearTitle: "Premium барои як сол", MonthDays: "30 рӯз", YearDays: "365 рӯз"},
	"hy": {FreeVoiceUnavailable: "ձայնային հաղորդագրությունները հասանելի չեն", LessonsPerDay: "օրական %d դաս", PracticeMessagesPerDay: "օրական %d պրակտիկայի հաղորդագրություն", PremiumVoicesPerDay: "օրական %d ձայնային հաղորդագրություն մինչեւ %d վայրկյան", VoiceTextAndTranslation: "ձայնը տեքստի վերածում եւ ասվածի թարգմանություն", ImageTextTranslation: "նկարի տեքստի թարգմանություն", VoicePhotoContextPractice: "պրակտիկա ձայնի կամ նկարի համատեքստով", MonthPrice: "30 օր: %d ₽ կամ %d Stars", YearPrice: "Տարի: %d ₽ կամ %d Stars (մոտ %d%% զեղչ)", InviteFree: "Հրավիրիր ընկերոջ եւ ստացիր 7 օր անվճար", InvoiceDescription: "Premium: օրական %d դաս, %d պրակտիկայի հաղորդագրություն, %d ձայնային մինչեւ %d վայրկյան, ձայնը տեքստի եւ նկարի թարգմանություն:", ActivatedTitle: "Premium-ը ակտիվացվեց %s ժամկետով", LimitsTitle: "Իմ սահմանաչափերը", PlanLabel: "Սակագին", LessonsToday: "Այսօրվա դասերը", PracticeToday: "Այսօրվա պրակտիկան", VoicesToday: "Այսօրվա ձայնայինները", InviteFriendButton: "Հրավիրել ընկերոջ", PaySBPButton: "Վճարել SBP-ով", PayStarsButton: "Վճարել %d Stars", MonthTitle: "Premium 30 օր", YearTitle: "Premium մեկ տարով", MonthDays: "30 օր", YearDays: "365 օր"},
}

var referralShareTextByInterfaceLanguage = map[string]string{
	"en": "Try Poliglot AI:\n%s\n\nPoliglot AI is a Telegram bot and web app for learning languages with AI: short lessons, conversation practice, vocabulary training, translator, voice-to-text, and photo text translation.\n\nOpen this link and you will receive 7 days of Premium in the bot.",
	"ru": "Попробуй Poliglot AI:\n%s\n\nPoliglot AI - это Telegram-бот и веб-версия для изучения языков с AI: короткие уроки, разговорная практика, тренировка слов, переводчик, голос в текст и перевод текста с фото.\n\nОткрой эту ссылку и получи 7 дней Premium-подписки в боте.",
	"es": "Prueba Poliglot AI:\n%s\n\nPoliglot AI es un bot de Telegram y una versión web para aprender idiomas con AI: lecciones cortas, práctica de conversación, vocabulario, traductor, voz a texto y traducción de texto en fotos.\n\nAbre este enlace y recibirás 7 días de suscripción Premium en el bot.",
	"de": "Probiere Poliglot AI aus:\n%s\n\nPoliglot AI ist ein Telegram-Bot und eine Web-Version zum Sprachenlernen mit AI: kurze Lektionen, Gesprächspraxis, Worttraining, Übersetzer, Sprache-zu-Text und Übersetzung von Text auf Fotos.\n\nÖffne diesen Link und du erhältst 7 Tage Premium im Bot.",
	"fr": "Essaie Poliglot AI :\n%s\n\nPoliglot AI est un bot Telegram et une version web pour apprendre les langues avec AI : courtes leçons, pratique de conversation, vocabulaire, traducteur, voix en texte et traduction du texte des photos.\n\nOuvre ce lien et tu recevras 7 jours d'abonnement Premium dans le bot.",
	"it": "Prova Poliglot AI:\n%s\n\nPoliglot AI è un bot Telegram e una versione web per imparare le lingue con AI: lezioni brevi, pratica di conversazione, vocaboli, traduttore, voce in testo e traduzione del testo nelle foto.\n\nApri questo link e riceverai 7 giorni di abbonamento Premium nel bot.",
	"pt": "Experimenta o Poliglot AI:\n%s\n\nO Poliglot AI é um bot do Telegram e uma versão web para aprender idiomas com AI: lições curtas, prática de conversa, vocabulário, tradutor, voz para texto e tradução de texto em fotos.\n\nAbre este link e receberás 7 dias de assinatura Premium no bot.",
	"pl": "Wypróbuj Poliglot AI:\n%s\n\nPoliglot AI to bot Telegram i wersja web do nauki języków z AI: krótkie lekcje, praktyka rozmowy, słownictwo, tłumacz, głos na tekst i tłumaczenie tekstu ze zdjęć.\n\nOtwórz ten link, a otrzymasz 7 dni subskrypcji Premium w bocie.",
	"ro": "Încearcă Poliglot AI:\n%s\n\nPoliglot AI este un bot Telegram și o versiune web pentru învățarea limbilor cu AI: lecții scurte, practică de conversație, vocabular, traducător, voce în text și traducerea textului din fotografii.\n\nDeschide acest link și vei primi 7 zile de abonament Premium în bot.",
	"uk": "Спробуй Poliglot AI:\n%s\n\nPoliglot AI - це Telegram-бот і веб-версія для вивчення мов з AI: короткі уроки, розмовна практика, тренування слів, перекладач, голос у текст і переклад тексту з фото.\n\nВідкрий це посилання й отримай 7 днів Premium-підписки в боті.",
	"kk": "Poliglot AI-ды байқап көр:\n%s\n\nPoliglot AI - AI арқылы тіл үйренуге арналған Telegram боты және web нұсқасы: қысқа сабақтар, сөйлесу практикасы, сөз жаттығулары, аудармашы, дауысты мәтінге айналдыру және фотодағы мәтінді аудару.\n\nОсы сілтемені ашсаң, ботта 7 күн Premium жазылымын аласың.",
	"ky": "Poliglot AI'ды байкап көр:\n%s\n\nPoliglot AI - AI менен тил үйрөнүү үчүн Telegram боту жана web нускасы: кыска сабактар, сүйлөшүү практикасы, сөз машыгуулары, котормочу, үндү текстке айлантуу жана сүрөттөгү текстти которуу.\n\nБул шилтемени ачсаң, ботто 7 күн Premium жазылуусун аласың.",
	"ka": "სცადე Poliglot AI:\n%s\n\nPoliglot AI არის Telegram ბოტი და web ვერსია AI-ით ენების სასწავლად: მოკლე გაკვეთილები, საუბრის პრაქტიკა, სიტყვების ვარჯიში, თარჯიმანი, ხმის ტექსტად გადაყვანა და ფოტოზე ტექსტის თარგმნა.\n\nგახსენი ეს ბმული და ბოტში მიიღებ 7 დღიან Premium გამოწერას.",
	"uz": "Poliglot AI'ni sinab ko‘ring:\n%s\n\nPoliglot AI - AI bilan til o‘rganish uchun Telegram boti va web versiya: qisqa darslar, suhbat mashqi, so‘zlar, tarjimon, ovozni matnga aylantirish va rasmdagi matnni tarjima qilish.\n\nUshbu havolani oching va botda 7 kunlik Premium obunasini oling.",
	"tt": "Poliglot AI'ны сынап кара:\n%s\n\nPoliglot AI - AI белән тел өйрәнү өчен Telegram боты һәм web версиясе: кыска дәресләр, сөйләшү күнегүе, сүзләр, тәрҗемәче, тавышны текстка әйләндерү һәм фотодагы текстны тәрҗемә итү.\n\nБу сылтаманы ачсаң, ботта 7 көнлек Premium язылуын аласың.",
	"tg": "Poliglot AI-ро санҷед:\n%s\n\nPoliglot AI - боти Telegram ва версияи web барои омӯзиши забонҳо бо AI аст: дарсҳои кӯтоҳ, машқи гуфтугӯ, луғат, тарҷумон, овоз ба матн ва тарҷумаи матн аз акс.\n\nИн пайвандро кушоед ва дар бот 7 рӯзи обунаи Premium мегиред.",
	"hy": "Փորձիր Poliglot AI-ը:\n%s\n\nPoliglot AI-ը Telegram բոտ և web տարբերակ է AI-ով լեզուներ սովորելու համար՝ կարճ դասեր, խոսակցական պրակտիկա, բառերի մարզում, թարգմանիչ, ձայնը տեքստի վերածում և լուսանկարների տեքստի թարգմանություն։\n\nԲացիր այս հղումը և բոտում կստանաս Premium բաժանորդագրություն 7 օրով։",
	"zh": "试试 Poliglot AI：\n%s\n\nPoliglot AI 是一个用 AI 学语言的 Telegram 机器人和 web 版：短课程、对话练习、单词训练、翻译器、语音转文字、照片文字翻译。\n\n打开此链接即可在机器人中获得 7 天 Premium 订阅。",
	"ja": "Poliglot AI を試してみてください:\n%s\n\nPoliglot AI は AI で言語学習を支える Telegram ボットと web 版です。短いレッスン、会話練習、単語トレーニング、翻訳、音声の文字起こし、写真内テキストの翻訳ができます。\n\nこのリンクを開くと、ボットで Premium サブスクリプションを7日間受け取れます。",
	"ko": "Poliglot AI를 사용해 보세요:\n%s\n\nPoliglot AI는 AI로 언어를 배우는 Telegram 봇과 웹 버전입니다. 짧은 레슨, 대화 연습, 단어 학습, 번역기, 음성 텍스트 변환, 사진 속 텍스트 번역을 제공합니다.\n\n이 링크를 열면 봇에서 Premium 구독 7일을 받을 수 있습니다.",
}

func localizedReferralShareText(interfaceLanguage string) string {
	return referralShareTextByInterfaceLanguage[normalizeInterfaceLanguage(interfaceLanguage)]
}

func premiumUI(user userState) premiumUICopy {
	code := normalizeInterfaceLanguage(user.InterfaceLanguage)
	copy := englishPremiumUICopy()
	if code != "en" {
		copy = localizedGenericPremiumUICopy(code)
	}
	if override, ok := premiumUICopyOverrides[code]; ok {
		copy = mergePremiumUICopy(copy, override)
	}
	if referralShareText := localizedReferralShareText(user.InterfaceLanguage); referralShareText != "" {
		copy.InviteLinkText = referralShareText
		copy.ReferralShareText = referralShareText
	}
	return copy
}

func localizedGenericPremiumUICopy(code string) premiumUICopy {
	copy := ui(userState{InterfaceLanguage: code})
	premium := copy.Premium
	referral := firstNonEmpty(copy.Referral, premium)
	lessons := copy.NewLesson
	practice := copy.Practice
	voice := firstNonEmpty(copy.Tool.VoiceToText, copy.Shadowing)
	image := firstNonEmpty(copy.Tool.ImageTranslate, copy.Tools)
	limits := copy.Limits
	inviteText := fmt.Sprintf("Poliglot AI\n%%s\n\n%s / %s / %s / %s.", premium, lessons, practice, voice)
	return premiumUICopy{
		FreeVoiceUnavailable:        voice + " - " + premium,
		LessonsPerDay:               "%d " + lessons,
		PracticeMessagesPerDay:      "%d " + practice,
		PremiumVoicesPerDay:         "%d " + voice + " / %d",
		VoiceTextAndTranslation:     voice,
		ImageTextTranslation:        image,
		VoicePhotoContextPractice:   practice + ": " + voice + " / " + image,
		BestValueLabel:              premium + " 70%",
		MaxAccessLabel:              limits + " Max",
		PlatinumPriority:            "Platinum " + limits,
		MonthPrice:                  "30: %d ₽ / %d Stars",
		YearPrice:                   "365: %d ₽ / %d Stars (%d%%)",
		PlatinumMonthPrice:          "Platinum 30: %d ₽ / %d Stars",
		PlatinumYearPrice:           "Platinum 365: %d ₽ / %d Stars (%d%%)",
		InviteFree:                  referral + ": Premium 7",
		InvoiceDescription:          premium + ": %d " + lessons + ", %d " + practice + ", %d " + voice + " / %d.",
		PaymentUnrecognized:         premium + ": %s",
		PreCheckoutFailed:           premium + ". " + copy.TryAgain,
		ActivatedTitle:              premium + ": %s",
		ActivatedAvailable:          premium + ": %d / %d / %d",
		ActiveUntil:                 premium + ": %s",
		CheckLimits:                 "%s -> %s",
		YooKassaUnavailable:         "YooKassa: " + copy.TryAgain,
		YooKassaCreateFailed:        "YooKassa: " + copy.TryAgain,
		YooKassaPaymentText:         "*%s*\n\n%d ₽\nYooKassa\n\n" + premium,
		InviteCreateFailed:          referral + ": %s",
		InviteLinkText:              inviteText,
		ReferralShareText:           inviteText,
		ReferralInviter:             referral + ": Premium 7",
		ReferralInviterLevel:        referral + ": %d / %s",
		ReferralInvitee:             referral + ": Premium 7",
		ReferralBalanceTitle:        referral,
		ReferralInvitationsLabel:    referral,
		ReferralBalanceHint:         referral + ": 20% / 5% / USDT",
		ReferralWithdrawButton:      referral,
		ReferralWithdrawUnavailable: referral + ": 1000 ₽ / USDT. %s (%s).",
		LimitsTitle:                 limits,
		PlanLabel:                   premium,
		LessonsToday:                lessons,
		PracticeToday:               practice,
		VoicesToday:                 voice,
		MonthStarsButton:            "30: %d Stars",
		MonthRubButton:              "30: %d ₽ YooKassa",
		YearStarsButton:             "365: %d Stars",
		YearRubButton:               "365: %d ₽ YooKassa",
		PlanStarsButton:             "%s: %d Stars",
		PlanRubButton:               "%s: %d ₽ YooKassa",
		InviteFriendButton:          referral,
		PaySBPButton:                "SBP",
		PayStarsButton:              "%d Stars",
		MonthTitle:                  premium + " 30",
		YearTitle:                   premium + " 365",
		PlatinumMonthTitle:          "Platinum 30",
		PlatinumYearTitle:           "Platinum 365",
		MonthDays:                   "30",
		YearDays:                    "365",
	}
}

func mergePremiumUICopy(base, override premiumUICopy) premiumUICopy {
	if override.FreeVoiceUnavailable != "" {
		base.FreeVoiceUnavailable = override.FreeVoiceUnavailable
	}
	if override.LessonsPerDay != "" {
		base.LessonsPerDay = override.LessonsPerDay
	}
	if override.PracticeMessagesPerDay != "" {
		base.PracticeMessagesPerDay = override.PracticeMessagesPerDay
	}
	if override.PremiumVoicesPerDay != "" {
		base.PremiumVoicesPerDay = override.PremiumVoicesPerDay
	}
	if override.VoiceTextAndTranslation != "" {
		base.VoiceTextAndTranslation = override.VoiceTextAndTranslation
	}
	if override.ImageTextTranslation != "" {
		base.ImageTextTranslation = override.ImageTextTranslation
	}
	if override.VoicePhotoContextPractice != "" {
		base.VoicePhotoContextPractice = override.VoicePhotoContextPractice
	}
	if override.BestValueLabel != "" {
		base.BestValueLabel = override.BestValueLabel
	}
	if override.MaxAccessLabel != "" {
		base.MaxAccessLabel = override.MaxAccessLabel
	}
	if override.PlatinumPriority != "" {
		base.PlatinumPriority = override.PlatinumPriority
	}
	if override.MonthPrice != "" {
		base.MonthPrice = override.MonthPrice
	}
	if override.YearPrice != "" {
		base.YearPrice = override.YearPrice
	}
	if override.PlatinumMonthPrice != "" {
		base.PlatinumMonthPrice = override.PlatinumMonthPrice
	}
	if override.PlatinumYearPrice != "" {
		base.PlatinumYearPrice = override.PlatinumYearPrice
	}
	if override.InviteFree != "" {
		base.InviteFree = override.InviteFree
	}
	if override.InvoiceDescription != "" {
		base.InvoiceDescription = override.InvoiceDescription
	}
	if override.PaymentUnrecognized != "" {
		base.PaymentUnrecognized = override.PaymentUnrecognized
	}
	if override.PreCheckoutFailed != "" {
		base.PreCheckoutFailed = override.PreCheckoutFailed
	}
	if override.ActivatedTitle != "" {
		base.ActivatedTitle = override.ActivatedTitle
	}
	if override.ActivatedAvailable != "" {
		base.ActivatedAvailable = override.ActivatedAvailable
	}
	if override.ActiveUntil != "" {
		base.ActiveUntil = override.ActiveUntil
	}
	if override.CheckLimits != "" {
		base.CheckLimits = override.CheckLimits
	}
	if override.YooKassaUnavailable != "" {
		base.YooKassaUnavailable = override.YooKassaUnavailable
	}
	if override.YooKassaCreateFailed != "" {
		base.YooKassaCreateFailed = override.YooKassaCreateFailed
	}
	if override.YooKassaPaymentText != "" {
		base.YooKassaPaymentText = override.YooKassaPaymentText
	}
	if override.InviteCreateFailed != "" {
		base.InviteCreateFailed = override.InviteCreateFailed
	}
	if override.InviteLinkText != "" {
		base.InviteLinkText = override.InviteLinkText
	}
	if override.ReferralShareText != "" {
		base.ReferralShareText = override.ReferralShareText
	}
	if override.ReferralInviter != "" {
		base.ReferralInviter = override.ReferralInviter
	}
	if override.ReferralInviterLevel != "" {
		base.ReferralInviterLevel = override.ReferralInviterLevel
	}
	if override.ReferralInvitee != "" {
		base.ReferralInvitee = override.ReferralInvitee
	}
	if override.ReferralBalanceTitle != "" {
		base.ReferralBalanceTitle = override.ReferralBalanceTitle
	}
	if override.ReferralInvitationsLabel != "" {
		base.ReferralInvitationsLabel = override.ReferralInvitationsLabel
	}
	if override.ReferralBalanceHint != "" {
		base.ReferralBalanceHint = override.ReferralBalanceHint
	}
	if override.ReferralWithdrawButton != "" {
		base.ReferralWithdrawButton = override.ReferralWithdrawButton
	}
	if override.ReferralWithdrawUnavailable != "" {
		base.ReferralWithdrawUnavailable = override.ReferralWithdrawUnavailable
	}
	if override.LimitsTitle != "" {
		base.LimitsTitle = override.LimitsTitle
	}
	if override.PlanLabel != "" {
		base.PlanLabel = override.PlanLabel
	}
	if override.LessonsToday != "" {
		base.LessonsToday = override.LessonsToday
	}
	if override.PracticeToday != "" {
		base.PracticeToday = override.PracticeToday
	}
	if override.VoicesToday != "" {
		base.VoicesToday = override.VoicesToday
	}
	if override.MonthStarsButton != "" {
		base.MonthStarsButton = override.MonthStarsButton
	}
	if override.MonthRubButton != "" {
		base.MonthRubButton = override.MonthRubButton
	}
	if override.YearStarsButton != "" {
		base.YearStarsButton = override.YearStarsButton
	}
	if override.YearRubButton != "" {
		base.YearRubButton = override.YearRubButton
	}
	if override.PlanStarsButton != "" {
		base.PlanStarsButton = override.PlanStarsButton
	}
	if override.PlanRubButton != "" {
		base.PlanRubButton = override.PlanRubButton
	}
	if override.InviteFriendButton != "" {
		base.InviteFriendButton = override.InviteFriendButton
	}
	if override.PaySBPButton != "" {
		base.PaySBPButton = override.PaySBPButton
	}
	if override.PayStarsButton != "" {
		base.PayStarsButton = override.PayStarsButton
	}
	if override.MonthTitle != "" {
		base.MonthTitle = override.MonthTitle
	}
	if override.YearTitle != "" {
		base.YearTitle = override.YearTitle
	}
	if override.PlatinumMonthTitle != "" {
		base.PlatinumMonthTitle = override.PlatinumMonthTitle
	}
	if override.PlatinumYearTitle != "" {
		base.PlatinumYearTitle = override.PlatinumYearTitle
	}
	if override.MonthDays != "" {
		base.MonthDays = override.MonthDays
	}
	if override.YearDays != "" {
		base.YearDays = override.YearDays
	}
	return base
}

func localizedPremiumPlan(user userState, plan premiumPlan) premiumPlan {
	copy := premiumUI(user)
	switch plan.Product {
	case premiumMonthlyProduct:
		plan.Title = copy.MonthTitle
		plan.DaysLabel = copy.MonthDays
	case premiumYearlyProduct:
		plan.Title = copy.YearTitle
		plan.DaysLabel = copy.YearDays
	case platinumMonthlyProduct:
		plan.Title = copy.PlatinumMonthTitle
		plan.DaysLabel = copy.MonthDays
	case platinumYearlyProduct:
		plan.Title = copy.PlatinumYearTitle
		plan.DaysLabel = copy.YearDays
	}
	return plan
}

func localizedPremiumPlans(user userState, plans []premiumPlan) []premiumPlan {
	localized := make([]premiumPlan, len(plans))
	for i, plan := range plans {
		localized[i] = localizedPremiumPlan(user, plan)
	}
	return localized
}

func premiumActivationMarkdown(user userState, plan premiumPlan, until string, includeAvailable bool) string {
	copy := premiumUI(user)
	text := "*" + escapeMarkdownV2(fmt.Sprintf(copy.ActivatedTitle, plan.DaysLabel)) + "* ⭐\n\n"
	if includeAvailable {
		text += escapeMarkdownV2(fmt.Sprintf(copy.ActivatedAvailable, premiumLessonLimit, premiumPracticeLimit, premiumVoiceLimit)) + "\n\n"
	}
	text += escapeMarkdownV2(fmt.Sprintf(copy.ActiveUntil, until)) + "\n\n"
	text += escapeMarkdownV2(fmt.Sprintf(copy.CheckLimits, ui(user).MenuButton, ui(user).Limits))
	return text
}
