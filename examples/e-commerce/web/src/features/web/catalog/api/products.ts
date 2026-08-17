import type { Category, Product } from "@/shared/types/catalog.types";

import { API_BASE, INTERNAL_API_BASE } from "@/shared/helpers/api/client";

import type {
  CatalogQuery,
  Facets,
  Pagination,
  ProductListPage,
} from "@/features/web/catalog/types/catalogQuery.types";

import { buildCatalogSearchParams } from "@/features/web/catalog/helpers/catalogQuery";

export const fetchProducts = async (categorySlug?: string): Promise<Product[]> => {
  const query = categorySlug && categorySlug !== "all" ? `?category=${categorySlug}` : "";
  const res = await fetch(`${INTERNAL_API_BASE}/products${query}`, { cache: "no-store" });
  if (!res.ok) throw new Error("failed to load products");
  const data = await res.json();
  return data.items ?? [];
};

export const fetchProduct = async (slug: string): Promise<Product | null> => {
  const res = await fetch(`${INTERNAL_API_BASE}/products/${slug}`, { cache: "no-store" });
  if (res.status === 404) return null;
  if (!res.ok) throw new Error("failed to load product");
  return res.json();
};

export const fetchCategories = async (): Promise<Category[]> => {
  const res = await fetch(`${INTERNAL_API_BASE}/categories`, { cache: "no-store" });
  if (!res.ok) return [];
  const data = await res.json();
  return data.items ?? [];
};

export const searchProducts = async (query: CatalogQuery): Promise<ProductListPage> => {
  const qs = buildCatalogSearchParams(query).toString();
  const res = await fetch(`${API_BASE}/products${qs ? `?${qs}` : ""}`, { cache: "no-store" });
  if (!res.ok) throw new Error("failed to load products");
  const data = (await res.json()) as { items: Product[] | null; pagination: Pagination };
  return { items: data.items ?? [], pagination: data.pagination };
};

export const fetchFacets = async (query: CatalogQuery): Promise<Facets> => {
  const qs = buildCatalogSearchParams(query).toString();
  const res = await fetch(`${API_BASE}/products/facets${qs ? `?${qs}` : ""}`, { cache: "no-store" });
  if (!res.ok) throw new Error("failed to load facets");
  return (await res.json()) as Facets;
};
