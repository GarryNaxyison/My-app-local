package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTelegramCallAcceptsBooleanAndMessageResults(t *testing.T) {
	for _, tc := range []struct {
		method, result string
		id             int64
	}{
		{"answerCallbackQuery", "true", 0},
		{"deleteMessage", "true", 0},
		{"sendChatAction", "true", 0},
		{"editMessageText", "true", 0},
		{"sendMessage", `{"message_id":42}`, 42},
	} {
		t.Run(tc.method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, `{"ok":true,"result":%s}`, tc.result)
			}))
			defer server.Close()
			c := &telegramClient{baseURL: server.URL, http: server.Client()}
			id, err := c.callMessageID(context.Background(), tc.method, nil)
			if err != nil || id != tc.id {
				t.Fatalf("id=%d err=%v; want %d, nil", id, err, tc.id)
			}
		})
	}
}
