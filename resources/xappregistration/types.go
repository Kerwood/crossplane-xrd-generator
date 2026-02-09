package xappregistration

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type XAppRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XAppRegistrationSpec   `json:"spec"`
	Status XAppRegistrationStatus `json:"status,omitempty"`
}

type XAppRegistrationSpec struct {
	RedirectURLs               []*string                   `json:"redirectURLs,omitempty" required:"true"`
	ServiceAccountName         string                      `json:"serviceAccountName" required:"true"`
	RoleAssignments            []RoleAssignments           `json:"roleAssignments,omitempty"`
	WriteConnectionSecretToRef *WriteConnectionSecretToRef `json:"writeConnectionSecretToRef,omitempty"`
}

type XAppRegistrationStatus struct {
	ClientID    string `json:"clientId,omitempty"`
	TenantID    string `json:"tenantId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type RoleAssignments struct {
	RoleName       string   `json:"roleName" required:"true"`
	Description    string   `json:"description" required:"true"`
	AssignedGroups []string `json:"assignedGroups" required:"true"`
}

type WriteConnectionSecretToRef struct {
	Name string `json:"name" required:"true"`
}
