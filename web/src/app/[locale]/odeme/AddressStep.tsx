"use client";

import { type FormEvent, useState } from "react";
import { useTranslations } from "next-intl";
import { useToastStore } from "@/stores/toastStore";
import styles from "./odeme.module.scss";

export type Address = {
  name: string;
  surname: string;
  email: string;
  line: string;
  district: string;
  city: string;
};

type AddressStepProps = {
  onContinue: (address: Address) => void;
};

const AddressStep = ({ onContinue }: AddressStepProps) => {
  const t = useTranslations("checkout");
  const showToast = useToastStore((s) => s.show);
  const [form, setForm] = useState<Address>({
    name: "",
    surname: "",
    email: "",
    line: "",
    district: "",
    city: "İstanbul",
  });

  const update = (key: keyof Address, value: string) => setForm((f) => ({ ...f, [key]: value }));

  const handleContinueSubmit = (event: FormEvent) => {
    event.preventDefault();
    if (!form.name || !form.surname || !form.email || !form.line) {
      showToast(t("fillAddress"));
      return;
    }
    onContinue(form);
  };

  return (
    <form onSubmit={handleContinueSubmit}>
      <div className={styles.group}>
        <label>{t("name")}</label>
        <input value={form.name} onChange={(e) => update("name", e.target.value)} />
      </div>
      <div className={styles.group}>
        <label>{t("surname")}</label>
        <input value={form.surname} onChange={(e) => update("surname", e.target.value)} />
      </div>
      <div className={styles.group}>
        <label>{t("email")}</label>
        <input type="email" value={form.email} onChange={(e) => update("email", e.target.value)} />
      </div>
      <div className={styles.group}>
        <label>{t("address")}</label>
        <input value={form.line} onChange={(e) => update("line", e.target.value)} />
      </div>
      <div className={styles.row}>
        <div className={styles.group}>
          <label>{t("district")}</label>
          <input value={form.district} onChange={(e) => update("district", e.target.value)} />
        </div>
        <div className={styles.group}>
          <label>{t("city")}</label>
          <select value={form.city} onChange={(e) => update("city", e.target.value)}>
            <option>İstanbul</option>
            <option>Ankara</option>
            <option>İzmir</option>
            <option>Bursa</option>
            <option>Antalya</option>
          </select>
        </div>
      </div>
      <button className="btn btnPrimary" type="submit">
        {t("continue")} →
      </button>
    </form>
  );
};

export default AddressStep;
