import { create } from "zustand";

export const THEME = { LIGHT: "light", DARK: "dark" } as const;
export type Theme = (typeof THEME)[keyof typeof THEME];

type ThemeState = {
  theme: Theme;
  setTheme: (theme: Theme) => void;
  toggle: () => void;
};

const applyTheme = (theme: Theme) => {
  if (typeof document !== "undefined") {
    document.documentElement.setAttribute("data-theme", theme);
    window.localStorage.setItem("vcTheme", theme);
  }
};

const initialTheme = (): Theme => {
  if (typeof window === "undefined") return THEME.LIGHT;
  return (window.localStorage.getItem("vcTheme") as Theme) ?? THEME.LIGHT;
};

export const useThemeStore = create<ThemeState>((set, get) => ({
  theme: initialTheme(),
  setTheme: (theme) => {
    applyTheme(theme);
    set({ theme });
  },
  toggle: () => {
    const next = get().theme === THEME.DARK ? THEME.LIGHT : THEME.DARK;
    applyTheme(next);
    set({ theme: next });
  },
}));
