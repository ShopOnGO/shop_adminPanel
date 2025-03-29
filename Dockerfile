FROM golang:1.23.3 AS builder

WORKDIR /admin

# Устанавливаем pg_isready и очищаем кеш
RUN apt-get update && apt-get install -y postgresql-client \
    && rm -rf /var/lib/apt/lists/* && apt-get clean

# Отключаем CGO для статической компиляции
 ENV CGO_ENABLED=0

# Копируем файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download && go mod verify

# Копируем весь код
COPY . .

# Компилируем бинарник
RUN go build -o /admin/admin_panel ./cmd/server.go



# Второй этап: финальный образ (без лишних инструментов)
FROM alpine:latest

WORKDIR /admin

# Устанавливаем postgresql-client и dos2unix
RUN apk add --no-cache postgresql-client dos2unix

COPY .env /admin/.env

# Копируем бинарный файл из предыдущего этапа
COPY --from=builder /admin/admin_panel /admin/admin_panel

# Копируем wait-for-db.sh и делаем исполняемым
COPY --from=builder /admin/wait-for-db.sh /admin/wait-for-db.sh
RUN chmod +x /admin/wait-for-db.sh

# Преобразуем формат строки в скрипте wait-for-db.sh в Unix-формат
RUN dos2unix /admin/wait-for-db.sh

# Запуск приложения
CMD ["/admin/admin_panel"]
