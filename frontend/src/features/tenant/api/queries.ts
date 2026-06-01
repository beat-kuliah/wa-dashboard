import {
  useMutation,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";

import {
  addMember,
  fetchCurrentTenant,
  listMembers,
  updateTenant,
} from "./tenantApi";
import type { MembersListParams } from "../model/schema";

export const tenantKeys = {
  all: ["tenant"] as const,
  current: () => [...tenantKeys.all, "current"] as const,
  members: () => [...tenantKeys.all, "members"] as const,
  membersList: (params: MembersListParams) =>
    [...tenantKeys.members(), params] as const,
};

export function useCurrentTenant() {
  return useQuery({
    queryKey: tenantKeys.current(),
    queryFn: fetchCurrentTenant,
  });
}

export function useUpdateTenant() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (input: { name: string }) => updateTenant(input),
    onSuccess: (tenant) => {
      queryClient.setQueryData(tenantKeys.current(), tenant);
    },
  });
}

export function useMembers(params: MembersListParams = {}) {
  return useQuery({
    queryKey: tenantKeys.membersList(params),
    queryFn: () => listMembers(params),
  });
}

export function useAddMember() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: addMember,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: tenantKeys.members() });
    },
  });
}
