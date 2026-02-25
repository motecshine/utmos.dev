package wpml

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionRequest_MarshalJSON(t *testing.T) {
	action := &TakePhotoAction{
		PayloadPositionIndex:      PayloadPosition0,
		FileSuffix:                "test",
		UseGlobalPayloadLensIndex: true,
	}

	actionReq := &ActionRequest{
		Type:   ActionTypeTakePhoto,
		Action: action,
	}

	jsonBytes, err := json.Marshal(actionReq)
	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), ActionTypeTakePhoto)
}

func TestActionRequest_UnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"type": "takePhoto",
		"action": {
			"payload_position_index": 0,
			"file_suffix": "test",
			"use_global_payload_lens_index": true
		}
	}`

	var actionReq ActionRequest
	err := json.Unmarshal([]byte(jsonStr), &actionReq)
	assert.NoError(t, err)
	assert.Equal(t, ActionTypeTakePhoto, actionReq.Type)

	takePhotoAction, ok := actionReq.Action.(*TakePhotoAction)
	require.True(t, ok)
	assert.Equal(t, PayloadPosition0, takePhotoAction.PayloadPositionIndex)
	assert.Equal(t, "test", takePhotoAction.FileSuffix)
	assert.True(t, takePhotoAction.UseGlobalPayloadLensIndex)
}

func TestActionRequest_UnmarshalJSON_InvalidType(t *testing.T) {
	jsonStr := `{
		"type": "unknownAction",
		"action": {}
	}`

	var actionReq ActionRequest
	err := json.Unmarshal([]byte(jsonStr), &actionReq)
	assert.Error(t, err)
}

func TestNewActionRequest(t *testing.T) {
	action := &HoverAction{
		HoverTime: 5.0,
	}

	actionReq := NewActionRequest(action)
	assert.NotNil(t, actionReq)
	assert.Equal(t, ActionTypeHover, actionReq.Type)
	assert.Equal(t, action, actionReq.Action)
}

func TestTypedActionRequest(t *testing.T) {
	action := &ZoomAction{
		PayloadPositionIndex: PayloadPosition0,
		FocalLength:          50.0,
	}

	actionReq := TypedActionRequest(action)
	assert.NotNil(t, actionReq)
	assert.Equal(t, ActionTypeZoom, actionReq.Type)
	assert.Equal(t, action, actionReq.Action)
}

func TestActionRequestFromJSON(t *testing.T) {
	jsonStr := `{
		"type": "gimbalRotate",
		"action": {
			"payload_position_index": 0,
			"gimbal_heading_yaw_base": "north",
			"gimbal_rotate_mode": "absolute",
			"gimbal_pitch_rotate_enable": true,
			"gimbal_pitch_rotate_angle": -30.0,
			"gimbal_yaw_rotate_enable": true,
			"gimbal_yaw_rotate_angle": 90.0,
			"gimbal_rotate_time_enable": true,
			"gimbal_rotate_time": 2.0
		}
	}`

	actionReq, err := ActionRequestFromJSON([]byte(jsonStr))
	assert.NoError(t, err)
	assert.NotNil(t, actionReq)
	assert.Equal(t, ActionTypeGimbalRotate, actionReq.Type)

	gimbalAction, ok := actionReq.Action.(*GimbalRotateAction)
	require.True(t, ok)
	assert.Equal(t, PayloadPosition0, gimbalAction.PayloadPositionIndex)
	assert.Equal(t, "north", gimbalAction.GimbalHeadingYawBase)
}

func TestActionRequest_ToJSON(t *testing.T) {
	action := &CustomDirNameAction{
		PayloadPositionIndex: PayloadPosition0,
		DirectoryName:        "my_photos",
	}

	actionReq := &ActionRequest{
		Type:   ActionTypeCustomDirName,
		Action: action,
	}

	jsonBytes, err := actionReq.ToJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(jsonBytes), "customDirName")
	assert.Contains(t, string(jsonBytes), "my_photos")
}

func TestActionRequest_GetConcreteAction(t *testing.T) {
	action := &FocusAction{
		PayloadPositionIndex: PayloadPosition0,
		IsPointFocus:         true,
		FocusX:               0.5,
		FocusY:               0.5,
	}

	actionReq := &ActionRequest{
		Type:   ActionTypeFocus,
		Action: action,
	}

	result := actionReq.GetConcreteAction()
	assert.Equal(t, action, result)
}

func TestGetTypedAction(t *testing.T) {
	action := &RotateYawAction{
		AircraftHeading: 180.0,
	}

	actionReq := &ActionRequest{
		Type:   ActionTypeRotateYaw,
		Action: action,
	}

	typedAction, err := GetTypedAction[*RotateYawAction](actionReq)
	assert.NoError(t, err)
	assert.Equal(t, action, typedAction)
	assert.Equal(t, 180.0, typedAction.AircraftHeading)
}

func TestGetTypedAction_TypeError(t *testing.T) {
	action := &TakePhotoAction{
		PayloadPositionIndex: PayloadPosition0,
	}

	actionReq := &ActionRequest{
		Type:   ActionTypeTakePhoto,
		Action: action,
	}

	// Try to get wrong type
	_, err := GetTypedAction[*HoverAction](actionReq)
	assert.Error(t, err)
}

func TestActionRequest_Validate(t *testing.T) {
	action := &TakePhotoAction{
		PayloadPositionIndex:      PayloadPosition0,
		FileSuffix:                "valid",
		UseGlobalPayloadLensIndex: true,
	}

	actionReq := &ActionRequest{
		Type:   ActionTypeTakePhoto,
		Action: action,
	}

	err := actionReq.Validate()
	assert.NoError(t, err)
}

func TestActionRequest_GetActionType(t *testing.T) {
	// Test with Action set
	action := &TakePhotoAction{
		PayloadPositionIndex: PayloadPosition0,
		FileSuffix:           "test",
	}

	actionReq := &ActionRequest{
		Type:   ActionTypeTakePhoto,
		Action: action,
	}

	actionType := actionReq.GetActionType()
	assert.Equal(t, ActionTypeTakePhoto, actionType)

	// Test with nil Action - should return Type field
	actionReq2 := &ActionRequest{
		Type:   ActionTypeTakePhoto,
		Action: nil,
	}

	actionType2 := actionReq2.GetActionType()
	assert.Equal(t, ActionTypeTakePhoto, actionType2)
}

func TestActionTypes_GetActionType(t *testing.T) {
	// Test all action types that have 0% coverage on GetActionType

	tests := []struct {
		name     string
		action   ActionInterface
		expected string
	}{
		{
			name:     "GimbalAngleLockAction",
			action:   &GimbalAngleLockAction{PayloadPositionIndex: PayloadPosition0},
			expected: ActionTypeGimbalAngleLock,
		},
		{
			name:     "GimbalAngleUnlockAction",
			action:   &GimbalAngleUnlockAction{PayloadPositionIndex: PayloadPosition0},
			expected: ActionTypeGimbalAngleUnlock,
		},
		{
			name:     "StartSmartObliqueAction",
			action:   &StartSmartObliqueAction{PayloadPositionIndex: PayloadPosition0},
			expected: ActionTypeStartSmartOblique,
		},
		{
			name:     "StartTimeLapseAction",
			action:   &StartTimeLapseAction{PayloadPositionIndex: PayloadPosition0},
			expected: ActionTypeStartTimeLapse,
		},
		{
			name:     "StopTimeLapseAction",
			action:   &StopTimeLapseAction{PayloadPositionIndex: PayloadPosition0},
			expected: ActionTypeStopTimeLapse,
		},
		{
			name:     "SetFocusTypeAction",
			action:   &SetFocusTypeAction{PayloadPositionIndex: PayloadPosition0},
			expected: ActionTypeSetFocusType,
		},
		{
			name:     "TargetDetectionAction",
			action:   &TargetDetectionAction{PayloadPositionIndex: PayloadPosition0},
			expected: ActionTypeTargetDetection,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actionType := tt.action.GetActionType()
			assert.Equal(t, tt.expected, actionType)
		})
	}
}
