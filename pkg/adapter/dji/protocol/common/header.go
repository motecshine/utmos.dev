package common

import (
	"time"

	"github.com/google/uuid"
)

// Header represents the standard message header for all DJI IoT SDK packets
type Header struct {
	TID       string `json:"tid"`       // Transaction ID
	BID       string `json:"bid"`       // Business ID
	Timestamp int64  `json:"timestamp"` // Message send timestamp in milliseconds
	Gateway   string `json:"gateway"`   // Gateway serial number
}

// NewHeader creates a new header with auto-generated TID/BID and current timestamp
// Gateway is empty and should be set by msgbox
func NewHeader() Header {
	return Header{
		TID:       uuid.New().String(),
		BID:       uuid.New().String(),
		Timestamp: time.Now().UnixMilli(),
	}
}

// NewHeaderWithBID creates a new header with specified BID (for business flow)
func NewHeaderWithBID(bid string) Header {
	return Header{
		TID:       uuid.New().String(),
		BID:       bid,
		Timestamp: time.Now().UnixMilli(),
	}
}

// GetTID returns the transaction ID
func (h *Header) GetTID() string {
	return h.TID
}

// GetBID returns the business ID
func (h *Header) GetBID() string {
	return h.BID
}

// GetTimestamp returns the timestamp
func (h *Header) GetTimestamp() int64 {
	return h.Timestamp
}

// GetGateway returns the gateway serial number
func (h *Header) GetGateway() string {
	return h.Gateway
}

// SetGateway sets the gateway serial number
func (h *Header) SetGateway(gateway string) {
	h.Gateway = gateway
}

// SetTimestamp sets the timestamp to current time
func (h *Header) SetTimestamp() {
	h.Timestamp = time.Now().UnixMilli()
}

// IsValid checks if the header has required fields
func (h *Header) IsValid() bool {
	return h.TID != "" && h.BID != "" && h.Timestamp > 0 && h.Gateway != ""
}
