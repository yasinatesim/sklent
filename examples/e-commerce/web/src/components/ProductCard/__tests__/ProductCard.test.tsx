import { describe, expect, it } from "vitest";
import { formatTRY } from "@/lib/api";

describe("formatTRY", () => {
  it("formats cents as Turkish lira by default", () => {
    const out = formatTRY(129900, "tr");
    expect(out).toContain("1.299");
    expect(out).toContain("₺");
  });

  it("formats with the en locale", () => {
    const out = formatTRY(249900, "en");
    expect(out).toContain("2,499");
  });
});
