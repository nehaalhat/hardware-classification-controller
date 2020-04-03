package validate

import (
	"encoding/json"
	"fmt"
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	"net"
)

type NICS map[string]interface{}
type Disks map[string]interface{}

// Validation fucntion to validate the parameters
func Validation(myMap map[string]map[string]interface{}) map[string]map[string]interface{} {
	var listOfNICS []NICS
	var listOfDisks []Disks
	CPUs := make(map[string]int)
	SystemVendor := make(map[string]string)
	var RAMGb int

	fmt.Println("Ashu : Validation defination in validation.go")
	data := bmh.HardwareDetails{}
	ValMap := make(map[string]map[string]interface{})
	for key, value := range myMap {
		myHWMap := make(map[string]interface{})
		//fmt.Println("  Ashu Key*******", key)
		//fmt.Println("  Ashu Value *******", value)
		jsonbody, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
		}
		err = json.Unmarshal(jsonbody, &data)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Ashu Data : ", data)
		fmt.Println("Ashu CPU count", data.CPU.Count)
		for i := 0; i < len(data.NIC); i++ {
			if data.NIC[i].PXE == true && CheckValidIP(data.NIC[i].IP) {
				nic := NICS{}
				nic["Name"] = data.NIC[i].Name
				nic["PXE"] = data.NIC[i].PXE
				//		nic["IP"] = data.NIC[i].IP
				// save host in maps
				listOfNICS = append(listOfNICS, nic)
			} else {
				continue
			}
		}
		fmt.Println("nics : ", listOfNICS)
		myHWMap["NICS"] = listOfNICS
		for j := 0; j < len(data.Storage); j++ {
			if data.Storage[j].SizeBytes != 0 {
				disk := Disks{}
				disk["Name"] = data.Storage[j].Name
				disk["SizeGb"] = ConvertBytesToGb(int64(data.Storage[j].SizeBytes))
				listOfDisks = append(listOfDisks, disk)

			}
		}
		fmt.Println("Ashu Disks list : ", listOfDisks)
		myHWMap["Storage"] = listOfDisks
		if data.CPU.Count > 0 {
			CPUs["Count"] = data.CPU.Count
		} else {
			fmt.Println("Return error")
		}
		myHWMap["CPU"] = CPUs
		fmt.Println("Ashu CPU Count : ", CPUs)
		if data.RAMMebibytes > 0 {
			RAMGb = data.RAMMebibytes / 954
			// Convert mebibyte into Gb
		} else {
			fmt.Println(" Ram Empty Return error")
		}
		fmt.Println("Ashu Ram : ", RAMGb)
		myHWMap["RAM"] = RAMGb
		if data.SystemVendor.Manufacturer == "Dell Inc." {
			SystemVendor["ManuFacturer"] = data.SystemVendor.Manufacturer
		}
		myHWMap["SystemVendor"] = SystemVendor
		fmt.Println("Ashu system vendor : ", SystemVendor)
		fmt.Println("--------------- Node Ends ------------------------------------------------")
		ValMap[key] = myHWMap

	}
	//fmt.Println("Ashu Final map  :: ", ValMap)
	return ValMap

}

func CheckValidIP(NICIp string) bool {
	return net.ParseIP(NICIp) != nil
}

func ConvertBytesToGb(inBytes int64) float64 {
	inGb := float64(inBytes / 1024 / 1024 / 1024)
	//	inGb = inGb / 1024
	//	inGb = inGb / 1024
	return inGb
}

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
