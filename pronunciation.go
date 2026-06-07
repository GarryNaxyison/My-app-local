package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"strings"
	"time"
	"unicode"
)

type pronunciationMode string

const (
	pronunciationModeExact pronunciationMode = "exact_repeat"
	pronunciationModeFree  pronunciationMode = "free_speech"
)

type pronunciationProblemWord struct {
	Word       string  `json:"word"`
	Spoken     string  `json:"spoken,omitempty"`
	Confidence float64 `json:"confidence"`
	Issue      string  `json:"issue,omitempty"`
	Tip        string  `json:"tip,omitempty"`
}

type pronunciationPhonemeIssue struct {
	Word          string  `json:"word"`
	ExpectedSound string  `json:"expected_sound,omitempty"`
	HeardSound    string  `json:"heard_sound,omitempty"`
	Confidence    float64 `json:"confidence,omitempty"`
	Tip           string  `json:"tip,omitempty"`
}

type pronunciationAssessment struct {
	Expected          string                      `json:"expected,omitempty"`
	Spoken            string                      `json:"spoken"`
	Score             int                         `json:"score"`
	AccentStrength    int                         `json:"accent_strength"`
	Fluency           int                         `json:"fluency"`
	Similarity        int                         `json:"similarity,omitempty"`
	AverageConfidence float64                     `json:"average_confidence"`
	ConfidenceSource  string                      `json:"confidence_source"`
	SpeechRateWPM     float64                     `json:"speech_rate_wpm,omitempty"`
	ProblemWords      []pronunciationProblemWord  `json:"problem_words,omitempty"`
	PhonemeIssues     []pronunciationPhonemeIssue `json:"phoneme_issues,omitempty"`
	Tips              []string                    `json:"tips,omitempty"`
	Feedback          string                      `json:"feedback"`
	Mode              string                      `json:"mode"`
	Stress            string                      `json:"stress,omitempty"`
	Rhythm            string                      `json:"rhythm,omitempty"`
	Intonation        string                      `json:"intonation,omitempty"`
	CorrectedText     string                      `json:"corrected_text,omitempty"`
	AudioModel        string                      `json:"audio_model,omitempty"`
}

type pronunciationTechnicalReport struct {
	Expected          string                     `json:"expected,omitempty"`
	Spoken            string                     `json:"spoken"`
	Mode              string                     `json:"mode"`
	Level             string                     `json:"level"`
	Language          string                     `json:"language"`
	InterfaceLanguage string                     `json:"interface_language"`
	Similarity        int                        `json:"similarity"`
	AverageConfidence float64                    `json:"average_confidence"`
	ConfidenceSource  string                     `json:"confidence_source"`
	Fluency           int                        `json:"fluency"`
	SpeechRateWPM     float64                    `json:"speech_rate_wpm,omitempty"`
	DurationSeconds   float64                    `json:"duration_seconds,omitempty"`
	MissingWords      []string                   `json:"missing_words,omitempty"`
	Substitutions     []map[string]string        `json:"substitutions,omitempty"`
	ProblemWords      []pronunciationProblemWord `json:"problem_words,omitempty"`
	BaseScore         float64                    `json:"base_score_before_model_correction"`
}

type pronunciationCoachJSON struct {
	ModelCorrection int                        `json:"model_correction"`
	AccentStrength  int                        `json:"accent_strength"`
	Fluency         int                        `json:"fluency"`
	ProblemWords    []pronunciationProblemWord `json:"problem_words"`
	Tips            []string                   `json:"tips"`
	OverallFeedback string                     `json:"overall_feedback"`
}

type pronunciationAudioCoachJSON struct {
	Transcript        string                      `json:"transcript"`
	Score             int                         `json:"score"`
	AccentStrength    int                         `json:"accent_strength"`
	Fluency           int                         `json:"fluency"`
	Similarity        int                         `json:"similarity"`
	AverageConfidence float64                     `json:"average_confidence"`
	ProblemWords      []pronunciationProblemWord  `json:"problem_words"`
	PhonemeIssues     []pronunciationPhonemeIssue `json:"phoneme_issues"`
	Stress            string                      `json:"stress"`
	Rhythm            string                      `json:"rhythm"`
	Intonation        string                      `json:"intonation"`
	Tips              []string                    `json:"tips"`
	OverallFeedback   string                      `json:"overall_feedback"`
	CorrectedText     string                      `json:"corrected_text"`
}

type pronunciationAlignment struct {
	Target      string
	Spoken      string
	Op          string
	SpokenIndex int
	Similarity  float64
}

func normalizeAudioForSTT(ctx context.Context, audioBytes []byte, format string) ([]byte, string, bool) {
	if len(audioBytes) == 0 {
		return audioBytes, format, false
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return audioBytes, format, false
	}
	normalizeCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(normalizeCtx, "ffmpeg",
		"-hide_banner", "-loglevel", "error",
		"-i", "pipe:0",
		"-ac", "1",
		"-ar", "16000",
		"-af", "loudnorm=I=-18:TP=-1.5:LRA=11",
		"-f", "wav",
		"pipe:1",
	)
	cmd.Stdin = bytes.NewReader(audioBytes)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil || out.Len() == 0 {
		return audioBytes, format, false
	}
	return out.Bytes(), "wav", true
}

