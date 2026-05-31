import { useQuery } from "@tanstack/react-query";

import { fetchAnalyticsSummary } from "./analyticsApi";
import type { AnalyticsSummaryParams } from "../model/schema";

export const analyticsKeys = {
  all: ["analytics"] as const,
  summary: (params: AnalyticsSummaryParams) =>
    [...analyticsKeys.all, "summary", params] as const,
};

export function useAnalyticsSummary(params: AnalyticsSummaryParams = {}) {
  return useQuery({
    queryKey: analyticsKeys.summary(params),
    queryFn: () => fetchAnalyticsSummary(params),
  });
}
