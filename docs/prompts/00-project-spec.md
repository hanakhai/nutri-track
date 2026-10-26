# Project Specification — Vibe Coding

> **Tugas 1:** Spesifikasi Aplikasi Vibe Coding
>
> Tentukan kebutuhan aplikasi, teknologi, desain UI/UX, dan batasan pengembangan sebelum memulai implementasi. Spesifikasi ini menjadi acuan utama pengembangan NutriTrack.

## 1. Informasi Aplikasi

**Nama aplikasi:**  
NutriTrack — Nutrition Tracker

**Deskripsi aplikasi:**  
NutriTrack adalah aplikasi web untuk membantu pengguna mencatat makanan yang dikonsumsi, memantau asupan kalori dan nutrisi harian, serta membandingkan asupan dengan target nutrisi pribadi. Aplikasi ini menyediakan antarmuka yang sederhana, elegan, mudah dipahami, dan nyaman digunakan untuk memantau kebiasaan makan sehari-hari.

**Target pengguna:**  
Individu yang ingin memantau pola makan dan asupan nutrisi harian secara mandiri, termasuk pengguna yang ingin mengatur asupan kalori, protein, karbohidrat, dan lemak.

**Masalah yang ingin diselesaikan:**  
Pengguna sering kesulitan mencatat makanan yang dikonsumsi dan menghitung asupan nutrisi secara teratur. NutriTrack membantu mengorganisasi catatan makanan, menghitung nutrisi berdasarkan data makanan dan jumlah konsumsi, serta menampilkan ringkasan asupan harian agar lebih mudah dipantau.

## 2. Kebutuhan Aplikasi

Fitur utama aplikasi:
1. Registrasi, login, dan logout — pengguna dapat membuat akun, masuk, dan keluar dari aplikasi dengan aman.
2. Dashboard nutrisi — menampilkan ringkasan kalori, protein, karbohidrat, dan lemak yang dikonsumsi pada tanggal yang dipilih, beserta perbandingannya terhadap target harian.
3. Pengelolaan makanan — pengguna dapat mencari dan memilih makanan dari data makanan yang tersedia, melihat informasi nutrisinya, serta memasukkan jumlah konsumsi.
4. Catatan makanan harian — pengguna dapat menambahkan, melihat, mengedit, dan menghapus catatan makanan berdasarkan waktu makan, seperti sarapan, makan siang, makan malam, dan camilan.
5. Perhitungan nutrisi otomatis — sistem menghitung estimasi kalori dan zat gizi berdasarkan data nutrisi makanan serta jumlah atau berat yang dikonsumsi.
6. Target nutrisi pribadi — pengguna dapat melihat dan mengatur target kalori, protein, karbohidrat, dan lemak harian.
7. Riwayat konsumsi — pengguna dapat melihat catatan makanan berdasarkan tanggal dan meninjau ringkasan asupan pada tanggal sebelumnya.
8. Profil pengguna — pengguna dapat melihat dan memperbarui informasi profil yang memang didukung oleh struktur database.

Ketentuan fitur:
- Seluruh data catatan konsumsi harus tersimpan di Supabase melalui backend.
- Data pribadi hanya boleh diakses oleh pengguna yang memiliki data tersebut.
- Perhitungan nutrisi harus menggunakan data yang tersedia di database, bukan angka acak atau data contoh yang tidak ditandai.
- Jika data makanan tidak ditemukan, tampilkan pesan yang jelas.
- Jika target nutrisi belum tersedia, tampilkan status belum diatur dan sediakan cara untuk mengaturnya apabila struktur database mendukung.
- Fitur tambahan seperti grafik, ekspor laporan, pengelolaan resep, dan pelacakan berat badan tidak termasuk implementasi awal kecuali disetujui setelah pemeriksaan database.

