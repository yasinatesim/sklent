"use client";

import { useLocale, useTranslations } from "next-intl";
import { usePathname } from "next/navigation";
import type { Category } from "@/lib/api";
import styles from "./CategoryBar.module.scss";

type CategoryBarProps = {
  categories: Category[];
};

const CategoryBar = ({ categories }: CategoryBarProps) => {
  const t = useTranslations("common");
  const locale = useLocale();
  const pathname = usePathname();

  const isActive = (slug: string): boolean => {
    if (slug === "all") return pathname === `/${locale}` || pathname === `/${locale}/kategori/all`;
    return pathname === `/${locale}/kategori/${slug}`;
  };

  const label = (c: Category): string => (locale === "en" ? c.nameEn : c.nameTr);

  if (pathname.includes("/admin")) {
    return null;
  }

  return (
    <nav className={styles.bar}>
      <div className={`container ${styles.inner}`}>
        <a
          className={`${styles.link} ${isActive("all") ? styles.active : ""}`}
          href={`/${locale}/kategori/all`}
        >
          {t("allProducts")}
        </a>
        {categories.map((c) => (
          <a
            key={c.id}
            className={`${styles.link} ${isActive(c.slug) ? styles.active : ""}`}
            href={`/${locale}/kategori/${c.slug}`}
          >
            {label(c)}
          </a>
        ))}
      </div>
    </nav>
  );
};

export default CategoryBar;
