package main

import (
	"context"
	"testing"
	"time"
)

func TestPremiumTransferExperimentVariant(t *testing.T) {
	if got := premiumTransferVariant(false, false, 7); got != "" {
		t.Fatalf("disabled experiment assigned %q", got)
	}
	if got := premiumTransferVariant(true, true, 7); got != "" {
		t.Fatalf("paid account assigned %q", got)
	}
	groups := map[string]int{}
	for id := int64(-100); id <= 100; id++ {
		if id == 0 {
			continue
		}
		got := premiumTransferVariant(true, false, id)
		if got != premiumTransferVariant(true, false, id) {
			t.Fatalf("account %d changed variants", id)
		}
		groups[got]++
	}
	if len(groups) != 2 || groups["control"] < 60 || groups["offer"] < 60 {
		t.Fatalf("unexpected allocation: %v", groups)
	}
}

func TestPremiumTransferEnrollmentIsIdempotent(t *testing.T) {
	api, storage, _ := newTestWebAPI(t)
	user := userState{TelegramID: -42}
	now := time.Now()
	api.recordPremiumTransferEnrollment(context.Background(), user, now.Add(-time.Hour))
	api.cfg.PremiumTransferExperiment = true
	api.recordPremiumTransferEnrollment(context.Background(), user, now)
	api.recordPremiumTransferEnrollment(context.Background(), user, now.Add(time.Hour))
	var count int
	var enrolledAt string
	err := storage.db.QueryRow(`SELECT COUNT(*), COALESCE(MIN(created_at), '') FROM learning_events
		WHERE telegram_id=? AND event_id='premium-transfer-v1' AND stage=?`,
		user.TelegramID, "premium_transfer_"+premiumTransferVariant(true, false, user.TelegramID)).Scan(&count, &enrolledAt)
	if err != nil || count != 1 || enrolledAt != formatDBTime(now) {
		t.Fatalf("count=%d enrolled=%q error=%v", count, enrolledAt, err)
	}
}
