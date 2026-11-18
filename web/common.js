// Helpers shared by the shop page (app.js) and the admin page (admin.js).

// api calls our Go backend and returns the JSON reply as a JavaScript object.
// If Go replies with an error status (400, 404, 409...), it throws an Error
// carrying Go's message from {"error": "..."}.
async function api(path, options = {}) {
  const response = await fetch(path, {
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

const categoryIcons = { 1: "💻", 2: "🔧", 3: "🔌" };
function categoryIcon(categoryId) {
  return categoryIcons[categoryId] || "📦";
}

// toast shows a short message at the bottom of the screen.
function toast(message, kind = "info") {
  const box = document.getElementById("toast");
  box.textContent = message;
  box.className = `toast show ${kind}`;
  clearTimeout(toast.timer);
  toast.timer = setTimeout(() => (box.className = "toast"), 3500);
}
