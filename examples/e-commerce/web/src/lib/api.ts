export const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8100";
export const INTERNAL_API_BASE = process.env.INTERNAL_API_URL ?? API_BASE;

export type Product = {
  id: string;
  slug: string;
  titleTr: string;
  titleEn: string;
  descriptionTr: string;
  priceCents: number;
  oldPriceCents: number;
  stock: number;
  categorySlug: string;
  badge: string;
  seller: string;
  imageUrl: string;
  published: boolean;
};

export type Category = {
  id: string;
  slug: string;
  nameTr: string;
  nameEn: string;
  icon: string;
  descTr: string;
  descEn: string;
};

export const formatTRY = (cents: number, locale = "tr"): string =>
  new Intl.NumberFormat(locale === "en" ? "en-US" : "tr-TR", {
    style: "currency",
    currency: "TRY",
    maximumFractionDigits: 0,
  }).format(cents / 100);

const readCookie = (name: string): string => {
  if (typeof document === "undefined") return "";
  const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : "";
};

const primeCsrf = async (): Promise<string> => {
  let csrf = readCookie("csrf_token");
  if (csrf) return csrf;
  await fetch(`${API_BASE}/healthz`, { credentials: "include" });
  csrf = readCookie("csrf_token");
  return csrf;
};

export const fetchProducts = async (categorySlug?: string): Promise<Product[]> => {
  const query = categorySlug && categorySlug !== "all" ? `?category=${categorySlug}` : "";
  const res = await fetch(`${INTERNAL_API_BASE}/products${query}`, { cache: "no-store" });
  if (!res.ok) throw new Error("failed to load products");
  const data = await res.json();
  return data.items ?? [];
};

export const fetchProduct = async (slug: string): Promise<Product | null> => {
  const res = await fetch(`${INTERNAL_API_BASE}/products/${slug}`, { cache: "no-store" });
  if (res.status === 404) return null;
  if (!res.ok) throw new Error("failed to load product");
  return res.json();
};

export const fetchCategories = async (): Promise<Category[]> => {
  const res = await fetch(`${INTERNAL_API_BASE}/categories`, { cache: "no-store" });
  if (!res.ok) return [];
  const data = await res.json();
  return data.items ?? [];
};

export const login = async (email: string, password: string): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/auth/login`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ email, password }),
  });
};

export const register = async (
  email: string,
  password: string,
  fullName: string,
): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/auth/register`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ email, password, fullName }),
  });
};

export type PlaceOrderItem = {
  productId: string;
  titleTr: string;
  unitCents: number;
  quantity: number;
};

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

export const createProduct = async (payload: Record<string, unknown>): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/admin/products`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify(payload),
  });
};
