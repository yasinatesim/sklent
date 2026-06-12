"use client";

import { useSearchParams } from "next/navigation";
import { useLocale, useTranslations } from "next-intl";
import { formatTRY } from "@/lib/api";
import CheckCircleIcon from "@/components/icons/CheckCircleIcon";
import styles from "./page.module.scss";

const SuccessContent = () => {
  const t = useTranslations("success");
  const locale = useLocale();
  const params = useSearchParams();
  const orderId = params.get("order") ?? "";
  const amount = Number(params.get("amount") ?? 0);

  return (
    <div className={styles.box}>
      <CheckCircleIcon size={72} />
      <h1 className={styles.title}>{t("title")}</h1>
      <p className={styles.text}>{t("orderNo")}</p>
      <div className={styles.orderNo}>{orderId || "—"}</div>
      <p className={styles.amount}>
        {t("amount")}: <strong>{formatTRY(amount, locale)}</strong>
      </p>
      <a className="btn btnPrimary" href={`/${locale}`}>
        {t("backHome")}
      </a>
    </div>
  );
};

export default SuccessContent;
