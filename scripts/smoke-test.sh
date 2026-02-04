#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://localhost:3000"
CREATE_USER="true"
RENEW_MEMBERSHIP="true"
RENEW_MONTHS="1"
PLAN="standard"
EMAIL=""
PASSWORD="password123"
NAME="Test User"
PHONE="08123456789"
ADDRESS="Jakarta"

section() {
  echo
  echo "=== $1 ==="
}

to_bool() {
  # accepts: true/false, 1/0, yes/no, y/n (case-insensitive)
  local v
  v="${1:-}"
  v="$(printf '%s' "$v" | tr '[:upper:]' '[:lower:]')"
  case "$v" in
    true|1|yes|y) echo "true";;
    false|0|no|n) echo "false";;
    *) echo "Invalid value for --$2: '$1'. Use true/false or 1/0." >&2; exit 2;;
  esac
}

json_get() {
  local path="$1"
  local auth_header="${2:-}"
  if [[ -n "$auth_header" ]]; then
    curl -fsS "$BASE_URL$path" -H "Authorization: Bearer $auth_header"
  else
    curl -fsS "$BASE_URL$path"
  fi
}

json_post() {
  local path="$1"
  local body="$2"
  local auth_header="${3:-}"
  if [[ -n "$auth_header" ]]; then
    curl -fsS "$BASE_URL$path" \
      -H "Authorization: Bearer $auth_header" \
      -H "Content-Type: application/json" \
      -d "$body"
  else
    curl -fsS "$BASE_URL$path" \
      -H "Content-Type: application/json" \
      -d "$body"
  fi
}

json_put() {
  local path="$1"
  local body="$2"
  local auth_header="${3:-}"
  if [[ -n "$auth_header" ]]; then
    curl -fsS -X PUT "$BASE_URL$path" \
      -H "Authorization: Bearer $auth_header" \
      -H "Content-Type: application/json" \
      -d "$body"
  else
    curl -fsS -X PUT "$BASE_URL$path" \
      -H "Content-Type: application/json" \
      -d "$body"
  fi
}

# Parse JSON using python3 (avoids jq dependency)
json_extract() {
  local expr="$1"
  python3 - "$expr" <<'PY'
import json,sys
expr=sys.argv[1]
data=json.load(sys.stdin)

# very small extractor: supports dot paths like data.access_token
cur=data
for part in expr.split('.'):
    if part=='':
        continue
    if isinstance(cur, dict) and part in cur:
        cur=cur[part]
    else:
        cur=None
        break
if cur is None:
    sys.exit(3)
if isinstance(cur,(dict,list)):
    print(json.dumps(cur))
else:
    print(cur)
PY
}

usage() {
  cat <<EOF
Usage: ./scripts/smoke-test.sh [options]

Options:
  --base-url URL              (default: $BASE_URL)
  --create-user true|false    (default: $CREATE_USER)
  --renew-membership true|false (default: $RENEW_MEMBERSHIP)
  --renew-months N            (default: $RENEW_MONTHS)
  --plan NAME                 (default: $PLAN)
  --email EMAIL               (default: empty; auto-generated if create-user=true)
  --password PASS             (default: $PASSWORD)
  --name NAME                 (default: $NAME)
  --phone PHONE               (default: $PHONE)
  --address ADDRESS           (default: $ADDRESS)
  -h, --help                  Show help

Examples:
  ./scripts/smoke-test.sh
  ./scripts/smoke-test.sh --renew-membership 0
  ./scripts/smoke-test.sh --create-user 0 --email you@example.com --password password123
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --base-url) BASE_URL="$2"; shift 2;;
    --create-user) CREATE_USER="$2"; shift 2;;
    --renew-membership) RENEW_MEMBERSHIP="$2"; shift 2;;
    --renew-months) RENEW_MONTHS="$2"; shift 2;;
    --plan) PLAN="$2"; shift 2;;
    --email) EMAIL="$2"; shift 2;;
    --password) PASSWORD="$2"; shift 2;;
    --name) NAME="$2"; shift 2;;
    --phone) PHONE="$2"; shift 2;;
    --address) ADDRESS="$2"; shift 2;;
    -h|--help) usage; exit 0;;
    *) echo "Unknown argument: $1" >&2; usage; exit 2;;
  esac
done

CREATE_USER="$(to_bool "$CREATE_USER" create-user)"
RENEW_MEMBERSHIP="$(to_bool "$RENEW_MEMBERSHIP" renew-membership)"

if [[ "$CREATE_USER" == "true" && -z "$EMAIL" ]]; then
  stamp="$(date +%Y%m%d%H%M%S)"
  EMAIL="test+$stamp@example.com"
fi

if [[ "$CREATE_USER" == "false" && -z "$EMAIL" ]]; then
  echo "--email is required when --create-user false (so the script can login)." >&2
  exit 2
fi

section "Health"
json_get "/health" | python3 -m json.tool

if [[ "$CREATE_USER" == "true" ]]; then
  section "Create User"
  create_body=$(python3 - <<PY
import json
print(json.dumps({
  "name": "${NAME}",
  "email": "${EMAIL}",
  "password": "${PASSWORD}",
  "phone": "${PHONE}",
  "address": "${ADDRESS}"
}))
PY
)
  set +e
  resp=$(json_post "/api/v1/users" "$create_body" 2>&1)
  status=$?
  set -e
  if [[ $status -eq 0 ]]; then
    echo "$resp" | python3 -m json.tool
  else
    echo "Create user failed (maybe already exists). Continuing..." >&2
    echo "$resp" >&2
  fi
