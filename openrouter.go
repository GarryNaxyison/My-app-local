package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type openRouterClient struct {
	apiKey  string
	model   string
	appURL  string
	appName string
	http    *http.Client
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type visionChatMessage struct {
	Role    string `json:"role"`
	Content []any  `json:"content"`
}

type audioTranscription struct {
	Text            string
	Tokens          []transcriptionToken
	DurationSeconds float64
	Normalized      bool
}

type transcriptionToken struct {
	Token   string
	LogProb float64
}

func newOpenRouterClient(apiKey, model, appURL, appName string, httpClient *http.Client) *openRouterClient {
	return &openRouterClient{
		apiKey:  apiKey,
		model:   model,
		appURL:  appURL,
		appName: appName,
		http:    httpClient,
	}
}

func (c *openRouterClient) complete(ctx context.Context, messages []chatMessage, temperature float64, maxTokens int) (string, error) {
	return c.completeWithModel(ctx, c.model, messages, temperature, maxTokens)
}

func (c *openRouterClient) completeWithModel(ctx context.Context, model string, messages []chatMessage, temperature float64, maxTokens int) (string, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		model = c.model
	}
	payload := map[string]any{
		"model":       model,
		"messages":    messages,
		"temperature": temperature,
		"max_tokens":  maxTokens,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", c.appURL)
	req.Header.Set("X-OpenRouter-Title", c.appName)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		if decoded.Error != nil && decoded.Error.Message != "" {
			return "", fmt.Errorf("openrouter: %s", decoded.Error.Message)
		}
		return "", fmt.Errorf("openrouter: status %d", resp.StatusCode)
	}
	if len(decoded.Choices) == 0 || decoded.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("openrouter returned an empty response")
	}
	return cleanModelReply(decoded.Choices[0].Message.Content), nil
}

func (c *openRouterClient) completeAudioWithModel(ctx context.Context, model string, audioBytes []byte, format string, systemPrompt string, userPrompt string, temperature float64, maxTokens int) (string, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		model = c.model
	}
	audioBytes, format, _ = normalizeAudioForSTT(ctx, audioBytes, format)
	if strings.TrimSpace(format) == "" {
		format = "wav"
	}
	base64Audio := base64.StdEncoding.EncodeToString(audioBytes)
	payload := map[string]any{
		"model": model,
		"messages": []map[string]any{
			{
				"role": "system",
				"content": []map[string]any{
					{"type": "text", "text": systemPrompt},
				},
			},
			{
				"role": "user",
				"content": []map[string]any{
					{"type": "text", "text": userPrompt},
					{
						"type": "input_audio",
						"inputAudio": map[string]any{
							"data":   base64Audio,
							"format": format,
						},
					},
				},
			},
		},
		"temperature": temperature,
		"max_tokens":  maxTokens,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", c.appURL)
	req.Header.Set("X-OpenRouter-Title", c.appName)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		if decoded.Error != nil && decoded.Error.Message != "" {
			return "", fmt.Errorf("openrouter audio: %s", decoded.Error.Message)
		}
		return "", fmt.Errorf("openrouter audio: status %d", resp.StatusCode)
	}
	if len(decoded.Choices) == 0 || strings.TrimSpace(decoded.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("openrouter audio returned an empty response")
	}
	return cleanModelReply(decoded.Choices[0].Message.Content), nil
}

