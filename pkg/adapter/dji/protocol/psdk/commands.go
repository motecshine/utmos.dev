package psdk

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// PSDK Payload Commands
// ===============================

// WidgetValueSetData represents the PSDK widget value set data
type WidgetValueSetData struct {
	PsdkIndex int `json:"psdk_index"` // PSDK payload device index (0-3)
	Index     int `json:"index"`      // Widget index
	Value     int `json:"value"`      // Widget value (defined by developer)
}

// WidgetValueSetCommand represents the PSDK widget value set request
type WidgetValueSetCommand struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  WidgetValueSetData `json:"data"`
}

// NewWidgetValueSetCommand creates a new PSDK widget value set request
func NewWidgetValueSetCommand(data WidgetValueSetData) *WidgetValueSetCommand {
	return &WidgetValueSetCommand{
		Header:     common.NewHeader(),
		MethodName: "psdk_widget_value_set",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *WidgetValueSetCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *WidgetValueSetCommand) Data() any { return c.DataValue }

// InputBoxTextSetData represents the PSDK input box text set data
type InputBoxTextSetData struct {
	PsdkIndex int    `json:"psdk_index"` // PSDK payload device index (0-3)
	Value     string `json:"value"`      // Text content (max 128 bytes)
}

// InputBoxTextSetCommand represents the PSDK input box text set request
type InputBoxTextSetCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  InputBoxTextSetData `json:"data"`
}

