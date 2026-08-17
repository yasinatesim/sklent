import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import CatalogFilters from "@/features/web/catalog/components/CatalogFilters";

vi.mock("next-intl", () => ({ useTranslations: () => (key: string) => key }));

const facets = {
  categories: [
    { slug: "rings", name: "rings", count: 4 },
    { slug: "bracelets", name: "bracelets", count: 2 },
  ],
  minPriceCents: 1000,
  maxPriceCents: 90000,
};

const renderFilters = (query = {}) => {
  const onChange = vi.fn();
  render(<CatalogFilters facets={facets} query={query} onChange={onChange} />);
  return { onChange };
};

describe("CatalogFilters", () => {
  it("shows each category with its result count", () => {
    renderFilters();
    expect(screen.getByText("rings")).toBeInTheDocument();
    expect(screen.getByText("4")).toBeInTheDocument();
  });

  it("selecting a category resets the page — page 3 of the old filter does not exist", () => {
    const { onChange } = renderFilters({ page: 3, category: "bracelets" });
    fireEvent.click(screen.getByText("rings"));
    expect(onChange).toHaveBeenCalledWith({ category: "rings", page: 1 });
  });

  it("clicking the selected category clears it rather than reselecting it", () => {
    const { onChange } = renderFilters({ category: "rings" });
    fireEvent.click(screen.getByText("rings"));
    expect(onChange).toHaveBeenCalledWith({ category: undefined, page: 1 });
  });

  it("toggling in-stock reports a boolean and resets the page", () => {
    const { onChange } = renderFilters({ page: 2 });
    fireEvent.click(screen.getByRole("checkbox"));
    expect(onChange).toHaveBeenCalledWith({ inStock: true, page: 1 });
  });

  it("seeds the price inputs from the facet bounds so the range is discoverable", () => {
    renderFilters();
    const [min, max] = screen.getAllByRole("spinbutton") as HTMLInputElement[];
    expect(min.placeholder).toBe("10");
    expect(max.placeholder).toBe("900");
  });
});