func (c *openRouterClient) translateImageText(ctx context.Context, imageBytes []byte, mimeType string, learningLanguage learningLanguage, interfaceLanguage learningLanguage) (string, error) {
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	dataURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(imageBytes)
	imagePrompt := "Read the image carefully. If there is visible text, extract it exactly, detect the source language, and translate it into natural " + interfaceLanguage.NativeName + ". If the image has no readable text, briefly describe what you see in " + interfaceLanguage.NativeName + ". Return plain Telegram text only, no Markdown. Use three short sections in " + interfaceLanguage.NativeName + " equivalent to: 1) What I see, 2) Text in the image, 3) Translation."
	payload := map[string]any{
		"model": c.model,
		"messages": []visionChatMessage{
			{
				Role: "system",
				Content: []any{
					map[string]any{"type": "text", "text": coachSystemPrompt(learningLanguage, interfaceLanguage)},
				},
			},
			{
				Role: "user",
				Content: []any{
					map[string]any{
						"type": "text",
						"text": renderAppPrompt("tools.image_translate_inline.user", imagePrompt, mergePromptVars(commonPromptVars(learningLanguage, interfaceLanguage), nil)),
					},
					map[string]any{
						"type": "image_url",
						"image_url": map[string]any{
							"url": dataURL,
						},
					},
				},
			},
		},
		"temperature": 0.2,
		"max_tokens":  900,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", c.appURL)
	req.Header.Set("X-OpenRouter-Title", c.appName)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		if decoded.Error != nil && decoded.Error.Message != "" {
			return "", fmt.Errorf("openrouter vision: %s", decoded.Error.Message)
		}
		return "", fmt.Errorf("openrouter vision: status %d", resp.StatusCode)
	}
	if len(decoded.Choices) == 0 || strings.TrimSpace(decoded.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("openrouter vision returned an empty response")
	}
	return cleanModelReply(decoded.Choices[0].Message.Content), nil
}

func (c *openRouterClient) describePracticeImage(ctx context.Context, imageBytes []byte, mimeType string, learningLanguage learningLanguage, interfaceLanguage learningLanguage, learnerText string) (string, error) {
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	dataURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(imageBytes)
	note := strings.TrimSpace(learnerText)
	if note == "" {
		note = "No learner note was provided."
	}
	payload := map[string]any{
		"model": c.model,
		"messages": []visionChatMessage{
			{
				Role: "system",
				Content: []any{
					map[string]any{
						"type": "text",
						"text": "You prepare compact visual context for a language-learning coach. You do not answer the learner. You only describe the attached image so another model can continue the practice chat.",
					},
				},
			},
			{
				Role: "user",
				Content: []any{
					map[string]any{
						"type": "text",
						"text": "Target learning language: " + learningLanguage.NativeName + ". Interface language: " + interfaceLanguage.NativeName + ". Learner note: " + note + "\n\nReturn plain text only, no Markdown tables. Keep it concise. Include: 1) what is visible, 2) any readable text, 3) 6 useful " + learningLanguage.NativeName + " words or phrases with short " + interfaceLanguage.NativeName + " meanings, 4) one natural practice situation for this image. Do not invent details that are not visible.",
					},
					map[string]any{
						"type": "image_url",
						"image_url": map[string]any{
							"url": dataURL,
						},
					},
				},
			},
		},
		"temperature": 0.2,
		"max_tokens":  700,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", c.appURL)
	req.Header.Set("X-OpenRouter-Title", c.appName)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		if decoded.Error != nil && decoded.Error.Message != "" {
			return "", fmt.Errorf("openrouter practice vision: %s", decoded.Error.Message)
		}
		return "", fmt.Errorf("openrouter practice vision: status %d", resp.StatusCode)
	}
	if len(decoded.Choices) == 0 || strings.TrimSpace(decoded.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("openrouter practice vision returned an empty response")
	}
	return cleanModelReply(decoded.Choices[0].Message.Content), nil
}

func (c *openRouterClient) translateImageForTool(ctx context.Context, imageBytes []byte, mimeType string, sourceCode string, targetCode string, interfaceLanguage learningLanguage) (translationToolResult, error) {
	return c.translateImageForToolWithModel(ctx, c.model, imageBytes, mimeType, sourceCode, targetCode, interfaceLanguage)
}

