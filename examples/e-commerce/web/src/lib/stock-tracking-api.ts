import { API_BASE, primeCsrf } from "./api";

export type StockTrackingItem = {
  id: string;
  productName: string;
  quantity: number;
  createdAt: string;
};

export const fetchStockTracking = async (): Promise<StockTrackingItem[]> => {
  const res = await fetch(`${API_BASE}/admin/stock-tracking`, { credentials: "include", cache: "no-store" });
  if (!res.ok) throw new Error("failed to load stock tracking");
  const data = await res.json();
  return data.items ?? [];
};

export const createStockTrackingItem = async (productName: string, quantity: number): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/stock-tracking`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ productName, quantity }),
  });
};

export const updateStockTrackingItem = async (
  id: string,
  productName: string,
  quantity: number,
): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/stock-tracking/${id}`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ productName, quantity }),
  });
};

export const deleteStockTrackingItem = async (id: string): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/stock-tracking/${id}`, {
    method: "DELETE",
    credentials: "include",
    headers: { "X-CSRF-Token": csrf },
  });
};