**Halaman aplikasi:**
- Login — halaman masuk ke akun.
- Register — halaman pembuatan akun.
- Dashboard — ringkasan asupan nutrisi dan target harian.
- Makanan — halaman untuk mencari dan melihat informasi nutrisi makanan.
- Catatan Makan — halaman untuk menambah, melihat, mengedit, dan menghapus catatan konsumsi.
- Riwayat — halaman untuk melihat catatan konsumsi berdasarkan tanggal.
- Target Nutrisi — halaman untuk melihat dan mengatur target nutrisi.
- Profil — halaman informasi dan pengaturan profil pengguna.
Halaman dan fitur di atas merupakan rancangan awal. Sebelum implementasi, periksa tabel, kolom, relasi, dan data yang sudah tersedia di Supabase. Sesuaikan implementasi dengan database yang benar-benar ada, tanpa membuat asumsi bahwa seluruh tabel atau kolom sudah tersedia.

## 3. Pilihan Teknologi

Pilih satu opsi untuk setiap komponen yang diperlukan. Teknologi dipilih untuk menjaga aplikasi tetap sederhana, mudah dipelajari, dan sesuai dengan database Supabase yang sudah tersedia.

Pilih satu opsi untuk setiap komponen yang diperlukan. Sesuaikan dengan kebutuhan aplikasi. isiskan [x] hanya satu saja

### A. Backend Framework

- [ ] 1. Go Standard Library (`net/http`)
- [x] 2. Gin
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** 2

### B. ORM (Object-Relational Mapping)

- [X] 1. Tidak menggunakan ORM (SQL manual)
- [ ] 2. GORM
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** 1

**Ketentuan:** Gunakan SQL manual dengan library PostgreSQL yang sesuai, misalnya pgx. Pilih versi dependensi yang kompatibel dan jelaskan alasannya sebelum menambahkannya. Gunakan query berparameter untuk mencegah SQL injection.

### C. Database

- [X] 1. PostgreSQL
- [ ] 2. MySQL
- [ ] 3. SQLite

**Pilihan saya:** 1

Ketentuan:
- Gunakan database PostgreSQL yang sudah tersedia di Supabase.
- Jangan membuat database Supabase baru atau menghapus database yang ada.
- Lakukan pemeriksaan struktur tabel, kolom, primary key, foreign key, constraint, index, dan relasi yang relevan sebelum menyusun model dan query.
- Jangan mengubah struktur database secara langsung tanpa menjelaskan perubahan yang diperlukan dan mendapatkan persetujuan.
- Jangan menghapus, menimpa, atau mereset data yang sudah ada.
- Gunakan koneksi melalui backend dengan DATABASE_URL yang dibaca dari environment variable.
- Jangan menghubungkan frontend secara langsung ke database PostgreSQL.
- Jangan mengekspos kredensial database kepada browser.

### D. Frontend

- [x] 1. HTML, CSS, JavaScript murni
- [ ] 2. React
- [ ] 3. Vue

**Pilihan saya:** 1

### E. Sistem Autentikasi

- [ ] 1. Session
- [x] 2. JWT

**Pilihan saya:** 2

**Ketentuan:** 
Gunakan autentikasi JWT yang dikelola backend untuk memverifikasi identitas pengguna pada endpoint yang memerlukan login.
- Password harus di-hash menggunakan algoritma yang sesuai, misalnya Argon2id atau bcrypt.
- Backend harus memvalidasi token dan mengambil identitas pengguna dari token yang telah diverifikasi.
- Jangan mempercayai user_id yang dikirim oleh frontend sebagai bukti kepemilikan data.
- Terapkan middleware autentikasi pada endpoint yang memerlukan login.
- Berikan respons yang sesuai untuk token tidak valid, kedaluwarsa, atau tidak tersedia.
- Tentukan mekanisme penyimpanan dan pengiriman token sebelum implementasi.

### F. Penyimpanan Sesi atau Token

- [x] 1. Cookie
- [ ] 2. localStorage
- [ ] 3. Tidak diperlukan

**Pilihan saya:** 1

