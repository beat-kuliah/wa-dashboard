package broadcastmod

import (
	"context"
	"time"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/google/uuid"
)

type BroadcastDTO struct {
	ID             uuid.UUID  `json:"id"`
	TenantID       uuid.UUID  `json:"tenant_id"`
	Name           string     `json:"name"`
	TemplateID     uuid.UUID  `json:"template_id"`
	Status         string     `json:"status"`
	ScheduledAt    *time.Time `json:"scheduled_at"`
	SentAt         *time.Time `json:"sent_at"`
	RecipientCount int32      `json:"recipient_count"`
	DeliveredCount int32      `json:"delivered_count"`
	ReadCount      int32      `json:"read_count"`
	FailedCount    int32      `json:"failed_count"`
	ErrorMessage   *string    `json:"error_message"`
	CreatedBy      uuid.UUID  `json:"created_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func BroadcastFromRow(b sqlc.Broadcast) BroadcastDTO {
	var scheduledAt, sentAt *time.Time
	if b.ScheduledAt != nil {
		t := b.ScheduledAt.UTC()
		scheduledAt = &t
	}
	if b.SentAt != nil {
		t := b.SentAt.UTC()
		sentAt = &t
	}
	return BroadcastDTO{
		ID:             b.ID,
		TenantID:       b.TenantID,
		Name:           b.Name,
		TemplateID:     b.TemplateID,
		Status:         b.Status,
		ScheduledAt:    scheduledAt,
		SentAt:         sentAt,
		RecipientCount: b.RecipientCount,
		DeliveredCount: b.DeliveredCount,
		ReadCount:      b.ReadCount,
		FailedCount:    b.FailedCount,
		ErrorMessage:   b.ErrorMessage,
		CreatedBy:      b.CreatedBy,
		CreatedAt:      b.CreatedAt.UTC(),
		UpdatedAt:      b.UpdatedAt.UTC(),
	}
}

func tenantContext(ctx context.Context) (uuid.UUID, uuid.UUID, error) {
	rc, ok := auth.RequestContextFrom(ctx)
	if !ok {
		return uuid.Nil, uuid.Nil, apperrors.Unauthorized("")
	}
	return rc.TenantID, rc.UserID, nil
}

type CreateBroadcastInput struct {
	Name         string
	TemplateID   uuid.UUID
	ScheduledAt  *time.Time
}

func ValidateCreateBroadcast(name string, templateID uuid.UUID) error {
	if name == "" {
		return apperrors.Validation("name is required")
	}
	if templateID == uuid.Nil {
		return apperrors.Validation("template_id is required")
	}
	return nil
}
