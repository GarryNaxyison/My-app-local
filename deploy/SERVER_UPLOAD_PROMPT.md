# Server Deploy Upload Prompt

Use this as a persistent instruction for future Codex deploys.

## Сохранение прогресса — обязательное правило от 2026-11-14

Не сбрасывать прогресс существующих пользователей. Для обновления React-интерфейса публиковать только актуальные JS/CSS, service worker и HTML через постоянный upload-сервис. Сначала размещать ассеты, затем атомарно заменять HTML. Не менять backend binary, `.env`, базы SQLite, словари и изображения; не перезапускать сервис без отдельной необходимости. До и после сравнивать записи прогресса в режиме чтения и проверять неизменность бинарного файла и процесса. Резервные копии не создавать.

Приведённый ниже старый полный деплой содержит историческое описание миграции со сбросом прогресса. Это не разрешение выполнять такую миграцию. Скрипт полного деплоя нельзя использовать для обычного обновления интерфейса.

```text
Для деплоя NERIVA на Ubuntu-сервер используй постоянный deploy-upload сервис, а не временный upload route.

Сервер: root@186.246.45.123.
Upload endpoint: https://poliglotai.online/__codex_deploy_upload/<filename> (use the `.online` host; large PUTs to `.ru` can end with an incomplete upload/Broken pipe).
Systemd service: poliglot-deploy-upload.service
Service files: /opt/aibot/deploy-upload/upload_server.py and /opt/aibot/deploy-upload/deploy-upload.env
Token storage: /opt/aibot/deploy-upload/deploy-upload.env, key DEPLOY_UPLOAD_TOKEN. Не печатай токен в ответ пользователю.

Перед upload:
1. Собери backend Linux binary и web archive. `poliglot-app-web.tgz` должен включать только `web/`; словари JSON больше не входят в пакет.
2. Получи token через SSH, сохрани его только в переменную текущего процесса и не выводи:
   $token = (& tmp\putty\plink.exe -batch -hostkey "SHA256:SCw7oIAHX/xB9QJueAc9L7OvbXtDyVqA++ciDMxefI0" -pw <password> root@186.246.45.123 "awk -F= '/DEPLOY_UPLOAD_TOKEN/{print `$2}' /opt/aibot/deploy-upload/deploy-upload.env").Trim()
3. Загружай только allowlist-имена:
   - aibot-linux-amd64
   - poliglot-app-web.tgz
   - poliglot-public-site.tgz
   - Caddyfile.deploy
   - deploy_react_server.sh

Пример upload:
curl.exe --fail --show-error --silent -H "X-Deploy-Token: $token" --upload-file tmp\aibot-linux-amd64-current-web https://poliglotai.online/__codex_deploy_upload/aibot-linux-amd64
curl.exe --fail --show-error --silent -H "X-Deploy-Token: $token" --upload-file tmp\poliglot-app-web-current.tgz https://poliglotai.online/__codex_deploy_upload/poliglot-app-web.tgz
curl.exe --fail --show-error --silent -H "X-Deploy-Token: $token" --upload-file tmp\deploy_react_server.sh https://poliglotai.online/__codex_deploy_upload/deploy_react_server.sh

После upload запускай стандартный deploy:
tmp\putty\plink.exe -batch -hostkey "SHA256:SCw7oIAHX/xB9QJueAc9L7OvbXtDyVqA++ciDMxefI0" -pw <password> root@186.246.45.123 "chmod 700 /tmp/deploy_react_server.sh && bash /tmp/deploy_react_server.sh"

Deploy script must not create backups or reset user progress. Historical binaries included a destructive one-time progress migration: inspect migration behavior before any backend replacement. Do not execute that migration for an interface release. Preserve accounts, learning sessions, rewards, learned words, mistakes and user preferences.

После deploy всегда проверяй:
- https://neriva.ru/healthz
- https://neriva.ru/app/
- https://neriva.ru/neriva.html
- https://poliglotai.online/__codex_deploy_upload/healthz

Если upload endpoint не отвечает, проверь:
systemctl status poliglot-deploy-upload.service
systemctl status caddy
grep -n deployUpload /etc/caddy/Caddyfile
```
