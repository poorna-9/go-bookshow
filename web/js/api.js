const API_BASE = "http://localhost:8080/api/v1";

function getToken() {
  return localStorage.getItem("gs_token");
}

function isLoggedIn() {
  return !!getToken();
}

function logout() {
  localStorage.removeItem("gs_token");
  localStorage.removeItem("gs_user");
  window.location.href = "/";
}

async function apiFetch(path, options = {}) {
  const headers = { "Content-Type": "application/json" };
  const token = getToken();
  if (token) headers["Authorization"] = `Bearer ${token}`;

  const res = await fetch(`${API_BASE}${path}`, { headers, ...options });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || data.message || "Something went wrong");
  return data;
}