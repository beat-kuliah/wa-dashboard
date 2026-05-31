package inboxmod

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

func (r *Repository) ListConversations(ctx context.Context, tenantID uuid.UUID, limit, offset int32, status *string) ([]sqlc.Conversation, int64, error) {
	items, err := r.q.ListConversations(ctx, sqlc.ListConversationsParams{
		TenantID: tenantID,
		Limit:    limit,
		Offset:   offset,
		Status:   status,
	})
	if err != nil {
		return nil, 0, err
	}
	total, err := r.q.CountConversations(ctx, sqlc.CountConversationsParams{TenantID: tenantID, Status: status})
	return items, total, err
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListConversations(ctx context.Context, limit, offset int32, status *string) ([]any, int64, error) {
	rc, ok := auth.RequestContextFrom(ctx)
	if !ok {
		return nil, 0, apperrors.Unauthorized("")
	}
	_, total, err := s.repo.ListConversations(ctx, rc.TenantID, limit, offset, status)
	if err != nil {
		return nil, 0, err
	}
	return []any{}, total, nil
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListConversations(c echo.Context) error {
	limit, offset := httpx.ParsePagination(c)
	status := httpx.OptionalQueryString(c, "status")
	_, total, err := h.svc.ListConversations(c.Request().Context(), limit, offset, status)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, httpx.PaginatedResponse{
		Data: []any{},
		Page: httpx.PageMeta{Limit: limit, Offset: offset, Total: total},
	})
}

func (h *Handler) Stream(c echo.Context) error {
	rc := httpx.RequestContext(c)
	if rc == nil {
		return apperrors.Unauthorized("")
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	ctx := c.Request().Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if _, err := fmt.Fprintf(c.Response(), "event: ping\ndata: {}\n\n"); err != nil {
				return nil
			}
			flusher.Flush()
		}
	}
}

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.Use(auth.RequireAuth)
	g.GET("/conversations", h.ListConversations)
	g.GET("/stream", h.Stream)
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
