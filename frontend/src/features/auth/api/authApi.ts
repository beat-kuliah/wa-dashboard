import { userResponseSchema } from "@/entities/user";
import { apiClient } from "@/shared/lib/apiClient";
import { clearTokens, setTokens } from "@/shared/lib/authStorage";

import { loginResponseSchema, type LoginFormValues } from "../model/schema";

export async function login(values: LoginFormValues) {
  const json = await apiClient<unknown>("/auth/login", {
    method: "POST",
    body: JSON.stringify(values),
    skipAuth: true,
    skipRefresh: true,
  });

  const parsed = loginResponseSchema.parse(json);
  setTokens(parsed.tokens.access_token, parsed.tokens.refresh_token);
  return parsed;
}

export async function logout(refreshToken: string | null) {
  if (refreshToken) {
    try {
      await apiClient<void>("/auth/logout", {
        method: "POST",
        body: JSON.stringify({ refresh_token: refreshToken }),
      });
    } catch {
      // Best-effort logout
    }
  }
  clearTokens();
}

export async function fetchCurrentUser() {
  const json = await apiClient<unknown>("/auth/me");
  return userResponseSchema.parse(json).user;
}
