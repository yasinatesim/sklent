import type { AdminOrder } from "@/shared/types/adminOrder.types";

export type OrderDraft = {
  status: AdminOrder["status"];
  trackingNumber: string;
};
