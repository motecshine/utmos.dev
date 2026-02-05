package firmware

import "github.com/utmos/utmos/pkg/adapter/dji/protocol/common"

// ===============================
// Firmware Upgrade Commands
// ===============================

// FirmwareUpgradeDevice represents a single device firmware upgrade entry
type FirmwareUpgradeDevice struct {
	SN                  string `json:"sn"`                    // Device serial number (SN)
	ProductVersion      string `json:"product_version"`       // Firmware version
	FileURL             string `json:"file_url"`              // Firmware file download URL
	MD5                 string `json:"md5"`                   // Firmware file MD5 checksum
	FileSize            int64  `json:"file_size"`             // Firmware file size in bytes
	FileName            string `json:"file_name"`             // Firmware file name
	FirmwareUpgradeType int    `json:"firmware_upgrade_type"` // Firmware upgrade type: 2=consistent upgrade, 3=normal upgrade
}

// OTACreateData represents the OTA create data
type OTACreateData struct {
	Devices []FirmwareUpgradeDevice `json:"devices"` // Array of devices to upgrade (max 2 devices)
}

// OTACreateCommand represents the OTA create command
type OTACreateCommand struct {
	common.Header
	MethodName string        `json:"method"`
	DataValue  OTACreateData `json:"data"`
}

// NewOTACreateCommand creates a new OTA create command
func NewOTACreateCommand(data OTACreateData) *OTACreateCommand {
	return &OTACreateCommand{
		Header:     common.NewHeader(),
		MethodName: "ota_create",
		DataValue:  data,
	}
}

func (c *OTACreateCommand) Method() string { return c.MethodName }
func (c *OTACreateCommand) Data() any      { return c.DataValue }

// GetHeader implements Command.GetHeader
func (c *OTACreateCommand) GetHeader() *common.Header {
	return &c.Header
}
