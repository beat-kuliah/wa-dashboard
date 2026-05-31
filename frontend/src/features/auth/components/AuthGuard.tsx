"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";

import { ROUTES } from "@/shared/config/constants";
import { hasAccessToken } from "@/shared/lib/authStorage";
import { Spinner } from "@/shared/ui/Spinner";

import { useCurrentUser } from "../api/queries";

type AuthGuardProps = {
  children: React.ReactNode;
};

export function AuthGuard({ children }: AuthGuardProps) {
  const router = useRouter();
  const hasToken = hasAccessToken();
  const { isLoading, isError } = useCurrentUser(hasToken);

  useEffect(() => {
    if (!hasToken) {
      router.replace(ROUTES.login);
    }
  }, [hasToken, router]);

  useEffect(() => {
    if (isError) {
      router.replace(ROUTES.login);
    }
  }, [isError, router]);

  if (!hasToken || isLoading) {
    return (
      <div className="flex min-h-[50vh] items-center justify-center">
        <Spinner size="lg" />
      </div>
    );
  }

  if (isError) {
    return null;
  }

  return <>{children}</>;
}
