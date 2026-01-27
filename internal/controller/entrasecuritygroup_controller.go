package controller

import (
	"context"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	entraGroup "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	"github.com/vimal-vijayan/entra-governance/internal/entra/groups"
)

// EntraSecurityGroupReconciler reconciles a EntraSecurityGroup object
type EntraSecurityGroupReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	GroupService *groups.Service
}

const (
	defaultRequeueDuration           = 10 * time.Minute
	faildStatusUpdateRequeueDuration = 10 * time.Second
	entraSecurityGroupFinalizer      = "finalizer.entraSecurityGroup.iam.entra.governance.com"
)

// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entrasecuritygroups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entrasecuritygroups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entrasecuritygroups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *EntraSecurityGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("------------------ Reconciling EntraSecurityGroup --------------------", "name", req.Name, "namespace", req.Namespace)

	entraGroup := &entraGroup.EntraSecurityGroup{}
	if err := r.Get(ctx, req.NamespacedName, entraGroup); err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("EntraSecurityGroup resource not found. skipping reconciliation.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get EntraSecurityGroup")
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	// Ensure finalizer is present
	if err := r.ensureFinalizer(ctx, entraGroup); err != nil {
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	if !entraGroup.DeletionTimestamp.IsZero() {
		logger.Info("EntraSecurityGroup resource is being deleted. skipping reconciliation.")
		return r.deleteResource(ctx, entraGroup)
	}

	//TODO: Reconciliation logic goes here, later move to helper functions/services
	if entraGroup.Status.ID != "" {
		// Group already exists in status, checking if group exists in Entra
		if err := r.CheckAndUpdateGroupExists(ctx, entraGroup); err != nil {
			return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
		}

		// Group exists, check members and owners
		if entraGroup.Status.ManagedMemberGroups != nil {
			// check if members are in sync, if not create
			if err := r.CheckAndUpdateMembers(ctx, *entraGroup); err != nil {
				return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
			}
		}

		if entraGroup.Status.Owners != nil {
			// check if owners are in sync, if not create
		}

		// Group already exists and is synced, requeue after default duration
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
	}

	// Group doesn't exist yet, create it
	return r.createResource(ctx, entraGroup)
}

func (r *EntraSecurityGroupReconciler) CheckAndUpdateMembers(ctx context.Context, entraGroup entraGroup.EntraSecurityGroup) error {
	// check if members are in sync with the spec
	logger := log.FromContext(ctx)

	unManagedMemberIds, err := r.GroupService.CheckMemberIds(ctx, entraGroup)
	if err != nil {
		logger.Error(err, "failed to check member IDs for Entra Security Group", "GroupID", entraGroup.Status.ID)
		return err
	}

	if len(unManagedMemberIds) > 0 {
		// add missing members
		// missing members need to be added from the spec
		logger.Info("groups members are not in sync. adding missing members.", "missingMemberIDs", unManagedMemberIds)
	}

	return nil
}

func (r *EntraSecurityGroupReconciler) CheckAndUpdateGroupExists(ctx context.Context, entraGroup *entraGroup.EntraSecurityGroup) error {
	logger := log.FromContext(ctx)

	_, err := r.GroupService.Get(ctx, *entraGroup, entraGroup.Status.ID)
	if err != nil {
		logger.Error(err, "failed to get Entra Security Group by ID from status", "GroupID", entraGroup.Status.ID)
		entraGroup.Status.ID = ""
		entraGroup.Status.DisplayName = ""
		entraGroup.Status.Phase = "Pending"
		if err := r.Status().Update(ctx, entraGroup); err != nil {
			logger.Error(err, "failed to clear EntraSecurityGroup status after failed get")
		}
		return err
	}

	return nil
}

// setupWithManager sets up the controller with the Manager.
func (r *EntraSecurityGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&entraGroup.EntraSecurityGroup{}).
		Complete(r)
}

// create security group in Entra and update status
func (r *EntraSecurityGroupReconciler) createResource(ctx context.Context, entraGroup *entraGroup.EntraSecurityGroup) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	groupId, groupName, err := r.GroupService.Create(ctx, *entraGroup)
	if err != nil {
		logger.Error(err, "failed to create Entra Security Group")
		entraGroup.Status.Phase = "Failed"
		if err := r.Status().Update(ctx, entraGroup); err != nil {
			logger.Error(err, "failed to update EntraSecurityGroup status after creation failure")
		}
		return ctrl.Result{RequeueAfter: faildStatusUpdateRequeueDuration}, err
	}

	// Update status with the created group ID
	entraGroup.Status.ID = groupId
	entraGroup.Status.DisplayName = groupName
	entraGroup.Status.ObservedGeneration = entraGroup.Generation
	entraGroup.Status.Phase = "Success"
	if err := r.Status().Update(ctx, entraGroup); err != nil {
		logger.Error(err, "failed to update EntraSecurityGroup status with GroupID")
		return ctrl.Result{Requeue: true}, err
	}

	logger.Info("Successfully created Entra Security Group", "GroupID", groupId, "DisplayName", groupName)
	return ctrl.Result{Requeue: true}, nil
}

// ensure finalizer is present on the resource
func (r *EntraSecurityGroupReconciler) ensureFinalizer(ctx context.Context, entraGroup *entraGroup.EntraSecurityGroup) error {
	logger := log.FromContext(ctx)
	if !controllerutil.ContainsFinalizer(entraGroup, entraSecurityGroupFinalizer) {
		controllerutil.AddFinalizer(entraGroup, entraSecurityGroupFinalizer)
		if err := r.Update(ctx, entraGroup); err != nil {
			logger.Error(err, "failed to add finalizer to EntraSecurityGroup")
			return err
		}
		logger.Info("finalizer added to EntraSecurityGroup")
	}
	return nil
}

// Delete resource and remove finalizer
func (r *EntraSecurityGroupReconciler) deleteResource(ctx context.Context, entraGroup *entraGroup.EntraSecurityGroup) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	statusCode, err := r.GroupService.Get(ctx, *entraGroup, entraGroup.Status.ID)
	if err != nil {
		if statusCode == "404" {
			logger.Info("Entra Security Group not found in Entra. Removing finalizer.")
			return r.removeFinalizer(ctx, entraGroup)
		}
		logger.Error(err, "failed to get Entra Security Group in Entra during deletion")
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	if entraGroup.Status.ID == "" {
		logger.Info("entra security group id is empty in status. skipping deletion in Entra.")
		return r.removeFinalizer(ctx, entraGroup)
	}

	err = r.GroupService.Delete(ctx, *entraGroup, entraGroup.Status.ID)
	if err != nil {
		logger.Error(err, "failed to delete Entra Security Group in Entra")
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}
	// remove finalizer
	return r.removeFinalizer(ctx, entraGroup)
}

// Remove finalizer
func (r *EntraSecurityGroupReconciler) removeFinalizer(ctx context.Context, entraGroup *entraGroup.EntraSecurityGroup) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	controllerutil.RemoveFinalizer(entraGroup, entraSecurityGroupFinalizer)
	if err := r.Update(ctx, entraGroup); err != nil {
		logger.Error(err, "failed to remove finalizer from EntraSecurityGroup")
		return ctrl.Result{RequeueAfter: 60 * time.Second}, err
	}
	logger.Info("finalizer removed from EntraSecurityGroup. deletion complete.")
	return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
}
