package authmod

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

type registerRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	FullName     string `json:"full_name"`
	BusinessName string `json:"business_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Register(c echo.Context) error {
	var req registerRequest
	if err := httpx.BindAndValidate(c, &req); err != nil {
		return err
	}
	result, err := h.svc.Register(c.Request().Context(), RegisterInput{
		Email:        req.Email,
		Password:     req.Password,
		FullName:     req.FullName,
		BusinessName: req.BusinessName,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]any{
		"user":   result.User,
		"tenant": result.Tenant,
		"tokens": result.Tokens,
	})
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
		"user":   result.User,
		"tokens": result.Tokens,
	})
}

func (h *Handler) Refresh(c echo.Context) error {
	var req refreshRequest
	if err := httpx.BindAndValidate(c, &req); err != nil {
		return err
	}
	tokens, err := h.svc.Refresh(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"tokens": tokens})
}

func (h *Handler) Logout(c echo.Context) error {
	var req logoutRequest
	if err := httpx.BindAndValidate(c, &req); err != nil {
		return err
	}
	if err := h.svc.Logout(c.Request().Context(), req.RefreshToken); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Me(c echo.Context) error {
	user, err := h.svc.Me(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]any{"user": user})
}

func RegisterRoutes(g *echo.Group, h *Handler) {
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.POST("/refresh", h.Refresh)
	g.POST("/logout", h.Logout)
	g.GET("/me", h.Me, auth.RequireAuth)
}
