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
	"fmt"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	capi "sigs.k8s.io/cluster-api/api/v1alpha2"
	capierrors "sigs.k8s.io/cluster-api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// "fmt"

	metal3iov1alpha1 "hardware-classification-controller/api/v1alpha1"
)

// HardwareClassificationControllerReconciler reconciles a HardwareClassificationController object
type HardwareClassificationControllerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// MachineManager is responsible for performing machine reconciliation
type MachineManager struct {
	client  client.Client
	Cluster *capi.Cluster
	Machine *capi.Machine
	Log     logr.Logger
}

// Reconcile reconcile function
// +kubebuilder:rbac:groups=metal3.io.sigs.k8s.io,resources=hardwareclassificationcontrollers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=metal3.io.sigs.k8s.io,resources=hardwareclassificationcontrollers/status,verbs=get;update;patch
func (r *HardwareClassificationControllerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	machineLog := r.Log.WithValues("hardwareclassificationcontroller", req.NamespacedName)

	// Fetch the BareMetalMachine instance.
	capbmMachine := &metal3iov1alpha1.HardwareClassificationController{}

	if err := r.Client.Get(ctx, req.NamespacedName, capbmMachine); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	helper, err := patch.NewHelper(capbmMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, errors.Wrap(err, "failed to init patch helper")
	}

	defer func() {
		err := helper.Patch(ctx, capbmMachine)
		if err != nil {
			machineLog.Info("failed to Patch capbmMachine")
		}
	}()

	//clear an error if one was previously set
	clearErrorBMMachine(capbmMachine)

	// Fetch the Machine.
	capiMachine, err := util.GetOwnerMachine(ctx, r.Client, capbmMachine.ObjectMeta)

	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "BareMetalMachine's owner Machine could not be retrieved")
	}
	if capiMachine == nil {
		machineLog.Info("Waiting for Machine Controller to set OwnerRef on BareMetalMachine")
		return ctrl.Result{}, nil
	}

	machineLog = machineLog.WithValues("machine", capiMachine.Name)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, capiMachine.ObjectMeta)
	if err != nil {
		machineLog.Info("BareMetalMachine's owner Machine is missing cluster label or cluster does not exist")
		setErrorBMMachine(capbmMachine, "BareMetalMachine's owner Machine is missing cluster label or cluster does not exist", capierrors.InvalidConfigurationMachineError)

		return ctrl.Result{}, errors.Wrapf(err, "BareMetalMachine's owner Machine is missing label or the cluster does not exist")
	}

	if cluster == nil {
		setErrorBMMachine(capbmMachine, fmt.Sprintf(
			"The machine is NOT associated with a cluster using the label %s: <name of cluster>",
			capi.MachineClusterLabelName,
		), capierrors.InvalidConfigurationMachineError)
		machineLog.Info(fmt.Sprintf("The machine is NOT associated with a cluster using the label %s: <name of cluster>", capi.MachineClusterLabelName))
		return ctrl.Result{}, nil
	}

	machineLog = machineLog.WithValues("cluster", cluster.Name)

	// Make sure infrastructure is ready
	if !cluster.Status.InfrastructureReady {
		machineLog.Info("Waiting for BareMetalCluster Controller to create cluster infrastructure")
		return ctrl.Result{}, nil
	}

	machineMgr, err := NewMachineManager(r.Client, cluster, capiMachine, machineLog)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to create helper for managing the machineMgr")
	}

	fetchHostList(ctx, machineMgr)

	hardwareClassification := &metal3iov1alpha1.HardwareClassificationController{}

	if err := r.Client.Get(ctx, req.NamespacedName, hardwareClassification); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	// your logic here

	fmt.Println("OUTPUT************************", hardwareClassification.Spec.ExpectedHardwareConfiguration)
	return ctrl.Result{}, nil
}

func fetchHostList(ctx context.Context, mgr *MachineManager) {
	// get list of BMH
	// hosts := bmh.BareMetalHostList{}
}

// NewMachineManager returns a new helper for managing a machine
func NewMachineManager(client client.Client, cluster *capi.Cluster,
	machine *capi.Machine,
	machineLog logr.Logger) (*MachineManager, error) {

	return &MachineManager{
		client:  client,
		Cluster: cluster,
		Machine: machine,
		Log:     machineLog,
	}, nil
}

// setError sets the ErrorMessage and ErrorReason fields on the baremetalmachine
func setErrorBMMachine(bmm *metal3iov1alpha1.HardwareClassificationController, message string, reason capierrors.MachineStatusError) {

	bmm.Status.ErrorMessage = pointer.StringPtr(message)

}

// clearError removes the ErrorMessage from the baremetalmachine's Status if set.
func clearErrorBMMachine(bmm *metal3iov1alpha1.HardwareClassificationController) {

	if bmm.Status.ErrorMessage != nil {
		bmm.Status.ErrorMessage = nil
	}

}

func (r *HardwareClassificationControllerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&metal3iov1alpha1.HardwareClassificationController{}).
		Complete(r)
}
