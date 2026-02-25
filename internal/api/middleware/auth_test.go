package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthTestRouter(middleware *AuthMiddleware) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.Handler())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/api/v1/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	return router
}

func TestDefaultAuthConfig(t *testing.T) {
	config := DefaultAuthConfig()

	assert.Equal(t, "X-API-Key", config.HeaderName)
	assert.Contains(t, config.SkipPaths, "/health")
	assert.Contains(t, config.SkipPaths, "/ready")
	assert.Contains(t, config.SkipPaths, "/metrics")
}

func TestNewAuthMiddleware(t *testing.T) {
	config := &AuthConfig{
		APIKeys:    []string{"key1", "key2"},
		HeaderName: "X-API-Key",
	}
	middleware := NewAuthMiddleware(config, nil)

	require.NotNil(t, middleware)
	assert.True(t, middleware.apiKeys["key1"])
	assert.True(t, middleware.apiKeys["key2"])
}

func TestAuthMiddleware_SkipPaths(t *testing.T) {
	config := &AuthConfig{
		APIKeys:   []string{"valid-key"},
		SkipPaths: []string{"/health"},
	}
	middleware := NewAuthMiddleware(config, nil)
	router := setupAuthTestRouter(middleware)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_ValidAPIKey(t *testing.T) {
	config := &AuthConfig{
		APIKeys:    []string{"valid-key"},
		HeaderName: "X-API-Key",
	}
	middleware := NewAuthMiddleware(config, nil)
	router := setupAuthTestRouter(middleware)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/test", nil)
	r.Header.Set("X-API-Key", "valid-key")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_ValidBearerToken(t *testing.T) {
	config := &AuthConfig{
		APIKeys:    []string{"valid-key"},
		HeaderName: "X-API-Key",
	}
	middleware := NewAuthMiddleware(config, nil)
	router := setupAuthTestRouter(middleware)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/test", nil)
	r.Header.Set("Authorization", "Bearer valid-key")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_MissingAPIKey(t *testing.T) {
	config := &AuthConfig{
		APIKeys:    []string{"valid-key"},
		HeaderName: "X-API-Key",
	}
	middleware := NewAuthMiddleware(config, nil)
	router := setupAuthTestRouter(middleware)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/test", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidAPIKey(t *testing.T) {
	config := &AuthConfig{
		APIKeys:    []string{"valid-key"},
		HeaderName: "X-API-Key",
	}
	middleware := NewAuthMiddleware(config, nil)
	router := setupAuthTestRouter(middleware)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/v1/test", nil)
	r.Header.Set("X-API-Key", "invalid-key")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_AddRemoveAPIKey(t *testing.T) {
	config := &AuthConfig{
		APIKeys:    []string{"key1"},
		HeaderName: "X-API-Key",
	}
	middleware := NewAuthMiddleware(config, nil)

	// Add new key
	middleware.AddAPIKey("key2")
	assert.True(t, middleware.apiKeys["key2"])

	// Remove key
	middleware.RemoveAPIKey("key1")
	assert.False(t, middleware.apiKeys["key1"])
}

func TestMaskAPIKey(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"short", "****"},
		{"12345678", "****"},
		{"1234567890", "1234****7890"},
		{"abcdefghijklmnop", "abcd****mnop"},
	}

	for _, tc := range testCases {
		result := maskAPIKey(tc.input)
		assert.Equal(t, tc.expected, result, "input: %s", tc.input)
	}
}

func TestRequireAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequireAuth([]string{"valid-key"}))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	t.Run("valid key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Set("X-API-Key", "valid-key")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Set("X-API-Key", "invalid-key")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestOptionalAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(OptionalAuth())
	router.GET("/test", func(c *gin.Context) {
		apiKey, exists := c.Get("api_key")
		c.JSON(http.StatusOK, gin.H{
			"has_key": exists,
			"api_key": apiKey,
		})
	})

	t.Run("with key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		r.Header.Set("X-API-Key", "some-key")
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "some-key")
	})

	t.Run("without key", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
