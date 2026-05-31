import type { Metadata } from "next";

import { LoginForm } from "@/features/auth";

export const metadata: Metadata = {
  title: "Sign in · WA Dashboard",
};

export default function LoginPage() {
  return (
    <div className="flex min-h-screen">
      <div className="hidden w-[42%] flex-col justify-between bg-surface-sunken p-10 lg:flex">
        <div>
          <p className="font-display text-lg font-bold tracking-tight">
            WA<span className="text-accent">.</span>Dashboard
          </p>
        </div>
        <div className="max-w-sm space-y-4">
          <h1 className="font-display text-3xl font-semibold leading-tight tracking-tight">
            WhatsApp operations, one workspace.
          </h1>
          <p className="text-sm leading-relaxed text-ink-muted">
            Manage broadcasts, customer conversations, and delivery analytics for your
            business — scoped to your tenant, secured by role.
          </p>
        </div>
        <p className="text-xs text-ink-subtle">Multi-tenant · JWT auth · Realtime inbox</p>
      </div>

      <div className="flex flex-1 flex-col items-center justify-center px-6 py-12">
        <div className="w-full max-w-sm space-y-8">
          <div className="space-y-2 lg:hidden">
            <p className="font-display text-lg font-bold tracking-tight">
              WA<span className="text-accent">.</span>Dashboard
            </p>
          </div>
          <div className="space-y-1">
            <h2 className="font-display text-2xl font-semibold tracking-tight">Sign in</h2>
            <p className="text-sm text-ink-muted">Enter your credentials to continue.</p>
          </div>
          <LoginForm />
        </div>
      </div>
    </div>
  );
}
