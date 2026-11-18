// The shop page. Every piece of data comes from the Go API:
//   GET  /api/categories                      → the category tabs
//   GET  /api/products?search=&condition=...  → the product grid (all categories)
//   GET  /api/categories/{id}/products?...    → the product grid (one category)
//   GET  /api/products/{id}                   → the details dialog
//   POST /api/orders                          → checkout

const state = {
  categories: [],
  products: new Map(), // id → product, for the cart and details
  filters: { categoryId: 0, search: "", condition: "", inStock: false },
  cart: loadCart(), // [{ productId, quantity }]
};

// ---------- Loading data from the API ----------

async function loadCategories() {
  state.categories = await api("/api/categories");
  renderTabs();
}

async function loadProducts() {
  const { categoryId, search, condition, inStock } = state.filters;

  // Build the query string: ?search=air&condition=used&inStock=true
  const params = new URLSearchParams();
  if (search) params.set("search", search);
  if (condition) params.set("condition", condition);
  if (inStock) params.set("inStock", "true");

  const path = categoryId ? `/api/categories/${categoryId}/products` : "/api/products";
  const query = params.toString();

  setStatus("Loading…");
  try {
    const products = await api(query ? `${path}?${query}` : path);
    products.forEach((p) => state.products.set(p.id, p));
    renderGrid(products);
    setStatus(products.length ? "" : "No products match these filters.");
  } catch (error) {
    setStatus(`Could not load products: ${error.message}`);
  }
}

// refreshAll reloads every product so stock numbers are current
// (after an order, for example).
async function refreshAll() {
  const all = await api("/api/products");
  state.products = new Map(all.map((p) => [p.id, p]));
  await loadProducts();
  renderCart();
}

// ---------- Rendering ----------

function categoryName(id) {
  const category = state.categories.find((c) => c.id === id);
  return category ? category.name : "";
}

function stockText(p) {
  if (!p.inStock) return `<span class="stock out">Out of stock</span>`;
  if (p.stock <= 3) return `<span class="stock low">Only ${p.stock} left</span>`;
  return `<span class="stock">${p.stock} in stock</span>`;
}

function renderTabs() {
  const tabs = [{ id: 0, name: "All" }, ...state.categories];
  document.getElementById("category-tabs").innerHTML = tabs
    .map(
      (c) => `<button type="button" class="tab ${c.id === state.filters.categoryId ? "active" : ""}"
                data-category="${c.id}">${escapeHTML(c.name)}</button>`
    )
    .join("");
}

function renderGrid(products) {
  document.getElementById("product-grid").innerHTML = products
    .map(
      (p) => `
      <article class="card">
        <button type="button" class="card-image" data-view="${p.id}" aria-label="View ${escapeHTML(p.name)}">
          <span>${categoryIcon(p.categoryId)}</span>
        </button>
        <div class="card-body">
          <div class="card-meta">
            <span class="badge ${escapeHTML(p.condition)}">${escapeHTML(p.condition)}</span>
            <span class="muted">${escapeHTML(categoryName(p.categoryId))}</span>
          </div>
          <h3><a href="#" data-view="${p.id}">${escapeHTML(p.name)}</a></h3>
          <p class="price">${formatPrice(p.price)}</p>
          ${stockText(p)}
        </div>
        <button type="button" class="primary" data-add="${p.id}" ${canAdd(p) ? "" : "disabled"}>
          Add to cart
        </button>
      </article>`
    )
    .join("");
}

async function showProduct(id) {
  try {
    const p = await api(`/api/products/${id}`);
    state.products.set(p.id, p);

    const specs = Object.entries(p.specs || {})
      .map(([key, value]) => `<tr><th>${escapeHTML(key)}</th><td>${escapeHTML(value)}</td></tr>`)
      .join("");
    const fits = (p.compatibleModels || []).map((m) => `<li>${escapeHTML(m)}</li>`).join("");

    document.getElementById("product-detail").innerHTML = `
      <div class="dialog-head">
        <span class="badge ${escapeHTML(p.condition)}">${escapeHTML(p.condition)}</span>
        <button class="icon-button" data-close type="button" aria-label="Close">✕</button>
      </div>
      <div class="detail-image">${categoryIcon(p.categoryId)}</div>
      <h2>${escapeHTML(p.name)}</h2>
      <p class="muted">${escapeHTML(categoryName(p.categoryId))}</p>
      <p>${escapeHTML(p.description)}</p>
      ${specs ? `<h3>Specifications</h3><table class="specs">${specs}</table>` : ""}
      ${fits ? `<h3>Compatible with</h3><ul class="fits">${fits}</ul>` : ""}
      <div class="detail-buy">
        <div><p class="price">${formatPrice(p.price)}</p>${stockText(p)}</div>
        <button type="button" class="primary" data-add="${p.id}" ${canAdd(p) ? "" : "disabled"}>Add to cart</button>
      </div>`;
    document.getElementById("product-dialog").showModal();
  } catch (error) {
    toast(error.message, "error");
  }
}

// ---------- Cart (kept in the browser; the server only sees the final order) ----------

function loadCart() {
  try {
    return JSON.parse(localStorage.getItem("cart")) || [];
  } catch {
    return [];
  }
}

function saveCart() {
  try {
    localStorage.setItem("cart", JSON.stringify(state.cart));
  } catch {
    /* storage unavailable: the cart still works until the page reloads */
  }
}

function quantityInCart(productId) {
  const line = state.cart.find((l) => l.productId === productId);
  return line ? line.quantity : 0;
}

function canAdd(p) {
  return p.inStock && quantityInCart(p.id) < p.stock;
}

