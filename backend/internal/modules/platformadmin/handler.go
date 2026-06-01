package platformadmin

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

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type provisionRequest struct {
	BusinessName  string `json:"business_name"`
	OwnerEmail      string `json:"owner_email"`
	OwnerFullName   string `json:"owner_full_name"`
	OwnerPassword   string `json:"owner_password"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) Login(c echo.Context) error {
	var req loginRequest
	if err := httpx.BindAndValidate(c, &req); err != nil {
		return err
	}
	result, err := h.svc.Login(c.Request().Context(), LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{
		"admin":         result.Admin,
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

func (h *Handler) ListTenants(c echo.Context) error {
	limit, offset := httpx.ParsePagination(c)
	tenants, total, err := h.svc.ListTenants(c.Request().Context(), limit, offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, httpx.PaginatedResponse{
		Data: tenants,
		Page: httpx.PageMeta{Limit: limit, Offset: offset, Total: total},
	})
}

func (h *Handler) GetTenant(c echo.Context) error {
	id, err := httpx.ParseUUIDParam(c, "id")
	if err != nil {
		return err
	}
	tenant, err := h.svc.GetTenant(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tenant)
}

func (h *Handler) ProvisionTenant(c echo.Context) error {
	var req provisionRequest
	if err := httpx.BindAndValidate(c, &req); err != nil {
		return err
	}
	result, err := h.svc.ProvisionTenant(c.Request().Context(), ProvisionInput{
		BusinessName:  req.BusinessName,
		OwnerEmail:    req.OwnerEmail,
		OwnerFullName: req.OwnerFullName,
		OwnerPassword: req.OwnerPassword,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"tenant": result.Tenant,
		"owner":  result.Owner,
	})
}

func (h *Handler) UpdateTenantStatus(c echo.Context) error {
	id, err := httpx.ParseUUIDParam(c, "id")
	if err != nil {
		return err
	}
	var req updateStatusRequest
	if err := httpx.BindAndValidate(c, &req); err != nil {
		return err
	}
	tenant, err := h.svc.UpdateTenantStatus(c.Request().Context(), id, req.Status)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tenant)
}

func RegisterRoutes(g *echo.Group, h *Handler) {
	authGroup := g.Group("/auth")
	authGroup.POST("/login", h.Login)

	tenants := g.Group("/tenants", auth.RequirePlatformAdmin)
	tenants.GET("", h.ListTenants)
	tenants.POST("", h.ProvisionTenant)
	tenants.GET("/:id", h.GetTenant)
	tenants.PATCH("/:id", h.UpdateTenantStatus)
}
