import type { Product } from "@/shared/types/catalog.types";

export const PRODUCT_SORT = {
  NEWEST: "newest",
  PRICE_ASC: "price_asc",
  PRICE_DESC: "price_desc",
} as const;

export type ProductSort = (typeof PRODUCT_SORT)[keyof typeof PRODUCT_SORT];

export type CatalogQuery = {
  category?: string;
  q?: string;
  page?: number;
  pageSize?: number;
  minPrice?: number;
  maxPrice?: number;
  sort?: ProductSort;
  inStock?: boolean;
  minRating?: number;
};

export type Pagination = {
  page: number;
  pageSize: number;
  total: number;
};

export type CategoryFacet = {
  slug: string;
  name: string;
  count: number;
};

export type Facets = {
  categories: CategoryFacet[];
  minPriceCents: number;
  maxPriceCents: number;
};

export type ProductListPage = {
  items: Product[];
  pagination: Pagination;
};
