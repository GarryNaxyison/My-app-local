package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	webPasswordHashAlgorithm = "pbkdf2_sha256"
	webPasswordIterations    = 120000
	webPasswordSaltBytes     = 16
	webPasswordKeyBytes      = 32
	webReferralCodeLength    = 8
	webReferralAlphabet      = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
)

var (
	errWebLoginTaken         = errors.New("web login is already taken")
	errWebAccountUnavailable = errors.New("web accounts are available only with sqlite storage")
	errWebAuthRequired       = errors.New("нужно войти или создать аккаунт")
	errWebReferralNotFound   = errors.New("Реферальный код не найден")
	errWebReferralSelf       = errors.New("Нельзя использовать свой реферальный код")
	errWebAccountTGLocked    = errors.New("Эта учётка уже связана с Telegram. Отвязать или заменить Telegram нельзя.")
	errWebTelegramUsed       = errors.New("Этот Telegram уже связан с другим веб-аккаунтом. Войдите через Telegram или используйте привязанный логин.")
)

type webAccount struct {
	UserID      int64
	Login       string
	PasswordSet bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	LastLoginAt time.Time
}

type webProgressSummary struct {
	Source   string `json:"source"`
	Label    string `json:"label"`
	XP       int    `json:"xp"`
	Level    string `json:"level"`
	Lessons  int    `json:"lessons"`
	Practice int    `json:"practice"`
	Words    int    `json:"words"`
	Mistakes int    `json:"mistakes"`
	HasData  bool   `json:"has_data"`
}

type webProgressMergeChallenge struct {
	Required bool               `json:"required"`
	Web      webProgressSummary `json:"web"`
	Telegram webProgressSummary `json:"telegram"`
}

type webAccountStore interface {
	createWebAccount(userID int64, login string, passwordHash string) (webAccount, error)
	getWebAccountByLogin(login string) (webAccount, string, bool, error)
	getWebAccountByUserID(userID int64) (webAccount, bool, error)
	touchWebAccountLogin(userID int64) error
	updateWebAccountPassword(userID int64, passwordHash string) error
	completeWebAccountProfile(userID int64, login string, passwordHash string) (webAccount, error)
}

type telegramWebProfile struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date"`
}

type webTelegramAccountStore interface {
	webAccountStore
	authenticateTelegramWebAccount(currentUserID int64, profile telegramWebProfile) (webAccount, userState, error)
	authenticateTelegramWebAccountWithChoice(currentUserID int64, profile telegramWebProfile, mergeChoice string) (webAccount, userState, error)
	telegramWebProgressMergeChallenge(currentUserID int64, profile telegramWebProfile) (webProgressMergeChallenge, error)
}

type webAuthBroker struct {
	mu        sync.Mutex
	requests  map[string]webAuthRequest
	pendingBy map[int64]string
}

type webAuthRequest struct {
	Token         string
	Code          string
	CurrentUserID int64
	CreatedAt     time.Time
	ExpiresAt     time.Time
	Profile       telegramWebProfile
	Completed     bool
}

func newWebAuthBroker() *webAuthBroker {
	return &webAuthBroker{
		requests:  map[string]webAuthRequest{},
		pendingBy: map[int64]string{},
	}
}

