import { notFound } from "next/navigation";
import { getTranslations, setRequestLocale } from "next-intl/server";
import AddToCart from "@/components/AddToCart/AddToCart";
import BoxIcon from "@/components/icons/BoxIcon";
import { fetchProduct, formatTRY } from "@/lib/api";
import styles from "./page.module.scss";

type ProductPageProps = {
  params: Promise<{ locale: string; slug: string }>;
};

const ProductPage = async ({ params }: ProductPageProps) => {
  const { locale, slug } = await params;
  setRequestLocale(locale);
  const t = await getTranslations("product");
  const product = await fetchProduct(slug);

  if (!product) {
    notFound();
  }

  const title = locale === "en" && product.titleEn ? product.titleEn : product.titleTr;
  const hasOld = product.oldPriceCents > product.priceCents;

  return (
    <main className={`container ${styles.page}`}>
      <div className={styles.backBar}>
        <a className={styles.back} href={`/${locale}/kategori/all`}>
          ← {t("back")}
        </a>
      </div>
      <div className={styles.layout}>
        <div className={styles.image}>
          <BoxIcon size={120} />
        </div>
        <div>
          <h1 className={styles.title}>{title}</h1>
          <div className={styles.price}>
            {hasOld ? <span className={styles.old}>{formatTRY(product.oldPriceCents, locale)}</span> : null}
            {formatTRY(product.priceCents, locale)}
          </div>
          <p className={styles.desc}>{product.descriptionTr}</p>
          <div className={styles.meta}>
            <div>
              <strong>{t("stock")}:</strong> {product.stock}
            </div>
            <div>
              <strong>{t("seller")}:</strong> {product.seller || "Vela Commerce"}
            </div>
          </div>
          <AddToCart
            productId={product.id}
            slug={product.slug}
            title={title}
            unitCents={product.priceCents}
            oldUnitCents={product.oldPriceCents}
          />
        </div>
      </div>
    </main>
  );
};

export default ProductPage;
