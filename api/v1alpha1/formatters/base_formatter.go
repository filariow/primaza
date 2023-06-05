package formatters

import (
	corev1 "k8s.io/api/core/v1"
)

type BaseFormatter struct {
	// FilePath Represent the folder in which the file should be projected.
	// If it's a relative path or empty, it will be prefixed with the SERVICE_BINDING_ROOT.
	// +optional
	FilePath *string `json:"filepath,omitempty"`

	// FileName is the name the file will have when projected in the pod filesystem.
	// It will be used together with the Filepath or SERVICE_BINDING_ROOT.
	// The name of the property in the ServiceEndpointDefinition Secret will match with
	// this Filepath or SERVICE_BINDING_ROOT and Filename concatentation.
	// +required
	FileName string `json:"filename"`
}

// +kubebuilder:validation:MaxProperties:=1
// +kubebuilder:validation:MinProperties:=1
type ConfigurationRef struct {
	// Selects a key of a ConfigMap.
	// +optional
	ConfigMapKeyRef *corev1.ConfigMapKeySelector `json:"configMapKeyRef,omitempty"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}
