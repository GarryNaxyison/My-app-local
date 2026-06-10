package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	freeLessonLimit          = 5
	freePracticeLimit        = 15
	premiumLessonLimit       = 50
	premiumPracticeLimit     = 200
	premiumVoiceLimit        = 20
	platinumLessonLimit      = 100
	platinumPracticeLimit    = 500
	platinumVoiceLimit       = 60
	modeToolVoiceText        = "tool_voice_text"
	modeToolImageTranslate   = "tool_image_translate"
	modeToolTranslatorPrefix = "tool_translator:"
	webAuthStartPrefix       = "web_"
	modeSpellingPrefix       = "spelling:"
	modeMistakePrefix        = "mistake:"
	modeLevelTestPrefix      = "leveltest:"
	modeAITutorPrefix        = "ai_tutor:"
	spamWindow               = 5 * time.Second
	spamTimeout              = 30 * time.Second
	spamMaxMessages          = 5
	premiumDuration          = 30 * 24 * time.Hour
	premiumYearDuration      = 365 * 24 * time.Hour
	premiumMonthlyProduct    = "premium_30d"
	premiumYearlyProduct     = "premium_365d"
	platinumMonthlyProduct   = "platinum_30d"
	platinumYearlyProduct    = "platinum_365d"
	annualDiscountPercent    = 17
	telegramOpsRecipientID   = 185156683
)

type premiumPlan struct {
	Product    string
	Tier       string
	Title      string
	DaysLabel  string
	Duration   time.Duration
	RubPrice   int
	StarsPrice int
}

func paidProductIDs() []string {
	return []string{
		premiumMonthlyProduct,
		premiumYearlyProduct,
		platinumMonthlyProduct,
		platinumYearlyProduct,
	}
}

const welcomeText = `AI Polyglot Coach

First choose the bot language.
Сначала выбери язык бота.

Then I will ask for your time zone and learning language.`

const welcomeMarkdownText = `*AI Polyglot Coach* 👋

First choose the bot language\.
Сначала выбери язык бота\.

Then I will ask for your time zone and learning language\.`

const menuMarkdownText = `*Главное меню*

Выбери действие:`

func onboardingWelcomeText(user userState) string {
	if user.InterfaceSelected {
		copy := ui(user)
		return "Poliglot AI\n\n" +
			copy.ChooseBotLang + "\n" +
			copy.ChooseTimezone + "\n" +
			copy.ChooseLearnLang
	}

	var builder strings.Builder
	builder.WriteString("Poliglot AI\n\n")
	builder.WriteString("Choose the bot language first. Then I will ask for your time zone and learning language.\n\n")
	for _, language := range interfaceLanguages() {
		copy := ui(userState{InterfaceLanguage: language.Code, InterfaceSelected: true})
		builder.WriteString(language.InterfaceName)
		builder.WriteString(": ")
		builder.WriteString(copy.ChooseBotLang)
		builder.WriteString(" -> ")
		builder.WriteString(copy.ChooseTimezone)
		builder.WriteString(" -> ")
		builder.WriteString(copy.ChooseLearnLang)
		builder.WriteString("\n")
	}
	return strings.TrimSpace(builder.String())
}

var botRuntimeTexts = map[string]map[string]string{
	"en": {
		"lesson_task_failed":     "Could not get the task: %s",
		"word_game_not_found":    "Could not find a word for the game.",
		"word_task_missing":      "Could not find the word. Try opening the task again.",
		"word_game_missing":      "Could not find the word. Try opening the game again.",
		"spelling_not_found":     "Could not find a word for spelling practice.",
		"spelling_missing":       "Could not find the word to check. Open spelling again from the menu.",
		"mistake_missing":        "Could not find a mistake to practice. Open the mistake dictionary again.",
		"check_failed":           "Could not check the answer: %s",
		"reply_failed":           "Could not answer: %s",
		"level_promoted":         "New study level: %s 🎓\n\nCongratulations. Lessons, practice, and new words will now get gradually harder.",
		"translator_unavailable": "The translator is not configured yet.",
		"translator_empty":       "The translator returned an empty translation.",
		"translation_tts_failed": "Could not voice the translation: %s",
		"heard_label":            "I heard:",
		"web_auth_unavailable":   "Telegram login for the website is not configured yet. Try logging in with username and password.",
		"web_auth_failed":        "Could not confirm website login: %s",
		"web_auth_code":          "Website login code: %s\n\nEnter these 6 digits on the Poliglot AI website. The code is valid for 10 minutes.",
		"mistakes_empty_hint":    "Mistakes are saved automatically during lessons and practice.",
		"mistakes_empty_start":   "Start with %s -> %s",
	},
	"ru": {
		"lesson_task_failed":     "Не получилось получить задание: %s",
		"word_game_not_found":    "Не получилось найти слово для игры.",
		"word_task_missing":      "Не нашел слово. Попробуй открыть задание заново.",
		"word_game_missing":      "Не нашел слово. Попробуй открыть игру заново.",
		"spelling_not_found":     "Не получилось найти слово для правописания.",
		"spelling_missing":       "Не нашёл слово для проверки. Открой правописание заново через меню.",
		"mistake_missing":        "Не нашёл ошибку для тренировки. Открой словарь ошибок заново.",
		"check_failed":           "Не получилось проверить ответ: %s",
		"reply_failed":           "Не получилось ответить: %s",
		"level_promoted":         "Новый учебный уровень: %s 🎓\n\nПоздравляю! Теперь уроки, практика и новые слова будут постепенно сложнее.",
		"translator_unavailable": "Переводчик пока не настроен.",
		"translator_empty":       "Переводчик вернул пустой перевод.",
		"translation_tts_failed": "Не получилось озвучить перевод: %s",
		"heard_label":            "Я услышал:",
		"web_auth_unavailable":   "Вход через Telegram для сайта пока не настроен. Попробуйте войти по логину и паролю.",
		"web_auth_failed":        "Не получилось подтвердить вход на сайте: %s",
		"web_auth_code":          "Код для входа на сайт: %s\n\nВведите эти 6 цифр на сайте Poliglot AI. Код действует 10 минут.",
		"mistakes_empty_hint":    "Ошибки копятся автоматически во время уроков и практики.",
		"mistakes_empty_start":   "Начни с %s -> %s",
	},
	"zh": {
		"lesson_task_failed":     "无法获取任务：%s",
		"word_game_not_found":    "找不到用于游戏的单词。",
		"word_task_missing":      "找不到这个单词。请重新打开任务。",
		"word_game_missing":      "找不到这个单词。请重新打开游戏。",
		"spelling_not_found":     "找不到用于拼写练习的单词。",
		"spelling_missing":       "找不到要检查的单词。请从菜单重新打开拼写。",
		"mistake_missing":        "找不到要练习的错误。请重新打开错误词典。",
		"check_failed":           "无法检查答案：%s",
		"reply_failed":           "无法回复：%s",
		"level_promoted":         "新的学习水平：%s 🎓\n\n恭喜！课程、练习和新单词会逐步变难。",
		"translator_unavailable": "翻译器尚未配置。",
		"translator_empty":       "翻译器返回了空翻译。",
		"translation_tts_failed": "无法朗读翻译：%s",
		"heard_label":            "我听到：",
		"web_auth_unavailable":   "网站的 Telegram 登录尚未配置。请尝试用登录名和密码进入。",
		"web_auth_failed":        "无法确认网站登录：%s",
		"web_auth_code":          "网站登录代码：%s\n\n请在 Poliglot AI 网站输入这 6 位数字。代码有效期为 10 分钟。",
		"mistakes_empty_hint":    "错误会在课程和练习中自动保存。",
		"mistakes_empty_start":   "从 %s -> %s 开始",
	},
	"ja": {
		"lesson_task_failed":     "課題を取得できませんでした：%s",
		"word_game_not_found":    "ゲーム用の単語が見つかりませんでした。",
		"word_task_missing":      "単語が見つかりませんでした。課題を開き直してください。",
		"word_game_missing":      "単語が見つかりませんでした。ゲームを開き直してください。",
		"spelling_not_found":     "スペル練習用の単語が見つかりませんでした。",
		"spelling_missing":       "確認する単語が見つかりませんでした。メニューからスペル練習を開き直してください。",
		"mistake_missing":        "練習する間違いが見つかりませんでした。間違い辞書を開き直してください。",
		"check_failed":           "回答を確認できませんでした：%s",
		"reply_failed":           "返信できませんでした：%s",
		"level_promoted":         "新しい学習レベル：%s 🎓\n\nおめでとうございます。レッスン、練習、新しい単語が少しずつ難しくなります。",
		"translator_unavailable": "翻訳ツールはまだ設定されていません。",
		"translator_empty":       "翻訳ツールが空の翻訳を返しました。",
		"translation_tts_failed": "翻訳を音声化できませんでした：%s",
		"heard_label":            "聞き取った内容：",
		"web_auth_unavailable":   "Webサイト用のTelegramログインはまだ設定されていません。ログイン名とパスワードで試してください。",
		"web_auth_failed":        "Webサイトへのログインを確認できませんでした：%s",
		"web_auth_code":          "Webサイトのログインコード：%s\n\nPoliglot AIサイトでこの6桁を入力してください。コードは10分間有効です。",
		"mistakes_empty_hint":    "間違いはレッスンと練習中に自動で保存されます。",
		"mistakes_empty_start":   "%s -> %s から始めましょう",
	},
	"ko": {
		"lesson_task_failed":     "과제를 가져올 수 없습니다: %s",
		"word_game_not_found":    "게임에 사용할 단어를 찾을 수 없습니다.",
		"word_task_missing":      "단어를 찾을 수 없습니다. 과제를 다시 열어 주세요.",
		"word_game_missing":      "단어를 찾을 수 없습니다. 게임을 다시 열어 주세요.",
		"spelling_not_found":     "철자 연습용 단어를 찾을 수 없습니다.",
		"spelling_missing":       "확인할 단어를 찾을 수 없습니다. 메뉴에서 철자 연습을 다시 여세요.",
		"mistake_missing":        "연습할 오답을 찾을 수 없습니다. 오답 사전을 다시 여세요.",
		"check_failed":           "답변을 확인할 수 없습니다: %s",
		"reply_failed":           "답변할 수 없습니다: %s",
		"level_promoted":         "새 학습 레벨: %s 🎓\n\n축하합니다. 이제 수업, 연습, 새 단어가 조금씩 더 어려워집니다.",
		"translator_unavailable": "번역기가 아직 설정되지 않았습니다.",
		"translator_empty":       "번역기가 빈 번역을 반환했습니다.",
		"translation_tts_failed": "번역을 음성으로 만들 수 없습니다: %s",
		"heard_label":            "들은 내용:",
		"web_auth_unavailable":   "웹사이트용 Telegram 로그인이 아직 설정되지 않았습니다. 로그인과 비밀번호로 접속해 보세요.",
		"web_auth_failed":        "웹사이트 로그인을 확인할 수 없습니다: %s",
		"web_auth_code":          "웹사이트 로그인 코드: %s\n\nPoliglot AI 사이트에 이 6자리 숫자를 입력하세요. 코드는 10분 동안 유효합니다.",
		"mistakes_empty_hint":    "오답은 수업과 연습 중 자동으로 저장됩니다.",
		"mistakes_empty_start":   "%s -> %s 에서 시작하세요",
	},
}

func botRuntimeText(user userState, key string) string {
	code := normalizeInterfaceLanguage(user.InterfaceLanguage)
	if copy, ok := botRuntimeTexts[code]; ok {
		if text := copy[key]; text != "" {
			return text
		}
	}
	if text := botRuntimeTexts["en"][key]; text != "" {
		return text
	}
	return key
}

func botRuntimeMessage(user userState, key string, args ...any) string {
	return fmt.Sprintf(botRuntimeText(user, key), args...)
}

func mistakesEmptyMarkdown(user userState) string {
	copy := ui(user)
	return "*" + escapeMarkdownV2(copy.Mistakes) + "* 📖\n\n" +
		escapeMarkdownV2(botRuntimeText(user, "mistakes_empty_hint")) + "\n" +
		escapeMarkdownV2(botRuntimeMessage(user, "mistakes_empty_start", copy.MenuButton, copy.NewLesson))
}

type bot struct {
	cfg                   config
	store                 store
	telegram              *telegramClient
	openrouter            *openRouterClient
	aiTutor               *aiTutorEngine
	yookassa              *yooKassaClient
	webYooKassa           *yooKassaClient
	cryptoRates           *cryptoRateProvider
	activationKeys        *activationKeyStore
	webAuth               *webAuthBroker
	offset                int64
	throttleMu            sync.Mutex
	throttles             map[int64]messageThrottle
	pronunciationMu       sync.Mutex
	pronunciationMessages map[int64][]int64
	vocabularyHintMu      sync.Mutex
	vocabularyHints       map[string]string
}

func (b *bot) aiTutorEngine() *aiTutorEngine {
	if b.aiTutor != nil {
		return b.aiTutor
	}
	return newAITutorEngine(b.store, b.openrouter)
}

type messageThrottle struct {
	WindowStart  time.Time
	Count        int
	TimeoutUntil time.Time
	LastNotice   time.Time
}

func (b *bot) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		updates, err := b.telegram.getUpdates(ctx, b.offset)
		if err != nil {
			return err
		}

		for _, update := range updates {
			b.offset = update.UpdateID + 1
			if err := b.handleUpdate(ctx, update); err != nil {
				log.Printf("failed to handle update %d: %v", update.UpdateID, err)
			}
		}
	}
}

func (b *bot) checkMessageThrottle(userID int64, now time.Time) (bool, bool, time.Time) {
	if userID == 0 {
		return false, false, time.Time{}
	}
	b.throttleMu.Lock()
	defer b.throttleMu.Unlock()

	if b.throttles == nil {
		b.throttles = map[int64]messageThrottle{}
	}

	state := b.throttles[userID]
	if now.Before(state.TimeoutUntil) {
		notify := state.LastNotice.IsZero() || now.Sub(state.LastNotice) >= 5*time.Second
		if notify {
			state.LastNotice = now
			b.throttles[userID] = state
		}
		return true, notify, state.TimeoutUntil
	}

	if state.WindowStart.IsZero() || now.Sub(state.WindowStart) > spamWindow {
		state.WindowStart = now
		state.Count = 1
		state.TimeoutUntil = time.Time{}
		state.LastNotice = time.Time{}
		b.throttles[userID] = state
		return false, false, time.Time{}
	}

	state.Count++
	if state.Count > spamMaxMessages {
		state.TimeoutUntil = now.Add(spamTimeout)
		state.LastNotice = now
		b.throttles[userID] = state
		return true, true, state.TimeoutUntil
	}

	b.throttles[userID] = state
	return false, false, time.Time{}
}

func secondsUntil(until time.Time) int {
	seconds := int(time.Until(until).Seconds())
	if seconds < 1 {
		return 1
	}
	return seconds
}

func (b *bot) handleUpdate(ctx context.Context, update telegramUpdate) error {
	if update.CallbackQuery != nil {
		return b.handleCallbackQuery(ctx, *update.CallbackQuery)
	}
	if update.PreCheckoutQuery != nil {
		return b.handlePreCheckoutQuery(ctx, *update.PreCheckoutQuery)
	}
	if update.Message == nil {
		return nil
	}

	message := update.Message
	if message.SuccessfulPayment != nil {
		user, err := b.store.getOrCreateUser(message.From.ID, telegramDisplayName(message.From))
		if err != nil {
			return err
		}
		return b.handleSuccessfulPayment(ctx, message.Chat.ID, user, *message.SuccessfulPayment)
	}
	if blocked, notify, until := b.checkMessageThrottle(message.From.ID, time.Now()); blocked {
		if notify {
			user, _ := b.store.getOrCreateUser(message.From.ID, telegramDisplayName(message.From))
			copy := ui(user)
			return b.telegram.sendMessageWithCopy(ctx, message.Chat.ID, fmt.Sprintf(systemUI(user).SpamWait, secondsUntil(until)), copy)
		}
		return nil
	}
	if message.Voice != nil {
		return b.handleVoiceMessage(ctx, message)
	}
	if len(message.Photo) > 0 {
		return b.handlePhotoMessage(ctx, message)
	}
	if strings.TrimSpace(message.Text) == "" {
		return nil
	}

	user, err := b.store.getOrCreateUser(message.From.ID, telegramDisplayName(message.From))
	if err != nil {
		return err
	}

	text := strings.TrimSpace(message.Text)

	if handled, err := b.maybeBeginWebAuth(ctx, message.Chat.ID, user, message.From, text); handled || err != nil {
		return err
	}

	if handled, err := b.maybeApplyStartReferral(ctx, message.Chat.ID, user, text); handled || err != nil {
		return err
	}

	// Bottom panel buttons
	switch {
	case isMenuCommandText(text):
		return b.sendMainMenu(ctx, message.Chat.ID, user)
	case isStopCommandText(text):
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.telegram.sendMessageWithCopy(ctx, message.Chat.ID, ui(user).Stopped, ui(user))
	}

	if isStopIntent(text) {
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.telegram.sendMessageWithCopy(ctx, message.Chat.ID, ui(user).Stopped, ui(user))
	}

	if strings.HasPrefix(text, "/") {
		return b.handleCommand(ctx, message.Chat.ID, text, user)
	}

	if strings.HasPrefix(user.Mode, modeAITutorPrefix) {
		return b.handleAITutorText(ctx, message.Chat.ID, user, text)
	}
	if strings.HasPrefix(user.Mode, modeToolTranslatorPrefix) {
		return b.handleTranslatorText(ctx, message.Chat.ID, user, text)
	}
	if user.Mode == modeToolVoiceText || user.Mode == modeToolImageTranslate {
		return b.continueToolContextAsPractice(ctx, message.Chat.ID, message.From.ID, telegramDisplayName(message.From), user, text)
	}
	if strings.HasPrefix(user.Mode, modeSpellingPrefix) {
		return b.handleSpellingAnswer(ctx, message.Chat.ID, user, text)
	}
	if strings.HasPrefix(user.Mode, modeMistakePrefix) {
		return b.handleMistakePracticeAnswer(ctx, message.Chat.ID, user, text)
	}
	if strings.HasPrefix(user.Mode, modeShadowingPrefix) {
		return b.handleShadowingAnswer(ctx, message.Chat.ID, user, text, false, nil)
	}

	switch user.Mode {
	case "lesson":
		return b.handleLessonAnswer(ctx, message.Chat.ID, user, text, nil)
	case "practice":
		return b.handlePracticeMessage(ctx, message.Chat.ID, user, text, nil)
	case modeToolVoiceText:
		return b.telegram.sendMessage(ctx, message.Chat.ID, ui(user).Tool.VoicePrompt)
	case modeToolImageTranslate:
		return b.telegram.sendMessage(ctx, message.Chat.ID, ui(user).Tool.ImagePrompt)
	default:
		return b.telegram.sendMessage(ctx, message.Chat.ID, ui(user).UnknownButton)
	}
}

func (b *bot) handleCommand(ctx context.Context, chatID int64, text string, user userState) error {
	command := strings.ToLower(strings.Fields(text)[0])

	switch command {
	case "/start", "/help":
		if !user.InterfaceSelected {
			return b.sendPrivacyPrompt(ctx, chatID)
		}
		if err := b.telegram.sendMessage(ctx, chatID, onboardingWelcomeText(user)); err != nil {
			return err
		}
		return b.continueOnboarding(ctx, chatID, user)
	case "/menu":
		return b.sendMainMenu(ctx, chatID, user)
	case "/uilanguage", "/interface":
		return b.sendInterfaceLanguageSelection(ctx, chatID, user)
	case "/language":
		return b.sendLanguageSelection(ctx, chatID, user)
	case "/timezone", "/reminders":
		return b.sendReminderSettings(ctx, chatID, user)
	case "/reminderon":
		if err := b.store.setReminderEnabled(user.TelegramID, true); err != nil {
			return err
		}
		user.ReminderEnabled = true
		return b.sendReminderSettings(ctx, chatID, user)
	case "/reminderoff":
		if err := b.store.setReminderEnabled(user.TelegramID, false); err != nil {
			return err
		}
		user.ReminderEnabled = false
		return b.sendReminderSettings(ctx, chatID, user)
	case "/lesson":
		return b.startLesson(ctx, chatID, user)
	case "/practice":
		return b.startPractice(ctx, chatID, user)
	case "/shadowing", "/repeat":
		return b.startShadowing(ctx, chatID, user)
	case "/words":
		return b.startWordLesson(ctx, chatID, user)
	case "/game":
		return b.startWordGame(ctx, chatID, user)
	case "/phrasebook":
		return b.sendPhrasebook(ctx, chatID, user, 0)
	case "/spelling":
		return b.startSpellingPractice(ctx, chatID, user)
	case "/progress":
		return b.sendProgress(ctx, chatID, user)
	case "/limits":
		return b.telegram.sendMarkdownMessage(ctx, chatID, b.limitsMarkdownText(user))
	case "/premium":
		return b.sendPremiumMenu(ctx, chatID, user)
	case "/buy":
		return b.sendPremiumMenu(ctx, chatID, user)
	case "/invite":
		return b.sendInviteLink(ctx, chatID, user)
	case "/referral", "/referrals":
		return b.sendReferralMenu(ctx, chatID, user)
	case "/mistakes":
		return b.sendMistakes(ctx, chatID, user, 0)
	case "/clearmistakes":
		return b.clearMistakes(ctx, chatID, user)
	case "/leaderboard":
		return b.sendLeaderboard(ctx, chatID, user)
	case "/tools":
		return b.sendToolsMenu(ctx, chatID, user)
	case "/level", "/leveltest":
		return b.startLevelAssessment(ctx, chatID, user)
	default:
		return b.telegram.sendMessage(ctx, chatID, ui(user).UnknownButton)
	}
}

