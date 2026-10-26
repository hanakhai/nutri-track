/**
 * NutriTrack - Dashboard Controller
 * Mengelola logika interaktif dashboard, pencatatan makanan, dan target harian
 */

let currentDate = new Date().toISOString().split("T")[0];
let selectedFood = null;
let currentTarget = null;

document.addEventListener("DOMContentLoaded", async () => {
  await initUser();
  initDateControls();
  initSearch();
  initTargetModal();
  initMealForm();
  await loadDashboardData();
});

// 1. Verifikasi User & Profil
async function initUser() {
  try {
    const res = await API.get("/api/auth/me");
    if (res && res.user) {
      document.getElementById("userNameDisplay").innerText = res.user.name;
    }
  } catch (err) {
    // Sesi tidak valid, API helper akan redirect otomatis ke /login
  }

  document.getElementById("btnLogout").addEventListener("click", async () => {
    try {
      await API.post("/api/auth/logout", {});
    } finally {
      window.location.href = "/login";
    }
  });
}

// 2. Kontrol Tanggal & Navigasi Riwayat
function initDateControls() {
  const dateInput = document.getElementById("datePicker");
  dateInput.value = currentDate;

  dateInput.addEventListener("change", (e) => {
    currentDate = e.target.value;
    loadDashboardData();
  });

  document.getElementById("btnPrevDay").addEventListener("click", () => {
    changeDate(-1);
  });

  document.getElementById("btnNextDay").addEventListener("click", () => {
    changeDate(1);
  });

  document.getElementById("btnToday").addEventListener("click", () => {
    currentDate = new Date().toISOString().split("T")[0];
    dateInput.value = currentDate;
    loadDashboardData();
  });
}

function changeDate(days) {
  const d = new Date(currentDate);
  d.setDate(d.getDate() + days);
  currentDate = d.toISOString().split("T")[0];
  document.getElementById("datePicker").value = currentDate;
  loadDashboardData();
}

// 3. Muat Data Ringkasan Harian & Sesi Makan
async function loadDashboardData() {
  const container = document.getElementById("sessionsContainer");
  container.innerHTML = '<div class="empty-state"><p>Memuat catatan makanan...</p></div>';

  try {
    const res = await API.get(`/api/logs/summary?date=${currentDate}`);
    const summary = res.data;
    currentTarget = summary.target;

    renderNutritionCards(summary);
    renderSessions(summary.sessions || []);
  } catch (err) {
    container.innerHTML = `<div class="empty-state"><p style="color: var(--danger);">Gagal memuat data: ${err.message}</p></div>`;
  }
}

// 4. Render Kartu Nutrisi & Progress Bar Target
function renderNutritionCards(summary) {
  const target = summary.target || {
    calories_target: 2000,
    proteins_target: 50,
    carbohydrate_target: 250,
    fat_target: 70,
  };

  const calPct = Math.min(Math.round((summary.total_calories / target.calories_target) * 100), 100);
  const protPct = Math.min(Math.round((summary.total_proteins / target.proteins_target) * 100), 100);
  const carbPct = Math.min(Math.round((summary.total_carbohydrate / target.carbohydrate_target) * 100), 100);
  const fatPct = Math.min(Math.round((summary.total_fat / target.fat_target) * 100), 100);

  // Kalori
  document.getElementById("valCalories").innerText = summary.total_calories;
  document.getElementById("targetCalories").innerText = `/ ${target.calories_target} kkal`;
  document.getElementById("fillCalories").style.width = `${calPct}%`;
  document.getElementById("badgeCalories").innerText = `${calPct}%`;

  // Protein
  document.getElementById("valProtein").innerText = summary.total_proteins + "g";
  document.getElementById("targetProtein").innerText = `/ ${target.proteins_target}g`;
  document.getElementById("fillProtein").style.width = `${protPct}%`;
  document.getElementById("badgeProtein").innerText = `${protPct}%`;

  // Karbohidrat
  document.getElementById("valCarb").innerText = summary.total_carbohydrate + "g";
  document.getElementById("targetCarb").innerText = `/ ${target.carbohydrate_target}g`;
  document.getElementById("fillCarb").style.width = `${carbPct}%`;
  document.getElementById("badgeCarb").innerText = `${carbPct}%`;

  // Lemak
  document.getElementById("valFat").innerText = summary.total_fat + "g";
  document.getElementById("targetFat").innerText = `/ ${target.fat_target}g`;
  document.getElementById("fillFat").style.width = `${fatPct}%`;
  document.getElementById("badgeFat").innerText = `${fatPct}%`;
}

