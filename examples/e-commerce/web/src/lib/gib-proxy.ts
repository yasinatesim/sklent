const ALLOWED_HOSTS = new Set([
  "earsivportal.efatura.gov.tr",
  "earsivportaltest.efatura.gov.tr",
]);

const PROXY_ENDPOINT = "/api/gib-proxy";

export const isGibRequest = (target: string): boolean => {
  try {
    const parsed = new URL(target);
    return parsed.protocol === "https:" && ALLOWED_HOSTS.has(parsed.hostname);
  } catch {
    return false;
  }
};

// installGibProxy temporarily wraps window.fetch so only GIB e-Arşiv calls route through our proxy.
export const installGibProxy = (): (() => void) => {
  const original = window.fetch;
  window.fetch = (input: RequestInfo | URL, init?: RequestInit) => {
    const target = typeof input === "string" ? input : input.toString();
    if (isGibRequest(target)) {
      const proxied = `${PROXY_ENDPOINT}?target=${encodeURIComponent(target)}`;
      return original(proxied, { ...init, credentials: "include" });
    }
    return original(input, init);
  };
  return () => {
    window.fetch = original;
  };
};
