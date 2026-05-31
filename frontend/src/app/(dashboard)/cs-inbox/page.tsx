import type { Metadata } from "next";

import { ConversationList } from "@/features/cs-inbox";

export const metadata: Metadata = {
  title: "CS Inbox · WA Dashboard",
};

export default function CsInboxPage() {
  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <div className="space-y-1">
        <h1 className="font-display text-2xl font-semibold tracking-tight">CS Inbox</h1>
        <p className="text-sm text-ink-muted">
          Customer conversations from WhatsApp, in one queue.
        </p>
      </div>
      <ConversationList />
    </div>
  );
}
