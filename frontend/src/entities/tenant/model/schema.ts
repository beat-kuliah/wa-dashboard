import { z } from "zod";

export const tenantSchema = z.object({
  id: z.string().uuid(),
  name: z.string(),
  slug: z.string(),
  created_at: z.string(),
  updated_at: z.string(),
});

export type Tenant = z.infer<typeof tenantSchema>;

export const tenantResponseSchema = z.object({
  tenant: tenantSchema,
});
