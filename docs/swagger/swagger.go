package swagger

import (
	_ "embed"

	"github.com/swaggo/swag"
)

//go:embed openapi.json
var openAPIDoc string

type openAPISpec struct{}

// ReadDoc returns the registered OpenAPI document for Swagger UI.
func (openAPISpec) ReadDoc() string {
	return openAPIDoc
}

func init() {
	swag.Register(swag.Name, openAPISpec{})
}