func (b *bot) voiceRecognitionConfigured() bool {
	return b != nil && b.openrouter != nil && strings.TrimSpace(b.cfg.OpenRouterSTTModel) != ""
}

func (b *bot) transcribeLearningVoice(ctx context.Context, user userState, audioBytes []byte, format string, expected string) (audioTranscription, error) {
	if b == nil || b.openrouter == nil {
		return audioTranscription{}, fmt.Errorf("speech recognition is not configured")
	}
	language := userLearningLanguage(user)
	prompt := "The learner is speaking " + language.NativeName + " for a language-learning exercise. Transcribe exactly what they actually said. Do not translate, correct, normalize, or infer a target phrase."
	if strings.TrimSpace(expected) != "" {
		prompt = "The learner is repeating a short " + language.NativeName + " phrase. Transcribe exactly what they actually said. Do not correct toward the expected phrase. Keep missing, unclear, or mispronounced words as you hear them."
	}
	return b.openrouter.transcribeAudioDetailed(ctx, b.cfg.OpenRouterSTTModel, audioBytes, format, language.Code, prompt)
}

func (b *bot) buildPronunciationAssessment(ctx context.Context, user userState, expected string, transcription audioTranscription, mode pronunciationMode) pronunciationAssessment {
	report := buildPronunciationTechnicalReport(user, expected, transcription, mode)
	assessment := localPronunciationAssessment(user, report, -1)
	if b == nil || b.openrouter == nil {
		return assessment
	}
	model := strings.TrimSpace(b.cfg.OpenRouterPronunciationModel)
	if model == "" {
		model = "google/gemini-3.1-flash-lite"
	}
	raw, err := b.openrouter.completeWithModel(ctx, model, pronunciationCoachPrompt(report), 0.1, 700)
	if err != nil {
		return assessment
	}
	coach, ok := parsePronunciationCoachJSON(raw)
	if !ok {
		return assessment
	}
	return mergePronunciationCoach(user, report, coach)
}

func (b *bot) assessPronunciationFromAudio(ctx context.Context, user userState, audioBytes []byte, format string, expected string, mode pronunciationMode) (pronunciationAssessment, audioTranscription, bool) {
	if b == nil || b.openrouter == nil || len(audioBytes) == 0 {
		return pronunciationAssessment{}, audioTranscription{}, false
	}
	model := strings.TrimSpace(b.cfg.OpenRouterPronunciationAudioModel)
	if model == "" {
		model = "google/gemini-2.5-pro"
	}
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	raw, err := b.openrouter.completeAudioWithModel(ctx, model, audioBytes, format, pronunciationAudioSystemPrompt(), pronunciationAudioUserPrompt(user, language, interfaceLanguage, expected, mode), 0.05, 1100)
	if err != nil {
		return pronunciationAssessment{}, audioTranscription{}, false
	}
	coach, ok := parsePronunciationAudioCoachJSON(raw)
	if !ok || strings.TrimSpace(coach.Transcript) == "" {
		return pronunciationAssessment{}, audioTranscription{}, false
	}
	assessment := pronunciationAssessmentFromAudioCoach(user, expected, mode, model, coach)
	return assessment, audioTranscription{Text: assessment.Spoken}, true
}

func pronunciationAudioSystemPrompt() string {
	return "You are a strict pronunciation assessor for a language-learning app. You receive learner audio directly. Do not use transcription-only assumptions. Estimate phoneme-level issues carefully but mark uncertainty in confidence. Return only valid JSON."
}

func pronunciationAudioUserPrompt(user userState, language learningLanguage, interfaceLanguage learningLanguage, expected string, mode pronunciationMode) string {
	expected = strings.TrimSpace(expected)
	if expected == "" {
		expected = "(free speech; no exact target phrase)"
	}
	return "Target language: " + language.NativeName + "\nInterface language for feedback: " + interfaceLanguage.NativeName + "\nLearner CEFR level: " + strings.TrimSpace(user.Level) + "\nMode: " + string(mode) + "\nExpected phrase/text:\n" + expected + "\n\nListen to the attached learner audio and return JSON with exactly these keys: transcript, score, accent_strength, fluency, similarity, average_confidence, problem_words, phoneme_issues, stress, rhythm, intonation, tips, overall_feedback, corrected_text.\nRules:\n- transcript: what the learner actually said, no translation.\n- score/accent_strength/fluency/similarity: 0-100 integers.\n- average_confidence: 0-1.\n- problem_words: up to 5 objects {word, spoken, confidence, issue, tip}.\n- phoneme_issues: up to 5 approximate objects {word, expected_sound, heard_sound, confidence, tip}.\n- stress/rhythm/intonation: one short diagnostic sentence each in the interface language.\n- tips: 2-4 short actionable tips in the interface language.\n- overall_feedback: 1-2 concise tutor sentences in the interface language.\n- corrected_text: the best corrected target-language version to play back with TTS. For exact_repeat, use the expected phrase unless a tiny grammar cleanup is clearly needed.\nNo Markdown. No code fence. JSON only."
}

