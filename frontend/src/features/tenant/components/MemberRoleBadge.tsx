import type { UserRole } from "@/entities/user";
import { Badge } from "@/shared/ui/Badge";

const roleConfig: Record<
  UserRole,
  { label: string; variant: "default" | "secondary" | "outline" }
> = {
  admin: { label: "Admin", variant: "default" },
  supervisor: { label: "Supervisor", variant: "secondary" },
  agent: { label: "Agent", variant: "outline" },
};

export function MemberRoleBadge({ role }: { role: UserRole }) {
  const config = roleConfig[role] ?? { label: role, variant: "outline" as const };
  return <Badge variant={config.variant}>{config.label}</Badge>;
}
