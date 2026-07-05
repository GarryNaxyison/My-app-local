package main

import (
	"path/filepath"
	"testing"
	"time"
)

func TestSQLiteSimpleUserUpdatesDoNotRewriteCollections(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.db.Close() })

	if _, err := store.getOrCreateUser(123, "tester"); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.addLearnedWord(123, vocabWord{ID: "en:hello", Language: "en", English: "hello", Russian: "привет"}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(`CREATE TABLE audit_deletes (table_name TEXT NOT NULL)`); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(`CREATE TRIGGER audit_learned_word_deletes AFTER DELETE ON learned_words BEGIN INSERT INTO audit_deletes (table_name) VALUES ('learned_words'); END`); err != nil {
		t.Fatal(err)
	}

	for name, action := range map[string]func() error{
		"getOrCreateUser":      func() error { _, err := store.getOrCreateUser(123, "tester"); return err },
		"setMode":              func() error { return store.setMode(123, "practice") },
		"setInterfaceLanguage": func() error { return store.setInterfaceLanguage(123, "es") },
		"setLearningLanguage":  func() error { return store.setLearningLanguage(123, "de") },
		"setUserLevel":         func() error { return store.setUserLevel(123, "B1") },
		"setLearningFocus":     func() error { return store.setLearningFocus(123, "travel") },
		"setNavigationLayout": func() error {
			return store.setNavigationLayout(123, navigationLayout{MobilePinned: []string{"home", "tools"}})
		},
		"setReminderEnabled": func() error { return store.setReminderEnabled(123, false) },
		"setTimezone":        func() error { return store.setTimezone(123, 240) },
		"markReminderSent":   func() error { return store.markReminderSent(123, "2026-06-15") },
	} {
		if _, err := store.db.Exec(`DELETE FROM audit_deletes`); err != nil {
			t.Fatal(err)
		}
		if err := action(); err != nil {
			t.Fatalf("%s() error = %v", name, err)
		}
		var deletes int
		if err := store.db.QueryRow(`SELECT COUNT(*) FROM audit_deletes`).Scan(&deletes); err != nil {
			t.Fatal(err)
		}
		if deletes != 0 {
			t.Fatalf("%s rewrote learned_words with %d delete(s)", name, deletes)
		}
	}

	if err := store.setMode(123, "practice"); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.Exec(`DELETE FROM audit_deletes`); err != nil {
		t.Fatal(err)
	}
	if err := store.setMode(123, "practice"); err != nil {
		t.Fatal(err)
	}
	var deletes int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM audit_deletes`).Scan(&deletes); err != nil {
		t.Fatal(err)
	}
	if deletes != 0 {
		t.Fatalf("setMode with unchanged value rewrote learned_words with %d delete(s)", deletes)
	}
}

func TestSQLiteStoreUsesSingleConnectionPool(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.db.Close() })

	if got := store.db.Stats().MaxOpenConnections; got != 1 {
		t.Fatalf("sqlite max open connections = %d, want 1", got)
	}
}

func TestSQLiteResetLearningStatePreservesAccountAndCommercialData(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.db.Close() })

	const userID int64 = 123
	hash, err := hashWebPassword("strong-password")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.createWebAccount(userID, "tester", hash); err != nil {
		t.Fatal(err)
	}
	user, err := store.getOrCreateUser(userID, "tester")
	if err != nil {
		t.Fatal(err)
	}
	user.Plan = "premium"
	user.PremiumUntil = time.Now().UTC().Add(30 * 24 * time.Hour)
	user.InvitedBy = 77
	user.ReferralCode = "KEEP123"
	user.ReferralCount = 4
	user.ReferralBalanceKopecks = 12345
	user.ReferralLevelRewarded = true
	user.LastPaymentChargeID = "charge-keep"
	user.Mode = "lesson"
	user.Level = "B2"
	user.XP = 500
	user.LessonCount = 9
	user.PracticeCount = 8
	user.VoiceCount = 7
	user.WordLessonCount = 6
	user.WordGameCount = 5
	user.DailyDate = "2026-07-05"
	user.LessonsToday = 3
	user.PracticeToday = 2
	user.VoiceToday = 1
	user.LessonHistory = []string{"old lesson"}
	user.PracticeHistory = []string{"old practice"}
	user.LastLessonPrompt = "old prompt"
	if err := store.saveUser(user); err != nil {
		t.Fatal(err)
	}
	if _, _, err := store.addLearnedWord(userID, vocabWord{ID: "en:hello", Language: "en", English: "hello", Russian: "привет"}); err != nil {
		t.Fatal(err)
	}
	if err := store.addMistakes(userID, "en", []mistakeEntry{{Word: "helo", Language: "en", Correction: "hello", AddedAt: time.Now().UTC()}}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.addPhrasebookEntry(userID, phrasebookEntry{Phrase: "Hello there", Translation: "Привет", Source: "manual", Language: "en"}); err != nil {
		t.Fatal(err)
	}
	now := formatDBTime(time.Now().UTC())
	for _, stmt := range []string{
		`INSERT INTO daily_bonus_claims (telegram_id, claim_date, xp, created_at) VALUES (123, '2026-07-05', 20, '` + now + `')`,
		`INSERT INTO habit_days (telegram_id, habit_date, login, complete, claimed, claimed_at, lessons, practice, voice, updated_at) VALUES (123, '2026-07-05', 1, 1, 1, '` + now + `', 1, 1, 1, '` + now + `')`,
		`INSERT INTO tutor_lessons (id, learning_language, interface_language, level, lesson_number, topic, payload_json, created_at, updated_at) VALUES ('lesson-keep', 'en', 'ru', 'A1', 1, 'topic', '{}', '` + now + `', '` + now + `')`,
		`INSERT INTO tutor_user_lessons (telegram_id, lesson_id, assigned_at, completed_at, lesson_order) VALUES (123, 'lesson-keep', '` + now + `', '` + now + `', 1)`,
		`INSERT INTO ai_tutor_lessons (id, learning_language, interface_language, exact_level, level_band, theme, status, payload_json, fingerprint, created_at, updated_at) VALUES ('ai-lesson-keep', 'en', 'ru', 'A1', 'A1', 'topic', 'approved', '{}', 'fp', '` + now + `', '` + now + `')`,
		`INSERT INTO ai_tutor_sessions (id, telegram_id, lesson_id, surface, current_stage, status, started_at, updated_at) VALUES ('session-wipe', 123, 'ai-lesson-keep', 'web', 'story_intro', 'active', '` + now + `', '` + now + `')`,
		`INSERT INTO ai_tutor_answers (session_id, stage, answer_text, created_at) VALUES ('session-wipe', 'story_intro', 'answer', '` + now + `')`,
		`INSERT INTO ai_tutor_reviews (id, telegram_id, lesson_id, due_at, interval_code, status, created_at, updated_at) VALUES ('review-wipe', 123, 'ai-lesson-keep', '` + now + `', 'd1', 'scheduled', '` + now + `', '` + now + `')`,
		`INSERT INTO ai_tutor_word_reports (id, telegram_id, session_id, lesson_id, stage, word_index, original_word, original_translation, proposed_word, proposed_translation, comment, status, created_at, updated_at) VALUES ('report-wipe', 123, 'session-wipe', 'ai-lesson-keep', 'stage', 0, 'bad', 'bad', 'good', 'good', '', 'pending', '` + now + `', '` + now + `')`,
		`INSERT INTO crypto_payments (id, provider, currency, network, product, status, address, memo, amount, amount_nano, price_rub, telegram_id, expires_at, created_at, updated_at) VALUES ('payment-keep', 'test', 'TON', 'ton', 'premium', 'paid', 'addr', 'memo', '1', 1, 100, 123, '` + now + `', '` + now + `', '` + now + `')`,
		`INSERT INTO referral_purchase_rewards (payment_id, buyer_id, direct_inviter_id, paid_kopecks, direct_reward_kopecks, created_at) VALUES ('reward-keep', 123, 77, 10000, 1000, '` + now + `')`,
	} {
		if _, err := store.db.Exec(stmt); err != nil {
			t.Fatalf("seed statement failed: %v\n%s", err, stmt)
		}
	}
	beforeReset, ok, err := store.getUser(userID)
	if err != nil || !ok {
		t.Fatalf("getUser before reset ok=%v err=%v", ok, err)
	}

	if err := store.resetLearningState(userID); err != nil {
		t.Fatalf("resetLearningState() error = %v", err)
	}

	for _, table := range []string{"mistakes", "learned_words", "phrasebook_entries", "daily_bonus_claims", "habit_days", "tutor_user_lessons", "ai_tutor_sessions", "ai_tutor_reviews", "ai_tutor_word_reports"} {
		var count int
		if err := store.db.QueryRow(`SELECT COUNT(*) FROM `+table+` WHERE telegram_id = ?`, userID).Scan(&count); err != nil {
			t.Fatalf("count %s: %v", table, err)
		}
		if count != 0 {
			t.Fatalf("%s rows after reset = %d, want 0", table, count)
		}
	}
	var answerRows int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM ai_tutor_answers WHERE session_id = 'session-wipe'`).Scan(&answerRows); err != nil {
		t.Fatal(err)
	}
	if answerRows != 0 {
		t.Fatalf("ai_tutor_answers rows after reset = %d, want 0", answerRows)
	}
	refreshed, ok, err := store.getUser(userID)
	if err != nil || !ok {
		t.Fatalf("getUser after reset ok=%v err=%v", ok, err)
	}
	if refreshed.Plan != beforeReset.Plan ||
		!refreshed.PremiumUntil.Equal(beforeReset.PremiumUntil) ||
		refreshed.InvitedBy != beforeReset.InvitedBy ||
		refreshed.ReferralCode != beforeReset.ReferralCode ||
		refreshed.ReferralCount != beforeReset.ReferralCount ||
		refreshed.ReferralBalanceKopecks != beforeReset.ReferralBalanceKopecks ||
		refreshed.ReferralLevelRewarded != beforeReset.ReferralLevelRewarded ||
		refreshed.LastPaymentChargeID != beforeReset.LastPaymentChargeID {
		t.Fatalf("commercial fields were not preserved: %#v", refreshed)
	}
	if refreshed.XP != 0 || refreshed.Mode != "idle" || refreshed.LessonCount != 0 || refreshed.PracticeCount != 0 || refreshed.VoiceCount != 0 || refreshed.WordLessonCount != 0 || refreshed.WordGameCount != 0 || refreshed.LessonsToday != 0 || refreshed.PracticeToday != 0 || refreshed.VoiceToday != 0 || refreshed.LastLessonPrompt != "" || len(refreshed.LessonHistory) != 0 || len(refreshed.PracticeHistory) != 0 {
		t.Fatalf("learning fields were not reset: %#v", refreshed)
	}
	var webAccounts, payments, rewards int
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM web_accounts WHERE user_id = ?`, userID).Scan(&webAccounts); err != nil {
		t.Fatal(err)
	}
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM crypto_payments WHERE telegram_id = ?`, userID).Scan(&payments); err != nil {
		t.Fatal(err)
	}
	if err := store.db.QueryRow(`SELECT COUNT(*) FROM referral_purchase_rewards WHERE buyer_id = ?`, userID).Scan(&rewards); err != nil {
		t.Fatal(err)
	}
	if webAccounts != 1 || payments != 1 || rewards != 1 {
		t.Fatalf("preserved rows web=%d payments=%d rewards=%d, want 1 each", webAccounts, payments, rewards)
	}
}
