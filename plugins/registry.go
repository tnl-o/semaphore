package plugins

import (
	"fmt"
	"net/http"
	"sync"
)

// Plugin интерфейс для всех плагинов
type Plugin interface {
	// Name возвращает имя плагина
	Name() string
	
	// Version возвращает версию плагина
	Version() string
	
	// Description возвращает описание плагина
	Description() string
	
	// Initialize инициализирует плагин
	Initialize() error
	
	// RegisterAPI регистрирует API endpoints плагина
	RegisterAPI(router interface{}) error
	
	// RegisterWeb регистрирует web компоненты плагина
	RegisterWeb() error
	
	// GetMenuItems возвращает пункты меню для плагина
	GetMenuItems() []MenuItem
	
	// GetMiddleware возвращает middleware для применения к роутеру (опционально)
	// Может вернуть nil, если middleware не требуется
	GetMiddleware() func(http.Handler) http.Handler
}

// MenuItem представляет пункт меню плагина
type MenuItem struct {
	ID          string
	Title       string
	Icon        string
	Path        string
	Component   string
	Permissions []string
	Order       int
}

// Registry реестр плагинов
type Registry struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

var globalRegistry = &Registry{
	plugins: make(map[string]Plugin),
}

// Register регистрирует плагин
func Register(plugin Plugin) error {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()
	
	if _, exists := globalRegistry.plugins[plugin.Name()]; exists {
		return fmt.Errorf("plugin %s already registered", plugin.Name())
	}
	
	if err := plugin.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize plugin %s: %w", plugin.Name(), err)
	}
	
	globalRegistry.plugins[plugin.Name()] = plugin
	return nil
}

// Get возвращает плагин по имени
func Get(name string) (Plugin, error) {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	
	plugin, exists := globalRegistry.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	
	return plugin, nil
}

// GetAll возвращает все зарегистрированные плагины
func GetAll() []Plugin {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	
	plugins := make([]Plugin, 0, len(globalRegistry.plugins))
	for _, plugin := range globalRegistry.plugins {
		plugins = append(plugins, plugin)
	}
	
	return plugins
}

// GetAllMenuItems возвращает все пункты меню из всех плагинов
func GetAllMenuItems() []MenuItem {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()
	
	var items []MenuItem
	for _, plugin := range globalRegistry.plugins {
		items = append(items, plugin.GetMenuItems()...)
	}
	
	return items
}

