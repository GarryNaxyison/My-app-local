package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	privacyPolicyURL       = "https://neriva.ru/privacy.html"
	personalDataConsentURL = "https://neriva.ru/consent.html"
	userAgreementURL       = "https://neriva.ru/agreement.html"
)

type telegramEditTarget struct {
	ChatID    int64
	MessageID int64
}

type telegramEditTargetKey struct{}

func contextWithTelegramEditTarget(ctx context.Context, chatID int64, messageID int64) context.Context {
	if chatID == 0 || messageID == 0 {
		return ctx
	}
	return context.WithValue(ctx, telegramEditTargetKey{}, telegramEditTarget{
		ChatID:    chatID,
		MessageID: messageID,
	})
}

func telegramEditTargetFromContext(ctx context.Context) (telegramEditTarget, bool) {
	target, ok := ctx.Value(telegramEditTargetKey{}).(telegramEditTarget)
	return target, ok && target.ChatID != 0 && target.MessageID != 0
}

type telegramClient struct {
	baseURL     string
	fileBaseURL string
	http        *http.Client
}

type telegramUpdate struct {
	UpdateID         int64             `json:"update_id"`
	Message          *telegramMessage  `json:"message"`
	PreCheckoutQuery *preCheckoutQuery `json:"pre_checkout_query"`
	CallbackQuery    *callbackQuery    `json:"callback_query"`
}

type telegramMessage struct {
	MessageID         int64              `json:"message_id"`
	From              telegramUser       `json:"from"`
	Chat              telegramChat       `json:"chat"`
	Text              string             `json:"text"`
	Caption           string             `json:"caption"`
	ReplyToMessage    *telegramMessage   `json:"reply_to_message"`
	Voice             *telegramVoice     `json:"voice"`
	Photo             []telegramPhoto    `json:"photo"`
	MediaGroupID      string             `json:"media_group_id"`
	SuccessfulPayment *successfulPayment `json:"successful_payment"`
}

type telegramUser struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

func telegramDisplayName(user telegramUser) string {
	if username := strings.TrimSpace(user.Username); username != "" {
		return username
	}
	if firstName := strings.TrimSpace(user.FirstName); firstName != "" {
		return firstName
	}
	if user.ID != 0 {
		return "telegram_" + strconv.FormatInt(user.ID, 10)
	}
	return ""
}

type telegramChat struct {
	ID int64 `json:"id"`
}

type telegramVoice struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Duration     int    `json:"duration"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

type telegramPhoto struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

type telegramFile struct {
	FileID   string `json:"file_id"`
	FilePath string `json:"file_path"`
	FileSize int    `json:"file_size"`
}

type callbackQuery struct {
	ID      string           `json:"id"`
	From    telegramUser     `json:"from"`
	Message *telegramMessage `json:"message"`
	Data    string           `json:"data"`
}

type preCheckoutQuery struct {
	ID             string       `json:"id"`
	From           telegramUser `json:"from"`
	Currency       string       `json:"currency"`
	TotalAmount    int          `json:"total_amount"`
	InvoicePayload string       `json:"invoice_payload"`
}

type successfulPayment struct {
	Currency                string `json:"currency"`
	TotalAmount             int    `json:"total_amount"`
	InvoicePayload          string `json:"invoice_payload"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
}

func newTelegramClient(token string, httpClient *http.Client) *telegramClient {
	return &telegramClient{
		baseURL:     "https://api.telegram.org/bot" + token,
		fileBaseURL: "https://api.telegram.org/file/bot" + token,
		http:        httpClient,
	}
}

func (c *telegramClient) getMe(ctx context.Context) (telegramUser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/getMe", nil)
	if err != nil {
		return telegramUser{}, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return telegramUser{}, err
	}
	defer resp.Body.Close()

	var decoded struct {
		OK          bool         `json:"ok"`
		Result      telegramUser `json:"result"`
		Description string       `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return telegramUser{}, err
	}
	if !decoded.OK {
		return telegramUser{}, fmt.Errorf("telegram getMe failed: %s", decoded.Description)
	}
	if decoded.Result.Username == "" {
		return telegramUser{}, fmt.Errorf("bot username is empty")
	}
	return decoded.Result, nil
}

func (c *telegramClient) getUpdates(ctx context.Context, offset int64) ([]telegramUpdate, error) {
	values := url.Values{}
	values.Set("offset", strconv.FormatInt(offset, 10))
	values.Set("timeout", "30")
	values.Set("allowed_updates", `["message","pre_checkout_query","callback_query"]`)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/getUpdates?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var decoded struct {
		OK          bool             `json:"ok"`
		Result      []telegramUpdate `json:"result"`
		Description string           `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}
	if !decoded.OK {
		return nil, fmt.Errorf("telegram getUpdates failed: %s", decoded.Description)
	}
	return decoded.Result, nil
}

