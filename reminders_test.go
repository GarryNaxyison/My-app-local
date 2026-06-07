package main

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSQLiteDueReminderUsersIncludesSelectedTimezoneUser(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()

	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	user := newUserState(12345, now)
	user.FirstName = "Tester"
	user.InterfaceLanguage = "es"
	user.ReminderEnabled = true
	user.TimezoneSelected = true
	user.ReminderUTCOffset = 180
	user.ReminderHour = 19
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := store.saveUser(user); err != nil {
		t.Fatal(err)
	}

	targets, err := store.dueReminderUsers(time.Date(2026, 5, 15, 16, 0, 0, 0, time.UTC), 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected 1 due reminder target, got %d", len(targets))
	}
	if targets[0].TelegramID != user.TelegramID {
		t.Fatalf("expected telegram id %d, got %d", user.TelegramID, targets[0].TelegramID)
	}
	if targets[0].LocalDate != "2026-05-15" {
		t.Fatalf("expected local date 2026-05-15, got %q", targets[0].LocalDate)
	}
	if targets[0].InterfaceLanguage != "es" {
		t.Fatalf("expected interface language es, got %q", targets[0].InterfaceLanguage)
	}
}

func TestSQLiteDueReminderUsersUsesMoscowTimeByDefault(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()

	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	user := newUserState(67890, now)
	user.FirstName = "Default TZ"
	user.ReminderEnabled = true
	user.TimezoneSelected = false
	user.ReminderUTCOffset = 180
	user.ReminderHour = 19
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := store.saveUser(user); err != nil {
		t.Fatal(err)
	}

	targets, err := store.dueReminderUsers(time.Date(2026, 5, 15, 16, 0, 0, 0, time.UTC), 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected 1 due reminder target, got %d", len(targets))
	}
	if targets[0].TelegramID != user.TelegramID {
		t.Fatalf("expected telegram id %d, got %d", user.TelegramID, targets[0].TelegramID)
	}
}

func TestSQLiteSetTimezoneKeepsReminderDisabled(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()

	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	user := newUserState(24680, now)
	user.FirstName = "No reminders"
	user.ReminderEnabled = false
	user.TimezoneSelected = false
	user.ReminderUTCOffset = 180
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := store.saveUser(user); err != nil {
		t.Fatal(err)
	}
	if err := store.setTimezone(user.TelegramID, 300); err != nil {
		t.Fatal(err)
	}

	updated, ok, err := store.getUser(user.TelegramID)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected user to exist")
	}
	if updated.ReminderEnabled {
		t.Fatal("expected reminders to stay disabled after timezone change")
	}
	if !updated.TimezoneSelected {
		t.Fatal("expected timezone to be selected")
	}
	if updated.ReminderUTCOffset != 300 {
		t.Fatalf("expected UTC offset 300, got %d", updated.ReminderUTCOffset)
	}
}

func TestReminderTextDoesNotExposeCycleDay(t *testing.T) {
	text := reminderTextForDate("2026-05-15", "ru")
	if strings.Contains(text, "День ") || strings.Contains(text, "Day ") {
		t.Fatalf("reminder must not expose cycle day: %q", text)
	}
}

func TestReminderTextUsesInterfaceLanguage(t *testing.T) {
	text := reminderTextForDate("2026-05-15", "es")
	if strings.Contains(text, "Открой") || strings.Contains(text, "Напоминания") {
		t.Fatalf("expected non-Russian reminder for Spanish UI, got %q", text)
	}
	if !strings.Contains(text, "Abre") {
		t.Fatalf("expected Spanish reminder footer, got %q", text)
	}
}

