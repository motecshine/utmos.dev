// Package middleware provides HTTP middleware for iot-api
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	// APIKeys is a list of valid API keys
	APIKeys []string
	// HeaderName is the header name for API key
	HeaderName string
	// SkipPaths are paths that don't require authentication
	SkipPaths []string
}

// DefaultAuthConfig returns default auth configuration
func DefaultAuthConfig() *AuthConfig {
	return &AuthConfig{
		HeaderName: "X-API-Key",
		SkipPaths:  []string{"/health", "/ready", "/metrics"},
	}
}

// AuthMiddleware provides API key authentication
type AuthMiddleware struct {
	config *AuthConfig
	logger *logrus.Entry
	apiKeys map[string]bool
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(config *AuthConfig, logger *logrus.Entry) *AuthMiddleware {
	if config == nil {
		config = DefaultAuthConfig()
	}
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}

	// Build API key lookup map
	apiKeys := make(map[string]bool)
	for _, key := range config.APIKeys {
		apiKeys[key] = true
	}

	return &AuthMiddleware{
		config:  config,
		logger:  logger.WithField("middleware", "auth"),
		apiKeys: apiKeys,
	}
}

// Handler returns the Gin middleware handler
func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path should be skipped
		path := c.Request.URL.Path
		for _, skipPath := range m.config.SkipPaths {
			if strings.HasPrefix(path, skipPath) {
				c.Next()
				return
			}
		}

		// Get API key from header
		apiKey := c.GetHeader(m.config.HeaderName)
		if apiKey == "" {
			// Also check Authorization header with Bearer token
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if apiKey == "" {
			m.logger.WithField("path", path).Warn("Missing API key")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "API key is required",
			})
			return
		}

		// Validate API key
		if !m.apiKeys[apiKey] {
			m.logger.WithFields(logrus.Fields{
				"path":    path,
				"api_key": maskAPIKey(apiKey),
			}).Warn("Invalid API key")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid API key",
			})
			return
		}

		// Set API key in context for later use
		c.Set("api_key", apiKey)
		c.Next()
	}
}

// AddAPIKey adds an API key to the allowed list
func (m *AuthMiddleware) AddAPIKey(key string) {
	m.apiKeys[key] = true
}

// RemoveAPIKey removes an API key from the allowed list
func (m *AuthMiddleware) RemoveAPIKey(key string) {
	delete(m.apiKeys, key)
}

// maskAPIKey masks an API key for logging
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

// RequireAuth is a simple middleware that requires authentication
func RequireAuth(apiKeys []string) gin.HandlerFunc {
	keyMap := make(map[string]bool)
	for _, key := range apiKeys {
		keyMap[key] = true
	}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if apiKey == "" || !keyMap[apiKey] {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "UNAUTHORIZED",
				"message": "Invalid or missing API key",
			})
			return
		}

		c.Set("api_key", apiKey)
		c.Next()
	}
}

// OptionalAuth is a middleware that sets API key if present but doesn't require it
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if apiKey != "" {
			c.Set("api_key", apiKey)
		}

		c.Next()
	}
}
