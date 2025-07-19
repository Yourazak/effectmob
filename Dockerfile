# ---------- Этап сборки ----------
FROM golang:1.24.4-alpine AS builder

# Установка необходимых пакетов
RUN apk add --no-cache git ca-certificates

# Рабочая директория внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем всё остальное
COPY . .

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -o /subscription-service ./cmd/main.go

# ---------- Этап запуска ----------
FROM alpine:latest

# Установка сертификатов (для HTTPS-запросов, если нужны)
RUN apk --no-cache add ca-certificates

# Рабочая директория
WORKDIR /root/

# Копируем бинарник из builder stage
COPY --from=builder /subscription-service .

# Открываем порт
EXPOSE 8080

# Команда запуска
CMD ["./subscription-service"]
