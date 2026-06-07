package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestTONAmountParsing(t *testing.T) {
	cases := map[string]int64{
		"1":           1_000_000_000,
		"0.5":         500_000_000,
		"2.000000001": 2_000_000_001,
		"0,25":        250_000_000,
	}
	for input, want := range cases {
		got, err := parseTONAmountToNano(input)
		if err != nil {
			t.Fatalf("parseTONAmountToNano(%q): %v", input, err)
		}
		if got != want {
			t.Fatalf("parseTONAmountToNano(%q) = %d, want %d", input, got, want)
		}
	}
	if _, err := parseTONAmountToNano("1.1234567891"); err == nil {
		t.Fatal("expected too many decimal places to fail")
	}
}

func TestTONYearAmountFallsBackToMonthlyPriceRatio(t *testing.T) {
	cfg := config{
		PremiumRubPrice:              300,
		PremiumYearRubPrice:          3000,
		CryptoTONMonthAmount:         "1.5",
		CryptoTONYearAmount:          "",
		PlatinumRubPrice:             590,
		PlatinumYearRubPrice:         5900,
		CryptoTONPlatinumMonthAmount: "3",
	}
	amount, err := cfg.cryptoTONAmountForProduct(premiumYearlyProduct)
	if err != nil {
		t.Fatalf("premium yearly TON fallback: %v", err)
	}
	if amount != 15*tonNanoPerTON {
		t.Fatalf("premium yearly TON amount = %d, want %d", amount, 15*tonNanoPerTON)
	}
	platinumAmount, err := cfg.cryptoTONAmountForProduct(platinumYearlyProduct)
	if err != nil {
		t.Fatalf("platinum yearly TON fallback: %v", err)
	}
	if platinumAmount != 30*tonNanoPerTON {
		t.Fatalf("platinum yearly TON amount = %d, want %d", platinumAmount, 30*tonNanoPerTON)
	}
}

func TestUSDTAmountParsingAndUniqueOffset(t *testing.T) {
	units, err := parseDecimalAmountToUnits("2.50", 6, "USDT")
	if err != nil {
		t.Fatalf("parseDecimalAmountToUnits: %v", err)
	}
	if units != 2_500_000 {
		t.Fatalf("USDT units = %d, want 2500000", units)
	}
	unique := addUniqueUSDTAmountOffset(units, "usdt_trc20_test")
	if unique <= units || unique > units+9_999 {
		t.Fatalf("unique amount offset out of range: %d -> %d", units, unique)
	}
	if got := formatUnitsAmount(unique, 6); !strings.HasPrefix(got, "2.50") {
		t.Fatalf("formatted unique amount = %q", got)
	}
}

func TestCryptoUSDTAmountUsesTonAPIRateWhenExplicitAmountIsEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rates" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("tokens"); got != "USDT_MASTER" {
			t.Fatalf("tokens = %q", got)
		}
		if got := r.URL.Query().Get("currencies"); got != "rub" {
			t.Fatalf("currencies = %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer tonapi-token" {
			t.Fatalf("Authorization = %q", got)
		}
		_, _ = w.Write([]byte(`{"rates":{"USDT_MASTER":{"prices":{"RUB":75}}}}`))
	}))
	defer server.Close()

	cfg := config{
		PremiumRubPrice:           300,
		CryptoTONAPIBaseURL:       server.URL,
		CryptoTONAPIKey:           "tonapi-token",
		CryptoUSDTTONJettonMaster: "USDT_MASTER",
		CryptoUSDTRubRate:         72,
	}
	b := &bot{cfg: cfg, cryptoRates: newCryptoRateProvider(cfg, server.Client())}
	amount, err := b.cryptoUSDTAmountForProduct(context.Background(), premiumMonthlyProduct)
	if err != nil {
		t.Fatalf("cryptoUSDTAmountForProduct: %v", err)
	}
	if amount != 4*usdtUnitsPerUSDT {
		t.Fatalf("USDT amount = %d, want %d", amount, 4*usdtUnitsPerUSDT)
	}
}

