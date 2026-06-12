"use client";

import { useToastStore } from "@/stores/toastStore";
import styles from "./Toast.module.scss";

const Toast = () => {
  const message = useToastStore((s) => s.message);
  const visible = useToastStore((s) => s.visible);

  return (
    <div className={`${styles.toast} ${visible ? styles.show : ""}`} role="status" aria-live="polite">
      {message}
    </div>
  );
};

export default Toast;
