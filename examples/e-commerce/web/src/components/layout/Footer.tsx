import { getLocale, getTranslations } from "next-intl/server";
import styles from "./Footer.module.scss";

const Footer = async () => {
  const t = await getTranslations("footer");
  const locale = await getLocale();

  return (
    <footer className={styles.footer}>
      <div className={`container ${styles.grid}`}>
        <div className={styles.col}>
          <span className={styles.brand}>
            Vela<span> Commerce</span>
          </span>
          <p>{t("about")}</p>
        </div>
        <div className={styles.col}>
          <h4>{t("categories")}</h4>
          <a href={`/${locale}/kategori/all`}>{t("allProducts")}</a>
          <a href={`/${locale}/kategori/elektronik`}>Elektronik</a>
          <a href={`/${locale}/kategori/moda`}>Moda</a>
          <a href={`/${locale}/kategori/dogal-tas-bileklik`}>Doğal Taş</a>
        </div>
        <div className={styles.col}>
          <h4>{t("corporate")}</h4>
          <a href={`/${locale}`}>{t("aboutUs")}</a>
          <a href={`/${locale}`}>{t("contact")}</a>
        </div>
        <div className={styles.col}>
          <h4>{t("help")}</h4>
          <a href={`/${locale}`}>{t("faq")}</a>
          <a href={`/${locale}`}>{t("shipping")}</a>
          <a href={`/${locale}`}>{t("returns")}</a>
        </div>
        <div className={styles.copyright}>© 2026 Vela Commerce — {t("rights")}</div>
      </div>
    </footer>
  );
};

export default Footer;
