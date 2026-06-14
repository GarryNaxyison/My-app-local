package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"embed"
)

const (
	webSessionCookieName      = "poliglot_web_session"
	modeWebWordPrefix         = "webword:"
	modeWebWordGamePrefix     = "webgame:"
	modeWebSpellingPrefix     = "webspelling:"
	modeWebSpellingSkipPrefix = "webspelling-skip:"
	modeWebMistakePrefix      = "webmistake:"
	webRateSweepLimit         = 2048
	webMaxVoiceUploadBytes    = 24 << 20
	webMaxImageUploadBytes    = 10 << 20
	webMaxMultipartOverhead   = 1 << 20
	webTurnstileVerifyURL     = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
)

var errWebCaptchaRequired = errors.New("captcha verification required")

//go:embed web/assets/logo.jpg
var webLogoJPG []byte

//go:embed web/assets/logo2.jpg
var webLogo2JPG []byte

//go:embed web/assets/header-*.png web/assets/icon-*.png web/assets/plan-*-*.png web/assets/panel-*.png web/assets/app-background-*.png web/assets/award-*.png web/assets/brand-*.png web/assets/brand-assets-manifest.json
var webGeneratedAssets embed.FS

type webAPI struct {
	cfg                config
	bot                *bot
	sessionSecret      []byte
	allowedOrigins     map[string]bool
	allowAnyOrigin     bool
	rateMu             sync.Mutex
	rateBuckets        map[string]webRateBucket
	pronunciationMu    sync.Mutex
	pronunciationCache map[string][]byte
}

type webRateBucket struct {
	Count   int
	ResetAt time.Time
}

type webLanguageDTO struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	NativeName    string `json:"native_name"`
	InterfaceName string `json:"interface_name"`
}

type webLevelQuestionDTO struct {
	Index    int      `json:"index"`
	Total    int      `json:"total"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Score    int      `json:"score"`
}

type webWordOptionDTO struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type webMistakeDTO struct {
	Index       int    `json:"index"`
	Word        string `json:"word"`
	Correction  string `json:"correction"`
	Explanation string `json:"explanation"`
	AddedAt     string `json:"added_at,omitempty"`
}

func newWebAPI(cfg config, bot *bot) *webAPI {
	if bot != nil && bot.webAuth == nil {
		bot.webAuth = newWebAuthBroker()
	}
	secret := strings.TrimSpace(cfg.WebAPISessionSecret)
	if secret == "" {
		secret = strings.TrimSpace(cfg.YooKassaWebhookKey)
	}
	if secret == "" {
		secret = strings.TrimSpace(cfg.OpenRouterAPIKey)
	}
	if secret == "" {
		secret = strings.TrimSpace(cfg.TelegramBotToken)
	}
	if secret == "" {
		var randomSecret [32]byte
		if _, err := rand.Read(randomSecret[:]); err == nil {
			secret = base64.RawStdEncoding.EncodeToString(randomSecret[:])
		}
	}
	api := &webAPI{
		cfg:                cfg,
		bot:                bot,
		sessionSecret:      []byte(secret),
		allowedOrigins:     map[string]bool{},
		rateBuckets:        map[string]webRateBucket{},
		pronunciationCache: map[string][]byte{},
	}
	for _, origin := range cfg.WebCORSOrigins {
		origin = strings.TrimRight(strings.TrimSpace(origin), "/")
		if origin == "*" {
			api.allowAnyOrigin = true
			continue
		}
		if origin != "" {
			api.allowedOrigins[origin] = true
		}
	}
	return api
}

func (api *webAPI) register(mux *http.ServeMux) {
	mux.HandleFunc("/", api.handleRoot)
	mux.HandleFunc("/login", api.handleWebLogin)
	mux.HandleFunc("/app", api.handleWebApp)
	mux.HandleFunc("/app/login", api.handleWebLogin)
	mux.HandleFunc("/app/v1", api.handleWebVersionRedirect)
	mux.HandleFunc("/app/v1/", api.handleWebVersionRedirect)
	mux.HandleFunc("/app/v2", api.handleWebVersionRedirect)
	mux.HandleFunc("/app/v2/", api.handleWebVersionRedirect)
	mux.HandleFunc("/manifest.webmanifest", api.handleWebManifest)
	mux.HandleFunc("/app/manifest.webmanifest", api.handleWebManifest)
	mux.HandleFunc("/offline-deck-sw.js", api.handleWebOfflineDeckServiceWorker)
	mux.HandleFunc("/app/offline-deck-sw.js", api.handleWebOfflineDeckServiceWorker)
	mux.HandleFunc("/app/", api.handleWebApp)
	mux.HandleFunc("/assets/logo.jpg", api.handleLogoAsset)
	mux.HandleFunc("/assets/logo2.jpg", api.handleLogo2Asset)
	mux.HandleFunc("/app/assets/logo.jpg", api.handleLogoAsset)
	mux.HandleFunc("/app/assets/logo2.jpg", api.handleLogo2Asset)
	mux.HandleFunc("/assets/", api.handleGeneratedAsset)
	mux.HandleFunc("/app/assets/", api.handleGeneratedAsset)
	mux.HandleFunc("/favicon.ico", api.handleFaviconAsset)
	mux.HandleFunc("/api/health", api.handleAPIHealth)
	mux.HandleFunc("/api/auth/register", api.handleAuthRegister)
	mux.HandleFunc("/api/auth/login", api.handleAuthLogin)
	mux.HandleFunc("/api/auth/telegram/start", api.handleAuthTelegramStart)
	mux.HandleFunc("/api/auth/telegram/status", api.handleAuthTelegramStatus)
	mux.HandleFunc("/api/auth/telegram/merge", api.handleAuthTelegramMerge)
	mux.HandleFunc("/api/auth/telegram", api.handleAuthTelegram)
	mux.HandleFunc("/api/auth/profile", api.handleAuthProfile)
	mux.HandleFunc("/api/auth/password", api.handleAuthPassword)
	mux.HandleFunc("/api/auth/logout", api.handleAuthLogout)
	mux.HandleFunc("/api/session", api.handleSession)
	mux.HandleFunc("/api/settings", api.handleSettings)
	mux.HandleFunc("/api/navigation-layout", api.handleNavigationLayout)
	mux.HandleFunc("/api/daily/claim", api.handleDailyClaim)
	mux.HandleFunc("/api/leaderboard", api.handleLeaderboard)
	mux.HandleFunc("/api/tutor/start", api.handleTutorStart)
	mux.HandleFunc("/api/ai-tutor/start", api.handleAITutorStart)
	mux.HandleFunc("/api/ai-tutor/session", api.handleAITutorSession)
	mux.HandleFunc("/api/ai-tutor/completed", api.handleAITutorCompleted)
	mux.HandleFunc("/api/ai-tutor/restart", api.handleAITutorRestart)
	mux.HandleFunc("/api/ai-tutor/answer", api.handleAITutorAnswer)
	mux.HandleFunc("/api/ai-tutor/word-report", api.handleAITutorWordReport)
	mux.HandleFunc("/api/ai-tutor/review", api.handleAITutorReview)
	mux.HandleFunc("/api/ai-tutor/finish", api.handleAITutorFinish)
	mux.HandleFunc("/api/lesson/start", api.handleLessonStart)
	mux.HandleFunc("/api/lesson/answer", api.handleLessonAnswer)
	mux.HandleFunc("/api/practice", api.handlePractice)
	mux.HandleFunc("/api/shadowing/start", api.handleShadowingStart)
	mux.HandleFunc("/api/shadowing/answer", api.handleShadowingAnswer)
	mux.HandleFunc("/api/pronunciation/start", api.handlePronunciationStart)
	mux.HandleFunc("/api/pronunciation/check", api.handlePronunciationCheck)
	mux.HandleFunc("/api/level-test/start", api.handleLevelTestStart)
	mux.HandleFunc("/api/level-test/answer", api.handleLevelTestAnswer)
	mux.HandleFunc("/api/words/next", api.handleWordNext)
	mux.HandleFunc("/api/words/answer", api.handleWordAnswer)
	mux.HandleFunc("/api/words/pronunciation", api.handleWordPronunciation)
	mux.HandleFunc("/api/vocabulary", api.handleVocabulary)
	mux.HandleFunc("/api/phrasebook", api.handlePhrasebook)
	mux.HandleFunc("/api/bug-report", api.handleBugReport)
	mux.HandleFunc("/api/word-game/next", api.handleWordGameNext)
	mux.HandleFunc("/api/word-game/answer", api.handleWordGameAnswer)
	mux.HandleFunc("/api/spelling/start", api.handleSpellingStart)
	mux.HandleFunc("/api/spelling/answer", api.handleSpellingAnswer)
	mux.HandleFunc("/api/mistakes", api.handleMistakes)
	mux.HandleFunc("/api/mistakes/practice/start", api.handleMistakePracticeStart)
	mux.HandleFunc("/api/mistakes/practice/answer", api.handleMistakePracticeAnswer)
	mux.HandleFunc("/api/mistakes/clear", api.handleMistakesClear)
	mux.HandleFunc("/api/tools/voice-text", api.handleToolVoiceText)
	mux.HandleFunc("/api/tools/image-translate", api.handleToolImageTranslate)
	mux.HandleFunc("/api/tools/translator", api.handleToolTranslator)
	mux.HandleFunc("/api/tools/translator-speech", api.handleToolTranslatorSpeech)
	mux.HandleFunc("/api/progress", api.handleProgress)
	mux.HandleFunc("/api/premium/plans", api.handlePremiumPlans)
	mux.HandleFunc("/api/premium/payment", api.handlePremiumPayment)
	mux.HandleFunc("/api/premium/stars", api.handlePremiumStarsPayment)
	mux.HandleFunc("/api/premium/crypto/payment", api.handleCryptoPremiumPayment)
	mux.HandleFunc("/api/premium/crypto/check", api.handleCryptoPremiumCheck)
	mux.HandleFunc("/api/premium/activation-key", api.handlePremiumActivationKey)
}

func (api *webAPI) webYooKassaClient() *yooKassaClient {
	if api.bot == nil {
		return nil
	}
	if api.bot.webYooKassa != nil {
		return api.bot.webYooKassa
	}
	if !api.cfg.webYooKassaOverride() {
		return api.bot.yookassa
	}
	return nil
}

func (api *webAPI) webYooKassaReady() bool {
	return api.webYooKassaClient() != nil && api.cfg.webYooKassaEnabled()
}

func (api *webAPI) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			if !api.applyCORS(w, r) {
				return
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			if r.URL.Path != "/api/health" && !api.allowAPITraffic(r) {
				writeAPIError(w, http.StatusTooManyRequests, "Too many requests. Try again later.")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (api *webAPI) allowAPITraffic(r *http.Request) bool {
	if !api.allowRate(r, "api-ip:1m", 240, time.Minute) {
		return false
	}
	if !api.allowRate(r, "api-ip:10m", 1200, 10*time.Minute) {
		return false
	}
	if userID, ok := api.sessionUserID(r); ok {
		return api.allowRate(r, "api-user:"+strconv.FormatInt(userID, 10), 180, time.Minute)
	}
	return api.allowRate(r, "api-anon", 80, time.Minute)
}

func (api *webAPI) applyCORS(w http.ResponseWriter, r *http.Request) bool {
	origin := strings.TrimRight(strings.TrimSpace(r.Header.Get("Origin")), "/")
	if origin == "" {
		return true
	}
	if !api.originAllowed(origin, r.Host) {
		http.Error(w, "origin not allowed", http.StatusForbidden)
		return false
	}
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Vary", "Origin")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	return true
}

func (api *webAPI) originAllowed(origin string, requestHost string) bool {
	if api.allowAnyOrigin || api.allowedOrigins[origin] {
		return true
	}
	parsed, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := strings.TrimSpace(requestHost)
	return host != "" && strings.EqualFold(parsed.Host, host)
}

func (api *webAPI) webCaptchaEnabled() bool {
	return strings.TrimSpace(api.cfg.WebTurnstileSiteKey) != "" && strings.TrimSpace(api.cfg.WebTurnstileSecretKey) != ""
}

func (api *webAPI) webCaptchaDTO() map[string]any {
	if !api.webCaptchaEnabled() {
		return map[string]any{"enabled": false}
	}
	return map[string]any{
		"enabled":  true,
		"provider": "turnstile",
		"site_key": strings.TrimSpace(api.cfg.WebTurnstileSiteKey),
	}
}

func (api *webAPI) verifyWebCaptcha(ctx context.Context, r *http.Request, token string) error {
	if !api.webCaptchaEnabled() {
		return nil
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return errWebCaptchaRequired
	}
	form := url.Values{}
	form.Set("secret", strings.TrimSpace(api.cfg.WebTurnstileSecretKey))
	form.Set("response", token)
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil && host != "" {
		form.Set("remoteip", host)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webTurnstileVerifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var result struct {
		Success bool `json:"success"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&result); err != nil {
		return err
	}
	if !result.Success {
		return errWebCaptchaRequired
	}
	return nil
}

func (api *webAPI) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	api.handleWebApp(w, r)
}

func (api *webAPI) handleWebApp(w http.ResponseWriter, r *http.Request) {
	api.writeWebApp(w, r, webAppShellFromRequest(r))
}

func (api *webAPI) handleWebLogin(w http.ResponseWriter, r *http.Request) {
	api.writeWebAppWithMode(w, r, webAppShellFromRequest(r), true)
}

func (api *webAPI) handleWebVersionRedirect(w http.ResponseWriter, r *http.Request) {
	target := "/app"
	if strings.HasSuffix(strings.TrimRight(strings.ToLower(r.URL.Path), "/"), "/login") {
		target = "/app/login"
	}
	if r.URL.RawQuery != "" {
		target += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, target, http.StatusMovedPermanently)
}

func (api *webAPI) handleWebManifest(w http.ResponseWriter, r *http.Request) {
	api.writeWebStaticFile(w, r, "application/manifest+json; charset=utf-8", "web", "manifest.webmanifest")
}

func (api *webAPI) handleWebOfflineDeckServiceWorker(w http.ResponseWriter, r *http.Request) {
	api.writeWebStaticFile(w, r, "text/javascript; charset=utf-8", "web", "offline-deck-sw.js")
}

func (api *webAPI) writeWebStaticFile(w http.ResponseWriter, r *http.Request, contentType string, parts ...string) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, err := readWebFile(parts...)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	if r.Method != http.MethodHead {
		_, _ = w.Write(data)
	}
}

func (api *webAPI) writeWebApp(w http.ResponseWriter, r *http.Request, shell string) {
	api.writeWebAppWithMode(w, r, shell, false)
}

func (api *webAPI) writeWebAppWithMode(w http.ResponseWriter, r *http.Request, shell string, authStandalone bool) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'; base-uri 'self'; form-action 'self'")
	w.Header().Set("Cache-Control", "no-store")
	html, err := readWebFile("web", "index.html")
	if err != nil {
		http.Error(w, "web app file not found", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(prepareWebAppHTMLWithMode(html, shell, authStandalone))
}

func webAppShellFromRequest(r *http.Request) string {
	path := strings.TrimRight(strings.ToLower(r.URL.Path), "/")
	if strings.HasSuffix(path, "/mobile") {
		return "mobile"
	}
	if strings.HasSuffix(path, "/desktop") {
		return "desktop"
	}
	switch strings.ToLower(strings.TrimSpace(r.URL.Query().Get("shell"))) {
	case "mobile":
		return "mobile"
	case "desktop":
		return "desktop"
	}
	if isMobileUserAgent(r.UserAgent()) {
		return "mobile"
	}
	return "desktop"
}

func isMobileUserAgent(userAgent string) bool {
	ua := strings.ToLower(userAgent)
	for _, token := range []string{"android", "iphone", "ipod", "ipad", "mobile", "windows phone"} {
		if strings.Contains(ua, token) {
			return true
		}
	}
	return false
}

func prepareWebAppHTML(html []byte, shell string) []byte {
	return prepareWebAppHTMLWithMode(html, shell, false)
}

func prepareWebAppHTMLWithMode(html []byte, shell string, authStandalone bool) []byte {
	if shell != "mobile" {
		shell = "desktop"
	}
	device := shell
	text := string(html)
	const marker = `<html lang="ru" data-app-shell="auto" data-app-device="desktop">`
	authValue := "false"
	if authStandalone {
		authValue = "true"
	}
	replacement := `<html lang="ru" data-app-shell="` + shell + `" data-app-device="` + device + `" data-auth-standalone="` + authValue + `">`
	if strings.Contains(text, marker) {
		text = strings.Replace(text, marker, replacement, 1)
	} else {
		text = strings.Replace(text, `<html lang="ru">`, replacement, 1)
	}
	return []byte(text)
}

func (api *webAPI) handleLogoAsset(w http.ResponseWriter, r *http.Request) {
	api.writeAsset(w, r, "image/jpeg", "logo.jpg", webLogoJPG)
}

func (api *webAPI) handleLogo2Asset(w http.ResponseWriter, r *http.Request) {
	api.writeAsset(w, r, "image/jpeg", "logo2.jpg", webLogo2JPG)
}

func (api *webAPI) handleFaviconAsset(w http.ResponseWriter, r *http.Request) {
	if data, err := readWebAsset("brand-logo-mini.png"); err == nil {
		api.writeAssetBytes(w, r, "image/png", data)
		return
	}
	if data, err := webGeneratedAssets.ReadFile("web/assets/brand-logo-mini.png"); err == nil {
		api.writeAssetBytes(w, r, "image/png", data)
		return
	}
	api.writeAsset(w, r, "image/jpeg", "logo2.jpg", webLogo2JPG)
}

func (api *webAPI) handleGeneratedAsset(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet, http.MethodHead) {
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/app/assets/")
	if name == r.URL.Path {
		name = strings.TrimPrefix(r.URL.Path, "/assets/")
	}
	if name == "" || strings.Contains(name, "/") || strings.Contains(name, "\\") || strings.Contains(name, "..") {
		http.NotFound(w, r)
		return
	}
	contentType := webAssetContentType(name)
	data, err := readWebAsset(name)
	if err != nil {
		data, err = webGeneratedAssets.ReadFile("web/assets/" + name)
	}
	if err != nil {
		if legacy := legacyThemeAssetName(name); legacy != "" {
			data, err = readWebAsset(legacy)
			if err != nil {
				data, err = webGeneratedAssets.ReadFile("web/assets/" + legacy)
			}
		}
	}
	if err != nil {
		http.NotFound(w, r)
		return
	}
	api.writeAssetBytes(w, r, contentType, data)
}

func webAssetContentType(name string) string {
	switch strings.ToLower(strings.TrimSpace(path.Base(name))) {
	case "":
		return "application/octet-stream"
	}
	switch strings.ToLower(path.Ext(name)) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".js", ".mjs":
		return "text/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".woff2":
		return "font/woff2"
	case ".json":
		return "application/json; charset=utf-8"
	default:
		return "application/octet-stream"
	}
}

