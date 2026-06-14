import { getTranslations, setRequestLocale } from "next-intl/server";
import ProductCard from "@/components/ProductCard/ProductCard";
import CategoryCard from "@/components/CategoryCard/CategoryCard";
import BoxIcon from "@/components/icons/BoxIcon";
import TruckIcon from "@/components/icons/TruckIcon";
import SupportIcon from "@/components/icons/SupportIcon";
import { fetchProducts, fetchCategories, type Product, type Category } from "@/lib/api";
import styles from "./page.module.scss";

type HomePageProps = {
  params: Promise<{ locale: string }>;
};

const safeProducts = async (): Promise<Product[]> => {
  try {
    return await fetchProducts();
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

const bestSellers = (products: Product[]): Product[] => {
  const featured = products.filter((p) => p.badge === "Çok Satan" || p.badge.includes("İndirim"));
  return (featured.length > 0 ? featured : products).slice(0, 8);
};

const HomePage = async ({ params }: HomePageProps) => {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations("home");
  const products = await safeProducts();
  const categories = await safeCategories();

  return (
    <main>
      <section className={styles.hero}>
        <div className={`container ${styles.heroInner}`}>
          <div>
            <h1 className={styles.heroTitle}>{t.rich("heroTitle", { em: (c) => <em>{c}</em> })}</h1>
            <p className={styles.heroText}>{t("heroText")}</p>
            <div className={styles.heroButtons}>
              <a className="btn btnPrimary" href={`/${locale}/kategori/all`}>
                {t("heroCta")}
              </a>
              <a className="btn btnOutline" href={`/${locale}/kategori/all`}>
                {t("allProducts")}
              </a>
            </div>
          </div>
          <div className={styles.heroVisual}>
            <div className={styles.feature}>
              <BoxIcon size={48} />
              <span>{t("featureCollection")}</span>
            </div>
            <div className={styles.feature}>
              <TruckIcon />
              <span>{t("featureShipping")}</span>
            </div>
            <div className={styles.feature}>
              <SupportIcon />
              <span>{t("featureSupport")}</span>
            </div>
          </div>
        </div>
      </section>

      <section className={styles.section}>
        <div className="container">
          <div className={styles.label}>{t("categoriesLabel")}</div>
          <h2 className={styles.title}>{t("categoriesTitle")}</h2>
          <p className={styles.desc}>{t("categoriesDesc")}</p>
          <div className={styles.catGrid}>
            {categories.map((c) => (
              <CategoryCard key={c.id} category={c} locale={locale} />
            ))}
          </div>
        </div>
      </section>

      <section className={`${styles.section} ${styles.sectionSurface}`}>
        <div className="container">
          <div className={styles.label}>{t("bestLabel")}</div>
          <h2 className={styles.title}>{t("bestTitle")}</h2>
          <p className={styles.desc}>{t("bestDesc")}</p>
          <div className={styles.prodGrid}>
            {bestSellers(products).map((p) => (
              <ProductCard key={p.id} product={p} locale={locale} />
            ))}
          </div>
        </div>
      </section>
    </main>
  );
};

export default HomePage;
