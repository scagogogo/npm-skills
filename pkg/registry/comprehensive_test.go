package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/scagogogo/npm-skills/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// comprehensiveMockServer creates a full-featured mock NPM registry for testing write operations.
// It simulates: package CRUD, dist-tags, whoami, login, tokens, stars, orgs, hooks, audit, etc.
func comprehensiveMockServer() *httptest.Server {
	packages := map[string]map[string]interface{}{
		"test-pkg": {
			"_id":   "test-pkg",
			"_rev":  "1-abc123",
			"name":  "test-pkg",
			"description": "A test package",
			"dist-tags": map[string]interface{}{
				"latest": "1.0.0",
			},
			"versions": map[string]interface{}{
				"1.0.0": map[string]interface{}{
					"name":    "test-pkg",
					"version": "1.0.0",
					"dist": map[string]interface{}{
						"shasum":  "abc123",
						"tarball": "https://registry.npmjs.org/test-pkg/-/test-pkg-1.0.0.tgz",
					},
				},
			},
			"users": map[string]bool{},
		},
	}

	tokens := []map[string]interface{}{
		{
			"id":     "token-1",
			"key":    "npm_deadbeef",
			"token":  "npm_deadbeef",
			"created": "2024-01-01T00:00:00.000Z",
			"readonly": false,
		},
	}

	hooks := []map[string]interface{}{
		{
			"id":       "hook-1",
			"type":     "hook",
			"endpoint": "https://example.com/webhook",
			"package":  "test-pkg",
			"active":   true,
		},
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		// Root: registry info
		if path == "/" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"db_name": "registry", "doc_count": 1000,
			})
			return
		}

		// WhoAmI
		if path == "/-/whoami" {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{"error": "need auth"})
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"username": "testuser"})
			return
		}

		// Login
		if strings.HasPrefix(path, "/-/user/org.couchdb.user:") && method == http.MethodPut {
			body, _ := io.ReadAll(r.Body)
			var payload map[string]string
			json.Unmarshal(body, &payload)
			if payload["name"] == "testuser" && payload["password"] == "testpass" {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]string{"token": "npm_login_token_123"})
				return
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid credentials"})
			return
		}

		// Create user
		if strings.HasPrefix(path, "/-/user/org.couchdb.user:") && method == http.MethodPut {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"token": "npm_new_user_token"})
			return
		}

		// Get user
		if strings.HasPrefix(path, "/-/user/org.couchdb.user:") && method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"name": "testuser", "email": "test@example.com",
			})
			return
		}

		// Dist-tags endpoints
		if strings.HasPrefix(path, "/-/package/") && strings.Contains(path, "/dist-tags") {
			pkgName := strings.Split(strings.TrimPrefix(path, "/-/package/"), "/")[0]
			pkgName = strings.ReplaceAll(pkgName, "%2F", "/")

			if strings.HasSuffix(path, "/dist-tags") && method == http.MethodGet {
				if pkg, ok := packages[strings.TrimPrefix(pkgName, "@")]; ok {
					if tags, ok := pkg["dist-tags"]; ok {
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(tags)
						return
					}
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{"latest": "1.0.0"})
				return
			}

			if strings.HasSuffix(path, "/dist-tags") && method == http.MethodPost {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{})
				return
			}

			// Single tag: /-/package/{name}/dist-tags/{tag}
			parts := strings.Split(path, "/dist-tags/")
			if len(parts) == 2 {
				if method == http.MethodGet {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode("1.0.0")
					return
				}
				if method == http.MethodPut {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{})
					return
				}
				if method == http.MethodDelete {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{})
					return
				}
			}
		}

		// Token endpoints
		if path == "/-/npm/v1/tokens" && method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"objects": tokens, "total": len(tokens),
			})
			return
		}
		if path == "/-/npm/v1/tokens" && method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "token-new", "token": "npm_new_token", "readonly": true,
			})
			return
		}
		if strings.HasPrefix(path, "/-/npm/v1/tokens/") && method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "token-1", "token": "npm_deadbeef",
			})
			return
		}
		if strings.HasPrefix(path, "/-/npm/v1/tokens/") && method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}

		// Star/unstar — handled via package PUT (see below)

		// Hooks
		if path == "/-/npm/v1/hooks" && method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"objects": hooks, "total": len(hooks),
			})
			return
		}
		if path == "/-/npm/v1/hooks" && method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "hook-new", "endpoint": "https://example.com/new",
			})
			return
		}
		if strings.HasPrefix(path, "/-/npm/v1/hooks/") && method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "hook-1", "endpoint": "https://example.com/webhook",
			})
			return
		}
		if strings.HasPrefix(path, "/-/npm/v1/hooks/") && method == http.MethodPut {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": "hook-1", "endpoint": "https://example.com/updated",
			})
			return
		}
		if strings.HasPrefix(path, "/-/npm/v1/hooks/") && method == http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}

		// Access endpoints
		if strings.HasSuffix(path, "/access") && method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"package": "test-pkg",
				"access":  map[string]string{"read": "public", "write": "restricted"},
			})
			return
		}
		if strings.HasSuffix(path, "/access") && method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}
		if strings.HasSuffix(path, "/collaborators") && method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode([]map[string]interface{}{
				{"name": "testuser", "permissions": "write"},
			})
			return
		}
		if strings.HasSuffix(path, "/collaborators") && method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}
		if strings.Contains(path, "/collaborators/") && method == http.MethodDelete {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Org endpoints
		if strings.HasPrefix(path, "/-/org/") {
			orgPath := strings.TrimPrefix(path, "/-/org/")
			if method == http.MethodGet && !strings.Contains(orgPath, "/") {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{"name": orgPath})
				return
			}
			if method == http.MethodPut && !strings.Contains(orgPath, "/") {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]interface{}{"name": orgPath})
				return
			}
			if method == http.MethodDelete && !strings.Contains(orgPath, "/") {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{})
				return
			}
			if strings.HasSuffix(orgPath, "/member") && method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]string{"user1", "user2"})
				return
			}
			if strings.Contains(orgPath, "/member/") && method == http.MethodPut {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]string{})
				return
			}
			if strings.Contains(orgPath, "/member/") && method == http.MethodDelete {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{})
				return
			}
			if strings.HasSuffix(orgPath, "/package") && method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]string{"pkg1"})
				return
			}
			if strings.HasSuffix(orgPath, "/team") && method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"objects": []map[string]interface{}{{"name": "dev", "id": "org:dev"}},
					"total":   1,
				})
				return
			}
			if strings.Contains(orgPath, "/team/") && method == http.MethodPut {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]interface{}{"name": "newteam"})
				return
			}
			if strings.Contains(orgPath, "/team/") && method == http.MethodDelete {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{})
				return
			}
			if strings.HasSuffix(orgPath, "/member") && method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]string{"user1"})
				return
			}
			if strings.Contains(orgPath, "/member/") && method == http.MethodPut {
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]string{})
				return
			}
			if strings.Contains(orgPath, "/member/") && method == http.MethodDelete {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]string{})
				return
			}
			if strings.HasSuffix(orgPath, "/package") && method == http.MethodGet {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode([]string{"pkg1"})
				return
			}
		}

		// Audit endpoints
		if path == "/-/npm/v1/security/advisories/bulk" && method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}
		if path == "/-/npm/v1/security/audits/quick" && method == http.MethodPost {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"metadata": map[string]interface{}{
					"vulnerabilities": map[string]int{"low": 0, "moderate": 0, "high": 0, "critical": 0},
				},
			})
			return
		}
		if strings.HasPrefix(path, "/-/npm/v1/security/advisories/") && method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": 1234, "title": "Test Advisory", "severity": "moderate",
			})
			return
		}
		if path == "/-/npm/v1/security/advisories" && method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"objects": []map[string]interface{}{}, "total": 0,
			})
			return
		}

		// Search
		if path == "/-/v1/search" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"objects": []map[string]interface{}{
					{"package": map[string]interface{}{"name": "test-pkg", "version": "1.0.0"}},
				},
				"total": 1,
			})
			return
		}

		// StarredByUser / StarredByPackage views
		if strings.HasPrefix(path, "/-/_view/") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"rows": []map[string]interface{}{
					{"value": "test-pkg"},
				},
			})
			return
		}

		// Package PUT (publish, deprecate, star, unpublish-version)
		if method == http.MethodPut && !strings.HasPrefix(path, "/-/") {
			pkgName := strings.TrimPrefix(path, "/")
			pkgName = strings.Split(pkgName, "/")[0]
			body, _ := io.ReadAll(r.Body)
			var pkg map[string]interface{}
			json.Unmarshal(body, &pkg)
			packages[pkgName] = pkg
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "rev": "2-newrev"})
			return
		}

		// Package DELETE (unpublish)
		if method == http.MethodDelete && !strings.HasPrefix(path, "/-/") {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{"ok": true})
			return
		}

		// Package GET
		if method == http.MethodGet && !strings.HasPrefix(path, "/-/") {
			pkgName := strings.TrimPrefix(path, "/")
			// Handle version path: /pkg/1.0.0
			parts := strings.SplitN(pkgName, "/", 2)
			pkgName = parts[0]

			if pkg, ok := packages[pkgName]; ok {
				if len(parts) == 2 {
					// Version request
					version := parts[1]
					if versions, ok := pkg["versions"].(map[string]interface{}); ok {
						if ver, ok := versions[version]; ok {
							w.WriteHeader(http.StatusOK)
							json.NewEncoder(w).Encode(ver)
							return
						}
					}
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(pkg)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
			return
		}

		// Default: not found
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error":"not found: %s"}`, path)
	}))
}

