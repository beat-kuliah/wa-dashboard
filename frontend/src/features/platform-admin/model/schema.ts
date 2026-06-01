import { z } from "zod";

export const platformAdminSchema = z.object({
  id: z.string().uuid(),
  email: z.string().email(),
  full_name: z.string(),
  created_at: z.string(),
});

export type PlatformAdmin = z.infer<typeof platformAdminSchema>;

export const adminLoginFormSchema = z.object({
  email: z.string().email("Enter a valid email address"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

export type AdminLoginFormValues = z.infer<typeof adminLoginFormSchema>;

export const adminLoginResponseSchema = z.object({
  admin: platformAdminSchema,
  access_token: z.string(),
  refresh_token: z.string(),
});

export const tenantStatusSchema = z.enum(["active", "suspended"]);

export type TenantStatus = z.infer<typeof tenantStatusSchema>;

export const adminTenantSchema = z.object({
  id: z.string().uuid(),
  business_name: z.string(),
  status: tenantStatusSchema,
  settings: z.record(z.string(), z.unknown()).optional().default({}),
  ai_enabled: z.boolean().optional().default(false),
  features: z
    .object({
      broadcast: z.boolean(),
      cs_inbox: z.boolean(),
      analytics: z.boolean(),
      ai_chatbot: z.boolean(),
    })
    .optional(),
  created_at: z.string(),
});

export type AdminTenant = z.infer<typeof adminTenantSchema>;

export const adminTenantListResponseSchema = z.object({
  data: z.array(adminTenantSchema),
  page: z.object({
    limit: z.number(),
    offset: z.number(),
    total: z.number(),
  }),
});

export const provisionTenantFormSchema = z.object({
  business_name: z.string().min(1, "Business name is required"),
  owner_email: z.string().email("Enter a valid owner email"),
  owner_full_name: z.string().min(1, "Owner full name is required"),
  owner_password: z.string().min(8, "Password must be at least 8 characters"),
});

export type ProvisionTenantFormValues = z.infer<typeof provisionTenantFormSchema>;

export const provisionTenantResponseSchema = z.object({
  tenant: adminTenantSchema,
  owner: z.object({
    id: z.string().uuid(),
    email: z.string().email(),
    full_name: z.string(),
    role: z.enum(["admin", "supervisor", "agent"]),
    tenant_id: z.string().uuid(),
    created_at: z.string(),
  }),
});

export const updateTenantStatusSchema = z.object({
  status: tenantStatusSchema,
});

export type UpdateTenantStatusValues = z.infer<typeof updateTenantStatusSchema>;
