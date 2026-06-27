package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type config struct {
	TelegramBotToken                  string
	TelegramOpsRecipients             []telegramOpsRecipient
	OpenRouterAPIKey                  string
	OpenRouterModel                   string
	OpenRouterTranslatorModel         string
	OpenRouterVocabularyModel         string
	OpenRouterSTTModel                string
	OpenRouterPronunciationModel      string
	OpenRouterPronunciationAudioModel string
	OpenRouterTTSModel                string
	OpenRouterTTSVoice                string
	OpenRouterAppURL                  string
	OpenRouterAppName                 string
	DataPath                          string
	DatabasePath                      string
	AITutorDatabasePath               string
	ActivationKeysDatabasePath        string
	ActivationKeysImportPath          string
	VocabularyDatabasePath            string
	StorageDriver                     string
	PremiumStarsPrice                 int
	PremiumRubPrice                   int
	PremiumYearStarsPrice             int
	PremiumYearRubPrice               int
	PlatinumStarsPrice                int
	PlatinumRubPrice                  int
	PlatinumYearStarsPrice            int
	PlatinumYearRubPrice              int
	MaxVoiceSeconds                   int
	YooKassaShopID                    string
	YooKassaSecretKey                 string
	WebYooKassaShopID                 string
	WebYooKassaSecretKey              string
	YooKassaReturnURL                 string
	YooKassaWebhookKey                string
	RollyPayBotCashboxID              string
	RollyPayBotAPIKey                 string
	RollyPayWebCashboxID              string
	RollyPayWebAPIKey                 string
	RollyPayAPIBaseURL                string
	RollyPayCreatePaymentPath         string
	RollyPayWebhookSecret             string
	WebhookListenAddr                 string
	WebAPISessionSecret               string
	WebCORSOrigins                    []string
	WebCookieDomain                   string
	WebCookieSecure                   bool
	WebCookieSameSite                 string
	WebAppURL                         string
	WebTelegramLoginBot               string
	WebPaymentReturnURL               string
	WebTurnstileSiteKey               string
	WebTurnstileSecretKey             string
	TrustedProxyCIDRs                 []string
	CryptoTONWallet                   string
	CryptoTONMonthAmount              string
	CryptoTONYearAmount               string
	CryptoTONPlatinumMonthAmount      string
	CryptoTONPlatinumYearAmount       string
	CryptoTONCenterAPIKey             string
	CryptoTONCenterBaseURL            string
	CryptoTONAPIKey                   string
	CryptoTONAPIBaseURL               string
	CryptoTONAPIWebhookKey            string
	CryptoUSDTMonthAmount             string
	CryptoUSDTYearAmount              string
	CryptoUSDTPlatinumMonthAmount     string
	CryptoUSDTPlatinumYearAmount      string
	CryptoUSDTRubRate                 int
	CryptoUSDTTONWallet               string
	CryptoUSDTTONJettonMaster         string
	CryptoUSDTTRC20Wallet             string
	CryptoUSDTTRC20Contract           string
	CryptoTronGridAPIKey              string
	CryptoTronGridBaseURL             string
	CryptoPaymentTTLMinutes           int
	ReminderHour                      int
	ReminderUTCOffset                 int
	SQLiteBackupDir                   string
	SQLiteBackupHour                  int
	SQLiteBackupRetention             int
	MemoryLimitMB                     int
}

