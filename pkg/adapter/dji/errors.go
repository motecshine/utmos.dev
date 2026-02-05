package dji

import "errors"

// Error definitions for the DJI adapter.
var (
	ErrEmptyTopic       = errors.New("empty topic")
	ErrInvalidTopic     = errors.New("invalid topic format")
	ErrUnknownTopicType = errors.New("unknown topic type")
	ErrEmptyPayload     = errors.New("empty payload")
	ErrInvalidPayload   = errors.New("invalid payload format")
	ErrMissingTID       = errors.New("missing tid field")
	ErrMissingBID       = errors.New("missing bid field")
)
