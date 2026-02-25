package wpml

import "errors"

// Error format string constants for wrapping errors with contextual messages.
const (
	// ErrWaylinesValidationFailed is the error format for waylines validation failures.
	ErrWaylinesValidationFailed = "waylines validation failed: %w"
	// ErrConvertMissionConfig is the error format for mission config conversion failures.
	ErrConvertMissionConfig = "failed to convert mission config: %w"
	// ErrConvertTemplateFolder is the error format for template folder conversion failures.
	ErrConvertTemplateFolder = "failed to convert template folder: %w"
	// ErrConvertWaylineFolder is the error format for wayline folder conversion failures.
	ErrConvertWaylineFolder = "failed to convert wayline folder: %w"
	// ErrConvertWaypoint is the error format for waypoint conversion failures.
	ErrConvertWaypoint = "failed to convert waypoint %d: %w"

	// ErrGenerateKMZBuffer is the error format for KMZ buffer generation failures.
	ErrGenerateKMZBuffer = "failed to generate KMZ buffer: %w"
	// ErrCreateDirectory is the error format for directory creation failures.
	ErrCreateDirectory = "failed to create directory: %w"
	// ErrWriteKMZFile is the error format for KMZ file write failures.
	ErrWriteKMZFile = "failed to write KMZ file: %w"
	// ErrSerializeTemplate is the error format for template serialization failures.
	ErrSerializeTemplate = "failed to serialize template: %w"
	// ErrSerializeWaylines is the error format for waylines serialization failures.
	ErrSerializeWaylines = "failed to serialize waylines: %w"
	// ErrCreateTemplateEntry is the error format for template ZIP entry creation failures.
	ErrCreateTemplateEntry = "failed to create wpmz/template.kml entry: %w"
	// ErrWriteTemplate is the error format for template write failures.
	ErrWriteTemplate = "failed to write wpmz/template.kml: %w"
	// ErrCreateWaylinesEntry is the error format for waylines ZIP entry creation failures.
	ErrCreateWaylinesEntry = "failed to create wpmz/waylines.wpml entry: %w"
	// ErrWriteWaylines is the error format for waylines write failures.
	ErrWriteWaylines = "failed to write wpmz/waylines.wpml: %w"
	// ErrCloseZIPWriter is the error format for ZIP writer close failures.
	ErrCloseZIPWriter = "failed to close ZIP writer: %w"
	// ErrConvertWaylines is the error format for waylines conversion failures.
	ErrConvertWaylines = "failed to convert waylines: %w"
	// ErrParseZIP is the error format for ZIP parsing failures.
	ErrParseZIP = "failed to parse ZIP: %w"
	// ErrParseZIPFile is the error format for ZIP file parsing failures.
	ErrParseZIPFile = "failed to parse ZIP file: %w"
	// ErrReadTemplate is the error format for template.kml read failures.
	ErrReadTemplate = "failed to read template.kml: %w"
	// ErrReadWaylinesWPML is the error format for waylines.wpml read failures.
	ErrReadWaylinesWPML = "failed to read waylines.wpml: %w"
	// ErrParseTemplateKML is the error format for template.kml parse failures.
	ErrParseTemplateKML = "failed to parse template.kml: %w"
	// ErrParseWaylinesWPML is the error format for waylines.wpml parse failures.
	ErrParseWaylinesWPML = "failed to parse waylines.wpml: %w"

	// ErrMarshalDocument is the error format for document marshaling failures.
	ErrMarshalDocument = "failed to marshal document: %w"
	// ErrUnmarshalDocument is the error format for document unmarshaling failures.
	ErrUnmarshalDocument = "failed to unmarshal document: %w"
	// ErrReadData is the error format for data read failures.
	ErrReadData = "failed to read data: %w"
	// ErrParseXML is the error format for XML parsing failures.
	ErrParseXML = "failed to parse XML: %w"
	// ErrEncodeToken is the error format for XML token encoding failures.
	ErrEncodeToken = "failed to encode token: %w"
	// ErrFlushEncoder is the error format for XML encoder flush failures.
	ErrFlushEncoder = "failed to flush encoder: %w"
	// ErrInvalidXML is the error format for invalid XML detection.
	ErrInvalidXML = "invalid XML: %w"
	// ErrMarshalTemplateDocument is the error format for template document marshaling failures.
	ErrMarshalTemplateDocument = "failed to marshal template document: %w"
	// ErrMarshalWaylinesDocument is the error format for waylines document marshaling failures.
	ErrMarshalWaylinesDocument = "failed to marshal waylines document: %w"
	// ErrUnmarshalTemplateDocument is the error format for template document unmarshaling failures.
	ErrUnmarshalTemplateDocument = "failed to unmarshal template document: %w"
	// ErrUnmarshalWaylinesDocument is the error format for waylines document unmarshaling failures.
	ErrUnmarshalWaylinesDocument = "failed to unmarshal waylines document: %w"

	// ErrUnknownActionType is the error format for unknown action types.
	ErrUnknownActionType = "unknown action type: %s"
	// ErrTypeMismatch is the error format for action type mismatches.
	ErrTypeMismatch = "type mismatch: declared=%s, actual=%s"
	// ErrUnmarshalActionRequest is the error format for action request unmarshaling failures.
	ErrUnmarshalActionRequest = "failed to unmarshal action request: %w"
	// ErrInvalidActionRequest is the error format for invalid action request validation failures.
	ErrInvalidActionRequest = "invalid action request: %w"
	// ErrActionTypeMismatch is the error format for concrete action type mismatches.
	ErrActionTypeMismatch = "action type mismatch: expected %s, got %s"
	// ErrActionGroupValidationFailed is the error format for action group validation failures.
	ErrActionGroupValidationFailed = "action group validation failed: %w"
	// ErrActionValidationFailed is the error format for individual action validation failures.
	ErrActionValidationFailed = "action[%d] validation failed: %w"
	// ErrWaylinesDocumentValidationFailed is the error format for waylines document validation failures.
	ErrWaylinesDocumentValidationFailed = "waylines document validation failed: %w"
	// ErrExpectedStruct is the error format for non-struct type validation errors.
	ErrExpectedStruct = "expected struct, got %s"
	// ErrFieldValidationFailed is the error format for individual field validation failures.
	ErrFieldValidationFailed = "field %s validation failed: %w"
	// ErrFieldRequiredForDroneModel is the error format for fields required by specific drone models.
	ErrFieldRequiredForDroneModel = "field %s is required for drone model %d"
	// ErrFieldRequiredForPayloadModel is the error format for fields required by specific payload models.
	ErrFieldRequiredForPayloadModel = "field %s is required for payload model %d"
)

