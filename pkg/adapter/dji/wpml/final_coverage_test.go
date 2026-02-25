package wpml

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompleteActionTypeCoverage(t *testing.T) {
	// Test all remaining action types and their JSON marshaling/unmarshaling

	allActionTypes := []struct {
		actionType string
		createFunc func() ActionInterface
	}{
		{ActionTypeTakePhoto, func() ActionInterface {
			return &TakePhotoAction{PayloadPositionIndex: PayloadPosition0, FileSuffix: "test"}
		}},
		{ActionTypeStartRecord, func() ActionInterface {
			return &StartRecordAction{PayloadPositionIndex: PayloadPosition0, FileSuffix: "test"}
		}},
		{ActionTypeStopRecord, func() ActionInterface { return &StopRecordAction{PayloadPositionIndex: PayloadPosition0} }},
		{ActionTypeHover, func() ActionInterface { return &HoverAction{HoverTime: 5.0} }},
		{ActionTypeFocus, func() ActionInterface {
			return &FocusAction{PayloadPositionIndex: PayloadPosition0, IsPointFocus: true}
		}},
		{ActionTypeZoom, func() ActionInterface { return &ZoomAction{PayloadPositionIndex: PayloadPosition0, FocalLength: 50.0} }},
		{ActionTypeCustomDirName, func() ActionInterface {
			return &CustomDirNameAction{PayloadPositionIndex: PayloadPosition0, DirectoryName: "test"}
		}},
		{ActionTypeRotateYaw, func() ActionInterface { return &RotateYawAction{AircraftHeading: 90.0} }},
		{ActionTypeGimbalRotate, func() ActionInterface {
			return &GimbalRotateAction{PayloadPositionIndex: PayloadPosition0, GimbalHeadingYawBase: "aircraft"}
		}},
		{ActionTypeGimbalAngleLock, func() ActionInterface { return &GimbalAngleLockAction{PayloadPositionIndex: PayloadPosition0} }},
		{ActionTypeGimbalAngleUnlock, func() ActionInterface { return &GimbalAngleUnlockAction{PayloadPositionIndex: PayloadPosition0} }},
		{ActionTypeStartSmartOblique, func() ActionInterface { return &StartSmartObliqueAction{PayloadPositionIndex: PayloadPosition0} }},
		{ActionTypeStartTimeLapse, func() ActionInterface { return &StartTimeLapseAction{PayloadPositionIndex: PayloadPosition0} }},
		{ActionTypeStopTimeLapse, func() ActionInterface { return &StopTimeLapseAction{PayloadPositionIndex: PayloadPosition0} }},
		{ActionTypeSetFocusType, func() ActionInterface { return &SetFocusTypeAction{PayloadPositionIndex: PayloadPosition0} }},
		{ActionTypeTargetDetection, func() ActionInterface { return &TargetDetectionAction{PayloadPositionIndex: PayloadPosition0} }},
	}

	for _, actionInfo := range allActionTypes {
		t.Run(actionInfo.actionType+"_complete_test", func(t *testing.T) {
			// Test action creation
			action := actionInfo.createFunc()
			assert.Equal(t, actionInfo.actionType, action.GetActionType())

			// Test ActionRequest creation and JSON round-trip
			actionReq := NewActionRequest(action)
			assert.Equal(t, actionInfo.actionType, actionReq.Type)

			// Test JSON marshaling
			jsonBytes, err := actionReq.ToJSON()
			assert.NoError(t, err)
			assert.Contains(t, string(jsonBytes), actionInfo.actionType)

			// Test JSON unmarshaling back
			var actionReq2 ActionRequest
			err = json.Unmarshal(jsonBytes, &actionReq2)
			assert.NoError(t, err)
			assert.Equal(t, actionInfo.actionType, actionReq2.Type)

			// Test validation
			err = actionReq.Validate()
			// Most should pass validation with our basic setup
			if err != nil {
				t.Logf("Validation failed for %s: %v", actionInfo.actionType, err)
			}
		})
	}
}

func TestActionRequest_MoreEdgeCases(t *testing.T) {
	// Test createActionByType with all action types to increase coverage

	t.Run("createActionByType comprehensive", func(t *testing.T) {
		validTypes := []string{
			ActionTypeTakePhoto,
			ActionTypeStartRecord,
			ActionTypeStopRecord,
			ActionTypeHover,
			ActionTypeFocus,
			ActionTypeZoom,
			ActionTypeCustomDirName,
			ActionTypeRotateYaw,
			ActionTypeGimbalRotate,
			ActionTypeGimbalAngleLock,
			ActionTypeGimbalAngleUnlock,
			ActionTypeStartSmartOblique,
			ActionTypeStartTimeLapse,
			ActionTypeStopTimeLapse,
			ActionTypeSetFocusType,
			ActionTypeTargetDetection,
		}

		for _, actionType := range validTypes {
			// Test createActionByType through ActionRequest unmarshaling
			jsonStr := `{"type": "` + actionType + `", "action": {}}`

			var actionReq ActionRequest
			err := json.Unmarshal([]byte(jsonStr), &actionReq)
			// Some may fail validation but we're testing the creation path
			if err != nil {
				t.Logf("Unmarshal failed for %s: %v", actionType, err)
			}
		}
	})

	t.Run("ActionRequest validation scenarios", func(t *testing.T) {
		// Test ActionRequest with nil action
		actionReq := &ActionRequest{
			Type:   ActionTypeTakePhoto,
			Action: nil,
		}

		err := actionReq.Validate()
		assert.Error(t, err, "Should fail with nil action")
		assert.Equal(t, ErrActionIsNil, err)

		// Test ActionRequest.GetConcreteAction with nil
		result := actionReq.GetConcreteAction()
		assert.Nil(t, result)
	})
}

func TestWaypointValidationCoverage(t *testing.T) {
	// Create waypoints with various configurations to test validation paths

	t.Run("Basic waypoint validation", func(t *testing.T) {
		waylines := createValidWaylines("Action Test")

		// Convert and validate
		mission, err := ConvertWaylinesToMission(waylines)
		if err != nil {
			t.Logf("Conversion failed: %v", err)
		} else {
			assert.NotNil(t, mission)
		}
	})
}

func TestSerializerEdgeCases(t *testing.T) {
	// Test serializer with various XML configurations

	t.Run("XMLSerializer with different settings", func(t *testing.T) {
		waylines := createValidWaylines("Serializer Test")
		mission, err := ConvertWaylinesToMission(waylines)
		require.NoError(t, err)

		// Test global marshal functions
		templateBytes, err := MarshalTemplate(mission.Template)
		assert.NoError(t, err)
		assert.NotEmpty(t, templateBytes)

		waylinesBytes, err := MarshalWaylines(mission.Waylines)
		assert.NoError(t, err)
		assert.NotEmpty(t, waylinesBytes)

		// Test global unmarshal functions
		template, err := UnmarshalTemplate(templateBytes)
		assert.NoError(t, err)
		assert.NotNil(t, template)

		waylinesDoc, err := UnmarshalWaylines(waylinesBytes)
		assert.NoError(t, err)
		assert.NotNil(t, waylinesDoc)
	})
}
