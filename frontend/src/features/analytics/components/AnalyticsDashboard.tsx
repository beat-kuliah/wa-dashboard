"use client";

import { BarChart3 } from "lucide-react";

import { formatDate, formatPercent } from "@/shared/lib/format";
import { Card, CardContent, CardHeader, CardTitle } from "@/shared/ui/Card";
import { EmptyState } from "@/shared/ui/EmptyState";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Skeleton } from "@/shared/ui/Skeleton";

import { useAnalyticsSummary } from "../api/queries";

function MetricCard({
  label,
  value,
  isLoading,
}: {
  label: string;
  value: string;
  isLoading?: boolean;
}) {
  if (isLoading) {
    return <Skeleton className="h-28 rounded-lg" />;
  }

  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-ink-muted">{label}</CardTitle>
      </CardHeader>
      <CardContent>
        <p className="font-display text-2xl font-semibold tabular-nums tracking-tight">{value}</p>
      </CardContent>
    </Card>
  );
}

export function AnalyticsDashboard() {
  const { data, isLoading, isError, error, refetch } = useAnalyticsSummary();

  if (isError) {
    return (
      <ErrorState
        message={error instanceof Error ? error.message : "Failed to load analytics"}
        onRetry={() => void refetch()}
      />
    );
  }

  const hasData =
    data &&
    (data.conversation_volume > 0 ||
      data.delivery_rate > 0 ||
      data.open_rate > 0 ||
      data.resolution_rate > 0);

  if (!isLoading && data && !hasData) {
    return (
      <EmptyState
        title="No analytics data yet"
        description={`Metrics for ${formatDate(data.period_start)} – ${formatDate(data.period_end)} will populate once you start sending broadcasts and handling conversations.`}
        icon={<BarChart3 className="h-5 w-5" />}
      />
    );
  }

  return (
    <div className="space-y-6">
      {data ? (
        <p className="text-sm text-ink-muted">
          {formatDate(data.period_start)} – {formatDate(data.period_end)}
        </p>
      ) : null}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <MetricCard
          label="Delivery rate"
          value={data ? formatPercent(data.delivery_rate) : "—"}
          isLoading={isLoading}
        />
        <MetricCard
          label="Open rate"
          value={data ? formatPercent(data.open_rate) : "—"}
          isLoading={isLoading}
        />
        <MetricCard
          label="Resolution rate"
          value={data ? formatPercent(data.resolution_rate) : "—"}
          isLoading={isLoading}
        />
        <MetricCard
          label="Conversations"
          value={data ? String(data.conversation_volume) : "—"}
          isLoading={isLoading}
        />
        <MetricCard
          label="Avg response time"
          value={
            data
              ? data.avg_response_time_seconds > 0
                ? `${Math.round(data.avg_response_time_seconds)}s`
                : "—"
              : "—"
          }
          isLoading={isLoading}
        />
      </div>
    </div>
  );
}
