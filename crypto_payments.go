package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
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
	cryptoProviderDirect = "direct"
	cryptoCurrencyTON    = "TON"
	cryptoCurrencyUSDT   = "USDT"
	cryptoNetworkTON     = "TON"
	cryptoNetworkTRC20   = "TRC20"

	cryptoMethodTON       = "ton"
	cryptoMethodUSDTTON   = "usdt_ton"
	cryptoMethodUSDTTRC20 = "usdt_trc20"

	cryptoStatusPending = "pending"
	cryptoStatusPaid    = "paid"
	cryptoStatusExpired = "expired"

	tonNanoPerTON    = int64(1_000_000_000)
	usdtUnitsPerUSDT = int64(1_000_000)
)

type cryptoPayment struct {
	ID         string
	Provider   string
	Currency   string
	Network    string
	Product    string
	Status     string
	Address    string
	Memo       string
	Amount     string
	AmountNano int64
	PriceRub   int
	TelegramID int64
	TxHash     string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type cryptoPaymentStore interface {
	createCryptoPayment(payment cryptoPayment) error
	getCryptoPayment(id string) (cryptoPayment, bool, error)
	pendingCryptoPayments(network string, currency string, address string, now time.Time, limit int) ([]cryptoPayment, error)
	updateCryptoPayment(payment cryptoPayment) error
}

type cryptoPaymentConfirmation struct {
	Paid   bool
	TxHash string
}

type cryptoPaymentMethod struct {
	ID           string
	Label        string
	Currency     string
	Network      string
	Address      string
	Asset        string
	AmountUnits  int64
	Decimals     int
	RequiresMemo bool
	UniqueAmount bool
}

func (cfg config) cryptoPaymentTTL() time.Duration {
	if cfg.CryptoPaymentTTLMinutes <= 0 {
		return time.Hour
	}
	return time.Duration(cfg.CryptoPaymentTTLMinutes) * time.Minute
}

func (cfg config) cryptoTONEnabled() bool {
	if strings.TrimSpace(cfg.CryptoTONWallet) == "" {
		return false
	}
	for _, product := range paidProductIDs() {
		amount, err := cfg.cryptoTONAmountForProduct(product)
		if err == nil && amount > 0 {
			return true
		}
	}
	return false
}

func (cfg config) cryptoTONPlanEnabled(product string) bool {
	if strings.TrimSpace(cfg.CryptoTONWallet) == "" {
		return false
	}
	amount, err := cfg.cryptoTONAmountForProduct(product)
	return err == nil && amount > 0
}

func (cfg config) cryptoTONAmountForProduct(product string) (int64, error) {
	switch product {
	case premiumMonthlyProduct:
		return parseTONAmountToNano(cfg.CryptoTONMonthAmount)
	case premiumYearlyProduct:
		if strings.TrimSpace(cfg.CryptoTONYearAmount) != "" {
			return parseTONAmountToNano(cfg.CryptoTONYearAmount)
		}
		monthlyAmount, err := parseTONAmountToNano(cfg.CryptoTONMonthAmount)
		if err != nil {
			return 0, err
		}
		return deriveCryptoAmountByPriceRatio(monthlyAmount, cfg.PremiumRubPrice, cfg.PremiumYearRubPrice, "TON")
	case platinumMonthlyProduct:
		return parseTONAmountToNano(cfg.CryptoTONPlatinumMonthAmount)
	case platinumYearlyProduct:
		if strings.TrimSpace(cfg.CryptoTONPlatinumYearAmount) != "" {
			return parseTONAmountToNano(cfg.CryptoTONPlatinumYearAmount)
		}
		monthlyAmount, err := parseTONAmountToNano(cfg.CryptoTONPlatinumMonthAmount)
		if err != nil {
			return 0, err
		}
		return deriveCryptoAmountByPriceRatio(monthlyAmount, cfg.PlatinumRubPrice, cfg.PlatinumYearRubPrice, "TON")
	default:
		return 0, errors.New("unknown product")
	}
}

func deriveCryptoAmountByPriceRatio(baseAmount int64, basePrice int, targetPrice int, label string) (int64, error) {
	if baseAmount <= 0 || basePrice <= 0 || targetPrice <= 0 {
		return 0, fmt.Errorf("%s amount is not configured for this plan", label)
	}
	if baseAmount > math.MaxInt64/int64(targetPrice) {
		return 0, fmt.Errorf("%s amount is too large", label)
	}
	amount := baseAmount * int64(targetPrice) / int64(basePrice)
	if amount <= 0 {
		return 0, fmt.Errorf("%s amount is not configured for this plan", label)
	}
	return amount, nil
}

func (cfg config) cryptoUSDTAmountForProduct(product string) (int64, error) {
	return cfg.cryptoUSDTAmountForProductAtRate(product, float64(cfg.CryptoUSDTRubRate))
}

func (cfg config) cryptoUSDTAmountForProductAtRate(product string, usdtRubRate float64) (int64, error) {
	var configured string
	var priceRub int
	switch product {
	case premiumMonthlyProduct:
		configured = cfg.CryptoUSDTMonthAmount
		priceRub = cfg.PremiumRubPrice
	case premiumYearlyProduct:
		configured = cfg.CryptoUSDTYearAmount
		priceRub = cfg.PremiumYearRubPrice
	case platinumMonthlyProduct:
		configured = cfg.CryptoUSDTPlatinumMonthAmount
		priceRub = cfg.PlatinumRubPrice
	case platinumYearlyProduct:
		configured = cfg.CryptoUSDTPlatinumYearAmount
		priceRub = cfg.PlatinumYearRubPrice
	default:
		return 0, errors.New("unknown product")
	}
	if strings.TrimSpace(configured) != "" {
		return parseDecimalAmountToUnits(configured, 6, "USDT")
	}
	if usdtRubRate <= 0 {
		usdtRubRate = float64(defaultUSDTExchangeRateRubles)
	}
	if priceRub <= 0 {
		return 0, errors.New("USDT amount is not configured for this plan")
	}
	amount := int64(float64(priceRub) * float64(usdtUnitsPerUSDT) / usdtRubRate)
	if amount <= 0 {
		return 0, errors.New("USDT amount is not configured for this plan")
	}
	return amount, nil
}

func (cfg config) cryptoUSDTTONWallet() string {
	if wallet := strings.TrimSpace(cfg.CryptoUSDTTONWallet); wallet != "" {
		return wallet
	}
	return strings.TrimSpace(cfg.CryptoTONWallet)
}

func (cfg config) cryptoPaymentMethodsForProduct(product string) []cryptoPaymentMethod {
	return cfg.cryptoPaymentMethodsForProductAtRate(product, float64(cfg.CryptoUSDTRubRate))
}

func (cfg config) cryptoPaymentMethodsForProductAtRate(product string, usdtRubRate float64) []cryptoPaymentMethod {
	methods := []cryptoPaymentMethod{}
	if amount, err := cfg.cryptoTONAmountForProduct(product); err == nil && amount > 0 {
		if wallet := strings.TrimSpace(cfg.CryptoTONWallet); wallet != "" {
			methods = append(methods, cryptoPaymentMethod{
				ID:           cryptoMethodTON,
				Label:        "TON",
				Currency:     cryptoCurrencyTON,
				Network:      cryptoNetworkTON,
				Address:      wallet,
				AmountUnits:  amount,
				Decimals:     9,
				RequiresMemo: true,
			})
		}
	}
	if amount, err := cfg.cryptoUSDTAmountForProductAtRate(product, usdtRubRate); err == nil && amount > 0 {
		if wallet := cfg.cryptoUSDTTONWallet(); wallet != "" && strings.TrimSpace(cfg.CryptoUSDTTONJettonMaster) != "" {
			methods = append(methods, cryptoPaymentMethod{
				ID:           cryptoMethodUSDTTON,
				Label:        "USDT TON",
				Currency:     cryptoCurrencyUSDT,
				Network:      cryptoNetworkTON,
				Address:      wallet,
				Asset:        strings.TrimSpace(cfg.CryptoUSDTTONJettonMaster),
				AmountUnits:  amount,
				Decimals:     6,
				RequiresMemo: true,
			})
		}
		if wallet := strings.TrimSpace(cfg.CryptoUSDTTRC20Wallet); wallet != "" && strings.TrimSpace(cfg.CryptoUSDTTRC20Contract) != "" {
			methods = append(methods, cryptoPaymentMethod{
				ID:           cryptoMethodUSDTTRC20,
				Label:        "USDT TRC20",
				Currency:     cryptoCurrencyUSDT,
				Network:      cryptoNetworkTRC20,
				Address:      wallet,
				Asset:        strings.TrimSpace(cfg.CryptoUSDTTRC20Contract),
				AmountUnits:  amount,
				Decimals:     6,
				UniqueAmount: true,
			})
		}
	}
	return methods
}

func (cfg config) cryptoPaymentMethodForProduct(product string, methodID string) (cryptoPaymentMethod, bool) {
	methodID = strings.TrimSpace(methodID)
	if methodID == "" {
		methodID = cryptoMethodTON
	}
	for _, method := range cfg.cryptoPaymentMethodsForProduct(product) {
		if method.ID == methodID {
			return method, true
		}
	}
	return cryptoPaymentMethod{}, false
}

func (cfg config) cryptoPlanEnabled(product string) bool {
	return len(cfg.cryptoPaymentMethodsForProduct(product)) > 0
}

func (cfg config) cryptoPaymentsReady(store store) bool {
	enabled := false
	for _, product := range paidProductIDs() {
		if cfg.cryptoPlanEnabled(product) {
			enabled = true
			break
		}
	}
	if !enabled {
		return false
	}
	_, ok := store.(cryptoPaymentStore)
	return ok
}

func (b *bot) cryptoPaymentsReady() bool {
	return b != nil && b.cfg.cryptoPaymentsReady(b.store)
}

func (b *bot) currentUSDTRubRate(ctx context.Context) cryptoRateSnapshot {
	if b == nil {
		return fallbackCryptoRateSnapshot(config{})
	}
	if b.cryptoRates == nil {
		return fallbackCryptoRateSnapshot(b.cfg)
	}
	return b.cryptoRates.usdtRubRate(ctx)
}

func (b *bot) cryptoUSDTAmountForProduct(ctx context.Context, product string) (int64, error) {
	rate := b.currentUSDTRubRate(ctx)
	return b.cfg.cryptoUSDTAmountForProductAtRate(product, rate.USDTRub)
}

func (b *bot) cryptoPaymentMethodsForProduct(ctx context.Context, product string) []cryptoPaymentMethod {
	rate := b.currentUSDTRubRate(ctx)
	return b.cfg.cryptoPaymentMethodsForProductAtRate(product, rate.USDTRub)
}

func (b *bot) cryptoPaymentMethodForProduct(ctx context.Context, product string, methodID string) (cryptoPaymentMethod, bool) {
	methodID = strings.TrimSpace(methodID)
	if methodID == "" {
		methodID = cryptoMethodTON
	}
	for _, method := range b.cryptoPaymentMethodsForProduct(ctx, product) {
		if method.ID == methodID {
			return method, true
		}
	}
	return cryptoPaymentMethod{}, false
}

func (b *bot) createDirectTONPayment(ctx context.Context, product string, user userState) (cryptoPayment, string, error) {
	return b.createDirectCryptoPayment(ctx, product, user, cryptoMethodTON)
}

func (b *bot) createDirectCryptoPayment(ctx context.Context, product string, user userState, methodID string) (cryptoPayment, string, error) {
	return b.createDirectCryptoPaymentWithID(ctx, product, user, methodID, "")
}

func (b *bot) createDirectCryptoPaymentWithID(ctx context.Context, product string, user userState, methodID string, paymentID string) (cryptoPayment, string, error) {
	plan, ok := b.premiumPlan(product)
	if !ok {
		return cryptoPayment{}, "", errors.New("unknown product")
	}
	method, ok := b.cryptoPaymentMethodForProduct(ctx, product, methodID)
	if !ok {
		return cryptoPayment{}, "", errors.New("selected crypto payment method is not configured")
	}
	id := strings.TrimSpace(paymentID)
	if id == "" {
		var err error
		id, err = newCryptoPaymentID(method.ID)
		if err != nil {
			return cryptoPayment{}, "", err
		}
	}
	amountUnits := method.AmountUnits
	if method.UniqueAmount {
		amountUnits = addUniqueUSDTAmountOffset(amountUnits, id)
	}
	memo := ""
	if method.RequiresMemo {
		memo = "POLIGLOT:" + id
	}
	now := time.Now().UTC()
	payment := cryptoPayment{
		ID:         id,
		Provider:   cryptoProviderDirect,
		Currency:   method.Currency,
		Network:    method.Network,
		Product:    product,
		Status:     cryptoStatusPending,
		Address:    method.Address,
		Memo:       memo,
		Amount:     formatUnitsAmount(amountUnits, method.Decimals),
		AmountNano: amountUnits,
		PriceRub:   plan.RubPrice,
		TelegramID: user.TelegramID,
		ExpiresAt:  now.Add(b.cfg.cryptoPaymentTTL()).UTC(),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	return payment, b.cryptoPaymentURL(payment), nil
}

func (b *bot) createOrReuseDirectCryptoPayment(ctx context.Context, product string, user userState, methodID string, idempotencyKey string) (cryptoPayment, string, error) {
	idempotencyKey = normalizeIdempotencyKey(idempotencyKey)
	if idempotencyKey == "" {
		payment, invoiceURL, err := b.createDirectCryptoPayment(ctx, product, user, methodID)
		if err != nil {
			return cryptoPayment{}, "", err
		}
		if err := b.saveDirectCryptoPayment(payment); err != nil {
			return cryptoPayment{}, "", err
		}
		return payment, invoiceURL, nil
	}
	method, ok := b.cryptoPaymentMethodForProduct(ctx, product, methodID)
	if !ok {
		return cryptoPayment{}, "", errors.New("selected crypto payment method is not configured")
	}
	stableID := cryptoPaymentIDFromIdempotencyKey(method.ID, user.TelegramID, product, idempotencyKey)
	store, ok := b.store.(cryptoPaymentStore)
	if !ok {
		return cryptoPayment{}, "", errors.New("direct crypto payments require sqlite storage")
	}
	if existing, ok, err := store.getCryptoPayment(stableID); err != nil {
		return cryptoPayment{}, "", err
	} else if ok {
		if existing.TelegramID != user.TelegramID {
			return cryptoPayment{}, "", errors.New("payment belongs to another user")
		}
		return existing, b.cryptoPaymentURL(existing), nil
	}
	payment, invoiceURL, err := b.createDirectCryptoPaymentWithID(ctx, product, user, method.ID, stableID)
	if err != nil {
		return cryptoPayment{}, "", err
	}
	if err := store.createCryptoPayment(payment); err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "unique") || strings.Contains(lower, "constraint") {
			if existing, ok, lookupErr := store.getCryptoPayment(stableID); lookupErr != nil {
				return cryptoPayment{}, "", lookupErr
			} else if ok && existing.TelegramID == user.TelegramID {
				return existing, b.cryptoPaymentURL(existing), nil
			}
		}
		return cryptoPayment{}, "", err
	}
	return payment, invoiceURL, nil
}

