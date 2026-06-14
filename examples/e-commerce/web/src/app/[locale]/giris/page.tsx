"use client";

import { type FormEvent, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { login, register, API_BASE } from "@/lib/api";
import { useAuthStore, type SessionUser } from "@/stores/authStore";
import { useToastStore } from "@/stores/toastStore";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import styles from "./page.module.scss";

const MODE = { LOGIN: "login", REGISTER: "register" } as const;
type Mode = (typeof MODE)[keyof typeof MODE];

const AuthPage = () => {
  const t = useTranslations("auth");
  const locale = useLocale();
  const router = useRouter();
  const setUser = useAuthStore((s) => s.setUser);
  const showToast = useToastStore((s) => s.show);

  const [mode, setMode] = useState<Mode>(MODE.LOGIN);
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("admin@vela.test");
  const [password, setPassword] = useState("admin12345");
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);

  const loadSession = async (): Promise<SessionUser | null> => {
    const res = await fetch(`${API_BASE}/auth/me`, { credentials: "include" });
    if (!res.ok) return null;
    return res.json();
  };

  const handleAuthSubmit = async (event: FormEvent) => {
    event.preventDefault();
    setStatus(REQUEST_STATUS.LOADING);
    const res = mode === MODE.LOGIN ? await login(email, password) : await register(email, password, fullName);
    if (!res.ok) {
      setStatus(REQUEST_STATUS.ERROR);
      showToast(t("failed"));
      return;
    }
    const user = await loadSession();
    setUser(user);
    setStatus(REQUEST_STATUS.SUCCESS);
    showToast(t("welcome"));
    router.push(user?.role === "admin" ? `/${locale}/admin` : `/${locale}`);
  };

  const isRegister = mode === MODE.REGISTER;

  return (
    <main className="container">
      <div className={styles.box}>
        <h1 className={styles.title}>{isRegister ? t("registerTitle") : t("loginTitle")}</h1>
        <p className={styles.desc}>{isRegister ? t("registerDesc") : t("loginDesc")}</p>
        <form onSubmit={handleAuthSubmit}>
          {isRegister ? (
            <div className={styles.group}>
              <label>{t("fullName")}</label>
              <input value={fullName} onChange={(e) => setFullName(e.target.value)} />
            </div>
          ) : null}
          <div className={styles.group}>
            <label>{t("email")}</label>
            <input type="email" value={email} onChange={(e) => setEmail(e.target.value)} required />
          </div>
          <div className={styles.group}>
            <label>{t("password")}</label>
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} required />
          </div>
          <button className={`btn btnPrimary ${styles.full}`} type="submit" disabled={status === REQUEST_STATUS.LOADING}>
            {isRegister ? t("register") : t("login")}
          </button>
        </form>
        <div className={styles.switch}>
          {isRegister ? t("haveAccount") : t("noAccount")}{" "}
          <button onClick={() => setMode(isRegister ? MODE.LOGIN : MODE.REGISTER)}>
            {isRegister ? t("login") : t("register")}
          </button>
        </div>
        <p className={styles.hint}>{t("adminHint")}</p>
      </div>
    </main>
  );
};

export default AuthPage;
