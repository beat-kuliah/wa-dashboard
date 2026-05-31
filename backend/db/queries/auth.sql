-- name: CreateTenant :one
INSERT INTO tenants (name, slug)
VALUES ($1, $2)
RETURNING *;

-- name: GetTenantByID :one
SELECT * FROM tenants
WHERE id = $1;

-- name: GetTenantBySlug :one
SELECT * FROM tenants
WHERE slug = $1;

-- name: UpdateTenant :one
UPDATE tenants
SET name = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (tenant_id, email, password_hash, full_name, roles)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND tenant_id = $2;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByEmailAndTenant :one
SELECT * FROM users
WHERE tenant_id = $1 AND email = $2;

-- name: ListUsersByTenant :many
SELECT * FROM users
WHERE tenant_id = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: CountUsersByTenant :one
SELECT COUNT(*) FROM users
WHERE tenant_id = $1;

-- name: GetUserByIDOnly :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT * FROM refresh_tokens
WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token_hash = $1 AND revoked_at IS NULL;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE user_id = $1 AND revoked_at IS NULL;
