package main

import (
	"path/filepath"
	"testing"
	"time"
)

func TestJSONStoreDailyBonusRequiresTwentyFourHours(t *testing.T) {
	store, err := newJSONStore(filepath.Join(t.TempDir(), "users.json"))
	if err != nil {
		t.Fatalf("newJSONStore() error = %v", err)
	}
	user, err := store.getOrCreateUser(12345, "Daily")
	if err != nil {
		t.Fatalf("getOrCreateUser() error = %v", err)
	}
	if claimed, err := store.claimDailyBonus(user.TelegramID, "2026-05-27", 25); err != nil || !claimed {
		t.Fatalf("first claim = %v, %v; want claimed", claimed, err)
	}
	if claimed, err := store.claimDailyBonus(user.TelegramID, "2026-05-28", 25); err != nil || claimed {
		t.Fatalf("second claim inside 24h = %v, %v; want locked", claimed, err)
	}
	if err := store.update(user.TelegramID, func(current *userState) {
		current.DailyBonusLastClaimedAt = time.Now().UTC().Add(-25 * time.Hour)
	}); err != nil {
		t.Fatalf("store.update() error = %v", err)
	}
	if claimed, err := store.claimDailyBonus(user.TelegramID, "2026-05-28", 25); err != nil || !claimed {
		t.Fatalf("claim after 24h = %v, %v; want claimed", claimed, err)
	}
}
