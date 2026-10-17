#!/usr/bin/env bash
set -euo pipefail

APP_ROOT="/opt/aibot"
BINARY="/tmp/aibot-linux-amd64"
WEB_ARCHIVE="/tmp/poliglot-app-web.tgz"

test -s "$BINARY"
test -s "$WEB_ARCHIVE"
tar -tzf "$WEB_ARCHIVE" | awk '
  $0 !~ /^web(\/|$)/ || $0 ~ /(^|\/)\.\.(\/|$)/ || $0 ~ /^\// { bad = 1 }
  END { exit bad }
'

install -o root -g root -m 0755 "$BINARY" "$APP_ROOT/aibot"
tar -xzf "$WEB_ARCHIVE" -C "$APP_ROOT"
chown -R root:root "$APP_ROOT/web"
find "$APP_ROOT/web" -type d -exec chmod 0755 {} +
find "$APP_ROOT/web" -type f -exec chmod 0644 {} +

upsert_model() {
  local key="$1"
  local value="openai/gpt-6-luna"
  if grep -q "^${key}=" "$APP_ROOT/.env"; then
    sed -i "s|^${key}=.*|${key}=${value}|" "$APP_ROOT/.env"
  else
    printf '%s=%s\n' "$key" "$value" >> "$APP_ROOT/.env"
  fi
}

upsert_model OPENROUTER_MODEL
upsert_model OPENROUTER_TRANSLATOR_MODEL
upsert_model OPENROUTER_VOCABULARY_MODEL
upsert_model OPENROUTER_PRONUNCIATION_MODEL
if grep -q '^OPENROUTER_PRONUNCIATION_AUDIO_MODEL=' "$APP_ROOT/.env"; then
  sed -i 's|^OPENROUTER_PRONUNCIATION_AUDIO_MODEL=.*|OPENROUTER_PRONUNCIATION_AUDIO_MODEL=google/gemini-3-flash-preview|' "$APP_ROOT/.env"
else
  printf '%s\n' 'OPENROUTER_PRONUNCIATION_AUDIO_MODEL=google/gemini-3-flash-preview' >> "$APP_ROOT/.env"
fi

WEBHOOK_ADDR="$(sed -n 's/^WEBHOOK_LISTEN_ADDR=//p' "$APP_ROOT/.env" | tail -n 1)"
WEBHOOK_PORT="${WEBHOOK_ADDR##*:}"
if [[ ! "$WEBHOOK_PORT" =~ ^[0-9]+$ ]]; then
  echo "WEBHOOK_LISTEN_ADDR does not contain a valid port" >&2
  exit 1
fi

systemctl restart aibot.service
for attempt in $(seq 1 120); do
  if systemctl is-active --quiet aibot.service && curl --fail --silent "http://127.0.0.1:${WEBHOOK_PORT}/healthz" >/dev/null; then
    break
  fi
  if [[ "$attempt" == 120 ]]; then
    systemctl --no-pager --full status aibot.service || true
    exit 1
  fi
  sleep 2
done

LEGACY_VOCAB_DIR="$APP_ROOT/data/vocabulary"
if [[ -d "$LEGACY_VOCAB_DIR" ]]; then
  while IFS= read -r -d '' old_file; do
    case "$old_file" in
      "$LEGACY_VOCAB_DIR"/vocabulary_words*.json) rm -f -- "$old_file" ;;
      *) echo "Unexpected file in legacy vocabulary directory; leaving it in place" >&2; exit 1 ;;
    esac
  done < <(find "$LEGACY_VOCAB_DIR" -maxdepth 1 -type f -print0)
  if [[ -z "$(find "$LEGACY_VOCAB_DIR" -mindepth 1 -maxdepth 1 -print -quit)" ]]; then
    rmdir -- "$LEGACY_VOCAB_DIR"
  fi
fi

echo "deployment complete"
