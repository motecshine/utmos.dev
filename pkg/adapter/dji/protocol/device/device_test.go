package device

import (
	"testing"
)

func TestCoverOpenCommand(t *testing.T) {
	cmd := NewCoverOpenCommand()

	if cmd.Method() != "cover_open" {
		t.Errorf("Method() = %v, want cover_open", cmd.Method())
	}

	if cmd.GetHeader() == nil {
		t.Error("GetHeader() should not return nil")
	}
}

func TestCoverCloseCommand(t *testing.T) {
	cmd := NewCoverCloseCommand()

	if cmd.Method() != "cover_close" {
		t.Errorf("Method() = %v, want cover_close", cmd.Method())
	}
}

func TestDeviceRebootCommand(t *testing.T) {
	cmd := NewDeviceRebootCommand()

	if cmd.Method() != "device_reboot" {
		t.Errorf("Method() = %v, want device_reboot", cmd.Method())
	}
}

func TestDeviceFormatCommand(t *testing.T) {
	cmd := NewDeviceFormatCommand()

	if cmd.Method() != "device_format" {
		t.Errorf("Method() = %v, want device_format", cmd.Method())
	}
}

func TestDroneOpenCommand(t *testing.T) {
	cmd := NewDroneOpenCommand()

	if cmd.Method() != "drone_open" {
		t.Errorf("Method() = %v, want drone_open", cmd.Method())
	}
}

func TestDroneCloseCommand(t *testing.T) {
	cmd := NewDroneCloseCommand()

	if cmd.Method() != "drone_close" {
		t.Errorf("Method() = %v, want drone_close", cmd.Method())
	}
}
