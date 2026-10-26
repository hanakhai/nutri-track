package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// InitDB menginisialisasi connection pool ke database PostgreSQL
func InitDB(databaseURL string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Gagal mem-parsing konfigurasi database: %v", err)
	}

	// Konfigurasi pool
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Gagal membuat connection pool database: %v", err)
	}

	// Ping untuk memastikan koneksi berhasil
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	fmt.Println("✓ Berhasil terhubung ke database PostgreSQL Supabase.")
	return pool
}

