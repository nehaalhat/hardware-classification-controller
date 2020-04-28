package filter

import (
	"fmt"
	hwcc "hardware-classification-controller/api/v1alpha1"
	valTypes "hardware-classification-controller/validation/validationModel"

	"k8s.io/apimachinery/pkg/api/resource"
)

// MinMaxComparison it will compare the minimum and maximum comparison based on the value provided by the user and check for the valid host
func MinMaxComparison(ProfileName string, validatedHost map[string]map[string]interface{}, expectedHardwareprofile hwcc.ExpectedHardwareConfiguration) []string {

	fmt.Println("Extracted HWDetails", expectedHardwareprofile)
	fmt.Printf("\n\n\n")
	var comparedHost []string

	for hostname, details := range validatedHost {

		isHostValid := true

		for _, value := range details {

			isValid := false

			cpu, CPUOK := value.(valTypes.CPU)
			if CPUOK {
				if checkCPUCount(cpu, expectedHardwareprofile.CPU) {
					isValid = true
				}
			}

			ram, RAMOK := value.(valTypes.RAM)
			if RAMOK {
				if checkRAM(ram, expectedHardwareprofile.RAM) {
					isValid = true
				}
			}

			nics, NICSOK := value.(valTypes.NIC)
			if NICSOK {
				if checkNICS(nics, expectedHardwareprofile.NIC) {
					isValid = true
				}
			}

			disk, DISKOK := value.(valTypes.Storage)
			if DISKOK {
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
			fmt.Printf("\n\n\n")

		} else {
			fmt.Println(hostname, " Does not matches profile ", ProfileName)
			fmt.Printf("\n\n\n")

		}

	}

	return comparedHost

}

func checkNICS(nics valTypes.NIC, expectedNIC hwcc.NIC) bool {
	fmt.Printf("\n")
	if expectedNIC.MaximumCount > 0 && expectedNIC.MinimumCount > 0 {
		fmt.Println("Provided Minimum Count for NICS", expectedNIC.MinimumCount, " and fetched count ", nics.Count)
		fmt.Println("Provided Maximum count for NICS", expectedNIC.MaximumCount, " and fetched count ", nics.Count)
		if expectedNIC.MinimumCount > nics.Count && expectedNIC.MaximumCount < nics.Count {
			return false
		}
	} else if expectedNIC.MaximumCount > 0 {
		fmt.Println("Provided Maximum count for NICS", expectedNIC.MaximumCount, " and fetched count ", nics.Count)
		if expectedNIC.MaximumCount < nics.Count {
			return false
		}

	} else if expectedNIC.MinimumCount > 0 {
		fmt.Println("Provided Minimum Count for NICS", expectedNIC.MinimumCount, " and fetched count ", nics.Count)
		if expectedNIC.MinimumCount > nics.Count {
			return false
		}

	}
	return true
}

func checkRAM(ram valTypes.RAM, expectedRAM hwcc.RAM) bool {
	fmt.Printf("\n")
	if expectedRAM.MaximumSizeGB > 0 && expectedRAM.MinimumSizeGB > 0 {
		fmt.Println("Provided Minimum Size for RAM", expectedRAM.MinimumSizeGB, " and fetched SIZE ", ram.RAMGb)
		fmt.Println("Provided Maximum Size for RAM", expectedRAM.MaximumSizeGB, " and fetched SIZE ", ram.RAMGb)
		if expectedRAM.MinimumSizeGB > ram.RAMGb && expectedRAM.MaximumSizeGB < ram.RAMGb {
			return false
		}
	} else if expectedRAM.MaximumSizeGB > 0 {
		fmt.Println("Provided Maximum Size for RAM", expectedRAM.MaximumSizeGB, " and fetched SIZE ", ram.RAMGb)
		if expectedRAM.MaximumSizeGB < ram.RAMGb {
			return false
		}

	} else if expectedRAM.MinimumSizeGB > 0 {
		fmt.Println("Provided Minimum Size for RAM", expectedRAM.MinimumSizeGB, " and fetched SIZE ", ram.RAMGb)
		if expectedRAM.MinimumSizeGB > ram.RAMGb {
			return false
		}

	}
	return true
}

func checkCPUCount(cpu valTypes.CPU, expectedCPU hwcc.CPU) bool {

	fmt.Printf("\n")

	if expectedCPU.MaximumCount > 0 && expectedCPU.MinimumCount > 0 {
		fmt.Println("Provided Minimum count for CPU", expectedCPU.MinimumCount, " and fetched count ", cpu.Count)
		fmt.Println("Provided Maximum count for CPU", expectedCPU.MaximumCount, " and fetched count ", cpu.Count)
		if expectedCPU.MinimumCount > cpu.Count && expectedCPU.MaximumCount < cpu.Count {
			return false
		}

	} else if expectedCPU.MaximumCount > 0 {
		fmt.Println("Provided Maximum count for CPU", expectedCPU.MaximumCount, " and fetched count ", cpu.Count)
		if expectedCPU.MaximumCount < cpu.Count {
			return false
		}

	} else if expectedCPU.MinimumCount > 0 {
		fmt.Println("Provided Minimum count for CPU", expectedCPU.MinimumCount, " and fetched count ", cpu.Count)
		if expectedCPU.MinimumCount > cpu.Count {
			return false
		}

	}

	if expectedCPU.MaximumSpeed != (resource.Quantity{}) && expectedCPU.MinimumSpeed != (resource.Quantity{}) {
		MinSpeed := float64(expectedCPU.MinimumSpeed.AsDec().UnscaledBig().Int64()) / 10
		MaxSpeed := float64(expectedCPU.MaximumSpeed.AsDec().UnscaledBig().Int64()) / 10

		fmt.Println("Provided Minimum ClockSpeed for CPU", MinSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
		fmt.Println("Provided Maximum ClockSpeed for CPU", MaxSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
		if MinSpeed > 0 && MaxSpeed > 0 {
			if MinSpeed > cpu.ClockSpeed && MaxSpeed < cpu.ClockSpeed {
				return false
			}

		}
	} else if expectedCPU.MaximumSpeed != (resource.Quantity{}) {
		MaxSpeed := float64(expectedCPU.MaximumSpeed.AsDec().UnscaledBig().Int64()) / 10
		fmt.Println("Provided Maximum ClockSpeed for CPU", MaxSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
		if MaxSpeed < cpu.ClockSpeed {
			return false
		}

	} else if expectedCPU.MinimumSpeed != (resource.Quantity{}) {
		MinSpeed := float64(expectedCPU.MinimumSpeed.AsDec().UnscaledBig().Int64()) / 10
		fmt.Println("Provided Minimum ClockSpeed for CPU", MinSpeed, " and fetched ClockSpeed ", cpu.ClockSpeed)
		if MinSpeed > cpu.ClockSpeed {
			return false
		}

	}

	return true

}

func checkDiskDetailsl(storage valTypes.Storage, expectedDisk hwcc.Disk) bool {
	fmt.Printf("\n")
	fmt.Println("Extracted Storage details", expectedDisk)

	if expectedDisk.MaximumCount > 0 && expectedDisk.MinimumCount > 0 {
		fmt.Println("Provided Minimum count for Disk", expectedDisk.MinimumCount, " and fetched count ", storage.Count)
		fmt.Println("Provided Maximum count for Disk", expectedDisk.MaximumCount, " and fetched count ", storage.Count)
		if expectedDisk.MaximumCount <= storage.Count && expectedDisk.MinimumCount <= storage.Count {
			for _, disk := range storage.Disk {
				fmt.Println("Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
				fmt.Println("Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
				if expectedDisk.MaximumIndividualSizeGB < disk.SizeGb && expectedDisk.MinimumIndividualSizeGB > disk.SizeGb {
					return false
				}
			}

		}
	} else if expectedDisk.MaximumCount > 0 {
		fmt.Println("Provided Maximum count for Disk", expectedDisk.MaximumCount, " and fetched count ", storage.Count)
		for _, disk := range storage.Disk {
			fmt.Println("Provided Maximum Size for Disk", expectedDisk.MaximumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
			if expectedDisk.MaximumIndividualSizeGB < disk.SizeGb {
				return false
			}
		}
	} else if expectedDisk.MinimumCount > 0 {
		fmt.Println("Provided Minimum count for Disk", expectedDisk.MinimumCount, " and fetched count ", storage.Count)
		for _, disk := range storage.Disk {
			fmt.Println("Provided Minimum Size for Disk", expectedDisk.MinimumIndividualSizeGB, " and fetched Size ", disk.SizeGb)
			if expectedDisk.MinimumIndividualSizeGB > disk.SizeGb {
				return false
			}
		}
	}

	return true
}