// --- Request helpers tests ---

func TestRequirePackageName(t *testing.T) {
	// Empty
	err := requirePackageName("")
	assert.Error(t, err)

	// Valid simple names
	assert.Nil(t, requirePackageName("react"))
	assert.Nil(t, requirePackageName("lodash"))
	assert.Nil(t, requirePackageName("my-package"))
	assert.Nil(t, requirePackageName("my_package"))
	assert.Nil(t, requirePackageName("my.package"))
	assert.Nil(t, requirePackageName("a1"))

	// Valid scoped names
	assert.Nil(t, requirePackageName("@nestjs/core"))
	assert.Nil(t, requirePackageName("@babel/preset-env"))
	assert.Nil(t, requirePackageName("@angular/cli"))

	// Too long
	longName := strings.Repeat("a", 215)
	err = requirePackageName(longName)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "too long")

	// Max length (214) is OK
	maxName := strings.Repeat("a", 214)
	assert.Nil(t, requirePackageName(maxName))

	// Invalid: starts with dot
	err = requirePackageName(".bad-start")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must not start")

	// Invalid: starts with underscore
	err = requirePackageName("_bad-start")
	assert.Error(t, err)

	// Invalid: uppercase
	err = requirePackageName("React")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")

	// Invalid: spaces
	err = requirePackageName("my package")
	assert.Error(t, err)

	// Invalid scoped: missing name part
	err = requirePackageName("@scopeonly")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid scoped")

	// Invalid scoped: empty scope
	err = requirePackageName("@/name")
	assert.Error(t, err)

	// Invalid scoped: empty name
	err = requirePackageName("@scope/")
	assert.Error(t, err)

	// Invalid scoped: uppercase scope
	err = requirePackageName("@Scope/name")
	assert.Error(t, err)
}

