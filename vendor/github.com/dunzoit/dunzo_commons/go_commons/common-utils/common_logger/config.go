package common_logger

import "net/http"

type SourceType string

// Not supporting BODY as a SourceType  since it would
// require reflection on the request body
const (
	HEADER     SourceType = "HEADER"
	QUERYPARAM SourceType = "QUERYPARAM"
)

// A ContextFieldSource specifies the location (header or queryParam)
// and key to lookup by for the value to be passed into the logging context
type ContextFieldSource struct {
	Type SourceType `json:"type" mapstructure:"type"` // Location of context value, supports HEADER or QUERYPARAM for now
	Key  string     `json:"key" mapstructure:"key"`
}

// A ContextField denotes a source (where we extract from) as explained above and the field (or key)
// by which it will be put into the context
type ContextField struct {
	Field  string             `json:"field" mapstructure:"field"`
	Source ContextFieldSource `json:"source" mapstructure:"source"`
}

// An ApiContextField provides the path and method of the API
// and the corresponding contexts field to be extracted in it
type ApiContextField struct {
	Fields []ContextField `json:"fields" mapstructure:"fields"`
	Path   string         `json:"path" mapstructure:"path"`
	Method string         `json:"method" mapstructure:"method"`
}

// The LogConfig is the global logging config to be provided by a web service.
// - It specifies the key,values to be extracted and put into the logging context
//   at a global and per API level.
// - All fields, including nested ones, must contain `mapstructure` tags so that
//   `viper.UnmarshalKey` can correctly unmarshal from config
type LogConfig struct {
	Global []ContextField    `json:"global" mapstructure:"global"`
	Apis   []ApiContextField `json:"apis" mapstructure:"apis"`
	Level  int8              `json:"level" mapstructure:"level"`
}

// ExtractContextFieldForSource extracts values from the request provided the source
// which can be a key in either the request headers or query params
func ExtractContextFieldForSource(request *http.Request, source ContextFieldSource) string {
	switch source.Type {
	case HEADER:
		return request.Header.Get(source.Key)
	case QUERYPARAM:
		return request.URL.Query().Get(source.Key)
	default:
		return ""
	}
}

// MatchApiToExtractConfig exactly matches the METHOD & PATH (without host) of the API
// with the provided config
func MatchApiToExtractConfig(method string, path string, apiConfigs []ApiContextField) *ApiContextField {
	for _, apiConfig := range apiConfigs {
		// Method and path should exactly match : not implementing regex here since that would be costly
		if apiConfig.Method == method && apiConfig.Path == path {
			return &apiConfig
		}
	}
	//If no api is matched return nil
	return nil
}

// ExtractContextFieldFromRequest is a helper function provided to extract the key, value pairs
// from the request as specified by the config maintained globally in the packages
func ExtractContextFieldFromRequest(request *http.Request) map[string]string {
	contextfields := make(map[string]string)

	if logConfig == nil {
		log.Warn("SetUpLogging hasn't been called yet, so no context was added to the logger")
		return contextfields
	}

	//Extract globally configured fields
	for _, fieldSource := range logConfig.Global {
		contextfields[fieldSource.Field] = ExtractContextFieldForSource(request, fieldSource.Source)
	}
	// Extract API-level configured fields

	// Match API TODO: 1. Optimize by moving to init? 2. Support regex?
	apiConfig := MatchApiToExtractConfig(request.Method, request.URL.Path, logConfig.Apis)
	if apiConfig != nil {
		for _, fieldSource := range apiConfig.Fields {
			contextfields[fieldSource.Field] = ExtractContextFieldForSource(request, fieldSource.Source)
		}
	}

	return contextfields
}
