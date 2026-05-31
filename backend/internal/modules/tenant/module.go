package tenantmod

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Handler *Handler
	Service *Service
}

func NewModule(pool *pgxpool.Pool) *Module {
	repo := NewRepository(pool)
	svc := NewService(repo)
	handler := NewHandler(svc)
	return &Module{Handler: handler, Service: svc}
}

func (m *Module) RegisterRoutes(g *echo.Group) {
	RegisterRoutes(g, m.Handler)
}
