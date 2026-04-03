# Go EDI Processor

Микросервис для обработки EDI-документов (XML, JSON).

## Возможности

- Приём документов через gRPC и HTTP (REST)
- Поддержка форматов: XML, JSON
- Асинхронная обработка через Kafka
- Хранение документов в MongoDB
- Кэширование в Redis
- Отслеживание статусов документов
- Валидация документов
- Outbox-паттерн для надежной доставки событий
- OpenTelemetry для трейсинга

## Архитектура

Проект построен по принципам гексагональной архитектуры (Ports & Adapters):

```
┌─────────────────────────────────────────────┐
│                API Layer                     │
│  gRPC Controller │ HTTP Gateway (REST)      │
└──────────────────┼───────────────────────────┘
                   │
┌──────────────────▼───────────────────────────┐
│             Application Layer                 │
│           Document Service                    │
└──────────────────┬───────────────────────────┘
                   │
┌──────────────────▼───────────────────────────┐
│             Domain Layer                      │
│           Document Entity                     │
└──────────────────┬───────────────────────────┘
                   │
┌──────────────────▼───────────────────────────┐
│             Infrastructure Layer              │
│  MongoDB │ Redis │ Kafka │ PostgreSQL        │
└─────────────────────────────────────────────┘
```

## Структура проекта

```
├── api/proto/              # Protobuf определения
├── cmd/main.go             # Точка входа
├── internal/
│   ├── core/               # Ядро приложения
│   │   ├── domain/         # Доменные сущности
│   │   ├── ports/          # Интерфейсы (порты)
│   │   └── services/       # Бизнес-логика
│   └── adapters/           # Адаптеры
│       ├── primary/        # Входные адаптеры (контроллеры)
│       └── secondary/      # Выходные адаптеры (репозитории)
├── db/migrations/          # Миграции базы данных
└── docs/                   # Документация
```

## Стек

- **Язык**: Go 1.25+
- **gRPC**: Обработка запросов
- **HTTP Gateway**: REST API
- **Базы данных**: 
  - MongoDB (документы)
  - PostgreSQL (outbox, метаданные)
  - Redis (кэш)
- **Очереди**: Kafka (асинхронная обработка)
- **Валидация**: protoc-gen-validate
- **Логирование**: Zap
- **Трейсинг**: OpenTelemetry
- **Конфигурация**: Viper

## Быстрый старт

### Предварительные требования

- Go 1.25+
- Docker и Docker Compose
- Protobuf компилятор

### Запуск инфраструктуры

```bash
docker-compose up -d
```

### Генерация кода из Protobuf

```bash
make generate
```

### Сборка и запуск

```bash
make build
./edi-processor
```

### Тестирование

```bash
make test
make test-integration
```

## 📡 API

### gRPC

- `SendDocument` - отправка документа на обработку
- `GetDocumentByUUID` - получение документа по UUID
- `ReceiveDocument` - получение документов (возможно, для polling)

### REST (HTTP Gateway)

Доступно через Swagger UI: `http://localhost:8080/swagger/`

**HTTP эндпоинты:**
- `POST /api/v1/doc/send` - отправка документа
- `GET /api/v1/doc/get/{doc_id}` - получение документа по ID
- `GET /api/v1/doc/receive` - получение документов

## 🔧 Конфигурация

Настройки через переменные окружения или `.env` файл.

## 📄 Лицензия

MIT
