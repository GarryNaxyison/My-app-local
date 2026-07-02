package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	rollyPayChannelTelegram = "bot"
	rollyPayChannelWeb      = "site"
)

type rollyPayClient struct {
	cashboxID string
	apiKey    string
	baseURL   string
	path      string
	http      *http.Client
}

type rollyPayPayment struct {
	ID     string
	URL    string
	Status string
	Raw    map[string]any
}

type rollyPayWebhookPayload struct {
	ID        string         `json:"id"`
	PaymentID string         `json:"payment_id"`
	OrderID   string         `json:"order_id"`
	Status    string         `json:"status"`
	Amount    any            `json:"amount"`
	Currency  string         `json:"currency"`
	Metadata  map[string]any `json:"metadata"`
	Data      map[string]any `json:"data"`
	Object    map[string]any `json:"object"`
}

func newRollyPayClient(cashboxID string, apiKey string, baseURL string, path string, httpClient *http.Client) *rollyPayClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &rollyPayClient{
		cashboxID: strings.TrimSpace(cashboxID),
		apiKey:    strings.TrimSpace(apiKey),
		baseURL:   strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		path:      strings.TrimSpace(path),
		http:      httpClient,
	}
}

func (c *rollyPayClient) createPremiumPayment(ctx context.Context, user userState, plan premiumPlan, channel string, successURL string, failURL string, now time.Time, idempotencyKeys ...string) (rollyPayPayment, error) {
	if c == nil || c.cashboxID == "" || c.apiKey == "" {
		return rollyPayPayment{}, fmt.Errorf("RollyPay is not configured")
	}
	endpoint, err := c.endpoint()
	if err != nil {
		return rollyPayPayment{}, err
	}
	orderID := rollyPayOrderID(plan.Product, user.TelegramID, now)
	if len(idempotencyKeys) > 0 {
		if stableOrderID := paymentIdempotencyToken("rollypay", user.TelegramID, plan.Product, channel, idempotencyKeys[0]); stableOrderID != "" {
			orderID = stableOrderID
		}
	}
	payload := map[string]any{
		"amount":           fmt.Sprintf("%.2f", float64(plan.RubPrice)),
		"payment_currency": "RUB",
		"order_id":         orderID,
		"terminal_id":      c.cashboxID,
		"description":      plan.Title + " - NERIVA",
		"customer_id":      strconv.FormatInt(user.TelegramID, 10),
		"metadata": map[string]any{
			"telegram_id":        strconv.FormatInt(user.TelegramID, 10),
			"first_name":         user.FirstName,
			"interface_language": normalizeInterfaceLanguage(user.InterfaceLanguage),
			"product":            plan.Product,
			"channel":            channel,
		},
	}
	if strings.TrimSpace(successURL) != "" {
		payload["success_redirect_url"] = strings.TrimSpace(successURL)
	}
	if strings.TrimSpace(failURL) != "" {
		payload["fail_redirect_url"] = strings.TrimSpace(failURL)
	}
	var result map[string]any
	if err := c.do(ctx, http.MethodPost, endpoint, payload, &result); err != nil {
		return rollyPayPayment{}, err
	}
	payment := rollyPayPayment{
		ID:     firstString(result, "id", "payment_id", "uuid"),
		URL:    firstString(result, "url", "pay_url", "payment_url", "confirmation_url", "redirect_url"),
		Status: firstString(result, "status", "state"),
		Raw:    result,
	}
	if payment.ID == "" {
		payment.ID = orderID
	}
	if payment.URL == "" {
		return rollyPayPayment{}, fmt.Errorf("rollypay payment has no payment url")
	}
	return payment, nil
}

func (c *rollyPayClient) endpoint() (string, error) {
	base := c.baseURL
	if base == "" {
		base = "https://rollypay.io"
	}
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	p := c.path
	if p == "" {
		p = "/api/v1/payments"
	}
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p, nil
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	basePath := strings.TrimRight(u.Path, "/")
	if basePath != "" && strings.HasPrefix(p, basePath+"/") {
		u.Path = p
	} else {
		u.Path = basePath + p
	}
	return u.String(), nil
}

