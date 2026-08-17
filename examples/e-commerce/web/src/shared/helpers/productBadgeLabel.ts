import { PRODUCT_BADGE } from "@/shared/constants/productBadge";

const KNOWN_BADGES = new Set<string>(Object.values(PRODUCT_BADGE).filter(Boolean));

// Rows written before badges were keyed hold display text; render those verbatim.
export const productBadgeLabel = (badge: string | undefined, t: (key: string) => string): string => {
  if (!badge) return "";
  return KNOWN_BADGES.has(badge) ? t(badge) : badge;
};

export default productBadgeLabel;
