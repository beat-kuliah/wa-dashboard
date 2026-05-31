"use client";

import { ErrorState } from "@/shared/ui/ErrorState";

export default function BroadcastError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <div className="mx-auto max-w-4xl space-y-6">
      <div className="space-y-1">
        <h1 className="font-display text-2xl font-semibold tracking-tight">Broadcasts</h1>
      </div>
      <ErrorState message={error.message} onRetry={reset} />
    </div>
  );
}
