# WA Dashboard

WA Dashboard is a multi-tenant SaaS platform for managing WhatsApp operations: broadcast campaigns, CS inbox, analytics, WhatsApp Business API connections, and optional per-tenant AI.

## Monorepo layout

```text
wa-dashboard/
  backend/          # Go modular monolith (Echo, sqlc, pgx, Asynq)
  frontend/         # Next.js App Router (Tailwind, shadcn/ui)
  docs/             # Architecture and API contract
  docker-compose.yml
```

## Tech stack

| Layer | Technology |
| --- | --- |
| Backend | Go 1.22+, Echo, sqlc, pgx, golang-migrate, Asynq |
| Frontend | Next.js (App Router), TypeScript, Tailwind, shadcn/ui, TanStack Query |
| Database | PostgreSQL 18 |
| Queue / cache | Redis 8 |
| Auth | JWT (access + refresh), RBAC (admin / supervisor / agent) |
| Realtime inbox | Server-Sent Events (SSE) |

## Documentation

- [Architecture](docs/wa-dashboard-architecture.md)
- [API contract (frozen)](docs/api-contract.md) — shared agreement for backend and frontend

## Roles & access model

The platform has three layers of access. The first ("super admin") sits **above**
all tenants and is the planned platform-operator role; the other two live **inside**
a single tenant and are already implemented.

| Role | Scope | Can do | Status |
| --- | --- | --- | --- |
| **Super admin** (platform operator) | All tenants | Provision tenants + their owner, suspend/oversee every tenant, cross-tenant support | Available |
| **Admin / owner** (`admin`) | One tenant | Configure own tenant (`PATCH /tenant`), manage members, full feature access | Available |
| **Supervisor** (`supervisor`) | One tenant | View members, manage broadcasts, view all conversations | Available |
| **Agent** (`agent`) | One tenant | Handle assigned conversations only | Available |

Key rule: a tenant is **always resolved from the JWT `tid` claim**, never from the
request body or path — so a tenant admin can only ever touch their own tenant. The
super admin role is intentionally kept **outside** this tenant model (separate auth)
so tenant isolation is never weakened.

## Onboarding & auth flows

### Today (self-service via API only)

Tenants are created through the public `POST /auth/register` endpoint, which makes a
new tenant **and** its first `admin` user in one transaction. There is **no register
UI** — the only way to call it today is directly via API/curl. The web app exposes
**login only**.

```mermaid
flowchart LR
    Operator["Operator (via API/curl)"] -->|"POST /auth/register"| Owner["New tenant + admin/owner"]
    Owner -->|"POST /tenant/members"| Members["supervisor / agent"]
    Owner -->|"login"| Dashboard["Tenant dashboard"]
    Members -->|"login"| Dashboard
```

### Target (provisioned by super admin)

The chosen direction for this product is **provisioned onboarding**: the super admin
creates the tenant + owner and hands credentials to the client. Public self-service
registration stays disabled (or becomes an internal endpoint the super admin calls).
This fits a paid B2B WhatsApp product where each tenant needs an assisted WhatsApp
Business API connection.

```mermaid
flowchart LR
    SuperAdmin["Super admin (operator panel)"] -->|"provision tenant + owner"| Owner["Tenant + admin/owner"]
    SuperAdmin -->|"suspend / oversee"| AllTenants["All tenants"]
    Owner -->|"POST /tenant/members"| Members["supervisor / agent"]
    Owner -->|"login + PATCH /tenant"| Dashboard["Tenant dashboard"]
    Members -->|"login"| Dashboard
```

### Login / session

`POST /auth/login` returns a short-lived access token (JWT) plus an opaque refresh
token. Clients call `POST /auth/refresh` to rotate the pair and `POST /auth/logout`
to revoke. See the [API contract](docs/api-contract.md) for exact shapes.

## Feature status

What exists today versus what is planned. "Backend only" means the API works but
there is no UI yet.

