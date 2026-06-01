import { AdminGuestGuard } from "@/features/platform-admin";

export default function AdminLoginLayout({ children }: { children: React.ReactNode }) {
  return <AdminGuestGuard>{children}</AdminGuestGuard>;
}
