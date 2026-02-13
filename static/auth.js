async function login() {
  const email = document.getElementById("email").value.trim()
  const password = document.getElementById("password").value
  const errEl = document.getElementById("loginError")
  if (errEl) errEl.textContent = ""

  if (!email || !password) {
    if (errEl) errEl.textContent = "Email and password are required"
    return
  }

  const res = await fetch("/api/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })

  if (!res.ok) {
    const t = await res.text()
    if (errEl) errEl.textContent = t || "Login failed"
    return
  }

  const user = await res.json()
  setAuth(user)
  window.location = "/"
}

async function register() {
  const name = document.getElementById("name").value.trim()
  const email = document.getElementById("email").value.trim()
  const password = document.getElementById("password").value
  const errEl = document.getElementById("registerError")
  if (errEl) errEl.textContent = ""

  if (!name || !email || !password) {
    if (errEl) errEl.textContent = "Name, email and password are required"
    return
  }

  const res = await fetch("/api/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, email, password }),
  })

  if (!res.ok) {
    const t = await res.text()
    if (errEl) errEl.textContent = t || "Register failed"
    return
  }

  const user = await res.json()
  setAuth(user)
  window.location = "/"
}

renderHeader()
