import { act, renderHook, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { REQUEST_STATUS } from "@/shared/constants/requestStatus";

import { useDiscountRules } from "@/features/admin/_shared/hooks/useDiscountRules";

type Rule = { id: string; active: boolean };

const makeApi = (overrides: Partial<Parameters<typeof useDiscountRules<Rule>>[0]> = {}) => ({
  fetchAll: vi.fn(async (): Promise<Rule[]> => [{ id: "1", active: true }]),
  update: vi.fn(async () => new Response(null, { status: 200 })),
  remove: vi.fn(async () => new Response(null, { status: 200 })),
  ...overrides,
});

describe("useDiscountRules", () => {
  it("loads on mount and lands on SUCCESS", async () => {
    const api = makeApi();
    const { result } = renderHook(() => useDiscountRules(api));
    await waitFor(() => expect(result.current.status).toBe(REQUEST_STATUS.SUCCESS));
    expect(result.current.items).toHaveLength(1);
  });

  it("surfaces a failed load as ERROR instead of an empty list", async () => {
    const api = makeApi({ fetchAll: vi.fn(async () => { throw new Error("down"); }) });
    const { result } = renderHook(() => useDiscountRules(api));
    await waitFor(() => expect(result.current.status).toBe(REQUEST_STATUS.ERROR));
    expect(result.current.items).toEqual([]);
  });

  it("toggling flips active and refetches so the row reflects the server", async () => {
    const api = makeApi();
    const { result } = renderHook(() => useDiscountRules(api));
    await waitFor(() => expect(result.current.status).toBe(REQUEST_STATUS.SUCCESS));

    await act(async () => {
      await result.current.toggle({ id: "1", active: true });
    });
    expect(api.update).toHaveBeenCalledWith("1", { active: false });
    expect(api.fetchAll).toHaveBeenCalledTimes(2);
  });

  it("removing refetches", async () => {
    const api = makeApi();
    const { result } = renderHook(() => useDiscountRules(api));
    await waitFor(() => expect(result.current.status).toBe(REQUEST_STATUS.SUCCESS));

    await act(async () => {
      await result.current.remove("1");
    });
    expect(api.remove).toHaveBeenCalledWith("1");
    expect(api.fetchAll).toHaveBeenCalledTimes(2);
  });

  it("does not refetch when the server rejects the update", async () => {
    const api = makeApi({ update: vi.fn(async () => new Response(null, { status: 500 })) });
    const { result } = renderHook(() => useDiscountRules(api));
    await waitFor(() => expect(result.current.status).toBe(REQUEST_STATUS.SUCCESS));

    await act(async () => {
      await result.current.toggle({ id: "1", active: true });
    });
    expect(api.fetchAll).toHaveBeenCalledTimes(1);
  });
});
