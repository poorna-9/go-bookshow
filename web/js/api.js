// Central place to change the backend URL once, instead of in every file.
const API_BASE = "http://localhost:8080/api/v1";

// Small wrapper around fetch: adds the base URL, parses JSON,
// and throws a readable error if the backend returns one.
async function apiFetch(path, options = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(data.error || data.message || "Something went wrong");
  }
  return data;
}

// Very small "current user" placeholder until real auth exists.
// Every booking call needs a user_id — we generate one once per
// browser and reuse it, so the same visitor keeps the same identity.
function getUserId() {
  let id = localStorage.getItem("gs_user_id");
  if (!id) {
    id = crypto.randomUUID();
    localStorage.setItem("gs_user_id", id);
  }
  return id;
}