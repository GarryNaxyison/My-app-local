package main

import (
	"path/filepath"
	"testing"
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
