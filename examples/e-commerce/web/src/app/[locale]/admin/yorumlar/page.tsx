"use client";

import { useCallback, useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { fetchAdminReviews, updateReviewStatus, type Review } from "@/lib/reviews-api";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import styles from "../admin.module.scss";

const STATUS_CLASS: Record<Review["status"], string> = {
  pending: "statusPassive",
  approved: "statusActive",
  rejected: "statusPassive",
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

  return (
    <div>
      <h2>{t("reviews")}</h2>
      {status === REQUEST_STATUS.LOADING ? <p>{t("loading")}</p> : null}
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
                <span className={`${styles.status} ${styles[STATUS_CLASS[r.status]]}`}>{t(r.status)}</span>
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
