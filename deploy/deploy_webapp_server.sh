#!/usr/bin/env bash
set -euo pipefail

APP_ROOT="/opt/aibot"
BIN_SRC="/tmp/aibot-linux-amd64"
BIN_DST="${APP_ROOT}/aibot"
WEB_ARCHIVE="/tmp/poliglot-app-web.tgz"
WEB_ROOT="${APP_ROOT}/web"
BACKUP_ROOT="${APP_ROOT}/deploy-backups"
STAMP="$(date -u +%Y%m%dT%H%M%SZ)-webapp"
BACKUP_DIR="${BACKUP_ROOT}/${STAMP}"
STAGING="$(mktemp -d)"
BIN_NEW="${BIN_DST}.new.${STAMP}"
BIN_OLD="${BIN_DST}.old.${STAMP}"
WEB_NEW="${WEB_ROOT}.new.${STAMP}"
WEB_OLD="${WEB_ROOT}.old.${STAMP}"

cleanup() {
  rm -rf "${STAGING}" "${BIN_NEW}" "${WEB_NEW}"
}

rollback() {
  set +e
  if [[ -f "${BIN_OLD}" ]]; then
    rm -f "${BIN_DST}"
    mv "${BIN_OLD}" "${BIN_DST}"
  fi
  if [[ -d "${WEB_OLD}" ]]; then
    rm -rf "${WEB_ROOT}"
    mv "${WEB_OLD}" "${WEB_ROOT}"
  fi
  if compgen -G "${BACKUP_DIR}/vocabulary_words*.json" >/dev/null; then
    cp -a "${BACKUP_DIR}"/vocabulary_words*.json "${APP_ROOT}/"
  fi
  systemctl restart aibot.service || true
}

on_error() {
  local rc=$?
  echo "web app deployment failed; attempting rollback" >&2
  rollback
  cleanup
  exit "${rc}"
}

trap cleanup EXIT
trap on_error ERR

test -s "${BIN_SRC}"
test -s "${WEB_ARCHIVE}"
mkdir -p "${APP_ROOT}" "${BACKUP_DIR}"

if [[ -f "${BIN_DST}" ]]; then
  cp -a "${BIN_DST}" "${BACKUP_DIR}/aibot"
fi
if [[ -d "${WEB_ROOT}" ]]; then
  tar -czf "${BACKUP_DIR}/web.tgz" -C "${APP_ROOT}" web
fi
find "${APP_ROOT}" -maxdepth 1 -type f -name 'vocabulary_words*.json' -exec cp -a {} "${BACKUP_DIR}/" \;

mkdir -p "${STAGING}/package"
tar -xzf "${WEB_ARCHIVE}" -C "${STAGING}/package"

WEB_SRC="${STAGING}/package/web"
VOCAB_SRC="${STAGING}/package/data/vocabulary"
test -f "${WEB_SRC}/index.html"
test -f "${WEB_SRC}/manifest.webmanifest"
test -d "${WEB_SRC}/assets"
test -d "${VOCAB_SRC}"
compgen -G "${VOCAB_SRC}/vocabulary_words*.json" >/dev/null

install -m 0755 -o root -g root "${BIN_SRC}" "${BIN_NEW}"
mkdir -p "${WEB_NEW}"
cp -a "${WEB_SRC}/." "${WEB_NEW}/"
chown -R root:root "${WEB_NEW}"
find "${WEB_NEW}" -type d -exec chmod 0755 {} +
find "${WEB_NEW}" -type f -exec chmod 0644 {} +

if [[ -f "${BIN_DST}" ]]; then
  mv "${BIN_DST}" "${BIN_OLD}"
fi
mv "${BIN_NEW}" "${BIN_DST}"
if [[ -d "${WEB_ROOT}" ]]; then
  mv "${WEB_ROOT}" "${WEB_OLD}"
fi
mv "${WEB_NEW}" "${WEB_ROOT}"
find "${VOCAB_SRC}" -maxdepth 1 -type f -name 'vocabulary_words*.json' -exec install -m 0644 -o root -g root {} "${APP_ROOT}/" \;

systemctl restart aibot.service
systemctl is-active --quiet aibot.service
curl -fsS --max-time 15 http://127.0.0.1:8080/healthz >/dev/null
curl -fsS --max-time 15 https://neriva.ru/healthz >/dev/null
curl -fsS --max-time 15 https://neriva.ru/app/ >/dev/null

rm -f "${BIN_OLD}"
rm -rf "${WEB_OLD}"
echo "NERIVA web app deployment complete"
echo "backup: ${BACKUP_DIR}"
