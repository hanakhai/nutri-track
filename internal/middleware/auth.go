package middleware

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims struktur isi token JWT
type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// GenerateToken membuat token JWT baru dengan masa berlaku 7 hari
func GenerateToken(userID int, email string, secret string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken memverifikasi dan membaca payload token JWT
func ParseToken(tokenStr string, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode penandatanganan token tidak valid")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("token tidak valid atau kedaluwarsa")
}

// AuthRequired middleware untuk melindungi rute yang mewajibkan login
func AuthRequired(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("auth_token")
		if err != nil || tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Sesi tidak ditemukan. Silakan login terlebih dahulu.",
			})
			c.Abort()
			return
		}

		claims, err := ParseToken(tokenStr, jwtSecret)
		if err != nil {
			// Hapus cookie yang rusak/kedaluwarsa
			c.SetCookie("auth_token", "", -1, "/", "", false, true)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Sesi telah berakhir atau tidak valid. Silakan login kembali.",
			})
			c.Abort()
			return
		}

		// Simpan user_id dan email di context Gin agar bisa diakses oleh handler
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}

