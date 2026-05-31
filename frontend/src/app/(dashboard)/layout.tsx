"use client";

import { AuthGuard } from "@/features/auth";
import { DashboardShell } from "@/features/shell";

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <AuthGuard>
      <DashboardShell>{children}</DashboardShell>
    </AuthGuard>
  );
}