func buildPronunciationTechnicalReport(user userState, expected string, transcription audioTranscription, mode pronunciationMode) pronunciationTechnicalReport {
	expected = strings.TrimSpace(expected)
	spoken := strings.TrimSpace(transcription.Text)
	if mode == pronunciationModeExact && expected == "" {
		mode = pronunciationModeFree
	}
	spokenUnits := shadowingUnits(spoken)
	targetUnits := shadowingUnits(expected)
	duration := transcription.DurationSeconds
	spokenConf, confidenceSource := spokenWordConfidences(transcription, spokenUnits)
	avgConfidence := averageConfidence(spokenConf, confidenceSource)
	fluency, wpm := pronunciationFluency(spokenUnits, duration, user.Level)

	var similarity int
	var problems []pronunciationProblemWord
	var missing []string
	var substitutions []map[string]string
	if mode == pronunciationModeExact {
		similarity = shadowingScore(expected, spoken)
		alignments := alignPronunciationUnits(targetUnits, spokenUnits)
		problems, missing, substitutions = pronunciationProblemsFromAlignment(alignments, spokenConf, confidenceSource)
		if len(problems) == 0 && avgConfidence < 0.80 {
			for i, unit := range spokenUnits {
				if i >= len(spokenConf) || spokenConf[i] >= 0.80 {
					continue
				}
				problems = append(problems, pronunciationProblemWord{Word: unit, Confidence: roundConfidence(spokenConf[i]), Issue: "low_confidence"})
				if len(problems) >= 4 {
					break
				}
			}
		}
	} else {
		similarity = freeSpeechTextQuality(spokenUnits)
		problems = pronunciationProblemsFromFreeSpeech(spokenUnits, spokenConf, confidenceSource)
	}

	avgConfidence = clampFloat(avgConfidence, 0, 1)
	base := pronunciationBaseScore(mode, similarity, avgConfidence, fluency)
	return pronunciationTechnicalReport{
		Expected:          expected,
		Spoken:            spoken,
		Mode:              string(mode),
		Level:             strings.TrimSpace(user.Level),
		Language:          userLearningLanguage(user).NativeName,
		InterfaceLanguage: userInterfaceLanguage(user).NativeName,
		Similarity:        similarity,
		AverageConfidence: roundConfidence(avgConfidence),
		ConfidenceSource:  confidenceSource,
		Fluency:           fluency,
		SpeechRateWPM:     math.Round(wpm*10) / 10,
		DurationSeconds:   math.Round(duration*10) / 10,
		MissingWords:      missing,
		Substitutions:     substitutions,
		ProblemWords:      problems,
		BaseScore:         math.Round(base*10) / 10,
	}
}

func spokenWordConfidences(transcription audioTranscription, spokenUnits []string) ([]float64, string) {
	if len(spokenUnits) == 0 {
		return nil, "none"
	}
	confidences := make([]float64, len(spokenUnits))
	counts := make([]int, len(spokenUnits))
	if len(transcription.Tokens) > 0 {
		index := 0
		for _, token := range transcription.Tokens {
			units := shadowingUnits(token.Token)
			if len(units) == 0 {
				continue
			}
			prob := math.Exp(token.LogProb)
			prob = clampFloat(prob, 0.01, 0.99)
			for range units {
				if index >= len(spokenUnits) {
					break
				}
				confidences[index] += prob
				counts[index]++
				index++
			}
		}
		seen := 0
		for i := range confidences {
			if counts[i] > 0 {
				confidences[i] = confidences[i] / float64(counts[i])
				seen++
			}
		}
		if seen > 0 {
			fill := averageConfidence(confidences, "logprobs")
			for i := range confidences {
				if counts[i] == 0 {
					confidences[i] = fill
				}
			}
			return confidences, "logprobs"
		}
	}
	for i := range confidences {
		confidences[i] = 0.74
	}
	return confidences, "estimated"
}

func averageConfidence(confidences []float64, source string) float64 {
	if len(confidences) == 0 {
		return 0
	}
	total := 0.0
	for _, value := range confidences {
		if value <= 0 && source == "estimated" {
			value = 0.74
		}
		total += clampFloat(value, 0, 1)
	}
	return total / float64(len(confidences))
}

