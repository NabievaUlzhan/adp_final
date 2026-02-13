requireRole("admin")

function getProductIdFromQuery() {
  const p = new URLSearchParams(window.location.search)
  return p.get("id") || ""
}

async function createProduct() {
  const info = document.getElementById("createInfo")
  if (info) info.textContent = ""

  const product = {
    Name: document.getElementById("name").value.trim(),
    Category: document.getElementById("category").value.trim(),
    ImageURL: document.getElementById("imageUrl").value.trim(),
    Price: parseFloat(document.getElementById("price").value),
    Stock: parseInt(document.getElementById("stock").value, 10),
    Ingredients: document.getElementById("ingredients").value.trim(),
  }

  if (!product.Name || !product.Category || Number.isNaN(product.Price) || Number.isNaN(product.Stock)) {
    if (info) info.textContent = "Please fill all fields"
    return
  }

  const res = await apiFetch("/api/products", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(product),
  })

  if (!res.ok) {
    if (info) info.textContent = await res.text()
    return
  }

  alert("Product created")
  window.location = "/"
}

async function prefillUpdateForm() {
  const id = getProductIdFromQuery()
  if (!id) return

  const res = await fetch(`/api/products?search=`)
  const products = await res.json()

  const p = products.find((x) => x.ID === id)
  if (!p) return

  document.getElementById("name").value = p.Name || ""
  document.getElementById("category").value = p.Category || ""
  document.getElementById("imageUrl").value = p.ImageURL || ""
  document.getElementById("price").value = (p.Price ?? "")
  document.getElementById("stock").value = (p.Stock ?? "")
  document.getElementById("ingredients").value = p.Ingredients || ""
}

async function updateProduct() {
  const info = document.getElementById("updateInfo")
  if (info) info.textContent = ""

  const id = getProductIdFromQuery()
  if (!id) {
    if (info) info.textContent = "Missing product id"
    return
  }

  const product = {
    Name: document.getElementById("name").value.trim(),
    Category: document.getElementById("category").value.trim(),
    ImageURL: document.getElementById("imageUrl").value.trim(),
    Price: parseFloat(document.getElementById("price").value),
    Stock: parseInt(document.getElementById("stock").value, 10),
    Ingredients: document.getElementById("ingredients").value.trim(),
  }

  if (!product.Name || !product.Category || Number.isNaN(product.Price) || Number.isNaN(product.Stock)) {
    if (info) info.textContent = "Please fill all fields"
    return
  }

  const res = await apiFetch(`/api/products/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(product),
  })

  if (!res.ok) {
    if (info) info.textContent = await res.text()
    return
  }

  alert("Product updated")
  window.location = "/"
}

if (window.location.pathname === "/update-product") {
  prefillUpdateForm()
}

renderHeader()
