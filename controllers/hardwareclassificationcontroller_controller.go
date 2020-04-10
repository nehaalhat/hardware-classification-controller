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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"

	hwcc "hardware-classification-controller/api/v1alpha1"
	ironic "hardware-classification-controller/ironic"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HardwareClassificationControllerReconciler reconciles a HardwareClassificationController object
type HardwareClassificationControllerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile reconcile function
// +kubebuilder:rbac:groups=metal3.io.sigs.k8s.io,resources=hardwareclassificationcontrollers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io.sigs.k8s.io,resources=hardwareclassificationcontrollers/status,verbs=get;update;patch
func (r *HardwareClassificationControllerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	// Get HardwareClassificationController to get values for Namespace and ExpectedHardwareConfiguration
	hardwareClassification := &hwcc.HardwareClassificationController{}
	if err := r.Client.Get(ctx, req.NamespacedName, hardwareClassification); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Get ExpectedHardwareConfiguraton from hardwareClassification
	extractedProfile := hardwareClassification.Spec.ExpectedHardwareConfiguration
	r.Log.Info("Extracted expected hardware configuration successfully", "extractedProfile", extractedProfile)

	ironic_data := fetchHosts()

	validateMap := make(map[string]map[string]interface{})

	if extractedProfile != (hwcc.ExpectedHardwareConfiguration{}) {
		for _, host := range ironic_data.Host {
			hardwareDetailsMap := make(map[string]interface{})

			if extractedProfile.CPU != (hwcc.CPU{}) {
				hardwareDetailsMap["CPU"] = host.Status.HardwareDetails.CPU
			}
			if extractedProfile.Disk != (hwcc.Disk{}) {
				hardwareDetailsMap["Disk"] = host.Status.HardwareDetails.Storage
			}
			if extractedProfile.NICS != (hwcc.NICS{}) {
				hardwareDetailsMap["NICS"] = host.Status.HardwareDetails.NIC
			}
			if extractedProfile.SystemVendor != (hwcc.SystemVendor{}) {
				hardwareDetailsMap["SystemVendor"] = host.Status.HardwareDetails.SystemVendor
			}
			if extractedProfile.Firmware != (hwcc.Firmware{}) {
				hardwareDetailsMap["Firmware"] = host.Status.HardwareDetails.Firmware
			}
			if extractedProfile.RAM > 0 {
				hardwareDetailsMap["RAMMebibytes"] = host.Status.HardwareDetails.RAMMebibytes
			}

			validateMap[host.Metadata.Name] = hardwareDetailsMap
		}

	}
	fmt.Println("Map********", validateMap)

	return ctrl.Result{}, nil
}

// fetchHosts Retrive the introspection data
func fetchHosts() ironic.Data {
	jsonFile, err := os.Open("introspectionData.json")
	if err != nil {
		fmt.Println(err)
	}

	jsonString, _ := ioutil.ReadAll(jsonFile)

	ironicData := ironic.Data{}
	json.Unmarshal([]byte(jsonString), &ironicData)
	return ironicData
}

func (r *HardwareClassificationControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassificationController{}).
		Complete(r)
}
