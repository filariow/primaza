package formatters

type JSONMapperStrategy string

const (
	IncludeMissingStrategy JSONMapperStrategy = "IncludeMissing"
	ExcludeMissingStrategy JSONMapperStrategy = "ExcludeMissing"
)

type JSONFormatter struct {
	BaseFormatter `json:",inline"`

	// Mapper is the Mapper configuration for the JSON Formatter
	// +required
	Mapper JSONMapper `json:"mapper"`
}

// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type JSONMapper struct {
	// Configuration is an inline configuration
	// +optional
	Configuration *JSONMapperConfiguration `json:"configuration,omitempty"`
	// ConfigurationRef is used to retrieve the data from a ConfigMap or Secret
	// +optional
	ConfigurationRef *ConfigurationRef `json:"configurationRef,omitempty"`
}

type JSONMapperConfiguration struct {
	// Strategy defines the default strategy for fields not mapped by Mappings
	// +required
	//+kubebuilder:validation:Enum=IncludeMissing;ExcludeMissing
	//+kubebuilder:default=ExcludeMissing
	Strategy JSONMapperStrategy `json:"strategy"`

	// MissingFieldsRoot can be used to define the property that will contain
	// the missing fields data. As an example, a value of `/data/others` will produce
	// a `{ "data": { "others": [ ... ] }, ... }` section in the configuration file.
	//+optional
	MissingFieldsRoot *string `json:"missingFieldsRoot,omitempty"`

	// Mappings contains the specification for mapping ServiceEndpointDefinition
	// data to properties in the JSON Configuration file.
	// +required
	Mappings []JSONFieldMapping `json:"fieldsMapping"`
}

type JSONFieldMapping struct {
	// SEDKey is the Key from the ServiceEndpointDefinition Secret to be mapped
	// +required
	SEDKey string `json:"sedKey"`

	// PropertyPath is the path for the property in the JSON configuration file
	// +required
	PropertyPath string `json:"propertyPath"`
}
