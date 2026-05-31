# Folder Structure Reference

## Full annotated tree (Node.js + TypeScript)

```text
.
├── src/
│   ├── modules/                      # business capabilities (module-first)
│   │   ├── broadcast/
│   │   │   ├── broadcast.routes.ts       # registers HTTP routes on a router
│   │   │   ├── broadcast.controller.ts   # thin: validate → call service → serialize
│   │   │   ├── broadcast.service.ts      # business rules, transactions, orchestration
│   │   │   ├── broadcast.repository.ts   # SQL/ORM, returns domain types
│   │   │   ├── broadcast.schema.ts       # zod input/output schemas + inferred DTOs
│   │   │   ├── broadcast.types.ts        # domain entities/value objects
│   │   │   ├── broadcast.worker.ts       # queue consumer for async send
│   │   │   ├── broadcast.test.ts
│   │   │   └── index.ts                  # public API: export { BroadcastService, types }
│   │   ├── inbox/
│   │   ├── analytics/
│   │   └── tenant/
│   ├── shared/                       # cross-cutting infra; NO business logic
│   │   ├── config/
│   │   │   └── index.ts                  # load + zod-validate env, export typed config
│   │   ├── db/
│   │   │   ├── pool.ts                    # pg pool / ORM client
│   │   │   ├── base-repository.ts         # shared query helpers, tx helper
│   │   │   └── migrations/                # versioned, forward-only SQL
│   │   ├── http/
│   │   │   ├── server.ts                  # build app, mount middleware + module routers
│   │   │   ├── error-middleware.ts        # central error → HTTP mapping
│   │   │   ├── request-context.ts         # async-local-storage context
│   │   │   └── middleware/                # auth, rate-limit, request-id, body-limit
│   │   ├── auth/
│   │   │   ├── verify-token.ts
│   │   │   ├── tenant-context.ts          # resolve + scope tenant
│   │   │   └── rbac.ts                     # role/permission guards
│   │   ├── errors/
│   │   │   └── index.ts                    # AppError, NotFoundError, ValidationError...
│   │   ├── logger/
│   │   │   └── index.ts                    # pino/winston, JSON, bound request/tenant id
│   │   ├── queue/
│   │   │   ├── index.ts                    # BullMQ connection + registry
│   │   │   └── base-worker.ts
│   │   ├── events/
│   │   │   └── bus.ts                      # in-process domain event emitter
│   │   └── observability/
│   │       ├── metrics.ts                  # prom-client RED metrics
│   │       ├── tracing.ts                  # OpenTelemetry setup
│   │       └── health.ts                   # /healthz, /readyz
│   ├── app.ts                         # compose: register all module routers + shared
│   └── server.ts                      # entrypoint: load config, listen, graceful shutdown
├── test/
│   ├── integration/                   # routes/repos against a real test DB
│   └── e2e/                           # critical-flow smoke tests
├── migrations/                        # (or src/shared/db/migrations)
├── scripts/                           # seed, one-off ops scripts
├── .env.example                       # documents every config key
├── Dockerfile
├── docker-compose.yml                 # local: app + postgres + redis
├── package.json
└── tsconfig.json                      # path alias "@/*" -> "src/*"
```

## Module public API (`index.ts`) example

Only this file is importable by other modules. It hides repositories and internals.

```ts
// src/modules/broadcast/index.ts
export { BroadcastService } from "./broadcast.service";
export type { Broadcast, BroadcastStatus } from "./broadcast.types";
// NOTE: repository, controller, schema are intentionally NOT exported.
```

## Cross-module communication

Prefer, in order:
1. **Direct call to another module's public service** when you need a synchronous result (e.g. `analytics` reads via `TenantService.getTenant(id)`).
2. **Domain events** for fire-and-forget side effects, to keep modules decoupled:

```ts
// publisher (broadcast.service.ts)
events.emit("broadcast.sent", { tenantId, broadcastId, recipients });

// subscriber (analytics module registers at startup)
events.on("broadcast.sent", (e) => analyticsService.recordBroadcast(e));
```

Never let a module query another module's tables directly.

## Go variant

```text
cmd/
  api/main.go                 # entrypoint
internal/
  modules/
    broadcast/
      handler.go              # interface layer
      service.go              # business logic
      repository.go           # data access
      broadcast.go            # domain types
      module.go               # exported constructor / public funcs
  shared/
    config/                   # env via envconfig/viper
    db/                       # pgx pool, tx helpers
    httpx/                    # router, middleware, error mapping
    auth/                     # token + tenant context
    logx/                     # zap/slog structured logging
    queue/                    # async jobs
    obs/                      # metrics, tracing, health
```

`internal/` enforces that nothing outside the module tree can import internals — Go's compiler gives you boundary enforcement for free. Keep each module behind an exported interface defined in `module.go`.
