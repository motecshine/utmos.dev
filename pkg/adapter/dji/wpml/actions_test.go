package wpml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTakePhotoAction_GetActionType(t *testing.T) {
	action := &TakePhotoAction{
		PayloadPositionIndex:      PayloadPosition0,
		FileSuffix:                "test",
		UseGlobalPayloadLensIndex: true,
	}

	assert.Equal(t, ActionTypeTakePhoto, action.GetActionType())
}

func TestStartRecordAction_GetActionType(t *testing.T) {
	action := &StartRecordAction{
		PayloadPositionIndex:      PayloadPosition1,
		FileSuffix:                "video",
		UseGlobalPayloadLensIndex: false,
	}

	assert.Equal(t, ActionTypeStartRecord, action.GetActionType())
}

func TestStopRecordAction_GetActionType(t *testing.T) {
	action := &StopRecordAction{
		PayloadPositionIndex: PayloadPosition0,
	}

	assert.Equal(t, ActionTypeStopRecord, action.GetActionType())
}

func TestFocusAction_GetActionType(t *testing.T) {
	action := &FocusAction{
		PayloadPositionIndex: PayloadPosition0,
		IsPointFocus:         true,
		FocusX:               0.5,
		FocusY:               0.5,
	}

	assert.Equal(t, ActionTypeFocus, action.GetActionType())
}

func TestZoomAction_GetActionType(t *testing.T) {
	action := &ZoomAction{
		PayloadPositionIndex: PayloadPosition0,
		FocalLength:          50.0,
	}

	assert.Equal(t, ActionTypeZoom, action.GetActionType())
}

func TestCustomDirNameAction_GetActionType(t *testing.T) {
	action := &CustomDirNameAction{
		PayloadPositionIndex: PayloadPosition0,
		DirectoryName:        "custom_dir",
	}

	assert.Equal(t, ActionTypeCustomDirName, action.GetActionType())
}

func TestGimbalRotateAction_GetActionType(t *testing.T) {
	action := &GimbalRotateAction{
		PayloadPositionIndex:    PayloadPosition0,
		GimbalHeadingYawBase:    "north",
		GimbalRotateMode:        "absolute",
		GimbalPitchRotateEnable: true,
		GimbalPitchRotateAngle:  45.0,
		GimbalRollRotateEnable:  false,
		GimbalRollRotateAngle:   0.0,
		GimbalYawRotateEnable:   true,
		GimbalYawRotateAngle:    90.0,
		GimbalRotateTimeEnable:  true,
		GimbalRotateTime:        5.0,
	}

	assert.Equal(t, ActionTypeGimbalRotate, action.GetActionType())
}

func TestHoverAction_GetActionType(t *testing.T) {
	action := &HoverAction{
		HoverTime: 10.0,
	}

	assert.Equal(t, ActionTypeHover, action.GetActionType())
}

func TestRotateYawAction_GetActionType(t *testing.T) {
	action := &RotateYawAction{
		AircraftHeading: 180.0,
	}

	assert.Equal(t, ActionTypeRotateYaw, action.GetActionType())
}

func TestGimbalEvenlyRotateAction_GetActionType(t *testing.T) {
	action := &GimbalEvenlyRotateAction{
		GimbalPitchRotateAngle: 30.0,
		PayloadPositionIndex:   PayloadPosition0,
	}

	assert.Equal(t, ActionTypeGimbalEvenlyRotate, action.GetActionType())
}

func TestOrientedShootAction_GetActionType(t *testing.T) {
	action := &OrientedShootAction{
		PayloadPositionIndex:      PayloadPosition0,
		GimbalPitchRotateAngle:    -30.0,
		GimbalYawRotateAngle:      0.0,
		FocusX:                    480,
		FocusY:                    360,
		FocusRegionWidth:          100,
		FocusRegionHeight:         100,
		FocalLength:               24.0,
		AircraftHeading:           90.0,
		UseGlobalPayloadLensIndex: true,
		TargetAngle:               45.0,
		ActionUUID:                "test-uuid-123",
		ImageWidth:                4000,
		ImageHeight:               3000,
		OrientedPhotoMode:         "auto",
		OrientedFilePath:          "/path/to/file.jpg",
		OrientedFileMD5:           "abc123def456",
		OrientedFileSize:          1024000,
	}

	assert.Equal(t, ActionTypeOrientedShoot, action.GetActionType())
}

