package platformadmin

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

type LoginInput struct {
	Email    string
	Password string
}

type LoginResult struct {
	Admin        PlatformAdminDTO
	AccessToken  string
	RefreshToken string
}

type ProvisionInput struct {
	BusinessName   string
	OwnerEmail     string
	OwnerFullName  string
	OwnerPassword  string
}

type ProvisionResult struct {
	Tenant TenantAdminDTO
	Owner  OwnerDTO
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	if err := ValidateLoginInput(input.Email, input.Password); err != nil {
		return nil, err
	}
	admin, err := s.repo.GetPlatformAdminByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.Unauthorized("invalid credentials")
		}
		return nil, err
	}
	if err := CheckPassword(admin.PasswordHash, input.Password); err != nil {
		return nil, apperrors.Unauthorized("invalid credentials")
	}
	access, refresh, err := s.issueTokens(ctx, admin.ID)
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		Admin:        PlatformAdminFromRow(admin),
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *Service) ListTenants(ctx context.Context, limit, offset int32) ([]TenantAdminDTO, int64, error) {
	tenants, total, err := s.repo.ListTenants(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out := make([]TenantAdminDTO, 0, len(tenants))
	for _, t := range tenants {
		out = append(out, TenantAdminFromRow(t))
	}
	return out, total, nil
}

func (s *Service) GetTenant(ctx context.Context, id uuid.UUID) (*TenantAdminDTO, error) {
	tenant, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound("Tenant not found")
		}
		return nil, err
	}
	dto := TenantAdminFromRow(tenant)
	return &dto, nil
}

func (s *Service) ProvisionTenant(ctx context.Context, input ProvisionInput) (*ProvisionResult, error) {
	if err := ValidateProvisionInput(input.BusinessName, input.OwnerEmail, input.OwnerFullName, input.OwnerPassword); err != nil {
		return nil, err
	}
	passwordHash, err := HashPassword(input.OwnerPassword)
	if err != nil {
		return nil, apperrors.Internal("failed to hash password", err)
	}
	slug := Slugify(input.BusinessName)
	tenant, owner, err := s.repo.CreateTenantWithOwner(ctx, input.BusinessName, slug, input.OwnerEmail, passwordHash, input.OwnerFullName)
	if err != nil {
		if IsUniqueViolation(err) {
			return nil, apperrors.Conflict("business name or owner email already exists")
		}
		return nil, apperrors.Internal("failed to provision tenant", err)
	}
	return &ProvisionResult{
		Tenant: TenantAdminFromRow(tenant),
		Owner:  OwnerFromUser(owner),
	}, nil
}

func (s *Service) UpdateTenantStatus(ctx context.Context, id uuid.UUID, status string) (*TenantAdminDTO, error) {
	if err := ValidateUpdateStatus(status); err != nil {
		return nil, err
	}
	tenant, err := s.repo.UpdateTenantStatus(ctx, id, status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.NotFound("Tenant not found")
		}
		return nil, err
	}
	dto := TenantAdminFromRow(tenant)
	return &dto, nil
}

func (s *Service) SeedPlatformAdmin(ctx context.Context, email, password, fullName string) error {
	if err := ValidateSeedInput(email, password); err != nil {
		return err
	}
	exists, err := s.repo.PlatformAdminExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	hash, err := HashPassword(password)
	if err != nil {
		return apperrors.Internal("failed to hash password", err)
	}
	_, err = s.repo.CreatePlatformAdmin(ctx, email, hash, fullName)
	if err != nil {
		if IsUniqueViolation(err) {
			return nil
		}
		return apperrors.Internal("failed to seed platform admin", err)
	}
	return nil
}

func (s *Service) issueTokens(ctx context.Context, adminID uuid.UUID) (access, refresh string, err error) {
	access, err = s.tokens.IssueAdminToken(adminID)
	if err != nil {
		return "", "", apperrors.Internal("failed to issue access token", err)
	}
	raw, hash, err := GenerateRefreshToken()
	if err != nil {
		return "", "", apperrors.Internal("failed to generate refresh token", err)
	}
	expiresAt := time.Now().UTC().Add(s.refreshTTL)
	if err := s.repo.CreatePlatformAdminRefreshToken(ctx, adminID, hash, expiresAt); err != nil {
		return "", "", apperrors.Internal("failed to store refresh token", err)
	}
	return access, raw, nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
