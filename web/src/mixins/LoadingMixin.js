/**
 * LoadingMixin provides loading state management for components
 *
 * Usage:
 * ```javascript
 * export default {
 *   mixins: [LoadingMixin],
 *   methods: {
 *     async loadData() {
 *       await this.withLoading(async () => {
 *         // Your async operation
 *       }, 'Loading data...');
 *     }
 *   }
 * }
 * ```
 */
export default {
  data() {
    return {
      loading: false,
      loadingMessage: null,
    };
  },

  methods: {
    /**
     * Executes an async function with loading state management
     * @param {Function} fn - Async function to execute
     * @param {string} message - Optional loading message
     * @returns {Promise} Result of the async function
     */
    async withLoading(fn, message = null) {
      this.loading = true;
      this.loadingMessage = message;
      try {
        return await fn();
      } finally {
        this.loading = false;
        this.loadingMessage = null;
      }
    },

    /**
     * Sets loading state
     * @param {boolean} value - Loading state
     * @param {string} message - Optional loading message
     */
    setLoading(value, message = null) {
      this.loading = value;
      this.loadingMessage = message;
    },
  },
};
