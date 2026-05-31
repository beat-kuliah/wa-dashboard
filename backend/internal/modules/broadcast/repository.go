package broadcastmod

import (
	"context"
	"errors"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{q: sqlc.NewFromPool(pool)}
}

func (r *Repository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int32, status *string) ([]sqlc.Broadcast, int64, error) {
	items, err := r.q.ListBroadcasts(ctx, sqlc.ListBroadcastsParams{
		TenantID: tenantID,
		Limit:    limit,
		Offset:   offset,
		Status:   status,
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountBroadcasts(ctx, sqlc.CountBroadcastsParams{TenantID: tenantID, Status: status})
	return items, total, err
}

func (r *Repository) GetByID(ctx context.Context, tenantID, id uuid.UUID) (sqlc.Broadcast, error) {
	return r.q.GetBroadcastByID(ctx, sqlc.GetBroadcastByIDParams{ID: id, TenantID: tenantID})
}

func (r *Repository) Create(ctx context.Context, params sqlc.CreateBroadcastParams) (sqlc.Broadcast, error) {
	return r.q.CreateBroadcast(ctx, params)
}

func (r *Repository) UpdateStatus(ctx context.Context, tenantID, id uuid.UUID, status string, errMsg *string) (sqlc.Broadcast, error) {
	return r.q.UpdateBroadcastStatus(ctx, sqlc.UpdateBroadcastStatusParams{
		ID: id, TenantID: tenantID, Status: status, ErrorMessage: errMsg,
	})
}

func (r *Repository) TemplateExists(ctx context.Context, tenantID, templateID uuid.UUID) (bool, error) {
	return r.q.TemplateExistsForTenant(ctx, sqlc.TemplateExistsForTenantParams{ID: templateID, TenantID: tenantID})
}

func MapNotFound(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.NotFound("Broadcast not found")
	}
	return err
}
