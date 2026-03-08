/**
 * Secret masking utilities
 * Masks secrets in strings for logging/debugging
 */

/**
 * Mask secrets in strings (for logging/debugging)
 * @param {string} text - Text that may contain secrets
 * @param {Array<RegExp>} patterns - Additional regex patterns to match secrets
 * @returns {string} - Text with secrets masked
 */
export function maskSecrets(text, patterns = []) {
  if (!text) return text;

  const defaultPatterns = [
    /password["\s:=]+([^\s"']+)/gi,
    /token["\s:=]+([^\s"']+)/gi,
    /secret["\s:=]+([^\s"']+)/gi,
    /api[_-]?key["\s:=]+([^\s"']+)/gi,
    /-----BEGIN.*?-----[\s\S]*?-----END.*?-----/gi, // SSH keys
    /aws_access_key_id["\s:=]+([^\s"']+)/gi,
    /aws_secret_access_key["\s:=]+([^\s"']+)/gi,
  ];

  let masked = text;
  [...defaultPatterns, ...patterns].forEach((pattern) => {
    masked = masked.replace(pattern, (match, secret) => {
      if (secret) {
        return match.replace(secret, '***MASKED***');
      }
      return '***MASKED***';
    });
  });

  return masked;
}

/**
 * Check if a string contains potential secrets
 * @param {string} text - Text to check
 * @returns {boolean} - True if text likely contains secrets
 */
export function containsSecrets(text) {
  if (!text) return false;

  const secretIndicators = [
    /password/i,
    /token/i,
    /secret/i,
    /api[_-]?key/i,
    /-----BEGIN/,
    /aws_access_key/i,
    /aws_secret/i,
  ];

  return secretIndicators.some((pattern) => pattern.test(text));
}
