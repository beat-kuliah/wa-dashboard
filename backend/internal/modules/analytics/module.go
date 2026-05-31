package analyticsmod

import (
	"net/http"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	"github.com/labstack/echo/v4"
)

type SummaryDTO struct {
	DeliveryRate           float64   `json:"delivery_rate"`
	OpenRate               float64   `json:"open_rate"`
	AvgResponseTimeSeconds float64   `json:"avg_response_time_seconds"`
	ConversationVolume     int       `json:"conversation_volume"`
	ResolutionRate         float64   `json:"resolution_rate"`
	PeriodStart            time.Time `json:"period_start"`
	PeriodEnd              time.Time `json:"period_end"`
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Summary(from, to *time.Time) SummaryDTO {
	now := time.Now().UTC()
	end := now
	start := now.AddDate(0, 0, -30)
	if from != nil {
		start = from.UTC()
	}
	if to != nil {
		end = to.UTC()
	}
	return SummaryDTO{
		DeliveryRate:           0,
		OpenRate:               0,
		AvgResponseTimeSeconds: 0,
		ConversationVolume:     0,
		ResolutionRate:         0,
		PeriodStart:            start,
		PeriodEnd:              end,
	}
}

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Summary(c echo.Context) error {
	var from, to *time.Time
	if v := c.QueryParam("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			from = &t
		}
	}
	if v := c.QueryParam("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			to = &t
		}
	}
	summary := h.svc.Summary(from, to)
	return c.JSON(http.StatusOK, map[string]any{"summary": summary})
}

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.Use(auth.RequireAuth)
	g.GET("/summary", h.Summary)
}

type Module struct {
	Handler *Handler
}

func NewModule() *Module {
	svc := NewService()
	handler := NewHandler(svc)
	return &Module{Handler: handler}
}

func (m *Module) RegisterRoutes(g *echo.Group) {
	RegisterRoutes(g, m.Handler)
}
