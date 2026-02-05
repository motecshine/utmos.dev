package wpml

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type ActionRequest struct {
	Type   string          `json:"type" validate:"required,action_type"`
	Action ActionInterface `json:"action"`
}

func (ar *ActionRequest) UnmarshalJSON(data []byte) error {

	var temp struct {
		Type   string          `json:"type"`
		Action json.RawMessage `json:"action"`
		Params json.RawMessage `json:"params"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	ar.Type = temp.Type

	action, err := createActionByType(temp.Type)
	if err != nil {
		return err
	}

	var actionData json.RawMessage
	if len(temp.Params) > 0 {
		actionData = temp.Params
	} else {
		actionData = temp.Action
	}

	if err := json.Unmarshal(actionData, action); err != nil {
		return err
	}

	ar.Action = action
	return nil
}

func (ar *ActionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string          `json:"type"`
		Action ActionInterface `json:"action"`
	}{
		Type:   ar.Type,
		Action: ar.Action,
	})
}

func createActionByType(actionType string) (ActionInterface, error) {
	switch actionType {
	case ActionTypeTakePhoto:
		return &TakePhotoAction{}, nil
	case ActionTypeStartRecord:
		return &StartRecordAction{}, nil
	case ActionTypeStopRecord:
		return &StopRecordAction{}, nil
	case ActionTypeFocus:
		return &FocusAction{}, nil
	case ActionTypeZoom:
		return &ZoomAction{}, nil
	case ActionTypeCustomDirName:
		return &CustomDirNameAction{}, nil
	case ActionTypeGimbalRotate:
		return &GimbalRotateAction{}, nil
	case ActionTypeRotateYaw:
		return &RotateYawAction{}, nil
	case ActionTypeHover:
		return &HoverAction{}, nil
	case ActionTypeGimbalEvenlyRotate:
		return &GimbalEvenlyRotateAction{}, nil
	case ActionTypeOrientedShoot:
		return &OrientedShootAction{}, nil
	case ActionTypePanoShot:
		return &PanoShotAction{}, nil
	case ActionTypeRecordPointCloud:
		return &RecordPointCloudAction{}, nil
	case ActionTypeAccurateShoot:
		return &AccurateShootAction{}, nil
	case ActionTypeGimbalAngleLock:
		return &GimbalAngleLockAction{}, nil
	case ActionTypeGimbalAngleUnlock:
		return &GimbalAngleUnlockAction{}, nil
	case ActionTypeStartSmartOblique:
		return &StartSmartObliqueAction{}, nil
	case ActionTypeStartTimeLapse:
		return &StartTimeLapseAction{}, nil
	case ActionTypeStopTimeLapse:
		return &StopTimeLapseAction{}, nil
	case ActionTypeSetFocusType:
		return &SetFocusTypeAction{}, nil
	case ActionTypeTargetDetection:
		return &TargetDetectionAction{}, nil
	default:
		return nil, fmt.Errorf(ErrUnknownActionType, actionType)
	}
}

func (ar *ActionRequest) GetActionType() string {
	if ar.Action != nil {
		return ar.Action.GetActionType()
	}
	return ar.Type
}

func (ar *ActionRequest) Validate() error {
	if ar.Action == nil {
		return ErrActionIsNil
	}

	expectedType := ar.Action.GetActionType()
	if ar.Type != expectedType {
		return fmt.Errorf(ErrTypeMismatch, ar.Type, expectedType)
	}

	return nil
}

func NewActionRequest(action ActionInterface) *ActionRequest {
	return &ActionRequest{
		Type:   action.GetActionType(),
		Action: action,
	}
}

func TypedActionRequest[T ActionInterface](action T) *ActionRequest {
	return &ActionRequest{
		Type:   action.GetActionType(),
		Action: action,
	}
}

func ActionRequestFromJSON(jsonData []byte) (*ActionRequest, error) {
	var actionRequest ActionRequest
	if err := json.Unmarshal(jsonData, &actionRequest); err != nil {
		return nil, fmt.Errorf(ErrUnmarshalActionRequest, err)
	}

	if err := actionRequest.Validate(); err != nil {
		return nil, fmt.Errorf(ErrInvalidActionRequest, err)
	}

	return &actionRequest, nil
}

func (ar *ActionRequest) ToJSON() ([]byte, error) {
	return json.Marshal(ar)
}

func (ar *ActionRequest) GetConcreteAction() interface{} {
	return ar.Action
}

func GetTypedAction[T ActionInterface](ar *ActionRequest) (T, error) {
	var zero T
	if ar.Action == nil {
		return zero, ErrActionIsNil
	}

	typedAction, ok := ar.Action.(T)
	if !ok {
		actualType := reflect.TypeOf(ar.Action).Elem().Name()
		expectedType := reflect.TypeOf(zero).Elem().Name()
		return zero, fmt.Errorf(ErrActionTypeMismatch, expectedType, actualType)
	}

	return typedAction, nil
}
