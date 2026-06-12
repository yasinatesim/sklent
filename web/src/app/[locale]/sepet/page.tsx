"use client";

import { use } from "react";
import { useLocale, useTranslations } from "next-intl";
import { useCartStore } from "@/stores/cartStore";
import { formatTRY } from "@/lib/api";
import BoxIcon from "@/components/icons/BoxIcon";
import CartIcon from "@/components/icons/CartIcon";
import styles from "./page.module.scss";

type CartPageProps = {
  params: Promise<{ locale: string }>;
};

const CartPage = ({ params }: CartPageProps) => {
  use(params);
  const locale = useLocale();
  const t = useTranslations("cart");
  const tc = useTranslations("common");
  const items = useCartStore((s) => s.items);
  const setQuantity = useCartStore((s) => s.setQuantity);
  const removeItem = useCartStore((s) => s.removeItem);
  const total = useCartStore((s) => s.totalCents());
  const savings = useCartStore((s) => s.savingsCents());

  return (
    <main className={`container ${styles.page}`}>
      <div className={styles.backBar}>
        <a className={styles.back} href={`/${locale}/kategori/all`}>
          ← {t("continue")}
        </a>
        <h1 className={styles.title}>{t("title")}</h1>
      </div>

      {items.length === 0 ? (
        <div className={styles.empty}>
          <CartIcon />
          <h2>{t("empty")}</h2>
          <p>{t("emptyHint")}</p>
          <a className="btn btnPrimary" href={`/${locale}/kategori/all`}>
            {t("startShopping")}
          </a>
        </div>
      ) : (
        <>
          <ul className={styles.items}>
            {items.map((item) => (
              <li key={item.productId} className={styles.item}>
                <div className={styles.thumb}>
                  <BoxIcon size={32} />
                </div>
                <div>
                  <div className={styles.itemTitle}>{item.title}</div>
                  <div className={styles.itemPrice}>{formatTRY(item.unitCents * item.quantity, locale)}</div>
                </div>
                <div className={styles.itemActions}>
                  <button onClick={() => setQuantity(item.productId, item.quantity - 1)} aria-label="-">
                    −
                  </button>
                  <span className={styles.qtyValue}>{item.quantity}</span>
                  <button onClick={() => setQuantity(item.productId, item.quantity + 1)} aria-label="+">
                    +
                  </button>
                  <button className={styles.remove} onClick={() => removeItem(item.productId)}>
                    {tc("remove")}
                  </button>
                </div>
              </li>
            ))}
          </ul>
          <div className={styles.summary}>
            <div>
              {savings > 0 ? <div className={styles.savings}>{t("savings")}: {formatTRY(savings, locale)}</div> : null}
              <div className={styles.total}>
                <small>{t("total")}</small> {formatTRY(total, locale)}
              </div>
            </div>
            <a className="btn btnPrimary" href={`/${locale}/odeme`}>
              {t("checkout")}
            </a>
          </div>
        </>
      )}
    </main>
  );
};

export default CartPage;
