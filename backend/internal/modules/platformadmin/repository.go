package platformadmin

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q    *sqlc.Queries
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{q: sqlc.NewFromPool(pool), pool: pool}
}

func (r *Repository) GetPlatformAdminByEmail(ctx context.Context, email string) (sqlc.PlatformAdmin, error) {
	return r.q.GetPlatformAdminByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
}

func (r *Repository) GetPlatformAdminByID(ctx context.Context, id uuid.UUID) (sqlc.PlatformAdmin, error) {
	return r.q.GetPlatformAdminByID(ctx, id)
}

func (r *Repository) CreatePlatformAdmin(ctx context.Context, email, passwordHash, fullName string) (sqlc.PlatformAdmin, error) {
	return r.q.CreatePlatformAdmin(ctx, sqlc.CreatePlatformAdminParams{
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		FullName:     fullName,
	})
}

func (r *Repository) CreatePlatformAdminRefreshToken(ctx context.Context, adminID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.q.CreatePlatformAdminRefreshToken(ctx, sqlc.CreatePlatformAdminRefreshTokenParams{
		PlatformAdminID: adminID,
		TokenHash:       tokenHash,
		ExpiresAt:       expiresAt,
	})
	return err
}

func (r *Repository) GetPlatformAdminRefreshTokenByHash(ctx context.Context, tokenHash string) (sqlc.PlatformAdminRefreshToken, error) {
	return r.q.GetPlatformAdminRefreshTokenByHash(ctx, tokenHash)
}

func (r *Repository) RevokePlatformAdminRefreshToken(ctx context.Context, tokenHash string) error {
	return r.q.RevokePlatformAdminRefreshToken(ctx, tokenHash)
}

func (r *Repository) ListTenants(ctx context.Context, limit, offset int32) ([]sqlc.Tenant, int64, error) {
	tenants, err := r.q.ListTenants(ctx, sqlc.ListTenantsParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountTenants(ctx)
	return tenants, total, err
}

func (r *Repository) GetTenantByID(ctx context.Context, id uuid.UUID) (sqlc.Tenant, error) {
	return r.q.GetTenantByID(ctx, id)
}

func (r *Repository) UpdateTenantStatus(ctx context.Context, id uuid.UUID, status string) (sqlc.Tenant, error) {
	return r.q.UpdateTenantStatus(ctx, sqlc.UpdateTenantStatusParams{ID: id, Status: status})
}

func (r *Repository) CreateTenantWithOwner(ctx context.Context, businessName, slug, ownerEmail, passwordHash, ownerFullName string) (sqlc.Tenant, sqlc.User, error) {
	var tenant sqlc.Tenant
	var user sqlc.User
	err := db.WithTx(ctx, r.pool, func(tx pgx.Tx) error {
		qtx := r.q.WithTx(tx)
		var err error
		tenant, err = qtx.CreateTenant(ctx, sqlc.CreateTenantParams{Name: businessName, Slug: slug})
		if err != nil {
			return err
		}
		user, err = qtx.CreateUser(ctx, sqlc.CreateUserParams{
			TenantID:     tenant.ID,
			Email:        strings.ToLower(strings.TrimSpace(ownerEmail)),
			PasswordHash: passwordHash,
			FullName:     ownerFullName,
			Roles:        []string{string(auth.RoleAdmin)},
		})
		return err
	})
	return tenant, user, err
}

func (r *Repository) PlatformAdminExists(ctx context.Context, email string) (bool, error) {
	_, err := r.q.GetPlatformAdminByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return false, err
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
