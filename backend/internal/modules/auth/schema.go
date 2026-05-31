package authmod

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"
	"unicode"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Roles     []string  `json:"roles"`
	TenantID  uuid.UUID `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TenantDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func UserFromRow(u sqlc.User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Roles:     u.Roles,
		TenantID:  u.TenantID,
		CreatedAt: u.CreatedAt.UTC(),
		UpdatedAt: u.UpdatedAt.UTC(),
	}
}

func TenantFromRow(t sqlc.Tenant) TenantDTO {
	return TenantDTO{
		ID:        t.ID,
		Name:      t.Name,
		Slug:      t.Slug,
		CreatedAt: t.CreatedAt.UTC(),
		UpdatedAt: t.UpdatedAt.UTC(),
	}
}

func Slugify(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && b.Len() > 0 {
			b.WriteRune('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "workspace"
	}
	return out
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateRefreshToken() (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	raw := base64.RawURLEncoding.EncodeToString(buf)
	sum := sha256.Sum256([]byte(raw))
	return raw, hex.EncodeToString(sum[:]), nil
}

func ValidateRegisterInput(email, password, fullName, businessName string) error {
	if strings.TrimSpace(email) == "" {
		return apperrors.Validation("email is required")
	}
	if len(password) < 8 {
		return apperrors.Validation("password must be at least 8 characters")
	}
	if strings.TrimSpace(fullName) == "" {
		return apperrors.Validation("full_name is required")
	}
	if strings.TrimSpace(businessName) == "" {
		return apperrors.Validation("business_name is required")
	}
	return nil
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

func ValidateRefreshInput(refreshToken string) error {
	if strings.TrimSpace(refreshToken) == "" {
		return apperrors.Validation("refresh_token is required")
	}
	return nil
}

func ValidateAddMemberInput(email, fullName string, roles []string) error {
	if strings.TrimSpace(email) == "" {
		return apperrors.Validation("email is required")
	}
	if strings.TrimSpace(fullName) == "" {
		return apperrors.Validation("full_name is required")
	}
	if len(roles) == 0 {
		return apperrors.Validation("roles is required")
	}
	for _, role := range roles {
		if !auth.HasRole([]string{role}, auth.RoleAdmin, auth.RoleSupervisor, auth.RoleAgent) {
			return apperrors.Validation("invalid role: " + role)
		}
	}
	return nil
}