func (b *bot) sendMainMenu(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, "*"+escapeMarkdownV2(copy.MainMenuTitle)+"*\n\n"+escapeMarkdownV2(copy.MainMenuBody), mainMenuInlineKeyboard(copy))
}

func (b *bot) sendToolsMenu(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	text := copy.Tools + "\n\n" +
		"🎙 " + copy.Tool.VoicePrompt + "\n\n" +
		"🖼 " + copy.Tool.ImagePrompt
	text += "\n\n🌐 " + copy.Tool.TranslatorPrompt
	return b.telegram.sendInlineMessage(ctx, chatID, text, toolsInlineKeyboard(b.cfg.WebAppURL, copy))
}

func (b *bot) sendLearningMenu(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	return b.telegram.sendInlineMessage(ctx, chatID, copy.Learning+"\n\n"+copy.MainMenuBody, learningMenuKeyboard(copy))
}

func (b *bot) sendPrivacyPrompt(ctx context.Context, chatID int64) error {
	text := "Перед началом подтвердите, что принимаете Политику обработки персональных данных."
	return b.telegram.sendInlineMessage(ctx, chatID, text, privacyPromptKeyboard())
}

func (b *bot) sendWordsMenu(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	return b.telegram.sendInlineMessage(ctx, chatID, copy.Words+"\n\n"+copy.MainMenuBody, wordsMenuKeyboard(copy))
}

func (b *bot) sendStatsMenu(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	return b.telegram.sendInlineMessage(ctx, chatID, copy.Stats+"\n\n"+copy.MainMenuBody, statsMenuKeyboard(copy))
}

func (b *bot) sendSettingsMenu(ctx context.Context, chatID int64, users ...userState) error {
	user := userState{}
	if len(users) > 0 {
		user = users[0]
	}
	copy := ui(user)
	text := copy.Settings + "\n\n"
	if len(users) > 0 {
		text += referralAccountText(user) + "\n\n"
	}
	text += copy.MainMenuBody
	return b.telegram.sendInlineMessage(ctx, chatID, text, settingsMenuKeyboard(user))
}

func (b *bot) sendReferralMenu(ctx context.Context, chatID int64, user userState) error {
	copy := ui(user)
	text := copy.Referral + "\n\n" + referralAccountText(user)
	return b.telegram.sendInlineMessage(ctx, chatID, text, referralMenuKeyboard(user))
}

func (b *bot) sendLanguageSelection(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	return b.telegram.sendInlineMessage(ctx, chatID, copy.ChooseLearnLang, languageSelectionKeyboard(copy))
}

func (b *bot) sendInterfaceLanguageSelection(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	text := copy.ChooseBotLang
	if len(users) == 0 || !users[0].InterfaceSelected {
		text = interfaceLanguageSelectionText()
	}
	return b.telegram.sendInlineMessage(ctx, chatID, text, interfaceLanguageSelectionKeyboard())
}

func (b *bot) continueOnboarding(ctx context.Context, chatID int64, user userState) error {
	if !user.InterfaceSelected {
		return b.sendInterfaceLanguageSelection(ctx, chatID, user)
	}
	if !user.TimezoneSelected {
		return b.sendTimezonePrompt(ctx, chatID, user)
	}
	if !user.LanguageSelected {
		return b.sendLanguageSelection(ctx, chatID, user)
	}
	return b.sendMainMenu(ctx, chatID, user)
}

func (b *bot) handleInterfaceLanguageChoice(ctx context.Context, chatID int64, user userState, data string) error {
	languageCode, ok := parseInterfaceLanguageCallback(data)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, systemUI(user).UnknownInterfaceLang, ui(user))
	}
	if err := b.store.setInterfaceLanguage(user.TelegramID, languageCode); err != nil {
		return err
	}
	user.InterfaceLanguage = languageCode
	user.InterfaceSelected = true
	language := interfaceLanguageByCode(languageCode)
	copy := ui(user)
	if err := b.telegram.sendMessageWithCopy(ctx, chatID, fmt.Sprintf(copy.BotLangSet, language.InterfaceName), copy); err != nil {
		return err
	}
	if !user.TimezoneSelected {
		return b.sendTimezonePrompt(ctx, chatID, user)
	}
	if !user.LanguageSelected {
		return b.sendLanguageSelection(ctx, chatID, user)
	}
	return b.sendMainMenu(ctx, chatID, user)
}

func (b *bot) handleLearningLanguageChoice(ctx context.Context, chatID int64, user userState, data string) error {
	languageCode, ok := parseLanguageCallback(data)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, systemUI(user).UnknownLearningLang, ui(user))
	}
	if err := b.store.setLearningLanguage(user.TelegramID, languageCode); err != nil {
		return err
	}
	user.LearningLanguage = languageCode
	user.LanguageSelected = true
	language := learningLanguageByCode(languageCode)
	copy := ui(user)
	text := fmt.Sprintf(copy.LearnLangSet, language.InterfaceName)
	if err := b.telegram.sendMessageWithCopy(ctx, chatID, text, copy); err != nil {
		return err
	}
	if !user.InterfaceSelected {
		return b.sendInterfaceLanguageSelection(ctx, chatID, user)
	}
	if !user.TimezoneSelected {
		return b.sendTimezonePrompt(ctx, chatID, user)
	}
	return b.startLevelAssessment(ctx, chatID, user)
}

func (b *bot) startLevelAssessment(ctx context.Context, chatID int64, user userState) error {
	if len(levelAssessmentQuestionsForUser(user)) == 0 {
		language := userLearningLanguage(user)
		return b.telegram.sendInlineMessage(ctx, chatID,
			fmt.Sprintf(systemUI(user).LevelUnavailable, language.InterfaceName),
			manualLevelSelectionKeyboard(ui(user)))
	}
	return b.telegram.sendInlineMessage(ctx, chatID, levelAssessmentStartTextForUser(user), levelAssessmentStartKeyboard(user))
}

func (b *bot) sendManualLevelSelection(ctx context.Context, chatID int64, user userState) error {
	return b.telegram.sendInlineMessage(ctx, chatID, manualLevelSelectionTextForUser(user), manualLevelSelectionKeyboard(ui(user)))
}

func (b *bot) sendLevelQuestion(ctx context.Context, chatID int64, user userState, index int, score int) error {
	questions := levelAssessmentQuestionsForUser(user)
	if index < 0 || index >= len(questions) {
		return b.telegram.sendMessageWithCopy(ctx, chatID, systemUI(user).LevelOpenFailed, ui(user))
	}
	return b.telegram.sendInlineMessage(ctx, chatID, levelQuestionTextForUser(user, index, score), levelAssessmentKeyboardForUser(user, index))
}

func (b *bot) handleManualLevelChoice(ctx context.Context, chatID int64, user userState, data string) error {
	level, ok := parseManualLevelCallback(data)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, systemUI(user).LevelUnknownChoice, ui(user))
	}
	if err := b.store.setUserLevel(user.TelegramID, level); err != nil {
		return err
	}
	if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
		return err
	}
	text := fmt.Sprintf(systemUI(user).LevelManualSet, level)
	return b.telegram.sendInlineMessage(ctx, chatID, text, levelAssessmentDoneKeyboard(ui(user)))
}

func (b *bot) handleLevelAssessmentChoice(ctx context.Context, chatID int64, user userState, data string) error {
	index, answer, ok := parseLevelTestCallback(data)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, systemUI(user).LevelUnknownAnswer, ui(user))
	}
	expectedIndex, score, ok := parseLevelTestMode(user.Mode)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, systemUI(user).LevelInactive, ui(user))
	}
	if index != expectedIndex {
		return b.telegram.sendMessageWithCopy(ctx, chatID, systemUI(user).LevelStale, ui(user))
	}
	questions := levelAssessmentQuestionsForUser(user)
	if index >= len(questions) {
		return b.telegram.sendMessageWithCopy(ctx, chatID, systemUI(user).LevelStale, ui(user))
	}
	question := questions[index]
	if answer == question.CorrectIndex {
		score += question.Points
	}
	nextIndex := index + 1
	if nextIndex < len(questions) {
		if err := b.store.setMode(user.TelegramID, levelTestMode(nextIndex, score)); err != nil {
			return err
		}
		return b.sendLevelQuestion(ctx, chatID, user, nextIndex, score)
	}

	level := levelFromAssessmentScore(score, levelAssessmentMaxScoreForUser(user))
	if err := b.store.setUserLevel(user.TelegramID, level); err != nil {
		return err
	}
	if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
		return err
	}
	text := fmt.Sprintf(systemUI(user).LevelTestSet, level)
	return b.telegram.sendInlineMessage(ctx, chatID, text, levelAssessmentDoneKeyboard(ui(user)))
}

func (b *bot) sendTimezonePrompt(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	text := copy.ChooseTimezone + " 🌍\n\n" + copy.TimezoneHint
	return b.telegram.sendInlineMessage(ctx, chatID, text, timezoneInlineKeyboard(copy))
}

func (b *bot) sendReminderSettings(ctx context.Context, chatID int64, user userState) error {
	return b.telegram.sendInlineMessage(ctx, chatID, reminderSettingsText(user), reminderSettingsInlineKeyboard(user))
}

func reminderSettingsText(user userState) string {
	copy := systemUI(user)
	status := copy.ReminderDisabled
	if user.ReminderEnabled {
		status = copy.ReminderEnabled
	}

	timezone := timezoneLabel(user.ReminderUTCOffset)
	if !user.TimezoneSelected {
		timezone += " (" + copy.ReminderDefault + ")"
	}

	return fmt.Sprintf("%s\n\n%s: %s\n%s: %s\n%s: %02d:00\n\n%s",
		copy.ReminderTitle,
		copy.ReminderStatus, status,
		copy.ReminderTimezone, timezone,
		copy.ReminderHour, user.ReminderHour,
		copy.ReminderHint)
}

func parseTimezoneCallback(data string) (int, bool) {
	parts := strings.Split(data, "|")
	if len(parts) != 2 {
		return 0, false
	}
	offsetMinutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, false
	}
	if offsetMinutes < -720 || offsetMinutes > 840 {
		return 0, false
	}
	return offsetMinutes, true
}

func timezoneLabel(offsetMinutes int) string {
	sign := "+"
	if offsetMinutes < 0 {
		sign = "-"
		offsetMinutes = -offsetMinutes
	}
	hours := offsetMinutes / 60
	minutes := offsetMinutes % 60
	if minutes == 0 {
		return fmt.Sprintf("UTC%s%d", sign, hours)
	}
	return fmt.Sprintf("UTC%s%d:%02d", sign, hours, minutes)
}

func userLocalDateTimeLabel(value time.Time, user userState) string {
	offsetMinutes := user.ReminderUTCOffset
	if offsetMinutes == 0 && !user.TimezoneSelected {
		offsetMinutes = 180
	}
	local := value.UTC().Add(time.Duration(offsetMinutes) * time.Minute)
	return fmt.Sprintf("%s (%s)", local.Format("2006-01-02 15:04"), timezoneLabel(offsetMinutes))
}

func (b *bot) handleCallbackQuery(ctx context.Context, query callbackQuery) error {
	_ = b.telegram.answerCallbackQuery(ctx, query.ID, "")

	user, err := b.store.getOrCreateUser(query.From.ID, telegramDisplayName(query.From))
	if err != nil {
		return err
	}

	chatID := query.From.ID
	if query.Message != nil {
		chatID = query.Message.Chat.ID
		ctx = contextWithTelegramEditTarget(ctx, chatID, query.Message.MessageID)
	}

	if strings.HasPrefix(query.Data, "tz|") {
		offsetMinutes, ok := parseTimezoneCallback(query.Data)
		if !ok {
			return b.telegram.sendMessageWithCopy(ctx, chatID, systemUI(user).UnknownTimezone, ui(user))
		}
		firstTimezoneSelection := !user.TimezoneSelected
		if err := b.store.setTimezone(user.TelegramID, offsetMinutes); err != nil {
			return err
		}
		user.ReminderUTCOffset = offsetMinutes
		user.TimezoneSelected = true
		if firstTimezoneSelection {
			if !user.InterfaceSelected {
				return b.sendInterfaceLanguageSelection(ctx, chatID, user)
			}
			if !user.LanguageSelected {
				return b.sendLanguageSelection(ctx, chatID, user)
			}
			return b.startLevelAssessment(ctx, chatID, user)
		}
		return b.sendReminderSettings(ctx, chatID, user)
	}

	switch query.Data {
	case "menu_learning":
		return b.sendLearningMenu(ctx, chatID, user)
	case "menu_words":
		return b.sendWordsMenu(ctx, chatID, user)
	case "menu_stats":
		return b.sendStatsMenu(ctx, chatID, user)
	case "menu_settings":
		return b.sendSettingsMenu(ctx, chatID, user)
	case "menu_lesson":
		return b.startLesson(ctx, chatID, user)
	case "menu_practice":
		return b.startPractice(ctx, chatID, user)
	case "menu_shadowing":
		return b.startShadowing(ctx, chatID, user)
	case "menu_tutor":
		return b.startTutorLesson(ctx, chatID, user)
	case "menu_word_lesson":
		return b.startWordLesson(ctx, chatID, user)
	case "menu_word_game":
		return b.startWordGame(ctx, chatID, user)
	case "menu_spelling":
		return b.startSpellingPractice(ctx, chatID, user)
	case "menu_vocabulary":
		return b.sendVocabulary(ctx, chatID, user, 0)
	case "menu_phrasebook":
		return b.sendPhrasebook(ctx, chatID, user, 0)
	case "menu_tools":
		return b.sendToolsMenu(ctx, chatID, user)
	case "privacy_continue":
		if err := b.telegram.sendMessage(ctx, chatID, onboardingWelcomeText(user)); err != nil {
			return err
		}
		return b.continueOnboarding(ctx, chatID, user)
	case "menu_interface_language":
		return b.sendInterfaceLanguageSelection(ctx, chatID, user)
	case "menu_language":
		return b.sendLanguageSelection(ctx, chatID, user)
	case "menu_level_test":
		return b.startLevelAssessment(ctx, chatID, user)
	case "level_manual_select":
		return b.sendManualLevelSelection(ctx, chatID, user)
	case "level_test_start":
		if err := b.store.setMode(user.TelegramID, levelTestMode(0, 0)); err != nil {
			return err
		}
		return b.sendLevelQuestion(ctx, chatID, user, 0, 0)
	case "menu_mistakes":
		if strings.HasPrefix(user.Mode, modeMistakePrefix) {
			if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
				return err
			}
		}
		return b.sendMistakes(ctx, chatID, user, 0)
	case "practice_mistakes":
		return b.startMistakePractice(ctx, chatID, user)
	case "menu_progress":
		return b.sendProgress(ctx, chatID, user)
	case "menu_limits":
		return b.telegram.sendMarkdownMessage(ctx, chatID, b.limitsMarkdownText(user))
	case "menu_premium":
		return b.sendPremiumMenu(ctx, chatID, user)
	case "menu_referral":
		return b.sendReferralMenu(ctx, chatID, user)
	case "menu_leaderboard":
		return b.sendLeaderboardMenu(ctx, chatID, user)
	case "menu_reminders", "menu_timezone":
		return b.sendReminderSettings(ctx, chatID, user)
	case "reminder_on":
		if err := b.store.setReminderEnabled(user.TelegramID, true); err != nil {
			return err
		}
		user.ReminderEnabled = true
		return b.sendReminderSettings(ctx, chatID, user)
	case "reminder_off":
		if err := b.store.setReminderEnabled(user.TelegramID, false); err != nil {
			return err
		}
		user.ReminderEnabled = false
		return b.sendReminderSettings(ctx, chatID, user)
	case "tool_voice_text":
		if err := b.store.setMode(user.TelegramID, modeToolVoiceText); err != nil {
			return err
		}
		return b.telegram.sendMessage(ctx, chatID, ui(user).Tool.VoiceModePrompt)
	case "tool_image_translate":
		if err := b.store.setMode(user.TelegramID, modeToolImageTranslate); err != nil {
			return err
		}
		return b.telegram.sendMessage(ctx, chatID, ui(user).Tool.ImageModePrompt)
	case "tool_translator":
		if err := b.store.setMode(user.TelegramID, translatorMode(defaultTranslatorSource(), defaultTranslatorTarget(user))); err != nil {
			return err
		}
		refreshed, err := b.store.getOrCreateUser(user.TelegramID, user.FirstName)
		if err != nil {
			return err
		}
		return b.sendTranslatorMenu(ctx, chatID, refreshed)
	case "translator_source_menu":
		return b.telegram.sendInlineMessage(ctx, chatID, ui(user).Tool.SourceLanguage, translatorLanguageKeyboard(user, false))
	case "translator_target_menu":
		return b.telegram.sendInlineMessage(ctx, chatID, ui(user).Tool.TargetLanguage, translatorLanguageKeyboard(user, true))
	case "translator_swap":
		return b.swapTranslatorLanguages(ctx, chatID, user)
	case "buy_premium", "buy_premium_month", "buy_" + premiumMonthlyProduct:
		return b.sendPremiumPaymentOptions(ctx, chatID, user, premiumMonthlyProduct)
	case "buy_premium_year", "buy_" + premiumYearlyProduct:
		return b.sendPremiumPaymentOptions(ctx, chatID, user, premiumYearlyProduct)
	case "buy_platinum_month", "buy_" + platinumMonthlyProduct:
		return b.sendPremiumPaymentOptions(ctx, chatID, user, platinumMonthlyProduct)
	case "buy_platinum_year", "buy_" + platinumYearlyProduct:
		return b.sendPremiumPaymentOptions(ctx, chatID, user, platinumYearlyProduct)
	case "buy_yookassa", "buy_yookassa_month":
		return b.sendYooKassaPayment(ctx, chatID, user, premiumMonthlyProduct)
	case "buy_yookassa_year":
		return b.sendYooKassaPayment(ctx, chatID, user, premiumYearlyProduct)
	case "buy_yookassa_platinum_month":
		return b.sendYooKassaPayment(ctx, chatID, user, platinumMonthlyProduct)
	case "buy_yookassa_platinum_year":
		return b.sendYooKassaPayment(ctx, chatID, user, platinumYearlyProduct)
	case "buy_crypto_month":
		return b.sendCryptoPayment(ctx, chatID, user, premiumMonthlyProduct, cryptoMethodTON)
	case "buy_crypto_year":
		return b.sendCryptoPayment(ctx, chatID, user, premiumYearlyProduct, cryptoMethodTON)
	case "invite_friend":
		return b.sendInviteLink(ctx, chatID, user)
	case "referral_withdraw":
		return b.sendReferralWithdrawInfo(ctx, chatID, user)
	case "show_limits":
		return b.telegram.sendMarkdownMessage(ctx, chatID, b.limitsMarkdownText(user))
	case "clear_mistakes":
		return b.clearMistakes(ctx, chatID, user)
	case "stop_mode_menu":
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.sendMainMenu(ctx, chatID, user)
	case "back_menu":
		if strings.HasPrefix(user.Mode, modeLevelTestPrefix) {
			if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
				return err
			}
		}
		return b.sendMainMenu(ctx, chatID, user)
	default:
		if strings.HasPrefix(query.Data, "ait|") {
			return b.handleAITutorCallback(ctx, chatID, user, query.Data)
		}
		if strings.HasPrefix(query.Data, "trsrc|") {
			return b.handleTranslatorLanguageCallback(ctx, chatID, user, query.Data, false)
		}
		if strings.HasPrefix(query.Data, "trdst|") {
			return b.handleTranslatorLanguageCallback(ctx, chatID, user, query.Data, true)
		}
		if strings.HasPrefix(query.Data, "ui|") {
			return b.handleInterfaceLanguageChoice(ctx, chatID, user, query.Data)
		}
		if strings.HasPrefix(query.Data, "lang|") {
			return b.handleLearningLanguageChoice(ctx, chatID, user, query.Data)
		}
		if strings.HasPrefix(query.Data, "level|") {
			return b.handleManualLevelChoice(ctx, chatID, user, query.Data)
		}
		if strings.HasPrefix(query.Data, "lt|") {
			return b.handleLevelAssessmentChoice(ctx, chatID, user, query.Data)
		}
		if strings.HasPrefix(query.Data, "wl|") {
			return b.handleWordLessonChoice(ctx, chatID, user, query.Data)
		}
		if strings.HasPrefix(query.Data, "wg|") {
			return b.handleWordGameChoice(ctx, chatID, user, query.Data)
		}
		if strings.HasPrefix(query.Data, "vocab|") {
			page, ok := parseVocabularyPageCallback(query.Data)
			if !ok {
				return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
			}
			return b.sendVocabulary(ctx, chatID, user, page)
		}
		if strings.HasPrefix(query.Data, "phrasebook|") {
			page, err := strconv.Atoi(strings.TrimPrefix(query.Data, "phrasebook|"))
			if err != nil {
				return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
			}
			return b.sendPhrasebook(ctx, chatID, user, page)
		}
		if strings.HasPrefix(query.Data, "mistakes|") {
			page, ok := parseMistakesPageCallback(query.Data)
			if !ok {
				return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
			}
			return b.sendMistakes(ctx, chatID, user, page)
		}
		if strings.HasPrefix(query.Data, "leaderboard|") {
			scope := strings.TrimPrefix(query.Data, "leaderboard|")
			if scope == "all" {
				return b.sendLeaderboard(ctx, chatID, user)
			}
			return b.sendLanguageLeaderboard(ctx, chatID, scope, user)
		}
		if strings.HasPrefix(query.Data, "check_crypto|") {
			paymentID := strings.TrimPrefix(query.Data, "check_crypto|")
			return b.checkCryptoPayment(ctx, chatID, user, paymentID)
		}
		if strings.HasPrefix(query.Data, "premium_plan|") {
			product := strings.TrimPrefix(query.Data, "premium_plan|")
			return b.sendPremiumPaymentOptions(ctx, chatID, user, product)
		}
		if strings.HasPrefix(query.Data, "buy_stars|") {
			product := strings.TrimPrefix(query.Data, "buy_stars|")
			return b.sendPremiumInvoice(ctx, chatID, user, product)
		}
		if strings.HasPrefix(query.Data, "buy_yookassa|") {
			product := strings.TrimPrefix(query.Data, "buy_yookassa|")
			return b.sendYooKassaPayment(ctx, chatID, user, product)
		}
		if strings.HasPrefix(query.Data, "buy_crypto|") {
			parts := strings.Split(query.Data, "|")
			if len(parts) != 3 {
				return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
			}
			return b.sendCryptoPayment(ctx, chatID, user, parts[1], parts[2])
		}
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
}

