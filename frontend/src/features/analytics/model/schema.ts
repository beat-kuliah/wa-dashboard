import { z } from "zod";

export const analyticsSummarySchema = z.object({
  delivery_rate: z.number(),
  open_rate: z.number(),
  avg_response_time_seconds: z.number(),
  conversation_volume: z.number(),
  resolution_rate: z.number(),
  period_start: z.string(),
  period_end: z.string(),
});

export type AnalyticsSummary = z.infer<typeof analyticsSummarySchema>;

export const analyticsSummaryResponseSchema = z.object({
  summary: analyticsSummarySchema,
});

export type AnalyticsSummaryParams = {
  from?: string;
  to?: string;
};
