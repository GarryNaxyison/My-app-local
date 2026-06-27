package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type yooKassaClient struct {
	shopID    string
	secretKey string
	http      *http.Client
}

type yooKassaPayment struct {
	ID           string                 `json:"id"`
	Status       string                 `json:"status"`
	Paid         bool                   `json:"paid"`
	Amount       yooKassaAmount         `json:"amount"`
	Confirmation yooKassaConfirmation   `json:"confirmation"`
	Metadata     map[string]interface{} `json:"metadata"`
	Description  string                 `json:"description"`
}

type yooKassaAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type yooKassaConfirmation struct {
	Type            string `json:"type,omitempty"`
	ReturnURL       string `json:"return_url,omitempty"`
	ConfirmationURL string `json:"confirmation_url,omitempty"`
}

type yooKassaNotification struct {
	Type   string          `json:"type"`
	Event  string          `json:"event"`
	Object yooKassaPayment `json:"object"`
}

func newYooKassaClient(shopID string, secretKey string, httpClient *http.Client) *yooKassaClient {
	return &yooKassaClient{
		shopID:    shopID,
		secretKey: secretKey,
		http:      httpClient,
	}
}

func (c *yooKassaClient) createPremiumPayment(ctx context.Context, telegramID int64, firstName string, interfaceLanguage string, plan premiumPlan, returnURL string, channel string, idempotencyKeys ...string) (yooKassaPayment, error) {
	channel = strings.TrimSpace(channel)
	if channel == "" {
		channel = "telegram"
	}
	payload := map[string]any{
		"amount": map[string]string{
			"value":    fmt.Sprintf("%d.00", plan.RubPrice),
			"currency": "RUB",
		},
		"payment_method_data": map[string]string{
			"type": "sbp",
		},
		"confirmation": map[string]string{
			"type":       "redirect",
			"return_url": returnURL,
		},
		"capture":     true,
		"description": plan.Title + " - POLYGLOT AI",
		"metadata": map[string]string{
			"telegram_id":        strconv.FormatInt(telegramID, 10),
			"first_name":         firstName,
			"interface_language": normalizeInterfaceLanguage(interfaceLanguage),
			"product":            plan.Product,
			"channel":            channel,
		},
	}

	var payment yooKassaPayment
	idempotenceKey := ""
	if len(idempotencyKeys) > 0 {
		idempotenceKey = paymentIdempotencyToken("yookassa", telegramID, plan.Product, channel, idempotencyKeys[0])
	}
	if err := c.do(ctx, http.MethodPost, "https://api.yookassa.ru/v3/payments", payload, &payment, idempotenceKey); err != nil {
		return yooKassaPayment{}, err
	}
	if payment.Confirmation.ConfirmationURL == "" {
		return yooKassaPayment{}, fmt.Errorf("yookassa payment has no confirmation_url")
	}
	return payment, nil
}

func (c *yooKassaClient) getPayment(ctx context.Context, paymentID string) (yooKassaPayment, error) {
	var payment yooKassaPayment
	err := c.do(ctx, http.MethodGet, "https://api.yookassa.ru/v3/payments/"+paymentID, nil, &payment)
	return payment, err
}

func (c *yooKassaClient) do(ctx context.Context, method string, url string, payload any, target any, idempotenceKeys ...string) error {
	var body io.Reader
	if payload != nil {
		bytesBody, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(bytesBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.shopID, c.secretKey)
	req.Header.Set("Content-Type", "application/json")
	if method == http.MethodPost {
		idempotenceKey := ""
		if len(idempotenceKeys) > 0 {
			idempotenceKey = strings.TrimSpace(idempotenceKeys[0])
		}
		if idempotenceKey == "" {
			idempotenceKey = fmt.Sprintf("premium-%d", time.Now().UnixNano())
		}
		req.Header.Set("Idempotence-Key", idempotenceKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("yookassa %s failed: status %d: %s", method, resp.StatusCode, string(respBody))
	}
	if target == nil {
		return nil
	}
	return json.Unmarshal(respBody, target)
}

func metadataString(metadata map[string]interface{}, key string) (string, bool) {
	value, ok := metadata[key]
	if !ok {
		return "", false
	}
	typed, ok := value.(string)
	return typed, ok && typed != ""
}

func metadataInt64(metadata map[string]interface{}, key string) (int64, bool) {
	value, ok := metadata[key]
	if !ok {
		return 0, false
	}
	switch typed := value.(type) {
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 64)
		return parsed, err == nil
	case float64:
		return int64(typed), true
	default:
		return 0, false
	}
}
