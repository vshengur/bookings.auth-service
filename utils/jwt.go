package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vshengur/bookings.auth-service/config"
	"github.com/vshengur/bookings.auth-service/models"
)

func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

// validateJWT проверяет JWT токен
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Используем секрет из конфигурации
	secret := []byte(config.AppConfig.JWTSecret)

	// Разбор токена
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверка метода подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("Unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	// Извлечение claims (данных из токена)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.NewValidationError("Invalid token claims", jwt.ValidationErrorClaimsInvalid)
	}

	// Дополнительные проверки:
	// 1. Проверка истечения срока действия
	if !validateTokenExpiry(claims) {
		return nil, jwt.NewValidationError("Token has expired", jwt.ValidationErrorClaimsInvalid)
	}

	return claims, nil
}

// validateTokenExpiry проверяет срок действия токена
func validateTokenExpiry(claims jwt.MapClaims) bool {
	exp, ok := claims["exp"].(float64) // exp хранится в виде Unix timestamp
	if !ok {
		return false
	}
	return time.Now().Unix() < int64(exp)
}