func pronunciationProblemsFromAlignment(alignments []pronunciationAlignment, spokenConf []float64, source string) ([]pronunciationProblemWord, []string, []map[string]string) {
	var problems []pronunciationProblemWord
	var missing []string
	var substitutions []map[string]string
	for _, item := range alignments {
		switch item.Op {
		case "delete":
			missing = append(missing, item.Target)
			problems = append(problems, pronunciationProblemWord{Word: item.Target, Confidence: 0, Issue: "missing"})
		case "substitute":
			conf := confidenceForAlignment(item, spokenConf, source)
			substitutions = append(substitutions, map[string]string{"expected": item.Target, "spoken": item.Spoken})
			if conf < 0.80 || item.Similarity < 0.86 {
				problems = append(problems, pronunciationProblemWord{Word: item.Target, Spoken: item.Spoken, Confidence: roundConfidence(conf), Issue: "substituted_or_unclear"})
			}
		case "match":
			conf := confidenceForAlignment(item, spokenConf, source)
			if conf < 0.80 {
				problems = append(problems, pronunciationProblemWord{Word: item.Target, Spoken: item.Spoken, Confidence: roundConfidence(conf), Issue: "low_confidence"})
			}
		}
		if len(problems) >= 6 {
			break
		}
	}
	return problems, missing, substitutions
}

func pronunciationProblemsFromFreeSpeech(spokenUnits []string, spokenConf []float64, source string) []pronunciationProblemWord {
	var problems []pronunciationProblemWord
	for i, unit := range spokenUnits {
		conf := 0.74
		if i < len(spokenConf) {
			conf = spokenConf[i]
		}
		if source == "estimated" && conf >= 0.74 {
			continue
		}
		if conf < 0.80 {
			problems = append(problems, pronunciationProblemWord{Word: unit, Spoken: unit, Confidence: roundConfidence(conf), Issue: "low_confidence"})
		}
		if len(problems) >= 6 {
			break
		}
	}
	return problems
}

func confidenceForAlignment(item pronunciationAlignment, spokenConf []float64, source string) float64 {
	if item.SpokenIndex >= 0 && item.SpokenIndex < len(spokenConf) {
		return clampFloat(spokenConf[item.SpokenIndex], 0, 1)
	}
	if source == "estimated" {
		if item.Op == "match" {
			return 0.90
		}
		return clampFloat(0.35+item.Similarity*0.45, 0.20, 0.78)
	}
	return 0
}

func pronunciationBaseScore(mode pronunciationMode, similarity int, avgConfidence float64, fluency int) float64 {
	confidencePct := avgConfidence * 100
	if mode == pronunciationModeExact {
		return float64(similarity)*0.72 + confidencePct*0.16 + float64(fluency)*0.07
	}
	return confidencePct*0.55 + float64(fluency)*0.25 + float64(similarity)*0.15
}

func localPronunciationAssessment(user userState, report pronunciationTechnicalReport, modelCorrection int) pronunciationAssessment {
	if modelCorrection < 0 {
		if report.ConfidenceSource == "estimated" || report.ConfidenceSource == "none" {
			modelCorrection = 0
		} else {
			modelCorrection = int(math.Round(report.BaseScore / 0.95))
		}
	}
	score := capPronunciationScoreForConfidence(report, finalPronunciationScore(report.BaseScore, modelCorrection))
	accent := pronunciationAccentStrength(report, score)
	feedback := fallbackPronunciationFeedback(user, report, score)
	if strict := strictPronunciationFeedback(user, report, score); strict != "" {
		feedback = strict
	}
	return pronunciationAssessment{
		Expected:          report.Expected,
		Spoken:            report.Spoken,
		Score:             score,
		AccentStrength:    accent,
		Fluency:           report.Fluency,
		Similarity:        report.Similarity,
		AverageConfidence: report.AverageConfidence,
		ConfidenceSource:  report.ConfidenceSource,
		SpeechRateWPM:     report.SpeechRateWPM,
		ProblemWords:      report.ProblemWords,
		Tips:              fallbackPronunciationTips(user, report),
		Feedback:          feedback,
		Mode:              report.Mode,
	}
}

func mergePronunciationCoach(user userState, report pronunciationTechnicalReport, coach pronunciationCoachJSON) pronunciationAssessment {
	modelCorrection := coach.ModelCorrection
	if modelCorrection <= 0 {
		modelCorrection = int(math.Round(report.BaseScore / 0.95))
	}
	assessment := localPronunciationAssessment(user, report, modelCorrection)
	if coach.AccentStrength > 0 {
		assessment.AccentStrength = clampInt(coach.AccentStrength, 0, 100)
	}
	if coach.Fluency > 0 {
		assessment.Fluency = clampInt(coach.Fluency, 0, 100)
	}
	if len(coach.ProblemWords) > 0 {
		assessment.ProblemWords = normalizeCoachProblemWords(coach.ProblemWords, report.ProblemWords)
	}
	if len(coach.Tips) > 0 {
		assessment.Tips = trimStringSlice(coach.Tips, 4)
	}
	if strings.TrimSpace(coach.OverallFeedback) != "" && (assessment.Mode != string(pronunciationModeExact) || assessment.Score >= 70) {
		assessment.Feedback = strings.TrimSpace(coach.OverallFeedback)
	}
	assessment.Score = capPronunciationScoreForConfidence(report, assessment.Score)
	assessment.AccentStrength = pronunciationAccentStrength(report, assessment.Score)
	return assessment
}

