package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestActivationKeyRedeemIsOneTimeEvenWhenImportFileStillContainsKey(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "activation_keys.sqlite")
	importPath := filepath.Join(dir, "activation_keys.txt")
	content := []byte("ABCD-EFGH-JKLM-NPQR month premium launch\nWXYZ-2345-6789-ABCD year premium annual\n")
	if err := os.WriteFile(importPath, content, 0o600); err != nil {
		t.Fatalf("write import file: %v", err)
	}

	store, err := newActivationKeyStore(dbPath, importPath)
	if err != nil {
		t.Fatalf("newActivationKeyStore: %v", err)
	}
	grant, err := store.redeem(context.Background(), "abcd efgh jklm npqr", 42)
	if err != nil {
		t.Fatalf("redeem month key: %v", err)
	}
	if grant.Key != "ABCD-EFGH-JKLM-NPQR" || grant.DurationDays != 30 || grant.Tier != "premium" {
		t.Fatalf("unexpected grant: %+v", grant)
	}
	if _, err := store.redeem(context.Background(), "ABCD-EFGH-JKLM-NPQR", 42); !errors.Is(err, errActivationKeyUsed) {
		t.Fatalf("second redeem error = %v, want errActivationKeyUsed", err)
	}
	if err := store.close(); err != nil {
		t.Fatalf("close store: %v", err)
	}

	reopened, err := newActivationKeyStore(dbPath, importPath)
	if err != nil {
		t.Fatalf("reopen activation store: %v", err)
	}
	defer reopened.close()
	if _, err := reopened.redeem(context.Background(), "ABCD-EFGH-JKLM-NPQR", 42); !errors.Is(err, errActivationKeyUsed) {
		t.Fatalf("redeem imported used key after reopen = %v, want errActivationKeyUsed", err)
	}
	yearGrant, err := reopened.redeem(context.Background(), "WXYZ-2345-6789-ABCD", 42)
	if err != nil {
		t.Fatalf("redeem year key: %v", err)
	}
	if yearGrant.DurationDays != 365 {
		t.Fatalf("year key duration = %d, want 365", yearGrant.DurationDays)
	}
	if _, err := reopened.redeem(context.Background(), "ZZZZ-ZZZZ-ZZZZ-ZZZZ", 42); !errors.Is(err, errActivationKeyMissing) {
		t.Fatalf("unknown key error = %v, want errActivationKeyMissing", err)
	}
}

func TestActivationKeyValidation(t *testing.T) {
	if _, err := normalizeActivationKey("ABCD-EFGH-JKLM-NPQR"); err != nil {
		t.Fatalf("valid key rejected: %v", err)
	}
	grant, ok, err := parseActivationKeyLine("WXYZ-2345-6789-ABCD platinum_year platinum vip")
	if err != nil || !ok {
		t.Fatalf("parse platinum key: ok=%v err=%v", ok, err)
	}
	if grant.DurationDays != 365 || grant.Tier != "platinum" || grant.Note != "vip" {
		t.Fatalf("unexpected platinum grant: %+v", grant)
	}
	if _, err := normalizeActivationKey("ABCD-EFGH-IJKL-NPQR"); !errors.Is(err, errActivationKeyInvalid) {
		t.Fatalf("ambiguous letter error = %v, want errActivationKeyInvalid", err)
	}
	if _, err := parseActivationDurationDays("90"); err == nil {
		t.Fatal("90 days accepted, want error")
	}
}
