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

// HardwareClassificationSpec defines the desired state of HardwareClassification
type HardwareClassificationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// HardwareCharacteristics defines expected hardware configurations for Cpu, Disk, Nic and Ram.
	HardwareCharacteristics HardwareCharacteristics `json:"hardwareCharacteristics,omitempty"`
}

// HardwareCharacteristics details to match with the host
type HardwareCharacteristics struct {
	// +optional
	Cpu *Cpu `json:"cpu,omitempty"`
	// +optional
	Disk *Disk `json:"disk,omitempty"`
	// +optional
	Nic *Nic `json:"nic,omitempty"`
	// +optional
	Ram *Ram `json:"ram,omitempty"`
}

// Cpu contains cpu details extracted from the hardware profile
type Cpu struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	// MinimumCount of cpu should be greater than 0
	// Ex. MinimumCount > 0
	MinimumCount int `json:"minimumCount,omitempty"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	// MaximumCount of cpu should be greater than 0 and greater than MinimumCount
	// Ex. MaximumCount > 0 && MaximumCount > MinimumCount
	MaximumCount int `json:"maximumCount,omitempty"`
	// +optional
	// +kubebuilder:validation:Minimum=1000
	// MinimumSpeed of cpu should be greater than 0
	// Ex. MinimumSpeed > 0
	// Ex. MinimumSpeed: 2600
	// User wants CPU speed 2.6 (in GHz), then s/he should specify as 2600 MHz
	MinimumSpeedMHz int32 `json:"minimumSpeedMHz,omitempty"`
	// +optional
	// +kubebuilder:validation:Minimum=1000
	// Maximum speed of cpu should be greater than 0 and greater than MinimumSpeed
	// Ex. MaximumSpeed > 0 && MaximumSpeed > MinimumSpeed
	// Ex. MaximumSpeed: 3200
	// User wants CPU speed 3.2 (in GHz), then he should specify as 3200 MHz
	MaximumSpeedMHz int32 `json:"maximumSpeedMHz,omitempty"`
}

// Disk contains disk details extracted from the hardware profile
type Disk struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	// MinimumCount of disk should be greater than 0
	// MinimumCount > 0
	MinimumCount int `json:"minimumCount,omitempty"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	// MinimumIndividualSizeGB should be greater than 0
	// Ex. MinimumIndividualSizeGB > 0
	MinimumIndividualSizeGB int64 `json:"minimumIndividualSizeGB,omitempty"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	// MaximumCount of disk should be greater than 0 and greater than MinimumCount
	// Ex. MaximumCount > 0 && MaximumCount > MinimumCount
	MaximumCount int `json:"maximumCount,omitempty"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	// Maximum individual size should be greater than 0 and greater than MinimumIndividualSizeGB
	// Ex. MaximumIndividualSizeGB > 0 && MaximumIndividualSizeGB > MinimumIndividualSizeGB
	MaximumIndividualSizeGB int64 `json:"maximumIndividualSizeGB,omitempty"`
}

