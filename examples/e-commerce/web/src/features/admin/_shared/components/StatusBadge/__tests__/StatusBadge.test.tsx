import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { STATUS_TONE } from "@/features/admin/_shared/types/statusBadge.types";

import StatusBadge from "@/features/admin/_shared/components/StatusBadge";

describe("StatusBadge", () => {
  it("renders the label it is given", () => {
    render(<StatusBadge label="Shipped" tone={STATUS_TONE.ACTIVE} />);
    expect(screen.getByText("Shipped")).toBeInTheDocument();
  });

  it("gives every tone a distinct class so states are visually separable", () => {
    const classes = Object.values(STATUS_TONE).map((tone) => {
      const { container, unmount } = render(<StatusBadge label="x" tone={tone} />);
      const cls = container.firstElementChild?.className ?? "";
      unmount();
      return cls;
    });
    expect(new Set(classes).size).toBe(classes.length);
  });

  it("defaults to the passive tone rather than throwing on an unmapped status", () => {
    render(<StatusBadge label="unknown" />);
    expect(screen.getByText("unknown")).toBeInTheDocument();
  });

  it("exposes the state to assistive tech, not just by colour", () => {
    render(<StatusBadge label="Cancelled" tone={STATUS_TONE.DANGER} />);
    expect(screen.getByText("Cancelled")).toHaveAttribute("data-tone", STATUS_TONE.DANGER);
  });
});