func (b *webAuthBroker) create(currentUserID int64, now time.Time) (webAuthRequest, error) {
	if b == nil {
		return webAuthRequest{}, errors.New("web auth broker is not configured")
	}
	token, err := newWebAuthToken()
	if err != nil {
		return webAuthRequest{}, err
	}
	code, err := newWebAuthCode()
	if err != nil {
		return webAuthRequest{}, err
	}
	request := webAuthRequest{
		Token:         token,
		Code:          code,
		CurrentUserID: currentUserID,
		CreatedAt:     now.UTC(),
		ExpiresAt:     now.UTC().Add(10 * time.Minute),
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.sweepLocked(now)
	b.requests[token] = request
	return request, nil
}

func (b *webAuthBroker) beginTelegram(token string, profile telegramWebProfile) (webAuthRequest, error) {
	if b == nil {
		return webAuthRequest{}, errors.New("web auth broker is not configured")
	}
	token = strings.TrimSpace(token)
	now := time.Now().UTC()
	b.mu.Lock()
	defer b.mu.Unlock()
	request, ok := b.requests[token]
	if !ok {
		return webAuthRequest{}, errors.New("запрос входа не найден или устарел")
	}
	if now.After(request.ExpiresAt) {
		delete(b.requests, token)
		return webAuthRequest{}, errors.New("срок действия входа истёк")
	}
	if profile.ID <= 0 {
		return webAuthRequest{}, errors.New("telegram id is required")
	}
	request.Profile = profile
	request.Profile.AuthDate = now.Unix()
	b.requests[token] = request
	b.pendingBy[profile.ID] = token
	return request, nil
}

func (b *webAuthBroker) verifySiteCode(token string, code string, now time.Time) (webAuthRequest, error) {
	if b == nil {
		return webAuthRequest{}, errors.New("web auth broker is not configured")
	}
	token = strings.TrimSpace(token)
	code = strings.TrimSpace(code)
	b.mu.Lock()
	defer b.mu.Unlock()
	request, ok := b.requests[token]
	if !ok {
		return webAuthRequest{}, errors.New("запрос входа не найден или устарел")
	}
	if now.UTC().After(request.ExpiresAt) {
		delete(b.requests, token)
		if request.Profile.ID != 0 {
			delete(b.pendingBy, request.Profile.ID)
		}
		return webAuthRequest{}, errors.New("срок действия входа истёк")
	}
	if request.Profile.ID == 0 {
		return webAuthRequest{}, errors.New("сначала откройте Telegram и получите код")
	}
	if request.Code != code {
		return webAuthRequest{}, errors.New("код не совпал")
	}
	request.Completed = true
	request.Profile.AuthDate = now.UTC().Unix()
	delete(b.requests, token)
	delete(b.pendingBy, request.Profile.ID)
	return request, nil
}

func (b *webAuthBroker) hasPending(telegramID int64) bool {
	if b == nil || telegramID == 0 {
		return false
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	token := b.pendingBy[telegramID]
	if token == "" {
		return false
	}
	request, ok := b.requests[token]
	if !ok || time.Now().UTC().After(request.ExpiresAt) {
		delete(b.requests, token)
		delete(b.pendingBy, telegramID)
		return false
	}
	return true
}

func (b *webAuthBroker) consume(token string, now time.Time) (webAuthRequest, bool, error) {
	if b == nil {
		return webAuthRequest{}, false, errors.New("web auth broker is not configured")
	}
	token = strings.TrimSpace(token)
	b.mu.Lock()
	defer b.mu.Unlock()
	request, ok := b.requests[token]
	if !ok {
		return webAuthRequest{}, false, errors.New("запрос входа не найден или устарел")
	}
	if now.UTC().After(request.ExpiresAt) {
		delete(b.requests, token)
		return webAuthRequest{}, false, errors.New("срок действия входа истёк")
	}
	if !request.Completed {
		return request, false, nil
	}
	delete(b.requests, token)
	return request, true, nil
}

func (b *webAuthBroker) sweepLocked(now time.Time) {
	for token, request := range b.requests {
		if now.UTC().After(request.ExpiresAt) {
			delete(b.requests, token)
			if request.Profile.ID != 0 {
				delete(b.pendingBy, request.Profile.ID)
			}
		}
	}
}

func newWebAuthToken() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw[:]), nil
}

func newWebAuthCode() (string, error) {
	var builder strings.Builder
	builder.Grow(6)
	max := big.NewInt(10)
	for builder.Len() < 6 {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		builder.WriteByte(byte('0' + n.Int64()))
	}
	return builder.String(), nil
}

type webReferralResult = referralApplication

type webReferralStore interface {
	getUserByReferralCode(code string) (userState, bool, error)
	applyWebReferral(newUserID int64, code string) (webReferralResult, error)
}

