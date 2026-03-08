package plugins

import (
	"net/http"
	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/plugins"
)

// GetPlugins возвращает список всех зарегистрированных плагинов
func GetPlugins(w http.ResponseWriter, r *http.Request) {
	allPlugins := plugins.GetAll()
	
	pluginsList := make([]map[string]interface{}, 0, len(allPlugins))
	for _, plugin := range allPlugins {
		pluginsList = append(pluginsList, map[string]interface{}{
			"name":        plugin.Name(),
			"version":     plugin.Version(),
			"description": plugin.Description(),
			"menu_items":  plugin.GetMenuItems(),
		})
	}
	
	helpers.WriteJSON(w, http.StatusOK, pluginsList)
}

// GetPluginMenuItems возвращает все пункты меню из всех плагинов
func GetPluginMenuItems(w http.ResponseWriter, r *http.Request) {
	menuItems := plugins.GetAllMenuItems()
	helpers.WriteJSON(w, http.StatusOK, menuItems)
}