func TestCryptoPaymentTextUsesUserTimezone(t *testing.T) {
	payment := cryptoPayment{
		Amount:    "1.25",
		Address:   "EQ_TEST",
		Memo:      "POLIGLOT:ton_test",
		ExpiresAt: time.Date(2026, 5, 19, 7, 44, 0, 0, time.UTC),
	}
	user := userState{
		InterfaceLanguage: "en",
		ReminderUTCOffset: 180,
		TimezoneSelected:  true,
	}
	text := cryptoPaymentText(user, premiumPlan{Title: "Premium"}, payment, "https://example.test/invoice")
	if !strings.Contains(text, "Invoice expires at: 2026-05-19 10:44 (UTC+3)") {
		t.Fatalf("expected user-local expiry, got:\n%s", text)
	}
}

func TestCryptoPaymentTextKeepsSelectedUTCOffset(t *testing.T) {
	payment := cryptoPayment{
		Amount:    "1.25",
		Address:   "EQ_TEST",
		Memo:      "POLIGLOT:ton_test",
		ExpiresAt: time.Date(2026, 5, 19, 7, 44, 0, 0, time.UTC),
	}
	user := userState{
		InterfaceLanguage: "en",
		ReminderUTCOffset: 0,
		TimezoneSelected:  true,
	}
	text := cryptoPaymentText(user, premiumPlan{Title: "Premium"}, payment, "https://example.test/invoice")
	if !strings.Contains(text, "Invoice expires at: 2026-05-19 07:44 (UTC+0)") {
		t.Fatalf("expected selected UTC offset to be preserved, got:\n%s", text)
	}
}

func TestSQLiteCryptoPaymentLifecycle(t *testing.T) {
	store, err := newSQLiteStore(filepath.Join(t.TempDir(), "test.sqlite"), "")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.db.Close() })

	now := time.Now().UTC()
	payment := cryptoPayment{
		ID:         "ton_test",
		Provider:   cryptoProviderDirect,
		Currency:   cryptoCurrencyTON,
		Network:    cryptoNetworkTON,
		Product:    premiumMonthlyProduct,
		Status:     cryptoStatusPending,
		Address:    "EQ_TEST",
		Memo:       "POLIGLOT:ton_test",
		Amount:     "1.25",
		AmountNano: 1_250_000_000,
		PriceRub:   300,
		TelegramID: -42,
		ExpiresAt:  now.Add(time.Hour),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := store.createCryptoPayment(payment); err != nil {
		t.Fatalf("createCryptoPayment: %v", err)
	}
	loaded, ok, err := store.getCryptoPayment(payment.ID)
	if err != nil || !ok {
		t.Fatalf("getCryptoPayment ok=%v err=%v", ok, err)
	}
	if loaded.Memo != payment.Memo || loaded.AmountNano != payment.AmountNano {
		t.Fatalf("loaded payment mismatch: %#v", loaded)
	}
	loaded.Status = cryptoStatusPaid
	loaded.TxHash = "hash"
	if err := store.updateCryptoPayment(loaded); err != nil {
		t.Fatalf("updateCryptoPayment: %v", err)
	}
	paid, ok, err := store.getCryptoPayment(payment.ID)
	if err != nil || !ok {
		t.Fatalf("get updated payment ok=%v err=%v", ok, err)
	}
	if paid.Status != cryptoStatusPaid || paid.TxHash != "hash" {
		t.Fatalf("updated payment mismatch: %#v", paid)
	}
}

