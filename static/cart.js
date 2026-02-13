requireRole("customer")

async function getProductsMap() {
  const res = await fetch("/api/products")
  const products = await res.json()
  const map = {}
  products.forEach((p) => (map[p.ID] = p))
  return map
}

async function loadCart() {
  const info = document.getElementById("cartInfo")
  if (info) info.textContent = ""

  const [map, res] = await Promise.all([
    getProductsMap(),
    apiFetch("/api/cart"),
  ])

  if (!res.ok) {
    const t = await res.text()
    if (info) info.textContent = t
    return
  }

  const cart = await res.json()
  const div = document.getElementById("cartItems")
  div.innerHTML = ""

  const items = cart.items || cart.Items || []
  if (items.length === 0) {
    if (info) info.textContent = "Cart is empty."
    return
  }

  items.forEach((item) => {
    const pid = item.product_id || item.ProductID
    const qty = item.quantity || item.Quantity
    const p = map[pid] || { Name: pid, Price: 0, Category: "?" }
    const img = p.ImageURL || "/static/images/veggie.svg"
    const priceNum = typeof p.Price === "number" ? p.Price : parseFloat(p.Price)
    const price = Number.isFinite(priceNum) ? priceNum.toFixed(2) : String(p.Price ?? "0.00")
    div.innerHTML += `
      <div class="card" style="margin:0">
        <div class="cart-row">
          <img class="cart-thumb" src="${img}" alt="${p.Name}" />
          <div class="row" style="justify-content:space-between;gap:12px;flex-wrap:wrap">
            <div>
              <div class="product-title">${p.Name}</div>
              <div class="product-meta" style="margin-top:6px">
                <span class="badge">${p.Category || "other"}</span>
                <span class="badge">$${price}</span>
              </div>
            </div>
            <div class="row row-gap">
              <button class="btn" onclick="changeQuantity('${pid}', ${qty - 1})">-</button>
              <span class="badge">qty: ${qty}</span>
              <button class="btn" onclick="changeQuantity('${pid}', ${qty + 1})">+</button>
              <button class="btn" onclick="removeFromCart('${pid}')">Delete</button>
            </div>
          </div>
        </div>
      </div>
    `
  })
}

async function changeQuantity(productId, newQty) {
  if (newQty < 1) return

  await apiFetch(`/api/cart/update/${productId}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ quantity: newQty }),
  })

  loadCart()
}

async function removeFromCart(productId) {
  await apiFetch(`/api/cart/remove/${productId}`, { method: "DELETE" })
  loadCart()
}

async function createOrder() {
  const res = await apiFetch("/api/orders", { method: "POST" })
  if (!res.ok) {
    alert(await res.text())
    return
  }
  alert("Order created")
  loadCart()
}

window.onload = loadCart

renderHeader()
