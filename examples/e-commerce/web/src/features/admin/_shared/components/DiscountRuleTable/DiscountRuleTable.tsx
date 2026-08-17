"use client";

import { useTranslations } from "next-intl";

import { STATUS_TONE } from "@/features/admin/_shared/types/statusBadge.types";

import StatusBadge from "@/features/admin/_shared/components/StatusBadge";

import styles from "@/features/admin/_shared/styles/admin.module.scss";

type ActivatableRule = {
  id: string;
  active: boolean;
};

type Props<T extends ActivatableRule> = {
  items: T[];
  labelHeader: string;
  labelOf: (item: T) => string;
  discountOf: (item: T) => number;
  onToggle: (item: T) => void;
  onRemove: (id: string) => void;
};

const DiscountRuleTable = <T extends ActivatableRule>({
  items,
  labelHeader,
  labelOf,
  discountOf,
  onToggle,
  onRemove,
}: Props<T>) => {
  const t = useTranslations("admin");

  return (
    <table className={styles.table}>
      <thead>
        <tr>
          <th>{t(labelHeader)}</th>
          <th>{t("discount")}</th>
          <th>{t("statusCol")}</th>
          <th />
        </tr>
      </thead>
      <tbody>
        {items.map((item) => (
          <tr key={item.id}>
            <td>{labelOf(item)}</td>
            <td>%{discountOf(item)}</td>
            <td>
              <StatusBadge
                tone={item.active ? STATUS_TONE.ACTIVE : STATUS_TONE.PASSIVE}
                label={item.active ? t("active") : t("passive")}
              />
            </td>
            <td>
              <button type="button" className="btn btnOutline btnSm" onClick={() => onToggle(item)}>
                {item.active ? t("makePassive") : t("makeActive")}
              </button>{" "}
              <button type="button" className="btn btnDanger btnSm" onClick={() => onRemove(item.id)}>
                {t("delete")}
              </button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default DiscountRuleTable;