func legacyThemeAssetName(name string) string {
	if !strings.HasSuffix(name, ".png") {
		return ""
	}
	stem := strings.TrimSuffix(name, ".png")
	if strings.HasSuffix(stem, "-light") || strings.HasSuffix(stem, "-dark") {
		return ""
	}
	if strings.HasPrefix(stem, "header-") || strings.HasPrefix(stem, "icon-") || strings.HasPrefix(stem, "plan-") {
		return stem + "-light.png"
	}
	return ""
}

func (api *webAPI) writeAsset(w http.ResponseWriter, r *http.Request, contentType string, name string, embedded []byte) {
	if !allowMethod(w, r, http.MethodGet, http.MethodHead) {
		return
	}
	data, err := readWebAsset(name)
	if err != nil {
		data = embedded
	}
	api.writeAssetBytes(w, r, contentType, data)
}

func (api *webAPI) writeAssetBytes(w http.ResponseWriter, r *http.Request, contentType string, data []byte) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "public, max-age=604800")
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func readWebAsset(name string) ([]byte, error) {
	return readWebFile("web", "assets", name)
}

func readWebFile(parts ...string) ([]byte, error) {
	paths := []string{filepath.Join(parts...)}
	if exe, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(append([]string{filepath.Dir(exe)}, parts...)...))
	}
	var lastErr error
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil && len(data) > 0 {
			return data, nil
		}
		if err != nil {
			lastErr = err
		}
	}
	if lastErr == nil {
		lastErr = os.ErrNotExist
	}
	return nil, lastErr
}

func (api *webAPI) handleAPIHealth(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (api *webAPI) handleAuthRegister(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	store, err := api.webAccountStore()
	if err != nil {
		writeAPIError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	var req struct {
		Login           string `json:"login"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"password_confirm"`
		ReferralCode    string `json:"referral_code"`
		CaptchaToken    string `json:"captcha_token"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	if !api.allowRate(r, "auth-register", 20, 10*time.Minute) ||
		!api.allowRate(r, "auth-register:"+strings.ToLower(strings.TrimSpace(req.Login)), 8, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, "Слишком много попыток. Попробуйте чуть позже.")
		return
	}
	if err := api.verifyWebCaptcha(r.Context(), r, req.CaptchaToken); err != nil {
		writeAPIError(w, http.StatusBadRequest, "Confirm that you are not a robot.")
		return
	}
	login, err := normalizeWebLogin(req.Login)
	if err != nil {
		writeAPIErrorCode(w, http.StatusBadRequest, webLoginValidationCode(req.Login), err.Error())
		return
	}
	if err := validateWebPassword(req.Password); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateWebPasswordConfirmation(req.Password, req.PasswordConfirm); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	referralCode, err := normalizeWebReferralCode(req.ReferralCode)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	passwordHash, err := hashWebPassword(req.Password)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	userID, ok := api.sessionUserID(r)
	if ok {
		if _, linked, err := store.getWebAccountByUserID(userID); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		} else if linked {
			ok = false
		}
	}
	if !ok || userID == 0 {
		userID, err = newWebUserID()
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	var referralStore webReferralStore
	if referralCode != "" {
		var supportsReferral bool
		referralStore, supportsReferral = api.bot.store.(webReferralStore)
		if !supportsReferral {
			writeAPIError(w, http.StatusServiceUnavailable, errWebAccountUnavailable.Error())
			return
		}
		inviter, found, err := referralStore.getUserByReferralCode(referralCode)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if !found {
			writeAPIError(w, http.StatusBadRequest, errWebReferralNotFound.Error())
			return
		}
		if inviter.TelegramID == userID {
			writeAPIError(w, http.StatusBadRequest, errWebReferralSelf.Error())
			return
		}
	}
	account, err := store.createWebAccount(userID, login, passwordHash)
	if errors.Is(err, errWebLoginTaken) {
		writeAPIError(w, http.StatusConflict, "Такой логин уже занят")
		return
	}
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := api.bot.store.getOrCreateUser(account.UserID, account.Login)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if referralCode != "" {
		if _, err := referralStore.applyWebReferral(account.UserID, referralCode); err != nil {
			if errors.Is(err, errWebReferralNotFound) || errors.Is(err, errWebReferralSelf) {
				writeAPIError(w, http.StatusBadRequest, err.Error())
				return
			}
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		user, err = api.bot.store.getOrCreateUser(account.UserID, account.Login)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	api.setSessionCookie(w, account.UserID)
	api.writeSession(w, user)
}

func (api *webAPI) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	store, err := api.webAccountStore()
	if err != nil {
		writeAPIError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	var req struct {
		Login        string `json:"login"`
		Password     string `json:"password"`
		CaptchaToken string `json:"captcha_token"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	if !api.allowRate(r, "auth-login", 60, 10*time.Minute) ||
		!api.allowRate(r, "auth-login:"+strings.ToLower(strings.TrimSpace(req.Login)), 12, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, "Слишком много попыток входа. Попробуйте чуть позже.")
		return
	}
	if err := api.verifyWebCaptcha(r.Context(), r, req.CaptchaToken); err != nil {
		writeAPIError(w, http.StatusBadRequest, "Confirm that you are not a robot.")
		return
	}
	login, err := normalizeWebLogin(req.Login)
	if err != nil {
		writeAPIError(w, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}
	account, passwordHash, ok, err := store.getWebAccountByLogin(login)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok || !verifyWebPassword(req.Password, passwordHash) {
		writeAPIError(w, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}
	if err := store.touchWebAccountLogin(account.UserID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	user, err := api.bot.store.getOrCreateUser(account.UserID, account.Login)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.setSessionCookie(w, account.UserID)
	api.writeSession(w, user)
}

func (api *webAPI) handleAuthTelegramStart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	if api.bot.webAuth == nil {
		api.bot.webAuth = newWebAuthBroker()
	}
	if !api.allowRate(r, "auth-telegram-start", 40, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, "Слишком много попыток входа. Попробуйте чуть позже.")
		return
	}
	currentUserID, _ := api.sessionUserID(r)
	request, err := api.bot.webAuth.create(currentUserID, time.Now().UTC())
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token":      request.Token,
		"status":     "pending",
		"bot_url":    api.telegramAuthBotURL(request.Token),
		"expires_at": request.ExpiresAt.Format(time.RFC3339),
	})
}

func (api *webAPI) handleAuthTelegramStatus(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	if api.bot.webAuth == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "Telegram login is not configured")
		return
	}
	var req struct {
		Token string `json:"token"`
		Code  string `json:"code"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	now := time.Now().UTC()
	var request webAuthRequest
	var ready bool
	var err error
	if strings.TrimSpace(req.Code) != "" {
		request, err = api.bot.webAuth.verifySiteCode(req.Token, req.Code, now)
		if err != nil {
			writeAPIError(w, http.StatusBadRequest, err.Error())
			return
		}
		ready = true
	} else {
		request, ready, err = api.bot.webAuth.consume(req.Token, now)
	}
	if err != nil {
		writeAPIError(w, http.StatusGone, err.Error())
		return
	}
	if !ready {
		writeJSON(w, http.StatusOK, map[string]any{
			"status":     "pending",
			"expires_at": request.ExpiresAt.Format(time.RFC3339),
		})
		return
	}
	store, ok := api.bot.store.(webTelegramAccountStore)
	if !ok {
		writeAPIError(w, http.StatusServiceUnavailable, errWebAccountUnavailable.Error())
		return
	}
	if api.writeTelegramProgressMergeChallenge(w, store, request.CurrentUserID, request.Profile) {
		return
	}
	account, user, err := store.authenticateTelegramWebAccount(request.CurrentUserID, request.Profile)
	if err != nil {
		if errors.Is(err, errWebTelegramUsed) {
			writeAPIErrorCode(w, http.StatusConflict, "telegram_already_linked", err.Error())
			return
		}
		if errors.Is(err, errWebAccountTGLocked) {
			writeAPIError(w, http.StatusConflict, err.Error())
			return
		}
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.setSessionCookie(w, account.UserID)
	api.writeSession(w, user)
}

func (api *webAPI) telegramAuthBotURL(token string) string {
	botName := strings.TrimPrefix(strings.TrimSpace(api.cfg.WebTelegramLoginBot), "@")
	if botName == "" {
		botName = "Poliglot_AI_bot"
	}
	values := url.Values{}
	values.Set("start", webAuthStartPrefix+token)
	return "https://t.me/" + botName + "?" + values.Encode()
}

func (api *webAPI) handleAuthTelegram(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	store, ok := api.bot.store.(webTelegramAccountStore)
	if !ok {
		writeAPIError(w, http.StatusServiceUnavailable, errWebAccountUnavailable.Error())
		return
	}
	var payload map[string]any
	if !decodeJSONRequest(w, r, &payload) {
		return
	}
	if !api.allowRate(r, "auth-telegram", 40, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, "Слишком много попыток входа. Попробуйте чуть позже.")
		return
	}
	profile, err := validateTelegramWebAuth(payload, api.cfg.TelegramBotToken, time.Now())
	if err != nil {
		writeAPIError(w, http.StatusUnauthorized, err.Error())
		return
	}
	currentUserID, _ := api.sessionUserID(r)
	if api.writeTelegramProgressMergeChallenge(w, store, currentUserID, profile) {
		return
	}
	account, user, err := store.authenticateTelegramWebAccount(currentUserID, profile)
	if err != nil {
		if errors.Is(err, errWebTelegramUsed) {
			writeAPIErrorCode(w, http.StatusConflict, "telegram_already_linked", err.Error())
			return
		}
		if errors.Is(err, errWebAccountTGLocked) {
			writeAPIError(w, http.StatusConflict, err.Error())
			return
		}
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.setSessionCookie(w, account.UserID)
	api.writeSession(w, user)
}

func (api *webAPI) handleAuthTelegramMerge(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	store, ok := api.bot.store.(webTelegramAccountStore)
	if !ok {
		writeAPIError(w, http.StatusServiceUnavailable, errWebAccountUnavailable.Error())
		return
	}
	var req struct {
		Token  string `json:"token"`
		Choice string `json:"choice"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	choice := strings.ToLower(strings.TrimSpace(req.Choice))
	if choice != "web" && choice != "telegram" {
		writeAPIError(w, http.StatusBadRequest, "Choose which progress to keep: web or telegram.")
		return
	}
	currentUserID, profile, err := api.parseTelegramMergeToken(req.Token)
	if err != nil {
		writeAPIError(w, http.StatusGone, err.Error())
		return
	}
	account, user, err := store.authenticateTelegramWebAccountWithChoice(currentUserID, profile, choice)
	if err != nil {
		if errors.Is(err, errWebTelegramUsed) {
			writeAPIErrorCode(w, http.StatusConflict, "telegram_already_linked", err.Error())
			return
		}
		if errors.Is(err, errWebAccountTGLocked) {
			writeAPIError(w, http.StatusConflict, err.Error())
			return
		}
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.setSessionCookie(w, account.UserID)
	api.writeSession(w, user)
}

func (api *webAPI) writeTelegramProgressMergeChallenge(w http.ResponseWriter, store webTelegramAccountStore, currentUserID int64, profile telegramWebProfile) bool {
	challenge, err := store.telegramWebProgressMergeChallenge(currentUserID, profile)
	if err != nil {
		if errors.Is(err, errWebTelegramUsed) {
			writeAPIErrorCode(w, http.StatusConflict, "telegram_already_linked", err.Error())
			return true
		}
		if errors.Is(err, errWebAccountTGLocked) {
			writeAPIError(w, http.StatusConflict, err.Error())
			return true
		}
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return true
	}
	if !challenge.Required {
		return false
	}
	token, err := api.newTelegramMergeToken(currentUserID, profile)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return true
	}
	writeJSON(w, http.StatusConflict, map[string]any{
		"error": map[string]any{
			"message": "\u0412 Telegram \u0438 \u043d\u0430 \u0441\u0430\u0439\u0442\u0435 \u0443\u0436\u0435 \u0435\u0441\u0442\u044c \u043f\u0440\u043e\u0433\u0440\u0435\u0441\u0441. \u0412\u044b\u0431\u0435\u0440\u0438\u0442\u0435, \u043a\u0430\u043a\u043e\u0439 \u043f\u0440\u043e\u0433\u0440\u0435\u0441\u0441 \u0441\u0434\u0435\u043b\u0430\u0442\u044c \u043e\u0441\u043d\u043e\u0432\u043d\u044b\u043c \u043f\u0435\u0440\u0435\u0434 \u043e\u0431\u044a\u0435\u0434\u0438\u043d\u0435\u043d\u0438\u0435\u043c.",
		},
		"merge_required": true,
		"merge_token":    token,
		"progress_merge": challenge,
	})
	return true
}

type telegramMergeTokenPayload struct {
	CurrentUserID int64              `json:"current_user_id"`
	Profile       telegramWebProfile `json:"profile"`
	ExpiresAt     int64              `json:"expires_at"`
}

func (api *webAPI) newTelegramMergeToken(currentUserID int64, profile telegramWebProfile) (string, error) {
	payload := telegramMergeTokenPayload{
		CurrentUserID: currentUserID,
		Profile:       profile,
		ExpiresAt:     time.Now().UTC().Add(10 * time.Minute).Unix(),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, api.sessionSecret)
	_, _ = mac.Write(body)
	signature := mac.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(body) + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

func (api *webAPI) parseTelegramMergeToken(token string) (int64, telegramWebProfile, error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 2 {
		return 0, telegramWebProfile{}, errors.New("Merge request expired. Start Telegram login again.")
	}
	body, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return 0, telegramWebProfile{}, errors.New("Merge request expired. Start Telegram login again.")
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, telegramWebProfile{}, errors.New("Merge request expired. Start Telegram login again.")
	}
	mac := hmac.New(sha256.New, api.sessionSecret)
	_, _ = mac.Write(body)
	if subtle.ConstantTimeCompare(signature, mac.Sum(nil)) != 1 {
		return 0, telegramWebProfile{}, errors.New("Merge request expired. Start Telegram login again.")
	}
	var payload telegramMergeTokenPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return 0, telegramWebProfile{}, errors.New("Merge request expired. Start Telegram login again.")
	}
	if time.Now().UTC().Unix() > payload.ExpiresAt || payload.Profile.ID <= 0 {
		return 0, telegramWebProfile{}, errors.New("Merge request expired. Start Telegram login again.")
	}
	return payload.CurrentUserID, payload.Profile, nil
}

