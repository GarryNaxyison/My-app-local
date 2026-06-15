package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (b *bot) handleAITutorWordReportCallback(ctx context.Context, query callbackQuery, chatID int64) error {
	if b == nil || b.store == nil {
		return errors.New("store is not configured")
	}
	if !telegramOpsUserAllowed(b.cfg, query.From.ID) {
		if b.telegram != nil {
			_ = b.telegram.answerCallbackQuery(ctx, query.ID, "Not allowed")
		}
		return nil
	}
	parts := strings.Split(strings.TrimSpace(query.Data), "|")
	if len(parts) != 3 || parts[0] != "ait_word_report" {
		return nil
	}
	reportID := strings.TrimSpace(parts[1])
	action := strings.TrimSpace(parts[2])
	report, ok, err := b.store.getAITutorWordReport(reportID)
	if err != nil {
		return err
	}
	if !ok {
		if b.telegram != nil {
			return b.telegram.sendMessageToChat(ctx, chatID, "AI Tutor word report not found.")
		}
		return nil
	}
	if report.Status != aiTutorWordReportPending {
		b.closeAITutorWordReportAdminCards(ctx, report)
		return b.sendAITutorWordReportResolvedNotice(ctx, chatID, report)
	}

	switch action {
	case "accept":
		return b.acceptAITutorWordReport(ctx, chatID, report)
	case "reject":
		return b.rejectAITutorWordReport(ctx, chatID, report)
	case "fix":
		return b.promptAITutorWordReportFix(ctx, chatID, report)
	default:
		if b.telegram != nil {
			return b.telegram.sendMessageToChat(ctx, chatID, "Unknown AI Tutor word report action.")
		}
		return nil
	}
}

func (b *bot) acceptAITutorWordReport(ctx context.Context, chatID int64, report aiTutorWordReportRecord) error {
	if report.Status != aiTutorWordReportPending {
		b.closeAITutorWordReportAdminCards(ctx, report)
		return b.sendAITutorWordReportResolvedNotice(ctx, chatID, report)
	}
	finalWord := strings.TrimSpace(report.ProposedWord)
	finalTranslation := strings.TrimSpace(report.ProposedTranslation)
	if err := applyAITutorWordReportCorrection(b.store, report, finalWord, finalTranslation); err != nil {
		return err
	}
	resolved, _, err := b.store.resolveAITutorWordReport(report.ID, aiTutorWordReportAccepted, finalWord, finalTranslation)
	if err != nil {
		return err
	}
	b.closeAITutorWordReportAdminCards(ctx, resolved)
	return b.sendAITutorWordReportResolvedNotice(ctx, chatID, resolved)
}

func (b *bot) rejectAITutorWordReport(ctx context.Context, chatID int64, report aiTutorWordReportRecord) error {
	if report.Status != aiTutorWordReportPending {
		b.closeAITutorWordReportAdminCards(ctx, report)
		return b.sendAITutorWordReportResolvedNotice(ctx, chatID, report)
	}
	resolved, _, err := b.store.resolveAITutorWordReport(report.ID, aiTutorWordReportRejected, "", "")
	if err != nil {
		return err
	}
	b.closeAITutorWordReportAdminCards(ctx, resolved)
	return b.sendAITutorWordReportResolvedNotice(ctx, chatID, resolved)
}

func (b *bot) promptAITutorWordReportFix(ctx context.Context, chatID int64, report aiTutorWordReportRecord) error {
	if report.Status != aiTutorWordReportPending {
		b.closeAITutorWordReportAdminCards(ctx, report)
		return b.sendAITutorWordReportResolvedNotice(ctx, chatID, report)
	}
	if b == nil || b.telegram == nil {
		return nil
	}
	formatHint := b.aiTutorWordReportFixFormat(report)
	text := "Ответьте на это сообщение.\n" + formatHint + "\n\nПосле исправления слово будет заменено в SQLite: AI Tutor lesson payload + vocabulary runtime tables."
	if aiTutorWordReportIsVocabulary(report) {
		text = "Ответьте на это сообщение.\n" + formatHint + "\n\nПосле исправления слово будет заменено в SQLite vocabulary tables: vocabulary_words, vocabulary_translations, vocabulary_ai_words, vocabulary_ai_translations."
	}
	messageID, err := b.telegram.sendInlineMessageToChat(ctx, chatID, text, nil)
	if err != nil {
		return err
	}
	if messageID != 0 {
		if err := b.store.setAITutorWordReportFixPrompt(report.ID, strconv.FormatInt(chatID, 10), messageID); err != nil {
			return err
		}
	}
	return nil
}

