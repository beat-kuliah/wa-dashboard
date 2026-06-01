import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { usePlatformAdminStore } from "../model/store";
import type { AdminLoginFormValues, ProvisionTenantFormValues } from "../model/schema";
import {
  adminLogin,
  adminLogout,
  fetchAdminTenants,
  provisionTenant,
  updateTenantStatus,
  type AdminTenantListParams,
} from "./platformAdminApi";

export const platformAdminKeys = {
  all: ["platform-admin"] as const,
  tenants: (params?: AdminTenantListParams) =>
    [...platformAdminKeys.all, "tenants", params] as const,
};

export function useAdminTenants(params?: AdminTenantListParams) {
  return useQuery({
    queryKey: platformAdminKeys.tenants(params),
    queryFn: () => fetchAdminTenants(params),
  });
}

export function useAdminLogin() {
  const setAdmin = usePlatformAdminStore((s) => s.setAdmin);

  return useMutation({
    mutationFn: (values: AdminLoginFormValues) => adminLogin(values),
    onSuccess: (data) => {
      setAdmin(data.admin);
    },
  });
}

export function useAdminLogout() {
  const queryClient = useQueryClient();
  const clearAdmin = usePlatformAdminStore((s) => s.clearAdmin);

  return useMutation({
    mutationFn: () => adminLogout(),
    onSettled: () => {
      clearAdmin();
      queryClient.removeQueries({ queryKey: platformAdminKeys.all });
    },
  });
}

export function useProvisionTenant() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (values: ProvisionTenantFormValues) => provisionTenant(values),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: platformAdminKeys.all });
    },
  });
}

export function useUpdateTenantStatus() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      tenantId,
      status,
    }: {
      tenantId: string;
      status: "active" | "suspended";
    }) => updateTenantStatus(tenantId, { status }),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: platformAdminKeys.all });
    },
  });
}
