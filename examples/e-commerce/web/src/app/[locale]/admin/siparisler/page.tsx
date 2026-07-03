"use client";

import { useCallback, useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useToastStore } from "@/stores/toastStore";
import { fetchAdminOrders, updateOrderStatus, type AdminOrder } from "@/lib/orders-api";
import { formatTRY } from "@/lib/api";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import styles from "../admin.module.scss";

const STATUS_LABEL_KEY: Record<AdminOrder["status"], string> = {
  pending: "preparing",
  paid: "preparing",
  shipped: "shipped",
  cancelled: "cancelled",
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

  return (
    <div>
      <h2>{t("orders")}</h2>
      {status === REQUEST_STATUS.LOADING ? <p>{t("loading")}</p> : null}
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
          {orders.map((o) => {
            const draft = drafts[o.id] ?? { status: o.status, trackingNumber: "" };
            return (
              <tr key={o.id}>
                <td>{o.id.slice(0, 8)}</td>
                <td>{o.email}</td>
                <td>{formatTRY(o.totalCents)}</td>
                <td>
                  <select
                    value={draft.status}
                    onChange={(e) =>
                      setDrafts((prev) => ({
                        ...prev,
                        [o.id]: { ...draft, status: e.target.value as AdminOrder["status"] },
                      }))
                    }
                  >
                    <option value="pending">{t("preparing")}</option>
                    <option value="paid">{t("preparing")}</option>
                    <option value="shipped">{t("shipped")}</option>
                    <option value="cancelled">{t("cancelled")}</option>
                  </select>
                  <span className={`${styles.status} ${draft.status === "cancelled" ? styles.statusPassive : styles.statusActive}`}>
                    {t(STATUS_LABEL_KEY[draft.status])}
                  </span>
                </td>
                <td>
                  <input
                    value={draft.trackingNumber}
                    onChange={(e) =>
                      setDrafts((prev) => ({ ...prev, [o.id]: { ...draft, trackingNumber: e.target.value } }))
                    }
                  />
                </td>
                <td>
                  <button className="btn btnPrimary btnSm" onClick={() => save(o.id)}>
                    {t("save")}
                  </button>
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};

export default AdminOrdersPage;
