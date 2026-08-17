"use client";

import { type ComponentProps, type ComponentType, type FormEvent, useMemo, useState } from "react";

import { useTranslations } from "next-intl";

import { REQUEST_STATUS, type RequestStatus } from "@/shared/constants/requestStatus";

import { useToastStore } from "@/shared/stores/toastStore";

import ErrorView from "@/shared/ui/ErrorView";
import LoadingView from "@/shared/ui/LoadingView";
import NullView from "@/shared/ui/NullView";

import type { Promotion } from "@/features/admin/promotions/types/promotions.types";

import { createPromotion, deletePromotion, fetchPromotions, updatePromotion } from "@/features/admin/promotions/api/promotions";

import { useDiscountRules } from "@/features/admin/_shared/hooks/useDiscountRules";

import DiscountRuleTable from "@/features/admin/_shared/components/DiscountRuleTable";

import styles from "@/features/admin/_shared/styles/admin.module.scss";

type StateViewProps = ComponentProps<typeof LoadingView> &
  ComponentProps<typeof ErrorView> &
  ComponentProps<typeof NullView>;

// Component references, never elements: an element map builds every branch on every render.
const STATE_VIEWS: Record<RequestStatus, ComponentType<StateViewProps>> = {
  [REQUEST_STATUS.IDLE]: LoadingView,
  [REQUEST_STATUS.LOADING]: LoadingView,
  [REQUEST_STATUS.ERROR]: ErrorView,
  [REQUEST_STATUS.SUCCESS]: NullView,
};

const DEFAULT_DISCOUNT_PCT = 10;

const AdminCampaignsPage = () => {
  const t = useTranslations("admin");
  const showToast = useToastStore((s) => s.show);
  const api = useMemo(
    () => ({ fetchAll: fetchPromotions, update: updatePromotion, remove: deletePromotion }),
    [],
  );
  const { status, items, load, toggle, remove } = useDiscountRules<Promotion>(api);
  const [name, setName] = useState("");
  const [discount, setDiscount] = useState(DEFAULT_DISCOUNT_PCT);

  const handleAddSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (!name || !discount) {
      showToast(t("campaignRequired"));
      return;
    }
    const res = await createPromotion({
      name,
      discountType: "percent",
      discountValue: discount,
      scopeType: "all",
      minCartCents: 0,
    });
    if (!res.ok) {
      showToast(t("createFailed"));
      return;
    }
    setName("");
    showToast(`${t("created")}: ${name}`);
    load();
  };

  const CurrentView = STATE_VIEWS[status];

  return (
    <div>
      <h2>{t("campaigns")}</h2>
      <form className={styles.toolbar} onSubmit={handleAddSubmit}>
        <label className={styles.field}>
          {t("campaignName")}
          <input value={name} onChange={(e) => setName(e.target.value)} />
        </label>
        <label className={styles.field}>
          {t("discountPct")}
          <input
            type="number"
            value={discount}
            onChange={(e) => setDiscount(Number(e.target.value))}
          />
        </label>
        <button className="btn btnPrimary btnSm" type="submit">
          + {t("create")}
        </button>
      </form>
      <CurrentView onRetry={load} />
      <DiscountRuleTable
        items={items}
        labelHeader="campaignName"
        labelOf={(item) => item.name}
        discountOf={(item) => item.discountValue}
        onToggle={toggle}
        onRemove={remove}
      />
    </div>
  );
};

export default AdminCampaignsPage;
