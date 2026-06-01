DROP TABLE IF EXISTS platform_admin_refresh_tokens;
ALTER TABLE tenants DROP COLUMN IF EXISTS status;
DROP TABLE IF EXISTS platform_admins;
