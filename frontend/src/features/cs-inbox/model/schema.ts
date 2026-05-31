import { z } from "zod";

export const conversationStatusSchema = z.enum(["open", "in_progress", "resolved"]);

export const conversationListResponseSchema = z.object({
  data: z.array(z.unknown()),
  page: z.object({
    limit: z.number(),
    offset: z.number(),
    total: z.number(),
  }),
});

export type ConversationListParams = {
  limit?: number;
  offset?: number;
  status?: z.infer<typeof conversationStatusSchema>;
};
