import { z } from "zod";

import { userRoleSchema, userSchema } from "@/entities/user";
import { paginationSchema } from "@/shared/lib/apiClient";

export const membersListResponseSchema = z.object({
  data: z.array(userSchema),
  page: paginationSchema,
});

export type MembersListResponse = z.infer<typeof membersListResponseSchema>;

export type MembersListParams = {
  limit?: number;
  offset?: number;
};

export const updateTenantFormSchema = z.object({
  name: z.string().min(1, "Workspace name is required").max(120, "Name is too long"),
});

export type UpdateTenantFormValues = z.infer<typeof updateTenantFormSchema>;

export const addMemberFormSchema = z.object({
  email: z.string().email("Enter a valid email address"),
  full_name: z.string().min(1, "Full name is required"),
  role: userRoleSchema,
});

export type AddMemberFormValues = z.infer<typeof addMemberFormSchema>;
