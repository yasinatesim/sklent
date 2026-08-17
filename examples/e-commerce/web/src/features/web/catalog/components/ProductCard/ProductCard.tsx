import { useTranslations } from "next-intl";

import { localePath,ROUTES } from "@/shared/constants/routes";

import type { Product } from "@/shared/types/catalog.types";

import { formatTRY } from "@/shared/helpers/money";
import { productBadgeLabel } from "@/shared/helpers/productBadgeLabel";

import BoxIcon from "@/shared/ui/icons/BoxIcon";

import styles from "./ProductCard.module.scss";

type ProductCardProps = {
  product: Product;
  locale: string;
};

const ProductCard = ({ product, locale }: ProductCardProps) => {
  const tBadge = useTranslations("productBadge");
  const badgeLabel = productBadgeLabel(product.badge, tBadge);
  const title = locale === "en" && product.titleEn ? product.titleEn : product.titleTr;
  const hasOld = product.oldPriceCents > product.priceCents;

  return (
    <a className={styles.card} href={localePath(locale, `${ROUTES.PRODUCT}/${product.slug}`)}>
      <div className={styles.img}>
        {badgeLabel && <span className={styles.badge}>{badgeLabel}</span>}
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
