"use client";

import { type FormEvent, useCallback, useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useToastStore } from "@/stores/toastStore";
import {
  createPromotion,
  deletePromotion,
  fetchPromotions,
  updatePromotion,
  type Promotion,
} from "@/lib/promotions-api";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import styles from "../admin.module.scss";

const AdminCampaignsPage = () => {
  const t = useTranslations("admin");
  const showToast = useToastStore((s) => s.show);
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);
  const [items, setItems] = useState<Promotion[]>([]);
  const [name, setName] = useState("");
  const [disc, setDisc] = useState(10);

  const load = useCallback(async () => {
    setStatus(REQUEST_STATUS.LOADING);
    try {
      setItems(await fetchPromotions());
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
    if (!name || !disc) {
      showToast(t("campaignRequired"));
      return;
    }
    const res = await createPromotion({
      name, discountType: "percent", discountValue: disc, scopeType: "all", minCartCents: 0,
    });
    if (!res.ok) {
      showToast(t("createFailed"));
      return;
    }
    setName("");
    showToast(`${t("created")}: ${name}`);
    load();
  };

  const toggle = async (p: Promotion) => {
    await updatePromotion(p.id, { active: !p.active });
    load();
  };

  const remove = async (id: string) => {
    await deletePromotion(id);
    load();
  };

  return (
    <div>
      <h2>{t("campaigns")}</h2>
      <form className={styles.toolbar} onSubmit={handleAddSubmit}>
        <label className={styles.field}>
          {t("campaignName")}
          <input value={name} onChange={(e) => setName(e.target.value)} />
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
            <th>{t("campaign")}</th>
            <th>{t("discount")}</th>
            <th>{t("statusCol")}</th>
            <th />
          </tr>
        </thead>
        <tbody>
          {items.map((p) => (
            <tr key={p.id}>
              <td>{p.name}</td>
              <td>%{p.discountValue}</td>
              <td>
                <span className={`${styles.status} ${p.active ? styles.statusActive : styles.statusPassive}`}>
                  {p.active ? t("active") : t("passive")}
                </span>
              </td>
              <td>
                <button className="btn btnOutline btnSm" onClick={() => toggle(p)}>
                  {p.active ? t("makePassive") : t("makeActive")}
                </button>{" "}
                <button className="btn btnDanger btnSm" onClick={() => remove(p.id)}>
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

export default AdminCampaignsPage;
