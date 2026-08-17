"use client";

import { type ComponentProps, type ComponentType, useCallback, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";

import { useLocale, useTranslations } from "next-intl";

import { REQUEST_STATUS, type RequestStatus } from "@/shared/constants/requestStatus";

import type { Product } from "@/shared/types/catalog.types";

import ErrorView from "@/shared/ui/ErrorView";
import LoadingView from "@/shared/ui/LoadingView";
import NullView from "@/shared/ui/NullView";

import type { CatalogQuery, Facets } from "@/features/web/catalog/types/catalogQuery.types";
import { PRODUCT_SORT } from "@/features/web/catalog/types/catalogQuery.types";

import { fetchFacets, searchProducts } from "@/features/web/catalog/api/products";

import { totalPages } from "@/features/web/catalog/helpers/catalogQuery";

import CatalogFilters from "@/features/web/catalog/components/CatalogFilters";
import ProductCard from "@/features/web/catalog/components/ProductCard";
import SortControl from "@/features/web/catalog/components/SortControl";

import styles from "@/features/web/catalog/styles/search.module.scss";

type StateViewProps = ComponentProps<typeof LoadingView> &
  ComponentProps<typeof ErrorView> &
  ComponentProps<typeof NullView>;

// Component references, never elements: an element map builds every branch on every render.
const STATE_VIEWS: Record<RequestStatus, ComponentType<StateViewProps>> = {
  [REQUEST_STATUS.IDLE]: LoadingView,
  [REQUEST_STATUS.LOADING]: LoadingView,
  [REQUEST_STATUS.ERROR]: ErrorView,
  [REQUEST_STATUS.SUCCESS]: NullView,
};

const EMPTY_FACETS: Facets = { categories: [], minPriceCents: 0, maxPriceCents: 0 };
const FIRST_PAGE = 1;

const SearchContent = () => {
  const t = useTranslations("search");
  const tCatalog = useTranslations("catalog");
  const locale = useLocale();
  const params = useSearchParams();

  const [query, setQuery] = useState<CatalogQuery>({
    q: params.get("q") ?? "",
    sort: PRODUCT_SORT.NEWEST,
    page: FIRST_PAGE,
  });
  const [products, setProducts] = useState<Product[]>([]);
  const [facets, setFacets] = useState<Facets>(EMPTY_FACETS);
  const [total, setTotal] = useState(0);
  const [pageSize, setPageSize] = useState(0);
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);

  const loadProducts = useCallback(async () => {
    setStatus(REQUEST_STATUS.LOADING);
    try {
      const [page, nextFacets] = await Promise.all([
        searchProducts(query),
        fetchFacets(query),
      ]);
      setProducts(page.items);
      setTotal(page.pagination.total);
      setPageSize(page.pagination.pageSize);
      setFacets(nextFacets);
      setStatus(REQUEST_STATUS.SUCCESS);
    } catch {
      setProducts([]);
      setStatus(REQUEST_STATUS.ERROR);
    }
  }, [query]);

  // The query object is the single source of truth; every control patches it and the effect refetches.
  useEffect(() => {
    loadProducts();
  }, [loadProducts]);

  const patchQuery = (patch: Partial<CatalogQuery>) => setQuery((prev) => ({ ...prev, ...patch }));

  const CurrentView = STATE_VIEWS[status];
  const isSuccess = status === REQUEST_STATUS.SUCCESS;
  const pageCount = pageSize ? totalPages({ page: query.page ?? FIRST_PAGE, pageSize, total }) : FIRST_PAGE;

  return (
    <>
      <div className={styles.backBar}>
        <a className={styles.back} href={`/${locale}`}>
          ←
        </a>
        <input
          className={styles.input}
          placeholder={t("placeholder")}
          value={query.q ?? ""}
          onChange={(e) => patchQuery({ q: e.target.value, page: FIRST_PAGE })}
        />
        <SortControl
          value={query.sort ?? PRODUCT_SORT.NEWEST}
          onChange={(sort) => patchQuery({ sort, page: FIRST_PAGE })}
        />
      </div>

      <div className={styles.layout}>
        <CatalogFilters facets={facets} query={query} onChange={patchQuery} />

        <div className={styles.results}>
          <CurrentView onRetry={loadProducts} />

          {isSuccess && products.length > 0 && (
            <>
              <p className={styles.hint}>{tCatalog("resultCount", { count: total })}</p>
              <div className={styles.grid}>
                {products.map((p) => (
                  <ProductCard key={p.id} product={p} locale={locale} />
                ))}
              </div>
              {pageCount > FIRST_PAGE && (
                <nav className={styles.pager} aria-label={tCatalog("filters")}>
                  {Array.from({ length: pageCount }, (_, i) => i + FIRST_PAGE).map((page) => (
                    <button
                      key={page}
                      type="button"
                      className="btn btnOutline btnSm"
                      aria-current={page === (query.page ?? FIRST_PAGE) ? "page" : undefined}
                      onClick={() => patchQuery({ page })}
                    >
                      {page}
                    </button>
                  ))}
                </nav>
              )}
            </>
          )}

          {isSuccess && products.length === 0 && (
            <p className={styles.hint}>{tCatalog("noResults")}</p>
          )}
        </div>
      </div>
    </>
  );
};

export default SearchContent;
