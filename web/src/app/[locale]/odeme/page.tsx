"use client";

import { useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import { useCartStore } from "@/stores/cartStore";
import OrderSummary from "@/components/checkout/OrderSummary";
import AddressStep, { type Address } from "./AddressStep";
import PaymentStep from "./PaymentStep";
import styles from "./odeme.module.scss";

const STEP = { ADDRESS: "address", PAYMENT: "payment" } as const;
type Step = (typeof STEP)[keyof typeof STEP];

const CheckoutPage = () => {
  const locale = useLocale();
  const t = useTranslations("checkout");
  const items = useCartStore((s) => s.items);
  const [step, setStep] = useState<Step>(STEP.ADDRESS);
  const [address, setAddress] = useState<Address | null>(null);

  const goToPayment = (next: Address) => {
    setAddress(next);
    setStep(STEP.PAYMENT);
  };

  const steps = {
    [STEP.ADDRESS]: <AddressStep onContinue={goToPayment} />,
    [STEP.PAYMENT]: <PaymentStep email={address?.email ?? ""} />,
  };

  if (items.length === 0) {
    return (
      <main className={`container ${styles.page}`}>
        <h1 className={styles.title}>{t("empty")}</h1>
        <a className="btn btnPrimary" href={`/${locale}/kategori/all`}>
          {t("continue")}
        </a>
      </main>
    );
  }

  return (
    <main className={`container ${styles.page}`}>
      <div className={styles.backBar}>
        <a className={styles.back} href={`/${locale}/sepet`}>
          ← {t("backToCart")}
        </a>
        <h1 className={styles.title}>{t("title")}</h1>
      </div>
      <div className={styles.layout}>
        <div>{steps[step]}</div>
        <OrderSummary />
      </div>
    </main>
  );
};

export default CheckoutPage;
