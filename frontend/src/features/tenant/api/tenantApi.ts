import { userResponseSchema, type UserRole } from "@/entities/user";
import { tenantResponseSchema } from "@/entities/tenant";
import { apiClient, buildQueryString } from "@/shared/lib/apiClient";

import {
  membersListResponseSchema,
  type MembersListParams,
} from "../model/schema";

export async function fetchCurrentTenant() {
  const json = await apiClient<unknown>("/tenant");
  return tenantResponseSchema.parse(json).tenant;
}

export async function updateTenant(input: { name: string }) {
  const json = await apiClient<unknown>("/tenant", {
    method: "PATCH",
    body: JSON.stringify({ name: input.name }),
  });
  return tenantResponseSchema.parse(json).tenant;
}

export async function listMembers(params: MembersListParams = {}) {
  const qs = buildQueryString({
    limit: params.limit ?? 20,
    offset: params.offset ?? 0,
  });

  const json = await apiClient<unknown>(`/tenant/members${qs}`);
  return membersListResponseSchema.parse(json);
}

export async function addMember(input: {
  email: string;
  full_name: string;
  role: UserRole;
}) {
  const json = await apiClient<unknown>("/tenant/members", {
    method: "POST",
    body: JSON.stringify({
      email: input.email,
      full_name: input.full_name,
      roles: [input.role],
    }),
  });
  return userResponseSchema.parse(json).user;
}
