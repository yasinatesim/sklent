"use client";

import { type FormEvent, useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useToastStore } from "@/stores/toastStore";
import styles from "../admin.module.scss";

type Coupon = {
  code: string;
  disc: number;
  active: boolean;
};

const STORAGE_KEY = "vcAdminCoups";
const SEED: Coupon[] = [
  { code: "VELA15", disc: 15, active: true },
  { code: "BAHAR10", disc: 10, active: true },
];

const AdminCouponsPage = () => {
  const t = useTranslations("admin");
  const showToast = useToastStore((s) => s.show);
  const [items, setItems] = useState<Coupon[]>([]);
  const [code, setCode] = useState("");
  const [disc, setDisc] = useState(10);

  useEffect(() => {
    const raw = window.localStorage.getItem(STORAGE_KEY);
    setItems(raw ? JSON.parse(raw) : SEED);
  }, []);

  const persist = (next: Coupon[]) => {
    setItems(next);
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
  };

  const handleAddSubmit = (event: FormEvent) => {
    event.preventDefault();
    if (!code || !disc) {
      showToast(t("couponRequired"));
      return;
    }
    persist([...items, { code: code.toUpperCase(), disc, active: true }]);
    setCode("");
    showToast(`${t("created")}: ${code.toUpperCase()}`);
  };

  const toggle = (index: number) =>
    persist(items.map((c, i) => (i === index ? { ...c, active: !c.active } : c)));

  const remove = (index: number) => persist(items.filter((_, i) => i !== index));

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
          {items.map((c, i) => (
            <tr key={c.code}>
              <td>{c.code}</td>
              <td>%{c.disc}</td>
              <td>
                <span className={`${styles.status} ${c.active ? styles.statusActive : styles.statusPassive}`}>
                  {c.active ? t("active") : t("passive")}
                </span>
              </td>
              <td>
                <button className="btn btnOutline btnSm" onClick={() => toggle(i)}>
                  {c.active ? t("makePassive") : t("makeActive")}
                </button>{" "}
                <button className="btn btnDanger btnSm" onClick={() => remove(i)}>
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
