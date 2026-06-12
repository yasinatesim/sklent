"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { useCartStore } from "@/stores/cartStore";
import { useToastStore } from "@/stores/toastStore";
import styles from "./AddToCart.module.scss";

type AddToCartProps = {
  productId: string;
  slug: string;
  title: string;
  unitCents: number;
  oldUnitCents: number;
};

const AddToCart = ({ productId, slug, title, unitCents, oldUnitCents }: AddToCartProps) => {
  const t = useTranslations("common");
  const [quantity, setQuantity] = useState(1);
  const addItem = useCartStore((s) => s.addItem);
  const showToast = useToastStore((s) => s.show);

  const handleAddClick = () => {
    addItem({ productId, slug, title, unitCents, oldUnitCents, quantity });
    showToast(`${title} ${t("addedToCart")}`);
  };

  return (
    <div className={styles.actions}>
      <div className={styles.qty}>
        <button onClick={() => setQuantity((q) => Math.max(1, q - 1))} aria-label="-">
          −
        </button>
        <span>{quantity}</span>
        <button onClick={() => setQuantity((q) => q + 1)} aria-label="+">
          +
        </button>
      </div>
      <button className="btn btnPrimary" onClick={handleAddClick}>
        {t("addToCart")}
      </button>
    </div>
  );
};

export default AddToCart;
