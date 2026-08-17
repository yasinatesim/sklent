"use client";

import { type ComponentProps, type ComponentType, useCallback, useEffect, useState } from "react";

import { useTranslations } from "next-intl";

import { REQUEST_STATUS, type RequestStatus } from "@/shared/constants/requestStatus";

import ErrorView from "@/shared/ui/ErrorView";
import LoadingView from "@/shared/ui/LoadingView";
import NullView from "@/shared/ui/NullView";

import { STATUS_TONE, type StatusTone } from "@/features/admin/_shared/types/statusBadge.types";
import type { Review } from "@/features/admin/reviews/types/reviews.types";

import { fetchAdminReviews, updateReviewStatus } from "@/features/admin/reviews/api/reviews";

import StatusBadge from "@/features/admin/_shared/components/StatusBadge";

import styles from "@/features/admin/_shared/styles/admin.module.scss";

const REVIEW_STATUS_TONE: Record<Review["status"], StatusTone> = {
  pending: STATUS_TONE.WARNING,
  approved: STATUS_TONE.ACTIVE,
  rejected: STATUS_TONE.DANGER,
};

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

const AdminReviewsPage = () => {
  const t = useTranslations("admin");
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);
  const [items, setItems] = useState<Review[]>([]);

  const load = useCallback(async () => {
    setStatus(REQUEST_STATUS.LOADING);
    try {
      setItems(await fetchAdminReviews());
      setStatus(REQUEST_STATUS.SUCCESS);
    } catch {
      setStatus(REQUEST_STATUS.ERROR);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  const decide = async (id: string, next: "approved" | "rejected") => {
    await updateReviewStatus(id, next);
    load();
  };

  const CurrentView = STATE_VIEWS[status];

  return (
    <div>
      <h2>{t("reviews")}</h2>
      <CurrentView onRetry={load} />
      <table className={styles.table}>
        <thead>
          <tr>
            <th>{t("reviewAuthor")}</th>
            <th>{t("reviewRating")}</th>
            <th>{t("reviewComment")}</th>
            <th>{t("statusCol")}</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {items.map((r) => (
            <tr key={r.id}>
              <td>{r.authorName}</td>
              <td>{"★".repeat(r.rating)}</td>
              <td>{r.comment}</td>
              <td>
                <StatusBadge label={t(r.status)} tone={REVIEW_STATUS_TONE[r.status]} />
              </td>
              <td>
                {r.status === "pending" ? (
                  <>
                    <button className="btn btnOutline btnSm" onClick={() => decide(r.id, "approved")}>
                      {t("approve")}
                    </button>{" "}
                    <button className="btn btnDanger btnSm" onClick={() => decide(r.id, "rejected")}>
                      {t("reject")}
                    </button>
                  </>
                ) : null}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminReviewsPage;
