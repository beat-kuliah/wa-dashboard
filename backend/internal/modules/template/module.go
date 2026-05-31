package templatemod

import (
	"context"
	"net/http"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/httpx"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Repository struct {
	q *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{q: sqlc.NewFromPool(pool)}
}

func (r *Repository) List(ctx context.Context, tenantID uuid.UUID, limit, offset int32) ([]sqlc.Template, int64, error) {
	items, err := r.q.ListTemplates(ctx, sqlc.ListTemplatesParams{
		TenantID: tenantID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountTemplates(ctx, tenantID)
	return items, total, err
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, limit, offset int32) ([]sqlc.Template, int64, error) {
	rc, ok := auth.RequestContextFrom(ctx)
	if !ok {
		return nil, 0, apperrors.Unauthorized("")
	}
	return s.repo.List(ctx, rc.TenantID, limit, offset)
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(c echo.Context) error {
	limit, offset := httpx.ParsePagination(c)
	items, total, err := h.svc.List(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}
	data := make([]any, 0, len(items))
	for _, item := range items {
		data = append(data, map[string]any{
			"id":         item.ID,
			"tenant_id":  item.TenantID,
			"name":       item.Name,
			"status":     item.Status,
			"created_at": item.CreatedAt.UTC(),
			"updated_at": item.UpdatedAt.UTC(),
		})
	}
	return c.JSON(http.StatusOK, httpx.PaginatedResponse{
		Data: data,
		Page: httpx.PageMeta{Limit: limit, Offset: offset, Total: total},
	})
}

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.Use(auth.RequireAuth)
	g.GET("", h.List)
}

type Module struct {
	Handler *Handler
}

func NewModule(pool *pgxpool.Pool) *Module {
	repo := NewRepository(pool)
	svc := NewService(repo)
	handler := NewHandler(svc)
	return &Module{Handler: handler}
}

func (m *Module) RegisterRoutes(g *echo.Group) {
	RegisterRoutes(g, m.Handler)
}
