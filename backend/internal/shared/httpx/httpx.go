package httpx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/logx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Server struct {
	Echo *echo.Echo
}

type ServerConfig struct {
	CORSOrigins []string
	Logger      *slog.Logger
	Tokens      *auth.TokenService
}

func NewServer(cfg ServerConfig) *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = errorHandler

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("1M"))
	e.Use(corsMiddleware(cfg.CORSOrigins))
	e.Use(requestLogger(cfg.Logger))
	if cfg.Tokens != nil {
		e.Use(auth.JWTAuth(cfg.Tokens))
	}

	return &Server{Echo: e}
}

func corsMiddleware(origins []string) echo.MiddlewareFunc {
	allowAll := len(origins) == 0
	allowed := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		allowed[o] = struct{}{}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			origin := c.Request().Header.Get("Origin")
			if origin != "" {
				if allowAll {
					c.Response().Header().Set("Access-Control-Allow-Origin", origin)
				} else if _, ok := allowed[origin]; ok {
					c.Response().Header().Set("Access-Control-Allow-Origin", origin)
				}
				c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
				c.Response().Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept")
				c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			}
			if c.Request().Method == http.MethodOptions {
				return c.NoContent(http.StatusNoContent)
			}
			return next(c)
		}
	}
}

func requestLogger(base *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			logger := base.With("request_id", reqID, "method", c.Request().Method, "path", c.Request().URL.Path)
			ctx := logx.WithContext(c.Request().Context(), logger)
			c.SetRequest(c.Request().WithContext(ctx))

			err := next(c)
			status := c.Response().Status
			if err != nil {
				if he, ok := err.(*echo.HTTPError); ok {
					status = he.Code
				} else if appErr := apperrors.AsAppError(err); appErr != nil {
					status = appErr.Status
				}
			}
			logger.Info("request completed",
				"status", status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
			return err
		}
	}
}

func errorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	appErr := apperrors.AsAppError(err)
	if appErr.Status >= http.StatusInternalServerError {
		logx.FromContext(c.Request().Context()).Error("unhandled error", "error", err)
	}

	_ = c.JSON(appErr.Status, ErrorResponse{
		Error: ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}

func BindAndValidate(c echo.Context, dst any) error {
	if err := c.Bind(dst); err != nil {
		return apperrors.Validation("invalid request body")
	}
	return nil
}

type PaginatedResponse struct {
	Data any       `json:"data"`
	Page PageMeta  `json:"page"`
}

type PageMeta struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
	Total  int64 `json:"total"`
}

func ParsePagination(c echo.Context) (limit, offset int32) {
	limit = 20
	offset = 0
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := parseInt32(l); err == nil && parsed >= 1 && parsed <= 100 {
			limit = parsed
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := parseInt32(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	return limit, offset
}

func parseInt32(s string) (int32, error) {
	var v int32
	_, err := fmt.Sscanf(s, "%d", &v)
	return v, err
}

func RequestContext(c echo.Context) *auth.RequestContext {
	if rc, ok := c.Get("request_context").(*auth.RequestContext); ok {
		return rc
	}
	return nil
}

func EchoContext(ctx context.Context) context.Context {
	return ctx
}

func ParseUUIDParam(c echo.Context, name string) (uuid.UUID, error) {
	raw := c.Param(name)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, apperrors.Validation(name + " must be a valid UUID")
	}
	return id, nil
}

func OptionalQueryString(c echo.Context, name string) *string {
	v := strings.TrimSpace(c.QueryParam(name))
	if v == "" {
		return nil
	}
	return &v
}
