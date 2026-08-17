import type { CatalogQuery, Pagination } from "@/features/web/catalog/types/catalogQuery.types";

const FIRST_PAGE = 1;

// Only meaningful filters reach the URL: a default page or a zero bound would read as a real
// filter and would fragment the cache for no reason.
export const buildCatalogSearchParams = (query: CatalogQuery): URLSearchParams => {
  const params = new URLSearchParams();
  if (query.category) params.set("category", query.category);
  if (query.q) params.set("q", query.q);
  if (query.page && query.page > FIRST_PAGE) params.set("page", String(query.page));
  if (query.pageSize) params.set("pageSize", String(query.pageSize));
  if (query.minPrice) params.set("minPrice", String(query.minPrice));
  if (query.maxPrice) params.set("maxPrice", String(query.maxPrice));
  if (query.sort) params.set("sort", query.sort);
  if (query.inStock) params.set("inStock", "true");
  if (query.minRating) params.set("minRating", String(query.minRating));
  return params;
};

export const totalPages = ({ pageSize, total }: Pagination): number =>
  Math.max(FIRST_PAGE, Math.ceil(total / pageSize));
