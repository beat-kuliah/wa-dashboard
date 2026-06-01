#!/usr/bin/env bash
# Smoke test: platform admin login → provision tenant → owner login → suspend → owner blocked
set -euo pipefail

BASE="${BASE:-http://localhost:8080/api/v1}"
ADMIN_EMAIL="${PLATFORM_ADMIN_EMAIL:-ops@wa-dashboard.local}"
ADMIN_PASSWORD="${PLATFORM_ADMIN_PASSWORD:-change-me-platform-admin}"
OWNER_EMAIL="${SMOKE_OWNER_EMAIL:-smoke-owner-$(date +%s)@example.com}"
OWNER_PASSWORD="${SMOKE_OWNER_PASSWORD:-s3cr3tpassword}"
BUSINESS_NAME="${SMOKE_BUSINESS_NAME:-Smoke Test Co}"

echo "==> Platform admin login"
ADMIN_RESP=$(curl -sS -X POST "$BASE/admin/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")
ADMIN_TOKEN=$(echo "$ADMIN_RESP" | jq -r '.access_token')
if [[ "$ADMIN_TOKEN" == "null" || -z "$ADMIN_TOKEN" ]]; then
  echo "Admin login failed: $ADMIN_RESP" >&2
  exit 1
fi

echo "==> Provision tenant + owner"
PROVISION_RESP=$(curl -sS -X POST "$BASE/admin/tenants" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d "{\"business_name\":\"$BUSINESS_NAME\",\"owner_email\":\"$OWNER_EMAIL\",\"owner_full_name\":\"Smoke Owner\",\"owner_password\":\"$OWNER_PASSWORD\"}")
TENANT_ID=$(echo "$PROVISION_RESP" | jq -r '.tenant.id')
if [[ "$TENANT_ID" == "null" || -z "$TENANT_ID" ]]; then
  echo "Provision failed: $PROVISION_RESP" >&2
  exit 1
fi
echo "Tenant id: $TENANT_ID"

echo "==> Owner login (active tenant)"
OWNER_RESP=$(curl -sS -X POST "$BASE/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"email\":\"$OWNER_EMAIL\",\"password\":\"$OWNER_PASSWORD\"}")
OWNER_TOKEN=$(echo "$OWNER_RESP" | jq -r '.tokens.access_token // .access_token // empty')
if [[ -z "$OWNER_TOKEN" ]]; then
  # handler may nest tokens
  OWNER_TOKEN=$(echo "$OWNER_RESP" | jq -r '.tokens.access_token')
fi
if [[ "$OWNER_TOKEN" == "null" || -z "$OWNER_TOKEN" ]]; then
  echo "Owner login failed: $OWNER_RESP" >&2
  exit 1
fi

echo "==> Suspend tenant"
curl -sS -X PATCH "$BASE/admin/tenants/$TENANT_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"status":"suspended"}' | jq -e '.status == "suspended"' >/dev/null

echo "==> Owner API call should return TENANT_SUSPENDED"
HTTP_CODE=$(curl -sS -o /tmp/smoke-tenant.json -w '%{http_code}' \
  "$BASE/tenant" -H "Authorization: Bearer $OWNER_TOKEN")
BODY=$(cat /tmp/smoke-tenant.json)
if [[ "$HTTP_CODE" != "403" ]] || ! echo "$BODY" | jq -e '.error.code == "TENANT_SUSPENDED"' >/dev/null; then
  echo "Expected 403 TENANT_SUSPENDED, got HTTP $HTTP_CODE: $BODY" >&2
  exit 1
fi

echo "==> All smoke checks passed"
