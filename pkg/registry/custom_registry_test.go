package registry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCustomRegistry(t *testing.T) {
	// Verdaccio with Basic Auth
	client := NewCustomRegistry("http://verdaccio.internal:4873",
		func(o *Options) { o.SetBasicAuth("admin", "secret") })
	assert.Equal(t, "http://verdaccio.internal:4873", client.GetOptions().RegistryURL)
	assert.Equal(t, "admin", client.GetOptions().Username)
	assert.Equal(t, "secret", client.GetOptions().Password)

	// GitHub Packages with Token
	client = NewCustomRegistry("https://npm.pkg.github.com",
		func(o *Options) { o.SetToken("ghp_xxxxx") })
	assert.Equal(t, "https://npm.pkg.github.com", client.GetOptions().RegistryURL)
	assert.Equal(t, "ghp_xxxxx", client.GetOptions().Token)

	// Artifactory with self-signed cert
	client = NewCustomRegistry("https://artifactory.corp.com/artifactory/api/npm/npm-local",
		func(o *Options) {
			o.SetBasicAuth("user", "pass")
			o.SetInsecureSkipVerify(true)
			o.SetTimeout(30 * time.Second)
		})
	assert.True(t, client.GetOptions().InsecureSkipVerify)
	assert.Equal(t, 30*time.Second, client.GetOptions().Timeout)
}

func TestIsPrivateRegistry(t *testing.T) {
	// Official registry is not private
	client := NewRegistry()
	assert.False(t, client.IsPrivateRegistry())

	// China mirrors are not private
	client = NewNpmMirrorRegistry()
	assert.False(t, client.IsPrivateRegistry())

	// Custom URL is private
	client = NewCustomRegistry("http://verdaccio.internal:4873")
	assert.True(t, client.IsPrivateRegistry())

	// GitHub Packages is private
	client = NewCustomRegistry("https://npm.pkg.github.com")
	assert.True(t, client.IsPrivateRegistry())
}

func TestDownloadStatsNotAvailableForPrivateRegistry(t *testing.T) {
	client := NewCustomRegistry("http://verdaccio.internal:4873")
	ctx := context.Background()

	_, err := client.GetDownloadStats(ctx, "my-pkg", "last-week")
	assert.ErrorIs(t, err, ErrDownloadStatsNotAvailable)

	_, err = client.GetDownloadRangeStats(ctx, "my-pkg", "last-week")
	assert.ErrorIs(t, err, ErrDownloadStatsNotAvailable)

	// But if user explicitly sets download stats URL, it should work
	client = NewCustomRegistry("http://verdaccio.internal:4873",
		func(o *Options) { o.SetDownloadStatsURL("http://verdaccio.internal:4873/downloads") })
	assert.NoError(t, client.requireDownloadStatsURL())
}

func TestSetBasicAuth(t *testing.T) {
	options := NewOptions().
		SetRegistryURL("http://verdaccio.internal:4873").
		SetBasicAuth("admin", "secret")

	assert.Equal(t, "admin", options.Username)
	assert.Equal(t, "secret", options.Password)
	assert.True(t, options.HasAuth())

	// No auth
	options2 := NewOptions()
	assert.False(t, options2.HasAuth())

	// Token auth
	options3 := NewOptions().SetToken("npm_xxxxx")
	assert.True(t, options3.HasAuth())
}

func TestSetInsecureSkipVerify(t *testing.T) {
	options := NewOptions()
	assert.False(t, options.InsecureSkipVerify)

	options.SetInsecureSkipVerify(true)
	assert.True(t, options.InsecureSkipVerify)

	// Verify HTTP client is created with TLS config
	client, err := options.GetHttpClient()
	assert.Nil(t, err)
	assert.NotNil(t, client.Transport)
}

func TestRegistryHealthCheckUnreachable(t *testing.T) {
	// Test with a non-routable IP to ensure timeout/failure
	client := NewCustomRegistry("http://192.0.2.1:59999", // TEST-NET-1, guaranteed non-routable
		func(o *Options) { o.SetTimeout(2 * time.Second) })
	ctx := context.Background()
	ok, err := client.RegistryHealthCheck(ctx)
	assert.False(t, ok)
	assert.Error(t, err)
}