func (b *bot) maybeHandleAITutorWordReportFixReply(ctx context.Context, message *telegramMessage) (bool, error) {
	if b == nil || b.store == nil || message == nil {
		return false, nil
	}
	if !telegramOpsUserAllowed(b.cfg, message.From.ID) {
		return false, nil
	}
	text := strings.TrimSpace(message.Text)
	if text == "" {
		return false, nil
	}
	if isStopCommandText(text) || isStopIntent(text) {
		return false, nil
	}
	replyMessageID := int64(0)
	if message.ReplyToMessage != nil {
		replyMessageID = message.ReplyToMessage.MessageID
	}
	report, ok, err := b.store.findPendingAITutorWordReportByFixPrompt(strconv.FormatInt(message.Chat.ID, 10), replyMessageID)
	if err != nil || !ok {
		return ok, err
	}
	finalWord, finalTranslation, ok := parseAITutorWordReportFixText(text)
	if !ok {
		if b.telegram != nil {
			return true, b.telegram.sendMessageToChat(ctx, message.Chat.ID, b.aiTutorWordReportFixFormat(report))
		}
		return true, nil
	}
	if err := applyAITutorWordReportCorrection(b.store, report, finalWord, finalTranslation); err != nil {
		return true, err
	}
	resolved, _, err := b.store.resolveAITutorWordReport(report.ID, aiTutorWordReportFixed, finalWord, finalTranslation)
	if err != nil {
		return true, err
	}
	b.closeAITutorWordReportAdminCards(ctx, resolved)
	return true, b.sendAITutorWordReportResolvedNotice(ctx, message.Chat.ID, resolved)
}

func (b *bot) sendAITutorWordReportResolvedNotice(ctx context.Context, chatID int64, report aiTutorWordReportRecord) error {
	if b == nil || b.telegram == nil {
		return nil
	}
	return b.telegram.sendMessageToChat(ctx, chatID, aiTutorWordReportResolvedText(report))
}

func (b *bot) closeAITutorWordReportAdminCards(ctx context.Context, report aiTutorWordReportRecord) {
	if b == nil || b.store == nil || b.telegram == nil || strings.TrimSpace(report.ID) == "" {
		return
	}
	messages, err := b.store.listAITutorWordReportAdminMessages(report.ID)
	if err != nil || len(messages) == 0 {
		return
	}
	keyboard := map[string]any{"inline_keyboard": []any{}}
	text := aiTutorWordReportResolvedText(report)
	for _, message := range messages {
		chatID := telegramOpsRecipient{ChatID: message.ChatID}.telegramChatIDValue()
		if err := b.telegram.editInlineMessageToChat(ctx, chatID, message.MessageID, text, "", keyboard); err != nil && !isTelegramMessageNotModified(err) {
			continue
		}
	}
}

func aiTutorWordReportResolvedText(report aiTutorWordReportRecord) string {
	label := "AI Tutor word report"
	if aiTutorWordReportIsVocabulary(report) {
		label = "Vocabulary word report"
	}
	text := fmt.Sprintf("%s %s: %s", label, strings.TrimSpace(report.ID), strings.TrimSpace(string(report.Status)))
	if report.FinalWord != "" || report.FinalTranslation != "" {
		text += "\n" + strings.TrimSpace(report.FinalWord) + " - " + strings.TrimSpace(report.FinalTranslation)
	}
	return text
}

func (b *bot) aiTutorWordReportFixFormat(report aiTutorWordReportRecord) string {
	return ui(userState{InterfaceLanguage: b.aiTutorWordReportInterfaceLanguage(report)}).WordReportFixFormat
}

func (b *bot) aiTutorWordReportInterfaceLanguage(report aiTutorWordReportRecord) string {
	if interfaceLanguage, _, ok := parseVocabularyWordReportSessionID(report.SessionID); ok {
		return interfaceLanguage
	}
	if b != nil && b.store != nil && strings.TrimSpace(report.LessonID) != "" {
		lesson, ok, err := b.store.getAITutorLesson(report.LessonID)
		if err == nil && ok {
			return normalizeInterfaceLanguage(firstNonEmpty(lesson.InterfaceLanguage, lesson.Payload.InterfaceLanguage, "ru"))
		}
	}
	return "ru"
}

