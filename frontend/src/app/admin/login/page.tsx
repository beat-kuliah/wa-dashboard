import type { Metadata } from "next";

import { AdminLoginForm } from "@/features/platform-admin";

export const metadata: Metadata = {
  title: "Platform sign in · WA Dashboard",
};

export default function AdminLoginPage() {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-surface px-6 py-12">
      <div className="w-full max-w-sm space-y-8">
        <div className="space-y-1">
          <h1 className="font-display text-2xl font-semibold tracking-tight">Platform console</h1>
          <p className="text-sm text-ink-muted">
            Sign in as a platform operator to manage tenants.
          </p>
        </div>
        <AdminLoginForm />
      </div>
    </div>
  );
}
