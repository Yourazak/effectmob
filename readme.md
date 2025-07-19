# 📦 Subscription Service

Сервис для управления онлайн-подписками. Позволяет добавлять, получать, обновлять и удалять подписки пользователей.

## 🚀 Технологии

- Язык: Go
- База данных: PostgreSQL
- Docker + Docker Compose
- Swagger для документации API

## 📁 Структура проекта


## ⚙️ Установка и запуск

### 1. Клонируй репозиторий

```bash
git clone <repo-url>
cd service
docker-compose up --build
Приложение будет доступно на http://localhost:8080

Swagger UI: http://localhost:8080/swagger/index.html

PostgreSQL: localhost:5432
Выполни вручную или через инструмент вроде migrate, goose, sqlc и т.д.

🔌 REST API
GET /subscriptions
Получить все подписки.

Response:

json
Копировать
Редактировать
[
  {
    "id": 1,
    "user_id": 101,
    "service_name": "Netflix",
    "price": 9.99,
    "renewal_date": "2025-08-01"
  }
]
POST /subscriptions
Добавить новую подписку.

Request:

json
Копировать
Редактировать
{
  "user_id": 101,
  "service_name": "Spotify",
  "price": 4.99,
  "renewal_date": "2025-09-01"
}
PUT /subscriptions?id=1
Обновить существующую подписку по ID.

DELETE /subscriptions?id=1
Удалить подписку по ID.

📚 Swagger
Swagger UI доступен по адресу:

bash
Копировать
Редактировать
http://localhost:8080/swagger/index.html
Документация автоматически генерируется по аннотациям в main.go.

🧪 Тесты
Юнит-тесты не включены, так как не требуются по заданию.
Могут быть легко добавлены для слоя service и handler.

🐳 Docker
Сборка образа вручную:

bash
Копировать
Редактировать
docker build -t subscription-service .
Запуск контейнера:

bash
Копировать
Редактировать
docker run -p 8080:8080 subscription-service