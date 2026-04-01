package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRealtimeRegistrar struct {
	ack *RealtimeSubscribeAck
	err error
}

func (f *fakeRealtimeRegistrar) RegisterSubscriptions(_ context.Context, _ string, _ []string) (*RealtimeSubscribeAck, error) {
	return f.ack, f.err
}

func TestRealtimeHandler_Subscribe(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("missing session id", func(t *testing.T) {
		h := NewRealtimeHandler(nil, &fakeRealtimeRegistrar{})

		body := bytes.NewBufferString(`{"topics":["device.device-001.property.report"]}`)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/realtime/subscriptions", body)
		c.Request.Header.Set("Content-Type", "application/json")

		h.Subscribe(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("successful registration", func(t *testing.T) {
		h := NewRealtimeHandler(nil, &fakeRealtimeRegistrar{
			ack: &RealtimeSubscribeAck{
				SessionID: "session-1",
				Topics:    []string{"device.device-001.property.report"},
				Accepted:  []string{"device.device-001.property.report"},
			},
		})

		reqBody, err := json.Marshal(RealtimeSubscribeRequest{
			SessionID: "session-1",
			Topics:    []string{"device.device-001.property.report"},
		})
		require.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/realtime/subscriptions", bytes.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		h.Subscribe(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp RealtimeSubscribeAck
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "session-1", resp.SessionID)
		assert.Equal(t, []string{"device.device-001.property.report"}, resp.Accepted)
	})
}
