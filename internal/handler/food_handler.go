package handler

import (
	"net/http"
	"strconv"

	"go_vibe/internal/repository"

	"github.com/gin-gonic/gin"
)

type FoodHandler struct {
	foodRepo *repository.FoodRepository
}

func NewFoodHandler(foodRepo *repository.FoodRepository) *FoodHandler {
	return &FoodHandler{foodRepo: foodRepo}
}

// Search mencari daftar makanan berdasarkan kata kunci pencarian
func (h *FoodHandler) Search(c *gin.Context) {
	q := c.Query("q")
	limitStr := c.DefaultQuery("limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	foods, err := h.foodRepo.SearchFoods(c.Request.Context(), q, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data makanan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  foods,
		"total": len(foods),
	})
}

// GetDetail mengambil detail nutrisi dari satu makanan
func (h *FoodHandler) GetDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID makanan tidak valid."})
		return
	}

	food, err := h.foodRepo.GetFoodByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data makanan."})
		return
	}
	if food == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Makanan tidak ditemukan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": food})
}

