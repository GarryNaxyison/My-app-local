package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type sqliteStore struct {
	db *sql.DB
}

func openSQLiteDatabase(name string, databasePath string) (*sql.DB, error) {
	if err := requireNonEmpty(name, databasePath); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)
	return db, nil
}

func newSQLiteStore(databasePath string, importJSONPath string) (*sqliteStore, error) {
	db, err := openSQLiteDatabase("DATABASE_PATH", databasePath)
	if err != nil {
		return nil, err
	}
	store := &sqliteStore{db: db}
	if err := store.init(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.importJSONIfEmpty(importJSONPath); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *sqliteStore) init() error {
	statements := []string{
		`PRAGMA journal_mode=WAL`,
		`PRAGMA synchronous=NORMAL`,
		`PRAGMA busy_timeout=5000`,
		`PRAGMA temp_store=MEMORY`,
		`PRAGMA foreign_keys=ON`,
		`CREATE TABLE IF NOT EXISTS users (
			telegram_id INTEGER PRIMARY KEY,
			first_name TEXT NOT NULL DEFAULT '',
			plan TEXT NOT NULL DEFAULT 'free',
			premium_until TEXT,
			invited_by INTEGER NOT NULL DEFAULT 0,
			referral_code TEXT NOT NULL DEFAULT '',
			referral_count INTEGER NOT NULL DEFAULT 0,
			referral_balance_kopecks INTEGER NOT NULL DEFAULT 0,
			referral_level_rewarded INTEGER NOT NULL DEFAULT 0,
			mode TEXT NOT NULL DEFAULT 'idle',
			interface_language TEXT NOT NULL DEFAULT 'ru',
			interface_selected INTEGER NOT NULL DEFAULT 0,
			learning_language TEXT NOT NULL DEFAULT 'en',
			language_selected INTEGER NOT NULL DEFAULT 0,
			level TEXT NOT NULL DEFAULT 'A2',
			learning_focus TEXT NOT NULL DEFAULT '',
			lesson_count INTEGER NOT NULL DEFAULT 0,
			practice_count INTEGER NOT NULL DEFAULT 0,
			voice_count INTEGER NOT NULL DEFAULT 0,
			word_lesson_count INTEGER NOT NULL DEFAULT 0,
			word_game_count INTEGER NOT NULL DEFAULT 0,
			xp INTEGER NOT NULL DEFAULT 0,
			daily_date TEXT NOT NULL DEFAULT '',
			lessons_today INTEGER NOT NULL DEFAULT 0,
			practice_today INTEGER NOT NULL DEFAULT 0,
			voice_today INTEGER NOT NULL DEFAULT 0,
			lesson_history TEXT NOT NULL DEFAULT '',
			practice_history TEXT NOT NULL DEFAULT '',
			last_payment_charge_id TEXT NOT NULL DEFAULT '',
			last_lesson_prompt TEXT NOT NULL DEFAULT '',
			reminder_enabled INTEGER NOT NULL DEFAULT 1,
			reminder_utc_offset INTEGER NOT NULL DEFAULT 180,
			reminder_hour INTEGER NOT NULL DEFAULT 19,
			timezone_selected INTEGER NOT NULL DEFAULT 0,
			last_reminder_date TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS mistakes (
			telegram_id INTEGER NOT NULL,
			word TEXT NOT NULL,
			language TEXT NOT NULL DEFAULT 'en',
			correction TEXT NOT NULL,
			explanation TEXT NOT NULL DEFAULT '',
			added_at TEXT NOT NULL,
			PRIMARY KEY (telegram_id, word, correction)
		)`,
		`CREATE TABLE IF NOT EXISTS learned_words (
			telegram_id INTEGER NOT NULL,
			id TEXT NOT NULL,
			language TEXT NOT NULL DEFAULT 'en',
			russian TEXT NOT NULL,
			english TEXT NOT NULL,
			context TEXT NOT NULL DEFAULT '',
			review_correct_count INTEGER NOT NULL DEFAULT 0,
			spelling_correct_count INTEGER NOT NULL DEFAULT 0,
			learned_at TEXT NOT NULL,
			PRIMARY KEY (telegram_id, id)
		)`,
		`CREATE TABLE IF NOT EXISTS phrasebook_entries (
			telegram_id INTEGER NOT NULL,
			id TEXT NOT NULL,
			phrase TEXT NOT NULL,
			translation TEXT NOT NULL DEFAULT '',
			note TEXT NOT NULL DEFAULT '',
			source TEXT NOT NULL DEFAULT 'manual',
			language TEXT NOT NULL DEFAULT 'en',
			created_at TEXT NOT NULL,
			PRIMARY KEY (telegram_id, id)
		)`,
		`CREATE TABLE IF NOT EXISTS daily_bonus_claims (
			telegram_id INTEGER NOT NULL,
			claim_date TEXT NOT NULL,
			xp INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			PRIMARY KEY (telegram_id, claim_date)
		)`,
		`CREATE TABLE IF NOT EXISTS habit_days (
			telegram_id INTEGER NOT NULL,
			habit_date TEXT NOT NULL,
			login INTEGER NOT NULL DEFAULT 0,
			complete INTEGER NOT NULL DEFAULT 0,
			claimed INTEGER NOT NULL DEFAULT 0,
			claimed_at TEXT NOT NULL DEFAULT '',
			lessons INTEGER NOT NULL DEFAULT 0,
			practice INTEGER NOT NULL DEFAULT 0,
			voice INTEGER NOT NULL DEFAULT 0,
			updated_at TEXT NOT NULL,
			PRIMARY KEY (telegram_id, habit_date)
		)`,
		`CREATE TABLE IF NOT EXISTS navigation_layouts (
			telegram_id INTEGER PRIMARY KEY,
			function_ribbon TEXT NOT NULL DEFAULT '',
			mobile_pinned TEXT NOT NULL DEFAULT '',
			mobile_more TEXT NOT NULL DEFAULT '',
			mobile_rail TEXT NOT NULL DEFAULT '',
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tutor_lessons (
			id TEXT PRIMARY KEY,
			learning_language TEXT NOT NULL,
			interface_language TEXT NOT NULL,
			level TEXT NOT NULL,
			lesson_number INTEGER NOT NULL,
			topic TEXT NOT NULL DEFAULT '',
			payload_json TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tutor_user_lessons (
			telegram_id INTEGER NOT NULL,
			lesson_id TEXT NOT NULL,
			assigned_at TEXT NOT NULL,
			completed_at TEXT NOT NULL DEFAULT '',
			lesson_order INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (telegram_id, lesson_id)
		)`,
		`CREATE TABLE IF NOT EXISTS referral_purchase_rewards (
			payment_id TEXT PRIMARY KEY,
			buyer_id INTEGER NOT NULL,
			direct_inviter_id INTEGER NOT NULL DEFAULT 0,
			indirect_inviter_id INTEGER NOT NULL DEFAULT 0,
			paid_kopecks INTEGER NOT NULL DEFAULT 0,
			direct_reward_kopecks INTEGER NOT NULL DEFAULT 0,
			indirect_reward_kopecks INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS crypto_payments (
			id TEXT PRIMARY KEY,
			provider TEXT NOT NULL,
			currency TEXT NOT NULL,
			network TEXT NOT NULL,
			product TEXT NOT NULL,
			status TEXT NOT NULL,
			address TEXT NOT NULL,
			memo TEXT NOT NULL,
			amount TEXT NOT NULL,
			amount_nano INTEGER NOT NULL,
			price_rub INTEGER NOT NULL DEFAULT 0,
			telegram_id INTEGER NOT NULL,
			tx_hash TEXT NOT NULL DEFAULT '',
			expires_at TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS web_accounts (
			user_id INTEGER PRIMARY KEY,
			login TEXT NOT NULL,
			login_normalized TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			last_login_at TEXT NOT NULL DEFAULT ''
		)`,
	}
	for _, stmt := range statements {
		if _, err := s.db.Exec(stmt); err != nil {
			return err
		}
	}
	migrations := []string{
		`ALTER TABLE users ADD COLUMN reminder_enabled INTEGER NOT NULL DEFAULT 1`,
		`ALTER TABLE users ADD COLUMN reminder_utc_offset INTEGER NOT NULL DEFAULT 180`,
		`ALTER TABLE users ADD COLUMN reminder_hour INTEGER NOT NULL DEFAULT 19`,
		`ALTER TABLE users ADD COLUMN timezone_selected INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN last_reminder_date TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN lesson_history TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN practice_history TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN interface_language TEXT NOT NULL DEFAULT 'ru'`,
		`ALTER TABLE users ADD COLUMN interface_selected INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN learning_language TEXT NOT NULL DEFAULT 'en'`,
		`ALTER TABLE users ADD COLUMN language_selected INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN learning_focus TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE learned_words ADD COLUMN language TEXT NOT NULL DEFAULT 'en'`,
		`ALTER TABLE learned_words ADD COLUMN context TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE learned_words ADD COLUMN review_correct_count INTEGER NOT NULL DEFAULT 10`,
		`ALTER TABLE learned_words ADD COLUMN spelling_correct_count INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE mistakes ADD COLUMN language TEXT NOT NULL DEFAULT 'en'`,
		`ALTER TABLE users ADD COLUMN referral_code TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN referral_balance_kopecks INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN referral_level_rewarded INTEGER NOT NULL DEFAULT 0`,
	}
	for _, stmt := range migrations {
		if _, err := s.db.Exec(stmt); err != nil && !strings.Contains(err.Error(), "duplicate column name") {
			return err
		}
	}
	if _, err := s.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_referral_code ON users(referral_code) WHERE referral_code <> ''`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_crypto_payments_user_created ON crypto_payments(telegram_id, created_at)`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE UNIQUE INDEX IF NOT EXISTS idx_crypto_payments_tx_hash ON crypto_payments(tx_hash) WHERE tx_hash <> ''`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_daily_bonus_claims_user_created ON daily_bonus_claims(telegram_id, created_at DESC)`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_habit_days_user_date ON habit_days(telegram_id, habit_date DESC)`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_learned_words_user_language_review ON learned_words(telegram_id, language, review_correct_count, spelling_correct_count, learned_at)`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_learned_words_user_language_learned ON learned_words(telegram_id, language, learned_at)`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_mistakes_user_language_added ON mistakes(telegram_id, language, added_at DESC)`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_tutor_lessons_context ON tutor_lessons(learning_language, interface_language, level, lesson_number, id)`); err != nil {
		return err
	}
	if _, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_tutor_user_lessons_user_order ON tutor_user_lessons(telegram_id, lesson_order, assigned_at)`); err != nil {
		return err
	}
	return nil
}

