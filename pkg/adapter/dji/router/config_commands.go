package router

import (
	"encoding/json"
)

// Config command method names.
const (
	MethodConfigGet = "config_get"
	MethodConfigSet = "config_set"
)

// ConfigGetData represents config get request data.
// Note: This type is specific to router as it's not defined in protocol packages.
type ConfigGetData struct {
	ConfigType string `json:"config_type"`
}

// ConfigSetData represents config set request data.
// Note: This type is specific to router as it's not defined in protocol packages.
type ConfigSetData struct {
	ConfigType  string          `json:"config_type"`
	ConfigValue json.RawMessage `json:"config_value"`
}

// RegisterConfigCommands registers all config commands to the router.
// Returns an error if any handler registration fails.
func RegisterConfigCommands(r *ServiceRouter) error {
	handlers := map[string]ServiceHandlerFunc{
		MethodConfigGet: SimpleCommandHandler[ConfigGetData](MethodConfigGet),
		MethodConfigSet: SimpleCommandHandler[ConfigSetData](MethodConfigSet),
	}

	return RegisterHandlers(r, handlers)
}
