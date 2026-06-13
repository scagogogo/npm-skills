package registry

import (
	"fmt"
	"net/http"
)

// Sentinel errors for programmatic error handling.
// These allow SDK consumers to check error types using errors.Is().
var (
	// ErrUnauthorized is returned when the API returns HTTP 401 (authentication required or invalid token).
	ErrUnauthorized = NewAPIError(401, "unauthorized: authentication required or invalid token")

	// ErrForbidden is returned when the API returns HTTP 403 (insufficient permissions).
	ErrForbidden = NewAPIError(403, "forbidden: insufficient permissions")

	// ErrNotFound is returned when the API returns HTTP 404 (resource not found).
	ErrNotFound = NewAPIError(404, "not found: the requested resource does not exist")

	// ErrConflict is returned when the API returns HTTP 409 (conflict, typically a revision mismatch).
	ErrConflict = NewAPIError(409, "conflict: the resource was modified by another request (revision mismatch)")

	// ErrRateLimited is returned when the API returns HTTP 429 (too many requests).
	ErrRateLimited = NewAPIError(429, "rate limited: too many requests, please retry later")

	// ErrInternalServerError is returned when the API returns HTTP 500.
	ErrInternalServerError = NewAPIError(500, "internal server error")

	// ErrBadGateway is returned when the API returns HTTP 502.
	ErrBadGateway = NewAPIError(502, "bad gateway")

	// ErrServiceUnavailable is returned when the API returns HTTP 503.
	ErrServiceUnavailable = NewAPIError(503, "service unavailable")

	// ErrPackageNotFound is returned when a specific package cannot be found.
	ErrPackageNotFound = fmt.Errorf("package not found")

	// ErrVersionNotFound is returned when a specific version cannot be found.
	ErrVersionNotFound = fmt.Errorf("version not found")

	// ErrTokenRequired is returned when an operation requires authentication but no token is set.
	ErrTokenRequired = fmt.Errorf("authentication required: no token set")
)

// APIError represents an error returned by the NPM Registry API with an HTTP status code.
// Use errors.Is() to check against specific error types:
//
//	_, err := registry.GetPackageInformation(ctx, "nonexistent")
//	if errors.Is(err, registry.ErrNotFound) {
//	    // handle 404
//	}
type APIError struct {
	StatusCode int
	Message    string
}

// NewAPIError creates a new APIError with the given status code and message.
func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

// Error implements the error interface.
func (e *APIError) Error() string {
	return fmt.Sprintf("npm registry error (HTTP %d): %s", e.StatusCode, e.Message)
}

// Is implements errors.Is so that errors.Is(err, ErrNotFound) works.
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	return e.StatusCode == t.StatusCode
}

// statusCodeToError maps an HTTP status code to a sentinel APIError.
// Returns nil for non-error status codes (2xx, 3xx).
func statusCodeToError(statusCode int) error {
	switch statusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusConflict:
		return ErrConflict
	case http.StatusTooManyRequests:
		return ErrRateLimited
	case http.StatusInternalServerError:
		return ErrInternalServerError
	case http.StatusBadGateway:
		return ErrBadGateway
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable
	default:
		if statusCode >= 400 {
			return NewAPIError(statusCode, fmt.Sprintf("unexpected HTTP error"))
		}
		return nil
	}
}
