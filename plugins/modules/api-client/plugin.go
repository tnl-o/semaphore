package apiclient

import (
	"fmt"
	"net/http"

	"github.com/semaphoreui/semaphore/plugins"
)

// APIClientPlugin плагин для улучшенного API клиента с кэшированием
type APIClientPlugin struct {
	name        string
	version     string
	description string
}

// NewAPIClientPlugin создает новый экземпляр плагина API клиента
func NewAPIClientPlugin() *APIClientPlugin {
	return &APIClientPlugin{
		name:        "api-client",
		version:     "1.0.0",
		description: "Enhanced API client with caching and error handling",
	}
}

// Name возвращает имя плагина
func (p *APIClientPlugin) Name() string {
	return p.name
}

// Version возвращает версию плагина
func (p *APIClientPlugin) Version() string {
	return p.version
}

// Description возвращает описание плагина
func (p *APIClientPlugin) Description() string {
	return p.description
}

// Initialize инициализирует плагин
func (p *APIClientPlugin) Initialize() error {
	return nil
}

// RegisterAPI регистрирует API endpoints плагина
func (p *APIClientPlugin) RegisterAPI(router interface{}) error {
	// API клиент работает на frontend, не требует backend endpoints
	return nil
}

// RegisterWeb регистрирует web компоненты плагина
func (p *APIClientPlugin) RegisterWeb() error {
	// Frontend компоненты будут загружены автоматически через систему плагинов
	return nil
}

// GetMenuItems возвращает пункты меню для плагина
func (p *APIClientPlugin) GetMenuItems() []plugins.MenuItem {
	// API клиент не требует пунктов меню
	return []plugins.MenuItem{}
}

// GetMiddleware возвращает middleware для плагина (не требуется)
func (p *APIClientPlugin) GetMiddleware() func(http.Handler) http.Handler {
	return nil
}

func init() {
	if err := plugins.Register(NewAPIClientPlugin()); err != nil {
		panic(fmt.Sprintf("Failed to register api-client plugin: %v", err))
	}
}

