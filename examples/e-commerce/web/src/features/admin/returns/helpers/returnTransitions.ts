import { RETURN_STATUS, type ReturnStatus } from "@/features/admin/returns/types/returns.types";

// Mirrors api/internal/returnreq/statemachine.go — the UI must never offer a move the server
// will reject.
const ALLOWED_TRANSITIONS: Record<ReturnStatus, ReturnStatus[]> = {
  [RETURN_STATUS.REQUESTED]: [RETURN_STATUS.APPROVED, RETURN_STATUS.REJECTED],
  [RETURN_STATUS.APPROVED]: [RETURN_STATUS.REFUNDED],
  [RETURN_STATUS.REJECTED]: [],
  [RETURN_STATUS.REFUNDED]: [],
};

export const nextStatuses = (from: ReturnStatus): ReturnStatus[] => ALLOWED_TRANSITIONS[from] ?? [];

export const isTerminal = (status: ReturnStatus): boolean => nextStatuses(status).length === 0;
