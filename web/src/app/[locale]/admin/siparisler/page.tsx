"use client";

import { useTranslations } from "next-intl";
import styles from "../admin.module.scss";

const ORDERS = [
  { id: "#VC-2026-042", customer: "Ayşe Yılmaz", amount: "₺2.499", status: "active", date: "12.06.2026" },
  { id: "#VC-2026-041", customer: "Mehmet Demir", amount: "₺599", status: "active", date: "11.06.2026" },
  { id: "#VC-2026-040", customer: "Zeynep Kaya", amount: "₺1.299", status: "active", date: "10.06.2026" },
  { id: "#VC-2026-039", customer: "Ali Öztürk", amount: "₺349", status: "passive", date: "09.06.2026" },
  { id: "#VC-2026-038", customer: "Elif Şahin", amount: "₺1.848", status: "active", date: "09.06.2026" },
];

const AdminOrdersPage = () => {
  const t = useTranslations("admin");

  return (
    <div>
      <h2>{t("orders")}</h2>
      <table className={styles.table}>
        <thead>
          <tr>
            <th>{t("order")}</th>
            <th>{t("customer")}</th>
            <th>{t("amount")}</th>
            <th>{t("statusCol")}</th>
            <th>{t("date")}</th>
          </tr>
        </thead>
        <tbody>
          {ORDERS.map((o) => (
            <tr key={o.id}>
              <td>{o.id}</td>
              <td>{o.customer}</td>
              <td>{o.amount}</td>
              <td>
                <span className={`${styles.status} ${o.status === "active" ? styles.statusActive : styles.statusPassive}`}>
                  {o.status === "active" ? t("preparing") : t("cancelled")}
                </span>
              </td>
              <td>{o.date}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default AdminOrdersPage;