func (c *telegramClient) downloadFile(ctx context.Context, fileID string) ([]byte, error) {
	file, err := c.getFile(ctx, fileID)
	if err != nil {
		return nil, err
	}
	if file.FilePath == "" {
		return nil, fmt.Errorf("telegram file path is empty")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.fileBaseURL+"/"+file.FilePath, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("telegram file download failed: status %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

func (c *telegramClient) getFile(ctx context.Context, fileID string) (telegramFile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/getFile?file_id="+url.QueryEscape(fileID), nil)
	if err != nil {
		return telegramFile{}, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return telegramFile{}, err
	}
	defer resp.Body.Close()

	var decoded struct {
		OK          bool         `json:"ok"`
		Result      telegramFile `json:"result"`
		Description string       `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return telegramFile{}, err
	}
	if !decoded.OK {
		return telegramFile{}, fmt.Errorf("telegram getFile failed: %s", decoded.Description)
	}
	return decoded.Result, nil
}

func (c *telegramClient) sendMessage(ctx context.Context, chatID int64, text string) error {
	return c.sendMessageWithCopy(ctx, chatID, text, englishUICopy())
}

func (c *telegramClient) sendMessageToChat(ctx context.Context, chatID any, text string) error {
	for _, chunk := range splitTelegramText(text, 3900) {
		if err := c.call(ctx, "sendMessage", map[string]any{
			"chat_id": chatID,
			"text":    chunk,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (c *telegramClient) sendInlineMessageToChat(ctx context.Context, chatID any, text string, keyboard map[string]any) (int64, error) {
	var firstMessageID int64
	for _, chunk := range splitTelegramText(text, 3900) {
		payload := map[string]any{
			"chat_id": chatID,
			"text":    chunk,
		}
		if keyboard != nil {
			payload["reply_markup"] = keyboard
		}
		messageID, err := c.callMessageID(ctx, "sendMessage", payload)
		if err != nil {
			return firstMessageID, err
		}
		if firstMessageID == 0 {
			firstMessageID = messageID
		}
	}
	return firstMessageID, nil
}

func (c *telegramClient) sendPhotoFile(ctx context.Context, chatID int64, photoPath string, caption string) error {
	return c.sendPhotoFileToChat(ctx, chatID, photoPath, caption)
}

func (c *telegramClient) sendPhotoFileToChat(ctx context.Context, chatID any, photoPath string, caption string) error {
	file, err := os.Open(photoPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("chat_id", fmt.Sprint(chatID)); err != nil {
		return err
	}
	caption = strings.TrimSpace(caption)
	if len([]rune(caption)) > 1024 {
		caption = string([]rune(caption)[:1021]) + "..."
	}
	if caption != "" {
		if err := writer.WriteField("caption", caption); err != nil {
			return err
		}
	}
	part, err := writer.CreateFormFile("photo", filepath.Base(photoPath))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return c.callMultipart(ctx, "sendPhoto", writer.FormDataContentType(), &body)
}

// sendMessageWithCopy sends a plain-text message with the persistent bottom panel (/menu /stop).
func (c *telegramClient) sendMessageWithCopy(ctx context.Context, chatID int64, text string, copy uiCopy) error {
	for _, chunk := range splitTelegramText(text, 3900) {
		if err := c.call(ctx, "sendMessage", map[string]any{
			"chat_id":      chatID,
			"text":         chunk,
			"reply_markup": bottomReplyKeyboard(copy),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (c *telegramClient) sendMarkdownMessage(ctx context.Context, chatID int64, text string) error {
	return c.sendMarkdownMessageWithCopy(ctx, chatID, text, englishUICopy())
}

// sendMarkdownMessageWithCopy sends a MarkdownV2 message with the persistent bottom panel.
func (c *telegramClient) sendMarkdownMessageWithCopy(ctx context.Context, chatID int64, text string, copy uiCopy) error {
	for _, chunk := range splitTelegramText(text, 3900) {
		if err := c.call(ctx, "sendMessage", map[string]any{
			"chat_id":      chatID,
			"text":         chunk,
			"parse_mode":   "MarkdownV2",
			"reply_markup": bottomReplyKeyboard(copy),
		}); err != nil {
			if isTelegramMarkdownParseError(err) {
				if fallbackErr := c.sendMessageWithCopy(ctx, chatID, plainTextFromMarkdownV2(chunk), copy); fallbackErr != nil {
					return fallbackErr
				}
				continue
			}
			return err
		}
	}
	return nil
}

func (c *telegramClient) sendInlineMessage(ctx context.Context, chatID int64, text string, keyboard map[string]any) error {
	if target, ok := telegramEditTargetFromContext(ctx); ok {
		if err := c.editInlineMessage(ctx, target.ChatID, target.MessageID, text, "", keyboard); err != nil {
			if isTelegramMessageNotModified(err) {
				return nil
			}
		} else {
			return nil
		}
	}

	payload := map[string]any{
		"chat_id": chatID,
		"text":    text,
	}
	if keyboard != nil {
		payload["reply_markup"] = keyboard
	}
	return c.call(ctx, "sendMessage", payload)
}

func (c *telegramClient) sendInlineMarkdownMessage(ctx context.Context, chatID int64, text string, keyboard map[string]any) error {
	if target, ok := telegramEditTargetFromContext(ctx); ok {
		if err := c.editInlineMessage(ctx, target.ChatID, target.MessageID, text, "MarkdownV2", keyboard); err != nil {
			if isTelegramMessageNotModified(err) {
				return nil
			}
			if isTelegramMarkdownParseError(err) {
				if fallbackErr := c.editInlineMessage(ctx, target.ChatID, target.MessageID, plainTextFromMarkdownV2(text), "", keyboard); fallbackErr == nil || isTelegramMessageNotModified(fallbackErr) {
					return nil
				}
			}
		} else {
			return nil
		}
	}

	payload := map[string]any{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "MarkdownV2",
	}
	if keyboard != nil {
		payload["reply_markup"] = keyboard
	}
	if err := c.call(ctx, "sendMessage", payload); err != nil {
		if isTelegramMarkdownParseError(err) {
			payload["text"] = plainTextFromMarkdownV2(text)
			delete(payload, "parse_mode")
			return c.call(ctx, "sendMessage", payload)
		}
		return err
	}
	return nil
}

func (c *telegramClient) editInlineMessage(ctx context.Context, chatID int64, messageID int64, text string, parseMode string, keyboard map[string]any) error {
	return c.editInlineMessageToChat(ctx, chatID, messageID, text, parseMode, keyboard)
}

func (c *telegramClient) editInlineMessageToChat(ctx context.Context, chatID any, messageID int64, text string, parseMode string, keyboard map[string]any) error {
	payload := map[string]any{
		"chat_id":    chatID,
		"message_id": messageID,
		"text":       text,
	}
	if parseMode != "" {
		payload["parse_mode"] = parseMode
	}
	if keyboard != nil {
		payload["reply_markup"] = keyboard
	}
	return c.call(ctx, "editMessageText", payload)
}

func isTelegramMessageNotModified(err error) bool {
	return err != nil && strings.Contains(err.Error(), "message is not modified")
}

func isTelegramMarkdownParseError(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "can't parse entities") ||
		strings.Contains(message, "parse entities") ||
		strings.Contains(message, "markdown")
}

func plainTextFromMarkdownV2(text string) string {
	text = strings.NewReplacer(
		"\\_", "_",
		"\\*", "*",
		"\\[", "[",
		"\\]", "]",
		"\\(", "(",
		"\\)", ")",
		"\\~", "~",
		"\\`", "`",
		"\\>", ">",
		"\\#", "#",
		"\\+", "+",
		"\\-", "-",
		"\\=", "=",
		"\\|", "|",
		"\\{", "{",
		"\\}", "}",
		"\\.", ".",
		"\\!", "!",
	).Replace(text)
	text = strings.ReplaceAll(text, "*", "")
	text = strings.ReplaceAll(text, "_", "")
	text = strings.ReplaceAll(text, "`", "")
	return text
}

func (c *telegramClient) sendSmoothMessage(ctx context.Context, chatID int64, text string) error {
	return c.sendSmoothMessageWithCopy(ctx, chatID, text, englishUICopy())
}

func (c *telegramClient) sendSmoothMessageWithCopy(ctx context.Context, chatID int64, text string, copy uiCopy) error {
	if len(text) > 1200 {
		return c.sendMessageWithCopy(ctx, chatID, text, copy)
	}

	draftID := time.Now().UnixNano()
	for _, draft := range smoothDrafts(text) {
		_ = c.call(ctx, "sendMessageDraft", map[string]any{
			"chat_id":  chatID,
			"draft_id": draftID,
			"text":     draft,
		})
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(90 * time.Millisecond):
		}
	}

	_ = c.call(ctx, "deleteMessageDraft", map[string]any{
		"chat_id":  chatID,
		"draft_id": draftID,
	})
	return c.sendMessageWithCopy(ctx, chatID, text, copy)
}

func (c *telegramClient) sendChatAction(ctx context.Context, chatID int64, action string) error {
	return c.call(ctx, "sendChatAction", map[string]any{
		"chat_id": chatID,
		"action":  action,
	})
}

func (c *telegramClient) sendAudioBytes(ctx context.Context, chatID int64, audio []byte, filename string, title string) (int64, error) {
	return c.sendAudioBytesWithCopy(ctx, chatID, audio, filename, title, englishUICopy())
}

func (c *telegramClient) sendAudioBytesWithCopy(ctx context.Context, chatID int64, audio []byte, filename string, title string, copy uiCopy) (int64, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if err := writer.WriteField("chat_id", strconv.FormatInt(chatID, 10)); err != nil {
		return 0, err
	}
	if title != "" {
		if err := writer.WriteField("title", title); err != nil {
			return 0, err
		}
	}
	if err := writer.WriteField("reply_markup", mustJSONText(bottomReplyKeyboard(copy))); err != nil {
		return 0, err
	}
	part, err := writer.CreateFormFile("audio", filename)
	if err != nil {
		return 0, err
	}
	if _, err := part.Write(audio); err != nil {
		return 0, err
	}
	if err := writer.Close(); err != nil {
		return 0, err
	}
	return c.callMultipartMessageID(ctx, "sendAudio", writer.FormDataContentType(), &body)
}

func (c *telegramClient) deleteMessage(ctx context.Context, chatID int64, messageID int64) error {
	return c.call(ctx, "deleteMessage", map[string]any{
		"chat_id":    chatID,
		"message_id": messageID,
	})
}

func (c *telegramClient) sendInvoice(ctx context.Context, chatID int64, title, description, payload string, stars int, copy uiCopy, payButtonTemplate string) error {
	return c.call(ctx, "sendInvoice", map[string]any{
		"chat_id":        chatID,
		"title":          title,
		"description":    description,
		"payload":        payload,
		"provider_token": "",
		"currency":       "XTR",
		"prices": []map[string]any{
			{
				"label":  title,
				"amount": stars,
			},
		},
		"reply_markup": invoiceInlineKeyboard(stars, payButtonTemplate, copy),
	})
}

func mustJSONText(value any) string {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func (c *telegramClient) answerPreCheckoutQuery(ctx context.Context, queryID string, ok bool, errorMessage string) error {
	payload := map[string]any{
		"pre_checkout_query_id": queryID,
		"ok":                    ok,
	}
	if !ok && errorMessage != "" {
		payload["error_message"] = errorMessage
	}
	return c.call(ctx, "answerPreCheckoutQuery", payload)
}

func (c *telegramClient) answerCallbackQuery(ctx context.Context, queryID string, text string) error {
	payload := map[string]any{
		"callback_query_id": queryID,
	}
	if text != "" {
		payload["text"] = text
	}
	return c.call(ctx, "answerCallbackQuery", payload)
}

// bottomReplyKeyboard is the persistent bottom panel that stays available in text modes.
func bottomReplyKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"keyboard": [][]map[string]string{
			{
				{"text": copy.MenuButton},
				{"text": copy.StopButton},
			},
		},
		"resize_keyboard":   true,
		"is_persistent":     true,
		"one_time_keyboard": false,
	}
}

// mainMenuInlineKeyboard — inline buttons sent as a chat message when the
// user opens the menu. Each button fires a callback_query.
func keyboardCopy(copies []uiCopy) uiCopy {
	if len(copies) > 0 {
		return copies[0]
	}
	return englishUICopy()
}

func mainMenuInlineKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	webAppURL := strings.TrimSpace(os.Getenv("WEB_APP_URL"))
	if webAppURL == "" {
		webAppURL = "https://neriva.ru/app"
	}
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{
				{"text": "📝 " + copy.NewLesson, "callback_data": "menu_lesson"},
				{"text": "💬 " + copy.Practice, "callback_data": "menu_practice"},
			},
			{
				{"text": "\U0001F3A7 " + copy.Shadowing, "callback_data": "menu_shadowing"},
				{"text": "🗣 " + copy.Pronunciation, "callback_data": "menu_pronunciation"},
			},
			{
				{"text": "🤖 " + copy.AITutor, "callback_data": "menu_tutor"},
			},
			{
				{"text": "🌐 " + copy.Tool.WebApp, "url": webAppURL},
			},
			{
				{"text": "🧠 " + copy.LearnWords, "callback_data": "menu_word_lesson"},
				{"text": "🎮 " + copy.WordGame, "callback_data": "menu_word_game"},
			},
			{
				{"text": "✍️ " + copy.Spelling, "callback_data": "menu_spelling"},
				{"text": "📚 " + copy.Vocabulary, "callback_data": "menu_vocabulary"},
			},
			{
				{"text": "🔖 " + copy.Phrasebook, "callback_data": "menu_phrasebook"},
			},
			{
				{"text": "📖 " + copy.Mistakes, "callback_data": "menu_mistakes"},
				{"text": "🎯 " + copy.LevelTest, "callback_data": "menu_level_test"},
			},
			{
				{"text": "📊 " + copy.Progress, "callback_data": "menu_progress"},
				{"text": "🏆 " + copy.Leaders, "callback_data": "menu_leaderboard"},
			},
			{
				{"text": "⚡ " + copy.Limits, "callback_data": "menu_limits"},
				{"text": "🧰 " + copy.Tools, "callback_data": "menu_tools"},
			},
			{
				{"text": "🔔 " + copy.Notifications, "callback_data": "menu_reminders"},
				{"text": "⚙️ " + copy.Settings, "callback_data": "menu_settings"},
			},
			{
				{"text": "👥 " + copy.Referral, "callback_data": "menu_referral"},
				{"text": "💎 " + copy.Premium, "callback_data": "menu_premium"},
			},
		},
	}
}

func learningMenuKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "🤖 " + copy.AITutor, "callback_data": "menu_tutor"}},
			{{"text": "📝 " + copy.NewLesson, "callback_data": "menu_lesson"}},
			{{"text": "💬 " + copy.Practice, "callback_data": "menu_practice"}},
			{{"text": "\U0001F3A7 " + copy.Shadowing, "callback_data": "menu_shadowing"}},
			{{"text": "🗣 " + copy.Pronunciation, "callback_data": "menu_pronunciation"}},
			{{"text": "🎓 " + copy.LevelTest, "callback_data": "menu_level_test"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func wordsMenuKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "🧠 " + copy.LearnWords, "callback_data": "menu_word_lesson"}},
			{{"text": "🎮 " + copy.WordGame, "callback_data": "menu_word_game"}},
			{{"text": "✍️ " + copy.Spelling, "callback_data": "menu_spelling"}},
			{{"text": "📚 " + copy.Vocabulary, "callback_data": "menu_vocabulary"}},
			{{"text": "🔖 " + copy.Phrasebook, "callback_data": "menu_phrasebook"}},
			{{"text": "📖 " + copy.Mistakes, "callback_data": "menu_mistakes"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func statsMenuKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "📊 " + copy.Progress, "callback_data": "menu_progress"}},
			{{"text": "🏆 " + copy.Leaders, "callback_data": "menu_leaderboard"}},
			{{"text": "⚡ " + copy.Limits, "callback_data": "menu_limits"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func settingsMenuKeyboard(user userState) map[string]any {
	copy := ui(user)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "🌐 " + copy.BotLanguage, "callback_data": "menu_interface_language"}},
			{{"text": "🎯 " + copy.LearningLanguage, "callback_data": "menu_language"}},
			{{"text": "🔔 " + copy.Notifications, "callback_data": "menu_reminders"}},
			{{"text": "💎 " + copy.Premium, "callback_data": "menu_premium"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func leaderboardMenuKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	rows := [][]map[string]any{
		{{"text": "🌐 " + copy.Leaders, "callback_data": "leaderboard|all"}},
	}
	for i := 0; i < len(learningLanguages); i += 2 {
		row := []map[string]any{}
		for j := i; j < i+2 && j < len(learningLanguages); j++ {
			language := learningLanguages[j]
			row = append(row, map[string]any{
				"text":          language.InterfaceName,
				"callback_data": "leaderboard|" + language.Code,
			})
		}
		rows = append(rows, row)
	}
	rows = append(rows, []map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}})
	return map[string]any{
		"inline_keyboard": rows,
	}
}

