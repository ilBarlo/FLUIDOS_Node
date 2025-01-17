/*
Copyright 2023.

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

package discoverymanager

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	advertisementv1alpha1 "github.com/fluidos-project/node/apis/advertisement/v1alpha1"
	nodecorev1alpha1 "github.com/fluidos-project/node/apis/nodecore/v1alpha1"
	gateway "github.com/fluidos-project/node/pkg/rear-controller/gateway"
	resourceforge "github.com/fluidos-project/node/pkg/utils/resourceforge"
	"github.com/fluidos-project/node/pkg/utils/tools"
)

// DiscoveryReconciler reconciles a Discovery object
type DiscoveryReconciler struct {
	client.Client
	Scheme  *runtime.Scheme
	Gateway *gateway.Gateway
}

// clusterRole
//+kubebuilder:rbac:groups=advertisement.fluidos.eu,resources=peeringcandidates,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=advertisement.fluidos.eu,resources=peeringcandidates/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=advertisement.fluidos.eu,resources=peeringcandidates/finalizers,verbs=update
//+kubebuilder:rbac:groups=advertisement.fluidos.eu,resources=discoveries,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=advertisement.fluidos.eu,resources=discoveries/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=advertisement.fluidos.eu,resources=discoveries/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Discovery object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *DiscoveryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx, "discovery", req.NamespacedName)
	ctx = ctrl.LoggerInto(ctx, log)

	var peeringCandidate *advertisementv1alpha1.PeeringCandidate
	var peeringCandidateReserved advertisementv1alpha1.PeeringCandidate

	var discovery advertisementv1alpha1.Discovery
	if err := r.Get(ctx, req.NamespacedName, &discovery); client.IgnoreNotFound(err) != nil {
		klog.Errorf("Error when getting Discovery %s before reconcile: %s", req.NamespacedName, err)
		return ctrl.Result{}, err
	} else if err != nil {
		klog.Infof("Discovery %s not found, probably deleted", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	klog.Infof("Discovery %s started", discovery.Name)

	if discovery.Status.Phase.Phase != nodecorev1alpha1.PhaseSolved &&
		discovery.Status.Phase.Phase != nodecorev1alpha1.PhaseTimeout &&
		discovery.Status.Phase.Phase != nodecorev1alpha1.PhaseFailed &&
		discovery.Status.Phase.Phase != nodecorev1alpha1.PhaseRunning &&
		discovery.Status.Phase.Phase != nodecorev1alpha1.PhaseIdle {

		discovery.Status.Phase.StartTime = tools.GetTimeNow()
		discovery.SetPhase(nodecorev1alpha1.PhaseRunning, "Discovery started")

		if err := r.updateDiscoveryStatus(ctx, &discovery); err != nil {
			klog.Errorf("Error when updating Discovery %s status before reconcile: %s", req.NamespacedName, err)
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	switch discovery.Status.Phase.Phase {
	case nodecorev1alpha1.PhaseRunning:
		flavours, err := r.Gateway.DiscoverFlavours(discovery.Spec.Selector)
		if err != nil {
			klog.Errorf("Error when getting Flavour: %s", err)
			discovery.SetPhase(nodecorev1alpha1.PhaseFailed, "Error when getting Flavour")
			if err := r.updateDiscoveryStatus(ctx, &discovery); err != nil {
				klog.Errorf("Error when updating Discovery %s status: %s", req.NamespacedName, err)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		if len(flavours) == 0 {
			klog.Infof("No Flavours found")
			discovery.SetPhase(nodecorev1alpha1.PhaseFailed, "No Flavours found")
			if err := r.updateDiscoveryStatus(ctx, &discovery); err != nil {
				klog.Errorf("Error when updating Discovery %s status: %s", req.NamespacedName, err)
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		klog.Infof("Flavours found: %d", len(flavours))

		// TODO: check if a corresponding PeeringCandidate already exists!!
		var first bool = true
		for _, flavour := range flavours {
			if first {
				// We refer to the first peering candidate as the one that is reserved
				peeringCandidate = resourceforge.ForgePeeringCandidate(flavour, discovery.Spec.SolverID, true)
				peeringCandidateReserved = *peeringCandidate
				first = false
			} else {
				// the others are created as not reserved
				peeringCandidate = resourceforge.ForgePeeringCandidate(flavour, discovery.Spec.SolverID, false)
			}

			err = r.Create(context.Background(), peeringCandidate)
			if err != nil {
				klog.Infof("Discovery %s failed: error while creating Peering Candidate", discovery.Name)
				return ctrl.Result{}, err
			}
		}

		// Update the Discovery with the PeeringCandidate
		discovery.Status.PeeringCandidate = nodecorev1alpha1.GenericRef{
			Name:      peeringCandidateReserved.Name,
			Namespace: peeringCandidateReserved.Namespace,
		}

		discovery.SetPhase(nodecorev1alpha1.PhaseSolved, "Discovery Solved: Peering Candidate found")
		if err := r.updateDiscoveryStatus(ctx, &discovery); err != nil {
			klog.Errorf("Error when updating Discovery %s: %s", discovery.Name, err)
			return ctrl.Result{}, err
		}
		klog.Infof("Discovery %s updated", discovery.Name)

		return ctrl.Result{}, nil

	case nodecorev1alpha1.PhaseSolved:
		klog.Infof("Discovery %s solved", discovery.Name)
	case nodecorev1alpha1.PhaseFailed:
		klog.Infof("Discovery %s failed", discovery.Name)
	}

	return ctrl.Result{}, nil
}

// updateDiscoveryStatus updates the status of the discovery
func (r *DiscoveryReconciler) updateDiscoveryStatus(ctx context.Context, discovery *advertisementv1alpha1.Discovery) error {
	return r.Status().Update(ctx, discovery)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DiscoveryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&advertisementv1alpha1.Discovery{}).
		Complete(r)
}
