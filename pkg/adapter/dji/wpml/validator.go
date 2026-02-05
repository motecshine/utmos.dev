package wpml

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type WPMLValidator struct {
	validator *validator.Validate
}

func NewWPMLValidator() *WPMLValidator {
	validate := validator.New()

	w := &WPMLValidator{validator: validate}
	w.registerCustomValidators()

	return w
}

func (w *WPMLValidator) registerCustomValidators() {
	w.validator.RegisterValidation("payload_position", w.validatePayloadPosition)
	w.validator.RegisterValidation("drone_model", w.validateDroneModel)
	w.validator.RegisterValidation("payload_model", w.validatePayloadModel)
	w.validator.RegisterValidation("action_type", w.validateActionType)
	w.validator.RegisterValidation("required_for_drone", w.validateRequiredForDrone)
	w.validator.RegisterValidation("required_for_payload", w.validateRequiredForPayload)
}

func (w *WPMLValidator) validatePayloadPosition(fl validator.FieldLevel) bool {
	value := fl.Field().Int()

	return value == 0 || value == 1 || value == 2 || value == 7
}

func (w *WPMLValidator) validateDroneModel(fl validator.FieldLevel) bool {
	value := fl.Field().Int()

	validDroneModels := []int{
		int(DroneM300RTK),
		int(DroneM350RTK),
		int(DroneM30),
		int(DroneM3Series),
		int(DroneM3DSeries),
		int(DroneM4Series),
		int(DroneM4DSeries),
		int(DroneM400),
		int(DroneDahua),
	}

	for _, validModel := range validDroneModels {
		if int(value) == validModel {
			return true
		}
	}
	return false
}

func (w *WPMLValidator) validatePayloadModel(fl validator.FieldLevel) bool {
	value := fl.Field().Int()

	validPayloadModels := []int{
		int(PayloadZ30),
		int(PayloadXT2),
		int(PayloadXTS),
		int(PayloadH20),
		int(PayloadH20T),
		int(PayloadH20N),
		int(PayloadH30),
		int(PayloadH30T),
		int(PayloadM30Camera),
		int(PayloadM30TCamera),
		int(PayloadMavic3ECamera),
		int(PayloadMavic3TCamera),
		int(PayloadMavic3ACamera),
		int(PayloadMatrice3DCamera),
		int(PayloadMatrice3TDCamera),
		int(PayloadMatrice4ECamera),
		int(PayloadMatrice4TCamera),
		int(PayloadMatrice4DCamera),
		int(PayloadMatrice4TDCamera),
		int(PayloadFPVCamera),
		int(PayloadDockCamera),
		int(PayloadAuxiliaryCamera),
		int(PayloadPSDK),
		int(PayloadX900),
	}

	for _, validModel := range validPayloadModels {
		if int(value) == validModel {
			return true
		}
	}
	return false
}

func (w *WPMLValidator) validateActionType(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	validActionTypes := []string{
		ActionTypeTakePhoto,
		ActionTypeStartRecord,
		ActionTypeStopRecord,
		ActionTypeFocus,
		ActionTypeZoom,
		ActionTypeCustomDirName,
		ActionTypeGimbalRotate,
		ActionTypeRotateYaw,
		ActionTypeHover,
		ActionTypeGimbalEvenlyRotate,
		ActionTypeOrientedShoot,
		ActionTypePanoShot,
		ActionTypeRecordPointCloud,
		ActionTypeAccurateShoot,
		ActionTypeGimbalAngleLock,
		ActionTypeGimbalAngleUnlock,
		ActionTypeStartSmartOblique,
		ActionTypeStartTimeLapse,
		ActionTypeStopTimeLapse,
		ActionTypeSetFocusType,
		ActionTypeTargetDetection,
	}

	for _, validType := range validActionTypes {
		if value == validType {
			return true
		}
	}
	return false
}

func (w *WPMLValidator) validateRequiredForDrone(fl validator.FieldLevel) bool {

	return true
}

func (w *WPMLValidator) validateRequiredForPayload(fl validator.FieldLevel) bool {

	return true
}

