# TodoKy

> REST API Todo-приложение на Go — учебный проект с продакшн-архитектурой.

Цель проекта — реализовать todo-сервис, применив паттерны и инструменты, которые используются в реальных Go-бэкендах: чистую архитектуру, Repository Pattern, Dependency Injection, graceful shutdown, миграции БД, структурированное логирование и порт-форвардинг через socat.

---

## 🌐 Live Demo

Приложение задеплоено и доступно прямо сейчас:

| | Адрес |
|---|---|
| **API** | http://45.131.42.160:5050/ |
| **Swagger UI** | http://45.131.42.160:5050/swagger/index.html |

---

## 📦 Стек технологий

| Категория | Инструмент |
|-----------|-----------|
| **Язык** | Go 1.26+ |
| **База данных** | PostgreSQL (через `pgx/v5`) |
| **Документация API** | Swagger / OpenAPI (`swaggo/swag`, `swaggo/http-swagger`) |
| **Валидация** | `go-playground/validator/v10` |
| **Логирование** | `go.uber.org/zap` (структурированные логи) |
| **Конфигурация** | `kelseyhightower/envconfig` (env → struct) |
| **UUID** | `google/uuid` |
| **Контейнеризация** | Docker, Docker Compose |
| **Миграции** | golang-migrate (через Docker Compose) |
| **Порт-форвардинг** | `socat` (в виде отдельного Docker-сервиса) |
| **Сборка** | Makefile |

---

## 🏗 Архитектура

Проект следует принципам **Clean Architecture**. Код разбит на три уровня, каждый из которых не знает о деталях реализации следующего:

```
cmd/
└── todoky/
    ├── main.go        ← точка входа, composition root
    └── Dockerfile

internal/
├── core/              ← ядро, не зависит от внешних деталей
│   ├── domain/        ← доменные модели и бизнес-логика
│   ├── errors/        ← типизированные ошибки
│   ├── logger/        ← абстракция логгера
│   ├── repository/    ← интерфейсы и пул соединений (pgx pool)
│   └── transport/
│       └── http/
│           ├── server/      ← HTTP-сервер
│           └── middleware/  ← middleware (CORS, логирование и т.д.)
└── features/          ← бизнес-фичи
    ├── tasks/
    │   ├── repository/postgres/  ← реализация репозитория для PostgreSQL
    │   ├── service/              ← бизнес-логика (use cases)
    │   └── transport/http/       ← HTTP-хендлеры
    ├── users/
    │   ├── repository/postgres/
    │   ├── service/
    │   └── transport/
    ├── statistics/
    │   ├── repository/postgres/
    │   ├── service/
    │   └── transport/
    └── web/           ← статические ресурсы

migrations/            ← SQL-миграции (up/down)
```

---

## 🧩 Реализованные паттерны

### 1. Clean Architecture (Чистая архитектура)

Код разделён на три кольца зависимости:

- **Domain** (`internal/core/domain`) — центральное кольцо. Содержит чистые доменные модели (`Task`, `User`) и бизнес-правила (`CompletionDuration()`). Не зависит ни от HTTP, ни от PostgreSQL.
- **Service / Use Cases** (`internal/features/*/service`) — слой бизнес-логики. Каждая операция вынесена в отдельный файл: `create_task.go`, `get_task.go`, `patch_task.go`, `delete_task.go`.
- **Transport & Repository** — внешние адаптеры. Transport — HTTP-хендлеры, Repository — конкретные реализации для PostgreSQL.

Правило зависимостей: **стрелки зависимостей направлены только внутрь** — transport зависит от service, service зависит от domain, но не наоборот.

---

### 2. Repository Pattern (Паттерн Репозиторий)

Сервисный слой работает не с конкретными реализациями БД, а с **интерфейсом**:

```go
// internal/features/tasks/service/service.go
type TasksRepository interface {
    CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
    GetTasks(ctx context.Context, userID int, ...) ([]domain.Task, error)
    GetTask(ctx context.Context, taskID int) (domain.Task, error)
    PatchTask(ctx context.Context, task domain.Task) (domain.Task, error)
    DeleteTask(ctx context.Context, taskID int) error
}

type TasksService struct {
    tasksRepository TasksRepository
}
```

Конкретная реализация (`internal/features/tasks/repository/postgres`) подключается снаружи (в `main.go`). Это позволяет подменить PostgreSQL на другое хранилище без изменения бизнес-логики.

> **В реальных проектах:** такой же паттерн используется повсеместно — это основа тестируемости (легко подменить репозиторий моком).

---

### 3. Dependency Injection (Внедрение зависимостей)

Все зависимости (репозиторий, логгер) **явно передаются** через конструктор. В `main.go` собирается граф зависимостей вручную — без фреймворков:

