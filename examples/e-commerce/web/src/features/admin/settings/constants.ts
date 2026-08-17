// Mirrors ORDER_SOURCE_* in api/internal/constants; the marketplace names are brands, not translatable copy.
export const ORDER_NOTIFY_SOURCES = [
  { value: "site", label: "Site" },
  { value: "HB", label: "Hepsiburada" },
  { value: "TY", label: "Trendyol" },
] as const;

export const ORDER_NOTIFY_SOURCE_LABELS: Record<string, string> = Object.fromEntries(
  ORDER_NOTIFY_SOURCES.map((source) => [source.value, source.label]),
);
