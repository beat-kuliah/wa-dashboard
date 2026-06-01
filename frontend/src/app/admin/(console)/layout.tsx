"use client";

import { AdminAuthGuard, AdminShell } from "@/features/platform-admin";

export default function AdminConsoleLayout({ children }: { children: React.ReactNode }) {
  return (
    <AdminAuthGuard>
      <AdminShell>{children}</AdminShell>
    </AdminAuthGuard>
  );
}
