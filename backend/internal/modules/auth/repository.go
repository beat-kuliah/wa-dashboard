package authmod

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/db"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
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

func (r *Repository) CreateTenantWithAdmin(ctx context.Context, tenantName, slug, email, passwordHash, fullName string) (sqlc.Tenant, sqlc.User, error) {
	var tenant sqlc.Tenant
	var user sqlc.User
	err := db.WithTx(ctx, r.pool, func(tx pgx.Tx) error {
		qtx := r.q.WithTx(tx)
		var err error
		tenant, err = qtx.CreateTenant(ctx, sqlc.CreateTenantParams{Name: tenantName, Slug: slug})
		if err != nil {
			return err
		}
		user, err = qtx.CreateUser(ctx, sqlc.CreateUserParams{
			TenantID:     tenant.ID,
			Email:        strings.ToLower(strings.TrimSpace(email)),
			PasswordHash: passwordHash,
			FullName:     fullName,
			Roles:        []string{string(auth.RoleAdmin)},
		})
		return err
	})
	return tenant, user, err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	return r.q.GetUserByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
}

func (r *Repository) GetUserByID(ctx context.Context, userID, tenantID uuid.UUID) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, sqlc.GetUserByIDParams{ID: userID, TenantID: tenantID})
}

func (r *Repository) GetUserByEmailAndTenant(ctx context.Context, tenantID uuid.UUID, email string) (sqlc.User, error) {
	return r.q.GetUserByEmailAndTenant(ctx, sqlc.GetUserByEmailAndTenantParams{
		TenantID: tenantID,
		Email:    strings.ToLower(strings.TrimSpace(email)),
	})
}

func (r *Repository) CreateUser(ctx context.Context, tenantID uuid.UUID, email, passwordHash, fullName string, roles []string) (sqlc.User, error) {
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{
		TenantID:     tenantID,
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		FullName:     fullName,
		Roles:        roles,
	})
}

func (r *Repository) ListUsersByTenant(ctx context.Context, tenantID uuid.UUID, limit, offset int32) ([]sqlc.User, int64, error) {
	users, err := r.q.ListUsersByTenant(ctx, sqlc.ListUsersByTenantParams{
		TenantID: tenantID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountUsersByTenant(ctx, tenantID)
	return users, total, err
}

func (r *Repository) GetUserByIDOnly(ctx context.Context, userID uuid.UUID) (sqlc.User, error) {
	return r.q.GetUserByIDOnly(ctx, userID)
}

func (r *Repository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.q.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
	return err
}

func (r *Repository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (sqlc.RefreshToken, error) {
	return r.q.GetRefreshTokenByHash(ctx, tokenHash)
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return r.q.RevokeRefreshToken(ctx, tokenHash)
}

func (r *Repository) GetTenantByID(ctx context.Context, tenantID uuid.UUID) (sqlc.Tenant, error) {
	return r.q.GetTenantByID(ctx, tenantID)
}

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}

func MapNotFound(err error, resource string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.NotFound(resource + " not found")
	}
	return err
}