func (b *bot) startLesson(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	if !canUseLessons(user) {
		return b.telegram.sendMessageWithCopy(ctx, chatID, b.limitReachedText(systemUI(user).LessonKind, user), ui(user))
	}
	_ = b.telegram.sendChatAction(ctx, chatID, "typing")
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	task, err := b.openrouter.complete(ctx, lessonPrompt(language, interfaceLanguage, user.Level, user.LearningFocus, user.LessonCount, user.LessonHistory), 0.7, 500)
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeMessage(user, "lesson_task_failed", err.Error()), ui(user))
	}
	if err := b.store.saveLesson(user.TelegramID, task); err != nil {
		return err
	}
	return b.telegram.sendSmoothMessageWithCopy(ctx, chatID, task+"\n\n"+systemUI(user).LessonAnswerInstruction, ui(user))
}

func (b *bot) startPractice(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	if !canUsePractice(user) {
		return b.telegram.sendMessageWithCopy(ctx, chatID, b.limitReachedText(systemUI(user).PracticeKind, user), ui(user))
	}
	if err := b.store.setMode(user.TelegramID, "practice"); err != nil {
		return err
	}
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	return b.telegram.sendMessageWithCopy(ctx, chatID,
		fmt.Sprintf(systemUI(user).PracticeStarted, language.PracticeDirection, interfaceLanguage.InterfaceName, ui(user).StopButton),
		ui(user))
}

func (b *bot) startTutorLesson(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	_ = b.telegram.sendChatAction(ctx, chatID, "typing")
	result, err := b.aiTutorEngine().Start(ctx, user, "telegram")
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeMessage(user, "lesson_task_failed", err.Error()), ui(user))
	}
	if err := b.store.setMode(user.TelegramID, modeAITutorPrefix+result.Session.ID); err != nil {
		return err
	}
	return b.sendAITutorStep(ctx, chatID, user, result)
}

func (b *bot) handleAITutorText(ctx context.Context, chatID int64, user userState, text string) error {
	sessionID := strings.TrimPrefix(user.Mode, modeAITutorPrefix)
	result, err := b.aiTutorEngine().Submit(ctx, user, sessionID, aiTutorSubmitInput{Text: text})
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, err.Error(), ui(user))
	}
	if result.Session.Status == aiTutorSessionComplete {
		_ = b.store.setMode(user.TelegramID, "idle")
	}
	return b.sendAITutorStep(ctx, chatID, user, result)
}

func (b *bot) handleAITutorCallback(ctx context.Context, chatID int64, user userState, data string) error {
	parts := strings.Split(data, "|")
	if len(parts) < 4 {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
	sessionID := strings.TrimSpace(parts[1])
	value := strings.TrimSpace(parts[3])
	result, err := b.aiTutorEngine().Submit(ctx, user, sessionID, aiTutorSubmitInput{Choice: value})
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, err.Error(), ui(user))
	}
	if result.Session.Status == aiTutorSessionComplete {
		_ = b.store.setMode(user.TelegramID, "idle")
	}
	return b.sendAITutorStep(ctx, chatID, user, result)
}

func (b *bot) sendAITutorStep(ctx context.Context, chatID int64, user userState, result aiTutorResult) error {
	text := formatTelegramAITutorStep(result.NextStep, result.Feedback, user)
	keyboard := aiTutorTelegramKeyboard(result.Session.ID, result.NextStep, ui(user))
	if keyboard != nil {
		return b.telegram.sendInlineMessage(ctx, chatID, text, keyboard)
	}
	return b.telegram.sendMessageWithCopy(ctx, chatID, text, ui(user))
}

func formatTelegramAITutorStep(step aiTutorStep, feedback aiTutorFeedback, user userState) string {
	var builder strings.Builder
	if strings.TrimSpace(feedback.Message) != "" {
		builder.WriteString(strings.TrimSpace(feedback.Message))
		builder.WriteString("\n\n")
	}
	if strings.TrimSpace(step.Title) != "" {
		builder.WriteString(strings.TrimSpace(step.Title))
		builder.WriteString("\n\n")
	}
	if strings.TrimSpace(step.Instruction) != "" {
		builder.WriteString(strings.TrimSpace(step.Instruction))
	}
	if step.Kind == "story" && strings.TrimSpace(step.Lesson.Story.TextTarget) != "" {
		builder.WriteString("\n\n")
		builder.WriteString(strings.TrimSpace(step.Lesson.Story.TextTarget))
	}
	if step.Word != nil {
		builder.WriteString("\n\n")
		builder.WriteString(strings.TrimSpace(step.Word.Target))
		builder.WriteString(" - ")
		builder.WriteString(strings.TrimSpace(step.Word.InterfaceTranslation))
	}
	if step.Question != nil && strings.TrimSpace(step.Question.QuestionTarget) != "" && !strings.Contains(builder.String(), step.Question.QuestionTarget) {
		builder.WriteString("\n\n")
		builder.WriteString(strings.TrimSpace(step.Question.QuestionTarget))
	}
	text := strings.TrimSpace(builder.String())
	if text == "" {
		text = ui(user).AITutor
	}
	return text
}

func aiTutorTelegramKeyboard(sessionID string, step aiTutorStep, copy uiCopy) map[string]any {
	callback := func(value string) string {
		return "ait|" + sessionID + "|choice|" + value
	}
	switch step.Kind {
	case "story", "word_learn":
		return map[string]any{"inline_keyboard": [][]map[string]any{
			{{"text": "Continue", "callback_data": callback("continue")}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		}}
	case "rating":
		return map[string]any{"inline_keyboard": [][]map[string]any{
			{
				{"text": "Easy", "callback_data": callback("easy")},
				{"text": "Good", "callback_data": callback("good")},
			},
			{
				{"text": "Hard", "callback_data": callback("hard")},
				{"text": "Bad", "callback_data": callback("bad")},
			},
		}}
	case "review":
		return map[string]any{"inline_keyboard": [][]map[string]any{
			{
				{"text": "Tomorrow", "callback_data": callback("tomorrow")},
				{"text": "3 days", "callback_data": callback("3_days")},
			},
			{
				{"text": "1 week", "callback_data": callback("1_week")},
				{"text": "No review", "callback_data": callback("no_review")},
			},
		}}
	default:
		return nil
	}
}

func tutorLessonKeyboard(copy uiCopy) map[string]any {
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "🤖 " + copy.AITutor, "callback_data": "menu_tutor"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func formatTelegramTutorLesson(lesson tutorLesson, user userState) string {
	copy := ui(user)
	var builder strings.Builder
	title := strings.TrimSpace(lesson.Title)
	if title == "" {
		title = copy.AITutor
	}
	builder.WriteString("*" + escapeMarkdownV2(title) + "*")
	if lesson.LessonNumber > 0 && lesson.CourseSize > 0 {
		builder.WriteString(" · " + escapeMarkdownV2(fmt.Sprintf("%d/%d", lesson.LessonNumber, lesson.CourseSize)))
	}
	appendTelegramTutorLine(&builder, "🎯", lesson.Goal)
	appendTelegramTutorLine(&builder, "📍", lesson.Scenario)
	appendTelegramTutorLine(&builder, "🧩", lesson.MiniExplanation)
	if strings.TrimSpace(lesson.TeachingPoint.ModelAnswer) != "" {
		appendTelegramTutorLine(&builder, "Model:", lesson.TeachingPoint.ModelAnswer)
	}
	target := strings.TrimSpace(lesson.PronunciationText)
	if target == "" {
		target = strings.TrimSpace(lesson.ScenarioSlots.ModelAnswer)
	}
	appendTelegramTutorLine(&builder, "🗣", target)
	appendTelegramTutorChoice(&builder, lesson)
	appendTelegramTutorDialogue(&builder, lesson)
	appendTelegramTutorFinalCheck(&builder, lesson)
	return builder.String()
}

func appendTelegramTutorLine(builder *strings.Builder, prefix string, text string) {
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	builder.WriteString("\n\n")
	builder.WriteString(prefix)
	builder.WriteString(" ")
	builder.WriteString(escapeMarkdownV2(text))
}

func appendTelegramTutorChoice(builder *strings.Builder, lesson tutorLesson) {
	checks := lesson.Checks
	if len(checks) == 0 && len(lesson.Choice.Options) > 0 {
		checks = []tutorLessonChoice{lesson.Choice}
	}
	if len(checks) == 0 {
		return
	}
	check := checks[0]
	prompt := strings.TrimSpace(check.Prompt)
	if prompt == "" {
		return
	}
	builder.WriteString("\n\n")
	builder.WriteString("✅ ")
	builder.WriteString(escapeMarkdownV2(prompt))
	for _, option := range check.Options {
		text := strings.TrimSpace(option.Text)
		if text == "" {
			continue
		}
		marker := "•"
		if option.ID == check.CorrectAnswerID {
			marker = "→"
		}
		builder.WriteString("\n")
		builder.WriteString(escapeMarkdownV2(marker + " " + text))
		if option.ID == check.CorrectAnswerID && strings.TrimSpace(option.Why) != "" {
			builder.WriteString("\n")
			builder.WriteString(escapeMarkdownV2("  " + strings.TrimSpace(option.Why)))
		}
	}
	if strings.TrimSpace(check.Feedback) != "" {
		builder.WriteString("\n")
		builder.WriteString(escapeMarkdownV2(strings.TrimSpace(check.Feedback)))
	}
}

func appendTelegramTutorDialogue(builder *strings.Builder, lesson tutorLesson) {
	if len(lesson.Dialogue) == 0 && len(lesson.DialogueVariants) == 0 {
		return
	}
	builder.WriteString("\n\n")
	builder.WriteString("💬 ")
	if len(lesson.Dialogue) > 0 {
		builder.WriteString(escapeMarkdownV2(strings.TrimSpace(lesson.Dialogue[0])))
	}
	for _, variant := range lesson.DialogueVariants {
		if variant.Avoid || strings.TrimSpace(variant.Text) == "" {
			continue
		}
		builder.WriteString("\n")
		builder.WriteString(escapeMarkdownV2("→ " + strings.TrimSpace(variant.Text)))
		if strings.TrimSpace(variant.Why) != "" {
			builder.WriteString("\n")
			builder.WriteString(escapeMarkdownV2("  " + strings.TrimSpace(variant.Why)))
		}
		return
	}
}

func appendTelegramTutorFinalCheck(builder *strings.Builder, lesson tutorLesson) {
	if len(lesson.FinalWordCheck.Items) == 0 {
		return
	}
	builder.WriteString("\n\n")
	builder.WriteString(escapeMarkdownV2("Final word check"))
	limit := tutorMinInt(3, len(lesson.FinalWordCheck.Items))
	for index := 0; index < limit; index++ {
		item := lesson.FinalWordCheck.Items[index]
		text := strings.TrimSpace(item.Prompt)
		if text == "" {
			continue
		}
		builder.WriteString("\n")
		builder.WriteString(escapeMarkdownV2("- " + text))
	}
}

func (b *bot) startWordLesson(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	word, ok := nextUnlearnedWord(user)
	prompt := ""
	promptLanguage := vocabularyGenerationPromptLanguage(user)
	if ok {
		prompt = wordTranslation(word, promptLanguage.Code)
		if prompt == "" {
			prompt = b.vocabularyPromptInLanguage(ctx, user, word, promptLanguage.Code, "learn word")
		}
		if prompt == "" {
			if generated, generatedOK := b.generateVocabularyWord(ctx, user, "learn word"); generatedOK {
				word = generated
				prompt = wordTranslation(word, promptLanguage.Code)
			}
		}
	}
	if !ok {
		if generated, generatedOK := b.generateVocabularyWord(ctx, user, "learn word"); generatedOK {
			word = generated
			ok = true
			prompt = wordTranslation(word, promptLanguage.Code)
		}
	}
	if ok && prompt == "" {
		ok = false
	}
	if !ok {
		if !vocabularyLanguageHasAny(user.LearningLanguage) {
			language := userLearningLanguage(user)
			copy := ui(user)
			return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
				"*"+escapeMarkdownV2(copy.LearnWords)+"* 🧠\n\n"+escapeMarkdownV2(copy.LearningLanguage)+": *"+escapeMarkdownV2(language.InterfaceName)+"*",
				backToMenuKeyboard(copy))
		}
		copy := ui(user)
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			"*"+escapeMarkdownV2(copy.LearnWords)+"* 🧠\n\n"+escapeMarkdownV2(copy.WordGame),
			backToMenuKeyboard(copy))
	}

	if prompt == "" {
		prompt = wordTranslation(word, promptLanguage.Code)
	}

	language := userLearningLanguage(user)
	copy := ui(user)
	hint := ""
	if strings.EqualFold(strings.TrimSpace(word.Source), "ai") {
		hint = normalizeVocabularyContext(word.Context)
	}
	if hint == "" {
		hint = b.vocabularyHint(ctx, user, word, "learn word")
	}
	if hint == "" {
		hint = wordContext(word, promptLanguage.Code)
	}
	text := "*" + escapeMarkdownV2(copy.LearnWords) + "* 🧠\n\n" +
		escapeMarkdownV2(fmt.Sprintf(copy.WordQuestion, learningDirectionForUI(user, language))) +
		"\n" + studyMarkdownBlock(prompt, hint) + "\n\n" +
		escapeMarkdownV2(copy.ChooseAnswer)
	if err := b.telegram.sendInlineMarkdownMessage(ctx, chatID, text, wordChoiceKeyboard("wl", word.ID, wordOptions(word), user.InterfaceLanguage)); err != nil {
		return err
	}
	return nil
}

func wordStudyMarkdownBlock(word vocabWord, interfaceLanguage string) string {
	return studyMarkdownBlock(wordTranslation(word, interfaceLanguage), wordContext(word, interfaceLanguage))
}

func wordStudyMarkdownBlockWithHint(word vocabWord, interfaceLanguage string, hint string) string {
	return studyMarkdownBlock(wordTranslation(word, interfaceLanguage), hint)
}

func learnedWordStudyMarkdownBlock(word learnedWordEntry, interfaceLanguage string) string {
	return studyMarkdownBlock(learnedWordTranslation(word, interfaceLanguage), learnedWordContext(word, interfaceLanguage))
}

func studyMarkdownBlock(translation string, context string) string {
	translation = cleanDictionaryDisplay(translation)
	context = normalizeVocabularyContext(context)
	block := "`" + escapeMarkdownV2(translation) + "`"
	if context != "" && context != translation {
		block += "\n" + escapeMarkdownV2(context)
	}
	return block
}

func (b *bot) vocabularyHint(ctx context.Context, user userState, word vocabWord, mode string) string {
	if b == nil || b.openrouter == nil || strings.TrimSpace(b.cfg.OpenRouterVocabularyModel) == "" {
		return ""
	}
	prompt := wordTranslation(word, user.InterfaceLanguage)
	if prompt == "" {
		return ""
	}
	learningLanguage := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	cacheKey := strings.Join([]string{
		b.cfg.OpenRouterVocabularyModel,
		normalizeInterfaceLanguage(user.InterfaceLanguage),
		normalizeLearningLanguage(word.Language),
		mode,
		word.ID,
		prompt,
	}, "|")

	b.vocabularyHintMu.Lock()
	if b.vocabularyHints != nil {
		if cached := b.vocabularyHints[cacheKey]; cached != "" {
			b.vocabularyHintMu.Unlock()
			return cached
		}
	}
	b.vocabularyHintMu.Unlock()

	hint, err := b.openrouter.completeWithModel(ctx, b.cfg.OpenRouterVocabularyModel, vocabularyHintPrompt(word, prompt, learningLanguage, interfaceLanguage, mode), 0.1, 80)
	if err != nil {
		log.Printf("failed to build vocabulary hint for %s: %v", word.ID, err)
		return ""
	}
	hint = sanitizeVocabularyHint(hint, word)
	if hint == "" {
		return ""
	}

	b.vocabularyHintMu.Lock()
	if b.vocabularyHints == nil {
		b.vocabularyHints = map[string]string{}
	}
	if len(b.vocabularyHints) > 2048 {
		b.vocabularyHints = map[string]string{}
	}
	b.vocabularyHints[cacheKey] = hint
	b.vocabularyHintMu.Unlock()
	return hint
}

func (b *bot) generateVocabularyWord(ctx context.Context, user userState, mode string) (vocabWord, bool) {
	if b == nil || b.openrouter == nil || strings.TrimSpace(b.cfg.OpenRouterVocabularyModel) == "" || currentSQLiteVocabularyDB() == nil {
		return vocabWord{}, false
	}
	language := normalizeLearningLanguage(user.LearningLanguage)
	promptLanguage := vocabularyGenerationPromptLanguage(user)
	extraForbidden := []string{}
	for attempt := 0; attempt < 4; attempt++ {
		forbidden, forbiddenList := vocabularyGenerationForbiddenWords(user, extraForbidden)
		messages := vocabularyGenerationPrompt(user, promptLanguage, forbiddenList)
		raw, err := b.openrouter.completeWithModel(ctx, b.cfg.OpenRouterVocabularyModel, messages, 0.2, 240)
		if err != nil {
			log.Printf("failed to generate vocabulary word for %s/%s: %v", language, normalizeCEFRLevel(user.Level), err)
			return vocabWord{}, false
		}
		word, ok := parseGeneratedVocabularyWord(raw, language, promptLanguage.Code, user.Level, forbidden)
		if !ok {
			log.Printf("rejected generated vocabulary word for %s/%s", language, normalizeCEFRLevel(user.Level))
			continue
		}
		exists, err := sqliteVocabularyWordExists(language, word.English)
		if err != nil {
			log.Printf("failed to check generated vocabulary duplicate %s: %v", word.ID, err)
			return vocabWord{}, false
		}
		if exists {
			extraForbidden = append(extraForbidden, word.English)
			continue
		}
		inserted, err := sqliteVocabularyAIWordSet(word, b.cfg.OpenRouterVocabularyModel, chatMessagesPromptText(messages))
		if err != nil {
			log.Printf("failed to persist generated vocabulary word %s: %v", word.ID, err)
			return vocabWord{}, false
		}
		if !inserted {
			extraForbidden = append(extraForbidden, word.English)
			continue
		}
		return word, true
	}
	return vocabWord{}, false
}

