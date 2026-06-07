package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (b *bot) runSQLiteBackupScheduler(ctx context.Context) {
	if strings.ToLower(strings.TrimSpace(b.cfg.StorageDriver)) == "json" {
		return
	}
	sqlite, ok := b.store.(*sqliteStore)
	if !ok {
		log.Printf("sqlite backup disabled: store is not sqlite")
		return
	}
	if strings.TrimSpace(b.cfg.SQLiteBackupDir) == "" {
		log.Printf("sqlite backup disabled: SQLITE_BACKUP_DIR is empty")
		return
	}

	b.runDueSQLiteBackup(ctx, sqlite)
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			b.runDueSQLiteBackup(ctx, sqlite)
		}
	}
}

func (b *bot) runDueSQLiteBackup(ctx context.Context, sqlite *sqliteStore) {
	now := time.Now()
	if now.Hour() != b.cfg.SQLiteBackupHour {
		return
	}
	markerPath := filepath.Join(b.cfg.SQLiteBackupDir, ".last_backup_date")
	today := now.Format("2006-01-02")
	if lastBackupDate(markerPath) == today {
		return
	}

	if err := sqlite.backup(ctx, b.cfg.DatabasePath, b.cfg.SQLiteBackupDir, now); err != nil {
		log.Printf("sqlite backup failed: %v", err)
		return
	}
	if err := os.WriteFile(markerPath, []byte(today), 0o644); err != nil {
		log.Printf("sqlite backup marker failed: %v", err)
	}
	if err := cleanupSQLiteBackups(b.cfg.SQLiteBackupDir, b.cfg.SQLiteBackupRetention, now); err != nil {
		log.Printf("sqlite backup cleanup failed: %v", err)
	}
}

func (s *sqliteStore) backup(ctx context.Context, databasePath string, backupDir string, now time.Time) error {
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}
	name := sqliteBackupName(databasePath, now)
	destination := filepath.Join(backupDir, name)
	if _, err := os.Stat(destination); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	escaped := strings.ReplaceAll(destination, "'", "''")
	if _, err := s.db.ExecContext(ctx, "VACUUM INTO '"+escaped+"'"); err != nil {
		return err
	}
	log.Printf("sqlite backup created: %s", destination)
	return nil
}

func sqliteBackupName(databasePath string, now time.Time) string {
	base := filepath.Base(databasePath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	if name == "" {
		name = "sqlite"
	}
	if ext == "" {
		ext = ".sqlite"
	}
	return fmt.Sprintf("%s-%s%s", name, now.Format("2006-01-02"), ext)
}

func lastBackupDate(markerPath string) string {
	bytes, err := os.ReadFile(markerPath)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(bytes))
}

func cleanupSQLiteBackups(backupDir string, retentionDays int, now time.Time) error {
	if retentionDays <= 0 {
		return nil
	}
	cutoff := now.AddDate(0, 0, -retentionDays)
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sqlite") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.ModTime().Before(cutoff) {
			if err := os.Remove(filepath.Join(backupDir, entry.Name())); err != nil {
				return err
			}
		}
	}
	return nil
}
