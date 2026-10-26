package handler

import (
	"net/http"
	"strconv"
	"time"

	"go_vibe/internal/model"
	"go_vibe/internal/repository"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logRepo *repository.LogRepository
}

func NewLogHandler(logRepo *repository.LogRepository) *LogHandler {
	return &LogHandler{logRepo: logRepo}
}

// getDateParam membaca tanggal dari query param atau menggunakan hari ini (YYYY-MM-DD)
func getDateParam(c *gin.Context) string {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}
	return dateStr
}

// GetDailySummary mengambil ringkasan harian (total nutrisi, target perbandingan, dan sesi makan)
func (h *LogHandler) GetDailySummary(c *gin.Context) {
	userID := c.GetInt("user_id")
	dateStr := getDateParam(c)

	summary, err := h.logRepo.GetDailySummary(c.Request.Context(), userID, dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil ringkasan nutrisi harian."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

// GetDailyLogs mengambil daftar sesi makan pada tanggal tertentu
func (h *LogHandler) GetDailyLogs(c *gin.Context) {
	userID := c.GetInt("user_id")
	dateStr := getDateParam(c)

	sessions, err := h.logRepo.GetSessionsByDate(c.Request.Context(), userID, dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil catatan makan."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"date":     dateStr,
		"sessions": sessions,
	})
}

// CreateSession membuat sesi makan baru beserta makanannya
func (h *LogHandler) CreateSession(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req model.CreateMealSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data sesi makan tidak lengkap: pastikan tanggal, waktu makan, dan makanan telah dipilih."})
		return
	}

	session, err := h.logRepo.CreateMealSession(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan catatan makan: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Catatan makan berhasil disimpan!",
		"data":    session,
	})
}

// AddItemToSession menambahkan makanan ke sesi yang sudah ada
func (h *LogHandler) AddItemToSession(c *gin.Context) {
	userID := c.GetInt("user_id")
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sesi tidak valid."})
		return
	}

	var req model.AddMealItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data makanan tidak lengkap."})
		return
	}

	err = h.logRepo.AddMealItem(c.Request.Context(), userID, sessionID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Makanan berhasil ditambahkan ke sesi makan!"})
}

// DeleteSession menghapus satu sesi makan
func (h *LogHandler) DeleteSession(c *gin.Context) {
	userID := c.GetInt("user_id")
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID sesi tidak valid."})
		return
	}

	err = h.logRepo.DeleteMealSession(c.Request.Context(), userID, sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sesi makan berhasil dihapus."})
}

// DeleteItem menghapus satu item makanan dari sesi makan
func (h *LogHandler) DeleteItem(c *gin.Context) {
	userID := c.GetInt("user_id")
	itemIDStr := c.Param("itemId")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid."})
		return
	}

	err = h.logRepo.DeleteMealItem(c.Request.Context(), userID, itemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item makanan berhasil dihapus."})
}

