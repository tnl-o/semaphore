# План реализации модулей/плагинов для Semaphore

## Цель документа

Этот документ предназначен для ИИ-ассистента, который будет реализовывать новые модули или плагины для проекта Semaphore. Документ содержит пошаговую инструкцию, архитектуру системы и примеры кода.

## Архитектура системы плагинов

### Основные принципы

1. **Минимальные изменения в основном проекте** - все новые функции реализуются только через плагины
2. **Автоматическая регистрация** - плагины регистрируются автоматически при импорте через `init()`
3. **Разделение ответственности** - каждый плагин независим и может быть включен/отключен
4. **Единый интерфейс** - все плагины реализуют интерфейс `plugins.Plugin`

### Структура директорий

```
plugins/
├── registry.go              # Интерфейс Plugin и реестр плагинов
├── api/                     # API endpoints для управления плагинами
│   └── plugins.go           # GET /api/plugins, GET /api/plugins/menu
├── modules/                 # Go модули плагинов (бэкенд)
│   ├── example/            # Пример плагина
│   │   └── plugin.go       # Реализация интерфейса Plugin
│   ├── cache/              # Плагин управления кэшем
│   ├── events/             # Плагин работы с событиями
│   └── your-plugin/        # Ваш новый плагин
│       └── plugin.go
└── web/                     # Frontend компоненты плагинов
    ├── api-client/         # JavaScript модули
    │   ├── apiClient.js
    │   └── apiCache.js
    └── your-plugin/        # Vue компоненты вашего плагина
        └── YourPlugin.vue
```

## Интерфейс Plugin

### Определение интерфейса

```go
// plugins/registry.go
type Plugin interface {
    // Name возвращает уникальное имя плагина
    Name() string
    
    // Version возвращает версию плагина
    Version() string
    
    // Description возвращает описание плагина
    Description() string
    
    // Initialize инициализирует плагин (подключение к БД, загрузка конфига и т.д.)
    Initialize() error
    
    // RegisterAPI регистрирует API endpoints плагина
    // router имеет тип *mux.Router из github.com/gorilla/mux
    RegisterAPI(router interface{}) error
    
    // RegisterWeb регистрирует web компоненты (пока не используется)
    RegisterWeb() error
    
    // GetMenuItems возвращает пункты меню для отображения в сайдбаре
    GetMenuItems() []MenuItem
    
    // GetMiddleware возвращает HTTP middleware (опционально, может вернуть nil)
    GetMiddleware() func(http.Handler) http.Handler
}
```

### Структура MenuItem

```go
type MenuItem struct {
    ID          string   // Уникальный идентификатор пункта меню
    Title       string   // Отображаемое название
    Icon        string   // Иконка Material Design Icons (например, "mdi-puzzle")
    Path        string   // Путь маршрута (например, "/plugins/your-plugin")
    Component   string   // Имя Vue компонента (без расширения .vue)
    Permissions []string // Требуемые права доступа (пока не используется)
    Order       int      // Порядок сортировки (0-99: основные, 100-999: плагины, 1000+: разделители)
}
```

## Пошаговая инструкция создания плагина

### Шаг 1: Создание структуры директорий

```bash
# Создайте директорию для Go модуля
mkdir -p plugins/modules/your-plugin

# Создайте директорию для Vue компонентов (если нужен frontend)
mkdir -p plugins/web/your-plugin
```

### Шаг 2: Создание основного файла плагина

Создайте файл `plugins/modules/your-plugin/plugin.go`:

