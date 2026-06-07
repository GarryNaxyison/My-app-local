#!/usr/bin/env bash
set -euo pipefail

if [[ "${EUID:-$(id -u)}" -ne 0 ]]; then
	echo "poliglot-keys must be run as root" >&2
	exit 1
fi

BOT_DIR="${BOT_DIR:-/opt/aibot}"
BOT_USER="${BOT_USER:-botservice}"
MASTER_FILE="${ACTIVATION_KEYS_FILE:-$BOT_DIR/activation_keys.txt}"
LOCK_FILE="$BOT_DIR/activation_keys.lock"
ALPHABET='ABCDEFGHJKLMNPQRSTUVWXYZ23456789'

usage() {
  cat <<'USAGE'
Usage:
  poliglot-keys <count> [month|year|30|365] [premium|platinum] [note]

Examples:
  poliglot-keys 50 month
  poliglot-keys 20 year premium launch-may
  poliglot-keys 10 month platinum
  poliglot-keys 5 year platinum vip

The command creates a batch file near the bot:
  /opt/aibot/activation_keys_YYYYMMDD_HHMMSS.txt

It also appends the same keys to:
  /opt/aibot/activation_keys.txt
USAGE
}

count="${1:-}"
period="${2:-month}"
tier="${3:-premium}"
note="${4:-}"

if [[ -z "$count" || "$count" == "-h" || "$count" == "--help" ]]; then
  usage
  exit 0
fi

if ! [[ "$count" =~ ^[0-9]+$ ]] || (( count < 1 || count > 100000 )); then
  echo "count must be a number from 1 to 100000" >&2
  exit 1
fi

case "${period,,}" in
  month|monthly|30|30d|1m) days=30 ;;
  year|yearly|365|365d|1y) days=365 ;;
  *) echo "period must be month/year or 30/365" >&2; exit 1 ;;
esac

case "${tier,,}" in
  platinum) tier="platinum" ;;
  premium|"") tier="premium" ;;
  *) echo "tier must be premium or platinum" >&2; exit 1 ;;
esac

mkdir -p "$BOT_DIR"
touch "$MASTER_FILE"
chmod 640 "$MASTER_FILE" 2>/dev/null || true

timestamp="$(date +%Y%m%d_%H%M%S)"
batch_file="$BOT_DIR/activation_keys_${timestamp}.txt"

make_key() {
	local compact
	compact=""
	while (( ${#compact} < 16 )); do
		local need chunk
		need=$((16 - ${#compact}))
		chunk="$(LC_ALL=C tr -dc "$ALPHABET" < /dev/urandom | head -c "$need" || true)"
		compact+="$chunk"
	done
	printf '%s-%s-%s-%s' "${compact:0:4}" "${compact:4:4}" "${compact:8:4}" "${compact:12:4}"
}

(
  flock -x 9
  generated=0
  while (( generated < count )); do
    key="$(make_key)"
    if grep -qE "^${key}([[:space:]]|$)" "$MASTER_FILE" "$BOT_DIR"/activation_keys_*.txt 2>/dev/null; then
      continue
    fi
    line="$key $days $tier"
    if [[ -n "$note" ]]; then
      line="$line $note"
    fi
    printf '%s\n' "$line" >> "$MASTER_FILE"
    printf '%s\n' "$line" >> "$batch_file"
    generated=$((generated + 1))
  done
) 9>"$LOCK_FILE"

chmod 640 "$batch_file" "$MASTER_FILE" 2>/dev/null || true
if id "$BOT_USER" >/dev/null 2>&1; then
	chown "$BOT_USER:$BOT_USER" "$batch_file" "$MASTER_FILE" 2>/dev/null || true
fi

echo "Generated $count keys"
echo "Batch file: $batch_file"
echo "Import file: $MASTER_FILE"
