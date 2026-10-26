/**
 * NutriTrack - API Client Helper
 * Mengelola komunikasi fetch() dengan REST API Backend
 */

const API = {
  async request(endpoint, options = {}) {
    const defaultHeaders = {
      "Content-Type": "application/json",
      "X-Requested-With": "XMLHttpRequest",
    };

    const config = {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
      credentials: "same-origin", // Mengirimkan cookie HttpOnly ke backend
    };

    try {
      const response = await fetch(endpoint, config);
      const data = await response.json().catch(() => ({}));

      if (!response.ok) {
        // Jika sesi berakhir / unauthorized, arahkan ke login jika bukan di halaman auth
        if (response.status === 401) {
          const currentPath = window.location.pathname;
          if (currentPath !== "/login" && currentPath !== "/register") {
            window.location.href = "/login";
            return;
          }
        }
        throw new Error(data.error || "Terjadi kesalahan pada server.");
      }

      return data;
    } catch (err) {
      console.error(`API Error [${endpoint}]:`, err);
      throw err;
    }
  },

  get(endpoint) {
    return this.request(endpoint, { method: "GET" });
  },

  post(endpoint, body) {
    return this.request(endpoint, {
      method: "POST",
      body: JSON.stringify(body),
    });
  },

  put(endpoint, body) {
    return this.request(endpoint, {
      method: "PUT",
      body: JSON.stringify(body),
    });
  },

  delete(endpoint) {
    return this.request(endpoint, { method: "DELETE" });
  },
};