func (s *sqliteStore) createWebAccount(userID int64, login string, passwordHash string) (webAccount, error) {
	login, err := normalizeWebLogin(login)
	if err != nil {
		return webAccount{}, err
	}
	if strings.TrimSpace(passwordHash) == "" {
		return webAccount{}, errors.New("password hash must not be empty")
	}
	if userID == 0 {
		return webAccount{}, errors.New("web account user id must not be empty")
	}
	if existing, _, ok, err := s.getWebAccountByLogin(login); err != nil {
		return webAccount{}, err
	} else if ok {
		_ = existing
		return webAccount{}, errWebLoginTaken
	}
	now := time.Now().UTC()
	_, err = s.db.Exec(`INSERT INTO web_accounts (user_id, login, login_normalized, password_hash, created_at, updated_at, last_login_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userID, login, login, passwordHash, formatDBTime(now), formatDBTime(now), formatDBTime(now))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "constraint") {
			return webAccount{}, errWebLoginTaken
		}
		return webAccount{}, err
	}
	account, ok, err := s.getWebAccountByUserID(userID)
	if err != nil {
		return webAccount{}, err
	}
	if !ok {
		return webAccount{}, errors.New("web account not found after create")
	}
	return account, nil
}

func (s *sqliteStore) getWebAccountByLogin(login string) (webAccount, string, bool, error) {
	login, err := normalizeWebLogin(login)
	if err != nil {
		return webAccount{}, "", false, err
	}
	return s.getWebAccount(`WHERE login_normalized = ?`, login)
}

func (s *sqliteStore) getWebAccountByUserID(userID int64) (webAccount, bool, error) {
	if userID == 0 {
		return webAccount{}, false, nil
	}
	account, _, ok, err := s.getWebAccount(`WHERE user_id = ?`, userID)
	return account, ok, err
}

func (s *sqliteStore) touchWebAccountLogin(userID int64) error {
	if userID == 0 {
		return nil
	}
	now := time.Now().UTC()
	_, err := s.db.Exec(`UPDATE web_accounts SET last_login_at = ?, updated_at = ? WHERE user_id = ?`,
		formatDBTime(now), formatDBTime(now), userID)
	return err
}

func (s *sqliteStore) updateWebAccountPassword(userID int64, passwordHash string) error {
	if userID == 0 {
		return errWebAuthRequired
	}
	if strings.TrimSpace(passwordHash) == "" {
		return errors.New("password hash must not be empty")
	}
	now := time.Now().UTC()
	result, err := s.db.Exec(`UPDATE web_accounts SET password_hash = ?, updated_at = ? WHERE user_id = ?`,
		passwordHash, formatDBTime(now), userID)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return errWebAuthRequired
	}
	return nil
}

func (s *sqliteStore) completeWebAccountProfile(userID int64, login string, passwordHash string) (webAccount, error) {
	if userID == 0 {
		return webAccount{}, errWebAuthRequired
	}
	login, err := normalizeWebLogin(login)
	if err != nil {
		return webAccount{}, err
	}
	if strings.TrimSpace(passwordHash) == "" {
		return webAccount{}, errors.New("password hash must not be empty")
	}
	existing, ok, err := s.getWebAccountByUserID(userID)
	if err != nil {
		return webAccount{}, err
	}
	if !ok {
		return webAccount{}, errWebAuthRequired
	}
	if other, _, ok, err := s.getWebAccountByLogin(login); err != nil {
		return webAccount{}, err
	} else if ok && other.UserID != userID {
		return webAccount{}, errWebLoginTaken
	}
	now := time.Now().UTC()
	result, err := s.db.Exec(`UPDATE web_accounts
		SET login = ?, login_normalized = ?, password_hash = ?, updated_at = ?
		WHERE user_id = ?`,
		login, login, passwordHash, formatDBTime(now), userID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "constraint") {
			return webAccount{}, errWebLoginTaken
		}
		return webAccount{}, err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return webAccount{}, err
	}
	if changed == 0 {
		return webAccount{}, errWebAuthRequired
	}
	user, ok, err := s.getUser(userID)
	if err != nil {
		return webAccount{}, err
	}
	if ok && strings.EqualFold(strings.TrimSpace(user.FirstName), strings.TrimSpace(existing.Login)) {
		user.FirstName = login
		user.UpdatedAt = now
		if err := s.saveUser(user); err != nil {
			return webAccount{}, err
		}
	}
	account, ok, err := s.getWebAccountByUserID(userID)
	if err != nil {
		return webAccount{}, err
	}
	if !ok {
		return webAccount{}, errWebAuthRequired
	}
	return account, nil
}

func (s *sqliteStore) authenticateTelegramWebAccount(currentUserID int64, profile telegramWebProfile) (webAccount, userState, error) {
	return s.authenticateTelegramWebAccountWithChoice(currentUserID, profile, "")
}

func (s *sqliteStore) authenticateTelegramWebAccountWithChoice(currentUserID int64, profile telegramWebProfile, mergeChoice string) (webAccount, userState, error) {
	if profile.ID <= 0 {
		return webAccount{}, userState{}, errors.New("telegram id is required")
	}
	now := time.Now().UTC()

	var currentAccount webAccount
	var currentPasswordHash string
	var currentAccountOK bool
	if currentUserID != 0 {
		account, passwordHash, ok, err := s.getWebAccount(`WHERE user_id = ?`, currentUserID)
		if err != nil {
			return webAccount{}, userState{}, err
		}
		currentAccount, currentPasswordHash, currentAccountOK = account, passwordHash, ok
	}

	telegramAccount, telegramPasswordHash, telegramAccountOK, err := s.getWebAccount(`WHERE user_id = ?`, profile.ID)
	if err != nil {
		return webAccount{}, userState{}, err
	}
	if currentUserID > 0 && currentUserID != profile.ID {
		return webAccount{}, userState{}, errWebAccountTGLocked
	}
	if currentAccountOK && currentAccount.UserID > 0 && currentAccount.UserID != profile.ID {
		return webAccount{}, userState{}, errWebAccountTGLocked
	}
	if currentAccountOK && currentAccount.UserID != profile.ID && telegramAccountOK &&
		!isTelegramPlaceholderWebAccount(telegramAccount, telegramPasswordHash) {
		return webAccount{}, userState{}, errWebTelegramUsed
	}

	var webUser userState
	var webUserOK bool
	if currentAccountOK && currentAccount.UserID != profile.ID {
		webUser, webUserOK, err = s.getUser(currentAccount.UserID)
		if err != nil {
			return webAccount{}, userState{}, err
		}
	}
	telegramUser, telegramUserOK, err := s.getUser(profile.ID)
	if err != nil {
		return webAccount{}, userState{}, err
	}

	finalUser := newUserState(profile.ID, now)
	webHasData := webUserOK && userHasMeaningfulTelegramData(webUser)
	telegramHasData := telegramUserOK && userHasMeaningfulTelegramData(telegramUser)
	switch {
	case mergeChoice == "web" && webHasData:
		finalUser = webUser
		if telegramUserOK {
			finalUser = mergeUserProgress(finalUser, telegramUser, now)
		}
	case mergeChoice == "telegram" && telegramHasData:
		finalUser = telegramUser
		if webUserOK {
			finalUser = mergeUserProgress(finalUser, webUser, now)
		}
	case telegramUserOK && userHasMeaningfulTelegramData(telegramUser):
		finalUser = telegramUser
		if webUserOK {
			finalUser = mergeUserProgress(finalUser, webUser, now)
		}
	case webUserOK:
		finalUser = webUser
		if telegramUserOK {
			finalUser = mergeUserProgress(finalUser, telegramUser, now)
		}
	case telegramUserOK:
		finalUser = telegramUser
	}
	finalUser.TelegramID = profile.ID
	finalUser.FirstName = telegramProfileDisplayName(profile, finalUser.FirstName)
	if finalUser.CreatedAt.IsZero() {
		finalUser.CreatedAt = now
	}
	finalUser.UpdatedAt = now
	if webUserOK && currentAccount.UserID != profile.ID && finalUser.ReferralCode == webUser.ReferralCode {
		finalUser.ReferralCode = ""
	}
	normalizeUser(&finalUser, now)
	if err := s.saveUser(finalUser); err != nil {
		return webAccount{}, userState{}, err
	}
	if webUserOK && currentAccount.UserID != profile.ID {
		if err := s.deleteUserStateRows(currentAccount.UserID); err != nil {
			return webAccount{}, userState{}, err
		}
	}

	account, err := s.ensureTelegramWebAccount(profile, currentAccount, currentPasswordHash, currentAccountOK, telegramAccount, telegramPasswordHash, telegramAccountOK)
	if err != nil {
		return webAccount{}, userState{}, err
	}
	if err := s.touchWebAccountLogin(account.UserID); err != nil {
		return webAccount{}, userState{}, err
	}
	user, err := s.mustGetUser(profile.ID)
	if err != nil {
		return webAccount{}, userState{}, err
	}
	return account, user, nil
}

func (s *sqliteStore) telegramWebProgressMergeChallenge(currentUserID int64, profile telegramWebProfile) (webProgressMergeChallenge, error) {
	if profile.ID <= 0 {
		return webProgressMergeChallenge{}, errors.New("telegram id is required")
	}
	var currentAccount webAccount
	var currentAccountOK bool
	if currentUserID != 0 {
		account, _, ok, err := s.getWebAccount(`WHERE user_id = ?`, currentUserID)
		if err != nil {
			return webProgressMergeChallenge{}, err
		}
		currentAccount, currentAccountOK = account, ok
	}
	telegramAccount, telegramPasswordHash, telegramAccountOK, err := s.getWebAccount(`WHERE user_id = ?`, profile.ID)
	if err != nil {
		return webProgressMergeChallenge{}, err
	}
	if currentUserID > 0 && currentUserID != profile.ID {
		return webProgressMergeChallenge{}, errWebAccountTGLocked
	}
	if currentAccountOK && currentAccount.UserID > 0 && currentAccount.UserID != profile.ID {
		return webProgressMergeChallenge{}, errWebAccountTGLocked
	}
	if currentAccountOK && currentAccount.UserID != profile.ID && telegramAccountOK &&
		!isTelegramPlaceholderWebAccount(telegramAccount, telegramPasswordHash) {
		return webProgressMergeChallenge{}, errWebTelegramUsed
	}
	if !currentAccountOK || currentAccount.UserID == profile.ID {
		return webProgressMergeChallenge{}, nil
	}

	webUser, webUserOK, err := s.getUser(currentAccount.UserID)
	if err != nil {
		return webProgressMergeChallenge{}, err
	}
	telegramUser, telegramUserOK, err := s.getUser(profile.ID)
	if err != nil {
		return webProgressMergeChallenge{}, err
	}
	webSummary := webProgressSummaryForUser("web", currentAccount.Login, webUser)
	telegramSummary := webProgressSummaryForUser("telegram", telegramProfileDisplayName(profile, telegramUser.FirstName), telegramUser)
	webSummary.HasData = webUserOK && userHasMeaningfulTelegramData(webUser)
	telegramSummary.HasData = telegramUserOK && userHasMeaningfulTelegramData(telegramUser)
	return webProgressMergeChallenge{
		Required: webSummary.HasData && telegramSummary.HasData,
		Web:      webSummary,
		Telegram: telegramSummary,
	}, nil
}

func (s *sqliteStore) ensureTelegramWebAccount(profile telegramWebProfile, currentAccount webAccount, currentPasswordHash string, currentAccountOK bool, telegramAccount webAccount, telegramPasswordHash string, telegramAccountOK bool) (webAccount, error) {
	now := time.Now().UTC()
	if currentAccountOK && currentAccount.UserID > 0 && currentAccount.UserID != profile.ID {
		return webAccount{}, errWebAccountTGLocked
	}
	if currentAccountOK && currentAccount.UserID != profile.ID && telegramAccountOK &&
		!isTelegramPlaceholderWebAccount(telegramAccount, telegramPasswordHash) {
		return webAccount{}, errWebTelegramUsed
	}
	desiredLogin := telegramWebLoginBase(profile)
	if telegramAccountOK && telegramAccount.Login != "" {
		desiredLogin = telegramAccount.Login
	}
	if currentAccountOK && currentAccount.Login != "" {
		desiredLogin = currentAccount.Login
	}
	var allowedUserIDs []int64
	if currentAccountOK {
		allowedUserIDs = append(allowedUserIDs, currentAccount.UserID)
	}
	login, err := s.uniqueWebLogin(desiredLogin, profile.ID, allowedUserIDs...)
	if err != nil {
		return webAccount{}, err
	}

	if telegramAccountOK {
		if currentAccountOK && currentAccount.UserID != profile.ID {
			if _, err := s.db.Exec(`DELETE FROM web_accounts WHERE user_id = ?`, currentAccount.UserID); err != nil {
				return webAccount{}, err
			}
		}
		if strings.TrimSpace(currentPasswordHash) == "" {
			currentPasswordHash = telegramPasswordHash
		}
		if strings.TrimSpace(currentPasswordHash) == "" {
			currentPasswordHash = telegramDisabledPasswordHash(profile.ID)
		}
		_, err = s.db.Exec(`UPDATE web_accounts
			SET login = ?, login_normalized = ?, password_hash = ?, updated_at = ?, last_login_at = ?
			WHERE user_id = ?`,
			login, login, currentPasswordHash, formatDBTime(now), formatDBTime(now), profile.ID)
		if err != nil {
			return webAccount{}, err
		}
		account, ok, err := s.getWebAccountByUserID(profile.ID)
		if err != nil {
			return webAccount{}, err
		}
		if !ok {
			return webAccount{}, errors.New("telegram web account not found after update")
		}
		return account, nil
	}

	if currentAccountOK {
		_, err = s.db.Exec(`UPDATE web_accounts
			SET user_id = ?, login = ?, login_normalized = ?, updated_at = ?, last_login_at = ?
			WHERE user_id = ?`,
			profile.ID, login, login, formatDBTime(now), formatDBTime(now), currentAccount.UserID)
		if err != nil {
			return webAccount{}, err
		}
		account, ok, err := s.getWebAccountByUserID(profile.ID)
		if err != nil {
			return webAccount{}, err
		}
		if !ok {
			return webAccount{}, errors.New("telegram web account not found after link")
		}
		return account, nil
	}

	return s.createWebAccount(profile.ID, login, telegramDisabledPasswordHash(profile.ID))
}

func (s *sqliteStore) uniqueWebLogin(base string, reservedUserID int64, allowedUserIDs ...int64) (string, error) {
	base, err := normalizeWebLogin(base)
	if err != nil {
		base = "tg_" + strconv.FormatInt(reservedUserID, 10)
	}
	base = strings.Trim(base, ".-_")
	if len(base) < 3 {
		base = "tg_" + strconv.FormatInt(reservedUserID, 10)
	}
	if len(base) > 32 {
		base = base[:32]
	}
	for attempt := 0; attempt < 100; attempt++ {
		candidate := base
		if attempt > 0 {
			suffix := "_" + strconv.Itoa(attempt+1)
			cut := 32 - len(suffix)
			if cut < 3 {
				cut = 3
			}
			if len(candidate) > cut {
				candidate = strings.Trim(candidate[:cut], ".-_")
			}
			candidate += suffix
		}
		candidate, err = normalizeWebLogin(candidate)
		if err != nil {
			continue
		}
		existing, _, ok, err := s.getWebAccountByLogin(candidate)
		if err != nil {
			return "", err
		}
		if !ok || existing.UserID == reservedUserID || int64In(existing.UserID, allowedUserIDs) {
			return candidate, nil
		}
	}
	return "", errors.New("failed to allocate telegram web login")
}

func int64In(value int64, values []int64) bool {
	for _, item := range values {
		if value == item {
			return true
		}
	}
	return false
}

func (s *sqliteStore) deleteUserStateRows(userID int64) error {
	if userID == 0 {
		return nil
	}
	if _, err := s.db.Exec(`DELETE FROM mistakes WHERE telegram_id = ?`, userID); err != nil {
		return err
	}
	if _, err := s.db.Exec(`DELETE FROM learned_words WHERE telegram_id = ?`, userID); err != nil {
		return err
	}
	if _, err := s.db.Exec(`DELETE FROM users WHERE telegram_id = ?`, userID); err != nil {
		return err
	}
	return nil
}

func telegramWebLoginBase(profile telegramWebProfile) string {
	username := strings.TrimSpace(profile.Username)
	if username != "" {
		return username
	}
	return "tg_" + strconv.FormatInt(profile.ID, 10)
}

func telegramProfileDisplayName(profile telegramWebProfile, fallback string) string {
	if strings.TrimSpace(profile.Username) != "" {
		return strings.TrimSpace(profile.Username)
	}
	parts := []string{strings.TrimSpace(profile.FirstName), strings.TrimSpace(profile.LastName)}
	name := strings.TrimSpace(strings.Join(parts, " "))
	if name != "" {
		return name
	}
	if strings.TrimSpace(fallback) != "" {
		return strings.TrimSpace(fallback)
	}
	return "telegram_" + strconv.FormatInt(profile.ID, 10)
}

func telegramDisabledPasswordHash(userID int64) string {
	return "telegram_login_disabled:" + strconv.FormatInt(userID, 10)
}

func isTelegramPlaceholderWebAccount(account webAccount, passwordHash string) bool {
	return account.UserID > 0 && !isWebPasswordHash(passwordHash)
}

func userHasMeaningfulTelegramData(user userState) bool {
	return user.XP > 0 ||
		user.LessonCount > 0 ||
		user.PracticeCount > 0 ||
		user.VoiceCount > 0 ||
		user.WordLessonCount > 0 ||
		user.WordGameCount > 0 ||
		len(user.LessonHistory) > 0 ||
		len(user.PracticeHistory) > 0 ||
		len(user.Mistakes) > 0 ||
		len(user.LearnedWords) > 0 ||
		strings.TrimSpace(user.LastLessonPrompt) != "" ||
		user.Plan == "premium" ||
		!user.PremiumUntil.IsZero() ||
		user.InvitedBy != 0 ||
		user.ReferralCount > 0 ||
		user.ReferralBalanceKopecks > 0
}

func webProgressSummaryForUser(source string, label string, user userState) webProgressSummary {
	label = strings.TrimSpace(label)
	if label == "" {
		label = strings.TrimSpace(user.FirstName)
	}
	if label == "" {
		label = source
	}
	return webProgressSummary{
		Source:   source,
		Label:    label,
		XP:       user.XP,
		Level:    user.Level,
		Lessons:  user.LessonCount,
		Practice: user.PracticeCount,
		Words:    masteredWordCount(user),
		Mistakes: len(user.Mistakes),
		HasData:  userHasMeaningfulTelegramData(user),
	}
}

func preservePremiumEntitlement(primary userState, secondary userState, now time.Time) userState {
	if !secondary.isPremium(now) {
		return primary
	}
	if !primary.isPremium(now) || secondary.PremiumUntil.After(primary.PremiumUntil) {
		primary.Plan = "premium"
		primary.PremiumUntil = secondary.PremiumUntil
	}
	if primary.LastPaymentChargeID == "" {
		primary.LastPaymentChargeID = secondary.LastPaymentChargeID
	}
	return primary
}

func mergeUserProgress(primary userState, secondary userState, now time.Time) userState {
	primary = preservePremiumEntitlement(primary, secondary, now)
	if primary.InvitedBy == 0 {
		primary.InvitedBy = secondary.InvitedBy
	}
	primary.ReferralCount += secondary.ReferralCount
	primary.ReferralBalanceKopecks += secondary.ReferralBalanceKopecks
	primary.LessonCount += secondary.LessonCount
	primary.PracticeCount += secondary.PracticeCount
	primary.VoiceCount += secondary.VoiceCount
	primary.WordLessonCount += secondary.WordLessonCount
	primary.WordGameCount += secondary.WordGameCount
	primary.XP += secondary.XP
	primary.LessonsToday = maxInt(primary.LessonsToday, secondary.LessonsToday)
	primary.PracticeToday = maxInt(primary.PracticeToday, secondary.PracticeToday)
	primary.VoiceToday = maxInt(primary.VoiceToday, secondary.VoiceToday)
	if primary.LastPaymentChargeID == "" {
		primary.LastPaymentChargeID = secondary.LastPaymentChargeID
	}
	if primary.LastLessonPrompt == "" {
		primary.LastLessonPrompt = secondary.LastLessonPrompt
	}
	primary.LessonHistory = trimLessonHistory(append(primary.LessonHistory, secondary.LessonHistory...))
	if primary.Mode == "" || primary.Mode == "idle" {
		primary.Mode = secondary.Mode
	}
	if primary.InterfaceLanguage == "" || !primary.InterfaceSelected {
		primary.InterfaceLanguage = secondary.InterfaceLanguage
		primary.InterfaceSelected = secondary.InterfaceSelected
	}
	if primary.LearningLanguage == "" || !primary.LanguageSelected {
		primary.LearningLanguage = secondary.LearningLanguage
		primary.LanguageSelected = secondary.LanguageSelected
	}
	if primary.Level == "" || primary.Level == "A2" {
		primary.Level = secondary.Level
	}
	if !secondary.CreatedAt.IsZero() && (primary.CreatedAt.IsZero() || secondary.CreatedAt.Before(primary.CreatedAt)) {
		primary.CreatedAt = secondary.CreatedAt
	}
	if secondary.UpdatedAt.After(primary.UpdatedAt) {
		primary.UpdatedAt = secondary.UpdatedAt
	}
	primary.PracticeHistory = trimPracticeHistory(append(primary.PracticeHistory, secondary.PracticeHistory...), practiceMemoryLimit)
	primary.Mistakes = mergeMistakeEntries(primary.Mistakes, secondary.Mistakes)
	primary.LearnedWords = mergeLearnedWordEntries(primary.LearnedWords, secondary.LearnedWords)
	return primary
}

func mergeMistakeEntries(primary []mistakeEntry, secondary []mistakeEntry) []mistakeEntry {
	seen := map[string]bool{}
	result := make([]mistakeEntry, 0, len(primary)+len(secondary))
	for _, item := range append(primary, secondary...) {
		key := normalizeLearningLanguage(item.Language) + "\x00" + item.Word + "\x00" + item.Correction
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, item)
	}
	return result
}

func mergeLearnedWordEntries(primary []learnedWordEntry, secondary []learnedWordEntry) []learnedWordEntry {
	seen := map[string]int{}
	result := make([]learnedWordEntry, 0, len(primary)+len(secondary))
	for _, item := range append(primary, secondary...) {
		key := normalizeLearningLanguage(item.Language) + "\x00" + legacyVocabID(item.ID)
		if index, ok := seen[key]; ok {
			if item.ReviewCorrectCount > result[index].ReviewCorrectCount {
				result[index].ReviewCorrectCount = item.ReviewCorrectCount
			}
			if item.SpellingCorrectCount > result[index].SpellingCorrectCount {
				result[index].SpellingCorrectCount = item.SpellingCorrectCount
			}
			continue
		}
		seen[key] = len(result)
		result = append(result, item)
	}
	return result
}

func maxInt(a int, b int) int {
	if b > a {
		return b
	}
	return a
}

func (s *sqliteStore) getWebAccount(where string, arg any) (webAccount, string, bool, error) {
	row := s.db.QueryRow(`SELECT user_id, login, password_hash, created_at, updated_at, last_login_at FROM web_accounts `+where, arg)
	var account webAccount
	var passwordHash, createdAt, updatedAt, lastLoginAt string
	err := row.Scan(&account.UserID, &account.Login, &passwordHash, &createdAt, &updatedAt, &lastLoginAt)
	if errors.Is(err, sql.ErrNoRows) {
		return webAccount{}, "", false, nil
	}
	if err != nil {
		return webAccount{}, "", false, err
	}
	account.CreatedAt = parseDBTime(createdAt)
	account.UpdatedAt = parseDBTime(updatedAt)
	account.LastLoginAt = parseDBTime(lastLoginAt)
	account.PasswordSet = isWebPasswordHash(passwordHash)
	return account, passwordHash, true, nil
}

func (s *sqliteStore) importJSONIfEmpty(path string) error {
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count); err != nil {
		return err
	}
	if count > 0 || path == "" {
		return nil
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	jsonStore, err := newJSONStore(path)
	if err != nil {
		return err
	}
	for _, user := range jsonStore.users {
		normalizeUser(&user, time.Now().UTC())
		if err := s.saveUser(user); err != nil {
			return err
		}
	}
	return nil
}

func (s *sqliteStore) getOrCreateUser(telegramID int64, firstName string) (userState, error) {
	user, ok, err := s.getUser(telegramID)
	if err != nil {
		return userState{}, err
	}
	now := time.Now().UTC()
	if !ok {
		user = newUserState(telegramID, now)
	}
	normalizeUser(&user, now)
	displayName := strings.TrimSpace(firstName)
	if telegramID > 0 && ok && displayName != "" {
		if account, accountOK, err := s.getWebAccountByUserID(telegramID); err == nil && accountOK &&
			strings.EqualFold(strings.TrimSpace(account.Login), displayName) &&
			strings.TrimSpace(user.FirstName) != "" {
			displayName = user.FirstName
		}
	}
	user.FirstName = displayName
	user.UpdatedAt = now
	if err := s.saveUser(user); err != nil {
		return userState{}, err
	}
	return s.mustGetUser(telegramID)
}

func (s *sqliteStore) applyReferral(newUserID int64, inviterID int64) (referralApplication, error) {
	if inviterID == 0 || inviterID == newUserID {
		return referralApplication{}, nil
	}
	newUser, _, err := s.getOrNewUser(newUserID)
	if err != nil {
		return referralApplication{}, err
	}
	if newUser.InvitedBy != 0 {
		return referralApplication{}, nil
	}
	inviter, ok, err := s.getUser(inviterID)
	if err != nil || !ok {
		return referralApplication{}, err
	}
	now := time.Now().UTC()
	normalizeUser(&newUser, now)
	normalizeUser(&inviter, now)

	inviteeDays := referralInviteePremiumDays

	newUser.InvitedBy = inviterID
	newUser.Plan = "premium"
	newUser.PremiumUntil = premiumBase(newUser, now).Add(referralPremiumDuration(inviteeDays)).UTC()
	newUser.UpdatedAt = now
	inviter.ReferralCount++
	inviter.UpdatedAt = now

	if err := s.saveUser(newUser); err != nil {
		return referralApplication{}, err
	}
	if err := s.saveUser(inviter); err != nil {
		return referralApplication{}, err
	}
	levelReward, err := s.rewardReferralLevel(newUserID)
	if err != nil {
		return referralApplication{}, err
	}
	return referralApplication{
		Applied:            true,
		InviterID:          inviterID,
		InviteeID:          newUserID,
		InviteePremiumDays: inviteeDays,
		InviterPremiumDays: levelReward.InviterPremiumDays,
		NewUserUntil:       newUser.PremiumUntil,
		InviteeUntil:       newUser.PremiumUntil,
		InviterUntil:       levelReward.InviterUntil,
		ReferralCount:      inviter.ReferralCount,
	}, nil
}

func (s *sqliteStore) getUserByReferralCode(code string) (userState, bool, error) {
	code, err := normalizeWebReferralCode(code)
	if err != nil || code == "" {
		return userState{}, false, err
	}
	var userID int64
	err = s.db.QueryRow(`SELECT telegram_id FROM users WHERE referral_code = ?`, code).Scan(&userID)
	if errors.Is(err, sql.ErrNoRows) {
		return userState{}, false, nil
	}
	if err != nil {
		return userState{}, false, err
	}
	return s.getUser(userID)
}

func (s *sqliteStore) applyWebReferral(newUserID int64, code string) (webReferralResult, error) {
	code, err := normalizeWebReferralCode(code)
	if err != nil || code == "" {
		return webReferralResult{}, err
	}
	newUser, _, err := s.getOrNewUser(newUserID)
	if err != nil {
		return webReferralResult{}, err
	}
	if newUser.InvitedBy != 0 {
		return webReferralResult{Applied: false, ReferralCode: code}, nil
	}
	inviter, ok, err := s.getUserByReferralCode(code)
	if err != nil {
		return webReferralResult{}, err
	}
	if !ok {
		return webReferralResult{}, errWebReferralNotFound
	}
	if inviter.TelegramID == newUserID {
		return webReferralResult{}, errWebReferralSelf
	}

	now := time.Now().UTC()
	normalizeUser(&newUser, now)
	normalizeUser(&inviter, now)

	inviteeDays := referralInviteePremiumDays

	newUser.InvitedBy = inviter.TelegramID
	newUser.Plan = "premium"
	newUser.PremiumUntil = premiumBase(newUser, now).Add(referralPremiumDuration(inviteeDays)).UTC()
	newUser.UpdatedAt = now

	inviter.ReferralCount++
	inviter.UpdatedAt = now

	if err := s.saveUser(newUser); err != nil {
		return webReferralResult{}, err
	}
	if err := s.saveUser(inviter); err != nil {
		return webReferralResult{}, err
	}
	levelReward, err := s.rewardReferralLevel(newUserID)
	if err != nil {
		return webReferralResult{}, err
	}
	return webReferralResult{
		Applied:            true,
		InviterID:          inviter.TelegramID,
		InviteeID:          newUserID,
		InviteePremiumDays: inviteeDays,
		InviterPremiumDays: levelReward.InviterPremiumDays,
		NewUserUntil:       newUser.PremiumUntil,
		InviteeUntil:       newUser.PremiumUntil,
		InviterUntil:       levelReward.InviterUntil,
		ReferralCode:       code,
		ReferralCount:      inviter.ReferralCount,
	}, nil
}

func (s *sqliteStore) rewardReferralLevel(telegramID int64) (referralLevelReward, error) {
	invitee, ok, err := s.getUser(telegramID)
	if err != nil || !ok {
		return referralLevelReward{}, err
	}
	now := time.Now().UTC()
	normalizeUser(&invitee, now)
	if invitee.InvitedBy == 0 || invitee.ReferralLevelRewarded || !referralLevelReached(invitee.XP) {
		return referralLevelReward{}, nil
	}
	inviter, ok, err := s.getUser(invitee.InvitedBy)
	if err != nil || !ok {
		return referralLevelReward{}, err
	}
	normalizeUser(&inviter, now)

	inviter.Plan = paidTierForExtension(inviter.Plan, "premium", now)
	inviter.PremiumUntil = premiumBase(inviter, now).Add(referralPremiumDuration(referralInviterLevelPremiumDays)).UTC()
	inviter.UpdatedAt = now
	invitee.ReferralLevelRewarded = true
	invitee.UpdatedAt = now

	if err := s.saveUser(invitee); err != nil {
		return referralLevelReward{}, err
	}
	if err := s.saveUser(inviter); err != nil {
		return referralLevelReward{}, err
	}
	return referralLevelReward{
		Applied:            true,
		InviteeID:          telegramID,
		InviterID:          inviter.TelegramID,
		InviterPremiumDays: referralInviterLevelPremiumDays,
		InviterUntil:       inviter.PremiumUntil,
	}, nil
}

func (s *sqliteStore) creditReferralPurchase(telegramID int64, chargeID string, paidKopecks int64) (referralPurchaseReward, error) {
	chargeID = strings.TrimSpace(chargeID)
	if chargeID == "" || paidKopecks <= 0 {
		return referralPurchaseReward{}, nil
	}
	buyer, ok, err := s.getUser(telegramID)
	if err != nil || !ok || buyer.InvitedBy == 0 {
		return referralPurchaseReward{}, err
	}
	direct, ok, err := s.getUser(buyer.InvitedBy)
	if err != nil || !ok {
		return referralPurchaseReward{}, err
	}
	directID := direct.TelegramID
	directReward := referralDirectRewardKopecks(paidKopecks)

	indirectID := direct.InvitedBy
	indirectReward := int64(0)
	if indirectID == telegramID || indirectID == directID {
		indirectID = 0
	}
	if indirectID != 0 {
		if _, ok, err := s.getUser(indirectID); err != nil {
			return referralPurchaseReward{}, err
		} else if ok {
			indirectReward = referralIndirectRewardKopecks(paidKopecks)
		} else {
			indirectID = 0
		}
	}

	now := time.Now().UTC()
	tx, err := s.db.Begin()
	if err != nil {
		return referralPurchaseReward{}, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`INSERT OR IGNORE INTO referral_purchase_rewards (
		payment_id, buyer_id, direct_inviter_id, indirect_inviter_id, paid_kopecks,
		direct_reward_kopecks, indirect_reward_kopecks, created_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		chargeID, telegramID, directID, indirectID, paidKopecks, directReward, indirectReward, formatDBTime(now))
	if err != nil {
		return referralPurchaseReward{}, err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return referralPurchaseReward{}, err
	}
	if changed == 0 {
		if err := tx.Commit(); err != nil {
			return referralPurchaseReward{}, err
		}
		return referralPurchaseReward{BuyerID: telegramID, PaymentID: chargeID, PaidKopecks: paidKopecks}, nil
	}

	if directReward > 0 {
		if _, err := tx.Exec(`UPDATE users SET referral_balance_kopecks = referral_balance_kopecks + ?, updated_at = ? WHERE telegram_id = ?`,
			directReward, formatDBTime(now), directID); err != nil {
			return referralPurchaseReward{}, err
		}
	}
	if indirectID != 0 && indirectReward > 0 {
		if _, err := tx.Exec(`UPDATE users SET referral_balance_kopecks = referral_balance_kopecks + ?, updated_at = ? WHERE telegram_id = ?`,
			indirectReward, formatDBTime(now), indirectID); err != nil {
			return referralPurchaseReward{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return referralPurchaseReward{}, err
	}
	return referralPurchaseReward{
		Applied:               true,
		BuyerID:               telegramID,
		PaymentID:             chargeID,
		PaidKopecks:           paidKopecks,
		DirectInviterID:       directID,
		IndirectInviterID:     indirectID,
		DirectRewardKopecks:   directReward,
		IndirectRewardKopecks: indirectReward,
	}, nil
}

func premiumBase(user userState, now time.Time) time.Time {
	if user.isPremium(now) {
		return user.PremiumUntil
	}
	return now
}

func (s *sqliteStore) setMode(telegramID int64, mode string) error {
	return s.updateUser(telegramID, func(user *userState) { user.Mode = mode })
}

func (s *sqliteStore) setInterfaceLanguage(telegramID int64, language string) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.InterfaceLanguage = normalizeInterfaceLanguage(language)
		user.InterfaceSelected = true
		user.Mode = "idle"
	})
}

