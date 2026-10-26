# NutriTrack — Nutrition Tracker

Aplikasi web pemantau nutrisi dan makanan harian dengan antarmuka **Quiet Elegance (Minimalist Nutrition Dashboard)** yang dibangun menggunakan **Go (Gin)** dan **HTML, CSS, JavaScript murni**, terhubung langsung ke **PostgreSQL Supabase**.

---

## 🌿 Fitur Utama

1. **Autentikasi Aman:**
   * Registrasi akun dan login menggunakan password hash `bcrypt`.
   * Sesi dikelola via **JWT dalam HttpOnly Cookie** (terlindungi dari serangan XSS).
2. **Dashboard Nutrisi Harian:**
   * Ringkasan total kalori, protein, karbohidrat, dan lemak per tanggal yang dipilih.
   * Indikator progress bar interaktif terhadap target harian.
3. **Pencarian Katalog Makanan:**
   * Pencarian interaktif dengan *live search* dari database Supabase (100+ jenis makanan).
   * Rincian kalori dan makronutrisi per 100 gram.
4. **Pencatatan Makanan (Meal Log):**
   * Pengelompokan sesi makan berdasarkan waktu: Sarapan, Makan Siang, Makan Malam, Camilan.
   * Perhitungan nutrisi otomatis proporsional terhadap berat (gram) atau porsi makanan.
5. **Target Nutrisi Pribadi:**
   * Setiap pengguna dapat mengatur target kalori dan makronutrisi hariannya sendiri.
6. **Riwayat & Navigasi Tanggal:**
   * Navigasi tanggal yang fleksibel (Hari sebelumnya, Hari berikutnya, atau pemilih kalender).
7. **Isolasi Data Penuh:**
   * Query database terproteksi `WHERE user_id = $1` untuk menjamin privasi catatan antar pengguna.

---

## 📁 Struktur Folder Proyek

```text
go_vibe/
├── .env.example              # Template konfigurasi environment
├── .gitignore                 # Daftar file rahasia yang diabaikan Git
├── go.mod                     # Definisi modul Go
├── go.sum                     # Checksum dependensi Go
├── README.md                  # Dokumentasi panduan
│
├── cmd/
│   └── server/
│       └── main.go            # Titik masuk utama aplikasi (Entry Point)
│
├── internal/
│   ├── config/
│   │   └── config.go          # Membaca konfigurasi .env
│   ├── database/
│   │   └── database.go        # Koneksi pool PostgreSQL Supabase (pgx)
│   ├── handler/
│   │   ├── auth_handler.go    # Handler REST API: Login, Register, Logout, Profile
│   │   ├── food_handler.go    # Handler REST API: Pencarian & Detail Makanan
│   │   ├── log_handler.go     # Handler REST API: Catatan Makan & Ringkasan Harian
│   │   └── target_handler.go  # Handler REST API: Target Nutrisi Pribadi
│   ├── middleware/
│   │   └── auth.go            # Middleware verifikasi token JWT
│   ├── model/
│   │   ├── food.go            # Model data makanan
│   │   ├── log.go             # Model data sesi & item makanan
│   │   ├── target.go          # Model data target nutrisi
│   │   └── user.go            # Model data pengguna
│   └── repository/
│       ├── food_repo.go       # Query SQL manual katalog makanan
│       ├── log_repo.go        # Query SQL manual catatan & ringkasan nutrisi
│       ├── target_repo.go     # Query SQL manual target nutrisi
│       └── user_repo.go       # Query SQL manual pengguna & autentikasi
│
└── web/                       # Frontend HTML, CSS, JavaScript Murni
    ├── css/
    │   └── style.css          # Desain Quiet Elegance Minimalist
    ├── js/
    │   ├── api.js             # API Client fetch helper
    │   └── dashboard.js       # Logika interaktif dashboard
    ├── index.html             # Tampilan dashboard utama
    ├── login.html             # Tampilan halaman masuk
    └── register.html          # Tampilan halaman pendaftaran
```

---

## 🚀 Cara Menjalankan Aplikasi

### 1. Prasyarat
* Telah terinstal **Go** (versi 1.20 atau yang lebih baru).
* File `.env` lokal sudah terisi dengan kredensial PostgreSQL Supabase Anda:
  ```dotenv
  APP_PORT=8080
  DATABASE_URL=postgresql://postgres.xxxx:your_password@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres?sslmode=require
  JWT_SECRET=rahasia-string-acak-minimal-32-karakter
  ENVIRONMENT=development
  ```

### 2. Jalankan Server
Buka terminal PowerShell pada folder `go_vibe`, lalu jalankan perintah:

```powershell
go run ./cmd/server
```

### 3. Buka di Browser
Buka peramban (browser) dan akses:
👉 **`http://localhost:8080`**

* Jika belum login, Anda akan otomatis diarahkan ke halaman login (`/login`).
* Anda dapat mendaftarkan akun baru di menu **Daftar Akun** (`/register`).

