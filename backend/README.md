# WA Dashboard Backend

Go modular monolith for the WA Dashboard multi-tenant WhatsApp SaaS platform.

## Stack

- **Go 1.22** + Echo HTTP framework
- **PostgreSQL** via pgx connection pool
- **sqlc** for type-safe SQL queries
- **golang-migrate** for schema migrations
- **Asynq** (Redis) for async jobs
- **slog** for structured JSON logging

## Prerequisites

- Go 1.22+
- PostgreSQL 18 (local)
- Redis 8 (local)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI
- [sqlc](https://sqlc.dev/) CLI (optional — generated code is committed)

## Setup

```bash
cp .env.example .env
# Defaults: postgres/Admin123, Redis password Admin123 (see .env.example)

createdb wa_dashboard   # if needed
make migrate-up
make run-api            # terminal 1
make run-worker         # terminal 2
```

API base URL: `http://localhost:8080/api/v1`

Health checks (non-versioned):

- `GET /healthz` — liveness
- `GET /readyz` — readiness (PostgreSQL + Redis)

## Makefile targets

| Target | Description |
| --- | --- |
| `make migrate-up` | Apply pending migrations |
| `make migrate-down` | Roll back one migration |
| `make sqlc` | Regenerate sqlc code from `db/queries/` |
| `make run-api` | Start HTTP API server |
| `make run-worker` | Start Asynq worker |
| `make build` | Compile all packages |
| `make tidy` | Tidy Go modules |

## Project layout

```text
backend/
  cmd/api/          HTTP entrypoint
  cmd/worker/       Asynq worker entrypoint
  internal/
    app/            Application wiring
    modules/        Business modules (auth, tenant, broadcast, inbox, analytics, template)
    shared/         Cross-cutting infra (config, db, httpx, auth, errors, logx, queue, obs)
  db/
    migrations/     SQL migrations (golang-migrate)
    queries/        sqlc query definitions
    sqlc/           Generated pgx code
```

## API contract

All endpoints conform to the frozen contract in [`docs/api-contract.md`](../docs/api-contract.md).

## Modules

| Module | Endpoints | Notes |
| --- | --- | --- |
| auth | `/auth/*` | Register, login, refresh, logout, me |
| tenant | `/tenant/*` | Workspace CRUD, member management |
| broadcast | `/broadcasts/*` | List, create, get; Asynq `broadcast.send` job |
| inbox | `/inbox/*` | Conversation list stub, SSE stream |
| analytics | `/analytics/summary` | Dashboard metrics stub |
| template | `/templates` | Template list stub |
