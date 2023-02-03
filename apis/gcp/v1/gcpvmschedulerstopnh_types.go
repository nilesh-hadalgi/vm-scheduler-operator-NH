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

// GCPVMSchedulerStopNHSpec defines the desired state of GCPVMSchedulerStopNH
type GCPVMSchedulerStopNHSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of GCPVMSchedulerStopNH. Edit gcpvmschedulerstopnh_types.go to remove/update
	Zone      string `json:"zone"`
	Instance  string `json:"instance"`
	ProjectId string `json:"projectId"`
	Command   string `json:"command"`

	// Schedule period for the CronJob.
	// This spec allows you to setup the start schedule for VM
	StopSchedule string `json:"stopSchedule"`
	Image        string `json:"image"`
}

// GCPVMSchedulerStopNHStatus defines the observed state of GCPVMSchedulerStopNH
type GCPVMSchedulerStopNHStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	VMStopStatus string `json:"vmStopStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GCPVMSchedulerStopNH is the Schema for the gcpvmschedulerstopnhs API
type GCPVMSchedulerStopNH struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPVMSchedulerStopNHSpec   `json:"spec,omitempty"`
	Status GCPVMSchedulerStopNHStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GCPVMSchedulerStopNHList contains a list of GCPVMSchedulerStopNH
type GCPVMSchedulerStopNHList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPVMSchedulerStopNH `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GCPVMSchedulerStopNH{}, &GCPVMSchedulerStopNHList{})
}
