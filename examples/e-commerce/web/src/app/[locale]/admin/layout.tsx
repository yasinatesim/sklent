"use client";

import { useLocale, useTranslations } from "next-intl";
import type { ReactNode } from "react";

import { localePath,ROUTES } from "@/shared/constants/routes";

import { useAuthStore } from "@/features/web/auth/store/authStore";

import AdminSidebar from "@/features/admin/_shared/components/AdminSidebar";

import styles from "@/features/admin/_shared/styles/admin.module.scss";

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
        <a className="btn btnPrimary" href={localePath(locale, ROUTES.LOGIN)}>
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
