package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"
)

type tutorEvaluation struct {
	falseRejections, falseAcceptances, invalid, unstable int
}

func (e *tutorEvaluation) add(expected bool, responses []string) {
	accepted, rejected := false, false
	for _, raw := range responses {
		var fields map[string]any
		err := json.Unmarshal([]byte(raw), &fields)
		_, ok := fields["ok"].(bool)
		_, correct := fields["correct"].(bool)
		if err != nil || (!ok && !correct) {
			e.invalid++
			continue
		}
		actual := aiTutorCheckerAccepted(raw)
		accepted = accepted || actual
		rejected = rejected || !actual
		if expected && !actual {
			e.falseRejections++
		}
		if !expected && actual {
			e.falseAcceptances++
		}
	}
	if accepted && rejected {
		e.unstable++
	}
}

func TestTutorEvaluationCountsDisagreements(t *testing.T) {
	var result tutorEvaluation
	result.add(true, []string{`{"ok":true}`, `{"ok":false}`, `{}`})
	result.add(false, []string{`{"correct":true}`, `{"correct":true}`})
	if result != (tutorEvaluation{falseRejections: 1, falseAcceptances: 2, invalid: 1, unstable: 1}) {
		t.Fatalf("unexpected evaluation: %+v", result)
	}
}

func TestTutorQualityLive(t *testing.T) {
	if os.Getenv("NERIVA_RUN_TUTOR_EVAL") != "1" {
		t.Skip("set NERIVA_RUN_TUTOR_EVAL=1 to evaluate the live model")
	}
	key, model := os.Getenv("OPENROUTER_API_KEY"), os.Getenv("NERIVA_TUTOR_EVAL_MODEL")
	if key == "" || model == "" {
		t.Fatal("set OPENROUTER_API_KEY and NERIVA_TUTOR_EVAL_MODEL")
	}
	client := newOpenRouterClient(key, model, "https://neriva.ru", "NERIVA evaluation", &http.Client{Timeout: 20 * time.Second})
	engine := newAITutorEngine(nil, client)
	user := userState{LearningLanguage: "en", InterfaceLanguage: "ru", Level: "A2"}
	lesson := aiTutorLessonPayload{
		Title: "Train to York", Level: "A2", TargetLanguage: "en", InterfaceLanguage: "ru",
		Story: aiTutorStory{TextTarget: "Mia took the train to York. It left at nine in the morning."},
		ComprehensionQuestions: []aiTutorQuestion{{
			ID: "departure", QuestionTarget: "When did the train leave?",
			ExpectedPoints: []string{"at nine in the morning"},
		}},
	}
	cases := []struct {
		name, answer string
		accept       bool
	}{
		{"exact", "It left at nine in the morning.", true},
		{"paraphrase", "The train departed at 9 a.m.", true},
		{"short", "At nine in the morning.", true},
		{"wrong-time", "It left at nine in the evening.", false},
		{"off-topic", "Mia went to York.", false},
		{"contradiction", "It did not leave at nine in the morning.", false},
	}
	var result tutorEvaluation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	for _, item := range cases {
		var responses []string
		for i := 0; i < 3; i++ {
			raw, err := engine.checkFreeTextStage(ctx, user, lesson, aiTutorStageQuestion1, item.answer)
			if err != nil {
				var schemaError *json.UnmarshalTypeError
				if errors.As(err, &schemaError) {
					t.Fatalf("%s run %d: provider response field %q has type %s, expected %s", item.name, i+1, schemaError.Field, schemaError.Value, schemaError.Type)
				}
				t.Fatalf("%s run %d: model request failed (%T)", item.name, i+1, err)
			}
			responses = append(responses, raw)
		}
		result.add(item.accept, responses)
	}
	t.Logf("cases=%d runs=%d false_rejections=%d false_acceptances=%d invalid=%d unstable_cases=%d",
		len(cases), 3*len(cases), result.falseRejections, result.falseAcceptances, result.invalid, result.unstable)
	if result.invalid+result.falseAcceptances+result.falseRejections+result.unstable > 0 {
		t.Error("the checker disagreed with the reference answers")
	}
}
