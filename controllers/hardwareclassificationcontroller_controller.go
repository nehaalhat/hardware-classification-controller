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

//type for map
type M map[string]interface{}

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

	ironic_data := fetchHosts()
	// Get a list of BaremetalHost from Baremetal-Operator and metal3 namespace
	// bmhHostList, err := fetchBmhHostList(ctx, r, hardwareClassification.Spec.Namespace)
	// if err != nil {
	// 	return ctrl.Result{}, err
	// }
	r.Log.Info("Fetched Baremetal host list successfully", "BareMetalHostList", ironic_data)

	var mySlice []M
	//var myValues []interface{}
	myMap := make(map[string]interface{})
	myHWMap := make(map[string]interface{})

	for _, host := range ironic_data.Host {
		fmt.Println("hi********************")
		if extractedProfile.CPU != (hwcc.CPU{}) {
			myHWMap["CPU"] = host.Status.HardwareDetails.CPU
			//myValues = append(myValues, host.Status.HardwareDetails.CPU)
		}

		if extractedProfile.Disk != (hwcc.Disk{}) {
			myHWMap["Disk"] = host.Status.HardwareDetails.Storage
			//myValues = append(myValues, host.Status.HardwareDetails.Storage)
		}

		if extractedProfile.NICS != (hwcc.NICS{}) {
			myHWMap["NICS"] = host.Status.HardwareDetails.NIC
			//myValues = append(myValues, host.Status.HardwareDetails.NIC)
		}

		if extractedProfile.SystemVendor != (hwcc.SystemVendor{}) {
			myHWMap["SystemVendor"] = host.Status.HardwareDetails.SystemVendor
			//myValues = append(myValues, host.Status.HardwareDetails.SystemVendor)
		}

		if extractedProfile.Firmware != (hwcc.Firmware{}) {
			myHWMap["Firmware"] = host.Status.HardwareDetails.Firmware
			//myValues = append(myValues, host.Status.HardwareDetails.Firmware)
		}

		if extractedProfile.RAM > 0 {
			myHWMap["RAM"] = host.Status.HardwareDetails.RAMMebibytes
			//myValues = append(myValues, host.Status.HardwareDetails.RAMMebibytes)
		}

		myMap[host.Metadata.Name] = myHWMap
	}
	mySlice = append(mySlice, myMap)
	fmt.Println("My Slice**********************", mySlice)

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

// func fetchBmhHostList(ctx context.Context, r *HardwareClassificationControllerReconciler, namespace string) ([]bmh.BareMetalHost, error) {

// 	bmhHostList := bmh.BareMetalHostList{}
// 	validHostList := []bmh.BareMetalHost{}
// 	hardwareClassification := &hwcc.HardwareClassificationController{}

// 	opts := &client.ListOptions{
// 		Namespace: namespace,
// 	}

// 	// Get list of BareMetalHost
// 	err := r.Client.List(ctx, &bmhHostList, opts)
// 	if err != nil {
// 		setError(hardwareClassification, "Failed to get BareMetalHost List")
// 		return validHostList, err
// 	}

// 	// Get hosts in ready and inspecting status from bmhHostList
// 	for _, host := range bmhHostList.Items {
// 		if host.Status.Provisioning.State == "ready" || host.Status.Provisioning.State == "inspecting" {
// 			validHostList = append(validHostList, host)
// 		}
// 	}

// 	return validHostList, nil
//}

// setError sets the ErrorMessage field on the HardwareClassificationController
// func setError(hwcc *hwcc.HardwareClassificationController, message string) {
// 	hwcc.Status.ErrorMessage = pointer.StringPtr(message)
// }

func (r *HardwareClassificationControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassificationController{}).
		Complete(r)
}
