function renderNavbar() {
  const city = localStorage.getItem("gs_city") || "";
  const user = JSON.parse(localStorage.getItem("gs_user") || "null");

  const nav = document.createElement("nav");
  nav.className = "gs-nav";
  nav.innerHTML = `
    <a href="/" class="brand">GoShow</a>
    <div style="display:flex; align-items:center; gap:16px;">
      ${city ? `<span class="city-pill">📍 ${city}</span>` : ""}
      ${user
        ? `<span class="mono" style="font-size:0.85rem;">${user.name}</span>
           <button class="btn btn-outline" onclick="logout()" style="padding:8px 16px;">Log out</button>`
        : `<a href="/login" class="btn btn-marquee" style="padding:8px 16px;">Log in</a>`}
    </div>
  `;
  document.body.prepend(nav);
}
document.addEventListener("DOMContentLoaded", renderNavbar);