function getAuth() {
  return {
    userID: localStorage.getItem("userID") || "",
    role: localStorage.getItem("role") || "",
  }
}

function setAuth(user) {
  localStorage.setItem("userID", user.ID || "")
  localStorage.setItem("role", user.Role || "")
}

function clearAuth() {
  localStorage.removeItem("userID")
  localStorage.removeItem("role")
}

async function apiFetch(path, options = {}) {
  const { userID, role } = getAuth()

  const headers = new Headers(options.headers || {})
  if (userID) headers.set("UserID", userID)
  if (role) headers.set("Role", role)

  return fetch(path, {
    ...options,
    headers,
    credentials: "include", // ✅ cookies тоже отправляем
  })
}

function renderHeader() {
  const el = document.getElementById("appHeader")
  if (!el) return

  const { role } = getAuth()
  const loggedIn = Boolean(role)

  const right = []
  if (!loggedIn) {
    right.push(`<a class="btn" href="/login">Login</a>`)
    right.push(`<a class="btn" href="/register">Register</a>`)
  } else {
    right.push(`<a class="btn" href="/cart">Cart</a>`)
    right.push(`<a class="btn" href="/orders">Orders</a>`)
    if (role === "admin") {
      right.push(`<a class="btn" href="/create-product">Create Product</a>`)
    }
    right.push(`<button class="btn" id="logoutBtn">Logout</button>`)
  }

  el.innerHTML = `
    <div class="container row" style="justify-content:space-between;gap:12px">
      <a class="row" href="/" style="gap:10px">
        <span class="logo">🥐 BakeryStore</span>
        ${loggedIn ? `<span class="badge">${role}</span>` : ``}
      </a>
      <nav class="nav">${right.join("")}</nav>
    </div>
  `

  const logout = document.getElementById("logoutBtn")
  if (logout) {
    logout.onclick = () => {
      clearAuth()
      window.location = "/"
    }
  }

  const footerBtn = document.getElementById("footerLoginBtn")
  if (footerBtn) footerBtn.style.display = loggedIn ? "none" : "inline-flex"
}

function requireLogin() {
  const { role } = getAuth()
  if (!role) window.location = "/login"
}

function requireRole(expected) {
  const { role } = getAuth()
  if (!role) return (window.location = "/login")
  if (role !== expected) return (window.location = "/")
}

renderHeader()
