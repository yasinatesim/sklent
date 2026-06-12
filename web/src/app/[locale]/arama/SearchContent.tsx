"use client";

import { useEffect, useMemo, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import { useSearchParams } from "next/navigation";
import ProductCard from "@/components/ProductCard/ProductCard";
import { fetchProducts, type Product } from "@/lib/api";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import styles from "./page.module.scss";

const SearchContent = () => {
  const t = useTranslations("search");
  const locale = useLocale();
  const params = useSearchParams();
  const [query, setQuery] = useState(params.get("q") ?? "");
  const [products, setProducts] = useState<Product[]>([]);
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);

  useEffect(() => {
    const loadProducts = async () => {
      setStatus(REQUEST_STATUS.LOADING);
      try {
        setProducts(await fetchProducts());
        setStatus(REQUEST_STATUS.SUCCESS);
      } catch {
        setStatus(REQUEST_STATUS.ERROR);
      }
    };
    loadProducts();
  }, []);

  const results = useMemo(() => {
    const q = query.toLowerCase().trim();
    if (!q) return [];
    return products.filter(
      (p) => p.titleTr.toLowerCase().includes(q) || p.descriptionTr.toLowerCase().includes(q),
    );
  }, [products, query]);

  const showResults = query.trim().length > 0 && status === REQUEST_STATUS.SUCCESS;

  return (
    <>
      <div className={styles.backBar}>
        <a className={styles.back} href={`/${locale}`}>
          ←
        </a>
        <input
          className={styles.input}
          placeholder={t("placeholder")}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </div>
      {showResults && results.length > 0 ? (
        <div className={styles.grid}>
          {results.map((p) => (
            <ProductCard key={p.id} product={p} locale={locale} />
          ))}
        </div>
      ) : (
        <p className={styles.hint}>{query.trim() ? t("noResults", { q: query }) : t("startTyping")}</p>
      )}
    </>
  );
};

export default SearchContent;
