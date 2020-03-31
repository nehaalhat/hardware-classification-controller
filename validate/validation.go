package validate

import (
	"fmt"
	//	"hardware-classification-controller/api/v1alpha1"
	ironic "hardware-classification-controller/ironic"
	"net"
)

type NICS map[string]interface{}
type Disks map[string]interface{}

// Validation fucntion to validate the parameters

func Validation(ironic_data ironic.Data) {
	var listOfNICS []NICS
	var listOfDisks []Disks
	CPUs := make(map[string]int)
	var RAMGb int
	fmt.Println("Ashu : Validation defination in validation.go")
	for _, host := range ironic_data.Host {
		for i := 0; i < len(host.Status.HardwareDetails.NIC); i++ {
			if host.Status.HardwareDetails.NIC[i].PXE == true && CheckValidIP(host.Status.HardwareDetails.NIC[i].IP) {
				nic := NICS{}
				nic["Name"] = host.Status.HardwareDetails.NIC[i].Name
				nic["PXE"] = host.Status.HardwareDetails.NIC[i].PXE
				nic["IP"] = host.Status.HardwareDetails.NIC[i].IP
				// save host in maps
				listOfNICS = append(listOfNICS, nic)
			} else {
				continue
			}
		}
		fmt.Println("nics : ", listOfNICS)
		for j := 0; j < len(host.Status.HardwareDetails.Storage); j++ {
			if host.Status.HardwareDetails.Storage[j].SizeBytes != 0 {
				disk := Disks{}
				disk["Name"] = host.Status.HardwareDetails.Storage[j].Name
				disk["SizeGb"] = ConvertBytesToGb(int64(host.Status.HardwareDetails.Storage[j].SizeBytes))
				listOfDisks = append(listOfDisks, disk)

			}
		}
		fmt.Println("Ashu Disks list : ", listOfDisks)
		if host.Status.HardwareDetails.CPU.Count > 0 {
			CPUs["Count"] = host.Status.HardwareDetails.CPU.Count
		} else {
			fmt.Println("Return error")
		}

		fmt.Println("Ashu CPU Count : ", CPUs)
		if host.Status.HardwareDetails.RAMMebibytes > 0 {
			RAMGb = host.Status.HardwareDetails.RAMMebibytes / 954
			// Convert mebibyte into Gb
		} else {
			fmt.Println("Return error")
		}
		fmt.Println("Ashu RAM in Gb: ", RAMGb)
	}
}

func CheckValidIP(NICIp string) bool {
	return net.ParseIP(NICIp) != nil
}

func ConvertBytesToGb(inBytes int64) float64 {
	inGb := float64(inBytes / 1024)
	inGb = inGb / 1024
	inGb = inGb / 1024
	return inGb
}
