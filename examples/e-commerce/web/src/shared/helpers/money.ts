export const formatTRY = (cents: number, locale = "tr"): string =>
  new Intl.NumberFormat(locale === "en" ? "en-US" : "tr-TR", {
    style: "currency",
    currency: "TRY",
    maximumFractionDigits: 0,
  }).format(cents / 100);
