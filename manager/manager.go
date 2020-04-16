package manager

import (
	hwcc "hardware-classification-controller/api/v1alpha1"
	"hardware-classification-controller/filter"
)

//Manager function call the comaprison algorithm as specified by the user
func Manager(profileName string, customFilter string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) {
	if customFilter == "maximum" {
		filter.MaximumFieldComparison(profileName, validatedHost, expectedHardwareprofile)
	} else {
		filter.MinimumFieldComparison(profileName, validatedHost, expectedHardwareprofile)
	}

}
