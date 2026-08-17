
export type Promotion = {
  id: string;
  name: string;
  discountType: "percent" | "fixed_try";
  discountValue: number;
  scopeType: "all" | "products" | "categories";
  productIds: string[] | null;
  categoryIds: string[] | null;
  minCartCents: number;
  active: boolean;
};

export type Coupon = {
  id: string;
  code: string;
  discountType: "percent" | "fixed_try";
  discountValue: number;
  scopeType: "all" | "products" | "categories";
  productIds: string[] | null;
  categoryIds: string[] | null;
  minCartCents: number;
  active: boolean;
};
