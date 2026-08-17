import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import type { AdminOrder } from "@/shared/types/adminOrder.types";

import OrderRow from "@/features/admin/orders/components/OrderRow";

vi.mock("next-intl", () => ({ useTranslations: () => (key: string) => key }));

const order: AdminOrder = {
  id: "abcdef1234567890",
  email: "buyer@example.com",
  status: "pending",
  paymentMethod: "card",
  totalCents: 12345,
  items: [],
  createdAt: "2026-01-01T00:00:00Z",
};

const renderRow = (props: Partial<React.ComponentProps<typeof OrderRow>> = {}) => {
  const onSave = vi.fn();
  const onDraftChange = vi.fn();
  const view = render(
    <table>
      <tbody>
        <OrderRow
          order={order}
          draft={{ status: order.status, trackingNumber: "" }}
          onDraftChange={onDraftChange}
          onSave={onSave}
          {...props}
        />
      </tbody>
    </table>,
  );
  return { onSave, onDraftChange, container: view.container };
};

describe("OrderRow", () => {
  it("shortens the order id so the column stays readable", () => {
    renderRow();
    expect(screen.getByText("abcdef12")).toBeInTheDocument();
    expect(screen.queryByText(order.id)).toBeNull();
  });

  it("reports a status change as a draft instead of saving immediately", () => {
    const { onDraftChange, onSave } = renderRow();
    fireEvent.change(screen.getByRole("combobox"), { target: { value: "shipped" } });
    expect(onDraftChange).toHaveBeenCalledWith(order.id, { status: "shipped", trackingNumber: "" });
    expect(onSave).not.toHaveBeenCalled();
  });

  it("saves only when the save control is used", () => {
    const { onSave } = renderRow();
    fireEvent.click(screen.getByRole("button"));
    expect(onSave).toHaveBeenCalledWith(order.id);
  });

  it("renders the status badge from the draft, not the persisted order", () => {
    const { container } = renderRow({ draft: { status: "cancelled", trackingNumber: "" } });
    const badge = container.querySelector("[data-tone]");
    expect(badge).toHaveAttribute("data-tone", "danger");
    expect(badge).toHaveTextContent("cancelled");
  });
});
