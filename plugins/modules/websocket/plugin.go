package websocket

import (
	"fmt"
	"net/http"

	"github.com/semaphoreui/semaphore/plugins"
)

// WebSocketPlugin плагин для улучшенного WebSocket клиента
type WebSocketPlugin struct {
	name        string
	version     string
	description string
}

// NewWebSocketPlugin создает новый экземпляр плагина WebSocket
func NewWebSocketPlugin() *WebSocketPlugin {
	return &WebSocketPlugin{
		name:        "websocket",
		version:     "1.0.0",
		description: "Enhanced WebSocket client with automatic reconnection and error handling",
	}
}

// Name возвращает имя плагина
func (p *WebSocketPlugin) Name() string {
	return p.name
}

// Version возвращает версию плагина
func (p *WebSocketPlugin) Version() string {
	return p.version
}

// Description возвращает описание плагина
func (p *WebSocketPlugin) Description() string {
	return p.description
}

// Initialize инициализирует плагин
func (p *WebSocketPlugin) Initialize() error {
	return nil
}

// RegisterAPI регистрирует API endpoints плагина
func (p *WebSocketPlugin) RegisterAPI(router interface{}) error {
	// WebSocket работает на frontend, не требует backend endpoints
	return nil
}

// RegisterWeb регистрирует web компоненты плагина
func (p *WebSocketPlugin) RegisterWeb() error {
	// Frontend компоненты будут загружены автоматически через систему плагинов
	return nil
}

// GetMenuItems возвращает пункты меню для плагина
func (p *WebSocketPlugin) GetMenuItems() []plugins.MenuItem {
	// WebSocket не требует пунктов меню
	return []plugins.MenuItem{}
}

// GetMiddleware возвращает middleware для плагина (не требуется)
func (p *WebSocketPlugin) GetMiddleware() func(http.Handler) http.Handler {
	return nil
}

func init() {
	if err := plugins.Register(NewWebSocketPlugin()); err != nil {
		panic(fmt.Sprintf("Failed to register websocket plugin: %v", err))
	}
}

