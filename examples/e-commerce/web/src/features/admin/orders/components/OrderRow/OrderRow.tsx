"use client";

import { useTranslations } from "next-intl";

import type { AdminOrder } from "@/shared/types/adminOrder.types";

import { formatTRY } from "@/shared/helpers/money";

import { STATUS_TONE, type StatusTone } from "@/features/admin/_shared/types/statusBadge.types";
import type { OrderDraft } from "@/features/admin/orders/types/orderRow.types";

import StatusBadge from "@/features/admin/_shared/components/StatusBadge";


type Props = {
  order: AdminOrder;
  draft: OrderDraft;
  onDraftChange: (orderId: string, draft: OrderDraft) => void;
  onSave: (orderId: string) => void;
};

const ORDER_STATUS_TONE: Record<AdminOrder["status"], StatusTone> = {
  pending: STATUS_TONE.WARNING,
  paid: STATUS_TONE.ACTIVE,
  shipped: STATUS_TONE.ACTIVE,
  cancelled: STATUS_TONE.DANGER,
};

const STATUS_LABEL_KEY: Record<AdminOrder["status"], string> = {
  pending: "preparing",
  paid: "preparing",
  shipped: "shipped",
  cancelled: "cancelled",
};

const SELECTABLE_STATUSES: AdminOrder["status"][] = ["pending", "paid", "shipped", "cancelled"];

const ID_DISPLAY_LENGTH = 8;

const OrderRow = ({ order, draft, onDraftChange, onSave }: Props) => {
  const t = useTranslations("admin");

  return (
    <tr>
      <td>{order.id.slice(0, ID_DISPLAY_LENGTH)}</td>
      <td>{order.email}</td>
      <td>{formatTRY(order.totalCents)}</td>
      <td>
        <select
          aria-label={t("statusCol")}
          value={draft.status}
          onChange={(e) =>
            onDraftChange(order.id, { ...draft, status: e.target.value as AdminOrder["status"] })
          }
        >
          {SELECTABLE_STATUSES.map((status) => (
            <option key={status} value={status}>
              {t(STATUS_LABEL_KEY[status])}
            </option>
          ))}
        </select>
        <StatusBadge label={t(STATUS_LABEL_KEY[draft.status])} tone={ORDER_STATUS_TONE[draft.status]} />
      </td>
      <td>
        <input
          aria-label={t("trackingNumber")}
          value={draft.trackingNumber}
          onChange={(e) => onDraftChange(order.id, { ...draft, trackingNumber: e.target.value })}
        />
      </td>
      <td>
        <button type="button" className="btn btnPrimary btnSm" onClick={() => onSave(order.id)}>
          {t("save")}
        </button>
      </td>
    </tr>
  );
};

export default OrderRow;
