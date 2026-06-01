import type { Metadata } from "next";

import { ProvisionTenantForm, TenantList } from "@/features/platform-admin";

export const metadata: Metadata = {
  title: "Tenants · Platform Console",
};

export default function AdminTenantsPage() {
  return (
    <div className="mx-auto max-w-4xl space-y-8">
      <div className="space-y-1">
        <h1 className="font-display text-2xl font-semibold tracking-tight">Tenants</h1>
        <p className="text-sm text-ink-muted">
          Provision workspaces and suspend or reactivate tenant access.
        </p>
      </div>
      <ProvisionTenantForm />
      <TenantList />
    </div>
  );
}
