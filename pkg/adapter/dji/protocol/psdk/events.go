package psdk

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// PSDK Payload Events
// ===============================

// UIResourceUploadResultData represents the PSDK UI resource upload result data
type UIResourceUploadResultData struct {
	PsdkIndex int    `json:"psdk_index"` // PSDK payload device index (0-3)
	ObjectKey string `json:"object_key"` // OSS object key
	Size      int    `json:"size"`       // File size (bytes)
	Result    int    `json:"result"`     // Error code
}

// UIResourceUploadResultEvent represents the PSDK UI resource upload result event
type UIResourceUploadResultEvent struct {
	common.Header
	MethodName string                         `json:"method"`
	DataValue  UIResourceUploadResultData `json:"data"`
}

// Method returns the method name.
func (e *UIResourceUploadResultEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *UIResourceUploadResultEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *UIResourceUploadResultEvent) GetHeader() *common.Header {
	return &e.Header
}

// FloatingWindowTextData represents the PSDK floating window text data
type FloatingWindowTextData struct {
	PsdkIndex int    `json:"psdk_index"` // PSDK payload device index (0-3)
	Value     string `json:"value"`      // Floating window content
}

// FloatingWindowTextEvent represents the PSDK floating window text event
type FloatingWindowTextEvent struct {
	common.Header
	MethodName string                     `json:"method"`
	DataValue  FloatingWindowTextData `json:"data"`
}

// Method returns the method name.
func (e *FloatingWindowTextEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *FloatingWindowTextEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *FloatingWindowTextEvent) GetHeader() *common.Header { return &e.Header }

// SpeakerPlayProgress represents the speaker play progress information
type SpeakerPlayProgress struct {
	Percent int    `json:"percent"`  // Progress percentage (0-100)
	StepKey string `json:"step_key"` // Current step
}

// SpeakerAudioPlayStartProgressOutput represents the speaker audio play output
type SpeakerAudioPlayStartProgressOutput struct {
	PsdkIndex int                 `json:"psdk_index"` // PSDK payload device index
	Status    string              `json:"status"`     // Current stage (in_progress, ok)
	MD5       string              `json:"md5"`        // File content MD5 checksum
	Progress  SpeakerPlayProgress `json:"progress"`   // Progress information
}

// SpeakerAudioPlayStartProgressData represents the speaker audio play start progress data
type SpeakerAudioPlayStartProgressData struct {
	Result int                                 `json:"result"` // Return code
	Output SpeakerAudioPlayStartProgressOutput `json:"output"` // Output data
}

// SpeakerAudioPlayStartProgressEvent represents the speaker audio play start progress event
type SpeakerAudioPlayStartProgressEvent struct {
	common.Header
	MethodName string                            `json:"method"`
	DataValue  SpeakerAudioPlayStartProgressData `json:"data"`
}

// Method returns the method name.
func (e *SpeakerAudioPlayStartProgressEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *SpeakerAudioPlayStartProgressEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *SpeakerAudioPlayStartProgressEvent) GetHeader() *common.Header {
	return &e.Header
}

// SpeakerTtsPlayStartProgressOutput represents the speaker TTS play output
type SpeakerTtsPlayStartProgressOutput struct {
	PsdkIndex int                 `json:"psdk_index"` // PSDK payload device index
	Status    string              `json:"status"`     // Current stage (in_progress, ok)
	MD5       string              `json:"md5"`        // File content MD5 checksum
	Progress  SpeakerPlayProgress `json:"progress"`   // Progress information
}

// SpeakerTtsPlayStartProgressData represents the speaker TTS play start progress data
type SpeakerTtsPlayStartProgressData struct {
	Result int                               `json:"result"` // Return code
	Output SpeakerTtsPlayStartProgressOutput `json:"output"` // Output data
}

// SpeakerTtsPlayStartProgressEvent represents the speaker TTS play start progress event
type SpeakerTtsPlayStartProgressEvent struct {
	common.Header
	MethodName string                          `json:"method"`
	DataValue  SpeakerTtsPlayStartProgressData `json:"data"`
}

// Method returns the method name.
func (e *SpeakerTtsPlayStartProgressEvent) Method() string { return e.MethodName }

// Data returns the command/event data.
func (e *SpeakerTtsPlayStartProgressEvent) Data() any { return e.DataValue }

// GetHeader returns the event header.
func (e *SpeakerTtsPlayStartProgressEvent) GetHeader() *common.Header {
	return &e.Header
}