func TestEncodePackageName(t *testing.T) {
	assert.Equal(t, "@nestjs%2Fcore", encodePackageName("@nestjs/core"))
	assert.Equal(t, "react", encodePackageName("react"))
	assert.Equal(t, "@babel%2Fpreset-env", encodePackageName("@babel/preset-env"))
}

func TestRequireAuth(t *testing.T) {
	// No auth
	opts := NewOptions()
	reg := NewRegistry(opts)
	err := reg.requireAuth()
	assert.Error(t, err)

	// With token
	opts = NewOptions().SetToken("npm_test")
	reg = NewRegistry(opts)
	err = reg.requireAuth()
	assert.Nil(t, err)

	// With basic auth
	opts = NewOptions().SetBasicAuth("user", "pass")
	reg = NewRegistry(opts)
	err = reg.requireAuth()
	assert.Nil(t, err)
}

func TestRequireToken(t *testing.T) {
	// No token
	opts := NewOptions()
	reg := NewRegistry(opts)
	err := reg.requireToken()
	assert.Error(t, err)

	// With token
	opts = NewOptions().SetToken("npm_test")
	reg = NewRegistry(opts)
	err = reg.requireToken()
	assert.Nil(t, err)

	// Basic auth doesn't satisfy requireToken
	opts = NewOptions().SetBasicAuth("user", "pass")
	reg = NewRegistry(opts)
	err = reg.requireToken()
	assert.Error(t, err)
}

func TestSendRequestMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test").
		SetTimeout(5 * time.Second))

	// Test DELETE request
	err := reg.DeleteDistTag(context.Background(), "test-pkg", "beta")
	assert.Nil(t, err)
}

func TestSendJSONMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	// Test POST request (set dist-tags)
	err := reg.SetDistTags(context.Background(), "test-pkg", map[string]string{
		"beta": "2.0.0-beta.1",
	})
	assert.Nil(t, err)
}

// --- WhoAmI tests ---

func TestWhoAmIMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	// With token
	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))
	username, err := reg.WhoAmI(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "testuser", username)

	// Without token
	reg = NewRegistry(NewOptions().SetRegistryURL(server.URL))
	_, err = reg.WhoAmI(context.Background())
	assert.Error(t, err)
}

// --- Dist-tags write tests ---

func TestSetDistTagMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.SetDistTag(context.Background(), "test-pkg", "next", "2.0.0-rc.1")
	assert.Nil(t, err)
}

func TestDeleteDistTagMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.DeleteDistTag(context.Background(), "test-pkg", "beta")
	assert.Nil(t, err)
}

func TestSetDistTagsMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.SetDistTags(context.Background(), "test-pkg", map[string]string{
		"next": "2.0.0-rc.1",
		"beta": "1.9.0-beta.3",
	})
	assert.Nil(t, err)
}

// --- Token tests ---

func TestListTokensMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	tokens, err := reg.ListTokens(context.Background())
	assert.Nil(t, err)
	assert.Len(t, tokens, 1)
}

func TestGetTokenMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	token, err := reg.GetToken(context.Background(), "token-1")
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestCreateTokenMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	token, err := reg.CreateToken(context.Background(), &models.TokenCreation{
		Password: "testpass",
		Readonly: true,
	})
	assert.Nil(t, err)
	assert.NotNil(t, token)
}

func TestDeleteTokenMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.DeleteToken(context.Background(), "token-1")
	assert.Nil(t, err)
}

// --- User tests ---

func TestLoginMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	result, err := reg.Login(context.Background(), "testuser", "testpass")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "npm_login_token_123", result.Token)
}

func TestLoginFailureMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	_, err := reg.Login(context.Background(), "baduser", "badpass")
	assert.Error(t, err)
}

func TestGetUserMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	user, err := reg.GetUser(context.Background(), "testuser")
	assert.Nil(t, err)
	assert.NotNil(t, user)
}

// --- Hook tests ---

func TestListHooksMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	hooks, err := reg.ListHooks(context.Background(), models.HookListOptions{})
	assert.Nil(t, err)
	assert.Len(t, hooks, 1)
}

func TestGetHookMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	hook, err := reg.GetHook(context.Background(), "hook-1")
	assert.Nil(t, err)
	assert.NotNil(t, hook)
}

func TestCreateHookMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	hook, err := reg.CreateHook(context.Background(), &models.HookCreation{
		Endpoint: "https://example.com/new",
		Package:  "test-pkg",
	})
	assert.Nil(t, err)
	assert.NotNil(t, hook)
}

func TestUpdateHookMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	hook, err := reg.UpdateHook(context.Background(), "hook-1", &models.HookUpdate{
		Endpoint: "https://example.com/updated",
	})
	assert.Nil(t, err)
	assert.NotNil(t, hook)
}

func TestDeleteHookMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.DeleteHook(context.Background(), "hook-1")
	assert.Nil(t, err)
}

// --- Access tests ---

func TestGetPackageAccessMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	access, err := reg.GetPackageAccess(context.Background(), "test-pkg")
	assert.Nil(t, err)
	assert.NotNil(t, access)
}

func TestSetPackageAccessMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.SetPackageAccess(context.Background(), "test-pkg", &models.PackageAccessUpdate{Access: "restricted"})
	assert.Nil(t, err)
}

func TestListCollaboratorsMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	collabs, err := reg.ListCollaborators(context.Background(), "test-pkg")
	assert.Nil(t, err)
	assert.NotEmpty(t, collabs)
}

func TestGrantAccessMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.GrantAccess(context.Background(), "test-pkg", "newuser", models.PermissionWrite)
	assert.Nil(t, err)
}

func TestRevokeAccessMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.RevokeAccess(context.Background(), "test-pkg", "olduser")
	assert.Nil(t, err)
}

// --- Org tests ---

func TestGetOrgMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	org, err := reg.GetOrg(context.Background(), "testorg")
	assert.Nil(t, err)
	assert.NotNil(t, org)
}

func TestCreateOrgMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	org, err := reg.CreateOrg(context.Background(), "neworg")
	assert.Nil(t, err)
	assert.NotNil(t, org)
}

func TestDeleteOrgMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.DeleteOrg(context.Background(), "oldorg")
	assert.Nil(t, err)
}

func TestListOrgMembersMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	members, err := reg.ListOrgMembers(context.Background(), "testorg")
	assert.Nil(t, err)
	assert.NotEmpty(t, members)
}

func TestAddOrgMemberMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.AddOrgMember(context.Background(), "testorg", "newuser")
	assert.Nil(t, err)
}

func TestRemoveOrgMemberMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.RemoveOrgMember(context.Background(), "testorg", "olduser")
	assert.Nil(t, err)
}

func TestListOrgPackagesMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	pkgs, err := reg.ListOrgPackages(context.Background(), "testorg")
	assert.Nil(t, err)
	assert.NotEmpty(t, pkgs)
}

func TestListTeamsMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	teams, err := reg.ListTeams(context.Background(), "testorg")
	assert.Nil(t, err)
	assert.NotEmpty(t, teams)
}

func TestCreateTeamMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	team, err := reg.CreateTeam(context.Background(), "testorg", "newteam")
	assert.Nil(t, err)
	assert.NotNil(t, team)
}

func TestDeleteTeamMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.DeleteTeam(context.Background(), "testorg", "oldteam")
	assert.Nil(t, err)
}

func TestAddTeamMemberMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.AddTeamMember(context.Background(), "testorg", "dev", "newuser")
	assert.Nil(t, err)
}

func TestRemoveTeamMemberMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	err := reg.RemoveTeamMember(context.Background(), "testorg", "dev", "olduser")
	assert.Nil(t, err)
}

func TestListTeamPackagesMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	pkgs, err := reg.ListTeamPackages(context.Background(), "testorg", "dev")
	assert.Nil(t, err)
	assert.NotEmpty(t, pkgs)
}

// --- Star tests ---

func TestGetStarredByUserMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	pkgs, err := reg.GetStarredByUser(context.Background(), "testuser")
	assert.Nil(t, err)
	assert.NotEmpty(t, pkgs)
}

func TestGetStarredByPackageMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	users, err := reg.GetStarredByPackage(context.Background(), "test-pkg")
	assert.Nil(t, err)
	assert.NotEmpty(t, users)
}

// --- Audit tests ---

func TestGetAdvisoryMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	advisory, err := reg.GetAdvisory(context.Background(), 1234)
	assert.Nil(t, err)
	assert.NotNil(t, advisory)
}

func TestListAdvisoriesMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	advisories, err := reg.ListAdvisories(context.Background(), models.AdvisoryListOptions{PerPage: 20})
	assert.Nil(t, err)
	assert.NotNil(t, advisories)
}

func TestBulkAuditMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	result, err := reg.BulkAudit(context.Background(), map[string][]string{
		"lodash": {"<4.17.12"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestQuickAuditMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	result, err := reg.QuickAudit(context.Background(), &models.QuickAuditRequest{
		Dependencies: map[string]string{"lodash": "4.17.11"},
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

// --- Version tests ---

func TestGetPackageVersionsMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	versions, err := reg.GetPackageVersions(context.Background(), "test-pkg")
	assert.Nil(t, err)
	assert.Contains(t, versions, "1.0.0")
}

func TestGetPackageVersionCountMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	count, err := reg.GetPackageVersionCount(context.Background(), "test-pkg")
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
}

func TestGetPackageLatestVersionMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	latest, err := reg.GetPackageLatestVersion(context.Background(), "test-pkg")
	assert.Nil(t, err)
	assert.Equal(t, "1.0.0", latest)
}

func TestGetPackageInformationSummaryMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))

	pkg, err := reg.GetPackageInformationSummary(context.Background(), "test-pkg")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
}

// --- Auth requirement tests ---

