# Трассировка в go-edi-document-processor

## Что такое трассировка?

Трассировка (tracing) — это метод наблюдения за выполнением запроса в распределённой системе.
Каждый запрос получает уникальный `trace_id`, который проходит через все компоненты системы.
Внутри запроса могут создаваться `span`'ы (промежутки) для операций (например, вызов gRPC, работа с БД, обработка документа).

## Как работает в нашем проекте

Мы используем OpenTelemetry — стандарт для инструментирования приложений.

### Компоненты

1. **TracerProvider** — глобальный менеджер трассировки, создаёт и экспортирует spans.
2. **Exporter** — отправляет spans во внешнюю систему (у нас `stdout` для простоты).
3. **Interceptor** — перехватывает gRPC вызовы, создаёт span для каждого вызова.
4. **Logger** — логирует с привязкой к `trace_id`.

### Поток данных

```
Клиент
  |
  | HTTP/gRPC запрос
  v
[gRPC сервер]
  |
  | RecoveryInterceptor (ловит паники)
  |
  | TracingInterceptor (создаёт span, добавляет в контекст)
  |   - trace_id = aa9d8e8964b720a85e2ba9872a966c3b
  |   - span_id = 504e47b8596a3b93
  |
  | LoggingInterceptor (логирует с trace_id)
  |   - пишет: "gRPC method called", trace_id=...
  |
  | Обработчик (DocumentService)
  |   - бизнес-логика
  |
  | Span завершается
  |
  | Exporter (stdout) выводит span в консоль
  v
Ответ клиенту
```

### Пример span (из stdout)

```json
{
  "Name": "/go_edi_document_processor.DocumentService/SendDocument",
  "SpanContext": {
    "TraceID": "aa9d8e8964b720a85e2ba9872a966c3b",
    "SpanID": "504e47b8596a3b93",
    ...
  },
  "Attributes": [
    {"Key": "rpc.method", "Value": "/go_edi_document_processor.DocumentService/SendDocument"},
    {"Key": "rpc.system", "Value": "grpc"}
  ]
}
```

### Как добавить свою трассировку в сервисах

Если нужно отследить операцию внутри `document_service` или `outbox_processor`, используйте функцию `tracing.StartSpan`:

```go
import "github.com/go-edi-document-processor/internal/bootstrap/tracing"

func (s *DocumentService) SomeMethod(ctx context.Context) error {
    ctx, span := tracing.StartSpan(ctx, "operation_name")
    defer span.End()
    // ... ваш код
}
```

Это создаст дочерний span, который будет связан с родительским trace_id.

### Настройка

Трассировка инициализируется в `cmd/main.go`:

```go
tracing.InitTracing("go-edi-document-processor")
```

Экспортер настроен на вывод в консоль (`stdout`). В продакшене можно заменить на Jaeger, Zipkin, Grafana Tempo и т.д.

### Логи с trace_id

В логах gRPC методов теперь есть поле `trace_id`. Пример:

```
{"level":"INFO","time":"...","caller":"interceptors.go:45","msg":"gRPC method called","method":"/go_edi_document_processor.DocumentService/SendDocument","trace_id":"aa9d8e8964b720a85e2ba9872a966c3b"}
```

Это позволяет связать логи с конкретным запросом.

## Визуализация в Jaeger/Grafana

Если подключить Jaeger, можно увидеть графическое представление spans на временной шкале.

1. Установите Jaeger (например, через Docker).
2. Замените экспортер на `jaeger` или `otlp`.
3. Запустите приложение и отправьте запрос.
4. Откройте Jaeger UI (http://localhost:16686) и найдите trace по trace_id.

## Полезные ссылки

- [OpenTelemetry Go](https://opentelemetry.io/docs/instrumentation/go/)
- [Jaeger](https://www.jaegertracing.io/)
- [Визуализация трассировок](https://opentelemetry.io/docs/concepts/observability-primer/#distributed-traces)

## Заключение

Трассировка даёт полную картину выполнения запроса, помогает находить узкие места и debug в production.
Начальная реализация готова, можно расширять добавлением spans в ключевые операции.