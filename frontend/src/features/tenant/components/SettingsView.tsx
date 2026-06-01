"use client";

import { canViewTenantMembers, isTenantAdmin, useAuthStore } from "@/features/auth";
import { Alert, AlertDescription } from "@/shared/ui/Alert";

import { InviteMemberForm } from "./InviteMemberForm";
import { MembersList } from "./MembersList";
import { TenantSettingsForm } from "./TenantSettingsForm";

export function SettingsView() {
  const user = useAuthStore((s) => s.user);
  const isAdmin = isTenantAdmin(user);
  const showMembers = canViewTenantMembers(user);

  return (
    <div className="space-y-8">
      {!isAdmin ? (
        <Alert>
          <AlertDescription>
            You have read-only access. Only admins can edit the workspace or invite members.
          </AlertDescription>
        </Alert>
      ) : null}

      <section className="space-y-4">
        <h2 className="font-display text-lg font-semibold tracking-tight">Workspace</h2>
        <TenantSettingsForm canEdit={isAdmin} />
      </section>

      {showMembers ? (
        <section className="space-y-4">
          <h2 className="font-display text-lg font-semibold tracking-tight">Team</h2>
          {isAdmin ? <InviteMemberForm /> : null}
          <MembersList />
        </section>
      ) : (
        <p className="text-sm text-ink-muted">
          Member management is available to admins and supervisors.
        </p>
      )}
    </div>
  );
}
