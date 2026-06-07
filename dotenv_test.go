package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvTrimsUTF8BOM(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(path, []byte("\xef\xbb\xbfTELEGRAM_BOT_TOKEN=123:abc\n"), 0o600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	oldValue, hadValue := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if err := os.Unsetenv("TELEGRAM_BOT_TOKEN"); err != nil {
		t.Fatalf("unset env: %v", err)
	}
	t.Cleanup(func() {
		if hadValue {
			_ = os.Setenv("TELEGRAM_BOT_TOKEN", oldValue)
			return
		}
		_ = os.Unsetenv("TELEGRAM_BOT_TOKEN")
	})

	if err := loadDotEnv(path); err != nil {
		t.Fatalf("load env: %v", err)
	}
	if got := os.Getenv("TELEGRAM_BOT_TOKEN"); got != "123:abc" {
		t.Fatalf("expected token from BOM env, got %q", got)
	}
}
