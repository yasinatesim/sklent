import { API_BASE, primeCsrf } from "@/shared/helpers/api/client";

export const login = async (email: string, password: string): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/auth/login`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ email, password }),
  });
};

export const register = async (
  email: string,
  password: string,
  fullName: string,
): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/auth/register`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ email, password, fullName }),
  });
};
