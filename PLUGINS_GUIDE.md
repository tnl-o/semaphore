# Руководство по системе плагинов Semaphore

## Обзор

Система плагинов позволяет расширять функциональность Semaphore без изменения основного кода. Все новые функции должны создаваться как отдельные плагины.

## Структура плагина

```
plugins/
├── modules/              # Go модули плагинов
│   └── your-plugin/
│       └── plugin.go    # Основной файл плагина
├── api/                 # API endpoints для плагинов
│   └── plugins.go       # API для управления плагинами
└── web/                 # Frontend компоненты
    └── your-plugin/
        └── YourPlugin.vue
```

## Создание нового плагина

### 1. Создайте структуру плагина

```bash
mkdir -p plugins/modules/your-plugin
mkdir -p plugins/web/your-plugin
```

### 2. Реализуйте интерфейс Plugin

Создайте файл `plugins/modules/your-plugin/plugin.go`:

```go
package yourplugin

import (
    "fmt"
    "github.com/semaphoreui/semaphore/plugins"
)

type YourPlugin struct {
    name        string
    version     string
    description string
}

func NewYourPlugin() *YourPlugin {
    return &YourPlugin{
        name:        "your-plugin",
        version:     "1.0.0",
        description: "Описание вашего плагина",
    }
}

func (p *YourPlugin) Name() string {
    return p.name
}

func (p *YourPlugin) Version() string {
    return p.version
}

func (p *YourPlugin) Description() string {
    return p.description
}

func (p *YourPlugin) Initialize() error {
    // Инициализация плагина
    return nil
}

func (p *YourPlugin) RegisterAPI(router interface{}) error {
    // Регистрация API endpoints
    return nil
}

func (p *YourPlugin) RegisterWeb() error {
    // Регистрация web компонентов
    return nil
}

func (p *YourPlugin) GetMenuItems() []plugins.MenuItem {
    return []plugins.MenuItem{
        {
            ID:          "your-plugin",
            Title:       "Your Plugin",
            Icon:        "mdi-puzzle",
            Path:        "/plugins/your-plugin",
            Component:   "YourPlugin",
            Permissions: []string{},
            Order:       100,
        },
    }
}

func init() {
    if err := plugins.Register(NewYourPlugin()); err != nil {
        panic(fmt.Sprintf("Failed to register your-plugin: %v", err))
    }
}
```

### 3. Зарегистрируйте плагин в router.go

Добавьте импорт в `api/router.go`:

```go
_ "github.com/semaphoreui/semaphore/plugins/modules/your-plugin"
```

### 4. Создайте Vue компонент

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
          <!-- Ваш контент здесь -->
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script>
import EventBus from '@/event-bus';

export default {
  name: 'YourPlugin',
  methods: {
    showDrawer() {
      EventBus.$emit('i-show-drawer');
    },
  },
};
</script>
```

## API Endpoints

Плагины автоматически получают доступ к следующим API endpoints:

- `GET /api/plugins` - список всех плагинов
- `GET /api/plugins/menu` - пункты меню всех плагинов

## Порядок пунктов меню

Пункты меню сортируются по полю `Order`:
- 0-99: Основные пункты меню
- 100-999: Плагины
- 1000+: Разделитель и дополнительные пункты

## Примеры использования

См. `plugins/modules/example/` для примера базового плагина.

## Миграция старых функций

Старые функции из `old_project/` можно мигрировать в плагины:

1. Скопируйте код в `plugins/modules/your-feature/`
2. Создайте плагин, реализующий интерфейс `Plugin`
3. Зарегистрируйте плагин
4. Создайте Vue компоненты в `plugins/web/your-feature/`
5. Добавьте пункты меню через `GetMenuItems()`

