"use client";

import { type FormEvent, useCallback, useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useToastStore } from "@/stores/toastStore";
import { createCoupon, deleteCoupon, fetchCoupons, updateCoupon, type Coupon } from "@/lib/promotions-api";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import styles from "../admin.module.scss";

const AdminCouponsPage = () => {
  const t = useTranslations("admin");
  const showToast = useToastStore((s) => s.show);
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);
  const [items, setItems] = useState<Coupon[]>([]);
  const [code, setCode] = useState("");
  const [disc, setDisc] = useState(10);

  const load = useCallback(async () => {
    setStatus(REQUEST_STATUS.LOADING);
    try {
      setItems(await fetchCoupons());
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
    if (!code || !disc) {
      showToast(t("couponRequired"));
      return;
    }
    const res = await createCoupon({
      code, discountType: "percent", discountValue: disc, scopeType: "all", minCartCents: 0,
    });
    if (!res.ok) {
      showToast(t("createFailed"));
      return;
    }
    setCode("");
    showToast(`${t("created")}: ${code.toUpperCase()}`);
    load();
  };

  const toggle = async (c: Coupon) => {
    await updateCoupon(c.id, { active: !c.active });
    load();
  };

  const remove = async (id: string) => {
    await deleteCoupon(id);
    load();
  };

  return (
    <div>
      <h2>{t("coupons")}</h2>
      <form className={styles.toolbar} onSubmit={handleAddSubmit}>
        <label className={styles.field}>
          {t("couponCode")}
          <input value={code} onChange={(e) => setCode(e.target.value)} />
        </label>
        <label className={styles.field}>
          {t("discountPct")}
          <input type="number" value={disc} onChange={(e) => setDisc(Number(e.target.value))} />
        </label>
        <button className="btn btnPrimary btnSm" type="submit">
          + {t("create")}
        </button>
      </form>
      {status === REQUEST_STATUS.LOADING ? <p>{t("loading")}</p> : null}
      <table className={styles.table}>
        <thead>
          <tr>
            <th>{t("code")}</th>
            <th>{t("discount")}</th>
            <th>{t("statusCol")}</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {items.map((c) => (
            <tr key={c.id}>
              <td>{c.code}</td>
              <td>%{c.discountValue}</td>
              <td>
                <span className={`${styles.status} ${c.active ? styles.statusActive : styles.statusPassive}`}>
                  {c.active ? t("active") : t("passive")}
                </span>
              </td>
              <td>
                <button className="btn btnOutline btnSm" onClick={() => toggle(c)}>
                  {c.active ? t("makePassive") : t("makeActive")}
                </button>{" "}
                <button className="btn btnDanger btnSm" onClick={() => remove(c.id)}>
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

export default AdminCouponsPage;
