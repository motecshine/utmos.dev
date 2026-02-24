// Package bridge provides MQTT to/from RabbitMQ bridging functionality
package bridge

import (
	"encoding/json"
	"fmt"
)

// parseRawMessage unmarshals data into a value of type T
func parseRawMessage[T any](data []byte, typeName string) (*T, error) {
	var msg T
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal raw %s message: %w", typeName, err)
	}
	return &msg, nil
}