func configFromEnv() (config, error) {
	cfg := config{
		TelegramBotToken:                  strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")),
		TelegramOpsRecipients:             parseTelegramOpsRecipients(os.Getenv("TELEGRAM_OPS_RECIPIENTS")),
		OpenRouterAPIKey:                  strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY")),
		OpenRouterModel:                   envOrDefault("OPENROUTER_MODEL", "google/gemini-3.1-flash"),
		OpenRouterTranslatorModel:         envOrDefault("OPENROUTER_TRANSLATOR_MODEL", envOrDefault("OPENROUTER_MODEL", "google/gemini-3.1-flash")),
		OpenRouterVocabularyModel:         envOrDefault("OPENROUTER_VOCABULARY_MODEL", "google/gemini-3.1-flash-lite"),
		OpenRouterSTTModel:                envOrDefault("OPENROUTER_STT_MODEL", "openai/gpt-4o-transcribe"),
		OpenRouterPronunciationModel:      envOrDefault("OPENROUTER_PRONUNCIATION_MODEL", "google/gemini-3.1-flash-lite"),
		OpenRouterPronunciationAudioModel: envOrDefault("OPENROUTER_PRONUNCIATION_AUDIO_MODEL", "google/gemini-2.5-pro"),
		OpenRouterTTSModel:                normalizeOpenRouterTTSModel(envOrDefault("OPENROUTER_TTS_MODEL", "google/gemini-3.1-flash-tts-preview")),
		OpenRouterTTSVoice:                envOrDefault("OPENROUTER_TTS_VOICE", "Kore"),
		OpenRouterAppURL:                  envOrDefault("OPENROUTER_APP_URL", "http://localhost"),
		OpenRouterAppName:                 envOrDefault("OPENROUTER_APP_NAME", "AI Polyglot Coach"),
		DataPath:                          envOrDefault("DATA_PATH", "english_coach_data.json"),
		DatabasePath:                      envOrDefault("DATABASE_PATH", "english_coach.sqlite"),
		AITutorDatabasePath:               envOrDefault("AI_TUTOR_DATABASE_PATH", "ai_tutor_lessons.sqlite"),
		ActivationKeysDatabasePath:        envOrDefault("ACTIVATION_KEYS_DATABASE_PATH", "activation_keys.sqlite"),
		ActivationKeysImportPath:          envOrDefault("ACTIVATION_KEYS_FILE", "activation_keys.txt"),
		VocabularyDatabasePath:            envOrDefault("VOCABULARY_DATABASE_PATH", "vocabulary.sqlite"),
		StorageDriver:                     envOrDefault("STORAGE_DRIVER", "sqlite"),
		PremiumStarsPrice:                 envIntOrDefault("PREMIUM_STARS_PRICE", 150),
		PremiumRubPrice:                   envIntOrDefault("PREMIUM_RUB_PRICE", 300),
		PremiumYearStarsPrice:             envIntOrDefault("PREMIUM_YEAR_STARS_PRICE", 1500),
		PremiumYearRubPrice:               envIntOrDefault("PREMIUM_YEAR_RUB_PRICE", 3000),
		PlatinumStarsPrice:                envIntOrDefault("PLATINUM_STARS_PRICE", 300),
		PlatinumRubPrice:                  envIntOrDefault("PLATINUM_RUB_PRICE", 590),
		PlatinumYearStarsPrice:            envIntOrDefault("PLATINUM_YEAR_STARS_PRICE", 3000),
		PlatinumYearRubPrice:              envIntOrDefault("PLATINUM_YEAR_RUB_PRICE", 5900),
		MaxVoiceSeconds:                   envIntOrDefault("MAX_VOICE_SECONDS", 30),
		YooKassaShopID:                    strings.TrimSpace(os.Getenv("YOOKASSA_SHOP_ID")),
		YooKassaSecretKey:                 strings.TrimSpace(os.Getenv("YOOKASSA_SECRET_KEY")),
		WebYooKassaShopID:                 strings.TrimSpace(os.Getenv("WEB_YOOKASSA_SHOP_ID")),
		WebYooKassaSecretKey:              strings.TrimSpace(os.Getenv("WEB_YOOKASSA_SECRET_KEY")),
		YooKassaReturnURL:                 strings.TrimSpace(os.Getenv("YOOKASSA_RETURN_URL")),
		YooKassaWebhookKey:                strings.TrimSpace(os.Getenv("YOOKASSA_WEBHOOK_KEY")),
		RollyPayBotCashboxID:              strings.TrimSpace(os.Getenv("ROLLYPAY_BOT_CASHBOX_ID")),
		RollyPayBotAPIKey:                 strings.TrimSpace(os.Getenv("ROLLYPAY_BOT_API_KEY")),
		RollyPayWebCashboxID:              strings.TrimSpace(os.Getenv("ROLLYPAY_WEB_CASHBOX_ID")),
		RollyPayWebAPIKey:                 strings.TrimSpace(os.Getenv("ROLLYPAY_WEB_API_KEY")),
		RollyPayAPIBaseURL:                envOrDefault("ROLLYPAY_API_BASE_URL", "https://rollypay.io"),
		RollyPayCreatePaymentPath:         envOrDefault("ROLLYPAY_CREATE_PAYMENT_PATH", "/api/v1/payments"),
		RollyPayWebhookSecret:             strings.TrimSpace(os.Getenv("ROLLYPAY_WEBHOOK_SECRET")),
		WebhookListenAddr:                 envOrDefault("WEBHOOK_LISTEN_ADDR", ":8080"),
		WebAPISessionSecret:               strings.TrimSpace(os.Getenv("WEB_API_SESSION_SECRET")),
		WebCORSOrigins:                    envListOrDefault("WEB_CORS_ORIGINS", []string{"https://poliglot.ai", "https://www.poliglot.ai", "https://poliglotai.ru", "https://www.poliglotai.ru", "https://poliglotai.online", "https://www.poliglotai.online"}),
		WebCookieDomain:                   strings.TrimSpace(os.Getenv("WEB_COOKIE_DOMAIN")),
		WebCookieSecure:                   envBoolOrDefault("WEB_COOKIE_SECURE", false),
		WebCookieSameSite:                 strings.ToLower(envOrDefault("WEB_COOKIE_SAMESITE", "lax")),
		WebAppURL:                         envOrDefault("WEB_APP_URL", "https://poliglotai.ru/app"),
		WebTelegramLoginBot:               strings.TrimPrefix(envOrDefault("WEB_TELEGRAM_LOGIN_BOT", "Poliglot_AI_bot"), "@"),
		WebPaymentReturnURL:               strings.TrimSpace(os.Getenv("WEB_PAYMENT_RETURN_URL")),
		WebTurnstileSiteKey:               strings.TrimSpace(os.Getenv("WEB_TURNSTILE_SITE_KEY")),
		WebTurnstileSecretKey:             strings.TrimSpace(os.Getenv("WEB_TURNSTILE_SECRET_KEY")),
		TrustedProxyCIDRs:                 envListOrDefault("TRUSTED_PROXY_CIDRS", nil),
		CryptoTONWallet:                   strings.TrimSpace(os.Getenv("CRYPTO_TON_WALLET")),
		CryptoTONMonthAmount:              strings.TrimSpace(os.Getenv("CRYPTO_TON_MONTH_AMOUNT")),
		CryptoTONYearAmount:               strings.TrimSpace(os.Getenv("CRYPTO_TON_YEAR_AMOUNT")),
		CryptoTONPlatinumMonthAmount:      strings.TrimSpace(os.Getenv("CRYPTO_TON_PLATINUM_MONTH_AMOUNT")),
		CryptoTONPlatinumYearAmount:       strings.TrimSpace(os.Getenv("CRYPTO_TON_PLATINUM_YEAR_AMOUNT")),
		CryptoTONCenterAPIKey:             strings.TrimSpace(os.Getenv("CRYPTO_TONCENTER_API_KEY")),
		CryptoTONCenterBaseURL:            envOrDefault("CRYPTO_TONCENTER_BASE_URL", "https://toncenter.com/api/v2"),
		CryptoTONAPIKey:                   strings.TrimSpace(os.Getenv("CRYPTO_TONAPI_KEY")),
		CryptoTONAPIBaseURL:               envOrDefault("CRYPTO_TONAPI_BASE_URL", "https://tonapi.io/v2"),
		CryptoTONAPIWebhookKey:            strings.TrimSpace(os.Getenv("CRYPTO_TONAPI_WEBHOOK_KEY")),
		CryptoUSDTMonthAmount:             strings.TrimSpace(os.Getenv("CRYPTO_USDT_MONTH_AMOUNT")),
		CryptoUSDTYearAmount:              strings.TrimSpace(os.Getenv("CRYPTO_USDT_YEAR_AMOUNT")),
		CryptoUSDTPlatinumMonthAmount:     strings.TrimSpace(os.Getenv("CRYPTO_USDT_PLATINUM_MONTH_AMOUNT")),
		CryptoUSDTPlatinumYearAmount:      strings.TrimSpace(os.Getenv("CRYPTO_USDT_PLATINUM_YEAR_AMOUNT")),
		CryptoUSDTRubRate:                 envIntOrDefault("CRYPTO_USDT_RUB_RATE", 72),
		CryptoUSDTTONWallet:               strings.TrimSpace(os.Getenv("CRYPTO_USDT_TON_WALLET")),
		CryptoUSDTTONJettonMaster:         envOrDefault("CRYPTO_USDT_TON_JETTON_MASTER", "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"),
		CryptoUSDTTRC20Wallet:             strings.TrimSpace(os.Getenv("CRYPTO_USDT_TRC20_WALLET")),
		CryptoUSDTTRC20Contract:           envOrDefault("CRYPTO_USDT_TRC20_CONTRACT", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"),
		CryptoTronGridAPIKey:              strings.TrimSpace(os.Getenv("CRYPTO_TRONGRID_API_KEY")),
		CryptoTronGridBaseURL:             envOrDefault("CRYPTO_TRONGRID_BASE_URL", "https://api.trongrid.io"),
		CryptoPaymentTTLMinutes:           envIntOrDefault("CRYPTO_PAYMENT_TTL_MINUTES", 60),
		ReminderHour:                      envIntOrDefault("REMINDER_HOUR", 19),
		ReminderUTCOffset:                 envIntOrDefault("REMINDER_UTC_OFFSET_MINUTES", 180),
		SQLiteBackupDir:                   envOrDefault("SQLITE_BACKUP_DIR", "backups"),
		SQLiteBackupHour:                  envHourOrDefault("SQLITE_BACKUP_HOUR", 3),
		SQLiteBackupRetention:             envIntOrDefault("SQLITE_BACKUP_RETENTION_DAYS", 14),
		MemoryLimitMB:                     envIntOrDefault("MEMORY_LIMIT_MB", 0),
	}
	cfg.OpenRouterTTSVoice = normalizeOpenRouterTTSVoice(cfg.OpenRouterTTSModel, cfg.OpenRouterTTSVoice)

	var missing []string
	if cfg.TelegramBotToken == "" {
		missing = append(missing, "TELEGRAM_BOT_TOKEN")
	}
	if cfg.OpenRouterAPIKey == "" {
		missing = append(missing, "OPENROUTER_API_KEY")
	}
	if len(missing) > 0 {
		return config{}, fmt.Errorf("missing required environment variable(s): %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

type telegramOpsRecipient struct {
	ChatID string
}

func (cfg config) telegramOpsRecipients() []telegramOpsRecipient {
	recipients := make([]telegramOpsRecipient, 0, len(cfg.TelegramOpsRecipients)+len(telegramDefaultOpsRecipientIDs))
	seen := map[string]bool{}
	add := func(recipient telegramOpsRecipient) {
		chatID := normalizeTelegramOpsChatID(recipient.ChatID)
		if chatID == "" || seen[strings.ToLower(chatID)] {
			return
		}
		seen[strings.ToLower(chatID)] = true
		recipients = append(recipients, telegramOpsRecipient{ChatID: chatID})
	}
	for _, recipient := range cfg.TelegramOpsRecipients {
		add(recipient)
	}
	for _, recipientID := range telegramDefaultOpsRecipientIDs {
		add(telegramOpsRecipient{ChatID: strconv.FormatInt(recipientID, 10)})
	}
	return recipients
}

func parseTelegramOpsRecipients(raw string) []telegramOpsRecipient {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n' || r == '\r' || r == '\t' || r == ' '
	})
	recipients := make([]telegramOpsRecipient, 0, len(parts))
	seen := map[string]bool{}
	for _, part := range parts {
		chatID := normalizeTelegramOpsChatID(part)
		if chatID == "" || seen[strings.ToLower(chatID)] {
			continue
		}
		seen[strings.ToLower(chatID)] = true
		recipients = append(recipients, telegramOpsRecipient{ChatID: chatID})
	}
	return recipients
}

func normalizeTelegramOpsChatID(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "https://t.me/")
	value = strings.TrimPrefix(value, "http://t.me/")
	value = strings.TrimPrefix(value, "t.me/")
	value = strings.Trim(value, "/")
	if value == "" {
		return ""
	}
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return value
	}
	if !strings.HasPrefix(value, "@") {
		value = "@" + value
	}
	return value
}

