"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { ROUTES } from "@/shared/config/constants";
import { hasAdminAccessToken } from "@/shared/lib/adminAuthStorage";
import { Spinner } from "@/shared/ui/Spinner";

type AdminGuestGuardProps = {
  children: React.ReactNode;
};

export function AdminGuestGuard({ children }: AdminGuestGuardProps) {
  const router = useRouter();
  const hasToken = hasAdminAccessToken();

  useEffect(() => {
    if (hasToken) {
      router.replace(ROUTES.adminHome);
    }
  }, [hasToken, router]);

  if (hasToken) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  return <>{children}</>;
}
