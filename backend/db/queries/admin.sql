-- name: CreatePlatformAdmin :one
INSERT INTO platform_admins (email, password_hash, full_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPlatformAdminByEmail :one
SELECT * FROM platform_admins
WHERE email = $1;

-- name: GetPlatformAdminByID :one
SELECT * FROM platform_admins
WHERE id = $1;

-- name: ListTenants :many
SELECT * FROM tenants
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountTenants :one
SELECT COUNT(*) FROM tenants;

-- name: UpdateTenantStatus :one
UPDATE tenants
SET status = $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreatePlatformAdminRefreshToken :one
INSERT INTO platform_admin_refresh_tokens (platform_admin_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPlatformAdminRefreshTokenByHash :one
SELECT * FROM platform_admin_refresh_tokens
WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > NOW();

-- name: RevokePlatformAdminRefreshToken :exec
UPDATE platform_admin_refresh_tokens
SET revoked_at = NOW()
WHERE token_hash = $1 AND revoked_at IS NULL;