```go
// cmd/todoky/main.go (composition root)
pool := core_pgx_pool.New(...)

userRepo    := users_postgres_repository.New(pool)
userService := users_service.New(userRepo)
userHandler := users_transport_http.New(userService)

taskRepo    := task_postgres_repository.New(pool)
taskService := task_service.New(taskRepo)
taskHandler := tasks_transport_http.New(taskService)

// ...подключение роутов
```

> **В реальных проектах:** Google Wire, Uber fx, или такой же ручной DI — всё зависит от команды. Ручной DI прозрачен и легко дебажится.

---

### 4. Graceful Shutdown (Плавная остановка)

Приложение перехватывает POSIX-сигналы `SIGINT` и `SIGTERM` и корректно завершает работу:

```go
// cmd/todoky/main.go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
defer cancel()
server.Shutdown(ctx)
```

Сервер ждёт завершения уже принятых запросов, не обрывая соединений. Паттерн — **стандартный** для любого продакшн Go-сервиса.

---

### 5. Optimistic Locking (Оптимистическая блокировка)

Доменные модели содержат поле `Version`:

```go
// internal/core/domain/task.go
type Task struct {
    ID      int
    Version int  // ← версия для оптимистической блокировки
    // ...
}
```

```sql
-- migrations/000001_init.up.sql
CREATE TABLE todoky.users (
    id       SERIAL PRIMARY KEY,
    version  BIGINT NOT NULL DEFAULT 1,
    ...
);
```

При обновлении запись WHERE-условие проверяет совпадение версии. Если другой процесс уже изменил запись — версия не совпадёт, возникнет конфликт. Это предотвращает потерю данных при конкурентных изменениях без дорогостоящих эксклюзивных блокировок.

> **В реальных проектах:** Optimistic Locking широко применяется в PostgreSQL-сервисах вместо `SELECT FOR UPDATE` когда конкуренция изменений невысока.

---

### 6. Nullable / Optional Fields Pattern

В доменных моделях используются указатели для опциональных полей:

```go
type Task struct {
    Description *string    // nil = не задано
    CompletedAt *time.Time // nil = задача не завершена
}
```

Вспомогательный тип `Nullable` в `internal/core/domain/nullabel.go` абстрагирует работу с такими полями. Это идиоматичный Go-способ отличить «пустое значение» от «значение не задано», что критично для PATCH-запросов (частичное обновление).

---

### 7. Feature-Sliced структура проекта

Вместо разбивки по техническому слою (`controllers/`, `models/`, `repositories/`) проект организован по **бизнес-фичам**:

```
internal/features/
├── tasks/       ← всё, что касается задач
├── users/       ← всё, что касается пользователей
└── statistics/  ← статистика
```

Каждая фича самодостаточна: содержит свой репозиторий, сервис и транспорт. Добавление новой фичи не затрагивает существующий код.

> **В реальных проектах:** Feature Slicing (или Vertical Slice Architecture) всё чаще предпочитают горизонтальному разбиению в крупных командах.

---

### 8. Middleware Chain (Цепочка Middleware)

HTTP-сервер использует цепочку middleware (`internal/core/transport/http/middleware`):

- CORS — управление разрешёнными источниками через `ALLOWED_ORIGINS`
- Структурированное логирование входящих запросов через `zap`
- Единый формат ответа на ошибки

Каждый middleware — отдельная функция-обёртка над `http.Handler`, что позволяет комбинировать их произвольно.

---

### 9. Environment-Based Configuration (12-Factor App)

Конфигурация полностью вынесена в переменные окружения и загружается при старте через `envconfig`:

```bash
# .env.example
HTTP_ADDR=
HTTP_SHUTDOWN_TIMEOUT=
ALLOWED_ORIGINS=

POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_TIMEOUT=

LOGGER_LEVEL=
```

`Makefile` автоматически экспортирует переменные из `.env` (`include .env; export`), так что разработчик просто копирует `.env.example → .env` и стартует.

---

### 10. Multi-Stage Docker Build (Многоэтапная сборка)

```dockerfile
# Этап 1: сборка
FROM golang:1.26.2-bookworm AS builder
RUN CGO_ENABLED=0 go build -o /app/todoky ./cmd/todoky

# Этап 2: минимальный рантайм-образ
FROM alpine:3.23
COPY --from=builder /app/todoky .
```

Финальный образ содержит только бинарник без Go-тулчейна, что уменьшает его размер в ~10 раз и снижает поверхность атаки.

---

## 🔌 socat: Port Forwarder как паттерн

В `docker-compose.yaml` описан отдельный сервис `port-forwarder` на базе **socat**.

**socat** (SOcket CAT) — утилита, которая работает как **коммутационный переключатель (switch) для потоков данных**. В данном контексте она позволяет пробросить порт PostgreSQL из Docker-сети на localhost разработчика без изменения конфигурации контейнера БД.