func (cfg config) yooKassaEnabled() bool {
	return cfg.YooKassaShopID != "" && cfg.YooKassaSecretKey != "" && cfg.YooKassaReturnURL != ""
}

func (cfg config) webYooKassaOverride() bool {
	return cfg.WebYooKassaShopID != "" || cfg.WebYooKassaSecretKey != ""
}

func (cfg config) webYooKassaShopID() string {
	if cfg.webYooKassaOverride() {
		return cfg.WebYooKassaShopID
	}
	return cfg.YooKassaShopID
}

func (cfg config) webYooKassaSecretKey() string {
	if cfg.webYooKassaOverride() {
		return cfg.WebYooKassaSecretKey
	}
	return cfg.YooKassaSecretKey
}

func (cfg config) webPaymentReturnURL() string {
	if cfg.WebPaymentReturnURL != "" {
		return cfg.WebPaymentReturnURL
	}
	return cfg.YooKassaReturnURL
}

func (cfg config) webYooKassaEnabled() bool {
	return cfg.webYooKassaShopID() != "" && cfg.webYooKassaSecretKey() != "" && cfg.webPaymentReturnURL() != ""
}

func (cfg config) rollyPayBotEnabled() bool {
	return strings.TrimSpace(cfg.RollyPayBotCashboxID) != "" && strings.TrimSpace(cfg.RollyPayBotAPIKey) != ""
}

