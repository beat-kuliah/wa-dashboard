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

# Register (creates tenant + admin user)
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
