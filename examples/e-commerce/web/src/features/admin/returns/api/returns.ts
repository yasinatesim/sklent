import { API_BASE, primeCsrf } from "@/shared/helpers/api/client";

import type { ReturnRequest, ReturnStatus } from "@/features/admin/returns/types/returns.types";

export const fetchReturns = async (status?: ReturnStatus): Promise<ReturnRequest[]> => {
  const qs = status ? `?status=${encodeURIComponent(status)}` : "";
  const res = await fetch(`${API_BASE}/admin/returns${qs}`, { credentials: "include" });
  if (!res.ok) throw new Error("failed to load returns");
  const data = (await res.json()) as { returns: ReturnRequest[] | null };
  return data.returns ?? [];
};

export const updateReturnStatus = async (
  id: string,
  status: ReturnStatus,
  adminNote: string,
): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/returns/${encodeURIComponent(id)}`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ status, adminNote }),
  });
};
