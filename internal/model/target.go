package model

import "time"

// NutritionTarget merepresentasikan target nutrisi harian pengguna
type NutritionTarget struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	CaloriesTarget     float64   `json:"calories_target"`
	ProteinsTarget     float64   `json:"proteins_target"`
	CarbohydrateTarget float64   `json:"carbohydrate_target"`
	FatTarget          float64   `json:"fat_target"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// UpdateTargetRequest payload untuk mengubah target nutrisi
type UpdateTargetRequest struct {
	CaloriesTarget     float64 `json:"calories_target" binding:"required,gt=0"`
	ProteinsTarget     float64 `json:"proteins_target" binding:"required,gte=0"`
	CarbohydrateTarget float64 `json:"carbohydrate_target" binding:"required,gte=0"`
	FatTarget          float64 `json:"fat_target" binding:"required,gte=0"`
}

