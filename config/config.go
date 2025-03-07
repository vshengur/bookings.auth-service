package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/vshengur/bookings.auth-service/services"
)

type Config struct {
	RunMode           string
	DbUser            string
	DbPassword        string
	DbHost            string
	DbPort            string
	DbName            string
	GoogleRedirectURL string
	GoogleClientID    string
	GoogleSecret      string
	JWTSecret         string
}

var AppConfig *Config

// LoadConfig загружает конфигурацию в порядке приоритета: .env → ENV → Consul
func LoadConfig() {
	// Шаг 1: Загрузка из .env файла (если он существует)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, skipping...")
	}

	// Шаг 2: Чтение из ENV переменных
	viper.AutomaticEnv()

	// Шаг 3: Загрузка конфигурации Consul
	services.LoadConsulServiceConfig()

	// Чтение из Consul (если не найдено в ENV или .env)
	AppConfig = &Config{
		RunMode: getConfigValue("RUN_MODE"),

		DbUser:     getConfigValue("DB_USER"),
		DbPassword: getConfigValue("DB_PASSWORD"),
		DbHost:     getConfigValue("DB_HOST"),
		DbPort:     getConfigValue("DB_PORT"),
		DbName:     getConfigValue("DB_NAME"),

		GoogleRedirectURL: getConfigValue("GOOGLE_REDIRECT_URL"),
		GoogleClientID:    getConfigValue("GOOGLE_CLIENT_ID"),
		GoogleSecret:      getConfigValue("GOOGLE_CLIENT_SECRET"),
		JWTSecret:         getConfigValue("JWT_SECRET"),
	}
}

// getConfigValue получает значение из ENV или Consul
func getConfigValue(key string) string {
	// Проверяем ENV переменные
	value := viper.GetString(key)
	if value != "" {
		return value
	}

	// Если не найдено, читаем из Consul
	value, err := services.GetConsulSecret(key)
	if err != nil {
		log.Printf("Error fetching key %s from Consul: %v", key, err)
	}
	return value
}
