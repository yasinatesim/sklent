import { describe, expect, it } from "vitest";

import { ADMIN_ROUTES, ROUTES } from "../routes";

const LOCALE_SEGMENTS = ["en", "tr"];

describe("ROUTES", () => {
  const allPaths = [...Object.values(ROUTES), ...Object.values(ADMIN_ROUTES)];

  it("exposes every public and admin destination as a constant", () => {
    expect(allPaths.length).toBeGreaterThan(0);
    expect(new Set(allPaths).size).toBe(allPaths.length);
  });

  it("starts every path with a single leading slash", () => {
    for (const path of allPaths) {
      expect(path.startsWith("/")).toBe(true);
      expect(path.startsWith("//")).toBe(false);
    }
  });

  it("uses only lowercase ASCII kebab-case segments", () => {
    for (const path of allPaths) {
      for (const segment of path.split("/").filter(Boolean)) {
        expect(segment).toMatch(/^[a-z0-9]+(-[a-z0-9]+)*$/);
      }
    }
  });

  it("keeps route segments locale-neutral — the URL is not a translation surface", () => {
    const turkishSegments = [
      "sepet", "odeme", "basarili", "hata", "kategori", "urun", "arama",
      "giris", "kayit", "siparis", "urunler", "siparisler", "kuponlar",
      "kampanyalar", "yorumlar", "stok-takibi", "sifremi-unuttum", "sifre-sifirla",
    ];
    for (const path of allPaths) {
      for (const segment of path.split("/").filter(Boolean)) {
        expect(turkishSegments).not.toContain(segment);
      }
    }
  });

  it("does not bake the locale into the path — the [locale] segment owns that", () => {
    for (const path of allPaths) {
      const first = path.split("/").filter(Boolean)[0];
      expect(LOCALE_SEGMENTS).not.toContain(first);
    }
  });

  it("nests every admin route under the admin root", () => {
    for (const path of Object.values(ADMIN_ROUTES)) {
      expect(path.startsWith(`${ROUTES.ADMIN}/`)).toBe(true);
    }
  });
});
