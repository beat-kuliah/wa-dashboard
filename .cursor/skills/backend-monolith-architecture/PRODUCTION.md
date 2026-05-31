# Production Patterns

Concrete, copy-adaptable patterns for the non-negotiable concerns. TypeScript + Postgres (`pg`) + BullMQ assumed; adapt to the project's stack.

## 1. Validated config (12-factor, fail-fast)

```ts
// src/shared/config/index.ts
import { z } from "zod";

const schema = z.object({
  NODE_ENV: z.enum(["development", "test", "production"]).default("development"),
  PORT: z.coerce.number().default(3000),
  DATABASE_URL: z.string().url(),
  REDIS_URL: z.string().url(),
  JWT_PUBLIC_KEY: z.string().min(1),
  LOG_LEVEL: z.enum(["debug", "info", "warn", "error"]).default("info"),
});

const parsed = schema.safeParse(process.env);
if (!parsed.success) {
  console.error("Invalid config:", parsed.error.flatten().fieldErrors);
  process.exit(1); // refuse to boot with bad config
}
export const config = Object.freeze(parsed.data);
```

Access config only via this module. Document every key in `.env.example`. Never log secrets.

## 2. Typed error hierarchy + central middleware

```ts
// src/shared/errors/index.ts
export class AppError extends Error {
  constructor(
    message: string,
    readonly status = 500,
    readonly code = "INTERNAL",
    readonly details?: unknown,
  ) { super(message); }
}
export class NotFoundError extends AppError {
  constructor(what = "Resource") { super(`${what} not found`, 404, "NOT_FOUND"); }
}
export class ValidationError extends AppError {
  constructor(details: unknown) { super("Validation failed", 422, "VALIDATION", details); }
}
export class ForbiddenError extends AppError {
  constructor() { super("Forbidden", 403, "FORBIDDEN"); }
}
```

```ts
// src/shared/http/error-middleware.ts
export function errorMiddleware(err, req, res, _next) {
  const e = err instanceof AppError ? err : new AppError("Internal error");
  if (e.status >= 500) req.log.error({ err }, "unhandled error"); // full log
  res.status(e.status).json({ error: { code: e.code, message: e.message, details: e.details } });
}
```

Clients get a stable shape `{ error: { code, message } }`; stack traces stay server-side.

## 3. Repository + service-owned transactions

```ts
// src/shared/db/base-repository.ts
import { pool } from "./pool";
export async function withTransaction<T>(fn: (tx) => Promise<T>): Promise<T> {
  const client = await pool.connect();
  try {
    await client.query("BEGIN");
    const result = await fn(client);
    await client.query("COMMIT");
    return result;
  } catch (e) {
    await client.query("ROLLBACK");
    throw e;
  } finally {
    client.release();
  }
}
```

```ts
// broadcast.service.ts — service owns the transaction boundary, repo runs queries
async function schedule(input: ScheduleDTO, ctx: RequestContext) {
  return withTransaction(async (tx) => {
    const b = await broadcastRepo.create(tx, { ...input, tenantId: ctx.tenantId });
    await broadcastRepo.attachRecipients(tx, b.id, input.recipientIds);
    await queue.add("broadcast.send", { broadcastId: b.id }); // enqueue after commit-safe work
    return b;
  });
}
```

Always parameterize queries (`$1, $2`). Never string-concat SQL. Every query in a multi-tenant app includes `WHERE tenant_id = $1`.

## 4. Auth + tenant context (propagated, not passed everywhere)

```ts
// src/shared/http/request-context.ts
import { AsyncLocalStorage } from "node:async_hooks";
export const als = new AsyncLocalStorage<RequestContext>();
export const getContext = () => als.getStore()!;
```

```ts
// middleware: verify token, resolve tenant, run rest of request inside context
export async function authMiddleware(req, res, next) {
  const claims = verifyToken(req.headers.authorization);          // throws -> 401
  const ctx = { userId: claims.sub, tenantId: claims.tid, roles: claims.roles, requestId: req.id };
  als.run(ctx, () => next());
}
```

```ts
// RBAC enforced in the service layer (defense in depth, not just routes)
function assertCan(ctx: RequestContext, permission: string) {
  if (!hasPermission(ctx.roles, permission)) throw new ForbiddenError();
}
```

## 5. Idempotent, retryable async jobs

```ts
// broadcast.worker.ts
new Worker("broadcast.send", async (job) => {
  const { broadcastId } = job.data;
  // idempotency: skip if already processed
  if (await broadcastRepo.isSent(broadcastId)) return;
  await sendViaWhatsApp(broadcastId);
  await broadcastRepo.markSent(broadcastId);
}, {
  connection,
  concurrency: 10,
}).on("failed", (job, err) => logger.error({ jobId: job?.id, err }, "job failed"));
// Configure attempts + backoff on enqueue: queue.add(name, data, { attempts: 5, backoff: { type: "exponential", delay: 1000 } })
```

Jobs must be safe to run twice. Use a dead-letter queue or `failed` handler for poison messages.

## 6. Observability

```ts
// request logging middleware: bind a child logger with ids
req.log = logger.child({ requestId: req.id, tenantId: getContext()?.tenantId });

// RED metrics (prom-client)
httpDuration.observe({ method, route, status }, durationSeconds);

// health endpoints
app.get("/healthz", () => res.sendStatus(200));                 // liveness: process up
app.get("/readyz", async () => {                                 // readiness: deps reachable
  await pool.query("SELECT 1"); await redis.ping();
  res.sendStatus(200);
});
```

Use OpenTelemetry for tracing; propagate trace context into queue jobs so async work stays in the same trace.

## 7. Resilience on outbound calls

```ts
const res = await fetchWithTimeout(url, { timeoutMs: 5000 });   // always set a timeout
// retry transient failures with exponential backoff + jitter
// wrap flaky third parties (WhatsApp API) in a circuit breaker (e.g. opossum)
```

Add rate limiting on public/auth endpoints. Set request body size limits. Idempotency keys on POST endpoints that create resources.

## 8. Graceful shutdown

```ts
// src/server.ts
const server = app.listen(config.PORT);
async function shutdown(signal: string) {
  logger.info({ signal }, "shutting down");
  server.close();                  // stop accepting new connections
  await worker.close();            // finish/quiesce in-flight jobs
  await pool.end();                // close DB pool
  await redis.quit();
  process.exit(0);
}
process.on("SIGTERM", () => shutdown("SIGTERM"));
process.on("SIGINT", () => shutdown("SIGINT"));
```

## 9. Migrations

- Forward-only, versioned, checked into git (node-pg-migrate, Prisma Migrate, golang-migrate).
- Run migrations as a separate step in deploy/CI — never auto-migrate on app boot in production.
- Migrations must be backward-compatible for zero-downtime deploys (expand → migrate data → contract).

## 10. CI gates (block merge on failure)

```
lint  →  typecheck  →  unit tests  →  integration tests (test DB)  →  build  →  dependency audit
```

Run integration tests against a real Postgres/Redis via testcontainers or a CI service container, not mocks.
