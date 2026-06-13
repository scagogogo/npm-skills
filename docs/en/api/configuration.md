# Configuration Options

This document describes all configuration options available for customizing the NPM Crawler client behavior.

## Options Structure

### Options
The main configuration structure for customizing client behavior.

```go
type Options struct {
    registryURL string
    httpClient  *http.Client
    proxy       string
    userAgent   string
    timeout     time.Duration
}
```

## Creating Options

### NewOptions
```go
func NewOptions() *Options
```

Creates a new Options instance with default values.

**Returns:**
- `*Options` - New options instance with defaults

**Example:**
```go
options := registry.NewOptions()
```

## Configuration Methods

All configuration methods return the Options instance for method chaining.

### SetRegistryURL
```go
func (o *Options) SetRegistryURL(url string) *Options
```

Sets the NPM registry URL to use.

**Parameters:**
- `url` - Registry URL (e.g., "https://registry.npmjs.org")

**Example:**
```go
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org")
```

### SetHTTPClient
```go
func (o *Options) SetHTTPClient(client *http.Client) *Options
```

Sets a custom HTTP client for all requests.

**Parameters:**
- `client` - Custom HTTP client instance

**Example:**
```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}

options := registry.NewOptions().
    SetHTTPClient(httpClient)
```

### SetProxy
```go
func (o *Options) SetProxy(proxyURL string) *Options
```

Sets a proxy server for all requests.

**Parameters:**
- `proxyURL` - Proxy server URL (e.g., "http://proxy.example.com:8080")

**Example:**
```go
options := registry.NewOptions().
    SetProxy("http://proxy.example.com:8080")
```

### SetUserAgent
```go
func (o *Options) SetUserAgent(userAgent string) *Options
```

Sets a custom User-Agent header for requests.

**Parameters:**
- `userAgent` - User-Agent string

**Example:**
```go
options := registry.NewOptions().
    SetUserAgent("MyApp/1.0 npm-skills")
```

### SetTimeout
```go
func (o *Options) SetTimeout(timeout time.Duration) *Options
```

Sets the default timeout for all requests.

**Parameters:**
- `timeout` - Request timeout duration

**Example:**
```go
options := registry.NewOptions().
    SetTimeout(30 * time.Second)
```

## Configuration Examples

### Basic Configuration
```go
package main

import (
    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    options := registry.NewOptions().
        SetRegistryURL("https://registry.npmjs.org").
        SetTimeout(30 * time.Second)
    
    client := registry.NewRegistry(options)
}
```

### Advanced HTTP Client Configuration
```go
package main

import (
    "crypto/tls"
    "net/http"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Create custom HTTP client
    httpClient := &http.Client{
        Timeout: 60 * time.Second,
        Transport: &http.Transport{
            MaxIdleConns:        100,
            MaxIdleConnsPerHost: 10,
            IdleConnTimeout:     90 * time.Second,
            TLSHandshakeTimeout: 10 * time.Second,
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: false,
                MinVersion:         tls.VersionTLS12,
            },
        },
    }
    
    options := registry.NewOptions().
        SetHTTPClient(httpClient).
        SetUserAgent("MyApp/1.0")
    
    client := registry.NewRegistry(options)
}
```

### Proxy Configuration
```go
package main

import (
    "net/http"
    "net/url"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Method 1: Using SetProxy
    options1 := registry.NewOptions().
        SetProxy("http://proxy.example.com:8080")
    
    client1 := registry.NewRegistry(options1)
    
    // Method 2: Custom HTTP client with proxy
    proxyURL, _ := url.Parse("http://proxy.example.com:8080")
    httpClient := &http.Client{
        Transport: &http.Transport{
            Proxy: http.ProxyURL(proxyURL),
        },
    }
    
    options2 := registry.NewOptions().
        SetHTTPClient(httpClient)
    
    client2 := registry.NewRegistry(options2)
}
```

### Authentication Configuration
```go
package main

import (
    "net/http"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// Custom transport for authentication
type authTransport struct {
    token string
    base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    req.Header.Set("Authorization", "Bearer "+t.token)
    return t.base.RoundTrip(req)
}

func main() {
    // Create HTTP client with authentication
    httpClient := &http.Client{
        Transport: &authTransport{
            token: "your-auth-token",
            base:  http.DefaultTransport,
        },
    }
    
    options := registry.NewOptions().
        SetRegistryURL("https://private-registry.com").
        SetHTTPClient(httpClient)
    
    client := registry.NewRegistry(options)
}
```

### Method Chaining
```go
package main

import (
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    // Chain multiple configuration methods
    client := registry.NewRegistry(
        registry.NewOptions().
            SetRegistryURL("https://registry.npmjs.org").
            SetTimeout(30*time.Second).
            SetUserAgent("MyApp/1.0").
            SetProxy("http://proxy.example.com:8080"),
    )
}
```

