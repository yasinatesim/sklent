"use client";

import type { ReactNode } from "react";
import { useLocale, useTranslations } from "next-intl";
import { useAuthStore } from "@/stores/authStore";
import AdminSidebar from "./AdminSidebar";
import styles from "./admin.module.scss";

type AdminLayoutProps = {
  children: ReactNode;
};

const AdminLayout = ({ children }: AdminLayoutProps) => {
  const t = useTranslations("admin");
  const locale = useLocale();
  const isAdmin = useAuthStore((s) => s.isAdmin());

  if (!isAdmin) {
    return (
      <div className={styles.guard}>
        <h2>{t("guardTitle")}</h2>
        <p>{t("guardDesc")}</p>
        <a className="btn btnPrimary" href={`/${locale}/giris`}>
          {t("goLogin")}
        </a>
      </div>
    );
  }

  return (
    <div className={styles.layout}>
      <AdminSidebar />
      <div className={styles.content}>{children}</div>
    </div>
  );
};

export default AdminLayout;
