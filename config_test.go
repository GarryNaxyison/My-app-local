package main

import "testing"

func TestLoadConfigUsesNerivaProductionDefaults(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "test-token")
	t.Setenv("OPENROUTER_API_KEY", "test-openrouter-key")
	for _, key := range []string{"WEB_CORS_ORIGINS", "WEB_APP_URL", "WEB_PAYMENT_RETURN_URL"} {
		t.Setenv(key, "")
	}

	cfg, err := configFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.WebAppURL != "https://neriva.ru/app" {
		t.Fatalf("default WEB_APP_URL = %q, want Neriva app URL", cfg.WebAppURL)
	}
	for _, want := range []string{"https://neriva.ru", "https://www.neriva.ru", "https://api.neriva.ru", "https://poliglotai.ru", "https://poliglotai.online"} {
		if !stringSliceContains(cfg.WebCORSOrigins, want) {
			t.Fatalf("default WEB_CORS_ORIGINS missing %q: %#v", want, cfg.WebCORSOrigins)
		}
	}
	if got := cfg.rollyPayBotReturnURL(); got != "https://neriva.ru/app" {
		t.Fatalf("RollyPay bot fallback URL = %q, want Neriva app URL", got)
	}
	if got := cfg.rollyPayWebReturnURL(); got != "https://poliglotai.ru/app?payment=success&provider=rollypay" {
		t.Fatalf("RollyPay web fallback URL = %q, want Poliglot success URL", got)
	}
}

func stringSliceContains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

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
