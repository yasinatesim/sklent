"use client";

import { useTranslations } from "next-intl";

import type { CatalogQuery, Facets } from "@/features/web/catalog/types/catalogQuery.types";

import styles from "./CatalogFilters.module.scss";

type Props = {
  facets: Facets;
  query: CatalogQuery;
  onChange: (patch: Partial<CatalogQuery>) => void;
};

const CENTS_IN_LIRA = 100;
const FIRST_PAGE = 1;

const toLira = (cents: number): string => String(Math.round(cents / CENTS_IN_LIRA));

const CatalogFilters = ({ facets, query, onChange }: Props) => {
  const t = useTranslations("catalog");

  // Any filter change invalidates the current page: page 3 of the previous result set does not
  // exist in the new one.
  const patch = (next: Partial<CatalogQuery>) => onChange({ ...next, page: FIRST_PAGE });

  return (
    <aside className={styles.root} aria-label={t("filters")}>
      <section className={styles.group}>
        <h3 className={styles.title}>{t("categories")}</h3>
        <ul className={styles.list}>
          {facets.categories.map((facet) => (
            <li key={facet.slug}>
              <button
                type="button"
                className={facet.slug === query.category ? styles.rowActive : styles.row}
                aria-pressed={facet.slug === query.category}
                onClick={() =>
                  patch({ category: facet.slug === query.category ? undefined : facet.slug })
                }
              >
                <span>{facet.name}</span>
                <span className={styles.count}>{facet.count}</span>
              </button>
            </li>
          ))}
        </ul>
      </section>

      <section className={styles.group}>
        <h3 className={styles.title}>{t("price")}</h3>
        <div className={styles.priceRow}>
          <input
            type="number"
            className={styles.price}
            aria-label={t("minPrice")}
            placeholder={toLira(facets.minPriceCents)}
            value={query.minPrice ? toLira(query.minPrice) : ""}
            onChange={(e) =>
              patch({ minPrice: e.target.value ? Number(e.target.value) * CENTS_IN_LIRA : undefined })
            }
          />
          <span aria-hidden="true">–</span>
          <input
            type="number"
            className={styles.price}
            aria-label={t("maxPrice")}
            placeholder={toLira(facets.maxPriceCents)}
            value={query.maxPrice ? toLira(query.maxPrice) : ""}
            onChange={(e) =>
              patch({ maxPrice: e.target.value ? Number(e.target.value) * CENTS_IN_LIRA : undefined })
            }
          />
        </div>
      </section>

      <label className={styles.checkbox}>
        <input
          type="checkbox"
          checked={Boolean(query.inStock)}
          onChange={(e) => patch({ inStock: e.target.checked })}
        />
        {t("inStockOnly")}
      </label>
    </aside>
  );
};

export default CatalogFilters;
