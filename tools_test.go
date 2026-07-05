package main

import (
	"testing"
	"time"
)

func TestLargestTelegramPhoto(t *testing.T) {
	photo := largestTelegramPhoto([]telegramPhoto{
		{FileID: "small", Width: 90, Height: 90},
		{FileID: "wide", Width: 200, Height: 40},
		{FileID: "large", Width: 160, Height: 160},
	})
	if photo.FileID != "large" {
		t.Fatalf("expected largest photo, got %q", photo.FileID)
	}
}

func TestMessageThrottleTimeout(t *testing.T) {
	b := &bot{}
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)

	for i := 0; i < spamMaxMessages; i++ {
		blocked, _, _ := b.checkMessageThrottle(123, now.Add(time.Duration(i)*time.Millisecond))
		if blocked {
			t.Fatalf("message %d should not be blocked", i+1)
		}
	}

	blocked, notify, until := b.checkMessageThrottle(123, now.Add(time.Second))
	if !blocked {
		t.Fatal("expected spam message to be blocked")
	}
	if !notify {
		t.Fatal("expected first timeout message to notify")
	}
	if until.Sub(now.Add(time.Second)) != spamTimeout {
		t.Fatalf("expected timeout %s, got %s", spamTimeout, until.Sub(now.Add(time.Second)))
	}
}

func TestTotalLeaderboardScoreUsesOnlyLanguageActivity(t *testing.T) {
	user := userState{
		XP: 100,
		LearnedWords: []learnedWordEntry{
			{ID: "en:hello", Language: "en", ReviewCorrectCount: 10},
			{ID: "de:hallo", Language: "de", ReviewCorrectCount: 10},
		},
		Mistakes: []mistakeEntry{
			{Word: "helo", Language: "en"},
			{Word: "gutten", Language: "de"},
			{Word: "bonjor", Language: "fr"},
		},
	}

	got := totalLeaderboardScore(user)
	want := user.XP
	if got != want {
		t.Fatalf("expected total leaderboard score %d, got %d", want, got)
	}
}

func TestTotalLeaderboardScoreNeverDropsBelowXP(t *testing.T) {
	user := userState{
		XP: 100,
		LearnedWords: []learnedWordEntry{
			{ID: "en:hello", Language: "en"},
		},
	}

	if got := totalLeaderboardScore(user); got != user.XP {
		t.Fatalf("totalLeaderboardScore() = %d, want XP floor %d", got, user.XP)
	}
}

func TestWebLeaderboardDTOUsesMonotonicRatingScore(t *testing.T) {
	got := webLeaderboardDTO([]leaderboardEntry{{
		FirstName: "Maria",
		Score:     2007,
		XP:        2007,
		Words:     2,
		Mistakes:  1,
		Level:     5,
		Title:     "Explorer",
	}})
	if len(got) != 1 {
		t.Fatalf("expected one leaderboard DTO entry, got %#v", got)
	}
	entry := got[0]
	if entry["score"] != 2007 || entry["rating_points"] != 2007 {
		t.Fatalf("expected rating score fields to be 2007, got %#v", entry)
	}
	if entry["xp"] != 2007 {
		t.Fatalf("expected raw XP to stay separate from rating score, got %#v", entry)
	}
}

func TestLeaderboardLanguagesIncludesCurrentAndStudiedLanguages(t *testing.T) {
	user := userState{
		LearningLanguage: "fr",
		LearnedWords: []learnedWordEntry{
			{ID: "en:hello", Language: "en"},
			{ID: "de:hallo", Language: "de"},
		},
		Mistakes: []mistakeEntry{
			{Word: "hola", Language: "es"},
			{Word: "bonjor", Language: "fr"},
		},
	}

	got := leaderboardLanguages(user)
	want := []string{"English", "Spanish", "German", "French"}
	if len(got) != len(want) {
		t.Fatalf("expected languages %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected languages %v, got %v", want, got)
		}
	}
}

func TestLanguageLeaderboardDoesNotUseGlobalXP(t *testing.T) {
	store := &jsonStore{users: map[int64]userState{
		1: {
			TelegramID:       1,
			FirstName:        "Arseniy",
			LearningLanguage: "it",
			XP:               2007,
		},
		2: {
			TelegramID:       2,
			FirstName:        "Maria",
			LearningLanguage: "it",
			LearnedWords: []learnedWordEntry{
				{ID: "it:casa", Language: "it", Russian: "дом", English: "casa", ReviewCorrectCount: 10},
			},
		},
	}}

	entries := store.languageLeaderboard("it", 10)
	if len(entries) != 1 {
		t.Fatalf("expected only the user with Italian activity in leaderboard, got %#v", entries)
	}
	if entries[0].FirstName != "Maria" || entries[0].Score != 15 {
		t.Fatalf("unexpected Italian leaderboard entry: %#v", entries[0])
	}
}
