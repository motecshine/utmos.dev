package psdk

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// PSDK Payload Commands
// ===============================

// PsdkWidgetValueSetData represents the PSDK widget value set data
type PsdkWidgetValueSetData struct {
	PsdkIndex int `json:"psdk_index"` // PSDK payload device index (0-3)
	Index     int `json:"index"`      // Widget index
	Value     int `json:"value"`      // Widget value (defined by developer)
}

// PsdkWidgetValueSetRequest represents the PSDK widget value set request
type PsdkWidgetValueSetCommand struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  PsdkWidgetValueSetData `json:"data"`
}

// NewPsdkWidgetValueSetRequest creates a new PSDK widget value set request
func NewPsdkWidgetValueSetCommand(data PsdkWidgetValueSetData) *PsdkWidgetValueSetCommand {
	return &PsdkWidgetValueSetCommand{
		Header:     common.NewHeader(),
		MethodName: "psdk_widget_value_set",
		DataValue:  data,
	}
}

func (c *PsdkWidgetValueSetCommand) Method() string { return c.MethodName }
func (c *PsdkWidgetValueSetCommand) Data() any      { return c.DataValue }

// PsdkInputBoxTextSetData represents the PSDK input box text set data
type PsdkInputBoxTextSetData struct {
	PsdkIndex int    `json:"psdk_index"` // PSDK payload device index (0-3)
	Value     string `json:"value"`      // Text content (max 128 bytes)
}

// PsdkInputBoxTextSetRequest represents the PSDK input box text set request
type PsdkInputBoxTextSetCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  PsdkInputBoxTextSetData `json:"data"`
}

// NewPsdkInputBoxTextSetRequest creates a new PSDK input box text set request
func NewPsdkInputBoxTextSetCommand(data PsdkInputBoxTextSetData) *PsdkInputBoxTextSetCommand {
	return &PsdkInputBoxTextSetCommand{
		Header:     common.NewHeader(),
		MethodName: "psdk_input_box_text_set",
		DataValue:  data,
	}
}

func (c *PsdkInputBoxTextSetCommand) Method() string { return c.MethodName }
func (c *PsdkInputBoxTextSetCommand) Data() any      { return c.DataValue }

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

// SpeakerAudioPlayStartRequest represents the speaker audio play start request
type SpeakerAudioPlayStartCommand struct {
	common.Header
	MethodName string                    `json:"method"`
	DataValue  SpeakerAudioPlayStartData `json:"data"`
}

// NewSpeakerAudioPlayStartRequest creates a new speaker audio play start request
func NewSpeakerAudioPlayStartCommand(data SpeakerAudioPlayStartData) *SpeakerAudioPlayStartCommand {
	return &SpeakerAudioPlayStartCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_audio_play_start",
		DataValue:  data,
	}
}

func (c *SpeakerAudioPlayStartCommand) Method() string { return c.MethodName }
func (c *SpeakerAudioPlayStartCommand) Data() any      { return c.DataValue }

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

// SpeakerTtsPlayStartRequest represents the speaker TTS play start request
type SpeakerTtsPlayStartCommand struct {
	common.Header
	MethodName string                  `json:"method"`
	DataValue  SpeakerTtsPlayStartData `json:"data"`
}

// NewSpeakerTtsPlayStartRequest creates a new speaker TTS play start request
func NewSpeakerTtsPlayStartCommand(data SpeakerTtsPlayStartData) *SpeakerTtsPlayStartCommand {
	return &SpeakerTtsPlayStartCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_tts_play_start",
		DataValue:  data,
	}
}

func (c *SpeakerTtsPlayStartCommand) Method() string { return c.MethodName }
func (c *SpeakerTtsPlayStartCommand) Data() any      { return c.DataValue }

// SpeakerReplayData represents the speaker replay data
type SpeakerReplayData struct {
	PsdkIndex int `json:"psdk_index"` // PSDK payload device index (min: 0)
}

// SpeakerReplayRequest represents the speaker replay request
type SpeakerReplayCommand struct {
	common.Header
	MethodName string            `json:"method"`
	DataValue  SpeakerReplayData `json:"data"`
}

// NewSpeakerReplayRequest creates a new speaker replay request
func NewSpeakerReplayCommand(data SpeakerReplayData) *SpeakerReplayCommand {
	return &SpeakerReplayCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_replay",
		DataValue:  data,
	}
}

func (c *SpeakerReplayCommand) Method() string { return c.MethodName }
func (c *SpeakerReplayCommand) Data() any      { return c.DataValue }

// SpeakerPlayStopData represents the speaker play stop data
type SpeakerPlayStopData struct {
	PsdkIndex int `json:"psdk_index"` // PSDK payload device index (min: 0)
}

// SpeakerPlayStopRequest represents the speaker play stop request
type SpeakerPlayStopCommand struct {
	common.Header
	MethodName string              `json:"method"`
	DataValue  SpeakerPlayStopData `json:"data"`
}

// NewSpeakerPlayStopRequest creates a new speaker play stop request
func NewSpeakerPlayStopCommand(data SpeakerPlayStopData) *SpeakerPlayStopCommand {
	return &SpeakerPlayStopCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_play_stop",
		DataValue:  data,
	}
}

