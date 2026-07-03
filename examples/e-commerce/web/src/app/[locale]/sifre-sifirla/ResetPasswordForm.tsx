"use client";

import { type FormEvent, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import { useRouter, useSearchParams } from "next/navigation";
import { resetPassword } from "@/lib/auth-extra-api";
import { useToastStore } from "@/stores/toastStore";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import styles from "../giris/page.module.scss";

const ResetPasswordForm = () => {
  const t = useTranslations("auth");
  const locale = useLocale();
  const router = useRouter();
  const searchParams = useSearchParams();
  const showToast = useToastStore((s) => s.show);
  const token = searchParams.get("token") ?? "";
  const [newPassword, setNewPassword] = useState("");
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setStatus(REQUEST_STATUS.LOADING);
    const res = await resetPassword(token, newPassword);
    if (!res.ok) {
      setStatus(REQUEST_STATUS.ERROR);
      showToast(t("resetPasswordFailed"));
      return;
    }
    setStatus(REQUEST_STATUS.SUCCESS);
    showToast(t("resetPasswordSuccess"));
    router.push(`/${locale}/giris`);
  };

  return (
    <div className={styles.box}>
      <h1 className={styles.title}>{t("resetPasswordTitle")}</h1>
      <form onSubmit={handleSubmit}>
        <div className={styles.group}>
          <label>{t("newPassword")}</label>
          <input
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            required
            minLength={8}
          />
        </div>
        <button className={`btn btnPrimary ${styles.full}`} type="submit" disabled={status === REQUEST_STATUS.LOADING}>
          {t("resetPasswordSubmit")}
        </button>
      </form>
    </div>
  );
};

export default ResetPasswordForm;
