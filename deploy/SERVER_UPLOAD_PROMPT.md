# Server Deploy Upload Prompt

Use this as a persistent instruction for future Codex deploys.

```text
Для деплоя PoliglotAI на Ubuntu-сервер используй постоянный deploy-upload сервис, а не временный upload route.

Сервер: root@186.246.45.123.
Upload endpoint: https://poliglotai.ru/__codex_deploy_upload/<filename>
Systemd service: poliglot-deploy-upload.service
Service files: /opt/aibot/deploy-upload/upload_server.py and /opt/aibot/deploy-upload/deploy-upload.env
Token storage: /opt/aibot/deploy-upload/deploy-upload.env, key DEPLOY_UPLOAD_TOKEN. Не печатай токен в ответ пользователю.

Перед upload:
1. Собери backend Linux binary и web/public archives локально.
2. Получи token через SSH:
   $token = (& tmp\putty\plink.exe -batch -hostkey "SHA256:SCw7oIAHX/xB9QJueAc9L7OvbXtDyVqA++ciDMxefI0" -pw <password> root@186.246.45.123 "awk -F= '/DEPLOY_UPLOAD_TOKEN/{print `$2}' /opt/aibot/deploy-upload/deploy-upload.env").Trim()
3. Загружай только allowlist-имена:
   - aibot-linux-amd64
   - poliglot-app-web.tgz
   - poliglot-public-site.tgz
   - Caddyfile.deploy
   - deploy_react_server.sh

Пример upload:
curl.exe --fail --show-error --silent -H "X-Deploy-Token: $token" --upload-file tmp\aibot-linux-amd64-v2-auth-privacy-awards https://poliglotai.ru/__codex_deploy_upload/aibot-linux-amd64
curl.exe --fail --show-error --silent -H "X-Deploy-Token: $token" --upload-file tmp\poliglot-app-web-20260528-current.tgz https://poliglotai.ru/__codex_deploy_upload/poliglot-app-web.tgz
curl.exe --fail --show-error --silent -H "X-Deploy-Token: $token" --upload-file tmp\poliglot-public-site-20260528-current.tgz https://poliglotai.ru/__codex_deploy_upload/poliglot-public-site.tgz
curl.exe --fail --show-error --silent -H "X-Deploy-Token: $token" --upload-file tmp\Caddyfile.deploy.v2-auth-privacy-awards https://poliglotai.ru/__codex_deploy_upload/Caddyfile.deploy
curl.exe --fail --show-error --silent -H "X-Deploy-Token: $token" --upload-file tmp\deploy_react_server.sh https://poliglotai.ru/__codex_deploy_upload/deploy_react_server.sh

После upload запускай стандартный deploy:
tmp\putty\plink.exe -batch -hostkey "SHA256:SCw7oIAHX/xB9QJueAc9L7OvbXtDyVqA++ciDMxefI0" -pw <password> root@186.246.45.123 "chmod 700 /tmp/deploy_react_server.sh && bash /tmp/deploy_react_server.sh"

После deploy всегда проверяй:
- https://poliglotai.ru/healthz
- https://poliglotai.ru/app/v2/
- https://poliglotai.ru/poliglot-ai.html
- https://poliglotai.ru/__codex_deploy_upload/healthz

Если upload endpoint не отвечает, проверь:
systemctl status poliglot-deploy-upload.service
systemctl status caddy
grep -n deployUpload /etc/caddy/Caddyfile
```
