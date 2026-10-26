package main

import (
	"fmt"
	"log"

	"go_vibe/internal/config"
	"go_vibe/internal/database"
	"go_vibe/internal/handler"
	"go_vibe/internal/middleware"
	"go_vibe/internal/repository"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Muat konfigurasi
	cfg := config.LoadConfig()

	// Atur mode Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 2. Inisialisasi Database
	db := database.InitDB(cfg.DatabaseURL)
	defer db.Close()

	// 3. Inisialisasi Repository
	userRepo := repository.NewUserRepository(db)
	foodRepo := repository.NewFoodRepository(db)
	targetRepo := repository.NewTargetRepository(db)
	logRepo := repository.NewLogRepository(db, targetRepo)

	// 4. Inisialisasi Handler
	authHandler := handler.NewAuthHandler(userRepo, cfg)
	foodHandler := handler.NewFoodHandler(foodRepo)
	targetHandler := handler.NewTargetHandler(targetRepo)
	logHandler := handler.NewLogHandler(logRepo)

	// 5. Setup Router Gin
	r := gin.Default()

	// Sajikan file statis frontend (HTML, CSS, JS, Gambar)
	r.Static("/static", "./web")
	r.StaticFile("/", "./web/index.html")
	r.StaticFile("/login", "./web/login.html")
	r.StaticFile("/register", "./web/register.html")

	// 6. Rute REST API
	api := r.Group("/api")
	{
		// Rute Autentikasi Publik
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/logout", authHandler.Logout)
		}

		// Rute Makanan Publik (pencarian katalog)
		foodGroup := api.Group("/foods")
		{
			foodGroup.GET("", foodHandler.Search)
			foodGroup.GET("/:id", foodHandler.GetDetail)
		}

		// Rute yang Dilindungi Autentikasi JWT
		protected := api.Group("")
		protected.Use(middleware.AuthRequired(cfg.JWTSecret))
		{
			// Profil Pengguna
			protected.GET("/auth/me", authHandler.GetProfile)

			// Target Nutrisi Pribadi
			protected.GET("/targets", targetHandler.GetTarget)
			protected.PUT("/targets", targetHandler.UpdateTarget)

			// Catatan Makan & Ringkasan Harian
			protected.GET("/logs/summary", logHandler.GetDailySummary)
			protected.GET("/logs", logHandler.GetDailyLogs)
			protected.POST("/logs", logHandler.CreateSession)
			protected.POST("/logs/:id/items", logHandler.AddItemToSession)
			protected.DELETE("/logs/:id", logHandler.DeleteSession)
			protected.DELETE("/logs/:id/items/:itemId", logHandler.DeleteItem)
		}
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	fmt.Printf("\n==================================================\n")
	fmt.Printf(" NutriTrack berjalan di http://localhost:%s\n", cfg.AppPort)
	fmt.Printf("==================================================\n\n")

	if err := r.Run(addr); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

