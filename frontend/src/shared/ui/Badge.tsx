import { cva, type VariantProps } from "class-variance-authority";
import * as React from "react";

import { cn } from "@/shared/lib/cn";

const badgeVariants = cva(
  "inline-flex items-center rounded-sm border px-2 py-0.5 text-xs font-medium transition-colors",
  {
    variants: {
      variant: {
        default: "border-transparent bg-accent-100 text-accent-800 dark:bg-accent-900/40 dark:text-accent-200",
        secondary: "border-border-subtle bg-surface-sunken text-ink-muted",
        outline: "border-border-strong text-ink-muted",
        destructive: "border-transparent bg-destructive/10 text-destructive",
        success: "border-transparent bg-accent-50 text-accent-700 dark:bg-accent-900/30 dark:text-accent-300",
        warning: "border-transparent bg-amber-50 text-amber-800 dark:bg-amber-900/30 dark:text-amber-200",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  },
);

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

export function Badge({ className, variant, ...props }: BadgeProps) {
  return <div className={cn(badgeVariants({ variant }), className)} {...props} />;
}