func referralMenuKeyboard(user userState) map[string]any {
	copy := ui(user)
	premiumCopy := premiumUI(user)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": premiumCopy.InviteFriendButton, "callback_data": "invite_friend"}},
			{{"text": premiumCopy.ReferralWithdrawButton, "callback_data": "referral_withdraw"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func toolsInlineKeyboard(_ string, copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "🎙 " + copy.Tool.VoiceToText, "callback_data": "tool_voice_text"}},
			{{"text": "🖼 " + copy.Tool.ImageTranslate, "callback_data": "tool_image_translate"}},
			{{"text": "🌐 " + copy.Tool.Translator, "callback_data": "tool_translator"}},
			{{"text": "🤖 AI Router", "url": aiRouterTelegramURL}},
			{{"text": "🤖 " + copy.Tool.GPTAgent, "url": "https://t.me/Ton_Rivals_AI_bot"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func backToToolsKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": copy.Back, "callback_data": "menu_tools"}},
		},
	}
}

func privacyPromptKeyboard() map[string]any {
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "Продолжить", "callback_data": "privacy_continue"}},
			{{"text": "Политика обработки персональных данных", "url": privacyPolicyURL}},
			{{"text": "Согласие на обработку ПДн", "url": personalDataConsentURL}},
			{{"text": "Пользовательское соглашение", "url": userAgreementURL}},
		},
	}
}

