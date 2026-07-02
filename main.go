package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

func main() {
	if err := loadDotEnv(".env"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to load .env: %v", err)
	}

	cfg, err := configFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.MemoryLimitMB > 0 {
		debug.SetMemoryLimit(int64(cfg.MemoryLimitMB) * 1024 * 1024)
		debug.SetGCPercent(75)
		log.Printf("Go memory limit set to %d MB", cfg.MemoryLimitMB)
	}

	httpClient := &http.Client{Timeout: 60 * time.Second}
	store, err := newStore(cfg)
	if err != nil {
		log.Fatalf("failed to open store: %v", err)
	}
	activationKeys, err := newActivationKeyStore(cfg.ActivationKeysDatabasePath, cfg.ActivationKeysImportPath)
	if err != nil {
		log.Fatalf("failed to open activation key store: %v", err)
	}
	defer activationKeys.close()
	vocabularyDB, err := initializeVocabularyRuntime(cfg)
	if err != nil {
		log.Fatalf("failed to prepare vocabulary store: %v", err)
	}
	defer vocabularyDB.Close()

	bot := &bot{
		cfg:            cfg,
		store:          store,
		telegram:       newTelegramClient(cfg.TelegramBotToken, httpClient),
		cryptoRates:    newCryptoRateProvider(cfg, httpClient),
		activationKeys: activationKeys,
		openrouter: newOpenRouterClient(
			cfg.OpenRouterAPIKey,
			cfg.OpenRouterModel,
			cfg.OpenRouterAppURL,
			cfg.OpenRouterAppName,
			httpClient,
		),
	}
	if cfg.yooKassaEnabled() {
		bot.yookassa = newYooKassaClient(cfg.YooKassaShopID, cfg.YooKassaSecretKey, httpClient)
	}
	if cfg.webYooKassaEnabled() {
		bot.webYooKassa = newYooKassaClient(cfg.webYooKassaShopID(), cfg.webYooKassaSecretKey(), httpClient)
	}
	if cfg.rollyPayBotEnabled() {
		bot.rollyPayBot = newRollyPayClient(cfg.RollyPayBotCashboxID, cfg.RollyPayBotAPIKey, cfg.RollyPayAPIBaseURL, cfg.RollyPayCreatePaymentPath, httpClient)
	}
	if cfg.rollyPayWebEnabled() {
		bot.rollyPayWeb = newRollyPayClient(cfg.RollyPayWebCashboxID, cfg.RollyPayWebAPIKey, cfg.RollyPayAPIBaseURL, cfg.RollyPayCreatePaymentPath, httpClient)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := newWebhookHTTPServer(cfg, bot)
	go func() {
		log.Printf("HTTP server listening on %s", cfg.WebhookListenAddr)
		if bot.yookassa == nil {
			log.Printf("Telegram YooKassa is disabled: set YOOKASSA_SHOP_ID, YOOKASSA_SECRET_KEY and YOOKASSA_RETURN_URL")
		}
		if bot.webYooKassa == nil {
			log.Printf("Website YooKassa is disabled: set WEB_YOOKASSA_SHOP_ID, WEB_YOOKASSA_SECRET_KEY and WEB_PAYMENT_RETURN_URL")
		}
		if cfg.YooKassaWebhookKey == "" {
			log.Printf("YooKassa webhook is disabled: set YOOKASSA_WEBHOOK_KEY")
		}
		if bot.rollyPayBot == nil {
			log.Printf("Telegram RollyPay is disabled: set ROLLYPAY_BOT_CASHBOX_ID and ROLLYPAY_BOT_API_KEY")
		}
		if bot.rollyPayWeb == nil {
			log.Printf("Website RollyPay is disabled: set ROLLYPAY_WEB_CASHBOX_ID and ROLLYPAY_WEB_API_KEY")
		}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server failed: %v", err)
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	log.Printf("NERIVA bot started with model %s (translator %s)", cfg.OpenRouterModel, cfg.OpenRouterTranslatorModel)
	go bot.runReminderScheduler(ctx)
	go bot.runSQLiteBackupScheduler(ctx)
	if err := bot.run(ctx); err != nil && ctx.Err() == nil {
		log.Fatal(err)
	}
}

func yookassaWebhookMux(webhookKey string, bot *bot) http.Handler {
	mux := http.NewServeMux()
	web := newWebAPI(bot.cfg, bot)
	web.register(mux)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/payment/success", paymentSuccessHandler)
	webYooKassaWebhook := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/yookassa/webhook/" && r.URL.Path != "/yookassa/webhook" {
			http.NotFound(w, r)
			return
		}
		bot.handleWebYooKassaWebhook(w, r)
	}
	mux.HandleFunc("/yookassa/webhook/", webYooKassaWebhook)
	mux.HandleFunc("/yookassa/webhook", webYooKassaWebhook)
	if webhookKey != "" {
		mux.HandleFunc("/yookassa/webhook/"+webhookKey, bot.handleYooKassaWebhook)
	}
	rollyPayWebhook := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rollypay/webhook/site":
			bot.handleRollyPayWebhook(w, r, rollyPayChannelWeb)
		case "/rollypay/webhook/bot":
			bot.handleRollyPayWebhook(w, r, rollyPayChannelTelegram)
		default:
			http.NotFound(w, r)
		}
	}
	mux.HandleFunc("/rollypay/webhook/site", rollyPayWebhook)
	mux.HandleFunc("/rollypay/webhook/bot", rollyPayWebhook)
	tonAPIWebhook := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tonapi/webhook" && !strings.HasPrefix(r.URL.Path, "/tonapi/webhook/") {
			http.NotFound(w, r)
			return
		}
		bot.handleTONAPIWebhook(w, r)
	}
	mux.HandleFunc("/tonapi/webhook", tonAPIWebhook)
	mux.HandleFunc("/tonapi/webhook/", tonAPIWebhook)
	return web.withCORS(mux)
}

