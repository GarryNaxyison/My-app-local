package main

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

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