func (b *bot) saveDirectTONPayment(payment cryptoPayment) error {
	return b.saveDirectCryptoPayment(payment)
}

func (b *bot) saveDirectCryptoPayment(payment cryptoPayment) error {
	store, ok := b.store.(cryptoPaymentStore)
	if !ok {
		return errors.New("direct crypto payments require sqlite storage")
	}
	return store.createCryptoPayment(payment)
}

func (b *bot) refreshDirectTONPayment(ctx context.Context, paymentID string, user userState) (cryptoPayment, userState, bool, error) {
	return b.refreshDirectCryptoPayment(ctx, paymentID, user)
}

func (b *bot) refreshDirectCryptoPayment(ctx context.Context, paymentID string, user userState) (cryptoPayment, userState, bool, error) {
	store, ok := b.store.(cryptoPaymentStore)
	if !ok {
		return cryptoPayment{}, userState{}, false, errors.New("direct crypto payments require sqlite storage")
	}
	payment, ok, err := store.getCryptoPayment(paymentID)
	if err != nil {
		return cryptoPayment{}, userState{}, false, err
	}
	if !ok {
		return cryptoPayment{}, userState{}, false, errors.New("payment not found")
	}
	if user.TelegramID != 0 && payment.TelegramID != user.TelegramID {
		return cryptoPayment{}, userState{}, false, errors.New("payment belongs to another user")
	}
	if payment.Status == cryptoStatusPaid {
		refreshed, err := b.store.getOrCreateUser(payment.TelegramID, user.FirstName)
		return payment, refreshed, true, err
	}
	now := time.Now().UTC()
	expired := now.After(payment.ExpiresAt)
	confirmation, err := b.confirmDirectCryptoPayment(ctx, payment)
	if err != nil {
		return payment, userState{}, false, err
	}
	if !confirmation.Paid {
		if expired {
			payment.Status = cryptoStatusExpired
			payment.UpdatedAt = now
			if err := store.updateCryptoPayment(payment); err != nil {
				return cryptoPayment{}, userState{}, false, err
			}
		}
		return payment, userState{}, false, nil
	}
	plan, ok := b.premiumPlan(payment.Product)
	if !ok {
		return payment, userState{}, false, errors.New("payment product is no longer available")
	}
	payment.Status = cryptoStatusPaid
	payment.TxHash = confirmation.TxHash
	payment.UpdatedAt = now
	if err := store.updateCryptoPayment(payment); err != nil {
		return cryptoPayment{}, userState{}, false, err
	}
	until, err := b.store.extendPremium(payment.TelegramID, "crypto:"+payment.ID, plan.Duration, plan.Tier)
	if err != nil {
		return cryptoPayment{}, userState{}, false, err
	}
	if _, err := b.store.creditReferralPurchase(payment.TelegramID, "crypto:"+payment.ID, rubToKopecks(payment.PriceRub)); err != nil {
		return cryptoPayment{}, userState{}, false, err
	}
	refreshed, err := b.store.getOrCreateUser(payment.TelegramID, user.FirstName)
	if err != nil {
		return cryptoPayment{}, userState{}, false, err
	}
	refreshed.PremiumUntil = until
	b.notifyOpsPayment(ctx, "Direct crypto", refreshed, plan, payment.ID, strings.TrimSpace(payment.Amount+" "+payment.Currency), until)
	return payment, refreshed, true, nil
}

