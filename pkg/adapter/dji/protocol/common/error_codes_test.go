package common

import (
	"strings"
	"testing"

	pkgerrors "github.com/utmos/utmos/pkg/errors"
)

func TestIsSuccess(t *testing.T) {
	tests := []struct {
		name string
		code DJIErrorCode
		want bool
	}{
		{name: "zero is success", code: DJI_ERR_SUCCESS, want: true},
		{name: "general failure is not success", code: DJI_ERR_GENERAL_FAILURE, want: false},
		{name: "command not supported is not success", code: DJI_ERR_COMMAND_NOT_SUPPORTED, want: false},
		{name: "arbitrary negative code is not success", code: DJIErrorCode(-1), want: false},
		{name: "arbitrary positive code is not success", code: DJIErrorCode(999999), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSuccess(tt.code); got != tt.want {
				t.Errorf("IsSuccess(%d) = %v, want %v", tt.code, got, tt.want)
			}
		})
	}
}

func TestMapDJIError_KnownCodes(t *testing.T) {
	tests := []struct {
		name         string
		djiCode      DJIErrorCode
		wantPlatform pkgerrors.ErrorCode
		wantContains string
	}{
		{
			name:         "success maps to zero platform code",
			djiCode:      DJI_ERR_SUCCESS,
			wantPlatform: PLATFORM_ERR_SUCCESS,
			wantContains: "success",
		},
		{
			name:         "general failure maps to device error",
			djiCode:      DJI_ERR_GENERAL_FAILURE,
			wantPlatform: PLATFORM_ERR_DEVICE_GENERAL,
			wantContains: "general failure",
		},
		{
			name:         "command not supported maps to protocol error",
			djiCode:      DJI_ERR_COMMAND_NOT_SUPPORTED,
			wantPlatform: PLATFORM_ERR_COMMAND_NOT_SUPPORTED,
			wantContains: "command not supported",
		},
		{
			name:         "command timeout maps to protocol error",
			djiCode:      DJI_ERR_COMMAND_TIMEOUT,
			wantPlatform: PLATFORM_ERR_COMMAND_TIMEOUT,
			wantContains: "command timeout",
		},
		{
			name:         "device busy maps to device error",
			djiCode:      DJI_ERR_DEVICE_BUSY,
			wantPlatform: PLATFORM_ERR_DEVICE_BUSY,
			wantContains: "device busy",
		},
		{
			name:         "parameter error maps to protocol error",
			djiCode:      DJI_ERR_PARAMETER_ERROR,
			wantPlatform: PLATFORM_ERR_PARAMETER_ERROR,
			wantContains: "parameter error",
		},
		{
			name:         "device offline maps to device error",
			djiCode:      DJI_ERR_DEVICE_OFFLINE,
			wantPlatform: PLATFORM_ERR_DEVICE_OFFLINE,
			wantContains: "device offline",
		},
		{
			name:         "flight control error maps to device error",
			djiCode:      DJI_ERR_FLIGHT_CONTROL,
			wantPlatform: PLATFORM_ERR_FLIGHT_CONTROL,
			wantContains: "flight control error",
		},
		{
			name:         "battery low maps to device error",
			djiCode:      DJI_ERR_BATTERY_LOW,
			wantPlatform: PLATFORM_ERR_BATTERY_LOW,
			wantContains: "battery low",
		},
		{
			name:         "camera error maps to device error",
			djiCode:      DJI_ERR_CAMERA_ERROR,
			wantPlatform: PLATFORM_ERR_CAMERA_ERROR,
			wantContains: "camera error",
		},
		{
			name:         "wayline error maps to device error",
			djiCode:      DJI_ERR_WAYLINE_ERROR,
			wantPlatform: PLATFORM_ERR_WAYLINE_ERROR,
			wantContains: "wayline error",
		},
		{
			name:         "file upload failed maps to resource error",
			djiCode:      DJI_ERR_FILE_UPLOAD_FAILED,
			wantPlatform: PLATFORM_ERR_FILE_UPLOAD_FAILED,
			wantContains: "file upload failed",
		},
		{
			name:         "firmware check failed maps to resource error",
			djiCode:      DJI_ERR_FIRMWARE_CHECK_FAILED,
			wantPlatform: PLATFORM_ERR_FIRMWARE_CHECK_FAILED,
			wantContains: "firmware check failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pe := MapDJIError(tt.djiCode)
			if pe == nil {
				t.Fatal("MapDJIError returned nil")
			}
			if pe.Code != tt.wantPlatform {
				t.Errorf("Code = %d, want %d", pe.Code, tt.wantPlatform)
			}
			if pe.DJICode != tt.djiCode {
				t.Errorf("DJICode = %d, want %d", pe.DJICode, tt.djiCode)
			}
			if !strings.Contains(strings.ToLower(pe.Message), tt.wantContains) {
				t.Errorf("Message %q does not contain %q", pe.Message, tt.wantContains)
			}
		})
	}
}

