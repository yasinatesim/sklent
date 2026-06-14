import { getRequestConfig } from "next-intl/server";
import { routing } from "./routing";

const loadConfig = getRequestConfig(async ({ requestLocale }) => {
  const requested = await requestLocale;
  const locale = routing.locales.includes(requested as never)
    ? (requested as string)
    : routing.defaultLocale;

  const messages = (await import(`./messages/${locale}.json`)).default;
  return { locale, messages };
});

export default loadConfig;