func (b *bot) confirmDirectCryptoPayment(ctx context.Context, payment cryptoPayment) (cryptoPaymentConfirmation, error) {
	switch {
	case payment.Network == cryptoNetworkTON && payment.Currency == cryptoCurrencyTON:
		return b.confirmDirectTONPayment(ctx, payment)
	case payment.Network == cryptoNetworkTON && payment.Currency == cryptoCurrencyUSDT:
		return b.confirmDirectUSDTTONPayment(ctx, payment)
	case payment.Network == cryptoNetworkTRC20 && payment.Currency == cryptoCurrencyUSDT:
		return b.confirmDirectUSDTTRC20Payment(ctx, payment)
	default:
		return cryptoPaymentConfirmation{}, errors.New("unsupported crypto payment network")
	}
}

func (b *bot) confirmDirectTONPayment(ctx context.Context, payment cryptoPayment) (cryptoPaymentConfirmation, error) {
	if payment.Network != cryptoNetworkTON || payment.Currency != cryptoCurrencyTON {
		return cryptoPaymentConfirmation{}, errors.New("unsupported crypto payment network")
	}
	var providerErrors []string
	providerSuccess := false
	tryTONAPI := func() (cryptoPaymentConfirmation, bool) {
		if strings.TrimSpace(b.cfg.CryptoTONAPIBaseURL) == "" {
			return cryptoPaymentConfirmation{}, false
		}
		confirmation, err := b.confirmDirectTONPaymentWithTONAPI(ctx, payment)
		if err == nil {
			providerSuccess = true
			return confirmation, confirmation.Paid
		}
		providerErrors = append(providerErrors, err.Error())
		return cryptoPaymentConfirmation{}, false
	}
	tryTONCenter := func() (cryptoPaymentConfirmation, bool) {
		if strings.TrimSpace(b.cfg.CryptoTONCenterBaseURL) == "" {
			return cryptoPaymentConfirmation{}, false
		}
		confirmation, err := b.confirmDirectTONPaymentWithTONCenter(ctx, payment, true)
		if err == nil {
			providerSuccess = true
			return confirmation, confirmation.Paid
		}
		providerErrors = append(providerErrors, err.Error())
		if isTONAuthError(err) && strings.TrimSpace(b.cfg.CryptoTONCenterAPIKey) != "" {
			confirmation, retryErr := b.confirmDirectTONPaymentWithTONCenter(ctx, payment, false)
			if retryErr == nil {
				providerSuccess = true
				return confirmation, confirmation.Paid
			}
			providerErrors = append(providerErrors, retryErr.Error())
		}
		return cryptoPaymentConfirmation{}, false
	}
	if strings.TrimSpace(b.cfg.CryptoTONAPIKey) != "" {
		if confirmation, paid := tryTONAPI(); paid {
			return confirmation, nil
		}
		if confirmation, paid := tryTONCenter(); paid {
			return confirmation, nil
		}
	} else {
		if confirmation, paid := tryTONCenter(); paid {
			return confirmation, nil
		}
		if confirmation, paid := tryTONAPI(); paid {
			return confirmation, nil
		}
	}
	if providerSuccess {
		return cryptoPaymentConfirmation{}, nil
	}
	if len(providerErrors) == 0 {
		return cryptoPaymentConfirmation{}, errors.New("TON transaction providers are not configured")
	}
	return cryptoPaymentConfirmation{}, fmt.Errorf("could not check TON payment: %s", strings.Join(providerErrors, "; "))
}