func (c *openRouterClient) translateImageForToolWithModel(ctx context.Context, model string, imageBytes []byte, mimeType string, sourceCode string, targetCode string, interfaceLanguage learningLanguage) (translationToolResult, error) {
	if mimeType == "" {
		mimeType = "image/jpeg"
	}
	model = strings.TrimSpace(model)
	if model == "" {
		model = c.model
	}
	dataURL := "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(imageBytes)
	payload := map[string]any{
		"model": model,
		"messages": []visionChatMessage{
			{
				Role: "system",
				Content: []any{
					map[string]any{"type": "text", "text": "You are a precise OCR and translation engine. Follow the user's requested source and target languages exactly."},
				},
			},
			{
				Role: "user",
				Content: []any{
					map[string]any{
						"type": "text",
						"text": imageTranslationToolInstruction(sourceCode, targetCode, interfaceLanguage),
					},
					map[string]any{
						"type": "image_url",
						"image_url": map[string]any{
							"url": dataURL,
						},
					},
				},
			},
		},
		"temperature": 0.1,
		"max_tokens":  900,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return translationToolResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return translationToolResult{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", c.appURL)
	req.Header.Set("X-OpenRouter-Title", c.appName)

	resp, err := c.http.Do(req)
	if err != nil {
		return translationToolResult{}, err
	}
	defer resp.Body.Close()

	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return translationToolResult{}, err
	}
	if resp.StatusCode >= 400 {
		if decoded.Error != nil && decoded.Error.Message != "" {
			return translationToolResult{}, fmt.Errorf("openrouter translator vision: %s", decoded.Error.Message)
		}
		return translationToolResult{}, fmt.Errorf("openrouter translator vision: status %d", resp.StatusCode)
	}
	if len(decoded.Choices) == 0 || strings.TrimSpace(decoded.Choices[0].Message.Content) == "" {
		return translationToolResult{}, fmt.Errorf("openrouter translator vision returned an empty response")
	}
	result := parseTranslationToolResult(cleanModelReply(decoded.Choices[0].Message.Content))
	if strings.TrimSpace(result.SourceText) == "" && strings.TrimSpace(result.Translation) == "" {
		return translationToolResult{}, fmt.Errorf("openrouter translator vision returned an empty result")
	}
	return result, nil
}

func cleanModelReply(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}
	lines := strings.Split(text, "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "you are an encouraging ") ||
			strings.HasPrefix(lower, "act like an expert ") ||
			strings.HasPrefix(lower, "the learner's interface language is ") ||
			strings.HasPrefix(lower, "base explanations on standard language teaching practice") ||
			strings.HasPrefix(lower, "when you teach vocabulary") ||
			strings.HasPrefix(lower, "when checking answers") ||
			strings.HasPrefix(lower, "keep replies practical") ||
			strings.HasPrefix(lower, "never use markdown syntax") {
			continue
		}
		filtered = append(filtered, line)
	}
	return strings.TrimSpace(strings.Join(filtered, "\n"))
}

func (c *openRouterClient) transcribeAudio(ctx context.Context, model string, audioBytes []byte, format string) (string, error) {
	transcription, err := c.transcribeAudioDetailed(ctx, model, audioBytes, format, "", "")
	if err != nil {
		return "", err
	}
	return transcription.Text, nil
}

func (c *openRouterClient) transcribeAudioDetailed(ctx context.Context, model string, audioBytes []byte, format string, languageCode string, prompt string) (audioTranscription, error) {
	audioBytes, format, normalized := normalizeAudioForSTT(ctx, audioBytes, format)
	model = strings.TrimSpace(model)
	if strings.Contains(model, "transcribe") || strings.Contains(model, "whisper") {
		transcription, err := c.transcribeAudioWithSTTEndpointDetailed(ctx, model, audioBytes, format, languageCode, prompt, true)
		if err != nil {
			transcription, err = c.transcribeAudioWithSTTEndpointDetailed(ctx, model, audioBytes, format, languageCode, prompt, false)
		}
		transcription.Normalized = normalized
		return transcription, err
	}
	transcription, err := c.transcribeAudioWithChatDetailed(ctx, model, audioBytes, format)
	transcription.Normalized = normalized
	return transcription, err
}

func (c *openRouterClient) synthesizeSpeech(ctx context.Context, model string, voice string, text string) ([]byte, error) {
	model = normalizeOpenRouterTTSModel(model)
	voice = normalizeOpenRouterTTSVoice(model, voice)
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("openrouter tts input is empty")
	}
	responseFormat := openRouterTTSResponseFormat(model)
	payload := map[string]any{
		"model":           model,
		"voice":           voice,
		"input":           text,
		"response_format": responseFormat,
		"speed":           0.95,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/audio/speech", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", c.appURL)
	req.Header.Set("X-OpenRouter-Title", c.appName)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		var decoded struct {
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if json.Unmarshal(respBody, &decoded) == nil && decoded.Error != nil && decoded.Error.Message != "" {
			return nil, fmt.Errorf("openrouter tts: %s", decoded.Error.Message)
		}
		return nil, fmt.Errorf("openrouter tts: status %d", resp.StatusCode)
	}
	if len(respBody) == 0 {
		return nil, fmt.Errorf("openrouter tts returned empty audio")
	}
	if responseFormat == "pcm" || strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "audio/pcm") {
		rate, channels := parseOpenRouterPCMContentType(resp.Header.Get("Content-Type"))
		return wrapPCM16LEAsWAV(respBody, rate, channels), nil
	}
	return respBody, nil
}

