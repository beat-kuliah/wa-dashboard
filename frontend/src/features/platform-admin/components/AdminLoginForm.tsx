"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import type { KeyboardEvent } from "react";
import { useForm } from "react-hook-form";

import { ROUTES } from "@/shared/config/constants";
import { ApiError } from "@/shared/lib/apiClient";
import { Alert, AlertDescription } from "@/shared/ui/Alert";
import { Button } from "@/shared/ui/Button";
import { Input } from "@/shared/ui/Input";
import { Label } from "@/shared/ui/Label";
import { Spinner } from "@/shared/ui/Spinner";

import { useAdminLogin } from "../api/queries";
import { adminLoginFormSchema, type AdminLoginFormValues } from "../model/schema";

export function AdminLoginForm() {
  const router = useRouter();
  const loginMutation = useAdminLogin();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<AdminLoginFormValues>({
    resolver: zodResolver(adminLoginFormSchema),
    defaultValues: { email: "", password: "" },
  });

  const submit = handleSubmit(async (values) => {
    try {
      await loginMutation.mutateAsync(values);
      router.replace(ROUTES.adminHome);
    } catch {
      // surfaced via mutation state
    }
  });

  const handleEnterSubmit = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key !== "Enter" || event.nativeEvent.isComposing) {
      return;
    }
    event.preventDefault();
    event.currentTarget.form?.requestSubmit();
  };

  const apiError =
    loginMutation.error instanceof ApiError
      ? loginMutation.error.message
      : loginMutation.error
        ? "Unable to sign in. Please try again."
        : null;

  return (
    <form noValidate onSubmit={submit} className="space-y-5">
      {apiError ? (
        <Alert variant="destructive">
          <AlertDescription>{apiError}</AlertDescription>
        </Alert>
      ) : null}

      <div className="space-y-2">
        <Label htmlFor="admin-email">Email</Label>
        <Input
          id="admin-email"
          type="email"
          autoComplete="email"
          placeholder="ops@wa-dashboard.com"
          {...register("email")}
          onKeyDown={handleEnterSubmit}
        />
        {errors.email ? (
          <p className="text-sm text-destructive">{errors.email.message}</p>
        ) : null}
      </div>

      <div className="space-y-2">
        <Label htmlFor="admin-password">Password</Label>
        <Input
          id="admin-password"
          type="password"
          autoComplete="current-password"
          placeholder="••••••••"
          {...register("password")}
          onKeyDown={handleEnterSubmit}
        />
        {errors.password ? (
          <p className="text-sm text-destructive">{errors.password.message}</p>
        ) : null}
      </div>

      <Button type="submit" className="w-full" disabled={loginMutation.isPending}>
        {loginMutation.isPending ? (
          <>
            <Spinner size="sm" />
            Signing in…
          </>
        ) : (
          "Sign in to console"
        )}
      </Button>
    </form>
  );
}
