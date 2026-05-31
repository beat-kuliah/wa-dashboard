import { AlertCircle, CheckCircle2, Info } from "lucide-react";
import * as React from "react";

import { cn } from "@/shared/lib/cn";

type AlertVariant = "default" | "destructive" | "success";

const variantStyles: Record<AlertVariant, string> = {
  default: "border-border-subtle bg-surface-raised text-ink",
  destructive: "border-destructive/30 bg-destructive/5 text-destructive",
  success: "border-accent-300/40 bg-accent-50 text-accent-800 dark:bg-accent-900/20 dark:text-accent-200",
};

const icons: Record<AlertVariant, React.ReactNode> = {
  default: <Info className="h-4 w-4 shrink-0" />,
  destructive: <AlertCircle className="h-4 w-4 shrink-0" />,
  success: <CheckCircle2 className="h-4 w-4 shrink-0" />,
};

type AlertProps = React.HTMLAttributes<HTMLDivElement> & {
  variant?: AlertVariant;
};

export function Alert({ className, variant = "default", children, ...props }: AlertProps) {
  return (
    <div
      role="alert"
      className={cn(
        "flex gap-3 rounded-md border px-4 py-3 text-sm",
        variantStyles[variant],
        className,
      )}
      {...props}
    >
      {icons[variant]}
      <div className="flex-1">{children}</div>
    </div>
  );
}

export function AlertTitle({
  className,
  ...props
}: React.HTMLAttributes<HTMLParagraphElement>) {
  return <p className={cn("font-medium leading-none", className)} {...props} />;
}

export function AlertDescription({
  className,
  ...props
}: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={cn("mt-1 text-sm opacity-90", className)} {...props} />;
}
