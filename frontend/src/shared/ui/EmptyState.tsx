import { Inbox } from "lucide-react";
import * as React from "react";

import { cn } from "@/shared/lib/cn";

type EmptyStateProps = {
  title: string;
  description?: string;
  icon?: React.ReactNode;
  action?: React.ReactNode;
  className?: string;
};

export function EmptyState({
  title,
  description,
  icon,
  action,
  className,
}: EmptyStateProps) {
  return (
    <div
      className={cn(
        "flex flex-col items-start gap-4 rounded-lg border border-dashed border-border-subtle bg-surface-raised/50 px-8 py-12",
        className,
      )}
    >
      <div className="flex h-10 w-10 items-center justify-center rounded-md bg-surface-sunken text-ink-muted">
        {icon ?? <Inbox className="h-5 w-5" />}
      </div>
      <div className="space-y-1">
        <h3 className="font-display text-base font-semibold text-ink">{title}</h3>
        {description ? <p className="max-w-md text-sm text-ink-muted">{description}</p> : null}
      </div>
      {action}
    </div>
  );
}
