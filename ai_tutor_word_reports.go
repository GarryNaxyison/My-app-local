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
	return b.sendAITutorWordReportResolvedNotice(ctx, chatID, resolved)
}

func (b *bot) rejectAITutorWordReport(ctx context.Context, chatID int64, report aiTutorWordReportRecord) error {
	if report.Status != aiTutorWordReportPending {
		return b.sendAITutorWordReportResolvedNotice(ctx, chatID, report)
	}
	resolved, _, err := b.store.resolveAITutorWordReport(report.ID, aiTutorWordReportRejected, "", "")
	if err != nil {
		return err
	}
	return b.sendAITutorWordReportResolvedNotice(ctx, chatID, resolved)
}

func (b *bot) promptAITutorWordReportFix(ctx context.Context, chatID int64, report aiTutorWordReportRecord) error {
	if report.Status != aiTutorWordReportPending {
		return b.sendAITutorWordReportResolvedNotice(ctx, chatID, report)
	}
	if b == nil || b.telegram == nil {
		return nil
	}
	text := "Ответьте на это сообщение в формате: word - translation"
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
	if b == nil || b.store == nil || message == nil || message.ReplyToMessage == nil {
		return false, nil
	}
	if !telegramOpsUserAllowed(b.cfg, message.From.ID) {
		return false, nil
	}
	text := strings.TrimSpace(message.Text)
	if text == "" {
		return false, nil
	}
	report, ok, err := b.store.findPendingAITutorWordReportByFixPrompt(strconv.FormatInt(message.Chat.ID, 10), message.ReplyToMessage.MessageID)
	if err != nil || !ok {
		return ok, err
	}
	finalWord, finalTranslation, ok := parseAITutorWordReportFixText(text)
	if !ok {
		if b.telegram != nil {
			return true, b.telegram.sendMessageToChat(ctx, message.Chat.ID, "Формат: word - translation")
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
	return true, b.sendAITutorWordReportResolvedNotice(ctx, message.Chat.ID, resolved)
}

func (b *bot) sendAITutorWordReportResolvedNotice(ctx context.Context, chatID int64, report aiTutorWordReportRecord) error {
	if b == nil || b.telegram == nil {
		return nil
	}
	text := fmt.Sprintf("AI Tutor word report %s: %s", strings.TrimSpace(report.ID), strings.TrimSpace(string(report.Status)))
	if report.FinalWord != "" || report.FinalTranslation != "" {
		text += "\n" + strings.TrimSpace(report.FinalWord) + " - " + strings.TrimSpace(report.FinalTranslation)
	}
	return b.telegram.sendMessageToChat(ctx, chatID, text)
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
