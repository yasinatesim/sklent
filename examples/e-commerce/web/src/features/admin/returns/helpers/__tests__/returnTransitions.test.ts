import { describe, expect, it } from "vitest";

import { RETURN_STATUS } from "@/features/admin/returns/types/returns.types";

import { isTerminal, nextStatuses } from "@/features/admin/returns/helpers/returnTransitions";

describe("returnTransitions", () => {
  it("offers approve or reject on a fresh request", () => {
    expect(nextStatuses(RETURN_STATUS.REQUESTED)).toEqual([
      RETURN_STATUS.APPROVED,
      RETURN_STATUS.REJECTED,
    ]);
  });

  it("offers only refund once approved", () => {
    expect(nextStatuses(RETURN_STATUS.APPROVED)).toEqual([RETURN_STATUS.REFUNDED]);
  });

  it("offers nothing on a terminal status so the UI cannot propose an illegal move", () => {
    expect(nextStatuses(RETURN_STATUS.REJECTED)).toEqual([]);
    expect(nextStatuses(RETURN_STATUS.REFUNDED)).toEqual([]);
  });

  it("mirrors the server's terminal set", () => {
    expect(isTerminal(RETURN_STATUS.REQUESTED)).toBe(false);
    expect(isTerminal(RETURN_STATUS.APPROVED)).toBe(false);
    expect(isTerminal(RETURN_STATUS.REJECTED)).toBe(true);
    expect(isTerminal(RETURN_STATUS.REFUNDED)).toBe(true);
  });
});
