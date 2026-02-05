package psdk

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// PSDK Payload Events
// ===============================

// PsdkUIResourceUploadResultData represents the PSDK UI resource upload result data
type PsdkUIResourceUploadResultData struct {
	PsdkIndex int    `json:"psdk_index"` // PSDK payload device index (0-3)
	ObjectKey string `json:"object_key"` // OSS object key
	Size      int    `json:"size"`       // File size (bytes)
	Result    int    `json:"result"`     // Error code
}

// PsdkUIResourceUploadResultEvent represents the PSDK UI resource upload result event
type PsdkUIResourceUploadResultEvent struct {
	common.Header
	MethodName string                         `json:"method"`
	DataValue  PsdkUIResourceUploadResultData `json:"data"`
}

func (e *PsdkUIResourceUploadResultEvent) Method() string { return e.MethodName }
func (e *PsdkUIResourceUploadResultEvent) Data() any      { return e.DataValue }
func (e *PsdkUIResourceUploadResultEvent) GetHeader() *common.Header {
	return &e.Header
}

// PsdkFloatingWindowTextData represents the PSDK floating window text data
type PsdkFloatingWindowTextData struct {
	PsdkIndex int    `json:"psdk_index"` // PSDK payload device index (0-3)
	Value     string `json:"value"`      // Floating window content
}

// PsdkFloatingWindowTextEvent represents the PSDK floating window text event
type PsdkFloatingWindowTextEvent struct {
	common.Header
	MethodName string                     `json:"method"`
	DataValue  PsdkFloatingWindowTextData `json:"data"`
}

func (e *PsdkFloatingWindowTextEvent) Method() string            { return e.MethodName }
func (e *PsdkFloatingWindowTextEvent) Data() any                 { return e.DataValue }
func (e *PsdkFloatingWindowTextEvent) GetHeader() *common.Header { return &e.Header }

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

func (e *SpeakerAudioPlayStartProgressEvent) Method() string { return e.MethodName }
func (e *SpeakerAudioPlayStartProgressEvent) Data() any      { return e.DataValue }
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

func (e *SpeakerTtsPlayStartProgressEvent) Method() string { return e.MethodName }
func (e *SpeakerTtsPlayStartProgressEvent) Data() any      { return e.DataValue }
func (e *SpeakerTtsPlayStartProgressEvent) GetHeader() *common.Header {
	return &e.Header
}
