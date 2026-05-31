package broadcastmod

import (
	"context"
	"encoding/json"

	"github.com/beatfraps/wa-dashboard/backend/db/sqlc"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/auth"
	apperrors "github.com/beatfraps/wa-dashboard/backend/internal/shared/errors"
	"github.com/beatfraps/wa-dashboard/backend/internal/shared/queue"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type Service struct {
	repo   *Repository
	queue  *queue.Client
}

func NewService(repo *Repository, q *queue.Client) *Service {
	return &Service{repo: repo, queue: q}
}

func (s *Service) List(ctx context.Context, limit, offset int32, status *string) ([]BroadcastDTO, int64, error) {
	tenantID, _, err := tenantContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	items, total, err := s.repo.List(ctx, tenantID, limit, offset, status)
	if err != nil {
		return nil, 0, err
	}
	out := make([]BroadcastDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BroadcastFromRow(item))
	}
	return out, total, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*BroadcastDTO, error) {
	tenantID, _, err := tenantContext(ctx)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, MapNotFound(err)
	}
	dto := BroadcastFromRow(item)
	return &dto, nil
}

func (s *Service) Create(ctx context.Context, input CreateBroadcastInput) (*BroadcastDTO, error) {
	if err := auth.AssertRole(ctx, auth.RoleAdmin, auth.RoleSupervisor); err != nil {
		return nil, err
	}
	if err := ValidateCreateBroadcast(input.Name, input.TemplateID); err != nil {
		return nil, err
	}
	tenantID, userID, err := tenantContext(ctx)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.TemplateExists(ctx, tenantID, input.TemplateID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, apperrors.NotFound("Template not found")
	}

	status := "draft"
	if input.ScheduledAt != nil {
		status = "scheduled"
	}

	item, err := s.repo.Create(ctx, sqlc.CreateBroadcastParams{
		TenantID:    tenantID,
		Name:        input.Name,
		TemplateID:  input.TemplateID,
		Status:      status,
		ScheduledAt: input.ScheduledAt,
		CreatedBy:   userID,
	})
	if err != nil {
		return nil, apperrors.Internal("failed to create broadcast", err)
	}

	if status == "scheduled" {
		payload, _ := json.Marshal(SendJobPayload{
			BroadcastID: item.ID.String(),
			TenantID:    tenantID.String(),
		})
		task := asynq.NewTask(queue.TaskBroadcastSend, payload)
		if _, err := s.queue.Enqueue(ctx, task, asynq.TaskID("broadcast:"+item.ID.String())); err != nil {
			return nil, apperrors.Internal("failed to enqueue broadcast", err)
		}
	}

	dto := BroadcastFromRow(item)
	return &dto, nil
}

func (s *Service) ProcessSend(ctx context.Context, broadcastID, tenantID uuid.UUID) error {
	item, err := s.repo.GetByID(ctx, tenantID, broadcastID)
	if err != nil {
		return MapNotFound(err)
	}
	if item.Status == "sent" {
		return nil
	}
	_, err = s.repo.UpdateStatus(ctx, tenantID, broadcastID, "sent", nil)
	return err
}

type SendJobPayload struct {
	BroadcastID string `json:"broadcast_id"`
	TenantID    string `json:"tenant_id"`
}

func ParseSendJobPayload(data []byte) (SendJobPayload, error) {
	var p SendJobPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return p, err
	}
	return p, nil
}
