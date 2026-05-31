package obs

import (
	"net/http"

	"github.com/beatfraps/wa-dashboard/backend/internal/shared/db"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/queue"
	"github.com/labstack/echo/v4"
)

type Health struct {
	pool     *db.Pool
	redisURL string
}

func NewHealth(pool *db.Pool, redisURL string) *Health {
	return &Health{pool: pool, redisURL: redisURL}
}

func (h *Health) Register(e *echo.Echo) {
	e.GET("/healthz", h.liveness)
	e.GET("/readyz", h.readiness)
}

func (h *Health) liveness(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Health) readiness(c echo.Context) error {
	ctx := c.Request().Context()
	if err := h.pool.Ping(ctx); err != nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}
	if err := queue.PingRedis(h.redisURL); err != nil {
		return c.NoContent(http.StatusServiceUnavailable)
	}
	return c.NoContent(http.StatusOK)
}
