// Helpers shared by the shop page (app.js) and the admin page (admin.js).

// api calls our Go backend and returns the JSON reply as a JavaScript object.
// If Go replies with an error status (400, 404, 409...), it throws an Error
// carrying Go's message from {"error": "..."}.
async function api(path, options = {}) {
  // location.origin is "http://localhost:8080" without any login in it,
  // so this also works if the page was opened as http://user:pass@host/...
  const response = await fetch(new URL(path, location.origin), {
    ...options,
    headers: { "Content-Type": "application/json" },
  });
  if (response.status === 204) return null; // e.g. DELETE: success, no body
  const data = await response.json();
  if (!response.ok) {
    throw new Error(data.error || `Request failed (${response.status})`);
  }
  return data;
}

// Prices travel as paise (integers); show them as rupees: 12999900 → ₹1,29,999
const rupees = new Intl.NumberFormat("en-IN", {
  style: "currency",
  currency: "INR",
  minimumFractionDigits: 0,
  maximumFractionDigits: 2,
});
function formatPrice(paise) {
  return rupees.format(paise / 100);
}

// escapeHTML makes text safe to put inside HTML, so a product named
// "<script>" is shown as text instead of being run.
function escapeHTML(value) {
  return String(value ?? "")
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

const categoryIcons = { 1: "💻", 2: "📱", 3: "🔧", 4: "🔌" };
function categoryIcon(categoryId) {
  return categoryIcons[categoryId] || "📦";
}

// productImage returns the product's photo, or the category icon if the
// product has no image (or the image fails to load).
function productImage(p, className) {
  const icon = categoryIcon(p.categoryId);
  if (!p.image) return `<span class="${className} placeholder">${icon}</span>`;
  return `<img class="${className}" src="${escapeHTML(p.image)}" alt="${escapeHTML(p.name)}"
    loading="lazy" data-icon="${icon}">`;
}

// Tall (portrait) photos, like a phone standing up, are shown whole instead
// of being cropped to the middle. Load events don't bubble either.
document.addEventListener(
  "load",
  (event) => {
    const img = event.target;
    if (img.tagName === "IMG" && img.naturalHeight > img.naturalWidth * 1.15) {
      img.classList.add("tall");
    }
  },
  true
);

// If a product photo fails to load, swap in the category icon.
// Image errors don't bubble up, so this listens in the capture phase (true).
document.addEventListener(
  "error",
  (event) => {
    const img = event.target;
    if (img.tagName !== "IMG" || !img.dataset.icon) return;
    const span = document.createElement("span");
    span.className = `${img.className} placeholder`;
    span.textContent = img.dataset.icon;
    img.replaceWith(span);
  },
  true
);

// toast shows a short message at the bottom of the screen.
function toast(message, kind = "info") {
  const box = document.getElementById("toast");
  box.textContent = message;
  box.className = `toast show ${kind}`;
  clearTimeout(toast.timer);
  toast.timer = setTimeout(() => (box.className = "toast"), 3500);
}
