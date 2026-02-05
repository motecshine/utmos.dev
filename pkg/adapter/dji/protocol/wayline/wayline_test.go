package wayline

import (
	"testing"
)

func TestCreateCommand(t *testing.T) {
	data := CreateData{
		FlighttaskID:       "task-123",
		TaskType:           "immediate",
		WaylineType:        "wayline",
		OutOfControlAction: "execute_go_home",
	}
	cmd := NewCreateCommand(data)

	if cmd.Method() != "flighttask_create" {
		t.Errorf("Method() = %v, want flighttask_create", cmd.Method())
	}

	if cmd.GetHeader() == nil {
		t.Error("GetHeader() should not return nil")
	}

	cmdData := cmd.Data()
	if cmdData == nil {
		t.Error("Data() should not return nil")
	}
}

func TestPrepareCommand(t *testing.T) {
	data := PrepareData{
		FlightId: "task-123",
	}
	cmd := NewPrepareCommand(data)

	if cmd.Method() != "flighttask_prepare" {
		t.Errorf("Method() = %v, want flighttask_prepare", cmd.Method())
	}
}

func TestExecuteCommand(t *testing.T) {
	data := ExecuteData{
		FlighttaskID: "task-123",
	}
	cmd := NewExecuteCommand(data)

	if cmd.Method() != "flighttask_execute" {
		t.Errorf("Method() = %v, want flighttask_execute", cmd.Method())
	}
}

func TestPauseCommand(t *testing.T) {
	cmd := NewPauseCommand()

	if cmd.Method() != "flighttask_pause" {
		t.Errorf("Method() = %v, want flighttask_pause", cmd.Method())
	}
}

func TestRecoveryCommand(t *testing.T) {
	cmd := NewRecoveryCommand()

	if cmd.Method() != "flighttask_recovery" {
		t.Errorf("Method() = %v, want flighttask_recovery", cmd.Method())
	}
}

func TestUndoCommand(t *testing.T) {
	data := UndoData{
		FlightIds: []string{"task-123"},
	}
	cmd := NewUndoCommand(data)

	if cmd.Method() != "flighttask_undo" {
		t.Errorf("Method() = %v, want flighttask_undo", cmd.Method())
	}
}

func TestReturnHomeCommand(t *testing.T) {
	cmd := NewReturnHomeCommand()

	if cmd.Method() != "return_home" {
		t.Errorf("Method() = %v, want return_home", cmd.Method())
	}
}

func TestCancelReturnHomeCommand(t *testing.T) {
	cmd := NewCancelReturnHomeCommand()

	if cmd.Method() != "return_home_cancel" {
		t.Errorf("Method() = %v, want return_home_cancel", cmd.Method())
	}
}
