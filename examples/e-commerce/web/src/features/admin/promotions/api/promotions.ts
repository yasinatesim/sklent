import { API_BASE, primeCsrf } from "@/shared/helpers/api/client";

import type { Coupon, Promotion } from "@/features/admin/promotions/types/promotions.types";

const authedFetch = async (path: string, init: RequestInit = {}): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}${path}`, {
    ...init,
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf, ...init.headers },
  });
};

export const fetchPromotions = async (): Promise<Promotion[]> => {
  const res = await fetch(`${API_BASE}/admin/promotions`, { credentials: "include", cache: "no-store" });
  if (!res.ok) throw new Error("failed to load promotions");
  const data = await res.json();
  return data.items ?? [];
};

export const createPromotion = (payload: Record<string, unknown>): Promise<Response> =>
  authedFetch("/admin/promotions", { method: "POST", body: JSON.stringify(payload) });

export const updatePromotion = (id: string, payload: Record<string, unknown>): Promise<Response> =>
  authedFetch(`/admin/promotions/${id}`, { method: "PATCH", body: JSON.stringify(payload) });

export const deletePromotion = (id: string): Promise<Response> =>
  authedFetch(`/admin/promotions/${id}`, { method: "DELETE" });

export const fetchCoupons = async (): Promise<Coupon[]> => {
  const res = await fetch(`${API_BASE}/admin/coupons`, { credentials: "include", cache: "no-store" });
  if (!res.ok) throw new Error("failed to load coupons");
  const data = await res.json();
  return data.items ?? [];
};

export const createCoupon = (payload: Record<string, unknown>): Promise<Response> =>
  authedFetch("/admin/coupons", { method: "POST", body: JSON.stringify(payload) });

export const updateCoupon = (id: string, payload: Record<string, unknown>): Promise<Response> =>
  authedFetch(`/admin/coupons/${id}`, { method: "PATCH", body: JSON.stringify(payload) });

export const deleteCoupon = (id: string): Promise<Response> =>
  authedFetch(`/admin/coupons/${id}`, { method: "DELETE" });
