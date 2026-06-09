# Telegram AI English Coach

Poliglot AI — Telegram-бот и web app для изучения языков с AI-уроками, практикой, голосом, фото-переводом, словарем, ошибками, прогрессом и тарифами Free, Premium и Platinum.

## Возможности

- `/lesson` - короткое задание по английскому.
- `/practice` - свободная разговорная практика.
- `/progress` - статистика занятий.
- `/limits` - остаток лимитов на сегодня.
- `/premium` - описание тарифов Free, Premium и Platinum.
- `/buy` - покупка Premium или Platinum через Telegram Stars, YooKassa/СБП, TON или USDT.
- `/invite` - персональная ссылка для приглашения друга.
- `/stop` - остановить текущий режим.
- `/resetmode` - выход из текущего режима.
- Основные действия доступны через кнопки Telegram-меню: web app, уроки, практика, Premium/Platinum, лимиты, инструменты и остановка режима.
- В Premium-меню есть inline-кнопки покупки, приглашения друга, лимитов и возврата назад.
- OpenRouter API для выбора модели под бюджет.
- Без внешних Go-зависимостей: только стандартная библиотека.

В режиме `/lesson` бот дает одно короткое задание, проверяет один ответ и автоматически завершает урок. Для длинной беседы используется `/practice`.

## Монетизация

Тарифы подаются как простая лестница ценности: Free показывает пользу бота, Premium закрывает ежедневное обучение, Platinum продает максимальные лимиты для интенсивной подготовки.

| Тариф | Лимиты и доступ | Месяц | Год |
| --- | --- | --- | --- |
| Free | 5 уроков, 15 сообщений практики, без голосовых функций | 0 ₽ | 0 ₽ |
| Premium | 50 уроков, 200 сообщений практики, 20 голосовых до 30 секунд, голос в текст, перевод услышанного, перевод текста с картинки | 300 ₽ / 150 Stars / около 4.15 USDT | 3000 ₽ / 1500 Stars / около 41.50 USDT |
| Platinum | 100 уроков, 500 сообщений практики, 60 голосовых до 30 секунд, максимальный доступ к AI-диалогам, voice/photo-практике и переводчику | 590 ₽ / 300 Stars / около 8.15 USDT | 5900 ₽ / 3000 Stars / около 81.50 USDT |

Месячная цена Premium указана с учетом стартовой скидки 70%. Оплата доступна в рублях через YooKassa/СБП, Telegram Stars, TON и USDT. Годовая оплата включена для Premium и Platinum.

Цена задается в `.env`:

```env
PREMIUM_STARS_PRICE=150
PREMIUM_YEAR_STARS_PRICE=1500
PREMIUM_RUB_PRICE=300
PREMIUM_YEAR_RUB_PRICE=3000
PLATINUM_STARS_PRICE=300
PLATINUM_YEAR_STARS_PRICE=3000
PLATINUM_RUB_PRICE=590
PLATINUM_YEAR_RUB_PRICE=5900
CRYPTO_USDT_MONTH_AMOUNT=
CRYPTO_USDT_YEAR_AMOUNT=
CRYPTO_USDT_PLATINUM_MONTH_AMOUNT=
CRYPTO_USDT_PLATINUM_YEAR_AMOUNT=
CRYPTO_USDT_RUB_RATE=72
STORAGE_DRIVER=sqlite
DATABASE_PATH=english_coach.sqlite
VOCABULARY_DIR=data/vocabulary
VOCABULARY_DATABASE_PATH=vocabulary.sqlite
ACTIVATION_KEYS_DATABASE_PATH=activation_keys.sqlite
ACTIVATION_KEYS_FILE=activation_keys.txt
MEMORY_LIMIT_MB=512
OPENROUTER_STT_MODEL=openai/gpt-4o-transcribe
OPENROUTER_PRONUNCIATION_MODEL=google/gemini-3.1-flash-lite
OPENROUTER_TTS_MODEL=google/gemini-3.1-flash-tts-preview
OPENROUTER_TTS_VOICE=Kore
MAX_VOICE_SECONDS=30
TELEGRAM_OPS_RECIPIENTS=185156683,@your_admin_username
```

If the explicit `CRYPTO_USDT_*_AMOUNT` values are empty, the bot now asks TonAPI `/v2/rates` for the current USDT/RUB rate using `CRYPTO_TONAPI_BASE_URL`, `CRYPTO_TONAPI_KEY`, and `CRYPTO_USDT_TON_JETTON_MASTER`. `CRYPTO_USDT_RUB_RATE` stays as a fallback when TonAPI is unavailable.

`TELEGRAM_OPS_RECIPIENTS` controls extra owner/admin Telegram notifications, including web bug reports and payment alerts. Use a comma-separated list of numeric chat IDs and optional `@username` channel/group targets. The built-in AsaselD chat ID (`185156683`) is always included and duplicated values are ignored.

### Serial Premium Keys

Activation keys are kept in a separate SQLite database, controlled by `ACTIVATION_KEYS_DATABASE_PATH`, so they do not mix with user progress in `DATABASE_PATH`. The app imports available keys from `ACTIVATION_KEYS_FILE` on startup and before redemption. File lines use:

