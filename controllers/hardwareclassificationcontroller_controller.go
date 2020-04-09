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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

	// Get expression rules from hardwareClassification
	expressionRules := hardwareClassification.Spec.Rules
	r.Log.Info("Extracted expression rules successfully", "expressionRules", expressionRules)

	// Check oneOf customFilter and expression rules is passed
	if extractedProfile.CustomFilter != "" && len(expressionRules) > 0 {
		validationError := errors.New("oneOf customFilter and expression rules can be passed")
		return ctrl.Result{}, validationError
	}

	ironic_data := fetchHosts()

	r.Log.Info("Fetched Baremetal host list successfully", "BareMetalHostList", ironic_data)

	myMap := make(map[string]map[string]interface{})
	myHWMap := make(map[string]interface{})

	if len(expressionRules) > 0 {
		for _, rules := range expressionRules {
			for _, host := range ironic_data.Host {
				if strings.ToUpper(rules.Field) == "CPU" {
					myHWMap["CPU"] = host.Status.HardwareDetails.CPU
				} else if strings.ToUpper(rules.Field) == "STORAGE" || strings.ToUpper(rules.Field) == "DISK" || strings.ToUpper(rules.Field) == "NUMBEROFDISK" {
					myHWMap["Storage"] = host.Status.HardwareDetails.Storage
				} else if strings.ToUpper(rules.Field) == "RAM" {
					myHWMap["RAMMebibytes"] = host.Status.HardwareDetails.RAMMebibytes
				} else if strings.ToUpper(rules.Field) == "NICS" || strings.ToUpper(rules.Field) == "NUMBEROFNICS" {
					myHWMap["NICS"] = host.Status.HardwareDetails.NIC
				} else if strings.ToUpper(rules.Field) == "SYSTEMVENDOR" {
					myHWMap["SystemVendor"] = host.Status.HardwareDetails.SystemVendor
				} else if strings.ToUpper(rules.Field) == "FIRMWARE" {
					myHWMap["Firmware"] = host.Status.HardwareDetails.NIC
				}
				myMap[host.Metadata.Name] = myHWMap
			}
		}

		fmt.Println("My Map********", myMap)
	}

	if extractedProfile != (hwcc.ExpectedHardwareConfiguration{}) {
		for _, host := range ironic_data.Host {
			//myHWMap := make(map[string]interface{})

			if extractedProfile.CPU != (hwcc.CPU{}) {
				myHWMap["CPU"] = host.Status.HardwareDetails.CPU
			} else if extractedProfile.Disk != (hwcc.Disk{}) {
				myHWMap["Disk"] = host.Status.HardwareDetails.Storage
			} else if extractedProfile.NICS != (hwcc.NICS{}) {
				myHWMap["NICS"] = host.Status.HardwareDetails.NIC
			} else if extractedProfile.SystemVendor != (hwcc.SystemVendor{}) {
				myHWMap["SystemVendor"] = host.Status.HardwareDetails.SystemVendor
			} else if extractedProfile.Firmware != (hwcc.Firmware{}) {
				myHWMap["Firmware"] = host.Status.HardwareDetails.Firmware
			} else if extractedProfile.RAM > 0 {
				myHWMap["RAMMebibytes"] = host.Status.HardwareDetails.RAMMebibytes
			}

			myMap[host.Metadata.Name] = myHWMap
		}
		fmt.Println("My Map********", myMap)
	}

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
