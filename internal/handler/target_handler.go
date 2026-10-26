package handler

import (
	"net/http"

	"go_vibe/internal/model"
	"go_vibe/internal/repository"

	"github.com/gin-gonic/gin"
)

type TargetHandler struct {
	targetRepo *repository.TargetRepository
}

func NewTargetHandler(targetRepo *repository.TargetRepository) *TargetHandler {
	return &TargetHandler{targetRepo: targetRepo}
}

// GetTarget mengambil target nutrisi pengguna
func (h *TargetHandler) GetTarget(c *gin.Context) {
	userID := c.GetInt("user_id")

	target, err := h.targetRepo.GetTargetByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data target nutrisi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": target})
}

// UpdateTarget memperbarui target kalori dan makronutrisi harian pengguna
func (h *TargetHandler) UpdateTarget(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req model.UpdateTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format target tidak valid. Kalori harus lebih dari 0."})
		return
	}

	updated, err := h.targetRepo.UpsertTarget(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan target nutrisi."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Target nutrisi berhasil diperbarui!",
		"data":    updated,
	})
}

