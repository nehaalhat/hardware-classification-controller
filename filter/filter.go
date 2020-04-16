package filter

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"
	valTypes "hardware-classification-controller/validate/validateModel"
)

//MinimumFieldComparison check for the minimum validation
func MinimumFieldComparison(ProfileName string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) {
	for hostname, details := range validatedHost {
		fmt.Println(hostname)
		isHostValid := true

		for _, value := range details {
			isValid := false

			cpu, ok := value.(valTypes.CPU)
			if ok {
				fmt.Println("Fetched CPU Count", cpu)
				if cpu.Count >= expectedHardwareprofile.CPU.Count {
					isValid = true
				} else {
					fmt.Println(hostname, "\n", "CPU count did not match")
				}
			}

			ram, ok := value.(valTypes.RAM)
			if ok {
				fmt.Println("Fetched RAM Size", ram)
				if ram.RAMGb >= expectedHardwareprofile.RAM {
					isValid = true
				} else {
					fmt.Println(hostname, "\n", "RAM size did not match")
				}
			}

			nics, ok := value.(valTypes.NIC)
			if ok {
				fmt.Println("Fetched NICS Count", nics)
				if nics.Count >= expectedHardwareprofile.NICS.Count {
					isValid = true
				} else {
					fmt.Println(hostname, "\n", "NICS count did not match")
				}
			}

			storage, ok := value.(valTypes.Storage)
			if ok {
				fmt.Println("Fetched Disk Details", storage)
				if checkValidStorage(true, hostname, storage, expectedHardwareprofile.Disk) {
					isValid = true
				}
			}

			if !isValid {
				isHostValid = false
				break
			}
		}

		if isHostValid {
			fmt.Println(hostname, " Matches profile ", ProfileName)
		} else {
			fmt.Println(hostname, " Does not matches profile ", ProfileName)
		}

	}
}

//MaximumFieldComparison check for the maximum validation
func MaximumFieldComparison(ProfileName string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) {
	for hostname, details := range validatedHost {
		fmt.Println(hostname)
		isHostValid := true

		for _, value := range details {

			isValid := false

			cpu, ok := value.(valTypes.CPU)
			if ok {
				fmt.Println("Fetched CPU Count", cpu)
				if cpu.Count <= expectedHardwareprofile.CPU.Count {
					isValid = true
				} else {
					fmt.Println(hostname, "\n", "CPU count did not match")
				}
			}

			ram, ok := value.(valTypes.RAM)
			if ok {
				fmt.Println("Fetched RAM Size", ram)
				if ram.RAMGb <= expectedHardwareprofile.RAM {
					isValid = true
				} else {
					fmt.Println(hostname, "\n", "RAM size did not match")
				}
			}

			nics, ok := value.(valTypes.NIC)
			if ok {
				fmt.Println("Fetched NICS Count", nics)
				if nics.Count <= expectedHardwareprofile.NICS.Count {
					isValid = true
				} else {
					fmt.Println(hostname, "\n", "NICS count did not match")
				}
			}

			storage, ok := value.(valTypes.Storage)
			if ok {
				fmt.Println("Fetched Disk Details", storage)
				if checkValidStorage(false, hostname, storage, expectedHardwareprofile.Disk) {
					isValid = true
				}
			}

			if !isValid {
				isHostValid = false
				break
			}
		}

		if isHostValid {
			fmt.Println(hostname, " Matches profile ", ProfileName)
		} else {
			fmt.Println(hostname, " Does not matches profile ", ProfileName)
		}

	}
}

func checkValidStorage(filter bool, hostname string, storage valTypes.Storage, expectedStorage hwcc.Disk) bool {

	if filter {
		if storage.Count >= expectedStorage.Count {
			for _, disk := range storage.Disk {
				if disk.SizeGb < expectedStorage.SizeGB {
					fmt.Println(hostname, "\n", "Disk size did not match")
					return false
				}
			}
		} else {
			fmt.Println(hostname, "\n", "Disk count did not match")
			return false
		}
	} else {
		if storage.Count <= expectedStorage.Count {
			for _, disk := range storage.Disk {
				if disk.SizeGb > expectedStorage.SizeGB {
					fmt.Println(hostname, "\n", "Disk size did not match")
					return false
				}
			}
		} else {
			fmt.Println(hostname, "\n", "Disk count did not match")
			return false
		}
	}

	return true
}