func (api *webAPI) handleAuthProfile(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	store, err := api.webAccountStore()
	if err != nil {
		writeAPIError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		Login           string `json:"login"`
		Password        string `json:"password"`
		PasswordConfirm string `json:"password_confirm"`
		CaptchaToken    string `json:"captcha_token"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	if !api.allowRate(r, "auth-profile:"+strconv.FormatInt(user.TelegramID, 10), 8, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, "Слишком много попыток. Попробуйте чуть позже.")
		return
	}
	if err := api.verifyWebCaptcha(r.Context(), r, req.CaptchaToken); err != nil {
		writeAPIError(w, http.StatusBadRequest, "Confirm that you are not a robot.")
		return
	}
	login, err := normalizeWebLogin(req.Login)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateWebPassword(req.Password); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateWebPasswordConfirmation(req.Password, req.PasswordConfirm); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	passwordHash, err := hashWebPassword(req.Password)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if _, err := store.completeWebAccountProfile(user.TelegramID, login, passwordHash); errors.Is(err, errWebLoginTaken) {
		existingAccount, existingPasswordHash, ok, lookupErr := store.getWebAccountByLogin(login)
		if lookupErr != nil {
			writeAPIError(w, http.StatusInternalServerError, lookupErr.Error())
			return
		}
		if !ok || !verifyWebPassword(req.Password, existingPasswordHash) {
			writeAPIError(w, http.StatusConflict, "Такой логин уже занят. Если это ваш аккаунт, введите пароль от него, и мы предложим привязку.")
			return
		}
		telegramStore, ok := store.(webTelegramAccountStore)
		if !ok {
			writeAPIError(w, http.StatusServiceUnavailable, errWebAccountUnavailable.Error())
			return
		}
		profile := telegramWebProfile{
			ID:        user.TelegramID,
			FirstName: user.FirstName,
			AuthDate:  time.Now().UTC().Unix(),
		}
		if api.writeTelegramProgressMergeChallenge(w, telegramStore, existingAccount.UserID, profile) {
			return
		}
		account, linkedUser, err := telegramStore.authenticateTelegramWebAccount(existingAccount.UserID, profile)
		if err != nil {
			if errors.Is(err, errWebTelegramUsed) {
				writeAPIErrorCode(w, http.StatusConflict, "telegram_already_linked", err.Error())
				return
			}
			if errors.Is(err, errWebAccountTGLocked) {
				writeAPIError(w, http.StatusConflict, err.Error())
				return
			}
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		api.setSessionCookie(w, account.UserID)
		api.writeSession(w, linkedUser)
		return
	} else if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, err := api.bot.store.getOrCreateUser(user.TelegramID, login)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.writeSession(w, refreshed)
}

func validateTelegramWebAuth(payload map[string]any, botToken string, now time.Time) (telegramWebProfile, error) {
	if strings.TrimSpace(botToken) == "" {
		return telegramWebProfile{}, errors.New("Telegram login is not configured")
	}
	hashValue := strings.ToLower(strings.TrimSpace(telegramAuthString(payload["hash"])))
	if hashValue == "" {
		return telegramWebProfile{}, errors.New("Telegram login signature is missing")
	}
	var pairs []string
	for key, value := range payload {
		if key == "hash" {
			continue
		}
		text := telegramAuthString(value)
		if text == "" {
			continue
		}
		pairs = append(pairs, key+"="+text)
	}
	sort.Strings(pairs)
	if len(pairs) == 0 {
		return telegramWebProfile{}, errors.New("Telegram login payload is empty")
	}
	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	_, _ = mac.Write([]byte(strings.Join(pairs, "\n")))
	expected := hex.EncodeToString(mac.Sum(nil))
	if subtle.ConstantTimeCompare([]byte(expected), []byte(hashValue)) != 1 {
		return telegramWebProfile{}, errors.New("Telegram login signature is invalid")
	}
	id, err := strconv.ParseInt(telegramAuthString(payload["id"]), 10, 64)
	if err != nil || id <= 0 {
		return telegramWebProfile{}, errors.New("Telegram user id is invalid")
	}
	authDate, err := strconv.ParseInt(telegramAuthString(payload["auth_date"]), 10, 64)
	if err != nil || authDate <= 0 {
		return telegramWebProfile{}, errors.New("Telegram auth date is invalid")
	}
	if now.Sub(time.Unix(authDate, 0)) > 24*time.Hour || time.Unix(authDate, 0).After(now.Add(5*time.Minute)) {
		return telegramWebProfile{}, errors.New("Telegram login has expired")
	}
	return telegramWebProfile{
		ID:        id,
		FirstName: telegramAuthString(payload["first_name"]),
		LastName:  telegramAuthString(payload["last_name"]),
		Username:  telegramAuthString(payload["username"]),
		PhotoURL:  telegramAuthString(payload["photo_url"]),
		AuthDate:  authDate,
	}, nil
}

func telegramAuthString(value any) string {
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	case float64:
		return strconv.FormatInt(int64(typed), 10)
	case int64:
		return strconv.FormatInt(typed, 10)
	case int:
		return strconv.Itoa(typed)
	case json.Number:
		return typed.String()
	default:
		return ""
	}
}

func (api *webAPI) handleAuthPassword(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	store, err := api.webAccountStore()
	if err != nil {
		writeAPIError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.allowRate(r, "auth-password:"+strconv.FormatInt(user.TelegramID, 10), 6, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, "Слишком много попыток. Попробуйте чуть позже.")
		return
	}
	var req struct {
		CurrentPassword    string `json:"current_password"`
		NewPassword        string `json:"new_password"`
		NewPasswordConfirm string `json:"new_password_confirm"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	account, ok, err := store.getWebAccountByUserID(user.TelegramID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeAPIError(w, http.StatusUnauthorized, errWebAuthRequired.Error())
		return
	}
	_, passwordHash, ok, err := store.getWebAccountByLogin(account.Login)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok || !verifyWebPassword(req.CurrentPassword, passwordHash) {
		writeAPIError(w, http.StatusUnauthorized, "Неверный текущий пароль")
		return
	}
	if err := validateWebPassword(req.NewPassword); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateWebPasswordConfirmation(req.NewPassword, req.NewPasswordConfirm); err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if verifyWebPassword(req.NewPassword, passwordHash) {
		writeAPIError(w, http.StatusBadRequest, "Новый пароль должен отличаться от текущего")
		return
	}
	newHash, err := hashWebPassword(req.NewPassword)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := store.updateWebAccountPassword(user.TelegramID, newHash); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (api *webAPI) handleAuthLogout(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	api.clearSessionCookie(w)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (api *webAPI) handleSession(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		if errors.Is(err, errWebAuthRequired) {
			api.writeAnonymousSession(w)
			return
		}
		api.writeCurrentUserError(w, err)
		return
	}
	if err := api.bot.store.recordHabitLogin(user.TelegramID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if refreshed, err := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName); err == nil {
		user = refreshed
	}
	api.writeSession(w, user)
}

func (api *webAPI) handleSettings(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		InterfaceLanguage string  `json:"interface_language"`
		LearningLanguage  string  `json:"learning_language"`
		Level             string  `json:"level"`
		LearningFocus     *string `json:"learning_focus"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	if strings.TrimSpace(req.InterfaceLanguage) != "" {
		if err := api.bot.store.setInterfaceLanguage(user.TelegramID, req.InterfaceLanguage); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if strings.TrimSpace(req.LearningLanguage) != "" {
		if err := api.bot.store.setLearningLanguage(user.TelegramID, req.LearningLanguage); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if strings.TrimSpace(req.Level) != "" {
		if err := api.bot.store.setUserLevel(user.TelegramID, req.Level); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if req.LearningFocus != nil {
		if err := api.bot.store.setLearningFocus(user.TelegramID, *req.LearningFocus); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	refreshed, err := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.writeSession(w, refreshed)
}

func (api *webAPI) handleNavigationLayout(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req navigationLayout
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	next := user.NavigationLayout
	if len(req.FunctionRibbon) > 0 {
		next.FunctionRibbon = req.FunctionRibbon
	}
	if len(req.MobilePinned) > 0 {
		next.MobilePinned = req.MobilePinned
	}
	if len(req.MobileMore) > 0 {
		next.MobileMore = req.MobileMore
	}
	if len(req.MobileRail) > 0 {
		next.MobileRail = req.MobileRail
	}
	if err := api.bot.store.setNavigationLayout(user.TelegramID, next); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, err := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	api.writeSession(w, refreshed)
}

func (api *webAPI) handleDailyClaim(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		Date string `json:"date"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	date := strings.TrimSpace(req.Date)
	if date == "" {
		date = localDateForUser(user, time.Now().UTC())
	}
	const amount = 25
	claimed, err := api.bot.store.claimDailyBonus(user.TelegramID, date, amount)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, err := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	claimedAt := ""
	if !refreshed.DailyBonusLastClaimedAt.IsZero() {
		claimedAt = refreshed.DailyBonusLastClaimedAt.Format(time.RFC3339)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"claimed":    claimed,
		"xp":         amount,
		"claimed_at": claimedAt,
		"user":       api.userDTO(refreshed),
	})
}

func (api *webAPI) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	languageCode := strings.TrimSpace(r.URL.Query().Get("language"))
	if languageCode == "" || strings.EqualFold(languageCode, "all") {
		languageCode = user.LearningLanguage
	}
	languageCode = normalizeLearningLanguage(languageCode)
	overall := api.bot.store.leaderboard(10)
	languageEntries := api.bot.store.languageLeaderboard(languageCode, 10)
	writeJSON(w, http.StatusOK, map[string]any{
		"overall":          webLeaderboardDTO(overall),
		"language":         languageCode,
		"language_name":    learningLanguageByCode(languageCode).NativeName,
		"language_entries": webLanguageLeaderboardDTO(languageEntries),
		"languages":        webLanguageDTOs(learningLanguages),
	})
}

func (api *webAPI) handleTutorStart(w http.ResponseWriter, r *http.Request) {
	api.handleAITutorStart(w, r)
}

func (api *webAPI) requireAITutorPremium(w http.ResponseWriter, user userState) bool {
	if user.isPremium(time.Now()) {
		return true
	}
	writeAPIErrorCode(w, http.StatusPaymentRequired, "premium_required", "AI Tutor is available with Premium.")
	return false
}

func (api *webAPI) handleAITutorStart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.requireAITutorPremium(w, user) {
		return
	}
	if api.bot == nil || api.bot.store == nil {
		writeAPIError(w, http.StatusInternalServerError, "AI Tutor is not configured")
		return
	}
	result, err := api.bot.aiTutorEngine().Start(r.Context(), user, "web")
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, api.aiTutorDTO(result, user))
}

func (api *webAPI) handleAITutorSession(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.requireAITutorPremium(w, user) {
		return
	}
	sessionID := strings.TrimSpace(r.URL.Query().Get("session_id"))
	if r.Method == http.MethodPost {
		var req struct {
			SessionID string `json:"session_id"`
		}
		if !decodeJSONRequest(w, r, &req) {
			return
		}
		sessionID = strings.TrimSpace(req.SessionID)
	}
	result, err := api.loadAITutorSessionResult(user, sessionID)
	if err != nil {
		writeAPIError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, api.aiTutorDTO(result, user))
}

func (api *webAPI) handleAITutorCompleted(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.requireAITutorPremium(w, user) {
		return
	}
	if api.bot == nil || api.bot.store == nil {
		writeAPIError(w, http.StatusInternalServerError, "AI Tutor is not configured")
		return
	}
	limit := 100
	if raw := strings.TrimSpace(r.URL.Query().Get("limit")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	items, err := api.bot.store.completedAITutorLessons(user.TelegramID, limit)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"total": len(items),
	})
}

func (api *webAPI) handleAITutorRestart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if api.bot == nil || api.bot.store == nil {
		writeAPIError(w, http.StatusInternalServerError, "AI Tutor is not configured")
		return
	}
	var req struct {
		LessonID string `json:"lesson_id"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	lessonID := strings.TrimSpace(req.LessonID)
	if lessonID == "" {
		writeAPIError(w, http.StatusBadRequest, "lesson_id is required")
		return
	}
	allowed, err := api.bot.store.userCanRestartAITutorLesson(user.TelegramID, lessonID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !allowed {
		writeAPIError(w, http.StatusForbidden, "This AI Tutor lesson is not available in your history.")
		return
	}
	lesson, ok, err := api.bot.store.getAITutorLesson(lessonID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeAPIError(w, http.StatusNotFound, "AI Tutor lesson not found")
		return
	}
	result, err := api.bot.aiTutorEngine().createSessionForLesson(user, "web", lesson)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, api.aiTutorDTO(result, user))
}

func (api *webAPI) handleAITutorAnswer(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.requireAITutorPremium(w, user) {
		return
	}
	var req struct {
		SessionID string `json:"session_id"`
		Text      string `json:"text"`
		Choice    string `json:"choice"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	result, err := api.bot.aiTutorEngine().Submit(r.Context(), user, req.SessionID, aiTutorSubmitInput{Text: req.Text, Choice: req.Choice})
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, api.aiTutorDTO(result, user))
}

func (api *webAPI) handleAITutorWordReport(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.requireAITutorPremium(w, user) {
		return
	}
	if api.bot == nil || api.bot.store == nil {
		writeAPIError(w, http.StatusInternalServerError, "AI Tutor is not configured")
		return
	}
	var req struct {
		SessionID           string `json:"session_id"`
		Stage               string `json:"stage"`
		ProposedWord        string `json:"proposed_word"`
		ProposedTranslation string `json:"proposed_translation"`
		Comment             string `json:"comment"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 4096)
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "bad JSON request")
		return
	}
	sessionID := strings.TrimSpace(req.SessionID)
	stage := strings.TrimSpace(req.Stage)
	if sessionID == "" || stage == "" {
		writeAPIError(w, http.StatusBadRequest, "session_id and stage are required")
		return
	}
	session, ok, err := api.bot.store.getAITutorSession(sessionID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok || session.TelegramID != user.TelegramID {
		writeAPIError(w, http.StatusNotFound, "AI Tutor session not found")
		return
	}
	if strings.TrimSpace(session.CurrentStage) != stage {
		writeAPIError(w, http.StatusConflict, "This word is no longer active.")
		return
	}
	wordIndex, ok := aiTutorStageIndex(stage, "word_learn_")
	if !ok {
		writeAPIError(w, http.StatusBadRequest, "Current AI Tutor step is not a word-learning step.")
		return
	}
	lesson, ok, err := api.bot.store.getAITutorLesson(session.LessonID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !ok {
		writeAPIError(w, http.StatusNotFound, "AI Tutor lesson not found")
		return
	}
	if wordIndex < 0 || wordIndex >= len(lesson.Payload.Words) {
		writeAPIError(w, http.StatusBadRequest, "Word stage is out of range.")
		return
	}
	proposedWord, ok := validateAITutorWordReportField(w, req.ProposedWord, 80, "proposed_word")
	if !ok {
		return
	}
	proposedTranslation, ok := validateAITutorWordReportField(w, req.ProposedTranslation, 160, "proposed_translation")
	if !ok {
		return
	}
	comment := strings.TrimSpace(req.Comment)
	if len([]rune(comment)) > 500 {
		writeAPIError(w, http.StatusBadRequest, "comment is too long")
		return
	}
	now := time.Now().UTC()
	userBurst, _, err := api.bot.store.countAITutorWordReports(user.TelegramID, session.ID, now.Add(-10*time.Minute))
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	userDaily, sessionTotal, err := api.bot.store.countAITutorWordReports(user.TelegramID, session.ID, now.Add(-24*time.Hour))
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if userBurst >= 3 || userDaily >= 10 || sessionTotal >= 2 {
		writeAPIError(w, http.StatusTooManyRequests, "Too many word reports. Try again later.")
		return
	}

	word := lesson.Payload.Words[wordIndex]
	report, _, err := api.bot.store.createAITutorWordReport(aiTutorWordReportRecord{
		ID:                  aiTutorNewID("aitwr"),
		TelegramID:          user.TelegramID,
		SessionID:           session.ID,
		LessonID:            lesson.ID,
		Stage:               stage,
		WordIndex:           wordIndex,
		OriginalWord:        strings.TrimSpace(word.Target),
		OriginalTranslation: strings.TrimSpace(word.InterfaceTranslation),
		ProposedWord:        proposedWord,
		ProposedTranslation: proposedTranslation,
		Comment:             comment,
		Status:              aiTutorWordReportPending,
	})
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	globalCount, err := api.bot.store.countAITutorWordReportsSince(now.Add(-10 * time.Minute))
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if globalCount <= 30 {
		api.sendAITutorWordReportNotification(r.Context(), report, user)
	}

	nextStage := aiTutorNextStage(session.CurrentStage)
	status := aiTutorSessionActive
	completedAt := ""
	if nextStage == aiTutorStageComplete {
		status = aiTutorSessionComplete
		completedAt = formatDBTime(now)
	}
	if err := api.bot.store.updateAITutorSessionStage(session.ID, nextStage, status, completedAt); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	result, err := api.loadAITutorSessionResult(user, session.ID)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	payload := api.aiTutorDTO(result, user)
	payload["word_report"] = report
	writeJSON(w, http.StatusOK, payload)
}

func validateAITutorWordReportField(w http.ResponseWriter, value string, maxRunes int, field string) (string, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		writeAPIError(w, http.StatusBadRequest, field+" is required")
		return "", false
	}
	if len([]rune(value)) > maxRunes {
		writeAPIError(w, http.StatusBadRequest, field+" is too long")
		return "", false
	}
	return value, true
}

func (api *webAPI) sendAITutorWordReportNotification(ctx context.Context, report aiTutorWordReportRecord, user userState) {
	if api == nil || api.bot == nil || api.bot.telegram == nil {
		return
	}
	text := aiTutorWordReportTelegramText(report, user)
	keyboard := aiTutorWordReportTelegramKeyboard(report.ID)
	for _, recipient := range api.telegramOpsRecipients() {
		messageID, err := api.bot.telegram.sendInlineMessageToChat(ctx, recipient.telegramChatIDValue(), text, keyboard)
		if err != nil {
			log.Printf("send ai tutor word report to %s: %v", recipient.ChatID, err)
			continue
		}
		if messageID != 0 {
			if err := api.bot.store.setAITutorWordReportAdminMessage(report.ID, fmt.Sprint(recipient.telegramChatIDValue()), messageID); err != nil {
				log.Printf("store ai tutor word report admin message %s: %v", report.ID, err)
			}
		}
	}
}

func aiTutorWordReportTelegramText(report aiTutorWordReportRecord, user userState) string {
	lines := []string{
		"AI Tutor word report",
		"Report: " + strings.TrimSpace(report.ID),
		"User: " + strings.TrimSpace(user.FirstName) + " (" + strconv.FormatInt(user.TelegramID, 10) + ")",
		"Session: " + strings.TrimSpace(report.SessionID),
		"Stage: " + strings.TrimSpace(report.Stage),
		"Original: " + strings.TrimSpace(report.OriginalWord) + " - " + strings.TrimSpace(report.OriginalTranslation),
		"User suggests: " + strings.TrimSpace(report.ProposedWord) + " - " + strings.TrimSpace(report.ProposedTranslation),
	}
	if strings.TrimSpace(report.Comment) != "" {
		lines = append(lines, "Comment: "+strings.TrimSpace(report.Comment))
	}
	lines = append(lines, "", "Accept/Fix will replace this word in SQLite: ai_tutor_lessons payload, vocabulary_words, vocabulary_translations, vocabulary_ai_words, and vocabulary_ai_translations.")
	lines = append(lines, "Reply with buttons: accept, reject, or fix.")
	return strings.Join(lines, "\n")
}

func aiTutorWordReportTelegramKeyboard(reportID string) map[string]any {
	reportID = strings.TrimSpace(reportID)
	return map[string]any{
		"inline_keyboard": [][]map[string]any{
			{
				{"text": "Принять", "callback_data": "ait_word_report|" + reportID + "|accept"},
				{"text": "Отклонить", "callback_data": "ait_word_report|" + reportID + "|reject"},
				{"text": "Исправить", "callback_data": "ait_word_report|" + reportID + "|fix"},
			},
		},
	}
}

func (api *webAPI) handleAITutorReview(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.requireAITutorPremium(w, user) {
		return
	}
	var req struct {
		SessionID string `json:"session_id"`
		Choice    string `json:"choice"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	result, err := api.bot.aiTutorEngine().Submit(r.Context(), user, req.SessionID, aiTutorSubmitInput{Choice: req.Choice})
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, api.aiTutorDTO(result, user))
}