// 5. Render Daftar Sesi Makan
function renderSessions(sessions) {
  const container = document.getElementById("sessionsContainer");

  if (!sessions || sessions.length === 0) {
    container.innerHTML = `
      <div class="empty-state">
        <div class="empty-state-icon">🍽️</div>
        <p style="font-weight: 500;">Belum ada makanan yang dicatat untuk tanggal ini.</p>
        <p style="font-size: 0.85rem; margin-top: 0.25rem;">Gunakan panel di sebelah kanan untuk mencari dan mencatat konsumsi harian Anda.</p>
      </div>
    `;
    return;
  }

  let html = "";
  sessions.forEach((s) => {
    html += `
      <div class="card session-card">
        <div class="session-header">
          <div>
            <span style="font-weight: 600; font-size: 1rem;">${escapeHtml(s.name)}</span>
            <span class="session-time-pill">${escapeHtml(s.meal_time)}</span>
            <div class="session-macros">
              Total: <strong>${s.session_calories} kkal</strong> · P: ${s.session_proteins}g · K: ${s.session_carbohydrate}g · L: ${s.session_fat}g
            </div>
          </div>
          <button class="btn btn-danger-outline" onclick="deleteSession(${s.id})">
            Hapus Sesi
          </button>
        </div>
        <div class="session-items-list">
    `;

    if (s.items && s.items.length > 0) {
      s.items.forEach((item) => {
        html += `
          <div class="session-item">
            <div>
              <div class="item-name">${escapeHtml(item.makanan_name)}</div>
              <div class="item-meta">
                ${item.gram}g (${escapeHtml(item.serving_unit || "Porsi")}) ${item.serving_note ? `· ${escapeHtml(item.serving_note)}` : ""}
                · <span class="pill-cal">${item.calories} kkal</span>
                · <span class="pill-p">P: ${item.proteins}g</span>
                · <span class="pill-c">K: ${item.carbohydrate}g</span>
                · <span class="pill-f">L: ${item.fat}g</span>
              </div>
            </div>
            <button class="btn btn-danger-outline btn-sm" onclick="deleteItem(${s.id}, ${item.id})" title="Hapus makanan ini">
              ✕
            </button>
          </div>
        `;
      });
    }

    html += `
        </div>
      </div>
    `;
  });

  container.innerHTML = html;
}

// 6. Pencarian Makanan Interaktif (Live Search dengan Debounce)
let searchTimeout = null;

function initSearch() {
  const searchInput = document.getElementById("foodSearchInput");
  searchInput.addEventListener("input", (e) => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      searchFoods(e.target.value.trim());
    }, 250);
  });

  // Muat rekomendasi awal
  searchFoods("");
}

async function searchFoods(keyword) {
  const resultsBox = document.getElementById("searchResults");
  resultsBox.innerHTML = '<div style="padding: 1rem; text-align: center; color: var(--text-muted); font-size: 0.85rem;">Mencari makanan...</div>';

  try {
    const res = await API.get(`/api/foods?q=${encodeURIComponent(keyword)}&limit=30`);
    const foods = res.data || [];

    if (foods.length === 0) {
      resultsBox.innerHTML = '<div style="padding: 1.5rem; text-align: center; color: var(--text-muted); font-size: 0.85rem;">Makanan tidak ditemukan.</div>';
      return;
    }

    let html = "";
    foods.forEach((f) => {
      const isSel = selectedFood && selectedFood.id === f.id ? "selected" : "";
      html += `
        <div class="food-row ${isSel}" onclick="selectFood(${f.id}, '${escapeHtml(f.name)}', ${f.calories}, ${f.proteins}, ${f.fat}, ${f.carbohydrate})">
          <div>
            <div style="font-weight: 500; font-size: 0.875rem;">${escapeHtml(f.name)}</div>
            <div class="food-nutri-pills">
              <span class="pill-cal">${f.calories} kkal/100g</span>
              <span class="pill-p">P: ${f.proteins}g</span>
              <span class="pill-c">K: ${f.carbohydrate}g</span>
              <span class="pill-f">L: ${f.fat}g</span>
            </div>
          </div>
          <span style="font-size: 0.8rem; color: var(--accent-primary); font-weight: 500;">Pilih</span>
        </div>
      `;
    });

    resultsBox.innerHTML = html;
  } catch (err) {
    resultsBox.innerHTML = `<div style="padding: 1rem; color: var(--danger); font-size: 0.85rem;">Gagal memuat makanan: ${err.message}</div>`;
  }
}

function selectFood(id, name, cal, prot, fat, carb) {
  selectedFood = { id, name, cal, prot, fat, carb };

  // Update tampilan makanan terpilih
  document.getElementById("selectedFoodBox").style.display = "block";
  document.getElementById("selectedFoodName").innerText = name;
  updateCalculatedNutrition();

  // Scroll form ke view jika di mobile
  document.getElementById("mealForm").scrollIntoView({ behavior: "smooth", block: "nearest" });
}

