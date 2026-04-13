# go-doc-history

Микросервис для обработки xml/json документов

## Стек

- Go, gRPC, REST API, PostgreSQL, Redis, Kafka

## API

### gRPC

- `SendDocument` - отправка документа
- `GetDocumentByUUID` - получение документа по UUID
- `ReceiveDocument` - получение документов

**HTTP эндпоинты:**
- `POST /api/v1/doc/send` - отправка документа
- `GET /api/v1/doc/get/{doc_id}` - получение документа по UUID
- `GET /api/v1/doc/receive` - получение документов

## Конфигурация

Настройки через переменные окружения или `.env` файл.

## Лицензия

MIT
