import EventBus from '@/event-bus';

/**
 * Toast notification helper
 * Provides convenient methods for showing toast notifications
 *
 * Usage:
 * ```javascript
 * import { toast } from '@/lib/toast';
 *
 * toast.success('Operation completed successfully');
 * toast.error('Something went wrong');
 * toast.info('Information message');
 * toast.warning('Warning message');
 * ```
 */

/**
 * Show a toast notification
 * @param {string} message - Message to display
 * @param {string} type - Type of notification: 'success', 'error', 'info', 'warning'
 * @param {Object} options - Additional options
 * @param {number} options.duration - Duration in milliseconds (default: 5000)
 */
function show(message, type = 'info', options = {}) {
  const colorMap = {
    success: 'success',
    error: 'error',
    info: 'info',
    warning: 'warning',
  };

  EventBus.$emit('i-snackbar', {
    color: colorMap[type] || 'info',
    text: message,
    ...options,
  });
}

/**
 * Show success notification
 * @param {string} message - Success message
 * @param {Object} options - Additional options
 */
export function success(message, options = {}) {
  show(message, 'success', options);
}

/**
 * Show error notification
 * @param {string|Error} message - Error message or Error object
 * @param {Object} options - Additional options
 */
export function error(message, options = {}) {
  const errorMessage = message instanceof Error ? message.message : message;
  show(errorMessage, 'error', options);
}

/**
 * Show info notification
 * @param {string} message - Info message
 * @param {Object} options - Additional options
 */
export function info(message, options = {}) {
  show(message, 'info', options);
}

/**
 * Show warning notification
 * @param {string} message - Warning message
 * @param {Object} options - Additional options
 */
export function warning(message, options = {}) {
  show(message, 'warning', options);
}

/**
 * Show error from API response
 * Uses getErrorMessage to extract user-friendly error message
 * @param {Error} err - Error object from axios
 * @param {Object} options - Additional options
 */
export function apiError(err, options = {}) {
  // Dynamic import to avoid circular dependency
  import('@/lib/error').then(({ getErrorMessage }) => {
    error(getErrorMessage(err), options);
  });
}

// Default export with all methods
export default {
  success,
  error,
  info,
  warning,
  apiError,
};

// Named export for convenience
export const toast = {
  success,
  error,
  info,
  warning,
  apiError,
};
