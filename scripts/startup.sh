DB_DRIVER="postgres"
DB_STRING="postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable"
MIGRATIONS_DIR="./migrations"
DEFAULT_DB="postgres"

export PGPASSWORD=$POSTGRES_PASSWORD
TARGET_USER="$POSTGRES_USER"
TARGET_PASSWORD="$POSTGRES_PASSWORD"

./scripts/wait-for-it.sh "$POSTGRES_HOST:$POSTGRES_PORT" -t 60 -- echo "Postgres доступен - выполнение команд"

echo "Выполняем миграции по созданию базы данных '$POSTGRES_DB'..."
psql -h "$POSTGRES_HOST" -U "$DEFAULT_DB" -c "
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_database WHERE datname = '$POSTGRES_DB') THEN
      EXECUTE format('CREATE DATABASE %I', '$POSTGRES_DB');
   END IF;
   -- Создаем пользователя, если его ещё нет
   IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = '$TARGET_USER') THEN
      EXECUTE format('CREATE USER %I WITH PASSWORD %L', '$TARGET_USER', '$TARGET_PASSWORD');
   END IF;
   -- Выдаем права на схему public и все таблицы пользователю
   EXECUTE format('
       GRANT ALL ON SCHEMA public TO %I;
       GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public to %I;
       ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO %I;
   ', '$TARGET_USER', '$TARGET_USER', '$TARGET_USER');
END
\$\$;"

if [ $? -eq 0 ]; then
    echo "Миграции по созданию '$POSTGRES_DB' успешно выполнены."
else
    echo "Ошибка при выполнении миграции по созданию '$POSTGRES_DB'."
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
echo "Выполнение миграций на базе '$POSTGRES_DB'..."
goose -dir "$MIGRATIONS_DIR" "$DB_DRIVER" "$DB_STRING" up

# Проверка статуса миграций
if [ $? -eq 0 ]; then
    echo "Миграции на базе '$POSTGRES_DB' успешно выполнены."
else
    echo "Ошибка при выполнении миграций на базе '$POSTGRES_DB'."
    exit 1
fi

echo "Запуск приложения..."
./sad