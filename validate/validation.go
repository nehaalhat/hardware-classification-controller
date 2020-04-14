package validate

import (
	valTypes "hardware-classification-controller/validate/validateModel"
	"net"

	//ironic "hardware-classification-controller/ironic"
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
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
			cpu, ok := value.(bmh.CPU)
			if ok {
				validCPU := valTypes.CPU{
					Count: cpu.Count,
				}
				hardwareDetails[key] = validCPU
			}

			// Get the RAM details from the ironic host and validate it into new structure
			ram, ok := value.(int64)
			if ok {
				validRAM := valTypes.RAM{
					RAMGb: int(ConvertBytesToGb(ram)),
				}
				hardwareDetails[key] = validRAM
			}

			// Get the NICS details from the ironic host and validate it into new structure
			nics, ok := value.([]bmh.NIC)
			if ok {
				var validNICS valTypes.NIC
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
			storage, ok := value.([]bmh.Storage)
			if ok {
				var disks []valTypes.Disk

				for _, disk := range storage {
					disks = append(disks, valTypes.Disk{Name: disk.Name, SizeGb: ConvertBytesToGb(int64(disk.SizeBytes))})
				}
				validStorage := valTypes.Storage{
					Count: len(disks),
					Disk:  disks,
				}

				hardwareDetails[key] = validStorage
			}

		}

		validatedHostMap[hostName] = hardwareDetails

	}

	return validatedHostMap

}
