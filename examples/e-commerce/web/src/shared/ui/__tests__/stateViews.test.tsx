import { render, screen } from "@testing-library/react";
import type { ComponentType } from "react";
import { describe, expect, it, vi } from "vitest";

import { REQUEST_STATUS, type RequestStatus } from "@/shared/constants/requestStatus";

import EmptyView from "@/shared/ui/EmptyView";
import ErrorView from "@/shared/ui/ErrorView";
import LoadingView from "@/shared/ui/LoadingView";
import NullView from "@/shared/ui/NullView";

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => key,
}));

describe("state views", () => {
  it("NullView renders nothing so IDLE costs no markup", () => {
    const { container } = render(<NullView />);
    expect(container).toBeEmptyDOMElement();
  });

  it("LoadingView announces itself to assistive tech", () => {
    render(<LoadingView />);
    expect(screen.getByRole("status")).toBeInTheDocument();
  });

  it("ErrorView is an alert and shows the message it was given", () => {
    render(<ErrorView message="boom" />);
    const alert = screen.getByRole("alert");
    expect(alert).toHaveTextContent("boom");
  });

  it("ErrorView renders a retry control only when a handler is passed", () => {
    const onRetry = vi.fn();
    const { rerender } = render(<ErrorView />);
    expect(screen.queryByRole("button")).toBeNull();

    rerender(<ErrorView onRetry={onRetry} />);
    screen.getByRole("button").click();
    expect(onRetry).toHaveBeenCalledOnce();
  });

  it("EmptyView shows the caller's message", () => {
    render(<EmptyView message="nothing here" />);
    expect(screen.getByText("nothing here")).toBeInTheDocument();
  });
});

describe("STATE_VIEWS dispatch contract", () => {
  const STATE_VIEWS: Record<RequestStatus, ComponentType<{ message?: string }>> = {
    [REQUEST_STATUS.IDLE]: NullView,
    [REQUEST_STATUS.LOADING]: LoadingView,
    [REQUEST_STATUS.ERROR]: ErrorView,
    [REQUEST_STATUS.SUCCESS]: NullView,
  };

  it("covers every status so no state can fall through unrendered", () => {
    for (const status of Object.values(REQUEST_STATUS)) {
      expect(STATE_VIEWS[status]).toBeTypeOf("function");
    }
  });

  it("stores component references, never elements — the map must not build every branch", () => {
    for (const View of Object.values(STATE_VIEWS)) {
      expect(View).toBeTypeOf("function");
      expect(View).not.toHaveProperty("$$typeof");
    }
  });
});
