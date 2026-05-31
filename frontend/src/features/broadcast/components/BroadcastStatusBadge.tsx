import type { BroadcastStatus } from "../model/schema";
import { Badge } from "@/shared/ui/Badge";

const statusConfig: Record<
  BroadcastStatus,
  { label: string; variant: "secondary" | "default" | "success" | "warning" | "destructive" }
> = {
  draft: { label: "Draft", variant: "secondary" },
  scheduled: { label: "Scheduled", variant: "default" },
  sending: { label: "Sending", variant: "warning" },
  sent: { label: "Sent", variant: "success" },
  failed: { label: "Failed", variant: "destructive" },
};

export function BroadcastStatusBadge({ status }: { status: BroadcastStatus }) {
  const config = statusConfig[status];
  return <Badge variant={config.variant}>{config.label}</Badge>;
}