func (b *bot) confirmDirectTONPaymentWithTONCenter(ctx context.Context, payment cryptoPayment, useAPIKey bool) (cryptoPaymentConfirmation, error) {
	endpoint, err := url.Parse(strings.TrimRight(strings.TrimSpace(b.cfg.CryptoTONCenterBaseURL), "/") + "/getTransactions")
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	values := endpoint.Query()
	values.Set("address", payment.Address)
	values.Set("limit", "50")
	values.Set("to_lt", "0")
	values.Set("archival", "true")
	endpoint.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	request.Header.Set("Accept", "application/json")
	if apiKey := strings.TrimSpace(b.cfg.CryptoTONCenterAPIKey); useAPIKey && apiKey != "" {
		request.Header.Set("X-API-Key", apiKey)
	}
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 8<<20))
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	if response.StatusCode >= 400 {
		return cryptoPaymentConfirmation{}, tonProviderStatusError{Provider: "TON Center", StatusCode: response.StatusCode, Body: string(body)}
	}
	var decoded struct {
		OK     bool              `json:"ok"`
		Result []json.RawMessage `json:"result"`
		Error  string            `json:"error"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	if !decoded.OK {
		if decoded.Error == "" {
			decoded.Error = "unknown TON Center error"
		}
		return cryptoPaymentConfirmation{}, errors.New(decoded.Error)
	}
	for _, raw := range decoded.Result {
		tx, err := parseTONTransaction(raw)
		if err != nil {
			continue
		}
		if !tx.UnixTime.IsZero() && tx.UnixTime.Before(payment.CreatedAt.Add(-10*time.Minute)) {
			continue
		}
		if !tx.UnixTime.IsZero() && !payment.ExpiresAt.IsZero() && tx.UnixTime.After(payment.ExpiresAt.Add(10*time.Minute)) {
			continue
		}
		if tx.ValueNano < payment.AmountNano {
			continue
		}
		if strings.TrimSpace(tx.Comment) != payment.Memo {
			continue
		}
		return cryptoPaymentConfirmation{Paid: true, TxHash: tx.Hash}, nil
	}
	return cryptoPaymentConfirmation{}, nil
}

func (b *bot) confirmDirectTONPaymentWithTONAPI(ctx context.Context, payment cryptoPayment) (cryptoPaymentConfirmation, error) {
	endpoint, err := url.Parse(strings.TrimRight(strings.TrimSpace(b.cfg.CryptoTONAPIBaseURL), "/") + "/blockchain/accounts/" + url.PathEscape(payment.Address) + "/transactions")
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	values := endpoint.Query()
	values.Set("limit", "50")
	endpoint.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	request.Header.Set("Accept", "application/json")
	if apiKey := strings.TrimSpace(b.cfg.CryptoTONAPIKey); apiKey != "" {
		request.Header.Set("Authorization", "Bearer "+apiKey)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 8<<20))
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	if response.StatusCode >= 400 {
		return cryptoPaymentConfirmation{}, tonProviderStatusError{Provider: "TonAPI", StatusCode: response.StatusCode, Body: string(body)}
	}
	var decoded struct {
		Transactions []json.RawMessage `json:"transactions"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	for _, raw := range decoded.Transactions {
		tx, err := parseTONTransaction(raw)
		if err != nil {
			continue
		}
		if !tx.UnixTime.IsZero() && tx.UnixTime.Before(payment.CreatedAt.Add(-10*time.Minute)) {
			continue
		}
		if !tx.UnixTime.IsZero() && !payment.ExpiresAt.IsZero() && tx.UnixTime.After(payment.ExpiresAt.Add(10*time.Minute)) {
			continue
		}
		if tx.ValueNano < payment.AmountNano {
			continue
		}
		if strings.TrimSpace(tx.Comment) != payment.Memo {
			continue
		}
		return cryptoPaymentConfirmation{Paid: true, TxHash: tx.Hash}, nil
	}
	return cryptoPaymentConfirmation{}, nil
}