fi

section "Login"
login_body=$(python3 - <<PY
import json
print(json.dumps({"email":"${EMAIL}","password":"${PASSWORD}"}))
PY
)
login_resp=$(json_post "/api/v1/auth/login" "$login_body")
# Extract token (expects: { data: { access_token: "..." } })
TOKEN=$(echo "$login_resp" | json_extract "data.access_token")
if [[ -z "$TOKEN" ]]; then
  echo "Login succeeded but access_token is missing." >&2
  exit 1
fi
printf 'Token prefix: %.20s...\n' "$TOKEN"

section "Me"
json_get "/api/v1/me" "$TOKEN" | python3 -m json.tool

section "Membership (before renew)"
json_get "/api/v1/membership" "$TOKEN" | python3 -m json.tool

if [[ "$RENEW_MEMBERSHIP" == "true" ]]; then
  section "Membership Renew"
  renew_body=$(python3 - <<PY
import json
print(json.dumps({"months": int("${RENEW_MONTHS}"), "plan": "${PLAN}"}))
PY
)
  json_post "/api/v1/membership/renew" "$renew_body" "$TOKEN" | python3 -m json.tool
fi

section "QR Code"
set +e
qr_resp=$(json_get "/api/v1/qr/code" "$TOKEN" 2>&1)
qr_status=$?
set -e
if [[ $qr_status -eq 0 ]]; then
  echo "$qr_resp" | python3 -m json.tool
  QR_TOKEN=$(echo "$qr_resp" | json_extract "data.token" || true)
  if [[ -n "${QR_TOKEN:-}" ]]; then
    section "QR Scan (gate)"
    scan_body=$(python3 - <<PY
import json
print(json.dumps({"token":"${QR_TOKEN}"}))
PY
)
    set +e
    scan_resp=$(json_post "/api/v1/qr/scan" "$scan_body" 2>&1)
    scan_status=$?
    set -e
    if [[ $scan_status -eq 0 ]]; then
      echo "$scan_resp" | python3 -m json.tool
    else
      echo "QR scan failed (may be expected if membership expired)." >&2
      echo "$scan_resp" >&2
    fi
  else
    echo "Skipping QR scan because no QR token was returned." >&2
  fi
else
  echo "QR code request failed. This is expected if membership is not ACTIVE." >&2
  echo "$qr_resp" >&2
fi

section "Attendance History"
set +e
att_resp=$(json_get "/api/v1/attendance/history?limit=10" "$TOKEN" 2>&1)
att_status=$?
set -e
if [[ $att_status -eq 0 ]]; then
  echo "$att_resp" | python3 -m json.tool
else
  echo "Attendance history failed (continuing)." >&2
  echo "$att_resp" >&2
fi

section "Workout Progress (create + list)"
workout_body=$(python3 - <<PY
import json
print(json.dumps({"title":"Push day","notes":"Bench + triceps"}))
PY
)
json_post "/api/v1/workouts/progress" "$workout_body" "$TOKEN" | python3 -m json.tool
json_get "/api/v1/workouts/progress?limit=5" "$TOKEN" | python3 -m json.tool

section "Settings (get + update)"
json_get "/api/v1/settings" "$TOKEN" | python3 -m json.tool
settings_body=$(python3 - <<PY
import json
print(json.dumps({"push_enabled": False}))
PY
)
json_put "/api/v1/settings" "$settings_body" "$TOKEN" | python3 -m json.tool || true

section "Notifications"
set +e
notif_resp=$(json_get "/api/v1/notifications?limit=5" "$TOKEN" 2>&1)
notif_status=$?
set -e
if [[ $notif_status -eq 0 ]]; then
  echo "$notif_resp" | python3 -m json.tool
else
  echo "Notifications failed (continuing)." >&2
  echo "$notif_resp" >&2
fi

section "Program/PT (may be empty without seed)"
set +e
templates_resp=$(json_get "/api/v1/templates/followed" "$TOKEN" 2>&1)
templates_status=$?
set -e
if [[ $templates_status -eq 0 ]]; then
  echo "$templates_resp" | python3 -m json.tool
else
  echo "Templates followed failed (continuing)." >&2
  echo "$templates_resp" >&2
fi

set +e
targets_resp=$(json_get "/api/v1/targets?period=weekly" "$TOKEN" 2>&1)
targets_status=$?
set -e
if [[ $targets_status -eq 0 ]]; then
  echo "$targets_resp" | python3 -m json.tool
else
  echo "Targets list failed (continuing)." >&2
  echo "$targets_resp" >&2
fi

set +e
pt_resp=$(json_get "/api/v1/pt" "$TOKEN" 2>&1)
pt_status=$?
set -e
if [[ $pt_status -eq 0 ]]; then
  echo "$pt_resp" | python3 -m json.tool
else
  echo "PT list failed (continuing)." >&2
  echo "$pt_resp" >&2
fi

echo
echo "Smoke test completed."
echo "User email: $EMAIL"
