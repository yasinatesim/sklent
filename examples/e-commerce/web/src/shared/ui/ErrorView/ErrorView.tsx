import { useTranslations } from "next-intl";

import styles from "./ErrorView.module.scss";

type Props = {
  message?: string;
  onRetry?: () => void;
};

const ErrorView = ({ message, onRetry }: Props) => {
  const t = useTranslations("common");
  return (
    <div className={`${styles.box} ${styles.error}`} role="alert">
      <span>{message ?? t("errorGeneric")}</span>
      {onRetry && (
        <button type="button" className="btn btnOutline btnSm" onClick={onRetry}>
          {t("retry")}
        </button>
      )}
    </div>
  );
};

export default ErrorView;