func pronunciationAssessmentFromAudioCoach(user userState, expected string, mode pronunciationMode, model string, coach pronunciationAudioCoachJSON) pronunciationAssessment {
	spoken := strings.TrimSpace(coach.Transcript)
	expected = strings.TrimSpace(expected)
	score := clampInt(coach.Score, 0, 100)
	if score == 0 && spoken != "" {
		score = 50
	}
	avgConfidence := clampFloat(coach.AverageConfidence, 0, 1)
	if avgConfidence == 0 && len(coach.ProblemWords) > 0 {
		var total float64
		var count int
		for _, item := range coach.ProblemWords {
			if item.Confidence > 0 {
				total += clampFloat(item.Confidence, 0, 1)
				count++
			}
		}
		if count > 0 {
			avgConfidence = total / float64(count)
		}
	}
	if avgConfidence == 0 {
		avgConfidence = 0.74
	}
	feedback := strings.TrimSpace(coach.OverallFeedback)
	if feedback == "" {
		report := buildPronunciationTechnicalReport(user, expected, audioTranscription{Text: spoken}, mode)
		feedback = fallbackPronunciationFeedback(user, report, score)
	}
	corrected := strings.TrimSpace(coach.CorrectedText)
	if corrected == "" {
		corrected = expected
	}
	return pronunciationAssessment{
		Expected:          expected,
		Spoken:            spoken,
		Score:             score,
		AccentStrength:    clampInt(coach.AccentStrength, 0, 100),
		Fluency:           clampInt(coach.Fluency, 0, 100),
		Similarity:        clampInt(coach.Similarity, 0, 100),
		AverageConfidence: roundConfidence(avgConfidence),
		ConfidenceSource:  "gemini_audio",
		ProblemWords:      normalizeCoachProblemWords(coach.ProblemWords, nil),
		PhonemeIssues:     normalizePhonemeIssues(coach.PhonemeIssues),
		Tips:              trimStringSlice(coach.Tips, 4),
		Feedback:          feedback,
		Mode:              string(mode),
		Stress:            strings.TrimSpace(coach.Stress),
		Rhythm:            strings.TrimSpace(coach.Rhythm),
		Intonation:        strings.TrimSpace(coach.Intonation),
		CorrectedText:     corrected,
		AudioModel:        model,
	}
}

func normalizePhonemeIssues(items []pronunciationPhonemeIssue) []pronunciationPhonemeIssue {
	out := make([]pronunciationPhonemeIssue, 0, minInt(len(items), 5))
	seen := map[string]bool{}
	for _, item := range items {
		item.Word = strings.TrimSpace(item.Word)
		item.ExpectedSound = strings.TrimSpace(item.ExpectedSound)
		item.HeardSound = strings.TrimSpace(item.HeardSound)
		item.Tip = strings.TrimSpace(item.Tip)
		item.Confidence = roundConfidence(clampFloat(item.Confidence, 0, 1))
		key := strings.ToLower(item.Word + "|" + item.ExpectedSound + "|" + item.HeardSound)
		if item.Word == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, item)
		if len(out) >= 5 {
			break
		}
	}
	return out
}

func pronunciationCoachPrompt(report pronunciationTechnicalReport) []chatMessage {
	payload, _ := json.Marshal(report)
	systemPrompt := "You are an expert pronunciation coach. Return JSON only."
	userPrompt := "Technical pronunciation report:\n" + string(payload) + "\n\n" +
		"Rules:\n" +
		"- Compare expected sentence and transcript only when mode is exact_repeat.\n" +
		"- For free_speech, do not claim missing words because there is no exact target phrase.\n" +
		"- In exact_repeat mode, text match and missing/substituted words are strict gates; do not praise a bad transcript.\n" +
		"- If similarity is below 70 or missing_words/substitutions are present, overall_feedback must clearly say to repeat and must not say excellent, perfect, great, or well done.\n" +
		"- If similarity is below 88, do not praise the learner; explain the concrete mismatch and ask for another attempt.\n" +
		"- Missing words and substitutions are more important than fluency. A fluent wrong phrase is still a low score.\n" +
		"- Use confidence as the main pronunciation signal only after the exact target phrase is recognizably matched.\n" +
		"- Do not invent phoneme mistakes.\n" +
		"- Explain simply in the interface language from the report.\n" +
		"- Give practical speaking advice.\n" +
		"- Return JSON only with keys: model_correction (0-100, contributes only 5 percent), accent_strength (0-100), fluency (0-100), problem_words [{word, spoken, confidence, issue, tip}], tips [strings], overall_feedback.\n"
	return []chatMessage{
		{Role: "system", Content: renderAppPrompt("pronunciation.coach_report.system", systemPrompt, nil)},
		{
			Role: "user",
			Content: renderAppPrompt("pronunciation.coach_report.user", userPrompt, map[string]string{
				"technical_report_json": string(payload),
			}),
		},
	}
}

func parsePronunciationCoachJSON(raw string) (pronunciationCoachJSON, bool) {
	raw = normalizePronunciationJSON(raw)
	var parsed pronunciationCoachJSON
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return pronunciationCoachJSON{}, false
	}
	return parsed, true
}