func (s *sqliteStore) setLearningLanguage(telegramID int64, language string) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.LearningLanguage = normalizeLearningLanguage(language)
		user.LanguageSelected = true
		user.Mode = "idle"
	})
}

func (s *sqliteStore) setUserLevel(telegramID int64, level string) error {
	return s.updateUser(telegramID, func(user *userState) { user.Level = normalizeCEFRLevel(level) })
}

func (s *sqliteStore) setLearningFocus(telegramID int64, focus string) error {
	return s.updateUser(telegramID, func(user *userState) { user.LearningFocus = strings.TrimSpace(focus) })
}

func (s *sqliteStore) setNavigationLayout(telegramID int64, layout navigationLayout) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.NavigationLayout = normalizeNavigationLayout(layout)
	})
}

func (s *sqliteStore) addXP(telegramID int64, amount int) error {
	if amount <= 0 {
		return nil
	}
	return s.updateUser(telegramID, func(user *userState) { user.XP += amount })
}

func (s *sqliteStore) recordHabitLogin(telegramID int64) error {
	return s.updateUser(telegramID, func(user *userState) {
		recordHabitDay(user, time.Now().UTC(), true)
	})
}

func (s *sqliteStore) claimDailyBonus(telegramID int64, date string, amount int) (bool, error) {
	date = strings.TrimSpace(date)
	if date == "" || amount <= 0 {
		return false, nil
	}
	now := time.Now().UTC()
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	var lastCreatedAt string
	err = tx.QueryRow(`SELECT created_at FROM daily_bonus_claims WHERE telegram_id = ? ORDER BY created_at DESC LIMIT 1`, telegramID).Scan(&lastCreatedAt)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if lastCreatedAt != "" {
		if last := parseDBTime(lastCreatedAt); !last.IsZero() && now.Sub(last) < 24*time.Hour {
			return false, tx.Commit()
		}
	}
	result, err := tx.Exec(
		`INSERT OR IGNORE INTO daily_bonus_claims (telegram_id, claim_date, xp, created_at) VALUES (?, ?, ?, ?)`,
		telegramID,
		date,
		amount,
		formatDBTime(now),
	)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows == 0 {
		return false, tx.Commit()
	}
	if _, err := tx.Exec(`UPDATE users SET xp = xp + ?, updated_at = ? WHERE telegram_id = ?`, amount, formatDBTime(now), telegramID); err != nil {
		return false, err
	}
	if _, err := tx.Exec(
		`INSERT INTO habit_days (telegram_id, habit_date, login, complete, claimed, claimed_at, lessons, practice, voice, updated_at)
		VALUES (?, ?, 1, 0, 1, ?, 0, 0, 0, ?)
		ON CONFLICT(telegram_id, habit_date) DO UPDATE SET
			login=1, complete=MAX(habit_days.complete, excluded.complete), claimed=1,
			claimed_at=excluded.claimed_at, updated_at=excluded.updated_at`,
		telegramID,
		date,
		formatDBTime(now),
		formatDBTime(now),
	); err != nil {
		return false, err
	}
	return true, tx.Commit()
}