func (api *webAPI) handleAITutorFinish(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.requireAITutorPremium(w, user) {
		return
	}
	var req struct {
		SessionID string `json:"session_id"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	result, err := api.bot.aiTutorEngine().Submit(r.Context(), user, req.SessionID, aiTutorSubmitInput{Choice: "no_review"})
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, api.aiTutorDTO(result, user))
}

func (api *webAPI) loadAITutorSessionResult(user userState, sessionID string) (aiTutorResult, error) {
	if strings.TrimSpace(sessionID) == "" {
		return aiTutorResult{}, errors.New("AI Tutor session not found")
	}
	session, ok, err := api.bot.store.getAITutorSession(sessionID)
	if err != nil {
		return aiTutorResult{}, err
	}
	if !ok || session.TelegramID != user.TelegramID {
		return aiTutorResult{}, errors.New("AI Tutor session not found")
	}
	lesson, ok, err := api.bot.store.getAITutorLesson(session.LessonID)
	if err != nil {
		return aiTutorResult{}, err
	}
	if !ok {
		return aiTutorResult{}, errors.New("AI Tutor lesson not found")
	}
	return aiTutorResult{Session: session, Lesson: lesson, NextStep: aiTutorBuildStep(lesson.Payload, session.CurrentStage)}, nil
}

func (api *webAPI) aiTutorDTO(result aiTutorResult, user userState) map[string]any {
	rewardXP := 0
	if result.Session.Status == aiTutorSessionComplete || result.NextStep.Stage == aiTutorStageComplete {
		rewardXP = aiTutorCompletionXP
	}
	if api != nil && api.bot != nil && api.bot.store != nil && user.TelegramID != 0 {
		if refreshed, err := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName); err == nil {
			user = refreshed
		}
	}
	payload := map[string]any{
		"session":       result.Session,
		"lesson":        result.Lesson.Payload,
		"lesson_status": result.Lesson.Status,
		"current_stage": result.Session.CurrentStage,
		"next_step":     result.NextStep,
		"feedback":      result.Feedback,
		"user":          api.userDTO(user),
	}
	if rewardXP > 0 {
		payload["xp"] = rewardXP
		payload["reward_xp"] = rewardXP
		payload["reward_title"] = "AI Tutor"
	}
	return payload
}

func (api *webAPI) handleLessonStart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !canUseLessons(user) {
		writeAPIError(w, http.StatusTooManyRequests, api.bot.limitReachedText(systemUI(user).LessonKind, user))
		return
	}
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	task, err := api.bot.openrouter.complete(r.Context(), lessonPrompt(language, interfaceLanguage, user.Level, user.LearningFocus, user.LessonCount, user.LessonHistory), 0.7, 500)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "Не получилось получить задание: "+err.Error())
		return
	}
	if err := api.bot.store.saveLesson(user.TelegramID, task); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"lesson":      task,
		"instruction": systemUI(user).LessonAnswerInstruction,
		"user":        api.userDTO(refreshed),
	})
}

func (api *webAPI) handleLessonAnswer(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if strings.TrimSpace(user.LastLessonPrompt) == "" {
		writeAPIError(w, http.StatusConflict, "No active lesson. Start a new lesson first.")
		return
	}
	text, fromVoice, pronunciation, ok := api.readLearningAnswerInput(w, r, user, "")
	if !ok {
		return
	}
	if text == "" {
		writeAPIError(w, http.StatusBadRequest, "empty answer")
		return
	}
	awardLessonXP := strings.TrimSpace(user.Mode) == "lesson"
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	raw, err := api.bot.openrouter.complete(r.Context(), feedbackPrompt(language, interfaceLanguage, user.Level, user.LastLessonPrompt, text), 0.3, 900)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "Не получилось проверить ответ: "+err.Error())
		return
	}
	feedback, mistakesJSON := splitFeedbackAndMistakes(raw)
	entries := parseMistakesJSON(mistakesJSON, time.Now().UTC())
	if len(entries) == 0 && strings.TrimSpace(mistakesJSON) == "" {
		entries = api.extractWebMistakes(r.Context(), user, language, interfaceLanguage, "lesson answer", user.LastLessonPrompt, text, feedback)
	}
	if err := api.bot.store.addMistakes(user.TelegramID, user.LearningLanguage, entries); err != nil {
		log.Printf("addMistakes (web lesson) for %d: %v", user.TelegramID, err)
	}
	if err := api.bot.store.completeLesson(user.TelegramID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if awardLessonXP {
		if err := api.bot.store.addXP(user.TelegramID, 20); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	promotedTo, err := api.bot.maybePromoteLearningLevelForWeb(r.Context(), user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"feedback":              feedback,
		"transcript":            text,
		"from_voice":            fromVoice,
		"pronunciation":         pronunciationAssessmentDTO(pronunciation),
		"mistakes":              entries,
		"correction_audio_text": lessonCorrectionPronunciationText(feedback, entries),
		"promoted_to":           promotedTo,
		"user":                  api.userDTO(refreshed),
	})
}

func (api *webAPI) handlePractice(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !canUsePractice(user) {
		_ = api.bot.store.setMode(user.TelegramID, "idle")
		writeAPIError(w, http.StatusTooManyRequests, api.bot.limitReachedText(systemUI(user).PracticeKind, user))
		return
	}
	input, ok := api.readPracticeInput(w, r, user)
	if !ok {
		return
	}
	if input.text == "" && len(input.imageBytes) == 0 {
		writeAPIError(w, http.StatusBadRequest, "empty message")
		return
	}
	language := userLearningLanguage(user)
	interfaceLanguage := userInterfaceLanguage(user)
	imageContext := ""
	if len(input.imageBytes) > 0 {
		if api.bot.openrouter == nil {
			writeAPIError(w, http.StatusServiceUnavailable, "Image practice is not configured.")
			return
		}
		copy := ui(user)
		if !user.isPremium(time.Now()) {
			writeAPIError(w, http.StatusPaymentRequired, copy.Tool.ImagePremiumRequired)
			return
		}
		mimeType := imageUploadMime(input.imageName, input.imageBytes)
		if !strings.HasPrefix(mimeType, "image/") {
			writeAPIError(w, http.StatusBadRequest, "Upload an image for practice.")
			return
		}
		var err error
		imageContext, err = api.bot.openrouter.describePracticeImage(r.Context(), input.imageBytes, mimeType, language, interfaceLanguage, input.text)
		if err != nil {
			writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(copy.Tool.ImageReadFailed, err.Error()))
			return
		}
	}
	practiceMessage := practiceImageMessage(input.text, imageContext)
	practiceHistory := trimPracticeHistory(append(user.PracticeHistory, practiceMessage), practiceMemoryLimit)
	raw, err := api.bot.openrouter.complete(r.Context(), practicePrompt(language, interfaceLanguage, user.Level, practiceHistory, user.LearningFocus, user.PracticeCount), 0.6, 900)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "Не получилось ответить: "+err.Error())
		return
	}
	reply, mistakesJSON := splitFeedbackAndMistakes(raw)
	entries := parseMistakesJSON(mistakesJSON, time.Now().UTC())
	if len(entries) == 0 && strings.TrimSpace(mistakesJSON) == "" {
		entries = api.extractWebMistakes(r.Context(), user, language, interfaceLanguage, "practice message", practiceMessage, input.text, reply)
	}
	if err := api.bot.store.addMistakes(user.TelegramID, user.LearningLanguage, entries); err != nil {
		log.Printf("addMistakes (web practice) for %d: %v", user.TelegramID, err)
	}
	if err := api.bot.store.incrementPractice(user.TelegramID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := api.bot.store.savePracticeHistory(user.TelegramID, practiceHistory); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	promotedTo, err := api.bot.maybePromoteLearningLevelForWeb(r.Context(), user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"reply":                 reply,
		"transcript":            input.text,
		"from_voice":            input.fromVoice,
		"image_attached":        len(input.imageBytes) > 0,
		"pronunciation":         pronunciationAssessmentDTO(input.pronunciation),
		"mistakes":              entries,
		"correction_audio_text": practiceCorrectionPronunciationText(reply, entries),
		"question_audio_text":   practiceQuestionPronunciationText(reply),
		"promoted_to":           promotedTo,
		"user":                  api.userDTO(refreshed),
	})
}

func (api *webAPI) extractWebMistakes(ctx context.Context, user userState, language learningLanguage, interfaceLanguage learningLanguage, source string, task string, learnerAnswer string, feedback string) []mistakeEntry {
	if api == nil || api.bot == nil || api.bot.openrouter == nil {
		return nil
	}
	raw, err := api.bot.openrouter.complete(ctx, mistakeExtractionPrompt(language, interfaceLanguage, source, task, learnerAnswer, feedback), 0.1, 500)
	if err != nil {
		log.Printf("extract web mistakes for %d: %v", user.TelegramID, err)
		return nil
	}
	return parseMistakesJSON(raw, time.Now().UTC())
}

func (api *webAPI) handleLevelTestStart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if len(levelAssessmentQuestionsForUser(user)) == 0 {
		writeAPIError(w, http.StatusConflict, "Level assessment is unavailable for this language.")
		return
	}
	if err := api.bot.store.setMode(user.TelegramID, levelTestMode(0, 0)); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"complete": false,
		"question": api.levelQuestionDTO(user, 0, 0),
	})
}

func (api *webAPI) handleLevelTestAnswer(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		Index  int `json:"index"`
		Answer int `json:"answer"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	expectedIndex, score, ok := parseLevelTestMode(user.Mode)
	if !ok {
		writeAPIError(w, http.StatusConflict, "Level test is not active.")
		return
	}
	questions := levelAssessmentQuestionsForUser(user)
	if req.Index != expectedIndex || req.Index < 0 || req.Index >= len(questions) || req.Answer < -1 || req.Answer > 3 {
		writeAPIError(w, http.StatusBadRequest, "stale or invalid answer")
		return
	}
	if req.Answer == questions[req.Index].CorrectIndex {
		score += questions[req.Index].Points
	}
	nextIndex := req.Index + 1
	if nextIndex < len(questions) {
		if err := api.bot.store.setMode(user.TelegramID, levelTestMode(nextIndex, score)); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"complete": false,
			"question": api.levelQuestionDTO(user, nextIndex, score),
		})
		return
	}
	level := levelFromAssessmentScore(score, levelAssessmentMaxScoreForUser(user))
	if err := api.bot.store.setUserLevel(user.TelegramID, level); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := api.bot.store.setMode(user.TelegramID, "idle"); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"complete":  true,
		"level":     level,
		"score":     score,
		"max_score": levelAssessmentMaxScoreForUser(user),
		"user":      api.userDTO(refreshed),
	})
}

func (api *webAPI) handleWordNext(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	word, ok := nextUnlearnedWord(user)
	if !ok {
		writeJSON(w, http.StatusOK, map[string]any{"empty": true})
		return
	}
	if err := api.bot.store.setMode(user.TelegramID, modeWebWordPrefix+word.ID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	options := wordOptions(word)
	dto := make([]webWordOptionDTO, 0, len(options))
	for _, option := range options {
		dto = append(dto, webWordOptionDTO{ID: option.ID, Text: option.English})
	}
	language := userLearningLanguage(user)
	writeJSON(w, http.StatusOK, map[string]any{
		"empty":       false,
		"prompt":      wordTranslation(word, user.InterfaceLanguage),
		"context":     api.bot.vocabularyHint(r.Context(), user, word, "learn word"),
		"direction":   learningDirectionForUI(user, language),
		"options":     dto,
		"instruction": ui(user).ChooseAnswer,
	})
}

func (api *webAPI) handleWordAnswer(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		AnswerID string `json:"answer_id"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	correctID := strings.TrimPrefix(user.Mode, modeWebWordPrefix)
	if correctID == user.Mode || correctID == "" {
		writeAPIError(w, http.StatusConflict, "No active word question.")
		return
	}
	word, ok := findVocabWord(correctID)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "Word not found.")
		return
	}
	if req.AnswerID != correctID {
		writeJSON(w, http.StatusOK, map[string]any{
			"correct": false,
			"message": ui(user).TryAgain,
			"prompt":  wordTranslation(word, user.InterfaceLanguage),
			"context": api.bot.vocabularyHint(r.Context(), user, word, "learn word"),
		})
		return
	}
	learned, total, err := api.bot.store.addLearnedWord(user.TelegramID, word)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = api.bot.store.setMode(user.TelegramID, "idle")
	promotedTo, err := api.bot.maybePromoteLearningLevelForWeb(r.Context(), user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"correct":     true,
		"learned":     learned,
		"total":       total,
		"word_id":     word.ID,
		"word":        word.English,
		"translation": wordTranslation(word, user.InterfaceLanguage),
		"context":     api.bot.vocabularyHint(r.Context(), user, word, "learn word"),
		"example":     api.bot.vocabularyExample(r.Context(), user, word, "learn word"),
		"promoted_to": promotedTo,
		"user":        api.userDTO(refreshed),
	})
}

func (api *webAPI) handleWordPronunciation(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if api.bot.openrouter == nil || strings.TrimSpace(api.cfg.OpenRouterTTSModel) == "" {
		writeAPIError(w, http.StatusServiceUnavailable, "Озвучка пока не настроена")
		return
	}
	var req struct {
		WordID string `json:"word_id"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	wordID := strings.TrimSpace(req.WordID)
	if wordID == "" {
		writeAPIError(w, http.StatusBadRequest, "word_id is required")
		return
	}
	word, ok := findVocabWord(wordID)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "Word not found.")
		return
	}
	if normalizeLearningLanguage(word.Language) != normalizeLearningLanguage(user.LearningLanguage) {
		writeAPIError(w, http.StatusBadRequest, "Word language does not match current learning language.")
		return
	}
	if !userLearnedWord(user, wordID) {
		writeAPIError(w, http.StatusForbidden, "Озвучка доступна после правильного ответа на слово")
		return
	}
	cacheKey := strings.Join([]string{api.cfg.OpenRouterTTSModel, api.cfg.OpenRouterTTSVoice, word.ID}, "|")
	if audio, ok := api.cachedPronunciation(cacheKey); ok {
		writePronunciationAudio(w, word, audio)
		return
	}
	audio, err := api.bot.openrouter.synthesizeSpeech(r.Context(), api.cfg.OpenRouterTTSModel, api.cfg.OpenRouterTTSVoice, word.English)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "Не получилось озвучить слово: "+err.Error())
		return
	}
	api.cachePronunciation(cacheKey, audio)
	writePronunciationAudio(w, word, audio)
}

func (api *webAPI) handleVocabulary(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	const pageSize = 10
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 0 {
		page = 0
	}
	words := learnedWordsForLanguage(user)
	total := len(words)
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}
	if page >= totalPages {
		page = totalPages - 1
	}
	start := page * pageSize
	end := start + pageSize
	if end > total {
		end = total
	}
	pageWords := words[start:end]
	lookup := vocabularyLookupForLearnedWords(pageWords)
	items := []map[string]any{}
	for _, word := range pageWords {
		items = append(items, map[string]any{
			"id":                     word.ID,
			"word":                   word.English,
			"translation":            learnedWordTranslationWithLookup(word, user.InterfaceLanguage, lookup),
			"context":                learnedWordContextWithLookup(word, user.InterfaceLanguage, lookup),
			"review_correct_count":   word.ReviewCorrectCount,
			"spelling_correct_count": word.SpellingCorrectCount,
			"training_count":         learnedWordTrainingCount(word),
		})
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"items":       items,
		"page":        page,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": totalPages,
		"user":        api.userDTO(refreshed),
	})
}

