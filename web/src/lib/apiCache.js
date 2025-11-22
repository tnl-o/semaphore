/**
 * API Response Caching with Invalidation
 */

class ApiCache {
  constructor(ttl = 5 * 60 * 1000) { // 5 minutes default
    this.cache = new Map();
    this.ttl = ttl;
  }

  get(key) {
    const entry = this.cache.get(key);
    if (!entry) return null;

    if (Date.now() - entry.timestamp > this.ttl) {
      this.cache.delete(key);
      return null;
    }

    return entry.data;
  }

  set(key, data) {
    this.cache.set(key, {
      data,
      timestamp: Date.now(),
    });
  }

  invalidate(pattern) {
    if (typeof pattern === 'string') {
      // Exact match
      this.cache.delete(pattern);
    } else if (pattern instanceof RegExp) {
      // Pattern match
      const keys = Array.from(this.cache.keys());
      keys.forEach((key) => {
        if (pattern.test(key)) {
          this.cache.delete(key);
        }
      });
    }
  }

  clear() {
    this.cache.clear();
  }
}

const apiCache = new ApiCache();

export default apiCache;
export { apiCache };
