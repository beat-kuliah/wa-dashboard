import { apiClient, buildQueryString } from "@/shared/lib/apiClient";

import {
  analyticsSummaryResponseSchema,
  type AnalyticsSummaryParams,
} from "../model/schema";

export async function fetchAnalyticsSummary(params: AnalyticsSummaryParams = {}) {
  const qs = buildQueryString({
    from: params.from,
    to: params.to,
  });

  const json = await apiClient<unknown>(`/analytics/summary${qs}`);
  return analyticsSummaryResponseSchema.parse(json).summary;
}