func (api *webAPI) handlePhrasebook(w http.ResponseWriter, r *http.Request) {
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	switch r.Method {
	case http.MethodPost:
		var req phrasebookEntry
		if !decodeJSONRequest(w, r, &req) {
			return
		}
		if strings.TrimSpace(req.Language) == "" {
			req.Language = user.LearningLanguage
		}
		if strings.TrimSpace(req.Translation) == "" && strings.TrimSpace(req.Phrase) != "" && api.bot != nil && api.bot.openrouter != nil {
			model := strings.TrimSpace(api.cfg.OpenRouterTranslatorModel)
			if model == "" {
				model = strings.TrimSpace(api.bot.cfg.OpenRouterTranslatorModel)
			}
			translation, err := api.bot.openrouter.completeWithModel(
				r.Context(),
				model,
				translationToolPrompt(req.Phrase, normalizeTranslatorSource(req.Language), defaultTranslatorTarget(user), userInterfaceLanguage(user)),
				0.1,
				500,
			)
			if err == nil {
				req.Translation = strings.TrimSpace(translation)
			} else {
				log.Printf("phrasebook auto-translate: %v", err)
			}
		}
		items, err := api.bot.store.addPhrasebookEntry(user.TelegramID, req)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
	case http.MethodDelete:
		id := strings.TrimSpace(r.URL.Query().Get("id"))
		if id == "" {
			writeAPIError(w, http.StatusBadRequest, "missing phrasebook id")
			return
		}
		items, err := api.bot.store.removePhrasebookEntry(user.TelegramID, id)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (api *webAPI) handleBugReport(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, webMaxImageUploadBytes*5+webMaxMultipartOverhead)
	if err := r.ParseMultipartForm(webMaxImageUploadBytes); err != nil {
		writeAPIError(w, http.StatusBadRequest, "bad bug report form")
		return
	}
	message := strings.TrimSpace(r.FormValue("message"))
	files := bugReportUploadFiles(r)
	if len([]rune(message)) < 8 && len(files) == 0 {
		writeAPIError(w, http.StatusBadRequest, "Describe the problem in more detail.")
		return
	}
	if message == "" && len(files) > 0 {
		message = "Image attachment only."
	}
	view := strings.TrimSpace(r.FormValue("view"))
	userAgent := strings.TrimSpace(r.FormValue("user_agent"))
	reportID := "bug-" + time.Now().UTC().Format("20060102-150405.000000000")
	caption := bugReportCaption(reportID, user, view, message, userAgent)
	if err := appendBugReportFile(reportID, user, view, message, userAgent); err != nil {
		log.Printf("write bug report: %v", err)
	}
	photoCount := 0
	for index, header := range files {
		if header.Size > webMaxImageUploadBytes {
			continue
		}
		file, err := header.Open()
		if err != nil {
			continue
		}
		data, readErr := io.ReadAll(io.LimitReader(file, webMaxImageUploadBytes+1))
		_ = file.Close()
		if readErr != nil || len(data) == 0 || len(data) > webMaxImageUploadBytes {
			continue
		}
		if !strings.HasPrefix(imageUploadMime(header.Filename, data), "image/") {
			continue
		}
		photoCount++
		if api.bot != nil && api.bot.telegram != nil {
			path, err := writeBugReportScreenshot(reportID, index, header.Filename, data)
			if err != nil {
				log.Printf("save bug screenshot: %v", err)
				continue
			}
			photoCaption := caption
			if len(files) > 1 {
				photoCaption += fmt.Sprintf("\nScreenshot: %d/%d", index+1, len(files))
			}
			api.sendBugReportPhoto(r.Context(), path, photoCaption)
		}
	}
	if photoCount == 0 && api.bot != nil && api.bot.telegram != nil {
		api.sendBugReportText(r.Context(), caption)
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "id": reportID, "screenshots": photoCount})
}

func (api *webAPI) telegramOpsRecipients() []telegramOpsRecipient {
	if api != nil && len(api.cfg.TelegramOpsRecipients) > 0 {
		return api.cfg.telegramOpsRecipients()
	}
	if api != nil && api.bot != nil && len(api.bot.cfg.TelegramOpsRecipients) > 0 {
		return api.bot.cfg.telegramOpsRecipients()
	}
	return (config{}).telegramOpsRecipients()
}

func (api *webAPI) sendBugReportText(ctx context.Context, caption string) {
	if api == nil || api.bot == nil || api.bot.telegram == nil {
		return
	}
	for _, recipient := range api.telegramOpsRecipients() {
		if err := api.bot.telegram.sendMessageToChat(ctx, recipient.telegramChatIDValue(), caption); err != nil {
			log.Printf("send bug report to %s: %v", recipient.ChatID, err)
		}
	}
}

func (api *webAPI) sendBugReportPhoto(ctx context.Context, path string, caption string) {
	if api == nil || api.bot == nil || api.bot.telegram == nil {
		return
	}
	for _, recipient := range api.telegramOpsRecipients() {
		if err := api.bot.telegram.sendPhotoFileToChat(ctx, recipient.telegramChatIDValue(), path, caption); err != nil {
			log.Printf("send bug screenshot to %s: %v", recipient.ChatID, err)
		}
	}
}

func (api *webAPI) handleWordGameNext(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	ids := reviewWordIDs(user)
	if len(ids) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{
			"empty":   true,
			"message": ui(user).LearnWords,
		})
		return
	}
	word, ok := findVocabWord(ids[(user.WordGameCount+len(ids))%len(ids)])
	if !ok {
		writeAPIError(w, http.StatusNotFound, "Не получилось найти слово для игры.")
		return
	}
	if err := api.bot.store.setMode(user.TelegramID, modeWebWordGamePrefix+word.ID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	language := userLearningLanguage(user)
	writeJSON(w, http.StatusOK, map[string]any{
		"empty":       false,
		"prompt":      wordTranslation(word, user.InterfaceLanguage),
		"context":     api.bot.vocabularyHint(r.Context(), user, word, "review word"),
		"direction":   learningDirectionForUI(user, language),
		"options":     webWordOptionsDTO(learnedWordOptions(word, ids)),
		"instruction": ui(user).ChooseAnswer,
	})
}

func (api *webAPI) handleWordGameAnswer(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		AnswerID string `json:"answer_id"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	correctID := strings.TrimPrefix(user.Mode, modeWebWordGamePrefix)
	if correctID == user.Mode || correctID == "" {
		writeAPIError(w, http.StatusConflict, "No active review question.")
		return
	}
	word, ok := findVocabWord(correctID)
	if !ok {
		writeAPIError(w, http.StatusNotFound, "Word not found.")
		return
	}
	ids := reviewWordIDs(user)
	if req.AnswerID != correctID {
		writeJSON(w, http.StatusOK, map[string]any{
			"correct": false,
			"message": ui(user).TryAgain,
			"prompt":  wordTranslation(word, user.InterfaceLanguage),
			"context": api.bot.vocabularyHint(r.Context(), user, word, "review word"),
			"options": webWordOptionsDTO(learnedWordOptions(word, ids)),
		})
		return
	}
	masteredNow := wordWillBeMasteredAfterCorrect(user, word.ID)
	if err := api.bot.store.markWordGameCorrect(user.TelegramID, word.ID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = api.bot.store.setMode(user.TelegramID, "idle")
	promotedTo, err := api.bot.maybePromoteLearningLevelForWeb(r.Context(), user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"correct":     true,
		"word_id":     word.ID,
		"word":        word.English,
		"translation": wordTranslation(word, user.InterfaceLanguage),
		"context":     api.bot.vocabularyHint(r.Context(), user, word, "review word"),
		"example":     api.bot.vocabularyExample(r.Context(), user, word, "review word"),
		"mastered":    masteredNow,
		"xp":          10,
		"promoted_to": promotedTo,
		"user":        api.userDTO(refreshed),
	})
}

func parseWebSpellingMode(mode string) (string, int, bool) {
	value := strings.TrimPrefix(mode, modeWebSpellingPrefix)
	if value == mode || value == "" {
		return "", 0, false
	}
	wordID := value
	attempts := 0
	if cut := strings.LastIndex(value, "|"); cut >= 0 {
		wordID = value[:cut]
		if parsed, err := strconv.Atoi(value[cut+1:]); err == nil && parsed > 0 {
			attempts = parsed
		}
	}
	if wordID == "" {
		return "", 0, false
	}
	return wordID, attempts, true
}

func webSpellingMode(wordID string, attempts int) string {
	if attempts <= 0 {
		return modeWebSpellingPrefix + wordID
	}
	return modeWebSpellingPrefix + wordID + "|" + strconv.Itoa(attempts)
}

func selectWebSpellingWordID(user userState, ids []string) string {
	if len(ids) == 0 {
		return ""
	}
	if skipped := strings.TrimPrefix(user.Mode, modeWebSpellingSkipPrefix); skipped != user.Mode && skipped != "" {
		for index, id := range ids {
			if id == skipped {
				return ids[(index+1)%len(ids)]
			}
		}
	}
	return ids[(user.WordGameCount+user.XP+len(ids))%len(ids)]
}

func (api *webAPI) handleSpellingStart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	ids := spellingWordIDs(user)
	if len(ids) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{
			"empty":   true,
			"message": ui(user).LearnWords,
		})
		return
	}
	word, ok := findVocabWord(selectWebSpellingWordID(user, ids))
	if !ok {
		writeAPIError(w, http.StatusNotFound, "Не получилось найти слово для правописания.")
		return
	}
	if err := api.bot.store.setMode(user.TelegramID, webSpellingMode(word.ID, 0)); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	language := userLearningLanguage(user)
	writeJSON(w, http.StatusOK, map[string]any{
		"empty":     false,
		"word_id":   word.ID,
		"prompt":    wordTranslation(word, user.InterfaceLanguage),
		"context":   api.bot.vocabularyHint(r.Context(), user, word, "spelling practice"),
		"direction": learningDirectionForUI(user, language),
	})
}

func (api *webAPI) handleSpellingAnswer(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		Text   string `json:"text"`
		GiveUp bool   `json:"give_up"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	answer := strings.TrimSpace(req.Text)
	if answer == "" && !req.GiveUp {
		writeAPIError(w, http.StatusBadRequest, "empty answer")
		return
	}
	wordID, attempts, ok := parseWebSpellingMode(user.Mode)
	if !ok {
		writeAPIError(w, http.StatusConflict, "No active spelling question.")
		return
	}
	word, ok := learnedWordByID(user, wordID)
	if !ok {
		_ = api.bot.store.setMode(user.TelegramID, "idle")
		writeAPIError(w, http.StatusNotFound, "Не нашёл слово для проверки. Откройте правописание заново.")
		return
	}
	if req.GiveUp {
		entry := spellingMistakeEntry(user, word, time.Now().UTC())
		if err := api.bot.store.addMistakes(user.TelegramID, user.LearningLanguage, []mistakeEntry{entry}); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := api.bot.store.setMode(user.TelegramID, modeWebSpellingSkipPrefix+word.ID); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"correct":        false,
			"gave_up":        true,
			"word_id":        word.ID,
			"word":           word.English,
			"correct_answer": word.English,
			"translation":    learnedWordTranslation(word, user.InterfaceLanguage),
			"context":        learnedWordContext(word, user.InterfaceLanguage),
			"example": api.bot.vocabularyExample(r.Context(), user, vocabWord{
				ID:       word.ID,
				Language: word.Language,
				Russian:  word.Russian,
				English:  word.English,
				Context:  word.Context,
			}, "spelling practice"),
			"attempts": attempts,
			"message":  "Правильный ответ показан без начисления очков. Слово останется в правописании.",
		})
		return
	}
	if normalizeAnswer(answer) != normalizeAnswer(word.English) {
		attempts++
		if err := api.bot.store.setMode(user.TelegramID, webSpellingMode(word.ID, attempts)); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		language := userLearningLanguage(user)
		writeJSON(w, http.StatusOK, map[string]any{
			"correct":     false,
			"message":     ui(user).TryAgain,
			"prompt":      learnedWordTranslation(word, user.InterfaceLanguage),
			"context":     learnedWordContext(word, user.InterfaceLanguage),
			"direction":   learningDirectionForUI(user, language),
			"hint":        spellingHint(word.English),
			"attempts":    attempts,
			"can_give_up": attempts >= 3,
		})
		return
	}
	if err := api.bot.store.addXP(user.TelegramID, 7); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	masteredNow := wordWillBeMasteredAfterCorrect(user, word.ID)
	if err := api.bot.store.markSpellingCorrect(user.TelegramID, word.ID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := api.bot.store.setMode(user.TelegramID, "idle"); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	promotedTo, err := api.bot.maybePromoteLearningLevelForWeb(r.Context(), user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"correct":        true,
		"word_id":        word.ID,
		"word":           word.English,
		"correct_answer": word.English,
		"translation":    learnedWordTranslation(word, user.InterfaceLanguage),
		"context":        learnedWordContext(word, user.InterfaceLanguage),
		"example": api.bot.vocabularyExample(r.Context(), user, vocabWord{
			ID:       word.ID,
			Language: word.Language,
			Russian:  word.Russian,
			English:  word.English,
			Context:  word.Context,
		}, "spelling practice"),
		"mastered":    masteredNow,
		"xp":          7,
		"promoted_to": promotedTo,
		"user":        api.userDTO(refreshed),
	})
}

func (api *webAPI) handleMistakes(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	mistakes := mistakesForLanguage(user)
	page := 0
	if raw := strings.TrimSpace(r.URL.Query().Get("page")); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil {
			page = parsed
		}
	}
	totalPages := mistakeTotalPages(len(mistakes))
	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}
	start := page * mistakesPageSize
	if start > len(mistakes) {
		start = len(mistakes)
	}
	end := start + mistakesPageSize
	if end > len(mistakes) {
		end = len(mistakes)
	}
	items := webMistakeDTOsNewestFirst(mistakes)
	pageItems := items[start:end]
	writeJSON(w, http.StatusOK, map[string]any{
		"empty":       len(mistakes) == 0,
		"items":       items,
		"page_items":  pageItems,
		"page":        page,
		"total_pages": totalPages,
		"total":       len(mistakes),
		"user":        api.userDTO(user),
	})
}

func (api *webAPI) handleMistakePracticeStart(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		Index *int `json:"index"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	mistakes := mistakesForLanguage(user)
	if len(mistakes) == 0 {
		writeJSON(w, http.StatusOK, map[string]any{
			"empty":   true,
			"message": ui(user).Practice,
		})
		return
	}
	index := len(mistakes) - 1
	if req.Index != nil {
		index = *req.Index
	}
	if index < 0 || index >= len(mistakes) {
		writeAPIError(w, http.StatusBadRequest, "Mistake index is out of range.")
		return
	}
	if err := api.bot.store.setMode(user.TelegramID, modeWebMistakePrefix+itoa(index)); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"empty":       false,
		"mistake":     webMistakeDTOFromEntry(mistakes[index], index),
		"instruction": systemUI(user).LessonAnswerInstruction,
	})
}

func (api *webAPI) handleMistakePracticeAnswer(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		Text string `json:"text"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	answer := strings.TrimSpace(req.Text)
	if answer == "" {
		writeAPIError(w, http.StatusBadRequest, "empty answer")
		return
	}
	rawIndex := strings.TrimPrefix(user.Mode, modeWebMistakePrefix)
	if rawIndex == user.Mode || rawIndex == "" {
		writeAPIError(w, http.StatusConflict, "No active mistake practice.")
		return
	}
	index, err := strconv.Atoi(rawIndex)
	mistakes := mistakesForLanguage(user)
	if err != nil || index < 0 || index >= len(mistakes) {
		_ = api.bot.store.setMode(user.TelegramID, "idle")
		writeAPIError(w, http.StatusNotFound, "Не нашёл ошибку для тренировки. Откройте словарь ошибок заново.")
		return
	}
	mistake := mistakes[index]
	if normalizeAnswer(answer) != normalizeAnswer(mistake.Correction) {
		writeJSON(w, http.StatusOK, map[string]any{
			"correct": false,
			"message": ui(user).TryAgain,
			"mistake": webMistakeDTOFromEntry(mistake, index),
		})
		return
	}
	if err := api.bot.store.removeMistake(user.TelegramID, user.LearningLanguage, mistake.Word, mistake.Correction); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := api.bot.store.addXP(user.TelegramID, 8); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := api.bot.store.setMode(user.TelegramID, "idle"); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	promotedTo, err := api.bot.maybePromoteLearningLevelForWeb(r.Context(), user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"correct":               true,
		"mistake":               webMistakeDTOFromEntry(mistake, index),
		"correction_audio_text": mistake.Correction,
		"xp":                    8,
		"promoted_to":           promotedTo,
		"user":                  api.userDTO(refreshed),
	})
}

func (api *webAPI) handleMistakesClear(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if err := api.bot.store.clearMistakes(user.TelegramID, user.LearningLanguage); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":   true,
		"user": api.userDTO(refreshed),
	})
}

