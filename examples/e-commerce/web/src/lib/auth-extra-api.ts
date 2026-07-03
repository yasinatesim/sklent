import { API_BASE, primeCsrf } from "./api";

export const forgotPassword = async (email: string): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/auth/password/forgot`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ email }),
  });
};

export const resetPassword = async (token: string, newPassword: string): Promise<Response> => {
  const csrf = await primeCsrf();
  return fetch(`${API_BASE}/auth/password/reset`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json", "X-CSRF-Token": csrf },
    body: JSON.stringify({ token, newPassword }),
  });
};
