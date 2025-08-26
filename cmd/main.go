package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/vshengur/bookings.auth-service/config"
	"github.com/vshengur/bookings.auth-service/database"
	"github.com/vshengur/bookings.auth-service/routes"
)

func main() {
	// Загрузка конфигурации
	config.LoadConfig()

	if config.AppConfig.RunMode == "dev" {
		// Логирование загруженных данных (безопасно для проверки, но не для секрета!)
		log.Printf("Loaded configuration: %+v\n", config.AppConfig)
	}

	// Использование конфигурации
	if config.AppConfig.JWTSecret == "" {
		log.Fatal("JWT Secret is required but not found!")
	}

	// Подключение к базе данных
	db := database.ConnectDB()
	sqlDB, err := db.DB() // Получение *sql.DB из GORM
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}()

	// Настройка роутера
	r := gin.Default()
	routes.SetupRoutes(r, db)

	// Запуск сервера
	log.Println("Starting server on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
