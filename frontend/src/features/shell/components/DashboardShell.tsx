"use client";

import {
  BarChart3,
  Building2,
  ChevronDown,
  LogOut,
  Megaphone,
  MessageSquare,
  Menu,
  Settings,
  X,
} from "lucide-react";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";

import { useLogout } from "@/features/auth";
import { useCurrentTenant } from "@/features/tenant";
import { ROUTES } from "@/shared/config/constants";
import { useIsMobile } from "@/shared/hooks/useMediaQuery";
import { cn } from "@/shared/lib/cn";
import { useShellStore } from "@/shared/model/shellStore";
import { Button } from "@/shared/ui/Button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/shared/ui/DropdownMenu";
import { Separator } from "@/shared/ui/Separator";
import { Skeleton } from "@/shared/ui/Skeleton";
import { ThemeToggle } from "@/shared/ui/ThemeToggle";

const navItems = [
  { href: ROUTES.broadcast, label: "Broadcasts", icon: Megaphone },
  { href: ROUTES.csInbox, label: "CS Inbox", icon: MessageSquare },
  { href: ROUTES.analytics, label: "Analytics", icon: BarChart3 },
  { href: ROUTES.settings, label: "Settings", icon: Settings },
] as const;

function SidebarNav({ onNavigate }: { onNavigate?: () => void }) {
  const pathname = usePathname();

  return (
    <nav className="flex flex-col gap-1 px-3">
      {navItems.map(({ href, label, icon: Icon }) => {
        const isActive = pathname.startsWith(href);
        return (
          <Link
            key={href}
            href={href}
            onClick={onNavigate}
            className={cn(
              "flex items-center gap-3 rounded-md px-3 py-2.5 text-sm font-medium transition-colors",
              isActive
                ? "bg-accent-50 text-accent-800 dark:bg-accent-900/30 dark:text-accent-200"
                : "text-ink-muted hover:bg-surface-sunken hover:text-ink",
            )}
          >
            <Icon className="h-4 w-4 shrink-0" />
            {label}
          </Link>
        );
      })}
    </nav>
  );
}

function TenantSwitcher() {
  const { data: tenant, isLoading } = useCurrentTenant();

  if (isLoading) {
    return <Skeleton className="h-9 w-36" />;
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" size="sm" className="max-w-[200px] justify-between gap-2">
          <Building2 className="h-4 w-4 shrink-0 text-ink-muted" />
          <span className="truncate">{tenant?.name ?? "Workspace"}</span>
          <ChevronDown className="h-3.5 w-3.5 shrink-0 opacity-60" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-56">
        <DropdownMenuLabel>Current workspace</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem disabled>
          <div className="flex flex-col">
            <span className="font-medium">{tenant?.name}</span>
            <span className="text-xs text-ink-muted">{tenant?.slug}</span>
          </div>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

function TopBar({ onMenuClick }: { onMenuClick?: () => void }) {
  const router = useRouter();
  const logout = useLogout();

  const handleLogout = () => {
    logout.mutate(undefined, {
      onSettled: () => router.replace(ROUTES.login),
    });
  };

  return (
    <header className="flex h-14 shrink-0 items-center justify-between border-b border-border-subtle bg-surface-raised px-4 lg:px-6">
      <div className="flex items-center gap-3">
        {onMenuClick ? (
          <Button variant="ghost" size="icon" onClick={onMenuClick} className="lg:hidden">
            <Menu className="h-5 w-5" />
          </Button>
        ) : null}
        <span className="font-display text-sm font-semibold tracking-tight lg:hidden">
          WA Dashboard
        </span>
      </div>
      <div className="flex items-center gap-2">
        <TenantSwitcher />
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
  );
}

type DashboardShellProps = {
  children: React.ReactNode;
};

export function DashboardShell({ children }: DashboardShellProps) {
  const isMobile = useIsMobile();
  const isSidebarOpen = useShellStore((s) => s.isSidebarOpen);
  const setSidebarOpen = useShellStore((s) => s.setSidebarOpen);
  const toggleSidebar = useShellStore((s) => s.toggleSidebar);

  return (
    <div className="flex min-h-screen bg-surface">
      {/* Desktop sidebar */}
      <aside className="hidden w-60 shrink-0 flex-col border-r border-border-subtle bg-surface-raised lg:flex">
        <div className="flex h-14 items-center px-6">
          <Link href={ROUTES.broadcast} className="font-display text-base font-bold tracking-tight">
            WA<span className="text-accent">.</span>Dashboard
          </Link>
        </div>
        <Separator />
        <div className="flex-1 py-4">
          <SidebarNav />
        </div>
      </aside>

      {/* Mobile sidebar overlay */}
      {isMobile && isSidebarOpen ? (
        <div className="fixed inset-0 z-40 lg:hidden">
          <button
            type="button"
            className="absolute inset-0 bg-ink/20 backdrop-blur-sm"
            aria-label="Close menu"
            onClick={() => setSidebarOpen(false)}
          />
          <aside className="relative flex h-full w-72 flex-col bg-surface-raised shadow-lg">
            <div className="flex h-14 items-center justify-between px-4">
              <span className="font-display text-base font-bold tracking-tight">
                WA<span className="text-accent">.</span>Dashboard
              </span>
              <Button variant="ghost" size="icon" onClick={() => setSidebarOpen(false)}>
                <X className="h-5 w-5" />
              </Button>
            </div>
            <Separator />
            <div className="flex-1 py-4">
              <SidebarNav onNavigate={() => setSidebarOpen(false)} />
            </div>
          </aside>
        </div>
      ) : null}

      <div className="flex min-w-0 flex-1 flex-col">
        <TopBar onMenuClick={isMobile ? toggleSidebar : undefined} />
        <main className="flex-1 overflow-auto p-4 lg:p-8">{children}</main>
      </div>
    </div>
  );
}
