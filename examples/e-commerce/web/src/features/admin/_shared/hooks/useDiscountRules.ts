"use client";

import { useCallback, useEffect, useState } from "react";

import { REQUEST_STATUS, type RequestStatus } from "@/shared/constants/requestStatus";

type ActivatableRule = {
  id: string;
  active: boolean;
};

type DiscountRuleApi<T extends ActivatableRule> = {
  fetchAll: () => Promise<T[]>;
  update: (id: string, patch: Partial<T>) => Promise<Response>;
  remove: (id: string) => Promise<Response>;
};

// Promotions and coupons are the same rule with a different identifying field; one hook owns the
// list state so the two screens cannot drift apart.
export const useDiscountRules = <T extends ActivatableRule>(api: DiscountRuleApi<T>) => {
  const [status, setStatus] = useState<RequestStatus>(REQUEST_STATUS.IDLE);
  const [items, setItems] = useState<T[]>([]);

  const load = useCallback(async () => {
    setStatus(REQUEST_STATUS.LOADING);
    try {
      setItems(await api.fetchAll());
      setStatus(REQUEST_STATUS.SUCCESS);
    } catch {
      setItems([]);
      setStatus(REQUEST_STATUS.ERROR);
    }
  }, [api]);

  useEffect(() => {
    load();
  }, [load]);

  const toggle = useCallback(
    async (rule: T) => {
      const res = await api.update(rule.id, { active: !rule.active } as Partial<T>);
      if (!res.ok) return false;
      await load();
      return true;
    },
    [api, load],
  );

  const remove = useCallback(
    async (id: string) => {
      const res = await api.remove(id);
      if (!res.ok) return false;
      await load();
      return true;
    },
    [api, load],
  );

  return { status, items, load, toggle, remove };
};

export default useDiscountRules;
