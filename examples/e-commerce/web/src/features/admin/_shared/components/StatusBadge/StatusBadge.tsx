import { STATUS_TONE, type StatusTone } from "@/features/admin/_shared/types/statusBadge.types";

import styles from "./StatusBadge.module.scss";

type Props = {
  label: string;
  tone?: StatusTone;
};

const TONE_CLASS: Record<StatusTone, string> = {
  [STATUS_TONE.ACTIVE]: styles.active,
  [STATUS_TONE.PASSIVE]: styles.passive,
  [STATUS_TONE.WARNING]: styles.warning,
  [STATUS_TONE.DANGER]: styles.danger,
};

const StatusBadge = ({ label, tone = STATUS_TONE.PASSIVE }: Props) => (
  <span className={`${styles.badge} ${TONE_CLASS[tone]}`} data-tone={tone}>
    {label}
  </span>
);

export default StatusBadge;
