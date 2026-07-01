# Web API And Site App

The bot process now also serves a browser app and JSON API for the website.

## URLs

- `GET /app` - primary React V2 web interface for learners.
- `GET /app/v2` - direct React V2 alias for QA and old links.
- `GET /app/v1` - legacy V1 web interface.
- `GET /app?shell=mobile` / `GET /app?shell=desktop` - force the mobile or desktop web-app shell for QA; without the query the server chooses by user agent.
- `GET /app/mobile` / `GET /app/desktop` - same forced shell modes under the `/app/*` reverse-proxy path.
- `GET /api/health` - JSON health check.
- `GET /api/session` - creates/returns the current web session.
- `POST /api/settings` - update interface language, learning language, or CEFR level.
- `GET /api/leaderboard?language=en` - website-only top 10 overall and for one learning language.
- `POST /api/lesson/start` - create a new AI lesson.
- `POST /api/lesson/answer` - check the current lesson answer.
- `POST /api/practice` - send one free-practice message.
- `POST /api/level-test/start` - start the 36-question placement test.
- `POST /api/level-test/answer` - submit one placement-test answer.
- `POST /api/words/next` - get the next vocabulary question.
- `POST /api/words/answer` - submit a vocabulary answer.
- `POST /api/words/pronunciation` - returns generated `audio/mpeg` pronunciation for the active learning-language word.
- `GET /api/vocabulary?page=0` - paginated learned-word vocabulary for the current learning language.
- `GET /api/phrasebook` - saved Phrasebook entries for the current learner.
- `POST /api/phrasebook` - save or update one Phrasebook entry from lesson/practice/roleplay/manual input.
- `DELETE /api/phrasebook?id=...` - remove one Phrasebook entry.
- `POST /api/word-game/next` and `POST /api/word-game/answer` - review words that were added from word learning.
- `POST /api/spelling/start` and `POST /api/spelling/answer` - spelling practice for learned words.
- `GET /api/mistakes` - paginated mistake dictionary.
- `GET /api/progress` - current progress and limits.
- `GET /api/premium/plans` - available Premium plans.
- `POST /api/premium/payment` - create a YooKassa payment and return `confirmation_url`.

All web sessions use an HttpOnly signed cookie. The same SQLite store is used for Telegram and website users; web users receive negative internal IDs so they do not collide with Telegram IDs. Website leaderboards are filtered through `web_accounts`, so Telegram users stay in the Telegram bot leaderboards and login/password users stay in the site leaderboard.

## Environment

```env
WEBHOOK_LISTEN_ADDR=:8080
WEB_API_SESSION_SECRET=change_me_to_a_long_random_secret_for_web_sessions
WEB_CORS_ORIGINS=https://neriva.ru,https://www.neriva.ru,https://api.neriva.ru,https://poliglotai.ru,https://www.poliglotai.ru,https://api.poliglotai.ru,https://poliglotai.online,https://www.poliglotai.online,https://api.poliglotai.online
WEB_COOKIE_SECURE=true
WEB_COOKIE_SAMESITE=lax
WEB_YOOKASSA_SHOP_ID=1359880
WEB_YOOKASSA_SECRET_KEY=site_shop_secret_key
WEB_PAYMENT_RETURN_URL=https://neriva.ru/app?payment=success
```

If the site is on another origin, add it to `WEB_CORS_ORIGINS`.
If the site is on a different registrable domain than `api.neriva.ru`, use `WEB_COOKIE_SAMESITE=none` together with HTTPS and `WEB_COOKIE_SECURE=true`.
If `WEB_YOOKASSA_SHOP_ID` and `WEB_YOOKASSA_SECRET_KEY` are not set, website payments fall back to the Telegram YooKassa shop.
For the separate website shop, set YooKassa notifications to `https://api.neriva.ru/yookassa/webhook/`.

## Caddy

```caddyfile
neriva.ru, www.neriva.ru, poliglotai.ru, www.poliglotai.ru, poliglotai.online, www.poliglotai.online {
    @webapp path /app /app/*
    handle @webapp {
        reverse_proxy 127.0.0.1:8080
    }

    @botapi path /api/*
    handle @botapi {
        reverse_proxy 127.0.0.1:8080
    }

    @botcallbacks path /payment/success /payment/success* /tonapi/webhook /tonapi/webhook/* /yookassa/webhook /yookassa/webhook/*
    handle @botcallbacks {
        reverse_proxy 127.0.0.1:8080
    }

    handle {
        root * /var/www/poliglotai
        try_files {path} /poliglot-ai.html
        file_server
    }
}

api.neriva.ru, api.poliglotai.ru, api.poliglotai.online {
    reverse_proxy 127.0.0.1:8080
}
```

Open the app at `https://neriva.ru/app`, or use the old Poliglot hosts while they remain accepted aliases.
