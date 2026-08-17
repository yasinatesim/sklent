"use client";

import { usePathname, useRouter } from "next/navigation";

import { useLocale, useTranslations } from "next-intl";

import { ADMIN_ROUTES, localePath,ROUTES } from "@/shared/constants/routes";

import { API_BASE } from "@/shared/helpers/api/client";
import { useToastStore } from "@/shared/stores/toastStore";

import { useAuthStore } from "@/features/web/auth/store/authStore";

import styles from "@/features/admin/_shared/styles/admin.module.scss";

const AdminSidebar = () => {
  const t = useTranslations("admin");
  const tReturn = useTranslations("returns");
  const locale = useLocale();
  const pathname = usePathname();
  const router = useRouter();
  const setUser = useAuthStore((s) => s.setUser);
  const showToast = useToastStore((s) => s.show);

  const items = [
    { href: localePath(locale, ROUTES.ADMIN), label: `📊 ${t("dashboard")}`, exact: true },
    { href: localePath(locale, ADMIN_ROUTES.PRODUCTS), label: `📦 ${t("products")}`, exact: false },
    { href: localePath(locale, ADMIN_ROUTES.PROMOTIONS), label: `🏷️ ${t("campaigns")}`, exact: false },
    { href: localePath(locale, ADMIN_ROUTES.COUPONS), label: `🎫 ${t("coupons")}`, exact: false },
    { href: localePath(locale, ADMIN_ROUTES.ORDERS), label: `📋 ${t("orders")}`, exact: false },
    { href: localePath(locale, ADMIN_ROUTES.STOCK_TRACKING), label: `📈 ${t("stockTracking")}`, exact: false },
    { href: localePath(locale, ADMIN_ROUTES.REVIEWS), label: `💬 ${t("reviews")}`, exact: false },
    { href: localePath(locale, ADMIN_ROUTES.RETURNS), label: `↩️ ${tReturn("title")}`, exact: false },
    { href: localePath(locale, ADMIN_ROUTES.SETTINGS), label: `⚙️ ${t("settings")}`, exact: false },
  ];

  const isActive = (href: string, exact: boolean): boolean =>
    exact ? pathname === href : pathname.startsWith(href);

  const handleLogoutClick = async () => {
    await fetch(`${API_BASE}/auth/logout`, { method: "POST", credentials: "include" });
    setUser(null);
    showToast(t("loggedOut"));
    router.push(`/${locale}`);
  };

  return (
    <nav className={styles.sidebar}>
      {items.map((item) => (
        <a
          key={item.href}
          className={isActive(item.href, item.exact) ? styles.active : ""}
          href={item.href}
        >
          {item.label}
        </a>
      ))}
      <a className={styles.logout} onClick={handleLogoutClick}>
        🚪 {t("logout")}
      </a>
    </nav>
  );
};

export default AdminSidebar;
