/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TalosSpdSpec defines the desired state of TalosSpd
type TalosSpdSpec struct {
	Version string `json:"version"`
}

// TalosSpdStatus defines the observed state of TalosSpd
type TalosSpdStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TalosSpd is the Schema for the talosspds API
type TalosSpd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TalosSpdSpec   `json:"spec,omitempty"`
	Status TalosSpdStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TalosSpdList contains a list of TalosSpd
type TalosSpdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TalosSpd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TalosSpd{}, &TalosSpdList{})
}
