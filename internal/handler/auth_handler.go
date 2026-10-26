package handler

import (
	"net/http"
	"strings"

	"go_vibe/internal/config"
	"go_vibe/internal/middleware"
	"go_vibe/internal/model"
	"go_vibe/internal/repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthHandler(userRepo *repository.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, cfg: cfg}
}

// setAuthCookie mengatur cookie HttpOnly untuk token JWT
func (h *AuthHandler) setAuthCookie(c *gin.Context, token string, maxAge int) {
	isSecure := strings.ToLower(h.cfg.Environment) == "production"
	// Menggunakan SameSite Lax untuk keamanan navigasi browser lokal maupun domain produksi
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, maxAge, "/", "", isSecure, true)
}

// Register menangani registrasi akun baru
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format data pendaftaran tidak valid: periksa nama, email, dan minimal password 6 karakter."})
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Name = strings.TrimSpace(req.Name)

	// Cek apakah email sudah terdaftar
	existing, err := h.userRepo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat memeriksa akun."})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah digunakan oleh akun lain."})
		return
	}

	// Hash password menggunakan bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses kata sandi."})
		return
	}

	user, err := h.userRepo.CreateUser(c.Request.Context(), req.Name, req.Email, string(hash))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendaftarkan akun baru."})
		return
	}

	// Buat token JWT otomatis untuk langsung login
	token, err := middleware.GenerateToken(user.ID, user.Email, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Akun berhasil dibuat namun gagal membuat sesi login."})
		return
	}

	// Simpan dalam HttpOnly cookie (masa berlaku 7 hari)
	h.setAuthCookie(c, token, 7*24*3600)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pendaftaran berhasil!",
		"user": model.UserProfileResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Login menangani proses masuk pengguna
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format email atau kata sandi tidak valid."})
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	user, err := h.userRepo.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil || user == nil {
		// Pesan umum untuk menjaga keamanan
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau kata sandi tidak sesuai."})
		return
	}

	// Verifikasi hash password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau kata sandi tidak sesuai."})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat sesi login."})
		return
	}

	h.setAuthCookie(c, token, 7*24*3600)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil!",
		"user": model.UserProfileResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

// Logout menghapus sesi dan cookie login
func (h *AuthHandler) Logout(c *gin.Context) {
	h.setAuthCookie(c, "", -1)
	c.JSON(http.StatusOK, gin.H{
		"message": "Anda berhasil keluar dari akun.",
	})
}

// GetProfile mengambil data profil pengguna yang sedang login
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak terautentikasi."})
		return
	}

	user, err := h.userRepo.GetUserByID(c.Request.Context(), userID.(int))
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data pengguna tidak ditemukan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": model.UserProfileResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}

