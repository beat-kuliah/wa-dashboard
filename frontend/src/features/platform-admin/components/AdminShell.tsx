"use client";

import { LogOut, Shield } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { ROUTES } from "@/shared/config/constants";
import { Button } from "@/shared/ui/Button";
import { Separator } from "@/shared/ui/Separator";
import { ThemeToggle } from "@/shared/ui/ThemeToggle";

import { useAdminLogout } from "../api/queries";
import { usePlatformAdminStore } from "../model/store";

type AdminShellProps = {
  children: React.ReactNode;
};

export function AdminShell({ children }: AdminShellProps) {
  const router = useRouter();
  const admin = usePlatformAdminStore((s) => s.admin);
  const logout = useAdminLogout();

  const handleLogout = () => {
    logout.mutate(undefined, {
      onSettled: () => router.replace(ROUTES.adminLogin),
    });
  };

  return (
    <div className="flex min-h-screen flex-col bg-surface">
      <header className="flex h-14 shrink-0 items-center justify-between border-b border-border-subtle bg-surface-raised px-4 lg:px-8">
        <div className="flex items-center gap-3">
          <Shield className="h-5 w-5 text-accent" />
          <Link href={ROUTES.adminHome} className="font-display text-sm font-bold tracking-tight">
            Platform<span className="text-accent">.</span>Console
          </Link>
        </div>
        <div className="flex items-center gap-2">
          {admin ? (
            <span className="hidden text-sm text-ink-muted sm:inline">{admin.email}</span>
          ) : null}
          <ThemeToggle />
          <Separator orientation="vertical" className="mx-1 h-6" />
          <Button
            variant="ghost"
            size="icon"
            aria-label="Sign out"
            onClick={handleLogout}
            disabled={logout.isPending}
          >
            <LogOut className="h-4 w-4" />
          </Button>
        </div>
      </header>
      <main className="flex-1 overflow-auto p-4 lg:p-8">{children}</main>
    </div>
  );
}
