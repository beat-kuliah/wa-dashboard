package tenantmod

import (
	"net/http"

	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/httpx"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type updateTenantRequest struct {
	Name string `json:"name"`
}

type addMemberRequest struct {
	Email    string   `json:"email"`
	FullName string   `json:"full_name"`
	Roles    []string `json:"roles"`
}

func (h *Handler) GetTenant(c echo.Context) error {
	tenant, err := h.svc.GetCurrentTenant(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"tenant": tenant})
}

func (h *Handler) UpdateTenant(c echo.Context) error {
	var req updateTenantRequest
	if err := httpx.BindAndValidate(c, &req); err != nil {
		return err
	}
	tenant, err := h.svc.UpdateTenant(c.Request().Context(), UpdateTenantInput{Name: req.Name})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"tenant": tenant})
}

func (h *Handler) ListMembers(c echo.Context) error {
	limit, offset := httpx.ParsePagination(c)
	members, total, err := h.svc.ListMembers(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, httpx.PaginatedResponse{
		Data: members,
		Page: httpx.PageMeta{Limit: limit, Offset: offset, Total: total},
	})
}

func (h *Handler) AddMember(c echo.Context) error {
	var req addMemberRequest
	if err := httpx.BindAndValidate(c, &req); err != nil {
		return err
	}
	user, err := h.svc.AddMember(c.Request().Context(), AddMemberInput{
		Email:    req.Email,
		FullName: req.FullName,
		Roles:    req.Roles,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]any{"user": user})
}

func TenantContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rc := httpx.RequestContext(c)
			if rc == nil {
				return next(c)
			}
			ctx := auth.WithRequestContext(c.Request().Context(), rc)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.Use(auth.RequireAuth, TenantContext())
	g.GET("", h.GetTenant)
	g.PATCH("", h.UpdateTenant, auth.RequireRole(auth.RoleAdmin))
	g.GET("/members", h.ListMembers, auth.RequireRole(auth.RoleAdmin, auth.RoleSupervisor))
	g.POST("/members", h.AddMember, auth.RequireRole(auth.RoleAdmin))
}
