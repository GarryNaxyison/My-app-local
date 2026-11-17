package main

import (
	"net/http"
	"regexp"
	"slices"
	"time"
)

var learningEventID = regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)

func (api *webAPI) handleLearningEvent(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var input struct {
		ID        string `json:"id"`
		AttemptID string `json:"attempt_id"`
		Stage     string `json:"stage"`
		Format    string `json:"format"`
		ElapsedMS int    `json:"elapsed_ms"`
	}
	if !decodeJSONRequest(w, r, &input) {
		return
	}
	stages := []string{"today_opened", "lesson_started", "answer_submitted", "answer_checked", "retry_submitted", "review_due", "transfer_completed"}
	formats := []string{"", "text", "voice"}
	if !learningEventID.MatchString(input.ID) || (input.AttemptID != "" && !learningEventID.MatchString(input.AttemptID)) || !slices.Contains(stages, input.Stage) ||
		!slices.Contains(formats, input.Format) || input.ElapsedMS < 0 || input.ElapsedMS > 86400000 {
		writeAPIError(w, http.StatusBadRequest, "Invalid learning event.")
		return
	}
	storage, ok := api.bot.store.(*sqliteStore)
	if !ok {
		writeJSON(w, http.StatusOK, map[string]bool{"stored": false})
		return
	}
	_, err = storage.db.ExecContext(r.Context(), `INSERT INTO learning_events
		(telegram_id,event_id,attempt_id,stage,format,elapsed_ms,created_at) VALUES (?,?,?,?,?,?,?)
		ON CONFLICT(telegram_id,event_id) DO NOTHING`, user.TelegramID, input.ID, input.AttemptID, input.Stage, input.Format, input.ElapsedMS, formatDBTime(time.Now()))
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, "Could not store learning event.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"stored": true})
}
