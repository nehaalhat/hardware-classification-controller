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
		isHostValid := false

		for _, value := range details {

			cpu, ok := value.(valTypes.CPU)
			if ok {
				fmt.Println("CPU*************", cpu)
				if cpu.Count <= expectedHardwareprofile.CPU.Count {
					isHostValid = true
				}
			}

			ram, ok := value.(valTypes.RAM)
			if ok {
				fmt.Println("RAM*************", ram)
				if ram.RAMGb <= expectedHardwareprofile.RAM {
					isHostValid = true
				}
			}

			nics, ok := value.(valTypes.NIC)
			if ok {
				fmt.Println("NICS*************", nics)
				if nics.Count <= expectedHardwareprofile.NICS.Count {
					isHostValid = true
				}
			}

			storage, ok := value.(valTypes.Storage)
			if ok {
				fmt.Println("storage*************", storage)
				if checkValidStorage(true, storage, expectedHardwareprofile.Disk) {
					isHostValid = true
				}
			}

			if !isHostValid {
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

			cpu, ok := value.(valTypes.CPU)
			if ok {
				fmt.Println("CPU*************", cpu)
				if cpu.Count >= expectedHardwareprofile.CPU.Count {
					isHostValid = true
				}
			}

			ram, ok := value.(valTypes.RAM)
			if ok {
				fmt.Println("RAM*************", ram)
				if ram.RAMGb >= expectedHardwareprofile.RAM {
					isHostValid = true
				}
			}

			nics, ok := value.(valTypes.NIC)
			if ok {
				fmt.Println("NICS*************", nics)
				if nics.Count >= expectedHardwareprofile.NICS.Count {
					isHostValid = true
				}
			}

			storage, ok := value.(valTypes.Storage)
			if ok {
				fmt.Println("storage*************", storage)
				if checkValidStorage(false, storage, expectedHardwareprofile.Disk) {
					isHostValid = true
				}
			}

			if !isHostValid {
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
		if storage.Count <= expectedStorage.Count {
			for _, disk := range storage.Disk {
				if disk.SizeGB <= expectedStorage.SizeGB {
					continue
				} else {
					return false
				}
			}
		} else {
			return false
		}
	} else {
		if storage.Count >= expectedStorage.Count {
			for _, disk := range storage.Disk {
				if disk.SizeGB >= expectedStorage.SizeGB {
					continue
				} else {
					return false
				}
			}
		} else {
			return false
		}
	}

	return true
}
