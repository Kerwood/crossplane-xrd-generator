package xdeployment

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type XDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XDeploymentSpec   `json:"spec"`
	Status XDeploymentStatus `json:"status,omitempty"`
}

type XDeploymentSpec struct {
	Image              string            `json:"image" required:"true"`
	Tag                string            `json:"tag" required:"true"`
	Replicas           *int32            `json:"replicas,omitempty"`
	Port               *int              `json:"port,omitempty"`
	Hostname           *string           `json:"hostname,omitempty"`
	ServiceAccountName *string           `json:"serviceAccountName,omitempty"`
	Env                map[string]string `json:"env,omitempty"`
	SingleSignOn       *SingleSignOn     `json:"singleSignOn,omitempty"`
}

type XDeploymentStatus struct{}

type SingleSignOn struct {
	EnableAuthProxy            bool                        `json:"enableAuthProxy,omitempty"`
	ConnectionDetailsSecretRef *ConnectionDetailsSecretRef `json:"connectionDetailsSecretRef,omitempty"`
}

type ConnectionDetailsSecretRef struct {
	Name string `json:"name" required:"true"`
}