**Ketentuan:** 
- Simpan access token JWT dalam cookie yang dibuat oleh backend, bukan melalui localStorage.
- Atur atribut HttpOnly agar token tidak dapat dibaca langsung oleh JavaScript frontend.
- Gunakan atribut Secure pada lingkungan produksi yang menggunakan HTTPS.
- Atur atribut SameSite sesuai arsitektur aplikasi. Gunakan Lax sebagai pilihan awal jika sesuai, dan evaluasi kembali jika frontend dan backend berada pada site yang berbeda.
- Tentukan atribut Path, masa berlaku cookie (Max-Age atau Expires), dan domain seperlunya. Hindari menetapkan domain yang lebih luas daripada kebutuhan aplikasi.
- Backend harus memverifikasi JWT pada setiap request yang memerlukan autentikasi.
- Terapkan perlindungan CSRF yang sesuai untuk endpoint yang mengubah data, seperti menambahkan, mengedit, dan menghapus catatan.
- Jika frontend dan backend menggunakan origin berbeda, atur CORS secara spesifik. Jangan mengizinkan semua origin ketika menggunakan credential.
- Pada logout, backend harus menghapus cookie menggunakan atribut yang konsisten dengan cookie saat login. Jika sistem menggunakan mekanisme pencabutan atau rotasi token, tangani proses tersebut sesuai implementasinya.
- Jangan mencatat token atau kredensial autentikasi ke log maupun menampilkannya pada respons API.
- Jangan menganggap penghapusan cookie di browser otomatis membatalkan JWT yang masih berlaku. Jika pencabutan langsung diperlukan, tentukan mekanisme server-side yang sesuai.
- Sebelum implementasi, jelaskan konfigurasi cookie, alur login dan logout, CORS, serta perlindungan CSRF yang akan digunakan.

### G. API

- [x] 1. REST API
- [ ] 2. Tidak membutuhkan API terpisah

**Pilihan saya:** 1

**Ketentuan:** 
- Backend menyediakan REST API untuk autentikasi, data makanan, catatan konsumsi, target nutrisi, riwayat, dan profil sesuai dukungan database.
- Frontend mengakses API menggunakan `fetch()`.
- Gunakan format JSON untuk request dan response API.
- Gunakan HTTP method dan status code yang sesuai, seperti GET, POST, PUT atau PATCH, dan DELETE.
- Validasi input pada backend meskipun validasi juga tersedia pada frontend.
- Terapkan autentikasi dan pemeriksaan kepemilikan data pada setiap endpoint yang relevan.
- Jangan mengembalikan password hash, kredensial, token rahasia, atau informasi internal yang tidak perlu kepada frontend.
- Dokumentasikan endpoint dan contoh request serta response secara sederhana.

### H. Version Control

- [x] 1. Git
- [ ] 2. Tidak menggunakan Git

**Pilihan saya:** 1

**Ketentuan:** Gunakan Git untuk mengelola perubahan kode. Pastikan .env dan file lain yang berisi rahasia masuk ke .gitignore. File .env.example tanpa kredensial asli boleh disimpan di Git.

## 4. Detail Autentikasi

**Registrasi akun:** Ya

**Login:** Ya

**Logout:** Ya

**Role pengguna:** User saja. Tidak ada role Admin.

