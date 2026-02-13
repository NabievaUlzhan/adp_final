function parseDate(d) {
  try {
    return new Date(d)
  } catch {
    return new Date(0)
  }
}
function normId(x) {
  if (!x) return "";
  if (typeof x === "string") {
    const m = x.match(/[a-f0-9]{24}/i);
    return m ? m[0] : x;
  }
  if (typeof x === "object") {
    if (x.$oid) return x.$oid;
  }
  return String(x);
}

async function getProductsMap() {
  const res = await fetch("/api/products")
  const products = await res.json()
  const map = {}
  products.forEach((p) => (map[p.ID] = p))
  return map
}

async function loadOrders() {
  const { role } = getAuth()
  if (!role) return requireLogin()

  const [map, res] = await Promise.all([
    getProductsMap(),
    apiFetch("/api/orders"),
  ])

  const info = document.getElementById("ordersInfo")
  if (info) info.textContent = ""

  if (!res.ok) {
    const t = await res.text()
    if (info) info.textContent = t
    return
  }

  const orders = await res.json()
  const box = document.getElementById("orders")
  box.innerHTML = ""

  if (!orders || orders.length === 0) {
    if (info) info.textContent = "No orders yet."
    return
  }

  const now = Date.now()

  orders.forEach((o) => {
    const created = parseDate(o.CreatedAt)
    const ms = now - created.getTime()
    const canCancel = ms <= 24 * 60 * 60 * 1000

    const items = o.Items || []
    const itemsHtml = items
    .map((it) => {
    const pid = normId(it.ProductID || it.product_id);
    const qty = it.Quantity ?? it.quantity ?? 0;

    const p = map[pid] || {};
    const name = p.Name || p.name || pid || "Unknown product";

    return `<li>${name} × ${qty}</li>`;
  })
  .join("");



    box.innerHTML += `
      <div class="card" style="margin:0">
        <div class="row" style="justify-content:space-between;gap:12px;flex-wrap:wrap">
          <div>
            <div class="product-title">Order #${o.ID}</div>
            <div class="muted">${created.toLocaleString()}</div>
          </div>
          <div class="row row-gap">
            ${role === "customer" ? `<button class="btn" ${canCancel ? "" : "disabled"} onclick="cancelOrder('${o.ID}')">Cancel order</button>` : ``}
          </div>
        </div>
        <div class="hr"></div>
        <ul style="margin:0;padding-left:18px">${itemsHtml}</ul>
        ${role === "customer" && !canCancel ? `<div class="muted" style="margin-top:10px">Too late to cancel 😈</div>` : ``}
      </div>
    `
  })
}

async function cancelOrder(orderId) {
  const ok = confirm("Cancel this order?")
  if (!ok) return

  const res = await apiFetch(`/api/orders/${orderId}`, { method: "DELETE" })
  if (!res.ok) {
    alert(await res.text())
    return
  }
  loadOrders()
}

window.onload = loadOrders

renderHeader()
