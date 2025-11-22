// eslint-disable-next-line import/prefer-default-export
export function getErrorMessage(err) {
  // Network errors (no response)
  if (!err.response) {
    if (err.code === 'ECONNABORTED') {
      return 'Request timeout. Please try again.';
    }
    if (err.message === 'Network Error') {
      return 'Network error. Please check your connection.';
    }
    return 'An unexpected error occurred. Please try again.';
  }

  // HTTP errors
  const status = err.response.status;
  const data = err.response.data;

  // Validation errors
  if (status === 400 && data && data.error) {
    return data.error;
  }

  // Authentication errors
  if (status === 401) {
    return 'Your session has expired. Please log in again.';
  }

  // Authorization errors
  if (status === 403) {
    return 'You do not have permission to perform this action.';
  }

  // Not found errors
  if (status === 404) {
    return 'The requested resource was not found.';
  }

  // Server errors
  if (status >= 500) {
    return 'Server error. Please try again later.';
  }

  // Default: use error message from response or generic message
  if (data && data.error) {
    return data.error;
  }

  if (err.message && !err.message.startsWith('Request failed with status code ')) {
    return err.message;
  }

  return 'An error occurred';
}
