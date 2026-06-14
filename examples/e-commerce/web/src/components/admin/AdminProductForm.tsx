"use client";

import { type FormEvent, useState } from "react";
import { useTranslations } from "next-intl";
import { createProduct, type Category } from "@/lib/api";
import { useToastStore } from "@/stores/toastStore";
import { REQUEST_STATUS, type RequestStatus } from "@/constants/requestStatus";
import styles from "@/app/[locale]/admin/admin.module.scss";

type AdminProductFormProps = {
  categories: Category[];
  onCreated: () => void;
};

const BADGES = ["", "Yeni", "Çok Satan", "Fırsat", "%10 İndirim", "%20 İndirim"];

const AdminProductForm = ({ categories, onCreated }: AdminProductFormProps) => {
  const t = useTranslations("admin");
  const showToast = useToastStore((s) => s.show);

  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [priceTl, setPriceTl] = useState(199);
  const [oldPriceTl, setOldPriceTl] = useState(0);
  const [stock, setStock] = useState(10);
  const [categorySlug, setCategorySlug] = useState(categories[0]?.slug ?? "");
  const [badge, setBadge] = useState("");
  const [seller, setSeller] = useState("Vela Commerce");
  const [published, setPublished] = useState(true);
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);

  const handleCreateSubmit = async (event: FormEvent) => {
    event.preventDefault();
    if (!title || priceTl <= 0) {
      showToast(t("nameAndPriceRequired"));
      return;
    }
    setStatus(REQUEST_STATUS.LOADING);
    const res = await createProduct({
      titleTr: title,
      descriptionTr: description,
      priceCents: Math.round(priceTl * 100),
      oldPriceCents: Math.round(oldPriceTl * 100),
      stock,
      categorySlug,
      badge,
      seller,
      published,
    });
    if (!res.ok) {
      setStatus(REQUEST_STATUS.ERROR);
      showToast(t("createFailed"));
      return;
    }
    setStatus(REQUEST_STATUS.SUCCESS);
    setTitle("");
    setDescription("");
    showToast(`${t("created")}: ${title}`);
    onCreated();
  };

  return (
    <form onSubmit={handleCreateSubmit}>
      <div className={styles.formGrid}>
        <label className={`${styles.field} ${styles.span2}`}>
          {t("title")}
          <input value={title} onChange={(e) => setTitle(e.target.value)} required />
        </label>
        <label className={styles.field}>
          {t("category")}
          <select value={categorySlug} onChange={(e) => setCategorySlug(e.target.value)}>
            {categories.map((c) => (
              <option key={c.slug} value={c.slug}>
                {c.nameTr}
              </option>
            ))}
          </select>
        </label>
        <label className={`${styles.field} ${styles.span2}`}>
          {t("description")}
          <input value={description} onChange={(e) => setDescription(e.target.value)} />
        </label>
        <label className={styles.field}>
          {t("seller")}
          <input value={seller} onChange={(e) => setSeller(e.target.value)} />
        </label>
        <label className={styles.field}>
          {t("priceTl")}
          <input type="number" value={priceTl} onChange={(e) => setPriceTl(Number(e.target.value))} />
        </label>
        <label className={styles.field}>
          {t("oldPriceTl")}
          <input type="number" value={oldPriceTl} onChange={(e) => setOldPriceTl(Number(e.target.value))} />
        </label>
        <label className={styles.field}>
          {t("stock")}
          <input type="number" value={stock} onChange={(e) => setStock(Number(e.target.value))} />
        </label>
        <label className={styles.field}>
          {t("badge")}
          <select value={badge} onChange={(e) => setBadge(e.target.value)}>
            {BADGES.map((b) => (
              <option key={b} value={b}>
                {b || t("noBadge")}
              </option>
            ))}
          </select>
        </label>
        <label className={styles.field}>
          {t("published")}
          <span>
            <input type="checkbox" checked={published} onChange={(e) => setPublished(e.target.checked)} />
          </span>
        </label>
      </div>
      <div className={styles.formActions}>
        <button className="btn btnPrimary btnSm" type="submit" disabled={status === REQUEST_STATUS.LOADING}>
          + {t("create")}
        </button>
      </div>
    </form>
  );
};

export default AdminProductForm;