```go
package yourplugin

import (
    "fmt"
    "net/http"
    
    "github.com/gorilla/mux"
    "github.com/semaphoreui/semaphore/plugins"
)

// YourPlugin структура плагина
type YourPlugin struct {
    name        string
    version     string
    description string
}

// NewYourPlugin создает новый экземпляр плагина
func NewYourPlugin() *YourPlugin {
    return &YourPlugin{
        name:        "your-plugin",
        version:     "1.0.0",
        description: "Описание вашего плагина",
    }
}

// Name возвращает имя плагина
func (p *YourPlugin) Name() string {
    return p.name
}

// Version возвращает версию плагина
func (p *YourPlugin) Version() string {
    return p.version
}

// Description возвращает описание плагина
func (p *YourPlugin) Description() string {
    return p.description
}

// Initialize инициализирует плагин
func (p *YourPlugin) Initialize() error {
    // Здесь можно:
    // - Подключиться к базе данных
    // - Загрузить конфигурацию
    // - Инициализировать зависимости
    // - Проверить доступность ресурсов
    
    // Пример:
    // if err := p.loadConfig(); err != nil {
    //     return fmt.Errorf("failed to load config: %w", err)
    // }
    
    return nil
}

// RegisterAPI регистрирует API endpoints плагина
func (p *YourPlugin) RegisterAPI(router interface{}) error {
    // Приведение типа к *mux.Router
    r, ok := router.(*mux.Router)
    if !ok {
        return fmt.Errorf("router is not *mux.Router")
    }
    
    // Регистрация API endpoints
    // Пример: GET /api/your-plugin/items
    r.Path("/api/your-plugin/items").
        HandlerFunc(p.handleGetItems).
        Methods("GET", "HEAD")
    
    // Пример: POST /api/your-plugin/items
    r.Path("/api/your-plugin/items").
        HandlerFunc(p.handleCreateItem).
        Methods("POST")
    
    return nil
}

// RegisterWeb регистрирует web компоненты (пока не используется)
func (p *YourPlugin) RegisterWeb() error {
    // В будущем здесь можно будет регистрировать Vue компоненты
    return nil
}

// GetMenuItems возвращает пункты меню для плагина
func (p *YourPlugin) GetMenuItems() []plugins.MenuItem {
    return []plugins.MenuItem{
        {
            ID:          "your-plugin-main",
            Title:       "Your Plugin",
            Icon:        "mdi-puzzle",
            Path:        "/plugins/your-plugin",
            Component:   "YourPlugin",
            Permissions: []string{},
            Order:       100, // Плагины обычно имеют Order >= 100
        },
        // Можно добавить несколько пунктов меню
        {
            ID:          "your-plugin-settings",
            Title:       "Plugin Settings",
            Icon:        "mdi-cog",
            Path:        "/plugins/your-plugin/settings",
            Component:   "YourPluginSettings",
            Permissions: []string{"admin"},
            Order:       101,
        },
    }
}

// GetMiddleware возвращает HTTP middleware (опционально)
func (p *YourPlugin) GetMiddleware() func(http.Handler) http.Handler {
    // Если middleware не требуется, вернуть nil
    return nil
    
    // Пример middleware:
    // return func(next http.Handler) http.Handler {
    //     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    //         // Логика middleware
    //         next.ServeHTTP(w, r)
    //     })
    // }
}

// Обработчики HTTP запросов

func (p *YourPlugin) handleGetItems(w http.ResponseWriter, r *http.Request) {
    // Получение данных из контекста запроса (если нужно)
    // store := helpers.GetContextValue(r, "store").(db.Store)
    
    // Ваша логика обработки запроса
    items := []map[string]interface{}{
        {"id": 1, "name": "Item 1"},
        {"id": 2, "name": "Item 2"},
    }
    
    // Отправка ответа
    // helpers.WriteJSON(w, http.StatusOK, items)
}

func (p *YourPlugin) handleCreateItem(w http.ResponseWriter, r *http.Request) {
    // Ваша логика создания элемента
    // helpers.WriteJSON(w, http.StatusCreated, result)
}

// init автоматически регистрирует плагин при импорте пакета
func init() {
    if err := plugins.Register(NewYourPlugin()); err != nil {
        panic(fmt.Sprintf("Failed to register your-plugin: %v", err))
    }
}
```

### Шаг 3: Регистрация плагина в основном роутере

Добавьте импорт в файл `api/router.go`:

```go
import (
    // ... другие импорты
    
    // Импорт плагинов (автоматическая регистрация через init())
    _ "github.com/semaphoreui/semaphore/plugins/modules/your-plugin"
)
```

