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

package controllers

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	bmoapis "github.com/metal3-io/baremetal-operator/pkg/apis"
	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/klogr"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ = Describe("Test Fech BMH List", func() {

	//Creating the hosts
	host1 := bmh.BareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host1",
			Namespace: "metal3",
		},
		Status: bmh.BareMetalHostStatus{
			Provisioning: bmh.ProvisionStatus{
				State: bmh.StateProvisioned,
			},
		},
	}
	host2 := bmh.BareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host2",
			Namespace: "metal3",
		},
		Status: bmh.BareMetalHostStatus{
			Provisioning: bmh.ProvisionStatus{
				State: bmh.StateReady,
			},
		},
	}
	host3 := bmh.BareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host3",
			Namespace: "metal3",
		},
		Status: bmh.BareMetalHostStatus{
			Provisioning: bmh.ProvisionStatus{
				State: bmh.StateInspecting,
			},
		},
	}
	host4 := bmh.BareMetalHost{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "host4",
			Namespace: "test",
		},
		Status: bmh.BareMetalHostStatus{
			Provisioning: bmh.ProvisionStatus{
				State: bmh.StateReady,
			},
		},
	}

	type testCasefetchBmhHostList struct {
		Hosts         []runtime.Object
		ExpectedHosts []bmh.BareMetalHost
		Namespace     string
	}

	DescribeTable("Test fetchBmhHostList",
		func(tc testCasefetchBmhHostList) {

			c := fakeclient.NewFakeClientWithScheme(setupSchemeMm(), tc.Hosts...)
			r := HardwareClassificationControllerReconciler{
				Client: c,
				Log:    klogr.New(),
			}

			result := fetchBmhHostList(context.TODO(), &r, tc.Namespace)

			if len(tc.ExpectedHosts) == 0 {
				Expect(len(result)).To(Equal(len(tc.ExpectedHosts)))
			} else {
				for i, host := range tc.ExpectedHosts {
					Expect(result[i].Name).To(Equal(host.Name))
				}
			}
		},

		Entry("Pick hosts in ready status and in metal3 namespace from BMH list", testCasefetchBmhHostList{
			Hosts:         []runtime.Object{&host1, &host2, &host3, &host4},
			Namespace:     "metal3",
			ExpectedHosts: []bmh.BareMetalHost{host2, host3},
		}),
		Entry("Empty list as no host in ready status from BMH list", testCasefetchBmhHostList{
			Hosts:         []runtime.Object{&host1},
			Namespace:     "metal3",
			ExpectedHosts: []bmh.BareMetalHost{},
		}),
		Entry("Empty list as no host under metal3 namespace", testCasefetchBmhHostList{
			Hosts:         []runtime.Object{&host4},
			Namespace:     "metal3",
			ExpectedHosts: []bmh.BareMetalHost{},
		}),
	)
})

//-----------------
// Helper functions
//-----------------
func setupSchemeMm() *runtime.Scheme {
	s := runtime.NewScheme()
	if err := bmoapis.AddToScheme(s); err != nil {
		panic(err)
	}
	return s
}