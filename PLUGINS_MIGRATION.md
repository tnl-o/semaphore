# Миграция функций из old_project в плагины

## Выполнено

### ✅ Плагин Cache (plugins/modules/cache)
- Функция очистки кэша из `old_project/api/cache.go`
- API endpoint: `DELETE /api/cache` (только для админов)
- Статус: Готов к использованию

### ✅ Плагин Events (plugins/modules/events)
- Улучшенная работа с событиями из `old_project/api/events.go`
- API endpoints:
  - `GET /api/events/last` - последние 200 событий
  - `GET /api/events/all` - все события
- Статус: Готов к использованию

### ⚠️ Плагин Metrics (plugins/modules/metrics)
- Метрики Prometheus из `old_project/api/middleware/metrics.go`
- Требует установки: `go get github.com/prometheus/client_golang/prometheus`
- Статус: Создан, но отключен по умолчанию (см. `plugins/modules/metrics/README.md`)

## В процессе

### 🔄 Плагин API Client (plugins/modules/api-client)
- Улучшенный API клиент с кэшированием из `old_project/web/src/lib/apiClient.js`
- Кэш API ответов из `old_project/web/src/lib/apiCache.js`
- Статус: Планируется

### 🔄 Плагин WebSocket (plugins/modules/websocket)
- Улучшенный WebSocket клиент из `old_project/web/src/lib/Socket.js`
- Статус: Планируется

## Правила работы

1. **Минимальные изменения в основном проекте** - все новые функции только через плагины
2. **Не изменять основной код** без крайней необходимости
3. **Все плагины регистрируются автоматически** при импорте
4. **Middleware из плагинов применяется автоматически** в `api/router.go`
5. **API endpoints регистрируются через `RegisterAPI()`** в каждом плагине

## Структура плагина

```
plugins/modules/your-plugin/
├── plugin.go          # Основной файл плагина
├── api.go             # API handlers (опционально)
├── middleware.go      # Middleware (опционально)
└── README.md          # Документация плагина
```

## Пример использования

См. `plugins/modules/example/` для базового примера плагина.

