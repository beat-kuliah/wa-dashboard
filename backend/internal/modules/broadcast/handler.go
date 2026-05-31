package broadcastmod

import (
	"net/http"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/httpx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type createBroadcastRequest struct {
	Name        string     `json:"name"`
	TemplateID  uuid.UUID  `json:"template_id"`
	ScheduledAt *time.Time `json:"scheduled_at"`
}

func (h *Handler) List(c echo.Context) error {
	limit, offset := httpx.ParsePagination(c)
	status := httpx.OptionalQueryString(c, "status")
	items, total, err := h.svc.List(c.Request().Context(), limit, offset, status)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, httpx.PaginatedResponse{
		Data: items,
		Page: httpx.PageMeta{Limit: limit, Offset: offset, Total: total},
	})
}

func (h *Handler) Get(c echo.Context) error {
	id, err := httpx.ParseUUIDParam(c, "id")
	if err != nil {
		return err
	}
	item, err := h.svc.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"broadcast": item})
}

func (h *Handler) Create(c echo.Context) error {
	var req createBroadcastRequest
	if err := httpx.BindAndValidate(c, &req); err != nil {
		return err
	}
	item, err := h.svc.Create(c.Request().Context(), CreateBroadcastInput{
		Name:        req.Name,
		TemplateID:  req.TemplateID,
		ScheduledAt: req.ScheduledAt,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]any{"broadcast": item})
}

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.Use(auth.RequireAuth)
	g.GET("", h.List)
	g.POST("", h.Create, auth.RequireRole(auth.RoleAdmin, auth.RoleSupervisor))
	g.GET("/:id", h.Get)
}
