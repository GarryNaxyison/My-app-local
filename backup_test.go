package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSQLiteBackupCreatesFile(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.sqlite")
	backupDir := filepath.Join(dir, "backups")

	store, err := newSQLiteStore(dbPath, filepath.Join(dir, "missing.json"))
	if err != nil {
		t.Fatalf("newSQLiteStore: %v", err)
	}
	t.Cleanup(func() { _ = store.db.Close() })

	now := time.Date(2026, 5, 16, 3, 0, 0, 0, time.Local)
	if err := store.backup(context.Background(), dbPath, backupDir, now); err != nil {
		t.Fatalf("backup: %v", err)
	}

	backupPath := filepath.Join(backupDir, "test-2026-05-16.sqlite")
	info, err := os.Stat(backupPath)
	if err != nil {
		t.Fatalf("stat backup: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("backup file is empty")
	}
}