**Важно:** Импорт должен быть с `_` (blank import), чтобы выполнилась функция `init()` плагина.

### Шаг 4: Создание Vue компонента (если нужен frontend)

Создайте файл `plugins/web/your-plugin/YourPlugin.vue`:

```vue
<template>
  <div>
    <v-toolbar flat>
      <v-app-bar-nav-icon @click="showDrawer()"></v-app-bar-nav-icon>
      <v-toolbar-title>Your Plugin</v-toolbar-title>
    </v-toolbar>
    <v-divider/>
    <v-container>
      <v-card>
        <v-card-title>Your Plugin Content</v-card-title>
        <v-card-text>
          <v-data-table
            :headers="headers"
            :items="items"
            :loading="loading"
          >
            <template v-slot:item.actions="{ item }">
              <v-btn small @click="editItem(item)">Edit</v-btn>
            </template>
          </v-data-table>
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script>
import EventBus from '@/event-bus';
import apiClient from '@/lib/apiClient';

export default {
  name: 'YourPlugin',
  data() {
    return {
      loading: false,
      items: [],
      headers: [
        { text: 'ID', value: 'id' },
        { text: 'Name', value: 'name' },
        { text: 'Actions', value: 'actions', sortable: false },
      ],
    };
  },
  created() {
    this.loadItems();
  },
  methods: {
    showDrawer() {
      EventBus.$emit('i-show-drawer');
    },
    async loadItems() {
      this.loading = true;
      try {
        const response = await apiClient.get('/api/your-plugin/items');
        this.items = response.data;
      } catch (error) {
        console.error('Failed to load items:', error);
      } finally {
        this.loading = false;
      }
    },
    editItem(item) {
      // Логика редактирования
    },
  },
};
</script>
```

### Шаг 5: Регистрация маршрута во frontend (если нужен)

Маршруты плагинов загружаются динамически через API `/api/plugins/menu`. Компоненты должны быть доступны через динамический импорт.

**Примечание:** В текущей реализации динамическая загрузка компонентов плагинов отключена. Компоненты можно добавить вручную в `web/src/router/index.js`:

```javascript
import YourPlugin from '@/plugins/web/your-plugin/YourPlugin.vue';

// В массиве routes:
{
  path: '/plugins/your-plugin',
  component: YourPlugin,
}
```

## Примеры реализации

### Пример 1: Простой плагин без API

```go
package simple

import (
    "fmt"
    "github.com/semaphoreui/semaphore/plugins"
)

type SimplePlugin struct {
    name string
}

func NewSimplePlugin() *SimplePlugin {
    return &SimplePlugin{name: "simple"}
}

func (p *SimplePlugin) Name() string { return p.name }
func (p *SimplePlugin) Version() string { return "1.0.0" }
func (p *SimplePlugin) Description() string { return "Simple plugin" }
func (p *SimplePlugin) Initialize() error { return nil }
func (p *SimplePlugin) RegisterAPI(router interface{}) error { return nil }
func (p *SimplePlugin) RegisterWeb() error { return nil }
func (p *SimplePlugin) GetMenuItems() []plugins.MenuItem { return nil }
func (p *SimplePlugin) GetMiddleware() func(http.Handler) http.Handler { return nil }

func init() {
    plugins.Register(NewSimplePlugin())
}
```

### Пример 2: Плагин с API endpoints

```go
func (p *YourPlugin) RegisterAPI(router interface{}) error {
    r, ok := router.(*mux.Router)
    if !ok {
        return fmt.Errorf("router is not *mux.Router")
    }
    
    // Защищенные endpoints (требуют аутентификации)
    authenticatedAPI := r.PathPrefix("/api").Subrouter()
    authenticatedAPI.Use(authMiddleware) // Если нужно
    
    authenticatedAPI.Path("/your-plugin/data").
        HandlerFunc(p.handleGetData).
        Methods("GET")
    
    authenticatedAPI.Path("/your-plugin/data").
        HandlerFunc(p.handleCreateData).
        Methods("POST")
    
    // Админские endpoints
    adminAPI := authenticatedAPI.PathPrefix("/admin").Subrouter()
    adminAPI.Use(adminMiddleware) // Если нужно
    
    adminAPI.Path("/your-plugin/config").
        HandlerFunc(p.handleGetConfig).
        Methods("GET")
    
    return nil
}
```

