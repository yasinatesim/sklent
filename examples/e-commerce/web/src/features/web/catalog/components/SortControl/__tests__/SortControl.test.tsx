import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { PRODUCT_SORT } from "@/features/web/catalog/types/catalogQuery.types";

import SortControl from "@/features/web/catalog/components/SortControl";

vi.mock("next-intl", () => ({ useTranslations: () => (key: string) => key }));

describe("SortControl", () => {
  it("offers every supported order", () => {
    render(<SortControl value={PRODUCT_SORT.NEWEST} onChange={vi.fn()} />);
    const options = screen.getAllByRole("option").map((o) => (o as HTMLOptionElement).value);
    expect(options).toEqual(Object.values(PRODUCT_SORT));
  });

  it("reports the chosen order", () => {
    const onChange = vi.fn();
    render(<SortControl value={PRODUCT_SORT.NEWEST} onChange={onChange} />);
    fireEvent.change(screen.getByRole("combobox"), { target: { value: PRODUCT_SORT.PRICE_DESC } });
    expect(onChange).toHaveBeenCalledWith(PRODUCT_SORT.PRICE_DESC);
  });

  it("is labelled, so the control is reachable without sighted context", () => {
    render(<SortControl value={PRODUCT_SORT.NEWEST} onChange={vi.fn()} />);
    expect(screen.getByRole("combobox")).toHaveAccessibleName();
  });
});
