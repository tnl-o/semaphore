import apiClient from '@/lib/apiClient';

let cachedMenuItems = null;

/**
 * Загружает список плагинов с сервера
 * @returns {Promise<Array>}
 */
export async function loadPlugins() {
  try {
    const response = await apiClient({
      method: 'get',
      url: '/api/plugins',
      responseType: 'json',
    });
    return response.data || [];
  } catch (err) {
    console.error('Failed to load plugins:', err);
    return [];
  }
}

/**
 * Загружает пункты меню из всех плагинов
 * @returns {Promise<Array>}
 */
export async function loadPluginMenuItems() {
  if (cachedMenuItems) {
    console.log('[Plugins] Using cached menu items:', cachedMenuItems);
    return cachedMenuItems;
  }

  try {
    console.log('[Plugins] Loading menu items from API...');
    const response = await apiClient({
      method: 'get',
      url: '/api/plugins/menu',
      responseType: 'json',
    });
    console.log('[Plugins] API response:', response);
    cachedMenuItems = response.data || [];
    console.log('[Plugins] Parsed menu items:', cachedMenuItems);
    return cachedMenuItems;
  } catch (err) {
    console.error('[Plugins] Failed to load plugin menu items:', err);
    console.error('[Plugins] Error details:', {
      message: err.message,
      response: err.response,
      status: err.response?.status,
      data: err.response?.data,
    });
    return [];
  }
}

/**
 * Очищает кэш пунктов меню
 */
export function clearPluginMenuCache() {
  cachedMenuItems = null;
}