func parsePronunciationAudioCoachJSON(raw string) (pronunciationAudioCoachJSON, bool) {
	raw = normalizePronunciationJSON(raw)
	var parsed pronunciationAudioCoachJSON
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return pronunciationAudioCoachJSON{}, false
	}
	return parsed, true
}

func normalizePronunciationJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)
	if start := strings.Index(raw, "{"); start >= 0 {
		raw = raw[start:]
	}
	if end := strings.LastIndex(raw, "}"); end >= 0 {
		raw = raw[:end+1]
	}
	return raw
}

func finalPronunciationScore(base float64, modelCorrection int) int {
	score := int(math.Round(base + float64(clampInt(modelCorrection, 0, 100))*0.05))
	return clampInt(score, 0, 100)
}

func capPronunciationScoreForConfidence(report pronunciationTechnicalReport, score int) int {
	limit := 100
	if report.Mode == string(pronunciationModeExact) {
		switch {
		case report.Similarity < 45:
			limit = 25
		case report.Similarity < 60:
			limit = 40
		case report.Similarity < 75:
			limit = 55
		case report.Similarity < 88:
			limit = 68
		case report.Similarity < 96:
			limit = 80
		case len(report.ProblemWords) > 0 || len(report.MissingWords) > 0 || len(report.Substitutions) > 0:
			limit = 84
		default:
			limit = 94
		}
		expectedWords := len(strings.Fields(report.Expected))
		if expectedWords > 1 && len(report.MissingWords)*2 >= expectedWords && limit > 55 {
			limit = 45
		}
		if expectedWords > 0 && len(report.Substitutions)*3 >= expectedWords && limit > 60 {
			limit = 60
		}
		if expectedWords > 0 && len(report.ProblemWords)*2 >= expectedWords && limit > 72 {
			limit = 72
		}
		if report.ConfidenceSource == "estimated" || report.ConfidenceSource == "none" {
			switch {
			case report.Similarity >= 98 && len(report.ProblemWords) == 0 && len(report.MissingWords) == 0:
				if limit > 78 {
					limit = 78
				}
			case report.Similarity >= 90:
				if limit > 68 {
					limit = 68
				}
			case report.Similarity >= 75:
				if limit > 55 {
					limit = 55
				}
			default:
				if limit > 35 {
					limit = 35
				}
			}
		}
	} else if report.ConfidenceSource == "estimated" || report.ConfidenceSource == "none" {
		limit = 78
	}
	if score > limit {
		return limit
	}
	return score
}

func strictPronunciationFeedback(user userState, report pronunciationTechnicalReport, score int) string {
	if report.Mode != string(pronunciationModeExact) || score >= 75 {
		return ""
	}
	if report.ConfidenceSource == "estimated" || report.ConfidenceSource == "none" {
		return ""
	}
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		if score < 50 {
			return "Оценка: " + itoa(score) + "/100. Фраза распознана слишком далеко от образца: повторите медленнее и отдельно проговорите пропущенные слова."
		}
		return "Оценка: " + itoa(score) + "/100. Есть заметные расхождения с образцом, поэтому попытку лучше повторить без длинных пауз и замен слов."
	}
	if score < 50 {
		return fmt.Sprintf("Score: %d/100. The transcript is too far from the target phrase, so repeat it more slowly and practice the missing words first.", score)
	}
	return fmt.Sprintf("Score: %d/100. There are clear differences from the target phrase, so repeat it with steadier rhythm and fewer word substitutions.", score)
}

func pronunciationAccentStrength(report pronunciationTechnicalReport, score int) int {
	confPenalty := (1 - clampFloat(report.AverageConfidence, 0, 1)) * 78
	simPenalty := float64(100-report.Similarity) * 0.18
	fluencyPenalty := float64(100-report.Fluency) * 0.12
	if report.Mode == string(pronunciationModeFree) {
		simPenalty = float64(100-score) * 0.12
	}
	return clampInt(int(math.Round(confPenalty+simPenalty+fluencyPenalty)), 0, 100)
}

func pronunciationFluency(spokenUnits []string, duration float64, level string) (int, float64) {
	if len(spokenUnits) == 0 {
		return 0, 0
	}
	if duration <= 0 {
		return 74, 0
	}
	wpm := float64(len(spokenUnits)) / duration * 60
	targetLow, targetHigh := 60.0, 155.0
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "A1":
		targetLow, targetHigh = 45, 120
	case "A2":
		targetLow, targetHigh = 55, 135
	case "B1", "B2":
		targetLow, targetHigh = 70, 165
	default:
		targetLow, targetHigh = 80, 190
	}
	score := 100.0
	if wpm < targetLow {
		score -= (targetLow - wpm) * 1.05
	}
	if wpm > targetHigh {
		score -= (wpm - targetHigh) * 0.75
	}
	return clampInt(int(math.Round(score)), 35, 100), wpm
}

func freeSpeechTextQuality(spokenUnits []string) int {
	switch {
	case len(spokenUnits) == 0:
		return 0
	case len(spokenUnits) == 1:
		return 58
	case len(spokenUnits) == 2:
		return 72
	default:
		return 88
	}
}

