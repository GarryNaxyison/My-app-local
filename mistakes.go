package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const mistakesSentinel = "---MISTAKES---"
const mistakesPageSize = 10

// splitFeedbackAndMistakes splits the AI response into the human-readable
// feedback part and the raw JSON mistakes block. If the sentinel is absent,
// the whole text is returned as feedback and mistakes is empty.
func splitFeedbackAndMistakes(raw string) (feedback string, mistakesJSON string) {
	idx := strings.Index(raw, mistakesSentinel)
	if idx < 0 {
		return strings.TrimSpace(raw), ""
	}
	feedback = strings.TrimSpace(raw[:idx])
	mistakesJSON = strings.TrimSpace(raw[idx+len(mistakesSentinel):])
	return
}

// parseMistakesJSON parses the JSON array produced by the AI into mistakeEntry
// values. Malformed JSON is silently ignored (returns nil slice).
func parseMistakesJSON(raw string, now time.Time) []mistakeEntry {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "[]" {
		return nil
	}
	// Strip possible markdown fences the model might add
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var items []struct {
		Word        string `json:"word"`
		Correction  string `json:"correction"`
		Explanation string `json:"explanation"`
	}
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil
	}

	var out []mistakeEntry
	for _, item := range items {
		word := strings.TrimSpace(item.Word)
		correction := strings.TrimSpace(item.Correction)
		if word == "" || correction == "" {
			continue
		}
		out = append(out, mistakeEntry{
			Word:        word,
			Correction:  correction,
			Explanation: strings.TrimSpace(item.Explanation),
			AddedAt:     now,
		})
	}
	return out
}

// formatMistakesMessage renders one page of the user's mistake dictionary as a
// Telegram MarkdownV2 string. Returns an empty string if there are no mistakes.
func formatMistakesMessage(mistakes []mistakeEntry, page int, copies ...uiCopy) string {
	if len(mistakes) == 0 {
		return ""
	}
	copy := keyboardCopy(copies)
	totalPages := mistakeTotalPages(len(mistakes))
	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}
	start := page * mistakesPageSize
	end := start + mistakesPageSize
	if end > len(mistakes) {
		end = len(mistakes)
	}

	var sb strings.Builder
	sb.WriteString("*" + escapeMarkdownV2(copy.Mistakes) + "*\n\n")
	sb.WriteString(escapeMarkdownV2(fmt.Sprintf("%d/%d · %d\n\n", page+1, totalPages, len(mistakes))))
	for i, m := range mistakes[start:end] {
		sb.WriteString(escapeMarkdownV2(itoa(i+1) + ". "))
		sb.WriteString("~~" + escapeMarkdownV2(m.Word) + "~~ → `" + escapeMarkdownV2(m.Correction) + "`\n")
		if m.Explanation != "" {
			sb.WriteString("_" + escapeMarkdownV2(m.Explanation) + "_\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("`/clearmistakes`")
	return sb.String()
}

func mistakeTotalPages(total int) int {
	if total <= 0 {
		return 1
	}
	return (total + mistakesPageSize - 1) / mistakesPageSize
}
