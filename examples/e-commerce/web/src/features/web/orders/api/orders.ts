import type { PlaceOrderItem } from "@/shared/types/catalog.types";

import { API_BASE, primeCsrf } from "@/shared/helpers/api/client";

export const placeOrder = async (
  email: string,
  items: PlaceOrderItem[],
): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/orders`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ email, paymentMethod: "card", items }),
  });
};
