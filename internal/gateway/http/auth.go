// Package http provides HTTP handlers for the gateway service
package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/utmos/utmos/internal/gateway/mqtt"
	"github.com/utmos/utmos/pkg/repository"
)

// VerneMQ Auth Webhook Request
// See: https://docs.vernemq.comConfiguring%20the%20built-in%20webhook%20authenticator
type vernemqAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	// ClientID is optional and may not be present
	ClientID string `json:"client_id,omitempty"`
}

// VerneMQ Auth Webhook Response
type vernemqAuthResponse struct {
	Result   string `json:"result"`   // "ok" or "error"
	Kick     bool   `json:"kick"`     // whether to disconnect existing session
	Username string `json:"username"` // optional modified username
}

// AuthWebhookHandler handles VerneMQ authentication webhook requests
type AuthWebhookHandler struct {
	authenticator *mqtt.Authenticator
	msgLogRepo   *repository.MessageLogRepository
	logger       *logrus.Entry
}

// NewAuthWebhookHandler creates a new auth webhook handler
func NewAuthWebhookHandler(authenticator *mqtt.Authenticator, msgLogRepo *repository.MessageLogRepository, logger *logrus.Entry) *AuthWebhookHandler {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &AuthWebhookHandler{
		authenticator: authenticator,
		msgLogRepo:   msgLogRepo,
		logger:       logger.WithField("component", "auth-webhook"),
	}
}

// Handle processes VerneMQ authentication webhook requests
// POST /api/v1/auth/mqtt
func (h *AuthWebhookHandler) Handle(c *gin.Context) {
	if h.authenticator == nil {
		h.logger.Warn("Authentication request received but authenticator is not available")
		c.JSON(http.StatusServiceUnavailable, vernemqAuthResponse{
			Result: "error",
		})
		return
	}

	var req vernemqAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Debug("Invalid auth request body")
		c.JSON(http.StatusBadRequest, vernemqAuthResponse{
			Result: "error",
		})
		return
	}

	// Authenticate the device
	cred, err := h.authenticator.Authenticate(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		h.logger.WithFields(logrus.Fields{
			"username": req.Username,
			"error":    err.Error(),
		}).Debug("Authentication failed")

		status := http.StatusUnauthorized
		if errors.Is(err, mqtt.ErrDeviceDisabled) {
			status = http.StatusForbidden
		} else if errors.Is(err, mqtt.ErrDeviceNotFound) {
			status = http.StatusNotFound
		}

		// Log failed auth attempt to MessageLog
		h.logAuthAttempt(req.Username, "", false, err.Error())

		c.JSON(status, vernemqAuthResponse{
			Result: "error",
		})
		return
	}

	h.logger.WithField("username", req.Username).Debug("Authentication succeeded")

	// Log successful auth attempt to MessageLog
	h.logAuthAttempt(req.Username, cred.DeviceSN, true, "")

	c.JSON(http.StatusOK, vernemqAuthResponse{
		Result: "ok",
	})
}

// logAuthAttempt logs an authentication attempt to MessageLog
func (h *AuthWebhookHandler) logAuthAttempt(username, deviceSN string, success bool, errorMsg string) {
	if h.msgLogRepo == nil {
		return
	}
	ctx := context.Background()
	err := h.msgLogRepo.CreateAuthAttempt(ctx, username, deviceSN, "iot-gateway", success, errorMsg)
	if err != nil {
		h.logger.WithError(err).Warn("Failed to log auth attempt to MessageLog")
	}
}

// HealthCheck handles a simple health check for the auth webhook endpoint
// GET /api/v1/auth/health
func (h *AuthWebhookHandler) HealthCheck(c *gin.Context) {
	if h.authenticator == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unavailable",
			"reason": "authenticator not initialized",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// RegisterRoutes registers the auth webhook routes with a Gin router
func (h *AuthWebhookHandler) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/mqtt", h.Handle)
		auth.GET("/health", h.HealthCheck)
	}
}
