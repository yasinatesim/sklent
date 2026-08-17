"use client";

import { type ComponentProps, type ComponentType, useCallback, useEffect, useState } from "react";

import { useTranslations } from "next-intl";

import { REQUEST_STATUS, type RequestStatus } from "@/shared/constants/requestStatus";

import type { AdminOrder } from "@/shared/types/adminOrder.types";

import { useToastStore } from "@/shared/stores/toastStore";

import ErrorView from "@/shared/ui/ErrorView";
import LoadingView from "@/shared/ui/LoadingView";
import NullView from "@/shared/ui/NullView";

import { fetchAdminOrders, updateOrderStatus } from "@/features/web/orders/api/ordersExtra";

import OrderRow from "@/features/admin/orders/components/OrderRow";

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

const AdminOrdersPage = () => {
  const t = useTranslations("admin");
  const showToast = useToastStore((s) => s.show);
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);
  const [orders, setOrders] = useState<AdminOrder[]>([]);
  const [drafts, setDrafts] = useState<Record<string, { status: AdminOrder["status"]; trackingNumber: string }>>({});

  const load = useCallback(async () => {
    setStatus(REQUEST_STATUS.LOADING);
    try {
      const items = await fetchAdminOrders();
      setOrders(items);
      setDrafts(
        Object.fromEntries(
          items.map((o) => [o.id, { status: o.status, trackingNumber: o.trackingNumber ?? "" }]),
        ),
      );
      setStatus(REQUEST_STATUS.SUCCESS);
    } catch {
      setStatus(REQUEST_STATUS.ERROR);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const save = async (orderId: string) => {
    const draft = drafts[orderId];
    if (!draft) return;
    const res = await updateOrderStatus(orderId, draft.status, draft.trackingNumber);
    if (!res.ok) {
      showToast(t("updateFailed"));
      return;
    }
    showToast(t("orderUpdated"));
    load();
  };

  const CurrentView = STATE_VIEWS[status];

  return (
    <div>
      <h2>{t("orders")}</h2>
      <CurrentView onRetry={load} />
      <table className={styles.table}>
        <thead>
          <tr>
            <th>{t("order")}</th>
            <th>{t("customer")}</th>
            <th>{t("amount")}</th>
            <th>{t("statusCol")}</th>
            <th>{t("trackingNumber")}</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {orders.map((o) => (
            <OrderRow
              key={o.id}
              order={o}
              draft={drafts[o.id] ?? { status: o.status, trackingNumber: "" }}
              onDraftChange={(orderId, draft) => setDrafts((prev) => ({ ...prev, [orderId]: draft }))}
              onSave={save}
            />
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminOrdersPage;
