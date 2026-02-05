package camera

import (
	"testing"
)

func TestCameraPhotoTakeCommand(t *testing.T) {
	data := CameraPhotoTakeData{
		PayloadIndex: "0-0-0",
	}
	cmd := NewCameraPhotoTakeCommand(data)

	if cmd.Method() != "camera_photo_take" {
		t.Errorf("Method() = %v, want camera_photo_take", cmd.Method())
	}

	if cmd.GetHeader() == nil {
		t.Error("GetHeader() should not return nil")
	}

	cmdData := cmd.Data()
	if cmdData == nil {
		t.Error("Data() should not return nil")
	}
}

func TestCameraRecordingStartCommand(t *testing.T) {
	data := CameraRecordingStartData{
		PayloadIndex: "0-0-0",
	}
	cmd := NewCameraRecordingStartCommand(data)

	if cmd.Method() != "camera_recording_start" {
		t.Errorf("Method() = %v, want camera_recording_start", cmd.Method())
	}
}

func TestCameraRecordingStopCommand(t *testing.T) {
	data := CameraRecordingStopData{
		PayloadIndex: "0-0-0",
	}
	cmd := NewCameraRecordingStopCommand(data)

	if cmd.Method() != "camera_recording_stop" {
		t.Errorf("Method() = %v, want camera_recording_stop", cmd.Method())
	}
}

func TestCameraModeSwitchCommand(t *testing.T) {
	data := CameraModeSwitchData{
		PayloadIndex: "0-0-0",
		CameraMode:   0,
	}
	cmd := NewCameraModeSwitchCommand(data)

	if cmd.Method() != "camera_mode_switch" {
		t.Errorf("Method() = %v, want camera_mode_switch", cmd.Method())
	}

	cmdData := cmd.Data()
	if cmdData == nil {
		t.Error("Data() should not return nil")
	}
}

func TestGimbalResetCommand(t *testing.T) {
	data := GimbalResetData{
		PayloadIndex: "0-0-0",
		ResetMode:    0,
	}
	cmd := NewGimbalResetCommand(data)

	if cmd.Method() != "gimbal_reset" {
		t.Errorf("Method() = %v, want gimbal_reset", cmd.Method())
	}
}

func TestIRMeteringModeSetCommand(t *testing.T) {
	data := IRMeteringModeSetData{
		PayloadIndex: "0-0-0",
		Mode:         1,
	}
	cmd := NewIRMeteringModeSetCommand(data)

	if cmd.Method() != "ir_metering_mode_set" {
		t.Errorf("Method() = %v, want ir_metering_mode_set", cmd.Method())
	}

	if cmd.GetHeader() == nil {
		t.Error("GetHeader() should not return nil")
	}
}

func TestIRMeteringPointSetCommand(t *testing.T) {
	data := IRMeteringPointSetData{
		PayloadIndex: "0-0-0",
		X:            0.5,
		Y:            0.5,
	}
	cmd := NewIRMeteringPointSetCommand(data)

	if cmd.Method() != "ir_metering_point_set" {
		t.Errorf("Method() = %v, want ir_metering_point_set", cmd.Method())
	}
}

func TestIRMeteringAreaSetCommand(t *testing.T) {
	data := IRMeteringAreaSetData{
		PayloadIndex: "0-0-0",
		X:            0.1,
		Y:            0.1,
		Width:        0.5,
		Height:       0.5,
	}
	cmd := NewIRMeteringAreaSetCommand(data)

	if cmd.Method() != "ir_metering_area_set" {
		t.Errorf("Method() = %v, want ir_metering_area_set", cmd.Method())
	}
}
