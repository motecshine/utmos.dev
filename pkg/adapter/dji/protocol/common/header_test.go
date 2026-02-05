package common

import (
	"testing"
	"time"
)

func TestNewHeader(t *testing.T) {
	before := time.Now().UnixMilli()
	header := NewHeader()
	after := time.Now().UnixMilli()

	// Check TID is generated
	if header.TID == "" {
		t.Error("TID should not be empty")
	}

	// Check BID is generated
	if header.BID == "" {
		t.Error("BID should not be empty")
	}

	// Check TID and BID are different
	if header.TID == header.BID {
		t.Error("TID and BID should be different")
	}

	// Check timestamp is in reasonable range
	if header.Timestamp < before || header.Timestamp > after {
		t.Errorf("Timestamp %d not in range [%d, %d]", header.Timestamp, before, after)
	}

	// Check gateway is empty
	if header.Gateway != "" {
		t.Errorf("Gateway should be empty, got %s", header.Gateway)
	}
}

func TestNewHeaderWithBID(t *testing.T) {
	customBID := "custom-business-id"
	header := NewHeaderWithBID(customBID)

	// Check TID is generated
	if header.TID == "" {
		t.Error("TID should not be empty")
	}

	// Check BID is set to custom value
	if header.BID != customBID {
		t.Errorf("Expected BID %s, got %s", customBID, header.BID)
	}

	// Check timestamp is set
	if header.Timestamp <= 0 {
		t.Error("Timestamp should be positive")
	}

	// Check gateway is empty
	if header.Gateway != "" {
		t.Errorf("Gateway should be empty, got %s", header.Gateway)
	}
}

func TestHeader_Getters(t *testing.T) {
	header := Header{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1234567890,
		Gateway:   "test-gateway",
	}

	if header.GetTID() != "test-tid" {
		t.Errorf("Expected TID test-tid, got %s", header.GetTID())
	}

	if header.GetBID() != "test-bid" {
		t.Errorf("Expected BID test-bid, got %s", header.GetBID())
	}

	if header.GetTimestamp() != 1234567890 {
		t.Errorf("Expected timestamp 1234567890, got %d", header.GetTimestamp())
	}

	if header.GetGateway() != "test-gateway" {
		t.Errorf("Expected gateway test-gateway, got %s", header.GetGateway())
	}
}

func TestHeader_SetGateway(t *testing.T) {
	header := NewHeader()
	if header.Gateway != "" {
		t.Error("Initial gateway should be empty")
	}

	header.SetGateway("gateway-123")
	if header.Gateway != "gateway-123" {
		t.Errorf("Expected gateway gateway-123, got %s", header.Gateway)
	}
}

func TestHeader_SetTimestamp(t *testing.T) {
	header := Header{
		TID:       "test-tid",
		BID:       "test-bid",
		Timestamp: 1000,
		Gateway:   "test-gateway",
	}

	before := time.Now().UnixMilli()
	header.SetTimestamp()
	after := time.Now().UnixMilli()

	if header.Timestamp < before || header.Timestamp > after {
		t.Errorf("Timestamp %d not in range [%d, %d]", header.Timestamp, before, after)
	}
}

func TestHeader_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		header Header
		valid  bool
	}{
		{
			name: "valid header",
			header: Header{
				TID:       "test-tid",
				BID:       "test-bid",
				Timestamp: 1234567890,
				Gateway:   "gateway-123",
			},
			valid: true,
		},
		{
			name: "missing TID",
			header: Header{
				TID:       "",
				BID:       "test-bid",
				Timestamp: 1234567890,
				Gateway:   "gateway-123",
			},
			valid: false,
		},
		{
			name: "missing BID",
			header: Header{
				TID:       "test-tid",
				BID:       "",
				Timestamp: 1234567890,
				Gateway:   "gateway-123",
			},
			valid: false,
		},
		{
			name: "missing timestamp",
			header: Header{
				TID:       "test-tid",
				BID:       "test-bid",
				Timestamp: 0,
				Gateway:   "gateway-123",
			},
			valid: false,
		},
		{
			name: "missing gateway",
			header: Header{
				TID:       "test-tid",
				BID:       "test-bid",
				Timestamp: 1234567890,
				Gateway:   "",
			},
			valid: false,
		},
		{
			name: "all fields missing",
			header: Header{
				TID:       "",
				BID:       "",
				Timestamp: 0,
				Gateway:   "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.header.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestNewHeader_UniqueIDs(t *testing.T) {
	// Generate multiple headers and ensure TIDs and BIDs are unique
	headers := make([]Header, 100)
	tids := make(map[string]bool)
	bids := make(map[string]bool)

	for i := 0; i < 100; i++ {
		headers[i] = NewHeader()

		if tids[headers[i].TID] {
			t.Errorf("Duplicate TID found: %s", headers[i].TID)
		}
		tids[headers[i].TID] = true

		if bids[headers[i].BID] {
			t.Errorf("Duplicate BID found: %s", headers[i].BID)
		}
		bids[headers[i].BID] = true
	}
}
