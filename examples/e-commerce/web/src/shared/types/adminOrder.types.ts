
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
