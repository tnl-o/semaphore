package example

import (
	"fmt"
	"net/http"
	"github.com/semaphoreui/semaphore/plugins"
)

// ExamplePlugin пример плагина
type ExamplePlugin struct {
	name        string
	version     string
	description string
}

// NewExamplePlugin создает новый экземпляр плагина
func NewExamplePlugin() *ExamplePlugin {
	return &ExamplePlugin{
		name:        "example",
		version:     "1.0.0",
		description: "Example plugin for Semaphore",
	}
}

// Name возвращает имя плагина
func (p *ExamplePlugin) Name() string {
	return p.name
}

// Version возвращает версию плагина
func (p *ExamplePlugin) Version() string {
	return p.version
}

// Description возвращает описание плагина
func (p *ExamplePlugin) Description() string {
	return p.description
}

// Initialize инициализирует плагин
func (p *ExamplePlugin) Initialize() error {
	// Здесь можно выполнить инициализацию плагина
	// Например, подключиться к базе данных, загрузить конфигурацию и т.д.
	return nil
}

// RegisterAPI регистрирует API endpoints плагина
func (p *ExamplePlugin) RegisterAPI(router interface{}) error {
	// Здесь можно зарегистрировать API endpoints
	// router должен быть совместим с используемым роутером (gorilla/mux, chi и т.д.)
	return nil
}

// RegisterWeb регистрирует web компоненты плагина
func (p *ExamplePlugin) RegisterWeb() error {
	// Здесь можно зарегистрировать web компоненты
	// Например, добавить Vue компоненты, стили и т.д.
	return nil
}

// GetMenuItems возвращает пункты меню для плагина
func (p *ExamplePlugin) GetMenuItems() []plugins.MenuItem {
	return []plugins.MenuItem{
		{
			ID:          "example",
			Title:       "Example Plugin",
			Icon:        "mdi-puzzle",
			Path:        "/plugins/example",
			Component:   "ExamplePlugin",
			Permissions: []string{},
			Order:       100,
		},
	}
}

// GetMiddleware возвращает middleware для плагина (не требуется)
func (p *ExamplePlugin) GetMiddleware() func(http.Handler) http.Handler {
	return nil
}

// init автоматически регистрирует плагин при импорте
func init() {
	if err := plugins.Register(NewExamplePlugin()); err != nil {
		panic(fmt.Sprintf("Failed to register example plugin: %v", err))
	}
}

