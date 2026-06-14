import { type Product, formatTRY } from "@/lib/api";
import BoxIcon from "@/components/icons/BoxIcon";
import styles from "./ProductCard.module.scss";

type ProductCardProps = {
  product: Product;
  locale: string;
};

const ProductCard = ({ product, locale }: ProductCardProps) => {
  const title = locale === "en" && product.titleEn ? product.titleEn : product.titleTr;
  const hasOld = product.oldPriceCents > product.priceCents;

  return (
    <a className={styles.card} href={`/${locale}/urun/${product.slug}`}>
      <div className={styles.img}>
        {product.badge ? <span className={styles.badge}>{product.badge}</span> : null}
        <BoxIcon size={48} />
      </div>
      <div className={styles.info}>
        <div className={styles.title}>{title}</div>
        <div className={styles.price}>
          {hasOld ? <span className={styles.old}>{formatTRY(product.oldPriceCents, locale)}</span> : null}
          {formatTRY(product.priceCents, locale)}
        </div>
      </div>
    </a>
  );
};

export default ProductCard;
