import type { Category } from "@/lib/api";
import styles from "./CategoryCard.module.scss";

type CategoryCardProps = {
  category: Category;
  locale: string;
  showDesc?: boolean;
};

const CategoryCard = ({ category, locale, showDesc = true }: CategoryCardProps) => {
  const name = locale === "en" ? category.nameEn : category.nameTr;
  return (
    <a className={styles.card} href={`/${locale}/kategori/${category.slug}`}>
      <div className={styles.icon}>{category.icon || "📦"}</div>
      <h3 className={styles.name}>{name}</h3>
      {showDesc && category.descTr ? <p className={styles.desc}>{category.descTr}</p> : null}
    </a>
  );
};

export default CategoryCard;
