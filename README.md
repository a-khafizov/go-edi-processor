# go edi processor

Микросервис для обработки EDI-документов (XML, JSON).

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
│       ├── input/        # Входные адаптеры (контроллеры)
│       └── output/      # Выходные адаптеры (репозитории)
├── db/migrations/          # Миграции базы данных
└── docs/                   # Документация
```

## Стек

- **Язык**: Go 1.25+
- **gRPC**
- **HTTP Gateway**
- **Базы данных**: 
  - MongoDB (документы)
  - PostgreSQL (outbox, documents)
  - Redis (кэш)
- **Очереди**: Kafka (асинхронная обработка)
- **Валидация**: protoc-gen-validate
- **Логирование**: Zap
- **Конфигурация**: Viper

## API

### gRPC

- `SendDocument` - отправка документа на обработку
- `GetDocumentByUUID` - получение документа по UUID
- `ReceiveDocument` - получение документов

**HTTP эндпоинты:**
- `POST /api/v1/doc/send` - отправка документа
- `GET /api/v1/doc/get/{doc_id}` - получение документа по ID
- `GET /api/v1/doc/receive` - получение документов

## Конфигурация

Настройки через переменные окружения или `.env` файл.

## Лицензия

MIT