func (c *rollyPayClient) do(ctx context.Context, method string, endpoint string, payload any, target any) error {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
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
		return fmt.Errorf("rollypay %s failed: status %d: %s", method, resp.StatusCode, string(respBody))
	}
	if target == nil {
		return nil
	}
	return json.Unmarshal(respBody, target)
}

func rollyPayOrderID(product string, telegramID int64, now time.Time) string {
	return fmt.Sprintf("%s_%d_%d", strings.TrimSpace(product), telegramID, now.Unix())
}

func (b *bot) premiumPlanFromRollyPayOrderID(orderID string) (premiumPlan, int64, int64, bool) {
	return b.premiumPlanPaymentFromPayload(orderID)
}

func (b *bot) handleRollyPayWebhook(w http.ResponseWriter, r *http.Request, channel string) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if !rollyPayWebhookSignatureAllowed(r, body, b.cfg.RollyPayWebhookSecret) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var payload rollyPayWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	flat := payload.flatten()
	status := strings.ToLower(firstString(flat, "status", "state", "payment_status"))
	if !rollyPayPaidStatus(status) {
		w.WriteHeader(http.StatusOK)
		return
	}
	orderID := firstString(flat, "order_id", "external_id", "merchant_order_id", "invoice_id")
	plan, telegramID, ok := b.premiumPlanFromRollyPayPayload(flat)
	if !ok || telegramID == 0 {
		log.Printf("rollypay webhook has unknown payment identity: order_id=%q telegram_id=%q product=%q", orderID, firstString(flat, "telegram_id", "user_id", "customer_id"), firstString(flat, "product", "plan"))
		w.WriteHeader(http.StatusOK)
		return
	}
	if amount := amountAsFloat(firstAny(flat, "amount", "amount_rub", "sum", "value")); amount > 0 && amount < float64(plan.RubPrice) {
		log.Printf("rollypay payment for order %s has unexpected amount %.2f", orderID, amount)
		w.WriteHeader(http.StatusOK)
		return
	}
	paymentID := firstString(flat, "payment_id", "id", "uuid", "transaction_id")
	if paymentID == "" {
		paymentID = orderID
	}
	if paymentID == "" {
		paymentID = fmt.Sprintf("%d:%s", telegramID, plan.Product)
	}
	chargeID := "rollypay:" + paymentID
	until, err := b.store.extendPremium(telegramID, chargeID, plan.Duration, plan.Tier)
	if err != nil {
		log.Printf("failed to activate premium for rollypay payment %s: %v", paymentID, err)
		http.Error(w, "activation failed", http.StatusInternalServerError)
		return
	}
	if _, err := b.store.creditReferralPurchase(telegramID, chargeID, rubToKopecks(plan.RubPrice)); err != nil {
		log.Printf("failed to credit referral rewards for rollypay payment %s: %v", paymentID, err)
		http.Error(w, "referral reward failed", http.StatusInternalServerError)
		return
	}
	user, _ := b.store.getOrCreateUser(telegramID, "")
	plan = localizedPremiumPlan(user, plan)
	if channel != rollyPayChannelWeb && b.telegram != nil {
		_ = b.telegram.sendMarkdownMessageWithCopy(r.Context(), telegramID, premiumActivationMarkdown(user, plan, until.Format("2006-01-02"), false), ui(user))
	}
	b.notifyOpsPayment(r.Context(), "RollyPay "+channel, user, plan, paymentID, fmt.Sprintf("%d RUB", plan.RubPrice), until)
	w.WriteHeader(http.StatusOK)
}

