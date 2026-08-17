import { create } from "zustand";

import type { SessionUser } from "@/features/web/auth/types/auth.types";

type AuthState = {
  user: SessionUser | null;
  setUser: (user: SessionUser | null) => void;
  isAdmin: () => boolean;
};

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  setUser: (user) => set({ user }),
  isAdmin: () => get().user?.role === "admin",
}));
