package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

const (
	aiTutorStatusDraft             = "draft"
	aiTutorStatusPreflightRejected = "preflight_rejected"
	aiTutorStatusInTrial           = "in_trial"
	aiTutorStatusApproved          = "approved"
	aiTutorStatusRejected          = "rejected"
	aiTutorStatusArchived          = "archived"

	aiTutorSessionActive   = "active"
	aiTutorSessionComplete = "complete"

	aiTutorStageStoryIntro     = "story_intro"
	aiTutorStageRetell         = "retell"
	aiTutorStageQuestion1      = "question_1"
	aiTutorStageQuestion2      = "question_2"
	aiTutorStageQuestion3      = "question_3"
	aiTutorStageProduction     = "production"
	aiTutorStageLessonFeedback = "lesson_feedback"
	aiTutorStageReviewSchedule = "review_schedule"
	aiTutorStageComplete       = "complete"
)

type aiTutorLessonPayload struct {
	Title                  string                  `json:"title"`
	Level                  string                  `json:"level"`
	LevelBand              string                  `json:"level_band"`
	TargetLanguage         string                  `json:"target_language"`
	InterfaceLanguage      string                  `json:"interface_language"`
	Theme                  string                  `json:"theme"`
	LessonGoal             string                  `json:"lesson_goal"`
	GrammarFocus           aiTutorGrammarFocus     `json:"grammar_focus"`
	Story                  aiTutorStory            `json:"story"`
	Words                  []aiTutorWord           `json:"words"`
	RetellTask             aiTutorRetellTask       `json:"retell_task"`
	ComprehensionQuestions []aiTutorQuestion       `json:"comprehension_questions"`
	WordLearning           []aiTutorWordTask       `json:"word_learning"`
	WordRecall             []aiTutorWordTask       `json:"word_recall"`
	ProductionTask         aiTutorProductionTask   `json:"production_task"`
	ReviewOptions          []string                `json:"review_options"`
	QualitySelfCheck       aiTutorQualitySelfCheck `json:"quality_self_check"`
}

type aiTutorGrammarFocus struct {
	Name                      string `json:"name"`
	ShortExplanationInterface string `json:"short_explanation_interface"`
	ModelSentenceTarget       string `json:"model_sentence_target"`
}

type aiTutorStory struct {
	TextTarget      string `json:"text_target"`
	AudioTextTarget string `json:"audio_text_target"`
	SentenceCount   int    `json:"sentence_count"`
}

type aiTutorWord struct {
	ID                          string `json:"id"`
	Target                      string `json:"target"`
	InterfaceTranslation        string `json:"interface_translation"`
	PartOfSpeech                string `json:"part_of_speech"`
	ExampleSentenceTarget       string `json:"example_sentence_target"`
	ExampleTranslationInterface string `json:"example_translation_interface"`
	AudioTextTarget             string `json:"audio_text_target"`
	ExampleAudioTextTarget      string `json:"example_audio_text_target"`
	DifficultyNoteInterface     string `json:"difficulty_note_interface,omitempty"`
}

type aiTutorRetellTask struct {
	InstructionInterface string   `json:"instruction_interface"`
	MinSentences         int      `json:"min_sentences"`
	FeedbackRubric       []string `json:"feedback_rubric"`
}

type aiTutorQuestion struct {
	ID                    string   `json:"id"`
	QuestionTarget        string   `json:"question_target"`
	ExpectedPoints        []string `json:"expected_points"`
	FeedbackRuleInterface string   `json:"feedback_rule_interface"`
}

type aiTutorWordTask struct {
	WordID string `json:"word_id"`
}

type aiTutorProductionTask struct {
	InstructionInterface     string   `json:"instruction_interface"`
	RequiredWordCount        int      `json:"required_word_count"`
	SentenceCount            string   `json:"sentence_count"`
	EvaluationCriteria       []string `json:"evaluation_criteria"`
	RecommendationsInterface []string `json:"recommendations_interface"`
}

