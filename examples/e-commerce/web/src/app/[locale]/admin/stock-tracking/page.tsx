"use client";

import { type ComponentProps, type ComponentType, type FormEvent, useCallback, useEffect, useMemo, useState } from "react";

import { useTranslations } from "next-intl";

import { REQUEST_STATUS, type RequestStatus } from "@/shared/constants/requestStatus";

import { useToastStore } from "@/shared/stores/toastStore";

import ErrorView from "@/shared/ui/ErrorView";
import LoadingView from "@/shared/ui/LoadingView";
import NullView from "@/shared/ui/NullView";

import type { StockTrackingItem } from "@/features/admin/stocktracking/types/stockTracking.types";

import { createStockTrackingItem, deleteStockTrackingItem, fetchStockTracking, updateStockTrackingItem } from "@/features/admin/stocktracking/api/stockTracking";

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

const AdminStockTrackingPage = () => {
  const t = useTranslations("admin");
  const showToast = useToastStore((s) => s.show);
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);
  const [items, setItems] = useState<StockTrackingItem[]>([]);
  const [productName, setProductName] = useState("");
  const [quantity, setQuantity] = useState(0);
  const [searchQuery, setSearchQuery] = useState("");

  const load = useCallback(async () => {
    setStatus(REQUEST_STATUS.LOADING);
    try {
      setItems(await fetchStockTracking());
      setStatus(REQUEST_STATUS.SUCCESS);
    } catch {
      setStatus(REQUEST_STATUS.ERROR);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const filtered = useMemo(
    () => items.filter((i) => i.productName.toLowerCase().includes(searchQuery.toLowerCase())),
    [items, searchQuery],
  );

  const totalQuantity = useMemo(() => filtered.reduce((sum, i) => sum + i.quantity, 0), [filtered]);

  const handleAddSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (!productName) {
      showToast(t("nameAndPriceRequired"));
      return;
    }
    const res = await createStockTrackingItem(productName, quantity);
    if (!res.ok) {
      showToast(t("createFailed"));
      return;
    }
    setProductName("");
    setQuantity(0);
    load();
  };

  const saveQuantity = async (item: StockTrackingItem, nextQuantity: number) => {
    await updateStockTrackingItem(item.id, item.productName, nextQuantity);
    load();
  };

  const remove = async (id: string) => {
    await deleteStockTrackingItem(id);
    load();
  };

  const CurrentView = STATE_VIEWS[status];

  return (
    <div>
      <h2>{t("stockTracking")}</h2>
      <form className={styles.toolbar} onSubmit={handleAddSubmit}>
        <label className={styles.field}>
          {t("productName")}
          <input value={productName} onChange={(e) => setProductName(e.target.value)} />
        </label>
        <label className={styles.field}>
          {t("quantity")}
          <input type="number" value={quantity} onChange={(e) => setQuantity(Number(e.target.value))} />
        </label>
        <button className="btn btnPrimary btnSm" type="submit">
          + {t("addRow")}
        </button>
      </form>
      <div className={styles.toolbar}>
        <input
          placeholder={t("searchPlaceholder")}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>
      <CurrentView onRetry={load} />
      <table className={styles.table}>
        <thead>
          <tr>
            <th>{t("productName")}</th>
            <th>{t("quantity")}</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {filtered.map((item) => (
            <tr key={item.id}>
              <td>{item.productName}</td>
              <td>
                <input
                  type="number"
                  defaultValue={item.quantity}
                  onBlur={(e) => saveQuantity(item, Number(e.target.value))}
                />
              </td>
              <td>
                <button className="btn btnDanger btnSm" onClick={() => remove(item.id)}>
                  {t("delete")}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
      <p>
        {t("totalItems")}: {filtered.length} · {t("totalQuantity")}: {totalQuantity}
      </p>
    </div>
  );
};

export default AdminStockTrackingPage;
