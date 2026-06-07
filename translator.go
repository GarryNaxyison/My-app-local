package main

import (
	"encoding/json"
	"strings"
)

const translationAutoCode = "auto"

type translationToolResult struct {
	SourceText  string `json:"source_text"`
	Translation string `json:"translation"`
}

func defaultTranslatorSource() string {
	return translationAutoCode
}

func defaultTranslatorTarget(user userState) string {
	target := normalizeTranslatorTarget(user.InterfaceLanguage)
	if target == "" {
		target = normalizeTranslatorTarget(user.LearningLanguage)
	}
	if target == "" {
		target = "ru"
	}
	return target
}

func translatorMode(sourceCode, targetCode string) string {
	return modeToolTranslatorPrefix + normalizeTranslatorSource(sourceCode) + ":" + normalizeTranslatorTarget(targetCode)
}

func parseTranslatorMode(mode string, user userState) (string, string) {
	sourceCode := defaultTranslatorSource()
	targetCode := defaultTranslatorTarget(user)
	raw := strings.TrimPrefix(mode, modeToolTranslatorPrefix)
	if raw != mode && raw != "" {
		parts := strings.Split(raw, ":")
		if len(parts) >= 1 {
			sourceCode = normalizeTranslatorSource(parts[0])
		}
		if len(parts) >= 2 {
			targetCode = normalizeTranslatorTarget(parts[1])
		}
	}
	if targetCode == "" {
		targetCode = defaultTranslatorTarget(user)
	}
	return sourceCode, targetCode
}

func normalizeTranslatorSource(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "" || code == translationAutoCode {
		return translationAutoCode
	}
	return normalizeLearningLanguage(code)
}

func normalizeTranslatorTarget(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "" || code == translationAutoCode {
		return ""
	}
	return normalizeLearningLanguage(code)
}

func translatorLanguageName(code string) string {
	code = normalizeTranslatorSource(code)
	if code == translationAutoCode {
		return "Auto"
	}
	return learningLanguageByCode(code).NativeName
}

func translatorLanguageInterfaceName(code string) string {
	code = normalizeTranslatorSource(code)
	if code == translationAutoCode {
		return "Auto"
	}
	language := learningLanguageByCode(code)
	if strings.TrimSpace(language.InterfaceName) != "" {
		return language.InterfaceName
	}
	return language.NativeName
}

func translationToolPrompt(text string, sourceCode string, targetCode string, interfaceLanguage learningLanguage) []chatMessage {
	sourceName := translatorLanguageName(sourceCode)
	targetName := translatorLanguageName(targetCode)
	sourceInstruction := "Detect the source language automatically."
	if normalizeTranslatorSource(sourceCode) != translationAutoCode {
		sourceInstruction = "The source language is " + sourceName + "."
	}
	systemPrompt := "You are a precise AI translator. Translate faithfully, preserve meaning, names, numbers, line breaks when useful, and adapt idioms naturally. " +
		"Do not explain unless the input is impossible to translate. Return only the translated text."
	userPrompt := sourceInstruction + "\nTranslate into natural " + targetName + ". " +
		"The user's interface language is " + interfaceLanguage.NativeName + ", but the output must be only in " + targetName + ".\n\nText:\n" + strings.TrimSpace(text)
	return []chatMessage{
		{
			Role: "system",
			Content: renderAppPrompt("tools.translation_pair.system", systemPrompt, map[string]string{
				"source_language_name":           sourceName,
				"target_language_name":           targetName,
				"interface_language_native_name": interfaceLanguage.NativeName,
				"translator_source_instruction":  sourceInstruction,
			}),
		},
		{
			Role: "user",
			Content: renderAppPrompt("tools.translation_pair.user", userPrompt, map[string]string{
				"source_language_name":           sourceName,
				"target_language_name":           targetName,
				"interface_language_native_name": interfaceLanguage.NativeName,
				"translator_source_instruction":  sourceInstruction,
				"text":                           strings.TrimSpace(text),
			}),
		},
	}
}

func imageTranslationToolInstruction(sourceCode string, targetCode string, interfaceLanguage learningLanguage) string {
	sourceName := translatorLanguageName(sourceCode)
	targetName := translatorLanguageName(targetCode)
	sourceInstruction := "Detect the source language automatically."
	if normalizeTranslatorSource(sourceCode) != translationAutoCode {
		sourceInstruction = "Treat the source language as " + sourceName + "."
	}
	fallback := sourceInstruction + " Read the image carefully, extract the visible text, and translate it into natural " + targetName + ". " +
		"If there is no readable text, describe the image briefly in " + interfaceLanguage.NativeName + " and leave translation empty. " +
		"Return valid JSON only with exactly these string fields: source_text, translation. Do not use Markdown."
	return renderAppPrompt("tools.image_translation_instruction", fallback, map[string]string{
		"source_language_name":           sourceName,
		"target_language_name":           targetName,
		"interface_language_native_name": interfaceLanguage.NativeName,
		"translator_source_instruction":  sourceInstruction,
	})
}

func parseTranslationToolResult(raw string) translationToolResult {
	raw = strings.TrimSpace(raw)
	result := translationToolResult{}
	if raw == "" {
		return result
	}
	candidate := stripJSONCodeFence(raw)
	if err := json.Unmarshal([]byte(candidate), &result); err == nil {
		result.SourceText = strings.TrimSpace(result.SourceText)
		result.Translation = strings.TrimSpace(result.Translation)
		if result.SourceText != "" || result.Translation != "" {
			return result
		}
	}
	result.Translation = raw
	return result
}

func stripJSONCodeFence(text string) string {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "```") {
		return text
	}
	lines := strings.Split(text, "\n")
	if len(lines) < 3 {
		return text
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[len(lines)-1]), "```") {
		return text
	}
	return strings.TrimSpace(strings.Join(lines[1:len(lines)-1], "\n"))
}

func translatorToolResultText(copy uiCopy, sourceText string, translation string) string {
	var builder strings.Builder
	sourceText = strings.TrimSpace(sourceText)
	translation = strings.TrimSpace(translation)
	if sourceText != "" {
		builder.WriteString(copy.Tool.TranscriptLabel + ":\n")
		builder.WriteString(sourceText)
		builder.WriteString("\n\n")
	}
	builder.WriteString(copy.Tool.TranslationLabel + ":\n")
	builder.WriteString(translation)
	return strings.TrimSpace(builder.String())
}
