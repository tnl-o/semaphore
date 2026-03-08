# Semaphore Plugins System

Эта директория содержит систему плагинов для расширения функциональности Semaphore.

## Структура

```
plugins/
├── modules/          # Модули плагинов (Go)
│   ├── example/      # Пример плагина
│   ├── cache/        # Плагин управления кэшем
│   ├── events/       # Плагин работы с событиями
│   ├── metrics/      # Плагин метрик Prometheus (требует prometheus)
│   ├── api-client/   # Плагин улучшенного API клиента
│   └── websocket/    # Плагин улучшенного WebSocket
├── api/             # API endpoints для плагинов
└── web/             # Frontend компоненты плагинов
    ├── api-client/  # Компоненты API клиента
    └── websocket/   # Компоненты WebSocket
```

## Установленные плагины

### ✅ Cache Plugin
- Управление кэшем приложения
- API: `DELETE /api/cache` (только для админов)
- Статус: Активен

### ✅ Events Plugin
- Улучшенная работа с событиями
- API: `GET /api/events/last`, `GET /api/events/all`
- Статус: Активен

### ⚠️ Metrics Plugin
- Метрики Prometheus
- Требует: `go get github.com/prometheus/client_golang/prometheus`
- Статус: Создан, но отключен по умолчанию
- См. `modules/metrics/README.md` для активации

### ✅ API Client Plugin
- Улучшенный API клиент с кэшированием
- Автоматическая обработка ошибок
- Retry логика
- Статус: Активен (frontend компоненты)

### ✅ WebSocket Plugin
- Улучшенный WebSocket клиент
- Автоматическое переподключение
- Обработка ошибок авторизации
- Статус: Активен (frontend компоненты)

## Создание нового плагина

1. Создайте директорию в `plugins/modules/your-plugin/`
2. Реализуйте интерфейс `Plugin` из `plugins/registry.go`
3. Зарегистрируйте плагин в `api/router.go` (добавьте импорт)
4. Добавьте frontend компоненты в `plugins/web/your-plugin/`
5. Добавьте маршруты в `plugins/api/your-plugin/` (если нужны)

## Пример плагина

См. `plugins/modules/example/` для примера базового плагина.

## Правила работы

1. **Минимальные изменения в основном проекте** - все новые функции только через плагины
2. **Не изменять основной код** без крайней необходимости
3. **Все плагины регистрируются автоматически** при импорте
4. **Middleware из плагинов применяется автоматически** в `api/router.go`
5. **API endpoints регистрируются через `RegisterAPI()`** в каждом плагине
