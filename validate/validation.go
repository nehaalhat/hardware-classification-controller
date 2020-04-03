package validate

import (
	"fmt"
	//	"net"
	//"hardware-classification-controller/api/v1alpha1"
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	//ironic "hardware-classification-controller/ironic"
)

// Validation fucntion to validate the parameters
/*
type CPU struct {
        Arch           string     `json:"arch"`
        Model          string     `json:"model"`
        ClockMegahertz ClockSpeed `json:"clockMegahertz"`
        Flags          []string   `json:"flags"`
        Count          int        `json:"count"`
}
*/
func Validation(myMap map[string]map[string]interface{}) {
	fmt.Println("Ashu : Validation defination in validation.go")
	//	fmt.Println("Ashu : node-1 cpu", myMap.)
	//fmt.Println("Ashu : myMap", myMap)
	//ironic_data = myMap
	data := bmh.HardwareDetails{}
	fmt.Println("Ashu  : BMH from Map :  ", data.CPU)
	fmt.Printf("Ashu  : type BMH from Map : %T\n:  ", data.CPU)
	for key, value := range myMap {
		fmt.Println("Key*******", key)
		for k, v := range value {
			if k == "CPU" {
				fmt.Println("Values*******", v)
				fmt.Printf("Ashu  : type V : %T\n:  ", v)
				cpu, _ := v.(bmh.CPU)
				//if ok {
				fmt.Println("Ashu : cpu details", cpu)
				//}
			}
			//fmt.Println("key*******", k)
			//fmt.Println("Values*******", v)
			//	fmt.Println("value for : ", k, "value : ", v)
		}
	}
	fmt.Println("------------------------------------------------")

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