func (cfg config) rollyPayWebEnabled() bool {
	return strings.TrimSpace(cfg.RollyPayWebCashboxID) != "" && strings.TrimSpace(cfg.RollyPayWebAPIKey) != ""
}

func (cfg config) rollyPayBotReturnURL() string {
	if cfg.WebAppURL != "" {
		return cfg.WebAppURL
	}
	return "https://poliglotai.online/app"
}

func (cfg config) rollyPayWebReturnURL() string {
	if cfg.WebPaymentReturnURL != "" {
		return cfg.WebPaymentReturnURL
	}
	return "https://poliglotai.online/app?payment=success&provider=rollypay"
}

func envOrDefault(name string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	return value
}

func envListOrDefault(name string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return append([]string{}, fallback...)
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return append([]string{}, fallback...)
	}
	return result
}

func envBoolOrDefault(name string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(name)))
	if value == "" {
		return fallback
	}
	switch value {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}

func normalizeOpenRouterTTSModel(model string) string {
	model = strings.TrimSpace(model)
	if model == "openai/gpt-4o-mini-tts" || model == "openai/gpt-4o-mini-tts-2025-12-15" {
		return "google/gemini-3.1-flash-tts-preview"
	}
	return model
}

func normalizeOpenRouterTTSVoice(model string, voice string) string {
	voice = strings.TrimSpace(voice)
	if voice == "" {
		voice = "Kore"
	}
	if strings.Contains(strings.ToLower(strings.TrimSpace(model)), "gemini") && strings.EqualFold(voice, "alloy") {
		return "Kore"
	}
	return voice
}

func envIntOrDefault(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func envHourOrDefault(name string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 || parsed > 23 {
		return fallback
	}
	return parsed
}

func requireNonEmpty(name string, value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(name + " must not be empty")
	}
	return nil
}
