package xdeployment

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type XDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              XDeploymentSpec   `json:"spec"`
	Status            XDeploymentStatus `json:"status,omitempty"`
}

type XDeploymentSpec struct {
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// +kubebuilder:validation:MinLength=1
	Tag string `json:"tag"`

	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=50
	Replicas *int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port *int32 `json:"port,omitempty"`

	// +kubebuilder:validation:Pattern=`^[a-z0-9][a-z0-9\-\.]+[a-z0-9]$`
	Hostname *string `json:"hostname,omitempty"`

	ServiceAccountName *string           `json:"serviceAccountName,omitempty"`
	Env                map[string]string `json:"env,omitempty"`
	SingleSignOn       *SingleSignOn     `json:"singleSignOn,omitempty"`
}

type XDeploymentStatus struct{}

type SingleSignOn struct {
	// +kubebuilder:default=false
	EnableAuthProxy            bool                        `json:"enableAuthProxy,omitempty"`
	ConnectionDetailsSecretRef *ConnectionDetailsSecretRef `json:"connectionDetailsSecretRef,omitempty"`
}

type ConnectionDetailsSecretRef struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}
