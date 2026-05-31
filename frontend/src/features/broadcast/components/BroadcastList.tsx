"use client";

import { Megaphone } from "lucide-react";

import { formatDate } from "@/shared/lib/format";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Skeleton } from "@/shared/ui/Skeleton";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/Card";

import { useBroadcasts } from "../api/queries";
import { BroadcastStatusBadge } from "./BroadcastStatusBadge";

export function BroadcastList() {
  const { data, isLoading, isError, error, refetch } = useBroadcasts();

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-24 w-full rounded-lg" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <ErrorState
        message={error instanceof Error ? error.message : "Failed to load broadcasts"}
        onRetry={() => void refetch()}
      />
    );
  }

  if (!data?.data.length) {
    return (
      <EmptyState
        title="No broadcasts yet"
        description="Create your first WhatsApp broadcast campaign to reach your audience."
        icon={<Megaphone className="h-5 w-5" />}
      />
    );
  }

  return (
    <div className="space-y-3">
      {data.data.map((broadcast) => (
        <Card key={broadcast.id} className="transition-shadow hover:shadow-md">
          <CardHeader className="flex-row items-start justify-between gap-4 space-y-0 pb-2">
            <div className="min-w-0 flex-1">
              <CardTitle className="truncate text-base">{broadcast.name}</CardTitle>
              <CardDescription className="mt-1">
                Created {formatDate(broadcast.created_at)}
                {broadcast.scheduled_at
                  ? ` · Scheduled ${formatDate(broadcast.scheduled_at)}`
                  : null}
              </CardDescription>
            </div>
            <BroadcastStatusBadge status={broadcast.status} />
          </CardHeader>
          <CardContent>
            <dl className="grid grid-cols-2 gap-4 text-sm sm:grid-cols-4">
              <div>
                <dt className="text-ink-subtle">Recipients</dt>
                <dd className="font-medium tabular-nums">{broadcast.recipient_count}</dd>
              </div>
              <div>
                <dt className="text-ink-subtle">Delivered</dt>
                <dd className="font-medium tabular-nums">{broadcast.delivered_count}</dd>
              </div>
              <div>
                <dt className="text-ink-subtle">Read</dt>
                <dd className="font-medium tabular-nums">{broadcast.read_count}</dd>
              </div>
              <div>
                <dt className="text-ink-subtle">Failed</dt>
                <dd className="font-medium tabular-nums">{broadcast.failed_count}</dd>
              </div>
            </dl>
            {broadcast.error_message ? (
              <p className="mt-3 text-sm text-destructive">{broadcast.error_message}</p>
            ) : null}
          </CardContent>
        </Card>
      ))}
      {data.page.total > data.data.length ? (
        <p className="text-center text-sm text-ink-muted">
          Showing {data.data.length} of {data.page.total} broadcasts
        </p>
      ) : null}
    </div>
  );
}
