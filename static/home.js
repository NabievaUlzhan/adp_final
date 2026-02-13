function productCard(p, role) {
  const rawPrice = (typeof p.Price === "number" || typeof p.Price === "string") ? p.Price : 0
  const priceNum = typeof rawPrice === "number" ? rawPrice : parseFloat(rawPrice)
  const price = Number.isFinite(priceNum) ? priceNum.toFixed(2) : String(rawPrice)

  const img = p.ImageURL || "/static/images/veggie.svg"

  return `
    <div class="product-card">
      <img class="product-image" src="${img}" alt="${p.Name || "Product"}" />
      <div class="product-body">
        <div>
          <div class="product-title-lg">${p.Name || "Unnamed product"}</div>
          <div class="product-meta" style="margin-top:8px">
            <span class="badge">${p.Category || "other"}</span>
            ${p.Stock != null ? `<span class="badge">In stock: ${p.Stock}</span>` : ``}
          </div>
        </div>

        <div class="product-actions">
          <div class="price">$${price}</div>
          <div class="row row-gap">
            <button class="btn btn-icon" title="Add to cart"
              onclick="addToCart('${p.ID}')"
              ${role === "customer" ? "" : "disabled"}>🛒</button>

            ${role === "admin" ? `
              <a class="btn" href="/update-product?id=${p.ID}">Update</a>
              <button class="btn" onclick="deleteProduct('${p.ID}')">Delete</button>
            ` : ``}
          </div>
        </div>
      </div>
    </div>
  `
}

function filterCat(cat) {
  const sel = document.getElementById("categorySelect");
  if (sel) sel.value = (cat === "All" ? "" : cat);
  loadProducts();
}

async function loadProducts() {
  const { role } = getAuth()

  const search = (document.getElementById("searchInput")?.value || "").trim()
  const category = document.getElementById("categorySelect")?.value || ""
  const sort = document.getElementById("sortSelect")?.value || ""

  const qs = new URLSearchParams()
  if (search) qs.set("search", search)
  if (category) qs.set("category", category)
  if (sort) qs.set("sort", sort)

  const res = await fetch(`/api/products?${qs.toString()}`)
  const products = await res.json()

  const div = document.getElementById("products")
  if (!div) return

  div.innerHTML = products.map((p) => productCard(p, role)).join("")

  renderTopProducts(products)
}

function renderTopProducts(products) {
  const box = document.getElementById("topProducts")
  if (!box) return

  const top = [...products]
    .filter((p) => (p.Stock ?? 0) > 0)
    .slice(0, 4)

  box.innerHTML = top.map((p) => {
    const rawPrice = (typeof p.Price === "number" || typeof p.Price === "string") ? p.Price : 0
    const priceNum = typeof rawPrice === "number" ? rawPrice : parseFloat(rawPrice)
    const price = Number.isFinite(priceNum) ? priceNum.toFixed(2) : String(rawPrice)
    const img = p.ImageURL || "/static/images/bread.svg"
    return `
      <div class="top-item">
        <img class="top-thumb" src="${img}" alt="${p.Name || "Product"}" />
        <div>
          <div class="top-name">${p.Name || "Product"}</div>
          <div class="top-price">$${price}</div>
        </div>
      </div>
    `
  }).join("")
}

async function addToCart(productId) {
  const { role } = getAuth()
  if (role !== "customer") {
    alert("Login as customer first")
    return
  }

  const res = await apiFetch("/api/cart/add", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ product_id: productId, quantity: 1 }),
  })

  if (!res.ok) {
    alert("Add to cart failed: " + (await res.text()))
    return
  }

  alert("Added to cart")
}


async function deleteProduct(productId) {
  const { role } = getAuth()
  if (role !== "admin") return

  const ok = confirm("Delete this product?")
  if (!ok) return

  const res = await apiFetch(`/api/products/${productId}`, { method: "DELETE" })
  if (!res.ok) {
    alert(await res.text())
    return
  }
  loadProducts()
}

async function loadRecommendations() {
  const { role } = getAuth()
  const recBox = document.getElementById("recommendations")
  const info = document.getElementById("recInfo")
  if (!recBox) return

  const defaults = [
    { Name: "Croissant", Category: "Pastry", Price: 3.5 },
    { Name: "Baguette", Category: "Bread", Price: 1.8 },
    { Name: "Cheesecake", Category: "Cake", Price: 16 },
  ]

  recBox.innerHTML = ""

  const renderDefaults = (message) => {
    if (info) info.textContent = message
    defaults.forEach((p) => {
      recBox.innerHTML += `
        <div class="card" style="margin:0">
          <div class="product-title">${p.Name}</div>
          <div class="product-meta">
            <span class="badge">${p.Category}</span>
            <span class="badge">$${p.Price}</span>
          </div>
        </div>
      `
    })
  }

  if (role !== "customer") {
    renderDefaults("Login as customer to get AI recommendations.")
    return
  }

  const res = await apiFetch("/api/recommendations")
  if (!res.ok) {
    renderDefaults("No AI key / no history yet — showing popular bakery items.")
    return
  }

  const ct = res.headers.get("content-type") || ""
  if (ct.includes("application/json")) {
    const data = await res.json()
    if (data.mode === "default") {
      if (info) info.textContent = "No history yet — showing basic products."
      recBox.innerHTML = ""
      data.items.forEach((name) => {
        recBox.innerHTML += `
          <div class="card" style="margin:0">
            <div class="product-title">${name}</div>
            <div class="muted">Recommended starter</div>
          </div>
        `
      })
      return
    }
  }

  const text = await res.text()
  if (text.includes("PERMISSION_DENIED") || text.includes("unregistered callers") || text.includes("403")) {
    renderDefaults("AI key not configured — showing popular bakery items.")
    return
  }

  if (info) info.textContent = "AI suggestion (Gemini)"
  recBox.innerHTML = `
    <div class="card" style="margin:0;grid-column:1/-1">
      <pre style="white-space:pre-wrap;margin:0">${text}</pre>
    </div>
  `
}

document.querySelectorAll('#catList li')?.forEach((li) => {
  li.addEventListener('click', () => {
    document.querySelectorAll('#catList li').forEach(x => x.classList.remove('active'))
    li.classList.add('active')
    const cat = li.getAttribute('data-cat') || ''
    const sel = document.getElementById('categorySelect')
    if (sel) sel.value = cat
    loadProducts()
  })
})

document.getElementById("applyBtn")?.addEventListener("click", () => loadProducts())

loadRecommendations()
loadProducts()
renderHeader()