func newWebhookHTTPServer(cfg config, bot *bot) *http.Server {
	return &http.Server{
		Addr:              cfg.WebhookListenAddr,
		Handler:           yookassaWebhookMux(cfg.YooKassaWebhookKey, bot),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

func paymentSuccessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<!doctype html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="description" content="Страница подтверждения оплаты Premium через YooKassa">
  <title>Оплата принята</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f6f7fb;
      --card: #ffffff;
      --text: #151922;
      --muted: #4b5567;
      --border: #e4e7ef;
      --primary: #2f6fed;
      --primary-dark: #2458c7;
      --ok: #17a673;
    }

    * { box-sizing: border-box; }

    body {
      margin: 0;
      min-height: 100vh;
      font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      background:
        radial-gradient(circle at top left, rgba(47, 111, 237, .12), transparent 30%),
        linear-gradient(180deg, #fbfcff 0%, var(--bg) 100%);
      color: var(--text);
    }

    main {
      min-height: 100vh;
      display: grid;
      place-items: center;
      padding: 24px;
    }

    section {
      width: min(100%, 520px);
      background: var(--card);
      border: 1px solid var(--border);
      border-radius: 18px;
      padding: 30px;
      box-shadow: 0 18px 45px rgba(20, 30, 55, .10);
    }

    .status {
      width: 52px;
      height: 52px;
      display: grid;
      place-items: center;
      border-radius: 50%;
      background: rgba(23, 166, 115, .12);
      color: var(--ok);
      font-size: 28px;
      margin-bottom: 18px;
    }

    h1 {
      margin: 0 0 12px;
      font-size: 28px;
      line-height: 1.15;
      letter-spacing: 0;
    }

    p {
      margin: 0 0 22px;
      color: var(--muted);
      font-size: 16px;
      line-height: 1.55;
    }

    .note {
      padding: 12px 14px;
      border-radius: 12px;
      background: #f4f7ff;
      border: 1px solid #dfe8ff;
      margin-bottom: 22px;
      color: #31415f;
      font-size: 14px;
      line-height: 1.45;
    }

    .actions {
      display: flex;
      gap: 10px;
      flex-wrap: wrap;
    }

    a {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 46px;
      padding: 0 18px;
      border-radius: 12px;
      background: var(--primary);
      color: #fff;
      text-decoration: none;
      font-weight: 700;
      transition: background .15s ease, transform .15s ease;
    }

    a:hover { background: var(--primary-dark); }
    a:active { transform: translateY(1px); }

    @media (max-width: 420px) {
      section { padding: 24px; border-radius: 16px; }
      h1 { font-size: 25px; }
      a { width: 100%; }
      .actions { display: grid; }
    }
  </style>
</head>
<body>
  <main>
    <section>
      <div class="status" aria-hidden="true">✓</div>
      <h1>Оплата принята</h1>
      <p>Premium активируется автоматически после подтверждения YooKassa. Вернитесь в приложение или Telegram-бота, статус обновится после webhook-подтверждения.</p>
      <div class="note">Если статус не изменился сразу, подождите несколько секунд: YooKassa может отправить подтверждение с небольшой задержкой.</div>
      <div class="actions">
        <a href="/app?payment=success">Открыть web-приложение</a>
        <a href="https://t.me/EnglishRobotAI_bot">Вернуться в Telegram</a>
      </div>
    </section>
  </main>
</body>
</html>`))
}