func TestSQLiteDueReminderUsersTracksChangedInterfaceLanguage(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	defer store.db.Close()

	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	user := newUserState(54321, now)
	user.FirstName = "Changed UI"
	user.InterfaceLanguage = "ru"
	user.ReminderEnabled = true
	user.TimezoneSelected = true
	user.ReminderUTCOffset = 180
	user.ReminderHour = 19
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := store.saveUser(user); err != nil {
		t.Fatal(err)
	}
	if err := store.setInterfaceLanguage(user.TelegramID, "es"); err != nil {
		t.Fatal(err)
	}

	targets, err := store.dueReminderUsers(time.Date(2026, 5, 15, 16, 0, 0, 0, time.UTC), 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected 1 due reminder target, got %d", len(targets))
	}
	if targets[0].InterfaceLanguage != "es" {
		t.Fatalf("expected changed interface language es, got %q", targets[0].InterfaceLanguage)
	}

	text := reminderTextForDate(targets[0].LocalDate, targets[0].InterfaceLanguage)
	if !strings.Contains(text, "Abre") {
		t.Fatalf("expected reminder text to use changed Spanish UI language, got %q", text)
	}
}

func TestReminderTextUsesLocalizedCompactFallback(t *testing.T) {
	language := "pt"
	text := reminderTextForDate("2026-05-15", language)

	english := ui(userState{InterfaceLanguage: "en"})
	if strings.Contains(text, english.Notifications+": "+english.NewLesson+" / "+english.Practice) {
		t.Fatalf("expected localized fallback instead of English fallback, got %q", text)
	}
	if !strings.Contains(text, "Abra Menu") || !strings.Contains(text, "lembretes") {
		t.Fatalf("expected Portuguese reminder text and footer, got %q", text)
	}
}

func TestReminderTextIsCleanAndLocalizedForEveryInterfaceLanguage(t *testing.T) {
	for _, language := range interfaceLanguages() {
		text := reminderTextForDate("2026-05-15", language.Code)
		if hasTestEncodingArtifact(text) {
			t.Fatalf("reminder for %s contains encoding artifact: %q", language.Code, text)
		}
		if language.Code != "en" && strings.Contains(text, "A small language step today") {
			t.Fatalf("reminder for %s fell back to English text: %q", language.Code, text)
		}
		if language.Code != "en" && strings.Contains(text, "Open Menu") {
			t.Fatalf("reminder for %s fell back to English footer: %q", language.Code, text)
		}
	}
}

func TestSystemUIReminderSettingsAreCleanForEveryInterfaceLanguage(t *testing.T) {
	for _, language := range interfaceLanguages() {
		copy := systemUI(userState{InterfaceLanguage: language.Code})
		values := []string{
			copy.ReminderTitle,
			copy.ReminderStatus,
			copy.ReminderEnabled,
			copy.ReminderDisabled,
			copy.ReminderTimezone,
			copy.ReminderDefault,
			copy.ReminderHour,
			copy.ReminderHint,
			copy.ReminderOnButton,
			copy.ReminderOffButton,
			copy.ReminderFooter,
		}
		values = append(values, copy.ReminderMessages...)
		for _, value := range values {
			if value == "" {
				t.Fatalf("empty reminder copy for %s", language.Code)
			}
			if hasTestEncodingArtifact(value) {
				t.Fatalf("reminder copy for %s contains encoding artifact: %q", language.Code, value)
			}
		}
	}
}

func hasTestEncodingArtifact(text string) bool {
	artifacts := []string{
		"Рђ", "Р‘", "Р’", "Р“", "Р”", "Р•", "Р–", "Р—", "Р", "Р™", "Рљ", "Р›", "Рњ", "Рќ", "Рћ", "Рџ",
		"Р°", "Р±", "Рі", "Рґ", "Рµ", "Р¶", "Р·", "Рё", "Р№", "Рє", "Р»", "Рј", "РЅ", "Рѕ", "Рї",
		"С€", "С‰", "СЊ", "С‹", "СЏ", "СЂ", "СЃ", "С‚", "Сѓ", "С„", "С…", "С†", "С‡", "СЋ",
		"вЂ", "вЏ", "в—", "в–", "в‚", "Г", "Ð", "Ñ", "Â", "Ã", "�", "???",
	}
	for _, artifact := range artifacts {
		if strings.Contains(text, artifact) {
			return true
		}
	}
	return false
}