func telegramOpsUserAllowed(cfg config, telegramID int64) bool {
	if telegramID == 0 {
		return false
	}
	for _, recipient := range cfg.telegramOpsRecipients() {
		chatID := strings.TrimSpace(recipient.ChatID)
		parsed, err := strconv.ParseInt(chatID, 10, 64)
		if err == nil && parsed == telegramID {
			return true
		}
	}
	return false
}

func parseAITutorWordReportFixText(text string) (string, string, bool) {
	text = strings.TrimSpace(text)
	for _, separator := range []string{" - ", " — ", "-"} {
		if !strings.Contains(text, separator) {
			continue
		}
		parts := strings.SplitN(text, separator, 2)
		if len(parts) != 2 {
			continue
		}
		word := strings.TrimSpace(parts[0])
		translation := strings.TrimSpace(parts[1])
		if word == "" || translation == "" || len([]rune(word)) > 80 || len([]rune(translation)) > 160 {
			return "", "", false
		}
		return word, translation, true
	}
	return "", "", false
}

func applyAITutorWordReportCorrection(st store, report aiTutorWordReportRecord, finalWord string, finalTranslation string) error {
	if st == nil {
		return errors.New("store is not configured")
	}
	finalWord = strings.TrimSpace(finalWord)
	finalTranslation = strings.TrimSpace(finalTranslation)
	if finalWord == "" || finalTranslation == "" {
		return errors.New("final word and translation are required")
	}
	if len([]rune(finalWord)) > 80 || len([]rune(finalTranslation)) > 160 {
		return errors.New("final word correction is too long")
	}
	if aiTutorWordReportIsVocabulary(report) {
		return applyVocabularyWordReportCorrection(report, finalWord, finalTranslation)
	}
	lesson, ok, err := st.getAITutorLesson(report.LessonID)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("AI Tutor lesson not found")
	}
	if report.WordIndex < 0 || report.WordIndex >= len(lesson.Payload.Words) {
		return errors.New("AI Tutor word report index is out of range")
	}
	word := lesson.Payload.Words[report.WordIndex]
	word.Target = finalWord
	word.InterfaceTranslation = finalTranslation
	word.AudioTextTarget = finalWord
	lesson.Payload.Words[report.WordIndex] = word
	lesson.Fingerprint = aiTutorFingerprint(lesson.Payload)
	if err := st.saveAITutorLesson(lesson); err != nil {
		return err
	}
	return sqliteVocabularyApplyAITutorWordCorrection(report, lesson, finalWord, finalTranslation)
}

func aiTutorWordReportIsVocabulary(report aiTutorWordReportRecord) bool {
	return strings.TrimSpace(report.Stage) == aiTutorStageVocabularyWordReport
}

func vocabularyWordReportSessionID(user userState, wordID string) string {
	interfaceLanguage := normalizeInterfaceLanguage(firstNonEmpty(user.InterfaceLanguage, "ru"))
	return "vocabulary|" + interfaceLanguage + "|" + strings.TrimSpace(wordID)
}

func parseVocabularyWordReportSessionID(sessionID string) (string, string, bool) {
	parts := strings.SplitN(strings.TrimSpace(sessionID), "|", 3)
	if len(parts) != 3 || parts[0] != "vocabulary" {
		return "", "", false
	}
	interfaceLanguage := normalizeInterfaceLanguage(parts[1])
	wordID := strings.TrimSpace(parts[2])
	return interfaceLanguage, wordID, interfaceLanguage != "" && wordID != ""
}

func applyVocabularyWordReportCorrection(report aiTutorWordReportRecord, finalWord string, finalTranslation string) error {
	interfaceLanguage, wordID, ok := parseVocabularyWordReportSessionID(report.SessionID)
	if !ok {
		interfaceLanguage = normalizeInterfaceLanguage("ru")
		wordID = strings.TrimSpace(report.LessonID)
	}
	if strings.TrimSpace(report.LessonID) != "" {
		wordID = strings.TrimSpace(report.LessonID)
	}
	if wordID == "" {
		return errors.New("vocabulary word report id is empty")
	}
	word, ok := findVocabWord(wordID)
	if !ok {
		return errors.New("vocabulary word not found")
	}
	return sqliteVocabularyApplyVocabularyWordCorrection(report, word, interfaceLanguage, finalWord, finalTranslation)
}
