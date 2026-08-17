import { describe, expect, it } from "vitest";

import { PRODUCT_SORT } from "@/features/web/catalog/types/catalogQuery.types";

import { buildCatalogSearchParams, totalPages } from "@/features/web/catalog/helpers/catalogQuery";

describe("buildCatalogSearchParams", () => {
  it("omits empty filters so the URL stays readable and cacheable", () => {
    expect(buildCatalogSearchParams({}).toString()).toBe("");
  });

  it("serialises every supplied filter", () => {
    const params = buildCatalogSearchParams({
      category: "rings",
      q: "gümüş",
      page: 2,
      minPrice: 1000,
      maxPrice: 5000,
      sort: PRODUCT_SORT.PRICE_ASC,
      inStock: true,
      minRating: 4,
    });
    expect(params.get("category")).toBe("rings");
    expect(params.get("q")).toBe("gümüş");
    expect(params.get("page")).toBe("2");
    expect(params.get("minPrice")).toBe("1000");
    expect(params.get("maxPrice")).toBe("5000");
    expect(params.get("sort")).toBe(PRODUCT_SORT.PRICE_ASC);
    expect(params.get("inStock")).toBe("true");
    expect(params.get("minRating")).toBe("4");
  });

  it("drops page 1 — the default must not clutter every link", () => {
    expect(buildCatalogSearchParams({ page: 1 }).has("page")).toBe(false);
  });

  it("omits inStock when false rather than sending inStock=false", () => {
    expect(buildCatalogSearchParams({ inStock: false }).has("inStock")).toBe(false);
  });

  it("never emits a zero price bound, which would read as a real filter", () => {
    const params = buildCatalogSearchParams({ minPrice: 0, maxPrice: 0 });
    expect(params.has("minPrice")).toBe(false);
    expect(params.has("maxPrice")).toBe(false);
  });
});

describe("totalPages", () => {
  it("rounds up so a partial last page still exists", () => {
    expect(totalPages({ page: 1, pageSize: 12, total: 13 })).toBe(2);
    expect(totalPages({ page: 1, pageSize: 12, total: 24 })).toBe(2);
  });

  it("is 1 for an empty result so the pager never renders zero pages", () => {
    expect(totalPages({ page: 1, pageSize: 12, total: 0 })).toBe(1);
  });
});
