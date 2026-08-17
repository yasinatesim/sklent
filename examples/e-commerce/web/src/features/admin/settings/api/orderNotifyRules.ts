import { API_BASE, primeCsrf } from "@/shared/helpers/api/client";

import type { OrderNotifyRule } from "@/features/admin/settings/types/orderNotifyRules.types";

export const fetchOrderNotifyRules = async (): Promise<OrderNotifyRule[]> => {
  const res = await fetch(`${API_BASE}/admin/order-notifications`, { credentials: "include", cache: "no-store" });
  if (!res.ok) throw new Error("failed to load order notification rules");
  const data = await res.json();
  return data.items ?? [];
};

export const createOrderNotifyRule = async (source: string, recipient: string): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/order-notifications`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ source, recipient }),
  });
};

export const updateOrderNotifyRule = async (id: string, enabled: boolean): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/order-notifications/${id}`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ enabled }),
  });
};

export const deleteOrderNotifyRule = async (id: string): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/order-notifications/${id}`, {
    method: "DELETE",
    credentials: "include",
    headers: { "X-CSRF-Token": csrf },
  });
};
