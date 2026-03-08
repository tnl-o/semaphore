/**
 * Performance monitoring utilities
 */

/**
 * Measure performance of a function
 * @param {string} name - Name of the performance measurement
 * @param {Function} fn - Function to measure
 * @returns {*} - Result of the function
 */
export default function measurePerformance(name, fn) {
  if (process.env.NODE_ENV === 'production' && 'performance' in window) {
    performance.mark(`${name}-start`);

    const result = fn();

    if (result instanceof Promise) {
      return result.finally(() => {
        performance.mark(`${name}-end`);
        performance.measure(name, `${name}-start`, `${name}-end`);

        const measure = performance.getEntriesByName(name)[0];
        if (measure) {
          console.log(`[Performance] ${name}: ${measure.duration.toFixed(2)}ms`);

          // Send to analytics if available
          if (window.gtag) {
            window.gtag('event', 'timing_complete', {
              name,
              value: Math.round(measure.duration),
            });
          }
        }
      });
    }
    performance.mark(`${name}-end`);
    performance.measure(name, `${name}-start`, `${name}-end`);

    const measure = performance.getEntriesByName(name)[0];
    if (measure) {
      console.log(`[Performance] ${name}: ${measure.duration.toFixed(2)}ms`);
    }

    return result;
  }

  return fn();
}