func (s *sqliteStore) saveLesson(telegramID int64, prompt string) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.Mode = "lesson"
		user.LastLessonPrompt = prompt
		user.LessonHistory = trimLessonHistory(append(user.LessonHistory, prompt))
		user.LessonCount++
		user.LessonsToday++
		recordHabitDay(user, time.Now().UTC(), true)
	})
}

func (s *sqliteStore) completeTutorLesson(telegramID int64, lessonID string, amount int) (bool, error) {
	lessonID = strings.TrimSpace(lessonID)
	if telegramID == 0 || lessonID == "" || amount <= 0 {
		return false, nil
	}
	now := time.Now().UTC()
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	result, err := tx.Exec(
		`UPDATE tutor_user_lessons
		SET completed_at = ?
		WHERE telegram_id = ? AND lesson_id = ? AND (completed_at = '' OR completed_at IS NULL)`,
		formatDBTime(now),
		telegramID,
		lessonID,
	)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	if rows == 0 {
		return false, nil
	}
	return true, s.updateUser(telegramID, func(user *userState) {
		user.XP += amount
		user.LessonCount++
		user.LessonsToday++
		recordHabitDay(user, now, true)
	})
}

func (s *sqliteStore) nextTutorLesson(user userState, factory tutorLessonFactory) (tutorLesson, error) {
	language, interfaceLanguage, level := tutorLessonStorageContext(user)
	tx, err := s.db.Begin()
	if err != nil {
		return tutorLesson{}, err
	}
	defer tx.Rollback()

	if lesson, ok, err := s.findUnassignedTutorLessonTx(tx, user.TelegramID, language, interfaceLanguage, level); err != nil {
		return tutorLesson{}, err
	} else if ok {
		if assigned, err := s.assignTutorLessonTx(tx, user.TelegramID, lesson.ID); err != nil {
			return tutorLesson{}, err
		} else if !assigned {
			return tutorLesson{}, errors.New("failed to assign reusable tutor lesson")
		}
		if err := tx.Commit(); err != nil {
			return tutorLesson{}, err
		}
		return tutorPersonalizeLessonForUser(lesson, user), nil
	}

	startSequence, err := s.tutorAssignedLessonCountTx(tx, user.TelegramID, language, interfaceLanguage, level)
	if err != nil {
		return tutorLesson{}, err
	}
	for attempt := 0; attempt < tutorLessonGenerationAttemptLimit(); attempt++ {
		lesson, err := factory(startSequence + attempt)
		if err != nil {
			return tutorLesson{}, err
		}
		if lesson.ID == "" || !tutorLessonMatchesContext(lesson, language, interfaceLanguage, level) {
			continue
		}
		if err := s.saveTutorLessonTx(tx, lesson); err != nil {
			return tutorLesson{}, err
		}
		assigned, err := s.assignTutorLessonTx(tx, user.TelegramID, lesson.ID)
		if err != nil {
			return tutorLesson{}, err
		}
		if !assigned {
			continue
		}
		if err := tx.Commit(); err != nil {
			return tutorLesson{}, err
		}
		return tutorPersonalizeLessonForUser(lesson, user), nil
	}
	return tutorLesson{}, errors.New("no unique tutor lesson available")
}

