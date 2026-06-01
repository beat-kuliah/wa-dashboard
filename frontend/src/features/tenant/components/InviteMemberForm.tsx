"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { ApiError } from "@/shared/lib/apiClient";
import { cn } from "@/shared/lib/cn";
import { Alert, AlertDescription } from "@/shared/ui/Alert";
import { Button } from "@/shared/ui/Button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/shared/ui/Card";
import { Input } from "@/shared/ui/Input";
import { Label } from "@/shared/ui/Label";
import { Spinner } from "@/shared/ui/Spinner";

import { useAddMember } from "../api/queries";
import {
  addMemberFormSchema,
  type AddMemberFormValues,
} from "../model/schema";

const roleOptions = [
  { value: "admin", label: "Admin" },
  { value: "supervisor", label: "Supervisor" },
  { value: "agent", label: "Agent" },
] as const;

export function InviteMemberForm() {
  const addMutation = useAddMember();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<AddMemberFormValues>({
    resolver: zodResolver(addMemberFormSchema),
    defaultValues: { email: "", full_name: "", role: "agent" },
  });

  const submit = handleSubmit(async (values) => {
    try {
      await addMutation.mutateAsync(values);
      reset({ email: "", full_name: "", role: "agent" });
    } catch {
      // Error surfaced via mutation state
    }
  });

  const apiError =
    addMutation.error instanceof ApiError
      ? addMutation.error.message
      : addMutation.error
        ? "Unable to invite member. Please try again."
        : null;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Invite member</CardTitle>
        <CardDescription>
          Add a teammate to this workspace and assign them a role.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form noValidate onSubmit={submit} className="space-y-5">
          {apiError ? (
            <Alert variant="destructive">
              <AlertDescription>{apiError}</AlertDescription>
            </Alert>
          ) : null}

          {addMutation.isSuccess ? (
            <Alert variant="success">
              <AlertDescription>Member invited.</AlertDescription>
            </Alert>
          ) : null}

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="member-full-name">Full name</Label>
              <Input
                id="member-full-name"
                placeholder="John Agent"
                autoComplete="name"
                {...register("full_name")}
              />
              {errors.full_name ? (
                <p className="text-sm text-destructive">{errors.full_name.message}</p>
              ) : null}
            </div>

            <div className="space-y-2">
              <Label htmlFor="member-email">Email</Label>
              <Input
                id="member-email"
                type="email"
                placeholder="teammate@company.com"
                autoComplete="off"
                {...register("email")}
              />
              {errors.email ? (
                <p className="text-sm text-destructive">{errors.email.message}</p>
              ) : null}
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="member-role">Role</Label>
            <select
              id="member-role"
              className={cn(
                "flex h-10 w-full rounded-md border border-border-subtle bg-surface-raised px-3 py-2 text-sm text-ink shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
              )}
              {...register("role")}
            >
              {roleOptions.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
            {errors.role ? (
              <p className="text-sm text-destructive">{errors.role.message}</p>
            ) : null}
          </div>

          <Button type="submit" disabled={addMutation.isPending}>
            {addMutation.isPending ? (
              <>
                <Spinner size="sm" />
                Inviting…
              </>
            ) : (
              "Invite member"
            )}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
