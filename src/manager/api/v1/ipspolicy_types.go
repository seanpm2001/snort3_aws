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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IpsPolicySpec defines the desired state of IpsPolicy
type IpsPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	PolicyName string `json:"policyName"`
}

// IpsPolicyStatus defines the observed state of IpsPolicy
type IpsPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Information when was the last time the job was successfully scheduled.
	// +optional
	LastUpdateTime *metav1.Time `json:"lastUpdateTime,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// IpsPolicy is the Schema for the ipspolicies API
type IpsPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpsPolicySpec   `json:"spec,omitempty"`
	Status IpsPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IpsPolicyList contains a list of IpsPolicy
type IpsPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpsPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpsPolicy{}, &IpsPolicyList{})
}
