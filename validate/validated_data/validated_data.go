package validation

type RAM struct {
	RAMGb        int                  `json:"ramMebibytes"`
}

type HardwareSystemVendor struct {
	Manufacturer string `json:"manufacturer"`
}

type NIC struct {
	Name string `json:"name"`
	PXE  bool   `json:"pxe"`
}

type Storage struct {
	Name      string   `json:"name"`
	SizeGb	  int64    `json:"sizeBytes"`
}

type CPU struct {
	Count int    `json:"count"`
}
