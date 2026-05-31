import { useQuery } from "@tanstack/react-query";

import { fetchBroadcasts } from "./broadcastApi";
import type { BroadcastListParams } from "../model/schema";

export const broadcastKeys = {
  all: ["broadcasts"] as const,
  lists: () => [...broadcastKeys.all, "list"] as const,
  list: (params: BroadcastListParams) => [...broadcastKeys.lists(), params] as const,
};

export function useBroadcasts(params: BroadcastListParams = {}) {
  return useQuery({
    queryKey: broadcastKeys.list(params),
    queryFn: () => fetchBroadcasts(params),
  });
}