func (b *bot) confirmDirectUSDTTONPayment(ctx context.Context, payment cryptoPayment) (cryptoPaymentConfirmation, error) {
	if payment.Network != cryptoNetworkTON || payment.Currency != cryptoCurrencyUSDT {
		return cryptoPaymentConfirmation{}, errors.New("unsupported crypto payment network")
	}
	var providerErrors []string
	providerSuccess := false
	tryTONCenterV3 := func(useAPIKey bool) (cryptoPaymentConfirmation, bool) {
		confirmation, err := b.confirmDirectUSDTTONPaymentWithTONCenterV3(ctx, payment, useAPIKey)
		if err == nil {
			providerSuccess = true
			return confirmation, confirmation.Paid
		}
		providerErrors = append(providerErrors, err.Error())
		return cryptoPaymentConfirmation{}, false
	}
	if confirmation, paid := tryTONCenterV3(true); paid {
		return confirmation, nil
	}
	if len(providerErrors) > 0 && strings.TrimSpace(b.cfg.CryptoTONCenterAPIKey) != "" {
		if lastErr := providerErrors[len(providerErrors)-1]; strings.Contains(lastErr, "status 401") || strings.Contains(lastErr, "status 403") {
			if confirmation, paid := tryTONCenterV3(false); paid {
				return confirmation, nil
			}
		}
	}
	if providerSuccess {
		return cryptoPaymentConfirmation{}, nil
	}
	if len(providerErrors) == 0 {
		return cryptoPaymentConfirmation{}, errors.New("TON jetton transaction provider is not configured")
	}
	return cryptoPaymentConfirmation{}, fmt.Errorf("could not check USDT TON payment: %s", strings.Join(providerErrors, "; "))
}

func (b *bot) confirmDirectUSDTTONPaymentWithTONCenterV3(ctx context.Context, payment cryptoPayment, useAPIKey bool) (cryptoPaymentConfirmation, error) {
	baseURL := tonCenterV3BaseURL(b.cfg.CryptoTONCenterBaseURL)
	if baseURL == "" {
		return cryptoPaymentConfirmation{}, errors.New("TON Center v3 is not configured")
	}
	endpoint, err := url.Parse(strings.TrimRight(baseURL, "/") + "/jetton/transfers")
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	start := payment.CreatedAt.Add(-10 * time.Minute).Unix()
	end := payment.ExpiresAt.Add(10 * time.Minute).Unix()
	values := endpoint.Query()
	values.Set("owner_address", payment.Address)
	values.Set("direction", "in")
	values.Set("jetton_master", strings.TrimSpace(b.cfg.CryptoUSDTTONJettonMaster))
	values.Set("start_utime", strconv.FormatInt(start, 10))
	values.Set("end_utime", strconv.FormatInt(end, 10))
	values.Set("limit", "100")
	values.Set("sort", "desc")
	endpoint.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	request.Header.Set("Accept", "application/json")
	if apiKey := strings.TrimSpace(b.cfg.CryptoTONCenterAPIKey); useAPIKey && apiKey != "" {
		request.Header.Set("X-API-Key", apiKey)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 8<<20))
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	if response.StatusCode >= 400 {
		return cryptoPaymentConfirmation{}, tonProviderStatusError{Provider: "TON Center v3", StatusCode: response.StatusCode, Body: string(body)}
	}
	var decoded struct {
		JettonTransfers []json.RawMessage `json:"jetton_transfers"`
		Error           string            `json:"error"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	for _, raw := range decoded.JettonTransfers {
		tx, err := parseTONJettonTransfer(raw)
		if err != nil {
			continue
		}
		if !tx.UnixTime.IsZero() && tx.UnixTime.Before(payment.CreatedAt.Add(-10*time.Minute)) {
			continue
		}
		if !tx.UnixTime.IsZero() && !payment.ExpiresAt.IsZero() && tx.UnixTime.After(payment.ExpiresAt.Add(10*time.Minute)) {
			continue
		}
		if tx.ValueNano < payment.AmountNano {
			continue
		}
		if strings.TrimSpace(payment.Memo) != "" && strings.TrimSpace(tx.Comment) != payment.Memo {
			continue
		}
		return cryptoPaymentConfirmation{Paid: true, TxHash: "TON:" + tx.Hash}, nil
	}
	return cryptoPaymentConfirmation{}, nil
}

func (b *bot) confirmDirectUSDTTRC20Payment(ctx context.Context, payment cryptoPayment) (cryptoPaymentConfirmation, error) {
	if payment.Network != cryptoNetworkTRC20 || payment.Currency != cryptoCurrencyUSDT {
		return cryptoPaymentConfirmation{}, errors.New("unsupported crypto payment network")
	}
	endpoint, err := url.Parse(strings.TrimRight(strings.TrimSpace(b.cfg.CryptoTronGridBaseURL), "/") + "/v1/accounts/" + url.PathEscape(payment.Address) + "/transactions/trc20")
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	values := endpoint.Query()
	values.Set("limit", "200")
	values.Set("only_confirmed", "true")
	values.Set("contract_address", strings.TrimSpace(b.cfg.CryptoUSDTTRC20Contract))
	endpoint.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	request.Header.Set("Accept", "application/json")
	if apiKey := strings.TrimSpace(b.cfg.CryptoTronGridAPIKey); apiKey != "" {
		request.Header.Set("TRON-PRO-API-KEY", apiKey)
		request.Header.Set("TRON_PRO_API_KEY", apiKey)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 8<<20))
	if err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	if response.StatusCode >= 400 {
		return cryptoPaymentConfirmation{}, tonProviderStatusError{Provider: "TronGrid", StatusCode: response.StatusCode, Body: string(body)}
	}
	var decoded struct {
		Data    []json.RawMessage `json:"data"`
		Success bool              `json:"success"`
		Error   string            `json:"error"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		return cryptoPaymentConfirmation{}, err
	}
	for _, raw := range decoded.Data {
		tx, err := parseTRC20Transfer(raw)
		if err != nil {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(tx.To), strings.TrimSpace(payment.Address)) {
			continue
		}
		if tx.Contract != "" && !strings.EqualFold(strings.TrimSpace(tx.Contract), strings.TrimSpace(b.cfg.CryptoUSDTTRC20Contract)) {
			continue
		}
		if !tx.UnixTime.IsZero() && tx.UnixTime.Before(payment.CreatedAt.Add(-10*time.Minute)) {
			continue
		}
		if !tx.UnixTime.IsZero() && !payment.ExpiresAt.IsZero() && tx.UnixTime.After(payment.ExpiresAt.Add(10*time.Minute)) {
			continue
		}
		if tx.ValueNano != payment.AmountNano {
			continue
		}
		return cryptoPaymentConfirmation{Paid: true, TxHash: "TRON:" + tx.Hash}, nil
	}
	return cryptoPaymentConfirmation{}, nil
}