func (w *WPMLValidator) ValidateStruct(s interface{}) error {
	return w.validator.Struct(s)
}

func (w *WPMLValidator) ValidateVar(field interface{}, tag string) error {
	return w.validator.Var(field, tag)
}

func (w *WPMLValidator) ValidateAction(action interface{}) error {
	if action == nil {
		return ErrActionCannotBeNil
	}

	return w.ValidateStruct(action)
}

func (w *WPMLValidator) ValidateActionGroup(actionGroup *ActionGroup) error {
	if actionGroup == nil {
		return ErrActionGroupCannotBeNil
	}

	if err := w.ValidateStruct(actionGroup); err != nil {
		return fmt.Errorf(ErrActionGroupValidationFailed, err)
	}

	for i, action := range actionGroup.Actions {
		if err := w.ValidateAction(action); err != nil {
			return fmt.Errorf(ErrActionValidationFailed, i, err)
		}
	}

	return nil
}

func (w *WPMLValidator) ValidateWaylinesDocument(waylineDoc *WaylinesDocument) error {
	if waylineDoc == nil {
		return ErrWaylineDocumentCannotBeNil
	}

	if err := w.ValidateStruct(waylineDoc); err != nil {
		return fmt.Errorf(ErrWaylinesDocumentValidationFailed, err)
	}

	return nil
}

func (w *WPMLValidator) ValidateTemplateDocument(template *TemplateDocument) error {
	if template == nil {
		return ErrTemplateCannotBeNil
	}

	return w.ValidateStruct(template)
}

func (w *WPMLValidator) ValidateWithContext(s interface{}, droneModel DroneModel, payloadModel PayloadModel) error {

	if err := w.ValidateStruct(s); err != nil {
		return err
	}

	return w.validateWithMachineContext(s, droneModel, payloadModel)
}

func (w *WPMLValidator) validateWithMachineContext(s interface{}, droneModel DroneModel, payloadModel PayloadModel) error {

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return fmt.Errorf(ErrExpectedStruct, val.Kind())
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if validateTag := fieldType.Tag.Get("validate"); validateTag != "" {
			if err := w.validateFieldWithMachineContext(field, fieldType, validateTag, droneModel, payloadModel); err != nil {
				return fmt.Errorf(ErrFieldValidationFailed, fieldType.Name, err)
			}
		}
	}

	return nil
}

func (w *WPMLValidator) validateFieldWithMachineContext(field reflect.Value, fieldType reflect.StructField, validateTag string, droneModel DroneModel, payloadModel PayloadModel) error {

	rules := strings.Split(validateTag, ",")

	for _, rule := range rules {
		if strings.HasPrefix(rule, "required_for_drone:") {

			requiredForDrone := strings.TrimPrefix(rule, "required_for_drone:")
			if w.isRequiredForDrone(requiredForDrone, droneModel) {
				if field.IsZero() {
					return fmt.Errorf(ErrFieldRequiredForDroneModel, fieldType.Name, droneModel)
				}
			}
		}

		if strings.HasPrefix(rule, "required_for_payload:") {

			requiredForPayload := strings.TrimPrefix(rule, "required_for_payload:")
			if w.isRequiredForPayload(requiredForPayload, payloadModel) {
				if field.IsZero() {
					return fmt.Errorf(ErrFieldRequiredForPayloadModel, fieldType.Name, payloadModel)
				}
			}
		}
	}

	return nil
}

func (w *WPMLValidator) isRequiredForDrone(dronePattern string, droneModel DroneModel) bool {

	requiredDrones := strings.Split(dronePattern, "|")
	for _, pattern := range requiredDrones {
		if w.matchesDronePattern(pattern, droneModel) {
			return true
		}
	}
	return false
}

func (w *WPMLValidator) isRequiredForPayload(payloadPattern string, payloadModel PayloadModel) bool {

	requiredPayloads := strings.Split(payloadPattern, "|")
	for _, pattern := range requiredPayloads {
		if w.matchesPayloadPattern(pattern, payloadModel) {
			return true
		}
	}
	return false
}