**Aturan akses:**
1. Pengguna harus login untuk mengakses dashboard pribadi, catatan konsumsi, target nutrisi, riwayat, dan data profil yang dilindungi.
2. Pengguna yang belum login hanya dapat mengakses halaman Login dan Register serta endpoint publik yang memang diperlukan.
3. Backend harus memverifikasi JWT pada setiap request ke endpoint yang dilindungi.
4. Setiap data milik pengguna harus dikaitkan dengan identitas pengguna yang sesuai melalui relasi database yang sudah tersedia atau yang telah disetujui.
5. Backend harus memeriksa kepemilikan data pada setiap operasi membaca detail, menambahkan, mengedit, dan menghapus data pribadi.
6. Pengguna tidak boleh melihat atau mengubah catatan maupun target pengguna lain, termasuk dengan mengubah ID pada URL atau request API.
7. Password harus disimpan dalam bentuk hash, bukan teks biasa.
8. Pesan error autentikasi tidak boleh membocorkan apakah password, token, atau informasi rahasia tertentu benar.
9. Validasi dan pemeriksaan akses wajib dilakukan pada backend; menyembunyikan tombol atau halaman pada frontend saja tidak cukup.
10. Jangan mengasumsikan bahwa autentikasi Supabase Auth sudah digunakan. Periksa database dan konfigurasi proyek terlebih dahulu. Jika ternyata database sudah terintegrasi dengan Supabase Auth, jelaskan kondisi tersebut dan usulkan pendekatan yang paling sesuai sebelum mengubah rancangan autentikasi.

## 5. Konfigurasi Environment dan Database

**Ketentuan wajib:**

1. Semua konfigurasi sensitif harus menggunakan environment variable.
2. Gunakan file .env untuk konfigurasi pengembangan lokal.
3. Sediakan file .env.example tanpa kredensial asli.
4. File .env tidak boleh dimasukkan ke Git.
5. Connection string PostgreSQL harus dibaca dari environment variable bernama DATABASE_URL.
6. Jangan menulis URL database, username, password database, JWT secret, atau kredensial lain secara langsung di kode.
7. Aplikasi harus menampilkan pesan error yang jelas jika konfigurasi wajib belum tersedia, tanpa membocorkan password atau rahasia.
8. Jangan pernah menampilkan nilai rahasia dari .env pada log, respons API, atau antarmuka.
9. Gunakan konfigurasi SSL yang sesuai dengan persyaratan koneksi Supabase.
10. Jangan mengganti nilai konfigurasi yang sudah digunakan proyek tanpa memeriksa implementasi dan memberikan alasan terlebih dahulu.

Contoh `.env.example`:

```dotenv
APP_PORT=8080
DATABASE_URL=postgres://username:password@localhost:5432/personal_notes?sslmode=disable
SESSION_SECRET=replace-with-a-random-secret
```

**Catatan:** 
- Contoh di atas hanya sebagai panduan. Sesuaikan format connection string dengan informasi koneksi Supabase yang sebenarnya.
- Isi .env lokal menggunakan kredensial database milik sendiri yang diperoleh dari dashboard Supabase.
- Jangan pernah memasukkan connection string yang mengandung password asli ke dalam file spesifikasi, source code, screenshot, atau repositori publik.
- Jangan menampilkan atau mengirim nilai asli DATABASE_URL maupun JWT_SECRET kepada AI atau menuliskannya dalam log.
- Jika nama environment variable yang dibutuhkan oleh implementasi ternyata berbeda, usulkan penyelarasan terlebih dahulu dan jelaskan dampaknya.
- Jangan membuat migrasi atau mengubah schema Supabase sebelum rencana perubahan diperiksa dan disetujui.

## 6. Batasan Pengembangan (BOLEH DIUBAH SESUAI KEINGINAN PESERTA PELATIHAN)

- [x] Struktur folder sederhana dan mudah dipahami pemula.
- [x] Kode mudah dibaca serta diberi penjelasan pada bagian penting.
- [x] Tidak menambahkan teknologi di luar pilihan tanpa persetujuan.
- [x] Memisahkan frontend dan backend.
- [x] Menggunakan environment variable untuk konfigurasi.
- [x] Membuat pengujian dasar.
- [x] Menjelaskan kode dan perubahan yang dibuat AI.
- [x] Memeriksa database Supabase sebelum membuat model, repository, dan query.
- [x] Tidak menghapus atau mengubah data yang sudah ada tanpa persetujuan.
- [x] Membangun aplikasi secara bertahap setelah struktur folder disetujui.

