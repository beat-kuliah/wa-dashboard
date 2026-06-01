import type { Metadata } from "next";

import { SettingsView } from "@/features/tenant";

export const metadata: Metadata = {
  title: "Settings · WA Dashboard",
};

export default function SettingsPage() {
  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <div className="space-y-1">
        <h1 className="font-display text-2xl font-semibold tracking-tight">Settings</h1>
        <p className="text-sm text-ink-muted">
          Manage your workspace and team members.
        </p>
      </div>
      <SettingsView />
    </div>
  );
}
