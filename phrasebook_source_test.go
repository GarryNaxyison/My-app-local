package main

import "testing"

func TestNormalizePhrasebookSourceKeepsMistake(t *testing.T) {
	if got := normalizePhrasebookSource("mistake"); got != "mistake" {
		t.Fatalf("normalizePhrasebookSource(%q) = %q, want %q", "mistake", got, "mistake")
	}
}

func TestNormalizePhrasebookSourceKeepsShadowing(t *testing.T) {
	if got := normalizePhrasebookSource("shadowing"); got != "shadowing" {
		t.Fatalf("normalizePhrasebookSource(%q) = %q, want %q", "shadowing", got, "shadowing")
	}
}