func normalizeWebLogin(login string) (string, error) {
	login = strings.ToLower(strings.TrimSpace(login))
	if len(login) < 3 || len(login) > 32 {
		return "", errors.New("Логин должен быть от 3 до 32 символов")
	}
	for _, ch := range login {
		ok := ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '_' || ch == '-' || ch == '.'
		if !ok {
			return "", errors.New("Логин может содержать только латинские буквы, цифры, точку, дефис и подчёркивание")
		}
	}
	if login[0] == '.' || login[0] == '-' || login[0] == '_' {
		return "", errors.New("Логин должен начинаться с буквы или цифры")
	}
	return login, nil
}

func validateWebPassword(password string) error {
	if len(password) < 8 {
		return errors.New("Пароль должен быть не короче 8 символов")
	}
	if len(password) > 128 {
		return errors.New("Пароль слишком длинный")
	}
	return nil
}

func validateWebPasswordConfirmation(password string, confirmation string) error {
	if confirmation == "" {
		return errors.New("Подтвердите пароль")
	}
	if password != confirmation {
		return errors.New("Пароли не совпадают")
	}
	return nil
}

func normalizeWebReferralCode(code string) (string, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, "_", "")
	code = strings.ReplaceAll(code, " ", "")
	if code == "" {
		return "", nil
	}
	if len(code) != webReferralCodeLength {
		return "", errors.New("Реферальный код должен состоять из 8 символов")
	}
	for _, ch := range code {
		if !strings.ContainsRune(webReferralAlphabet, ch) {
			return "", errors.New("Реферальный код может содержать только латинские буквы и цифры")
		}
	}
	return code, nil
}

func newWebReferralCode() (string, error) {
	var builder strings.Builder
	builder.Grow(webReferralCodeLength)
	max := big.NewInt(int64(len(webReferralAlphabet)))
	for builder.Len() < webReferralCodeLength {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		builder.WriteByte(webReferralAlphabet[index.Int64()])
	}
	return builder.String(), nil
}

func hashWebPassword(password string) (string, error) {
	salt := make([]byte, webPasswordSaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := pbkdf2SHA256([]byte(password), salt, webPasswordIterations, webPasswordKeyBytes)
	return fmt.Sprintf("%s$%d$%s$%s",
		webPasswordHashAlgorithm,
		webPasswordIterations,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func verifyWebPassword(password string, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 || parts[0] != webPasswordHashAlgorithm {
		return false
	}
	iterations, err := strconv.Atoi(parts[1])
	if err != nil || iterations <= 0 {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil || len(salt) == 0 {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil || len(expected) == 0 {
		return false
	}
	actual := pbkdf2SHA256([]byte(password), salt, iterations, len(expected))
	return subtle.ConstantTimeCompare(actual, expected) == 1
}

func isWebPasswordHash(encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 4 || parts[0] != webPasswordHashAlgorithm {
		return false
	}
	iterations, err := strconv.Atoi(parts[1])
	if err != nil || iterations <= 0 {
		return false
	}
	if _, err := base64.RawStdEncoding.DecodeString(parts[2]); err != nil {
		return false
	}
	if _, err := base64.RawStdEncoding.DecodeString(parts[3]); err != nil {
		return false
	}
	return true
}

func pbkdf2SHA256(password []byte, salt []byte, iterations int, keyLen int) []byte {
	hashLen := sha256.Size
	blocks := (keyLen + hashLen - 1) / hashLen
	key := make([]byte, 0, blocks*hashLen)
	for block := 1; block <= blocks; block++ {
		mac := hmac.New(sha256.New, password)
		_, _ = mac.Write(salt)
		_, _ = mac.Write([]byte{byte(block >> 24), byte(block >> 16), byte(block >> 8), byte(block)})
		u := mac.Sum(nil)
		t := append([]byte(nil), u...)
		for i := 1; i < iterations; i++ {
			mac = hmac.New(sha256.New, password)
			_, _ = mac.Write(u)
			u = mac.Sum(nil)
			for j := range t {
				t[j] ^= u[j]
			}
		}
		key = append(key, t...)
	}
	return key[:keyLen]
}
