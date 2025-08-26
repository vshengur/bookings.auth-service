package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vshengur/bookings.auth-service/utils"
)

type TokenHandler struct{}

// ValidateToken проверяет валидность токена
func (h *TokenHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Извлечение токена из заголовка
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Проверка токена
	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
		return
	}

	// Успешная проверка
	c.JSON(http.StatusOK, gin.H{
		"message": "Token is valid",
		"claims":  claims,
	})
}
