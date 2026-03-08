/**
 * Enhanced WebSocket client
 * Plugin: websocket
 * Provides:
 * - Automatic reconnection
 * - Error handling
 * - Connection timeout
 * - Authentication error handling
 */

import Listenable from '@/lib/Listenable';

export default class Socket extends Listenable {
  constructor(websocketCreator) {
    super();
    this.websocketCreator = websocketCreator;
  }

  start() {
    if (this.ws != null) {
      throw new Error('Websocket already started. Please stop it before starting.');
    }
    try {
      this.ws = this.websocketCreator();

      // Таймаут для pending соединений (10 секунд)
      const connectionTimeout = setTimeout(() => {
        if (this.ws && this.ws.readyState === WebSocket.CONNECTING) {
          console.warn('[WebSocket] Connection timeout, closing pending connection');
          this.ws.close();
        }
      }, 10000);

      this.ws.onerror = (error) => {
        clearTimeout(connectionTimeout);
        console.warn('[WebSocket] Connection error:', error);
      };

      this.ws.onopen = () => {
        clearTimeout(connectionTimeout);
        console.log('[WebSocket] Connection opened successfully');
      };

      this.ws.onclose = (event) => {
        clearTimeout(connectionTimeout);
        const wasRunning = this.isRunning();
        this.ws = null;

        if (!wasRunning) {
          return;
        }

        // Переподключаемся только если это не было принудительное закрытие
        // и не было ошибки авторизации (код 1008 или 1002)
        // Код 1006 означает abnormal closure (например, сетевые проблемы)
        if (event.code !== 1008 && event.code !== 1002 && event.code !== 1006) {
          setTimeout(() => {
            if (!this.isRunning()) {
              this.start();
            }
          }, 2000);
        } else {
          console.warn('[WebSocket] Connection closed due to authentication/network error, will not reconnect automatically');
        }
      };

      this.ws.onmessage = ({ data }) => {
        try {
          this.callListeners(JSON.parse(data));
        } catch (err) {
          console.error('[WebSocket] Error parsing message:', err);
        }
      };
    } catch (error) {
      console.error('[WebSocket] Failed to create connection:', error);
    }
  }

  isRunning() {
    return this.ws != null;
  }

  stop() {
    if (!this.ws) {
      return;
    }
    this.ws.close();
    delete this.ws;
  }
}

