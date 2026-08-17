import { useTranslations } from "next-intl";

import styles from "./LoadingView.module.scss";

type Props = {
  message?: string;
};

const LoadingView = ({ message }: Props) => {
  const t = useTranslations("common");
  return (
    <div className={styles.box} role="status" aria-live="polite">
      <span className={styles.spinner} aria-hidden="true" />
      <span>{message ?? t("loading")}</span>
    </div>
  );
};

export default LoadingView;
