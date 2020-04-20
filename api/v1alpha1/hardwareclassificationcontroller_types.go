/*
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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HardwareClassificationControllerSpec defines the desired state of HardwareClassificationController
type HardwareClassificationControllerSpec struct {
	// ExpectedHardwareConfiguration defines expected hardware configurations for CPU, RAM, Disk, NIC.
	ExpectedHardwareConfiguration ExpectedHardwareConfiguration `json:"expectedValidationConfiguration"`
}

// ExpectedHardwareConfiguration details to match with the host
type ExpectedHardwareConfiguration struct {
	Namespace string `json:"namespace"`
	// +optional
	CPU CPU `json:"CPU"`
	// +optional
	Disk Disk `json:"Disk"`
	// +optional
	NIC NIC `json:"NIC"`
	// +optional
	RAM RAM `json:"RAM"`
}

// ClockSpeed is a clock speed in MHz
//type ClockSpeed float32

// ClockSpeed multipliers
/*const (
	MegaHertz ClockSpeed = 1.0
	GigaHertz            = 1000 * MegaHertz
)*/

// CPU count
type CPU struct {
	// +optional
	MinimumCount int `json:"minimumCount"`
	// +optional
	MaximumCount int `json:"maximumCount"`
	// +optional
	//MinimumSpeed ClockSpeed `json:"minimumSpeed"`
	// +optional
	//	MaximumSpeed ClockSpeed `json:"maximumSpeed"`
}

// Disk size and number of disks
type Disk struct {
	// +optional
	MinimumCount int `json:"minimumCount"`
	// +optional
	MinimumIndividualSizeGB int64 `json:"minimumIndividualSizeGB"`
	// +optional
	MaximumCount int `json:"maximumCount"`
	// +optional
	MaximumIndividualSizeGB int64 `json:"maximumIndividualSizeGB"`
}

// NIC count of nics cards
type NIC struct {
	// +optional
	MinimumCount int `json:"minimumCount"`
	// +optional
	MaximumCount int `json:"maximumCount"`
}

// RAM size
type RAM struct {
	// +optional
	MinimumSizeGB int `json:"minimumSizeGB"`
	// +optional
	MaximumSizeGB int `json:"maximumSizeGB"`
}

// HardwareClassificationControllerStatus defines the observed state of HardwareClassificationController
type HardwareClassificationControllerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ErrorMessage will be set in the event that there is a terminal problem
	// reconciling the BaremetalHost and will contain a more verbose string suitable
	// for logging and human consumption.

	ErrorMessage *string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true

// HardwareClassificationController is the Schema for the hardwareclassificationcontrollers API
type HardwareClassificationController struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HardwareClassificationControllerSpec   `json:"spec,omitempty"`
	Status HardwareClassificationControllerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HardwareClassificationControllerList contains a list of HardwareClassificationController
type HardwareClassificationControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HardwareClassificationController `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HardwareClassificationController{}, &HardwareClassificationControllerList{})
}
