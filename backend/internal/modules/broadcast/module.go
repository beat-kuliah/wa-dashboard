package broadcastmod

import (
	"log/slog"

	"github.com/beatfraps/wa-dashboard/backend/internal/shared/queue"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Handler *Handler
	Service *Service
}

func NewModule(pool *pgxpool.Pool, q *queue.Client) *Module {
	repo := NewRepository(pool)
	svc := NewService(repo, q)
	handler := NewHandler(svc)
	return &Module{Handler: handler, Service: svc}
}

func (m *Module) RegisterRoutes(g *echo.Group) {
	RegisterRoutes(g, m.Handler)
}

func (m *Module) RegisterWorkerHandlers(worker *queue.Worker, logger *slog.Logger) {
	RegisterWorkerHandlers(worker, m.Service, logger)
}
