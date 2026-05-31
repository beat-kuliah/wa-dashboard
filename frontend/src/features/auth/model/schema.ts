import { z } from "zod";

export const loginFormSchema = z.object({
  email: z.string().email("Enter a valid email address"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

export type LoginFormValues = z.infer<typeof loginFormSchema>;

export const tokensSchema = z.object({
  access_token: z.string(),
  refresh_token: z.string(),
  expires_in: z.number(),
});

export const loginResponseSchema = z.object({
  user: z.object({
    id: z.string().uuid(),
    email: z.string().email(),
    full_name: z.string(),
    roles: z.array(z.enum(["admin", "supervisor", "agent"])),
    tenant_id: z.string().uuid(),
    created_at: z.string(),
    updated_at: z.string(),
  }),
  tokens: tokensSchema,
});

export type LoginResponse = z.infer<typeof loginResponseSchema>;