func (s *sqliteStore) findUnassignedTutorLessonTx(tx *sql.Tx, telegramID int64, language string, interfaceLanguage string, level string) (tutorLesson, bool, error) {
	var payload string
	err := tx.QueryRow(
		`SELECT payload_json
		FROM tutor_lessons
		WHERE learning_language = ? AND interface_language = ? AND level = ?
			AND NOT EXISTS (
				SELECT 1 FROM tutor_user_lessons
				WHERE tutor_user_lessons.telegram_id = ? AND tutor_user_lessons.lesson_id = tutor_lessons.id
			)
		ORDER BY lesson_number ASC, id ASC
		LIMIT 1`,
		normalizeLearningLanguage(language), normalizeInterfaceLanguage(interfaceLanguage), normalizeTutorLevel(level), telegramID,
	).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return tutorLesson{}, false, nil
	}
	if err != nil {
		return tutorLesson{}, false, err
	}
	var lesson tutorLesson
	if err := json.Unmarshal([]byte(payload), &lesson); err != nil {
		return tutorLesson{}, false, err
	}
	return lesson, lesson.ID != "", nil
}

func (s *sqliteStore) tutorAssignedLessonCountTx(tx *sql.Tx, telegramID int64, language string, interfaceLanguage string, level string) (int, error) {
	var count int
	err := tx.QueryRow(
		`SELECT COUNT(*)
		FROM tutor_user_lessons
		JOIN tutor_lessons ON tutor_lessons.id = tutor_user_lessons.lesson_id
		WHERE tutor_user_lessons.telegram_id = ?
			AND tutor_lessons.learning_language = ?
			AND tutor_lessons.interface_language = ?
			AND tutor_lessons.level = ?`,
		telegramID, normalizeLearningLanguage(language), normalizeInterfaceLanguage(interfaceLanguage), normalizeTutorLevel(level),
	).Scan(&count)
	return count, err
}

