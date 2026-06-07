package main

import "testing"

func TestNormalizeOpenRouterTTSModelMigratesRemovedOpenAIModel(t *testing.T) {
	got := normalizeOpenRouterTTSModel("openai/gpt-4o-mini-tts-2025-12-15")
	if got != "google/gemini-3.1-flash-tts-preview" {
		t.Fatalf("unexpected migrated TTS model: %q", got)
	}
}

func TestNormalizeOpenRouterTTSVoiceMigratesAlloyForGemini(t *testing.T) {
	got := normalizeOpenRouterTTSVoice("google/gemini-3.1-flash-tts-preview", "alloy")
	if got != "Kore" {
		t.Fatalf("unexpected migrated TTS voice: %q", got)
	}
}

func TestWrapPCM16LEAsWAV(t *testing.T) {
	wav := wrapPCM16LEAsWAV([]byte{0x01, 0x02, 0x03, 0x04}, 24000, 1)
	if len(wav) != 48 {
		t.Fatalf("unexpected wav length: %d", len(wav))
	}
	if string(wav[:4]) != "RIFF" || string(wav[8:12]) != "WAVE" || string(wav[36:40]) != "data" {
		t.Fatalf("missing wav markers: %q %q %q", wav[:4], wav[8:12], wav[36:40])
	}
	if got := audioContentType(wav); got != "audio/wav" {
		t.Fatalf("unexpected wav content type: %q", got)
	}
	if got := audioFilenameForBytes("pronunciation.mp3", wav); got != "pronunciation.wav" {
		t.Fatalf("unexpected wav filename: %q", got)
	}
}
