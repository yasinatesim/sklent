import { localePath,ROUTES } from "@/shared/constants/routes";

import type { Category } from "@/shared/types/catalog.types";

import styles from "./CategoryCard.module.scss";

type CategoryCardProps = {
  category: Category;
  locale: string;
  showDesc?: boolean;
};

const CategoryCard = ({ category, locale, showDesc = true }: CategoryCardProps) => {
  const name = locale === "en" ? category.nameEn : category.nameTr;
  return (
    <a className={styles.card} href={localePath(locale, `${ROUTES.CATEGORY}/category.slug`)}>
      <div className={styles.icon}>{category.icon || "📦"}</div>
      <h3 className={styles.name}>{name}</h3>
      {showDesc && category.descTr ? <p className={styles.desc}>{category.descTr}</p> : null}
    </a>
  );
};

export default CategoryCard;