func (c *SpeakerPlayStopCommand) Method() string { return c.MethodName }
func (c *SpeakerPlayStopCommand) Data() any      { return c.DataValue }

// SpeakerPlayModeSetData represents the speaker play mode set data
type SpeakerPlayModeSetData struct {
	PsdkIndex int `json:"psdk_index"` // PSDK payload device index (0-3)
	PlayMode  int `json:"play_mode"`  // Play mode: 0=single play, 1=loop (single track)
}

// SpeakerPlayModeSetRequest represents the speaker play mode set request
type SpeakerPlayModeSetCommand struct {
	common.Header
	MethodName string                 `json:"method"`
	DataValue  SpeakerPlayModeSetData `json:"data"`
}

// NewSpeakerPlayModeSetRequest creates a new speaker play mode set request
func NewSpeakerPlayModeSetCommand(data SpeakerPlayModeSetData) *SpeakerPlayModeSetCommand {
	return &SpeakerPlayModeSetCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_play_mode_set",
		DataValue:  data,
	}
}

func (c *SpeakerPlayModeSetCommand) Method() string { return c.MethodName }
func (c *SpeakerPlayModeSetCommand) Data() any      { return c.DataValue }

// SpeakerPlayVolumeSetData represents the speaker play volume set data
type SpeakerPlayVolumeSetData struct {
	PsdkIndex  int `json:"psdk_index"` // PSDK payload device index (0-3)
	Index      int `json:"index"`
	PlayVolume int `json:"play_volume"` // Speaker volume (0-100)
}

// SpeakerPlayVolumeSetRequest represents the speaker play volume set request
type SpeakerPlayVolumeSetCommand struct {
	common.Header
	MethodName string                   `json:"method"`
	DataValue  SpeakerPlayVolumeSetData `json:"data"`
}

// NewSpeakerPlayVolumeSetRequest creates a new speaker play volume set request
func NewSpeakerPlayVolumeSetCommand(data SpeakerPlayVolumeSetData) *SpeakerPlayVolumeSetCommand {
	return &SpeakerPlayVolumeSetCommand{
		Header:     common.NewHeader(),
		MethodName: "speaker_play_volume_set",
		DataValue:  data,
	}
}

func (c *SpeakerPlayVolumeSetCommand) Method() string { return c.MethodName }
func (c *SpeakerPlayVolumeSetCommand) Data() any      { return c.DataValue }

// CustomDataTransmissionToPSDKData represents the custom data transmission to PSDK data
type CustomDataTransmissionToPSDKData struct {
	PayloadIndex string      `json:"payload_index"` // Payload index
	DataValue    interface{} `json:"data"`          // Custom data
}

// CustomDataTransmissionToPSDKRequest represents the custom data transmission to PSDK request
type CustomDataTransmissionToPSDKCommand struct {
	common.Header
	MethodName string                           `json:"method"`
	DataValue  CustomDataTransmissionToPSDKData `json:"data"`
}

// NewCustomDataTransmissionToPSDKRequest creates a new custom data transmission to PSDK request
func NewCustomDataTransmissionToPSDKCommand(data CustomDataTransmissionToPSDKData) *CustomDataTransmissionToPSDKCommand {
	return &CustomDataTransmissionToPSDKCommand{
		Header:     common.NewHeader(),
		MethodName: "custom_data_transmission_to_psdk",
		DataValue:  data,
	}
}

func (c *CustomDataTransmissionToPSDKCommand) Method() string { return c.MethodName }
func (c *CustomDataTransmissionToPSDKCommand) Data() any      { return c.DataValue }

// CustomDataTransmissionToESDKData represents the custom data transmission to ESDK data
type CustomDataTransmissionToESDKData struct {
	Data interface{} `json:"data"` // Custom data
}

// CustomDataTransmissionToESDKRequest represents the custom data transmission to ESDK request
type CustomDataTransmissionToESDKCommand struct {
	common.Header
	MethodName string                           `json:"method"`
	DataValue  CustomDataTransmissionToESDKData `json:"data"`
}

// NewCustomDataTransmissionToESDKRequest creates a new custom data transmission to ESDK request
func NewCustomDataTransmissionToESDKCommand(data CustomDataTransmissionToESDKData) *CustomDataTransmissionToESDKCommand {
	return &CustomDataTransmissionToESDKCommand{
		Header:     common.NewHeader(),
		MethodName: "custom_data_transmission_to_esdk",
		DataValue:  data,
	}
}

func (c *CustomDataTransmissionToESDKCommand) Method() string { return c.MethodName }
func (c *CustomDataTransmissionToESDKCommand) Data() any      { return c.DataValue }

// GetHeader implements Command.GetHeader
func (c *CustomDataTransmissionToESDKCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *CustomDataTransmissionToPSDKCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *PsdkInputBoxTextSetCommand) GetHeader() *common.Header {
	return &c.Header
}

// GetHeader implements Command.GetHeader
func (c *PsdkWidgetValueSetCommand) GetHeader() *common.Header {
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
