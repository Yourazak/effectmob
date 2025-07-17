package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"service/config"
	"service/internal/subscriptions"
)

func main() {
	// Загружаем конфиг из .env
	cfg := config.LoadConfig()

	// Формируем строку подключения к базе
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// Подключаемся к базе
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()
	log.Println("Успешное подключение к базе данных!")

	// Создаем репозиторий, сервис и хендлер
	repo := subscriptions.NewRepository(db)
	service := subscriptions.NewService(repo)
	handler := subscriptions.NewHandler(service)

	// Регистрируем HTTP-ручки
	http.HandleFunc("/subscriptions", handler.SubscriptionsHandler)
	// для CRUDL

	// Запускаем HTTP сервер
	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
