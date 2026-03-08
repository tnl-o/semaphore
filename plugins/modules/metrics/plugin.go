package metrics

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/semaphoreui/semaphore/plugins"
	"github.com/semaphoreui/semaphore/plugins/modules/metrics/middleware"
)

// MetricsPlugin плагин для метрик Prometheus
type MetricsPlugin struct {
	name        string
	version     string
	description string
}

// NewMetricsPlugin создает новый экземпляр плагина метрик
func NewMetricsPlugin() *MetricsPlugin {
	return &MetricsPlugin{
		name:        "metrics",
		version:     "1.0.0",
		description: "Prometheus metrics middleware for HTTP requests",
	}
}

// Name возвращает имя плагина
func (p *MetricsPlugin) Name() string {
	return p.name
}

// Version возвращает версию плагина
func (p *MetricsPlugin) Version() string {
	return p.version
}

// Description возвращает описание плагина
func (p *MetricsPlugin) Description() string {
	return p.description
}

// Initialize инициализирует плагин
func (p *MetricsPlugin) Initialize() error {
	// Инициализация метрик уже происходит при импорте пакета middleware
	return nil
}

// RegisterAPI регистрирует API endpoints плагина
func (p *MetricsPlugin) RegisterAPI(router *mux.Router) error {
	// Метрики регистрируются через middleware, не требуют отдельных endpoints
	// Prometheus метрики доступны через стандартный /metrics endpoint
	return nil
}

// RegisterWeb регистрирует web компоненты плагина
func (p *MetricsPlugin) RegisterWeb() error {
	// Метрики не требуют frontend компонентов
	return nil
}

// GetMenuItems возвращает пункты меню для плагина
func (p *MetricsPlugin) GetMenuItems() []plugins.MenuItem {
	// Метрики не требуют пунктов меню, они доступны через /metrics
	return []plugins.MenuItem{}
}

// GetMiddleware возвращает middleware для метрик
func (p *MetricsPlugin) GetMiddleware() func(http.Handler) http.Handler {
	return middleware.MetricsMiddleware
}

func init() {
	if err := plugins.Register(NewMetricsPlugin()); err != nil {
		panic(fmt.Sprintf("Failed to register metrics plugin: %v", err))
	}
}
