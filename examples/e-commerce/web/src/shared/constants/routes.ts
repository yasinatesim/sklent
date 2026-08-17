// URL segments stay locale-neutral English; the [locale] segment carries the language,
// and next-intl translates the labels — never the path.
export const ROUTES = {
  HOME: "/",
  CART: "/cart",
  CHECKOUT: "/checkout",
  CHECKOUT_SUCCESS: "/checkout/success",
  CHECKOUT_ERROR: "/checkout/error",
  CATEGORY: "/category",
  PRODUCT: "/product",
  SEARCH: "/search",
  LOGIN: "/login",
  FORGOT_PASSWORD: "/forgot-password",
  RESET_PASSWORD: "/reset-password",
  ADMIN: "/admin",
} as const;

// The admin root lives in ROUTES.ADMIN; this map holds only its sub-pages.
export const ADMIN_ROUTES = {
  PRODUCTS: "/admin/products",
  ORDERS: "/admin/orders",
  PROMOTIONS: "/admin/promotions",
  COUPONS: "/admin/coupons",
  STOCK_TRACKING: "/admin/stock-tracking",
  REVIEWS: "/admin/reviews",
  RETURNS: "/admin/returns",
  SETTINGS: "/admin/settings",
} as const;

export type Route = (typeof ROUTES)[keyof typeof ROUTES];
export type AdminRoute = (typeof ADMIN_ROUTES)[keyof typeof ADMIN_ROUTES];

export const localePath = (locale: string, path: string): string =>
  path === ROUTES.HOME ? `/${locale}` : `/${locale}${path}`;
