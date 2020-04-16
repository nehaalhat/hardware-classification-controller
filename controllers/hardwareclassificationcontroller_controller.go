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
	"k8s.io/utils/pointer"

	hwcc "hardware-classification-controller/api/v1alpha1"
	"hardware-classification-controller/manager"
	validate "hardware-classification-controller/validate"

	bmh "github.com/metal3-io/baremetal-operator/pkg/apis/metal3/v1alpha1"
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

	fmt.Println(hardwareClassification)

	// Get ExpectedHardwareConfiguraton from hardwareClassification
	extractedProfile := hardwareClassification.Spec.ExpectedHardwareConfiguration
	fmt.Println("-----------------------------------------")
	fmt.Printf("Extracted expected hardware configuration successfully %+v", extractedProfile)
	fmt.Println("-----------------------------------------")

	bmhList, err := fetchBmhHostList(ctx, r, hardwareClassification.Spec.ExpectedHardwareConfiguration.Namespace)
	if err != nil {
		r.Log.Error(err, "unable to fetch baremetal host list", "error", err.Error())
		return ctrl.Result{}, nil
	}

	extractedHardwareDetails := make(map[string]map[string]interface{})

	if extractedProfile != (hwcc.ExpectedHardwareConfiguration{}) {
		for _, host := range bmhList {
			introspectionDetails := make(map[string]interface{})

			if (extractedProfile.CPU == (hwcc.CPU{})) && (extractedProfile.Disk == (hwcc.Disk{})) && (extractedProfile.NICS == (hwcc.NICS{})) && (extractedProfile.RAM == 0) {
				err = errors.New("Provided configurations are not valid")
				break
			}

			if extractedProfile.CPU != (hwcc.CPU{}) {
				if extractedProfile.CPU.Count > 0 {
					introspectionDetails["CPU"] = host.Status.HardwareDetails.CPU
				} else {
					err = errors.New("Enter valid CPU count")
					break
				}
			} else {
				err = errors.New("Enter valid CPU Details")
				break
			}

			if extractedProfile.Disk != (hwcc.Disk{}) {
				if extractedProfile.Disk.Count > 0 && extractedProfile.Disk.SizeGB > 0 {
					introspectionDetails["Disk"] = host.Status.HardwareDetails.Storage
				} else {
					err = errors.New("Enter valid Disk count and Disk Size")
					break
				}
			} else {
				err = errors.New("Enter valid Disk Details")
				break
			}

			if extractedProfile.NICS != (hwcc.NICS{}) {
				if extractedProfile.NICS.Count > 0 {
					introspectionDetails["NICS"] = host.Status.HardwareDetails.NIC
				} else {
					err = errors.New("Enter valid NIC count")
					break
				}
			} else {
				err = errors.New("Enter valid NICS Details")
				break
			}

			if extractedProfile.RAM > 0 {
				introspectionDetails["RAMMebibytes"] = host.Status.HardwareDetails.RAMMebibytes
			} else {
				err = errors.New("Enter valid RAM size")
				break
			}

			if len(introspectionDetails) > 0 {
				extractedHardwareDetails[host.ObjectMeta.Name] = introspectionDetails
			}
		}

	}

	if err != nil {
		r.Log.Error(nil, "Unable to extract details", "error", err.Error())
		return ctrl.Result{}, nil
	}

	fmt.Println("-----------------------------------------")
	fmt.Printf("Extracted Hardware Details %+v", extractedHardwareDetails)
	fmt.Println("-----------------------------------------")

	if len(extractedHardwareDetails) > 0 {
		validatedHardwareDetails := validate.Validation(extractedHardwareDetails)
		fmt.Println(validatedHardwareDetails)
		manager.Manager(extractedProfile.CustomFilter, validatedHardwareDetails, extractedProfile)
	} else {
		fmt.Println("Provided configurations are not valid")
	}

	return ctrl.Result{}, nil
}

func fetchBmhHostList(ctx context.Context, r *HardwareClassificationControllerReconciler, namespace string) ([]bmh.BareMetalHost, error) {

	bmhHostList := bmh.BareMetalHostList{}
	validHostList := []bmh.BareMetalHost{}
	hardwareClassification := &hwcc.HardwareClassificationController{}

	opts := &client.ListOptions{
		Namespace: namespace,
	}

	// Get list of BareMetalHost
	err := r.Client.List(ctx, &bmhHostList, opts)
	if err != nil {
		setError(hardwareClassification, "Failed to get BareMetalHost List")
		return validHostList, err
	}

	// Get hosts in ready and inspecting status from bmhHostList
	for _, host := range bmhHostList.Items {
		if host.Status.Provisioning.State == "ready" {
			validHostList = append(validHostList, host)
		}
	}

	return validHostList, nil
}

// setError sets the ErrorMessage field on the HardwareClassificationController
func setError(hwcc *hwcc.HardwareClassificationController, message string) {
	hwcc.Status.ErrorMessage = pointer.StringPtr(message)
}

func (r *HardwareClassificationControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hwcc.HardwareClassificationController{}).
		Complete(r)
}
