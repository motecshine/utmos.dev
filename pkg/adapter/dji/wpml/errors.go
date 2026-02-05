package wpml

import "errors"

const (
	ErrWaylinesValidationFailed = "waylines validation failed: %w"
	ErrConvertMissionConfig     = "failed to convert mission config: %w"
	ErrConvertTemplateFolder    = "failed to convert template folder: %w"
	ErrConvertWaylineFolder     = "failed to convert wayline folder: %w"
	ErrConvertWaypoint          = "failed to convert waypoint %d: %w"

	ErrGenerateKMZBuffer   = "failed to generate KMZ buffer: %w"
	ErrCreateDirectory     = "failed to create directory: %w"
	ErrWriteKMZFile        = "failed to write KMZ file: %w"
	ErrSerializeTemplate   = "failed to serialize template: %w"
	ErrSerializeWaylines   = "failed to serialize waylines: %w"
	ErrCreateTemplateEntry = "failed to create wpmz/template.kml entry: %w"
	ErrWriteTemplate       = "failed to write wpmz/template.kml: %w"
	ErrCreateWaylinesEntry = "failed to create wpmz/waylines.wpml entry: %w"
	ErrWriteWaylines       = "failed to write wpmz/waylines.wpml: %w"
	ErrCloseZIPWriter      = "failed to close ZIP writer: %w"
	ErrConvertWaylines     = "failed to convert waylines: %w"
	ErrParseZIP            = "failed to parse ZIP: %w"
	ErrParseZIPFile        = "failed to parse ZIP file: %w"
	ErrReadTemplate        = "failed to read template.kml: %w"
	ErrReadWaylinesWPML    = "failed to read waylines.wpml: %w"
	ErrParseTemplateKML    = "failed to parse template.kml: %w"
	ErrParseWaylinesWPML   = "failed to parse waylines.wpml: %w"

	ErrMarshalDocument           = "failed to marshal document: %w"
	ErrUnmarshalDocument         = "failed to unmarshal document: %w"
	ErrReadData                  = "failed to read data: %w"
	ErrParseXML                  = "failed to parse XML: %w"
	ErrEncodeToken               = "failed to encode token: %w"
	ErrFlushEncoder              = "failed to flush encoder: %w"
	ErrInvalidXML                = "invalid XML: %w"
	ErrMarshalTemplateDocument   = "failed to marshal template document: %w"
	ErrMarshalWaylinesDocument   = "failed to marshal waylines document: %w"
	ErrUnmarshalTemplateDocument = "failed to unmarshal template document: %w"
	ErrUnmarshalWaylinesDocument = "failed to unmarshal waylines document: %w"

	ErrUnknownActionType                = "unknown action type: %s"
	ErrTypeMismatch                     = "type mismatch: declared=%s, actual=%s"
	ErrUnmarshalActionRequest           = "failed to unmarshal action request: %w"
	ErrInvalidActionRequest             = "invalid action request: %w"
	ErrActionTypeMismatch               = "action type mismatch: expected %s, got %s"
	ErrActionGroupValidationFailed      = "action group validation failed: %w"
	ErrActionValidationFailed           = "action[%d] validation failed: %w"
	ErrWaylinesDocumentValidationFailed = "waylines document validation failed: %w"
	ErrExpectedStruct                   = "expected struct, got %s"
	ErrFieldValidationFailed            = "field %s validation failed: %w"
	ErrFieldRequiredForDroneModel       = "field %s is required for drone model %d"
	ErrFieldRequiredForPayloadModel     = "field %s is required for payload model %d"
)

var (
	ErrMissionCannotBeEmpty         = errors.New("mission cannot be empty")
	ErrMissionTemplateCannotBeEmpty = errors.New("mission.Template cannot be empty")
	ErrMissionWaylinesCannotBeEmpty = errors.New("mission.Waylines cannot be empty")
	ErrKMZFormatIncorrect           = errors.New("KMZ file format is incorrect: missing required template or waylines file")
	ErrActionIsNil                  = errors.New("action is nil")
	ErrActionCannotBeNil            = errors.New("action cannot be nil")
	ErrActionGroupCannotBeNil       = errors.New("action group cannot be nil")
	ErrWaylineDocumentCannotBeNil   = errors.New("wayline document cannot be nil")
	ErrTemplateCannotBeNil          = errors.New("template cannot be nil")
)
