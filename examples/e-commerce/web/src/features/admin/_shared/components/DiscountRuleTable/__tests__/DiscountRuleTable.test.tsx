import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import DiscountRuleTable from "@/features/admin/_shared/components/DiscountRuleTable";

vi.mock("next-intl", () => ({ useTranslations: () => (key: string) => key }));

const items = [
  { id: "1", active: true, code: "SUMMER" },
  { id: "2", active: false, code: "WINTER" },
];

const renderTable = () => {
  const onToggle = vi.fn();
  const onRemove = vi.fn();
  render(
    <DiscountRuleTable
      items={items}
      labelHeader="code"
      labelOf={(item) => item.code}
      discountOf={() => 10}
      onToggle={onToggle}
      onRemove={onRemove}
    />,
  );
  return { onToggle, onRemove };
};

describe("DiscountRuleTable", () => {
  it("renders one row per rule with its identifying label", () => {
    renderTable();
    expect(screen.getByText("SUMMER")).toBeInTheDocument();
    expect(screen.getByText("WINTER")).toBeInTheDocument();
  });

  it("shows the active state as a badge, not only as button text", () => {
    const { container } = render(
      <DiscountRuleTable
        items={items}
        labelHeader="code"
        labelOf={(i) => i.code}
        discountOf={() => 10}
        onToggle={vi.fn()}
        onRemove={vi.fn()}
      />,
    );
    const tones = [...container.querySelectorAll("[data-tone]")].map((n) => n.getAttribute("data-tone"));
    expect(tones).toEqual(["active", "passive"]);
  });

  it("hands the whole rule to onToggle so the caller knows the current state", () => {
    const { onToggle } = renderTable();
    fireEvent.click(screen.getAllByText("makePassive")[0]);
    expect(onToggle).toHaveBeenCalledWith(items[0]);
  });

  it("labels the toggle by what it will do, not by the current state", () => {
    renderTable();
    expect(screen.getByText("makePassive")).toBeInTheDocument();
    expect(screen.getByText("makeActive")).toBeInTheDocument();
  });

  it("removes by id", () => {
    const { onRemove } = renderTable();
    fireEvent.click(screen.getAllByText("delete")[1]);
    expect(onRemove).toHaveBeenCalledWith("2");
  });
});
