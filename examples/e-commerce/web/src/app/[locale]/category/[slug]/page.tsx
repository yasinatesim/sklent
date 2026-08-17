import { getTranslations, setRequestLocale } from "next-intl/server";

import type { Category, Product } from "@/shared/types/catalog.types";

import { fetchCategories, fetchProducts } from "@/features/web/catalog/api/products";

import CategoryCard from "@/features/web/catalog/components/CategoryCard";
import ProductCard from "@/features/web/catalog/components/ProductCard";

import styles from "./page.module.scss";

type CategoryPageProps = {
  params: Promise<{ locale: string; slug: string }>;
};

type CategoryTitleInput = {
  slug: string;
  current?: { nameEn: string; nameTr: string };
  locale: string;
  t: (key: string) => string;
};

const categoryTitle = ({ slug, current, locale, t }: CategoryTitleInput): string => {
  if (slug === "all") return t("allProducts");
  if (!current) return t("category");
  return locale === "en" ? current.nameEn : current.nameTr;
};

const safeProducts = async (slug: string): Promise<Product[]> => {
  try {
    return await fetchProducts(slug);
  } catch {
    return [];
  }
};

const safeCategories = async (): Promise<Category[]> => {
  try {
    return await fetchCategories();
  } catch {
    return [];
  }
};

const CategoryPage = async ({ params }: CategoryPageProps) => {
  const { locale, slug } = await params;
  setRequestLocale(locale);
  const t = await getTranslations("common");
  const products = await safeProducts(slug);
  const categories = await safeCategories();
  const current = categories.find((c) => c.slug === slug);
  const title = categoryTitle({ slug, current, locale, t });

  return (
    <main className={`container ${styles.page}`}>
      <div className={styles.backBar}>
        <a className={styles.back} href={`/${locale}`}>
          ← {t("home")}
        </a>
        <h1 className={styles.title}>{title}</h1>
      </div>

      {slug !== "all" ? (
        <div className={styles.subCats}>
          {categories
            .filter((c) => c.slug !== slug)
            .map((c) => (
              <CategoryCard key={c.id} category={c} locale={locale} showDesc={false} />
            ))}
        </div>
      ) : null}

      {products.length === 0 ? (
        <p className={styles.empty}>{t("categoryEmpty")}</p>
      ) : (
        <div className={styles.grid}>
          {products.map((p) => (
            <ProductCard key={p.id} product={p} locale={locale} />
          ))}
        </div>
      )}
    </main>
  );
};

export default CategoryPage;
