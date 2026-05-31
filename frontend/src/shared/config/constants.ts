export const AUTH_STORAGE_KEYS = {
  accessToken: "wa_access_token",
  refreshToken: "wa_refresh_token",
} as const;

export const ROUTES = {
  login: "/login",
  broadcast: "/broadcast",
  csInbox: "/cs-inbox",
  analytics: "/analytics",
} as const;
