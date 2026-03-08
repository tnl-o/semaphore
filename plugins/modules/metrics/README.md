# Metrics Plugin

Плагин для сбора метрик Prometheus.

## Требования

Для работы плагина требуется установить зависимость Prometheus:

```bash
go get github.com/prometheus/client_golang/prometheus
```

## Активация

Раскомментируйте импорт в `api/router.go`:

```go
_ "github.com/semaphoreui/semaphore/plugins/modules/metrics"
```

## Использование

После активации плагин автоматически собирает метрики:
- `http_request_duration_seconds` - длительность HTTP запросов
- `http_requests_total` - общее количество HTTP запросов
- `http_requests_in_flight` - количество активных запросов

Метрики доступны через стандартный endpoint `/metrics` (если настроен Prometheus handler).

