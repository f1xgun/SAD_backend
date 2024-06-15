# Указываем базовый образ для компиляции Go-кода
FROM golang:1.22-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копируем остальные файлы проекта
COPY ./ ./

# Компилируем проект
RUN go build -o sad ./cmd/main.go

## Создаем финальный сжатый образ
FROM alpine:latest

## Устанавливаем зависимости
RUN apk --no-cache add ca-certificates bash postgresql-client

# Устанавливаем рабочую директорию
WORKDIR /root/

COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копируем скомпилированное бинарное приложение из промежуточного образа
COPY --from=builder /app/config.env ./config.env
COPY --from=builder /app/sad ./sad
COPY ./internal/db/migrations/ ./migrations/
COPY ./scripts/startup.sh scripts/wait-for-it.sh ./scripts/

CMD ["/app/sad"]