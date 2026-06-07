package main

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"
	"sync"
)

type promptFile struct {
	Version int                   `json:"version"`
	Prompts map[string]promptSpec `json:"prompts"`
}

type promptSpec struct {
	Feature  string `json:"feature"`
	Tool     string `json:"tool"`
	Function string `json:"function"`
	Notes    string `json:"notes"`
	Template string `json:"template"`
}

type promptRegistry struct {
	mu      sync.RWMutex
	entries map[string]promptSpec
}

var appPromptRegistry = &promptRegistry{entries: map[string]promptSpec{}}

func loadAppPromptFile(path string) (int, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return 0, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, nil
		}
		return 0, err
	}
	var file promptFile
	if err := json.Unmarshal(data, &file); err != nil {
		return 0, err
	}
	entries := make(map[string]promptSpec, len(file.Prompts))
	for key, spec := range file.Prompts {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		entries[key] = spec
	}
	appPromptRegistry.mu.Lock()
	appPromptRegistry.entries = entries
	appPromptRegistry.mu.Unlock()
	return len(entries), nil
}

func renderAppPrompt(key string, fallback string, vars map[string]string) string {
	template := appPromptTemplate(key)
	if strings.TrimSpace(template) == "" {
		template = fallback
	}
	if len(vars) == 0 {
		return template
	}
	keys := make([]string, 0, len(vars))
	for key := range vars {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	replacements := make([]string, 0, len(keys)*2)
	for _, key := range keys {
		replacements = append(replacements, "{{"+key+"}}", vars[key])
	}
	return strings.NewReplacer(replacements...).Replace(template)
}

func appPromptTemplate(key string) string {
	appPromptRegistry.mu.RLock()
	defer appPromptRegistry.mu.RUnlock()
	return appPromptRegistry.entries[strings.TrimSpace(key)].Template
}

func commonPromptVars(language learningLanguage, interfaceLanguage learningLanguage) map[string]string {
	return map[string]string{
		"learning_language_native_name":  language.NativeName,
		"learning_language_teacher_noun": language.TeacherNoun,
		"interface_language_native_name": interfaceLanguage.NativeName,
	}
}

func mergePromptVars(base map[string]string, extra map[string]string) map[string]string {
	if len(base) == 0 && len(extra) == 0 {
		return nil
	}
	merged := make(map[string]string, len(base)+len(extra))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range extra {
		merged[key] = value
	}
	return merged
}
