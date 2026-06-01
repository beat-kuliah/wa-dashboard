import { env } from "@/shared/config/env";

import { apiErrorEnvelopeSchema, ApiError, type ApiErrorCode } from "./apiClient";
import { getAdminAccessToken } from "./adminAuthStorage";

export type AdminApiClientOptions = RequestInit & {
  skipAuth?: boolean;
};

export async function adminApiClient<T>(
  path: string,
  options: AdminApiClientOptions = {},
): Promise<T> {
  const { skipAuth = false, headers, ...rest } = options;

  const requestHeaders = new Headers(headers);
  if (!requestHeaders.has("Content-Type") && rest.body) {
    requestHeaders.set("Content-Type", "application/json");
  }

  if (!skipAuth) {
    const token = getAdminAccessToken();
    if (token) {
      requestHeaders.set("Authorization", `Bearer ${token}`);
    }
  }

  const url = path.startsWith("http")
    ? path
    : `${env.NEXT_PUBLIC_API_BASE_URL}${path}`;

  const response = await fetch(url, {
    ...rest,
    headers: requestHeaders,
  });

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
