import { useQuery } from "@tanstack/react-query";

import { fetchCurrentTenant } from "./tenantApi";

export const tenantKeys = {
  all: ["tenant"] as const,
  current: () => [...tenantKeys.all, "current"] as const,
};

export function useCurrentTenant() {
  return useQuery({
    queryKey: tenantKeys.current(),
    queryFn: fetchCurrentTenant,
  });
}