func (b *bot) handleTONAPIWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if strings.TrimSpace(b.cfg.CryptoTONAPIWebhookKey) == "" {
		http.Error(w, "TON API webhook is not configured", http.StatusServiceUnavailable)
		return
	}
	if !b.validTONAPIWebhookToken(r) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var payload struct {
		AccountID string `json:"account_id"`
		LT        any    `json:"lt"`
		TxHash    string `json:"tx_hash"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(&payload); err != nil && !errors.Is(err, io.EOF) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	store, ok := b.store.(cryptoPaymentStore)
	if !ok {
		http.Error(w, "direct crypto payments require sqlite storage", http.StatusServiceUnavailable)
		return
	}
	now := time.Now().UTC()
	payments := []cryptoPayment{}
	if wallet := strings.TrimSpace(b.cfg.CryptoTONWallet); wallet != "" {
		tonPayments, err := store.pendingCryptoPayments(cryptoNetworkTON, cryptoCurrencyTON, wallet, now, 100)
		if err != nil {
			http.Error(w, "payment lookup failed", http.StatusInternalServerError)
			return
		}
		payments = append(payments, tonPayments...)
	}
	if wallet := b.cfg.cryptoUSDTTONWallet(); wallet != "" {
		usdtPayments, err := store.pendingCryptoPayments(cryptoNetworkTON, cryptoCurrencyUSDT, wallet, now, 100)
		if err != nil {
			http.Error(w, "payment lookup failed", http.StatusInternalServerError)
			return
		}
		payments = append(payments, usdtPayments...)
	}
	activated := 0
	checked := 0
	for _, pending := range payments {
		checked++
		payment, user, paid, err := b.refreshDirectCryptoPayment(r.Context(), pending.ID, userState{TelegramID: pending.TelegramID})
		if err != nil {
			logCryptoPaymentError("tonapi webhook payment check failed", err)
			continue
		}
		if !paid {
			_ = payment
			continue
		}
		activated++
		if user.TelegramID > 0 && b.telegram != nil {
			plan, _ := b.premiumPlan(payment.Product)
			plan = localizedPremiumPlan(user, plan)
			_ = b.telegram.sendMarkdownMessageWithCopy(r.Context(), user.TelegramID, premiumActivationMarkdown(user, plan, user.PremiumUntil.Format("2006-01-02"), true), ui(user))
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"checked":   checked,
		"activated": activated,
		"tx_hash":   strings.TrimSpace(payload.TxHash),
	})
}

func (b *bot) validTONAPIWebhookToken(r *http.Request) bool {
	expected := strings.TrimSpace(b.cfg.CryptoTONAPIWebhookKey)
	if expected == "" {
		return false
	}
	if token := strings.TrimSpace(r.URL.Query().Get("token")); token != "" {
		return subtleConstantTimeString(token, expected)
	}
	prefix := "/tonapi/webhook/"
	if strings.HasPrefix(r.URL.Path, prefix) {
		token := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, prefix))
		return subtleConstantTimeString(token, expected)
	}
	return false
}

func subtleConstantTimeString(a string, b string) bool {
	if len(a) != len(b) {
		return false
	}
	result := byte(0)
	for i := range a {
		result |= a[i] ^ b[i]
	}
	return result == 0
}

type tonProviderStatusError struct {
	Provider   string
	StatusCode int
	Body       string
}

func (err tonProviderStatusError) Error() string {
	body := strings.TrimSpace(err.Body)
	if len(body) > 160 {
		body = body[:160]
	}
	if body == "" {
		return fmt.Sprintf("%s status %d", err.Provider, err.StatusCode)
	}
	return fmt.Sprintf("%s status %d: %s", err.Provider, err.StatusCode, body)
}

func isTONAuthError(err error) bool {
	var statusErr tonProviderStatusError
	return errors.As(err, &statusErr) && (statusErr.StatusCode == http.StatusUnauthorized || statusErr.StatusCode == http.StatusForbidden)
}

type parsedTONTransaction struct {
	Hash      string
	ValueNano int64
	Comment   string
	UnixTime  time.Time
}

type parsedTokenTransfer struct {
	Hash      string
	ValueNano int64
	Comment   string
	UnixTime  time.Time
	To        string
	Contract  string
}

func parseTONTransaction(raw json.RawMessage) (parsedTONTransaction, error) {
	var tx map[string]any
	if err := json.Unmarshal(raw, &tx); err != nil {
		return parsedTONTransaction{}, err
	}
	parsed := parsedTONTransaction{
		Hash: firstStringByPath(tx, []string{"transaction_id", "hash"}, []string{"hash"}),
		Comment: firstStringByPath(tx,
			[]string{"in_msg", "message"},
			[]string{"in_msg", "msg_data", "text"},
			[]string{"in_msg", "decoded", "comment"},
			[]string{"in_msg", "decoded_body", "comment"},
			[]string{"in_msg", "decoded_body", "text"},
			[]string{"in_msg", "message_content", "decoded", "comment"},
			[]string{"in_msg", "message_content", "decoded", "text"},
		),
	}
	parsed.ValueNano = firstInt64ByPath(tx, []string{"in_msg", "value"}, []string{"in_msg", "amount"})
	if unixTime := firstInt64ByPath(tx, []string{"utime"}, []string{"now"}); unixTime > 0 {
		parsed.UnixTime = time.Unix(unixTime, 0).UTC()
	}
	if parsed.Hash == "" {
		parsed.Hash = firstStringByPath(tx, []string{"transaction_id", "lt"})
	}
	return parsed, nil
}

func parseTONJettonTransfer(raw json.RawMessage) (parsedTokenTransfer, error) {
	var tx map[string]any
	if err := json.Unmarshal(raw, &tx); err != nil {
		return parsedTokenTransfer{}, err
	}
	parsed := parsedTokenTransfer{
		Hash: firstStringByPath(tx,
			[]string{"transaction_hash"},
			[]string{"trace_id"},
			[]string{"transaction_lt"},
		),
		Comment: firstStringByPath(tx,
			[]string{"comment"},
			[]string{"decoded_forward_payload", "comment"},
			[]string{"decoded_custom_payload", "comment"},
		),
		To:       firstStringByPath(tx, []string{"destination"}),
		Contract: firstStringByPath(tx, []string{"jetton_master"}),
	}
	parsed.ValueNano = firstInt64ByPath(tx, []string{"amount"})
	if unixTime := firstInt64ByPath(tx, []string{"transaction_now"}); unixTime > 0 {
		parsed.UnixTime = time.Unix(unixTime, 0).UTC()
	}
	if parsed.Hash == "" {
		parsed.Hash = fmt.Sprintf("%s:%d", parsed.To, parsed.UnixTime.Unix())
	}
	return parsed, nil
}

func parseTRC20Transfer(raw json.RawMessage) (parsedTokenTransfer, error) {
	var tx map[string]any
	if err := json.Unmarshal(raw, &tx); err != nil {
		return parsedTokenTransfer{}, err
	}
	parsed := parsedTokenTransfer{
		Hash: firstStringByPath(tx,
			[]string{"transaction_id"},
			[]string{"transaction_hash"},
			[]string{"hash"},
		),
		To: firstStringByPath(tx,
			[]string{"to"},
			[]string{"transferToAddress"},
		),
		Contract: firstStringByPath(tx,
			[]string{"token_info", "address"},
			[]string{"contract_address"},
		),
	}
	parsed.ValueNano = firstInt64ByPath(tx, []string{"value"}, []string{"quant"})
	if ms := firstInt64ByPath(tx, []string{"block_timestamp"}); ms > 0 {
		if ms > 10_000_000_000 {
			parsed.UnixTime = time.UnixMilli(ms).UTC()
		} else {
			parsed.UnixTime = time.Unix(ms, 0).UTC()
		}
	}
	return parsed, nil
}

func firstStringByPath(source map[string]any, paths ...[]string) string {
	for _, path := range paths {
		value, ok := valueByPath(source, path)
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case string:
			if strings.TrimSpace(typed) != "" {
				return strings.TrimSpace(typed)
			}
		case float64:
			return strconv.FormatInt(int64(typed), 10)
		}
	}
	return ""
}

func firstInt64ByPath(source map[string]any, paths ...[]string) int64 {
	for _, path := range paths {
		value, ok := valueByPath(source, path)
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case float64:
			return int64(typed)
		case string:
			parsed, err := strconv.ParseInt(strings.TrimSpace(typed), 10, 64)
			if err == nil {
				return parsed
			}
		}
	}
	return 0
}

func valueByPath(source map[string]any, path []string) (any, bool) {
	var current any = source
	for _, key := range path {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = object[key]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func newCryptoPaymentID(methodID string) (string, error) {
	var bytes [12]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	prefix := strings.NewReplacer("-", "_", " ", "_").Replace(strings.ToLower(strings.TrimSpace(methodID)))
	if prefix == "" {
		prefix = "crypto"
	}
	return prefix + "_" + hex.EncodeToString(bytes[:]), nil
}

func (b *bot) cryptoPaymentURL(payment cryptoPayment) string {
	switch {
	case payment.Network == cryptoNetworkTON && payment.Currency == cryptoCurrencyTON:
		return tonTransferURL(payment)
	case payment.Network == cryptoNetworkTON && payment.Currency == cryptoCurrencyUSDT:
		return usdtTONTransferURL(payment, strings.TrimSpace(b.cfg.CryptoUSDTTONJettonMaster))
	default:
		return ""
	}
}

func cryptoPaymentMethodID(payment cryptoPayment) string {
	switch {
	case payment.Network == cryptoNetworkTON && payment.Currency == cryptoCurrencyTON:
		return cryptoMethodTON
	case payment.Network == cryptoNetworkTON && payment.Currency == cryptoCurrencyUSDT:
		return cryptoMethodUSDTTON
	case payment.Network == cryptoNetworkTRC20 && payment.Currency == cryptoCurrencyUSDT:
		return cryptoMethodUSDTTRC20
	default:
		return "crypto"
	}
}

func cryptoPaymentMethodLabel(payment cryptoPayment) string {
	switch cryptoPaymentMethodID(payment) {
	case cryptoMethodTON:
		return "TON"
	case cryptoMethodUSDTTON:
		return "USDT TON"
	case cryptoMethodUSDTTRC20:
		return "USDT TRC20"
	default:
		return strings.TrimSpace(payment.Currency + " " + payment.Network)
	}
}

func tonTransferURL(payment cryptoPayment) string {
	values := url.Values{}
	values.Set("amount", strconv.FormatInt(payment.AmountNano, 10))
	values.Set("text", payment.Memo)
	values.Set("exp", strconv.FormatInt(payment.ExpiresAt.Unix(), 10))
	return "https://app.tonkeeper.com/transfer/" + url.PathEscape(payment.Address) + "?" + values.Encode()
}

func usdtTONTransferURL(payment cryptoPayment, jettonMaster string) string {
	values := url.Values{}
	values.Set("jetton", strings.TrimSpace(jettonMaster))
	values.Set("amount", strconv.FormatInt(payment.AmountNano, 10))
	if strings.TrimSpace(payment.Memo) != "" {
		values.Set("text", payment.Memo)
	}
	values.Set("exp", strconv.FormatInt(payment.ExpiresAt.Unix(), 10))
	return "https://app.tonkeeper.com/transfer/" + url.PathEscape(payment.Address) + "?" + values.Encode()
}

func tonCenterV3BaseURL(base string) string {
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if base == "" {
		return "https://toncenter.com/api/v3"
	}
	if strings.HasSuffix(base, "/api/v3") {
		return base
	}
	if strings.HasSuffix(base, "/api/v2") {
		return strings.TrimSuffix(base, "/api/v2") + "/api/v3"
	}
	if strings.Contains(base, "/api/v2/") {
		return strings.Replace(base, "/api/v2/", "/api/v3/", 1)
	}
	if strings.Contains(base, "/api/v2") {
		return strings.Replace(base, "/api/v2", "/api/v3", 1)
	}
	return base + "/api/v3"
}

func parseTONAmountToNano(value string) (int64, error) {
	return parseDecimalAmountToUnits(value, 9, "TON")
}

func parseDecimalAmountToUnits(value string, decimals int, label string) (int64, error) {
	value = strings.TrimSpace(strings.ReplaceAll(value, ",", "."))
	if value == "" {
		return 0, fmt.Errorf("empty %s amount", label)
	}
	if strings.HasPrefix(value, "-") {
		return 0, fmt.Errorf("%s amount must be positive", label)
	}
	parts := strings.Split(value, ".")
	if len(parts) > 2 {
		return 0, fmt.Errorf("invalid %s amount", label)
	}
	wholeText := parts[0]
	if wholeText == "" {
		wholeText = "0"
	}
	whole, err := strconv.ParseInt(wholeText, 10, 64)
	if err != nil {
		return 0, err
	}
	fractionText := ""
	if len(parts) == 2 {
		fractionText = parts[1]
	}
	if len(fractionText) > decimals {
		return 0, fmt.Errorf("%s amount supports up to %d decimal places", label, decimals)
	}
	for len(fractionText) < decimals {
		fractionText += "0"
	}
	fraction := int64(0)
	if fractionText != "" {
		fraction, err = strconv.ParseInt(fractionText, 10, 64)
		if err != nil {
			return 0, err
		}
	}
	scale := int64(1)
	for i := 0; i < decimals; i++ {
		scale *= 10
	}
	if whole > (math.MaxInt64-fraction)/scale {
		return 0, fmt.Errorf("%s amount is too large", label)
	}
	return whole*scale + fraction, nil
}

func formatTONAmount(nano int64) string {
	return formatUnitsAmount(nano, 9)
}

func formatUnitsAmount(units int64, decimals int) string {
	if units <= 0 {
		return "0"
	}
	scale := int64(1)
	for i := 0; i < decimals; i++ {
		scale *= 10
	}
	whole := units / scale
	fraction := units % scale
	if fraction == 0 {
		return strconv.FormatInt(whole, 10)
	}
	fractionText := fmt.Sprintf("%0*d", decimals, fraction)
	fractionText = strings.TrimRight(fractionText, "0")
	return strconv.FormatInt(whole, 10) + "." + fractionText
}

func addUniqueUSDTAmountOffset(units int64, id string) int64 {
	if units <= 0 {
		return units
	}
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(id))
	return units + int64(hash.Sum32()%9000) + 1000
}

func (s *sqliteStore) createCryptoPayment(payment cryptoPayment) error {
	now := time.Now().UTC()
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = now
	}
	if payment.UpdatedAt.IsZero() {
		payment.UpdatedAt = now
	}
	_, err := s.db.Exec(`INSERT INTO crypto_payments (
		id, provider, currency, network, product, status, address, memo, amount, amount_nano,
		price_rub, telegram_id, tx_hash, expires_at, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		payment.ID, payment.Provider, payment.Currency, payment.Network, payment.Product, payment.Status,
		payment.Address, payment.Memo, payment.Amount, payment.AmountNano, payment.PriceRub, payment.TelegramID,
		payment.TxHash, formatDBTime(payment.ExpiresAt), formatDBTime(payment.CreatedAt), formatDBTime(payment.UpdatedAt))
	return err
}

