package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

const aiTutorPromotionThreshold = 88

type aiTutorCompletionClient interface {
	complete(ctx context.Context, messages []chatMessage, temperature float64, maxTokens int) (string, error)
}

type aiTutorEngine struct {
	store store
	ai    aiTutorCompletionClient
	now   func() time.Time
}

type aiTutorStep struct {
	Stage       string               `json:"stage"`
	Kind        string               `json:"kind"`
	Title       string               `json:"title"`
	Instruction string               `json:"instruction"`
	Lesson      aiTutorLessonPayload `json:"lesson,omitempty"`
	Word        *aiTutorWord         `json:"word,omitempty"`
	Question    *aiTutorQuestion     `json:"question,omitempty"`
	Options     []string             `json:"options,omitempty"`
}

type aiTutorFeedback struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	JSON    string `json:"json,omitempty"`
}

type aiTutorResult struct {
	Session  aiTutorSessionRecord `json:"session"`
	Lesson   aiTutorLessonRecord  `json:"lesson"`
	NextStep aiTutorStep          `json:"next_step"`
	Feedback aiTutorFeedback      `json:"feedback,omitempty"`
}

type aiTutorSubmitInput struct {
	Text   string `json:"text"`
	Choice string `json:"choice"`
}

type aiTutorQualityResult struct {
	Approved bool     `json:"approved"`
	Score    int      `json:"score"`
	Issues   []string `json:"critical_issues"`
}

func newAITutorEngine(store store, ai aiTutorCompletionClient) *aiTutorEngine {
	return &aiTutorEngine{store: store, ai: ai, now: time.Now}
}

func (e *aiTutorEngine) Start(ctx context.Context, user userState, surface string) (aiTutorResult, error) {
	language := normalizeLearningLanguage(user.LearningLanguage)
	interfaceLanguage := normalizeInterfaceLanguage(user.InterfaceLanguage)
	level := normalizeCEFRLevel(user.Level)
	levelBand := aiTutorLevelBand(level)
	if approved, ok, err := e.store.findApprovedAITutorLesson(language, interfaceLanguage, levelBand, user.TelegramID); err != nil {
		return aiTutorResult{}, err
	} else if ok {
		return e.createSessionForLesson(user, surface, approved)
	}
	lesson, err := e.generateLesson(ctx, user)
	if err != nil {
		return aiTutorResult{}, err
	}
	if err := e.store.saveAITutorLesson(lesson); err != nil {
		return aiTutorResult{}, err
	}
	return e.createSessionForLesson(user, surface, lesson)
}

func (e *aiTutorEngine) Submit(ctx context.Context, user userState, sessionID string, input aiTutorSubmitInput) (aiTutorResult, error) {
	session, ok, err := e.store.getAITutorSession(sessionID)
	if err != nil {
		return aiTutorResult{}, err
	}
	if !ok || session.TelegramID != user.TelegramID {
		return aiTutorResult{}, errors.New("AI Tutor session not found")
	}
	lesson, ok, err := e.store.getAITutorLesson(session.LessonID)
	if err != nil {
		return aiTutorResult{}, err
	}
	if !ok {
		return aiTutorResult{}, errors.New("AI Tutor lesson not found")
	}
	if session.Status == aiTutorSessionComplete {
		return aiTutorResult{Session: session, Lesson: lesson, NextStep: aiTutorBuildStep(lesson.Payload, aiTutorStageComplete)}, nil
	}

	current := session.CurrentStage
	answerText := strings.TrimSpace(input.Text)
	if strings.TrimSpace(input.Choice) != "" {
		answerText = strings.TrimSpace(input.Choice)
	}
	if err := e.store.saveAITutorAnswer(aiTutorAnswerRecord{
		SessionID:  session.ID,
		Stage:      current,
		AnswerText: answerText,
		Correct:    true,
		CreatedAt:  formatDBTime(e.now().UTC()),
	}); err != nil {
		return aiTutorResult{}, err
	}

	next := aiTutorNextStage(current)
	status := aiTutorSessionActive
	completedAt := ""
	if current == aiTutorStageReviewSchedule || next == aiTutorStageComplete {
		status = aiTutorSessionComplete
		next = aiTutorStageComplete
		completedAt = formatDBTime(e.now().UTC())
		if input.Choice != "" && input.Choice != "no_review" {
			if err := e.store.scheduleAITutorReview(user.TelegramID, lesson.ID, aiTutorReviewDueAt(e.now().UTC(), input.Choice), input.Choice); err != nil {
				return aiTutorResult{}, err
			}
		}
	}
	if err := e.store.updateAITutorSessionStage(session.ID, next, status, completedAt); err != nil {
		return aiTutorResult{}, err
	}
	session.CurrentStage = next
	session.Status = status
	session.CompletedAt = completedAt
	session.UpdatedAt = formatDBTime(e.now().UTC())
	return aiTutorResult{
		Session:  session,
		Lesson:   lesson,
		NextStep: aiTutorBuildStep(lesson.Payload, next),
		Feedback: aiTutorFeedback{OK: true, Message: "Saved."},
	}, nil
}

