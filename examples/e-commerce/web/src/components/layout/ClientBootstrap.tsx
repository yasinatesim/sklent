"use client";

import { useEffect } from "react";
import { useThemeStore } from "@/stores/themeStore";
import { useAuthStore } from "@/stores/authStore";
import { API_BASE } from "@/lib/api";

const ClientBootstrap = () => {
  const setTheme = useThemeStore((s) => s.setTheme);
  const theme = useThemeStore((s) => s.theme);
  const setUser = useAuthStore((s) => s.setUser);

  useEffect(() => {
    setTheme(theme);
  }, [setTheme, theme]);

  useEffect(() => {
    const loadSession = async () => {
      try {
        const res = await fetch(`${API_BASE}/auth/me`, { credentials: "include" });
        if (res.ok) {
          setUser(await res.json());
        }
      } catch {
        setUser(null);
      }
    };
    loadSession();
  }, [setUser]);

  return null;
};

export default ClientBootstrap;