### Пример 3: Плагин с middleware

```go
func (p *YourPlugin) GetMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Логика middleware
            // Например, логирование, метрики, проверка прав и т.д.
            
            start := time.Now()
            next.ServeHTTP(w, r)
            duration := time.Since(start)
            
            // Логирование или отправка метрик
            log.Printf("Request to %s took %v", r.URL.Path, duration)
        })
    }
}
```

### Пример 4: Плагин с доступом к базе данных

```go
func (p *YourPlugin) handleGetData(w http.ResponseWriter, r *http.Request) {
    // Получение store из контекста
    store := helpers.GetContextValue(r, "store").(db.Store)
    
    // Использование store для работы с БД
    items, err := store.GetYourPluginItems()
    if err != nil {
        helpers.WriteError(w, err)
        return
    }
    
    helpers.WriteJSON(w, http.StatusOK, items)
}
```

## Интеграция с основным приложением

### Автоматическая регистрация плагинов

Плагины регистрируются автоматически при импорте в `api/router.go`:

```go
import (
    // ... другие импорты
    
    // Плагины (автоматическая регистрация)
    _ "github.com/semaphoreui/semaphore/plugins/modules/example"
    _ "github.com/semaphoreui/semaphore/plugins/modules/cache"
    _ "github.com/semaphoreui/semaphore/plugins/modules/events"
    _ "github.com/semaphoreui/semaphore/plugins/modules/your-plugin" // Ваш плагин
)
```

### Применение middleware из плагинов

В `api/router.go` middleware из плагинов применяется автоматически:

```go
// Применяем middleware из плагинов
for _, plugin := range plugins.GetAll() {
    if middleware := plugin.GetMiddleware(); middleware != nil {
        r.Use(middleware)
    }
}
```

### Регистрация API endpoints плагинов

В `api/router.go` вызывается `RegisterAPI` для каждого плагина:

```go
// Регистрация API endpoints плагинов
for _, plugin := range plugins.GetAll() {
    if err := plugin.RegisterAPI(r); err != nil {
        log.Errorf("Failed to register API for plugin %s: %v", plugin.Name(), err)
    }
}
```

### API endpoints для управления плагинами

Система автоматически предоставляет следующие endpoints:

- `GET /api/plugins` - список всех зарегистрированных плагинов
- `GET /api/plugins/menu` - все пункты меню из всех плагинов

## Frontend интеграция

### Загрузка пунктов меню

В `web/src/App.vue` пункты меню загружаются динамически:

```javascript
async loadPluginMenuItems() {
  try {
    const response = await apiClient.get('/api/plugins/menu');
    this.pluginMenuItems = response.data;
  } catch (error) {
    console.error('Failed to load plugin menu items:', error);
    this.pluginMenuItems = [];
  }
}
```

### Отображение пунктов меню

Пункты меню отображаются в сайдбаре:

```vue
<template v-for="item in pluginMenuItems">
  <v-divider v-if="item.order >= 1000" :key="`divider-${item.id}`"></v-divider>
  <v-list-item :to="item.path" :key="item.id">
    <v-list-item-icon>
      <v-icon>{{ item.icon }}</v-icon>
    </v-list-item-icon>
    <v-list-item-content>
      {{ item.title }}
    </v-list-item-content>
  </v-list-item>
</template>
```

## Лучшие практики

### 1. Именование

- Имя плагина должно быть уникальным и в формате `kebab-case` (например, `your-plugin`)
- Имя пакета Go должно совпадать с именем плагина
- Vue компоненты должны быть в PascalCase (например, `YourPlugin.vue`)

### 2. Обработка ошибок

