package main

import (
	"log"
	"net/http"
	"sort"
	"time"
)

type learningReviewItem struct {
	ID       string    `json:"id"`
	Source   string    `json:"source"`
	LessonID string    `json:"lesson_id"`
	Title    string    `json:"title"`
	DueAt    time.Time `json:"due_at"`
	Status   string    `json:"status"`
}

func listTutorLearningReviews(storage store, user userState, now time.Time) ([]learningReviewItem, error) {
	reviews, err := storage.listScheduledAITutorReviews(user.TelegramID)
	if err != nil {
		return nil, err
	}
	// Older records are keyed by lesson and interval. Use one active entry
	// per lesson and prefer the latest choice among valid schedules.
	sort.Slice(reviews, func(i, j int) bool {
		a, b := reviews[i], reviews[j]
		aUpdated, bUpdated := parseDBTime(a.UpdatedAt), parseDBTime(b.UpdatedAt)
		if !aUpdated.Equal(bUpdated) {
			return aUpdated.After(bUpdated)
		}
		aDue, bDue := parseDBTime(a.DueAt), parseDBTime(b.DueAt)
		if !aDue.Equal(bDue) {
			return aDue.Before(bDue)
		}
		return a.ID < b.ID
	})
	items := make([]learningReviewItem, 0, len(reviews))
	seen := make(map[string]bool, len(reviews))
	for _, review := range reviews {
		dueAt := parseDBTime(review.DueAt)
		if dueAt.IsZero() || seen[review.LessonID] {
			continue
		}
		seen[review.LessonID] = true
		lesson, found, err := storage.getAITutorLesson(review.LessonID)
		if err != nil {
			return nil, err
		}
		if !found || normalizeLearningLanguage(lesson.LearningLanguage) != normalizeLearningLanguage(user.LearningLanguage) {
			continue
		}
		status := "upcoming"
		if !dueAt.After(now) {
			status = "due"
		}
		items = append(items, learningReviewItem{
			ID: "ai_tutor:" + lesson.ID, Source: "ai_tutor", LessonID: lesson.ID,
			Title: firstNonEmpty(lesson.Payload.Title, lesson.Payload.Theme, lesson.Theme, "AI Tutor"),
			DueAt: dueAt.UTC(), Status: status,
		})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].DueAt.Equal(items[j].DueAt) {
			return items[i].ID < items[j].ID
		}
		return items[i].DueAt.Before(items[j].DueAt)
	})
	return items, nil
}

func (api *webAPI) handleLearningReviews(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	items, err := listTutorLearningReviews(api.bot.store, user, time.Now().UTC())
	if err != nil {
		log.Printf("list learning reviews: %v", err)
		writeAPIError(w, http.StatusInternalServerError, "Unable to load learning reviews.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
