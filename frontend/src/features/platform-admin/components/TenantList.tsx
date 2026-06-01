"use client";

import { formatDate } from "@/shared/lib/format";
import { Badge } from "@/shared/ui/Badge";
import { Button } from "@/shared/ui/Button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Skeleton } from "@/shared/ui/Skeleton";
import { Building2 } from "lucide-react";

import { useAdminTenants, useUpdateTenantStatus } from "../api/queries";
import type { AdminTenant } from "../model/schema";

function TenantStatusBadge({ status }: { status: AdminTenant["status"] }) {
  return (
    <Badge variant={status === "active" ? "default" : "secondary"}>
      {status}
    </Badge>
  );
}

function TenantRow({ tenant }: { tenant: AdminTenant }) {
  const updateStatus = useUpdateTenantStatus();
  const isPending = updateStatus.isPending && updateStatus.variables?.tenantId === tenant.id;
  const nextStatus = tenant.status === "active" ? "suspended" : "active";

  const handleToggle = () => {
    updateStatus.mutate({ tenantId: tenant.id, status: nextStatus });
  };

  return (
    <tr className="border-b border-border-subtle/60 last:border-0">
      <td className="py-3 pr-4 font-medium">{tenant.business_name}</td>
      <td className="py-3 pr-4">
        <TenantStatusBadge status={tenant.status} />
      </td>
      <td className="py-3 pr-4 font-mono text-xs text-ink-muted">{tenant.id.slice(0, 8)}…</td>
      <td className="py-3 pr-4 text-ink-muted">{formatDate(tenant.created_at)}</td>
      <td className="py-3 text-right">
        <Button
          variant="outline"
          size="sm"
          onClick={handleToggle}
          disabled={isPending}
        >
          {isPending
            ? "Updating…"
            : tenant.status === "active"
              ? "Suspend"
              : "Activate"}
        </Button>
      </td>
    </tr>
  );
}

export function TenantList() {
  const { data, isLoading, isError, error, refetch } = useAdminTenants();

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-16 w-full rounded-lg" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <ErrorState
        message={error instanceof Error ? error.message : "Failed to load tenants"}
        onRetry={() => void refetch()}
      />
    );
  }

  const tenants = data?.data ?? [];

  if (tenants.length === 0) {
    return (
      <EmptyState
        title="No tenants yet"
        description="Provision your first tenant workspace to get started."
        icon={<Building2 className="h-5 w-5" />}
      />
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Tenants</CardTitle>
        <CardDescription>
          {data?.page.total ?? tenants.length} workspace
          {(data?.page.total ?? tenants.length) === 1 ? "" : "s"} on the platform.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="w-full min-w-[640px] text-left text-sm">
            <thead>
              <tr className="border-b border-border-subtle text-ink-subtle">
                <th className="pb-2 pr-4 font-medium">Business</th>
                <th className="pb-2 pr-4 font-medium">Status</th>
                <th className="pb-2 pr-4 font-medium">ID</th>
                <th className="pb-2 pr-4 font-medium">Created</th>
                <th className="pb-2 text-right font-medium">Actions</th>
              </tr>
            </thead>
            <tbody>
              {tenants.map((tenant) => (
                <TenantRow key={tenant.id} tenant={tenant} />
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
}