func TestMapDJIError_UnknownCode(t *testing.T) {
	unknownCodes := []DJIErrorCode{
		DJIErrorCode(999999),
		DJIErrorCode(-1),
		DJIErrorCode(42),
		DJIErrorCode(600000),
	}

	for _, code := range unknownCodes {
		pe := MapDJIError(code)
		if pe == nil {
			t.Fatalf("MapDJIError(%d) returned nil", code)
		}
		if pe.Code != PLATFORM_ERR_DEVICE_GENERAL {
			t.Errorf("MapDJIError(%d).Code = %d, want %d (generic device error)", code, pe.Code, PLATFORM_ERR_DEVICE_GENERAL)
		}
		if pe.DJICode != code {
			t.Errorf("MapDJIError(%d).DJICode = %d, want %d", code, pe.DJICode, code)
		}
		if !strings.Contains(pe.Message, "unknown DJI error") {
			t.Errorf("MapDJIError(%d).Message = %q, want it to contain 'unknown DJI error'", code, pe.Message)
		}
	}
}

func TestPlatformError_ErrorInterface(t *testing.T) {
	pe := &PlatformError{
		Code:    PLATFORM_ERR_DEVICE_BUSY,
		Message: "device busy",
		DJICode: DJI_ERR_DEVICE_BUSY,
	}

	// PlatformError must implement the error interface.
	var err error = pe
	errStr := err.Error()

	if !strings.Contains(errStr, "1002") {
		t.Errorf("Error() = %q, want it to contain platform code '1002'", errStr)
	}
	if !strings.Contains(errStr, "device busy") {
		t.Errorf("Error() = %q, want it to contain message 'device busy'", errStr)
	}
	if !strings.Contains(errStr, "514003") {
		t.Errorf("Error() = %q, want it to contain DJI code '514003'", errStr)
	}
}

func TestPlatformError_ErrorFormatting(t *testing.T) {
	tests := []struct {
		name       string
		pe         *PlatformError
		wantSubstr []string
	}{
		{
			name: "success error formatting",
			pe: &PlatformError{
				Code:    PLATFORM_ERR_SUCCESS,
				Message: "success",
				DJICode: DJI_ERR_SUCCESS,
			},
			wantSubstr: []string{"0", "success"},
		},
		{
			name: "device error formatting includes DJI code",
			pe: &PlatformError{
				Code:    PLATFORM_ERR_COMMAND_TIMEOUT,
				Message: "command timeout",
				DJICode: DJI_ERR_COMMAND_TIMEOUT,
			},
			wantSubstr: []string{"2001", "command timeout", "514002"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.pe.Error()
			for _, substr := range tt.wantSubstr {
				if !strings.Contains(errStr, substr) {
					t.Errorf("Error() = %q, want it to contain %q", errStr, substr)
				}
			}
		})
	}
}

