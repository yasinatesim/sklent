import { describe, expect, it } from "vitest";

import { PRODUCT_BADGE } from "@/shared/constants/productBadge";

import { productBadgeLabel } from "@/shared/helpers/productBadgeLabel";

const t = (key: string) => `t:${key}`;

describe("productBadgeLabel", () => {
  it("translates every known badge key", () => {
    for (const key of Object.values(PRODUCT_BADGE)) {
      if (!key) continue;
      expect(productBadgeLabel(key, t)).toBe(`t:${key}`);
    }
  });

  it("returns nothing for an empty badge so the card renders no chip", () => {
    expect(productBadgeLabel("", t)).toBe("");
    expect(productBadgeLabel(undefined, t)).toBe("");
  });

  it("falls back to the stored text for rows written before badges were keyed", () => {
    expect(productBadgeLabel("Çok Satan", t)).toBe("Çok Satan");
  });

  it("keeps every badge key locale-neutral so the value is data, not display text", () => {
    for (const key of Object.values(PRODUCT_BADGE)) {
      if (!key) continue;
      expect(key).toMatch(/^[a-z][a-zA-Z0-9]*$/);
    }
  });
});
