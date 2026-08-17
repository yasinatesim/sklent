"use client";

import { type ComponentProps, type ComponentType, useCallback, useEffect, useState } from "react";

import { useTranslations } from "next-intl";

import { REQUEST_STATUS, type RequestStatus } from "@/shared/constants/requestStatus";

import { useToastStore } from "@/shared/stores/toastStore";

import ErrorView from "@/shared/ui/ErrorView";
import LoadingView from "@/shared/ui/LoadingView";
import NullView from "@/shared/ui/NullView";

import { STATUS_TONE, type StatusTone } from "@/features/admin/_shared/types/statusBadge.types";
import { RETURN_STATUS, type ReturnRequest, type ReturnStatus } from "@/features/admin/returns/types/returns.types";

import { fetchReturns, updateReturnStatus } from "@/features/admin/returns/api/returns";

import { nextStatuses } from "@/features/admin/returns/helpers/returnTransitions";

import StatusBadge from "@/features/admin/_shared/components/StatusBadge";

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

const RETURN_STATUS_TONE: Record<ReturnStatus, StatusTone> = {
  [RETURN_STATUS.REQUESTED]: STATUS_TONE.WARNING,
  [RETURN_STATUS.APPROVED]: STATUS_TONE.ACTIVE,
  [RETURN_STATUS.REJECTED]: STATUS_TONE.DANGER,
  [RETURN_STATUS.REFUNDED]: STATUS_TONE.PASSIVE,
};

const ORDER_ID_DISPLAY_LENGTH = 8;

const AdminReturnsPage = () => {
  const t = useTranslations("admin");
  const tReturn = useTranslations("returns");
  const showToast = useToastStore((s) => s.show);
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);
  const [items, setItems] = useState<ReturnRequest[]>([]);

  const load = useCallback(async () => {
    setStatus(REQUEST_STATUS.LOADING);
    try {
      setItems(await fetchReturns());
      setStatus(REQUEST_STATUS.SUCCESS);
    } catch {
      setItems([]);
      setStatus(REQUEST_STATUS.ERROR);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const applyStatus = async (id: string, next: ReturnStatus) => {
    const res = await updateReturnStatus(id, next, "");
    if (!res.ok) {
      showToast(t("updateFailed"));
      return;
    }
    showToast(tReturn(next));
    load();
  };

  const CurrentView = STATE_VIEWS[status];

  return (
    <div>
      <h2>{tReturn("title")}</h2>
      <CurrentView onRetry={load} />
      <table className={styles.table}>
        <thead>
          <tr>
            <th>{t("order")}</th>
            <th>{t("customer")}</th>
            <th>{tReturn("reason")}</th>
            <th>{t("statusCol")}</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {items.map((item) => (
            <tr key={item.id}>
              <td>{item.orderId.slice(0, ORDER_ID_DISPLAY_LENGTH)}</td>
              <td>{item.email}</td>
              <td>{tReturn(item.reason)}</td>
              <td>
                <StatusBadge label={tReturn(item.status)} tone={RETURN_STATUS_TONE[item.status]} />
              </td>
              <td>
                {nextStatuses(item.status).map((next) => (
                  <button
                    key={next}
                    type="button"
                    className="btn btnOutline btnSm"
                    onClick={() => applyStatus(item.id, next)}
                  >
                    {tReturn(`action_${next}`)}
                  </button>
                ))}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminReturnsPage;
