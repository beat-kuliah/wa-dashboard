"use client";

import { ErrorState } from "@/shared/ui/ErrorState";

export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <html lang="en">
      <body className="flex min-h-screen items-center justify-center bg-surface p-6">
        <ErrorState
          title="Application error"
          message={error.message}
          onRetry={reset}
        />
      </body>
    </html>
  );
}
