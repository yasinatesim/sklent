"use client";

import { useTranslations } from "next-intl";

import { STATUS_TONE } from "@/features/admin/_shared/types/statusBadge.types";

import StatusBadge from "@/features/admin/_shared/components/StatusBadge";
import WeeklyOrdersChart from "@/features/admin/dashboard/components/WeeklyOrdersChart";

import styles from "@/features/admin/_shared/styles/admin.module.scss";

const RECENT_ORDERS = [
  { id: "#VC-2026-042", customer: "Ayşe Yılmaz", amount: "₺2.499", status: "active" },
  { id: "#VC-2026-041", customer: "Mehmet Demir", amount: "₺599", status: "active" },
  { id: "#VC-2026-040", customer: "Zeynep Kaya", amount: "₺1.299", status: "active" },
  { id: "#VC-2026-039", customer: "Ali Öztürk", amount: "₺349", status: "passive" },
];

const DashboardPage = () => {
  const t = useTranslations("admin");

  const stats = [
    { value: "24", label: t("statOrders") },
    { value: "₺18.240", label: t("statRevenue") },
    { value: "25", label: t("statProducts") },
    { value: "189", label: t("statCustomers") },
  ];

  return (
    <div>
      <h2>{t("dashboard")}</h2>
      <div className={styles.stats}>
        {stats.map((s) => (
          <div key={s.label} className={styles.statCard}>
            <div className={styles.statValue}>{s.value}</div>
            <div className={styles.statLabel}>{s.label}</div>
          </div>
        ))}
      </div>
      <div className={styles.chartBox}>
        <h3>{t("weeklyChart")}</h3>
        <WeeklyOrdersChart />
      </div>
      <div className={styles.chartBox}>
        <h3>{t("recentOrders")}</h3>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>{t("order")}</th>
              <th>{t("customer")}</th>
              <th>{t("amount")}</th>
              <th>{t("statusCol")}</th>
            </tr>
          </thead>
          <tbody>
            {RECENT_ORDERS.map((o) => (
              <tr key={o.id}>
                <td>{o.id}</td>
                <td>{o.customer}</td>
                <td>{o.amount}</td>
                <td>
                  <StatusBadge
                    tone={o.status === "active" ? STATUS_TONE.ACTIVE : STATUS_TONE.PASSIVE}
                    label={o.status === "active" ? t("preparing") : t("cancelled")}
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DashboardPage;