func (s *sqliteStore) getCryptoPayment(id string) (cryptoPayment, bool, error) {
	row := s.db.QueryRow(`SELECT id, provider, currency, network, product, status, address, memo, amount, amount_nano,
		price_rub, telegram_id, tx_hash, expires_at, created_at, updated_at
		FROM crypto_payments WHERE id = ?`, strings.TrimSpace(id))
	payment, err := scanCryptoPayment(row)
	if errors.Is(err, sql.ErrNoRows) {
		return cryptoPayment{}, false, nil
	}
	if err != nil {
		return cryptoPayment{}, false, err
	}
	return payment, true, nil
}

func (s *sqliteStore) pendingCryptoPayments(network string, currency string, address string, now time.Time, limit int) ([]cryptoPayment, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := s.db.Query(`SELECT id, provider, currency, network, product, status, address, memo, amount, amount_nano,
		price_rub, telegram_id, tx_hash, expires_at, created_at, updated_at
		FROM crypto_payments
		WHERE status = ? AND network = ? AND currency = ? AND address = ? AND created_at >= ?
		ORDER BY created_at DESC
		LIMIT ?`,
		cryptoStatusPending, network, currency, strings.TrimSpace(address), formatDBTime(now.UTC().Add(-24*time.Hour)), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	payments := []cryptoPayment{}
	for rows.Next() {
		payment, err := scanCryptoPayment(rows)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, rows.Err()
}

func (s *sqliteStore) updateCryptoPayment(payment cryptoPayment) error {
	if payment.UpdatedAt.IsZero() {
		payment.UpdatedAt = time.Now().UTC()
	}
	_, err := s.db.Exec(`UPDATE crypto_payments SET
		status = ?, tx_hash = ?, updated_at = ?
		WHERE id = ?`,
		payment.Status, payment.TxHash, formatDBTime(payment.UpdatedAt), payment.ID)
	return err
}

type cryptoPaymentScanner interface {
	Scan(dest ...any) error
}

func scanCryptoPayment(scanner cryptoPaymentScanner) (cryptoPayment, error) {
	var payment cryptoPayment
	var expiresAt, createdAt, updatedAt string
	err := scanner.Scan(&payment.ID, &payment.Provider, &payment.Currency, &payment.Network, &payment.Product,
		&payment.Status, &payment.Address, &payment.Memo, &payment.Amount, &payment.AmountNano, &payment.PriceRub,
		&payment.TelegramID, &payment.TxHash, &expiresAt, &createdAt, &updatedAt)
	if err != nil {
		return cryptoPayment{}, err
	}
	payment.ExpiresAt = parseDBTime(expiresAt)
	payment.CreatedAt = parseDBTime(createdAt)
	payment.UpdatedAt = parseDBTime(updatedAt)
	return payment, nil
}

func logCryptoPaymentError(context string, err error) {
	if err != nil {
		log.Printf("%s: %v", context, err)
	}
}
