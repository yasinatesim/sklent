"use client";

import { useLocale, useTranslations } from "next-intl";
import { useCartStore } from "@/stores/cartStore";
import { formatTRY } from "@/lib/api";
import styles from "./OrderSummary.module.scss";

const OrderSummary = () => {
  const locale = useLocale();
  const t = useTranslations("checkout");
  const items = useCartStore((s) => s.items);
  const total = useCartStore((s) => s.totalCents());

  return (
    <div className={styles.summary}>
      <div className={styles.title}>{t("summary")}</div>
      {items.map((item) => (
        <div key={item.productId} className={styles.item}>
          <span>
            {item.title} x{item.quantity}
          </span>
          <span>{formatTRY(item.unitCents * item.quantity, locale)}</span>
        </div>
      ))}
      <div className={styles.total}>
        <span>{t("total")}</span>
        <span>{formatTRY(total, locale)}</span>
      </div>
    </div>
  );
};

export default OrderSummary;
