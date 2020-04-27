package filter

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"
	valTypes "hardware-classification-controller/validation/validationModel"

	"k8s.io/apimachinery/pkg/api/resource"
)

// MinMaxComparison it will compare the minimum and maximum comparison based on the value provided by the user and check for the valid host
func MinMaxComparison(ProfileName string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) []string {

	var comparedHost []string

	for hostname, details := range validatedHost {

		isHostValid := true

		for _, value := range details {

			isValid := false

			cpu, CPUOK := value.(valTypes.CPU)
			if CPUOK {
				fmt.Println("Fetched CPU Count", cpu)
				if checkCPUCount(cpu, expectedHardwareprofile.CPU) {
					isValid = true
				}
			}

			ram, RAMOK := value.(valTypes.RAM)
			if RAMOK {
				fmt.Println("Fetched RAM Count", ram)
				if checkRAM(ram, expectedHardwareprofile.RAM) {
					isValid = true
				}
			}

			nics, NICSOK := value.(valTypes.NIC)
			if NICSOK {
				fmt.Println("Fetched NICS Count", nics)
				if checkNICS(nics, expectedHardwareprofile.NIC) {
					isValid = true
				}
			}

			disk, DISKOK := value.(valTypes.Storage)
			if DISKOK {
				fmt.Println("Fetched Storage Details", disk)
				if checkDiskDetailsl(disk, expectedHardwareprofile.Disk) {
					isValid = true
				}
			}

			if !isValid {
				isHostValid = false
				break
			}

		}

		if isHostValid {
			comparedHost = append(comparedHost, hostname)
			fmt.Println(hostname, " Matches profile ", ProfileName)
		} else {
			fmt.Println(hostname, " Does not matches profile ", ProfileName)
		}

	}

	return comparedHost

}

func checkNICS(nics valTypes.NIC, expectedNIC hwcc.NIC) bool {
	if expectedNIC.MaximumCount > 0 && expectedNIC.MinimumCount > 0 {
		if expectedNIC.MinimumCount > nics.Count && expectedNIC.MaximumCount < nics.Count {
			return false
		}
	} else if expectedNIC.MaximumCount > 0 {
		if expectedNIC.MaximumCount < nics.Count {
			return false
		}

	} else if expectedNIC.MinimumCount > 0 {
		if expectedNIC.MinimumCount > nics.Count {
			return false
		}

	}
	return true
}

func checkRAM(ram valTypes.RAM, expectedRAM hwcc.RAM) bool {
	if expectedRAM.MaximumSizeGB > 0 && expectedRAM.MinimumSizeGB > 0 {
		if expectedRAM.MinimumSizeGB > ram.RAMGb && expectedRAM.MaximumSizeGB < ram.RAMGb {
			return false
		}
	} else if expectedRAM.MaximumSizeGB > 0 {
		if expectedRAM.MaximumSizeGB < ram.RAMGb {
			return false
		}

	} else if expectedRAM.MinimumSizeGB > 0 {
		if expectedRAM.MinimumSizeGB > ram.RAMGb {
			return false
		}

	}
	return true
}

func checkCPUCount(cpu valTypes.CPU, expectedCPU hwcc.CPU) bool {

	if expectedCPU.MaximumCount > 0 && expectedCPU.MinimumCount > 0 {
		if expectedCPU.MinimumCount > cpu.Count && expectedCPU.MaximumCount < cpu.Count {
			return false
		}

	} else if expectedCPU.MaximumCount > 0 {
		if expectedCPU.MaximumCount < cpu.Count {
			return false
		}

	} else if expectedCPU.MinimumCount > 0 {
		if expectedCPU.MinimumCount > cpu.Count {
			return false
		}

	}

	if expectedCPU.MaximumSpeed != (resource.Quantity{}) || expectedCPU.MinimumSpeed != (resource.Quantity{}) {
		MinSpeed := float64(expectedCPU.MinimumSpeed.AsDec().UnscaledBig().Int64()) / 10
		MaxSpeed := float64(expectedCPU.MaximumSpeed.AsDec().UnscaledBig().Int64()) / 10
		if MinSpeed > 0 && MaxSpeed > 0 {
			if MinSpeed > cpu.ClockSpeed && MaxSpeed < cpu.ClockSpeed {
				return false
			}

		}
	} else if expectedCPU.MaximumSpeed != (resource.Quantity{}) {
		MaxSpeed := float64(expectedCPU.MaximumSpeed.AsDec().UnscaledBig().Int64()) / 10
		if MaxSpeed < cpu.ClockSpeed {
			return false
		}

	} else if expectedCPU.MinimumSpeed != (resource.Quantity{}) {
		MinSpeed := float64(expectedCPU.MinimumSpeed.AsDec().UnscaledBig().Int64()) / 10
		if MinSpeed > cpu.ClockSpeed {
			return false
		}

	}

	return true

}

func checkDiskDetailsl(storage valTypes.Storage, expectedDisk hwcc.Disk) bool {

	if expectedDisk.MaximumCount > 0 && expectedDisk.MinimumCount > 0 {
		if expectedDisk.MaximumCount <= storage.Count && expectedDisk.MinimumCount <= storage.Count {
			for _, disk := range storage.Disk {
				if expectedDisk.MaximumIndividualSizeGB < disk.SizeGb && expectedDisk.MinimumIndividualSizeGB > disk.SizeGb {
					return false
				}
			}

		} else if expectedDisk.MaximumCount > 0 {
			for _, disk := range storage.Disk {
				if expectedDisk.MaximumIndividualSizeGB < disk.SizeGb {
					return false
				}
			}
		} else if expectedDisk.MinimumCount > 0 {
			for _, disk := range storage.Disk {
				if expectedDisk.MinimumIndividualSizeGB > disk.SizeGb {
					return false
				}
			}
		}

	}

	return true
}
