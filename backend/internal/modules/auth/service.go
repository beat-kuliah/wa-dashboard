package authmod

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo         *Repository
	tokens       *auth.TokenService
	refreshTTL   time.Duration
	accessExpiry int
}

func NewService(repo *Repository, tokens *auth.TokenService, refreshTTL time.Duration, accessExpiry int) *Service {
	return &Service{
		repo:         repo,
		tokens:       tokens,
		refreshTTL:   refreshTTL,
		accessExpiry: accessExpiry,
	}
}

type RegisterInput struct {
	Email        string
	Password     string
	FullName     string
	BusinessName string
}

type LoginInput struct {
	Email    string
	Password string
}

type RegisterResult struct {
	User   UserDTO
	Tenant TenantDTO
	Tokens TokensDTO
}

type LoginResult struct {
	User   UserDTO
	Tokens TokensDTO
}

func (s *Service) Register(ctx context.Context, input RegisterInput) (*RegisterResult, error) {
	if err := ValidateRegisterInput(input.Email, input.Password, input.FullName, input.BusinessName); err != nil {
		return nil, err
	}

	if _, err := s.repo.GetUserByEmail(ctx, input.Email); err == nil {
		return nil, apperrors.Conflict("email already exists")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, apperrors.Internal("failed to hash password", err)
	}

	slug := Slugify(input.BusinessName)
	tenant, user, err := s.repo.CreateTenantWithAdmin(ctx, input.BusinessName, slug, input.Email, passwordHash, input.FullName)
	if err != nil {
		if IsUniqueViolation(err) {
			return nil, apperrors.Conflict("email or slug already exists")
		}
		return nil, apperrors.Internal("failed to register", err)
	}

	tokens, err := s.issueTokens(ctx, user.ID, user.TenantID, user.Roles)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{
		User:   UserFromRow(user),
		Tenant: TenantFromRow(tenant),
		Tokens: tokens,
	}, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	if err := ValidateLoginInput(input.Email, input.Password); err != nil {
		return nil, err
	}

	user, err := s.repo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.Unauthorized("invalid credentials")
		}
		return nil, err
	}

	if err := CheckPassword(user.PasswordHash, input.Password); err != nil {
		return nil, apperrors.Unauthorized("invalid credentials")
	}

	tokens, err := s.issueTokens(ctx, user.ID, user.TenantID, user.Roles)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:   UserFromRow(user),
		Tokens: tokens,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokensDTO, error) {
	if err := ValidateRefreshInput(refreshToken); err != nil {
		return nil, err
	}

	hash := hashToken(refreshToken)
	stored, err := s.repo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.Unauthorized("invalid or revoked refresh token")
		}
		return nil, err
	}

	if err := s.repo.RevokeRefreshToken(ctx, hash); err != nil {
		return nil, apperrors.Internal("failed to revoke refresh token", err)
	}

	user, err := s.repo.GetUserByIDOnly(ctx, stored.UserID)
	if err != nil {
		return nil, apperrors.Unauthorized("invalid or revoked refresh token")
	}

	tokens, err := s.issueTokens(ctx, user.ID, user.TenantID, user.Roles)
	if err != nil {
		return nil, err
	}
	return &tokens, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	if err := ValidateRefreshInput(refreshToken); err != nil {
		return err
	}
	hash := hashToken(refreshToken)
	if err := s.repo.RevokeRefreshToken(ctx, hash); err != nil {
		return apperrors.Internal("failed to logout", err)
	}
	return nil
}

func (s *Service) Me(ctx context.Context) (*UserDTO, error) {
	rc, ok := auth.RequestContextFrom(ctx)
	if !ok {
		return nil, apperrors.Unauthorized("")
	}
	user, err := s.repo.GetUserByID(ctx, rc.UserID, rc.TenantID)
	if err != nil {
		return nil, MapNotFound(err, "User")
	}
	dto := UserFromRow(user)
	return &dto, nil
}

func (s *Service) issueTokens(ctx context.Context, userID, tenantID uuid.UUID, roles []string) (TokensDTO, error) {
	access, err := s.tokens.IssueAccessToken(userID, tenantID, roles)
	if err != nil {
		return TokensDTO{}, apperrors.Internal("failed to issue access token", err)
	}
	raw, hash, err := GenerateRefreshToken()
	if err != nil {
		return TokensDTO{}, apperrors.Internal("failed to generate refresh token", err)
	}
	expiresAt := time.Now().UTC().Add(s.refreshTTL)
	if err := s.repo.CreateRefreshToken(ctx, userID, hash, expiresAt); err != nil {
		return TokensDTO{}, apperrors.Internal("failed to store refresh token", err)
	}
	return TokensDTO{
		AccessToken:  access,
		RefreshToken: raw,
		ExpiresIn:    s.accessExpiry,
	}, nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
