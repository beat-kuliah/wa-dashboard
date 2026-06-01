import { create } from "zustand";

import type { PlatformAdmin } from "./schema";

type PlatformAdminState = {
  admin: PlatformAdmin | null;
  isAuthenticated: boolean;
  setAdmin: (admin: PlatformAdmin | null) => void;
  clearAdmin: () => void;
};

export const usePlatformAdminStore = create<PlatformAdminState>((set) => ({
  admin: null,
  isAuthenticated: false,
  setAdmin: (admin) => set({ admin, isAuthenticated: admin !== null }),
  clearAdmin: () => set({ admin: null, isAuthenticated: false }),
}));
