package tangerino

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError represents an HTTP error returned by the Tangerino API.
type APIError struct {
	// StatusCode is the HTTP status code from the response.
	StatusCode int
	// Body contains the raw response body, available for debugging or custom parsing.
	Body []byte
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("tangerino: API error (status %d)", e.StatusCode)
}

// IsNotFound reports whether err is a 404 Not Found API error.
func IsNotFound(err error) bool {
	return hasStatus(err, http.StatusNotFound)
}

// IsUnauthorized reports whether err is a 401 Unauthorized API error.
func IsUnauthorized(err error) bool {
	return hasStatus(err, http.StatusUnauthorized)
}

// IsForbidden reports whether err is a 403 Forbidden API error.
func IsForbidden(err error) bool {
	return hasStatus(err, http.StatusForbidden)
}

// IsRateLimited reports whether err is a 429 Too Many Requests API error.
func IsRateLimited(err error) bool {
	return hasStatus(err, http.StatusTooManyRequests)
}

// IsServerError reports whether err is a 5xx server-side API error.
func IsServerError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode >= 500
}

func hasStatus(err error, code int) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == code
}
