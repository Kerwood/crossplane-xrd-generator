package xexample

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
type XExample struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              XExampleSpec   `json:"spec"`
	Status            XExampleStatus `json:"status,omitempty"`
}

// +kubebuilder:validation:XValidation:rule="self.someInt <= 100 || self.someBool == true",message="someInt can only exceed 100 if someBool is true"
type XExampleSpec struct {
	// A required string with length and pattern constraints.
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9][a-z0-9-]*[a-z0-9]$`
	SomeString string `json:"someString"`

	// An integer with a default and range constraints.
	// +kubebuilder:default=10
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=100
	SomeInt int `json:"someInt,omitempty"`

	// A boolean with a default value.
	// +kubebuilder:default=false
	SomeBool bool `json:"someBool,omitempty"`

	// A list of strings with item constraints and uniqueness enforced.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=set
	SomeList []string `json:"someList,omitempty"`

	// A list of objects with size constraints.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=5
	SomeObjectList []SomeStruct `json:"someObjectList,omitempty"`
}

type XExampleStatus struct {
	// +kubebuilder:validation:Minimum=0
	StatusFieldOne int `json:"replicas,omitempty"`

	// +kubebuilder:validation:Pattern=`^(\d{1,3}\.){3}\d{1,3}$`
	StatusFieldTwo string `json:"address,omitempty"`
}

// +kubebuilder:validation:XValidation:rule="self.fieldOne != self.fieldTwo",message="fieldOne and fieldTwo must be different"
type SomeStruct struct {
	// +kubebuilder:validation:Enum=alpha;beta;gamma
	// +kubebuilder:default=alpha
	FieldOne string `json:"fieldOne,omitempty"`

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	FieldTwo string `json:"fieldTwo,omitempty"`
}
