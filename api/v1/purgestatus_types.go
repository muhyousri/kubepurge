/*
Copyright 2024.

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

// PurgeStatusSpec defines the desired state of PurgeStatus.
type PurgeStatusSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of PurgeStatus. Edit purgestatus_types.go to remove/update
	//Foo string `json:"foo,omitempty"`

	CleanedNamespace string            `json:"cleanedNamespace"`
	LastPurgeTime    metav1.Time       `json:"lastPurgeTime"`
	PurgedResources  map[string]string `json:"purgedResources"`
}

// PurgeStatusStatus defines the observed state of PurgeStatus.
type PurgeStatusStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PurgeStatus is the Schema for the purgestatuses API.
type PurgeStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PurgeStatusSpec   `json:"spec,omitempty"`
	Status PurgeStatusStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PurgeStatusList contains a list of PurgeStatus.
type PurgeStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PurgeStatus `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PurgeStatus{}, &PurgeStatusList{})
}
