"use client";

import { ErrorState } from "@/shared/ui/ErrorState";

export default function AnalyticsError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="mx-auto max-w-5xl space-y-6">
      <div className="space-y-1">
        <h1 className="font-display text-2xl font-semibold tracking-tight">Analytics</h1>
      </div>
      <ErrorState message={error.message} onRetry={reset} />
    </div>
  );
}
