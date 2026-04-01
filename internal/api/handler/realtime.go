package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RealtimeHandler handles realtime subscription API requests
type RealtimeHandler struct {
	logger    *logrus.Entry
	registrar RealtimeRegistrar
}

// RealtimeSubscribeRequest represents the request body for subscription
type RealtimeSubscribeRequest struct {
	SessionID string   `json:"sessionID,omitempty"`
	Topics    []string `json:"topics" binding:"required"`
}

// RealtimeSubscribeAck represents the acknowledgment response
type RealtimeSubscribeAck struct {
	SessionID string   `json:"sessionID"`
	Topics    []string `json:"topics"`
	Accepted  []string `json:"accepted"`
	Rejected  []string `json:"rejected,omitempty"`
}

// RealtimeRegistrar activates subscriptions for an active realtime session.
type RealtimeRegistrar interface {
	RegisterSubscriptions(ctx context.Context, sessionID string, topics []string) (*RealtimeSubscribeAck, error)
}

type realtimeControlRequest struct {
	SessionID string   `json:"sessionID"`
	Topics    []string `json:"topics"`
}

// NewRealtimeHandler creates a new realtime handler
func NewRealtimeHandler(logger *logrus.Entry, registrar RealtimeRegistrar) *RealtimeHandler {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return &RealtimeHandler{
		logger:    logger.WithField("handler", "realtime"),
		registrar: registrar,
	}
}

// Subscribe handles POST /api/v1/realtime/subscriptions
// @Summary Register realtime topic subscriptions for the active session
// @Description Register realtime topic subscriptions for the active session
// @Tags realtime
// @Accept json
// @Produce json
// @Param request body RealtimeSubscribeRequest true "Subscription request"
// @Success 200 {object} RealtimeSubscribeAck
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/realtime/subscriptions [post]
func (h *RealtimeHandler) Subscribe(c *gin.Context) {
	var req RealtimeSubscribeRequest
	if !bindJSON(c, &req) {
		return
	}

	if len(req.Topics) == 0 {
		respondError(c, http.StatusBadRequest, "INVALID_TOPICS", "At least one topic is required")
		return
	}

	sessionID := strings.TrimSpace(req.SessionID)
	if sessionID == "" {
		sessionID = strings.TrimSpace(c.GetHeader("X-Session-ID"))
	}
	if sessionID == "" {
		respondError(c, http.StatusBadRequest, "MISSING_SESSION", "An active realtime session ID is required")
		return
	}

	if h.registrar == nil {
		respondError(c, http.StatusServiceUnavailable, "REALTIME_UNAVAILABLE", "Realtime subscription registrar is not configured")
		return
	}

	// Log subscription request
	h.logger.WithFields(logrus.Fields{
		"session_id": sessionID,
		"topics":     req.Topics,
	}).Info("Realtime subscription request")

	ack, err := h.registrar.RegisterSubscriptions(c.Request.Context(), sessionID, req.Topics)
	if err != nil {
		var registrarErr *RealtimeRegistrarError
		if errors.As(err, &registrarErr) {
			switch registrarErr.StatusCode {
			case http.StatusBadRequest:
				respondError(c, http.StatusBadRequest, "INVALID_SUBSCRIPTION", registrarErr.Error())
			case http.StatusNotFound:
				respondError(c, http.StatusNotFound, "SESSION_NOT_FOUND", registrarErr.Error())
			default:
				respondError(c, http.StatusBadGateway, "REALTIME_REGISTRATION_FAILED", registrarErr.Error())
			}
			return
		}
		respondError(c, http.StatusBadGateway, "REALTIME_REGISTRATION_FAILED", err.Error())
		return
	}

	c.JSON(http.StatusOK, ack)
}

// HTTPRealtimeRegistrar forwards subscription activation to the iot-ws control endpoint.
type HTTPRealtimeRegistrar struct {
	baseURL string
	client  *http.Client
}

// RealtimeRegistrarError captures HTTP failure details from the realtime registrar.
type RealtimeRegistrarError struct {
	StatusCode int
	Message    string
}

func (e *RealtimeRegistrarError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("realtime registrar returned status %d", e.StatusCode)
}

// NewHTTPRealtimeRegistrar creates a registrar backed by the iot-ws HTTP control endpoint.
func NewHTTPRealtimeRegistrar(baseURL string) *HTTPRealtimeRegistrar {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		return nil
	}

	return &HTTPRealtimeRegistrar{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// RegisterSubscriptions calls the configured iot-ws control endpoint.
func (r *HTTPRealtimeRegistrar) RegisterSubscriptions(ctx context.Context, sessionID string, topics []string) (*RealtimeSubscribeAck, error) {
	payload, err := json.Marshal(realtimeControlRequest{
		SessionID: sessionID,
		Topics:    topics,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var ack RealtimeSubscribeAck
	if err := json.NewDecoder(resp.Body).Decode(&ack); err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &RealtimeRegistrarError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("realtime registrar returned status %d", resp.StatusCode),
		}
	}

	if ack.Topics == nil {
		ack.Topics = ack.Accepted
	}

	return &ack, nil
}
