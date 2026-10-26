package model

import "time"

// MealItem merepresentasikan item makanan pada tabel public.isi_piring
type MealItem struct {
	ID           int     `json:"id"`
	SesiMakanID  int     `json:"sesi_makan_id"`
	MakananID    int     `json:"makanan_id"`
	MakananName  string  `json:"makanan_name"`
	Gram         float64 `json:"gram"`
	ServingUnit  string  `json:"serving_unit"`
	ServingNote  *string `json:"serving_note"`
	Calories     float64 `json:"calories"`     // Dihitung: (gram / 100) * calories makanan
	Proteins     float64 `json:"proteins"`     // Dihitung: (gram / 100) * proteins makanan
	Fat          float64 `json:"fat"`          // Dihitung: (gram / 100) * fat makanan
	Carbohydrate float64 `json:"carbohydrate"` // Dihitung: (gram / 100) * carb makanan
}

// MealSession merepresentasikan sesi makan pada tabel public.sesi_makan
type MealSession struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Name      string     `json:"name"`
	MealDate  string     `json:"meal_date"` // YYYY-MM-DD
	MealTime  string     `json:"meal_time"` // Sarapan, Siang, Malam, Camilan
	Note      *string    `json:"note"`
	CreatedAt time.Time  `json:"created_at"`
	Items     []MealItem `json:"items"`

	// Total nutrisi sesi ini
	SessionCalories     float64 `json:"session_calories"`
	SessionProteins     float64 `json:"session_proteins"`
	SessionFat          float64 `json:"session_fat"`
	SessionCarbohydrate float64 `json:"session_carbohydrate"`
}

// CreateMealItemInput input makanan dalam sesi
type CreateMealItemInput struct {
	MakananID   int     `json:"makanan_id" binding:"required"`
	Gram        float64 `json:"gram" binding:"required,gt=0"`
	ServingUnit string  `json:"serving_unit"`
	ServingNote *string `json:"serving_note"`
}

// CreateMealSessionRequest payload untuk membuat sesi makan baru
type CreateMealSessionRequest struct {
	MealDate string                `json:"meal_date" binding:"required"` // Format: YYYY-MM-DD
	MealTime string                `json:"meal_time" binding:"required"` // Sarapan, Siang, Malam, Camilan
	Name     string                `json:"name"`
	Note     *string               `json:"note"`
	Items    []CreateMealItemInput `json:"items" binding:"required,min=1"`
}

// AddMealItemRequest payload untuk menambah makanan ke sesi yang sudah ada
type AddMealItemRequest struct {
	MakananID   int     `json:"makanan_id" binding:"required"`
	Gram        float64 `json:"gram" binding:"required,gt=0"`
	ServingUnit string  `json:"serving_unit"`
	ServingNote *string `json:"serving_note"`
}

// DailyNutritionSummary ringkasan nutrisi harian untuk dashboard
type DailyNutritionSummary struct {
	Date          string            `json:"date"`
	TotalCalories float64           `json:"total_calories"`
	TotalProteins float64           `json:"total_proteins"`
	TotalFat      float64           `json:"total_fat"`
	TotalCarb     float64           `json:"total_carbohydrate"`
	Target        *NutritionTarget  `json:"target"`
	Sessions      []MealSession     `json:"sessions"`
}

