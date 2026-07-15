function renderNavbar() {
  const city = localStorage.getItem("gs_city") || "";
  const nav = document.createElement("nav");
  nav.className = "gs-nav";
  nav.innerHTML = `
    <a href="index.html" class="brand">GoShow</a>
    ${city ? `<span class="city-pill">📍 ${city}</span>` : ""}
  `;
  document.body.prepend(nav);
}

document.addEventListener("DOMContentLoaded", renderNavbar);