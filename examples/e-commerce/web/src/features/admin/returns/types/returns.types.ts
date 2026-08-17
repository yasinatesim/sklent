export const RETURN_STATUS = {
  REQUESTED: "requested",
  APPROVED: "approved",
  REJECTED: "rejected",
  REFUNDED: "refunded",
} as const;

export type ReturnStatus = (typeof RETURN_STATUS)[keyof typeof RETURN_STATUS];

export const RETURN_REASON = {
  DAMAGED: "damaged",
  WRONG_ITEM: "wrong_item",
  NOT_AS_SHOWN: "not_as_shown",
  CHANGED_MIND: "changed_mind",
  OTHER: "other",
} as const;

export type ReturnReason = (typeof RETURN_REASON)[keyof typeof RETURN_REASON];

export type ReturnRequest = {
  id: string;
  orderId: string;
  email: string;
  reason: ReturnReason;
  comment?: string;
  status: ReturnStatus;
  adminNote?: string;
  createdAt: string;
};
