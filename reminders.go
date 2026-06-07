package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (b *bot) runReminderScheduler(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	b.sendDueReminders(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.sendDueReminders(ctx)
		}
	}
}

func (b *bot) sendDueReminders(ctx context.Context) {
	targets, err := b.store.dueReminderUsers(time.Now().UTC(), 100)
	if err != nil {
		log.Printf("reminder scheduler failed: %v", err)
		return
	}
	if len(targets) > 0 {
		log.Printf("reminder scheduler: due targets=%d", len(targets))
	}
	for _, target := range targets {
		user := userState{TelegramID: target.TelegramID, FirstName: target.FirstName, InterfaceLanguage: target.InterfaceLanguage}
		text := reminderTextForDate(target.LocalDate, target.InterfaceLanguage)
		if err := b.telegram.sendMessageWithCopy(ctx, target.TelegramID, text, ui(user)); err != nil {
			log.Printf("failed to send reminder to %d: %v", target.TelegramID, err)
			continue
		}
		if err := b.store.markReminderSent(target.TelegramID, target.LocalDate); err != nil {
			log.Printf("failed to mark reminder for %d: %v", target.TelegramID, err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(200 * time.Millisecond):
		}
	}
}

func reminderDue(user userState, now time.Time) (reminderTarget, bool) {
	if !user.ReminderEnabled {
		return reminderTarget{}, false
	}
	localNow := now.UTC().Add(time.Duration(user.ReminderUTCOffset) * time.Minute)
	if localNow.Hour() < user.ReminderHour {
		return reminderTarget{}, false
	}
	localDate := localNow.Format("2006-01-02")
	if user.LastReminderDate == localDate {
		return reminderTarget{}, false
	}
	return reminderTarget{TelegramID: user.TelegramID, FirstName: user.FirstName, InterfaceLanguage: user.InterfaceLanguage, LocalDate: localDate}, true
}

func reminderTextForDate(localDate string, interfaceLanguage string) string {
	user := userState{InterfaceLanguage: interfaceLanguage}
	copy := systemUI(user)
	messages := copy.ReminderMessages
	if len(messages) == 0 {
		messages = systemUICopies["en"].ReminderMessages
	}
	parsed, err := time.Parse("2006-01-02", localDate)
	if err != nil {
		return messages[0] + "\n\n" + fmt.Sprintf(copy.ReminderFooter, cleanReminderMenuButton(interfaceLanguage))
	}
	index := int(parsed.Unix()/86400) % len(messages)
	if index < 0 {
		index = 0
	}
	return messages[index] + "\n\n" + fmt.Sprintf(copy.ReminderFooter, cleanReminderMenuButton(interfaceLanguage))
}
