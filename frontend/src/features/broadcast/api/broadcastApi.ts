import {
  apiClient,
  buildQueryString,
  paginationSchema,
} from "@/shared/lib/apiClient";

import {
  broadcastListResponseSchema,
  type BroadcastListParams,
} from "../model/schema";

export async function fetchBroadcasts(params: BroadcastListParams = {}) {
  const qs = buildQueryString({
    limit: params.limit ?? 20,
    offset: params.offset ?? 0,
    status: params.status,
  });

  const json = await apiClient<unknown>(`/broadcasts${qs}`);
  return broadcastListResponseSchema.parse(json);
}

export { paginationSchema };