func (s *sqliteStore) saveTutorLessonTx(tx *sql.Tx, lesson tutorLesson) error {
	if strings.TrimSpace(lesson.ID) == "" {
		return errors.New("tutor lesson id is empty")
	}
	payload, err := json.Marshal(lesson)
	if err != nil {
		return err
	}
	now := formatDBTime(time.Now().UTC())
	_, err = tx.Exec(
		`INSERT INTO tutor_lessons (
			id, learning_language, interface_language, level, lesson_number, topic, payload_json, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			learning_language=excluded.learning_language,
			interface_language=excluded.interface_language,
			level=excluded.level,
			lesson_number=excluded.lesson_number,
			topic=excluded.topic,
			payload_json=excluded.payload_json,
			updated_at=excluded.updated_at`,
		lesson.ID,
		normalizeLearningLanguage(lesson.LearningLanguage),
		normalizeInterfaceLanguage(lesson.InterfaceLanguage),
		normalizeTutorLevel(lesson.Level),
		lesson.LessonNumber,
		strings.TrimSpace(lesson.Topic),
		string(payload),
		now,
		now,
	)
	return err
}

func (s *sqliteStore) assignTutorLessonTx(tx *sql.Tx, telegramID int64, lessonID string) (bool, error) {
	if telegramID == 0 || strings.TrimSpace(lessonID) == "" {
		return false, nil
	}
	now := formatDBTime(time.Now().UTC())
	result, err := tx.Exec(
		`INSERT OR IGNORE INTO tutor_user_lessons (telegram_id, lesson_id, assigned_at, completed_at, lesson_order)
		VALUES (
			?, ?, ?, '',
			COALESCE((SELECT MAX(lesson_order) + 1 FROM tutor_user_lessons WHERE telegram_id = ?), 1)
		)`,
		telegramID, lessonID, now, telegramID,
	)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (s *sqliteStore) referralInvitees(telegramID int64, limit int) ([]referralInviteeEntry, error) {
	if telegramID == 0 {
		return nil, nil
	}
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.Query(
		`SELECT telegram_id, first_name, level, xp, referral_level_rewarded, created_at, updated_at
		FROM users
		WHERE invited_by = ?
		ORDER BY created_at DESC
		LIMIT ?`,
		telegramID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []referralInviteeEntry
	for rows.Next() {
		var user userState
		var rewarded int
		var createdAt, updatedAt string
		if err := rows.Scan(&user.TelegramID, &user.FirstName, &user.Level, &user.XP, &rewarded, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		user.InvitedBy = telegramID
		user.ReferralLevelRewarded = rewarded != 0
		user.CreatedAt = parseDBTime(createdAt)
		user.UpdatedAt = parseDBTime(updatedAt)
		items = append(items, referralInviteeFromUser(user))
	}
	return items, rows.Err()
}

func (s *sqliteStore) incrementPractice(telegramID int64) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.PracticeCount++
		user.PracticeToday++
		user.XP += 5
		recordHabitDay(user, time.Now().UTC(), true)
	})
}

func (s *sqliteStore) savePracticeHistory(telegramID int64, history []string) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.PracticeHistory = trimPracticeHistory(history, practiceMemoryLimit)
	})
}

func (s *sqliteStore) incrementVoice(telegramID int64) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.VoiceCount++
		user.VoiceToday++
		user.XP += 3
		recordHabitDay(user, time.Now().UTC(), true)
	})
}

func (s *sqliteStore) addLearnedWord(telegramID int64, word vocabWord) (bool, int, error) {
	user, _, err := s.getOrNewUser(telegramID)
	if err != nil {
		return false, 0, err
	}
	for _, existing := range user.LearnedWords {
		if normalizeLearningLanguage(existing.Language) == normalizeLearningLanguage(word.Language) &&
			(existing.ID == word.ID || legacyVocabID(existing.ID) == legacyVocabID(word.ID)) {
			return false, masteredWordCountForLanguage(user, word.Language), nil
		}
	}
	user.LearnedWords = append(user.LearnedWords, learnedWordEntry{ID: word.ID, Language: word.Language, Russian: word.Russian, English: word.English, Context: word.Context, LearnedAt: time.Now().UTC()})
	user.WordLessonCount++
	user.XP += 15
	user.UpdatedAt = time.Now().UTC()
	return true, masteredWordCountForLanguage(user, word.Language), s.saveUser(user)
}

func (s *sqliteStore) markWordGameCorrect(telegramID int64, wordID string) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.WordGameCount++
		user.XP += 10
		incrementLearnedWordReview(user, wordID)
	})
}

func (s *sqliteStore) markSpellingCorrect(telegramID int64, wordID string) error {
	return s.updateUser(telegramID, func(user *userState) {
		incrementLearnedWordSpelling(user, wordID)
	})
}

func (s *sqliteStore) extendPremium(telegramID int64, chargeID string, duration time.Duration, tier string) (time.Time, error) {
	var until time.Time
	err := s.updateUser(telegramID, func(user *userState) {
		if chargeID != "" && user.LastPaymentChargeID == chargeID {
			until = user.PremiumUntil
			return
		}
		base := time.Now().UTC()
		if user.isPremium(base) {
			base = user.PremiumUntil
		}
		until = base.Add(duration)
		user.Plan = paidTierForExtension(user.Plan, tier, base)
		user.PremiumUntil = until.UTC()
		user.LastPaymentChargeID = chargeID
	})
	return until, err
}

