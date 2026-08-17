"use client";

import { type FormEvent, useState } from "react";
import { useRouter } from "next/navigation";

import { useLocale, useTranslations } from "next-intl";

import { MODAL } from "@/shared/constants/modal";
import { localePath,ROUTES } from "@/shared/constants/routes";
import { THEME } from "@/shared/constants/theme";

import { useModalStore } from "@/shared/stores/modalStore";
import { useThemeStore } from "@/shared/stores/themeStore";

import CartIcon from "@/shared/ui/icons/CartIcon";
import SearchIcon from "@/shared/ui/icons/SearchIcon";
import UserIcon from "@/shared/ui/icons/UserIcon";

import { useAuthStore } from "@/features/web/auth/store/authStore";
import { useCartStore } from "@/features/web/cart/store/cartStore";

import styles from "./Header.module.scss";

const Header = () => {
  const t = useTranslations("common");
  const locale = useLocale();
  const router = useRouter();
  const count = useCartStore((s) => s.count());
  const theme = useThemeStore((s) => s.theme);
  const toggle = useThemeStore((s) => s.toggle);
  const openModal = useModalStore((s) => s.open);
  const user = useAuthStore((s) => s.user);
  const [query, setQuery] = useState("");

  const handleSearchSubmit = (event: FormEvent) => {
    event.preventDefault();
    router.push(localePath(locale, ROUTES.SEARCH) + `?q=${encodeURIComponent(query)}`);
  };

  return (
    <header className={styles.header}>
      <div className={`container ${styles.inner}`}>
        <a className={styles.logo} href={`/${locale}`}>
          Vela<span> Commerce</span>
        </a>

        <form className={styles.center} onSubmit={handleSearchSubmit}>
          <input
            type="text"
            placeholder={t("searchPlaceholder")}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          <span className={styles.searchIcon}>
            <SearchIcon />
          </span>
        </form>

        <div className={styles.actions}>
          <button className={styles.emoji} onClick={() => openModal(MODAL.COUPON)} title={t("campaign")}>
            🎁
          </button>
          <button className={styles.emoji} onClick={toggle} title={t("theme")}>
            {theme === THEME.DARK ? "☀️" : "🌙"}
          </button>
          <button onClick={() => router.push(localePath(locale, ROUTES.LOGIN))}>
            <UserIcon />
            <span className={styles.label}>{user ? user.email : t("signIn")}</span>
          </button>
          <button onClick={() => router.push(localePath(locale, ROUTES.CART))}>
            <CartIcon />
            <span className={styles.label}>{t("cart")}</span>
            {count > 0 ? <span className={styles.badge}>{count}</span> : null}
          </button>
        </div>
      </div>
    </header>
  );
};

export default Header;
