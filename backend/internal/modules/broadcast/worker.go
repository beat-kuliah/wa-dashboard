package broadcastmod

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/beatfraps/wa-dashboard/backend/internal/shared/queue"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

func RegisterWorkerHandlers(worker *queue.Worker, svc *Service, logger *slog.Logger) {
	worker.Handle(queue.TaskBroadcastSend, func(ctx context.Context, task *asynq.Task) error {
		var payload SendJobPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			return err
		}
		broadcastID, err := uuid.Parse(payload.BroadcastID)
		if err != nil {
			return err
		}
		tenantID, err := uuid.Parse(payload.TenantID)
		if err != nil {
			return err
		}
		logger.Info("processing broadcast send", "broadcast_id", broadcastID, "tenant_id", tenantID)
		return svc.ProcessSend(ctx, broadcastID, tenantID)
	})
}