func translatorInlineKeyboard(user userState, _ string) map[string]any {
	copy := ui(user)
	sourceCode, targetCode := parseTranslatorMode(user.Mode, user)
	rows := [][]map[string]any{
		{{"text": "↤ " + translatorLanguageDisplayName(sourceCode, user), "callback_data": "translator_source_menu"}},
		{{"text": "↦ " + translatorLanguageDisplayName(targetCode, user), "callback_data": "translator_target_menu"}},
		{{"text": "⇄", "callback_data": "translator_swap"}},
	}
	rows = append(rows, []map[string]any{{"text": copy.BackMenu, "callback_data": "menu_tools"}})
	return map[string]any{"inline_keyboard": rows}
}

func translatorLanguageKeyboard(user userState, target bool) map[string]any {
	copy := ui(user)
	prefix := "trsrc|"
	if target {
		prefix = "trdst|"
	}
	rows := [][]map[string]any{}
	if !target {
		rows = append(rows, []map[string]any{{"text": "↤ " + copy.Tool.AutoDetect, "callback_data": prefix + translationAutoCode}})
	}
	arrow := "↤ "
	if target {
		arrow = "↦ "
	}
	for _, language := range learningLanguages {
		rows = append(rows, []map[string]any{{
			"text":          arrow + languageButtonLabel(language, user.InterfaceLanguage),
			"callback_data": prefix + language.Code,
		}})
	}
	rows = append(rows, []map[string]any{{"text": copy.BackMenu, "callback_data": "tool_translator"}})
	return map[string]any{"inline_keyboard": rows}
}

func languageSelectionKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	rows := languageRows(learningLanguages, "lang")
	rows = append(rows, []map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}})
	return map[string]any{
		"inline_keyboard": rows,
	}
}

func interfaceLanguageSelectionKeyboard() map[string]any {
	rows := languageRows(interfaceLanguages(), "ui")
	return map[string]any{"inline_keyboard": rows}
}

func languageRows(languages []learningLanguage, callbackPrefix string) [][]map[string]any {
	rows := [][]map[string]any{}
	for i := 0; i < len(languages); i += 2 {
		row := []map[string]any{}
		for j := i; j < i+2 && j < len(languages); j++ {
			language := languages[j]
			row = append(row, map[string]any{
				"text":          language.InterfaceName,
				"callback_data": callbackPrefix + "|" + language.Code,
			})
		}
		rows = append(rows, row)
	}
	return rows
}

func levelAssessmentStartKeyboard(user userState) map[string]any {
	copy := ui(user)
	systemCopy := systemUI(user)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "🎯 " + systemCopy.LevelManualButton, "callback_data": "level_manual_select"}},
			{{"text": "▶ " + systemCopy.LevelStartButton, "callback_data": "level_test_start"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func manualLevelSelectionKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{
				{"text": "A1", "callback_data": "level|A1"},
				{"text": "A2", "callback_data": "level|A2"},
			},
			{
				{"text": "B1", "callback_data": "level|B1"},
				{"text": "B2", "callback_data": "level|B2"},
			},
			{
				{"text": "C1", "callback_data": "level|C1"},
				{"text": "C2", "callback_data": "level|C2"},
			},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func levelAssessmentDoneKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "📝 " + copy.NewLesson, "callback_data": "menu_lesson"}},
			{{"text": "🧠 " + copy.LearnWords, "callback_data": "menu_word_lesson"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func levelPromotionKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "📝 " + copy.NewLesson, "callback_data": "menu_lesson"}},
			{{"text": "🧠 " + copy.LearnWords, "callback_data": "menu_word_lesson"}},
		},
	}
}

func lessonDoneKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": copy.NewLesson, "callback_data": "menu_lesson"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func vocabularyKeyboard(page int, totalPages int, copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	rows := [][]map[string]any{}
	if totalPages > 1 {
		row := []map[string]any{}
		if page > 0 {
			row = append(row, map[string]any{"text": copy.Back, "callback_data": "vocab|" + strconv.Itoa(page-1)})
		}
		if page+1 < totalPages {
			row = append(row, map[string]any{"text": copy.NextPage, "callback_data": "vocab|" + strconv.Itoa(page+1)})
		}
		rows = append(rows, row)
	}
	rows = append(rows,
		[]map[string]any{{"text": "🎮 " + copy.WordGame, "callback_data": "menu_word_game"}},
		[]map[string]any{{"text": "✍️ " + copy.Spelling, "callback_data": "menu_spelling"}},
		[]map[string]any{{"text": "🧠 " + copy.LearnWords, "callback_data": "menu_word_lesson"}},
		[]map[string]any{{"text": "🔖 " + copy.Phrasebook, "callback_data": "menu_phrasebook"}},
		[]map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}},
	)
	return map[string]any{"inline_keyboard": rows}
}

func phrasebookKeyboard(page int, totalPages int, copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	rows := [][]map[string]any{}
	if totalPages > 1 {
		row := []map[string]any{}
		if page > 0 {
			row = append(row, map[string]any{"text": copy.Back, "callback_data": "phrasebook|" + strconv.Itoa(page-1)})
		}
		if page+1 < totalPages {
			row = append(row, map[string]any{"text": copy.NextPage, "callback_data": "phrasebook|" + strconv.Itoa(page+1)})
		}
		rows = append(rows, row)
	}
	rows = append(rows,
		[]map[string]any{{"text": "📚 " + copy.Vocabulary, "callback_data": "menu_vocabulary"}},
		[]map[string]any{{"text": "🧠 " + copy.LearnWords, "callback_data": "menu_word_lesson"}},
		[]map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}},
	)
	return map[string]any{"inline_keyboard": rows}
}

func vocabularyEmptyKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "🧠 " + copy.LearnWords, "callback_data": "menu_word_lesson"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func wordChoiceKeyboard(prefix string, correctID string, options []vocabWord, interfaceLanguage string) map[string]any {
	copy := ui(userState{InterfaceLanguage: interfaceLanguage})
	rows := [][]map[string]any{}
	for i := 0; i < len(options); i += 2 {
		row := []map[string]any{
			{"text": options[i].English, "callback_data": prefix + "|" + correctID + "|" + options[i].ID},
		}
		if i+1 < len(options) {
			row = append(row, map[string]any{"text": options[i+1].English, "callback_data": prefix + "|" + correctID + "|" + options[i+1].ID})
		}
		rows = append(rows, row)
	}
	rows = append(rows, []map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}})
	return map[string]any{"inline_keyboard": rows}
}

func wordLessonDoneKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "🧠 " + copy.NextWord, "callback_data": "menu_word_lesson"}},
			{{"text": "🎮 " + copy.WordGame, "callback_data": "menu_word_game"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func wordGameDoneKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "🎮 " + copy.NextWord, "callback_data": "menu_word_game"}},
			{{"text": "✍️ " + copy.Spelling, "callback_data": "menu_spelling"}},
			{{"text": "🧠 " + copy.LearnWords, "callback_data": "menu_word_lesson"}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func premiumInlineKeyboard(yooKassaEnabled bool, _ map[string][]cryptoPaymentMethod, plans []premiumPlan, user userState) map[string]any {
	copy := ui(user)
	premiumCopy := premiumUI(user)
	rows := [][]map[string]any{}
	for _, plan := range plans {
		rows = append(rows, []map[string]any{{
			"text":          premiumPlanChoiceButtonText(plan, yooKassaEnabled, nil),
			"callback_data": "premium_plan|" + plan.Product,
		}})
	}
	rows = append(rows,
		[]map[string]any{{"text": premiumCopy.InviteFriendButton, "callback_data": "invite_friend"}},
		[]map[string]any{{"text": copy.Limits, "callback_data": "show_limits"}},
		[]map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}},
	)
	return map[string]any{"inline_keyboard": rows}
}

func premiumPlanChoiceButtonText(plan premiumPlan, yooKassaEnabled bool, _ []cryptoPaymentMethod) string {
	prices := []string{strconv.Itoa(plan.StarsPrice) + " Stars"}
	if yooKassaEnabled && plan.RubPrice > 0 {
		prices = append(prices, strconv.Itoa(plan.RubPrice)+" RUB")
	}
	return plan.Title + " - " + strings.Join(prices, " / ")
}

func premiumPaymentOptionsKeyboard(yooKassaEnabled bool, _ bool, _ []cryptoPaymentMethod, plan premiumPlan, user userState) map[string]any {
	copy := ui(user)
	premiumCopy := premiumUI(user)
	rows := [][]map[string]any{}
	if plan.StarsPrice > 0 {
		rows = append(rows, []map[string]any{{"text": fmt.Sprintf(premiumCopy.PayStarsButton, plan.StarsPrice), "callback_data": "buy_stars|" + plan.Product}})
	}
	if yooKassaEnabled && plan.RubPrice > 0 {
		rows = append(rows, []map[string]any{{"text": "ЮKassa", "callback_data": "buy_yookassa|" + plan.Product}})
	}
	rows = append(rows,
		[]map[string]any{{"text": copy.Back, "callback_data": "menu_premium"}},
		[]map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}},
	)
	return map[string]any{"inline_keyboard": rows}
}

func yookassaPaymentMethodKeyboard(product string, user userState) map[string]any {
	copy := ui(user)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": yooKassaPaymentMethodLabel("", user.InterfaceLanguage), "callback_data": "buy_yookassa_method|" + product + "|any"}},
			{{"text": yooKassaPaymentMethodLabel("bank_card", user.InterfaceLanguage), "callback_data": "buy_yookassa_method|" + product + "|bank_card"}},
			{{"text": yooKassaPaymentMethodLabel("sbp", user.InterfaceLanguage), "callback_data": "buy_yookassa_method|" + product + "|sbp"}},
			{{"text": yooKassaPaymentMethodLabel("yoo_money", user.InterfaceLanguage), "callback_data": "buy_yookassa_method|" + product + "|yoo_money"}},
			{{"text": copy.Back, "callback_data": "premium_plan|" + product}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func yookassaPaymentKeyboard(paymentURL string, user userState) map[string]any {
	copy := ui(user)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{{"text": "ЮKassa", "url": paymentURL}},
			{{"text": copy.BackMenu, "callback_data": "back_menu"}},
		},
	}
}

