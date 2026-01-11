package xexample

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type XExample struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XExampleSpec   `json:"spec"`
	Status XExampleStatus `json:"status,omitempty"`
}

type XExampleSpec struct {
	SomeString     string       `json:"someString" required:"true"`
	SomeInt        int          `json:"someInt,omitempty"`
	SomeBool       bool         `json:"someBool,omitempty"`
	SomeList       []string     `json:"someList,omitempty"`
	SomeObjectList []SomeStruct `json:"someObjectList,omitempty"`
}

type XExampleStatus struct {
	StatusFieldOne int    `json:"replicas,omitempty"`
	StatusFieldTwo string `json:"address,omitempty"`
}

type SomeStruct struct {
	FieldOne string `json:"fieldOne" required:"true"`
	FieldTwo string `json:"fieldTwo" required:"true"`
}
