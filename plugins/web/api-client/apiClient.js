/**
 * Enhanced API client with interceptors
 * Plugin: api-client
 * Provides:
 * - Automatic error handling
 * - Request retry logic
 * - Rate limiting handling
 * - Request/response logging
 * - Correlation IDs for tracing
 * - API response caching
 */

import axios from 'axios';
import router from '@/router';
import { getErrorMessage } from '@/lib/error';
import { apiCache } from './apiCache';

const CONFIG_FLAG = '__semaphoreAxiosConfigured';
const DEFAULT_TIMEOUT = 30000; // 30 seconds

/**
 * Generate unique request ID for correlation tracking
 * @returns {string} - Unique request ID
 */
function generateRequestId() {
  return `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
}

function ensureHeaders(instance) {
  if (!instance.defaults.headers) {
    instance.defaults.headers = {};
  }
}

function setupAxiosInstance(instance) {
  if (instance[CONFIG_FLAG]) {
    return instance;
  }

  instance.defaults.baseURL = instance.defaults.baseURL || document.baseURI;
  instance.defaults.timeout = instance.defaults.timeout || DEFAULT_TIMEOUT;
  instance.defaults.validateStatus = instance.defaults.validateStatus
    || ((status) => status >= 200 && status < 300);
  ensureHeaders(instance);

  /**
   * Request interceptor
   * - Add correlation ID
   * - Serve cached data
   * - Log requests in development
   */
  instance.interceptors.request.use(
    (config = {}) => {
      const requestConfig = config;
      const method = (requestConfig.method || 'get').toLowerCase();

      // Отключаем кэш для запросов, которые требуют авторизации
      const isAuthRequest = requestConfig.url?.includes('/api/auth/');
      const isOnAuthPage = router.currentRoute?.path?.startsWith('/auth/');

      if (method === 'get' && !requestConfig.skipCache && !isAuthRequest && !isOnAuthPage) {
        const cached = apiCache.get(requestConfig.url);
        if (cached) {
          const isValid = cached !== null
            && cached !== undefined
            && (Array.isArray(cached) || (typeof cached === 'object' && !cached.error && !cached.message));

          if (isValid) {
            const cachedError = new Error('Cached response');
            cachedError.cached = true;
            cachedError.data = cached;
            cachedError.config = requestConfig;
            return Promise.reject(cachedError);
          } else {
            apiCache.invalidate(requestConfig.url);
          }
        }
      }

      requestConfig.headers = requestConfig.headers || {};
      requestConfig.headers['X-Request-ID'] = generateRequestId();

      if (process.env.NODE_ENV === 'development') {
        console.log(`[API Request] ${method.toUpperCase()} ${requestConfig.url}`, {
          data: requestConfig.data,
          params: requestConfig.params,
          requestId: requestConfig.headers['X-Request-ID'],
        });
      }

      return requestConfig;
    },
    (error) => Promise.reject(error),
  );

  instance.interceptors.request.use(
    undefined,
    (error) => {
      if (error.cached && error.data) {
        return Promise.resolve({
          data: error.data,
          status: 200,
          statusText: 'OK (Cached)',
          headers: {},
          config: error.config,
        });
      }
      console.error('[API Request Error]', error);
      return Promise.reject(error);
    },
  );

  /**
   * Response interceptor
   * - Cache GET responses
   * - Handle 401/429/network errors
   */
  instance.interceptors.response.use(
    (response) => {
      const method = (response.config?.method || 'get').toLowerCase();
      if (method === 'get' && !response.config.skipCache && response.status >= 200 && response.status < 300) {
        apiCache.set(response.config.url, response.data);
      }

      if (process.env.NODE_ENV === 'development') {
        console.log(`[API Response] ${method.toUpperCase()} ${response.config.url}`, {
          status: response.status,
          data: response.data,
        });
      }

      return response;
    },
    async (error) => {
      const originalRequest = error.config || {};
      const method = (originalRequest.method || 'get').toLowerCase();

      if (error.response?.status === 401 && !originalRequest.retry) {
        originalRequest.retry = true;
        apiCache.clear();
        apiCache.invalidate(/^\/api\//);

        const currentPath = router.currentRoute?.path || window.location.pathname;
        const isOnAuthPage = currentPath.startsWith('/auth/');
        const isAuthEndpoint = originalRequest.url?.includes('/api/auth/');

        if (!isOnAuthPage && !isAuthEndpoint) {
          const returnPath = router.currentRoute?.fullPath || window.location.pathname;
          if (returnPath !== '/auth/login') {
            const redirectUrl = `/auth/login?redirect=${encodeURIComponent(returnPath)}`;
            setTimeout(() => {
              window.location.href = redirectUrl;
            }, 0);
          }
        }

        return Promise.reject(error);
      }

      if (error.response?.status === 429 && !originalRequest.retry) {
        const retryRequest = { ...originalRequest, retry: true };
        const retryAfter = parseInt(error.response.headers?.['retry-after'] || '1', 10);
        console.warn(`[API] Rate limited, retrying after ${retryAfter}s`);

        await new Promise((resolve) => {
          setTimeout(resolve, retryAfter * 1000);
        });

        return instance(retryRequest);
      }

      if (!error.response && !originalRequest.retry && method === 'get') {
        const retryRequest = { ...originalRequest, retry: true };
        console.warn('[API] Network error, retrying...');

        await new Promise((resolve) => {
          setTimeout(resolve, 1000);
        });

        return instance(retryRequest);
      }

      const errorMessage = getErrorMessage(error);
      console.error('[API Error]', {
        url: originalRequest?.url,
        method: originalRequest?.method,
        status: error.response?.status,
        message: errorMessage,
        requestId: originalRequest?.headers?.['X-Request-ID'],
      });

      return Promise.reject(error);
    },
  );

  instance[CONFIG_FLAG] = true;
  return instance;
}

// Export enhanced client (will be used by plugin system)
export const createEnhancedClient = () => {
  return setupAxiosInstance(axios.create({
    baseURL: document.baseURI,
    timeout: DEFAULT_TIMEOUT,
    headers: {
      'Content-Type': 'application/json',
    },
    validateStatus: (status) => status >= 200 && status < 300,
  }));
};

export { apiCache };

