import { getTranslations, setRequestLocale } from "next-intl/server";
import ProductCard from "@/components/ProductCard/ProductCard";
import CategoryCard from "@/components/CategoryCard/CategoryCard";
import { fetchProducts, fetchCategories, type Product, type Category } from "@/lib/api";
import styles from "./page.module.scss";

type CategoryPageProps = {
  params: Promise<{ locale: string; slug: string }>;
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
  const title = slug === "all" ? t("allProducts") : current ? (locale === "en" ? current.nameEn : current.nameTr) : t("category");

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