// NewInputBoxTextSetCommand creates a new PSDK input box text set request
func NewInputBoxTextSetCommand(data InputBoxTextSetData) *InputBoxTextSetCommand {
	return &InputBoxTextSetCommand{
		Header:     common.NewHeader(),
		MethodName: "psdk_input_box_text_set",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *InputBoxTextSetCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *InputBoxTextSetCommand) Data() any { return c.DataValue }

// AudioFile represents audio file information for speaker
type AudioFile struct {
	Name   string `json:"name"`   // File name
	URL    string `json:"url"`    // File download URL
	MD5    string `json:"md5"`    // Audio file MD5 checksum (unique identifier for dock)
	Format string `json:"format"` // Speaker input file format (pcm)
}

// SpeakerAudioPlayStartData represents the speaker audio play start data
type SpeakerAudioPlayStartData struct {
	PsdkIndex int       `json:"psdk_index"` // PSDK payload device index (min: 0)
	File      AudioFile `json:"file"`       // Audio file information
}

// SpeakerAudioPlayStartCommand represents the speaker audio play start request
type SpeakerAudioPlayStartCommand struct {
	common.Header
	MethodName string                    `json:"method"`
	DataValue  SpeakerAudioPlayStartData `json:"data"`
}

// NewSpeakerAudioPlayStartCommand creates a new speaker audio play start request
func NewSpeakerAudioPlayStartCommand(data SpeakerAudioPlayStartData) *SpeakerAudioPlayStartCommand {
	return &SpeakerAudioPlayStartCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_audio_play_start",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *SpeakerAudioPlayStartCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *SpeakerAudioPlayStartCommand) Data() any { return c.DataValue }

// TTSContent represents TTS text content for speaker
type TTSContent struct {
	Name string `json:"name"` // File name
	Text string `json:"text"` // Text content to convert to speech
	MD5  string `json:"md5"`  // File content MD5 checksum (unique identifier for dock)
}

// SpeakerTtsPlayStartData represents the speaker TTS play start data
type SpeakerTtsPlayStartData struct {
	PsdkIndex int        `json:"psdk_index"` // PSDK payload device index (min: 0)
	TTS       TTSContent `json:"tts"`        // TTS text information
}

// SpeakerTtsPlayStartCommand represents the speaker TTS play start request
type SpeakerTtsPlayStartCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  SpeakerTtsPlayStartData `json:"data"`
}

// NewSpeakerTtsPlayStartCommand creates a new speaker TTS play start request
func NewSpeakerTtsPlayStartCommand(data SpeakerTtsPlayStartData) *SpeakerTtsPlayStartCommand {
	return &SpeakerTtsPlayStartCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_tts_play_start",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *SpeakerTtsPlayStartCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *SpeakerTtsPlayStartCommand) Data() any { return c.DataValue }

// SpeakerReplayData represents the speaker replay data
type SpeakerReplayData struct {
	PsdkIndex int `json:"psdk_index"` // PSDK payload device index (min: 0)
}

// SpeakerReplayCommand represents the speaker replay request
type SpeakerReplayCommand struct {
	common.Header
	MethodName string            `json:"method"`
	DataValue  SpeakerReplayData `json:"data"`
}

// NewSpeakerReplayCommand creates a new speaker replay request
func NewSpeakerReplayCommand(data SpeakerReplayData) *SpeakerReplayCommand {
	return &SpeakerReplayCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_replay",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *SpeakerReplayCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *SpeakerReplayCommand) Data() any { return c.DataValue }

// SpeakerPlayStopData represents the speaker play stop data
type SpeakerPlayStopData struct {
	PsdkIndex int `json:"psdk_index"` // PSDK payload device index (min: 0)
}

// SpeakerPlayStopCommand represents the speaker play stop request
type SpeakerPlayStopCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  SpeakerPlayStopData `json:"data"`
}

// NewSpeakerPlayStopCommand creates a new speaker play stop request
func NewSpeakerPlayStopCommand(data SpeakerPlayStopData) *SpeakerPlayStopCommand {
	return &SpeakerPlayStopCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_play_stop",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *SpeakerPlayStopCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *SpeakerPlayStopCommand) Data() any { return c.DataValue }

// SpeakerPlayModeSetData represents the speaker play mode set data
type SpeakerPlayModeSetData struct {
	PsdkIndex int `json:"psdk_index"` // PSDK payload device index (0-3)
	PlayMode  int `json:"play_mode"`  // Play mode: 0=single play, 1=loop (single track)
}

// SpeakerPlayModeSetCommand represents the speaker play mode set request
type SpeakerPlayModeSetCommand struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  SpeakerPlayModeSetData `json:"data"`
}

// NewSpeakerPlayModeSetCommand creates a new speaker play mode set request
func NewSpeakerPlayModeSetCommand(data SpeakerPlayModeSetData) *SpeakerPlayModeSetCommand {
	return &SpeakerPlayModeSetCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_play_mode_set",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *SpeakerPlayModeSetCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *SpeakerPlayModeSetCommand) Data() any { return c.DataValue }

// SpeakerPlayVolumeSetData represents the speaker play volume set data
type SpeakerPlayVolumeSetData struct {
	PsdkIndex  int `json:"psdk_index"` // PSDK payload device index (0-3)
	Index      int `json:"index"`
	PlayVolume int `json:"play_volume"` // Speaker volume (0-100)
}

// SpeakerPlayVolumeSetCommand represents the speaker play volume set request
type SpeakerPlayVolumeSetCommand struct {
	common.Header
	MethodName string                   `json:"method"`
	DataValue  SpeakerPlayVolumeSetData `json:"data"`
}

// NewSpeakerPlayVolumeSetCommand creates a new speaker play volume set request
func NewSpeakerPlayVolumeSetCommand(data SpeakerPlayVolumeSetData) *SpeakerPlayVolumeSetCommand {
	return &SpeakerPlayVolumeSetCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_play_volume_set",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *SpeakerPlayVolumeSetCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *SpeakerPlayVolumeSetCommand) Data() any { return c.DataValue }

// CustomDataTransmissionToPSDKData represents the custom data transmission to PSDK data
type CustomDataTransmissionToPSDKData struct {
	PayloadIndex string      `json:"payload_index"` // Payload index
	DataValue    any `json:"data"`          // Custom data
}

// CustomDataTransmissionToPSDKCommand represents the custom data transmission to PSDK request
type CustomDataTransmissionToPSDKCommand struct {
	common.Header
	MethodName string                           `json:"method"`
	DataValue  CustomDataTransmissionToPSDKData `json:"data"`
}

// NewCustomDataTransmissionToPSDKCommand creates a new custom data transmission to PSDK request
func NewCustomDataTransmissionToPSDKCommand(data CustomDataTransmissionToPSDKData) *CustomDataTransmissionToPSDKCommand {
	return &CustomDataTransmissionToPSDKCommand{
		Header:     common.NewHeader(),
		MethodName: "custom_data_transmission_to_psdk",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *CustomDataTransmissionToPSDKCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *CustomDataTransmissionToPSDKCommand) Data() any { return c.DataValue }

// CustomDataTransmissionToESDKData represents the custom data transmission to ESDK data
type CustomDataTransmissionToESDKData struct {
	Data any `json:"data"` // Custom data
}

// CustomDataTransmissionToESDKCommand represents the custom data transmission to ESDK request
type CustomDataTransmissionToESDKCommand struct {
	common.Header
	MethodName string                           `json:"method"`
	DataValue  CustomDataTransmissionToESDKData `json:"data"`
}

// NewCustomDataTransmissionToESDKCommand creates a new custom data transmission to ESDK request
func NewCustomDataTransmissionToESDKCommand(data CustomDataTransmissionToESDKData) *CustomDataTransmissionToESDKCommand {
	return &CustomDataTransmissionToESDKCommand{
		Header:     common.NewHeader(),
		MethodName: "custom_data_transmission_to_esdk",
		DataValue:  data,
	}
}

// Method returns the method name.
func (c *CustomDataTransmissionToESDKCommand) Method() string { return c.MethodName }

// Data returns the command/event data.
func (c *CustomDataTransmissionToESDKCommand) Data() any { return c.DataValue }

// GetHeader implements Command.GetHeader
func (c *CustomDataTransmissionToESDKCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CustomDataTransmissionToPSDKCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *InputBoxTextSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *WidgetValueSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SpeakerAudioPlayStartCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SpeakerPlayModeSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SpeakerPlayStopCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SpeakerPlayVolumeSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SpeakerReplayCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *SpeakerTtsPlayStartCommand) GetHeader() *common.Header {
	return &c.Header
}