```text
XXXX-XXXX-XXXX-XXXX 30 premium optional-note
XXXX-XXXX-XXXX-XXXX 365 premium optional-note
```

The key alphabet avoids ambiguous characters (`I`, `O`, `0`, `1`). A redeemed key is deleted from the active key table and written to `activation_key_redemptions`, so it cannot be reused even if the same line remains in the import file.

Generate keys locally or on the server:

```bash
node tools/generate_activation_keys.mjs --count 100 --period month --out activation_keys.txt
node tools/generate_activation_keys.mjs --count 25 --period year --out activation_keys.txt
```

On Ubuntu production the helper can be installed as the root-only `poliglot-keys` command:

```bash
poliglot-keys 100 month
poliglot-keys 25 year
```

It writes `XXXX-XXXX-XXXX-XXXX` keys to a timestamped batch file next to the bot and appends them to `/opt/aibot/activation_keys.txt` for import.

### Web Auth Protection

Set Cloudflare Turnstile keys to add a “not a robot” check to web login, registration, and Telegram profile completion:

```env
WEB_TURNSTILE_SITE_KEY=your_site_key
WEB_TURNSTILE_SECRET_KEY=your_secret_key
```

The web API also applies per-IP and per-session request limits, and the browser briefly locks action buttons after a click to reduce accidental double-submits and button spam.

### Low-memory server mode

For a 1 GB VPS, keep multilingual vocabulary JSON files on disk instead of embedding them into the binary. The app imports those JSON files into indexed SQLite vocabulary tables in `VOCABULARY_DATABASE_PATH` on startup, then uses SQLite for hot lookups. JSON files remain read-only seed files for import/rebuilds.

Recommended Linux `.env` values:

```env
VOCABULARY_DIR=/opt/aibot
MEMORY_LIMIT_MB=512
```

Recommended production build:

```bash
go build -ldflags "-s -w" -o aibot .
```

Deploy the `data/vocabulary/vocabulary_words*.json` files together with the `aibot` binary in `VOCABULARY_DIR`.

On the first production start, allow extra time for the SQLite vocabulary import. Later starts skip the import while the JSON file size and imported SQLite word count still match; mtime-only changes are ignored so routine deploys do not trigger expensive reimports. The vocabulary database is separate from `DATABASE_PATH`, so user-data backups stay small and the static dictionary can be regenerated.

AI-filled vocabulary is mutable runtime data. Preserve `VOCABULARY_DATABASE_PATH` (`vocabulary.sqlite` by default) during deploys and include it in production backups, or export the `vocabulary_ai_words` and `vocabulary_ai_translations` tables before replacing the database. These tables prevent repeated AI token spending after restart, JSON reimport, or deploy.

По умолчанию данные хранятся в SQLite-файле `english_coach.sqlite`. Если рядом есть старый `english_coach_data.json`, пустая SQLite-база импортирует его при первом запуске.

Для YooKassa лучше указывать страницу возврата:

```env
YOOKASSA_RETURN_URL=https://poliglotai.ru/payment/success
WEB_YOOKASSA_SHOP_ID=1359880
WEB_YOOKASSA_SECRET_KEY=site_shop_secret_key
WEB_PAYMENT_RETURN_URL=https://poliglotai.ru/app?payment=success
```

Caddy должен прокидывать `/payment/success`, `/app`, `/api` и webhook на бота. Для отдельного магазина сайта URL уведомлений в YooKassa: `https://api.poliglotai.ru/yookassa/webhook/`.

## Быстрый старт

1. Создайте Telegram-бота через `@BotFather` и получите `TELEGRAM_BOT_TOKEN`.
2. Создайте ключ OpenRouter: https://openrouter.ai/keys
3. Скопируйте `.env.example` в `.env` и заполните значения.
4. Запустите:

```powershell
go run .
```

On Windows, if `go`, `gofmt`, or a real `python.exe` are installed but missing from `PATH`, run:

```powershell
npm run doctor:tools
```

To persist found tool directories in the User PATH:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File tools/enable-dev-toolchain.ps1 -Persist
```

## Выбор модели

Модель меняется через `OPENROUTER_MODEL`. По умолчанию используется `google/gemini-3.1-flash-lite`, а позже можно выбрать другую в каталоге OpenRouter.

Voice recognition uses `OPENROUTER_STT_MODEL=openai/gpt-4o-transcribe` through OpenRouter. Pronunciation coaching uses `OPENROUTER_PRONUNCIATION_MODEL=google/gemini-3.1-flash-lite`: the model receives only transcript/confidence/fluency JSON and returns learner-friendly feedback, not raw audio.
Text-to-speech uses `OPENROUTER_TTS_MODEL=google/gemini-3.1-flash-tts-preview` and `OPENROUTER_TTS_VOICE=Kore`; the backend wraps OpenRouter PCM output into WAV for Telegram and web playback.

## Как можно развивать до подписки

- Бесплатный лимит AI-ответов в день.
- Pro-режим: больше практики, словарь, повторение ошибок, голосовые ответы.
- Telegram Mini App: красивый прогресс, оплата, личный кабинет.
