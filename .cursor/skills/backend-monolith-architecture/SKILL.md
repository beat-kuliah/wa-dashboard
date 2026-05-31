---
name: backend-monolith-architecture
description: Design and build a scalable, maintainable, production-ready backend as a modular monolith following conventions used by large engineering orgs (Shopify, GitHub, Stripe, Basecamp). Use when scaffolding a new backend, adding a feature/module, deciding folder structure, layering, data access, multi-tenancy, or when the user mentions backend architecture, modular monolith, project structure, scalability, or maintainability for a server/API.
disable-model-invocation: true
---

# Backend Modular Monolith Architecture

Build backends as a **modular monolith**: a single deployable unit composed of strongly-isolated modules with explicit boundaries. This is what most large companies run in production because it gives you the operational simplicity of a monolith with the maintainability of microservices, and lets you extract a module into a service later only if you actually need to.

Default stack assumed here: **Node.js + TypeScript** (Express/Fastify/NestJS) with **PostgreSQL** and **Redis**. The patterns are language-agnostic; Go notes are included where relevant. Match the project's existing stack before introducing anything new.

## Core Principles

1. **Modules over layers at the top level.** Organize first by business capability (`billing`, `inbox`, `broadcast`), then by technical layer inside each module. Never create a global `controllers/` + `services/` + `models/` split for the whole app.
2. **Explicit boundaries.** Modules talk to each other only through a published interface (the module's `index.ts` / public API), never by reaching into another module's internals or tables.
3. **Dependencies point inward.** Routes → services (use cases) → repositories → DB. Domain logic never imports framework/HTTP/ORM types.
4. **Stateless app tier.** No in-process session state. Put state in Postgres/Redis so you can run N replicas behind a load balancer.
5. **Everything observable, everything configurable.** Structured logs, metrics, traces, health checks, and 12-factor config from day one.
6. **Fail loud, fail safe.** Centralized error handling, typed errors, input validation at the edge, no silent catches.

## Layering Inside a Module

Use 3 layers. Keep the dependency direction strict.

| Layer | Responsibility | Knows about |
| --- | --- | --- |
| **Interface** (routes/controllers, queue consumers, cron) | Parse/validate input, authn/authz, map to use case, serialize response | Services only |
| **Application/Domain** (services / use cases, domain entities) | Business rules, orchestration, transactions | Repositories (via interfaces) |
| **Infrastructure** (repositories, external clients) | DB queries, HTTP/3rd-party calls, queues | DB/SDKs |

Controllers must be thin. Business logic lives in services. Data access lives in repositories. A service must be testable without HTTP or a real DB (inject repository interfaces).

## Folder Structure

This is the canonical layout. See [STRUCTURE.md](STRUCTURE.md) for the full annotated tree and a Go variant.

```text
src/
  modules/
    <module>/
      <module>.routes.ts        # interface layer (HTTP)
      <module>.controller.ts
      <module>.service.ts       # application layer (business logic)
      <module>.repository.ts    # infrastructure layer (data access)
      <module>.schema.ts        # zod/valibot validation + DTO types
      <module>.types.ts         # domain types/entities
      <module>.test.ts
      index.ts                  # PUBLIC API of the module (the only thing others import)
  shared/                       # cross-cutting, no business logic
    config/                     # env loading + validation (12-factor)
    db/                         # connection pool, migrations runner, base repo
    http/                       # server bootstrap, middleware, error handler
    logger/                     # structured logging
    errors/                     # typed error classes (AppError, NotFound, ...)
    auth/                       # token verify, tenant context, RBAC guards
    queue/                      # job queue setup (BullMQ), base worker
    observability/              # metrics, tracing, health checks
    events/                     # in-process domain event bus
  app.ts                        # wire modules + shared into the app
  server.ts                     # entrypoint: load config, start, graceful shutdown
```

Rules:
- Other modules import `from "@/modules/billing"` (the `index.ts`), never `from "@/modules/billing/billing.repository"`.
- `shared/` may be imported by any module; `shared/` must NOT import from `modules/`.
- Cross-module calls go through the public service interface or the event bus, never via direct DB access to another module's tables.

## Scaffolding Workflow

When asked to create a new backend, follow this checklist:

```
- [ ] 1. Confirm stack (runtime, framework, DB, queue) — match existing project
- [ ] 2. Create shared/ foundation: config, logger, errors, db, http server + error middleware
- [ ] 3. Add health/readiness endpoints + graceful shutdown
- [ ] 4. Scaffold first module end-to-end (routes→service→repo→schema→test)
- [ ] 5. Add auth + tenant-context middleware (if multi-tenant)
- [ ] 6. Add migrations + seed; wire CI (lint, typecheck, test, migrate)
- [ ] 7. Add observability (request logging, metrics, tracing) + Dockerfile
```

When adding a **new feature**, create or extend a module — never spread its logic across global folders. Generate all of: routes, controller, service, repository, schema, types, test, and export the public API in `index.ts`.

## Non-Negotiable Production Concerns

These must exist before calling a backend "production ready". Details and code patterns are in [PRODUCTION.md](PRODUCTION.md).

- **Config & secrets**: validated env at boot (zod). App refuses to start if config is invalid. Never read `process.env` deep in the code.
- **Validation**: validate every external input at the interface layer (zod schema → typed DTO). Reject unknown fields.
- **Error handling**: one centralized error middleware. Typed `AppError` hierarchy with HTTP status + machine code. Never leak stack traces to clients.
- **AuthN/AuthZ**: verify tokens in middleware; attach `requestContext` (userId, tenantId, roles). Enforce RBAC in the service layer, not just routes.
- **Multi-tenancy**: tenant id resolved once in middleware and propagated; every query is tenant-scoped. Pick an isolation model (shared schema + `tenant_id` column is the default) and enforce it consistently.
- **Data access**: repository pattern; transactions owned by the service layer; parameterized queries only; connection pooling; explicit migrations (no auto-sync in prod).
- **Async work**: long/spiky work (broadcasts, webhooks, emails) goes to a queue (BullMQ/Redis), not the request path. Idempotent, retryable jobs with dead-letter handling.
- **Observability**: structured JSON logs with request id + tenant id; RED metrics (rate/errors/duration); distributed tracing; `/healthz` (liveness) + `/readyz` (deps reachable).
- **Resilience**: timeouts on every outbound call, retries with backoff, circuit breakers for flaky deps, rate limiting on public endpoints.
- **Graceful lifecycle**: handle SIGTERM, stop accepting new requests, drain in-flight, close DB/queue, then exit.
- **Security**: helmet/secure headers, CORS allowlist, body size limits, secrets never logged, dependency audit in CI.
- **Testing**: unit tests on services (mock repos), integration tests on repos/routes (real test DB via testcontainers), one e2e smoke per critical flow.

## Anti-Patterns (reject these)

- A global `controllers/`, `services/`, `models/` split for the whole app (layer-first instead of module-first).
- Business logic in controllers or in the ORM models.
- One module reading another module's tables or importing its internal files.
- `process.env.X` scattered across the codebase.
- Catch-and-ignore, or returning `200` on errors.
- ORM `synchronize: true` / auto-migrate in production.
- Doing heavy work (loops of API calls, large exports) inside an HTTP handler.
- Shared mutable singletons holding request/tenant state.

## When To Split Into Services

Stay a monolith until a specific module has a *different* scaling, deployment, or team-ownership need. Because boundaries are already explicit (public `index.ts` + events + tenant-scoped tables), extracting a module later is a mechanical refactor, not a rewrite. Do not start with microservices.

## References

- [STRUCTURE.md](STRUCTURE.md) — full annotated folder tree, module template, Go variant
- [PRODUCTION.md](PRODUCTION.md) — concrete code patterns for config, errors, repository/transactions, auth/tenant context, queue jobs, observability, and graceful shutdown
