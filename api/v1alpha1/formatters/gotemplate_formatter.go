package formatters

type GoTemplateFormatter struct {
	BaseFormatter `json:",inline"`

	// Mapper is the Mapper configuration for the GoTemplate Formatter
	// +required
	Mapper *GoTemplateMapper `json:"mapper"`
}

// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type GoTemplateMapper struct {
	// Configuration is an inline configuration
	// +optional
	Configuration *GoTemplateMapperConfiguration `json:"configuration,omitempty"`
	// ConfigurationRef is used to retrieve the data from a ConfigMap or Secret
	// +optional
	ConfigurationRef *ConfigurationRef `json:"configurationRef,omitempty"`
}

type GoTemplateMapperConfiguration struct {
	// Template is a GoTemplate in which the ServiceEndpointDefinition data will be substituted.
	// It's the most flexible implementation for generating your custom Configuration.
	// +required
	Template string `json:"template"`
}
