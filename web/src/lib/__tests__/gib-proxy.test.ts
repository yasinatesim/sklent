import { describe, expect, it } from "vitest";
import { isGibRequest } from "../gib-proxy";

describe("isGibRequest", () => {
  it("allows https GIB hosts", () => {
    expect(isGibRequest("https://earsivportal.efatura.gov.tr/login")).toBe(true);
    expect(isGibRequest("https://earsivportaltest.efatura.gov.tr/x")).toBe(true);
  });

  it("rejects http, other hosts and lookalikes", () => {
    expect(isGibRequest("http://earsivportal.efatura.gov.tr/login")).toBe(false);
    expect(isGibRequest("https://evil.example.com/")).toBe(false);
    expect(isGibRequest("https://earsivportal.efatura.gov.tr.evil.com/")).toBe(false);
    expect(isGibRequest("not-a-url")).toBe(false);
  });
});