func (api *webAPI) handleToolVoiceText(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	copy := ui(user)
	if api.bot.openrouter == nil || strings.TrimSpace(api.cfg.OpenRouterSTTModel) == "" {
		writeAPIError(w, http.StatusServiceUnavailable, "Голосовой инструмент пока не настроен.")
		return
	}
	if !user.isPremium(time.Now()) {
		writeAPIError(w, http.StatusPaymentRequired, copy.Tool.VoicePremiumRequired)
		return
	}
	voiceLimit := voiceLimitFor(user)
	if user.VoiceToday >= voiceLimit {
		writeAPIError(w, http.StatusTooManyRequests, fmt.Sprintf(copy.Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
		return
	}
	audioBytes, filename, err := readUploadedFile(w, r, "voice", webMaxVoiceUploadBytes)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	format := audioUploadFormat(filename, http.DetectContentType(audioBytes))
	transcript, err := api.bot.openrouter.transcribeAudio(r.Context(), api.cfg.OpenRouterSTTModel, audioBytes, format)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(copy.Tool.VoiceTranscribeFailed, err.Error()))
		return
	}
	translation, err := api.bot.openrouter.complete(r.Context(), translationPrompt(transcript, userInterfaceLanguage(user)), 0.1, 500)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(copy.Tool.VoiceTranslationFailed, err.Error()))
		return
	}
	if err := api.bot.store.incrementVoice(user.TelegramID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	result := voiceToolResult(copy, transcript, translation)
	history := trimPracticeHistory(append(user.PracticeHistory, toolContextEntry("voice", result)), practiceMemoryLimit)
	if err := api.bot.store.savePracticeHistory(user.TelegramID, history); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"result":      result,
		"transcript":  transcript,
		"translation": translation,
		"user":        api.userDTO(refreshed),
	})
}

func (api *webAPI) handleToolImageTranslate(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	copy := ui(user)
	if api.bot.openrouter == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "Инструмент картинок пока не настроен.")
		return
	}
	if !user.isPremium(time.Now()) {
		writeAPIError(w, http.StatusPaymentRequired, copy.Tool.ImagePremiumRequired)
		return
	}
	imageBytes, filename, err := readUploadedFile(w, r, "image", webMaxImageUploadBytes)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	mimeType := imageUploadMime(filename, imageBytes)
	if !strings.HasPrefix(mimeType, "image/") {
		writeAPIError(w, http.StatusBadRequest, "Загрузите изображение с текстом.")
		return
	}
	result, err := api.bot.openrouter.translateImageText(r.Context(), imageBytes, mimeType, userLearningLanguage(user), userInterfaceLanguage(user))
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(copy.Tool.ImageReadFailed, err.Error()))
		return
	}
	history := trimPracticeHistory(append(user.PracticeHistory, toolContextEntry("image", result)), practiceMemoryLimit)
	if err := api.bot.store.savePracticeHistory(user.TelegramID, history); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"result": result,
		"user":   api.userDTO(refreshed),
	})
}

func (api *webAPI) handleToolTranslator(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	copy := ui(user)
	if api.bot.openrouter == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "Переводчик пока не настроен.")
		return
	}
	if !api.allowRate(r, "tool-translator:"+strconv.FormatInt(user.TelegramID, 10), 40, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, "Слишком много переводов подряд. Попробуйте чуть позже.")
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, webMaxVoiceUploadBytes+webMaxImageUploadBytes+webMaxMultipartOverhead)
	if err := r.ParseMultipartForm(webMaxVoiceUploadBytes + webMaxImageUploadBytes); err != nil {
		writeAPIError(w, http.StatusBadRequest, "Не удалось прочитать форму перевода.")
		return
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	sourceCode := normalizeTranslatorSource(r.FormValue("source_language"))
	targetCode := normalizeTranslatorTarget(r.FormValue("target_language"))
	if targetCode == "" {
		targetCode = defaultTranslatorTarget(user)
	}
	text := strings.TrimSpace(r.FormValue("text"))

	var result translationToolResult
	imageBytes, imageName, err := readOptionalUploadedFile(r, "image", webMaxImageUploadBytes)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	voiceBytes, voiceName, err := readOptionalUploadedFile(r, "voice", webMaxVoiceUploadBytes)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch {
	case len(imageBytes) > 0:
		if !user.isPremium(time.Now()) {
			writeAPIError(w, http.StatusPaymentRequired, copy.Tool.ImagePremiumRequired)
			return
		}
		mimeType := imageUploadMime(imageName, imageBytes)
		if !strings.HasPrefix(mimeType, "image/") {
			writeAPIError(w, http.StatusBadRequest, "Загрузите изображение с текстом.")
			return
		}
		result, err = api.bot.openrouter.translateImageForToolWithModel(r.Context(), api.cfg.OpenRouterTranslatorModel, imageBytes, mimeType, sourceCode, targetCode, userInterfaceLanguage(user))
		if err != nil {
			writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(copy.Tool.ImageReadFailed, err.Error()))
			return
		}
	case len(voiceBytes) > 0:
		if !user.isPremium(time.Now()) {
			writeAPIError(w, http.StatusPaymentRequired, copy.Tool.VoicePremiumRequired)
			return
		}
		voiceLimit := voiceLimitFor(user)
		if user.VoiceToday >= voiceLimit {
			writeAPIError(w, http.StatusTooManyRequests, fmt.Sprintf(copy.Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
			return
		}
		if strings.TrimSpace(api.cfg.OpenRouterSTTModel) == "" {
			writeAPIError(w, http.StatusServiceUnavailable, "Голосовой ввод пока не настроен.")
			return
		}
		format := audioUploadFormat(voiceName, http.DetectContentType(voiceBytes))
		text, err = api.bot.openrouter.transcribeAudio(r.Context(), api.cfg.OpenRouterSTTModel, voiceBytes, format)
		if err != nil {
			writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(copy.Tool.VoiceTranscribeFailed, err.Error()))
			return
		}
		if err := api.bot.store.incrementVoice(user.TelegramID); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		fallthrough
	default:
		if strings.TrimSpace(text) == "" {
			writeAPIError(w, http.StatusBadRequest, "Введите текст, запишите голос или загрузите фото.")
			return
		}
		translation, err := api.bot.openrouter.completeWithModel(r.Context(), api.cfg.OpenRouterTranslatorModel, translationToolPrompt(text, sourceCode, targetCode, userInterfaceLanguage(user)), 0.1, 900)
		if err != nil {
			writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(copy.Tool.VoiceTranslationFailed, err.Error()))
			return
		}
		result = translationToolResult{SourceText: text, Translation: translation}
	}

	result.SourceText = strings.TrimSpace(result.SourceText)
	result.Translation = strings.TrimSpace(result.Translation)
	if result.Translation == "" {
		writeAPIError(w, http.StatusBadGateway, "Переводчик вернул пустой перевод.")
		return
	}
	resultText := translatorToolResultText(copy, result.SourceText, result.Translation)
	history := trimPracticeHistory(append(user.PracticeHistory, toolContextEntry("translator", resultText)), practiceMemoryLimit)
	if err := api.bot.store.savePracticeHistory(user.TelegramID, history); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, _ := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	writeJSON(w, http.StatusOK, map[string]any{
		"source_text":     result.SourceText,
		"translation":     result.Translation,
		"result":          resultText,
		"source_language": sourceCode,
		"target_language": targetCode,
		"user":            api.userDTO(refreshed),
	})
}

func (api *webAPI) handleToolTranslatorSpeech(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if api.bot.openrouter == nil || strings.TrimSpace(api.cfg.OpenRouterTTSModel) == "" {
		writeAPIError(w, http.StatusServiceUnavailable, "Озвучка пока не настроена.")
		return
	}
	if !api.allowRate(r, "tool-translator-speech:"+strconv.FormatInt(user.TelegramID, 10), 25, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, "Слишком много запросов озвучки подряд. Попробуйте чуть позже.")
		return
	}
	var req struct {
		Text           string `json:"text"`
		TargetLanguage string `json:"target_language"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	text := strings.TrimSpace(req.Text)
	if text == "" {
		writeAPIError(w, http.StatusBadRequest, "text is required")
		return
	}
	if len([]rune(text)) > 1200 {
		writeAPIError(w, http.StatusBadRequest, "Текст для озвучки слишком длинный.")
		return
	}
	copy := ui(user)
	if !user.isPremium(time.Now()) {
		writeAPIError(w, http.StatusPaymentRequired, copy.Tool.VoicePremiumRequired)
		return
	}
	voiceLimit := voiceLimitFor(user)
	if user.VoiceToday >= voiceLimit {
		writeAPIError(w, http.StatusTooManyRequests, fmt.Sprintf(copy.Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
		return
	}
	targetCode := normalizeTranslatorTarget(req.TargetLanguage)
	if targetCode == "" {
		targetCode = defaultTranslatorTarget(user)
	}
	sum := sha256.Sum256([]byte(api.cfg.OpenRouterTTSModel + "|" + api.cfg.OpenRouterTTSVoice + "|" + targetCode + "|" + text))
	cacheKey := "translator|" + fmt.Sprintf("%x", sum[:])
	if audio, ok := api.cachedPronunciation(cacheKey); ok {
		if err := api.bot.store.incrementVoice(user.TelegramID); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return
		}
		writeToolAudio(w, "translation.mp3", audio)
		return
	}
	audio, err := api.bot.openrouter.synthesizeSpeech(r.Context(), api.cfg.OpenRouterTTSModel, api.cfg.OpenRouterTTSVoice, text)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, "Не получилось озвучить перевод: "+err.Error())
		return
	}
	api.cachePronunciation(cacheKey, audio)
	if err := api.bot.store.incrementVoice(user.TelegramID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeToolAudio(w, "translation.mp3", audio)
}

func userLearnedWord(user userState, wordID string) bool {
	for _, word := range user.LearnedWords {
		if word.ID == wordID {
			return true
		}
	}
	return false
}

func (api *webAPI) cachedPronunciation(key string) ([]byte, bool) {
	api.pronunciationMu.Lock()
	defer api.pronunciationMu.Unlock()
	audio, ok := api.pronunciationCache[key]
	if !ok {
		return nil, false
	}
	return append([]byte(nil), audio...), true
}

func (api *webAPI) cachePronunciation(key string, audio []byte) {
	if len(audio) == 0 {
		return
	}
	api.pronunciationMu.Lock()
	defer api.pronunciationMu.Unlock()
	if len(api.pronunciationCache) >= 256 {
		for cachedKey := range api.pronunciationCache {
			delete(api.pronunciationCache, cachedKey)
			break
		}
	}
	api.pronunciationCache[key] = append([]byte(nil), audio...)
}

func writePronunciationAudio(w http.ResponseWriter, word vocabWord, audio []byte) {
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", `inline; filename="`+legacyVocabID(word.ID)+`.mp3"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(audio)
}

func writeToolAudio(w http.ResponseWriter, filename string, audio []byte) {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = "audio.mp3"
	}
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", `inline; filename="`+filename+`"`)
	w.Header().Set("Cache-Control", "private, max-age=3600")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(audio)
}

func webWordOptionsDTO(options []vocabWord) []webWordOptionDTO {
	dto := make([]webWordOptionDTO, 0, len(options))
	for _, option := range options {
		dto = append(dto, webWordOptionDTO{ID: option.ID, Text: option.English})
	}
	return dto
}

func spellingMistakeEntry(user userState, word learnedWordEntry, now time.Time) mistakeEntry {
	prompt := strings.TrimSpace(learnedWordTranslation(word, user.InterfaceLanguage))
	if prompt == "" || normalizeAnswer(prompt) == normalizeAnswer(word.English) {
		prompt = "spelling: " + strings.TrimSpace(word.English)
	}
	return mistakeEntry{
		Language:    normalizeLearningLanguage(user.LearningLanguage),
		Word:        prompt,
		Correction:  strings.TrimSpace(word.English),
		Explanation: "Spelling practice needs another review.",
		AddedAt:     now,
	}
}

func webMistakeDTOsNewestFirst(mistakes []mistakeEntry) []webMistakeDTO {
	indices := make([]int, 0, len(mistakes))
	for index := range mistakes {
		indices = append(indices, index)
	}
	sort.SliceStable(indices, func(i, j int) bool {
		left := mistakes[indices[i]].AddedAt
		right := mistakes[indices[j]].AddedAt
		if left.Equal(right) {
			return indices[i] > indices[j]
		}
		if left.IsZero() {
			return false
		}
		if right.IsZero() {
			return true
		}
		return left.After(right)
	})
	items := make([]webMistakeDTO, 0, len(indices))
	for _, index := range indices {
		items = append(items, webMistakeDTOFromEntry(mistakes[index], index))
	}
	return items
}

func webMistakeDTOs(mistakes []mistakeEntry, start int, end int) []webMistakeDTO {
	if start < 0 {
		start = 0
	}
	if end > len(mistakes) {
		end = len(mistakes)
	}
	if start > end {
		start = end
	}
	items := make([]webMistakeDTO, 0, end-start)
	for index := start; index < end; index++ {
		items = append(items, webMistakeDTOFromEntry(mistakes[index], index))
	}
	return items
}

func webMistakeDTOFromEntry(mistake mistakeEntry, index int) webMistakeDTO {
	addedAt := ""
	if !mistake.AddedAt.IsZero() {
		addedAt = mistake.AddedAt.Format(time.RFC3339)
	}
	return webMistakeDTO{
		Index:       index,
		Word:        mistake.Word,
		Correction:  mistake.Correction,
		Explanation: mistake.Explanation,
		AddedAt:     addedAt,
	}
}

func isMultipartRequest(r *http.Request) bool {
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type"))), "multipart/form-data")
}

type webPracticeInput struct {
	text          string
	fromVoice     bool
	pronunciation *pronunciationAssessment
	imageBytes    []byte
	imageName     string
}

func (api *webAPI) readPracticeInput(w http.ResponseWriter, r *http.Request, user userState) (webPracticeInput, bool) {
	if !isMultipartRequest(r) {
		var req struct {
			Text string `json:"text"`
		}
		if !decodeJSONRequest(w, r, &req) {
			return webPracticeInput{}, false
		}
		return webPracticeInput{text: strings.TrimSpace(req.Text)}, true
	}

	r.Body = http.MaxBytesReader(w, r.Body, webMaxVoiceUploadBytes+webMaxImageUploadBytes+webMaxMultipartOverhead)
	if err := r.ParseMultipartForm(webMaxVoiceUploadBytes + webMaxImageUploadBytes); err != nil {
		writeAPIError(w, http.StatusBadRequest, "Could not read practice form.")
		return webPracticeInput{}, false
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	input := webPracticeInput{text: strings.TrimSpace(r.FormValue("text"))}
	voiceBytes, voiceName, err := readOptionalUploadedFile(r, "voice", webMaxVoiceUploadBytes)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return webPracticeInput{}, false
	}
	if len(voiceBytes) > 0 {
		copy := ui(user)
		if !api.bot.voiceRecognitionConfigured() {
			writeAPIError(w, http.StatusServiceUnavailable, "Voice recognition is not configured.")
			return webPracticeInput{}, false
		}
		if !user.isPremium(time.Now()) {
			writeAPIError(w, http.StatusPaymentRequired, copy.Tool.VoicePremiumRequired)
			return webPracticeInput{}, false
		}
		voiceLimit := voiceLimitFor(user)
		if user.VoiceToday >= voiceLimit {
			writeAPIError(w, http.StatusTooManyRequests, fmt.Sprintf(copy.Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
			return webPracticeInput{}, false
		}
		format := audioUploadFormat(voiceName, http.DetectContentType(voiceBytes))
		transcription, err := api.bot.transcribeLearningVoice(r.Context(), user, voiceBytes, format, "")
		if err != nil {
			writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(copy.Tool.VoiceTranscribeFailed, err.Error()))
			return webPracticeInput{}, false
		}
		if err := api.bot.store.incrementVoice(user.TelegramID); err != nil {
			writeAPIError(w, http.StatusInternalServerError, err.Error())
			return webPracticeInput{}, false
		}
		assessment := api.bot.buildPronunciationAssessment(r.Context(), user, "", transcription, pronunciationModeFree)
		input.text = strings.TrimSpace(transcription.Text)
		input.fromVoice = true
		input.pronunciation = &assessment
	}

	imageBytes, imageName, err := readOptionalUploadedFile(r, "image", webMaxImageUploadBytes)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return webPracticeInput{}, false
	}
	input.imageBytes = imageBytes
	input.imageName = imageName
	return input, true
}