func openRouterTTSResponseFormat(model string) string {
	if strings.Contains(strings.ToLower(strings.TrimSpace(model)), "gemini") {
		return "pcm"
	}
	return "mp3"
}

func parseOpenRouterPCMContentType(contentType string) (int, int) {
	rate := 24000
	channels := 1
	for _, part := range strings.Split(contentType, ";") {
		part = strings.TrimSpace(part)
		key, value, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		parsed, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil || parsed <= 0 {
			continue
		}
		switch strings.ToLower(strings.TrimSpace(key)) {
		case "rate":
			rate = parsed
		case "channels":
			channels = parsed
		}
	}
	return rate, channels
}

func wrapPCM16LEAsWAV(pcm []byte, sampleRate int, channels int) []byte {
	if sampleRate <= 0 {
		sampleRate = 24000
	}
	if channels <= 0 {
		channels = 1
	}
	const bitsPerSample = 16
	blockAlign := channels * bitsPerSample / 8
	byteRate := sampleRate * blockAlign
	dataSize := uint32(len(pcm))

	var out bytes.Buffer
	out.Grow(44 + len(pcm))
	out.WriteString("RIFF")
	_ = binary.Write(&out, binary.LittleEndian, uint32(36)+dataSize)
	out.WriteString("WAVE")
	out.WriteString("fmt ")
	_ = binary.Write(&out, binary.LittleEndian, uint32(16))
	_ = binary.Write(&out, binary.LittleEndian, uint16(1))
	_ = binary.Write(&out, binary.LittleEndian, uint16(channels))
	_ = binary.Write(&out, binary.LittleEndian, uint32(sampleRate))
	_ = binary.Write(&out, binary.LittleEndian, uint32(byteRate))
	_ = binary.Write(&out, binary.LittleEndian, uint16(blockAlign))
	_ = binary.Write(&out, binary.LittleEndian, uint16(bitsPerSample))
	out.WriteString("data")
	_ = binary.Write(&out, binary.LittleEndian, dataSize)
	out.Write(pcm)
	return out.Bytes()
}

func audioContentType(audio []byte) string {
	if len(audio) >= 12 && string(audio[:4]) == "RIFF" && string(audio[8:12]) == "WAVE" {
		return "audio/wav"
	}
	return "audio/mpeg"
}

func audioFilenameForBytes(filename string, audio []byte) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = "audio.mp3"
	}
	if audioContentType(audio) == "audio/wav" {
		base := strings.TrimSuffix(filename, ".mp3")
		if base == filename {
			base = strings.TrimSuffix(filename, ".mpeg")
		}
		if base == "" {
			base = "audio"
		}
		return base + ".wav"
	}
	return filename
}

func (c *openRouterClient) transcribeAudioWithSTTEndpoint(ctx context.Context, model string, audioBytes []byte, format string) (string, error) {
	transcription, err := c.transcribeAudioWithSTTEndpointDetailed(ctx, model, audioBytes, format, "", "", false)
	if err != nil {
		return "", err
	}
	return transcription.Text, nil
}

