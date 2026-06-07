package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderAppPromptUsesExternalTemplate(t *testing.T) {
	previous := appPromptRegistry.entries
	t.Cleanup(func() {
		appPromptRegistry.mu.Lock()
		appPromptRegistry.entries = previous
		appPromptRegistry.mu.Unlock()
	})

	path := filepath.Join(t.TempDir(), "app_prompts.json")
	content := `{"version":1,"prompts":{"test.prompt":{"feature":"test","tool":"test","function":"test","notes":"test","template":"Hello {{name}}"}}}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	count, err := loadAppPromptFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 prompt, got %d", count)
	}
	got := renderAppPrompt("test.prompt", "fallback", map[string]string{"name": "Denis"})
	if got != "Hello Denis" {
		t.Fatalf("unexpected prompt: %q", got)
	}
}

func TestRenderAppPromptFallsBackWhenMissing(t *testing.T) {
	got := renderAppPrompt("missing.prompt", "Built in {{name}}", map[string]string{"name": "prompt"})
	if got != "Built in prompt" {
		t.Fatalf("unexpected fallback render: %q", got)
	}
}

func TestAppPromptsFileCoversCorePromptKeys(t *testing.T) {
	previous := appPromptRegistry.entries
	t.Cleanup(func() {
		appPromptRegistry.mu.Lock()
		appPromptRegistry.entries = previous
		appPromptRegistry.mu.Unlock()
	})
	count, err := loadAppPromptFile("app_prompts.json")
	if err != nil {
		t.Fatal(err)
	}
	if count < 15 {
		t.Fatalf("expected core prompts to be listed, got %d", count)
	}
	for _, key := range []string{
		"coach.system",
		"lesson.generate.user",
		"lesson.feedback.user",
		"practice.chat.user",
		"tools.translation_pair.user",
		"vocabulary.example.user",
		"listening.phrase.user",
		"listening.feedback.user",
		"pronunciation.coach_report.user",
	} {
		if strings.TrimSpace(appPromptTemplate(key)) == "" {
			t.Fatalf("missing prompt template for %s", key)
		}
	}
}
