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
	"errors"
	"fmt"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"

	hwcc "hardware-classification-controller/api/v1alpha1"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var validHostList = []bmh.BareMetalHost{}
var checkValidHost = make(map[string]bool)

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
	fmt.Println("-----------------------------------------")
	fmt.Printf("Extracted expected hardware configuration successfully %+v", extractedProfile)
	fmt.Println("-----------------------------------------")

	//bmhList, err := fetchBmhHostList(ctx, r, hardwareClassification.Spec.ExpectedHardwareConfiguration.Namespace)
	// if err != nil {
	// 	r.Log.Error(err, "unable to fetch baremetal host list", "error", err.Error())
	// 	return ctrl.Result{}, nil
	// }

	extractedHardwareDetails, err := extractHardwareDetails(extractedProfile, validHostList)

	if err != nil {
		r.Log.Error(nil, "Unable to extract details", "error", err.Error())
		return ctrl.Result{}, nil
	}

	fmt.Println("-----------------------------------------")
	fmt.Printf("Extracted Hardware Details %+v", extractedHardwareDetails)
	fmt.Println("-----------------------------------------")

	return ctrl.Result{}, nil
}

func fetchBmhHostList(ctx context.Context, r *HardwareClassificationControllerReconciler, namespace string) []bmh.BareMetalHost {

	bmhHostList := bmh.BareMetalHostList{}
	//validHostList := []bmh.BareMetalHost{}
	//hardwareClassification := &hwcc.HardwareClassificationController{}

	opts := &client.ListOptions{
		Namespace: namespace,
	}

	// Get list of BareMetalHost
	err := r.Client.List(ctx, &bmhHostList, opts)
	if err != nil {
		r.Log.Error(nil, "Unable to extract details", "error", err.Error())
		return validHostList
	}

	// Get hosts in ready status from bmhHostList
	for _, host := range bmhHostList.Items {
		if host.Status.Provisioning.State == "ready" {
			validHostList = append(validHostList, host)
			checkValidHost[host.ObjectMeta.Name] = true
		}
	}

	return validHostList
}

func extractHardwareDetails(extractedProfile hwcc.ExpectedHardwareConfiguration,
	bmhList []bmh.BareMetalHost) (map[string]map[string]interface{}, error) {

	var err error
	extractedHardwareDetails := make(map[string]map[string]interface{})

	if extractedProfile != (hwcc.ExpectedHardwareConfiguration{}) {
		for _, host := range bmhList {
			introspectionDetails := make(map[string]interface{})

			if (extractedProfile.CPU == (hwcc.CPU{})) && (extractedProfile.Disk == (hwcc.Disk{})) &&
				(extractedProfile.NIC == (hwcc.NIC{})) && (extractedProfile.RAM == (hwcc.RAM{})) {
				err = errors.New("Ateleast one of the configuration should be provided")
				break
			}

			if extractedProfile.CPU != (hwcc.CPU{}) {
				if extractedProfile.CPU.MinimumCount > 0 || extractedProfile.CPU.MaximumCount > 0 {
					introspectionDetails["CPU"] = host.Status.HardwareDetails.CPU
				} else {
					err = errors.New("Enter valid CPU Count")
					break
				}
			}

			if extractedProfile.Disk != (hwcc.Disk{}) {
				if (extractedProfile.Disk.MinimumCount > 0 || extractedProfile.Disk.MinimumIndividualSizeGB > 0) ||
					(extractedProfile.Disk.MaximumCount > 0 || extractedProfile.Disk.MaximumIndividualSizeGB > 0) {
					introspectionDetails["Disk"] = host.Status.HardwareDetails.Storage
				} else {
					err = errors.New("Enter valid Disk Details")
					break
				}
			}

			if extractedProfile.NIC != (hwcc.NIC{}) {
				if extractedProfile.NIC.MinimumCount > 0 || extractedProfile.NIC.MaximumCount > 0 {
					introspectionDetails["NIC"] = host.Status.HardwareDetails.NIC
				} else {
					err = errors.New("Enter valid NICS Count")
					break
				}
			}

			if extractedProfile.RAM != (hwcc.RAM{}) {
				if extractedProfile.RAM.MinimumSizeGB > 0 || extractedProfile.RAM.MaximumSizeGB > 0 {
					introspectionDetails["RAMMebibytes"] = host.Status.HardwareDetails.RAMMebibytes
				} else {
					err = errors.New("Enter valid RAM size in GB")
					break
				}
			}

			if len(introspectionDetails) > 0 {
				extractedHardwareDetails[host.ObjectMeta.Name] = introspectionDetails
			}
		}

	}
	if err != nil {
		return extractedHardwareDetails, err
	}

	return extractedHardwareDetails, nil
}

