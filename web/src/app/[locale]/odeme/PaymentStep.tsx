"use client";

import { useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useCartStore } from "@/stores/cartStore";
import { useToastStore } from "@/stores/toastStore";
import { placeOrder } from "@/lib/api";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import CardPreview from "@/components/checkout/CardPreview";
import styles from "./odeme.module.scss";

type PaymentStepProps = {
  email: string;
};

const formatCardNumber = (value: string): string =>
  value
    .replace(/\D/g, "")
    .slice(0, 16)
    .replace(/(.{4})/g, "$1 ")
    .trim();

const formatExp = (value: string): string => {
  const digits = value.replace(/\D/g, "").slice(0, 4);
  return digits.length > 2 ? `${digits.slice(0, 2)}/${digits.slice(2)}` : digits;
};

const PaymentStep = ({ email }: PaymentStepProps) => {
  const t = useTranslations("checkout");
  const locale = useLocale();
  const router = useRouter();
  const items = useCartStore((s) => s.items);
  const total = useCartStore((s) => s.totalCents());
  const clear = useCartStore((s) => s.clear);
  const showToast = useToastStore((s) => s.show);

  const [name, setName] = useState("");
  const [number, setNumber] = useState("");
  const [exp, setExp] = useState("");
  const [cvc, setCvc] = useState("");
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);

  const handlePayClick = async () => {
    if (!name || number.replace(/\s/g, "").length < 16 || !exp || cvc.length < 3) {
      showToast(t("fillCard"));
      return;
    }
    setStatus(REQUEST_STATUS.LOADING);
    const res = await placeOrder(
      email,
      items.map((i) => ({ productId: i.productId, titleTr: i.title, unitCents: i.unitCents, quantity: i.quantity })),
    );
    if (!res.ok) {
      setStatus(REQUEST_STATUS.ERROR);
      showToast(t("payError"));
      return;
    }
    const data = await res.json();
    clear();
    router.push(`/${locale}/odeme/basarili?order=${data.orderId}&amount=${total}`);
  };

  return (
    <div className={styles.narrow}>
      <CardPreview number={number} name={name} exp={exp} />
      <div className={styles.group}>
        <label>{t("cardName")}</label>
        <input value={name} onChange={(e) => setName(e.target.value)} />
      </div>
      <div className={styles.group}>
        <label>{t("cardNumber")}</label>
        <input
          value={number}
          maxLength={19}
          placeholder="0000 0000 0000 0000"
          onChange={(e) => setNumber(formatCardNumber(e.target.value))}
        />
      </div>
      <div className={styles.row3}>
        <div className={styles.group}>
          <label>{t("expiry")}</label>
          <input value={exp} maxLength={5} placeholder="AA/YY" onChange={(e) => setExp(formatExp(e.target.value))} />
        </div>
        <div className={styles.group}>
          <label>CVC</label>
          <input value={cvc} maxLength={3} placeholder="000" onChange={(e) => setCvc(e.target.value.replace(/\D/g, ""))} />
        </div>
        <div className={styles.group}>
          <label>&nbsp;</label>
          <button
            className={`btn btnPrimary ${styles.full}`}
            onClick={handlePayClick}
            disabled={status === REQUEST_STATUS.LOADING}
          >
            {t("pay")}
          </button>
        </div>
      </div>
    </div>
  );
};

export default PaymentStep;
