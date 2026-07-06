package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
)

const idempotencyKeyHeader = "Idempotency-Key"

func requestIdempotencyKey(r *http.Request) string {
	return normalizeIdempotencyKey(r.Header.Get(idempotencyKeyHeader))
}

func normalizeIdempotencyKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" || len(key) > 128 {
		return ""
	}
	return key
}

func paymentIdempotencyToken(provider string, telegramID int64, product string, channel string, key string) string {
	key = normalizeIdempotencyKey(key)
	if key == "" {
		return ""
	}
	seed := strings.Join([]string{
		strings.TrimSpace(provider),
		strconv.FormatInt(telegramID, 10),
		strings.TrimSpace(product),
		strings.TrimSpace(channel),
		key,
	}, "|")
	sum := sha256.Sum256([]byte(seed))
	return strings.TrimSpace(provider) + "-" + hex.EncodeToString(sum[:])[:32]
}