// SetupWithManager will add watches for this controller
func (r *HardwareClassificationControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassificationController{}).
		Watches(
			&source.Kind{Type: &bmh.BareMetalHost{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(r.BareMetalHostToHardwareClassification),
			},
		).
		Complete(r)
}

// BareMetalHostToHardwareClassification will return a reconcile request for a
// HardwareClassification if the event is for a BareMetalHost.
func (r *HardwareClassificationControllerReconciler) BareMetalHostToHardwareClassification(obj handler.MapObject) []ctrl.Request {
	result := []ctrl.Request{}

	if len(validHostList) == 0 {
		validHostList = fetchBmhHostList(context.Background(), r, "metal3")
		//fmt.Println("First time fetched BMH list", validHostList)
		fmt.Println("Valid Host Map**************", checkValidHost)
	}

	if host, ok := obj.Object.(*bmh.BareMetalHost); ok {
		//exist := true
		//isValid := true

		name := client.ObjectKey{Namespace: "default", Name: "hardwareclassificationcontroller-sample"}

		if checkValidHost[host.ObjectMeta.Name] && host.Status.Provisioning.State != "ready" {
			fmt.Println("Host Found and not in ready state**********", host.ObjectMeta.Name)
			for i, validHost := range validHostList {
				if validHost.ObjectMeta.Name == host.ObjectMeta.Name {
					validHostList = append(validHostList[:i], validHostList[i+1:]...)
					checkValidHost[validHost.ObjectMeta.Name] = false
					fmt.Println("Updated Map********", checkValidHost)
					fmt.Println("Updated List*******", validHostList)
					result = append(result, ctrl.Request{NamespacedName: name})
				}
			}
		} else if !checkValidHost[host.ObjectMeta.Name] && host.Status.Provisioning.State == "ready" {
			fmt.Println("Host Not found so Appending BMH LIST*************", host.ObjectMeta.Name)
			validHostList = append(validHostList, *host)
			checkValidHost[host.ObjectMeta.Name] = true
			result = append(result, ctrl.Request{NamespacedName: name})
		}

		/*for i := 0; i < len(validHostList); i++ {
		//exist = false
		//isValid := false

		if checkValidHost[host.ObjectMeta.Name] && host.Status.Provisioning.State != "ready" {
			//exist = true
			fmt.Println("Inside If condition**********", host.ObjectMeta.Name)
			fmt.Println("Length of validHostList*********", validHostList)
			validHostList = append(validHostList[:1], validHostList[i+1:]...)
			checkValidHost[host.ObjectMeta.Name] = false
			result = append(result, ctrl.Request{NamespacedName: name})
		} else if !checkValidHost[host.ObjectMeta.Name] && host.Status.Provisioning.State == "ready" {
			//isValid = true
			//exist = false
			//break
			fmt.Println("Appending BMH LIST*************", host.ObjectMeta.Name)
			validHostList = append(validHostList, *host)
			result = append(result, ctrl.Request{NamespacedName: name})
		}

		/*if !exist {
			isValid = false
		}*/

	}
	/*if !isValid && host.Status.Provisioning.State == "ready" {
		fmt.Println("Appending BMH LIST*************", host.ObjectMeta.Name)
		validHostList = append(validHostList, *host)
		result = append(result, ctrl.Request{NamespacedName: name})
	}*/

	return result
}