func (c *openRouterClient) transcribeAudioWithSTTEndpointDetailed(ctx context.Context, model string, audioBytes []byte, format string, languageCode string, prompt string, requestLogProbs bool) (audioTranscription, error) {
	base64Audio := base64.StdEncoding.EncodeToString(audioBytes)
	payload := map[string]any{
		"model": model,
		"input_audio": map[string]any{
			"data":   base64Audio,
			"format": format,
		},
		"temperature": 0,
	}
	if languageCode = strings.TrimSpace(languageCode); languageCode != "" {
		payload["language"] = languageCode
	}
	if prompt = strings.TrimSpace(prompt); prompt != "" {
		payload["prompt"] = prompt
	}
	if requestLogProbs && strings.Contains(strings.ToLower(model), "gpt-4o") {
		payload["include"] = []string{"logprobs"}
		payload["response_format"] = "json"
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return audioTranscription{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/audio/transcriptions", bytes.NewReader(body))
	if err != nil {
		return audioTranscription{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", c.appURL)
	req.Header.Set("X-OpenRouter-Title", c.appName)

	resp, err := c.http.Do(req)
	if err != nil {
		return audioTranscription{}, err
	}
	defer resp.Body.Close()

	var decoded struct {
		Text     string `json:"text"`
		LogProbs []struct {
			Token   string  `json:"token"`
			LogProb float64 `json:"logprob"`
		} `json:"logprobs"`
		Usage *struct {
			Seconds float64 `json:"seconds"`
		} `json:"usage"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return audioTranscription{}, err
	}
	if resp.StatusCode >= 400 {
		if decoded.Error != nil && decoded.Error.Message != "" {
			return audioTranscription{}, fmt.Errorf("openrouter stt: %s", decoded.Error.Message)
		}
		return audioTranscription{}, fmt.Errorf("openrouter stt: status %d", resp.StatusCode)
	}
	transcript := strings.TrimSpace(decoded.Text)
	if transcript == "" {
		return audioTranscription{}, fmt.Errorf("openrouter stt returned an empty transcript")
	}
	tokens := make([]transcriptionToken, 0, len(decoded.LogProbs))
	for _, item := range decoded.LogProbs {
		if strings.TrimSpace(item.Token) == "" {
			continue
		}
		tokens = append(tokens, transcriptionToken{Token: item.Token, LogProb: item.LogProb})
	}
	var seconds float64
	if decoded.Usage != nil {
		seconds = decoded.Usage.Seconds
	}
	return audioTranscription{Text: transcript, Tokens: tokens, DurationSeconds: seconds}, nil
}

func (c *openRouterClient) transcribeAudioWithChat(ctx context.Context, model string, audioBytes []byte, format string) (string, error) {
	transcription, err := c.transcribeAudioWithChatDetailed(ctx, model, audioBytes, format)
	if err != nil {
		return "", err
	}
	return transcription.Text, nil
}

func (c *openRouterClient) transcribeAudioWithChatDetailed(ctx context.Context, model string, audioBytes []byte, format string) (audioTranscription, error) {
	base64Audio := base64.StdEncoding.EncodeToString(audioBytes)
	payload := map[string]any{
		"model": model,
		"messages": []map[string]any{
			{
				"role": "user",
				"content": []map[string]any{
					{
						"type": "text",
						"text": "Transcribe this audio exactly. Return only the spoken text. Do not translate, correct, explain, or answer.",
					},
					{
						"type": "input_audio",
						"inputAudio": map[string]any{
							"data":   base64Audio,
							"format": format,
						},
					},
				},
			},
		},
		"temperature": 0,
		"max_tokens":  500,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return audioTranscription{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://openrouter.ai/api/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return audioTranscription{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", c.appURL)
	req.Header.Set("X-OpenRouter-Title", c.appName)

	resp, err := c.http.Do(req)
	if err != nil {
		return audioTranscription{}, err
	}
	defer resp.Body.Close()

	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return audioTranscription{}, err
	}
	if resp.StatusCode >= 400 {
		if decoded.Error != nil && decoded.Error.Message != "" {
			return audioTranscription{}, fmt.Errorf("openrouter stt: %s", decoded.Error.Message)
		}
		return audioTranscription{}, fmt.Errorf("openrouter stt: status %d", resp.StatusCode)
	}
	if len(decoded.Choices) == 0 || decoded.Choices[0].Message.Content == "" {
		return audioTranscription{}, fmt.Errorf("openrouter stt returned an empty response")
	}

	transcript := strings.TrimSpace(decoded.Choices[0].Message.Content)
	transcript = strings.Trim(transcript, "\"“”")
	if transcript == "" {
		return audioTranscription{}, fmt.Errorf("openrouter stt returned an empty transcript")
	}
	return audioTranscription{Text: transcript}, nil
}
