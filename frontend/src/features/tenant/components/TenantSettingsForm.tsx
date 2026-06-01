"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { useForm } from "react-hook-form";

import { ApiError } from "@/shared/lib/apiClient";
import { Alert, AlertDescription } from "@/shared/ui/Alert";
import { Button } from "@/shared/ui/Button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/Card";
import { ErrorState } from "@/shared/ui/ErrorState";
import { Input } from "@/shared/ui/Input";
import { Label } from "@/shared/ui/Label";
import { Skeleton } from "@/shared/ui/Skeleton";
import { Spinner } from "@/shared/ui/Spinner";

import { useCurrentTenant, useUpdateTenant } from "../api/queries";
import {
  updateTenantFormSchema,
  type UpdateTenantFormValues,
} from "../model/schema";

export function TenantSettingsForm({ canEdit }: { canEdit: boolean }) {
  const { data: tenant, isLoading, isError, error, refetch } = useCurrentTenant();
  const updateMutation = useUpdateTenant();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isDirty },
  } = useForm<UpdateTenantFormValues>({
    resolver: zodResolver(updateTenantFormSchema),
    defaultValues: { name: "" },
  });

  useEffect(() => {
    if (tenant) {
      reset({ name: tenant.name });
    }
  }, [tenant, reset]);

  const submit = handleSubmit(async (values) => {
    try {
      const updated = await updateMutation.mutateAsync({ name: values.name });
      reset({ name: updated.name });
    } catch {
      // Error surfaced via mutation state
    }
  });

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Workspace</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <Skeleton className="h-4 w-24" />
          <Skeleton className="h-10 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <ErrorState
        message={error instanceof Error ? error.message : "Failed to load workspace"}
        onRetry={() => void refetch()}
      />
    );
  }

  const apiError =
    updateMutation.error instanceof ApiError
      ? updateMutation.error.message
      : updateMutation.error
        ? "Unable to save changes. Please try again."
        : null;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Workspace</CardTitle>
        <CardDescription>
          {canEdit
            ? "Update your workspace details."
            : "Workspace details. Only admins can make changes."}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form noValidate onSubmit={submit} className="space-y-5">
          {apiError ? (
            <Alert variant="destructive">
              <AlertDescription>{apiError}</AlertDescription>
            </Alert>
          ) : null}

          {updateMutation.isSuccess && !isDirty ? (
            <Alert variant="success">
              <AlertDescription>Workspace updated.</AlertDescription>
            </Alert>
          ) : null}

          <div className="space-y-2">
            <Label htmlFor="tenant-name">Workspace name</Label>
            <Input
              id="tenant-name"
              placeholder="Acme Inc."
              disabled={!canEdit}
              {...register("name")}
            />
            {errors.name ? (
              <p className="text-sm text-destructive">{errors.name.message}</p>
            ) : null}
          </div>

          {tenant?.slug ? (
            <div className="space-y-2">
              <Label htmlFor="tenant-slug">Slug</Label>
              <Input id="tenant-slug" value={tenant.slug} disabled readOnly />
              <p className="text-xs text-ink-subtle">
                The workspace slug is generated automatically.
              </p>
            </div>
          ) : null}

          {canEdit ? (
            <Button
              type="submit"
              disabled={updateMutation.isPending || !isDirty}
            >
              {updateMutation.isPending ? (
                <>
                  <Spinner size="sm" />
                  Saving…
                </>
              ) : (
                "Save changes"
              )}
            </Button>
          ) : null}
        </form>
      </CardContent>
    </Card>
  );
}
