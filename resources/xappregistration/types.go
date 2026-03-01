package xappregistration

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type XAppRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              XAppRegistrationSpec   `json:"spec"`
	Status            XAppRegistrationStatus `json:"status,omitempty"`
}

type XAppRegistrationSpec struct {
	// List of redirect URLs for the app registration.
	// +kubebuilder:validation:MinItems=1
	RedirectURLs []string `json:"redirectURLs"`

	// Name of the service account to associate with the app registration.
	// +kubebuilder:validation:MinLength=1
	ServiceAccountName string `json:"serviceAccountName"`

	// Role assignments to apply to the app registration.
	// +kubebuilder:validation:MinItems=1
	RoleAssignments []RoleAssignment `json:"roleAssignments,omitempty"`

	// Reference to the secret to write connection details to.
	WriteConnectionSecretToRef *WriteConnectionSecretToRef `json:"writeConnectionSecretToRef,omitempty"`
}

type XAppRegistrationStatus struct {
	ClientID    string `json:"clientId,omitempty"`
	TenantID    string `json:"tenantId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type RoleAssignment struct {
	// +kubebuilder:validation:MinLength=1
	RoleName string `json:"roleName"`

	// +kubebuilder:validation:MinLength=1
	Description string `json:"description"`

	// +kubebuilder:validation:MinItems=1
	AssignedGroups []string `json:"assignedGroups"`
}

type WriteConnectionSecretToRef struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}