func TestDirectTONCheckFallsBackToTonAPIWhenTonCenterRejectsKey(t *testing.T) {
	now := time.Now().UTC()
	payment := cryptoPayment{
		ID:         "ton_test",
		Currency:   cryptoCurrencyTON,
		Network:    cryptoNetworkTON,
		Address:    "EQ_TEST",
		Memo:       "POLIGLOT:ton_test",
		AmountNano: 1_250_000_000,
		CreatedAt:  now.Add(-time.Minute),
		ExpiresAt:  now.Add(time.Hour),
	}
	tonCenter := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"bad api key"}`))
	}))
	defer tonCenter.Close()
	tonAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer tonconsole-token" {
			t.Fatalf("TonAPI Authorization header = %q", got)
		}
		_, _ = w.Write([]byte(`{"transactions":[{"hash":"tx-hash","utime":` + strconv.FormatInt(now.Unix(), 10) + `,"in_msg":{"value":"1250000000","decoded_body":{"comment":"POLIGLOT:ton_test"}}}]}`))
	}))
	defer tonAPI.Close()
	b := &bot{cfg: config{
		CryptoTONCenterBaseURL: tonCenter.URL,
		CryptoTONCenterAPIKey:  "wrong-provider-token",
		CryptoTONAPIBaseURL:    tonAPI.URL,
		CryptoTONAPIKey:        "tonconsole-token",
	}}
	confirmation, err := b.confirmDirectTONPayment(context.Background(), payment)
	if err != nil {
		t.Fatalf("confirmDirectTONPayment: %v", err)
	}
	if !confirmation.Paid || confirmation.TxHash != "tx-hash" {
		t.Fatalf("unexpected confirmation: %#v", confirmation)
	}
}

func TestDirectUSDTTONCheckWithTonCenterV3(t *testing.T) {
	now := time.Now().UTC()
	payment := cryptoPayment{
		ID:         "usdt_ton_test",
		Currency:   cryptoCurrencyUSDT,
		Network:    cryptoNetworkTON,
		Address:    "EQ_TEST",
		Memo:       "POLIGLOT:usdt_ton_test",
		AmountNano: 2_500_000,
		CreatedAt:  now.Add(-time.Minute),
		ExpiresAt:  now.Add(time.Hour),
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/jetton/transfers" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("jetton_master"); got != "USDT_MASTER" {
			t.Fatalf("jetton_master = %q", got)
		}
		_, _ = w.Write([]byte(`{"jetton_transfers":[{"transaction_hash":"jetton-hash","transaction_now":` + strconv.FormatInt(now.Unix(), 10) + `,"destination":"EQ_TEST","jetton_master":"USDT_MASTER","amount":"2500000","decoded_forward_payload":{"comment":"POLIGLOT:usdt_ton_test"}}]}`))
	}))
	defer server.Close()
	b := &bot{cfg: config{
		CryptoTONCenterBaseURL:    server.URL + "/api/v2",
		CryptoUSDTTONJettonMaster: "USDT_MASTER",
	}}
	confirmation, err := b.confirmDirectUSDTTONPayment(context.Background(), payment)
	if err != nil {
		t.Fatalf("confirmDirectUSDTTONPayment: %v", err)
	}
	if !confirmation.Paid || confirmation.TxHash != "TON:jetton-hash" {
		t.Fatalf("unexpected confirmation: %#v", confirmation)
	}
}

func TestDirectUSDTTRC20CheckMatchesExactUniqueAmount(t *testing.T) {
	now := time.Now().UTC()
	payment := cryptoPayment{
		ID:         "usdt_trc20_test",
		Currency:   cryptoCurrencyUSDT,
		Network:    cryptoNetworkTRC20,
		Address:    "TFoRoWRxgYWzmnJAwsom51DmgnNyYD368M",
		AmountNano: 2_501_234,
		CreatedAt:  now.Add(-time.Minute),
		ExpiresAt:  now.Add(time.Hour),
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("TRON-PRO-API-KEY"); got != "tron-key" {
			t.Fatalf("TRON-PRO-API-KEY = %q", got)
		}
		if got := r.URL.Query().Get("contract_address"); got != "TRC20_USDT" {
			t.Fatalf("contract_address = %q", got)
		}
		_, _ = w.Write([]byte(`{"success":true,"data":[{"transaction_id":"tron-hash","block_timestamp":` + strconv.FormatInt(now.UnixMilli(), 10) + `,"to":"TFoRoWRxgYWzmnJAwsom51DmgnNyYD368M","value":"2501234","token_info":{"address":"TRC20_USDT"}}]}`))
	}))
	defer server.Close()
	b := &bot{cfg: config{
		CryptoTronGridBaseURL:   server.URL,
		CryptoTronGridAPIKey:    "tron-key",
		CryptoUSDTTRC20Contract: "TRC20_USDT",
	}}
	confirmation, err := b.confirmDirectUSDTTRC20Payment(context.Background(), payment)
	if err != nil {
		t.Fatalf("confirmDirectUSDTTRC20Payment: %v", err)
	}
	if !confirmation.Paid || confirmation.TxHash != "TRON:tron-hash" {
		t.Fatalf("unexpected confirmation: %#v", confirmation)
	}
}