func alignPronunciationUnits(targetUnits []string, spokenUnits []string) []pronunciationAlignment {
	n, m := len(targetUnits), len(spokenUnits)
	if n == 0 {
		return nil
	}
	dp := make([][]float64, n+1)
	op := make([][]string, n+1)
	for i := range dp {
		dp[i] = make([]float64, m+1)
		op[i] = make([]string, m+1)
	}
	for i := 1; i <= n; i++ {
		dp[i][0] = float64(i)
		op[i][0] = "delete"
	}
	for j := 1; j <= m; j++ {
		dp[0][j] = float64(j)
		op[0][j] = "insert"
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			sim := wordSimilarity(targetUnits[i-1], spokenUnits[j-1])
			subCost := 1.0 - sim
			if targetUnits[i-1] != spokenUnits[j-1] {
				subCost += 0.18
			}
			best := dp[i-1][j-1] + subCost
			bestOp := "match"
			if targetUnits[i-1] != spokenUnits[j-1] {
				bestOp = "substitute"
			}
			if value := dp[i-1][j] + 1; value < best {
				best = value
				bestOp = "delete"
			}
			if value := dp[i][j-1] + 1; value < best {
				best = value
				bestOp = "insert"
			}
			dp[i][j] = best
			op[i][j] = bestOp
		}
	}
	var reversed []pronunciationAlignment
	i, j := n, m
	for i > 0 || j > 0 {
		current := op[i][j]
		switch current {
		case "match", "substitute":
			reversed = append(reversed, pronunciationAlignment{Target: targetUnits[i-1], Spoken: spokenUnits[j-1], Op: current, SpokenIndex: j - 1, Similarity: wordSimilarity(targetUnits[i-1], spokenUnits[j-1])})
			i--
			j--
		case "insert":
			reversed = append(reversed, pronunciationAlignment{Spoken: spokenUnits[j-1], Op: "insert", SpokenIndex: j - 1})
			j--
		default:
			reversed = append(reversed, pronunciationAlignment{Target: targetUnits[i-1], Op: "delete", SpokenIndex: -1})
			i--
		}
	}
	for left, right := 0, len(reversed)-1; left < right; left, right = left+1, right-1 {
		reversed[left], reversed[right] = reversed[right], reversed[left]
	}
	return reversed
}

func wordSimilarity(a string, b string) float64 {
	aRunes := compactWordRunes(a)
	bRunes := compactWordRunes(b)
	if len(aRunes) == 0 && len(bRunes) == 0 {
		return 1
	}
	longer := math.Max(float64(len(aRunes)), float64(len(bRunes)))
	if longer == 0 {
		return 0
	}
	return clampFloat(1-float64(editDistanceRunes(aRunes, bRunes))/longer, 0, 1)
}

func compactWordRunes(text string) []rune {
	var out []rune
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			out = append(out, r)
		}
	}
	return out
}

func fallbackPronunciationFeedback(user userState, report pronunciationTechnicalReport, score int) string {
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		if report.ConfidenceSource == "estimated" {
			return "\u041e\u0446\u0435\u043d\u043a\u0430: " + itoa(score) + "/100. \u0411\u0430\u043b\u043b \u043e\u0433\u0440\u0430\u043d\u0438\u0447\u0435\u043d: \u044f \u0441\u0432\u0435\u0440\u044f\u044e \u0444\u0430\u043a\u0442\u0438\u0447\u0435\u0441\u043a\u0438 \u0440\u0430\u0441\u043f\u043e\u0437\u043d\u0430\u043d\u043d\u044b\u0439 \u0442\u0435\u043a\u0441\u0442 \u0441 \u043e\u0431\u0440\u0430\u0437\u0446\u043e\u043c, \u0442\u0435\u043c\u043f, \u043f\u0440\u043e\u043f\u0443\u0441\u043a\u0438 \u0438 \u0437\u0430\u043c\u0435\u043d\u044b \u0441\u043b\u043e\u0432."
		}
		return "\u041e\u0446\u0435\u043d\u043a\u0430: " + itoa(score) + "/100. \u0413\u043b\u0430\u0432\u043d\u044b\u0439 \u0441\u0438\u0433\u043d\u0430\u043b - \u0441\u043e\u0432\u043f\u0430\u0434\u0435\u043d\u0438\u0435 \u0441 \u0444\u0440\u0430\u0437\u043e\u0439, \u0443\u0432\u0435\u0440\u0435\u043d\u043d\u043e\u0441\u0442\u044c \u0440\u0430\u0441\u043f\u043e\u0437\u043d\u0430\u0432\u0430\u043d\u0438\u044f, \u0440\u0438\u0442\u043c \u0438 \u043f\u043e\u043b\u043d\u043e\u0442\u0430 \u0444\u0440\u0430\u0437\u044b."
	}
	if report.ConfidenceSource == "estimated" {
		return fmt.Sprintf("Score: %d/100. The score is capped and uses the recognized transcript, rhythm, missing words, and substitutions.", score)
	}
	return fmt.Sprintf("Score: %d/100. The main signal is recognition confidence plus match to the target phrase.", score)
}

