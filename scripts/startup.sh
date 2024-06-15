#!/bin/bash

DB_DRIVER="postgres"
DB_STRING="postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable"
MIGRATIONS_DIR="./migrations"
DEFAULT_DB="postgres"

export PGPASSWORD=$POSTGRES_PASSWORD

./scripts/wait-for-it.sh "$POSTGRES_HOST:$POSTGRES_PORT" -t 60 -- echo "Postgres доступен - выполнение команд"

echo "Выполняем миграции по созданию базы данных 'sad'..."
psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$DEFAULT_DB" -c "
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = 'sad') THEN
      PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE sad');
   END IF;
END
\$\$;"

if [ $? -eq 0 ]; then
    echo "Миграции по созданию 'sad' успешно выполнены."
else
    echo "Ошибка при выполнении миграции по созданию 'sad'."
    exit 1
fi

# Проверяем, установлена ли утилита goose
if ! command -v goose &> /dev/null
then
    echo "Goose не установлен. Установите его с помощью команды:"
    echo "go install github.com/pressly/goose/v3/cmd/goose@latest"
    exit 1
fi

# Прогон миграций
echo "Выполнение миграций на базе 'sad'..."
goose -dir "$MIGRATIONS_DIR" "$DB_DRIVER" "$DB_STRING" up

# Проверка статуса миграций
if [ $? -eq 0 ]; then
    echo "Миграции на базе 'sad' успешно выполнены."
else
    echo "Ошибка при выполнении миграций на базе 'sad'."
    exit 1
fi

echo "Запуск приложения..."
./sad