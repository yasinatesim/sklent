import { useTranslations } from "next-intl";

import styles from "./EmptyView.module.scss";

type Props = {
  message?: string;
};

const EmptyView = ({ message }: Props) => {
  const t = useTranslations("common");
  return (
    <div className={styles.box} role="status">
      <span>{message ?? t("empty")}</span>
    </div>
  );
};

export default EmptyView;
