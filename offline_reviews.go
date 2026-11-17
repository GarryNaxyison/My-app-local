package main

import (
	"maps"
	"net/http"
	"strings"
)

func (s *sqliteStore) syncOfflineReviews(userID int64, ids []string) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	count := 0
	for _, id := range ids {
		result, err := tx.Exec(`INSERT INTO offline_reviews (telegram_id, card_id)
			VALUES (?, ?) ON CONFLICT (telegram_id, card_id) DO NOTHING`, userID, id)
		if err != nil {
			return 0, err
		}
		added, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		count += int(added)
	}
	return count, tx.Commit()
}

func (s *jsonStore) syncOfflineReviews(userID int64, ids []string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	previous := s.users[userID]
	user := previous
	user.OfflineReviewed = maps.Clone(previous.OfflineReviewed)
	if user.OfflineReviewed == nil {
		user.OfflineReviewed = make(map[string]bool)
	}
	count := 0
	for _, id := range ids {
		if !user.OfflineReviewed[id] {
			user.OfflineReviewed[id] = true
			count++
		}
	}
	if count == 0 {
		return 0, nil
	}
	s.users[userID] = user
	if err := s.saveLocked(); err != nil {
		s.users[userID] = previous
		return 0, err
	}
	return count, nil
}

func (api *webAPI) handleOfflineReviewSync(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		CardIDs []string `json:"card_ids"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	if len(req.CardIDs) > 100 {
		writeAPIError(w, http.StatusBadRequest, "Too many cards.")
		return
	}
	for _, id := range req.CardIDs {
		if strings.TrimSpace(id) == "" || len(id) > 4096 {
			writeAPIError(w, http.StatusBadRequest, "Invalid card ID.")
			return
		}
	}
	count, err := api.bot.store.syncOfflineReviews(user.TelegramID, req.CardIDs)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "Could not sync reviews.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"synced": count})
}
