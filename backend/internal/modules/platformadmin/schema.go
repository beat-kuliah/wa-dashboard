package platformadmin

import (
	"strings"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	authmod "github.com/beatfraps/wa-dashboard/backend/internal/modules/auth"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/google/uuid"
)

type PlatformAdminDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
}

type TenantAdminDTO struct {
	ID           uuid.UUID      `json:"id"`
	BusinessName string         `json:"business_name"`
	Status       string         `json:"status"`
	Settings     map[string]any `json:"settings"`
	AIEnabled    bool           `json:"ai_enabled"`
	Features     TenantFeatures `json:"features"`
	CreatedAt    time.Time      `json:"created_at"`
}

type TenantFeatures struct {
	Broadcast  bool `json:"broadcast"`
	CSInbox    bool `json:"cs_inbox"`
	Analytics  bool `json:"analytics"`
	AIChatbot  bool `json:"ai_chatbot"`
}

type OwnerDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	TenantID  uuid.UUID `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
}

func PlatformAdminFromRow(a sqlc.PlatformAdmin) PlatformAdminDTO {
	return PlatformAdminDTO{
		ID:        a.ID,
		Email:     a.Email,
		FullName:  a.FullName,
		CreatedAt: a.CreatedAt.UTC(),
	}
}

func TenantAdminFromRow(t sqlc.Tenant) TenantAdminDTO {
	return TenantAdminDTO{
		ID:           t.ID,
		BusinessName: t.Name,
		Status:       t.Status,
		Settings:     map[string]any{},
		AIEnabled:    false,
		Features: TenantFeatures{
			Broadcast: true,
			CSInbox:   true,
			Analytics: true,
			AIChatbot: false,
		},
		CreatedAt: t.CreatedAt.UTC(),
	}
}

func OwnerFromUser(u sqlc.User) OwnerDTO {
	role := "admin"
	if len(u.Roles) > 0 {
		role = u.Roles[0]
	}
	return OwnerDTO{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      role,
		TenantID:  u.TenantID,
		CreatedAt: u.CreatedAt.UTC(),
	}
}

func ValidateLoginInput(email, password string) error {
	if strings.TrimSpace(email) == "" {
		return apperrors.Validation("email is required")
	}
	if password == "" {
		return apperrors.Validation("password is required")
	}
	return nil
}

func ValidateProvisionInput(businessName, ownerEmail, ownerFullName, ownerPassword string) error {
	if strings.TrimSpace(businessName) == "" {
		return apperrors.Validation("business_name is required")
	}
	if strings.TrimSpace(ownerEmail) == "" {
		return apperrors.Validation("owner_email is required")
	}
	if strings.TrimSpace(ownerFullName) == "" {
		return apperrors.Validation("owner_full_name is required")
	}
	if len(ownerPassword) < 8 {
		return apperrors.Validation("owner_password must be at least 8 characters")
	}
	return nil
}

func ValidateUpdateStatus(status string) error {
	switch status {
	case "active", "suspended":
		return nil
	default:
		return apperrors.Validation("status must be active or suspended")
	}
}

func ValidateSeedInput(email, password string) error {
	if strings.TrimSpace(email) == "" || password == "" {
		return apperrors.Validation("platform admin seed email and password are required")
	}
	if len(password) < 8 {
		return apperrors.Validation("platform admin seed password must be at least 8 characters")
	}
	return nil
}

// Re-export password helpers from auth module.
var (
	HashPassword   = authmod.HashPassword
	CheckPassword  = authmod.CheckPassword
	Slugify        = authmod.Slugify
	GenerateRefreshToken = authmod.GenerateRefreshToken
)
