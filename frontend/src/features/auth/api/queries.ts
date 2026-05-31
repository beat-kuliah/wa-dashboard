import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { fetchCurrentUser, login, logout } from "./authApi";
import { useAuthStore } from "../model/store";
import { getRefreshToken } from "@/shared/lib/authStorage";
import type { LoginFormValues } from "../model/schema";

export const authKeys = {
  all: ["auth"] as const,
  me: () => [...authKeys.all, "me"] as const,
};

export function useCurrentUser(enabled = true) {
  const setUser = useAuthStore((s) => s.setUser);

  return useQuery({
    queryKey: authKeys.me(),
    queryFn: async () => {
      const user = await fetchCurrentUser();
      setUser(user);
      return user;
    },
    enabled,
    retry: false,
  });
}

export function useLogin() {
  const queryClient = useQueryClient();
  const setUser = useAuthStore((s) => s.setUser);

  return useMutation({
    mutationFn: (values: LoginFormValues) => login(values),
    onSuccess: (data) => {
      setUser(data.user);
      queryClient.setQueryData(authKeys.me(), data.user);
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();
  const clearAuth = useAuthStore((s) => s.clearAuth);

  return useMutation({
    mutationFn: () => logout(getRefreshToken()),
    onSettled: () => {
      clearAuth();
      queryClient.clear();
    },
  });
}
