/**
 * HTML sanitization utilities for safe rendering
 * Prevents XSS attacks in user-generated content (logs, task names, etc.)
 */

import DOMPurify from 'dompurify';

/**
 * Sanitize HTML content for safe rendering in logs
 * Preserves ANSI color classes while removing dangerous content
 * @param {string} html - HTML string to sanitize
 * @returns {string} - Sanitized HTML
 */
export function sanitizeHtml(html) {
  if (!html) return '';

  return DOMPurify.sanitize(html, {
    // Only allow safe tags for log output
    ALLOWED_TAGS: ['span', 'div', 'br', 'pre', 'code'],
    ALLOWED_ATTR: ['class', 'style'],
    ALLOW_DATA_ATTR: false,
    // Preserve ANSI color classes
    ALLOWED_CLASSES: {
      span: /^ansi-(black|red|green|yellow|blue|magenta|cyan|white|bright-.*)$/,
      div: /^ansi-.*$/,
    },
    // Remove script tags and event handlers
    FORBID_TAGS: ['script', 'iframe', 'object', 'embed', 'link', 'style'],
    FORBID_ATTR: ['onerror', 'onload', 'onclick', 'onmouseover'],
  });
}

/**
 * Sanitize plain text to prevent XSS
 * Escapes HTML special characters
 * @param {string} text - Text to sanitize
 * @returns {string} - Escaped text
 */
export function escapeHtml(text) {
  if (!text) return '';

  const map = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;',
  };

  return String(text).replace(/[&<>"']/g, (m) => map[m]);
}

/**
 * Sanitize task names and other user input
 * @param {string} input - User input to sanitize
 * @returns {string} - Sanitized input
 */
export function sanitizeUserInput(input) {
  if (!input) return '';

  // Remove HTML tags
  const withoutTags = input.replace(/<[^>]*>/g, '');

  // Escape remaining HTML
  return escapeHtml(withoutTags);
}