function updateCalculatedNutrition() {
  if (!selectedFood) return;
  const gram = parseFloat(document.getElementById("gramInput").value) || 0;
  const factor = gram / 100;

  const cal = Math.round(selectedFood.cal * factor * 10) / 10;
  const prot = Math.round(selectedFood.prot * factor * 10) / 10;
  const carb = Math.round(selectedFood.carb * factor * 10) / 10;
  const fat = Math.round(selectedFood.fat * factor * 10) / 10;

  document.getElementById("calcNutriInfo").innerHTML = `
    Estimasi nutrisi: <strong>${cal} kkal</strong> · P: ${prot}g · K: ${carb}g · L: ${fat}g
  `;
}

// 7. Form Pencatatan Makanan
function initMealForm() {
  document.getElementById("gramInput").addEventListener("input", updateCalculatedNutrition);

  document.getElementById("btnCancelSelect").addEventListener("click", () => {
    selectedFood = null;
    document.getElementById("selectedFoodBox").style.display = "none";
  });

  const form = document.getElementById("mealForm");
  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    if (!selectedFood) {
      alert("Silakan pilih makanan terlebih dahulu dari daftar di atas.");
      return;
    }

    const mealTime = document.getElementById("mealTimeSelect").value;
    const gram = parseFloat(document.getElementById("gramInput").value);
    const servingUnit = document.getElementById("servingUnitInput").value.trim() || "Porsi";
    const servingNote = document.getElementById("servingNoteInput").value.trim() || null;

    if (!gram || gram <= 0) {
      alert("Masukkan berat atau porsi makanan yang valid (lebih dari 0 gram).");
      return;
    }

    const btn = document.getElementById("btnSubmitMeal");
    btn.disabled = true;
    btn.innerText = "Menyimpan...";

    try {
      await API.post("/api/logs", {
        meal_date: currentDate,
        meal_time: mealTime,
        items: [
          {
            makanan_id: selectedFood.id,
            gram: gram,
            serving_unit: servingUnit,
            serving_note: servingNote,
          },
        ],
      });

      // Reset pilihan makanan
      selectedFood = null;
      document.getElementById("selectedFoodBox").style.display = "none";
      document.getElementById("gramInput").value = "100";
      document.getElementById("servingNoteInput").value = "";

      // Refresh data
      await loadDashboardData();
    } catch (err) {
      alert("Gagal mencatat makanan: " + err.message);
    } finally {
      btn.disabled = false;
      btn.innerText = "+ Catat ke Menu Hari Ini";
    }
  });
}

// 8. Hapus Sesi dan Hapus Item
async function deleteSession(sessionId) {
  if (!confirm("Apakah Anda yakin ingin menghapus seluruh sesi makan ini?")) return;

  try {
    await API.delete(`/api/logs/${sessionId}`);
    await loadDashboardData();
  } catch (err) {
    alert("Gagal menghapus sesi: " + err.message);
  }
}

async function deleteItem(sessionId, itemId) {
  if (!confirm("Hapus makanan ini dari sesi makan?")) return;

  try {
    await API.delete(`/api/logs/${sessionId}/items/${itemId}`);
    await loadDashboardData();
  } catch (err) {
    alert("Gagal menghapus makanan: " + err.message);
  }
}

// 9. Modal Target Nutrisi
function initTargetModal() {
  const modal = document.getElementById("targetModal");
  const btnOpen = document.getElementById("btnOpenTargetModal");
  const btnClose = document.getElementById("btnCloseTargetModal");
  const form = document.getElementById("targetForm");

  btnOpen.addEventListener("click", () => {
    if (currentTarget) {
      document.getElementById("targetCaloriesInput").value = currentTarget.calories_target;
      document.getElementById("targetProteinInput").value = currentTarget.proteins_target;
      document.getElementById("targetCarbInput").value = currentTarget.carbohydrate_target;
      document.getElementById("targetFatInput").value = currentTarget.fat_target;
    }
    modal.classList.add("active");
  });

  btnClose.addEventListener("click", () => {
    modal.classList.remove("active");
  });

  modal.addEventListener("click", (e) => {
    if (e.target === modal) modal.classList.remove("active");
  });

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    const btn = document.getElementById("btnSaveTarget");
    btn.disabled = true;
    btn.innerText = "Menyimpan...";

    const payload = {
      calories_target: parseFloat(document.getElementById("targetCaloriesInput").value),
      proteins_target: parseFloat(document.getElementById("targetProteinInput").value),
      carbohydrate_target: parseFloat(document.getElementById("targetCarbInput").value),
      fat_target: parseFloat(document.getElementById("targetFatInput").value),
    };

    try {
      const res = await API.put("/api/targets", payload);
      currentTarget = res.data;
      modal.classList.remove("active");
      await loadDashboardData();
    } catch (err) {
      alert("Gagal memperbarui target: " + err.message);
    } finally {
      btn.disabled = false;
      btn.innerText = "Simpan Perubahan";
    }
  });
}

function escapeHtml(str) {
  if (!str) return "";
  return String(str)
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
}

