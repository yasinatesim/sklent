import { API_BASE, primeCsrf } from "./api";

export type AdminOrderItem = {
  id: string;
  productId: string;
  titleTr: string;
  unitCents: number;
  quantity: number;
};

export type AdminOrder = {
  id: string;
  email: string;
  status: "pending" | "paid" | "shipped" | "cancelled";
  paymentMethod: string;
  totalCents: number;
  trackingNumber?: string;
  items: AdminOrderItem[];
  createdAt: string;
};

export const fetchAdminOrders = async (status?: string): Promise<AdminOrder[]> => {
  const query = status ? `?status=${status}` : "";
  const res = await fetch(`${API_BASE}/admin/orders${query}`, { credentials: "include", cache: "no-store" });
  if (!res.ok) throw new Error("failed to load orders");
  const data = await res.json();
  return data.items ?? [];
};

export const updateOrderStatus = async (
  orderId: string,
  status: string,
  trackingNumber: string,
): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/orders/${orderId}/status`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ status, trackingNumber }),
  });
};
