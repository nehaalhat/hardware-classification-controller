package hcmanager

import (
	hwcc "github.com/metal3-io/hardware-classification-controller/api/v1alpha1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	bmoapis "github.com/metal3-io/baremetal-operator/pkg/apis"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Hardware Classification Controller", func() {

	hostTest := getHosts()

	c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), hostTest...)
	hcManager := NewHardwareClassificationManager(c, klogr.New())

	It("Should fetch BaremetalHosts in ready state and under metal3 namespace", func() {
		result, _, err := hcManager.FetchBmhHostList(getNamespace())
		if err != nil {
			Expect(len(result)).To(BeZero())
		} else {
			Expect(len(result)).Should(Equal(2))
		}

	})

	It("Should return error fetching BaremetalHosts", func() {
		_, _, err := hcManager.FetchBmhHostList("sample")
		if err != nil {
			Expect(err).To(HaveOccurred())
		}

	})

	It("Should validate extracted hardware profile", func() {
		err := hcManager.ValidateExtractedHardwareProfile(getEmptyProfile())
		if err != nil {
			Expect(err).To(HaveOccurred())
		}
	})

	It("Should validate extracted hardware profile", func() {
		result := hcManager.ValidateExtractedHardwareProfile(getExtractedHardwareProfile())
		if result == nil {
			Expect(result).Should(BeNil())
		}
	})

	It("Should validate extracted hardware profile", func() {
		err := hcManager.ValidateExtractedHardwareProfile(getInvalidCPUProfile())
		if err != nil {
			Expect(err).To(HaveOccurred())
		}
	})

	It("Should validate extracted hardware profile", func() {
		err := hcManager.ValidateExtractedHardwareProfile(getInvalidDiskProfile())
		if err != nil {
			Expect(err).To(HaveOccurred())
		}
	})

	It("Should validate extracted hardware profile", func() {
		hardwareClassification := &hwcc.HardwareClassification{}
		err := hcManager.ValidateExtractedHardwareProfile(getInvalidNicProfile())
		SetStatus(hardwareClassification, hwcc.ProfileMatchStatusEmpty, hwcc.ProfileMisConfigured, err.Error())
		if err != nil {
			Expect(hardwareClassification.Status.ErrorType).Should(Equal(hwcc.ProfileMisConfigured))
			Expect(err).To(HaveOccurred())
		}
	})

	It("Should validate extracted hardware profile", func() {
		err := hcManager.ValidateExtractedHardwareProfile(getInvalidRAMProfile())
		if err != nil {
			Expect(err).To(HaveOccurred())
		}
	})

	It("Should validate extracted hardware profile", func() {
		err := hcManager.ValidateExtractedHardwareProfile(getMissingNicDetails())
		if err != nil {
			Expect(err).To(HaveOccurred())
		}
	})
})

//setupSchemeMm Add the bmoapi to our scheme
func setupSchemeMm() *runtime.Scheme {
	s := runtime.NewScheme()
	if err := bmoapis.AddToScheme(s); err != nil {
		panic(err)
	}
	if err := hwcc.AddToScheme(s); err != nil {
		panic(err)
	}
	return s
}