func chatMessagesPromptText(messages []chatMessage) string {
	var builder strings.Builder
	for _, message := range messages {
		if strings.TrimSpace(message.Content) == "" {
			continue
		}
		if builder.Len() > 0 {
			builder.WriteString("\n\n")
		}
		builder.WriteString(strings.TrimSpace(message.Role))
		builder.WriteString(": ")
		builder.WriteString(strings.TrimSpace(message.Content))
	}
	return builder.String()
}

func (b *bot) vocabularyPrompt(ctx context.Context, user userState, word vocabWord, mode string) string {
	targetLanguage := normalizeInterfaceLanguage(user.InterfaceLanguage)
	if targetLanguage == "" {
		targetLanguage = "en"
	}
	if prompt := b.vocabularyPromptInLanguage(ctx, user, word, targetLanguage, mode); prompt != "" {
		return prompt
	}
	return vocabularyFallbackPrompt(word, targetLanguage)
}

func (b *bot) vocabularyPromptInLanguage(ctx context.Context, user userState, word vocabWord, targetLanguage string, mode string) string {
	targetLanguage = normalizeInterfaceLanguage(targetLanguage)
	if targetLanguage == "" {
		targetLanguage = "en"
	}
	if cached, ok, err := sqliteVocabularyAITranslationGet(word.ID, targetLanguage); err == nil && ok {
		if value := sanitizeVocabularyTranslation(cached, word); value != "" {
			return firstDictionaryValue(value)
		}
	}
	if translation := wordTranslation(word, targetLanguage); translation != "" {
		return firstDictionaryValue(translation)
	}
	if b != nil && b.openrouter != nil && strings.TrimSpace(b.cfg.OpenRouterVocabularyModel) != "" {
		raw, err := b.openrouter.completeWithModel(ctx, b.cfg.OpenRouterVocabularyModel, vocabularyTranslationPrompt(word, userLearningLanguage(user), interfaceLanguageByCode(targetLanguage), mode), 0.1, 160)
		if err == nil {
			if value := sanitizeVocabularyTranslation(raw, word); value != "" {
				_ = sqliteVocabularyAITranslationSet(word, targetLanguage, b.cfg.OpenRouterVocabularyModel, mode, value)
				return firstDictionaryValue(value)
			}
		}
	}
	return ""
}

func vocabularyFallbackPrompt(word vocabWord, interfaceLanguage string) string {
	if translation := wordTranslation(word, interfaceLanguage); translation != "" {
		return firstDictionaryValue(translation)
	}
	for _, fallbackLanguage := range []string{"ru", "en"} {
		if fallbackLanguage == normalizeLearningLanguage(word.Language) {
			continue
		}
		if translation := wordTranslationForLanguage(word, fallbackLanguage); translation != "" {
			return firstDictionaryValue(translation)
		}
	}
	return ""
}

func (b *bot) vocabularyExample(ctx context.Context, user userState, word vocabWord, mode string) string {
	if b == nil || b.openrouter == nil || strings.TrimSpace(b.cfg.OpenRouterVocabularyModel) == "" {
		return ""
	}
	if strings.TrimSpace(word.English) == "" {
		return ""
	}
	learningLanguage := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	cacheKey := strings.Join([]string{
		b.cfg.OpenRouterVocabularyModel,
		normalizeInterfaceLanguage(user.InterfaceLanguage),
		normalizeLearningLanguage(word.Language),
		"example",
		mode,
		strings.TrimSpace(user.Level),
		word.ID,
		word.English,
	}, "|")

	b.vocabularyHintMu.Lock()
	if b.vocabularyHints != nil {
		if cached := b.vocabularyHints[cacheKey]; cached != "" {
			b.vocabularyHintMu.Unlock()
			return cached
		}
	}
	b.vocabularyHintMu.Unlock()

	example, err := b.openrouter.completeWithModel(ctx, b.cfg.OpenRouterVocabularyModel, vocabularyExamplePrompt(word, learningLanguage, interfaceLanguage, user.Level, mode), 0.2, 100)
	if err != nil {
		log.Printf("failed to build vocabulary example for %s: %v", word.ID, err)
		return ""
	}
	example = sanitizeVocabularyExample(example)
	if example == "" {
		return ""
	}

	b.vocabularyHintMu.Lock()
	if b.vocabularyHints == nil {
		b.vocabularyHints = map[string]string{}
	}
	if len(b.vocabularyHints) > 2048 {
		b.vocabularyHints = map[string]string{}
	}
	b.vocabularyHints[cacheKey] = example
	b.vocabularyHintMu.Unlock()
	return example
}

func sanitizeVocabularyHint(hint string, word vocabWord) string {
	hint = cleanDictionaryDisplay(hint)
	if hint == "" || dictionaryTextContainsTerm(hint, word.English) {
		return ""
	}
	hint = strings.Trim(hint, `"'“”„«»`)
	lines := strings.FieldsFunc(hint, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	if len(lines) > 0 {
		hint = strings.TrimSpace(lines[0])
	}
	words := strings.Fields(hint)
	if len(words) > 12 {
		hint = strings.Join(words[:12], " ")
	}
	return strings.TrimSpace(hint)
}

func sanitizeVocabularyTranslation(raw string, word vocabWord) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var payload struct {
		Translation string `json:"translation"`
	}
	if strings.HasPrefix(raw, "{") {
		if err := json.Unmarshal([]byte(raw), &payload); err == nil {
			raw = payload.Translation
		}
	}
	raw = cleanDictionaryDisplay(raw)
	if raw == "" {
		return ""
	}
	parts := strings.Split(raw, ";")
	kept := make([]string, 0, len(parts))
	for _, part := range parts {
		part = cleanDictionaryDisplay(part)
		if part == "" || dictionaryTextContainsTerm(part, word.English) {
			continue
		}
		kept = append(kept, part)
	}
	return strings.Join(kept, "; ")
}

func sanitizeVocabularyExample(example string) string {
	example = cleanDictionaryDisplay(cleanModelReply(example))
	example = strings.Trim(example, " \t\r\n\"'")
	if example == "" {
		return ""
	}
	lines := strings.FieldsFunc(example, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	if len(lines) > 0 {
		example = strings.TrimSpace(lines[0])
	}
	runes := []rune(example)
	if len(runes) > 140 {
		example = strings.TrimSpace(string(runes[:140]))
	}
	return example
}

func correctionPronunciationText(entries []mistakeEntry) string {
	seen := map[string]bool{}
	corrections := make([]string, 0, len(entries))
	for _, entry := range entries {
		correction := strings.TrimSpace(entry.Correction)
		if correction == "" || seen[correction] {
			continue
		}
		seen[correction] = true
		corrections = append(corrections, correction)
	}
	text := strings.Join(corrections, ". ")
	runes := []rune(text)
	if len(runes) > 240 {
		text = strings.TrimSpace(string(runes[:240]))
	}
	return text
}

func lessonCorrectionPronunciationText(feedback string, entries []mistakeEntry) string {
	if text := firstFeedbackSection(feedback); text != "" {
		return text
	}
	return correctionPronunciationText(entries)
}

func practiceCorrectionPronunciationText(reply string, entries []mistakeEntry) string {
	if text := firstStudyLabelAudioLine(reply, studyModelPhraseLabels()); text != "" {
		return text
	}
	if text := firstStudyLabelAudioLine(reply, studyCorrectionLabels()); text != "" {
		return text
	}
	return correctionPronunciationText(entries)
}

func studyModelPhraseLabels() []string {
	return studyLabels(func(labels studyPromptLabels) string { return labels.ModelPhrase }, "model phrase", "model answer", "example phrase")
}

func studyCorrectionLabels() []string {
	return studyLabels(func(labels studyPromptLabels) string { return labels.Correction }, "correction", "corrected version")
}

func studyAllAudioLabels() []string {
	labels := studyLabels(func(labels studyPromptLabels) string { return labels.Correction }, "correction", "corrected version")
	labels = append(labels, studyLabels(func(labels studyPromptLabels) string { return labels.ModelPhrase }, "model phrase", "model answer", "example phrase")...)
	labels = append(labels, studyLabels(func(labels studyPromptLabels) string { return labels.YourTurn }, "your turn", "question")...)
	return labels
}

func studyLabels(pick func(studyPromptLabels) string, extras ...string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(studyPromptLabelsByLanguage)+len(extras))
	add := func(value string) {
		value = strings.ToLower(strings.Trim(strings.TrimSpace(value), ":："))
		if value == "" || seen[value] {
			return
		}
		seen[value] = true
		result = append(result, value)
	}
	for _, labels := range studyPromptLabelsByLanguage {
		add(pick(labels))
	}
	for _, extra := range extras {
		add(extra)
	}
	return result
}

func firstStudyLabelAudioLine(reply string, labels []string) string {
	reply = strings.ReplaceAll(reply, "\r\n", "\n")
	reply = strings.ReplaceAll(reply, "\r", "\n")
	lines := strings.Split(reply, "\n")
	for i, line := range lines {
		if rest, ok := audioLineAfterLabel(line, labels); ok {
			if text := sanitizeAudioText(rest); text != "" {
				return text
			}
			for j := i + 1; j < len(lines) && j <= i+3; j++ {
				candidate := strings.TrimSpace(lines[j])
				if candidate == "" {
					continue
				}
				if _, isNextLabel := audioLineAfterLabel(candidate, studyAllAudioLabels()); isNextLabel {
					break
				}
				if text := sanitizeAudioText(candidate); text != "" {
					return text
				}
			}
		}
	}
	return ""
}

func audioLineAfterLabel(line string, labels []string) (string, bool) {
	cleaned := strings.TrimSpace(stripAudioLineBullet(line))
	lower := strings.ToLower(cleaned)
	for _, label := range labels {
		if label == "" {
			continue
		}
		if lower == label {
			return "", true
		}
		for _, sep := range []string{":", "：", "-", "—", "–"} {
			prefix := label + sep
			if strings.HasPrefix(lower, prefix) {
				cleanedRunes := []rune(cleaned)
				prefixLen := len([]rune(prefix))
				if prefixLen > len(cleanedRunes) {
					return "", true
				}
				return strings.TrimSpace(string(cleanedRunes[prefixLen:])), true
			}
		}
	}
	return "", false
}

func stripAudioLineBullet(line string) string {
	line = strings.TrimLeft(line, " \t-*•—–")
	runes := []rune(line)
	if len(runes) >= 2 && runes[0] >= '1' && runes[0] <= '9' && (runes[1] == ')' || runes[1] == '.' || runes[1] == ':' || runes[1] == '-') {
		return strings.TrimSpace(string(runes[2:]))
	}
	return strings.TrimSpace(line)
}

func firstFeedbackSection(feedback string) string {
	feedback = strings.ReplaceAll(feedback, "\r\n", "\n")
	feedback = strings.ReplaceAll(feedback, "\r", "\n")
	lines := strings.Split(feedback, "\n")
	collecting := false
	collected := []string{}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch numberedSectionMarker(trimmed) {
		case 1:
			collecting = true
			if rest := strings.TrimSpace(trimmed[sectionMarkerLen(trimmed):]); rest != "" {
				collected = append(collected, rest)
			}
			continue
		case 2, 3:
			if collecting {
				return sanitizeAudioText(strings.Join(collected, "\n"))
			}
		}
		if collecting {
			collected = append(collected, line)
		}
	}
	if !collecting {
		return ""
	}
	return sanitizeAudioText(strings.Join(collected, "\n"))
}

func numberedSectionMarker(line string) int {
	if line == "" {
		return 0
	}
	runes := []rune(line)
	if len(runes) < 2 || runes[0] < '1' || runes[0] > '3' {
		return 0
	}
	if runes[1] == ')' || runes[1] == '.' || runes[1] == ':' || runes[1] == '-' || runes[1] == ' ' {
		return int(runes[0] - '0')
	}
	return 0
}

func sectionMarkerLen(line string) int {
	runes := []rune(line)
	if len(runes) == 0 {
		return 0
	}
	if len(runes) >= 2 && (runes[1] == ')' || runes[1] == '.' || runes[1] == ':' || runes[1] == '-') {
		return len(string(runes[:2]))
	}
	return len(string(runes[:1]))
}

func practiceQuestionPronunciationText(reply string) string {
	reply = strings.ReplaceAll(reply, "\r\n", "\n")
	reply = strings.ReplaceAll(reply, "\r", "\n")
	lines := strings.Split(reply, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if question := lastQuestionSentence(line); question != "" {
			return sanitizeAudioText(question)
		}
	}
	return ""
}

func lessonQuestionPronunciationText(task string) string {
	if text := firstStudyLabelAudioLine(task, studyLabels(func(labels studyPromptLabels) string { return labels.ModelPhrase }, "example phrase", "phrase example", "target phrase")); text != "" {
		return text
	}
	return ""
}

func roleplayQuestionPronunciationText(reply string) string {
	reply = strings.ReplaceAll(reply, "\r\n", "\n")
	reply = strings.ReplaceAll(reply, "\r", "\n")
	lines := strings.Split(reply, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if _, isLabel := audioLineAfterLabel(line, studyAllAudioLabels()); isLabel {
			continue
		}
		if speaker, text, ok := strings.Cut(line, ":"); ok && strings.TrimSpace(speaker) != "" {
			if audio := sanitizeAudioText(text); audio != "" {
				return audio
			}
		}
	}
	return ""
}

func lastQuestionSentence(line string) string {
	runes := []rune(line)
	questionEnd := -1
	for i := len(runes) - 1; i >= 0; i-- {
		if isQuestionMark(runes[i]) {
			questionEnd = i
			break
		}
	}
	if questionEnd < 0 {
		return ""
	}
	start := 0
	for i := questionEnd - 1; i >= 0; i-- {
		if isSentenceBoundary(runes[i]) {
			start = i + 1
			break
		}
	}
	return string(runes[start : questionEnd+1])
}

func isQuestionMark(r rune) bool {
	return r == '?' || r == '؟' || r == '？'
}

func isSentenceBoundary(r rune) bool {
	return r == '.' || r == '!' || r == '?' || r == '؟' || r == '。' || r == '！' || r == '？'
}

func sanitizeAudioText(text string) string {
	text = cleanDictionaryDisplay(cleanModelReply(text))
	text = strings.Trim(text, " \t\r\n\"'`")
	lines := strings.FieldsFunc(text, func(r rune) bool {
		return r == '\n' || r == '\r'
	})
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = stripAudioLinePrefix(strings.TrimSpace(line))
		if line == "" || isAudioHeading(line) {
			continue
		}
		cleaned = append(cleaned, line)
	}
	text = strings.Join(cleaned, " ")
	text = stripTrailingTranslationParenthetical(text)
	runes := []rune(text)
	if len(runes) > 240 {
		text = strings.TrimSpace(string(runes[:240]))
	}
	return text
}

func stripTrailingTranslationParenthetical(text string) string {
	text = strings.TrimSpace(text)
	for {
		open := strings.LastIndex(text, "(")
		if open < 0 || !strings.HasSuffix(text, ")") {
			return text
		}
		before := strings.TrimSpace(text[:open])
		inside := strings.TrimSpace(text[open+1 : len(text)-1])
		if before == "" || inside == "" || !looksLikeTranslationParenthetical(before, inside) {
			return text
		}
		text = before
	}
}

func looksLikeTranslationParenthetical(before, inside string) bool {
	if !hasLetters(before) || !hasLetters(inside) {
		return false
	}
	beforeScript := dominantLetterScript(before)
	insideScript := dominantLetterScript(inside)
	if beforeScript != "" && insideScript != "" && beforeScript != insideScript {
		return true
	}
	return hasTerminalSentencePunctuation(before) && parentheticalLooksLikeSentence(inside)
}

func hasLetters(text string) bool {
	for _, r := range text {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func dominantLetterScript(text string) string {
	counts := map[string]int{}
	for _, r := range text {
		if !unicode.IsLetter(r) {
			continue
		}
		switch {
		case isLatinRune(r):
			counts["latin"]++
		case isRuneIn(r, 0x0400, 0x052F) || isRuneIn(r, 0x2DE0, 0x2DFF) || isRuneIn(r, 0xA640, 0xA69F):
			counts["cyrillic"]++
		case isRuneIn(r, 0x4E00, 0x9FFF) || isRuneIn(r, 0x3400, 0x4DBF):
			counts["cjk"]++
		case isRuneIn(r, 0x3040, 0x30FF):
			counts["kana"]++
		case isRuneIn(r, 0xAC00, 0xD7AF) || isRuneIn(r, 0x1100, 0x11FF):
			counts["hangul"]++
		case isRuneIn(r, 0x0600, 0x06FF) || isRuneIn(r, 0x0750, 0x077F) || isRuneIn(r, 0x08A0, 0x08FF):
			counts["arabic"]++
		case isRuneIn(r, 0x0590, 0x05FF):
			counts["hebrew"]++
		case isRuneIn(r, 0x0530, 0x058F):
			counts["armenian"]++
		case isRuneIn(r, 0x10A0, 0x10FF):
			counts["georgian"]++
		case isRuneIn(r, 0x0370, 0x03FF):
			counts["greek"]++
		case isRuneIn(r, 0x0E00, 0x0E7F):
			counts["thai"]++
		case isRuneIn(r, 0x0900, 0x097F):
			counts["devanagari"]++
		default:
			counts["other"]++
		}
	}
	bestScript := ""
	bestCount := 0
	for script, count := range counts {
		if count > bestCount {
			bestScript = script
			bestCount = count
		}
	}
	return bestScript
}

func isLatinRune(r rune) bool {
	return isRuneIn(r, 'A', 'Z') ||
		isRuneIn(r, 'a', 'z') ||
		isRuneIn(r, 0x00C0, 0x024F) ||
		isRuneIn(r, 0x1E00, 0x1EFF)
}

func isRuneIn(r rune, from, to rune) bool {
	return r >= from && r <= to
}

func hasTerminalSentencePunctuation(text string) bool {
	runes := []rune(strings.TrimSpace(text))
	for i := len(runes) - 1; i >= 0; i-- {
		r := runes[i]
		if unicode.IsSpace(r) || r == '"' || r == '\'' || r == '`' {
			continue
		}
		return r == '.' || r == '!' || r == '?' || r == '\u3002' || r == '\uFF01' || r == '\uFF1F' || r == '\u061F' || r == '\u055C'
	}
	return false
}

func parentheticalLooksLikeSentence(text string) bool {
	letterCount := 0
	wordCount := 0
	inWord := false
	for _, r := range text {
		if unicode.IsLetter(r) {
			letterCount++
			if !inWord {
				wordCount++
				inWord = true
			}
			continue
		}
		inWord = false
	}
	return letterCount >= 3 && (wordCount >= 2 || len([]rune(strings.TrimSpace(text))) > 12)
}

func stripAudioLinePrefix(line string) string {
	line = strings.TrimLeft(line, " \t-*•—–")
	runes := []rune(line)
	if len(runes) >= 2 && runes[0] >= '1' && runes[0] <= '9' && (runes[1] == ')' || runes[1] == '.' || runes[1] == ':' || runes[1] == '-') {
		line = strings.TrimSpace(string(runes[2:]))
	}
	if idx := strings.Index(line, ":"); idx >= 0 && idx < 80 {
		prefix := strings.ToLower(strings.TrimSpace(line[:idx]))
		if audioLabelPrefix(prefix) {
			line = strings.TrimSpace(line[idx+1:])
		}
	}
	return strings.TrimSpace(line)
}

func audioLabelPrefix(prefix string) bool {
	for _, token := range []string{
		"correct", "correction", "corrected", "version", "question",
		"исправ", "вариант", "вопрос",
		"pregunta", "frage", "domanda", "pergunta",
	} {
		if strings.Contains(prefix, token) {
			return true
		}
	}
	return false
}

func isAudioHeading(line string) bool {
	trimmed := strings.TrimSuffix(strings.ToLower(strings.TrimSpace(line)), ":")
	if len([]rune(trimmed)) > 80 {
		return false
	}
	return audioLabelPrefix(trimmed)
}

func wordExampleMarkdown(example string) string {
	example = strings.TrimSpace(example)
	if example == "" {
		return ""
	}
	return "\n\n_" + escapeMarkdownV2(example) + "_"
}

func addedToPracticePoolText(user userState) string {
	switch normalizeInterfaceLanguage(user.InterfaceLanguage) {
	case "ru":
		return "Слово добавлено в повторяйку и правописание. В словарь попадет после 10 успешных ответов."
	case "en":
		return "Word added to review and spelling practice. It enters the vocabulary after 10 correct answers."
	default:
		return "Word added to review and spelling practice. It enters the vocabulary after 10 correct answers."
	}
}

func alreadyInPracticePoolText(user userState) string {
	switch normalizeInterfaceLanguage(user.InterfaceLanguage) {
	case "ru":
		return "Это слово уже есть в тренировках."
	case "en":
		return "This word is already in practice."
	default:
		return "This word is already in practice."
	}
}

func learningDirectionForUI(user userState, language learningLanguage) string {
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		return language.PracticeDirection
	}
	return language.InterfaceName
}

func (b *bot) startWordGame(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	ids := reviewWordIDs(user)
	if len(ids) == 0 {
		copy := ui(user)
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			"*"+escapeMarkdownV2(copy.WordGame)+"* 🎮\n\n"+escapeMarkdownV2(copy.LearnWords),
			backToMenuKeyboard(copy))
	}

	word, ok := findVocabWord(ids[(user.WordGameCount+len(ids))%len(ids)])
	if !ok {
		word, ok = findVocabWord(ids[0])
	}
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeText(user, "word_game_not_found"), ui(user))
	}

	language := userLearningLanguage(user)
	copy := ui(user)
	text := "*" + escapeMarkdownV2(copy.WordGame) + "* 🎮\n\n" +
		escapeMarkdownV2(fmt.Sprintf(copy.WordQuestion, learningDirectionForUI(user, language))) +
		"\n" + wordStudyMarkdownBlockWithHint(word, user.InterfaceLanguage, b.vocabularyHint(ctx, user, word, "review word"))
	if err := b.telegram.sendInlineMarkdownMessage(ctx, chatID, text, wordChoiceKeyboard("wg", word.ID, learnedWordOptions(word, ids), user.InterfaceLanguage)); err != nil {
		return err
	}
	return nil
}

