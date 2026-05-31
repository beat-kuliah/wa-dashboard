package auth

import (
	"context"
	"strings"
	"time"

	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Role string

const (
	RoleAdmin      Role = "admin"
	RoleSupervisor Role = "supervisor"
	RoleAgent      Role = "agent"
)

type RequestContext struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
	Roles    []string
}

type ctxKey struct{}

type Claims struct {
	TID   string   `json:"tid"`
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

type TokenService struct {
	secret    []byte
	accessTTL time.Duration
}

func NewTokenService(secret string, accessTTL time.Duration) *TokenService {
	return &TokenService{secret: []byte(secret), accessTTL: accessTTL}
}

func (s *TokenService) IssueAccessToken(userID, tenantID uuid.UUID, roles []string) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		TID:   tenantID.String(),
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *TokenService) ParseAccessToken(tokenString string) (*RequestContext, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.Unauthorized("invalid token")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, apperrors.Unauthorized("invalid or expired token")
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, apperrors.Unauthorized("invalid token")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid token")
	}
	tenantID, err := uuid.Parse(claims.TID)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid token")
	}
	return &RequestContext{
		UserID:   userID,
		TenantID: tenantID,
		Roles:    claims.Roles,
	}, nil
}

func WithRequestContext(ctx context.Context, rc *RequestContext) context.Context {
	return context.WithValue(ctx, ctxKey{}, rc)
}

func RequestContextFrom(ctx context.Context) (*RequestContext, bool) {
	rc, ok := ctx.Value(ctxKey{}).(*RequestContext)
	return rc, ok
}

func RequireAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		rc, ok := c.Get("request_context").(*RequestContext)
		if !ok || rc == nil {
			return apperrors.Unauthorized("")
		}
		ctx := WithRequestContext(c.Request().Context(), rc)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

func JWTAuth(tokens *TokenService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return next(c)
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return apperrors.Unauthorized("invalid authorization header")
			}
			rc, err := tokens.ParseAccessToken(parts[1])
			if err != nil {
				return err
			}
			c.Set("request_context", rc)
			ctx := WithRequestContext(c.Request().Context(), rc)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func HasRole(roles []string, required ...Role) bool {
	allowed := make(map[string]struct{}, len(required))
	for _, r := range required {
		allowed[string(r)] = struct{}{}
	}
	for _, role := range roles {
		if _, ok := allowed[role]; ok {
			return true
		}
	}
	return false
}

func RequireRole(required ...Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			rc, ok := c.Get("request_context").(*RequestContext)
			if !ok || rc == nil {
				return apperrors.Unauthorized("")
			}
			if !HasRole(rc.Roles, required...) {
				return apperrors.Forbidden("")
			}
			return next(c)
		}
	}
}

func AssertRole(ctx context.Context, required ...Role) error {
	rc, ok := RequestContextFrom(ctx)
	if !ok {
		return apperrors.Unauthorized("")
	}
	if !HasRole(rc.Roles, required...) {
		return apperrors.Forbidden("")
	}
	return nil
}
