package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	cryptoRateCacheTTL      = 30 * time.Minute
	cryptoRateRequestTimout = 5 * time.Second
)

type cryptoRateSnapshot struct {
	USDTRub   float64
	Source    string
	FetchedAt time.Time
	Fallback  bool
}

type cryptoRateProvider struct {
	cfg    config
	http   *http.Client
	mu     sync.Mutex
	cached cryptoRateSnapshot
}

func newCryptoRateProvider(cfg config, httpClient *http.Client) *cryptoRateProvider {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &cryptoRateProvider{cfg: cfg, http: httpClient}
}

func fallbackCryptoRateSnapshot(cfg config) cryptoRateSnapshot {
	return cryptoRateSnapshot{
		USDTRub:   fallbackUSDTRubRate(cfg),
		Source:    "fallback",
		FetchedAt: time.Now().UTC(),
		Fallback:  true,
	}
}

func fallbackUSDTRubRate(cfg config) float64 {
	if cfg.CryptoUSDTRubRate > 0 {
		return float64(cfg.CryptoUSDTRubRate)
	}
	return float64(defaultUSDTExchangeRateRubles)
}

func (p *cryptoRateProvider) usdtRubRate(ctx context.Context) cryptoRateSnapshot {
	if p == nil {
		return fallbackCryptoRateSnapshot(config{})
	}
	now := time.Now().UTC()
	p.mu.Lock()
	cached := p.cached
	if cached.USDTRub > 0 && now.Sub(cached.FetchedAt) < cryptoRateCacheTTL {
		p.mu.Unlock()
		return cached
	}
	p.mu.Unlock()

	rate, err := p.fetchUSDTRubRate(ctx)
	if err != nil {
		p.mu.Lock()
		cached = p.cached
		p.mu.Unlock()
		if cached.USDTRub > 0 {
			return cached
		}
		return fallbackCryptoRateSnapshot(p.cfg)
	}

	snapshot := cryptoRateSnapshot{
		USDTRub:   rate,
		Source:    "tonapi",
		FetchedAt: time.Now().UTC(),
	}
	p.mu.Lock()
	p.cached = snapshot
	p.mu.Unlock()
	return snapshot
}

func (p *cryptoRateProvider) fetchUSDTRubRate(ctx context.Context) (float64, error) {
	baseURL := strings.TrimSpace(p.cfg.CryptoTONAPIBaseURL)
	if baseURL == "" {
		return 0, errors.New("TonAPI base URL is not configured")
	}
	token := strings.TrimSpace(p.cfg.CryptoUSDTTONJettonMaster)
	if token == "" {
		return 0, errors.New("USDT TON jetton master is not configured")
	}
	endpoint, err := url.Parse(strings.TrimRight(baseURL, "/") + "/rates")
	if err != nil {
		return 0, err
	}
	values := endpoint.Query()
	values.Set("tokens", token)
	values.Set("currencies", "rub")
	endpoint.RawQuery = values.Encode()

	requestCtx, cancel := context.WithTimeout(ctx, cryptoRateRequestTimout)
	defer cancel()
	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return 0, err
	}
	if apiKey := strings.TrimSpace(p.cfg.CryptoTONAPIKey); apiKey != "" {
		request.Header.Set("Authorization", "Bearer "+apiKey)
	}
	response, err := p.http.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 512))
		return 0, fmt.Errorf("TonAPI rates returned %d: %s", response.StatusCode, strings.TrimSpace(string(body)))
	}
	var decoded struct {
		Rates map[string]struct {
			Prices map[string]float64 `json:"prices"`
		} `json:"rates"`
	}
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return 0, err
	}
	for key, rate := range decoded.Rates {
		if key == token || strings.EqualFold(key, token) {
			return priceFromTonAPIRate(rate.Prices)
		}
	}
	for _, rate := range decoded.Rates {
		return priceFromTonAPIRate(rate.Prices)
	}
	return 0, errors.New("TonAPI rates response does not contain USDT")
}

func priceFromTonAPIRate(prices map[string]float64) (float64, error) {
	for currency, price := range prices {
		if strings.EqualFold(currency, "rub") {
			if price <= 0 || math.IsNaN(price) || math.IsInf(price, 0) {
				return 0, errors.New("TonAPI returned invalid USDT/RUB price")
			}
			return price, nil
		}
	}
	return 0, errors.New("TonAPI rates response does not contain RUB price")
}
