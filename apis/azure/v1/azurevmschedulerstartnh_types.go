/*
Copyright 2023.

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

// AzureVMSchedulerStartNHSpec defines the desired state of AzureVMSchedulerStartNH
type AzureVMSchedulerStartNHSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AzureVMSchedulerStartNH. Edit azurevmschedulerstartnh_types.go to remove/update
	AzureVMName string `json:"azureVMName"`
	AzureRGName string `json:"azureRGName"`
	Command     string `json:"command"`

	// Schedule period for the CronJob.
	// This spec allows you to setup the start schedule for VM
	StartSchedule string `json:"startSchedule"`
	Image         string `json:"image"`
}

// AzureVMSchedulerStartNHStatus defines the observed state of AzureVMSchedulerStartNH
type AzureVMSchedulerStartNHStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	VMStartStatus string `json:"vmStartStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AzureVMSchedulerStartNH is the Schema for the azurevmschedulerstartnhs API
type AzureVMSchedulerStartNH struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureVMSchedulerStartNHSpec   `json:"spec,omitempty"`
	Status AzureVMSchedulerStartNHStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AzureVMSchedulerStartNHList contains a list of AzureVMSchedulerStartNH
type AzureVMSchedulerStartNHList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureVMSchedulerStartNH `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzureVMSchedulerStartNH{}, &AzureVMSchedulerStartNHList{})
}