func (api *webAPI) readLearningAnswerInput(w http.ResponseWriter, r *http.Request, user userState, expected string) (string, bool, *pronunciationAssessment, bool) {
	if !isMultipartRequest(r) {
		var req struct {
			Text string `json:"text"`
		}
		if !decodeJSONRequest(w, r, &req) {
			return "", false, nil, false
		}
		return strings.TrimSpace(req.Text), false, nil, true
	}

	r.Body = http.MaxBytesReader(w, r.Body, webMaxVoiceUploadBytes+webMaxMultipartOverhead)
	if err := r.ParseMultipartForm(webMaxVoiceUploadBytes); err != nil {
		writeAPIError(w, http.StatusBadRequest, "Could not read voice answer form.")
		return "", false, nil, false
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}

	text := strings.TrimSpace(r.FormValue("text"))
	voiceBytes, voiceName, err := readOptionalUploadedFile(r, "voice", webMaxVoiceUploadBytes)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return "", false, nil, false
	}
	if len(voiceBytes) == 0 {
		return text, false, nil, true
	}

	copy := ui(user)
	if !api.bot.voiceRecognitionConfigured() {
		writeAPIError(w, http.StatusServiceUnavailable, "Voice recognition is not configured.")
		return "", false, nil, false
	}
	if !user.isPremium(time.Now()) {
		writeAPIError(w, http.StatusPaymentRequired, copy.Tool.VoicePremiumRequired)
		return "", false, nil, false
	}
	voiceLimit := voiceLimitFor(user)
	if user.VoiceToday >= voiceLimit {
		writeAPIError(w, http.StatusTooManyRequests, fmt.Sprintf(copy.Tool.VoiceLimitReached, user.VoiceToday, voiceLimit))
		return "", false, nil, false
	}

	format := audioUploadFormat(voiceName, http.DetectContentType(voiceBytes))
	transcription, err := api.bot.transcribeLearningVoice(r.Context(), user, voiceBytes, format, expected)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, fmt.Sprintf(copy.Tool.VoiceTranscribeFailed, err.Error()))
		return "", false, nil, false
	}
	if err := api.bot.store.incrementVoice(user.TelegramID); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return "", false, nil, false
	}
	mode := pronunciationModeFree
	if strings.TrimSpace(expected) != "" {
		mode = pronunciationModeExact
	}
	assessment := api.bot.buildPronunciationAssessment(r.Context(), user, expected, transcription, mode)
	return strings.TrimSpace(transcription.Text), true, &assessment, true
}

func readUploadedFile(w http.ResponseWriter, r *http.Request, field string, maxBytes int64) ([]byte, string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes+webMaxMultipartOverhead)
	file, header, err := r.FormFile(field)
	if err != nil {
		return nil, "", errors.New("Файл не найден в запросе.")
	}
	if r.MultipartForm != nil {
		defer r.MultipartForm.RemoveAll()
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(data)) > maxBytes {
		return nil, "", errors.New("Файл слишком большой.")
	}
	if len(data) == 0 {
		return nil, "", errors.New("Файл пустой.")
	}
	filename := ""
	if header != nil {
		filename = header.Filename
	}
	return data, filename, nil
}

func readOptionalUploadedFile(r *http.Request, field string, maxBytes int64) ([]byte, string, error) {
	if r.MultipartForm == nil {
		return nil, "", nil
	}
	files := r.MultipartForm.File[field]
	if len(files) == 0 || files[0] == nil {
		return nil, "", nil
	}
	header := files[0]
	file, err := header.Open()
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, maxBytes+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(data)) > maxBytes {
		return nil, "", errors.New("Файл слишком большой.")
	}
	if len(data) == 0 {
		return nil, "", nil
	}
	return data, header.Filename, nil
}

func audioUploadFormat(filename string, mimeType string) string {
	lowerName := strings.ToLower(strings.TrimSpace(filename))
	lowerMime := strings.ToLower(strings.TrimSpace(mimeType))
	switch {
	case strings.HasSuffix(lowerName, ".ogg"), strings.HasSuffix(lowerName, ".oga"), strings.HasSuffix(lowerName, ".opus"), strings.Contains(lowerMime, "ogg"), strings.Contains(lowerMime, "opus"):
		return "ogg"
	case strings.HasSuffix(lowerName, ".webm"), strings.Contains(lowerMime, "webm"):
		return "webm"
	case strings.HasSuffix(lowerName, ".wav"), strings.Contains(lowerMime, "wav"):
		return "wav"
	case strings.HasSuffix(lowerName, ".m4a"), strings.HasSuffix(lowerName, ".mp4"), strings.Contains(lowerMime, "mp4"), strings.Contains(lowerMime, "m4a"):
		return "m4a"
	default:
		return "mp3"
	}
}

func imageUploadMime(filename string, data []byte) string {
	mimeType := strings.ToLower(http.DetectContentType(data))
	if strings.HasPrefix(mimeType, "image/") {
		return mimeType
	}
	lowerName := strings.ToLower(strings.TrimSpace(filename))
	switch {
	case strings.HasSuffix(lowerName, ".gif"):
		return "image/gif"
	case strings.HasSuffix(lowerName, ".bmp"):
		return "image/bmp"
	case strings.HasSuffix(lowerName, ".webp"):
		return "image/webp"
	case strings.HasSuffix(lowerName, ".heic"):
		return "image/heic"
	case strings.HasSuffix(lowerName, ".heif"):
		return "image/heif"
	case strings.HasSuffix(lowerName, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(lowerName, ".png"):
		return "image/png"
	case strings.HasSuffix(lowerName, ".jpg"), strings.HasSuffix(lowerName, ".jpeg"):
		return "image/jpeg"
	default:
		return mimeType
	}
}

func bugReportUploadFiles(r *http.Request) []*multipart.FileHeader {
	if r == nil || r.MultipartForm == nil {
		return nil
	}
	files := append([]*multipart.FileHeader{}, r.MultipartForm.File["screenshot"]...)
	files = append(files, r.MultipartForm.File["attachment"]...)
	return files
}

func bugReportCaption(reportID string, user userState, view string, message string, userAgent string) string {
	parts := []string{
		"Bug report Poliglot AI",
		"ID: " + reportID,
		fmt.Sprintf("From: %s (%d)", strings.TrimSpace(user.FirstName), user.TelegramID),
	}
	if view != "" {
		parts = append(parts, "View: "+view)
	}
	if userAgent != "" {
		parts = append(parts, "User-Agent: "+userAgent)
	}
	parts = append(parts, "", message)
	return strings.Join(parts, "\n")
}

func appendBugReportFile(reportID string, user userState, view string, message string, userAgent string) error {
	if err := os.MkdirAll(filepath.Join("tmp", "bug-reports"), 0755); err != nil {
		return err
	}
	payload := map[string]any{
		"id":          reportID,
		"telegram_id": user.TelegramID,
		"first_name":  user.FirstName,
		"view":        view,
		"message":     message,
		"user_agent":  userAgent,
		"created_at":  time.Now().UTC().Format(time.RFC3339Nano),
	}
	line, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join("tmp", "bug-reports", "bug-reports.jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(append(line, '\n')); err != nil {
		return err
	}
	return nil
}

func writeBugReportScreenshot(reportID string, index int, filename string, data []byte) (string, error) {
	if err := os.MkdirAll(filepath.Join("tmp", "bug-reports"), 0755); err != nil {
		return "", err
	}
	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" || len(ext) > 8 {
		ext = ".img"
	}
	path := filepath.Join("tmp", "bug-reports", fmt.Sprintf("%s-%02d%s", reportID, index+1, ext))
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
}

func (api *webAPI) handleProgress(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": api.userDTO(user)})
}

func (api *webAPI) handlePremiumPlans(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"enabled":        api.webYooKassaReady(),
		"crypto_enabled": api.bot.cryptoPaymentsReady(),
		"plans":          api.premiumPlansDTO(user),
	})
}

func (api *webAPI) handlePremiumPayment(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	client := api.webYooKassaClient()
	if client == nil || !api.cfg.webYooKassaEnabled() {
		writeAPIError(w, http.StatusServiceUnavailable, premiumUI(user).YooKassaUnavailable)
		return
	}
	var req struct {
		Product string `json:"product"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	plan, ok := api.bot.premiumPlan(req.Product)
	if !ok {
		writeAPIError(w, http.StatusBadRequest, "unknown product")
		return
	}
	plan = localizedPremiumPlan(user, plan)
	returnURL := api.cfg.webPaymentReturnURL()
	payment, err := client.createPremiumPayment(r.Context(), user.TelegramID, user.FirstName, user.InterfaceLanguage, plan, returnURL, "web")
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, premiumUI(user).YooKassaCreateFailed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"payment_id":       payment.ID,
		"confirmation_url": payment.Confirmation.ConfirmationURL,
	})
}

func (api *webAPI) handlePremiumStarsPayment(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if api.bot == nil || api.bot.telegram == nil {
		writeAPIError(w, http.StatusServiceUnavailable, webPremiumStarsMessage(user, "unavailable"))
		return
	}
	if user.TelegramID <= 0 {
		writeAPIError(w, http.StatusConflict, webPremiumStarsMessage(user, "link_telegram"))
		return
	}
	if !api.allowRate(r, "premium-stars:"+strconv.FormatInt(user.TelegramID, 10), 8, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, webPremiumStarsMessage(user, "rate_limited"))
		return
	}
	var req struct {
		Product string `json:"product"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	if _, ok := api.bot.premiumPlan(req.Product); !ok {
		writeAPIError(w, http.StatusBadRequest, "unknown product")
		return
	}
	if err := api.bot.sendPremiumInvoice(r.Context(), user.TelegramID, user, req.Product); err != nil {
		writeAPIError(w, http.StatusBadGateway, webPremiumStarsMessage(user, "send_failed"))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"sent":    true,
		"message": webPremiumStarsMessage(user, "sent"),
		"bot_url": api.premiumStarsBotURL(req.Product),
	})
}

func (api *webAPI) handlePremiumActivationKey(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if api.bot == nil || api.bot.activationKeys == nil {
		writeAPIError(w, http.StatusServiceUnavailable, "Activation keys are not configured.")
		return
	}
	if !api.allowRate(r, "activation-key:"+strconv.FormatInt(user.TelegramID, 10), 8, 10*time.Minute) {
		writeAPIError(w, http.StatusTooManyRequests, "Too many activation attempts. Try again later.")
		return
	}
	var req struct {
		Key string `json:"key"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	grant, err := api.bot.activationKeys.redeem(r.Context(), req.Key, user.TelegramID)
	if err != nil {
		status := http.StatusInternalServerError
		message := err.Error()
		if errors.Is(err, errActivationKeyInvalid) {
			status = http.StatusBadRequest
			message = "Invalid key format. Use XXXX-XXXX-XXXX-XXXX."
		} else if errors.Is(err, errActivationKeyUsed) {
			status = http.StatusConflict
			message = "This key has already been activated."
		} else if errors.Is(err, errActivationKeyMissing) {
			status = http.StatusNotFound
			message = "Key not found."
		}
		writeAPIError(w, status, message)
		return
	}
	duration := time.Duration(grant.DurationDays) * 24 * time.Hour
	until, err := api.bot.store.extendPremium(user.TelegramID, "activation-key:"+grant.Key, duration, grant.Tier)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	refreshed, err := api.bot.store.getOrCreateUser(user.TelegramID, user.FirstName)
	if err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	payload := api.sessionPayload(refreshed)
	payload["activation"] = map[string]any{
		"duration_days": grant.DurationDays,
		"tier":          grant.Tier,
		"premium_until": until.Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, payload)
}

func (api *webAPI) handleCryptoPremiumPayment(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	if !api.bot.cryptoPaymentsReady() {
		writeAPIError(w, http.StatusServiceUnavailable, "Direct crypto payments are not configured yet")
		return
	}
	var req struct {
		Product string `json:"product"`
		Method  string `json:"method"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	payment, invoiceURL, err := api.bot.createDirectCryptoPayment(r.Context(), req.Product, user, req.Method)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := api.bot.saveDirectCryptoPayment(payment); err != nil {
		writeAPIError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cryptoPaymentDTO(payment, invoiceURL))
}

func (api *webAPI) handleCryptoPremiumCheck(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodPost) {
		return
	}
	user, err := api.currentUser(w, r)
	if err != nil {
		api.writeCurrentUserError(w, err)
		return
	}
	var req struct {
		PaymentID string `json:"payment_id"`
	}
	if !decodeJSONRequest(w, r, &req) {
		return
	}
	payment, refreshed, paid, err := api.bot.refreshDirectCryptoPayment(r.Context(), req.PaymentID, user)
	if err != nil {
		writeAPIError(w, http.StatusBadGateway, err.Error())
		return
	}
	payload := cryptoPaymentDTO(payment, api.bot.cryptoPaymentURL(payment))
	payload["paid"] = paid
	if paid {
		payload["user"] = api.userDTO(refreshed)
	}
	writeJSON(w, http.StatusOK, payload)
}

func (api *webAPI) currentUser(w http.ResponseWriter, r *http.Request) (userState, error) {
	_ = w
	store, err := api.webAccountStore()
	if err != nil {
		return userState{}, err
	}
	userID, ok := api.sessionUserID(r)
	if !ok {
		return userState{}, errWebAuthRequired
	}
	account, ok, err := store.getWebAccountByUserID(userID)
	if err != nil {
		return userState{}, err
	}
	if !ok {
		return userState{}, errWebAuthRequired
	}
	return api.bot.store.getOrCreateUser(account.UserID, account.Login)
}

