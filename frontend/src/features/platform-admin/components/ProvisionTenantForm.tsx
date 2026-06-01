"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { ApiError } from "@/shared/lib/apiClient";
import { Alert, AlertDescription } from "@/shared/ui/Alert";
import { Button } from "@/shared/ui/Button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/shared/ui/Card";
import { Input } from "@/shared/ui/Input";
import { Label } from "@/shared/ui/Label";
import { Spinner } from "@/shared/ui/Spinner";

import { useProvisionTenant } from "../api/queries";
import { provisionTenantFormSchema, type ProvisionTenantFormValues } from "../model/schema";

export function ProvisionTenantForm() {
  const provisionMutation = useProvisionTenant();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ProvisionTenantFormValues>({
    resolver: zodResolver(provisionTenantFormSchema),
    defaultValues: {
      business_name: "",
      owner_email: "",
      owner_full_name: "",
      owner_password: "",
    },
  });

  const submit = handleSubmit(async (values) => {
    try {
      await provisionMutation.mutateAsync(values);
      reset();
    } catch {
      // surfaced via mutation state
    }
  });

  const apiError =
    provisionMutation.error instanceof ApiError
      ? provisionMutation.error.message
      : provisionMutation.error
        ? "Unable to provision tenant."
        : null;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Provision tenant</CardTitle>
        <CardDescription>
          Create a new workspace and its initial admin owner in one step.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form noValidate onSubmit={submit} className="space-y-4">
          {apiError ? (
            <Alert variant="destructive">
              <AlertDescription>{apiError}</AlertDescription>
            </Alert>
          ) : null}
          {provisionMutation.isSuccess ? (
            <Alert>
              <AlertDescription>
                Tenant provisioned. Owner: {provisionMutation.data.owner.email}
              </AlertDescription>
            </Alert>
          ) : null}

          <div className="space-y-2">
            <Label htmlFor="business-name">Business name</Label>
            <Input id="business-name" {...register("business_name")} />
            {errors.business_name ? (
              <p className="text-sm text-destructive">{errors.business_name.message}</p>
            ) : null}
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="owner-full-name">Owner full name</Label>
              <Input id="owner-full-name" {...register("owner_full_name")} />
              {errors.owner_full_name ? (
                <p className="text-sm text-destructive">{errors.owner_full_name.message}</p>
              ) : null}
            </div>
            <div className="space-y-2">
              <Label htmlFor="owner-email">Owner email</Label>
              <Input id="owner-email" type="email" {...register("owner_email")} />
              {errors.owner_email ? (
                <p className="text-sm text-destructive">{errors.owner_email.message}</p>
              ) : null}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="owner-password">Owner password</Label>
            <Input
              id="owner-password"
              type="password"
              autoComplete="new-password"
              {...register("owner_password")}
            />
            {errors.owner_password ? (
              <p className="text-sm text-destructive">{errors.owner_password.message}</p>
            ) : null}
          </div>

          <Button type="submit" disabled={provisionMutation.isPending}>
            {provisionMutation.isPending ? (
              <>
                <Spinner size="sm" />
                Provisioning…
              </>
            ) : (
              "Provision tenant"
            )}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
