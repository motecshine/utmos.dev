package wpml

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// WPMLValidator provides validation for WPML structures using custom validation rules.
type WPMLValidator struct {
	validator *validator.Validate
}

// NewWPMLValidator creates a new WPMLValidator with all custom validation rules registered.
func NewWPMLValidator() (*WPMLValidator, error) {
	validate := validator.New()

	w := &WPMLValidator{validator: validate}
	if err := w.registerCustomValidators(); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WPMLValidator) registerCustomValidators() error {
	validations := []struct {
		tag string
		fn  validator.Func
	}{
		{"payload_position", w.validatePayloadPosition},
		{"drone_model", w.validateDroneModel},
		{"payload_model", w.validatePayloadModel},
		{"action_type", w.validateActionType},
		{"required_for_drone", w.validateRequiredForDrone},
		{"required_for_payload", w.validateRequiredForPayload},
	}

	for _, v := range validations {
		if err := w.validator.RegisterValidation(v.tag, v.fn); err != nil {
			return fmt.Errorf("failed to register validation %s: %w", v.tag, err)
		}
	}
	return nil
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

func (w *WPMLValidator) validateRequiredForDrone(_ validator.FieldLevel) bool {
	return true
}

func (w *WPMLValidator) validateRequiredForPayload(_ validator.FieldLevel) bool {
	return true
}

// ValidateStruct validates all fields of the given struct according to their validation tags.
func (w *WPMLValidator) ValidateStruct(s any) error {
	return w.validator.Struct(s)
}

// ValidateVar validates a single variable against the given validation tag.
func (w *WPMLValidator) ValidateVar(field any, tag string) error {
	return w.validator.Var(field, tag)
}

// ValidateAction validates an action, returning an error if it is nil or fails struct validation.
func (w *WPMLValidator) ValidateAction(action any) error {
	if action == nil {
		return ErrActionCannotBeNil
	}

	return w.ValidateStruct(action)
}

// ValidateActionGroup validates an action group and all of its contained actions.
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

// ValidateWaylinesDocument validates a WaylinesDocument, returning an error if it is nil or fails validation.
func (w *WPMLValidator) ValidateWaylinesDocument(waylineDoc *WaylinesDocument) error {
	if waylineDoc == nil {
		return ErrWaylineDocumentCannotBeNil
	}

	if err := w.ValidateStruct(waylineDoc); err != nil {
		return fmt.Errorf(ErrWaylinesDocumentValidationFailed, err)
	}

	return nil
}

// ValidateTemplateDocument validates a TemplateDocument, returning an error if it is nil or fails validation.
func (w *WPMLValidator) ValidateTemplateDocument(template *TemplateDocument) error {
	if template == nil {
		return ErrTemplateCannotBeNil
	}

	return w.ValidateStruct(template)
}

// ValidateWithContext validates a struct with additional drone and payload model context for conditional rules.
func (w *WPMLValidator) ValidateWithContext(s any, droneModel DroneModel, payloadModel PayloadModel) error {
	if err := w.ValidateStruct(s); err != nil {
		return err
	}

	return w.validateWithMachineContext(s, droneModel, payloadModel)
}

func (w *WPMLValidator) validateWithMachineContext(s any, droneModel DroneModel, payloadModel PayloadModel) error {
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

func (w *WPMLValidator) checkRequiredFor(rule, prefix string, field reflect.Value, fieldName string, matchFn func(string) bool, errFmt string, modelValue any) error {
	if strings.HasPrefix(rule, prefix) {
		pattern := strings.TrimPrefix(rule, prefix)
		if matchFn(pattern) && field.IsZero() {
			return fmt.Errorf(errFmt, fieldName, modelValue)
		}
	}
	return nil
}

// two checkRequiredFor calls differ only in parameters; further extraction would hurt readability
func (w *WPMLValidator) validateFieldWithMachineContext(field reflect.Value, fieldType reflect.StructField, validateTag string, droneModel DroneModel, payloadModel PayloadModel) error {
	rules := strings.Split(validateTag, ",")

	for _, rule := range rules {
		if err := w.checkRequiredFor(rule, "required_for_drone:", field, fieldType.Name,
			func(pattern string) bool { return w.isRequiredForDrone(pattern, droneModel) },
			ErrFieldRequiredForDroneModel, droneModel); err != nil {
			return err
		}

		if err := w.checkRequiredFor(rule, "required_for_payload:", field, fieldType.Name,
			func(pattern string) bool { return w.isRequiredForPayload(pattern, payloadModel) },
			ErrFieldRequiredForPayloadModel, payloadModel); err != nil {
			return err
		}
	}

	return nil
}

func matchesAnyPattern[T any](pattern string, model T, matcher func(string, T) bool) bool {
	for _, p := range strings.Split(pattern, "|") {
		if matcher(p, model) {
			return true
		}
	}
	return false
}

func (w *WPMLValidator) isRequiredForDrone(dronePattern string, droneModel DroneModel) bool {
	return matchesAnyPattern(dronePattern, droneModel, w.matchesDronePattern)
}

func (w *WPMLValidator) isRequiredForPayload(payloadPattern string, payloadModel PayloadModel) bool {
	return matchesAnyPattern(payloadPattern, payloadModel, w.matchesPayloadPattern)
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

// GetValidationErrors extracts and formats validation errors into a list of human-readable strings.
func (w *WPMLValidator) GetValidationErrors(err error) []string {
	var errs []string

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			errorMsg := w.formatValidationError(e)
			errs = append(errs, errorMsg)
		}
	} else {
		errs = append(errs, err.Error())
	}

	return errs
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

// InitGlobalValidator initializes the package-level global WPML validator.
func InitGlobalValidator() error {
	var err error
	globalValidator, err = NewWPMLValidator()
	return err
}

// MustInitGlobalValidator initializes the global validator and panics on error.
// This should only be used during program initialization (e.g., in init() or main()).
// For runtime initialization with error handling, use InitGlobalValidator().
func MustInitGlobalValidator() {
	if err := InitGlobalValidator(); err != nil {
		panic(fmt.Errorf("MustInitGlobalValidator: %w", err))
	}
}

// Validate validates the given struct using the global WPML validator, initializing it if necessary.
func Validate(s any) error {
	if globalValidator == nil {
		if err := InitGlobalValidator(); err != nil {
			return err
		}
	}
	return globalValidator.ValidateStruct(s)
}

// ValidateActionGlobal validates an action using the global WPML validator, initializing it if necessary.
func ValidateActionGlobal(action any) error {
	if globalValidator == nil {
		if err := InitGlobalValidator(); err != nil {
			return err
		}
	}
	return globalValidator.ValidateAction(action)
}

// ValidateWaylinesDocumentGlobal validates a WaylinesDocument using the global WPML validator, initializing it if necessary.
func ValidateWaylinesDocumentGlobal(waylineDoc *WaylinesDocument) error {
	if globalValidator == nil {
		if err := InitGlobalValidator(); err != nil {
			return err
		}
	}
	return globalValidator.ValidateWaylinesDocument(waylineDoc)
}
