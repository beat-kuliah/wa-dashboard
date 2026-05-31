import { z } from "zod";

export const broadcastStatusSchema = z.enum([
  "draft",
  "scheduled",
  "sending",
  "sent",
  "failed",
]);

export type BroadcastStatus = z.infer<typeof broadcastStatusSchema>;

export const broadcastSchema = z.object({
  id: z.string().uuid(),
  tenant_id: z.string().uuid(),
  name: z.string(),
  template_id: z.string().uuid(),
  status: broadcastStatusSchema,
  scheduled_at: z.string().nullable(),
  sent_at: z.string().nullable(),
  recipient_count: z.number(),
  delivered_count: z.number(),
  read_count: z.number(),
  failed_count: z.number(),
  error_message: z.string().nullable(),
  created_by: z.string().uuid(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Broadcast = z.infer<typeof broadcastSchema>;

export const broadcastListResponseSchema = z.object({
  data: z.array(broadcastSchema),
  page: z.object({
    limit: z.number(),
    offset: z.number(),
    total: z.number(),
  }),
});

export type BroadcastListParams = {
  limit?: number;
  offset?: number;
  status?: BroadcastStatus;
};