func (b *bot) deletePreviousWordPronunciation(ctx context.Context, chatID int64) {
	b.pronunciationMu.Lock()
	messageIDs := []int64(nil)
	if b.pronunciationMessages != nil {
		messageIDs = append(messageIDs, b.pronunciationMessages[chatID]...)
		delete(b.pronunciationMessages, chatID)
	}
	b.pronunciationMu.Unlock()

	if len(messageIDs) == 0 {
		return
	}
	for _, messageID := range messageIDs {
		if messageID == 0 {
			continue
		}
		if err := b.telegram.deleteMessage(ctx, chatID, messageID); err != nil {
			log.Printf("failed to delete previous pronunciation message %d: %v", messageID, err)
		}
	}
}

func (b *bot) sendWordPronunciation(ctx context.Context, chatID int64, user userState, word vocabWord) {
	b.sendTextPronunciation(ctx, chatID, user, word.English, legacyVocabID(word.ID)+".mp3", word.English, word.ID)
}

func (b *bot) sendTextPronunciation(ctx context.Context, chatID int64, user userState, text string, filename string, title string, logKey string) {
	if b.openrouter == nil || strings.TrimSpace(b.cfg.OpenRouterTTSModel) == "" {
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	audio, err := b.openrouter.synthesizeSpeech(ctx, b.cfg.OpenRouterTTSModel, b.cfg.OpenRouterTTSVoice, text)
	if err != nil {
		log.Printf("failed to synthesize pronunciation for %s: %v", logKey, err)
		return
	}
	if strings.TrimSpace(filename) == "" {
		filename = "pronunciation.mp3"
	}
	if strings.TrimSpace(title) == "" {
		title = text
	}
	messageID, err := b.telegram.sendAudioBytesWithCopy(ctx, chatID, audio, filename, title, ui(user))
	if err != nil {
		log.Printf("failed to send pronunciation for %s: %v", logKey, err)
		return
	}
	if messageID != 0 {
		b.pronunciationMu.Lock()
		if b.pronunciationMessages == nil {
			b.pronunciationMessages = map[int64][]int64{}
		}
		b.pronunciationMessages[chatID] = append(b.pronunciationMessages[chatID], messageID)
		b.pronunciationMu.Unlock()
	}
}

func (b *bot) handleWordLessonChoice(ctx context.Context, chatID int64, user userState, data string) error {
	correctID, answerID, ok := parseChoiceCallback(data)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
	word, ok := findVocabWord(correctID)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeText(user, "word_task_missing"), ui(user))
	}
	if answerID != correctID {
		copy := ui(user)
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			escapeMarkdownV2(copy.TryAgain)+"\n\n"+wordStudyMarkdownBlock(word, user.InterfaceLanguage),
			wordChoiceKeyboard("wl", word.ID, wordOptions(word), user.InterfaceLanguage))
	}

	learned, total, err := b.store.addLearnedWord(user.TelegramID, word)
	if err != nil {
		return err
	}
	copy := ui(user)
	example := b.vocabularyExample(ctx, user, word, "learn word")
	message := escapeMarkdownV2(copy.Correct) + " 🎉\n\n" + wordStudyMarkdownBlock(word, user.InterfaceLanguage) + "\n*" + escapeMarkdownV2(word.English) + "*"
	message += wordExampleMarkdown(example)
	if learned {
		message += "\n\n" + escapeMarkdownV2(addedToPracticePoolText(user)) + "\n\\+15 XP"
	} else {
		message += "\n\n" + escapeMarkdownV2(alreadyInPracticePoolText(user))
	}
	message += "\n\n" + escapeMarkdownV2(fmt.Sprintf(copy.TotalLearned, total))
	if err := b.telegram.sendInlineMarkdownMessage(ctx, chatID, message, wordLessonDoneKeyboard(ui(user))); err != nil {
		return err
	}
	b.sendWordPronunciation(ctx, chatID, user, word)
	return b.maybePromoteLearningLevel(ctx, chatID, user.TelegramID, user.FirstName)
}

func (b *bot) handleWordGameChoice(ctx context.Context, chatID int64, user userState, data string) error {
	correctID, answerID, ok := parseChoiceCallback(data)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
	word, ok := findVocabWord(correctID)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeText(user, "word_game_missing"), ui(user))
	}
	ids := reviewWordIDs(user)
	if answerID != correctID {
		copy := ui(user)
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			escapeMarkdownV2(copy.TryAgain)+"\n\n"+wordStudyMarkdownBlock(word, user.InterfaceLanguage),
			wordChoiceKeyboard("wg", word.ID, learnedWordOptions(word, ids), user.InterfaceLanguage))
	}

	masteredNow := wordWillBeMasteredAfterCorrect(user, word.ID)
	if err := b.store.markWordGameCorrect(user.TelegramID, word.ID); err != nil {
		return err
	}
	copy := ui(user)
	example := b.vocabularyExample(ctx, user, word, "review word")
	masteredText := ""
	if masteredNow {
		masteredText = "\n\n" + escapeMarkdownV2(copy.AddedToVocab)
	}
	if err := b.telegram.sendInlineMarkdownMessage(ctx, chatID,
		escapeMarkdownV2(copy.Correct)+" 🎯\n\n"+wordStudyMarkdownBlock(word, user.InterfaceLanguage)+"\n*"+escapeMarkdownV2(word.English)+"*"+wordExampleMarkdown(example)+masteredText+"\n\n\\+10 XP",
		wordGameDoneKeyboard(ui(user))); err != nil {
		return err
	}
	b.sendWordPronunciation(ctx, chatID, user, word)
	return b.maybePromoteLearningLevel(ctx, chatID, user.TelegramID, user.FirstName)
}

func (b *bot) startSpellingPractice(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	ids := spellingWordIDs(user)
	if len(ids) == 0 {
		copy := ui(user)
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			"*"+escapeMarkdownV2(copy.Spelling)+"* ✍️\n\n"+escapeMarkdownV2(copy.LearnWords),
			backToMenuKeyboard(copy))
	}

	wordID := ids[(user.WordGameCount+user.XP+len(ids))%len(ids)]
	word, ok := findVocabWord(wordID)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeText(user, "spelling_not_found"), ui(user))
	}
	if err := b.store.setMode(user.TelegramID, modeSpellingPrefix+word.ID); err != nil {
		return err
	}

	language := userLearningLanguage(user)
	copy := ui(user)
	text := "*" + escapeMarkdownV2(copy.Spelling) + "* ✍️\n\n" +
		escapeMarkdownV2(fmt.Sprintf(copy.WriteWord, learningDirectionForUI(user, language))) +
		"\n" + wordStudyMarkdownBlockWithHint(word, user.InterfaceLanguage, b.vocabularyHint(ctx, user, word, "spelling practice"))
	if err := b.telegram.sendInlineMarkdownMessage(ctx, chatID, text, spellingPracticeKeyboard(ui(user))); err != nil {
		return err
	}
	return nil
}

func (b *bot) handleSpellingAnswer(ctx context.Context, chatID int64, user userState, answer string) error {
	wordID := strings.TrimPrefix(user.Mode, modeSpellingPrefix)
	word, ok := learnedWordByID(user, wordID)
	if !ok {
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeText(user, "spelling_missing"), ui(user))
	}

	if normalizeAnswer(answer) != normalizeAnswer(word.English) {
		language := userLearningLanguage(user)
		copy := ui(user)
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			escapeMarkdownV2(copy.TryAgain)+"\n\n"+learnedWordStudyMarkdownBlock(word, user.InterfaceLanguage)+"\n"+
				escapeMarkdownV2(fmt.Sprintf(copy.WriteWord, learningDirectionForUI(user, language)))+
				"\n\n"+escapeMarkdownV2(copy.Hint)+": `"+escapeMarkdownV2(spellingHint(word.English))+"`",
			spellingPracticeKeyboard(ui(user)))
	}

	if err := b.store.addXP(user.TelegramID, 7); err != nil {
		return err
	}
	masteredNow := wordWillBeMasteredAfterCorrect(user, word.ID)
	if err := b.store.markSpellingCorrect(user.TelegramID, word.ID); err != nil {
		return err
	}
	if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
		return err
	}
	copy := ui(user)
	pronunciationWord := vocabWord{
		ID:       word.ID,
		Language: word.Language,
		Russian:  word.Russian,
		English:  word.English,
		Context:  word.Context,
	}
	example := b.vocabularyExample(ctx, user, pronunciationWord, "spelling practice")
	masteredText := ""
	if masteredNow {
		masteredText = "\n\n" + escapeMarkdownV2(copy.AddedToVocab)
	}
	if err := b.telegram.sendInlineMarkdownMessage(ctx, chatID,
		escapeMarkdownV2(copy.GoodSpelling)+" ✍️\n\n"+learnedWordStudyMarkdownBlock(word, user.InterfaceLanguage)+"\n*"+escapeMarkdownV2(word.English)+"*"+wordExampleMarkdown(example)+masteredText+"\n\n\\+7 XP",
		spellingDoneKeyboard(ui(user))); err != nil {
		return err
	}
	b.sendWordPronunciation(ctx, chatID, user, pronunciationWord)
	return b.maybePromoteLearningLevel(ctx, chatID, user.TelegramID, user.FirstName)
}

func (b *bot) startMistakePractice(ctx context.Context, chatID int64, user userState) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	mistakes := mistakesForLanguage(user)
	if len(mistakes) == 0 {
		copy := ui(user)
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			"*"+escapeMarkdownV2(copy.Mistakes)+"* 📖\n\n"+escapeMarkdownV2(copy.Practice),
			backToMenuKeyboard(copy))
	}

	index := len(mistakes) - 1
	mistake := mistakes[index]
	if err := b.store.setMode(user.TelegramID, modeMistakePrefix+itoa(index)); err != nil {
		return err
	}

	copy := ui(user)
	text := "*" + escapeMarkdownV2(copy.Mistakes) + "* 🛠\n\n~~" + escapeMarkdownV2(mistake.Word) + "~~\n\n"
	if mistake.Explanation != "" {
		text += "_" + escapeMarkdownV2(mistake.Explanation) + "_\n\n"
	}
	text += escapeMarkdownV2(systemUI(user).LessonAnswerInstruction)
	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, text, mistakePracticeKeyboard(copy))
}

func (b *bot) handleMistakePracticeAnswer(ctx context.Context, chatID int64, user userState, answer string) error {
	index, err := strconv.Atoi(strings.TrimPrefix(user.Mode, modeMistakePrefix))
	mistakes := mistakesForLanguage(user)
	if err != nil || index < 0 || index >= len(mistakes) {
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeText(user, "mistake_missing"), ui(user))
	}
	mistake := mistakes[index]

	if normalizeAnswer(answer) != normalizeAnswer(mistake.Correction) {
		copy := ui(user)
		text := escapeMarkdownV2(copy.TryAgain) + "\n\n~~" + escapeMarkdownV2(mistake.Word) + "~~\n\n"
		if mistake.Explanation != "" {
			text += "_" + escapeMarkdownV2(mistake.Explanation) + "_\n\n"
		}
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID, text, mistakePracticeKeyboard(copy))
	}

	if err := b.store.removeMistake(user.TelegramID, user.LearningLanguage, mistake.Word, mistake.Correction); err != nil {
		return err
	}
	if err := b.store.addXP(user.TelegramID, 8); err != nil {
		return err
	}
	if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
		return err
	}

	copy := ui(user)
	nextText := copy.NextWord
	if len(mistakes) <= 1 {
		nextText = copy.Mistakes
	}
	if err := b.telegram.sendInlineMarkdownMessage(ctx, chatID,
		escapeMarkdownV2(copy.Correct)+"\n\n~~"+escapeMarkdownV2(mistake.Word)+"~~ → `"+escapeMarkdownV2(mistake.Correction)+"`\n\n\\+8 XP",
		mistakePracticeDoneKeyboard(nextText, len(mistakes) > 1, copy)); err != nil {
		return err
	}
	b.sendTextPronunciation(ctx, chatID, user, mistake.Correction, "correction.mp3", mistake.Correction, "mistake correction")
	return b.maybePromoteLearningLevel(ctx, chatID, user.TelegramID, user.FirstName)
}

func (b *bot) sendVocabulary(ctx context.Context, chatID int64, user userState, page int) error {
	const pageSize = 10

	words := learnedWordsForLanguage(user)
	total := len(words)
	if total == 0 {
		copy := ui(user)
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			"*"+escapeMarkdownV2(copy.Vocabulary)+"* 📚\n\n"+escapeMarkdownV2(copy.LearnWords),
			vocabularyEmptyKeyboard(copy))
	}

	totalPages := (total + pageSize - 1) / pageSize
	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}

	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	pageWords := words[start:end]
	lookup := vocabularyLookupForLearnedWords(pageWords)

	var builder strings.Builder
	copy := ui(user)
	builder.WriteString("*" + escapeMarkdownV2(copy.Vocabulary) + "* 📚\n")
	builder.WriteString(escapeMarkdownV2(fmt.Sprintf(copy.TotalLearned, total)) + "\n")
	builder.WriteString("*" + itoa(page+1) + "/" + itoa(totalPages) + "*\n\n")
	for i, word := range pageWords {
		translation := learnedWordTranslationWithLookup(word, user.InterfaceLanguage, lookup)
		context := learnedWordContextWithLookup(word, user.InterfaceLanguage, lookup)
		builder.WriteString(escapeMarkdownV2(itoa(start+i+1) + ". "))
		builder.WriteString("`" + escapeMarkdownV2(translation) + "` — *" + escapeMarkdownV2(word.English) + "*\n")
		if context != "" && context != translation {
			builder.WriteString("   " + escapeMarkdownV2(context) + "\n")
		}
	}

	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, builder.String(), vocabularyKeyboard(page, totalPages, ui(user)))
}

func (b *bot) sendPhrasebook(ctx context.Context, chatID int64, user userState, page int) error {
	const pageSize = 10

	items := normalizePhrasebookEntries(user.Phrasebook, user.LearningLanguage, time.Now().UTC())
	copy := ui(user)
	if len(items) == 0 {
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			"*"+escapeMarkdownV2(copy.Phrasebook)+"* 🔖\n\n0",
			phrasebookKeyboard(0, 1, copy))
	}

	totalPages := (len(items) + pageSize - 1) / pageSize
	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}

	start := page * pageSize
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	var builder strings.Builder
	builder.WriteString("*" + escapeMarkdownV2(copy.Phrasebook) + "* 🔖\n")
	builder.WriteString("*" + itoa(page+1) + "/" + itoa(totalPages) + "*\n\n")
	for i, item := range items[start:end] {
		builder.WriteString(escapeMarkdownV2(itoa(start+i+1) + ". "))
		builder.WriteString("*" + escapeMarkdownV2(item.Phrase) + "*")
		if strings.TrimSpace(item.Translation) != "" {
			builder.WriteString("\n   `")
			builder.WriteString(escapeMarkdownV2(item.Translation))
			builder.WriteString("`")
		}
		if strings.TrimSpace(item.Note) != "" && strings.TrimSpace(item.Note) != strings.TrimSpace(item.Translation) {
			builder.WriteString("\n   ")
			builder.WriteString(escapeMarkdownV2(item.Note))
		}
		builder.WriteString("\n\n")
	}

	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, strings.TrimSpace(builder.String()), phrasebookKeyboard(page, totalPages, copy))
}

func (b *bot) handleLessonAnswer(ctx context.Context, chatID int64, user userState, text string, pronunciation *pronunciationAssessment) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	if user.LastLessonPrompt == "" {
		return b.startLesson(ctx, chatID, user)
	}

	_ = b.telegram.sendChatAction(ctx, chatID, "typing")
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	raw, err := b.openrouter.complete(ctx, feedbackPrompt(language, interfaceLanguage, user.Level, user.LastLessonPrompt, text), 0.3, 900)
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeMessage(user, "check_failed", err.Error()), ui(user))
	}

	feedback, mistakesJSON := splitFeedbackAndMistakes(raw)

	entries := parseMistakesJSON(mistakesJSON, time.Now().UTC())
	if len(entries) > 0 {
		if saveErr := b.store.addMistakes(user.TelegramID, user.LearningLanguage, entries); saveErr != nil {
			log.Printf("addMistakes (lesson) for %d: %v", user.TelegramID, saveErr)
		}
	}
	correctionAudio := lessonCorrectionPronunciationText(feedback, entries)

	if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
		return err
	}
	if err := b.store.addXP(user.TelegramID, 20); err != nil {
		return err
	}
	copy := ui(user)
	replyText := feedback
	if pronunciationText := pronunciationAssessmentText(user, pronunciation); pronunciationText != "" {
		replyText += "\n\n" + pronunciationText
	}
	if err := b.telegram.sendInlineMessage(ctx, chatID, replyText+"\n\n"+fmt.Sprintf(systemUI(user).LessonDone, copy.NewLesson), lessonDoneKeyboard(copy)); err != nil {
		return err
	}
	if correctionAudio != "" {
		b.sendTextPronunciation(ctx, chatID, user, correctionAudio, "correction.mp3", correctionAudio, "lesson correction")
	}
	return b.maybePromoteLearningLevel(ctx, chatID, user.TelegramID, user.FirstName)
}

func (b *bot) handlePracticeMessage(ctx context.Context, chatID int64, user userState, text string, pronunciation *pronunciationAssessment) error {
	b.deletePreviousWordPronunciation(ctx, chatID)
	if !canUsePractice(user) {
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.telegram.sendMessageWithCopy(ctx, chatID, b.limitReachedText(systemUI(user).PracticeKind, user), ui(user))
	}
	_ = b.telegram.sendChatAction(ctx, chatID, "typing")
	practiceHistory := trimPracticeHistory(append(user.PracticeHistory, text), practiceMemoryLimit)
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	raw, err := b.openrouter.complete(ctx, practicePrompt(language, interfaceLanguage, user.Level, practiceHistory, user.LearningFocus, user.PracticeCount), 0.6, 900)
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeMessage(user, "reply_failed", err.Error()), ui(user))
	}

	reply, mistakesJSON := splitFeedbackAndMistakes(raw)

	entries := parseMistakesJSON(mistakesJSON, time.Now().UTC())
	if len(entries) > 0 {
		if saveErr := b.store.addMistakes(user.TelegramID, user.LearningLanguage, entries); saveErr != nil {
			log.Printf("addMistakes (practice) for %d: %v", user.TelegramID, saveErr)
		}
	}
	correctionAudio := practiceCorrectionPronunciationText(reply, entries)
	questionAudio := practiceQuestionPronunciationText(reply)

	if err := b.store.incrementPractice(user.TelegramID); err != nil {
		return err
	}
	if err := b.store.savePracticeHistory(user.TelegramID, practiceHistory); err != nil {
		return err
	}
	replyText := reply
	if pronunciationText := pronunciationAssessmentText(user, pronunciation); pronunciationText != "" {
		replyText += "\n\n" + pronunciationText
	}
	if err := b.telegram.sendSmoothMessageWithCopy(ctx, chatID, replyText, ui(user)); err != nil {
		return err
	}
	if correctionAudio != "" {
		b.sendTextPronunciation(ctx, chatID, user, correctionAudio, "correction.mp3", correctionAudio, "practice correction")
	}
	if questionAudio != "" && questionAudio != correctionAudio {
		if correctionAudio != "" {
			time.Sleep(2 * time.Second)
		}
		b.sendTextPronunciation(ctx, chatID, user, questionAudio, "question.mp3", questionAudio, "practice question")
	}
	return b.maybePromoteLearningLevel(ctx, chatID, user.TelegramID, user.FirstName)
}