**Batasan tambahan:**
- Saat ini belum ada kode backend maupun frontend. Keduanya harus dirancang dan dibuat dari awal.
- Database Supabase sudah tersedia dan harus diperlakukan sebagai sumber data utama.
- Jangan membuat tabel berdasarkan asumsi. Periksa struktur database aktual terlebih dahulu.
- Jika fitur yang diminta membutuhkan tabel atau kolom yang belum tersedia, laporkan kebutuhan tersebut dan ajukan rencana perubahan untuk disetujui.
- Jangan membuat data dummy seolah-olah merupakan data asli.
- Jangan menambahkan library, framework, layanan eksternal, atau fitur di luar kebutuhan tanpa menjelaskan manfaat, risiko, dan dampaknya.
- Hindari abstraksi yang berlebihan. Pilih struktur kode yang sederhana tetapi tetap memisahkan tanggung jawab utama.
- Jangan langsung membuat seluruh aplikasi dalam satu tahap.
- Setiap tahap harus dapat diperiksa dan diuji sebelum melanjutkan ke tahap berikutnya.

## 6. Fitur yang Tidak Dibutuhkan

Tuliskan fitur yang tidak perlu dibuat agar ruang lingkup proyek tetap jelas.
- Role Admin dan dashboard administrator.
- Pembayaran, langganan, dan integrasi layanan berbayar.
- Integrasi otomatis dengan smartwatch atau perangkat kesehatan.
- Integrasi API nutrisi eksternal yang belum disetujui.
- Rekomendasi diet berbasis AI.
- Forum, komunitas, dan fitur sosial.
- Fitur resep, perencanaan menu mingguan, dan daftar belanja.
- Ekspor laporan ke PDF atau Excel.
- Pelacakan berat badan dan komposisi tubuh jika belum didukung database.
- Migrasi database otomatis yang dapat mengubah atau menghapus data tanpa persetujuan.

## 7. Pemeriksaan Sebelum Vibe Coding (BOLEH DIUBAH SESUAI KEINGINAN PESERTA PELATIHAN)

- [ ] Nama dan tujuan aplikasi sudah jelas.
- [ ] Fitur utama sudah ditentukan.
- [ ] Pilihan teknologi sudah diisi.
- [ ] Kebutuhan autentikasi sudah ditentukan.
- [ ] Batasan proyek sudah ditentukan.
- [ ] Penggunaan `.env` dan `DATABASE_URL` sudah ditentukan.
- [ ] Struktur tabel dan relasi database Supabase sudah diperiksa.
- [ ] Kecocokan fitur dengan struktur database sudah dievaluasi.
- [ ] Struktur folder backend dan frontend sudah diusulkan.
- [ ] Struktur folder sudah disetujui sebelum implementasi.

**Status spesifikasi:** Siap untuk diperiksa dan disetujui sebelum implementasi.

---

## 8. Instruksi untuk Agy

Baca dan pahami file `docs/prompts/00-project-spec.md` secara lengkap. Gunakan spesifikasi ini sebagai acuan utama pengembangan NutriTrack.

Kondisi proyek saat ini:
- Database PostgreSQL di Supabase sudah tersedia.
- Backend belum dibuat.
- Frontend belum dibuat.
- Aplikasi harus dibangun dari awal tanpa menghapus atau merusak database yang sudah ada.

Tahap 1 — Audit awal, tanpa mengubah file
1. Periksa isi folder proyek dan file yang sudah tersedia.
2. Jangan menganggap backend, frontend, konfigurasi, atau struktur folder tertentu sudah ada.
3. Periksa struktur database Supabase yang tersedia menggunakan metode akses yang aman dan disetujui.
4. Identifikasi tabel, kolom, tipe data, primary key, foreign key, relasi, constraint, index, serta tabel autentikasi yang relevan.
5. Jangan menampilkan password, token, atau kredensial database dalam output.
6. Jika akses database belum tersedia, jelaskan informasi atau akses apa yang diperlukan. Jangan meminta pengguna menuliskan password database ke dalam chat.
7. Buat pemetaan fitur NutriTrack terhadap tabel dan kolom yang tersedia.
8. Tandai fitur yang bisa diimplementasikan langsung, fitur yang memerlukan perubahan schema, dan fitur yang belum bisa didukung oleh database.
9. Jangan mengubah schema, menghapus data, atau membuat tabel pada tahap ini.

