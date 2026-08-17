import { API_BASE, primeCsrf } from "@/shared/helpers/api/client";

export const createProduct = async (payload: Record<string, unknown>): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/products`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify(payload),
  });
};