## Default Values

When no options are provided, the following defaults are used:

| Option | Default Value |
|--------|---------------|
| Registry URL | `https://registry.npmjs.org` |
| HTTP Client | `http.DefaultClient` |
| User-Agent | `npm-skills/1.0` |
| Timeout | `30 seconds` |
| Proxy | None |

## Predefined Configurations

The library provides predefined configurations for popular registries:

### Official NPM Registry
```go
// Equivalent to:
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmjs.org")

client := registry.NewRegistry(options)
// Or simply:
client := registry.NewRegistry()
```

### Taobao Mirror
```go
// Equivalent to:
options := registry.NewOptions().
    SetRegistryURL("https://registry.npmmirror.com")

client := registry.NewRegistry(options)
// Or simply:
client := registry.NewTaoBaoRegistry()
```

### NPM Mirror
```go
// Equivalent to:
options := registry.NewOptions().
    SetRegistryURL("https://skimdb.npmjs.com/registry")

client := registry.NewRegistry(options)
// Or simply:
client := registry.NewNpmMirrorRegistry()
```

### Huawei Cloud Mirror
```go
// Equivalent to:
options := registry.NewOptions().
    SetRegistryURL("https://mirrors.huaweicloud.com/repository/npm")

client := registry.NewRegistry(options)
// Or simply:
client := registry.NewHuaWeiCloudRegistry()
```

## Advanced Configurations

### Connection Pooling
```go
package main

import (
    "net/http"
    "time"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    httpClient := &http.Client{
        Transport: &http.Transport{
            MaxIdleConns:        100,  // Maximum idle connections
            MaxIdleConnsPerHost: 10,   // Maximum idle connections per host
            IdleConnTimeout:     90 * time.Second,
            DisableCompression:  false,
        },
        Timeout: 30 * time.Second,
    }
    
    options := registry.NewOptions().
        SetHTTPClient(httpClient)
    
    client := registry.NewRegistry(options)
}
```

### TLS Configuration
```go
package main

import (
    "crypto/tls"
    "net/http"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    httpClient := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: false,
                MinVersion:         tls.VersionTLS12,
                MaxVersion:         tls.VersionTLS13,
            },
        },
    }
    
    options := registry.NewOptions().
        SetHTTPClient(httpClient)
    
    client := registry.NewRegistry(options)
}
```

### Custom Headers
```go
package main

import (
    "net/http"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

// Custom transport for adding headers
type headerTransport struct {
    base    http.RoundTripper
    headers map[string]string
}

func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    for key, value := range t.headers {
        req.Header.Set(key, value)
    }
    return t.base.RoundTrip(req)
}

func main() {
    httpClient := &http.Client{
        Transport: &headerTransport{
            base: http.DefaultTransport,
            headers: map[string]string{
                "X-Custom-Header": "MyValue",
                "Accept-Language": "en-US,en;q=0.9",
            },
        },
    }
    
    options := registry.NewOptions().
        SetHTTPClient(httpClient)
    
    client := registry.NewRegistry(options)
}
```

## Environment Variable Support

You can also configure the client using environment variables:

```bash
export NPM_REGISTRY_URL="https://registry.npmjs.org"
export HTTP_PROXY="http://proxy.example.com:8080"
export HTTPS_PROXY="https://proxy.example.com:8080"
```

```go
package main

import (
    "os"

    "github.com/scagogogo/npm-skills/pkg/registry"
)

func main() {
    options := registry.NewOptions()
    
    // Override with environment variables if set
    if registryURL := os.Getenv("NPM_REGISTRY_URL"); registryURL != "" {
        options.SetRegistryURL(registryURL)
    }
    
    if proxy := os.Getenv("HTTP_PROXY"); proxy != "" {
        options.SetProxy(proxy)
    }
    
    client := registry.NewRegistry(options)
}
```

## Best Practices

### 1. Use Appropriate Timeouts
```go
// For quick operations
options := registry.NewOptions().
    SetTimeout(10 * time.Second)

// For search operations or large packages
options := registry.NewOptions().
    SetTimeout(60 * time.Second)
```

### 2. Configure Connection Pooling
```go
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        50,
        MaxIdleConnsPerHost: 5,
        IdleConnTimeout:     30 * time.Second,
    },
}
```

### 3. Set Meaningful User-Agent
```go
options := registry.NewOptions().
    SetUserAgent("MyApp/1.0.0 (contact@example.com)")
```

### 4. Handle Proxy Configuration Securely
```go
// Don't hardcode proxy credentials
proxyURL := os.Getenv("HTTP_PROXY") // Get from environment
if proxyURL != "" {
    options.SetProxy(proxyURL)
}
```

## Next Steps

- Review [Registry API](registry.md) for method documentation  
- Check [Data Models](models.md) for response structures
- Explore [Examples](../examples/) for practical usage patterns 