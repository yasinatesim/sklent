import { defineRouting } from "next-intl/routing";

export const LOCALES = ["tr", "en"] as const;
export const DEFAULT_LOCALE = "tr";

export const routing = defineRouting({
  locales: LOCALES,
  defaultLocale: DEFAULT_LOCALE,
});
