package controller

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	iamv1alpha1 "github.com/vimal-vijayan/entra-governance/api/v1alpha1"
	entraClient "github.com/vimal-vijayan/entra-governance/internal/entra/client"
)

// EntraSecurityGroupReconciler reconciles a EntraSecurityGroup object
type EntraSecurityGroupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	defaultRequeueDuration      = 60 * time.Second
	entraSecurityGroupFinalizer = "finalizer.entraSecurityGroup.iam.entra.governance.com"
)

// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entrasecuritygroups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entrasecuritygroups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=iam.entra.governance.com,resources=entrasecuritygroups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *EntraSecurityGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("------------------ Reconciling EntraSecurityGroup --------------------", "name", req.Name, "namespace", req.Namespace)

	entraGroup := &iamv1alpha1.EntraSecurityGroup{}
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
		//TODO: Handle deletion logic here if needed.. currenlty reconilation is skipped
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
	}

	// Reconciliation logic goes here
	if entraGroup.Status.ID != "" {
		// Group already exists in status, checking if group exists in Entra
		logger.Info("entra security group already has GroupID in status. Skipping creation.", "GroupID", entraGroup.Status.ID)
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
	}

	// Create Entra Security Group
	groupID, err := r.createEntraSecurityGroup(ctx, entraGroup)
	if err != nil {
		logger.Error(err, "failed to create Entra Security Group")
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}

	// Update status with the created group ID
	entraGroup.Status.ID = groupID.ID
	if err := r.Status().Update(ctx, entraGroup); err != nil {
		logger.Error(err, "failed to update EntraSecurityGroup status with GroupID")
		return ctrl.Result{RequeueAfter: defaultRequeueDuration}, err
	}
	
	logger.Info("Successfully created Entra Security Group", "GroupID", groupID.ID)

	return ctrl.Result{RequeueAfter: defaultRequeueDuration}, nil
}

// setupWithManager sets up the controller with the Manager.
func (r *EntraSecurityGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&iamv1alpha1.EntraSecurityGroup{}).
		Complete(r)
}

// ensure finalizer is present on the resource
func (r *EntraSecurityGroupReconciler) ensureFinalizer(ctx context.Context, entraGroup *iamv1alpha1.EntraSecurityGroup) error {
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

// create entra security group in Entra using the client
type createGroupResponse struct {
	ID string
	DisplayName string
}

func (r *EntraSecurityGroupReconciler) createEntraSecurityGroup(ctx context.Context, entraGroup *iamv1alpha1.EntraSecurityGroup) (*createGroupResponse, error) {
	logger := log.FromContext(ctx)
	logger.Info("creating entra security group in entra", "name", entraGroup.Spec.Name)

	secretRef := entraClient.SecretRef{
		Name:      entraGroup.Spec.ForProvider.CredentialSecretRef,
		Namespace: entraGroup.Namespace,
	}

	// Initialize the Entra client
	factory := entraClient.NewClientFactory(r.Client)
	if entraGroup.Spec.ForProvider.CredentialSecretRef != "" {
		sdk, err := factory.ForClientSecret(ctx, secretRef)

		if err != nil {
			logger.Error(err, "failed to initialize Entra client using secret reference")
			return 	nil, err
		}

		entraClient := entraClient.NewGraphClient(sdk)
		resp, err := entraClient.CreateEntraGroup(entraGroup.Spec)

		return &createGroupResponse{
			ID: resp.ID,
			DisplayName: resp.DisplayName,
		}, err
	}

	logger.Error(nil, "credentialSecretRef is not specified in the EntraSecurityGroup spec")
	return nil, fmt.Errorf("credentialSecretRef is not specified in the EntraSecurityGroup spec")
}