func (s *sqliteStore) dueReminderUsers(now time.Time, limit int) ([]reminderTarget, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query(`SELECT telegram_id, first_name, interface_language, reminder_utc_offset, reminder_hour, last_reminder_date
		FROM users
		WHERE reminder_enabled = 1
		ORDER BY telegram_id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	targets := make([]reminderTarget, 0, limit)
	for rows.Next() {
		user := userState{
			ReminderEnabled: true,
		}
		if err := rows.Scan(&user.TelegramID, &user.FirstName, &user.InterfaceLanguage, &user.ReminderUTCOffset, &user.ReminderHour, &user.LastReminderDate); err != nil {
			return nil, err
		}
		if target, ok := reminderDue(user, now); ok {
			targets = append(targets, target)
			if len(targets) >= limit {
				break
			}
		}
	}
	return targets, rows.Err()
}

func (s *sqliteStore) markReminderSent(telegramID int64, localDate string) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.LastReminderDate = localDate
	})
}

func (s *sqliteStore) setReminderEnabled(telegramID int64, enabled bool) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.ReminderEnabled = enabled
	})
}

func (s *sqliteStore) setTimezone(telegramID int64, offsetMinutes int) error {
	return s.updateUser(telegramID, func(user *userState) {
		user.ReminderUTCOffset = offsetMinutes
		user.TimezoneSelected = true
	})
}

func (s *sqliteStore) addMistakes(telegramID int64, language string, entries []mistakeEntry) error {
	if len(entries) == 0 {
		return nil
	}
	language = normalizeLearningLanguage(language)
	user, _, err := s.getOrNewUser(telegramID)
	if err != nil {
		return err
	}
	existing := map[string]bool{}
	for _, item := range user.Mistakes {
		existing[normalizeLearningLanguage(item.Language)+"|"+item.Word+"|"+item.Correction] = true
	}
	for _, entry := range entries {
		entry.Language = language
		key := entry.Language + "|" + entry.Word + "|" + entry.Correction
		if !existing[key] {
			user.Mistakes = append(user.Mistakes, entry)
			existing[key] = true
		}
	}
	if len(user.Mistakes) > 200 {
		user.Mistakes = user.Mistakes[len(user.Mistakes)-200:]
	}
	user.UpdatedAt = time.Now().UTC()
	return s.saveUser(user)
}

func (s *sqliteStore) clearMistakes(telegramID int64, language string) error {
	language = normalizeLearningLanguage(language)
	user, _, err := s.getOrNewUser(telegramID)
	if err != nil {
		return err
	}
	filtered := user.Mistakes[:0]
	for _, mistake := range user.Mistakes {
		if normalizeLearningLanguage(mistake.Language) != language {
			filtered = append(filtered, mistake)
		}
	}
	user.Mistakes = filtered
	user.UpdatedAt = time.Now().UTC()
	return s.saveUser(user)
}

func (s *sqliteStore) removeMistake(telegramID int64, language string, word string, correction string) error {
	language = normalizeLearningLanguage(language)
	user, _, err := s.getOrNewUser(telegramID)
	if err != nil {
		return err
	}
	filtered := user.Mistakes[:0]
	removed := false
	for _, mistake := range user.Mistakes {
		if !removed && normalizeLearningLanguage(mistake.Language) == language && mistake.Word == word && mistake.Correction == correction {
			removed = true
			continue
		}
		filtered = append(filtered, mistake)
	}
	user.Mistakes = filtered
	user.UpdatedAt = time.Now().UTC()
	return s.saveUser(user)
}

func (s *sqliteStore) leaderboard(limit int) []leaderboardEntry {
	if limit <= 0 {
		limit = 10
	}
	rows, err := s.db.Query(`SELECT u.telegram_id,
			COALESCE(CASE WHEN u.telegram_id > 0 THEN NULLIF(u.first_name, '') END, NULLIF(wa.login, ''), NULLIF(u.first_name, ''), ''),
			u.xp, u.learning_language
		FROM users u
		LEFT JOIN web_accounts wa ON wa.user_id = u.telegram_id
		ORDER BY u.telegram_id ASC`)
	if err != nil {
		return nil
	}

	users := []userState{}
	for rows.Next() {
		var user userState
		if err := rows.Scan(&user.TelegramID, &user.FirstName, &user.XP, &user.LearningLanguage); err != nil {
			continue
		}
		users = append(users, user)
	}
	rows.Close()

	entries := make([]leaderboardEntry, 0, limit)
	for _, user := range users {
		mistakes, err := s.getMistakes(user.TelegramID)
		if err != nil {
			continue
		}
		learned, err := s.getLearnedWords(user.TelegramID)
		if err != nil {
			continue
		}
		user.Mistakes = mistakes
		user.LearnedWords = learned
		score := totalLeaderboardScore(user)
		if score <= 0 {
			continue
		}
		name := user.FirstName
		if name == "" {
			name = "Ученик"
		}
		level, title, _, _ := knowledgeLevel(score)
		entries = append(entries, leaderboardEntry{
			FirstName: name,
			Score:     score,
			XP:        user.XP,
			Words:     masteredWordCount(user),
			Mistakes:  len(user.Mistakes),
			Level:     level,
			Title:     title,
			Languages: leaderboardLanguages(user),
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score == entries[j].Score {
			return entries[i].FirstName < entries[j].FirstName
		}
		return entries[i].Score > entries[j].Score
	})
	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries
}

func (s *sqliteStore) languageLeaderboard(language string, limit int) []languageLeaderboardEntry {
	if limit <= 0 {
		limit = 10
	}
	rows, err := s.db.Query(`SELECT u.telegram_id,
			COALESCE(CASE WHEN u.telegram_id > 0 THEN NULLIF(u.first_name, '') END, NULLIF(wa.login, ''), NULLIF(u.first_name, ''), ''),
			u.xp, u.level, u.learning_language
		FROM users u
		LEFT JOIN web_accounts wa ON wa.user_id = u.telegram_id
		ORDER BY u.telegram_id ASC`)
	if err != nil {
		return nil
	}

	language = normalizeLearningLanguage(language)
	users := []userState{}
	for rows.Next() {
		var user userState
		if err := rows.Scan(&user.TelegramID, &user.FirstName, &user.XP, &user.Level, &user.LearningLanguage); err != nil {
			continue
		}
		users = append(users, user)
	}
	rows.Close()

	entries := make([]languageLeaderboardEntry, 0, limit)
	for _, user := range users {
		mistakes, err := s.getMistakes(user.TelegramID)
		if err != nil {
			continue
		}
		learned, err := s.getLearnedWords(user.TelegramID)
		if err != nil {
			continue
		}
		user.Mistakes = mistakes
		user.LearnedWords = learned
		entry := languageLeaderboardEntry{
			FirstName: user.FirstName,
			Words:     len(learnedWordsForLanguageCode(user, language)),
			Mistakes:  len(mistakesForLanguageCode(user, language)),
			Level:     normalizeCEFRLevel(user.Level),
		}
		if entry.FirstName == "" {
			entry.FirstName = "Ученик"
		}
		entry.Score = leaderboardActivityScore(entry.Words, entry.Mistakes)
		if entry.Score > 0 {
			entries = append(entries, entry)
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Score == entries[j].Score {
			return entries[i].FirstName < entries[j].FirstName
		}
		return entries[i].Score > entries[j].Score
	})
	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries
}

func (s *sqliteStore) updateUser(telegramID int64, change func(*userState)) error {
	user, _, err := s.getOrNewUser(telegramID)
	if err != nil {
		return err
	}
	normalizeUser(&user, time.Now().UTC())
	change(&user)
	user.UpdatedAt = time.Now().UTC()
	if err := s.saveUser(user); err != nil {
		return err
	}
	_, err = s.rewardReferralLevel(telegramID)
	return err
}

func (s *sqliteStore) getOrNewUser(telegramID int64) (userState, bool, error) {
	user, ok, err := s.getUser(telegramID)
	if err != nil {
		return userState{}, false, err
	}
	if !ok {
		user = newUserState(telegramID, time.Now().UTC())
	}
	normalizeUser(&user, time.Now().UTC())
	return user, ok, nil
}

func (s *sqliteStore) mustGetUser(telegramID int64) (userState, error) {
	user, ok, err := s.getUser(telegramID)
	if err != nil {
		return userState{}, err
	}
	if !ok {
		return userState{}, errors.New("user not found after save")
	}
	return user, nil
}

func (s *sqliteStore) ensureReferralCode(user *userState) error {
	code, err := normalizeWebReferralCode(user.ReferralCode)
	if err == nil && code != "" {
		user.ReferralCode = code
		return nil
	}
	for attempt := 0; attempt < 24; attempt++ {
		code, err = newWebReferralCode()
		if err != nil {
			return err
		}
		var existingID int64
		err = s.db.QueryRow(`SELECT telegram_id FROM users WHERE referral_code = ?`, code).Scan(&existingID)
		if errors.Is(err, sql.ErrNoRows) {
			user.ReferralCode = code
			return nil
		}
		if err != nil {
			return err
		}
		if existingID == user.TelegramID {
			user.ReferralCode = code
			return nil
		}
	}
	return errors.New("failed to generate unique referral code")
}

func (s *sqliteStore) getUser(telegramID int64) (userState, bool, error) {
	row := s.db.QueryRow(`SELECT telegram_id, first_name, plan, premium_until, invited_by, referral_code, referral_count, referral_balance_kopecks, referral_level_rewarded, mode,
		interface_language, interface_selected, learning_language, language_selected, level, learning_focus,
		lesson_count, practice_count, voice_count, word_lesson_count, word_game_count, xp, daily_date,
		lessons_today, practice_today, voice_today, lesson_history, practice_history, last_payment_charge_id, last_lesson_prompt,
		reminder_enabled, reminder_utc_offset, reminder_hour, timezone_selected, last_reminder_date, created_at, updated_at
		FROM users WHERE telegram_id = ?`, telegramID)
	var user userState
	var premiumUntil, createdAt, updatedAt, lessonHistoryJSON, practiceHistoryJSON string
	var reminderEnabled, timezoneSelected, interfaceSelected, languageSelected, referralLevelRewarded int
	err := row.Scan(&user.TelegramID, &user.FirstName, &user.Plan, &premiumUntil, &user.InvitedBy, &user.ReferralCode, &user.ReferralCount, &user.ReferralBalanceKopecks,
		&referralLevelRewarded, &user.Mode, &user.InterfaceLanguage, &interfaceSelected, &user.LearningLanguage, &languageSelected, &user.Level, &user.LearningFocus, &user.LessonCount, &user.PracticeCount, &user.VoiceCount, &user.WordLessonCount,
		&user.WordGameCount, &user.XP, &user.DailyDate, &user.LessonsToday, &user.PracticeToday, &user.VoiceToday,
		&lessonHistoryJSON, &practiceHistoryJSON, &user.LastPaymentChargeID, &user.LastLessonPrompt, &reminderEnabled, &user.ReminderUTCOffset, &user.ReminderHour,
		&timezoneSelected, &user.LastReminderDate, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return userState{}, false, nil
	}
	if err != nil {
		return userState{}, false, err
	}
	user.PremiumUntil = parseDBTime(premiumUntil)
	user.InterfaceSelected = interfaceSelected != 0
	user.LanguageSelected = languageSelected != 0
	user.ReferralLevelRewarded = referralLevelRewarded != 0
	user.ReminderEnabled = reminderEnabled != 0
	user.TimezoneSelected = timezoneSelected != 0
	user.LessonHistory = trimLessonHistory(parseStringListJSON(lessonHistoryJSON))
	user.PracticeHistory = parseStringListJSON(practiceHistoryJSON)
	user.CreatedAt = parseDBTime(createdAt)
	user.UpdatedAt = parseDBTime(updatedAt)
	user.Mistakes, err = s.getMistakes(telegramID)
	if err != nil {
		return userState{}, false, err
	}
	user.LearnedWords, err = s.getLearnedWords(telegramID)
	if err != nil {
		return userState{}, false, err
	}
	user.Phrasebook, err = s.getPhrasebookEntries(telegramID)
	if err != nil {
		return userState{}, false, err
	}
	user.HabitLog, err = s.getHabitLog(telegramID)
	if err != nil {
		return userState{}, false, err
	}
	user.NavigationLayout, err = s.getNavigationLayout(telegramID)
	if err != nil {
		return userState{}, false, err
	}
	normalizeUser(&user, time.Now().UTC())
	return user, true, nil
}

func (s *sqliteStore) saveUser(user userState) error {
	if err := s.ensureReferralCode(&user); err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`INSERT INTO users (
		telegram_id, first_name, plan, premium_until, invited_by, referral_code, referral_count, referral_balance_kopecks, referral_level_rewarded, mode, interface_language, interface_selected, learning_language, language_selected, level, learning_focus,
		lesson_count, practice_count, voice_count, word_lesson_count, word_game_count, xp, daily_date,
		lessons_today, practice_today, voice_today, lesson_history, practice_history, last_payment_charge_id, last_lesson_prompt,
		reminder_enabled, reminder_utc_offset, reminder_hour, timezone_selected, last_reminder_date, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(telegram_id) DO UPDATE SET
		first_name=excluded.first_name, plan=excluded.plan, premium_until=excluded.premium_until,
		invited_by=excluded.invited_by, referral_code=excluded.referral_code, referral_count=excluded.referral_count,
		referral_balance_kopecks=MAX(users.referral_balance_kopecks, excluded.referral_balance_kopecks),
		referral_level_rewarded=excluded.referral_level_rewarded, mode=excluded.mode,
		interface_language=excluded.interface_language, interface_selected=excluded.interface_selected,
		learning_language=excluded.learning_language, language_selected=excluded.language_selected, level=excluded.level,
		learning_focus=excluded.learning_focus,
		lesson_count=excluded.lesson_count, practice_count=excluded.practice_count, voice_count=excluded.voice_count,
		word_lesson_count=excluded.word_lesson_count, word_game_count=excluded.word_game_count, xp=excluded.xp,
		daily_date=excluded.daily_date, lessons_today=excluded.lessons_today, practice_today=excluded.practice_today,
		voice_today=excluded.voice_today, lesson_history=excluded.lesson_history, practice_history=excluded.practice_history,
		last_payment_charge_id=excluded.last_payment_charge_id,
		last_lesson_prompt=excluded.last_lesson_prompt, reminder_enabled=excluded.reminder_enabled,
		reminder_utc_offset=excluded.reminder_utc_offset, reminder_hour=excluded.reminder_hour,
		timezone_selected=excluded.timezone_selected, last_reminder_date=excluded.last_reminder_date,
		updated_at=excluded.updated_at`,
		user.TelegramID, user.FirstName, user.Plan, formatDBTime(user.PremiumUntil), user.InvitedBy, user.ReferralCode, user.ReferralCount, user.ReferralBalanceKopecks, boolToInt(user.ReferralLevelRewarded),
		user.Mode, normalizeInterfaceLanguage(user.InterfaceLanguage), boolToInt(user.InterfaceSelected), normalizeLearningLanguage(user.LearningLanguage), boolToInt(user.LanguageSelected), user.Level, strings.TrimSpace(user.LearningFocus), user.LessonCount, user.PracticeCount, user.VoiceCount, user.WordLessonCount, user.WordGameCount,
		user.XP, user.DailyDate, user.LessonsToday, user.PracticeToday, user.VoiceToday, stringListJSON(trimLessonHistory(user.LessonHistory)), stringListJSON(user.PracticeHistory),
		user.LastPaymentChargeID, user.LastLessonPrompt, boolToInt(user.ReminderEnabled), user.ReminderUTCOffset, user.ReminderHour,
		boolToInt(user.TimezoneSelected), user.LastReminderDate, formatDBTime(user.CreatedAt), formatDBTime(user.UpdatedAt))
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM mistakes WHERE telegram_id = ?`, user.TelegramID); err != nil {
		return err
	}
	for _, item := range user.Mistakes {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO mistakes (telegram_id, word, language, correction, explanation, added_at) VALUES (?, ?, ?, ?, ?, ?)`,
			user.TelegramID, item.Word, normalizeLearningLanguage(item.Language), item.Correction, item.Explanation, formatDBTime(item.AddedAt)); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`DELETE FROM learned_words WHERE telegram_id = ?`, user.TelegramID); err != nil {
		return err
	}
	for _, item := range user.LearnedWords {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO learned_words (telegram_id, id, language, russian, english, context, review_correct_count, spelling_correct_count, learned_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			user.TelegramID, item.ID, normalizeLearningLanguage(item.Language), item.Russian, item.English, item.Context, item.ReviewCorrectCount, item.SpellingCorrectCount, formatDBTime(item.LearnedAt)); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`DELETE FROM phrasebook_entries WHERE telegram_id = ?`, user.TelegramID); err != nil {
		return err
	}
	for _, item := range normalizePhrasebookEntries(user.Phrasebook, user.LearningLanguage, time.Now().UTC()) {
		if _, err := tx.Exec(`INSERT OR REPLACE INTO phrasebook_entries (telegram_id, id, phrase, translation, note, source, language, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			user.TelegramID, item.ID, item.Phrase, item.Translation, item.Note, normalizePhrasebookSource(item.Source), normalizeLearningLanguage(item.Language), formatDBTime(item.CreatedAt)); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`DELETE FROM habit_days WHERE telegram_id = ?`, user.TelegramID); err != nil {
		return err
	}
	for date, item := range normalizeHabitLog(user.HabitLog, time.Now().UTC()) {
		if _, err := tx.Exec(`INSERT OR REPLACE INTO habit_days (telegram_id, habit_date, login, complete, claimed, claimed_at, lessons, practice, voice, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			user.TelegramID, date, boolToInt(item.Login), boolToInt(item.Complete), boolToInt(item.Claimed), formatDBTime(item.ClaimedAt), item.Lessons, item.Practice, item.Voice, formatDBTime(item.UpdatedAt)); err != nil {
			return err
		}
	}
	if err := saveNavigationLayoutTx(tx, user.TelegramID, user.NavigationLayout); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *sqliteStore) getMistakes(telegramID int64) ([]mistakeEntry, error) {
	rows, err := s.db.Query(`SELECT word, language, correction, explanation, added_at FROM mistakes WHERE telegram_id = ? ORDER BY added_at ASC`, telegramID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []mistakeEntry
	for rows.Next() {
		var item mistakeEntry
		var addedAt string
		if err := rows.Scan(&item.Word, &item.Language, &item.Correction, &item.Explanation, &addedAt); err != nil {
			return nil, err
		}
		item.AddedAt = parseDBTime(addedAt)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *sqliteStore) getLearnedWords(telegramID int64) ([]learnedWordEntry, error) {
	rows, err := s.db.Query(`SELECT id, language, russian, english, context, review_correct_count, spelling_correct_count, learned_at FROM learned_words WHERE telegram_id = ? ORDER BY learned_at ASC`, telegramID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []learnedWordEntry
	for rows.Next() {
		var item learnedWordEntry
		var learnedAt string
		if err := rows.Scan(&item.ID, &item.Language, &item.Russian, &item.English, &item.Context, &item.ReviewCorrectCount, &item.SpellingCorrectCount, &learnedAt); err != nil {
			return nil, err
		}
		item.LearnedAt = parseDBTime(learnedAt)
		items = append(items, item)
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].LearnedAt.Before(items[j].LearnedAt) })
	return items, rows.Err()
}

func (s *sqliteStore) getPhrasebookEntries(telegramID int64) ([]phrasebookEntry, error) {
	rows, err := s.db.Query(`SELECT id, phrase, translation, note, source, language, created_at FROM phrasebook_entries WHERE telegram_id = ? ORDER BY created_at DESC LIMIT 80`, telegramID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []phrasebookEntry
	for rows.Next() {
		var item phrasebookEntry
		var createdAt string
		if err := rows.Scan(&item.ID, &item.Phrase, &item.Translation, &item.Note, &item.Source, &item.Language, &createdAt); err != nil {
			return nil, err
		}
		item.CreatedAt = parseDBTime(createdAt)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return normalizePhrasebookEntries(items, "", time.Now().UTC()), nil
}

func (s *sqliteStore) getHabitLog(telegramID int64) (map[string]habitDayEntry, error) {
	cutoff := time.Now().UTC().AddDate(0, 0, -364).Format("2006-01-02")
	rows, err := s.db.Query(`SELECT habit_date, login, complete, claimed, claimed_at, lessons, practice, voice, updated_at FROM habit_days WHERE telegram_id = ? AND habit_date >= ? ORDER BY habit_date ASC`, telegramID, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := map[string]habitDayEntry{}
	for rows.Next() {
		var date, claimedAt, updatedAt string
		var login, complete, claimed int
		var item habitDayEntry
		if err := rows.Scan(&date, &login, &complete, &claimed, &claimedAt, &item.Lessons, &item.Practice, &item.Voice, &updatedAt); err != nil {
			return nil, err
		}
		item.Date = date
		item.Login = login != 0
		item.Complete = complete != 0
		item.Claimed = claimed != 0
		item.ClaimedAt = parseDBTime(claimedAt)
		item.UpdatedAt = parseDBTime(updatedAt)
		items[date] = item
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return normalizeHabitLog(items, time.Now().UTC()), nil
}

func (s *sqliteStore) getNavigationLayout(telegramID int64) (navigationLayout, error) {
	row := s.db.QueryRow(`SELECT function_ribbon, mobile_pinned, mobile_more, mobile_rail FROM navigation_layouts WHERE telegram_id = ?`, telegramID)
	var functionRibbonJSON, mobilePinnedJSON, mobileMoreJSON, mobileRailJSON string
	if err := row.Scan(&functionRibbonJSON, &mobilePinnedJSON, &mobileMoreJSON, &mobileRailJSON); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return navigationLayout{}, nil
		}
		return navigationLayout{}, err
	}
	return normalizeNavigationLayout(navigationLayout{
		FunctionRibbon: parseStringListJSON(functionRibbonJSON),
		MobilePinned:   parseStringListJSON(mobilePinnedJSON),
		MobileMore:     parseStringListJSON(mobileMoreJSON),
		MobileRail:     parseStringListJSON(mobileRailJSON),
	}), nil
}

func saveNavigationLayoutTx(tx *sql.Tx, telegramID int64, layout navigationLayout) error {
	layout = normalizeNavigationLayout(layout)
	_, err := tx.Exec(
		`INSERT INTO navigation_layouts (telegram_id, function_ribbon, mobile_pinned, mobile_more, mobile_rail, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(telegram_id) DO UPDATE SET
			function_ribbon=excluded.function_ribbon,
			mobile_pinned=excluded.mobile_pinned,
			mobile_more=excluded.mobile_more,
			mobile_rail=excluded.mobile_rail,
			updated_at=excluded.updated_at`,
		telegramID,
		stringListJSON(layout.FunctionRibbon),
		stringListJSON(layout.MobilePinned),
		stringListJSON(layout.MobileMore),
		stringListJSON(layout.MobileRail),
		formatDBTime(time.Now().UTC()),
	)
	return err
}

func (s *sqliteStore) addPhrasebookEntry(telegramID int64, entry phrasebookEntry) ([]phrasebookEntry, error) {
	now := time.Now().UTC()
	entries := normalizePhrasebookEntries([]phrasebookEntry{entry}, "en", now)
	if len(entries) == 0 {
		return s.getPhrasebookEntries(telegramID)
	}
	entry = entries[0]
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM phrasebook_entries WHERE telegram_id = ? AND lower(phrase) = lower(?)`, telegramID, entry.Phrase); err != nil {
		return nil, err
	}
	if _, err := tx.Exec(`INSERT OR REPLACE INTO phrasebook_entries (telegram_id, id, phrase, translation, note, source, language, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		telegramID, entry.ID, entry.Phrase, entry.Translation, entry.Note, normalizePhrasebookSource(entry.Source), normalizeLearningLanguage(entry.Language), formatDBTime(entry.CreatedAt)); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return s.getPhrasebookEntries(telegramID)
}

func (s *sqliteStore) removePhrasebookEntry(telegramID int64, id string) ([]phrasebookEntry, error) {
	if _, err := s.db.Exec(`DELETE FROM phrasebook_entries WHERE telegram_id = ? AND id = ?`, telegramID, strings.TrimSpace(id)); err != nil {
		return nil, err
	}
	return s.getPhrasebookEntries(telegramID)
}

func formatDBTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339Nano)
}

func parseDBTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func stringListJSON(values []string) string {
	if len(values) == 0 {
		return ""
	}
	bytes, err := json.Marshal(values)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func parseStringListJSON(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	var values []string
	if err := json.Unmarshal([]byte(value), &values); err != nil {
		return nil
	}
	return values
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
