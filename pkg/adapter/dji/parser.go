package dji

import (
	"encoding/json"
	"fmt"
)

// ParseMessage parses raw bytes into a Message.
func ParseMessage(payload []byte) (*Message, error) {
	if len(payload) == 0 {
		return nil, ErrEmptyPayload
	}

	var msg Message
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidPayload, err)
	}

	return &msg, nil
}

// ParseAndValidate parses raw bytes and validates the message.
func ParseAndValidate(payload []byte) (*Message, error) {
	msg, err := ParseMessage(payload)
	if err != nil {
		return nil, err
	}

	if err := msg.Validate(); err != nil {
		return nil, err
	}

	return msg, nil
}
