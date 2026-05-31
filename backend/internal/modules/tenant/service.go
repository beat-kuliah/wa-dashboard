package tenantmod

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	authmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/auth"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCurrentTenant(ctx context.Context) (*TenantDTO, error) {
	tenantID, err := tenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	tenant, err := s.repo.GetTenant(ctx, tenantID)
	if err != nil {
		return nil, MapNotFound(err)
	}
	dto := TenantFromRow(tenant)
	return &dto, nil
}

func (s *Service) UpdateTenant(ctx context.Context, input UpdateTenantInput) (*TenantDTO, error) {
	if err := auth.AssertRole(ctx, auth.RoleAdmin); err != nil {
		return nil, err
	}
	if err := ValidateUpdateTenant(input.Name); err != nil {
		return nil, err
	}
	tenantID, err := tenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	tenant, err := s.repo.UpdateTenant(ctx, tenantID, input.Name)
	if err != nil {
		return nil, MapNotFound(err)
	}
	dto := TenantFromRow(tenant)
	return &dto, nil
}

func (s *Service) ListMembers(ctx context.Context, limit, offset int32) ([]authmod.UserDTO, int64, error) {
	if err := auth.AssertRole(ctx, auth.RoleAdmin, auth.RoleSupervisor); err != nil {
		return nil, 0, err
	}
	tenantID, err := tenantIDFromContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	users, total, err := s.repo.ListMembers(ctx, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	out := make([]authmod.UserDTO, 0, len(users))
	for _, u := range users {
		out = append(out, authmod.UserFromRow(u))
	}
	return out, total, nil
}

func (s *Service) AddMember(ctx context.Context, input AddMemberInput) (*authmod.UserDTO, error) {
	if err := auth.AssertRole(ctx, auth.RoleAdmin); err != nil {
		return nil, err
	}
	if err := authmod.ValidateAddMemberInput(input.Email, input.FullName, input.Roles); err != nil {
		return nil, err
	}
	tenantID, err := tenantIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	exists, err := s.repo.MemberExists(ctx, tenantID, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperrors.Conflict("email already exists in tenant")
	}

	tempPassword, err := randomPassword()
	if err != nil {
		return nil, apperrors.Internal("failed to generate password", err)
	}
	passwordHash, err := authmod.HashPassword(tempPassword)
	if err != nil {
		return nil, apperrors.Internal("failed to hash password", err)
	}

	user, err := s.repo.CreateMember(ctx, tenantID, input.Email, passwordHash, input.FullName, input.Roles)
	if err != nil {
		if IsUniqueViolation(err) {
			return nil, apperrors.Conflict("email already exists in tenant")
		}
		return nil, apperrors.Internal("failed to add member", err)
	}
	dto := authmod.UserFromRow(user)
	return &dto, nil
}

func randomPassword() (string, error) {
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