func cryptoPaymentKeyboard(paymentURL string, paymentID string, user userState) map[string]any {
	copy := ui(user)
	rows := [][]map[string]any{}
	if strings.TrimSpace(paymentURL) != "" {
		rows = append(rows, []map[string]any{{"text": "Open wallet", "url": paymentURL}})
	}
	rows = append(rows,
		[]map[string]any{{"text": "I paid - check", "callback_data": "check_crypto|" + paymentID}},
		[]map[string]any{{"text": copy.BackMenu, "callback_data": "back_menu"}},
	)
	return map[string]any{
		"inline_keyboard": rows,
	}
}

func invoiceInlineKeyboard(stars int, payButtonTemplate string, copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{
				{"text": fmt.Sprintf(payButtonTemplate, stars), "pay": true},
			},
			{
				{"text": copy.BackMenu, "callback_data": "back_menu"},
			},
		},
	}
}

func timezoneInlineKeyboard(copies ...uiCopy) map[string]any {
	copy := keyboardCopy(copies)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{
				{"text": "UTC-8", "callback_data": "tz|-480"},
				{"text": "UTC-5", "callback_data": "tz|-300"},
				{"text": "UTC+0", "callback_data": "tz|0"},
			},
			{
				{"text": "UTC+1", "callback_data": "tz|60"},
				{"text": "UTC+2", "callback_data": "tz|120"},
				{"text": "UTC+3", "callback_data": "tz|180"},
			},
			{
				{"text": "UTC+4", "callback_data": "tz|240"},
				{"text": "UTC+5", "callback_data": "tz|300"},
				{"text": "UTC+6", "callback_data": "tz|360"},
			},
			{
				{"text": "UTC+7", "callback_data": "tz|420"},
				{"text": "UTC+8", "callback_data": "tz|480"},
				{"text": "UTC+9", "callback_data": "tz|540"},
			},
			{
				{"text": copy.BackMenu, "callback_data": "back_menu"},
			},
		},
	}
}

func reminderSettingsInlineKeyboard(user userState) map[string]any {
	copy := ui(user)
	systemCopy := systemUI(user)
	toggleText := "🔔 " + systemCopy.ReminderOnButton
	toggleCallback := "reminder_on"
	if user.ReminderEnabled {
		toggleText = "🔕 " + systemCopy.ReminderOffButton
		toggleCallback = "reminder_off"
	}

	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{
				{"text": toggleText, "callback_data": toggleCallback},
			},
			{
				{"text": "UTC-8", "callback_data": "tz|-480"},
				{"text": "UTC-5", "callback_data": "tz|-300"},
				{"text": "UTC+0", "callback_data": "tz|0"},
			},
			{
				{"text": "UTC+1", "callback_data": "tz|60"},
				{"text": "UTC+2", "callback_data": "tz|120"},
				{"text": "UTC+3", "callback_data": "tz|180"},
			},
			{
				{"text": "UTC+4", "callback_data": "tz|240"},
				{"text": "UTC+5", "callback_data": "tz|300"},
				{"text": "UTC+6", "callback_data": "tz|360"},
			},
			{
				{"text": "UTC+7", "callback_data": "tz|420"},
				{"text": "UTC+8", "callback_data": "tz|480"},
				{"text": "UTC+9", "callback_data": "tz|540"},
			},
			{
				{"text": copy.BackMenu, "callback_data": "back_menu"},
			},
		},
	}
}

func (c *telegramClient) call(ctx context.Context, method string, payload map[string]any) error {
	_, err := c.callMessageID(ctx, method, payload)
	return err
}

func (c *telegramClient) callMessageID(ctx context.Context, method string, payload map[string]any) (int64, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+method, bytes.NewReader(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var decoded struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageID int64 `json:"message_id"`
		} `json:"result"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return 0, err
	}
	if !decoded.OK {
		return 0, fmt.Errorf("telegram %s failed: %s", method, decoded.Description)
	}
	return decoded.Result.MessageID, nil
}

func (c *telegramClient) callMultipart(ctx context.Context, method string, contentType string, body io.Reader) error {
	_, err := c.callMultipartMessageID(ctx, method, contentType, body)
	return err
}

func (c *telegramClient) callMultipartMessageID(ctx context.Context, method string, contentType string, body io.Reader) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+method, body)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := c.http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var decoded struct {
		OK     bool `json:"ok"`
		Result struct {
			MessageID int64 `json:"message_id"`
		} `json:"result"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return 0, err
	}
	if !decoded.OK {
		return 0, fmt.Errorf("telegram %s failed: %s", method, decoded.Description)
	}
	return decoded.Result.MessageID, nil
}
