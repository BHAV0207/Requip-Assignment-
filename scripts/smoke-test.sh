#!/usr/bin/env bash
# End-to-end user journey through the real entry point (ingress), plus a
# check that the running version matches what Git says should be deployed.
set -euo pipefail

BASE="${BASE:-http://requip.localtest.me}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
pass() { printf '  \033[32mPASS\033[0m %s\n' "$1"; }
fail() { printf '  \033[31mFAIL\033[0m %s\n' "$1"; exit 1; }

echo "Smoke test against $BASE"

code=$(curl -s -o /dev/null -w '%{http_code}' "$BASE/")
[[ "$code" == 200 ]] && pass "frontend served (HTTP $code)" || fail "frontend returned $code"

version=$(curl -sf "$BASE/api/v1/version")
pass "API version: $version"

expected=$(grep -A2 '^backend:' "$ROOT/deploy/helm/requip/values-kind.yaml" | awk -F'"' '/tag:/ {print $2}')
running=$(echo "$version" | sed -E 's/.*"version":"([^"]+)".*/\1/')
[[ "$expected" == "$running" ]] && pass "running version matches Git ($expected)" \
  || echo "  WARN running '$running' != Git '$expected' (Argo CD may still be syncing)"

suffix=$(printf "%09d" $(( $(date +%s) % 1000000000 )))   # 9 digits, unique per run
body=$(cat <<JSON
{"name":"Smoke Test","email":"smoke-$suffix@example.com","primaryMobile":"8$suffix","aadhaar":"777$suffix","pan":"SMOKE${suffix:5:4}Z","dateOfBirth":"1990-01-01","placeOfBirth":"Pune","currentAddress":"1 Test Road","permanentAddress":"1 Test Road"}
JSON
)
created=$(curl -sf -X POST "$BASE/api/v1/users" -H 'Content-Type: application/json' -d "$body") || fail "create user"
id=$(echo "$created" | sed -E 's/.*"id":"([^"]+)".*/\1/')
pass "created user $id"

curl -sf "$BASE/api/v1/users/$id" | grep -q "smoke-$suffix@example.com" && pass "read user back" || fail "read user"

curl -sf -X PUT "$BASE/api/v1/users/$id" -H 'Content-Type: application/json' -d '{"placeOfBirth":"Mumbai"}' | grep -q Mumbai \
  && pass "updated user" || fail "update user"

curl -sf -X DELETE "$BASE/api/v1/users/$id" >/dev/null && pass "deleted user" || fail "delete user"
code=$(curl -s -o /dev/null -w '%{http_code}' "$BASE/api/v1/users/$id")
[[ "$code" == 404 ]] && pass "deleted user returns 404" || fail "expected 404, got $code"

echo "All checks passed."
