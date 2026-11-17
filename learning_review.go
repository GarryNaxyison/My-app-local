package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

func learningReviewWebURL(base, id string) (string, error) {
	target, err := url.Parse(base)
	if err != nil || target.Scheme != "https" || target.Hostname() == "" {
		return "", fmt.Errorf("invalid review web URL")
	}
	query := target.Query()
	query.Set("view", "home")
	query.Set("review", id)
	target.RawQuery = query.Encode()
	return target.String(), nil
}

func (b *bot) sendLearningReview(ctx context.Context, chatID int64, user userState) error {
	items, err := listLearningReviews(b.store, user, time.Now().UTC())
	if err != nil {
		return err
	}
	russian := normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru"
	if len(items) == 0 {
		message := "No reviews scheduled."
		if russian {
			message = "Повторений пока нет."
		}
		return b.telegram.sendInlineMessage(ctx, chatID, message, backToMenuKeyboard(ui(user)))
	}
	item := items[0]
	keyboard, err := learningReviewKeyboard(item, russian)
	if err != nil {
		return err
	}
	title := "Next review"
	if russian {
		title = "Ближайшее повторение"
	}
	return b.telegram.sendInlineMessage(ctx, chatID, title+"\n"+item.Title, keyboard)
}

func learningReviewKeyboard(item learningReviewItem, russian bool) (map[string]any, error) {
	base := strings.TrimSpace(os.Getenv("WEB_APP_URL"))
	if base == "" {
		base = "https://neriva.ru/app"
	}
	link, err := learningReviewWebURL(base, item.ID)
	if err != nil {
		return nil, err
	}
	action := "Open review"
	if russian {
		action = "Открыть повторение"
	}
	return map[string]any{
		"inline_keyboard": [][]map[string]string{{{"text": action, "url": link}}},
	}, nil
}

type learningReviewItem struct {
	ID       string    `json:"id"`
	Source   string    `json:"source"`
	LessonID string    `json:"lesson_id,omitempty"`
	Title    string    `json:"title"`
	DueAt    time.Time `json:"due_at"`
	Status   string    `json:"status"`
}

func mistakeLearningReviewID(entry mistakeEntry) string {
	key, _ := json.Marshal([4]string{normalizeLearningLanguage(entry.Language), entry.Word, entry.Correction, formatDBTime(entry.AddedAt)})
	return fmt.Sprintf("mistake:%x", sha256.Sum256(key))
}

func mistakeLearningReviewIndex(mistakes []mistakeEntry, reviewID string) int {
	for index, entry := range mistakes {
		if mistakeLearningReviewID(entry) == reviewID {
			return index
		}
	}
	return -1
}

func listLearningReviews(storage store, user userState, now time.Time) ([]learningReviewItem, error) {
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
	for _, mistake := range mistakesForLanguage(user) {
		if mistake.Word == "" || mistake.Correction == "" {
			continue
		}
		items = append(items, learningReviewItem{
			ID: mistakeLearningReviewID(mistake), Source: "mistake",
			Title: mistake.Word, DueAt: mistake.AddedAt.UTC(), Status: "due",
		})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Status != items[j].Status {
			return items[i].Status == "due"
		}
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
	items, err := listLearningReviews(api.bot.store, user, time.Now().UTC())
	if err != nil {
		log.Printf("list learning reviews: %v", err)
		writeAPIError(w, http.StatusInternalServerError, "Unable to load learning reviews.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