// Sentinel error variables for common validation and structural error conditions.
var (
	// ErrMissionCannotBeEmpty is returned when a nil mission is provided.
	ErrMissionCannotBeEmpty = errors.New("mission cannot be empty")
	// ErrMissionTemplateCannotBeEmpty is returned when the mission template is nil.
	ErrMissionTemplateCannotBeEmpty = errors.New("mission.Template cannot be empty")
	// ErrMissionWaylinesCannotBeEmpty is returned when the mission waylines document is nil.
	ErrMissionWaylinesCannotBeEmpty = errors.New("mission.Waylines cannot be empty")
	// ErrKMZFormatIncorrect is returned when a KMZ file is missing required template or waylines files.
	ErrKMZFormatIncorrect = errors.New("KMZ file format is incorrect: missing required template or waylines file")
	// ErrActionIsNil is returned when an action is nil.
	ErrActionIsNil = errors.New("action is nil")
	// ErrActionCannotBeNil is returned when a nil action is provided to validation.
	ErrActionCannotBeNil = errors.New("action cannot be nil")
	// ErrActionGroupCannotBeNil is returned when a nil action group is provided to validation.
	ErrActionGroupCannotBeNil = errors.New("action group cannot be nil")
	// ErrWaylineDocumentCannotBeNil is returned when a nil wayline document is provided to validation.
	ErrWaylineDocumentCannotBeNil = errors.New("wayline document cannot be nil")
	// ErrTemplateCannotBeNil is returned when a nil template document is provided to validation.
	ErrTemplateCannotBeNil = errors.New("template cannot be nil")
)
