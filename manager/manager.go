package manager

import (
	hwcc "hardware-classification-controller/api/v1alpha1"
	filt "hardware-classification-controller/filter"
)

//Manager function call the comaprison algorithm as specified by the user
func Manager(filter string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) {
	if filter == "maximum" {
		filt.MaximumFieldComparison(validatedHost, expectedHardwareprofile)
	} else {
		filt.MinimumFieldComparison(validatedHost, expectedHardwareprofile)
	}

}