// Nic contains nic details extracted from the hardware profile
type Nic struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	// Minimum count should be greater than 0
	// Ex. MinimumCount > 0
	MinimumCount int `json:"minimumCount,omitempty"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	// Maximum count should be greater than 0 and greater than MinimumCount
	// Ex. MaximumCount > 0 && MaximumCount > MinimumCount
	MaximumCount int `json:"maximumCount,omitempty"`
}

// Ram contains ram details extracted from the hardware profile
type Ram struct {
	// +optional
	// +kubebuilder:validation:Minimum=1
	// MinimumSizeGB of Ram should be greater than 0
	// Ex. MinimumSizeGB > 0
	MinimumSizeGB int `json:"minimumSizeGB,omitempty"`
	// +optional
	// +kubebuilder:validation:Minimum=1
	// MaximumSizeGB should be greater than 0 or greater than MinimumSizeGB
	// Ex. MaximumSizeGB > 0 && MaximumSizeGB > MinimumSizeGB
	MaximumSizeGB int `json:"maximumSizeGB,omitempty"`
}

// ProfileMatchStatus represents the state of the HardwareClassification
type ProfileMatchStatus string

const (
	// ProfileMatchStatusEmpty is the default status value
	ProfileMatchStatusEmpty ProfileMatchStatus = ""
	// ProfileMatchStatusMatched is the status value when the profile
	// matches to one of the BareMetalHost.
	ProfileMatchStatusMatched ProfileMatchStatus = "matched"
	// ProfileMatchStatusUnMatched is the status value when the profile
	// does not match to one of the BareMetalHost.
	ProfileMatchStatusUnMatched ProfileMatchStatus = "unmatched"
)

// ErrorType indicates the class of problem that has caused the HCC resource
// to enter an error state.
type ErrorType string

const (
	// LabelUpdateFailure is an error condition occurring when the
	// controller is unable to update label of BareMetalHost.
	LabelUpdateFailure ErrorType = "label update error"
	// LabelDeleteFailure is an error condition occurring when the
	// controller is unable to delete label of BareMetalHost.
	LabelDeleteFailure ErrorType = "label delete error"
	// FetchBMHListFailure is an error condition occurring when the
	// controller is unable to fetch BMH from BMO
	FetchBMHListFailure ErrorType = "fetch BMH error"
	// ProfileMisConfigured is an error condition occurring when the
	// extracted profile is empty.
	ProfileMisConfigured ErrorType = "Empty Profile Error"
	// Empty is an empty error
	Empty ErrorType = ""
)

// MatchedCount will provide matched count of Hosts per profile
type MatchedCount int

// UnmatchedCount will provide unmatched count of Hosts per profile
type UnmatchedCount int

const (
	// MatchedCount is the default status value
	MatchedCountEmpty MatchedCount = 0
	// MatchedCount is the default status value
	UnmatchedCountEmpty UnmatchedCount = 0
)

// Total hosts in error state
type ErrorHosts int

// Total hosts in registration error state
type RegistrationErrorHosts int

// Total hosts in Introspection error state
type IntrospectionErrorHosts int

const (
	// ErrorHosts is count of Hosts in error state
	ErrorHostsEmpty ErrorHosts = 0

	// RegistrationErrorHosts is count of Hosts in registration error state
	RegistrationErrorHostsEmpty RegistrationErrorHosts = 0

	// IntrospectionErrorHosts is count of Hosts in Introspection error state
	IntrospectionErrorHostsEmpty IntrospectionErrorHosts = 0
)

const (
	//NoBaremetalHost no bmo host found message
	NoBaremetalHost string = "No BareMetalHost Found in Ready State"
	//UpdateLabelError if failed to update labels of baremetal host
	UpdateLabelError string = "Failed to update labels of BareMetalHost"
	//NOError no error occurred
	NOError string = ""
)

// HardwareClassificationStatus defines the observed state of HardwareClassification
type HardwareClassificationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ErrorType indicates the type of failure encountered
	ErrorType ErrorType `json:"errorType,omitempty"`
	// ProfileMatchStatus identifies whether a applied profile is matches or not
	ProfileMatchStatus ProfileMatchStatus `json:"profileMatchStatus,omitempty"`
	// The count of matched Hosts per profile reported by hardwareclassification system
	MatchedCount MatchedCount `json:"matchedCount,omitempty"`
	// The count of unmatched Hosts per profile reported by hardwareclassification system
	UnmatchedCount UnmatchedCount `json:"unmatchedCount,omitempty"`
	// The count of Hosts in error state
	ErrorHosts ErrorHosts `json:"errorHosts,omitempty"`
	// The count of hosts in registration error state
	RegistrationErrorHosts RegistrationErrorHosts `json:"registrationErrorHosts,omitempty"`
	// The count of hosts in introspection error state
	IntrospectionErrorHosts IntrospectionErrorHosts `json:"introspectionErrorHosts,omitempty"`
	// The last error message reported by the hardwareclassification system
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=hwc;hc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="ProfileMatchStatus",type="string",JSONPath=".status.profileMatchStatus",description="Profile Match Status"
// +kubebuilder:printcolumn:name="MatchedCount",type="integer",JSONPath=".status.matchedCount",description="Total Matched hosts."
// +kubebuilder:printcolumn:name="UnmatchedCount",type="integer",JSONPath=".status.unmatchedCount",description="Total Unmatched hosts."
// +kubebuilder:printcolumn:name="ErrorHosts",type="integer",JSONPath=".status.errorHosts",description="Total error hosts."
// +kubebuilder:printcolumn:name="Registration-Error-Hosts",type="integer",JSONPath=".status.registrationErrorHosts",description="Total hosts in Registration error state."
// +kubebuilder:printcolumn:name="Introspection-Error-Hosts",type="integer",JSONPath=".status.introspectionErrorHosts",description="Total hosts in Introspection error state."
// +kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.errorMessage",description="Most recent error"

// HardwareClassification is the Schema for the hardwareclassifications API
type HardwareClassification struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HardwareClassificationSpec   `json:"spec,omitempty"`
	Status HardwareClassificationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HardwareClassificationList contains a list of HardwareClassification
type HardwareClassificationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HardwareClassification `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HardwareClassification{}, &HardwareClassificationList{})
}
