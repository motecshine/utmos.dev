package common

import (
	"fmt"

	pkgerrors "github.com/utmos/utmos/pkg/errors"
)

// DJIErrorCode represents an error code returned by DJI Cloud API devices.
type DJIErrorCode int

// DJI Cloud API error code constants.
const (
	// DJI_ERR_SUCCESS indicates the operation completed successfully.
	DJI_ERR_SUCCESS DJIErrorCode = 0
	// DJI_ERR_GENERAL_FAILURE indicates a general/unspecified failure.
	DJI_ERR_GENERAL_FAILURE DJIErrorCode = 1
	// DJI_ERR_COMMAND_NOT_SUPPORTED indicates the command is not supported by the device.
	DJI_ERR_COMMAND_NOT_SUPPORTED DJIErrorCode = 514001
	// DJI_ERR_COMMAND_TIMEOUT indicates the command execution timed out.
	DJI_ERR_COMMAND_TIMEOUT DJIErrorCode = 514002
	// DJI_ERR_DEVICE_BUSY indicates the device is currently busy.
	DJI_ERR_DEVICE_BUSY DJIErrorCode = 514003
	// DJI_ERR_PARAMETER_ERROR indicates invalid parameters were provided.
	DJI_ERR_PARAMETER_ERROR DJIErrorCode = 514004
	// DJI_ERR_DEVICE_OFFLINE indicates the device is offline.
	DJI_ERR_DEVICE_OFFLINE DJIErrorCode = 514005
	// DJI_ERR_FLIGHT_CONTROL indicates a flight control error.
	DJI_ERR_FLIGHT_CONTROL DJIErrorCode = 514100
	// DJI_ERR_BATTERY_LOW indicates the battery level is too low.
	DJI_ERR_BATTERY_LOW DJIErrorCode = 514101
	// DJI_ERR_CAMERA_ERROR indicates a camera-related error.
	DJI_ERR_CAMERA_ERROR DJIErrorCode = 514200
	// DJI_ERR_WAYLINE_ERROR indicates a wayline/route-related error.
	DJI_ERR_WAYLINE_ERROR DJIErrorCode = 514300
	// DJI_ERR_FILE_UPLOAD_FAILED indicates a file upload failure.
	DJI_ERR_FILE_UPLOAD_FAILED DJIErrorCode = 316001
	// DJI_ERR_FIRMWARE_CHECK_FAILED indicates a firmware check failure.
	DJI_ERR_FIRMWARE_CHECK_FAILED DJIErrorCode = 321000
)

// Platform error code constants organized by range:
//   - 0:         Success
//   - 1000-1999: Device errors
//   - 2000-2999: Protocol errors
//   - 3000-3999: Service errors
//   - 4000-4999: Resource errors
const (
	// PLATFORM_ERR_SUCCESS indicates the operation completed successfully.
	PLATFORM_ERR_SUCCESS pkgerrors.ErrorCode = 0

	// Device errors (1000-1999).

	// PLATFORM_ERR_DEVICE_GENERAL indicates a general device error.
	PLATFORM_ERR_DEVICE_GENERAL pkgerrors.ErrorCode = 1000
	// PLATFORM_ERR_DEVICE_BUSY indicates the device is currently busy.
	PLATFORM_ERR_DEVICE_BUSY pkgerrors.ErrorCode = 1002
	// PLATFORM_ERR_DEVICE_OFFLINE indicates the device is offline.
	PLATFORM_ERR_DEVICE_OFFLINE pkgerrors.ErrorCode = 1003
	// PLATFORM_ERR_FLIGHT_CONTROL indicates a flight control error.
	PLATFORM_ERR_FLIGHT_CONTROL pkgerrors.ErrorCode = 1100
	// PLATFORM_ERR_BATTERY_LOW indicates the battery level is too low.
	PLATFORM_ERR_BATTERY_LOW pkgerrors.ErrorCode = 1101
	// PLATFORM_ERR_CAMERA_ERROR indicates a camera-related error.
	PLATFORM_ERR_CAMERA_ERROR pkgerrors.ErrorCode = 1200
	// PLATFORM_ERR_WAYLINE_ERROR indicates a wayline/route-related error.
	PLATFORM_ERR_WAYLINE_ERROR pkgerrors.ErrorCode = 1300

	// Protocol errors (2000-2999).

	// PLATFORM_ERR_COMMAND_NOT_SUPPORTED indicates the command is not supported.
	PLATFORM_ERR_COMMAND_NOT_SUPPORTED pkgerrors.ErrorCode = 2000
	// PLATFORM_ERR_COMMAND_TIMEOUT indicates the command execution timed out.
	PLATFORM_ERR_COMMAND_TIMEOUT pkgerrors.ErrorCode = 2001
	// PLATFORM_ERR_PARAMETER_ERROR indicates invalid parameters were provided.
	PLATFORM_ERR_PARAMETER_ERROR pkgerrors.ErrorCode = 2002

	// Resource errors (4000-4999).

	// PLATFORM_ERR_FILE_UPLOAD_FAILED indicates a file upload failure.
	PLATFORM_ERR_FILE_UPLOAD_FAILED pkgerrors.ErrorCode = 4000
	// PLATFORM_ERR_FIRMWARE_CHECK_FAILED indicates a firmware check failure.
	PLATFORM_ERR_FIRMWARE_CHECK_FAILED pkgerrors.ErrorCode = 4001
)

