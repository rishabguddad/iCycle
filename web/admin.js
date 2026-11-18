// The admin page. It uses the admin and order APIs:
//   GET    /api/products, /api/categories     → the product table
//   POST   /api/products                      → add
//   PUT    /api/products/{id}                 → edit (details only, never stock)
//   POST   /api/products/{id}/restock         → add stock
//   DELETE /api/products/{id}                 → remove
//   GET    /api/orders                        → the order table
//   POST   /api/orders/{id}/cancel            → cancel (Go puts the stock back)

let categories = [];
let products = [];
let editingId = 0; // 0 = adding a new product

async function loadAll() {
  try {
    categories = await api("/api/categories");
    document.getElementById("category-select").innerHTML = categories
      .map((c) => `<option value="${c.id}">${escapeHTML(c.name)}</option>`)
      .join("");
    await Promise.all([loadProducts(), loadOrders()]);
  } catch (error) {
    toast(`Could not reach the server: ${error.message}`, "error");
  }
}

// ---------- Products ----------

async function loadProducts() {
  products = await api("/api/products");
  document.getElementById("product-count").textContent = `${products.length} products`;
  document.getElementById("product-rows").innerHTML = products
    .map((p) => {
      const category = categories.find((c) => c.id === p.categoryId);
      return `
        <tr>
          <td>${p.id}</td>
          <td>${escapeHTML(p.name)}</td>
          <td>${escapeHTML(category ? category.name : p.categoryId)}</td>
          <td><span class="badge ${escapeHTML(p.condition)}">${escapeHTML(p.condition)}</span></td>
          <td class="num">${formatPrice(p.price)}</td>
          <td class="num ${p.stock === 0 ? "danger" : ""}">${p.stock}</td>
          <td><div class="actions">
            <form class="restock" data-restock="${p.id}">
              <input name="quantity" type="number" min="1" value="5" aria-label="Units to add">
              <button type="submit">Restock</button>
            </form>
            <button type="button" data-edit="${p.id}">Edit</button>
            <button type="button" class="danger" data-delete="${p.id}">Delete</button>
          </div></td>
        </tr>`;
    })
    .join("");
}

// "Processor: M2" lines → { Processor: "M2" }
function parseSpecs(text) {
  const specs = {};
  text.split("\n").forEach((line) => {
    const at = line.indexOf(":");
    if (at > 0) specs[line.slice(0, at).trim()] = line.slice(at + 1).trim();
  });
  return specs;
}

function lines(text) {
  return text.split("\n").map((l) => l.trim()).filter(Boolean);
}

async function saveProduct(event) {
  event.preventDefault();
  const form = event.target;
  const value = (name) => form.elements[name].value;

  const details = {
    name: value("name").trim(),
    description: value("description").trim(),
    categoryId: Number(value("categoryId")),
    condition: value("condition"),
    price: Math.round(Number(value("price")) * 100), // rupees → paise
    specs: parseSpecs(value("specs")),
    compatibleModels: lines(value("compatibleModels")),
  };

  try {
    if (editingId) {
      // PUT sends details only: the API rejects a "stock" field.
      await api(`/api/products/${editingId}`, { method: "PUT", body: JSON.stringify(details) });
      toast("Product updated", "ok");
    } else {
      const body = { ...details, stock: Number(value("stock")) };
      const created = await api("/api/products", { method: "POST", body: JSON.stringify(body) });
      toast(`Added product #${created.id}`, "ok");
    }
    stopEditing();
    await loadProducts();
  } catch (error) {
    toast(error.message, "error");
  }
}

function startEditing(id) {
  const p = products.find((x) => x.id === id);
  if (!p) return;
  editingId = id;
  const form = document.getElementById("product-form");
  form.elements.name.value = p.name;
  form.elements.description.value = p.description;
  form.elements.categoryId.value = p.categoryId;
  form.elements.condition.value = p.condition;
  form.elements.price.value = p.price / 100;
  form.elements.specs.value = Object.entries(p.specs || {})
    .map(([k, v]) => `${k}: ${v}`)
    .join("\n");
  form.elements.compatibleModels.value = (p.compatibleModels || []).join("\n");

  document.getElementById("form-title").textContent = `Edit product #${id}`;
  document.getElementById("form-submit").textContent = "Save changes";
  document.getElementById("stock-field").hidden = true; // stock is never edited directly
  form.elements.stock.disabled = true;
  document.getElementById("cancel-edit").hidden = false;
  form.scrollIntoView({ behavior: "smooth" });
}

function stopEditing() {
  editingId = 0;
  const form = document.getElementById("product-form");
  form.reset();
  document.getElementById("form-title").textContent = "Add a product";
  document.getElementById("form-submit").textContent = "Add product";
  document.getElementById("stock-field").hidden = false;
  form.elements.stock.disabled = false;
  document.getElementById("cancel-edit").hidden = true;
}

// ---------- Orders ----------

async function loadOrders() {
  const orders = await api("/api/orders");
  const rows = orders
    .slice()
    .reverse() // newest first
    .map(
      (o) => `
      <tr>
        <td>${o.id}</td>
        <td>${escapeHTML(o.customer.name)}<br><span class="muted">${escapeHTML(o.customer.address.city)}</span></td>
        <td>${o.items.map((i) => `${i.quantity} × ${escapeHTML(i.name)}`).join("<br>")}</td>
        <td class="num">${formatPrice(o.total)}</td>
        <td><span class="badge status-${escapeHTML(o.status)}">${escapeHTML(o.status)}</span></td>
        <td class="muted">${new Date(o.createdAt).toLocaleString("en-IN")}</td>
        <td><div class="actions">${
          o.status === "placed" ? `<button type="button" class="danger" data-cancel="${o.id}">Cancel</button>` : ""
        }</div></td>
      </tr>`
    )
    .join("");
  document.getElementById("order-rows").innerHTML =
    rows || `<tr><td colspan="7" class="muted empty">No orders yet. Place one from the shop.</td></tr>`;
}

// ---------- Events ----------

document.addEventListener("click", async (event) => {
  const target = event.target.closest("[data-edit], [data-delete], [data-cancel]");
  if (!target) return;
  try {
    if (target.dataset.edit) {
      startEditing(Number(target.dataset.edit));
    } else if (target.dataset.delete) {
      if (!confirm(`Delete product #${target.dataset.delete}?`)) return;
      await api(`/api/products/${target.dataset.delete}`, { method: "DELETE" });
      toast("Product deleted", "ok");
      await loadProducts();
    } else if (target.dataset.cancel) {
      await api(`/api/orders/${target.dataset.cancel}/cancel`, { method: "POST" });
      toast(`Order #${target.dataset.cancel} cancelled, stock returned`, "ok");
      await Promise.all([loadOrders(), loadProducts()]);
    }
  } catch (error) {
    toast(error.message, "error");
  }
});

document.addEventListener("submit", async (event) => {
  const form = event.target.closest("[data-restock]");
  if (!form) return;
  event.preventDefault();
  try {
    const quantity = Number(form.elements.quantity.value);
    const p = await api(`/api/products/${form.dataset.restock}/restock`, {
      method: "POST",
      body: JSON.stringify({ quantity }),
    });
    toast(`${p.name}: stock is now ${p.stock}`, "ok");
    await loadProducts();
  } catch (error) {
    toast(error.message, "error");
  }
});

document.getElementById("product-form").addEventListener("submit", saveProduct);
document.getElementById("cancel-edit").addEventListener("click", stopEditing);
document.getElementById("reload-orders").addEventListener("click", () => loadOrders().catch((e) => toast(e.message, "error")));

loadAll();
