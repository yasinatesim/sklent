import { create } from "zustand";

export type SessionUser = {
  userId: string;
  email: string;
  role: string;
};

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