function addToCart(productId) {
  const p = state.products.get(productId);
  if (!p || !canAdd(p)) return;
  const line = state.cart.find((l) => l.productId === productId);
  if (line) line.quantity++;
  else state.cart.push({ productId, quantity: 1 });
  saveCart();
  renderCart();
  loadProducts(); // re-render so "Add to cart" disables at the stock limit
  toast(`Added ${p.name}`);
}

function changeQuantity(productId, change) {
  const line = state.cart.find((l) => l.productId === productId);
  if (!line) return;
  const p = state.products.get(productId);
  line.quantity += change;
  if (p && line.quantity > p.stock) line.quantity = p.stock;
  if (line.quantity <= 0) state.cart = state.cart.filter((l) => l.productId !== productId);
  saveCart();
  renderCart();
}

function renderCart() {
  // Drop lines for products that no longer exist.
  state.cart = state.cart.filter((l) => state.products.has(l.productId));

  const count = state.cart.reduce((sum, l) => sum + l.quantity, 0);
  document.getElementById("cart-count").textContent = count;

  const box = document.getElementById("cart-items");
  const form = document.getElementById("checkout-form");
  if (state.cart.length === 0) {
    box.innerHTML = `<p class="muted empty">Your cart is empty.</p>`;
    form.hidden = true;
    return;
  }

  let total = 0;
  const rows = state.cart
    .map((line) => {
      const p = state.products.get(line.productId);
      const lineTotal = p.price * line.quantity;
      total += lineTotal;
      return `
        <li class="cart-line">
          <span class="cart-icon">${categoryIcon(p.categoryId)}</span>
          <div class="cart-name">
            <strong>${escapeHTML(p.name)}</strong>
            <span class="muted">${formatPrice(p.price)} each</span>
          </div>
          <div class="stepper">
            <button type="button" data-qty="${p.id}" data-change="-1" aria-label="One less">−</button>
            <span>${line.quantity}</span>
            <button type="button" data-qty="${p.id}" data-change="1" aria-label="One more"
              ${line.quantity >= p.stock ? "disabled" : ""}>+</button>
          </div>
          <strong class="cart-total">${formatPrice(lineTotal)}</strong>
        </li>`;
    })
    .join("");

  box.innerHTML = `<ul class="cart-list">${rows}</ul>
    <p class="cart-sum"><span>Total</span><strong>${formatPrice(total)}</strong></p>`;
  form.hidden = false;
}

// ---------- Checkout: POST /api/orders ----------

async function placeOrder(event) {
  event.preventDefault(); // stop the browser's own form submit; we send JSON ourselves
  const form = event.target;
  const field = (name) => form.elements[name].value.trim();

  const body = {
    customer: {
      name: field("name"),
      phone: field("phone"),
      address: {
        line1: field("line1"),
        city: field("city"),
        state: field("state"),
        pincode: field("pincode"),
      },
    },
    items: state.cart.map((l) => ({ productId: l.productId, quantity: l.quantity })),
  };

  const button = form.querySelector("button[type=submit]");
  button.disabled = true;
  try {
    const order = await api("/api/orders", { method: "POST", body: JSON.stringify(body) });
    state.cart = [];
    saveCart();
    form.reset();
    document.getElementById("order-result").innerHTML = `
      <div class="success">
        <strong>Order #${order.id} placed</strong>
        <p>${order.items.length} item(s), total ${formatPrice(order.total)}. Thank you, ${escapeHTML(order.customer.name)}!</p>
      </div>`;
    await refreshAll(); // stock went down: show the new numbers
  } catch (error) {
    // Go's message, e.g. "not enough stock: only 1 of MacBook Pro ... left"
    toast(error.message, "error");
    await refreshAll();
  } finally {
    button.disabled = false;
  }
}

// ---------- Events ----------

function setStatus(text) {
  document.getElementById("status").textContent = text;
}

document.addEventListener("click", (event) => {
  const target = event.target.closest("[data-view], [data-add], [data-qty], [data-category], [data-close]");
  if (!target) return;
  if (target.dataset.view) {
    event.preventDefault();
    showProduct(Number(target.dataset.view));
  } else if (target.dataset.add) {
    addToCart(Number(target.dataset.add));
    if (target.closest("#product-dialog")) target.closest("dialog").close();
  } else if (target.dataset.qty) {
    changeQuantity(Number(target.dataset.qty), Number(target.dataset.change));
    loadProducts();
  } else if (target.dataset.category !== undefined) {
    state.filters.categoryId = Number(target.dataset.category);
    renderTabs();
    loadProducts();
  } else if (target.dataset.close !== undefined) {
    target.closest("dialog").close();
  }
});

let searchTimer;
document.getElementById("search").addEventListener("input", (event) => {
  clearTimeout(searchTimer); // wait until typing pauses, so we don't call the API on every key
  searchTimer = setTimeout(() => {
    state.filters.search = event.target.value.trim();
    loadProducts();
  }, 250);
});
document.getElementById("condition").addEventListener("change", (event) => {
  state.filters.condition = event.target.value;
  loadProducts();
});
document.getElementById("in-stock").addEventListener("change", (event) => {
  state.filters.inStock = event.target.checked;
  loadProducts();
});
document.getElementById("cart-button").addEventListener("click", () => {
  document.getElementById("order-result").innerHTML = "";
  renderCart();
  document.getElementById("cart-dialog").showModal();
});
document.getElementById("checkout-form").addEventListener("submit", placeOrder);

// Close a dialog when clicking the dark area around it.
document.querySelectorAll("dialog").forEach((dialog) =>
  dialog.addEventListener("click", (event) => {
    if (event.target === dialog) dialog.close();
  })
);

// ---------- Start ----------

(async function start() {
  try {
    await loadCategories();
    await refreshAll();
  } catch (error) {
    setStatus(`Could not reach the server: ${error.message}`);
  }
})();
