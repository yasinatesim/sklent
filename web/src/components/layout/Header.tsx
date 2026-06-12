"use client";

import { type FormEvent, useState } from "react";
import { useLocale, useTranslations } from "next-intl";
import { useRouter } from "next/navigation";
import { useCartStore } from "@/stores/cartStore";
import { useThemeStore, THEME } from "@/stores/themeStore";
import { useModalStore, MODAL } from "@/stores/modalStore";
import { useAuthStore } from "@/stores/authStore";
import SearchIcon from "@/components/icons/SearchIcon";
import UserIcon from "@/components/icons/UserIcon";
import CartIcon from "@/components/icons/CartIcon";
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
    router.push(`/${locale}/arama?q=${encodeURIComponent(query)}`);
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
          <button onClick={() => router.push(`/${locale}/giris`)}>
            <UserIcon />
            <span className={styles.label}>{user ? user.email : t("signIn")}</span>
          </button>
          <button onClick={() => router.push(`/${locale}/sepet`)}>
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