func TestWriteOpsRequireAuth(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))
	ctx := context.Background()

	// All write operations should fail without auth
	assert.Error(t, reg.SetDistTag(ctx, "test-pkg", "next", "2.0.0"))
	assert.Error(t, reg.DeleteDistTag(ctx, "test-pkg", "beta"))
	assert.Error(t, reg.SetDistTags(ctx, "test-pkg", map[string]string{"next": "2.0.0"}))
	_, err := reg.ListTokens(ctx)
	assert.Error(t, err)
	_, err = reg.CreateToken(ctx, &models.TokenCreation{})
	assert.Error(t, err)
	assert.Error(t, reg.DeleteToken(ctx, "token-1"))
	_, err = reg.ListHooks(ctx, models.HookListOptions{})
	assert.Error(t, err)
	_, err = reg.GetHook(ctx, "hook-1")
	assert.Error(t, err)
	_, err = reg.CreateHook(ctx, &models.HookCreation{})
	assert.Error(t, err)
	_, err = reg.UpdateHook(ctx, "hook-1", &models.HookUpdate{})
	assert.Error(t, err)
	assert.Error(t, reg.DeleteHook(ctx, "hook-1"))
	_, err = reg.GetPackageAccess(ctx, "test-pkg")
	assert.Error(t, err)
	assert.Error(t, reg.SetPackageAccess(ctx, "test-pkg", &models.PackageAccessUpdate{}))
	_, err = reg.ListCollaborators(ctx, "test-pkg")
	assert.Error(t, err)
	assert.Error(t, reg.GrantAccess(ctx, "test-pkg", "user", models.PermissionWrite))
	assert.Error(t, reg.RevokeAccess(ctx, "test-pkg", "user"))
	_, err = reg.GetOrg(ctx, "org")
	assert.Error(t, err)
	_, err = reg.CreateOrg(ctx, "org")
	assert.Error(t, err)
	assert.Error(t, reg.DeleteOrg(ctx, "org"))
	_, err = reg.ListOrgMembers(ctx, "org")
	assert.Error(t, err)
	assert.Error(t, reg.AddOrgMember(ctx, "org", "user"))
	assert.Error(t, reg.RemoveOrgMember(ctx, "org", "user"))
	_, err = reg.ListOrgPackages(ctx, "org")
	assert.Error(t, err)
	_, err = reg.ListTeams(ctx, "org")
	assert.Error(t, err)
	_, err = reg.CreateTeam(ctx, "org", "team")
	assert.Error(t, err)
	assert.Error(t, reg.DeleteTeam(ctx, "org", "team"))
	assert.Error(t, reg.AddTeamMember(ctx, "org", "team", "user"))
	assert.Error(t, reg.RemoveTeamMember(ctx, "org", "team", "user"))
	_, err = reg.ListTeamPackages(ctx, "org", "team")
	assert.Error(t, err)
}

// --- Download tarball with auth mock ---

func TestDownloadTarballWithAuthMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test"))

	tmpDir := t.TempDir()
	destPath := tmpDir + "/test-pkg-1.0.0.tgz"
	// This will fail because tarball URL points to real registry, but we test the auth flow
	err := reg.DownloadTarball(context.Background(), "test-pkg", "1.0.0", destPath)
	// Expected to fail since tarball URL in mock points to registry.npmjs.org
	assert.Error(t, err)
}

// --- Publish tests ---

func TestPublishPackageRequiresToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	err := reg.PublishPackage(context.Background(), &models.Package{Name: "test"})
	assert.Error(t, err)
}

func TestPublishPackageFromTarballRequiresToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	err := reg.PublishPackageFromTarball(context.Background(), "test", "1.0.0", []byte{}, &models.PublishMetadata{})
	assert.Error(t, err)
}

// --- Deprecate tests ---

func TestDeprecateVersionRequiresToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	err := reg.DeprecateVersion(context.Background(), "test", "1.0.0", "deprecated")
	assert.Error(t, err)
}

// --- Unpublish tests ---

func TestUnpublishPackageRequiresToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	err := reg.UnpublishPackage(context.Background(), "test")
	assert.Error(t, err)
}

func TestUnpublishPackageVersionRequiresToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	err := reg.UnpublishPackageVersion(context.Background(), "test", "1.0.0")
	assert.Error(t, err)
}