- Всегда возвращайте ошибки из `Initialize()` и `RegisterAPI()`
- Используйте `helpers.WriteError()` для отправки ошибок клиенту
- Логируйте ошибки с помощью `log.Errorf()`

### 3. Доступ к ресурсам

- Используйте контекст запроса для получения `store`, `task_pool` и других ресурсов
- Не создавайте глобальные переменные для доступа к ресурсам

### 4. Тестирование

- Создавайте unit-тесты для логики плагина
- Тестируйте API endpoints с помощью `httptest`
- Используйте моки для зависимостей

### 5. Документация

- Добавляйте комментарии к публичным методам
- Создавайте README.md для сложных плагинов
- Документируйте API endpoints

## Миграция существующего кода в плагин

Если нужно мигрировать существующую функцию в плагин:

1. **Скопируйте код** из `old_project/` в `plugins/modules/your-feature/`
2. **Создайте структуру плагина** с реализацией интерфейса `Plugin`
3. **Перенесите API endpoints** в метод `RegisterAPI()`
4. **Перенесите Vue компоненты** в `plugins/web/your-feature/`
5. **Добавьте пункты меню** через `GetMenuItems()`
6. **Удалите старый код** из основного проекта
7. **Зарегистрируйте плагин** в `api/router.go`

## Примеры существующих плагинов

### Cache Plugin (`plugins/modules/cache/`)

- Управление кэшем приложения
- API: `DELETE /api/cache` (только для админов)
- Без frontend компонентов

### Events Plugin (`plugins/modules/events/`)

- Работа с событиями
- API: `GET /api/events/last`, `GET /api/events/all`
- Без frontend компонентов

### API Client Plugin (`plugins/modules/api-client/`)

- Улучшенный API клиент с кэшированием
- Frontend компоненты: `plugins/web/api-client/apiClient.js`, `apiCache.js`
- Автоматическая обработка ошибок и retry логика

### WebSocket Plugin (`plugins/modules/websocket/`)

- Улучшенный WebSocket клиент
- Frontend компоненты: `plugins/web/websocket/Socket.js`
- Автоматическое переподключение

## Чеклист создания плагина

- [ ] Создана директория `plugins/modules/your-plugin/`
- [ ] Создан файл `plugin.go` с реализацией интерфейса `Plugin`
- [ ] Реализованы все обязательные методы интерфейса
- [ ] Добавлена функция `init()` для автоматической регистрации
- [ ] Добавлен импорт в `api/router.go`
- [ ] Созданы API endpoints (если нужны)
- [ ] Созданы Vue компоненты (если нужны)
- [ ] Добавлены пункты меню через `GetMenuItems()`
- [ ] Протестирована работа плагина
- [ ] Добавлена документация

## Решение проблем

### Плагин не регистрируется

- Проверьте, что импорт добавлен в `api/router.go` с `_`
- Убедитесь, что функция `init()` вызывает `plugins.Register()`
- Проверьте логи на наличие ошибок при инициализации

### API endpoints не работают

- Убедитесь, что `RegisterAPI()` правильно приводит тип к `*mux.Router`
- Проверьте, что endpoints регистрируются на правильном роутере
- Убедитесь, что пути не конфликтуют с существующими

### Пункты меню не отображаются

- Проверьте, что `GetMenuItems()` возвращает непустой массив
- Убедитесь, что `Order` установлен правильно (>= 100 для плагинов)
- Проверьте, что frontend загружает пункты меню из `/api/plugins/menu`

### Middleware не применяется

- Убедитесь, что `GetMiddleware()` возвращает не `nil`
- Проверьте, что middleware применяется в `api/router.go`
- Убедитесь, что middleware имеет правильную сигнатуру

## Дополнительные ресурсы

- `plugins/registry.go` - определение интерфейса `Plugin`
- `plugins/modules/example/plugin.go` - пример базового плагина
- `plugins/modules/cache/plugin.go` - пример плагина с API
- `plugins/api/plugins.go` - API для управления плагинами
- `api/router.go` - регистрация плагинов в основном роутере
- `web/src/App.vue` - отображение пунктов меню плагинов