func (e *aiTutorEngine) createSessionForLesson(user userState, surface string, lesson aiTutorLessonRecord) (aiTutorResult, error) {
	now := formatDBTime(e.now().UTC())
	session := aiTutorSessionRecord{
		ID:           aiTutorNewID("aits"),
		TelegramID:   user.TelegramID,
		LessonID:     lesson.ID,
		Surface:      strings.TrimSpace(surface),
		CurrentStage: aiTutorStageStoryIntro,
		Status:       aiTutorSessionActive,
		StartedAt:    now,
		UpdatedAt:    now,
	}
	if session.Surface == "" {
		session.Surface = "web"
	}
	if err := e.store.createAITutorSession(session); err != nil {
		return aiTutorResult{}, err
	}
	return aiTutorResult{
		Session:  session,
		Lesson:   lesson,
		NextStep: aiTutorBuildStep(lesson.Payload, aiTutorStageStoryIntro),
	}, nil
}

func (e *aiTutorEngine) generateLesson(ctx context.Context, user userState) (aiTutorLessonRecord, error) {
	if e.ai == nil {
		return aiTutorLessonRecord{}, errors.New("ai tutor client is not configured")
	}
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	raw, err := e.ai.complete(ctx, aiTutorLessonGenerationPrompt(language, interfaceLanguage, user.Level, user.LearningFocus, nil), 0.4, 4000)
	if err != nil {
		return aiTutorLessonRecord{}, err
	}
	var payload aiTutorLessonPayload
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &payload); err != nil {
		return aiTutorLessonRecord{}, err
	}
	if issues := validateAITutorLessonPayload(payload); len(issues) > 0 {
		return aiTutorLessonRecord{}, errors.New(strings.Join(issues, "; "))
	}
	qualityRaw, err := e.ai.complete(ctx, aiTutorQualityPrompt(language, interfaceLanguage, payload, nil, "preflight"), 0.1, 1200)
	if err != nil {
		return aiTutorLessonRecord{}, err
	}
	quality := aiTutorParseQualityResult(qualityRaw)
	status := aiTutorStatusInTrial
	if !quality.Approved || quality.Score < aiTutorPromotionThreshold {
		status = aiTutorStatusPreflightRejected
	}
	now := formatDBTime(e.now().UTC())
	record := aiTutorLessonRecord{
		ID:                aiTutorNewID("aitl"),
		LearningLanguage:  normalizeLearningLanguage(user.LearningLanguage),
		InterfaceLanguage: normalizeInterfaceLanguage(user.InterfaceLanguage),
		ExactLevel:        normalizeCEFRLevel(user.Level),
		LevelBand:         aiTutorLevelBand(user.Level),
		Theme:             payload.Theme,
		Status:            status,
		Payload:           payload,
		Fingerprint:       aiTutorFingerprint(payload),
		PreflightScore:    quality.Score,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := e.store.saveAITutorQualityCheck(aiTutorQualityCheckRecord{
		ID:         aiTutorNewID("aitq"),
		LessonID:   record.ID,
		Kind:       "preflight",
		Score:      quality.Score,
		Approved:   quality.Approved,
		IssuesJSON: aiTutorJSON(quality.Issues),
		ResultJSON: qualityRaw,
		CreatedAt:  now,
	}); err != nil {
		return aiTutorLessonRecord{}, err
	}
	if status == aiTutorStatusPreflightRejected {
		return aiTutorLessonRecord{}, errors.New("generated lesson failed preflight")
	}
	return record, nil
}

