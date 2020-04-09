package validate

import (
	valdata "hardware-classification-controller/validate/validated_data"
	"net"

	ironic "hardware-classification-controller/ironic"
)

//CheckValidIP uses net package to check if the IP is valid or not
func CheckValidIP(NICIp string) bool {
	return net.ParseIP(NICIp) != nil
}

//ConvertBytesToGb it converts the Byte into GB
func ConvertBytesToGb(inBytes int64) int64 {
	inGb := (inBytes / 1024 / 1024 / 1024)
	return inGb
}

//Validation this function will validate the host and create a new map with structered details
func Validation(hostDetails map[string]map[string]interface{}) map[string]map[string]interface{} {

	validatedHostMap := make(map[string]map[string]interface{})

	for hostName, details := range hostDetails {
		hardwareDetails := make(map[string]interface{})
		for key, value := range details {

			// Get the CPU details from the ironic host and validate it into new structure
			cpu, ok := value.(ironic.CPU)
			if ok {
				validCPU := valdata.CPU{
					Count: cpu.Count,
				}
				hardwareDetails[key] = validCPU
			}

			// Get the RAM details from the ironic host and validate it into new structure
			ram, ok := value.(int64)
			if ok {
				validRAM := valdata.RAM{
					RAMGb: int(ConvertBytesToGb(ram)),
				}
				hardwareDetails[key] = validRAM
			}

			// Get the NICS details from the ironic host and validate it into new structure
			nics, ok := value.([]ironic.NIC)
			if ok {
				var validNICS valdata.NIC
				for _, NIC := range nics {
					if NIC.PXE && CheckValidIP(NIC.IP) {
						validNICS.Name = NIC.Name
						validNICS.PXE = NIC.PXE
					}
				}

				validNICS.Count = len(nics)
				hardwareDetails[key] = validNICS
			}

			// Get the Storage details from the ironic host and validate it into new structure
			storage, ok := value.([]ironic.Storage)
			if ok {
				var sizeGB int64
				for _, disk := range storage {
					sizeGB += ConvertBytesToGb(int64(disk.SizeBytes))
				}
				validStorage := valdata.Storage{
					SizeGb: sizeGB,
				}

				hardwareDetails[key] = validStorage
			}

		}

		validatedHostMap[hostName] = hardwareDetails

	}

	return validatedHostMap

}

// // Validation fucntion to validate the parameters
// func Validation(myMap map[string]map[string]interface{}) map[string]map[string]interface{} {
// 	ValMap := make(map[string]map[string]interface{})
// 	data := bmh.HardwareDetails{}
// 	var validNICList []valdata.NIC
// 	var validStorageList []valdata.Storage
// 	var validCPU valdata.CPU
// 	var RAM valdata.RAM
// 	var validSystemVendor valdata.HardwareSystemVendor
// 	for key, value := range myMap {
// 		myHWMap := make(map[string]interface{})
// 		jsonbody, err := json.Marshal(value)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		err = json.Unmarshal(jsonbody, &data)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		for i := 0; i < len(data.NIC); i++ {
// 			if data.NIC[i].PXE == true && CheckValidIP(data.NIC[i].IP) {
// 				var validNIC valdata.NIC
// 				validNIC.Name = data.NIC[i].Name
// 				validNIC.PXE = data.NIC[i].PXE
// 				validNICList = append(validNICList, validNIC)
// 			} else {
// 				continue
// 			}
// 		}
// 		myHWMap["NICS"] = validNICList

// 		for j := 0; j < len(data.Storage); j++ {
// 			if data.Storage[j].SizeBytes != 0 {
// 				var validStorage valdata.Storage
// 				validStorage.Name = data.Storage[j].Name
// 				validStorage.SizeGb = ConvertBytesToGb(int64(data.Storage[j].SizeBytes))
// 				validStorageList = append(validStorageList, validStorage)

// 			} else {
// 				continue
// 			}
// 		}
// 		myHWMap["Storage"] = validStorageList
// 		if data.CPU.Count > 0 {
// 			validCPU.Count = data.CPU.Count
// 		} else {
// 			fmt.Println("Return error")
// 		}
// 		myHWMap["CPU"] = validCPU
// 		if data.RAMMebibytes > 0 {
// 			RAM.RAMGb = data.RAMMebibytes / 954
// 			// Convert mebibyte into Gb
// 		} else {
// 			fmt.Println("No valid RAM")
// 		}
// 		myHWMap["RAM"] = RAM
// 		if data.SystemVendor.Manufacturer == "Dell Inc." {
// 			validSystemVendor.Manufacturer = data.SystemVendor.Manufacturer
// 		}
// 		myHWMap["SystemVendor"] = validSystemVendor
// 		ValMap[key] = myHWMap

// 	}
// 	return ValMap

// }

/*//	fmt.Println("Ashu : Fetched Baremetal host list successfully", "BareMetalHostList", ironic_data)
	for _, host := range ironic_data.Host {
		for i := 0; i < 10; i++ {
			fmt.Println("Ashu : NIC's name : ", host.Status.HardwareDetails.NIC[i].Name)
			//			if CheckValidIP(host.Status.HardwareDetails.NIC[i].IP) {
			// save host in struct
			fmt.Println("Ashu : NIC's IP : ", host.Status.HardwareDetails.NIC[i].IP)
			//			} else {
			//				continue
			//			}
			fmt.Println("Ashu : NIC's PXE : ", host.Status.HardwareDetails.NIC[i].PXE)
		}
		for j := 0; j < 10; j++ {
			fmt.Println("Ashu : Storage name : ", host.Status.HardwareDetails.Storage[j].Name)
			fmt.Println("Ashu : Storage size : ", host.Status.HardwareDetails.Storage[j].SizeBytes)
		}
		//	fmt.Println("Ashu : NIC Details********", host.Status.HardwareDetails.NIC[0].Name)
	}
}
*/
/*
func CheckValidIP(NICIp string) bool {
	return net.ParseIP(NICIp) != nil
}
*/
/*
import (
	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
)

//Comparison function compare the host against the profile and filter the valid host
func Comparison(hosts []bmh.BareMetalHost, profiles []hwcc.ExpectedHardwareConfiguration) map[interface{}][]hwcc.ExpectedHardwareConfiguration {

	validHost := make(map[interface{}][]hwcc.ExpectedHardwareConfiguration)
	for _, host := range hosts {
		for _, profile := range profiles {
			if host.Status.HardwareDetails.CPU.Count >= profile.MinimumCPU.Count &&
				int64(host.Status.HardwareDetails.Storage[0].SizeBytes) >= (profile.MinimumDisk.SizeBytesGB*1024*1024) &&
				len(host.Status.HardwareDetails.NIC) >= profile.MinimumNICS.NumberOfNICS &&
				host.Status.HardwareDetails.RAMMebibytes >= (profile.MinimumRAM*1024) {

				hostDetails := &host
				newHost, ok := validHost[hostDetails]
				if ok {
					validHost[hostDetails] = append(newHost, profile)
				} else {
					var validProfile []hwcc.ExpectedHardwareConfiguration
					validHost[hostDetails] = append(validProfile, profile)
				}
			}
		}
	}

	return validHost

}
*/