func (b *bot) maybePromoteLearningLevel(ctx context.Context, chatID int64, telegramID int64, firstName string) error {
	user, err := b.store.getOrCreateUser(telegramID, firstName)
	if err != nil {
		return err
	}
	next, ok := shouldPromoteLearningLevel(user)
	if !ok {
		return nil
	}
	if err := b.store.setUserLevel(telegramID, next); err != nil {
		return err
	}
	return b.telegram.sendInlineMessage(ctx, chatID,
		botRuntimeMessage(user, "level_promoted", next),
		levelPromotionKeyboard())
}

func (b *bot) continueToolContextAsPractice(ctx context.Context, chatID int64, telegramID int64, firstName string, user userState, text string) error {
	if !canUsePractice(user) {
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.telegram.sendMessageWithCopy(ctx, chatID, b.limitReachedText(systemUI(user).PracticeKind, user), ui(user))
	}
	if err := b.store.setMode(user.TelegramID, "practice"); err != nil {
		return err
	}
	refreshed, err := b.store.getOrCreateUser(telegramID, firstName)
	if err != nil {
		return err
	}
	return b.handlePracticeMessage(ctx, chatID, refreshed, text, nil)
}

func (b *bot) sendTranslatorMenu(ctx context.Context, chatID int64, user userState) error {
	copy := ui(user)
	sourceCode, targetCode := parseTranslatorMode(user.Mode, user)
	text := copy.Tool.TranslatorModePrompt + "\n\n" +
		"↤ " + translatorLanguageInterfaceName(sourceCode) + "\n" +
		"↦ " + translatorLanguageInterfaceName(targetCode)
	return b.telegram.sendInlineMessage(ctx, chatID, text, translatorInlineKeyboard(user, b.cfg.WebAppURL))
}

func (b *bot) handleTranslatorLanguageCallback(ctx context.Context, chatID int64, user userState, data string, target bool) error {
	parts := strings.Split(data, "|")
	if len(parts) != 2 {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
	sourceCode, targetCode := parseTranslatorMode(user.Mode, user)
	if target {
		targetCode = normalizeTranslatorTarget(parts[1])
		if targetCode == "" {
			targetCode = defaultTranslatorTarget(user)
		}
	} else {
		sourceCode = normalizeTranslatorSource(parts[1])
	}
	if err := b.store.setMode(user.TelegramID, translatorMode(sourceCode, targetCode)); err != nil {
		return err
	}
	refreshed, err := b.store.getOrCreateUser(user.TelegramID, user.FirstName)
	if err != nil {
		return err
	}
	return b.sendTranslatorMenu(ctx, chatID, refreshed)
}

func (b *bot) swapTranslatorLanguages(ctx context.Context, chatID int64, user userState) error {
	sourceCode, targetCode := parseTranslatorMode(user.Mode, user)
	if sourceCode == translationAutoCode {
		return b.sendTranslatorMenu(ctx, chatID, user)
	}
	if err := b.store.setMode(user.TelegramID, translatorMode(targetCode, sourceCode)); err != nil {
		return err
	}
	refreshed, err := b.store.getOrCreateUser(user.TelegramID, user.FirstName)
	if err != nil {
		return err
	}
	return b.sendTranslatorMenu(ctx, chatID, refreshed)
}

func (b *bot) handleTranslatorText(ctx context.Context, chatID int64, user userState, text string) error {
	text = strings.TrimSpace(text)
	if text == "" {
		return b.sendTranslatorMenu(ctx, chatID, user)
	}
	return b.handleTranslatorInput(ctx, chatID, user, text, "text")
}

func (b *bot) handleTranslatorInput(ctx context.Context, chatID int64, user userState, text string, kind string) error {
	if b.openrouter == nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeText(user, "translator_unavailable"), ui(user))
	}
	sourceCode, targetCode := parseTranslatorMode(user.Mode, user)
	_ = b.telegram.sendChatAction(ctx, chatID, "typing")
	translation, err := b.openrouter.completeWithModel(ctx, b.cfg.OpenRouterTranslatorModel, translationToolPrompt(text, sourceCode, targetCode, userInterfaceLanguage(user)), 0.1, 900)
	if err != nil {
		return b.telegram.sendMessage(ctx, chatID, fmt.Sprintf(ui(user).Tool.VoiceTranslationFailed, err.Error()))
	}
	return b.sendTranslatorResult(ctx, chatID, user, text, translation)
}

func (b *bot) sendTranslatorResult(ctx context.Context, chatID int64, user userState, sourceText string, translation string) error {
	copy := ui(user)
	result := translatorToolResultText(copy, sourceText, translation)
	if strings.TrimSpace(result) == "" {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeText(user, "translator_empty"), copy)
	}
	history := trimPracticeHistory(append(user.PracticeHistory, toolContextEntry("translator", result)), practiceMemoryLimit)
	if err := b.store.savePracticeHistory(user.TelegramID, history); err != nil {
		return err
	}
	if err := b.telegram.sendMessage(ctx, chatID, result); err != nil {
		return err
	}
	if b.openrouter == nil || strings.TrimSpace(b.cfg.OpenRouterTTSModel) == "" || strings.TrimSpace(translation) == "" {
		return nil
	}
	audio, err := b.openrouter.synthesizeSpeech(ctx, b.cfg.OpenRouterTTSModel, b.cfg.OpenRouterTTSVoice, translation)
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeMessage(user, "translation_tts_failed", err.Error()), copy)
	}
	_, err = b.telegram.sendAudioBytesWithCopy(ctx, chatID, audio, "translation.mp3", copy.Tool.Translator, copy)
	return err
}

func toolContextEntry(kind string, result string) string {
	result = strings.TrimSpace(result)
	if result == "" {
		return ""
	}
	return "Tool context (" + kind + "):\n" + result
}

func voiceToolResult(copy uiCopy, transcript, translation string) string {
	var builder strings.Builder
	builder.WriteString(copy.Tool.TranscriptLabel + ":\n")
	builder.WriteString(strings.TrimSpace(transcript))
	translation = strings.TrimSpace(translation)
	if translation != "" {
		builder.WriteString("\n\n" + copy.Tool.TranslationLabel + ":\n")
		builder.WriteString(translation)
	}
	return builder.String()
}

func (b *bot) handleVoiceMessage(ctx context.Context, message *telegramMessage) error {
	user, err := b.store.getOrCreateUser(message.From.ID, telegramDisplayName(message.From))
	if err != nil {
		return err
	}
	copy := ui(user)
	if !user.isPremium(time.Now()) {
		return b.telegram.sendMessage(ctx, message.Chat.ID, copy.Tool.VoicePremiumRequired)
	}
	voiceLimit := voiceLimitFor(user)
	if user.VoiceToday >= voiceLimit {
		return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
	}
	if message.Voice.Duration > b.cfg.MaxVoiceSeconds {
		return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.VoiceTooLong, b.cfg.MaxVoiceSeconds, message.Voice.Duration))
	}
	if user.Mode != modeToolVoiceText && !strings.HasPrefix(user.Mode, modeToolTranslatorPrefix) && !canUsePractice(user) {
		return b.telegram.sendMessageWithCopy(ctx, message.Chat.ID, b.limitReachedText(systemUI(user).PracticeKind, user), ui(user))
	}

	_ = b.telegram.sendChatAction(ctx, message.Chat.ID, "typing")
	audioBytes, err := b.telegram.downloadFile(ctx, message.Voice.FileID)
	if err != nil {
		return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.VoiceDownloadFailed, err.Error()))
	}

	if strings.HasPrefix(user.Mode, modeToolTranslatorPrefix) {
		transcript, err := b.openrouter.transcribeAudio(ctx, b.cfg.OpenRouterSTTModel, audioBytes, "ogg")
		if err != nil {
			return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.VoiceTranscribeFailed, err.Error()))
		}
		if err := b.store.incrementVoice(user.TelegramID); err != nil {
			return err
		}
		return b.handleTranslatorInput(ctx, message.Chat.ID, user, transcript, "voice")
	}
	if user.Mode == modeToolVoiceText {
		transcript, err := b.openrouter.transcribeAudio(ctx, b.cfg.OpenRouterSTTModel, audioBytes, "ogg")
		if err != nil {
			return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.VoiceTranscribeFailed, err.Error()))
		}
		if err := b.store.incrementVoice(user.TelegramID); err != nil {
			return err
		}
		translation, err := b.openrouter.complete(ctx, translationPrompt(transcript, userInterfaceLanguage(user)), 0.1, 500)
		if err != nil {
			return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.VoiceTranslationFailed, err.Error()))
		}
		result := voiceToolResult(copy, transcript, translation)
		history := trimPracticeHistory(append(user.PracticeHistory, toolContextEntry("voice", result)), practiceMemoryLimit)
		if err := b.store.savePracticeHistory(user.TelegramID, history); err != nil {
			return err
		}
		return b.telegram.sendMessage(ctx, message.Chat.ID, result+"\n\n"+copy.Tool.VoiceDiscussPrompt)
	}
	if strings.HasPrefix(user.Mode, modeShadowingPrefix) {
		target, _ := parseShadowingMode(user.Mode)
		transcription, err := b.transcribeLearningVoice(ctx, user, audioBytes, "ogg", target)
		if err != nil {
			return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.VoiceTranscribeFailed, err.Error()))
		}
		if message.Voice.Duration > 0 && transcription.DurationSeconds <= 0 {
			transcription.DurationSeconds = float64(message.Voice.Duration)
		}
		if err := b.store.incrementVoice(user.TelegramID); err != nil {
			return err
		}
		return b.handleShadowingAnswer(ctx, message.Chat.ID, user, transcription.Text, true, &transcription)
	}

	transcription, err := b.transcribeLearningVoice(ctx, user, audioBytes, "ogg", "")
	if err != nil {
		return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.VoiceTranscribeFailed, err.Error()))
	}
	if message.Voice.Duration > 0 && transcription.DurationSeconds <= 0 {
		transcription.DurationSeconds = float64(message.Voice.Duration)
	}
	if err := b.store.incrementVoice(user.TelegramID); err != nil {
		return err
	}
	transcript := strings.TrimSpace(transcription.Text)
	pronunciation := b.buildPronunciationAssessment(ctx, user, "", transcription, pronunciationModeFree)

	if err := b.telegram.sendMarkdownMessageWithCopy(ctx, message.Chat.ID, "*"+escapeMarkdownV2(botRuntimeText(user, "heard_label"))+"*\n_"+escapeMarkdownV2(transcript)+"_", copy); err != nil {
		return err
	}

	refreshed, err := b.store.getOrCreateUser(message.From.ID, telegramDisplayName(message.From))
	if err != nil {
		return err
	}
	if refreshed.Mode == "lesson" {
		return b.handleLessonAnswer(ctx, message.Chat.ID, refreshed, transcript, &pronunciation)
	}
	return b.handlePracticeMessage(ctx, message.Chat.ID, refreshed, transcript, &pronunciation)
}

func (b *bot) handlePhotoMessage(ctx context.Context, message *telegramMessage) error {
	user, err := b.store.getOrCreateUser(message.From.ID, telegramDisplayName(message.From))
	if err != nil {
		return err
	}
	copy := ui(user)
	translatorModeActive := strings.HasPrefix(user.Mode, modeToolTranslatorPrefix)
	practiceModeActive := user.Mode == "practice"
	if user.Mode != modeToolImageTranslate && !translatorModeActive && !practiceModeActive {
		return b.telegram.sendInlineMessage(ctx, message.Chat.ID, copy.Tool.ImageOpenToolsPrompt, toolsInlineKeyboard(b.cfg.WebAppURL, copy))
	}
	if message.MediaGroupID != "" {
		if err := b.store.setMode(user.TelegramID, "idle"); err != nil {
			return err
		}
		return b.telegram.sendMessage(ctx, message.Chat.ID, copy.Tool.ImageSinglePhotoPrompt)
	}
	if !user.isPremium(time.Now()) {
		return b.telegram.sendMessage(ctx, message.Chat.ID, copy.Tool.ImagePremiumRequired)
	}

	photo := largestTelegramPhoto(message.Photo)
	if photo.FileID == "" {
		return b.telegram.sendMessage(ctx, message.Chat.ID, copy.Tool.ImageFileMissing)
	}

	_ = b.telegram.sendChatAction(ctx, message.Chat.ID, "typing")
	imageBytes, err := b.telegram.downloadFile(ctx, photo.FileID)
	if err != nil {
		return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.ImageDownloadFailed, err.Error()))
	}

	mimeType := http.DetectContentType(imageBytes)
	if translatorModeActive {
		sourceCode, targetCode := parseTranslatorMode(user.Mode, user)
		result, err := b.openrouter.translateImageForToolWithModel(ctx, b.cfg.OpenRouterTranslatorModel, imageBytes, mimeType, sourceCode, targetCode, userInterfaceLanguage(user))
		if err != nil {
			return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.ImageReadFailed, err.Error()))
		}
		return b.sendTranslatorResult(ctx, message.Chat.ID, user, result.SourceText, result.Translation)
	}
	if practiceModeActive {
		imageContext, err := b.openrouter.describePracticeImage(ctx, imageBytes, mimeType, userLearningLanguage(user), userInterfaceLanguage(user), message.Caption)
		if err != nil {
			return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.ImageReadFailed, err.Error()))
		}
		return b.handlePracticeMessage(ctx, message.Chat.ID, user, practiceImageMessage(message.Caption, imageContext), nil)
	}
	result, err := b.openrouter.translateImageText(ctx, imageBytes, mimeType, userLearningLanguage(user), userInterfaceLanguage(user))
	if err != nil {
		return b.telegram.sendMessage(ctx, message.Chat.ID, fmt.Sprintf(copy.Tool.ImageReadFailed, err.Error()))
	}
	history := trimPracticeHistory(append(user.PracticeHistory, toolContextEntry("image", result)), practiceMemoryLimit)
	if err := b.store.savePracticeHistory(user.TelegramID, history); err != nil {
		return err
	}
	return b.telegram.sendMessage(ctx, message.Chat.ID, result+"\n\n"+copy.Tool.ImageDiscussPrompt)
}

func largestTelegramPhoto(photos []telegramPhoto) telegramPhoto {
	if len(photos) == 0 {
		return telegramPhoto{}
	}
	largest := photos[0]
	for _, photo := range photos[1:] {
		if photo.Width*photo.Height > largest.Width*largest.Height {
			largest = photo
		}
	}
	return largest
}

func (b *bot) handlePreCheckoutQuery(ctx context.Context, query preCheckoutQuery) error {
	plan, payloadUserID, _, ok := b.premiumPlanPaymentFromPayload(query.InvoicePayload)
	ok = query.Currency == "XTR" && ok && payloadUserID == query.From.ID && query.TotalAmount == plan.StarsPrice
	if !ok {
		user, err := b.store.getOrCreateUser(query.From.ID, telegramDisplayName(query.From))
		if err != nil {
			return err
		}
		return b.telegram.answerPreCheckoutQuery(ctx, query.ID, false, premiumUI(user).PreCheckoutFailed)
	}
	return b.telegram.answerPreCheckoutQuery(ctx, query.ID, true, "")
}

func (b *bot) handleSuccessfulPayment(ctx context.Context, chatID int64, user userState, payment successfulPayment) error {
	plan, payloadUserID, _, ok := b.premiumPlanPaymentFromPayload(payment.InvoicePayload)
	if payment.Currency != "XTR" || !ok || payloadUserID != user.TelegramID || payment.TotalAmount != plan.StarsPrice || strings.TrimSpace(payment.TelegramPaymentChargeID) == "" {
		return b.telegram.sendMessageWithCopy(ctx, chatID, fmt.Sprintf(premiumUI(user).PaymentUnrecognized, ui(user).MenuButton), ui(user))
	}
	until, err := b.store.extendPremium(user.TelegramID, payment.TelegramPaymentChargeID, plan.Duration, plan.Tier)
	if err != nil {
		return err
	}
	if _, err := b.store.creditReferralPurchase(user.TelegramID, payment.TelegramPaymentChargeID, rubToKopecks(plan.RubPrice)); err != nil {
		return err
	}
	plan = localizedPremiumPlan(user, plan)
	b.notifyOpsPayment(ctx, "Telegram Stars", user, plan, payment.TelegramPaymentChargeID, fmt.Sprintf("%d XTR", payment.TotalAmount), until)
	return b.telegram.sendMarkdownMessageWithCopy(ctx, chatID, premiumActivationMarkdown(user, plan, until.Format("2006-01-02"), true), ui(user))
}

func (b *bot) notifyOpsPayment(ctx context.Context, source string, user userState, plan premiumPlan, paymentID string, amount string, until time.Time) {
	if b == nil || b.telegram == nil {
		return
	}
	text := strings.Join([]string{
		"Payment Poliglot AI",
		"Source: " + strings.TrimSpace(source),
		fmt.Sprintf("User: %s (%d)", strings.TrimSpace(user.FirstName), user.TelegramID),
		"Plan: " + strings.TrimSpace(plan.Title),
		"Tier: " + strings.TrimSpace(plan.Tier),
		"Payment: " + strings.TrimSpace(paymentID),
		"Amount: " + strings.TrimSpace(amount),
		"Premium until: " + until.Format("2006-01-02"),
	}, "\n")
	b.notifyOpsText(ctx, text)
}

func (b *bot) notifyOpsText(ctx context.Context, text string) {
	if b == nil || b.telegram == nil {
		return
	}
	for _, recipient := range b.cfg.telegramOpsRecipients() {
		if err := b.telegram.sendMessageToChat(ctx, recipient.telegramChatIDValue(), text); err != nil {
			log.Printf("send ops notification to %s: %v", recipient.ChatID, err)
		}
	}
}

func (b *bot) notifyOpsPhoto(ctx context.Context, photoPath string, caption string) {
	if b == nil || b.telegram == nil {
		return
	}
	for _, recipient := range b.cfg.telegramOpsRecipients() {
		if err := b.telegram.sendPhotoFileToChat(ctx, recipient.telegramChatIDValue(), photoPath, caption); err != nil {
			log.Printf("send ops photo to %s: %v", recipient.ChatID, err)
		}
	}
}

func (recipient telegramOpsRecipient) telegramChatIDValue() any {
	chatID := strings.TrimSpace(recipient.ChatID)
	if parsed, err := strconv.ParseInt(chatID, 10, 64); err == nil {
		return parsed
	}
	return chatID
}

func (b *bot) maybeBeginWebAuth(ctx context.Context, chatID int64, user userState, from telegramUser, text string) (bool, error) {
	parts := strings.Fields(text)
	if len(parts) < 2 || strings.ToLower(parts[0]) != "/start" || !strings.HasPrefix(parts[1], webAuthStartPrefix) {
		return false, nil
	}
	if b.webAuth == nil {
		return true, b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeText(user, "web_auth_unavailable"), ui(user))
	}
	token := strings.TrimPrefix(parts[1], webAuthStartPrefix)
	profile := telegramWebProfile{
		ID:        user.TelegramID,
		FirstName: from.FirstName,
		Username:  from.Username,
	}
	request, err := b.webAuth.beginTelegram(token, profile)
	if err != nil {
		return true, b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeMessage(user, "web_auth_failed", err.Error()), ui(user))
	}
	return true, b.telegram.sendMessageWithCopy(ctx, chatID, botRuntimeMessage(user, "web_auth_code", request.Code), ui(user))
}

func (b *bot) maybeVerifyWebAuthCode(ctx context.Context, chatID int64, user userState, text string) (bool, error) {
	return false, nil
}

