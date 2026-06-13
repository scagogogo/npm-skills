package registry

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError(t *testing.T) {
	err := NewAPIError(404, "not found")
	assert.Equal(t, 404, err.StatusCode)
	assert.Contains(t, err.Error(), "404")
	assert.Contains(t, err.Error(), "not found")
}

func TestAPIErrorIs(t *testing.T) {
	err := NewAPIError(404, "package not found")
	assert.True(t, errors.Is(err, ErrNotFound))
	assert.False(t, errors.Is(err, ErrUnauthorized))
	assert.False(t, errors.Is(err, ErrConflict))
}

func TestStatusCodeToError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       error
	}{
		{"401 Unauthorized", 401, ErrUnauthorized},
		{"403 Forbidden", 403, ErrForbidden},
		{"404 Not Found", 404, ErrNotFound},
		{"409 Conflict", 409, ErrConflict},
		{"429 Rate Limited", 429, ErrRateLimited},
		{"500 Internal Server Error", 500, ErrInternalServerError},
		{"502 Bad Gateway", 502, ErrBadGateway},
		{"503 Service Unavailable", 503, ErrServiceUnavailable},
		{"200 OK returns nil", 200, nil},
		{"201 Created returns nil", 201, nil},
		{"301 Redirect returns nil", 301, nil},
		{"418 Other 4xx returns generic error", 418, NewAPIError(418, "unexpected HTTP error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := statusCodeToError(tt.statusCode)
			if tt.want == nil {
				assert.Nil(t, got)
			} else {
				assert.True(t, errors.Is(got, tt.want), "expected errors.Is to match for status %d", tt.statusCode)
			}
		})
	}
}
