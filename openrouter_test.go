package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestOpenRouterGPT6ReservesBudgetForVisibleReply(t *testing.T) {
	for _, model := range []string{"openai/gpt-6-luna", "other/model"} {
		client := newOpenRouterClient("test", model, "", "", &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			reasoning, ok := payload["reasoning"].(map[string]any)
			if model == "openai/gpt-6-luna" {
				if !ok || reasoning["enabled"] != false {
					t.Error("GPT-6 text requests must disable reasoning to preserve short reply budgets")
				}
			} else if ok {
				t.Error("unrelated model reasoning policy changed")
			}
			return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{"choices":[{"message":{"content":"Clue"}}]}`))}, nil
		})})
		if _, err := client.complete(context.Background(), nil, 0.1, 80); err != nil {
			t.Fatal(err)
		}
	}
}

func TestOpenRouterCompletionHTTPErrorPreservesStatus(t *testing.T) {
	for _, body := range []string{`{"error":"denied"}`, `{"error":{"message":"denied"}}`, `Bad gateway`} {
		client := newOpenRouterClient("test", "test", "", "", &http.Client{Transport: roundTripFunc(func(_ *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 502, Body: io.NopCloser(strings.NewReader(body))}, nil
		})})
		_, err := client.complete(context.Background(), nil, 0, 20)
		if err == nil || !strings.Contains(err.Error(), "status 502") {
			t.Errorf("body %q: expected HTTP status in error, got %v", body, err)
		}
	}
}
