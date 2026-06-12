"use client";

import { useLocale, useTranslations } from "next-intl";
import { usePathname, useRouter } from "next/navigation";
import { API_BASE } from "@/lib/api";
import { useAuthStore } from "@/stores/authStore";
import { useToastStore } from "@/stores/toastStore";
import styles from "./admin.module.scss";

const AdminSidebar = () => {
  const t = useTranslations("admin");
  const locale = useLocale();
  const pathname = usePathname();
  const router = useRouter();
  const setUser = useAuthStore((s) => s.setUser);
  const showToast = useToastStore((s) => s.show);

  const items = [
    { href: `/${locale}/admin`, label: `📊 ${t("dashboard")}`, exact: true },
    { href: `/${locale}/admin/urunler`, label: `📦 ${t("products")}`, exact: false },
    { href: `/${locale}/admin/kampanyalar`, label: `🏷️ ${t("campaigns")}`, exact: false },
    { href: `/${locale}/admin/kuponlar`, label: `🎫 ${t("coupons")}`, exact: false },
    { href: `/${locale}/admin/siparisler`, label: `📋 ${t("orders")}`, exact: false },
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
