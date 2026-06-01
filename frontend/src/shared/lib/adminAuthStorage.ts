import { ADMIN_AUTH_STORAGE_KEYS } from "@/shared/config/constants";

export function getAdminAccessToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(ADMIN_AUTH_STORAGE_KEYS.accessToken);
}

export function getAdminRefreshToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem(ADMIN_AUTH_STORAGE_KEYS.refreshToken);
}

export function setAdminTokens(accessToken: string, refreshToken: string): void {
  localStorage.setItem(ADMIN_AUTH_STORAGE_KEYS.accessToken, accessToken);
  localStorage.setItem(ADMIN_AUTH_STORAGE_KEYS.refreshToken, refreshToken);
}

export function clearAdminTokens(): void {
  localStorage.removeItem(ADMIN_AUTH_STORAGE_KEYS.accessToken);
  localStorage.removeItem(ADMIN_AUTH_STORAGE_KEYS.refreshToken);
}

export function hasAdminAccessToken(): boolean {
  return Boolean(getAdminAccessToken());
}