func (b *bot) maybeApplyStartReferral(ctx context.Context, chatID int64, user userState, text string) (bool, error) {
	parts := strings.Fields(text)
	if len(parts) < 2 || strings.ToLower(parts[0]) != "/start" || !strings.HasPrefix(parts[1], "ref_") {
		return false, nil
	}

	inviterID, err := parseReferralPayload(parts[1])
	if err != nil {
		if err := b.telegram.sendMessage(ctx, chatID, onboardingWelcomeText(user)); err != nil {
			return true, err
		}
		return true, b.continueOnboarding(ctx, chatID, user)
	}

	referral, err := b.store.applyReferral(user.TelegramID, inviterID)
	if err != nil {
		return true, err
	}
	if referral.Applied {
		copy := premiumUI(user)
		_ = b.telegram.sendMessage(ctx, inviterID, fmt.Sprintf(copy.ReferralInviter, referral.InviterPremiumDays, referral.InviterUntil.Format("2006-01-02")))
		if err := b.telegram.sendMessage(ctx, chatID, onboardingWelcomeText(user)+"\n\n"+copy.ReferralInvitee); err != nil {
			return true, err
		}
		return true, b.continueOnboarding(ctx, chatID, user)
	}
	if err := b.telegram.sendMessage(ctx, chatID, onboardingWelcomeText(user)); err != nil {
		return true, err
	}
	return true, b.continueOnboarding(ctx, chatID, user)
}

func (b *bot) sendInviteLink(ctx context.Context, chatID int64, user userState) error {
	me, err := b.telegram.getMe(ctx)
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, fmt.Sprintf(premiumUI(user).InviteCreateFailed, err.Error()), ui(user))
	}
	link := fmt.Sprintf("https://t.me/%s?start=ref_%d", me.Username, user.TelegramID)
	copy := premiumUI(user)
	template := strings.TrimSpace(copy.ReferralShareText)
	if template == "" {
		template = copy.InviteLinkText
	}
	return b.telegram.sendMessageWithCopy(ctx, chatID, fmt.Sprintf(template, link), ui(user))
}

func referralAccountText(user userState) string {
	copy := premiumUI(user)
	return fmt.Sprintf("%s: %s (%s)\n%s: %d\n%s",
		copy.ReferralBalanceTitle,
		formatRubKopecks(user.ReferralBalanceKopecks),
		formatUSDTFromKopecks(user.ReferralBalanceKopecks),
		copy.ReferralInvitationsLabel,
		user.ReferralCount,
		copy.ReferralBalanceHint)
}

func (b *bot) sendReferralWithdrawInfo(ctx context.Context, chatID int64, user userState) error {
	copy := premiumUI(user)
	text := fmt.Sprintf(copy.ReferralWithdrawUnavailable,
		formatRubKopecks(user.ReferralBalanceKopecks),
		formatUSDTFromKopecks(user.ReferralBalanceKopecks))
	return b.telegram.sendInlineMessage(ctx, chatID, text, backToMenuKeyboard(ui(user)))
}

func (b *bot) sendPremiumInvoice(ctx context.Context, chatID int64, user userState, product string) error {
	plan, ok := b.premiumPlan(product)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
	plan = localizedPremiumPlan(user, plan)
	payload := premiumPaymentPayload(plan.Product, user.TelegramID, time.Now())
	return b.telegram.sendInvoice(
		ctx,
		chatID,
		plan.Title,
		b.premiumInvoiceDescription(user, plan),
		payload,
		plan.StarsPrice,
		ui(user),
		premiumUI(user).PayStarsButton,
	)
}

func (b *bot) premiumInvoiceDescription(user userState, plan premiumPlan) string {
	lessonLimit, practiceLimit, voiceLimit := paidPlanLimits(plan.Tier)
	return fmt.Sprintf(premiumUI(user).InvoiceDescription,
		lessonLimit,
		practiceLimit,
		voiceLimit,
		b.cfg.MaxVoiceSeconds,
	)
}

func (b *bot) sendYooKassaPayment(ctx context.Context, chatID int64, user userState, product string) error {
	if b.yookassa == nil || !b.cfg.yooKassaEnabled() {
		return b.telegram.sendMessageWithCopy(ctx, chatID, premiumUI(user).YooKassaUnavailable, ui(user))
	}
	plan, ok := b.premiumPlan(product)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
	plan = localizedPremiumPlan(user, plan)

	payment, err := b.yookassa.createPremiumPayment(ctx, user.TelegramID, user.FirstName, user.InterfaceLanguage, plan, b.cfg.YooKassaReturnURL, "telegram")
	if err != nil {
		log.Printf("failed to create yookassa payment for %d: %v", user.TelegramID, err)
		return b.telegram.sendMessageWithCopy(ctx, chatID, premiumUI(user).YooKassaCreateFailed, ui(user))
	}

	text := plainTextFromMarkdownV2(fmt.Sprintf(
		premiumUI(user).YooKassaPaymentText,
		plan.Title,
		plan.RubPrice,
	))
	return b.telegram.sendInlineMessage(ctx, chatID, text, yookassaPaymentKeyboard(payment.Confirmation.ConfirmationURL, user))
}

func (b *bot) sendCryptoPayment(ctx context.Context, chatID int64, user userState, product string, methodID string) error {
	if !b.cryptoPaymentsReady() {
		return b.telegram.sendMessageWithCopy(ctx, chatID, "Direct crypto payments are not configured yet.", ui(user))
	}
	plan, ok := b.premiumPlan(product)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
	payment, invoiceURL, err := b.createDirectCryptoPayment(ctx, product, user, methodID)
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, err.Error(), ui(user))
	}
	if err := b.saveDirectCryptoPayment(payment); err != nil {
		return err
	}
	plan = localizedPremiumPlan(user, plan)
	return b.telegram.sendInlineMessage(ctx, chatID, cryptoPaymentTextV2(user, plan, payment, invoiceURL), cryptoPaymentKeyboard(invoiceURL, payment.ID, user))
}

func (b *bot) checkCryptoPayment(ctx context.Context, chatID int64, user userState, paymentID string) error {
	payment, refreshed, paid, err := b.refreshDirectCryptoPayment(ctx, paymentID, user)
	if err != nil {
		return b.telegram.sendMessageWithCopy(ctx, chatID, err.Error(), ui(user))
	}
	if !paid {
		plan, _ := b.premiumPlan(payment.Product)
		plan = localizedPremiumPlan(user, plan)
		return b.telegram.sendInlineMessage(ctx, chatID, cryptoPaymentPendingTextV2(user, plan, payment), cryptoPaymentKeyboard(b.cryptoPaymentURL(payment), payment.ID, user))
	}
	plan, _ := b.premiumPlan(payment.Product)
	plan = localizedPremiumPlan(refreshed, plan)
	return b.telegram.sendMarkdownMessageWithCopy(ctx, chatID, premiumActivationMarkdown(refreshed, plan, refreshed.PremiumUntil.Format("2006-01-02"), true), ui(refreshed))
}

func cryptoPaymentText(user userState, plan premiumPlan, payment cryptoPayment, invoiceURL string) string {
	expiresAt := userLocalDateTimeLabel(payment.ExpiresAt, user)
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		return fmt.Sprintf(
			"%s\n\nСумма: %s TON\nСеть: TON\nАдрес: %s\nКомментарий: %s\n\nОткройте TON-кошелек кнопкой ниже или отправьте перевод вручную. Комментарий обязателен: по нему сервер найдет ваш платеж.\n\nПосле перевода нажмите \"I paid - check\".\n\nСсылка: %s\nСчет действителен до: %s",
			plan.Title,
			payment.Amount,
			payment.Address,
			payment.Memo,
			invoiceURL,
			expiresAt,
		)
	}
	return fmt.Sprintf(
		"%s\n\nAmount: %s TON\nNetwork: TON\nAddress: %s\nComment: %s\n\nOpen a TON wallet with the button below or send the transfer manually. The comment is required: the server uses it to match your payment.\n\nAfter sending, tap \"I paid - check\".\n\nLink: %s\nInvoice expires at: %s",
		plan.Title,
		payment.Amount,
		payment.Address,
		payment.Memo,
		invoiceURL,
		expiresAt,
	)
}

func cryptoPaymentPendingText(user userState, plan premiumPlan, payment cryptoPayment) string {
	if payment.Status == cryptoStatusExpired {
		if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
			return "Счет истек. Откройте Premium и создайте новый TON-счет."
		}
		return "This invoice has expired. Open Premium and create a new TON invoice."
	}
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		return fmt.Sprintf("%s\n\nПлатеж пока не найден.\n\nПроверьте сумму %s TON и обязательный комментарий: %s", plan.Title, payment.Amount, payment.Memo)
	}
	return fmt.Sprintf("%s\n\nPayment is not found yet.\n\nCheck the %s TON amount and required comment: %s", plan.Title, payment.Amount, payment.Memo)
}

func cryptoPaymentTextV2(user userState, plan premiumPlan, payment cryptoPayment, invoiceURL string) string {
	expiresAt := userLocalDateTimeLabel(payment.ExpiresAt, user)
	methodLabel := cryptoPaymentMethodLabel(payment)
	amountLabel := payment.Amount + " " + payment.Currency
	memo := strings.TrimSpace(payment.Memo)
	ruLinkLine := ""
	enLinkLine := ""
	if strings.TrimSpace(invoiceURL) != "" {
		ruLinkLine = "\n\nСсылка: " + invoiceURL
		enLinkLine = "\n\nLink: " + invoiceURL
	}
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		memoLine := "Комментарий: не нужен"
		instruction := "Отправьте точную сумму на адрес ниже. Для этого метода сервер отличает платеж по сумме, поэтому не округляйте ее."
		if memo != "" {
			memoLine = "Комментарий: " + memo
			instruction = "Откройте кошелек кнопкой ниже или отправьте перевод вручную. Комментарий обязателен: по нему сервер найдет ваш платеж."
		}
		return fmt.Sprintf(
			"%s\n\nСумма: %s\nСеть: %s\nАдрес: %s\n%s\n\n%s\n\nПосле перевода нажмите \"I paid - check\".%s\nСчет действителен до: %s",
			plan.Title,
			amountLabel,
			methodLabel,
			payment.Address,
			memoLine,
			instruction,
			ruLinkLine,
			expiresAt,
		)
	}
	memoLine := "Comment: not needed"
	instruction := "Send the exact amount to the address below. This method is matched by the unique amount, so do not round it."
	if memo != "" {
		memoLine = "Comment: " + memo
		instruction = "Open a wallet with the button below or send the transfer manually. The comment is required: the server uses it to match your payment."
	}
	return fmt.Sprintf(
		"%s\n\nAmount: %s\nNetwork: %s\nAddress: %s\n%s\n\n%s\n\nAfter sending, tap \"I paid - check\".%s\nInvoice expires at: %s",
		plan.Title,
		amountLabel,
		methodLabel,
		payment.Address,
		memoLine,
		instruction,
		enLinkLine,
		expiresAt,
	)
}

func cryptoPaymentPendingTextV2(user userState, plan premiumPlan, payment cryptoPayment) string {
	if payment.Status == cryptoStatusExpired {
		if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
			return "Счет истек. Откройте Premium и создайте новый crypto-счет."
		}
		return "This invoice has expired. Open Premium and create a new crypto invoice."
	}
	amountLabel := payment.Amount + " " + payment.Currency
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		if strings.TrimSpace(payment.Memo) != "" {
			return fmt.Sprintf("%s\n\nПлатеж пока не найден.\n\nПроверьте сумму %s и обязательный комментарий: %s", plan.Title, amountLabel, payment.Memo)
		}
		return fmt.Sprintf("%s\n\nПлатеж пока не найден.\n\nПроверьте точную сумму %s и сеть %s.", plan.Title, amountLabel, cryptoPaymentMethodLabel(payment))
	}
	if strings.TrimSpace(payment.Memo) != "" {
		return fmt.Sprintf("%s\n\nPayment is not found yet.\n\nCheck the %s amount and required comment: %s", plan.Title, amountLabel, payment.Memo)
	}
	return fmt.Sprintf("%s\n\nPayment is not found yet.\n\nCheck the exact %s amount and the %s network.", plan.Title, amountLabel, cryptoPaymentMethodLabel(payment))
}

func (b *bot) sendPremiumMenu(ctx context.Context, chatID int64, user userState) error {
	cryptoProducts := map[string][]cryptoPaymentMethod{}
	if b.cryptoPaymentsReady() {
		for _, plan := range b.premiumPlanList() {
			cryptoProducts[plan.Product] = b.cfg.cryptoPaymentMethodsForProduct(plan.Product)
		}
	}
	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, b.premiumText(user), premiumInlineKeyboard(b.cfg.yooKassaEnabled(), cryptoProducts, localizedPremiumPlans(user, b.premiumPlanList()), user))
}

func (b *bot) sendPremiumPaymentOptions(ctx context.Context, chatID int64, user userState, product string) error {
	plan, ok := b.premiumPlan(product)
	if !ok {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).UnknownButton, ui(user))
	}
	plan = localizedPremiumPlan(user, plan)
	methods := []cryptoPaymentMethod{}
	if b.cryptoPaymentsReady() {
		methods = b.cfg.cryptoPaymentMethodsForProduct(product)
	}
	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, premiumPaymentOptionsText(user, plan), premiumPaymentOptionsKeyboard(b.cfg.yooKassaEnabled(), methods, plan, user))
}

func premiumPaymentOptionsText(user userState, plan premiumPlan) string {
	prompt := "Choose a payment method for this plan."
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		prompt = "\u0412\u044b\u0431\u0435\u0440\u0438\u0442\u0435 \u0441\u043f\u043e\u0441\u043e\u0431 \u043e\u043f\u043b\u0430\u0442\u044b \u0434\u043b\u044f \u044d\u0442\u043e\u0433\u043e \u0442\u0430\u0440\u0438\u0444\u0430."
	}
	return "*" + escapeMarkdownV2(plan.Title) + "*\n" +
		escapeMarkdownV2(plan.DaysLabel) + " \\- " + escapeMarkdownV2(strconv.Itoa(plan.RubPrice)+" RUB / "+strconv.Itoa(plan.StarsPrice)+" Stars") + "\n\n" +
		escapeMarkdownV2(prompt)
}

func (b *bot) premiumPlanList() []premiumPlan {
	return []premiumPlan{
		{
			Product:    premiumMonthlyProduct,
			Tier:       "premium",
			Title:      "Premium 30 дней",
			DaysLabel:  "30 дней",
			Duration:   premiumDuration,
			RubPrice:   b.cfg.PremiumRubPrice,
			StarsPrice: b.cfg.PremiumStarsPrice,
		},
		{
			Product:    premiumYearlyProduct,
			Tier:       "premium",
			Title:      "Premium на год",
			DaysLabel:  "365 дней",
			Duration:   premiumYearDuration,
			RubPrice:   b.cfg.PremiumYearRubPrice,
			StarsPrice: b.cfg.PremiumYearStarsPrice,
		},
		{
			Product:    platinumMonthlyProduct,
			Tier:       "platinum",
			Title:      "Platinum 30 дней",
			DaysLabel:  "30 дней",
			Duration:   premiumDuration,
			RubPrice:   b.cfg.PlatinumRubPrice,
			StarsPrice: b.cfg.PlatinumStarsPrice,
		},
		{
			Product:    platinumYearlyProduct,
			Tier:       "platinum",
			Title:      "Platinum на год",
			DaysLabel:  "365 дней",
			Duration:   premiumYearDuration,
			RubPrice:   b.cfg.PlatinumYearRubPrice,
			StarsPrice: b.cfg.PlatinumYearStarsPrice,
		},
	}
}

func (b *bot) premiumPlan(product string) (premiumPlan, bool) {
	for _, plan := range b.premiumPlanList() {
		if plan.Product == product {
			return plan, true
		}
	}
	return premiumPlan{}, false
}

func (b *bot) premiumPlanFromPayload(payload string) (premiumPlan, bool) {
	plan, _, _, ok := b.premiumPlanPaymentFromPayload(payload)
	return plan, ok
}

func premiumPaymentPayload(product string, telegramID int64, now time.Time) string {
	return fmt.Sprintf("%s_%d_%d", product, telegramID, now.Unix())
}

func premiumStarsStartPayload(product string) string {
	return "buy_" + strings.TrimSpace(product)
}

func parsePremiumStarsStartPayload(payload string) (string, bool) {
	product, ok := strings.CutPrefix(strings.TrimSpace(payload), "buy_")
	if !ok || strings.TrimSpace(product) == "" {
		return "", false
	}
	for _, id := range paidProductIDs() {
		if product == id {
			return product, true
		}
	}
	return "", false
}

func telegramBotStartURL(botName string, payload string) string {
	botName = strings.TrimPrefix(strings.TrimSpace(botName), "@")
	if botName == "" {
		botName = "Poliglot_AI_bot"
	}
	return "https://t.me/" + botName + "?start=" + strings.TrimSpace(payload)
}

func (b *bot) premiumPlanPaymentFromPayload(payload string) (premiumPlan, int64, int64, bool) {
	for _, plan := range b.premiumPlanList() {
		rest, ok := strings.CutPrefix(payload, plan.Product+"_")
		if !ok {
			continue
		}
		userIDText, timestampText, ok := strings.Cut(rest, "_")
		if !ok {
			return premiumPlan{}, 0, 0, false
		}
		userID, err := strconv.ParseInt(userIDText, 10, 64)
		if err != nil || userID == 0 {
			return premiumPlan{}, 0, 0, false
		}
		timestamp, err := strconv.ParseInt(timestampText, 10, 64)
		if err != nil || timestamp <= 0 {
			return premiumPlan{}, 0, 0, false
		}
		return plan, userID, timestamp, true
	}
	return premiumPlan{}, 0, 0, false
}

func (b *bot) premiumPlanFromMetadata(metadata map[string]interface{}) (premiumPlan, bool) {
	product, ok := metadataString(metadata, "product")
	if !ok {
		return premiumPlan{}, false
	}
	return b.premiumPlan(product)
}

func (b *bot) handleYooKassaWebhook(w http.ResponseWriter, r *http.Request) {
	b.handleYooKassaWebhookWithClient(w, r, b.yookassa, "yookassa")
}

func (b *bot) handleWebYooKassaWebhook(w http.ResponseWriter, r *http.Request) {
	b.handleYooKassaWebhookWithClient(w, r, b.webYooKassa, "web yookassa")
}

func (b *bot) handleYooKassaWebhookWithClient(w http.ResponseWriter, r *http.Request, client *yooKassaClient, serviceName string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if client == nil {
		http.Error(w, serviceName+" is not configured", http.StatusServiceUnavailable)
		return
	}

	var notification yooKassaNotification
	if err := json.NewDecoder(r.Body).Decode(&notification); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if notification.Event != "payment.succeeded" || notification.Object.ID == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	payment, err := client.getPayment(r.Context(), notification.Object.ID)
	if err != nil {
		log.Printf("failed to verify yookassa payment %s: %v", notification.Object.ID, err)
		http.Error(w, "payment verification failed", http.StatusBadGateway)
		return
	}
	plan, ok := b.premiumPlanFromMetadata(payment.Metadata)
	if !ok {
		log.Printf("yookassa payment %s has unknown product metadata", payment.ID)
		w.WriteHeader(http.StatusOK)
		return
	}
	if payment.Status != "succeeded" || !payment.Paid || payment.Amount.Currency != "RUB" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var paidAmount float64
	if _, err := fmt.Sscanf(payment.Amount.Value, "%f", &paidAmount); err != nil || paidAmount < float64(plan.RubPrice) {
		log.Printf("yookassa payment %s has unexpected amount %q", payment.ID, payment.Amount.Value)
		w.WriteHeader(http.StatusOK)
		return
	}
	telegramID, ok := metadataInt64(payment.Metadata, "telegram_id")
	if !ok || telegramID == 0 {
		log.Printf("yookassa payment %s has no telegram_id metadata", payment.ID)
		w.WriteHeader(http.StatusOK)
		return
	}

	until, err := b.store.extendPremium(telegramID, payment.ID, plan.Duration, plan.Tier)
	if err != nil {
		log.Printf("failed to activate premium for yookassa payment %s: %v", payment.ID, err)
		http.Error(w, "activation failed", http.StatusInternalServerError)
		return
	}
	paidKopecks := int64(math.Round(paidAmount * 100))
	if _, err := b.store.creditReferralPurchase(telegramID, payment.ID, paidKopecks); err != nil {
		log.Printf("failed to credit referral rewards for yookassa payment %s: %v", payment.ID, err)
		http.Error(w, "referral reward failed", http.StatusInternalServerError)
		return
	}
	interfaceLanguage, _ := metadataString(payment.Metadata, "interface_language")
	user := userState{TelegramID: telegramID, InterfaceLanguage: interfaceLanguage}
	plan = localizedPremiumPlan(user, plan)
	channel, _ := metadataString(payment.Metadata, "channel")
	if channel != "web" && telegramID > 0 {
		_ = b.telegram.sendMarkdownMessageWithCopy(r.Context(), telegramID, premiumActivationMarkdown(user, plan, until.Format("2006-01-02"), false), ui(user))
	}
	w.WriteHeader(http.StatusOK)
}

