/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hcmanager

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Hcmanager", func() {

	hostTest := getHosts()

	c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), hostTest...)
	hcManager := NewHardwareClassificationManager(c, klogr.New())

	It("Should validate the BaremetalHosts hardware details and return a list", func() {
		result, _, err := hcManager.FetchBmhHostList(getNamespace())
		if err != nil {
			Expect(len(hostTest)).Should(Equal(0))
		} else {
			validatedHardwareDetails := hcManager.ExtractAndValidateHardwareDetails(getExtractedHardwareProfile(), result)

			if len(validatedHardwareDetails) != 0 {
				Expect(len(result)).Should(Equal(len(validatedHardwareDetails)))
			} else {
				Expect(len(result)).To(BeZero())
			}
		}
	})
})
