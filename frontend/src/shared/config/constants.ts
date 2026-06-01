export const AUTH_STORAGE_KEYS = {
  accessToken: "wa_access_token",
  refreshToken: "wa_refresh_token",
} as const;

export const ADMIN_AUTH_STORAGE_KEYS = {
  accessToken: "wa_admin_access_token",
  refreshToken: "wa_admin_refresh_token",
} as const;

export const ROUTES = {
  login: "/login",
  broadcast: "/broadcast",
  csInbox: "/cs-inbox",
  analytics: "/analytics",
  settings: "/settings",
  adminLogin: "/admin/login",
  adminHome: "/admin",
} as const;
