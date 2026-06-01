package authmod

import (
	"time"

	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Handler *Handler
	Service *Service
}

func NewModule(pool *pgxpool.Pool, tokens *auth.TokenService, refreshTTL time.Duration, accessExpiry int, publicRegistrationEnabled bool) *Module {
	repo := NewRepository(pool)
	svc := NewService(repo, tokens, refreshTTL, accessExpiry)
	handler := NewHandler(svc, publicRegistrationEnabled)
	return &Module{Handler: handler, Service: svc}
}

func (m *Module) RegisterRoutes(g *echo.Group) {
	RegisterRoutes(g, m.Handler)
}