func TestPanoShotAction_GetActionType(t *testing.T) {
	action := &PanoShotAction{
		PayloadPositionIndex:      PayloadPosition0,
		UseGlobalPayloadLensIndex: true,
		PanoShotSubMode:           "sphere",
	}

	assert.Equal(t, ActionTypePanoShot, action.GetActionType())
}

func TestRecordPointCloudAction_GetActionType(t *testing.T) {
	action := &RecordPointCloudAction{
		PayloadPositionIndex:    PayloadPosition0,
		RecordPointCloudOperate: "start",
	}

	assert.Equal(t, ActionTypeRecordPointCloud, action.GetActionType())
}

func TestAccurateShootAction_GetActionType(t *testing.T) {
	action := &AccurateShootAction{
		PayloadPositionIndex:      PayloadPosition0,
		GimbalPitchRotateAngle:    -30.0,
		GimbalYawRotateAngle:      0.0,
		FocusX:                    480,
		FocusY:                    360,
		FocusRegionWidth:          100,
		FocusRegionHeight:         100,
		FocalLength:               24.0,
		AircraftHeading:           90.0,
		UseGlobalPayloadLensIndex: true,
		TargetAngle:               45.0,
		ImageWidth:                4000,
		ImageHeight:               3000,
		AccurateFilePath:          "/path/to/file.jpg",
		AccurateFileMD5:           "abc123def456",
		AccurateFileSize:          1024000,
		AccurateCameraShutterTime: 1.0 / 250.0,
	}

	assert.Equal(t, ActionTypeAccurateShoot, action.GetActionType())
}

func TestCreateActionFromType(t *testing.T) {
	tests := []struct {
		name       string
		actionType string
		expected   ActionInterface
	}{
		{
			name:       "TakePhoto",
			actionType: ActionTypeTakePhoto,
			expected:   &TakePhotoAction{},
		},
		{
			name:       "StartRecord",
			actionType: ActionTypeStartRecord,
			expected:   &StartRecordAction{},
		},
		{
			name:       "StopRecord",
			actionType: ActionTypeStopRecord,
			expected:   &StopRecordAction{},
		},
		{
			name:       "Focus",
			actionType: ActionTypeFocus,
			expected:   &FocusAction{},
		},
		{
			name:       "Zoom",
			actionType: ActionTypeZoom,
			expected:   &ZoomAction{},
		},
		{
			name:       "CustomDirName",
			actionType: ActionTypeCustomDirName,
			expected:   &CustomDirNameAction{},
		},
		{
			name:       "GimbalRotate",
			actionType: ActionTypeGimbalRotate,
			expected:   &GimbalRotateAction{},
		},
		{
			name:       "Hover",
			actionType: ActionTypeHover,
			expected:   &HoverAction{},
		},
		{
			name:       "RotateYaw",
			actionType: ActionTypeRotateYaw,
			expected:   &RotateYawAction{},
		},
		{
			name:       "GimbalEvenlyRotate",
			actionType: ActionTypeGimbalEvenlyRotate,
			expected:   &GimbalEvenlyRotateAction{},
		},
		{
			name:       "OrientedShoot",
			actionType: ActionTypeOrientedShoot,
			expected:   &OrientedShootAction{},
		},
		{
			name:       "PanoShot",
			actionType: ActionTypePanoShot,
			expected:   &PanoShotAction{},
		},
		{
			name:       "AccurateShoot",
			actionType: ActionTypeAccurateShoot,
			expected:   &AccurateShootAction{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateActionFromType(tt.actionType)
			assert.IsType(t, tt.expected, result)
		})
	}
}

func TestCreateActionFromType_UnknownType(t *testing.T) {
	result := CreateActionFromType("unknown_action_type")
	assert.Nil(t, result)
}
