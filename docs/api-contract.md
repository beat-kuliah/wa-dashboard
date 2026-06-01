# WA Dashboard — Frozen API Contract (v1)

> **This document is the single source of truth** for the REST API shared by the
> backend (`backend/`) and the frontend (`frontend/`). It is **frozen**: both
> teams implement against exactly these shapes so they can build independently.
> The [Data Models](#data-models) section is canonical — endpoint examples
> reference those shapes rather than redefining them.

## Conventions

- **Base URL:** `http://localhost:8080/api/v1` (configurable via env; backend
  `PORT`, frontend `NEXT_PUBLIC_API_BASE_URL`). All paths below are relative to
  this base, e.g. `POST /auth/login` → `http://localhost:8080/api/v1/auth/login`.
- **Content type:** all requests and responses are `application/json`
  (exception: `GET /inbox/stream` returns `text/event-stream`).
- **JSON keys:** `snake_case` everywhere (matches Go `json` tags and the frontend).
- **Timestamps:** ISO-8601 strings in **UTC** (e.g. `2026-05-31T16:15:00Z`).
- **IDs:** UUID strings (e.g. `9f1c2d3e-4a5b-6c7d-8e9f-0a1b2c3d4e5f`).
- **Pagination:** list endpoints accept `limit` and `offset` query params.
  Default `limit` is `20` (max `100`), default `offset` is `0`.
- **Auth header:** protected endpoints require `Authorization: Bearer <access_token>`.

## Authentication & multi-tenancy

- Auth uses a **short-lived JWT access token** plus a longer-lived **refresh token**.
- The access token is sent on every authenticated request as
  `Authorization: Bearer <access_token>`.
- The platform is **multi-tenant**: every authenticated user belongs to exactly
  one tenant. The tenant is resolved **from the JWT claims**, not from the request
  body or path — the server scopes all queries to that tenant.
- **JWT access token claims (informative):**

```json
{
  "sub": "9f1c2d3e-4a5b-6c7d-8e9f-0a1b2c3d4e5f",  // user id
  "tid": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",  // tenant id
  "role": "admin",                                  // admin | supervisor | agent
  "iat": 1748707200,
  "exp": 1748708100                                 // short-lived (e.g. 15 min)
}
```

- **Roles (RBAC):** `admin`, `supervisor`, `agent` (tenant-scoped JWT with `tid`).
  - `admin` — full access, manage tenant + members.
  - `supervisor` — view all conversations, manage broadcasts, view members.
  - `agent` — handle assigned conversations only.
- **Platform operator JWT** (no `tid`): `role` is `platform_admin`. Used only on
  `/admin/*` endpoints. Tenant tokens must not be accepted on admin routes, and
  admin tokens must not grant tenant-scoped access.

```json
{
  "sub": "9f1c2d3e-4a5b-6c7d-8e9f-0a1b2c3d4e5f",
  "role": "platform_admin",
  "iat": 1748707200,
  "exp": 1748708100
}
```
- **Platform admin access token claims (informative):** super admins (platform
  operators) live outside any tenant, so their access token has **no `tid`
  claim** and a fixed `role` of `platform_admin`. A token without `tid` is
  rejected on tenant-scoped endpoints, and a tenant token (with `tid`) is
  rejected on `/admin/*` endpoints — the two token kinds are not interchangeable.

```json
{
  "sub": "0a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",  // platform admin id
  "role": "platform_admin",                         // fixed; no other roles
  "iat": 1748707200,
  "exp": 1748708100                                 // short-lived (e.g. 15 min)
}
```
- The **refresh token** is an opaque string persisted server-side; clients call
  `POST /auth/refresh` to rotate it for a new access/refresh pair. Refresh tokens
  are revoked on `POST /auth/logout`.

## Error envelope

**Every non-2xx response** uses this exact envelope:

```json
{
  "error": {
    "code": "STRING_CODE",
    "message": "human readable explanation"
  }
}
```

| HTTP status | `code`         | Meaning                                              |
| ----------- | -------------- | ---------------------------------------------------- |
| 401         | `UNAUTHORIZED` | Missing/invalid/expired access token.                |
| 403         | `FORBIDDEN`    | Authenticated but lacks the required role/permission.|
| 403         | `TENANT_SUSPENDED` | Tenant is suspended; tenant-scoped requests are blocked. |
| 403         | `TENANT_SUSPENDED` | The caller's tenant is `suspended`; authentication/API access is blocked until it is reactivated. |
| 404         | `NOT_FOUND`    | Resource does not exist (or not in caller's tenant). |
| 409         | `CONFLICT`     | Conflict, e.g. email already registered.             |
| 422         | `VALIDATION`   | Request body/params failed validation.               |
| 429         | `RATE_LIMITED` | Too many requests.                                   |
| 500         | `INTERNAL`     | Unexpected server error.                             |

`VALIDATION` errors MAY include a `details` array for field-level messages:

```json
{
  "error": {
    "code": "VALIDATION",
    "message": "Validation failed",
    "details": [
      { "field": "email", "message": "must be a valid email" },
      { "field": "password", "message": "must be at least 8 characters" }
    ]
  }
}
```

## Success envelopes

- **Single resource:** the resource object is returned **directly** (no wrapper).
- **Lists:** wrapped with pagination metadata:

```json
{
  "data": [ /* array of resource objects */ ],
  "page": {
    "limit": 20,
    "offset": 0,
    "total": 137
  }
}
```

---

# Data Models

These are the canonical object shapes. Endpoints below reference them by name.

### User

```json
{
  "id": "9f1c2d3e-4a5b-6c7d-8e9f-0a1b2c3d4e5f",
  "email": "owner@acme.com",
  "full_name": "Jane Doe",
  "role": "admin",
  "tenant_id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
  "created_at": "2026-05-31T16:15:00Z"
}
```

| Field        | Type   | Notes                                  |
| ------------ | ------ | -------------------------------------- |
| `id`         | string | UUID                                   |
| `email`      | string | unique                                 |
| `full_name`  | string |                                        |
| `role`       | string | `admin` \| `supervisor` \| `agent`     |
| `tenant_id`  | string | UUID of the owning tenant              |
| `created_at` | string | ISO-8601 UTC                           |

### Tenant

```json
{
  "id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
  "business_name": "Acme Inc.",
  "status": "active",
  "settings": {},
  "ai_enabled": false,
  "features": {
    "broadcast": true,
    "cs_inbox": true,
    "analytics": true,
    "ai_chatbot": false
  },
  "created_at": "2026-05-31T16:15:00Z"
}
```

| Field           | Type    | Notes                                                  |
| --------------- | ------- | ------------------------------------------------------ |
| `id`            | string  | UUID                                                   |
| `business_name` | string  |                                                        |
| `status`        | string  | `active` \| `suspended` (default `active`)             |
| `settings`      | object  | free-form per-tenant settings (default `{}`)           |
| `ai_enabled`    | boolean | tenant-level AI master switch                          |
| `features`      | object  | feature flags (see below)                              |
| `created_at`    | string  | ISO-8601 UTC                                           |

`features` object:

| Field        | Type    |
| ------------ | ------- |
| `broadcast`  | boolean |
| `cs_inbox`   | boolean |
| `analytics`  | boolean |
| `ai_chatbot` | boolean |

### PlatformAdmin

A platform operator (super admin). Lives **outside** the tenant model and is not
a [User](#user); it can provision and manage tenants but never belongs to one.

```json
{
  "id": "0a1b2c3d-4e5f-6a7b-8c9d-0e1f2a3b4c5d",
  "email": "ops@wa-dashboard.com",
  "full_name": "Platform Operator",
  "created_at": "2026-05-31T16:15:00Z"
}
```

| Field        | Type   | Notes        |
| ------------ | ------ | ------------ |
| `id`         | string | UUID         |
| `email`      | string | unique       |
| `full_name`  | string |              |
| `created_at` | string | ISO-8601 UTC |

### Broadcast

```json
{
  "id": "7c0d1e2f-3a4b-5c6d-7e8f-9a0b1c2d3e4f",
  "tenant_id": "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
  "name": "May Promo Blast",
  "template_id": "tmpl_2f3a4b5c",
  "segment": "all_contacts",
  "status": "scheduled",
  "scheduled_at": "2026-06-01T09:00:00Z",
  "stats": {
    "delivered": 0,
    "read": 0,
    "failed": 0
  },
  "created_at": "2026-05-31T16:15:00Z"
}
```

| Field          | Type           | Notes                                                          |
| -------------- | -------------- | -------------------------------------------------------------- |
| `id`           | string         | UUID                                                           |
| `tenant_id`    | string         | UUID                                                           |
| `name`         | string         |                                                                |
| `template_id`  | string         | WhatsApp template reference                                    |
| `segment`      | string         | contact segment identifier                                     |
| `status`       | string         | `draft` \| `scheduled` \| `sending` \| `sent` \| `failed`      |
| `scheduled_at` | string \| null | ISO-8601 UTC; `null` when not scheduled (send-now/draft)       |
| `stats`        | object         | delivery counters (see below)                                  |
| `created_at`   | string         | ISO-8601 UTC                                                   |

`stats` object:

| Field       | Type   |
| ----------- | ------ |
| `delivered` | number |
| `read`      | number |
| `failed`    | number |

**Broadcast status enum:** `draft`, `scheduled`, `sending`, `sent`, `failed`.

---

# Endpoints

## Auth

### POST /auth/register

Creates a new tenant and its initial **admin** user. Public when enabled via
`PUBLIC_REGISTRATION_ENABLED` (default `false`). When disabled, returns
`403 FORBIDDEN`.

Request:

```json
{
  "email": "owner@acme.com",
  "password": "s3cr3tpassword",
  "business_name": "Acme Inc.",
  "full_name": "Jane Doe"
}
```

`201 Created`:

```json
{
  "user": { /* User */ },
  "tenant": { /* Tenant */ },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "rt_4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c"
}
```

Errors: `422 VALIDATION`, `403 FORBIDDEN` (registration disabled),
`409 CONFLICT` (email already registered).

Tenant-scoped requests for a **suspended** tenant return `403` with code
`TENANT_SUSPENDED` (see error envelope).

### POST /auth/login

Authenticates an existing user. Public.

Request:

```json
{
  "email": "owner@acme.com",
  "password": "s3cr3tpassword"
}
```

`200 OK`:

```json
{
  "user": { /* User */ },
  "tenant": { /* Tenant */ },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "rt_4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c"
}
```

Errors: `422 VALIDATION`, `401 UNAUTHORIZED` (bad credentials).

### POST /auth/refresh

Rotates a refresh token for a new access/refresh pair. Public (the refresh token
is the credential).

Request:

```json
{
  "refresh_token": "rt_4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c"
}
```

`200 OK`:

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "rt_9c8b7a6f5e4d3c2b1a0f9e8d7c6b5a4f"
}
```

Errors: `422 VALIDATION`, `401 UNAUTHORIZED` (invalid/expired/revoked token).

### POST /auth/logout

Revokes the given refresh token. Public (idempotent).

Request:

```json
{
  "refresh_token": "rt_4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c"
}
```

`204 No Content` — empty body.

### GET /auth/me

Returns the current user and tenant. **Auth required.**

`200 OK`:

```json
{
  "user": { /* User */ },
  "tenant": { /* Tenant */ }
}
```

Errors: `401 UNAUTHORIZED`.

---

## Platform Admin

Platform operators manage tenants across the whole platform. These endpoints use
a **separate token kind** (see [Platform admin access token claims](#authentication--multi-tenancy)):
the access token has `role: "platform_admin"` and **no `tid`**. All endpoints
except login require `Authorization: Bearer <access_token>` issued for a platform
admin; a normal tenant token is rejected with `403 FORBIDDEN`.

### POST /admin/auth/login

Authenticates a platform admin. Public.

Request:

```json
{
  "email": "ops@wa-dashboard.com",
  "password": "s3cr3tpassword"
}
```

`200 OK`:

```json
{
  "admin": { /* PlatformAdmin */ },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "rt_4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c"
}
```

Errors: `422 VALIDATION`, `401 UNAUTHORIZED` (bad credentials).

### GET /admin/tenants

Lists all tenants across the platform. Paginated. **Platform admin only.**

Query params: `limit`, `offset`.

`200 OK`:

```json
{
  "data": [ /* array of Tenant */ ],
  "page": { "limit": 20, "offset": 0, "total": 137 }
}
```

Errors: `401 UNAUTHORIZED`, `403 FORBIDDEN`.

### POST /admin/tenants

Provisions a new tenant **and** its initial owner (`admin`) user in one call.
**Platform admin only.**

Request:

```json
{
  "business_name": "Acme Inc.",
  "owner_email": "owner@acme.com",
  "owner_full_name": "Jane Doe",
  "owner_password": "s3cr3tpassword"
}
```

`201 Created`:

```json
{
  "tenant": { /* Tenant */ },
  "owner": { /* User */ }
}
```

Errors: `422 VALIDATION`, `401 UNAUTHORIZED`, `403 FORBIDDEN`, `409 CONFLICT`
(business name / owner email already registered).

### GET /admin/tenants/:id

Returns a single tenant by id. **Platform admin only.**

`200 OK`: a [Tenant](#tenant) object.

Errors: `401 UNAUTHORIZED`, `403 FORBIDDEN`, `404 NOT_FOUND`.

### PATCH /admin/tenants/:id

Updates a tenant's `status` to suspend or reactivate it. While a tenant is
`suspended`, its users are blocked with `403 TENANT_SUSPENDED`. **Platform admin only.**

Request:

```json
{
  "status": "suspended"
}
```

`status` must be one of `active` | `suspended`.

`200 OK`: the updated [Tenant](#tenant) object.

Errors: `422 VALIDATION`, `401 UNAUTHORIZED`, `403 FORBIDDEN`, `404 NOT_FOUND`.

---

## Tenant

All tenant endpoints require auth and operate on the caller's own tenant
(resolved from `tid` in the JWT).

### GET /tenant

Returns the current workspace. **Auth required (any role).**

`200 OK`: a [Tenant](#tenant) object.

### PATCH /tenant

Updates the current tenant. **Admin only.**

Request (all fields optional):

```json
{
  "business_name": "Acme International",
  "settings": { "timezone": "Asia/Jakarta" }
}
```

`200 OK`: the updated [Tenant](#tenant) object.

Errors: `422 VALIDATION`, `401 UNAUTHORIZED`, `403 FORBIDDEN`.

### GET /tenant/members

Lists users in the current tenant. Paginated. **Admin or supervisor.**

Query params: `limit`, `offset`.

`200 OK`:

```json
{
  "data": [ /* array of User */ ],
  "page": { "limit": 20, "offset": 0, "total": 3 }
}
```

Errors: `401 UNAUTHORIZED`, `403 FORBIDDEN`.

### POST /tenant/members

Invites/creates a user in the current tenant with a role. **Admin only.**

Request:

```json
{
  "email": "agent@acme.com",
  "full_name": "John Agent",
  "role": "agent",
  "password": "temp-password-123"
}
```

`201 Created`: the created [User](#user) object.

Errors: `422 VALIDATION`, `401 UNAUTHORIZED`, `403 FORBIDDEN`, `409 CONFLICT`.

---

## Broadcasts (reference vertical slice)

All broadcast endpoints require auth and are tenant-scoped.

### GET /broadcasts

Lists broadcasts for the current tenant. Paginated. **Auth required.**

Query params: `limit`, `offset`, `status` (optional filter; one of the
broadcast status enum values).

`200 OK`:

```json
{
  "data": [ /* array of Broadcast */ ],
  "page": { "limit": 20, "offset": 0, "total": 12 }
}
```

Errors: `401 UNAUTHORIZED`, `422 VALIDATION` (invalid `status` filter).

### POST /broadcasts

Creates a broadcast. **Admin or supervisor.**

Request:

```json
{
  "name": "May Promo Blast",
  "template_id": "tmpl_2f3a4b5c",
  "segment": "all_contacts",
  "scheduled_at": "2026-06-01T09:00:00Z"
}
```

`scheduled_at` is optional; omit (or `null`) to create a `draft`. When present,
the broadcast is created with status `scheduled`.

`201 Created`: the created [Broadcast](#broadcast) object.

Errors: `422 VALIDATION`, `401 UNAUTHORIZED`, `403 FORBIDDEN`.

### GET /broadcasts/:id

Returns a single broadcast in the current tenant. **Auth required.**

`200 OK`: a [Broadcast](#broadcast) object.

Errors: `401 UNAUTHORIZED`, `404 NOT_FOUND`.

---

## Skeleton / stub endpoints

These are defined in the contract so the frontend can integrate against stable
paths. They are **stubs**: backend returns empty/placeholder data for now.

### GET /inbox/conversations

**Stub** — returns an empty paginated list for now. **Auth required.**

`200 OK`:

```json
{
  "data": [],
  "page": { "limit": 20, "offset": 0, "total": 0 }
}
```

### GET /inbox/stream

**Stub** — Server-Sent Events stream for realtime inbox updates. **Auth required.**

Response headers: `Content-Type: text/event-stream`. The server emits events
whose `data` is a JSON object of this shape:

```json
{
  "type": "message.created",
  "conversation_id": "3a4b5c6d-7e8f-9a0b-1c2d-3e4f5a6b7c8d",
  "message": {
    "id": "5c6d7e8f-9a0b-1c2d-3e4f-5a6b7c8d9e0f",
    "body": "Hello, I need help with my order",
    "direction": "inbound",
    "created_at": "2026-05-31T16:15:00Z"
  }
}
```

`type` is an event discriminator (e.g. `message.created`, `conversation.updated`).
For now the stub may send only periodic keep-alive comments and no real events.

### GET /analytics/summary

**Stub** — returns placeholder aggregate metrics for the current tenant. **Auth required.**

`200 OK` (placeholder values):

```json
{
  "delivery_rate": 0,
  "open_rate": 0,
  "avg_response_time_seconds": 0,
  "conversation_volume": 0,
  "resolution_rate": 0
}
```

### GET /templates

**Stub** — returns an empty paginated list of WhatsApp message templates. **Auth required.**

`200 OK`:

```json
{
  "data": [],
  "page": { "limit": 20, "offset": 0, "total": 0 }
}
```
