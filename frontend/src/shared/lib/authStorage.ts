import { AUTH_STORAGE_KEYS } from "@/shared/config/constants";

export function getAccessToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(AUTH_STORAGE_KEYS.accessToken);
}

export function getRefreshToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(AUTH_STORAGE_KEYS.refreshToken);
}

export function setTokens(accessToken: string, refreshToken: string): void {
  localStorage.setItem(AUTH_STORAGE_KEYS.accessToken, accessToken);
  localStorage.setItem(AUTH_STORAGE_KEYS.refreshToken, refreshToken);
}

export function clearTokens(): void {
  localStorage.removeItem(AUTH_STORAGE_KEYS.accessToken);
  localStorage.removeItem(AUTH_STORAGE_KEYS.refreshToken);
}

export function hasAccessToken(): boolean {
  return Boolean(getAccessToken());
}
