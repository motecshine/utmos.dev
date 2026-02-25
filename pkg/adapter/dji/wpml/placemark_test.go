package wpml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewActionGroup(t *testing.T) {
	actionGroup := NewActionGroup(1, 0, 5)

	assert.NotNil(t, actionGroup)
	assert.Equal(t, 1, actionGroup.ActionGroupID)
	assert.Equal(t, 0, actionGroup.ActionGroupStartIndex)
	assert.Equal(t, 5, actionGroup.ActionGroupEndIndex)
	assert.Equal(t, ActionGroupModeSequence, actionGroup.ActionGroupMode)
	assert.Empty(t, actionGroup.Actions)
}

func TestActionGroup_AddAction(t *testing.T) {
	actionGroup := NewActionGroup(1, 0, 5)
	action := Action{
		ActionID:           1,
		ActionActuatorFunc: ActionTypeTakePhoto,
	}

	actionGroup.AddAction(action)

	assert.Len(t, actionGroup.Actions, 1)
	assert.Equal(t, action, actionGroup.Actions[0])
}

func TestActionGroup_SetTrigger(t *testing.T) {
	actionGroup := NewActionGroup(1, 0, 5)
	param := 10.0

	actionGroup.SetTrigger("reach", &param)

	assert.Equal(t, "reach", actionGroup.ActionTrigger.ActionTriggerType)
	assert.Equal(t, &param, actionGroup.ActionTrigger.ActionTriggerParam)
}

func TestNewWaypoint(t *testing.T) {
	waypoint := NewWaypoint(116.3974, 39.9093, 1)

	assert.NotNil(t, waypoint)
	assert.Equal(t, 1, waypoint.Index)
	assert.NotNil(t, waypoint.Point)
	// Check that coordinates contain the expected values (floating point precision)
	assert.Contains(t, waypoint.Point.Coordinates, "116.3974")
	assert.Contains(t, waypoint.Point.Coordinates, "39.9093")
}

func TestPlacemark_SetHeight(t *testing.T) {
	waypoint := NewWaypoint(116.3974, 39.9093, 1)

	waypoint.SetHeight(100.0, 50.0)

	assert.NotNil(t, waypoint.EllipsoidHeight)
	assert.Equal(t, 100.0, *waypoint.EllipsoidHeight)
	assert.NotNil(t, waypoint.Height)
	assert.Equal(t, 50.0, *waypoint.Height)
}

func TestPlacemark_SetExecuteHeight(t *testing.T) {
	waypoint := NewWaypoint(116.3974, 39.9093, 1)

	waypoint.SetExecuteHeight(75.0)

	assert.NotNil(t, waypoint.ExecuteHeight)
	assert.Equal(t, 75.0, *waypoint.ExecuteHeight)
}

func TestPlacemark_AddActionGroup(t *testing.T) {
	waypoint := NewWaypoint(116.3974, 39.9093, 1)
	actionGroup := *NewActionGroup(1, 0, 5)

	waypoint.AddActionGroup(actionGroup)

	assert.Len(t, waypoint.ActionGroups, 1)
	assert.Equal(t, actionGroup, waypoint.ActionGroups[0])
}
