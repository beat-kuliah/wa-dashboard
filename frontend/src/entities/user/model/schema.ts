import { z } from "zod";

export const userRoleSchema = z.enum(["admin", "supervisor", "agent"]);

export const userSchema = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  full_name: z.string(),
  roles: z.array(userRoleSchema),
  tenant_id: z.string().uuid(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type UserRole = z.infer<typeof userRoleSchema>;
export type User = z.infer<typeof userSchema>;

export const userResponseSchema = z.object({
  user: userSchema,
});
