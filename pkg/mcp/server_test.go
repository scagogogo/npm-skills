package mcp

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/scagogogo/npm-skills/pkg/registry"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	cfg := Config{
		RegistryOptions: registry.NewOptions(),
		Timeout:         10 * time.Second,
	}
	server := NewServer(cfg)
	assert.NotNil(t, server)
}

func TestToolError(t *testing.T) {
	result := toolError("test error: %s", "detail")
	assert.NotNil(t, result)
	assert.True(t, result.IsError)
}

func TestToolResult(t *testing.T) {
	data := map[string]string{"name": "react", "version": "18.0.0"}
	result := toolResult(data)
	assert.NotNil(t, result)
	assert.False(t, result.IsError)
	// Should contain JSON
	assert.Contains(t, result.Content[0].(mcp.TextContent).Text, "react")
}

func TestFormatJSON(t *testing.T) {
	// Normal data
	data := map[string]string{"key": "value"}
	s := formatJSON(data)
	assert.Contains(t, s, "key")
	assert.Contains(t, s, "value")

	// Large data should be truncated
	largeData := make(map[string]string)
	for i := 0; i < 20000; i++ {
		largeData[fmt.Sprintf("key_%d", i)] = "value"
	}
	s = formatJSON(largeData)
	assert.Contains(t, s, "TRUNCATED")

	// Unserializable data
	s = formatJSON(func() {})
	assert.Contains(t, s, "error")
}

func TestWithTimeout(t *testing.T) {
	cfg := Config{Timeout: 5 * time.Second}
	parent := context.Background()

	ctx, cancel := withTimeout(parent, cfg)
	defer cancel()

	deadline, ok := ctx.Deadline()
	assert.True(t, ok, "should have a deadline")
	assert.True(t, !deadline.IsZero(), "deadline should not be zero")
}

func TestFormatJSONTruncation(t *testing.T) {
	// Test exact boundary: data just under 100KB should not be truncated
	smallData := map[string]string{"a": strings.Repeat("x", 1000)}
	s := formatJSON(smallData)
	assert.NotContains(t, s, "TRUNCATED")

	// Data over 100KB should be truncated
	bigData := make(map[string]string)
	for i := 0; i < 10000; i++ {
		bigData[fmt.Sprintf("k%d", i)] = strings.Repeat("x", 20)
	}
	s = formatJSON(bigData)
	assert.Contains(t, s, "TRUNCATED")
}

func TestGetOptionalFloat(t *testing.T) {
	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: map[string]interface{}{
				"float_val":  1.5,
				"string_val": "2.5",
				"int_val":    float64(3),
			},
		},
	}

	// float64 value
	v, ok := getOptionalFloat(req, "float_val")
	assert.True(t, ok)
	assert.Equal(t, 1.5, v)

	// string value
	v, ok = getOptionalFloat(req, "string_val")
	assert.True(t, ok)
	assert.Equal(t, 2.5, v)

	// int as float64
	v, ok = getOptionalFloat(req, "int_val")
	assert.True(t, ok)
	assert.Equal(t, float64(3), v)

	// missing key
	v, ok = getOptionalFloat(req, "nonexistent")
	assert.False(t, ok)
	assert.Equal(t, float64(0), v)
}
