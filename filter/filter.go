package filter

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"
	valdata "hardware-classification-controller/validate/validated_data"
)

//MinimumFieldComparison check for the minimum validation
func MinimumFieldComparison(validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) {
	for hostname, details := range validatedHost {
		fmt.Println(hostname)
		isHostValid := false

		for _, value := range details {

			cpu, ok := value.(valdata.CPU)
			if ok {
				fmt.Println("CPU*************", cpu)
				if cpu.Count <= expectedHardwareprofile.CPU.Count {
					isHostValid = true
				}
			}

			ram, ok := value.(valdata.RAM)
			if ok {
				fmt.Println("RAM*************", ram)
				if ram.RAMGb <= expectedHardwareprofile.RAM {
					isHostValid = true
				}
			}

			nics, ok := value.(valdata.NIC)
			if ok {
				fmt.Println("NICS*************", nics)
				if nics.Count <= expectedHardwareprofile.NICS.NumberOfNICS {
					isHostValid = true
				}
			}

			storage, ok := value.(valdata.Storage)
			if ok {
				fmt.Println("storage*************", storage)
				if storage.SizeGb <= (expectedHardwareprofile.Disk.SizeBytesGB * int64(expectedHardwareprofile.Disk.NumberOfDisks)) {
					isHostValid = true
				}
			}
		}

		if isHostValid {
			fmt.Println(hostname, " Is a Valid Host")
		} else {
			fmt.Println(hostname, " Is NOT a Valid Host")
		}

	}
}

//MaximumFieldComparison check for the maximum validation
func MaximumFieldComparison(validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) {
	for hostname, details := range validatedHost {
		fmt.Println(hostname)
		isHostValid := false

		for _, value := range details {

			cpu, ok := value.(valdata.CPU)
			if ok {
				fmt.Println("CPU*************", cpu)
				if cpu.Count >= expectedHardwareprofile.CPU.Count {
					isHostValid = true
				}
			}

			ram, ok := value.(valdata.RAM)
			if ok {
				fmt.Println("RAM*************", ram)
				if ram.RAMGb >= expectedHardwareprofile.RAM {
					isHostValid = true
				}
			}

			nics, ok := value.(valdata.NIC)
			if ok {
				fmt.Println("NICS*************", nics)
				if nics.Count >= expectedHardwareprofile.NICS.NumberOfNICS {
					isHostValid = true
				}
			}

			storage, ok := value.(valdata.Storage)
			if ok {
				fmt.Println("storage*************", storage)
				if storage.SizeGb >= expectedHardwareprofile.Disk.SizeBytesGB {
					isHostValid = true
				}
			}
		}

		if isHostValid {
			fmt.Println(hostname, " Is a Valid Host")
		} else {
			fmt.Println(hostname, " Is NOT a Valid Host")
		}

	}
}
