import { z } from "zod";

import { env } from "@/shared/config/env";
import {
  clearTokens,
  getAccessToken,
  getRefreshToken,
  setTokens,
} from "@/shared/lib/authStorage";

export const apiErrorEnvelopeSchema = z.object({
  error: z.object({
    code: z.string(),
    message: z.string(),
  }),
});

export type ApiErrorCode =
  | "VALIDATION"
  | "UNAUTHORIZED"
  | "FORBIDDEN"
  | "NOT_FOUND"
  | "CONFLICT"
  | "RATE_LIMITED"
  | "INTERNAL";

export class ApiError extends Error {
  readonly code: ApiErrorCode;
  readonly status: number;

  constructor(code: ApiErrorCode, message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.code = code;
    this.status = status;
  }
}

const tokensResponseSchema = z.object({
  tokens: z.object({
    access_token: z.string(),
    refresh_token: z.string(),
    expires_in: z.number(),
  }),
});

let refreshPromise: Promise<boolean> | null = null;

async function refreshAccessToken(): Promise<boolean> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) return false;

  const response = await fetch(`${env.NEXT_PUBLIC_API_BASE_URL}/auth/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refreshToken }),
  });

  if (!response.ok) {
    clearTokens();
    return false;
  }

  const json: unknown = await response.json();
  const parsed = tokensResponseSchema.safeParse(json);
  if (!parsed.success) {
    clearTokens();
    return false;
  }

  setTokens(parsed.data.tokens.access_token, parsed.data.tokens.refresh_token);
  return true;
}

async function ensureRefreshed(): Promise<boolean> {
  if (!refreshPromise) {
    refreshPromise = refreshAccessToken().finally(() => {
      refreshPromise = null;
    });
  }
  return refreshPromise;
}

export type ApiClientOptions = RequestInit & {
  skipAuth?: boolean;
  skipRefresh?: boolean;
};

export async function apiClient<T>(
  path: string,
  options: ApiClientOptions = {},
): Promise<T> {
  const { skipAuth = false, skipRefresh = false, headers, ...rest } = options;

  const requestHeaders = new Headers(headers);
  if (!requestHeaders.has("Content-Type") && rest.body) {
    requestHeaders.set("Content-Type", "application/json");
  }

  if (!skipAuth) {
    const token = getAccessToken();
    if (token) {
      requestHeaders.set("Authorization", `Bearer ${token}`);
    }
  }

  const url = path.startsWith("http")
    ? path
    : `${env.NEXT_PUBLIC_API_BASE_URL}${path}`;

  let response = await fetch(url, {
    ...rest,
    headers: requestHeaders,
  });

  if (response.status === 401 && !skipAuth && !skipRefresh) {
    const refreshed = await ensureRefreshed();
    if (refreshed) {
      const retryHeaders = new Headers(headers);
      if (!retryHeaders.has("Content-Type") && rest.body) {
        retryHeaders.set("Content-Type", "application/json");
      }
      const token = getAccessToken();
      if (token) {
        retryHeaders.set("Authorization", `Bearer ${token}`);
      }
      response = await fetch(url, { ...rest, headers: retryHeaders });
    }
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const json: unknown = await response.json().catch(() => null);

  if (!response.ok) {
    const parsed = apiErrorEnvelopeSchema.safeParse(json);
    if (parsed.success) {
      throw new ApiError(
        parsed.data.error.code as ApiErrorCode,
        parsed.data.error.message,
        response.status,
      );
    }
    throw new ApiError("INTERNAL", "An unexpected error occurred", response.status);
  }

  return json as T;
}

export const paginationSchema = z.object({
  limit: z.number(),
  offset: z.number(),
  total: z.number(),
});

export type Pagination = z.infer<typeof paginationSchema>;

export function buildQueryString(
  params: Record<string, string | number | undefined | null>,
): string {
  const search = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== null && value !== "") {
      search.set(key, String(value));
    }
  }
  const qs = search.toString();
  return qs ? `?${qs}` : "";
}
