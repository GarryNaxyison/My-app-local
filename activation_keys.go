package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const activationKeyAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

var (
	errActivationKeyInvalid = errors.New("activation key must use format XXXX-XXXX-XXXX-XXXX")
	errActivationKeyMissing = errors.New("activation key not found")
	errActivationKeyUsed    = errors.New("activation key already activated")
)

type activationKeyStore struct {
	db         *sql.DB
	importPath string
}

type activationKeyGrant struct {
	Key          string
	DurationDays int
	Tier         string
	Note         string
}

func newActivationKeyStore(databasePath string, importPath string) (*activationKeyStore, error) {
	db, err := openSQLiteDatabase("ACTIVATION_KEYS_DATABASE_PATH", databasePath)
	if err != nil {
		return nil, err
	}
	store := &activationKeyStore{db: db, importPath: strings.TrimSpace(importPath)}
	if err := store.init(); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := store.importFile(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *activationKeyStore) close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *activationKeyStore) init() error {
	statements := []string{
		`PRAGMA journal_mode=WAL`,
		`PRAGMA synchronous=NORMAL`,
		`PRAGMA busy_timeout=5000`,
		`CREATE TABLE IF NOT EXISTS activation_keys (
			key TEXT PRIMARY KEY,
			duration_days INTEGER NOT NULL,
			tier TEXT NOT NULL DEFAULT 'premium',
			note TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS activation_key_redemptions (
			key TEXT PRIMARY KEY,
			telegram_id INTEGER NOT NULL,
			duration_days INTEGER NOT NULL,
			tier TEXT NOT NULL,
			redeemed_at TEXT NOT NULL
		)`,
	}
	for _, stmt := range statements {
		if _, err := s.db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func normalizeActivationKey(raw string) (string, error) {
	text := strings.ToUpper(strings.TrimSpace(raw))
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\t", "")
	text = strings.ReplaceAll(text, "-", "")
	if len(text) != 16 {
		return "", errActivationKeyInvalid
	}
	for _, ch := range text {
		if !strings.ContainsRune(activationKeyAlphabet, ch) {
			return "", errActivationKeyInvalid
		}
	}
	return text[0:4] + "-" + text[4:8] + "-" + text[8:12] + "-" + text[12:16], nil
}

func parseActivationKeyLine(line string) (activationKeyGrant, bool, error) {
	if before, _, ok := strings.Cut(line, "#"); ok {
		line = before
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return activationKeyGrant{}, false, nil
	}
	fields := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == '\t' || r == ',' || r == ';' || r == '='
	})
	if len(fields) < 2 {
		return activationKeyGrant{}, false, fmt.Errorf("activation key line must contain key and duration: %q", line)
	}
	key, err := normalizeActivationKey(fields[0])
	if err != nil {
		return activationKeyGrant{}, false, err
	}
	days, err := parseActivationDurationDays(fields[1])
	if err != nil {
		return activationKeyGrant{}, false, err
	}
	tier := "premium"
	if len(fields) >= 3 {
		tier = paidTierForActivation(fields[2])
	}
	note := ""
	if len(fields) > 3 {
		note = strings.Join(fields[3:], " ")
	}
	return activationKeyGrant{Key: key, DurationDays: days, Tier: tier, Note: note}, true, nil
}

func parseActivationDurationDays(value string) (int, error) {
	text := strings.ToLower(strings.TrimSpace(value))
	text = strings.TrimSuffix(text, "days")
	text = strings.TrimSuffix(text, "day")
	text = strings.TrimSuffix(text, "d")
	switch text {
	case "month", "monthly", "premium_month", "premium_30", "premium_30d", "platinum_month", "platinum_30", "platinum_30d", "1m":
		return 30, nil
	case "year", "yearly", "premium_year", "premium_365", "premium_365d", "platinum_year", "platinum_365", "platinum_365d", "1y":
		return 365, nil
	}
	days, err := strconv.Atoi(text)
	if err != nil || (days != 30 && days != 365) {
		return 0, fmt.Errorf("activation duration must be 30 or 365 days")
	}
	return days, nil
}

func paidTierForActivation(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "platinum":
		return "platinum"
	default:
		return "premium"
	}
}

func (s *activationKeyStore) importFile(ctx context.Context) error {
	if s == nil || strings.TrimSpace(s.importPath) == "" {
		return nil
	}
	bytes, err := os.ReadFile(s.importPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	lines := strings.Split(string(bytes), "\n")
	grants := make([]activationKeyGrant, 0, len(lines))
	for index, line := range lines {
		grant, ok, err := parseActivationKeyLine(line)
		if err != nil {
			return fmt.Errorf("%s:%d: %w", s.importPath, index+1, err)
		}
		if ok {
			grants = append(grants, grant)
		}
	}
	return s.upsertActivationKeys(ctx, grants)
}

func (s *activationKeyStore) upsertActivationKeys(ctx context.Context, grants []activationKeyGrant) error {
	if s == nil || len(grants) == 0 {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := formatDBTime(time.Now().UTC())
	for _, grant := range grants {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO activation_keys (key, duration_days, tier, note, created_at, updated_at)
			SELECT ?, ?, ?, ?, ?, ?
			WHERE NOT EXISTS (SELECT 1 FROM activation_key_redemptions WHERE key = ?)
			ON CONFLICT(key) DO UPDATE SET
				duration_days = excluded.duration_days,
				tier = excluded.tier,
				note = excluded.note,
				updated_at = excluded.updated_at`,
			grant.Key, grant.DurationDays, grant.Tier, grant.Note, now, now, grant.Key)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *activationKeyStore) redeem(ctx context.Context, rawKey string, telegramID int64) (activationKeyGrant, error) {
	if s == nil {
		return activationKeyGrant{}, errActivationKeyMissing
	}
	key, err := normalizeActivationKey(rawKey)
	if err != nil {
		return activationKeyGrant{}, err
	}
	if err := s.importFile(ctx); err != nil {
		return activationKeyGrant{}, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return activationKeyGrant{}, err
	}
	defer tx.Rollback()
	grant := activationKeyGrant{Key: key}
	row := tx.QueryRowContext(ctx, `SELECT duration_days, tier, note FROM activation_keys WHERE key = ?`, key)
	if err := row.Scan(&grant.DurationDays, &grant.Tier, &grant.Note); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			var redeemed int
			redemptionRow := tx.QueryRowContext(ctx, `SELECT 1 FROM activation_key_redemptions WHERE key = ?`, key)
			if scanErr := redemptionRow.Scan(&redeemed); scanErr == nil {
				return activationKeyGrant{}, errActivationKeyUsed
			} else if !errors.Is(scanErr, sql.ErrNoRows) {
				return activationKeyGrant{}, scanErr
			}
			return activationKeyGrant{}, errActivationKeyMissing
		}
		return activationKeyGrant{}, err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM activation_keys WHERE key = ?`, key); err != nil {
		return activationKeyGrant{}, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO activation_key_redemptions (key, telegram_id, duration_days, tier, redeemed_at)
		VALUES (?, ?, ?, ?, ?)`,
		key, telegramID, grant.DurationDays, grant.Tier, formatDBTime(time.Now().UTC()))
	if err != nil {
		return activationKeyGrant{}, err
	}
	return grant, tx.Commit()
}
