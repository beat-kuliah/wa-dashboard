-- name: ListBroadcasts :many
SELECT * FROM broadcasts
WHERE tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountBroadcasts :one
SELECT COUNT(*) FROM broadcasts
WHERE tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'));

-- name: GetBroadcastByID :one
SELECT * FROM broadcasts
WHERE id = $1 AND tenant_id = $2;

-- name: CreateBroadcast :one
INSERT INTO broadcasts (
    tenant_id, name, template_id, status, scheduled_at, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateBroadcastStatus :one
UPDATE broadcasts
SET status = $3,
    sent_at = CASE WHEN $3 = 'sent' THEN NOW() ELSE sent_at END,
    error_message = $4,
    updated_at = NOW()
WHERE id = $1 AND tenant_id = $2
RETURNING *;

-- name: GetBroadcastForSend :one
SELECT * FROM broadcasts
WHERE id = $1 AND tenant_id = $2;

-- name: TemplateExistsForTenant :one
SELECT EXISTS(
    SELECT 1 FROM templates WHERE id = $1 AND tenant_id = $2
) AS exists;

-- name: CreateTemplate :one
INSERT INTO templates (tenant_id, name, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListTemplates :many
SELECT * FROM templates
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountTemplates :one
SELECT COUNT(*) FROM templates
WHERE tenant_id = $1;

-- name: ListConversations :many
SELECT * FROM conversations
WHERE tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountConversations :one
SELECT COUNT(*) FROM conversations
WHERE tenant_id = $1
  AND (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'));
