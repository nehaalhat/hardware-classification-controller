package filter

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"
	valTypes "hardware-classification-controller/validate/validateModel"
)

//MinimumFieldComparison check for the minimum validation
func MinimumFieldComparison(validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) {
	for hostname, details := range validatedHost {
		fmt.Println(hostname)
		isHostValid := true

		for _, value := range details {
			isValid := false

			cpu, ok := value.(valTypes.CPU)
			if ok {
				fmt.Println("CPU*************", cpu)
				if cpu.Count >= expectedHardwareprofile.CPU.Count {
					isValid = true
				} else {
					fmt.Println("CPU did not match")
				}
			}

			ram, ok := value.(valTypes.RAM)
			if ok {
				fmt.Println("RAM*************", ram)
				if ram.RAMGb >= expectedHardwareprofile.RAM {
					isValid = true
				} else {
					fmt.Println("RAM did not match")
				}
			}

			nics, ok := value.(valTypes.NIC)
			if ok {
				fmt.Println("NICS*************", nics)
				if nics.Count >= expectedHardwareprofile.NICS.Count {
					isValid = true
				} else {
					fmt.Println("NICS did not match")
				}
			}

			storage, ok := value.(valTypes.Storage)
			if ok {
				fmt.Println("storage*************", storage)
				if checkValidStorage(true, storage, expectedHardwareprofile.Disk) {
					isValid = true
				} else {
					fmt.Println("Storage did not match")
				}
			}

			if !isValid {
				isHostValid = false
				break
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

			isValid := false

			cpu, ok := value.(valTypes.CPU)
			if ok {
				fmt.Println("CPU*************", cpu)
				if cpu.Count <= expectedHardwareprofile.CPU.Count {
					isValid = true
				} else {
					fmt.Println("CPU did not match")
				}
			}

			ram, ok := value.(valTypes.RAM)
			if ok {
				fmt.Println("RAM*************", ram)
				if ram.RAMGb <= expectedHardwareprofile.RAM {
					isValid = true
				} else {
					fmt.Println("RAM did not match")
				}
			}

			nics, ok := value.(valTypes.NIC)
			if ok {
				fmt.Println("NICS*************", nics)
				if nics.Count <= expectedHardwareprofile.NICS.Count {
					isValid = true
				} else {
					fmt.Println("NICS did not match")
				}
			}

			storage, ok := value.(valTypes.Storage)
			if ok {
				fmt.Println("storage*************", storage)
				if checkValidStorage(false, storage, expectedHardwareprofile.Disk) {
					isValid = true
				} else {
					fmt.Println("Storage did not match")
				}
			}

			if !isValid {
				isHostValid = false
				break
			}
		}

		if isHostValid {
			fmt.Println(hostname, " Is a Valid Host")
		} else {
			fmt.Println(hostname, " Is NOT a Valid Host")
		}

	}
}

func checkValidStorage(filter bool, storage valTypes.Storage, expectedStorage hwcc.Disk) bool {

	if filter {
		if storage.Count >= expectedStorage.Count {
			for _, disk := range storage.Disk {
				if disk.SizeGb < expectedStorage.SizeGB {
					return false
				}
			}
		} else {
			return false
		}
	} else {
		if storage.Count <= expectedStorage.Count {
			for _, disk := range storage.Disk {
				if disk.SizeGb > expectedStorage.SizeGB {
					return false
				}
			}
		} else {
			return false
		}
	}

	return true
}
