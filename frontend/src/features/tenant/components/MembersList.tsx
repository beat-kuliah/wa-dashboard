"use client";

import { Users } from "lucide-react";

import type { UserRole } from "@/entities/user";
import { formatDate } from "@/shared/lib/format";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Skeleton } from "@/shared/ui/Skeleton";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/Card";

import { useMembers } from "../api/queries";
import { MemberRoleBadge } from "./MemberRoleBadge";

export function MembersList() {
  const { data, isLoading, isError, error, refetch } = useMembers();

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Members</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          {Array.from({ length: 4 }).map((_, i) => (
            <Skeleton key={i} className="h-12 w-full rounded-md" />
          ))}
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <ErrorState
        message={error instanceof Error ? error.message : "Failed to load members"}
        onRetry={() => void refetch()}
      />
    );
  }

  if (!data?.data.length) {
    return (
      <EmptyState
        title="No members yet"
        description="Invite teammates to collaborate in this workspace."
        icon={<Users className="h-5 w-5" />}
      />
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Members</CardTitle>
        <CardDescription>
          People with access to this workspace.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border-subtle text-left text-xs uppercase tracking-wide text-ink-subtle">
                <th className="px-3 py-2 font-medium">Name</th>
                <th className="px-3 py-2 font-medium">Email</th>
                <th className="px-3 py-2 font-medium">Role</th>
                <th className="px-3 py-2 font-medium">Joined</th>
              </tr>
            </thead>
            <tbody>
              {data.data.map((member) => (
                <tr
                  key={member.id}
                  className="border-b border-border-subtle/60 last:border-0"
                >
                  <td className="px-3 py-3 font-medium text-ink">{member.full_name}</td>
                  <td className="px-3 py-3 text-ink-muted">{member.email}</td>
                  <td className="px-3 py-3">
                    <div className="flex flex-wrap gap-1.5">
                      {member.roles.map((role) => (
                        <MemberRoleBadge key={role} role={role as UserRole} />
                      ))}
                    </div>
                  </td>
                  <td className="px-3 py-3 text-ink-muted">
                    {formatDate(member.created_at)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
        {data.page.total > data.data.length ? (
          <p className="mt-4 text-center text-sm text-ink-muted">
            Showing {data.data.length} of {data.page.total} members
          </p>
        ) : null}
      </CardContent>
    </Card>
  );
}
