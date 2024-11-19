# Используем официальный базовый образ Go
FROM golang:1.23 as builder

# Устанавливаем рабочий каталог
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код в контейнер
COPY . .

# Собираем приложение.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# Настройка финального образа
FROM alpine:latest

WORKDIR /root/

# Копируем исполняемый файл из билдера
COPY --from=builder /app/main .

# Открываем порт, который слушает ваше приложение
EXPOSE 8000

# Команда для запуска приложения
CMD ["./main"]
