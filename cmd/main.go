package main

import (
	"fmt"
	"log"
	"net/http"
	"service/config"
	"service/internal/subscriptions"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "service/docs" // импорт сгенерированных swagger-документов (путь подкорректируй!)
)

// @title Subscription Service API
// @version 1.0
// @description Сервис для управления онлайн-подписками
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	cfg := config.LoadConfig()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()
	log.Println("Успешное подключение к базе данных!")

	repo := subscriptions.NewRepository(db)
	service := subscriptions.NewService(repo)
	handler := subscriptions.NewHandler(service)

	// API endpoint для подписок
	http.HandleFunc("/subscriptions", handler.SubscriptionsHandler)

	// Swagger UI endpoint по адресу http://localhost:8080/swagger/index.html
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