type aiTutorQualitySelfCheck struct {
	CEFRReason    string `json:"cefr_reason"`
	WhyReusable   string `json:"why_reusable"`
	DuplicateRisk string `json:"duplicate_risk"`
}

type aiTutorLessonRecord struct {
	ID                string
	LearningLanguage  string
	InterfaceLanguage string
	ExactLevel        string
	LevelBand         string
	Theme             string
	Status            string
	Payload           aiTutorLessonPayload
	Fingerprint       string
	PreflightScore    int
	PostScore         int
	CompletionCount   int
	AverageRating     float64
	CreatedAt         string
	UpdatedAt         string
}

type aiTutorSessionRecord struct {
	ID           string
	TelegramID   int64
	LessonID     string
	Surface      string
	CurrentStage string
	Status       string
	StartedAt    string
	CompletedAt  string
	UpdatedAt    string
}

type aiTutorAnswerRecord struct {
	SessionID    string
	Stage        string
	AnswerText   string
	AnswerJSON   string
	FeedbackJSON string
	Correct      bool
	CreatedAt    string
}

type aiTutorReviewRecord struct {
	ID           string
	TelegramID   int64
	LessonID     string
	DueAt        string
	IntervalCode string
	Status       string
	CreatedAt    string
	UpdatedAt    string
}

type aiTutorQualityCheckRecord struct {
	ID         string
	LessonID   string
	SessionID  string
	Kind       string
	Score      int
	Approved   bool
	IssuesJSON string
	ResultJSON string
	CreatedAt  string
}

func aiTutorLevelBand(level string) string {
	switch strings.ToUpper(strings.TrimSpace(level)) {
	case "A1-A2", "B1-B2", "C1-C2":
		return strings.ToUpper(strings.TrimSpace(level))
	}
	switch normalizeCEFRLevel(level) {
	case "B1", "B2":
		return "B1-B2"
	case "C1", "C2":
		return "C1-C2"
	default:
		return "A1-A2"
	}
}

func aiTutorCanonicalStages() []string {
	stages := []string{
		aiTutorStageStoryIntro,
		aiTutorStageRetell,
		aiTutorStageQuestion1,
		aiTutorStageQuestion2,
		aiTutorStageQuestion3,
	}
	for index := 1; index <= 6; index++ {
		stages = append(stages, aiTutorWordLearnStage(index))
	}
	for index := 1; index <= 6; index++ {
		stages = append(stages, aiTutorWordRecallStage(index))
	}
	stages = append(stages, aiTutorStageProduction, aiTutorStageLessonFeedback, aiTutorStageReviewSchedule, aiTutorStageComplete)
	return stages
}

func aiTutorWordLearnStage(index int) string {
	return "word_learn_" + itoa(index)
}

func aiTutorWordRecallStage(index int) string {
	return "word_recall_" + itoa(index)
}