| Capability | Backend | Frontend | Status |
| --- | --- | --- | --- |
| Auth: login / refresh / logout / me | Done | Login UI done | Available |
| Tenant self-register (creates tenant + owner) | Done | No UI (API only) | Backend only |
| Tenant settings (`GET`/`PATCH /tenant`) | Done | Settings UI | Available |
| Member management (`list` / `add`) | Done | Settings UI | Available |
| Broadcasts (list / create / get + worker) | Done | List UI | Available |
| CS Inbox (conversations + SSE stream) | Stub | Page scaffolded | Coming soon |
| Analytics summary | Stub | Page scaffolded | Coming soon |
| WhatsApp message templates | Stub | No UI yet | Coming soon |
| **Super admin** (platform operator) | Done | Admin console at `/admin` | Available |
| **Provisioned onboarding** (super admin creates tenants) | Done | Provision form in admin console | Available |
| Public register UI | Intentionally absent | Intentionally absent | Not planned |
| WhatsApp Business API connection | Not started | Not started | Coming soon |
| Per-tenant AI chatbot | Not started | Not started | Coming soon |

## Getting started

Native PostgreSQL and Redis (no Docker required). Default URLs: API `http://localhost:8080`, frontend `http://localhost:3000`, REST base `http://localhost:8080/api/v1`.

### Prerequisites

| Tool | Version |
| --- | --- |
| Go | 1.22+ |
| Node.js | 24+ |
| PostgreSQL | 18 |
| Redis | 8 |
| [golang-migrate](https://github.com/golang-migrate/migrate) CLI | optional for `make migrate-up` (install if missing) |

Ensure PostgreSQL and Redis are running locally before starting the backend.

### 1. Database

Create the application database (once):

```bash
createdb wa_dashboard
# or: psql -U postgres -c "CREATE DATABASE wa_dashboard;"
```

### 2. Backend

```bash
cd backend
cp .env.example .env
# Default credentials (adjust if your local setup differs):
#   PostgreSQL user: postgres, password: Admin123
#   Redis password: Admin123
#   DATABASE_URL=postgres://postgres:Admin123@localhost:5432/wa_dashboard?sslmode=disable
#   REDIS_URL=redis://:Admin123@localhost:6379/0
# Set JWT_SECRET to a long random string in production.
```

Apply migrations:

```bash
make migrate-up
```

If `migrate` is not installed:

```bash
# macOS
brew install golang-migrate

# Linux (example)
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

Then run `make migrate-up` again.

Start the API and worker in **two terminals**:

```bash
# Terminal 1 — HTTP API (port 8080)
make run-api

# Terminal 2 — Asynq worker
make run-worker
```

Health checks (non-versioned): `GET http://localhost:8080/healthz`, `GET http://localhost:8080/readyz`.

See [backend/README.md](backend/README.md) for Makefile targets and layout.

### 3. Frontend

```bash
cd frontend
cp .env.example .env.local
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000). The app calls the API at `NEXT_PUBLIC_API_BASE_URL` (default `http://localhost:8080/api/v1`).

See [frontend/README.md](frontend/README.md) for lint, typecheck, and build commands.

### Quick API smoke test

With the API running and migrations applied:

```bash
BASE=http://localhost:8080/api/v1

# Register (creates tenant + admin/owner user).
# This is the provisioning path today — there is no public register UI.
# In the target model this call is made by the super admin, not end users.
curl -sS -X POST "$BASE/auth/register" \
  -H 'Content-Type: application/json' \
  -d '{"email":"owner@example.com","password":"s3cr3tpassword","business_name":"Acme Inc.","full_name":"Jane Doe"}'

# Login
RESP=$(curl -sS -X POST "$BASE/auth/login" \
  -H 'Content-Type: application/json' \
  -d '{"email":"owner@example.com","password":"s3cr3tpassword"}')
TOKEN=$(echo "$RESP" | jq -r '.tokens.access_token')

# List broadcasts (tenant-scoped)
curl -sS "$BASE/broadcasts" -H "Authorization: Bearer $TOKEN"
```

### Optional: Docker Compose

If you prefer containers for Postgres and Redis only:

```bash
docker compose up -d postgres redis
```

Then point `DATABASE_URL` and `REDIS_URL` in `backend/.env` at the compose services.

### Platform admin smoke test

With the API running, migrations applied, and `PLATFORM_ADMIN_EMAIL` / `PLATFORM_ADMIN_PASSWORD`
set in `backend/.env`:

```bash
chmod +x scripts/smoke-platform-admin.sh
./scripts/smoke-platform-admin.sh
```

Open the operator console at [http://localhost:3000/admin/login](http://localhost:3000/admin/login)
(tenant users sign in at `/login`).
