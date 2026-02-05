package common

// Command defines the interface for all DJI commands.
// All command structs must implement this interface.
type Command interface {
	// Method returns the command method name (e.g., "cover_open").
	Method() string

	// Data returns the command data payload.
	// Returns nil for commands without data.
	Data() any

	// GetHeader returns the message header (TID/BID/Timestamp/Gateway)
	GetHeader() *Header
}
