"use client";

import { useTranslations } from "next-intl";
import { useModalStore, MODAL } from "@/stores/modalStore";
import { useToastStore } from "@/stores/toastStore";
import styles from "./CouponModal.module.scss";

const COUPON_CODE = "VELA15";

const CouponModal = () => {
  const t = useTranslations("coupon");
  const active = useModalStore((s) => s.active);
  const close = useModalStore((s) => s.close);
  const showToast = useToastStore((s) => s.show);

  if (active !== MODAL.COUPON) {
    return null;
  }

  const handleCopyClick = async () => {
    try {
      await navigator.clipboard.writeText(COUPON_CODE);
      showToast(`${t("copied")}: ${COUPON_CODE}`);
    } catch {
      showToast(`${t("code")}: ${COUPON_CODE}`);
    }
    close();
  };

  return (
    <div className={styles.overlay} onClick={close}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <button className={styles.close} onClick={close} aria-label={t("close")}>
          ✕
        </button>
        <h2 className={styles.title}>🎉 {t("title")}</h2>
        <p className={styles.desc}>{t("desc")}</p>
        <div className={styles.box}>
          <div className={styles.label}>{t("yourCode")}</div>
          <div className={styles.code}>{COUPON_CODE}</div>
          <div className={styles.hint}>{t("hint")}</div>
        </div>
        <div className={styles.actions}>
          <button className={`btn btnPrimary ${styles.grow}`} onClick={handleCopyClick}>
            {t("copy")}
          </button>
          <button className="btn btnOutline" onClick={close}>
            {t("close")}
          </button>
        </div>
      </div>
    </div>
  );
};

export default CouponModal;
