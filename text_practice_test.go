package main

import "testing"

func TestTextPracticeWordsRequireWholePassageTerms(t *testing.T) {
	got := textPracticeWords("The train leaves at nine.", []textPracticeWord{
		{Word: "rain", Translation: "дождь"},
		{Word: "train", Translation: "поезд"},
		{Word: "Train", Translation: "поезд"},
		{Word: "at nine", Translation: "в девять"},
	})
	if len(got) != 2 || got[0].Word != "train" || got[1].Word != "at nine" {
		t.Fatalf("unexpected cards: %+v", got)
	}
}
