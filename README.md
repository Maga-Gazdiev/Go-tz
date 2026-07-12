# Subscriptions API

REST-сервис для учёта онлайн-подписок пользователей.

## Возможности

- создание, получение, изменение, удаление и вывод списка подписок;
- расчёт суммарной стоимости за период;
- фильтрация по UUID пользователя и названию сервиса;
- PostgreSQL-миграции через Goose;
- структурированные HTTP-логи с request ID.

## Swagger / OpenAPI

После запуска Swagger UI доступен по адресу `http://localhost:8099/docs/`.

OpenAPI-спецификация также доступна в репозитории: [docs/openapi.yaml](docs/openapi.yaml).

## API

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/v1/subscriptions` | создать |
| GET | `/api/v1/subscriptions` | список |
| GET | `/api/v1/subscriptions/{id}` | одна запись |
| PUT | `/api/v1/subscriptions/{id}` | обновить |
| DELETE | `/api/v1/subscriptions/{id}` | удалить |
| GET | `/api/v1/subscriptions/total` | сумма за период |

### Создание

```bash
curl -X POST http://localhost:8099/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
  }'
```

### Сумма за период

Считается как `цена × кол-во месяцев`, когда подписка была активна в указанном диапазоне.

```bash
curl "http://localhost:8099/api/v1/subscriptions/total?from=01-2025&to=12-2025&user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba"
```

Фильтры `user_id` и `service_name` опциональны.

## Конфигурация

Перед локальным запуском создайте `.env` из шаблона и заполните параметры подключения:

```bash
cp .env.example .env
```

Пример содержимого `.env` для локального PostgreSQL:

```dotenv
SERVER_HOST=0.0.0.0
SERVER_PORT=8099

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSLMODE=disable

LOG_LEVEL=info
```

`LOG_LEVEL` определяет минимальный уровень выводимых логов:

| Значение | Вывод |
|----------|-------|
| `debug` | все сообщения, включая отладочные |
| `info` | обычные сообщения, HTTP-запросы, предупреждения и ошибки |
| `warn` | только предупреждения и ошибки |
| `error` | только ошибки |

При запуске через Docker Compose адрес БД уже задан как `db`, поэтому дополнительное редактирование `.env` не требуется.

## Стек

Go 1.25, PostgreSQL, pgx, Goose, slog

## Структура

```text
cmd/api                 # запуск HTTP API
cmd/migrate             # запуск Goose-миграций
internal/app            # сборка зависимостей и HTTP-сервера
internal/handler/subscription
internal/service/subscription
internal/repository/subscription
internal/infrastructure # PostgreSQL и миграции
migrations
```

## Команды Makefile

| Команда | Назначение |
|---------|------------|
| `make run` | запустить API локально |
| `make migrate` | применить Goose-миграции |
| `make test` | собрать и проверить все Go-пакеты; будущие тесты будут запущены этой же командой |
| `make vet` | найти типичные ошибки с помощью `go vet` |
| `make format` | отформатировать Go-файлы через `gofmt` |
| `make lint` | выполнить `golangci-lint` |
| `make check` | последовательно выполнить `test`, `vet` и `lint` |
| `make up` / `make down` | запустить / остановить Docker Compose |

Для локального `make lint` установите линтер один раз:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

## Развёртывание

### Docker Compose

```bash
make up
```

После старта:

- API: `http://localhost:8099/api/v1`;
- health-check: `http://localhost:8099/health`.

### Локальный запуск

Требуется PostgreSQL 16+.

```bash
cp .env.example .env
# заполните параметры подключения к локальной БД в .env

make migrate
make run
```
