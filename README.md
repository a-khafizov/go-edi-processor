# go edi processor

Микросервис для обработки EDI-документов (XML, JSON)

## Стек

- Go 1.25+, gRPC, HTTP, MongoDB, PostgreSQL, Redis, Kafka, Zap, Viper

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
