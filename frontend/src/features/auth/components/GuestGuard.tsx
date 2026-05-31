"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { useCurrentUser } from "@/features/auth";
import { ROUTES } from "@/shared/config/constants";
import { hasAccessToken } from "@/shared/lib/authStorage";
import { Spinner } from "@/shared/ui/Spinner";

type GuestGuardProps = {
  children: React.ReactNode;
};

export function GuestGuard({ children }: GuestGuardProps) {
  const router = useRouter();
  const hasToken = hasAccessToken();
  const { isLoading, isSuccess } = useCurrentUser(hasToken);

  useEffect(() => {
    if (hasToken && isSuccess) {
      router.replace(ROUTES.broadcast);
    }
  }, [hasToken, isSuccess, router]);

  if (hasToken && isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  if (hasToken && isSuccess) {
    return null;
  }

  return <>{children}</>;
}