// PlatformError represents a platform-standard error mapped from a DJI error code.
type PlatformError struct {
	Code    pkgerrors.ErrorCode
	Message string
	DJICode DJIErrorCode
}

// Error implements the error interface for PlatformError.
func (e *PlatformError) Error() string {
	return fmt.Sprintf("[%d] %s (DJI code: %d)", e.Code, e.Message, e.DJICode)
}

// djiErrorMapping maps DJI error codes to their platform error code and message.
type djiErrorMapping struct {
	platformCode pkgerrors.ErrorCode
	message      string
}

// djiErrorMap holds the mapping from DJI error codes to platform errors.
var djiErrorMap = map[DJIErrorCode]djiErrorMapping{
	DJI_ERR_SUCCESS:              {platformCode: PLATFORM_ERR_SUCCESS, message: "success"},
	DJI_ERR_GENERAL_FAILURE:      {platformCode: PLATFORM_ERR_DEVICE_GENERAL, message: "general failure"},
	DJI_ERR_COMMAND_NOT_SUPPORTED: {platformCode: PLATFORM_ERR_COMMAND_NOT_SUPPORTED, message: "command not supported"},
	DJI_ERR_COMMAND_TIMEOUT:      {platformCode: PLATFORM_ERR_COMMAND_TIMEOUT, message: "command timeout"},
	DJI_ERR_DEVICE_BUSY:          {platformCode: PLATFORM_ERR_DEVICE_BUSY, message: "device busy"},
	DJI_ERR_PARAMETER_ERROR:      {platformCode: PLATFORM_ERR_PARAMETER_ERROR, message: "parameter error"},
	DJI_ERR_DEVICE_OFFLINE:       {platformCode: PLATFORM_ERR_DEVICE_OFFLINE, message: "device offline"},
	DJI_ERR_FLIGHT_CONTROL:       {platformCode: PLATFORM_ERR_FLIGHT_CONTROL, message: "flight control error"},
	DJI_ERR_BATTERY_LOW:          {platformCode: PLATFORM_ERR_BATTERY_LOW, message: "battery low"},
	DJI_ERR_CAMERA_ERROR:         {platformCode: PLATFORM_ERR_CAMERA_ERROR, message: "camera error"},
	DJI_ERR_WAYLINE_ERROR:        {platformCode: PLATFORM_ERR_WAYLINE_ERROR, message: "wayline error"},
	DJI_ERR_FILE_UPLOAD_FAILED:   {platformCode: PLATFORM_ERR_FILE_UPLOAD_FAILED, message: "file upload failed"},
	DJI_ERR_FIRMWARE_CHECK_FAILED: {platformCode: PLATFORM_ERR_FIRMWARE_CHECK_FAILED, message: "firmware check failed"},
}

// MapDJIError maps a DJI error code to a platform-standard PlatformError.
// Unknown DJI error codes are mapped to a generic device error with the original DJI code preserved.
func MapDJIError(code DJIErrorCode) *PlatformError {
	if mapping, ok := djiErrorMap[code]; ok {
		return &PlatformError{
			Code:    mapping.platformCode,
			Message: mapping.message,
			DJICode: code,
		}
	}

	return &PlatformError{
		Code:    PLATFORM_ERR_DEVICE_GENERAL,
		Message: fmt.Sprintf("unknown DJI error (code: %d)", code),
		DJICode: code,
	}
}

// IsSuccess returns true if the DJI error code indicates success.
func IsSuccess(code DJIErrorCode) bool {
	return code == DJI_ERR_SUCCESS
}
