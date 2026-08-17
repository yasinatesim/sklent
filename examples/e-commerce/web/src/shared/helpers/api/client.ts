export const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8100";

export const INTERNAL_API_BASE = process.env.INTERNAL_API_URL ?? API_BASE;

const readCookie = (name: string): string => {
  if (typeof document === "undefined") return "";
  const match = document.cookie.match(new RegExp(`(?:^|; )${name}=([^;]*)`));
  return match ? decodeURIComponent(match[1]) : "";
};

export const primeCsrf = async (): Promise<string> => {
  let csrf = readCookie("csrf_token");
  if (csrf) return csrf;
  await fetch(`${API_BASE}/healthz`, { credentials: "include" });
  csrf = readCookie("csrf_token");
  return csrf;
};
