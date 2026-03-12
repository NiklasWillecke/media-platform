// services/frontend/static/js/api.js
// Wird von allen Pages genutzt

const API_BASE = "http://localhost:8080";

const api = {
  async login(email, password) {
    const res = await fetch(`${API_BASE}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
      credentials: "include",
    });

    const data = await res.json();

    if (!res.ok) {
      throw new Error(data.error || "Login fehlgeschlagen");
    }

    return data;
  },

  async register(email, password) {
    const res = await fetch(`${API_BASE}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
      credentials: "include",
    });

    if (!res.ok) {
      const data = await res.json();
      throw new Error(data.error || "Registrierung fehlgeschlagen");
    }
  },

  getToken() {
    return localStorage.getItem("token");
  },

  logout() {
    localStorage.removeItem("token");
    window.location.href = "/login";
  },
};
