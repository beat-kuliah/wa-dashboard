import type { Metadata } from "next";

import { BroadcastList } from "@/features/broadcast";

export const metadata: Metadata = {
  title: "Broadcasts · WA Dashboard",
};

export default function BroadcastPage() {
  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <div className="space-y-1">
        <h1 className="font-display text-2xl font-semibold tracking-tight">Broadcasts</h1>
        <p className="text-sm text-ink-muted">
          WhatsApp broadcast campaigns for your workspace.
        </p>
      </div>
      <BroadcastList />
    </div>
  );
}
