"use client";

import { useTranslations } from "next-intl";

import { PRODUCT_SORT, type ProductSort } from "@/features/web/catalog/types/catalogQuery.types";

import styles from "./SortControl.module.scss";

type Props = {
  value: ProductSort;
  onChange: (sort: ProductSort) => void;
};

const SORT_OPTIONS: ProductSort[] = Object.values(PRODUCT_SORT);

const SortControl = ({ value, onChange }: Props) => {
  const t = useTranslations("catalog");

  return (
    <label className={styles.root}>
      <span className={styles.label}>{t("sortBy")}</span>
      <select
        className={styles.select}
        aria-label={t("sortBy")}
        value={value}
        onChange={(e) => onChange(e.target.value as ProductSort)}
      >
        {SORT_OPTIONS.map((option) => (
          <option key={option} value={option}>
            {t(`sort_${option}`)}
          </option>
        ))}
      </select>
    </label>
  );
};

export default SortControl;
