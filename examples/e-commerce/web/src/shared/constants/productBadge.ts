// The stored value is a key, not display text — a badge must read correctly in every locale.
export const PRODUCT_BADGE = {
  NONE: "",
  NEW: "new",
  BESTSELLER: "bestseller",
  DEAL: "deal",
  DISCOUNT10: "discount10",
  DISCOUNT20: "discount20",
} as const;

export type ProductBadge = (typeof PRODUCT_BADGE)[keyof typeof PRODUCT_BADGE];
