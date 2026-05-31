package tenantmod

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{q: sqlc.NewFromPool(pool)}
}

func (r *Repository) GetTenant(ctx context.Context, tenantID uuid.UUID) (sqlc.Tenant, error) {
	return r.q.GetTenantByID(ctx, tenantID)
}

func (r *Repository) UpdateTenant(ctx context.Context, tenantID uuid.UUID, name string) (sqlc.Tenant, error) {
	return r.q.UpdateTenant(ctx, sqlc.UpdateTenantParams{ID: tenantID, Name: name})
}

func (r *Repository) ListMembers(ctx context.Context, tenantID uuid.UUID, limit, offset int32) ([]sqlc.User, int64, error) {
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

func (r *Repository) CreateMember(ctx context.Context, tenantID uuid.UUID, email, passwordHash, fullName string, roles []string) (sqlc.User, error) {
	return r.q.CreateUser(ctx, sqlc.CreateUserParams{
		TenantID:     tenantID,
		Email:        strings.ToLower(strings.TrimSpace(email)),
		PasswordHash: passwordHash,
		FullName:     fullName,
		Roles:        roles,
	})
}

func (r *Repository) MemberExists(ctx context.Context, tenantID uuid.UUID, email string) (bool, error) {
	_, err := r.q.GetUserByEmailAndTenant(ctx, sqlc.GetUserByEmailAndTenantParams{
		TenantID: tenantID,
		Email:    strings.ToLower(strings.TrimSpace(email)),
	})
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

func MapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.NotFound("Tenant not found")
	}
	return err
}

type TenantDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func TenantFromRow(t sqlc.Tenant) TenantDTO {
	return TenantDTO{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt.UTC(),
		UpdatedAt: t.UpdatedAt.UTC(),
	}
}

type UpdateTenantInput struct {
	Name string
}

type AddMemberInput struct {
	Email    string
	FullName string
	Roles    []string
}

func ValidateUpdateTenant(name string) error {
	if strings.TrimSpace(name) == "" {
		return apperrors.Validation("name is required")
	}
	return nil
}

func tenantIDFromContext(ctx context.Context) (uuid.UUID, error) {
	rc, ok := auth.RequestContextFrom(ctx)
	if !ok {
		return uuid.Nil, apperrors.Unauthorized("")
	}
	return rc.TenantID, nil
}
