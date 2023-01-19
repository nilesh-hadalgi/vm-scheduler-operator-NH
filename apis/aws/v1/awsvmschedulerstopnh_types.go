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

// AWSVMSchedulerStopNHSpec defines the desired state of AWSVMSchedulerStopNH
type AWSVMSchedulerStopNHSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of AWSVMSchedulerStopNH. Edit awsvmschedulerstopnh_types.go to remove/update
	InstanceIds string `json:"instanceIds"`

	// Schedule period for the CronJob.
	// This spec allows you to setup the stop schedule for VM
	StopSchedule string `json:"stopSchedule"`

	// This spec allows you to supply image name which will start/stop VM
	Image string `json:"image"`
}

// AWSVMSchedulerStopNHStatus defines the observed state of AWSVMSchedulerStopNH
type AWSVMSchedulerStopNHStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Schedule period for the CronJob.
	// This will show the status of stop action for the VM(s)
	VMStopStatus string `json:"vmStopStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AWSVMSchedulerStopNH is the Schema for the awsvmschedulerstopnhs API
type AWSVMSchedulerStopNH struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSVMSchedulerStopNHSpec   `json:"spec,omitempty"`
	Status AWSVMSchedulerStopNHStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AWSVMSchedulerStopNHList contains a list of AWSVMSchedulerStopNH
type AWSVMSchedulerStopNHList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSVMSchedulerStopNH `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSVMSchedulerStopNH{}, &AWSVMSchedulerStopNHList{})
}