func aiTutorBuildStep(lesson aiTutorLessonPayload, stage string) aiTutorStep {
	step := aiTutorStep{Stage: stage, Lesson: lesson}
	switch stage {
	case aiTutorStageStoryIntro:
		step.Kind = "story"
		step.Title = lesson.Title
		step.Instruction = lesson.LessonGoal
	case aiTutorStageRetell:
		step.Kind = "free_text"
		step.Title = "Retell"
		step.Instruction = lesson.RetellTask.InstructionInterface
	case aiTutorStageQuestion1, aiTutorStageQuestion2, aiTutorStageQuestion3:
		index := map[string]int{aiTutorStageQuestion1: 0, aiTutorStageQuestion2: 1, aiTutorStageQuestion3: 2}[stage]
		step.Kind = "free_text"
		step.Title = "Question " + itoa(index+1)
		if index < len(lesson.ComprehensionQuestions) {
			question := lesson.ComprehensionQuestions[index]
			step.Question = &question
			step.Instruction = question.QuestionTarget
		}
	case aiTutorStageProduction:
		step.Kind = "free_text"
		step.Title = "Your sentences"
		step.Instruction = lesson.ProductionTask.InstructionInterface
	case aiTutorStageLessonFeedback:
		step.Kind = "rating"
		step.Title = "Lesson feedback"
		step.Instruction = "How did this lesson feel?"
		step.Options = []string{"easy", "good", "hard", "bad"}
	case aiTutorStageReviewSchedule:
		step.Kind = "review"
		step.Title = "Review"
		step.Instruction = "Choose review time."
		step.Options = []string{"tomorrow", "3_days", "1_week", "no_review"}
	case aiTutorStageComplete:
		step.Kind = "complete"
		step.Title = "Complete"
		step.Instruction = "Lesson complete."
	default:
		if index, ok := aiTutorStageIndex(stage, "word_learn_"); ok && index < len(lesson.Words) {
			word := lesson.Words[index]
			step.Kind = "word_learn"
			step.Title = word.Target
			step.Instruction = word.InterfaceTranslation
			step.Word = &word
		} else if index, ok := aiTutorStageIndex(stage, "word_recall_"); ok && index < len(lesson.Words) {
			word := lesson.Words[index]
			step.Kind = "word_recall"
			step.Title = word.InterfaceTranslation
			step.Instruction = "Recall the target word."
			step.Word = &word
		}
	}
	return step
}

func aiTutorNextStage(stage string) string {
	stages := aiTutorCanonicalStages()
	for index, current := range stages {
		if current == stage && index+1 < len(stages) {
			return stages[index+1]
		}
	}
	return aiTutorStageComplete
}

func aiTutorReviewDueAt(now time.Time, interval string) string {
	switch strings.TrimSpace(interval) {
	case "tomorrow":
		return formatDBTime(now.AddDate(0, 0, 1))
	case "3_days":
		return formatDBTime(now.AddDate(0, 0, 3))
	case "1_week":
		return formatDBTime(now.AddDate(0, 0, 7))
	default:
		return ""
	}
}

func aiTutorStageIndex(stage string, prefix string) (int, bool) {
	if !strings.HasPrefix(stage, prefix) {
		return 0, false
	}
	value, err := strconv.Atoi(strings.TrimPrefix(stage, prefix))
	if err != nil || value <= 0 {
		return 0, false
	}
	return value - 1, true
}

func aiTutorParseQualityResult(raw string) aiTutorQualityResult {
	var result aiTutorQualityResult
	_ = json.Unmarshal([]byte(strings.TrimSpace(raw)), &result)
	return result
}

func aiTutorNewID(prefix string) string {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return prefix + "-" + itoa(int(time.Now().UnixNano()))
	}
	return prefix + "-" + hex.EncodeToString(bytes[:])
}