func fallbackPronunciationTips(user userState, report pronunciationTechnicalReport) []string {
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		if len(report.ProblemWords) > 0 {
			return []string{"\u041f\u043e\u0432\u0442\u043e\u0440\u0438 \u0441\u043b\u0430\u0431\u044b\u0435 \u0441\u043b\u043e\u0432\u0430 \u043e\u0442\u0434\u0435\u043b\u044c\u043d\u043e, \u043f\u043e\u0442\u043e\u043c \u0441\u043a\u0430\u0436\u0438 \u0432\u0441\u044e \u0444\u0440\u0430\u0437\u0443 \u0431\u0435\u0437 \u0434\u043b\u0438\u043d\u043d\u044b\u0445 \u043f\u0430\u0443\u0437."}
		}
		return []string{"\u0414\u0435\u0440\u0436\u0438 \u0440\u043e\u0432\u043d\u044b\u0439 \u0442\u0435\u043c\u043f \u0438 \u043d\u0435 \u043f\u0440\u043e\u0433\u043b\u0430\u0442\u044b\u0432\u0430\u0439 \u043a\u043e\u043d\u0446\u044b \u0441\u043b\u043e\u0432."}
	}
	if len(report.ProblemWords) > 0 {
		return []string{"Repeat the weak words alone first, then say the full phrase without long pauses."}
	}
	return []string{"Keep a steady rhythm and avoid swallowing word endings."}
}

func pronunciationAssessmentText(user userState, p *pronunciationAssessment) string {
	if p == nil {
		return ""
	}
	ru := normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru"
	title, accent, fluency, problems, tips := "Pronunciation", "Accent strength", "Fluency", "Problem words", "Tips"
	if ru {
		title = "\u041e\u0446\u0435\u043d\u043a\u0430 \u043f\u0440\u043e\u0438\u0437\u043d\u043e\u0448\u0435\u043d\u0438\u044f"
		accent = "\u0421\u0438\u043b\u0430 \u0430\u043a\u0446\u0435\u043d\u0442\u0430"
		fluency = "\u041f\u043b\u0430\u0432\u043d\u043e\u0441\u0442\u044c"
		problems = "\u0421\u043b\u0430\u0431\u044b\u0435 \u0441\u043b\u043e\u0432\u0430"
		tips = "\u0421\u043e\u0432\u0435\u0442"
	}
	var builder strings.Builder
	builder.WriteString(title)
	builder.WriteString(": ")
	builder.WriteString(itoa(p.Score))
	builder.WriteString("/100\n")
	builder.WriteString(accent + ": " + itoa(p.AccentStrength) + "/100\n")
	builder.WriteString(fluency + ": " + itoa(p.Fluency) + "/100")
	if len(p.ProblemWords) > 0 {
		builder.WriteString("\n")
		builder.WriteString(problems + ": ")
		var parts []string
		for _, item := range p.ProblemWords {
			label := strings.TrimSpace(item.Word)
			if item.Confidence > 0 {
				label += fmt.Sprintf(" %.2f", item.Confidence)
			}
			parts = append(parts, label)
			if len(parts) >= 4 {
				break
			}
		}
		builder.WriteString(strings.Join(parts, ", "))
	}
	if p.Feedback != "" {
		builder.WriteString("\n")
		builder.WriteString(strings.TrimSpace(p.Feedback))
	} else if len(p.Tips) > 0 {
		builder.WriteString("\n")
		builder.WriteString(tips + ": " + strings.TrimSpace(p.Tips[0]))
	}
	return builder.String()
}

func pronunciationAssessmentDTO(p *pronunciationAssessment) any {
	if p == nil {
		return nil
	}
	return p
}

func normalizeCoachProblemWords(coach []pronunciationProblemWord, local []pronunciationProblemWord) []pronunciationProblemWord {
	if len(coach) == 0 {
		return local
	}
	out := trimProblemWords(coach, 6)
	for i := range out {
		out[i].Word = strings.TrimSpace(out[i].Word)
		out[i].Spoken = strings.TrimSpace(out[i].Spoken)
		out[i].Tip = strings.TrimSpace(out[i].Tip)
		out[i].Issue = strings.TrimSpace(out[i].Issue)
		out[i].Confidence = roundConfidence(out[i].Confidence)
	}
	return out
}

func trimProblemWords(items []pronunciationProblemWord, limit int) []pronunciationProblemWord {
	var out []pronunciationProblemWord
	for _, item := range items {
		if strings.TrimSpace(item.Word) == "" {
			continue
		}
		out = append(out, item)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func trimStringSlice(items []string, limit int) []string {
	var out []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		out = append(out, item)
		if len(out) >= limit {
			break
		}
	}
	return out
}

func roundConfidence(value float64) float64 {
	return math.Round(clampFloat(value, 0, 1)*100) / 100
}

func clampFloat(value float64, minValue float64, maxValue float64) float64 {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func clampInt(value int, minValue int, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