func (w *WPMLValidator) matchesDronePattern(pattern string, droneModel DroneModel) bool {
	pattern = strings.ToUpper(strings.TrimSpace(pattern))

	switch pattern {
	case "M300":
		return droneModel == DroneM300RTK
	case "M350":
		return droneModel == DroneM350RTK
	case "M30":
		return droneModel == DroneM30
	case "M3":
		return droneModel == DroneM3Series
	case "M3D":
		return droneModel == DroneM3DSeries
	case "M4":
		return droneModel == DroneM4Series
	case "M4D":
		return droneModel == DroneM4DSeries
	case "M400":
		return droneModel == DroneM400
	default:
		return false
	}
}

func (w *WPMLValidator) matchesPayloadPattern(pattern string, payloadModel PayloadModel) bool {
	pattern = strings.ToUpper(strings.TrimSpace(pattern))

	switch pattern {
	case "H20":
		return payloadModel == PayloadH20 || payloadModel == PayloadH20T || payloadModel == PayloadH20N
	case "H30":
		return payloadModel == PayloadH30 || payloadModel == PayloadH30T
	case "M30":
		return payloadModel == PayloadM30Camera || payloadModel == PayloadM30TCamera
	case "M3":
		return payloadModel == PayloadMavic3ECamera || payloadModel == PayloadMavic3TCamera || payloadModel == PayloadMavic3ACamera
	case "M3D":
		return payloadModel == PayloadMatrice3DCamera || payloadModel == PayloadMatrice3TDCamera
	case "M4":
		return payloadModel == PayloadMatrice4ECamera || payloadModel == PayloadMatrice4TCamera
	case "M4D":
		return payloadModel == PayloadMatrice4DCamera || payloadModel == PayloadMatrice4TDCamera
	case "XT":
		return payloadModel == PayloadXT2 || payloadModel == PayloadXTS
	case "Z30":
		return payloadModel == PayloadZ30
	case "FPV":
		return payloadModel == PayloadFPVCamera
	case "DOCK":
		return payloadModel == PayloadDockCamera
	case "AUX":
		return payloadModel == PayloadAuxiliaryCamera
	case "PSDK":
		return payloadModel == PayloadPSDK
	default:
		return false
	}
}

func (w *WPMLValidator) GetValidationErrors(err error) []string {
	var errors []string

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			errorMsg := w.formatValidationError(e)
			errors = append(errors, errorMsg)
		}
	} else {
		errors = append(errors, err.Error())
	}

	return errors
}

func (w *WPMLValidator) formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("field '%s' 是required", e.Field())
	case "min":
		return fmt.Sprintf("field '%s' 值必须大于等于 %s", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("field '%s' 值必须小于等于 %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("field '%s' 值必须大于等于 %s", e.Field(), e.Param())
	case "lte":
		return fmt.Sprintf("field '%s' 值必须小于等于 %s", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("field '%s' 值必须大于 %s", e.Field(), e.Param())
	case "lt":
		return fmt.Sprintf("field '%s' 值必须小于 %s", e.Field(), e.Param())
	case "payload_position":
		return fmt.Sprintf("field '%s' 负载position值invalid，valid值为：0(main gimbal), 1, 2, 7", e.Field())
	case "drone_model":
		return fmt.Sprintf("field '%s' Drone modelinvalid", e.Field())
	case "payload_model":
		return fmt.Sprintf("field '%s' 负载modelinvalid", e.Field())
	default:
		return fmt.Sprintf("field '%s' validationfailure: %s", e.Field(), e.Tag())
	}
}

var globalValidator *WPMLValidator

func InitGlobalValidator() {
	globalValidator = NewWPMLValidator()
}

func Validate(s interface{}) error {
	if globalValidator == nil {
		InitGlobalValidator()
	}
	return globalValidator.ValidateStruct(s)
}

func ValidateActionGlobal(action interface{}) error {
	if globalValidator == nil {
		InitGlobalValidator()
	}
	return globalValidator.ValidateAction(action)
}

func ValidateWaylinesDocumentGlobal(waylineDoc *WaylinesDocument) error {
	if globalValidator == nil {
		InitGlobalValidator()
	}
	return globalValidator.ValidateWaylinesDocument(waylineDoc)
}