Tahap 2 — Usulan struktur folder
Setelah audit selesai:
1. Usulkan struktur folder sederhana untuk backend Go Gin dan frontend HTML, CSS, serta JavaScript.
2. Jelaskan fungsi setiap folder dan file utama.
3. Usulkan pembagian kode backend, misalnya konfigurasi, koneksi database, model atau struktur data, repository, service jika diperlukan, handler, middleware, dan routes.
4. Usulkan struktur frontend yang memisahkan halaman, stylesheet, JavaScript, dan aset.
5. Jelaskan alur komunikasi frontend, REST API, backend, dan Supabase.
6. Tentukan cara menjalankan aplikasi dari terminal serta konfigurasi yang diperlukan.
7. Jelaskan strategi autentikasi JWT dan bagaimana kepemilikan data akan diperiksa.
8. Ajukan tahapan implementasi beserta pengujian yang akan dilakukan.
9. Jika terdapat beberapa pilihan arsitektur yang masuk akal, jelaskan pilihan yang direkomendasikan dan alasannya.
10. Berhenti dan tunggu persetujuan saya setelah Tahap 1 dan Tahap 2. Jangan membuat, mengubah, atau menghapus file sebelum saya menyetujui usulan struktur folder dan rencana implementasi.

Tahap 3 — Implementasi bertahap setelah persetujuan
Setelah mendapatkan persetujuan:
1. Buat struktur folder dan konfigurasi dasar.
2. Implementasikan koneksi ke Supabase melalui DATABASE_URL.
3. Implementasikan autentikasi dan middleware keamanan.
4. Implementasikan REST API sesuai kebutuhan dan struktur database yang telah diverifikasi.
5. Implementasikan frontend dengan desain yang disepakati.
6. Hubungkan frontend dengan REST API.
7. Uji alur registrasi, login, logout, akses data, pencatatan makanan, perhitungan nutrisi, target, dan riwayat sesuai fitur yang didukung database.
8. Uji bahwa pengguna tidak dapat membaca atau mengubah data milik pengguna lain.
9. Jalankan pengujian dan perbaiki error sebelum melanjutkan ke tahap berikutnya.
10. Pada akhir setiap tahap, jelaskan file yang dibuat atau diubah, tujuan perubahan, cara menjalankan, cara menguji, dan hasil pengujian.

Ketentuan tambahan untuk seluruh tahap:
- Jangan mengganti teknologi atau menambah dependensi tanpa menjelaskan alasan dan meminta persetujuan.
- Jangan meng-hardcode DATABASE_URL atau kredensial lain.
- Pastikan .env masuk ke .gitignore dan .env.example tidak berisi rahasia.
- Gunakan query berparameter dan validasi input.
- Jangan mengubah database Supabase tanpa persetujuan.
- Jangan mengklaim pengujian berhasil jika belum benar-benar dijalankan.
- Jika menemui error, jelaskan penyebab yang diketahui, bukti yang tersedia, dan solusi yang diusulkan.
- Prioritaskan implementasi yang sederhana, aman, mudah dipahami pemula, dan sesuai spesifikasi.
- Gunakan desain Quiet Elegance — Minimalist Nutrition Dashboard: tampilan minimalis, elegan, lapang, konsisten, responsif, dan mudah digunakan.
- Jika spesifikasi bertentangan dengan kondisi database aktual, hentikan bagian yang terdampak dan minta persetujuan sebelum mengambil keputusan yang mengubah data atau arsitektur.