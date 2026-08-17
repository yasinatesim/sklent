"use client";

import { type ComponentProps, type ComponentType, type FormEvent, useCallback, useEffect, useState } from "react";

import { useTranslations } from "next-intl";

import { REQUEST_STATUS, type RequestStatus } from "@/shared/constants/requestStatus";

import { useToastStore } from "@/shared/stores/toastStore";

import ErrorView from "@/shared/ui/ErrorView";
import LoadingView from "@/shared/ui/LoadingView";
import NullView from "@/shared/ui/NullView";

import type { OrderNotifyRule } from "@/features/admin/settings/types/orderNotifyRules.types";

import {
  createOrderNotifyRule,
  deleteOrderNotifyRule,
  fetchOrderNotifyRules,
  updateOrderNotifyRule,
} from "@/features/admin/settings/api/orderNotifyRules";

import { ORDER_NOTIFY_SOURCE_LABELS, ORDER_NOTIFY_SOURCES } from "@/features/admin/settings/constants";

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

const AdminSettingsPage = () => {
  const t = useTranslations("admin");
  const showToast = useToastStore((s) => s.show);
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);
  const [rules, setRules] = useState<OrderNotifyRule[]>([]);
  const [source, setSource] = useState<string>(ORDER_NOTIFY_SOURCES[0].value);
  const [recipient, setRecipient] = useState("");

  const load = useCallback(async () => {
    setStatus(REQUEST_STATUS.LOADING);
    try {
      setRules(await fetchOrderNotifyRules());
      setStatus(REQUEST_STATUS.SUCCESS);
    } catch {
      setStatus(REQUEST_STATUS.ERROR);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const handleAddSubmit = async (event: FormEvent) => {
    event.preventDefault();
    const res = await createOrderNotifyRule(source, recipient.trim());
    if (!res.ok) {
      showToast(t("createFailed"));
      return;
    }
    setRecipient("");
    load();
  };

  const toggleRule = async (rule: OrderNotifyRule) => {
    await updateOrderNotifyRule(rule.id, !rule.enabled);
    load();
  };

  const removeRule = async (id: string) => {
    await deleteOrderNotifyRule(id);
    load();
  };

  const CurrentView = STATE_VIEWS[status];

  return (
    <div>
      <h2>{t("orderNotifications")}</h2>
      <p>{t("orderNotificationsHint")}</p>
      <form className={styles.toolbar} onSubmit={handleAddSubmit}>
        <label className={styles.field}>
          {t("channel")}
          <select value={source} onChange={(e) => setSource(e.target.value)}>
            {ORDER_NOTIFY_SOURCES.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </label>
        <label className={styles.field}>
          {t("recipient")}
          <input type="email" required value={recipient} onChange={(e) => setRecipient(e.target.value)} />
        </label>
        <button className="btn btnPrimary btnSm" type="submit">
          + {t("addRow")}
        </button>
      </form>
      <CurrentView onRetry={load} />
      <table className={styles.table}>
        <thead>
          <tr>
            <th>{t("recipient")}</th>
            <th>{t("channel")}</th>
            <th>{t("statusCol")}</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {rules.map((rule) => (
            <tr key={rule.id}>
              <td>{rule.recipient}</td>
              <td>{ORDER_NOTIFY_SOURCE_LABELS[rule.source] ?? rule.source}</td>
              <td>{rule.enabled ? t("enabled") : t("disabled")}</td>
              <td>
                <button className="btn btnSm" onClick={() => toggleRule(rule)}>
                  {rule.enabled ? t("disable") : t("enable")}
                </button>
                <button className="btn btnDanger btnSm" onClick={() => removeRule(rule.id)}>
                  {t("delete")}
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminSettingsPage;
