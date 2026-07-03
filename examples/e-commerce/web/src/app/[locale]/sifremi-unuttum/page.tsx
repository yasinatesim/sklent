"use client";

import { type FormEvent, useState } from "react";
import { useTranslations } from "next-intl";
import { forgotPassword } from "@/lib/auth-extra-api";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import styles from "../giris/page.module.scss";

const ForgotPasswordPage = () => {
  const t = useTranslations("auth");
  const [email, setEmail] = useState("");
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setStatus(REQUEST_STATUS.LOADING);
    await forgotPassword(email);
    setStatus(REQUEST_STATUS.SUCCESS);
  };

  return (
    <main className="container">
      <div className={styles.box}>
        <h1 className={styles.title}>{t("forgotPasswordTitle")}</h1>
        <p className={styles.desc}>{t("forgotPasswordDesc")}</p>
        {status === REQUEST_STATUS.SUCCESS ? (
          <p className={styles.hint}>{t("forgotPasswordSent")}</p>
        ) : (
          <form onSubmit={handleSubmit}>
            <div className={styles.group}>
              <label>{t("email")}</label>
              <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
            </div>
            <button className={`btn btnPrimary ${styles.full}`} type="submit" disabled={status === REQUEST_STATUS.LOADING}>
              {t("forgotPasswordSubmit")}
            </button>
          </form>
        )}
      </div>
    </main>
  );
};

export default ForgotPasswordPage;
