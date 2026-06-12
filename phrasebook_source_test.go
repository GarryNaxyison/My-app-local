package main

import "testing"

func TestNormalizePhrasebookSourceKeepsMistake(t *testing.T) {
	if got := normalizePhrasebookSource("mistake"); got != "mistake" {
		t.Fatalf("normalizePhrasebookSource(%q) = %q, want %q", "mistake", got, "mistake")
	}
}
