import { adminApiClient } from "@/shared/lib/adminApiClient";
import { buildQueryString } from "@/shared/lib/apiClient";
import { clearAdminTokens, setAdminTokens } from "@/shared/lib/adminAuthStorage";

import {
  adminLoginResponseSchema,
  adminTenantListResponseSchema,
  adminTenantSchema,
  provisionTenantResponseSchema,
  type AdminLoginFormValues,
  type ProvisionTenantFormValues,
  type UpdateTenantStatusValues,
} from "../model/schema";

export type AdminTenantListParams = {
  limit?: number;
  offset?: number;
};

export async function adminLogin(values: AdminLoginFormValues) {
  const json = await adminApiClient<unknown>("/admin/auth/login", {
    method: "POST",
    body: JSON.stringify(values),
    skipAuth: true,
  });

  const parsed = adminLoginResponseSchema.parse(json);
  setAdminTokens(parsed.access_token, parsed.refresh_token);
  return parsed;
}

export async function adminLogout() {
  clearAdminTokens();
}

export async function fetchAdminTenants(params: AdminTenantListParams = {}) {
  const qs = buildQueryString({
    limit: params.limit ?? 50,
    offset: params.offset ?? 0,
  });
  const json = await adminApiClient<unknown>(`/admin/tenants${qs}`);
  return adminTenantListResponseSchema.parse(json);
}

export async function provisionTenant(values: ProvisionTenantFormValues) {
  const json = await adminApiClient<unknown>("/admin/tenants", {
    method: "POST",
    body: JSON.stringify(values),
  });
  return provisionTenantResponseSchema.parse(json);
}

export async function updateTenantStatus(tenantId: string, values: UpdateTenantStatusValues) {
  const json = await adminApiClient<unknown>(`/admin/tenants/${tenantId}`, {
    method: "PATCH",
    body: JSON.stringify(values),
  });
  return adminTenantSchema.parse(json);
}
