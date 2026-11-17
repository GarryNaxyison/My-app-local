package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"time"
)

func premiumTransferVariant(enabled, premium bool, accountID int64) string {
	if !enabled || premium || accountID == 0 {
		return ""
	}
	var input [8]byte
	binary.BigEndian.PutUint64(input[:], uint64(accountID))
	digest := sha256.Sum256(input[:])
	if digest[0]&1 == 0 {
		return "control"
	}
	return "offer"
}

func (api *webAPI) recordPremiumTransferEnrollment(ctx context.Context, user userState, now time.Time) {
	variant := premiumTransferVariant(api.cfg.PremiumTransferExperiment, user.isPremium(now), user.TelegramID)
	storage, ok := api.bot.store.(*sqliteStore)
	if !ok || variant == "" {
		return
	}
	_, err := storage.db.ExecContext(ctx, `INSERT INTO learning_events
		(telegram_id, event_id, stage, format, elapsed_ms, created_at)
		VALUES (?, 'premium-transfer-v1', ?, '', 0, ?)
		ON CONFLICT(telegram_id, event_id) DO NOTHING`,
		user.TelegramID, "premium_transfer_"+variant, formatDBTime(now))
	if err != nil {
		log.Printf("Could not record Premium experiment enrollment: %v", err)
	}
}

func recordPremiumTransferPurchase(ctx context.Context, storage store, userID int64, now time.Time) {
	db, ok := storage.(*sqliteStore)
	if !ok {
		return
	}
	_, err := db.db.ExecContext(ctx, `INSERT INTO learning_events
		(telegram_id, event_id, stage, format, elapsed_ms, created_at)
		SELECT ?, 'premium-transfer-purchase-v1', 'premium_transfer_purchased', '', 0, ?
		WHERE EXISTS (SELECT 1 FROM learning_events WHERE telegram_id=? AND event_id='premium-transfer-v1')
		ON CONFLICT(telegram_id, event_id) DO NOTHING`, userID, formatDBTime(now), userID)
	if err != nil {
		log.Printf("Could not record Premium experiment purchase: %v", err)
	}
}
