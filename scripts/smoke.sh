#!/usr/bin/env bash
# Smoke test against a running stack (docker compose up -d --build).
# Exits non-zero on any failure.
set -euo pipefail

BASE="${BASE:-http://localhost:8080}"
API="$BASE/api/v1"
fail() { echo "FAIL: $*" >&2; exit 1; }
expect() {
    local label="$1" code="$2" url="$3" body="${4:-}"
    local got
    got=$(curl -s -o /tmp/smoke.out -w "%{http_code}" -X POST -H 'Content-Type: application/json' "$url" ${body:+-d "$body"})
    [[ "$got" == "$code" ]] || fail "$label: want $code got $got — body: $(cat /tmp/smoke.out)"
    echo "PASS $label ($got)"
}

echo "==> healthz"
[[ "$(curl -s -o /dev/null -w "%{http_code}" "$API/healthz")" == "200" ]] || fail "healthz not 200"

echo "==> ua/infer"
expect "ua/infer ok" 200 "$API/ua/infer" '{"userAgent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"}'

echo "==> versions/lookup bad version"
expect "versions/lookup 400" 400 "$API/versions/lookup" '{"channel":"stable","version":"abc","platform":"darwin","architecture":"arm64"}'

echo "==> versions/lookup bad platform"
expect "versions/lookup 400 platform" 400 "$API/versions/lookup" '{"channel":"stable","version":"1.94.2","platform":"plan9","architecture":"arm64"}'

echo "==> releases ok"
code=$(curl -s -o /dev/null -w "%{http_code}" "$API/releases?channel=stable&page=1&pageSize=5")
[[ "$code" == "200" || "$code" == "502" ]] || fail "releases: want 200/502 got $code"
echo "PASS releases ($code)"

echo "==> extensions/search ok"
code=$(curl -s -o /dev/null -w "%{http_code}" "$API/extensions/search?q=python")
[[ "$code" == "200" || "$code" == "502" ]] || fail "extensions/search: want 200/502 got $code"
echo "PASS extensions/search ($code)"

echo "==> events accepted (degraded ok)"
expect "events 202" 202 "$API/events" '{"eventType":"SEARCH","targetType":"CLIENT","targetIdentifier":"1.94.2","platform":"darwin","architecture":"arm64","channel":"stable"}'

echo "==> trending empty/200"
code=$(curl -s -o /dev/null -w "%{http_code}" "$API/trending?targetType=CLIENT&window=24h")
[[ "$code" == "200" ]] || fail "trending: want 200 got $code"
echo "PASS trending ($code)"

echo "All smoke tests passed."