func (b *bot) premiumPlanFromRollyPayPayload(flat map[string]any) (premiumPlan, int64, bool) {
	orderID := firstString(flat, "order_id", "external_id", "merchant_order_id", "invoice_id")
	if orderID != "" {
		plan, telegramID, _, ok := b.premiumPlanFromRollyPayOrderID(orderID)
		if ok && telegramID != 0 {
			return plan, telegramID, true
		}
	}
	product := firstString(flat, "product", "plan", "tariff")
	telegramIDRaw := firstString(flat, "telegram_id", "user_id", "customer_id")
	if product == "" || telegramIDRaw == "" {
		return premiumPlan{}, 0, false
	}
	telegramID, err := strconv.ParseInt(telegramIDRaw, 10, 64)
	if err != nil || telegramID == 0 {
		return premiumPlan{}, 0, false
	}
	plan, ok := b.premiumPlan(product)
	if !ok {
		return premiumPlan{}, 0, false
	}
	return plan, telegramID, true
}

func rollyPayWebhookSignatureAllowed(r *http.Request, body []byte, secret string) bool {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return true
	}
	timestamp := strings.TrimSpace(r.Header.Get("X-Timestamp"))
	signature := strings.TrimSpace(r.Header.Get("X-Signature"))
	if timestamp != "" && signature != "" {
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(timestamp))
		mac.Write([]byte("."))
		mac.Write(body)
		expected := hex.EncodeToString(mac.Sum(nil))
		return hmac.Equal([]byte(expected), []byte(strings.ToLower(signature)))
	}
	if candidate := strings.TrimSpace(r.URL.Query().Get("secret")); candidate != "" {
		return hmac.Equal([]byte(candidate), []byte(secret))
	}
	return false
}

func rollyPayPaidStatus(status string) bool {
	status = strings.ToLower(strings.TrimSpace(status))
	return status == "paid" || status == "succeeded" || status == "success" || status == "completed" || status == "confirmed"
}

func (payload rollyPayWebhookPayload) flatten() map[string]any {
	result := map[string]any{}
	add := func(source map[string]any) {
		for key, value := range source {
			result[key] = value
		}
	}
	if payload.Data != nil {
		add(payload.Data)
	}
	if payload.Object != nil {
		add(payload.Object)
	}
	if payload.Metadata != nil {
		result["metadata"] = payload.Metadata
		add(payload.Metadata)
	}
	if payload.ID != "" {
		result["id"] = payload.ID
	}
	if payload.PaymentID != "" {
		result["payment_id"] = payload.PaymentID
	}
	if payload.OrderID != "" {
		result["order_id"] = payload.OrderID
	}
	if payload.Status != "" {
		result["status"] = payload.Status
	}
	if payload.Amount != nil {
		result["amount"] = payload.Amount
	}
	if payload.Currency != "" {
		result["currency"] = payload.Currency
	}
	return result
}

func firstString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		value := values[key]
		switch typed := value.(type) {
		case string:
			if strings.TrimSpace(typed) != "" {
				return strings.TrimSpace(typed)
			}
		case float64:
			if typed != 0 {
				return strconv.FormatInt(int64(typed), 10)
			}
		case int:
			if typed != 0 {
				return strconv.Itoa(typed)
			}
		case int64:
			if typed != 0 {
				return strconv.FormatInt(typed, 10)
			}
		case json.Number:
			if strings.TrimSpace(typed.String()) != "" {
				return typed.String()
			}
		}
	}
	return ""
}

func firstAny(values map[string]any, keys ...string) any {
	for _, key := range keys {
		if value, ok := values[key]; ok {
			return value
		}
	}
	return nil
}

func amountAsFloat(value any) float64 {
	switch typed := value.(type) {
	case float64:
		return typed
	case float32:
		return float64(typed)
	case int:
		return float64(typed)
	case int64:
		return float64(typed)
	case json.Number:
		parsed, _ := typed.Float64()
		return parsed
	case string:
		parsed, _ := strconv.ParseFloat(strings.ReplaceAll(strings.TrimSpace(typed), ",", "."), 64)
		return parsed
	case map[string]any:
		return amountAsFloat(firstAny(typed, "value", "amount", "sum"))
	default:
		return math.NaN()
	}
}
