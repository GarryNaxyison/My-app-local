package main

import (
	"net/http"
	"strings"
	"time"
)

type mistakeTransferRecord struct {
	Language    string `json:"language"`
	CompletedAt string `json:"completed_at,omitempty"`
}

type mistakeTransferProgress struct {
	Completed          int `json:"completed"`
	CompletedLast7Days int `json:"completed_last_7_days"`
}

func (p *mistakeTransferProgress) count(stamp string, cutoff time.Time) {
	completedAt := parseDBTime(stamp)
	if completedAt.IsZero() {
		return
	}
	p.Completed++
	if !completedAt.Before(cutoff) {
		p.CompletedLast7Days++
	}
}

func (s *sqliteStore) mistakeTransferProgress(userID int64, language string, cutoff time.Time) (mistakeTransferProgress, error) {
	var progress mistakeTransferProgress
	rows, err := s.db.Query(`SELECT completed_at FROM mistake_transfers
		WHERE telegram_id=? AND language=? AND completed_at<>''`, userID, language)
	if err != nil {
		return progress, err
	}
	defer rows.Close()
	for rows.Next() {
		var stamp string
		if err := rows.Scan(&stamp); err != nil {
			return progress, err
		}
		progress.count(stamp, cutoff)
	}
	return progress, rows.Err()
}

func (s *sqliteStore) finishMistakeTransfer(userID int64, language, id string) (bool, error) {
	result, err := s.db.Exec(`UPDATE mistake_transfers
		SET completed_at=CASE WHEN completed_at='' THEN ? ELSE completed_at END
		WHERE telegram_id=? AND language=? AND review_id=?`, formatDBTime(time.Now().UTC()), userID, language, id)
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count > 0, err
}

func (s *jsonStore) mistakeTransferProgress(userID int64, language string, cutoff time.Time) (mistakeTransferProgress, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var progress mistakeTransferProgress
	for _, item := range s.users[userID].MistakeTransfers {
		if item.Language == language {
			progress.count(item.CompletedAt, cutoff)
		}
	}
	return progress, nil
}

func (s *jsonStore) finishMistakeTransfer(userID int64, language, id string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user := s.users[userID]
	item, exists := user.MistakeTransfers[id]
	if !exists || item.Language != language {
		return false, nil
	}
	if item.CompletedAt != "" {
		return true, nil
	}
	item.CompletedAt = formatDBTime(time.Now().UTC())
	user.MistakeTransfers[id] = item
	if err := s.saveLocked(); err != nil {
		item.CompletedAt = ""
		user.MistakeTransfers[id] = item
		return false, err
	}
	return true, nil
}

func (api *webAPI) handleMistakeTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	language := normalizeLearningLanguage(user.LearningLanguage)
	if r.Method == http.MethodPost {
		var req struct {
			ReviewID string `json:"review_id"`
		}
		if !decodeJSONRequest(w, r, &req) {
			return
		}
		id := strings.TrimSpace(req.ReviewID)
		if !strings.HasPrefix(id, "mistake:") || len(id) != len("mistake:")+64 {
			writeAPIError(w, http.StatusBadRequest, "Invalid review ID.")
			return
		}
		found, err := api.bot.store.finishMistakeTransfer(user.TelegramID, language, id)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, "Unable to save practice progress.")
			return
		}
		if !found {
			writeAPIError(w, http.StatusNotFound, "Practice exercise not found.")
			return
		}
	}
	progress, err := api.bot.store.mistakeTransferProgress(user.TelegramID, language, time.Now().UTC().Add(-7*24*time.Hour))
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "Unable to load practice progress.")
		return
	}
	writeJSON(w, http.StatusOK, progress)
}
