# department-management-service

Сервис управления организационной структурой на Go.

> В проекте реализован сервис по управлению отделами и персоналом с запуском через `docker-compose`.

---

## Описание

- Сервис предоставляет HTTP API для управления деревом подразделений (создание, перемещение, удаление), управления сотрудниками внутри подразделений, получения подразделения с поддеревом и списком сотрудников.
- Архитектура: разделение на domain / usecase / adapters / transport, работа с БД через GORM, миграции через goose, контейнеризация — Docker.
- Конфигурация через `configs/config.yaml` и переменные окружения (.env.example).

---

## Технологический стек

- Go 1.24.3
- PostgreSQL
- Docker / docker-compose
- OpenAPI
- GORM
- goose
- mockery
- ogen
- Swagger UI

---

## Быстрый обзор структуры (файлы/директории важные для запуска)

- `cmd/market-parser/main.go` — точка входа, инициализация конфигурации и сервисов.
- `configs/config.yaml` — основной YAML-конфиг (base_url, selectors и т.д.).
- `.env.example` — пример переменных окружения (WebSocket URL браузера, прокси и т.п.).
- `internal/adapters/browser/chromium` — инициализация Chromium/rod, прокси, DTO и утилиты для работы с элементами страницы.
- `internal/adapters/parsers/kuper.go` — пример реализации парсера для магазина.
- `internal/usecase/parser_service.go` — сервисная логика: валидация параметров и вызов репозитория-парсера.
- `internal/transport/http` — HTTP-обёртка, handlers и OpenAPI-строка.
- `Dockerfile`, `docker-compose.yaml`, `Makefile` — способы запуска.

---

## Структура проекта

```
cmd/                # точка входа приложения (main.go)
config/             # конфиг (config.yaml)
internal/           # внутренние пакеты (/usecase, /domain, /repository, /adapters, /tests)
internal/adapters   # хранит адаптеры, в частности реализацию интерфейса репозитория Postgres
internal/config     # код для сбора конфига
internal/domain     # доменные сущности и ошибки
internal/repository # интерфейсы (порты)
internal/test       # сгенерированные моки при помощи утилиты Mockery
internal/transport  # трнаспортный уровень (http + сгенерированные файлы при помощи утилиты ogen)
internal/usecase    # бизнес-логика (сервисы)
migrations/         # SQL миграции
pkg/                # внение пакеты (logger на log/slog)
.env/.env.example   # файл с переменными окружения
dockerfile          # образ приложения (go + alpine)
docker-compose.yml  # docker-compose для поднятия приложения, базы данных, контейнера с миграциями и swagger ui
```

---

## Запуск приложения

1. Подготовьте `.env` на основе `.env.example` и `config.yaml`.

2. Запустите compose:

```bash
docker-compose up --build
```

3. После поднятия сервисов API будет доступен на порту, указанном в `configs/config.yaml` / `.env` (по умолчанию `localhost:8080`).

---

## Использование API

Для использования API, после запуска docker-compose пройдите по адресу http://localhost:8081 (swagger-ui) для тестирования отправки запросов в приложение.

---