// --- Star/unstar require token ---

func TestStarPackageRequiresToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	err := reg.StarPackage(context.Background(), "test")
	assert.Error(t, err)
}

func TestUnstarPackageRequiresToken(t *testing.T) {
	reg := NewRegistry(NewOptions())
	err := reg.UnstarPackage(context.Background(), "test")
	assert.Error(t, err)
}

// --- Integration-style test: full workflow with mock ---

func TestFullWorkflowMock(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetToken("npm_test").
		SetUserAgent("npm-skills-test/1.0").
		SetTimeout(5 * time.Second))
	ctx := context.Background()

	// Get package info
	pkg, err := reg.GetPackageInformation(ctx, "test-pkg")
	require.Nil(t, err)
	assert.Equal(t, "test-pkg", pkg.Name)

	// Get dist-tags
	tags, err := reg.GetDistTags(ctx, "test-pkg")
	require.Nil(t, err)
	assert.Equal(t, "1.0.0", tags["latest"])

	// Get latest version
	latest, err := reg.GetPackageLatestVersion(ctx, "test-pkg")
	require.Nil(t, err)
	assert.Equal(t, "1.0.0", latest)

	// WhoAmI
	username, err := reg.WhoAmI(ctx)
	require.Nil(t, err)
	assert.Equal(t, "testuser", username)

	// Set a dist-tag
	err = reg.SetDistTag(ctx, "test-pkg", "next", "2.0.0-rc.1")
	require.Nil(t, err)

	// List tokens
	tokens, err := reg.ListTokens(ctx)
	require.Nil(t, err)
	assert.NotEmpty(t, tokens)

	// List hooks
	hooks, err := reg.ListHooks(ctx, models.HookListOptions{})
	require.Nil(t, err)
	assert.NotEmpty(t, hooks)
}

// --- Timeout test ---

func TestRequestTimeout(t *testing.T) {
	// Server that sleeps
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetTimeout(100 * time.Millisecond))

	ctx := context.Background()
	_, err := reg.GetRegistryInformation(ctx)
	assert.Error(t, err)
}

// --- Empty package name validation ---

func TestEmptyPackageNameValidation(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))
	ctx := context.Background()

	_, err := reg.GetPackageInformation(ctx, "")
	assert.Error(t, err)

	_, err = reg.GetPackageVersion(ctx, "", "1.0.0")
	assert.Error(t, err)

	_, err = reg.GetDownloadStats(ctx, "", "last-week")
	assert.Error(t, err)
}

// --- File creation error in DownloadTarball ---

func TestDownloadTarballFileError(t *testing.T) {
	server := comprehensiveMockServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))
	ctx := context.Background()

	// Try to download to a non-existent directory
	err := reg.DownloadTarball(ctx, "test-pkg", "1.0.0", "/nonexistent/dir/file.tgz")
	assert.Error(t, err)
}

// --- Helper: verify os.ReadFile still works after download ---

func TestDownloadTarballAndReadBack(t *testing.T) {
	server := mockTestServer()
	defer server.Close()

	reg := NewRegistry(NewOptions().SetRegistryURL(server.URL))
	ctx := context.Background()

	tmpDir := t.TempDir()
	destPath := tmpDir + "/axios-1.0.0.tgz"

	err := reg.DownloadTarball(ctx, "axios", "1.0.0", destPath)
	require.Nil(t, err)

	data, err := os.ReadFile(destPath)
	require.Nil(t, err)
	assert.NotEmpty(t, data)
}

// --- Basic Auth in requests ---

func TestBasicAuthInRequests(t *testing.T) {
	var receivedAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"db_name":"registry","doc_count":0}`))
	}))
	defer server.Close()

	reg := NewRegistry(NewOptions().
		SetRegistryURL(server.URL).
		SetBasicAuth("admin", "secret"))

	_, err := reg.getBytes(context.Background(), server.URL)
	assert.Nil(t, err)
	assert.Contains(t, receivedAuth, "Basic ")
}

// --- requestSettingBasicAuth ---

func TestRequestSettingBasicAuth(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "http://example.com", nil)
	client := &http.Client{}
	setting := requestSettingBasicAuth("user", "pass")
	err := setting(client, req)
	assert.Nil(t, err)
	assert.Contains(t, req.Header.Get("Authorization"), "Basic ")
}
