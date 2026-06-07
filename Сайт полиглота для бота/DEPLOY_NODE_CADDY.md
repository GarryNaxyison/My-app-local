# Деплой Poliglot AI на Ubuntu + Caddy + Node.js

## Куда класть файлы

Папка сайта на сервере:

```bash
/var/www/poliglotai.ru
```

В эту папку нужно загрузить:

```text
server.mjs
package.json
.env
bot.html
poliglot-ai.html
privacy.html
terms.html
assets/
deploy/
```

PHP не нужен. Бот на сайте работает напрямую через `server.mjs` и OpenRouter.

## .env

Создайте файл:

```bash
sudo nano /var/www/poliglotai.ru/.env
```

Пример содержимого:

```env
OPENROUTER_API_KEY=sk-or-v1-your-openrouter-key
OPENROUTER_MODEL=openrouter/auto
PORT=3000
SITE_URL=https://poliglotai.ru
APP_TITLE=Poliglot AI
```

Файл `.env` нельзя класть в публичный репозиторий и нельзя отдавать как статику. `server.mjs` его не отдает наружу.

## Установка Node.js

На чистой Ubuntu:

```bash
sudo apt update
sudo apt install -y curl ca-certificates
curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -
sudo apt install -y nodejs
node -v
```

## Проверка вручную

```bash
cd /var/www/poliglotai.ru
node server.mjs
```

В другом окне:

```bash
curl http://127.0.0.1:3000/bot
curl -X POST http://127.0.0.1:3000/api/poliglot-chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Привет, проверь связь","mode":"dialog","interfaceLanguage":"Русский","learningLanguage":"English","history":[]}'
```

## systemd-сервис

```bash
sudo cp /var/www/poliglotai.ru/deploy/poliglot-ai.service.example /etc/systemd/system/poliglot-ai.service
sudo systemctl daemon-reload
sudo systemctl enable --now poliglot-ai
sudo systemctl status poliglot-ai
```

Логи:

```bash
journalctl -u poliglot-ai -f
```

## Caddy

Если Caddy уже обслуживает `poliglotai.ru`, добавьте в ваш существующий `/etc/caddy/Caddyfile` прокси на Node.

Минимальный вариант:

```caddyfile
poliglotai.ru, poliglotai.online {
    encode zstd gzip
    reverse_proxy 127.0.0.1:3000
}

www.poliglotai.ru {
    redir https://poliglotai.ru{uri} permanent
}

www.poliglotai.online {
    redir https://poliglotai.online{uri} permanent
}
```

Проверка и перезагрузка:

```bash
sudo caddy validate --config /etc/caddy/Caddyfile
sudo systemctl reload caddy
```

После этого открывайте:

```text
https://poliglotai.ru/bot
```
