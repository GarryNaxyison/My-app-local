package main

import "strings"

func splitTelegramText(text string, limit int) []string {
	if len(text) <= limit {
		return []string{text}
	}

	var chunks []string
	current := strings.TrimSpace(text)
	for len(current) > limit {
		splitAt := strings.LastIndex(current[:limit], "\n")
		if splitAt < 0 {
			splitAt = limit
		}
		chunks = append(chunks, strings.TrimSpace(current[:splitAt]))
		current = strings.TrimSpace(current[splitAt:])
	}
	if current != "" {
		chunks = append(chunks, current)
	}
	return chunks
}

func smoothDrafts(text string) []string {
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	step := 28
	if len(runes) > 500 {
		step = 45
	}

	var drafts []string
	for end := step; end < len(runes); end += step {
		drafts = append(drafts, string(runes[:end]))
	}
	drafts = append(drafts, text)
	return drafts
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", `\_`,
		"*", `\*`,
		"[", `\[`,
		"]", `\]`,
		"(", `\(`,
		")", `\)`,
		"~", `\~`,
		"`", "\\`",
		">", `\>`,
		"#", `\#`,
		"+", `\+`,
		"-", `\-`,
		"=", `\=`,
		"|", `\|`,
		"{", `\{`,
		"}", `\}`,
		".", `\.`,
		"!", `\!`,
	)
	return replacer.Replace(text)
}
