"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { ROUTES } from "@/shared/config/constants";
import { hasAdminAccessToken } from "@/shared/lib/adminAuthStorage";
import { Spinner } from "@/shared/ui/Spinner";

type AdminAuthGuardProps = {
  children: React.ReactNode;
};

export function AdminAuthGuard({ children }: AdminAuthGuardProps) {
  const router = useRouter();
  const hasToken = hasAdminAccessToken();

  useEffect(() => {
    if (!hasToken) {
      router.replace(ROUTES.adminLogin);
    }
  }, [hasToken, router]);

  if (!hasToken) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  return <>{children}</>;
}