func validateAITutorLessonPayload(lesson aiTutorLessonPayload) []string {
	var issues []string
	if strings.TrimSpace(lesson.Title) == "" {
		issues = append(issues, "title is required")
	}
	if strings.TrimSpace(lesson.Level) == "" {
		issues = append(issues, "level is required")
	}
	if strings.TrimSpace(lesson.TargetLanguage) == "" || strings.TrimSpace(lesson.InterfaceLanguage) == "" {
		issues = append(issues, "target and interface languages are required")
	}
	if strings.TrimSpace(lesson.LessonGoal) == "" {
		issues = append(issues, "lesson goal is required")
	}
	if len(lesson.Words) != 6 {
		issues = append(issues, "lesson must contain exactly 6 words")
	}
	if len(lesson.ComprehensionQuestions) != 3 {
		issues = append(issues, "lesson must contain exactly 3 questions")
	}
	storySentenceCount := countAITutorStorySentences(lesson.Story.TextTarget, lesson.Story.SentenceCount)
	if storySentenceCount < 5 || storySentenceCount > 7 {
		issues = append(issues, "story must contain 5-7 story sentences")
	}
	if strings.TrimSpace(lesson.Story.AudioTextTarget) == "" {
		issues = append(issues, "story audio is required")
	}
	if len(lesson.WordLearning) != 6 {
		issues = append(issues, "word learning must contain exactly 6 tasks")
	}
	if len(lesson.WordRecall) != 6 {
		issues = append(issues, "word recall must contain exactly 6 tasks")
	}
	if lesson.ProductionTask.RequiredWordCount < 3 {
		issues = append(issues, "production task must require at least 3 words")
	}
	if len(aiTutorTrimmedStrings(lesson.ProductionTask.RecommendationsInterface)) == 0 {
		issues = append(issues, "final recommendations are required")
	}
	if !hasAITutorReviewOptions(lesson.ReviewOptions) {
		issues = append(issues, "review options must include tomorrow, 3_days, 1_week, no_review")
	}
	for index, word := range lesson.Words {
		if strings.TrimSpace(word.ID) == "" || strings.TrimSpace(word.Target) == "" || strings.TrimSpace(word.InterfaceTranslation) == "" {
			issues = append(issues, "word "+itoa(index+1)+" requires id, target, and interface translation")
		}
		if strings.TrimSpace(word.AudioTextTarget) == "" {
			issues = append(issues, "word "+itoa(index+1)+" audio is required")
		}
		if strings.TrimSpace(word.ExampleAudioTextTarget) == "" {
			issues = append(issues, "word "+itoa(index+1)+" example audio is required")
		}
	}
	for index, question := range lesson.ComprehensionQuestions {
		if strings.TrimSpace(question.ID) == "" || strings.TrimSpace(question.QuestionTarget) == "" || len(question.ExpectedPoints) == 0 {
			issues = append(issues, "question "+itoa(index+1)+" requires id, question, and expected points")
		}
	}
	return issues
}

func aiTutorTrimmedStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func normalizeAITutorLessonPayload(lesson aiTutorLessonPayload) aiTutorLessonPayload {
	if len(lesson.Words) == 6 {
		if len(lesson.WordLearning) != 6 {
			lesson.WordLearning = aiTutorWordTasksFromWords(lesson.Words)
		}
		if len(lesson.WordRecall) != 6 {
			lesson.WordRecall = aiTutorWordTasksFromWords(lesson.Words)
		}
	}
	if !hasAITutorReviewOptions(lesson.ReviewOptions) {
		lesson.ReviewOptions = []string{"tomorrow", "3_days", "1_week", "no_review"}
	}
	return lesson
}

func aiTutorWordTasksFromWords(words []aiTutorWord) []aiTutorWordTask {
	tasks := make([]aiTutorWordTask, 0, len(words))
	for _, word := range words {
		tasks = append(tasks, aiTutorWordTask{WordID: strings.TrimSpace(word.ID)})
	}
	return tasks
}

func countAITutorStorySentences(text string, declaredCount int) int {
	text = strings.TrimSpace(text)
	count := 0
	for _, char := range text {
		switch char {
		case '.', '!', '?':
			count++
		}
	}
	if count > 0 {
		return count
	}
	if declaredCount > 0 {
		return declaredCount
	}
	if len(strings.Fields(text)) > 0 {
		return 1
	}
	return 0
}

func hasAITutorReviewOptions(options []string) bool {
	seen := map[string]bool{}
	for _, option := range options {
		seen[strings.TrimSpace(option)] = true
	}
	return seen["tomorrow"] && seen["3_days"] && seen["1_week"] && seen["no_review"]
}

func aiTutorFingerprint(lesson aiTutorLessonPayload) string {
	seed := strings.ToLower(strings.Join([]string{
		lesson.LevelBand,
		lesson.TargetLanguage,
		lesson.InterfaceLanguage,
		lesson.Theme,
		lesson.Story.TextTarget,
	}, "|"))
	sum := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(sum[:])
}
