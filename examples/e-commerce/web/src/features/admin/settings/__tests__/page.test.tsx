import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import AdminSettingsPage from "@/features/admin/settings/page";

import {
  createOrderNotifyRule,
  deleteOrderNotifyRule,
  fetchOrderNotifyRules,
  updateOrderNotifyRule,
} from "@/features/admin/settings/api/orderNotifyRules";

vi.mock("next-intl", () => ({ useTranslations: () => (key: string) => key }));

vi.mock("@/features/admin/settings/api/orderNotifyRules", () => ({
  fetchOrderNotifyRules: vi.fn(),
  createOrderNotifyRule: vi.fn(),
  updateOrderNotifyRule: vi.fn(),
  deleteOrderNotifyRule: vi.fn(),
}));

const RULE = { id: "rule-1", source: "HB", recipient: "ops@vela.test", enabled: true };

describe("AdminSettingsPage", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(fetchOrderNotifyRules).mockResolvedValue([]);
  });

  it("lists a rule with its channel label, not the raw code", async () => {
    vi.mocked(fetchOrderNotifyRules).mockResolvedValue([RULE]);

    render(<AdminSettingsPage />);

    expect(await screen.findByText("ops@vela.test")).toBeInTheDocument();
    expect(screen.getByRole("cell", { name: "Hepsiburada" })).toBeInTheDocument();
  });

  it("creates a rule for the selected channel and reloads", async () => {
    vi.mocked(createOrderNotifyRule).mockResolvedValue({ ok: true } as Response);

    render(<AdminSettingsPage />);
    await waitFor(() => expect(fetchOrderNotifyRules).toHaveBeenCalledTimes(1));

    fireEvent.change(screen.getByRole("combobox"), { target: { value: "TY" } });
    fireEvent.change(screen.getByRole("textbox"), { target: { value: "ops@vela.test" } });
    fireEvent.submit(screen.getByRole("button", { name: /addRow/ }));

    await waitFor(() => expect(createOrderNotifyRule).toHaveBeenCalledWith("TY", "ops@vela.test"));
    await waitFor(() => expect(fetchOrderNotifyRules).toHaveBeenCalledTimes(2));
  });

  it("defaults a new rule to the site channel", async () => {
    render(<AdminSettingsPage />);
    await waitFor(() => expect(fetchOrderNotifyRules).toHaveBeenCalled());

    expect(screen.getByRole("combobox")).toHaveValue("site");
  });

  it("disables a rule instead of deleting it", async () => {
    vi.mocked(fetchOrderNotifyRules).mockResolvedValue([RULE]);
    vi.mocked(updateOrderNotifyRule).mockResolvedValue({ ok: true } as Response);

    render(<AdminSettingsPage />);
    fireEvent.click(await screen.findByRole("button", { name: "disable" }));

    await waitFor(() => expect(updateOrderNotifyRule).toHaveBeenCalledWith("rule-1", false));
    expect(deleteOrderNotifyRule).not.toHaveBeenCalled();
  });

  it("shows the error view when the rules cannot be loaded", async () => {
    vi.mocked(fetchOrderNotifyRules).mockRejectedValue(new Error("boom"));

    render(<AdminSettingsPage />);

    expect(await screen.findByRole("button", { name: /retry/i })).toBeInTheDocument();
  });
});
