import { tenantResponseSchema } from "@/entities/tenant";
import { apiClient } from "@/shared/lib/apiClient";

export async function fetchCurrentTenant() {
  const json = await apiClient<unknown>("/tenant");
  return tenantResponseSchema.parse(json).tenant;
}