func (b *bot) sendProgress(ctx context.Context, chatID int64, user userState) error {
	level, levelName, currentXP, neededXP := knowledgeLevel(user.XP)
	language := userLearningLanguage(user)
	copy := ui(user)
	return b.telegram.sendMarkdownMessageWithCopy(ctx, chatID,
		"*"+escapeMarkdownV2(copy.Progress)+"* 📊\n"+
			escapeMarkdownV2(copy.Premium)+": "+escapeMarkdownV2(planName(user))+"\n"+
			escapeMarkdownV2(copy.LearningLanguage)+": *"+escapeMarkdownV2(language.InterfaceName)+"*\n"+
			"XP level: *"+itoa(level)+"/20* \\- "+escapeMarkdownV2(levelName)+"\n"+
			"XP: *"+itoa(user.XP)+"* \\("+itoa(currentXP)+"/"+itoa(neededXP)+"\\)\n"+
			escapeMarkdownV2(copy.LevelTest)+": `"+escapeMarkdownV2(user.Level)+"`\n"+
			escapeMarkdownV2(copy.NewLesson)+": *"+itoa(user.LessonCount)+"*\n"+
			escapeMarkdownV2(copy.Practice)+": *"+itoa(user.PracticeCount)+"*\n"+
			escapeMarkdownV2(copy.Tool.VoiceToText)+": *"+itoa(user.VoiceCount)+"*\n"+
			escapeMarkdownV2(fmt.Sprintf(copy.TotalLearned, len(learnedWordsForLanguage(user))))+"\n"+
			escapeMarkdownV2(copy.WordGame)+": *"+itoa(user.WordGameCount)+"*\n"+
			escapeMarkdownV2(copy.Mistakes)+": *"+itoa(len(mistakesForLanguage(user)))+"*",
		copy)
}

func (b *bot) sendLeaderboard(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	entries := b.store.leaderboard(10)
	if len(entries) == 0 {
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			"*"+escapeMarkdownV2(copy.Leaders)+"* 🏆\n\n"+escapeMarkdownV2(copy.MainMenuBody),
			leaderboardMenuKeyboard(copy))
	}

	var builder strings.Builder
	builder.WriteString("*" + escapeMarkdownV2(copy.Leaders) + "* 🏆\n\n")
	for i, entry := range entries {
		place := itoa(i + 1)
		if i == 0 {
			place = "🥇"
		} else if i == 1 {
			place = "🥈"
		} else if i == 2 {
			place = "🥉"
		} else {
			place += "\\."
		}
		builder.WriteString(place + " ")
		builder.WriteString("*" + escapeMarkdownV2(entry.FirstName) + "*")
		builder.WriteString(" \\- " + itoa(entry.Score) + " XP\n")
		if len(entry.Languages) > 0 {
			builder.WriteString("   " + escapeMarkdownV2(copy.LearningLanguage) + ": *" + escapeMarkdownV2(strings.Join(entry.Languages, ", ")) + "*\n")
		}
		builder.WriteString("   " + escapeMarkdownV2(copy.Words) + ": *" + itoa(entry.Words) + "* · " + escapeMarkdownV2(copy.Mistakes) + ": *" + itoa(entry.Mistakes) + "*\n")
		builder.WriteString("   lvl " + itoa(entry.Level) + "/20 · " + escapeMarkdownV2(entry.Title) + "\n\n")
	}
	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, builder.String(), leaderboardMenuKeyboard(copy))
}

func (b *bot) sendLeaderboardMenu(ctx context.Context, chatID int64, users ...userState) error {
	copy := uiFromOptionalUser(users)
	return b.telegram.sendInlineMessage(ctx, chatID, copy.Leaders+"\n\n"+copy.MainMenuBody, leaderboardMenuKeyboard(copy))
}

func (b *bot) sendLanguageLeaderboard(ctx context.Context, chatID int64, languageCode string, users ...userState) error {
	copy := uiFromOptionalUser(users)
	language := learningLanguageByCode(languageCode)
	entries := b.store.languageLeaderboard(language.Code, 10)
	if len(entries) == 0 {
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			"*"+escapeMarkdownV2(copy.Leaders)+"* 🏆\n\n"+escapeMarkdownV2(copy.LearningLanguage)+": *"+escapeMarkdownV2(language.InterfaceName)+"*",
			leaderboardMenuKeyboard(copy))
	}

	var builder strings.Builder
	builder.WriteString("*" + escapeMarkdownV2(copy.Leaders) + ": " + escapeMarkdownV2(language.InterfaceName) + "* 🏆\n\n")
	for i, entry := range entries {
		place := itoa(i + 1)
		if i == 0 {
			place = "🥇"
		} else if i == 1 {
			place = "🥈"
		} else if i == 2 {
			place = "🥉"
		} else {
			place += "\\."
		}
		builder.WriteString(place + " ")
		builder.WriteString("*" + escapeMarkdownV2(entry.FirstName) + "*")
		builder.WriteString(" \\- " + itoa(entry.Score) + " XP\n")
		builder.WriteString("   " + escapeMarkdownV2(copy.LevelTest) + " `" + escapeMarkdownV2(entry.Level) + "` · " + escapeMarkdownV2(copy.Words) + ": *" + itoa(entry.Words) + "* · " + escapeMarkdownV2(copy.Mistakes) + ": *" + itoa(entry.Mistakes) + "*\n\n")
	}
	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, builder.String(), leaderboardMenuKeyboard(copy))
}

func (b *bot) sendMistakes(ctx context.Context, chatID int64, user userState, page int) error {
	mistakes := mistakesForLanguage(user)
	if len(mistakes) == 0 {
		return b.telegram.sendInlineMarkdownMessage(ctx, chatID,
			mistakesEmptyMarkdown(user),
			backToMenuKeyboard(ui(user)))
	}
	totalPages := mistakeTotalPages(len(mistakes))
	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}
	return b.telegram.sendInlineMarkdownMessage(ctx, chatID, formatMistakesMessage(mistakes, page, ui(user)), mistakesActionsKeyboard(page, totalPages, ui(user)))
}

func (b *bot) clearMistakes(ctx context.Context, chatID int64, user userState) error {
	if len(mistakesForLanguage(user)) == 0 {
		return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).Mistakes, ui(user))
	}
	if err := b.store.clearMistakes(user.TelegramID, user.LearningLanguage); err != nil {
		return err
	}
	return b.telegram.sendMessageWithCopy(ctx, chatID, ui(user).Correct+" ✓", ui(user))
}

func (b *bot) premiumText(user userState) string {
	premiumMonthly, _ := b.premiumPlan(premiumMonthlyProduct)
	premiumYearly, _ := b.premiumPlan(premiumYearlyProduct)
	platinumMonthly, _ := b.premiumPlan(platinumMonthlyProduct)
	platinumYearly, _ := b.premiumPlan(platinumYearlyProduct)
	copy := premiumUI(user)
	return "*Free*\n" +
		"• " + escapeMarkdownV2(fmt.Sprintf(copy.LessonsPerDay, freeLessonLimit)) + "\n" +
		"• " + escapeMarkdownV2(fmt.Sprintf(copy.PracticeMessagesPerDay, freePracticeLimit)) + "\n" +
		"• " + escapeMarkdownV2(copy.FreeVoiceUnavailable) + "\n\n" +
		"*Premium* ⭐ " + escapeMarkdownV2(copy.BestValueLabel) + "\n" +
		"• " + escapeMarkdownV2(fmt.Sprintf(copy.LessonsPerDay, premiumLessonLimit)) + "\n" +
		"• " + escapeMarkdownV2(fmt.Sprintf(copy.PracticeMessagesPerDay, premiumPracticeLimit)) + "\n" +
		"• " + escapeMarkdownV2(fmt.Sprintf(copy.PremiumVoicesPerDay, premiumVoiceLimit, b.cfg.MaxVoiceSeconds)) + "\n" +
		"• " + escapeMarkdownV2(copy.VoiceTextAndTranslation) + "\n" +
		"• " + escapeMarkdownV2(copy.ImageTextTranslation) + "\n" +
		"• " + escapeMarkdownV2(copy.VoicePhotoContextPractice) + "\n\n" +
		"*" + escapeMarkdownV2(fmt.Sprintf(copy.MonthPrice, premiumMonthly.RubPrice, premiumMonthly.StarsPrice)) + "*\n" +
		"*" + escapeMarkdownV2(fmt.Sprintf(copy.YearPrice, premiumYearly.RubPrice, premiumYearly.StarsPrice, annualDiscountPercent)) + "*\n\n" +
		"*Platinum* 💠 " + escapeMarkdownV2(copy.MaxAccessLabel) + "\n" +
		"• " + escapeMarkdownV2(fmt.Sprintf(copy.LessonsPerDay, platinumLessonLimit)) + "\n" +
		"• " + escapeMarkdownV2(fmt.Sprintf(copy.PracticeMessagesPerDay, platinumPracticeLimit)) + "\n" +
		"• " + escapeMarkdownV2(fmt.Sprintf(copy.PremiumVoicesPerDay, platinumVoiceLimit, b.cfg.MaxVoiceSeconds)) + "\n" +
		"• " + escapeMarkdownV2(copy.PlatinumPriority) + "\n\n" +
		"*" + escapeMarkdownV2(fmt.Sprintf(copy.PlatinumMonthPrice, platinumMonthly.RubPrice, platinumMonthly.StarsPrice)) + "*\n" +
		"*" + escapeMarkdownV2(fmt.Sprintf(copy.PlatinumYearPrice, platinumYearly.RubPrice, platinumYearly.StarsPrice, annualDiscountPercent)) + "*\n\n" +
		"_" + escapeMarkdownV2(copy.InviteFree) + "_"
}

func (b *bot) limitsText(user userState) string {
	lessonLimit, practiceLimit := limitsFor(user)
	copy := premiumUI(user)
	return fmt.Sprintf(
		"%s: %s\n%s: %d/%d\n%s: %d/%d\n%s: %d/%d",
		copy.PlanLabel,
		planName(user),
		copy.LessonsToday,
		user.LessonsToday, lessonLimit,
		copy.PracticeToday,
		user.PracticeToday, practiceLimit,
		copy.VoicesToday,
		user.VoiceToday, voiceLimitFor(user),
	)
}

func (b *bot) limitsMarkdownText(user userState) string {
	lessonLimit, practiceLimit := limitsFor(user)
	copy := premiumUI(user)
	return fmt.Sprintf(
		"*%s* ⚡\n%s: %s\n%s: *%d/%d*\n%s: *%d/%d*\n%s: *%d/%d*",
		escapeMarkdownV2(copy.LimitsTitle),
		escapeMarkdownV2(copy.PlanLabel),
		escapeMarkdownV2(planName(user)),
		escapeMarkdownV2(copy.LessonsToday),
		user.LessonsToday, lessonLimit,
		escapeMarkdownV2(copy.PracticeToday),
		user.PracticeToday, practiceLimit,
		escapeMarkdownV2(copy.VoicesToday),
		user.VoiceToday, voiceLimitFor(user),
	)
}

func (b *bot) limitReachedText(kind string, user userState) string {
	return fmt.Sprintf(systemUI(user).LimitReached, kind, b.limitsText(user))
}

func canUseLessons(user userState) bool {
	limit, _ := limitsFor(user)
	return user.LessonsToday < limit
}

func canUsePractice(user userState) bool {
	_, limit := limitsFor(user)
	return user.PracticeToday < limit
}

func limitsFor(user userState) (int, int) {
	if user.isPlatinum(time.Now()) {
		return platinumLessonLimit, platinumPracticeLimit
	}
	if user.isPremium(time.Now()) {
		return premiumLessonLimit, premiumPracticeLimit
	}
	return freeLessonLimit, freePracticeLimit
}

func voiceLimitFor(user userState) int {
	if user.isPlatinum(time.Now()) {
		return platinumVoiceLimit
	}
	if user.isPremium(time.Now()) {
		return premiumVoiceLimit
	}
	return 0
}

func paidPlanLimits(tier string) (lessons int, practice int, voices int) {
	if tier == "platinum" {
		return platinumLessonLimit, platinumPracticeLimit, platinumVoiceLimit
	}
	return premiumLessonLimit, premiumPracticeLimit, premiumVoiceLimit
}

func planName(user userState) string {
	if user.isPlatinum(time.Now()) {
		return "Platinum " + user.PremiumUntil.Format("2006-01-02")
	}
	if user.isPremium(time.Now()) {
		return "Premium " + user.PremiumUntil.Format("2006-01-02")
	}
	return "Free"
}

func parseReferralPayload(payload string) (int64, error) {
	raw := strings.TrimPrefix(payload, "ref_")
	var id int64
	_, err := fmt.Sscanf(raw, "%d", &id)
	return id, err
}

func parseChoiceCallback(data string) (correctID string, answerID string, ok bool) {
	parts := strings.Split(data, "|")
	if len(parts) != 3 {
		return "", "", false
	}
	return parts[1], parts[2], parts[1] != "" && parts[2] != ""
}

func parseVocabularyPageCallback(data string) (int, bool) {
	raw := strings.TrimPrefix(data, "vocab|")
	var page int
	if _, err := fmt.Sscanf(raw, "%d", &page); err != nil {
		return 0, false
	}
	return page, page >= 0
}

func parseMistakesPageCallback(data string) (int, bool) {
	raw := strings.TrimPrefix(data, "mistakes|")
	var page int
	if _, err := fmt.Sscanf(raw, "%d", &page); err != nil {
		return 0, false
	}
	return page, page >= 0
}

func learnedWordByID(user userState, id string) (learnedWordEntry, bool) {
	for _, word := range user.LearnedWords {
		if word.ID == id || legacyVocabID(word.ID) == id || word.ID == legacyVocabID(id) {
			return word, true
		}
	}
	return learnedWordEntry{}, false
}

func wordWillBeMasteredAfterCorrect(user userState, id string) bool {
	word, ok := learnedWordByID(user, id)
	if !ok || learnedWordMastered(word) {
		return false
	}
	return learnedWordTrainingCount(word)+1 >= learnedWordMasteryThreshold
}

func learnedWordsForLanguage(user userState) []learnedWordEntry {
	return learnedWordsForLanguageCode(user, user.LearningLanguage)
}

func learnedWordsForLanguageCode(user userState, languageCode string) []learnedWordEntry {
	language := normalizeLearningLanguage(languageCode)
	words := make([]learnedWordEntry, 0, len(user.LearnedWords))
	for _, word := range user.LearnedWords {
		if normalizeLearningLanguage(word.Language) == language && learnedWordMastered(word) {
			words = append(words, word)
		}
	}
	sort.SliceStable(words, func(i, j int) bool {
		return words[i].LearnedAt.After(words[j].LearnedAt)
	})
	return words
}

func mistakesForLanguage(user userState) []mistakeEntry {
	return mistakesForLanguageCode(user, user.LearningLanguage)
}

func mistakesForLanguageCode(user userState, languageCode string) []mistakeEntry {
	language := normalizeLearningLanguage(languageCode)
	mistakes := make([]mistakeEntry, 0, len(user.Mistakes))
	for _, mistake := range user.Mistakes {
		if normalizeLearningLanguage(mistake.Language) == language {
			mistakes = append(mistakes, mistake)
		}
	}
	return mistakes
}

func normalizeAnswer(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.Trim(value, " \t\r\n.!?,;:\"'`")
	return strings.Join(strings.Fields(value), " ")
}

func spellingHint(value string) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) == 0 {
		return ""
	}
	for i := 1; i < len(runes); i++ {
		if runes[i] != ' ' && runes[i] != '-' && runes[i] != '\'' {
			runes[i] = '_'
		}
	}
	return string(runes)
}

func isStopIntent(text string) bool {
	normalized := strings.ToLower(strings.TrimSpace(text))
	normalized = strings.TrimLeftFunc(normalized, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= 'а' && r <= 'я') || r == 'ё')
	})
	normalized = strings.TrimSpace(normalized)
	normalized = strings.Trim(normalized, ".!?, ")
	normalized = strings.ReplaceAll(normalized, "ё", "е")

	stopPhrases := map[string]bool{
		"stop": true, "стоп": true, "хватит": true, "закончить": true,
		"я хочу закончить": true, "давай закончим": true,
		"все": true, "всё": true, "это все": true, "это всё": true,
		"its all": true, "it's all": true, "that is all": true, "that's all": true,
		"i want to finish": true, "i want to stop": true, "finish": true,
	}
	if stopPhrases[normalized] {
		return true
	}
	for _, m := range []string{
		"хочу закончить", "хочу остановиться", "хочу выйти",
		"давай закончим", "давай остановимся", "можем закончить",
		"можно закончить", "закончим практику", "останови практику",
		"остановить практику", "выйти из практики",
		"на сегодня хватит", "пока хватит", "хватит практики",
	} {
		if strings.Contains(normalized, m) {
			return true
		}
	}
	for _, m := range []string{
		"i want to finish", "i want to stop", "let's finish", "lets finish",
		"let's stop", "lets stop", "stop practice", "finish practice",
		"end practice", "enough for today", "that's enough", "that is enough",
	} {
		if strings.Contains(normalized, m) {
			return true
		}
	}
	return false
}

func isMenuCommandText(text string) bool {
	trimmed := strings.TrimSpace(text)
	if strings.EqualFold(trimmed, "/menu") {
		return true
	}
	for _, copy := range uiCopies {
		if trimmed == copy.MenuButton {
			return true
		}
	}
	for _, language := range interfaceLanguages() {
		if trimmed == ui(userState{InterfaceLanguage: language.Code}).MenuButton {
			return true
		}
	}
	return trimmed == "📋 Меню" || trimmed == "📋 Menu"
}

func isStopCommandText(text string) bool {
	trimmed := strings.TrimSpace(text)
	if strings.EqualFold(trimmed, "/stop") || strings.EqualFold(trimmed, "/resetmode") {
		return true
	}
	for _, copy := range uiCopies {
		if trimmed == copy.StopButton || strings.TrimPrefix(trimmed, "⏹ ") == strings.TrimPrefix(copy.StopButton, "⏹ ") {
			return true
		}
	}
	for _, language := range interfaceLanguages() {
		copy := ui(userState{InterfaceLanguage: language.Code})
		if trimmed == copy.StopButton || strings.TrimPrefix(trimmed, "⏹ ") == strings.TrimPrefix(copy.StopButton, "⏹ ") {
			return true
		}
	}
	return trimmed == "⏹ Стоп" || trimmed == "Стоп" || trimmed == "⏹ Stop" || strings.EqualFold(trimmed, "Stop")
}

func backToMenuKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func mistakesActionsKeyboard(page int, totalPages int, copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	rows := [][]map[string]any{}
	if totalPages > 1 {
		row := []map[string]any{}
		if page > 0 {
			row = append(row, map[string]any{"text": copy.Back, "callback_data": "mistakes|" + itoa(page-1)})
		}
		if page+1 < totalPages {
			row = append(row, map[string]any{"text": "▶", "callback_data": "mistakes|" + itoa(page+1)})
		}
		rows = append(rows, row)
	}
	rows = append(rows,
		[]map[string]any{{"text": "🛠 " + copy.Mistakes, "callback_data": "practice_mistakes"}},
		[]map[string]any{{"text": "🗑 " + copy.Mistakes, "callback_data": "clear_mistakes"}},
		[]map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}},
	)
	return map[string]any{
		"inline_keyboard": rows,
	}
}

func spellingPracticeKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": copy.BackMenu, "callback_data": "stop_mode_menu"}},
		},
	}
}

func spellingDoneKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "✍️ " + copy.NextWord, "callback_data": "menu_spelling"}},
			{{"text": "🎮 " + copy.WordGame, "callback_data": "menu_word_game"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func mistakePracticeKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "📖 " + copy.Mistakes, "callback_data": "menu_mistakes"}},
			{{"text": copy.BackMenu, "callback_data": "stop_mode_menu"}},
		},
	}
}

func mistakePracticeDoneKeyboard(nextText string, hasNext bool, copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	rows := [][]map[string]any{}
	if hasNext {
		rows = append(rows, []map[string]any{{"text": "🛠 " + nextText, "callback_data": "practice_mistakes"}})
	}
	rows = append(rows,
		[]map[string]any{{"text": "📖 " + copy.Mistakes, "callback_data": "menu_mistakes"}},
		[]map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}},
	)
	return map[string]any{"inline_keyboard": rows}
}
