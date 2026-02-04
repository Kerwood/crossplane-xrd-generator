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
	RedirectURLs       []string `json:"redirectURL,omitempty"`
	GroupAssignments   []string `json:"groupAssignments,omitempty"`
	ServiceAccountName string   `json:"serviceAcccountName,omitempty"`
}

type XAppRegistrationStatus struct {
	ClientID    string `json:"clientId,omitempty"`
	TenantID    string `json:"tenantId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}
