import { API_BASE, primeCsrf } from "@/shared/helpers/api/client";

import type { Review } from "@/features/admin/reviews/types/reviews.types";

export const fetchProductReviews = async (slug: string): Promise<Review[]> => {
  const res = await fetch(`${API_BASE}/products/${slug}/reviews`, { cache: "no-store" });
  if (!res.ok) return [];
  const data = await res.json();
  return data.items ?? [];
};

export const submitReview = async (
  slug: string,
  authorName: string,
  rating: number,
  comment: string,
): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/products/${slug}/reviews`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ authorName, rating, comment }),
  });
};

export const fetchAdminReviews = async (): Promise<Review[]> => {
  const res = await fetch(`${API_BASE}/admin/reviews`, { credentials: "include", cache: "no-store" });
  if (!res.ok) throw new Error("failed to load reviews");
  const data = await res.json();
  return data.items ?? [];
};

export const updateReviewStatus = async (id: string, status: "approved" | "rejected"): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/reviews/${id}/status`, {
    method: "PATCH",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ status }),
  });
};