```yaml
# docker-compose.yaml
services:
  port-forwarder:
    image: alpine/socat
    command: TCP-LISTEN:5432,fork TCP:todoky-postgres:5432
    ports:
      - "5432:5432"
```

**Как это работает:**
1. `TCP-LISTEN:5432` — socat слушает подключения на порту 5432 внутри контейнера
2. `fork` — для каждого входящего соединения создаётся отдельный процесс (как inetd)
3. `TCP:todoky-postgres:5432` — каждое соединение проксируется к PostgreSQL-контейнеру

**Зачем это нужно:**
- Разработчик может подключить любой GUI-клиент (DBeaver, TablePlus) к `localhost:5432`
- Сервис включается отдельной командой (`make env-port-forward`) и не потребляет ресурсы в продакшне
- Не нужно пробрасывать порты PostgreSQL напрямую (что небезопасно)

**Управление через Makefile:**
```makefile
env-port-forward:
    @docker compose up -d port-forwarder

env-port-close:
    @docker compose down port-forwarder
```

> **В реальных проектах:** socat используют для отладки сетевых соединений, туннелирования, реализации простых прокси и port-forwarding в локальных dev-окружениях. Аналогичный подход практикуется в Kubernetes через `kubectl port-forward`.

---

## 🗄 База данных

### Схема

```sql
CREATE SCHEMA todoky;

CREATE TABLE todoky.users (
    id           SERIAL PRIMARY KEY,
    version      BIGINT NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL CHECK (char_length(full_name) BETWEEN 3 AND 100),
    phone_number VARCHAR(15)  CHECK (phone_number ~ '^\+[0-9]+$' ...)
);

CREATE TABLE todoky.tasks (
    id           SERIAL PRIMARY KEY,
    version      BIGINT NOT NULL DEFAULT 1,
    title        TEXT NOT NULL,
    description  TEXT,
    completed    BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    author_user_id INT REFERENCES todoky.users(id)
);
```

Используется отдельная PostgreSQL-схема `todoky` — хорошая практика для изоляции таблиц одного сервиса от других в той же БД.

### Миграции

Миграции хранятся в `migrations/` в формате `golang-migrate` (`000001_init.up.sql` / `000001_init.down.sql`). Применяются через Docker Compose:

```bash
make migrate-create seq=add_tags  # создать новую миграцию
# make migrate-up / migrate-down через docker compose run
```

---

## 🚀 Быстрый старт

```bash
# 1. Клонировать репозиторий
git clone https://github.com/KyoshiBlame/TodoKy.git
cd TodoKy

# 2. Настроить окружение
cp .env.example .env
# Заполнить значения в .env

# 3. Поднять PostgreSQL
make env-up

# 4. Применить миграции
make migrate-up

# 5. Запустить приложение
go run ./cmd/todoky

# Опционально: открыть порт PostgreSQL для GUI-клиентов
make env-port-forward
```

### Swagger UI

После запуска документация API доступна по адресу:

```
http://localhost:5050/swagger/index.html
```

> Также доступен на живом сервере: http://45.131.42.160:5050/swagger/index.html

---

## 📁 Структура переменных окружения

| Переменная | Описание | Пример |
|------------|----------|--------|
| `HTTP_ADDR` | Адрес и порт HTTP-сервера | `:5050` |
| `HTTP_SHUTDOWN_TIMEOUT` | Таймаут graceful shutdown | `30s` |
| `ALLOWED_ORIGINS` | Разрешённые CORS-источники | `http://localhost:3000` |
| `POSTGRES_HOST` | Хост PostgreSQL | `localhost` |
| `POSTGRES_USER` | Пользователь БД | `postgres` |
| `POSTGRES_PASSWORD` | Пароль БД | `secret` |
| `POSTGRES_DB` | Имя базы данных | `todoky` |
| `POSTGRES_TIMEOUT` | Таймаут подключения к БД | `10s` |
| `LOGGER_LEVEL` | Уровень логирования | `info` / `debug` |

---

## 📖 Ключевые решения

| Решение | Обоснование |
|---------|-------------|
| Без ORM, чистый `pgx` | Максимальный контроль над SQL, нет магии, производительность |
| `zap` вместо `log/slog` | Структурированные JSON-логи, высокая производительность, zero-allocation |
| `envconfig` | Минималистичная загрузка конфига из ENV без лишних зависимостей |
| Feature Sliced Structure | Масштабируется при росте команды, каждая фича независима |
| Ручной DI вместо фреймворка | Прозрачность, нет магии, легко дебажить граф зависимостей |
| `socat` как port-forwarder | Лёгкий и гибкий инструмент без overhead Nginx/HAProxy |

---

## 📝 Лицензия

MIT