func TestDJIErrorCodeConstants(t *testing.T) {
	// Verify DJI error code constant values match the DJI Cloud API spec.
	tests := []struct {
		name string
		code DJIErrorCode
		want int
	}{
		{name: "success", code: DJI_ERR_SUCCESS, want: 0},
		{name: "general failure", code: DJI_ERR_GENERAL_FAILURE, want: 1},
		{name: "command not supported", code: DJI_ERR_COMMAND_NOT_SUPPORTED, want: 514001},
		{name: "command timeout", code: DJI_ERR_COMMAND_TIMEOUT, want: 514002},
		{name: "device busy", code: DJI_ERR_DEVICE_BUSY, want: 514003},
		{name: "parameter error", code: DJI_ERR_PARAMETER_ERROR, want: 514004},
		{name: "device offline", code: DJI_ERR_DEVICE_OFFLINE, want: 514005},
		{name: "flight control", code: DJI_ERR_FLIGHT_CONTROL, want: 514100},
		{name: "battery low", code: DJI_ERR_BATTERY_LOW, want: 514101},
		{name: "camera error", code: DJI_ERR_CAMERA_ERROR, want: 514200},
		{name: "wayline error", code: DJI_ERR_WAYLINE_ERROR, want: 514300},
		{name: "file upload failed", code: DJI_ERR_FILE_UPLOAD_FAILED, want: 316001},
		{name: "firmware check failed", code: DJI_ERR_FIRMWARE_CHECK_FAILED, want: 321000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.code) != tt.want {
				t.Errorf("%s = %d, want %d", tt.name, int(tt.code), tt.want)
			}
		})
	}
}

func TestPlatformErrorCodeRanges(t *testing.T) {
	// Verify platform error codes fall within their designated ranges.
	tests := []struct {
		name    string
		code    pkgerrors.ErrorCode
		minCode int
		maxCode int
	}{
		{name: "success is zero", code: PLATFORM_ERR_SUCCESS, minCode: 0, maxCode: 0},
		{name: "device general in device range", code: PLATFORM_ERR_DEVICE_GENERAL, minCode: 1000, maxCode: 1999},
		{name: "device busy in device range", code: PLATFORM_ERR_DEVICE_BUSY, minCode: 1000, maxCode: 1999},
		{name: "device offline in device range", code: PLATFORM_ERR_DEVICE_OFFLINE, minCode: 1000, maxCode: 1999},
		{name: "flight control in device range", code: PLATFORM_ERR_FLIGHT_CONTROL, minCode: 1000, maxCode: 1999},
		{name: "battery low in device range", code: PLATFORM_ERR_BATTERY_LOW, minCode: 1000, maxCode: 1999},
		{name: "camera error in device range", code: PLATFORM_ERR_CAMERA_ERROR, minCode: 1000, maxCode: 1999},
		{name: "wayline error in device range", code: PLATFORM_ERR_WAYLINE_ERROR, minCode: 1000, maxCode: 1999},
		{name: "command not supported in protocol range", code: PLATFORM_ERR_COMMAND_NOT_SUPPORTED, minCode: 2000, maxCode: 2999},
		{name: "command timeout in protocol range", code: PLATFORM_ERR_COMMAND_TIMEOUT, minCode: 2000, maxCode: 2999},
		{name: "parameter error in protocol range", code: PLATFORM_ERR_PARAMETER_ERROR, minCode: 2000, maxCode: 2999},
		{name: "file upload failed in resource range", code: PLATFORM_ERR_FILE_UPLOAD_FAILED, minCode: 4000, maxCode: 4999},
		{name: "firmware check failed in resource range", code: PLATFORM_ERR_FIRMWARE_CHECK_FAILED, minCode: 4000, maxCode: 4999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := int(tt.code)
			if code < tt.minCode || code > tt.maxCode {
				t.Errorf("platform code %d not in range [%d, %d]", code, tt.minCode, tt.maxCode)
			}
		})
	}
}

func TestMapDJIError_SuccessIsNotError(t *testing.T) {
	pe := MapDJIError(DJI_ERR_SUCCESS)
	if pe.Code != PLATFORM_ERR_SUCCESS {
		t.Errorf("success should map to platform success code 0, got %d", pe.Code)
	}
}
