import type { Metadata } from "next";

import { AnalyticsDashboard } from "@/features/analytics";

export const metadata: Metadata = {
  title: "Analytics · WA Dashboard",
};

export default function AnalyticsPage() {
  return (
    <div className="mx-auto max-w-5xl space-y-6">
      <div className="space-y-1">
        <h1 className="font-display text-2xl font-semibold tracking-tight">Analytics</h1>
        <p className="text-sm text-ink-muted">
          Delivery, engagement, and support metrics for the current period.
        </p>
      </div>
      <AnalyticsDashboard />
    </div>
  );
}