func (api *webAPI) writeCurrentUserError(w http.ResponseWriter, err error) {
	if errors.Is(err, errWebAuthRequired) {
		writeAPIError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if errors.Is(err, errWebAccountUnavailable) {
		writeAPIError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	writeAPIError(w, http.StatusInternalServerError, err.Error())
}

func (api *webAPI) allowRate(r *http.Request, scope string, limit int, window time.Duration) bool {
	if limit <= 0 || window <= 0 {
		return true
	}
	key := clientRateKey(r) + "|" + scope
	now := time.Now()
	api.rateMu.Lock()
	defer api.rateMu.Unlock()
	if len(api.rateBuckets) > webRateSweepLimit {
		for bucketKey, bucket := range api.rateBuckets {
			if now.After(bucket.ResetAt) {
				delete(api.rateBuckets, bucketKey)
			}
		}
	}
	bucket := api.rateBuckets[key]
	if bucket.ResetAt.IsZero() || !now.Before(bucket.ResetAt) {
		bucket = webRateBucket{ResetAt: now.Add(window)}
	}
	if bucket.Count >= limit {
		api.rateBuckets[key] = bucket
		return false
	}
	bucket.Count++
	api.rateBuckets[key] = bucket
	return true
}

func clientRateKey(r *http.Request) string {
	if forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); forwarded != "" {
		first, _, _ := strings.Cut(forwarded, ",")
		if ip := net.ParseIP(strings.TrimSpace(first)); ip != nil {
			return ip.String()
		}
	}
	if realIP := net.ParseIP(strings.TrimSpace(r.Header.Get("X-Real-IP"))); realIP != nil {
		return realIP.String()
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

func (api *webAPI) webAccountStore() (webAccountStore, error) {
	store, ok := api.bot.store.(webAccountStore)
	if !ok {
		return nil, errWebAccountUnavailable
	}
	return store, nil
}

func (api *webAPI) sessionUserID(r *http.Request) (int64, bool) {
	cookie, err := r.Cookie(webSessionCookieName)
	if err != nil {
		return 0, false
	}
	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		return 0, false
	}
	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || userID == 0 {
		return 0, false
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, false
	}
	expected := api.signSessionID(parts[0])
	if subtle.ConstantTimeCompare(signature, expected) != 1 {
		return 0, false
	}
	return userID, true
}

func (api *webAPI) setSessionCookie(w http.ResponseWriter, userID int64) {
	rawID := strconv.FormatInt(userID, 10)
	value := rawID + "." + base64.RawURLEncoding.EncodeToString(api.signSessionID(rawID))
	cookie := &http.Cookie{
		Name:     webSessionCookieName,
		Value:    value,
		Path:     "/",
		MaxAge:   int((365 * 24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   api.cfg.WebCookieSecure,
		SameSite: webCookieSameSite(api.cfg.WebCookieSameSite),
	}
	if strings.TrimSpace(api.cfg.WebCookieDomain) != "" {
		cookie.Domain = strings.TrimSpace(api.cfg.WebCookieDomain)
	}
	http.SetCookie(w, cookie)
}

func (api *webAPI) clearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     webSessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   api.cfg.WebCookieSecure,
		SameSite: webCookieSameSite(api.cfg.WebCookieSameSite),
	}
	if strings.TrimSpace(api.cfg.WebCookieDomain) != "" {
		cookie.Domain = strings.TrimSpace(api.cfg.WebCookieDomain)
	}
	http.SetCookie(w, cookie)
}

func webCookieSameSite(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

func (api *webAPI) signSessionID(rawID string) []byte {
	mac := hmac.New(sha256.New, api.sessionSecret)
	_, _ = mac.Write([]byte(rawID))
	return mac.Sum(nil)
}

func newWebUserID() (int64, error) {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return 0, err
	}
	value := int64(binary.BigEndian.Uint64(bytes[:]) & ((uint64(1) << 62) - 1))
	if value == 0 {
		value = time.Now().UnixNano()
	}
	if value < 0 {
		value = -value
	}
	return -value, nil
}

func (api *webAPI) sessionPayload(user userState) map[string]any {
	account := map[string]any{
		"login":         user.FirstName,
		"password_set":  true,
		"profile_ready": true,
	}
	if store, err := api.webAccountStore(); err == nil {
		if webAccount, ok, err := store.getWebAccountByUserID(user.TelegramID); err == nil && ok {
			if webAccount.PasswordSet {
				account["login"] = webAccount.Login
			} else {
				account["login"] = ""
			}
			account["password_set"] = webAccount.PasswordSet
			account["profile_ready"] = webAccount.PasswordSet
		}
	}
	return map[string]any{
		"authenticated":       true,
		"account":             account,
		"user":                api.userDTO(user),
		"copy":                webCopyDTO(ui(user)),
		"tool_copy":           webToolCopyDTO(ui(user).Tool),
		"system_copy":         webSystemCopyDTO(systemUI(user)),
		"premium_copy":        webPremiumCopyDTO(premiumUI(user)),
		"interface_languages": webLanguageDTOs(interfaceLanguages()),
		"learning_languages":  webLanguageDTOs(learningLanguages),
		"premium_plans":       api.premiumPlansDTO(user),
		"yookassa_enabled":    api.webYooKassaReady(),
		"crypto_enabled":      api.bot.cryptoPaymentsReady(),
		"captcha":             api.webCaptchaDTO(),
		"telegram_login_bot":  api.cfg.WebTelegramLoginBot,
		"web_app_url":         api.cfg.WebAppURL,
	}
}

func (api *webAPI) writeSession(w http.ResponseWriter, user userState) {
	writeJSON(w, http.StatusOK, api.sessionPayload(user))
}

func (api *webAPI) writeAnonymousSession(w http.ResponseWriter) {
	writeJSON(w, http.StatusOK, map[string]any{
		"authenticated":       false,
		"interface_languages": webLanguageDTOs(interfaceLanguages()),
		"captcha":             api.webCaptchaDTO(),
		"telegram_login_bot":  api.cfg.WebTelegramLoginBot,
		"web_app_url":         api.cfg.WebAppURL,
		"support": map[string]string{
			"telegram": "@AsaselD",
			"email":    "supportpoliglotai@gmail.com",
		},
	})
}

func (api *webAPI) userDTO(user userState) map[string]any {
	level, title, currentXP, neededXP := knowledgeLevel(user.XP)
	lessonLimit, practiceLimit := limitsFor(user)
	var referralInvitees []referralInviteeEntry
	if api != nil && api.bot != nil && api.bot.store != nil {
		var err error
		referralInvitees, err = api.bot.store.referralInvitees(user.TelegramID, 100)
		if err != nil {
			log.Printf("load referral invitees for %d: %v", user.TelegramID, err)
		}
	}
	premiumUntil := ""
	if !user.PremiumUntil.IsZero() {
		premiumUntil = user.PremiumUntil.Format(time.RFC3339)
	}
	dailyBonusLastClaimedAt := ""
	if !user.DailyBonusLastClaimedAt.IsZero() {
		dailyBonusLastClaimedAt = user.DailyBonusLastClaimedAt.Format(time.RFC3339)
	}
	createdAt := ""
	if !user.CreatedAt.IsZero() {
		createdAt = user.CreatedAt.Format(time.RFC3339)
	}
	var telegramAccount any
	if user.TelegramID > 0 {
		name := strings.TrimSpace(user.FirstName)
		telegramAccount = map[string]any{
			"id":    user.TelegramID,
			"name":  name,
			"label": telegramAccountLabel(user.TelegramID, name),
		}
	}
	return map[string]any{
		"interface_language":          user.InterfaceLanguage,
		"learning_language":           user.LearningLanguage,
		"telegram_linked":             user.TelegramID > 0,
		"telegram_account":            telegramAccount,
		"created_at":                  createdAt,
		"level":                       user.Level,
		"learning_focus":              strings.TrimSpace(user.LearningFocus),
		"plan":                        planName(user),
		"premium":                     user.isPremium(time.Now()),
		"premium_until":               premiumUntil,
		"referral_code":               user.ReferralCode,
		"referral_count":              user.ReferralCount,
		"referral_balance_kopecks":    user.ReferralBalanceKopecks,
		"referral_balance":            formatRubKopecks(user.ReferralBalanceKopecks),
		"referral_balance_usdt":       formatUSDTFromKopecks(user.ReferralBalanceKopecks),
		"referral_withdraw_min":       formatRubKopecks(referralWithdrawalMinKopecks),
		"referral_withdraw_min_usdt":  formatUSDTFromKopecks(referralWithdrawalMinKopecks),
		"referral_invitees":           referralInvitees,
		"invited_by":                  user.InvitedBy,
		"xp":                          user.XP,
		"xp_level":                    level,
		"xp_title":                    title,
		"xp_current":                  currentXP,
		"xp_needed":                   neededXP,
		"lessons_today":               user.LessonsToday,
		"lesson_limit":                lessonLimit,
		"practice_today":              user.PracticeToday,
		"practice_limit":              practiceLimit,
		"voice_today":                 user.VoiceToday,
		"voice_limit":                 voiceLimitFor(user),
		"daily_bonus_claims":          user.DailyBonusClaims,
		"daily_bonus_last_claimed_at": dailyBonusLastClaimedAt,
		"lesson_count":                user.LessonCount,
		"practice_count":              user.PracticeCount,
		"word_game_count":             user.WordGameCount,
		"learned_words":               len(learnedWordsForLanguage(user)),
		"mistakes":                    len(mistakesForLanguage(user)),
		"phrasebook":                  user.Phrasebook,
		"habit_log":                   user.HabitLog,
		"navigation_layout":           normalizeNavigationLayout(user.NavigationLayout),
	}
}

func webLeaderboardDTO(entries []leaderboardEntry) []map[string]any {
	result := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		result = append(result, map[string]any{
			"name":      entry.FirstName,
			"score":     entry.Score,
			"xp":        entry.Score,
			"words":     entry.Words,
			"mistakes":  entry.Mistakes,
			"level":     entry.Level,
			"title":     entry.Title,
			"languages": entry.Languages,
		})
	}
	return result
}

func telegramAccountLabel(id int64, name string) string {
	name = strings.TrimSpace(name)
	if name != "" && id != 0 {
		return name + " (" + strconv.FormatInt(id, 10) + ")"
	}
	if id != 0 {
		return "Telegram " + strconv.FormatInt(id, 10)
	}
	if name != "" {
		return name
	}
	return ""
}

func webLanguageLeaderboardDTO(entries []languageLeaderboardEntry) []map[string]any {
	result := make([]map[string]any, 0, len(entries))
	for _, entry := range entries {
		result = append(result, map[string]any{
			"name":     entry.FirstName,
			"score":    entry.Score,
			"words":    entry.Words,
			"mistakes": entry.Mistakes,
			"level":    entry.Level,
		})
	}
	return result
}

func (api *webAPI) levelQuestionDTO(user userState, index int, score int) webLevelQuestionDTO {
	questions := levelAssessmentQuestionsForUser(user)
	question := questions[index]
	return webLevelQuestionDTO{
		Index:    index,
		Total:    len(questions),
		Question: question.Question,
		Options:  append([]string{}, question.Options...),
		Score:    score,
	}
}

func (api *webAPI) premiumPlansDTO(user userState) []map[string]any {
	plans := localizedPremiumPlans(user, api.bot.premiumPlanList())
	result := make([]map[string]any, 0, len(plans))
	for _, plan := range plans {
		result = append(result, map[string]any{
			"product":        plan.Product,
			"tier":           plan.Tier,
			"title":          plan.Title,
			"days_label":     plan.DaysLabel,
			"rub_price":      plan.RubPrice,
			"stars_price":    plan.StarsPrice,
			"usdt_price":     api.usdtPriceLabel(plan.Product),
			"crypto_enabled": api.cfg.cryptoPlanEnabled(plan.Product),
			"crypto_methods": api.cryptoPaymentMethodsDTO(plan.Product),
		})
	}
	return result
}

func (api *webAPI) usdtPriceLabel(product string) string {
	amount, err := api.cfg.cryptoUSDTAmountForProduct(product)
	if err != nil || amount <= 0 {
		return ""
	}
	return formatUnitsAmount(amount, 6) + " USDT"
}

func (api *webAPI) cryptoPaymentMethodsDTO(product string) []map[string]any {
	methods := api.cfg.cryptoPaymentMethodsForProduct(product)
	result := make([]map[string]any, 0, len(methods))
	for _, method := range methods {
		result = append(result, map[string]any{
			"id":       method.ID,
			"label":    method.Label,
			"currency": method.Currency,
			"network":  method.Network,
		})
	}
	return result
}

func webPremiumStarsMessage(user userState, key string) string {
	ru := map[string]string{
		"unavailable":   "Оплата Telegram Stars сейчас недоступна.",
		"link_telegram": "Чтобы оплатить Stars, привяжите Telegram в настройках веб-приложения.",
		"rate_limited":  "Слишком много запросов на оплату Stars. Попробуйте чуть позже.",
		"send_failed":   "Не удалось отправить счёт в Telegram. Откройте бота и попробуйте ещё раз.",
		"sent":          "Счёт в Telegram Stars отправлен в ваш Telegram-чат с ботом.",
	}
	en := map[string]string{
		"unavailable":   "Telegram Stars payments are unavailable right now.",
		"link_telegram": "Link Telegram in the web app settings to pay with Stars.",
		"rate_limited":  "Too many Stars payment requests. Try again a little later.",
		"send_failed":   "Could not send the invoice to Telegram. Open the bot and try again.",
		"sent":          "The Telegram Stars invoice was sent to your Telegram chat with the bot.",
	}
	if normalizeInterfaceLanguage(user.InterfaceLanguage) == "ru" {
		return ru[key]
	}
	return en[key]
}

func (api *webAPI) premiumStarsBotURL(product string) string {
	botName := strings.TrimPrefix(strings.TrimSpace(api.cfg.WebTelegramLoginBot), "@")
	if botName == "" {
		return ""
	}
	return "https://t.me/" + url.PathEscape(botName) + "?start=" + url.QueryEscape("buy_"+strings.TrimSpace(product))
}

func webLanguageDTOs(languages []learningLanguage) []webLanguageDTO {
	result := make([]webLanguageDTO, 0, len(languages))
	for _, language := range languages {
		nativeName := strings.TrimSpace(language.InterfaceName)
		if nativeName == "" {
			nativeName = strings.TrimSpace(language.NativeName)
		}
		result = append(result, webLanguageDTO{
			Code:          language.Code,
			Name:          language.Name,
			NativeName:    nativeName,
			InterfaceName: language.InterfaceName,
		})
	}
	return result
}

func cryptoPaymentDTO(payment cryptoPayment, invoiceURL string) map[string]any {
	return map[string]any{
		"payment_id":    payment.ID,
		"provider":      payment.Provider,
		"currency":      payment.Currency,
		"network":       payment.Network,
		"method":        cryptoPaymentMethodID(payment),
		"method_label":  cryptoPaymentMethodLabel(payment),
		"product":       payment.Product,
		"status":        payment.Status,
		"address":       payment.Address,
		"memo":          payment.Memo,
		"memo_required": strings.TrimSpace(payment.Memo) != "",
		"amount":        payment.Amount,
		"amount_nano":   payment.AmountNano,
		"amount_units":  payment.AmountNano,
		"invoice_url":   invoiceURL,
		"tx_hash":       payment.TxHash,
		"expires_at":    payment.ExpiresAt.Format(time.RFC3339),
	}
}

func webCopyDTO(copy uiCopy) map[string]string {
	return map[string]string{
		"main_menu_title":   copy.MainMenuTitle,
		"main_menu_body":    copy.MainMenuBody,
		"back":              copy.Back,
		"back_menu":         copy.BackMenu,
		"learning":          copy.Learning,
		"words":             copy.Words,
		"stats":             copy.Stats,
		"settings":          copy.Settings,
		"tools":             copy.Tools,
		"new_lesson":        copy.NewLesson,
		"practice":          copy.Practice,
		"shadowing":         copy.Shadowing,
		"level_test":        copy.LevelTest,
		"learn_words":       copy.LearnWords,
		"word_game":         copy.WordGame,
		"spelling":          copy.Spelling,
		"vocabulary":        copy.Vocabulary,
		"mistakes":          copy.Mistakes,
		"progress":          copy.Progress,
		"leaders":           copy.Leaders,
		"limits":            copy.Limits,
		"bot_language":      copy.BotLanguage,
		"learning_language": copy.LearningLanguage,
		"premium":           copy.Premium,
		"choose_bot_lang":   copy.ChooseBotLang,
		"choose_learn_lang": copy.ChooseLearnLang,
		"choose_answer":     copy.ChooseAnswer,
		"word_question":     copy.WordQuestion,
		"correct":           copy.Correct,
		"try_again":         copy.TryAgain,
		"next_word":         copy.NextWord,
		"next_page":         copy.NextPage,
		"already_learned":   copy.AlreadyLearned,
		"added_to_vocab":    copy.AddedToVocab,
		"total_learned":     copy.TotalLearned,
		"write_word":        copy.WriteWord,
		"hint":              copy.Hint,
		"good_spelling":     copy.GoodSpelling,
		"stopped":           copy.Stopped,
		"unknown_button":    copy.UnknownButton,
	}
}

func webToolCopyDTO(copy toolUICopy) map[string]string {
	return map[string]string{
		"voice_to_text":          copy.VoiceToText,
		"image_translate":        copy.ImageTranslate,
		"translator":             copy.Translator,
		"gpt_agent":              copy.GPTAgent,
		"voice_prompt":           copy.VoicePrompt,
		"image_prompt":           copy.ImagePrompt,
		"translator_prompt":      copy.TranslatorPrompt,
		"voice_mode_prompt":      copy.VoiceModePrompt,
		"image_mode_prompt":      copy.ImageModePrompt,
		"translator_mode_prompt": copy.TranslatorModePrompt,
		"transcript_label":       copy.TranscriptLabel,
		"translation_label":      copy.TranslationLabel,
		"source_language":        copy.SourceLanguage,
		"target_language":        copy.TargetLanguage,
		"auto_detect":            copy.AutoDetect,
		"web_app":                copy.WebApp,
		"voice_discuss_prompt":   copy.VoiceDiscussPrompt,
		"image_discuss_prompt":   copy.ImageDiscussPrompt,
		"voice_premium_required": copy.VoicePremiumRequired,
		"image_premium_required": copy.ImagePremiumRequired,
	}
}

func webSystemCopyDTO(copy systemUICopy) map[string]string {
	return map[string]string{
		"level_start_button":  copy.LevelStartButton,
		"level_manual_button": copy.LevelManualButton,
		"level_dont_know":     copy.LevelDontKnow,
		"lesson_answer":       copy.LessonAnswerInstruction,
		"level_test_set":      copy.LevelTestSet,
		"level_manual_set":    copy.LevelManualSet,
		"level_question_text": copy.LevelQuestionText,
		"level_current_score": copy.LevelCurrentScore,
		"lesson_kind":         copy.LessonKind,
		"practice_kind":       copy.PracticeKind,
	}
}

func webPremiumCopyDTO(copy premiumUICopy) map[string]string {
	return map[string]string{
		"limits_title":                  copy.LimitsTitle,
		"plan_label":                    copy.PlanLabel,
		"lessons_today":                 copy.LessonsToday,
		"practice_today":                copy.PracticeToday,
		"voices_today":                  copy.VoicesToday,
		"pay_sbp":                       copy.PaySBPButton,
		"invite_free":                   copy.InviteFree,
		"invite_friend":                 copy.InviteFriendButton,
		"referral_share":                copy.ReferralShareText,
		"referral_balance_title":        copy.ReferralBalanceTitle,
		"referral_invitations_label":    copy.ReferralInvitationsLabel,
		"referral_balance_hint":         copy.ReferralBalanceHint,
		"referral_withdraw":             copy.ReferralWithdrawButton,
		"referral_withdraw_unavailable": copy.ReferralWithdrawUnavailable,
		"month_price":                   copy.MonthPrice,
		"year_price":                    copy.YearPrice,
		"create_failed":                 copy.YooKassaCreateFailed,
		"unavailable":                   copy.YooKassaUnavailable,
	}
}

func (b *bot) maybePromoteLearningLevelForWeb(ctx context.Context, telegramID int64, firstName string) (string, error) {
	_ = ctx
	user, err := b.store.getOrCreateUser(telegramID, firstName)
	if err != nil {
		return "", err
	}
	next, ok := shouldPromoteLearningLevel(user)
	if !ok {
		return "", nil
	}
	if err := b.store.setUserLevel(telegramID, next); err != nil {
		return "", err
	}
	return next, nil
}

func allowMethod(w http.ResponseWriter, r *http.Request, methods ...string) bool {
	for _, method := range methods {
		if r.Method == method {
			return true
		}
	}
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	return false
}

func decodeJSONRequest(w http.ResponseWriter, r *http.Request, target any) bool {
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		writeAPIError(w, http.StatusBadRequest, "bad JSON request")
		return false
	}
	return true
}

func writeAPIError(w http.ResponseWriter, status int, message string) {
	writeAPIErrorCode(w, status, "", message)
}

func webLoginValidationCode(login string) string {
	login = strings.TrimSpace(login)
	if len(login) < 3 || len(login) > 32 {
		return "invalid_login_length"
	}
	first := login[0]
	if first == '.' || first == '-' || first == '_' {
		return "invalid_login_start"
	}
	for _, ch := range strings.ToLower(login) {
		ok := ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '_' || ch == '-' || ch == '.'
		if !ok {
			return "invalid_login_chars"
		}
	}
	return "invalid_login"
}

func writeAPIErrorCode(w http.ResponseWriter, status int, code string, message string) {
	apiError := map[string]any{"message": message}
	if code != "" {
		apiError["code"] = code
	}
	writeJSON(w, status, map[string]any{
		"error": apiError,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("write json response: %v", err)
	}
}
